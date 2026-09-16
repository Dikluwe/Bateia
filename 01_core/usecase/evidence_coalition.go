/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_coalition.md
 * @layer L1
 * @updated 2026-09-11
 */

package usecase

import (
	"sort"
	"strings"

	"bateia/01_core/domain"
)

const (
	maxCoalitionRules             = 128
	maxCoalitionClaims            = 128
	maxCoalitionWitnessesPerClaim = 64
	maxCoalitionWitnessesTotal    = 4096
	maxCoalitionRolesPerRule      = 4
	maxCoalitionOperations        = 8192
)

type coalitionRule struct {
	rule  domain.EvidenceCoalitionRule
	valid bool
}

type coalitionSourceKey struct {
	kind        domain.EvidenceKind
	owner       string
	repo        string
	commit      string
	path        string
	symbol      string
	line        int
	ruleID      string
	contentHash string
}

// AnalyzeEvidenceCoalitions derives a canonical shadow report without changing any
// existing evidence decision or score.
func AnalyzeEvidenceCoalitions(input domain.EvidenceCoalitionInput) domain.EvidenceCoalitionReport {
	report := emptyCoalitionReport()
	if input.SchemaVersion != domain.EvidenceCoalitionInputSchemaVersion {
		report.Reason = "invalid_schema"
		return report
	}
	if len(input.Claims) == 0 {
		report.Reason = "no_claims"
		return report
	}
	if exceedsCoalitionLimits(input) {
		report.Reason = "budget_exceeded"
		report.Metrics.LimitExceeded = true
		return report
	}

	rules, ambiguous := indexCoalitionRules(input.Rules)
	claims := append([]domain.EvidenceCoalitionClaim(nil), input.Claims...)
	sort.Slice(claims, func(i, j int) bool { return claims[i].ID < claims[j].ID })
	if ambiguous || hasAmbiguousCoalitionClaims(claims) {
		report.Reason = "ambiguous_input"
		return report
	}

	report.Status = domain.EvidenceCoalitionReportComplete
	report.Reason = "analyzed"
	report.Metrics.Claims = len(claims)
	for _, claim := range claims {
		assessment, witnesses := assessEvidenceCoalition(claim, rules)
		report.Assessments = append(report.Assessments, assessment)
		report.Metrics.Witnesses += witnesses
		report.Metrics.Operations += witnesses
		if indexed, ok := rules[claim.RuleID]; ok && indexed.valid {
			report.Metrics.Operations += len(indexed.rule.RequiredRoles)
		}
		countCoalitionOutcome(&report.Metrics, assessment.Outcome)
	}
	if report.Metrics.Invalid != 0 {
		report.Reason = "analyzed_with_invalid_claims"
	}
	if report.Metrics.Operations > maxCoalitionOperations {
		return budgetExceededCoalitionReport()
	}
	return report
}

func emptyCoalitionReport() domain.EvidenceCoalitionReport {
	return domain.EvidenceCoalitionReport{
		SchemaVersion: domain.EvidenceCoalitionReportSchemaVersion,
		Status:        domain.EvidenceCoalitionReportUnknown,
		Assessments:   []domain.EvidenceCoalitionAssessment{},
	}
}

func budgetExceededCoalitionReport() domain.EvidenceCoalitionReport {
	report := emptyCoalitionReport()
	report.Reason = "budget_exceeded"
	report.Metrics.LimitExceeded = true
	return report
}

func exceedsCoalitionLimits(input domain.EvidenceCoalitionInput) bool {
	if len(input.Rules) > maxCoalitionRules || len(input.Claims) > maxCoalitionClaims {
		return true
	}
	for _, rule := range input.Rules {
		if len(rule.RequiredRoles) > maxCoalitionRolesPerRule {
			return true
		}
	}
	total := 0
	for _, claim := range input.Claims {
		if len(claim.Witnesses) > maxCoalitionWitnessesPerClaim {
			return true
		}
		total += len(claim.Witnesses)
		if total > maxCoalitionWitnessesTotal {
			return true
		}
	}
	return false
}

func indexCoalitionRules(input []domain.EvidenceCoalitionRule) (map[string]coalitionRule, bool) {
	rules := make(map[string]coalitionRule, len(input))
	for _, candidate := range input {
		if _, exists := rules[candidate.ID]; exists {
			return nil, true
		}
		rule := domain.EvidenceCoalitionRule{
			ID:                     candidate.ID,
			RequiredRoles:          append([]domain.EvidenceCoalitionRole(nil), candidate.RequiredRoles...),
			MinimumDistinctOrigins: candidate.MinimumDistinctOrigins,
		}
		sort.Slice(rule.RequiredRoles, func(i, j int) bool { return rule.RequiredRoles[i] < rule.RequiredRoles[j] })
		rules[rule.ID] = coalitionRule{rule: rule, valid: validCoalitionRule(rule)}
	}
	return rules, false
}

func validCoalitionRule(rule domain.EvidenceCoalitionRule) bool {
	if !validCoalitionID(rule.ID) || len(rule.RequiredRoles) == 0 || len(rule.RequiredRoles) > maxCoalitionRolesPerRule {
		return false
	}
	if rule.MinimumDistinctOrigins < 1 || rule.MinimumDistinctOrigins > len(rule.RequiredRoles) {
		return false
	}
	seen := make(map[domain.EvidenceCoalitionRole]struct{}, len(rule.RequiredRoles))
	for _, role := range rule.RequiredRoles {
		if !validCoalitionRole(role) {
			return false
		}
		if _, exists := seen[role]; exists {
			return false
		}
		seen[role] = struct{}{}
	}
	return true
}

func hasAmbiguousCoalitionClaims(claims []domain.EvidenceCoalitionClaim) bool {
	for index, claim := range claims {
		if !validCoalitionID(claim.ID) || (index > 0 && claims[index-1].ID == claim.ID) {
			return true
		}
	}
	return false
}

func assessEvidenceCoalition(claim domain.EvidenceCoalitionClaim, rules map[string]coalitionRule) (domain.EvidenceCoalitionAssessment, int) {
	assessment := domain.EvidenceCoalitionAssessment{
		ClaimID:               claim.ID,
		RuleID:                claim.RuleID,
		BaselineStatus:        claim.BaselineStatus,
		DecisionImpact:        domain.EvidenceCoalitionDecisionImpactNone,
		Outcome:               domain.EvidenceCoalitionInvalidInput,
		Independence:          domain.EvidenceOriginUnknown,
		SupportingWitnessIDs:  []string{},
		ConflictingWitnessIDs: []string{},
		MissingRoles:          []domain.EvidenceCoalitionRole{},
		DistinctOrigins:       []string{},
		Reason:                "invalid_claim",
	}
	indexed, exists := rules[claim.RuleID]
	if !validCoalitionID(claim.RuleID) || !exists || !indexed.valid || !validEvidenceStatus(claim.BaselineStatus) {
		return assessment, 0
	}

	witnesses, valid := normalizeCoalitionWitnesses(claim.Witnesses)
	if !valid {
		return assessment, len(witnesses)
	}
	required := make(map[domain.EvidenceCoalitionRole]struct{}, len(indexed.rule.RequiredRoles))
	for _, role := range indexed.rule.RequiredRoles {
		required[role] = struct{}{}
	}
	supportByRole := make(map[domain.EvidenceCoalitionRole]bool, len(required))
	recordedByRole := make(map[domain.EvidenceCoalitionRole]bool, len(required))
	origins := make(map[string]struct{})
	for _, witness := range witnesses {
		if witness.Polarity == domain.EvidenceCoalitionPolarityConflict {
			assessment.ConflictingWitnessIDs = append(assessment.ConflictingWitnessIDs, witness.ID)
			continue
		}
		assessment.SupportingWitnessIDs = append(assessment.SupportingWitnessIDs, witness.ID)
		if _, relevant := required[witness.Role]; !relevant {
			continue
		}
		supportByRole[witness.Role] = true
		if witness.OriginStatus == domain.EvidenceOriginRecorded {
			recordedByRole[witness.Role] = true
			origins[witness.OriginID] = struct{}{}
		}
	}
	for _, role := range [...]domain.EvidenceCoalitionRole{
		domain.EvidenceCoalitionRoleDocumentation,
		domain.EvidenceCoalitionRoleExecution,
		domain.EvidenceCoalitionRoleImplementation,
		domain.EvidenceCoalitionRoleIndependentTest,
	} {
		if _, needed := required[role]; needed && !supportByRole[role] {
			assessment.MissingRoles = append(assessment.MissingRoles, role)
		}
	}
	for origin := range origins {
		assessment.DistinctOrigins = append(assessment.DistinctOrigins, origin)
	}
	sort.Strings(assessment.DistinctOrigins)

	assessment.Independence = coalitionIndependence(indexed.rule, assessment.MissingRoles, recordedByRole, len(origins))
	switch {
	case len(assessment.ConflictingWitnessIDs) != 0:
		assessment.Outcome = domain.EvidenceCoalitionConflict
		assessment.Reason = "conflicting_evidence"
	case len(assessment.MissingRoles) != 0:
		assessment.Outcome = domain.EvidenceCoalitionMissingRoles
		assessment.Reason = "missing_required_roles"
	case assessment.Independence == domain.EvidenceOriginUnknown:
		assessment.Outcome = domain.EvidenceCoalitionUnknownIndependence
		assessment.Reason = "origin_unrecorded"
	case assessment.Independence == domain.EvidenceOriginSame:
		assessment.Outcome = domain.EvidenceCoalitionInsufficientIndependence
		assessment.Reason = "distinct_origins_below_minimum"
	default:
		assessment.Outcome = domain.EvidenceCoalitionSatisfied
		assessment.Reason = "coalition_satisfied"
	}
	return assessment, len(witnesses)
}

func normalizeCoalitionWitnesses(input []domain.EvidenceCoalitionWitness) ([]domain.EvidenceCoalitionWitness, bool) {
	bySource := make(map[coalitionSourceKey]domain.EvidenceCoalitionWitness, len(input))
	ids := make(map[string]struct{}, len(input))
	valid := true
	for _, witness := range input {
		if !validCoalitionWitness(witness) {
			valid = false
		}
		if _, exists := ids[witness.ID]; exists {
			valid = false
		}
		ids[witness.ID] = struct{}{}
		key := coalitionWitnessSourceKey(witness)
		if previous, exists := bySource[key]; exists {
			if !sameCoalitionWitnessMeaning(previous, witness) {
				valid = false
			}
			if witness.ID < previous.ID {
				bySource[key] = witness
			}
			continue
		}
		bySource[key] = witness
	}
	witnesses := make([]domain.EvidenceCoalitionWitness, 0, len(bySource))
	for _, witness := range bySource {
		witnesses = append(witnesses, witness)
	}
	sort.Slice(witnesses, func(i, j int) bool { return witnesses[i].ID < witnesses[j].ID })
	return witnesses, valid
}

func coalitionWitnessSourceKey(witness domain.EvidenceCoalitionWitness) coalitionSourceKey {
	return coalitionSourceKey{
		kind:        witness.Source.Kind,
		owner:       witness.Source.Ref.Owner,
		repo:        witness.Source.Ref.Repo,
		commit:      witness.Source.Ref.Commit,
		path:        witness.Source.Path,
		symbol:      witness.Source.Symbol,
		line:        witness.Source.Line,
		ruleID:      witness.Source.RuleID,
		contentHash: witness.Source.ContentHash,
	}
}

func sameCoalitionWitnessMeaning(left, right domain.EvidenceCoalitionWitness) bool {
	return left.Role == right.Role &&
		left.Polarity == right.Polarity &&
		left.OriginStatus == right.OriginStatus &&
		left.OriginID == right.OriginID
}

func validCoalitionWitness(witness domain.EvidenceCoalitionWitness) bool {
	if !validCoalitionID(witness.ID) || !validCoalitionRole(witness.Role) {
		return false
	}
	if witness.Polarity != domain.EvidenceCoalitionPolaritySupport && witness.Polarity != domain.EvidenceCoalitionPolarityConflict {
		return false
	}
	switch witness.OriginStatus {
	case domain.EvidenceOriginRecorded:
		return witness.OriginID != ""
	case domain.EvidenceOriginUnknown:
		return witness.OriginID == ""
	default:
		return false
	}
}

func coalitionIndependence(rule domain.EvidenceCoalitionRule, missing []domain.EvidenceCoalitionRole, recorded map[domain.EvidenceCoalitionRole]bool, origins int) domain.EvidenceOriginStatus {
	if len(missing) != 0 {
		return domain.EvidenceOriginUnknown
	}
	for _, role := range rule.RequiredRoles {
		if !recorded[role] {
			return domain.EvidenceOriginUnknown
		}
	}
	if origins < rule.MinimumDistinctOrigins {
		return domain.EvidenceOriginSame
	}
	if rule.MinimumDistinctOrigins == 1 {
		return domain.EvidenceOriginRecorded
	}
	return domain.EvidenceOriginDistinct
}

func validCoalitionRole(role domain.EvidenceCoalitionRole) bool {
	switch role {
	case domain.EvidenceCoalitionRoleImplementation,
		domain.EvidenceCoalitionRoleIndependentTest,
		domain.EvidenceCoalitionRoleExecution,
		domain.EvidenceCoalitionRoleDocumentation:
		return true
	default:
		return false
	}
}

func validEvidenceStatus(status domain.EvidenceStatus) bool {
	switch status {
	case domain.EvidenceProven, domain.EvidencePartial, domain.EvidenceAbsent, domain.EvidenceInconclusive, domain.EvidenceConflict:
		return true
	default:
		return false
	}
}

func validCoalitionID(id string) bool {
	if id == "" || strings.TrimSpace(id) != id {
		return false
	}
	for index := 0; index < len(id); index++ {
		if id[index] >= 0x80 {
			return false
		}
	}
	return true
}

func countCoalitionOutcome(metrics *domain.EvidenceCoalitionMetrics, outcome domain.EvidenceCoalitionOutcome) {
	switch outcome {
	case domain.EvidenceCoalitionSatisfied:
		metrics.Satisfied++
	case domain.EvidenceCoalitionConflict:
		metrics.Conflicts++
	case domain.EvidenceCoalitionMissingRoles:
		metrics.MissingRoles++
	case domain.EvidenceCoalitionUnknownIndependence:
		metrics.UnknownIndependence++
	case domain.EvidenceCoalitionInsufficientIndependence:
		metrics.InsufficientIndependence++
	case domain.EvidenceCoalitionInvalidInput:
		metrics.Invalid++
	}
}
