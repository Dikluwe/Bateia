/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_analysis.md
 * @layer L1
 * @updated 2026-10-06
 */

package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// improvementAnalysisBasis contains source facts resolved and validated by the
// caller. Analysis owns conclusions, but never resolves artifacts or runs jobs.
type improvementAnalysisBasis struct {
	mechanisms map[string]domain.ArchitecturalMechanism
	degraded   bool
}

// improvementAnalysisResult is a business result, independent of run lifecycle.
// A completed job may contain inconclusive hypotheses and a partial result.
type improvementAnalysisResult struct {
	comparisons []domain.ArchitectureComparison
	hypotheses  []domain.ImprovementHypothesis
	gaps        []domain.HypothesisGap
	plans       []domain.ValidationPlanEntry
	provenance  []domain.SynthesisProvenance
	status      domain.ImprovementSynthesisStatus
}

func synthesizeImprovement(request domain.ImprovementSynthesisRequest, basis improvementAnalysisBasis) improvementAnalysisResult {
	comparisons := []domain.ArchitectureComparison{}
	hypotheses := []domain.ImprovementHypothesis{}
	gaps := []domain.HypothesisGap{}
	plans := []domain.ValidationPlanEntry{}
	provenance := []domain.SynthesisProvenance{}
	scopeRefs := make([]string, len(request.SourceArtifacts))
	for index, source := range request.SourceArtifacts {
		scopeRefs[index] = string(source.RunID) + ":" + string(source.Artifact) + ":" + source.SHA256
		provenance = append(provenance, domain.SynthesisProvenance{SourceRunID: source.RunID, Artifact: source.Artifact, SHA256: source.SHA256, RuleIDs: []string{}, ScopeRefs: []string{scopeRefs[index]}})
	}
	sort.Strings(scopeRefs)
	proposals := append([]domain.ImprovementProposal{}, request.Proposals...)
	if len(proposals) == 0 {
		ids := make([]string, 0, len(basis.mechanisms))
		for id := range basis.mechanisms {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			mechanism := basis.mechanisms[id]
			proposals = append(proposals, domain.ImprovementProposal{ProposalID: "rule:" + id, Title: mechanism.Name, TargetProperty: mechanism.Purpose, ExpectedEffect: mechanism.Strategy, MechanismRefs: []string{id}, ProductObservationRefs: []string{}, Assumptions: []string{"the target can expose a measurable baseline"}, Derivation: domain.DerivationRuleDerived, Rationale: "derived from a catalogued mechanism and an explicit target validation gap"})
		}
	}
	if len(proposals) > request.Policy.MaxHypotheses {
		proposals = proposals[:request.Policy.MaxHypotheses]
	}
	for _, proposal := range proposals {
		step := domain.ValidationStep{Question: "What is the current measurable state of " + proposal.TargetProperty + "?", Observable: proposal.TargetProperty, SafeProcedure: "measure the current target without modifying it", PassCriterion: "a repeatable baseline is recorded", RejectCriterion: "the baseline contradicts the expected effect", RelativeCost: "LOW", MissingEvidence: "target baseline"}
		status := domain.HypothesisProposed
		if basis.degraded {
			status = domain.HypothesisInconclusive
		}
		for _, ref := range proposal.MechanismRefs {
			mechanism := basis.mechanisms[ref]
			comparisons = append(comparisons, domain.ArchitectureComparison{TargetProperty: proposal.TargetProperty, TargetState: domain.TargetUnknown, TargetRef: "", MechanismID: ref, SourceRunID: sourceRunForMechanism(request.SourceArtifacts), Novelty: mechanism.Novelty, IntegrationUnknowns: compositionUnknowns(len(proposal.MechanismRefs))})
			if mechanism.Status != domain.MechanismProven {
				gaps = append(gaps, domain.HypothesisGap{Code: "MECHANISM_NOT_PROVEN", SourceRef: ref})
			}
		}
		hypothesis := domain.ImprovementHypothesis{
			Title: proposal.Title, TargetProperty: proposal.TargetProperty, TargetState: domain.TargetUnknown, TargetGapRefs: []string{}, ExpectedEffect: proposal.ExpectedEffect,
			MechanismRefs: append([]string{}, proposal.MechanismRefs...), ProductObservationRefs: append([]string{}, proposal.ProductObservationRefs...), TransferMode: domain.TransferReferenceOnly,
			License: domain.ImprovementLicenseUnrecognized, Assumptions: nonEmptySet(proposal.Assumptions, "the target baseline is not yet known"), Preconditions: []string{},
			IntegrationUnknowns: compositionUnknowns(len(proposal.MechanismRefs)), Risks: []string{"the expected effect may not transfer to the target"}, Tradeoffs: []string{},
			EvidenceRefs: append([]string{}, proposal.MechanismRefs...), CounterEvidenceRefs: []string{}, CounterEvidenceSearch: domain.CounterEvidenceSearch{Status: domain.CounterEvidenceNoneInScope, ScopeRefs: append([]string{}, scopeRefs...), ResultRefs: []string{}},
			ValidationPlan: []domain.ValidationStep{step}, Derivation: proposal.Derivation, Status: status, AdoptionDisposition: "NOT_AN_ADOPTION_DECISION",
		}
		hypothesis.ID = hypothesisID(hypothesis, scopeRefs)
		for index := range gaps {
			if gaps[index].HypothesisID == "" {
				gaps[index].HypothesisID = hypothesis.ID
			}
		}
		gaps = append(gaps, domain.HypothesisGap{HypothesisID: hypothesis.ID, Code: "LICENSE_NOT_RECOGNIZED", SourceRef: firstRef(proposal.MechanismRefs)})
		if basis.degraded {
			gaps = append(gaps, domain.HypothesisGap{HypothesisID: hypothesis.ID, Code: "SOURCE_RUN_NOT_SUCCEEDED", SourceRef: firstRef(scopeRefs)})
		}
		hypotheses = append(hypotheses, hypothesis)
		plans = append(plans, domain.ValidationPlanEntry{HypothesisID: hypothesis.ID, StepIndex: 0, Step: step})
		provenance = append(provenance, domain.SynthesisProvenance{HypothesisID: hypothesis.ID, Derivation: hypothesis.Derivation, RuleIDs: []string{}, ScopeRefs: append([]string{}, scopeRefs...)})
	}
	sort.Slice(comparisons, func(i, j int) bool {
		return strings.Join([]string{comparisons[i].TargetProperty, comparisons[i].MechanismID, string(comparisons[i].SourceRunID)}, "\x00") < strings.Join([]string{comparisons[j].TargetProperty, comparisons[j].MechanismID, string(comparisons[j].SourceRunID)}, "\x00")
	})
	sort.Slice(hypotheses, func(i, j int) bool { return hypotheses[i].ID < hypotheses[j].ID })
	sort.Slice(gaps, func(i, j int) bool {
		return strings.Join([]string{gaps[i].HypothesisID, gaps[i].Code, gaps[i].SourceRef}, "\x00") < strings.Join([]string{gaps[j].HypothesisID, gaps[j].Code, gaps[j].SourceRef}, "\x00")
	})
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].HypothesisID < plans[j].HypothesisID || plans[i].HypothesisID == plans[j].HypothesisID && plans[i].StepIndex < plans[j].StepIndex
	})
	status := domain.ImprovementComplete
	if len(hypotheses) == 0 {
		status = domain.ImprovementInconclusive
	} else if len(gaps) > 0 || basis.degraded {
		status = domain.ImprovementPartial
	}
	return improvementAnalysisResult{comparisons: comparisons, hypotheses: hypotheses, gaps: gaps, plans: plans, provenance: provenance, status: status}
}

func hypothesisID(value domain.ImprovementHypothesis, sourceHashes []string) string {
	preimage := struct {
		SchemaVersion          string                      `json:"schema_version"`
		Title                  string                      `json:"title"`
		TargetProperty         string                      `json:"target_property"`
		TargetGapRefs          []string                    `json:"target_gap_refs"`
		ExpectedEffect         string                      `json:"expected_effect"`
		MechanismRefs          []string                    `json:"mechanism_refs"`
		ProductObservationRefs []string                    `json:"product_observation_refs"`
		Derivation             domain.HypothesisDerivation `json:"derivation"`
		SourceHashes           []string                    `json:"source_hashes"`
	}{"1", value.Title, value.TargetProperty, canonicalMechanismSet(value.TargetGapRefs), value.ExpectedEffect, canonicalMechanismSet(value.MechanismRefs), canonicalMechanismSet(value.ProductObservationRefs), value.Derivation, canonicalMechanismSet(sourceHashes)}
	digest := sha256.Sum256(canonicalJSON(preimage))
	return "hyp_" + hex.EncodeToString(digest[:])
}

func compositionUnknowns(count int) []string {
	if count < 2 {
		return []string{}
	}
	return []string{"failure model compatibility", "interface compatibility", "license compatibility", "performance compatibility", "state model compatibility"}
}

func nonEmptySet(values []string, fallback string) []string {
	values = canonicalMechanismSet(values)
	if len(values) == 0 {
		return []string{fallback}
	}
	return values
}

func firstRef(values []string) string {
	if len(values) == 0 {
		return "target"
	}
	return values[0]
}

func sourceRunForMechanism(sources []domain.SourceArtifactRef) domain.RunID {
	for _, source := range sources {
		if source.Artifact == domain.ArtifactMechanismCatalog {
			return source.RunID
		}
	}
	return ""
}
