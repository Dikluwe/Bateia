/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_diagnosis.md
 * @layer L1
 * @updated 2026-08-26
 */

package usecase

import (
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type coverageUnit struct {
	repository string
	ruleID     string
	aliases    []string
	observed   domain.CapabilityEvidence
	found      bool
}

func DiagnoseEvidenceCoverage(report domain.AssessmentReport, bundles []domain.EvidenceBundle, judgments []domain.EvidenceJudgment, catalog domain.CapabilityRuleCatalog) domain.EvidenceCoverageReport {
	units := expectedProvenUnits(report, judgments, catalog)
	bundleByRepository := make(map[string]domain.EvidenceBundle, len(bundles))
	for _, bundle := range bundles {
		bundleByRepository[bundle.Repository.Owner+"/"+bundle.Repository.Repo] = bundle
	}
	diagnoses := make([]domain.EvidenceMissDiagnosis, 0)
	counts := map[domain.EvidenceMissCause]int{}
	emitted := 0
	for _, unit := range units {
		if unit.found && unit.observed.Status == domain.EvidenceProven {
			emitted++
			continue
		}
		diagnosis := diagnoseCoverageUnit(unit, bundleByRepository[unit.repository])
		diagnoses = append(diagnoses, diagnosis)
		counts[diagnosis.PrimaryCause]++
	}
	sort.Slice(diagnoses, func(i, j int) bool {
		if diagnoses[i].Repository != diagnoses[j].Repository {
			return diagnoses[i].Repository < diagnoses[j].Repository
		}
		if diagnoses[i].RuleID != diagnoses[j].RuleID {
			return diagnoses[i].RuleID < diagnoses[j].RuleID
		}
		return diagnoses[i].Capability < diagnoses[j].Capability
	})
	base := EvaluateEvidenceMetrics(report, bundles, judgments)
	base.ProvenCount = emitted
	base.SampleSufficient = emitted >= 10
	metrics := domain.EvidenceCoverageMetrics{ExpectedProvenCount: len(units), EmittedProvenCount: emitted, EvidenceMetrics: base, MissCountByCause: copyCauseCounts(counts)}
	if len(units) > 0 {
		recall := float64(emitted) / float64(len(units))
		metrics.ProvenRecall = &recall
	}
	return domain.EvidenceCoverageReport{SchemaVersion: "1", RulesVersion: report.RulesVersion, RulesHash: report.RulesHash, Diagnoses: diagnoses, CauseCounts: copyCauseCounts(counts), Metrics: metrics}
}

func expectedProvenUnits(report domain.AssessmentReport, judgments []domain.EvidenceJudgment, catalog domain.CapabilityRuleCatalog) []coverageUnit {
	observed := observedEvidenceByRepository(report)
	units := map[string]*coverageUnit{}
	sorted := append([]domain.EvidenceJudgment(nil), judgments...)
	sort.Slice(sorted, func(i, j int) bool {
		if repositoryIdentityKey(sorted[i].Repository) != repositoryIdentityKey(sorted[j].Repository) {
			return repositoryIdentityKey(sorted[i].Repository) < repositoryIdentityKey(sorted[j].Repository)
		}
		return normalizeEvidenceText(sorted[i].Capability) < normalizeEvidenceText(sorted[j].Capability)
	})
	for _, judgment := range sorted {
		if judgment.Expected != domain.EvidenceProven {
			continue
		}
		repository := repositoryIdentityKey(judgment.Repository)
		rule := findRule(judgment.Capability, catalog)
		ruleID := "unmapped:" + normalizeEvidenceText(judgment.Capability)
		if rule != nil {
			ruleID = rule.ID
		}
		key := repository + "|" + ruleID
		unit := units[key]
		if unit == nil {
			unit = &coverageUnit{repository: repository, ruleID: ruleID}
			units[key] = unit
		}
		if !containsNormalized(unit.aliases, judgment.Capability) {
			unit.aliases = append(unit.aliases, judgment.Capability)
		}
		if evidence, ok := observed[repository+"|"+normalizeEvidenceText(judgment.Capability)]; ok {
			if !unit.found || evidenceRank(evidence.Status) > evidenceRank(unit.observed.Status) {
				unit.observed, unit.found = evidence, true
			}
		}
	}
	out := make([]coverageUnit, 0, len(units))
	for _, unit := range units {
		sort.Slice(unit.aliases, func(i, j int) bool {
			return normalizeEvidenceText(unit.aliases[i]) < normalizeEvidenceText(unit.aliases[j])
		})
		out = append(out, *unit)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].repository != out[j].repository {
			return out[i].repository < out[j].repository
		}
		return out[i].ruleID < out[j].ruleID
	})
	return out
}

func observedEvidenceByRepository(report domain.AssessmentReport) map[string]domain.CapabilityEvidence {
	out := map[string]domain.CapabilityEvidence{}
	for _, assessment := range report.Assessments {
		repository := assessment.Repository.Owner + "/" + assessment.Repository.Repo
		set := append([]domain.CapabilityEvidence(nil), assessment.Capabilities...)
		if assessment.ProblemAlignment.Capability != "" {
			set = append(set, assessment.ProblemAlignment)
		}
		for _, evidence := range set {
			out[repository+"|"+normalizeEvidenceText(evidence.Capability)] = evidence
		}
	}
	return out
}

func diagnoseCoverageUnit(unit coverageUnit, bundle domain.EvidenceBundle) domain.EvidenceMissDiagnosis {
	capability := unit.ruleID
	if len(unit.aliases) > 0 {
		capability = unit.aliases[0]
	}
	observed := domain.EvidenceInconclusive
	var sources []domain.EvidenceSource
	if unit.found {
		observed, sources = unit.observed.Status, append([]domain.EvidenceSource(nil), unit.observed.Sources...)
	}
	diagnosis := domain.EvidenceMissDiagnosis{Repository: unit.repository, Capability: capability, Aliases: append([]string(nil), unit.aliases...), Expected: domain.EvidenceProven, Observed: observed, SupportingData: sources, RuleID: unit.ruleID}
	if bundle.Repository.Owner == "" {
		return withMiss(diagnosis, domain.MissTargetNotCollected, "bundle do repositório não foi fornecido")
	}
	for _, target := range bundle.Targets {
		if target.Status == domain.TargetSkippedBudget {
			diagnosis.MissingTargets = append(diagnosis.MissingTargets, target.Target.Path)
			return withMiss(diagnosis, domain.MissBudgetExhausted, "orçamento interrompeu alvo planejado")
		}
	}
	relevant := relevantTargets(bundle, unit.aliases)
	if len(relevant) == 0 {
		return withMiss(diagnosis, domain.MissTargetNotPlanned, "nenhum alvo do plano declara cobertura para a capacidade")
	}
	for _, target := range relevant {
		switch target.Status {
		case domain.TargetUnsupported:
			diagnosis.MissingTargets = append(diagnosis.MissingTargets, target.Target.Path)
			return withMiss(diagnosis, domain.MissLanguageUnsupported, "alvo relevante não possui suporte estrutural")
		case domain.TargetFailed, domain.TargetPlanned:
			diagnosis.MissingTargets = append(diagnosis.MissingTargets, target.Target.Path)
			return withMiss(diagnosis, domain.MissTargetNotCollected, "alvo relevante não terminou a coleta")
		}
	}
	if observed == domain.EvidenceConflict {
		return withMiss(diagnosis, domain.MissNegativeConflict, "evidência negativa explícita bloqueou a capacidade")
	}
	if !unit.found {
		return withMiss(diagnosis, domain.MissGroundTruthAmbiguous, "assessment não contém a capacidade julgada")
	}
	if len(sources) == 0 {
		return withMiss(diagnosis, domain.MissPattern, "coleta terminou sem fonte reconhecida pela regra")
	}
	for _, source := range sources {
		if (source.Kind == domain.EvidenceFile || source.Kind == domain.EvidenceExample) && source.Symbol == "" {
			return withMiss(diagnosis, domain.MissSymbolLinkage, "sinal em código não foi ligado a símbolo de implementação")
		}
	}
	if hasStrongKindWithWeakStrength(sources) {
		return withMiss(diagnosis, domain.MissStrengthUnderestimated, "fonte estrutural reconhecida recebeu força apenas auxiliar")
	}
	return withMiss(diagnosis, domain.MissInsufficientIndependence, "fontes reconhecidas não atingem independência e força exigidas")
}

func relevantTargets(bundle domain.EvidenceBundle, aliases []string) []domain.CollectionTargetResult {
	var out []domain.CollectionTargetResult
	for _, target := range bundle.Targets {
		for _, covered := range target.Target.Capabilities {
			for _, alias := range aliases {
				if normalizeEvidenceText(covered) == normalizeEvidenceText(alias) {
					out = append(out, target)
					break
				}
			}
		}
	}
	return out
}

func hasStrongKindWithWeakStrength(sources []domain.EvidenceSource) bool {
	for _, source := range sources {
		if (source.Kind == domain.EvidenceSymbol || source.Kind == domain.EvidenceTestImplementation || source.Kind == domain.EvidenceExample) && source.Strength < .65 {
			return true
		}
	}
	return false
}

func withMiss(diagnosis domain.EvidenceMissDiagnosis, cause domain.EvidenceMissCause, explanation string) domain.EvidenceMissDiagnosis {
	diagnosis.PrimaryCause, diagnosis.Explanation = cause, explanation
	sort.Strings(diagnosis.MissingTargets)
	return diagnosis
}

func containsNormalized(values []string, candidate string) bool {
	for _, value := range values {
		if normalizeEvidenceText(value) == normalizeEvidenceText(candidate) {
			return true
		}
	}
	return false
}

func evidenceRank(status domain.EvidenceStatus) int {
	switch status {
	case domain.EvidenceProven:
		return 5
	case domain.EvidenceConflict:
		return 4
	case domain.EvidencePartial:
		return 3
	case domain.EvidenceInconclusive:
		return 2
	case domain.EvidenceAbsent:
		return 1
	default:
		return 0
	}
}

func copyCauseCounts(in map[domain.EvidenceMissCause]int) map[domain.EvidenceMissCause]int {
	out := make(map[domain.EvidenceMissCause]int, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func ValidateEvidenceCoverageInputs(report domain.AssessmentReport, judgments []domain.EvidenceJudgment) error {
	if strings.TrimSpace(report.SchemaVersion) == "" {
		return fmt.Errorf("assessment sem schema_version")
	}
	if len(judgments) == 0 {
		return fmt.Errorf("julgamentos vazios")
	}
	return nil
}
