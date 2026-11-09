/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/software_alternative.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type alternativeGroup struct {
	key         string
	snapshots   []domain.RepositoryCandidate
	assessments []domain.CandidateAssessment
	locales     []string
}

const unknownLicenseWarning = "licença não reconhecida; valide os termos antes da adoção"

func RankSoftwareAlternatives(target domain.SoftwareAlternativeTarget, idea domain.IdeaSpec, discovery domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, policy domain.ReuseDecisionPolicy) (domain.SoftwareAlternativeReport, error) {
	capabilities := normalizedSorted(idea.Capabilities)
	if err := validateSoftwareAlternativeInputs(target, capabilities, policy); err != nil {
		return domain.SoftwareAlternativeReport{}, err
	}
	target = canonicalSoftwareTarget(target)
	report := domain.SoftwareAlternativeReport{
		SchemaVersion: "1",
		Target:        target,
		Filters: domain.SoftwareAlternativeFilters{
			AllowedLicenses:        normalizedLicenses(policy.AllowedLicenses),
			RequireLocal:           policy.RequireLocal,
			AllowArchived:          policy.AllowArchived,
			RecommendedMinCoverage: policy.ExtendMinCoverage,
			PartialMinCoverage:     policy.IntegrationMinValue,
		},
		Alternatives: []domain.SoftwareAlternative{},
		Excluded:     []domain.SoftwareAlternative{},
	}

	groups := map[string]*alternativeGroup{}
	for _, discovered := range discovery.Candidates {
		candidate := discovered.Repository
		key := alternativeRepositoryKey(candidate.Owner, candidate.Repo)
		group := ensureAlternativeGroup(groups, key)
		group.snapshots = append(group.snapshots, candidate)
		for _, origin := range discovered.Origins {
			if origin.SearchLocale != "" {
				group.locales = append(group.locales, origin.SearchLocale)
			}
		}
	}
	for _, assessment := range assessments {
		assessment = SanitizeAssessmentEvidence(assessment)
		key := alternativeRepositoryKey(assessment.Repository.Owner, assessment.Repository.Repo)
		group := ensureAlternativeGroup(groups, key)
		group.assessments = append(group.assessments, assessment)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		alternative := inspectSoftwareAlternative(*groups[key], capabilities, report.Filters, target)
		switch alternative.Fit {
		case domain.AlternativeRecommended, domain.AlternativePartial, domain.AlternativeReference:
			report.Alternatives = append(report.Alternatives, alternative)
		default:
			report.Excluded = append(report.Excluded, alternative)
		}
	}
	sort.Slice(report.Alternatives, func(i, j int) bool {
		return softwareAlternativeLess(report.Alternatives[i], report.Alternatives[j], false)
	})
	sort.Slice(report.Excluded, func(i, j int) bool { return softwareAlternativeLess(report.Excluded[i], report.Excluded[j], true) })
	for _, alternative := range report.Alternatives {
		if alternative.Fit == domain.AlternativeRecommended || alternative.Fit == domain.AlternativePartial {
			choice := alternative.Repository
			report.TopChoice = &choice
			break
		}
	}
	return report, nil
}

func validateSoftwareAlternativeInputs(target domain.SoftwareAlternativeTarget, capabilities []string, policy domain.ReuseDecisionPolicy) error {
	if strings.TrimSpace(target.Name) == "" {
		return fmt.Errorf("target.name é obrigatório")
	}
	if len(capabilities) == 0 {
		return fmt.Errorf("idea exige ao menos uma capacidade obrigatória")
	}
	if policy.SchemaVersion == "" {
		return fmt.Errorf("policy.schema_version é obrigatório")
	}
	if invalidCoverageThreshold(policy.ExtendMinCoverage) || invalidCoverageThreshold(policy.IntegrationMinValue) {
		return fmt.Errorf("limiares da política devem estar em (0,1]")
	}
	if policy.IntegrationMinValue > policy.ExtendMinCoverage {
		return fmt.Errorf("policy.integration_min_value não pode exceder extend_min_coverage")
	}
	if target.Repository != nil {
		owner := strings.TrimSpace(target.Repository.Owner)
		repo := strings.TrimSpace(target.Repository.Repo)
		if owner == "" || repo == "" || strings.Contains(owner, "/") || strings.Contains(repo, "/") {
			return fmt.Errorf("target.repository deve usar owner/repo")
		}
	}
	return nil
}

func invalidCoverageThreshold(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 || value > 1
}

func canonicalSoftwareTarget(target domain.SoftwareAlternativeTarget) domain.SoftwareAlternativeTarget {
	target.Name = strings.Join(strings.Fields(target.Name), " ")
	if target.Repository != nil {
		repository := domain.RepositoryIdentity{Owner: strings.ToLower(strings.TrimSpace(target.Repository.Owner)), Repo: strings.ToLower(strings.TrimSpace(target.Repository.Repo))}
		target.Repository = &repository
	}
	return target
}

func normalizedLicenses(values []string) []string {
	byNormalized := map[string]string{}
	for _, value := range values {
		value = strings.Join(strings.Fields(value), " ")
		normalized := strings.ToLower(value)
		if value != "" && (byNormalized[normalized] == "" || value < byNormalized[normalized]) {
			byNormalized[normalized] = value
		}
	}
	keys := make([]string, 0, len(byNormalized))
	for key := range byNormalized {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, byNormalized[key])
	}
	return out
}

func alternativeRepositoryKey(owner, repo string) string {
	return strings.ToLower(strings.TrimSpace(owner)) + "/" + strings.ToLower(strings.TrimSpace(repo))
}

func ensureAlternativeGroup(groups map[string]*alternativeGroup, key string) *alternativeGroup {
	if group := groups[key]; group != nil {
		return group
	}
	group := &alternativeGroup{key: key}
	groups[key] = group
	return group
}

func inspectSoftwareAlternative(group alternativeGroup, required []string, filters domain.SoftwareAlternativeFilters, target domain.SoftwareAlternativeTarget) domain.SoftwareAlternative {
	sort.Slice(group.snapshots, func(i, j int) bool { return directCandidateMetadataLess(group.snapshots[i], group.snapshots[j]) })
	sort.Slice(group.assessments, func(i, j int) bool {
		return commitRefLess(group.assessments[i].Repository, group.assessments[j].Repository)
	})
	repository := canonicalGroupRepository(group)
	alternative := domain.SoftwareAlternative{Repository: repository, Capabilities: []domain.SoftwareAlternativeCapability{}, Gaps: []string{}, Reasons: []string{}, Blockers: []string{}, Evidence: []domain.EvidenceSource{}}
	if locales := canonicalAlternativeLocales(group.locales); len(locales) > 0 {
		alternative.DiscoveryLocales = locales
	}
	opaque := false
	if len(group.snapshots) > 0 {
		metadata := group.snapshots[0]
		alternative.URL = metadata.URL
		alternative.Description = metadata.Description
		alternative.Stars = metadata.Stars
		if metadata.Stale {
			opaque = true
			alternative.Reasons = append(alternative.Reasons, "identidade local não verificável para "+group.key)
		}
	}
	isTarget := target.Repository != nil && group.key == target.Repository.Key()
	if len(group.snapshots) != 1 {
		opaque = true
		if len(group.snapshots) == 0 {
			alternative.Reasons = append(alternative.Reasons, "snapshot ausente para "+group.key)
		} else {
			alternative.Reasons = append(alternative.Reasons, "snapshot duplicado para "+group.key)
			if candidateMetadataDiverges(group.snapshots) {
				alternative.Reasons = append(alternative.Reasons, "metadados divergentes para "+group.key)
			}
		}
	}
	if len(group.assessments) != 1 {
		opaque = true
		if len(group.assessments) == 0 {
			alternative.Reasons = append(alternative.Reasons, "assessment ausente para "+group.key)
		} else {
			alternative.Reasons = append(alternative.Reasons, "assessment duplicado para "+group.key)
			if alternativeAssessmentsDiverge(group.assessments) {
				alternative.Reasons = append(alternative.Reasons, "assessments divergentes para "+group.key)
			}
		}
	}
	if repository.Commit == "" || strings.EqualFold(repository.Commit, "HEAD") {
		opaque = true
		alternative.Reasons = append(alternative.Reasons, "identidade mutável para "+group.key)
	}
	if strings.TrimSpace(repository.Owner) == "" || strings.TrimSpace(repository.Repo) == "" {
		opaque = true
		alternative.Reasons = append(alternative.Reasons, "identidade de repositório incompleta")
	}
	if len(group.snapshots) > 0 && len(group.assessments) > 0 && alternativeCommitSet(group.snapshots) != assessmentCommitSet(group.assessments) {
		opaque = true
		alternative.Reasons = append(alternative.Reasons, fmt.Sprintf("commit divergente para %s: snapshot [%s], assessment [%s]", group.key, alternativeCommitSet(group.snapshots), assessmentCommitSet(group.assessments)))
	}

	if len(group.assessments) > 0 {
		assessment := canonicalAssessment(group.assessments)
		alternative.License = canonicalLicenseAssessment(assessment.License)
		alternative.Maintenance = canonicalMaintenanceAssessment(assessment.Maintenance)
		alternative.Operation = canonicalOperationAssessment(assessment.Operation)
		alternative.ProblemAlignment = assessment.ProblemAlignment.Status
		alternative.Evidence = append(alternative.Evidence, assessment.ProblemAlignment.Sources...)
		alternative.Evidence = append(alternative.Evidence, assessment.License.Sources...)
		alternative.Evidence = append(alternative.Evidence, assessment.Maintenance.Sources...)
		alternative.Evidence = append(alternative.Evidence, assessment.Operation.Sources...)
		capabilities, gaps, covered, assessable, signal, capabilityOpaque, blockers, evidence, reasons := inspectAlternativeCapabilities(assessment, required)
		alternative.Capabilities = capabilities
		alternative.Gaps = gaps
		if assessable > 0 {
			alternative.Coverage = float64(covered) / float64(assessable)
			alternative.Reasons = append(alternative.Reasons, fmt.Sprintf("cobertura comprovada: %d/%d", covered, assessable))
		} else {
			alternative.Coverage = 0
			capabilityOpaque = true
			alternative.Reasons = append(alternative.Reasons, "cobertura não pontuável: o catálogo não possui regra declarada segura para as capacidades requeridas")
		}
		alternative.Blockers = append(alternative.Blockers, blockers...)
		alternative.Reasons = append(alternative.Reasons, reasons...)
		alternative.Evidence = append(alternative.Evidence, evidence...)
		opaque = opaque || capabilityOpaque
		gateOpaque, gateBlockers, gateReasons := inspectAlternativeGates(assessment, filters)
		opaque = opaque || gateOpaque
		alternative.Blockers = append(alternative.Blockers, gateBlockers...)
		alternative.Reasons = append(alternative.Reasons, gateReasons...)
		if len(assessment.Failures) > 0 {
			opaque = true
			for _, failure := range assessment.Failures {
				alternative.Reasons = append(alternative.Reasons, fmt.Sprintf("falha de coleta %s: %s", failure.Stage, failure.Message))
			}
		}
		problemSignal := false
		if catalogLimitedEvidence(assessment.ProblemAlignment) {
			opaque = true
			alternative.Reasons = append(alternative.Reasons, "alinhamento não avaliável pela cobertura atual do catálogo; isso não descreve o candidato")
		} else {
			problemSignal = assessment.ProblemAlignment.Status == domain.EvidenceProven || assessment.ProblemAlignment.Status == domain.EvidencePartial
			switch assessment.ProblemAlignment.Status {
			case domain.EvidenceProven, domain.EvidencePartial, domain.EvidenceAbsent:
			case domain.EvidenceConflict:
				alternative.Blockers = append(alternative.Blockers, "alinhamento do problema em conflito")
			case domain.EvidenceInconclusive:
				opaque = true
				alternative.Reasons = append(alternative.Reasons, "alinhamento do problema inconclusivo")
			default:
				opaque = true
				alternative.Reasons = append(alternative.Reasons, "estado de alinhamento do problema desconhecido")
			}
		}
		alternative.Fit = deriveSoftwareAlternativeFit(isTarget, alternative, opaque, signal || problemSignal, filters)
	} else {
		alternative.ProblemAlignment = domain.EvidenceInconclusive
		for _, capability := range required {
			alternative.Capabilities = append(alternative.Capabilities, domain.SoftwareAlternativeCapability{Capability: capability, Status: domain.EvidenceInconclusive, Evidence: []domain.EvidenceSource{}})
			alternative.Gaps = append(alternative.Gaps, capability)
		}
		alternative.Fit = deriveSoftwareAlternativeFit(isTarget, alternative, true, false, filters)
	}
	canonicalizeSoftwareAlternative(&alternative)
	return alternative
}

func canonicalAlternativeLocales(values []string) []string {
	seen := map[string]bool{}
	for _, value := range values {
		seen[strings.TrimSpace(value)] = true
	}
	out := make([]string, 0, len(seen))
	for _, locale := range canonicalSearchLocales {
		if seen[locale] {
			out = append(out, locale)
			delete(seen, locale)
		}
	}
	unknown := make([]string, 0, len(seen))
	for locale := range seen {
		if locale != "" {
			unknown = append(unknown, locale)
		}
	}
	sort.Strings(unknown)
	return append(out, unknown...)
}

func canonicalGroupRepository(group alternativeGroup) domain.CommitRef {
	references := make([]domain.CommitRef, 0, len(group.snapshots)+len(group.assessments))
	for _, candidate := range group.snapshots {
		references = append(references, domain.CommitRef{Owner: strings.ToLower(strings.TrimSpace(candidate.Owner)), Repo: strings.ToLower(strings.TrimSpace(candidate.Repo)), Commit: candidate.CommitID})
	}
	for _, assessment := range group.assessments {
		references = append(references, domain.CommitRef{Owner: strings.ToLower(strings.TrimSpace(assessment.Repository.Owner)), Repo: strings.ToLower(strings.TrimSpace(assessment.Repository.Repo)), Commit: assessment.Repository.Commit})
	}
	sort.Slice(references, func(i, j int) bool { return commitRefLess(references[i], references[j]) })
	if len(references) == 0 {
		parts := strings.SplitN(group.key, "/", 2)
		return domain.CommitRef{Owner: parts[0], Repo: parts[1]}
	}
	return references[0]
}

func canonicalAssessment(assessments []domain.CandidateAssessment) domain.CandidateAssessment {
	copyAssessments := append([]domain.CandidateAssessment(nil), assessments...)
	for i := range copyAssessments {
		copyAssessments[i].Capabilities = append([]domain.CapabilityEvidence(nil), copyAssessments[i].Capabilities...)
		copyAssessments[i].ProblemAlignment.Sources = canonicalEvidenceSources(copyAssessments[i].ProblemAlignment.Sources)
		for j := range copyAssessments[i].Capabilities {
			copyAssessments[i].Capabilities[j].Sources = canonicalEvidenceSources(copyAssessments[i].Capabilities[j].Sources)
		}
		sort.Slice(copyAssessments[i].Capabilities, func(left, right int) bool {
			leftKey := normalizeDecisionText(copyAssessments[i].Capabilities[left].Capability)
			rightKey := normalizeDecisionText(copyAssessments[i].Capabilities[right].Capability)
			if leftKey != rightKey {
				return leftKey < rightKey
			}
			leftJSON, _ := json.Marshal(copyAssessments[i].Capabilities[left])
			rightJSON, _ := json.Marshal(copyAssessments[i].Capabilities[right])
			return string(leftJSON) < string(rightJSON)
		})
		copyAssessments[i].License = canonicalLicenseAssessment(copyAssessments[i].License)
		copyAssessments[i].Maintenance = canonicalMaintenanceAssessment(copyAssessments[i].Maintenance)
		copyAssessments[i].Operation = canonicalOperationAssessment(copyAssessments[i].Operation)
		copyAssessments[i].Failures = canonicalExecutionFailures(copyAssessments[i].Failures)
	}
	sort.Slice(copyAssessments, func(i, j int) bool {
		left, _ := json.Marshal(alternativeAssessmentKey(copyAssessments[i]))
		right, _ := json.Marshal(alternativeAssessmentKey(copyAssessments[j]))
		return string(left) < string(right)
	})
	return copyAssessments[0]
}

func candidateMetadataDiverges(candidates []domain.RepositoryCandidate) bool {
	if len(candidates) < 2 {
		return false
	}
	first, _ := json.Marshal(candidates[0])
	for _, candidate := range candidates[1:] {
		encoded, _ := json.Marshal(candidate)
		if string(encoded) != string(first) {
			return true
		}
	}
	return false
}

func alternativeAssessmentsDiverge(assessments []domain.CandidateAssessment) bool {
	if len(assessments) < 2 {
		return false
	}
	first, _ := json.Marshal(alternativeAssessmentKey(canonicalAssessment([]domain.CandidateAssessment{assessments[0]})))
	for _, assessment := range assessments[1:] {
		encoded, _ := json.Marshal(alternativeAssessmentKey(canonicalAssessment([]domain.CandidateAssessment{assessment})))
		if string(encoded) != string(first) {
			return true
		}
	}
	return false
}

func alternativeCommitSet(candidates []domain.RepositoryCandidate) string {
	values := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		values = append(values, candidate.CommitID)
	}
	return canonicalCommitSet(values)
}

func assessmentCommitSet(assessments []domain.CandidateAssessment) string {
	values := make([]string, 0, len(assessments))
	for _, assessment := range assessments {
		values = append(values, assessment.Repository.Commit)
	}
	return canonicalCommitSet(values)
}

func canonicalCommitSet(values []string) string {
	seen := map[string]bool{}
	for _, value := range values {
		seen[value] = true
	}
	canonical := make([]string, 0, len(seen))
	for value := range seen {
		canonical = append(canonical, fmt.Sprintf("%q", value))
	}
	sort.Strings(canonical)
	return strings.Join(canonical, ", ")
}

func alternativeAssessmentKey(assessment domain.CandidateAssessment) any {
	return struct {
		Repository       domain.CommitRef             `json:"repository"`
		ProblemAlignment domain.CapabilityEvidence    `json:"problem_alignment"`
		Capabilities     []domain.CapabilityEvidence  `json:"capabilities"`
		License          domain.LicenseAssessment     `json:"license"`
		Maintenance      domain.MaintenanceAssessment `json:"maintenance"`
		Operation        domain.OperationAssessment   `json:"operation"`
		Failures         []domain.ExecutionFailure    `json:"failures"`
	}{
		Repository:       assessment.Repository,
		ProblemAlignment: assessment.ProblemAlignment,
		Capabilities:     assessment.Capabilities,
		License:          assessment.License,
		Maintenance:      assessment.Maintenance,
		Operation:        assessment.Operation,
		Failures:         assessment.Failures,
	}
}

func canonicalExecutionFailures(values []domain.ExecutionFailure) []domain.ExecutionFailure {
	byKey := map[string]domain.ExecutionFailure{}
	for _, value := range values {
		encoded, _ := json.Marshal(value)
		byKey[string(encoded)] = value
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.ExecutionFailure, 0, len(keys))
	for _, key := range keys {
		out = append(out, byKey[key])
	}
	return out
}

func inspectAlternativeCapabilities(assessment domain.CandidateAssessment, required []string) ([]domain.SoftwareAlternativeCapability, []string, int, int, bool, bool, []string, []domain.EvidenceSource, []string) {
	byCapability := map[string][]domain.CapabilityEvidence{}
	for _, evidence := range assessment.Capabilities {
		key := normalizeDecisionText(evidence.Capability)
		if key != "" {
			byCapability[key] = append(byCapability[key], evidence)
		}
	}
	statuses := make([]domain.SoftwareAlternativeCapability, 0, len(required))
	gaps := make([]string, 0, len(required))
	covered := 0
	assessable := 0
	signal := false
	opaque := false
	var blockers, reasons []string
	var allEvidence []domain.EvidenceSource
	for _, capability := range required {
		entries := byCapability[capability]
		status := domain.EvidenceInconclusive
		var sources []domain.EvidenceSource
		switch len(entries) {
		case 0:
			opaque = true
			reasons = append(reasons, "capacidade omitida: "+capability)
		case 1:
			status = entries[0].Status
			sources = entries[0].Sources
		default:
			opaque = true
			reasons = append(reasons, "capacidade duplicada: "+capability)
			for _, entry := range entries {
				sources = append(sources, entry.Sources...)
			}
		}
		sources = canonicalEvidenceSources(sources)
		statuses = append(statuses, domain.SoftwareAlternativeCapability{Capability: capability, Status: status, Evidence: sources})
		allEvidence = append(allEvidence, sources...)
		if len(entries) == 1 && catalogLimitedEvidence(entries[0]) {
			opaque = true
			reasons = append(reasons, "capacidade não avaliável pela cobertura atual do catálogo: "+capability+"; isso não descreve lacuna do candidato")
			continue
		}
		assessable++
		switch status {
		case domain.EvidenceProven:
			covered++
			signal = true
		case domain.EvidencePartial:
			signal = true
		case domain.EvidenceAbsent:
		case domain.EvidenceConflict:
			blockers = append(blockers, "capacidade em conflito: "+capability)
		case domain.EvidenceInconclusive:
			opaque = true
		default:
			opaque = true
			reasons = append(reasons, "estado de capacidade desconhecido: "+capability)
		}
		if status != domain.EvidenceProven {
			gaps = append(gaps, capability)
		}
	}
	return statuses, gaps, covered, assessable, signal, opaque, blockers, allEvidence, reasons
}

func inspectAlternativeGates(assessment domain.CandidateAssessment, filters domain.SoftwareAlternativeFilters) (bool, []string, []string) {
	opaque := false
	var blockers, reasons []string
	switch assessment.License.Status {
	case domain.EvidenceProven:
		if len(filters.AllowedLicenses) > 0 && !containsFold(filters.AllowedLicenses, assessment.License.Identifier) {
			blockers = append(blockers, "licença não permitida: "+assessment.License.Identifier)
		}
	case domain.EvidenceConflict:
		blockers = append(blockers, "licença em conflito")
	case domain.EvidencePartial, domain.EvidenceAbsent, domain.EvidenceInconclusive:
		reasons = append(reasons, unknownLicenseWarning)
	default:
		opaque = true
		reasons = append(reasons, "estado de licença desconhecido")
	}
	switch assessment.Maintenance.Status {
	case domain.EvidenceProven:
		if assessment.Maintenance.Archived && !filters.AllowArchived {
			blockers = append(blockers, "repositório arquivado não permitido")
		}
	case domain.EvidenceAbsent:
		blockers = append(blockers, "manutenção ausente")
	case domain.EvidenceConflict:
		blockers = append(blockers, "manutenção em conflito")
	case domain.EvidencePartial, domain.EvidenceInconclusive:
		opaque = true
		reasons = append(reasons, "manutenção inconclusiva")
	default:
		opaque = true
		reasons = append(reasons, "estado de manutenção desconhecido")
	}
	if filters.RequireLocal {
		switch assessment.Operation.LocalStatus {
		case domain.EvidenceProven:
		case domain.EvidenceAbsent:
			blockers = append(blockers, "operação local ausente")
		case domain.EvidenceConflict:
			blockers = append(blockers, "operação local em conflito")
		case domain.EvidencePartial, domain.EvidenceInconclusive:
			opaque = true
			reasons = append(reasons, "operação local inconclusiva")
		default:
			opaque = true
			reasons = append(reasons, "estado de operação local desconhecido")
		}
	}
	return opaque, blockers, reasons
}

func deriveSoftwareAlternativeFit(isTarget bool, alternative domain.SoftwareAlternative, opaque, signal bool, filters domain.SoftwareAlternativeFilters) domain.SoftwareAlternativeFit {
	if isTarget {
		return domain.AlternativeNotAlternative
	}
	if len(alternative.Blockers) > 0 {
		return domain.AlternativeBlocked
	}
	if opaque {
		return domain.AlternativeInconclusive
	}
	if alternative.ProblemAlignment == domain.EvidenceAbsent {
		return domain.AlternativeNotAlternative
	}
	if alternative.ProblemAlignment == domain.EvidenceProven && alternative.Coverage >= filters.RecommendedMinCoverage {
		return domain.AlternativeRecommended
	}
	if alternative.ProblemAlignment == domain.EvidenceProven && alternative.Coverage >= filters.PartialMinCoverage {
		return domain.AlternativePartial
	}
	if (alternative.ProblemAlignment == domain.EvidenceProven || alternative.ProblemAlignment == domain.EvidencePartial) && signal {
		return domain.AlternativeReference
	}
	return domain.AlternativeNotAlternative
}

func canonicalizeSoftwareAlternative(alternative *domain.SoftwareAlternative) {
	alternative.Reasons = uniqueSortedAlternativeStrings(alternative.Reasons)
	alternative.Blockers = uniqueSortedAlternativeStrings(alternative.Blockers)
	alternative.Gaps = uniqueSortedAlternativeStrings(alternative.Gaps)
	alternative.Evidence = canonicalEvidenceSources(alternative.Evidence)
	alternative.License = canonicalLicenseAssessment(alternative.License)
	alternative.Maintenance = canonicalMaintenanceAssessment(alternative.Maintenance)
	alternative.Operation = canonicalOperationAssessment(alternative.Operation)
}

func canonicalLicenseAssessment(value domain.LicenseAssessment) domain.LicenseAssessment {
	value.Sources = canonicalEvidenceSources(value.Sources)
	return value
}

func canonicalMaintenanceAssessment(value domain.MaintenanceAssessment) domain.MaintenanceAssessment {
	value.Sources = canonicalEvidenceSources(value.Sources)
	value.Warnings = uniqueSortedAlternativeStrings(value.Warnings)
	return value
}

func canonicalOperationAssessment(value domain.OperationAssessment) domain.OperationAssessment {
	value.Sources = canonicalEvidenceSources(value.Sources)
	value.Requirements = uniqueSortedAlternativeStrings(value.Requirements)
	return value
}

func canonicalEvidenceSources(values []domain.EvidenceSource) []domain.EvidenceSource {
	byKey := map[string]domain.EvidenceSource{}
	for _, value := range values {
		encoded, _ := json.Marshal(value)
		byKey[string(encoded)] = value
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.EvidenceSource, 0, len(keys))
	for _, key := range keys {
		out = append(out, byKey[key])
	}
	return out
}

func uniqueSortedAlternativeStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.Join(strings.Fields(value), " ")
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func candidateMetadataLess(left, right domain.RepositoryCandidate) bool {
	leftText := []string{strings.TrimSpace(left.Owner), strings.TrimSpace(left.Repo), left.CommitID, left.URL, left.Description}
	rightText := []string{strings.TrimSpace(right.Owner), strings.TrimSpace(right.Repo), right.CommitID, right.URL, right.Description}
	for i := range leftText {
		if leftText[i] != rightText[i] {
			return leftText[i] < rightText[i]
		}
	}
	return left.Stars < right.Stars
}

func commitRefLess(left, right domain.CommitRef) bool {
	leftKey := strings.ToLower(strings.TrimSpace(left.Owner)) + "/" + strings.ToLower(strings.TrimSpace(left.Repo)) + "@" + left.Commit
	rightKey := strings.ToLower(strings.TrimSpace(right.Owner)) + "/" + strings.ToLower(strings.TrimSpace(right.Repo)) + "@" + right.Commit
	return leftKey < rightKey
}

func softwareAlternativeLess(left, right domain.SoftwareAlternative, excluded bool) bool {
	leftFit := alternativeFitRank(left.Fit, excluded)
	rightFit := alternativeFitRank(right.Fit, excluded)
	if leftFit != rightFit {
		return leftFit < rightFit
	}
	if left.Coverage != right.Coverage {
		return left.Coverage > right.Coverage
	}
	leftProblem := alternativeEvidenceRank(left.ProblemAlignment)
	rightProblem := alternativeEvidenceRank(right.ProblemAlignment)
	if leftProblem != rightProblem {
		return leftProblem < rightProblem
	}
	if left.Stars != right.Stars {
		return left.Stars > right.Stars
	}
	return commitRefLess(left.Repository, right.Repository)
}

func alternativeFitRank(fit domain.SoftwareAlternativeFit, excluded bool) int {
	if excluded {
		switch fit {
		case domain.AlternativeBlocked:
			return 0
		case domain.AlternativeInconclusive:
			return 1
		default:
			return 2
		}
	}
	switch fit {
	case domain.AlternativeRecommended:
		return 0
	case domain.AlternativePartial:
		return 1
	default:
		return 2
	}
}

func alternativeEvidenceRank(status domain.EvidenceStatus) int {
	switch status {
	case domain.EvidenceProven:
		return 0
	case domain.EvidencePartial:
		return 1
	case domain.EvidenceAbsent:
		return 2
	case domain.EvidenceConflict:
		return 3
	default:
		return 4
	}
}
