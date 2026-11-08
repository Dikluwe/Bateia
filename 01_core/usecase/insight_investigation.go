/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_investigation.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// InvestigateInsights examines frozen, evaluated evidence without a reuse run.
// It verifies bindings and integrity, not the truth of caller-supplied judgments.
func InvestigateInsights(ctx context.Context, request domain.InsightInvestigationRequest) (domain.InsightInvestigationReport, error) {
	if ctx == nil {
		return domain.InsightInvestigationReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.InsightInvestigationReport{}, err
	}
	if request.SchemaVersion != "1" {
		return domain.InsightInvestigationReport{}, domain.ErrInvalidRunInput
	}
	catalog, catalogHash := ActiveMechanismRuleCatalog()
	rules, err := decodeMechanismCatalog(catalog, catalogHash)
	if err != nil {
		return domain.InsightInvestigationReport{}, domain.ErrIntegrity
	}
	request, err = normalizeInsightRequest(request, rules)
	if err != nil {
		return domain.InsightInvestigationReport{}, err
	}
	inputHash := insightHash(request)
	report := domain.InsightInvestigationReport{
		SchemaVersion: "insight-investigation-v1", ID: "insight_" + insightHash(struct{ Input, Catalog string }{inputHash, catalogHash}),
		Scope: "FROZEN_EVALUATED_EVIDENCE", Status: domain.ImprovementInconclusive, Question: request.Question, Target: request.Target,
		InputSHA256: inputHash, CatalogSHA256: catalogHash, Mechanisms: []domain.ArchitecturalMechanism{}, Hypotheses: []domain.ImprovementHypothesis{}, ExtractionGaps: []domain.MechanismGap{}, Provenance: []domain.InsightSourceProvenance{},
		Limitations: []string{"As avaliações fornecidas não foram refeitas nem autenticadas externamente.", "Hashes e localizadores verificam vínculo com os bytes fornecidos, não comportamento executado.", "O alvo, suas restrições e o efeito desejado ainda precisam de validação.", "Nenhuma decisão de reúso ou autorização de adoção foi produzida.", "A busca de contraevidência se limita às avaliações e bytes fornecidos; o ecossistema não foi pesquisado."},
	}
	classes := map[string]bool{}
	for _, class := range request.MechanismClasses {
		classes[class] = true
	}
	incomplete := map[string]bool{}
	extractor := NewDeterministicMechanismExtractor()
	for _, source := range request.Sources {
		if err := ctx.Err(); err != nil {
			return domain.InsightInvestigationReport{}, err
		}
		provenance := domain.InsightSourceProvenance{Source: immutableRef(source.Bundle.Repository), BundleSHA256: insightHash(source.Bundle), AssessmentsSHA256: insightHash(source.Assessments), Complete: source.Bundle.Complete, AssessmentOrigin: "DECLARED_NOT_REASSESSED"}
		report.Provenance = append(report.Provenance, provenance)
		extraction, err := extractor.ExtractMechanisms(ctx, domain.MechanismExtractionInput{Bundle: source.Bundle, Assessments: source.Assessments, RuleCatalogJSON: catalog, RuleCatalogSHA256: catalogHash})
		if err != nil {
			return domain.InsightInvestigationReport{}, err
		}
		report.ExtractionGaps = append(report.ExtractionGaps, extraction.Gaps...)
		for _, mechanism := range extraction.Mechanisms {
			if classes[mechanism.MechanismClass] {
				report.Mechanisms = append(report.Mechanisms, mechanism)
				incomplete[mechanism.ID] = !source.Bundle.Complete || len(source.Bundle.Failures) > 0
			}
		}
	}
	sort.Slice(report.Mechanisms, func(i, j int) bool { return report.Mechanisms[i].ID < report.Mechanisms[j].ID })
	sortMechanismGaps(report.ExtractionGaps)
	if len(report.Mechanisms) == 0 {
		return report, nil
	}
	report.Status = domain.ImprovementPartial // Target validation is still missing.
	for _, mechanism := range report.Mechanisms {
		if len(report.Hypotheses) == request.MaxHypotheses {
			report.Limitations = append(report.Limitations, "O teto de hipóteses limitou a projeção; todos os mecanismos selecionados permanecem no relatório.")
			break
		}
		report.Hypotheses = append(report.Hypotheses, insightHypothesis(request, mechanism, incomplete[mechanism.ID], inputHash, catalogHash))
	}
	return report, nil
}

func normalizeInsightRequest(request domain.InsightInvestigationRequest, rules mechanismCatalog) (domain.InsightInvestigationRequest, error) {
	invalid := func() (domain.InsightInvestigationRequest, error) {
		return domain.InsightInvestigationRequest{}, domain.ErrInvalidRunInput
	}
	if (request.SchemaVersion != "1" && request.SchemaVersion != "2") || !validSynthesisText(request.Question, 1000) || !validSynthesisText(request.Target.Name, 200) || !validSynthesisText(request.Target.Context, 1000) || !validSynthesisText(request.Target.Property, 200) || !validSynthesisText(request.Target.DesiredEffect, 1000) || !validSynthesisTextSet(request.Target.Constraints, 20, 200) {
		return invalid()
	}
	if len(request.Sources) < 1 || len(request.Sources) > 8 || len(request.MechanismClasses) < 1 || len(request.MechanismClasses) > 4 {
		return invalid()
	}
	raw, marshalErr := json.Marshal(request)
	if marshalErr != nil || len(raw) > 2<<20 {
		return invalid()
	}
	if request.MaxHypotheses == 0 {
		request.MaxHypotheses = 5
	}
	if request.MaxHypotheses < 1 || request.MaxHypotheses > 12 {
		return invalid()
	}
	request.Question = normalizeMechanismText(request.Question)
	request.Target.Name = normalizeMechanismText(request.Target.Name)
	request.Target.Context = normalizeMechanismText(request.Target.Context)
	request.Target.Property = normalizeMechanismText(request.Target.Property)
	request.Target.DesiredEffect = normalizeMechanismText(request.Target.DesiredEffect)
	if request.Question == "" || request.Target.Name == "" || request.Target.Context == "" || request.Target.Property == "" || request.Target.DesiredEffect == "" {
		return invalid()
	}
	request.Target.Constraints = canonicalMechanismSet(request.Target.Constraints)
	allowed := map[string]bool{}
	for _, rule := range rules.Rules {
		allowed[rule.MechanismClass] = true
	}
	seenClasses := map[string]bool{}
	for _, class := range request.MechanismClasses {
		if !allowed[class] || seenClasses[class] {
			return invalid()
		}
		seenClasses[class] = true
	}
	request.MechanismClasses = append([]string(nil), request.MechanismClasses...)
	sort.Strings(request.MechanismClasses)
	request.Sources = append([]domain.InsightEvidenceInput(nil), request.Sources...)
	identities := map[string]bool{}
	items, assessments, witnesses := 0, 0, 0
	for _, source := range request.Sources {
		if request.SchemaVersion == "1" && len(source.Observations) != 0 {
			return invalid()
		}
		bundle := source.Bundle
		ref := bundle.Repository
		if bundle.SchemaVersion != "1" || !immutableCommitPattern.MatchString(ref.Commit) || !validSynthesisText(ref.Owner, 100) || !validSynthesisText(ref.Repo, 100) || strings.ContainsAny(ref.Owner+ref.Repo, "/\\@") || ref.Owner != normalizeMechanismText(ref.Owner) || ref.Repo != normalizeMechanismText(ref.Repo) {
			return invalid()
		}
		identity := strings.ToLower(ref.Owner+"/"+ref.Repo) + "@" + ref.Commit
		if identities[identity] {
			return invalid()
		}
		identities[identity] = true
		items += len(bundle.Items)
		assessments += len(source.Assessments)
		if items > 512 || assessments > 128 {
			return invalid()
		}
		if err := validateEvidenceBundle(bundle); err != nil {
			return domain.InsightInvestigationRequest{}, domain.ErrIntegrity
		}
		paths := map[string]bool{}
		for _, item := range bundle.Items {
			if !validMechanismPath(item.Path) || paths[item.Path] {
				return invalid()
			}
			paths[item.Path] = true
		}
		for _, assessment := range source.Assessments {
			if assessment.Repository != ref {
				return domain.InsightInvestigationRequest{}, domain.ErrIntegrity
			}
			witnesses += len(assessment.ProblemAlignment.Sources)
			for _, capability := range assessment.Capabilities {
				witnesses += len(capability.Sources)
			}
			if witnesses > 512 {
				return invalid()
			}
		}
	}
	sort.Slice(request.Sources, func(i, j int) bool {
		return request.Sources[i].Bundle.Repository.Key() < request.Sources[j].Bundle.Repository.Key()
	})
	encoded, err := json.Marshal(request)
	if err != nil || len(encoded) > 2<<20 {
		return invalid()
	}
	return request, nil
}

func insightHash(value any) string {
	hash := sha256.Sum256(canonicalJSON(value))
	return hex.EncodeToString(hash[:])
}

func insightHypothesis(request domain.InsightInvestigationRequest, mechanism domain.ArchitecturalMechanism, incomplete bool, inputHash, catalogHash string) domain.ImprovementHypothesis {
	status := domain.HypothesisProposed
	if incomplete || mechanism.Status != domain.MechanismProven {
		status = domain.HypothesisInconclusive
	}
	counterStatus := domain.CounterEvidenceNoneInScope
	counterRefs := []string{}
	for _, witness := range mechanism.CounterEvidence {
		counterRefs = append(counterRefs, witness.ContentHash)
	}
	if len(counterRefs) > 0 {
		counterStatus = domain.CounterEvidenceFound
	}
	return domain.ImprovementHypothesis{
		ID:    "hyp_" + insightHash(struct{ Input, Catalog, Mechanism string }{inputHash, catalogHash, mechanism.ID}),
		Title: "Investigar " + mechanism.Name + " em " + request.Target.Name, TargetProperty: request.Target.Property, TargetState: domain.TargetUnknown, TargetGapRefs: []string{"TARGET_BASELINE_NOT_OBSERVED", "TARGET_APPLICABILITY_NOT_VERIFIED"}, ExpectedEffect: request.Target.DesiredEffect,
		MechanismRefs: []string{mechanism.ID}, ProductObservationRefs: []string{}, TransferMode: domain.TransferReferenceOnly, License: domain.ImprovementLicenseUnrecognized,
		Assumptions: []string{"O efeito desejado é uma intenção do pedido, não uma observação.", "A estratégia do mecanismo pode ser investigada sem presumir adoção de código."}, Preconditions: append([]string{}, mechanism.Preconditions...),
		IntegrationUnknowns: []string{"Compatibilidade com o contexto e as restrições do alvo.", "Licença e direitos de eventual transferência de código não avaliados."}, Risks: []string{"O mecanismo pode não ser relevante à pergunta ou não produzir o efeito desejado no alvo."}, Tradeoffs: []string{"O custo de adaptação e a utilidade no alvo ainda precisam de medição."}, EvidenceRefs: []string{mechanism.ID}, CounterEvidenceRefs: counterRefs,
		CounterEvidenceSearch: domain.CounterEvidenceSearch{Status: counterStatus, ScopeRefs: []string{mechanism.ID}, ResultRefs: append([]string{}, counterRefs...)},
		ValidationPlan: []domain.ValidationStep{
			{Question: "Qual é o estado atual de " + request.Target.Property + "?", Observable: request.Target.Property, SafeProcedure: "Medir o alvo atual sem modificá-lo, registrando casos e configuração.", PassCriterion: "Baseline repetível registrado.", RejectCriterion: "A medição não é repetível ou não representa a propriedade desejada.", RelativeCost: "LOW", MissingEvidence: "baseline do alvo"},
			{Question: request.Question, Observable: "Compatibilidade da estratégia com contexto e restrições", SafeProcedure: "Comparar a estratégia observada com o contexto fornecido e conferir cada restrição antes de propor transferência.", PassCriterion: "Premissas e restrições verificadas com evidência localizada.", RejectCriterion: "Uma premissa ou restrição necessária é incompatível.", RelativeCost: "MEDIUM", MissingEvidence: "aplicabilidade ao alvo"},
			{Question: "O mecanismo produz o efeito desejado?", Observable: request.Target.DesiredEffect, SafeProcedure: "Planejar experimento isolado comparando baseline e variante; não executar nem modificar o alvo nesta investigação.", PassCriterion: "Definir antes do experimento uma métrica e um limiar de ganho, mantendo as restrições.", RejectCriterion: "Efeito abaixo do limiar, regressão ou violação de restrição.", RelativeCost: "MEDIUM", MissingEvidence: "experimento comparativo e limiares pré-definidos"},
		}, Derivation: domain.DerivationRuleDerived, Status: status, AdoptionDisposition: "NOT_AN_ADOPTION_DECISION",
	}
}
