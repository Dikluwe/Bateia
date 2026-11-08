/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/source_insight.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"sort"
	"strings"
	"unicode/utf8"
)

const sourceInsightContract = "insight-source-analysis-v2"

// InvestigateSourceInsights locates supplied observations without re-assessing their meaning.
func InvestigateSourceInsights(ctx context.Context, request domain.InsightInvestigationRequest) (domain.SourceInsightInvestigationReport, error) {
	if ctx == nil || request.SchemaVersion != "2" {
		return domain.SourceInsightInvestigationReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.SourceInsightInvestigationReport{}, err
	}
	catalog, catalogHash := ActiveMechanismRuleCatalog()
	rules, err := decodeMechanismCatalog(catalog, catalogHash)
	if err != nil {
		return domain.SourceInsightInvestigationReport{}, domain.ErrIntegrity
	}
	request, err = normalizeInsightRequest(request, rules)
	if err == nil {
		request, err = normalizeSourceObservations(request)
	}
	if err != nil {
		return domain.SourceInsightInvestigationReport{}, err
	}
	classificationRequest := request
	classificationRequest.SchemaVersion = "1"
	classificationRequest.Sources = append([]domain.InsightEvidenceInput(nil), request.Sources...)
	for i := range classificationRequest.Sources {
		classificationRequest.Sources[i].Observations = nil
	}
	base, err := InvestigateInsights(ctx, classificationRequest)
	if err != nil {
		return domain.SourceInsightInvestigationReport{}, err
	}
	inputHash := insightHash(request)
	report := domain.SourceInsightInvestigationReport{
		SchemaVersion: sourceInsightContract, ID: "insight_" + sourceInsightIdentity(inputHash, catalogHash, ""),
		Scope: "FROZEN_DECLARED_SOURCE_OBSERVATIONS", Status: base.Status, Question: request.Question, Target: request.Target,
		InputSHA256: inputHash, CatalogSHA256: catalogHash, SourceAnalyses: []domain.SourceInsightAnalysis{},
		Hypotheses: []domain.ImprovementHypothesis{}, ExtractionGaps: base.ExtractionGaps, Provenance: base.Provenance,
		Limitations: append(append([]string(nil), base.Limitations...), "As observações são declarações fornecidas: referências e trechos foram conferidos, mas a interpretação não foi reavaliada.", "A definição do catálogo é uma classificação genérica; condições operacionais dependem das observações declaradas de cada fonte."),
	}
	mechanisms := map[string]domain.ArchitecturalMechanism{}
	for _, mechanism := range base.Mechanisms {
		mechanisms[sourceAnalysisKey(mechanism.Source, mechanism.MechanismClass)] = mechanism
	}
	for _, source := range request.Sources {
		if err := ctx.Err(); err != nil {
			return domain.SourceInsightInvestigationReport{}, err
		}
		ref := immutableRef(source.Bundle.Repository)
		for _, class := range request.MechanismClasses {
			observations := projectSourceObservations(source, class)
			mechanism, observed := mechanisms[sourceAnalysisKey(ref, class)]
			if !observed && len(observations) == 0 {
				continue
			}
			analysis := domain.SourceInsightAnalysis{
				Source: ref, MechanismStatus: domain.MechanismInconclusive, AssessmentOrigin: "DECLARED_NOT_REASSESSED",
				Evidence: []domain.MechanismWitness{}, CounterEvidence: []domain.MechanismWitness{}, Observations: observations, Gaps: []string{},
			}
			for _, rule := range rules.Rules {
				if rule.MechanismClass == class {
					analysis.CatalogDefinition = domain.InsightCatalogDefinition{Origin: "CATALOG_TEMPLATE", MechanismClass: class, Name: normalizeMechanismText(rule.Name), Purpose: normalizeMechanismText(rule.Purpose), Strategy: normalizeMechanismText(rule.Strategy)}
					analysis.ExtractionScope = domain.InsightExtractionScope{Languages: canonicalMechanismSet(rule.LanguageLimits), Ecosystems: canonicalMechanismSet(rule.EcosystemLimits)}
					break
				}
			}
			if observed {
				analysis.MechanismRef, analysis.MechanismStatus = mechanism.ID, mechanism.Status
				analysis.Evidence = append(analysis.Evidence, mechanism.Evidence...)
				analysis.CounterEvidence = append(analysis.CounterEvidence, mechanism.CounterEvidence...)
			} else {
				analysis.Gaps = append(analysis.Gaps, "MECHANISM_NOT_OBSERVED")
			}
			if len(observations) == 0 {
				analysis.Gaps = append(analysis.Gaps, "SOURCE_SPECIFIC_OBSERVATIONS_NOT_PROVIDED")
			}
			if !hasObservationKind(observations, domain.InsightPrecondition) {
				analysis.Gaps = append(analysis.Gaps, "OPERATIONAL_PRECONDITIONS_NOT_DECLARED")
			}
			report.SourceAnalyses = append(report.SourceAnalyses, analysis)
		}
	}
	for _, hypothesis := range base.Hypotheses {
		for _, analysis := range report.SourceAnalyses {
			if analysis.MechanismRef == hypothesis.MechanismRefs[0] {
				report.Hypotheses = append(report.Hypotheses, projectSourceHypothesis(hypothesis, analysis, inputHash, catalogHash))
				break
			}
		}
	}
	return report, nil
}

func normalizeSourceObservations(request domain.InsightInvestigationRequest) (domain.InsightInvestigationRequest, error) {
	classes := map[string]bool{}
	for _, class := range request.MechanismClasses {
		classes[class] = true
	}
	observationCount, referenceCount := 0, 0
	for i := range request.Sources {
		source := &request.Sources[i]
		observationCount += len(source.Observations)
		if len(source.Observations) > 32 || observationCount > 128 {
			return request, domain.ErrInvalidRunInput
		}
		source.Observations = append([]domain.SourceInsightObservation(nil), source.Observations...)
		items := map[string]domain.CollectedEvidence{}
		for _, item := range source.Bundle.Items {
			items[item.Path] = item
		}
		for j := range source.Observations {
			observation := &source.Observations[j]
			if !classes[observation.MechanismClass] || !validSynthesisText(observation.Statement, 1000) || len(observation.References) < 1 || len(observation.References) > 8 {
				return request, domain.ErrInvalidRunInput
			}
			switch observation.Kind {
			case domain.InsightStrategy, domain.InsightPrecondition, domain.InsightFailureMode, domain.InsightTradeoff, domain.InsightLimitation:
			default:
				return request, domain.ErrInvalidRunInput
			}
			observation.Statement = normalizeMechanismText(observation.Statement)
			if observation.Statement == "" {
				return request, domain.ErrInvalidRunInput
			}
			observation.References = append([]domain.InsightObservationReference(nil), observation.References...)
			referenceCount += len(observation.References)
			if referenceCount > 256 {
				return request, domain.ErrInvalidRunInput
			}
			for k := range observation.References {
				reference := &observation.References[k]
				item, found := items[reference.Path]
				if !found || !validMechanismPath(reference.Path) {
					return request, domain.ErrIntegrity
				}
				digest := contentHash(item.Content)
				if strings.TrimPrefix(reference.ContentHash, "sha256:") != digest {
					return request, domain.ErrIntegrity
				}
				reference.ContentHash = "sha256:" + digest
				if !validObservationRange(item.Content, *reference) {
					return request, domain.ErrInvalidRunInput
				}
			}
			sort.Slice(observation.References, func(a, b int) bool {
				return string(canonicalJSON(observation.References[a])) < string(canonicalJSON(observation.References[b]))
			})
			for k := 1; k < len(observation.References); k++ {
				if observation.References[k] == observation.References[k-1] {
					return request, domain.ErrInvalidRunInput
				}
			}
		}
		sort.Slice(source.Observations, func(a, b int) bool {
			return string(canonicalJSON(source.Observations[a])) < string(canonicalJSON(source.Observations[b]))
		})
		for j := 1; j < len(source.Observations); j++ {
			if string(canonicalJSON(source.Observations[j])) == string(canonicalJSON(source.Observations[j-1])) {
				return request, domain.ErrInvalidRunInput
			}
		}
	}
	encoded, err := json.Marshal(request)
	if err != nil || len(encoded) > 2<<20 {
		return request, domain.ErrInvalidRunInput
	}
	return request, nil
}

func observationLines(content string) []string {
	if content == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(content, "\n"), "\n")
}

func validObservationRange(content string, reference domain.InsightObservationReference) bool {
	lines := observationLines(content)
	if reference.StartLine < 1 || reference.EndLine < reference.StartLine || reference.EndLine > len(lines) || reference.EndLine-reference.StartLine >= 20 {
		return false
	}
	excerpt := strings.Join(lines[reference.StartLine-1:reference.EndLine], "\n")
	return len(excerpt) <= 4096 && utf8.ValidString(excerpt) && strings.TrimSpace(excerpt) != ""
}

func projectSourceObservations(source domain.InsightEvidenceInput, class string) []domain.LocatedInsightObservation {
	items := map[string]domain.CollectedEvidence{}
	for _, item := range source.Bundle.Items {
		items[item.Path] = item
	}
	out := []domain.LocatedInsightObservation{}
	for _, observation := range source.Observations {
		if observation.MechanismClass != class {
			continue
		}
		located := domain.LocatedInsightObservation{
			ID: "obs_" + insightHash(struct {
				Source      domain.ImmutableCandidateRef
				Observation domain.SourceInsightObservation
			}{immutableRef(source.Bundle.Repository), observation}),
			Kind: observation.Kind, Statement: observation.Statement, Origin: "DECLARED_NOT_REASSESSED", References: []domain.LocatedInsightReference{},
		}
		for _, reference := range observation.References {
			lines := observationLines(items[reference.Path].Content)
			located.References = append(located.References, domain.LocatedInsightReference{InsightObservationReference: reference, Excerpt: strings.Join(lines[reference.StartLine-1:reference.EndLine], "\n")})
		}
		out = append(out, located)
	}
	return out
}

func sourceAnalysisKey(source domain.ImmutableCandidateRef, class string) string {
	return source.Repository + "@" + source.Commit + "\x00" + class
}

func sourceInsightIdentity(input, catalog, mechanism string) string {
	return insightHash(struct{ Contract, Input, Catalog, Mechanism string }{sourceInsightContract, input, catalog, mechanism})
}

func hasObservationKind(observations []domain.LocatedInsightObservation, kind domain.InsightObservationKind) bool {
	for _, observation := range observations {
		if observation.Kind == kind {
			return true
		}
	}
	return false
}

func projectSourceHypothesis(hypothesis domain.ImprovementHypothesis, analysis domain.SourceInsightAnalysis, inputHash, catalogHash string) domain.ImprovementHypothesis {
	hypothesis.ID = "hyp_" + sourceInsightIdentity(inputHash, catalogHash, analysis.MechanismRef)
	hypothesis.Preconditions = []string{}
	hypothesis.EvidenceRefs = []string{analysis.MechanismRef}
	tradeoffs := []string{}
	for _, observation := range analysis.Observations {
		hypothesis.EvidenceRefs = append(hypothesis.EvidenceRefs, observation.ID)
		switch observation.Kind {
		case domain.InsightPrecondition:
			hypothesis.Preconditions = append(hypothesis.Preconditions, observation.Statement)
		case domain.InsightTradeoff:
			tradeoffs = append(tradeoffs, observation.Statement)
		}
	}
	if len(tradeoffs) > 0 {
		hypothesis.Tradeoffs = tradeoffs
	}
	return hypothesis
}
