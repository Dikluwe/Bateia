/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_comparison.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
)

const insightComparisonContract = "insight-baseline-comparison-v1"

// CompareInsightBaseline groups supplied observations without inferring a semantic difference.
func CompareInsightBaseline(ctx context.Context, request domain.InsightComparisonRequest) (domain.InsightComparisonReport, error) {
	if ctx == nil || request.SchemaVersion != "1" {
		return domain.InsightComparisonReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.InsightComparisonReport{}, err
	}
	raw, err := json.Marshal(request)
	if err != nil || len(raw) > 2<<20 {
		return domain.InsightComparisonReport{}, domain.ErrInvalidRunInput
	}
	source, err := InvestigateSourceInsights(ctx, request.Investigation)
	if err != nil {
		return domain.InsightComparisonReport{}, err
	}
	baseline, located, err := normalizeInsightBaseline(request.TargetBaseline, request.Investigation.MechanismClasses)
	if err != nil {
		return domain.InsightComparisonReport{}, err
	}
	identity := insightHash(struct {
		Contract, SourceInput, Catalog string
		Baseline                       *domain.InsightTargetBaseline
	}{insightComparisonContract, source.InputSHA256, source.CatalogSHA256, baseline})
	report := domain.InsightComparisonReport{
		SchemaVersion: insightComparisonContract, ID: "comparison_" + identity, Status: domain.ImprovementInconclusive,
		Question: source.Question, Target: source.Target, SourceReportID: source.ID, SourceInputSHA256: source.InputSHA256, CatalogSHA256: source.CatalogSHA256,
		SourceAnalyses: source.SourceAnalyses, SourceProvenance: source.Provenance, SourceExtractionGaps: source.ExtractionGaps,
		TargetBaseline: located, Comparisons: []domain.InsightComparisonRow{}, EffectState: "UNVERIFIED", AdoptionDisposition: "NOT_AN_ADOPTION_DECISION",
		Limitations: append(append([]string(nil), source.Limitations...), "O snapshot do alvo cobre somente os arquivos fornecidos; sua associação ao alvo é declarada, sem autenticação externa ou commit Git inferido.", "Agrupar declarações por classe e tipo não comprova equivalência, diferença, presença ou ausência de uma capacidade no alvo.", "Antes de propor alteração, comparar a política já declarada e validar uma diferença concreta; nenhum efeito foi medido."),
	}
	kinds := []domain.InsightObservationKind{domain.InsightStrategy, domain.InsightPrecondition, domain.InsightFailureMode, domain.InsightTradeoff, domain.InsightLimitation}
	for _, provenance := range source.Provenance {
		for _, class := range canonicalMechanismSet(request.Investigation.MechanismClasses) {
			for _, kind := range kinds {
				if err := ctx.Err(); err != nil {
					return domain.InsightComparisonReport{}, err
				}
				row := domain.InsightComparisonRow{Source: provenance.Source, MechanismClass: class, Kind: kind, SourceObservationRefs: []string{}, TargetObservationRefs: []string{}, Gaps: []string{"SEMANTIC_DIFFERENCE_NOT_ASSESSED", "TARGET_EFFECT_NOT_VERIFIED"}}
				for _, analysis := range source.SourceAnalyses {
					if analysis.Source == provenance.Source && analysis.CatalogDefinition.MechanismClass == class {
						for _, observation := range analysis.Observations {
							if observation.Kind == kind {
								row.SourceObservationRefs = append(row.SourceObservationRefs, observation.ID)
							}
						}
					}
				}
				if located != nil {
					for _, observation := range located.Observations {
						if observation.MechanismClass == class && observation.Kind == kind {
							row.TargetObservationRefs = append(row.TargetObservationRefs, observation.ID)
						}
					}
				}
				row.DeclarationStatus = comparisonDeclarationStatus(len(row.SourceObservationRefs), len(row.TargetObservationRefs))
				if row.DeclarationStatus == "BOTH_DECLARED" {
					report.Status = domain.ImprovementPartial
				}
				if len(row.SourceObservationRefs) == 0 {
					row.Gaps = append(row.Gaps, "SOURCE_OBSERVATION_NOT_DECLARED")
				}
				if len(row.TargetObservationRefs) == 0 {
					row.Gaps = append(row.Gaps, "TARGET_OBSERVATION_NOT_DECLARED")
				}
				row.NextStep = insightComparisonStep(kind, provenance.Source.Repository, source.Target.Name)
				row.ID = "cmp_" + insightHash(struct {
					Input string
					Row   domain.InsightComparisonRow
				}{identity, row})
				report.Comparisons = append(report.Comparisons, row)
			}
		}
	}
	return report, nil
}

func comparisonDeclarationStatus(source, target int) string {
	if source > 0 && target > 0 {
		return "BOTH_DECLARED"
	}
	if source > 0 {
		return "SOURCE_ONLY"
	}
	if target > 0 {
		return "TARGET_ONLY"
	}
	return "NO_DECLARATIONS"
}

func insightComparisonStep(kind domain.InsightObservationKind, source, target string) domain.ValidationStep {
	question := "Como as estratégias declaradas se relacionam"
	switch kind {
	case domain.InsightPrecondition:
		question = "Quais condições operacionais precisam ser comparadas"
	case domain.InsightFailureMode:
		question = "Como falhas, saturação e cancelamento são tratados"
	case domain.InsightTradeoff:
		question = "Quais custos e compromissos precisam ser comparados"
	case domain.InsightLimitation:
		question = "Quais limites das abordagens e evidências precisam ser conferidos"
	}
	return domain.ValidationStep{Question: question + " entre " + source + " e " + target + "?", Observable: string(kind), SafeProcedure: "Conferir os trechos referenciados e a política já declarada no alvo antes de propor alteração; registrar uma diferença concreta e planejar sua medição separadamente.", PassCriterion: "Diferença concreta e relevante à pergunta documentada com referências e plano de baseline/experimento.", RejectCriterion: "Declarações insuficientes, diferença irrelevante ou capacidade já atendida sem benefício adicional demonstrado.", RelativeCost: "LOW", MissingEvidence: "Interpretação comparativa, aplicabilidade e efeito no alvo ainda não verificados."}
}
