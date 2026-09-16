/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_history.md
 * @layer L1
 * @updated 2026-09-12
 */

package usecase

import (
	"sort"
	"strings"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const (
	algorithmHistoryMaxAlgorithms = 256
	algorithmHistoryMaxProvenance = 1024
	algorithmHistoryMaxClaims     = 512
	algorithmHistoryMaxBases      = 2048
	algorithmHistoryMaxReferences = 64
	algorithmHistoryMaxIDBytes    = 128
	algorithmHistoryMaxTextBytes  = 4096
	algorithmHistoryMaxTotalText  = 1024 * 1024
	algorithmHistoryMaxOperations = 16384
)

type canonicalHistoryProvenance struct {
	id, kind, origin, citation, locator string
}

type canonicalHistoryAlgorithm struct {
	id, name, chronology, reason string
	start, end                   *int
	provenance                   []string
}

type canonicalHistoryClaim struct {
	id, subject, relation, object, statement string
}

type canonicalHistoryBasis struct {
	id, claim, kind, stance, provenance, summary string
	derived                                      []string
}

// AnalyzeAlgorithmHistory classifies only explicitly supplied historical claims.
func AnalyzeAlgorithmHistory(input domain.AlgorithmHistoryInput) domain.AlgorithmHistoryReport {
	report := newAlgorithmHistoryReport(input)
	if input.SchemaVersion != domain.AlgorithmHistoryInputSchemaVersion {
		return unknownAlgorithmHistoryReport(report, "invalid_schema", false)
	}
	if algorithmHistoryStructuralBudgetExceeded(input) {
		return unknownAlgorithmHistoryReport(report, "budget_exceeded", true)
	}
	if len(input.Algorithms) == 0 {
		return unknownAlgorithmHistoryReport(report, "no_algorithms", false)
	}
	if len(input.Provenance) == 0 {
		return unknownAlgorithmHistoryReport(report, "no_provenance", false)
	}
	if len(input.Claims) == 0 {
		return unknownAlgorithmHistoryReport(report, "no_claims", false)
	}
	if len(input.Bases) == 0 {
		return unknownAlgorithmHistoryReport(report, "no_bases", false)
	}

	provenance, provenanceAliases, duplicates, reason := canonicalizeHistoryProvenance(input.Provenance)
	if reason != "" {
		return unknownAlgorithmHistoryReport(report, reason, false)
	}
	report.Metrics.DuplicateProvenance = duplicates

	algorithms, duplicates, reason := canonicalizeHistoryAlgorithms(input.Algorithms, provenanceAliases)
	if reason != "" {
		return unknownAlgorithmHistoryReport(report, reason, false)
	}
	report.Metrics.DuplicateAlgorithms = duplicates

	claims, claimAliases, duplicates, reason := canonicalizeHistoryClaims(input.Claims, algorithms)
	if reason != "" {
		return unknownAlgorithmHistoryReport(report, reason, false)
	}
	report.Metrics.DuplicateClaims = duplicates

	bases, duplicates, reason := canonicalizeHistoryBases(input.Bases, claims, claimAliases, provenance, provenanceAliases)
	if reason != "" {
		return unknownAlgorithmHistoryReport(report, reason, false)
	}
	report.Metrics.DuplicateBases = duplicates

	report.Metrics.Algorithms = len(algorithms)
	report.Metrics.Provenance = len(provenance)
	report.Metrics.Claims = len(claims)
	report.Metrics.Bases = len(bases)
	report.Metrics.Operations = algorithmHistoryOperationCount(input, len(claims), len(bases))
	if report.Metrics.Operations > algorithmHistoryMaxOperations {
		return unknownAlgorithmHistoryReport(report, "budget_exceeded", true)
	}

	byClaim := make(map[string][]canonicalHistoryBasis, len(claims))
	for _, basis := range bases {
		byClaim[basis.claim] = append(byClaim[basis.claim], basis)
		if basis.kind == string(domain.AlgorithmHistoryEvidence) {
			report.Metrics.Evidence++
		} else {
			report.Metrics.Inferences++
		}
	}
	claimIDs := sortedHistoryKeys(claims)
	for _, claimID := range claimIDs {
		assessment, ok := assessAlgorithmHistoryClaim(claims[claimID], byClaim[claimID], provenance)
		if !ok {
			return unknownAlgorithmHistoryReport(report, "invalid_input", false)
		}
		report.Assessments = append(report.Assessments, assessment)
		countAlgorithmHistoryOutcome(&report.Metrics, assessment.Outcome)
	}
	report.Status = domain.AlgorithmHistoryReportComplete
	report.Reason = "analyzed"
	return report
}

func newAlgorithmHistoryReport(input domain.AlgorithmHistoryInput) domain.AlgorithmHistoryReport {
	return domain.AlgorithmHistoryReport{
		SchemaVersion: domain.AlgorithmHistoryReportSchemaVersion,
		Status:        domain.AlgorithmHistoryReportUnknown,
		Assessments:   make([]domain.AlgorithmHistoryAssessment, 0),
		Metrics: domain.AlgorithmHistoryMetrics{
			InputAlgorithms: len(input.Algorithms), InputProvenance: len(input.Provenance),
			InputClaims: len(input.Claims), InputBases: len(input.Bases),
		},
	}
}

func unknownAlgorithmHistoryReport(report domain.AlgorithmHistoryReport, reason string, limited bool) domain.AlgorithmHistoryReport {
	report.Status, report.Reason = domain.AlgorithmHistoryReportUnknown, reason
	report.Assessments = make([]domain.AlgorithmHistoryAssessment, 0)
	report.Metrics.LimitExceeded = limited
	return report
}

func algorithmHistoryStructuralBudgetExceeded(input domain.AlgorithmHistoryInput) bool {
	if len(input.Algorithms) > algorithmHistoryMaxAlgorithms || len(input.Provenance) > algorithmHistoryMaxProvenance ||
		len(input.Claims) > algorithmHistoryMaxClaims || len(input.Bases) > algorithmHistoryMaxBases {
		return true
	}
	total := 0
	addText := func(value string, max int) bool {
		if !utf8.ValidString(value) || len(value) > max {
			return false
		}
		total += len(value)
		return total <= algorithmHistoryMaxTotalText
	}
	for _, item := range input.Algorithms {
		if len(item.ProvenanceIDs) > algorithmHistoryMaxReferences || !validHistoryIDBudget(item.ID) ||
			!addText(item.Name, algorithmHistoryMaxTextBytes) || !addText(item.ChronologyReason, algorithmHistoryMaxTextBytes) {
			return true
		}
		for _, id := range item.ProvenanceIDs {
			if !validHistoryIDBudget(id) {
				return true
			}
		}
	}
	for _, item := range input.Provenance {
		if !validHistoryIDBudget(item.ID) || !addText(item.OriginID, algorithmHistoryMaxIDBytes) ||
			!addText(item.Citation, algorithmHistoryMaxTextBytes) || !addText(item.Locator, algorithmHistoryMaxTextBytes) {
			return true
		}
	}
	for _, item := range input.Claims {
		if !validHistoryIDBudget(item.ID) || !validHistoryIDBudget(item.SubjectID) || !validHistoryIDBudget(item.ObjectID) ||
			!addText(item.Statement, algorithmHistoryMaxTextBytes) {
			return true
		}
	}
	for _, item := range input.Bases {
		if len(item.DerivedFromIDs) > algorithmHistoryMaxReferences || !validHistoryIDBudget(item.ID) ||
			!validHistoryIDBudget(item.ClaimID) || !validHistoryIDBudget(item.ProvenanceID) ||
			!addText(item.Summary, algorithmHistoryMaxTextBytes) {
			return true
		}
		for _, id := range item.DerivedFromIDs {
			if !validHistoryIDBudget(id) {
				return true
			}
		}
	}
	operations := len(input.Algorithms) + len(input.Provenance) + len(input.Claims) + len(input.Bases) + len(input.Claims)*2 + len(input.Bases)*3
	for _, item := range input.Algorithms {
		operations += len(item.ProvenanceIDs)
	}
	for _, item := range input.Bases {
		operations += len(item.DerivedFromIDs)
	}
	return operations > algorithmHistoryMaxOperations
}

func validHistoryIDBudget(value string) bool { return len(value) <= algorithmHistoryMaxIDBytes }

func canonicalHistoryID(value string) (string, bool) {
	value = strings.Trim(value, " \t\n\r\v\f")
	if value == "" || len(value) > algorithmHistoryMaxIDBytes {
		return "", false
	}
	bytes := []byte(value)
	for i, c := range bytes {
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || strings.ContainsRune("._:/-", rune(c))) {
			return "", false
		}
		bytes[i] = c
	}
	return string(bytes), true
}

func validOpaqueHistoryText(value string) bool {
	return utf8.ValidString(value) && value != "" && value == strings.Trim(value, " \t\n\r\v\f")
}

func canonicalHistorySet(values []string) ([]string, bool) {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		id, ok := canonicalHistoryID(value)
		if !ok {
			return nil, false
		}
		seen[id] = struct{}{}
	}
	result := sortedHistoryKeys(seen)
	return result, true
}

func canonicalizeHistoryProvenance(input []domain.AlgorithmHistoryProvenance) (map[string]canonicalHistoryProvenance, map[string]string, int, string) {
	byID := make(map[string]canonicalHistoryProvenance, len(input))
	for _, value := range input {
		id, ok := canonicalHistoryID(value.ID)
		if !ok || !validHistoryProvenanceKind(value.Kind) || !validOpaqueHistoryText(value.OriginID) || !validOpaqueHistoryText(value.Citation) || !validOpaqueHistoryText(value.Locator) {
			return nil, nil, 0, "invalid_input"
		}
		item := canonicalHistoryProvenance{id: id, kind: string(value.Kind), origin: value.OriginID, citation: value.Citation, locator: value.Locator}
		if prior, exists := byID[id]; exists && prior != item {
			return nil, nil, 0, "ambiguous_input"
		} else {
			byID[id] = item
		}
	}
	aliases := make(map[string]string, len(byID))
	representative := make(map[string]string)
	for _, id := range sortedHistoryKeys(byID) {
		item := byID[id]
		key := item.kind + "\x00" + item.origin + "\x00" + item.citation + "\x00" + item.locator
		if prior, exists := representative[key]; exists {
			aliases[id] = prior
		} else {
			representative[key], aliases[id] = id, id
		}
	}
	result := make(map[string]canonicalHistoryProvenance, len(representative))
	for id, item := range byID {
		if aliases[id] == id {
			result[id] = item
		}
	}
	return result, aliases, len(input) - len(result), ""
}

func canonicalizeHistoryAlgorithms(input []domain.AlgorithmHistoryAlgorithm, provenanceAliases map[string]string) (map[string]canonicalHistoryAlgorithm, int, string) {
	result := make(map[string]canonicalHistoryAlgorithm, len(input))
	for _, value := range input {
		id, ok := canonicalHistoryID(value.ID)
		if !ok || !validOpaqueHistoryText(value.Name) || !validOpaqueHistoryText(value.ChronologyReason) || !validHistoryChronology(value) {
			return nil, 0, "invalid_input"
		}
		refs, ok := canonicalHistorySet(value.ProvenanceIDs)
		if !ok || len(refs) == 0 {
			return nil, 0, "invalid_input"
		}
		for i, ref := range refs {
			mapped, exists := provenanceAliases[ref]
			if !exists {
				return nil, 0, "invalid_input"
			}
			refs[i] = mapped
		}
		refs = uniqueSortedHistoryStrings(refs)
		item := canonicalHistoryAlgorithm{id: id, name: value.Name, chronology: string(value.Chronology), reason: value.ChronologyReason, start: copyHistoryInt(value.YearStart), end: copyHistoryInt(value.YearEnd), provenance: refs}
		if prior, exists := result[id]; exists && !sameHistoryAlgorithm(prior, item) {
			return nil, 0, "ambiguous_input"
		} else {
			result[id] = item
		}
	}
	return result, len(input) - len(result), ""
}

func canonicalizeHistoryClaims(input []domain.AlgorithmHistoryClaim, algorithms map[string]canonicalHistoryAlgorithm) (map[string]canonicalHistoryClaim, map[string]string, int, string) {
	byID := make(map[string]canonicalHistoryClaim, len(input))
	for _, value := range input {
		id, ok := canonicalHistoryID(value.ID)
		subject, okS := canonicalHistoryID(value.SubjectID)
		object, okO := canonicalHistoryID(value.ObjectID)
		if !ok || !okS || !okO || subject == object || !validHistoryRelation(value.Relation) || !validOpaqueHistoryText(value.Statement) {
			return nil, nil, 0, "invalid_input"
		}
		if _, exists := algorithms[subject]; !exists {
			return nil, nil, 0, "invalid_input"
		}
		if _, exists := algorithms[object]; !exists {
			return nil, nil, 0, "invalid_input"
		}
		item := canonicalHistoryClaim{id: id, subject: subject, relation: string(value.Relation), object: object, statement: value.Statement}
		if prior, exists := byID[id]; exists && prior != item {
			return nil, nil, 0, "ambiguous_input"
		} else {
			byID[id] = item
		}
	}
	aliases := make(map[string]string, len(byID))
	representatives := make(map[string]string)
	for _, id := range sortedHistoryKeys(byID) {
		item := byID[id]
		key := item.subject + "\x00" + item.relation + "\x00" + item.object + "\x00" + item.statement
		if prior, exists := representatives[key]; exists {
			aliases[id] = prior
		} else {
			representatives[key] = id
			aliases[id] = id
		}
	}
	result := make(map[string]canonicalHistoryClaim, len(representatives))
	for id, item := range byID {
		if aliases[id] == id {
			result[id] = item
		}
	}
	return result, aliases, len(input) - len(result), ""
}

func canonicalizeHistoryBases(input []domain.AlgorithmHistoryBasis, claims map[string]canonicalHistoryClaim, claimAliases map[string]string, provenance map[string]canonicalHistoryProvenance, provenanceAliases map[string]string) ([]canonicalHistoryBasis, int, string) {
	byID := make(map[string]canonicalHistoryBasis, len(input))
	for _, value := range input {
		id, ok := canonicalHistoryID(value.ID)
		claim, okC := canonicalHistoryID(value.ClaimID)
		prov, okP := canonicalHistoryID(value.ProvenanceID)
		derived, okD := canonicalHistorySet(value.DerivedFromIDs)
		if !ok || !okC || !okP || !okD || !validHistoryBasisKind(value.Kind) || !validHistoryStance(value.Stance) || !validOpaqueHistoryText(value.Summary) {
			return nil, 0, "invalid_input"
		}
		claim, mappedC := claimAliases[claim]
		prov, mappedP := provenanceAliases[prov]
		if !mappedC || !mappedP {
			return nil, 0, "invalid_input"
		}
		if value.Kind == domain.AlgorithmHistoryEvidence && len(derived) != 0 || value.Kind == domain.AlgorithmHistoryInference && len(derived) == 0 {
			return nil, 0, "invalid_input"
		}
		item := canonicalHistoryBasis{id: id, claim: claim, kind: string(value.Kind), stance: string(value.Stance), provenance: prov, summary: value.Summary, derived: derived}
		if prior, exists := byID[id]; exists && !sameHistoryBasis(prior, item) {
			return nil, 0, "ambiguous_input"
		} else {
			byID[id] = item
		}
	}
	evidenceAliases := make(map[string]string)
	evidenceRepresentative := make(map[string]string)
	for _, id := range sortedHistoryKeys(byID) {
		item := byID[id]
		if item.kind != string(domain.AlgorithmHistoryEvidence) {
			continue
		}
		key := historyBasisKey(item)
		if prior, exists := evidenceRepresentative[key]; exists {
			evidenceAliases[id] = prior
		} else {
			evidenceRepresentative[key] = id
			evidenceAliases[id] = id
		}
	}
	for id, item := range byID {
		if item.kind != string(domain.AlgorithmHistoryInference) {
			continue
		}
		for i, ref := range item.derived {
			mapped, exists := evidenceAliases[ref]
			if !exists {
				return nil, 0, "invalid_input"
			}
			evidence := byID[mapped]
			if evidence.claim != item.claim {
				return nil, 0, "invalid_input"
			}
			item.derived[i] = mapped
		}
		item.derived = uniqueSortedHistoryStrings(item.derived)
		byID[id] = item
	}
	representatives := make(map[string]string)
	aliases := make(map[string]string, len(byID))
	for _, id := range sortedHistoryKeys(byID) {
		item := byID[id]
		key := historyBasisKey(item)
		if prior, exists := representatives[key]; exists {
			aliases[id] = prior
		} else {
			representatives[key] = id
			aliases[id] = id
		}
	}
	result := make([]canonicalHistoryBasis, 0, len(representatives))
	evidenceByClaim := make(map[string]bool)
	for id, item := range byID {
		if aliases[id] != id {
			continue
		}
		result = append(result, item)
		if item.kind == string(domain.AlgorithmHistoryEvidence) {
			evidenceByClaim[item.claim] = true
		}
	}
	for claim := range claims {
		if !evidenceByClaim[claim] {
			return nil, 0, "invalid_input"
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].id < result[j].id })
	_ = provenance
	return result, len(input) - len(result), ""
}

func assessAlgorithmHistoryClaim(claim canonicalHistoryClaim, bases []canonicalHistoryBasis, provenance map[string]canonicalHistoryProvenance) (domain.AlgorithmHistoryAssessment, bool) {
	a := domain.AlgorithmHistoryAssessment{ClaimID: claim.id, SubjectID: claim.subject, Relation: domain.AlgorithmHistoryRelation(claim.relation), ObjectID: claim.object, Statement: claim.statement, SupportingEvidenceIDs: []string{}, ConflictingEvidenceIDs: []string{}, AmbiguousEvidenceIDs: []string{}, SupportingInferenceIDs: []string{}, ConflictingInferenceIDs: []string{}, AmbiguousInferenceIDs: []string{}, DistinctSupportingOrigins: []string{}}
	origins := map[string]struct{}{}
	evidence := 0
	for _, basis := range bases {
		target := historyAssessmentList(&a, basis.kind, basis.stance)
		*target = append(*target, basis.id)
		if basis.kind == string(domain.AlgorithmHistoryEvidence) {
			evidence++
			if basis.stance == string(domain.AlgorithmHistorySupport) {
				origins[provenance[basis.provenance].origin] = struct{}{}
			}
		}
	}
	if evidence == 0 {
		return a, false
	}
	a.DistinctSupportingOrigins = sortedHistoryKeys(origins)
	support, conflict, ambiguous := len(a.SupportingEvidenceIDs), len(a.ConflictingEvidenceIDs), len(a.AmbiguousEvidenceIDs)
	switch {
	case support > 0 && conflict > 0:
		a.Outcome, a.Reason = domain.AlgorithmHistoryContested, "evidence_contested"
	case ambiguous > 0:
		a.Outcome, a.Reason = domain.AlgorithmHistoryUnknown, "evidence_ambiguous"
	case conflict > 0:
		a.Outcome, a.Reason = domain.AlgorithmHistoryContradicted, "evidence_contradicts_claim"
	case support > 0 && len(origins) >= 2:
		a.Outcome, a.Reason = domain.AlgorithmHistoryCorroborated, "independent_origins_corroborate"
	case support > 0:
		a.Outcome, a.Reason = domain.AlgorithmHistorySupported, "single_origin_support"
	default:
		a.Outcome, a.Reason = domain.AlgorithmHistoryUnknown, "insufficient_knowledge"
	}
	return a, true
}

func historyAssessmentList(a *domain.AlgorithmHistoryAssessment, kind, stance string) *[]string {
	if kind == string(domain.AlgorithmHistoryEvidence) {
		switch stance {
		case string(domain.AlgorithmHistorySupport):
			return &a.SupportingEvidenceIDs
		case string(domain.AlgorithmHistoryConflict):
			return &a.ConflictingEvidenceIDs
		default:
			return &a.AmbiguousEvidenceIDs
		}
	}
	switch stance {
	case string(domain.AlgorithmHistorySupport):
		return &a.SupportingInferenceIDs
	case string(domain.AlgorithmHistoryConflict):
		return &a.ConflictingInferenceIDs
	default:
		return &a.AmbiguousInferenceIDs
	}
}

func countAlgorithmHistoryOutcome(m *domain.AlgorithmHistoryMetrics, outcome domain.AlgorithmHistoryOutcome) {
	switch outcome {
	case domain.AlgorithmHistoryCorroborated:
		m.Corroborated++
	case domain.AlgorithmHistorySupported:
		m.Supported++
	case domain.AlgorithmHistoryContested:
		m.Contested++
	case domain.AlgorithmHistoryContradicted:
		m.Contradicted++
	default:
		m.Unknown++
	}
}
func algorithmHistoryOperationCount(input domain.AlgorithmHistoryInput, claims, bases int) int {
	n := len(input.Algorithms) + len(input.Provenance) + len(input.Claims) + len(input.Bases) + len(input.Claims)*2 + len(input.Bases)*2 + claims + bases
	for _, v := range input.Algorithms {
		n += len(v.ProvenanceIDs)
	}
	for _, v := range input.Bases {
		n += len(v.DerivedFromIDs)
	}
	return n
}
func validHistoryChronology(v domain.AlgorithmHistoryAlgorithm) bool {
	switch v.Chronology {
	case domain.AlgorithmChronologyKnown, domain.AlgorithmChronologyApproximate:
		return v.YearStart != nil && v.YearEnd != nil && *v.YearStart >= -9999 && *v.YearEnd <= 9999 && *v.YearStart <= *v.YearEnd
	case domain.AlgorithmChronologyUnknown:
		return v.YearStart == nil && v.YearEnd == nil
	default:
		return false
	}
}
func validHistoryProvenanceKind(v domain.AlgorithmHistoryProvenanceKind) bool {
	return v == domain.AlgorithmHistoryPrimarySource || v == domain.AlgorithmHistorySecondarySource || v == domain.AlgorithmHistoryEmpiricalArtifact
}
func validHistoryRelation(v domain.AlgorithmHistoryRelation) bool {
	switch v {
	case domain.AlgorithmHistoryDerivedFrom, domain.AlgorithmHistoryInfluencedBy, domain.AlgorithmHistoryImproves, domain.AlgorithmHistoryGeneralizes, domain.AlgorithmHistorySpecializes, domain.AlgorithmHistoryCombines, domain.AlgorithmHistoryRediscovers, domain.AlgorithmHistoryResembles:
		return true
	}
	return false
}
func validHistoryBasisKind(v domain.AlgorithmHistoryBasisKind) bool {
	return v == domain.AlgorithmHistoryEvidence || v == domain.AlgorithmHistoryInference
}
func validHistoryStance(v domain.AlgorithmHistoryStance) bool {
	return v == domain.AlgorithmHistorySupport || v == domain.AlgorithmHistoryConflict || v == domain.AlgorithmHistoryAmbiguous
}
func copyHistoryInt(v *int) *int {
	if v == nil {
		return nil
	}
	n := *v
	return &n
}
func sameHistoryAlgorithm(a, b canonicalHistoryAlgorithm) bool {
	return a.id == b.id && a.name == b.name && a.chronology == b.chronology && a.reason == b.reason && equalHistoryInts(a.start, b.start) && equalHistoryInts(a.end, b.end) && strings.Join(a.provenance, "\x00") == strings.Join(b.provenance, "\x00")
}
func equalHistoryInts(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
func sameHistoryBasis(a, b canonicalHistoryBasis) bool {
	return a.id == b.id && a.claim == b.claim && a.kind == b.kind && a.stance == b.stance && a.provenance == b.provenance && a.summary == b.summary && strings.Join(a.derived, "\x00") == strings.Join(b.derived, "\x00")
}
func historyBasisKey(v canonicalHistoryBasis) string {
	return v.claim + "\x00" + v.kind + "\x00" + v.stance + "\x00" + v.provenance + "\x00" + v.summary + "\x00" + strings.Join(v.derived, "\x00")
}
func uniqueSortedHistoryStrings(v []string) []string {
	sort.Strings(v)
	out := v[:0]
	for _, s := range v {
		if len(out) == 0 || out[len(out)-1] != s {
			out = append(out, s)
		}
	}
	return out
}
func sortedHistoryKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
