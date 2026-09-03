/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/mechanism_extraction.md
 * @layer L1
 * @updated 2026-09-02
 */

package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
	"golang.org/x/text/unicode/norm"
)

const activeMechanismRulesSHA256 = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"

var immutableCommitPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var contentHashPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

type DeterministicMechanismExtractor struct{}

func NewDeterministicMechanismExtractor() domain.MechanismExtractor {
	return DeterministicMechanismExtractor{}
}

type mechanismCatalog struct {
	SchemaVersion         string                    `json:"schema_version"`
	CatalogID             string                    `json:"catalog_id"`
	CatalogVersion        int                       `json:"catalog_version"`
	ClosedWorld           bool                      `json:"closed_world"`
	UnknownClassification string                    `json:"unknown_classification"`
	Canonicalization      mechanismCanonicalization `json:"canonicalization"`
	SourceContract        mechanismSourceContract   `json:"source_contract"`
	Rules                 []mechanismRule           `json:"rules"`
}

type mechanismCanonicalization struct {
	Text            string `json:"text"`
	Sets            string `json:"sets"`
	WitnessIdentity string `json:"witness_identity"`
	Hash            string `json:"hash"`
}

type mechanismSourceContract struct {
	EvaluatedContainer   string   `json:"evaluated_container"`
	EvaluatedSource      string   `json:"evaluated_source"`
	RuleIdentityField    string   `json:"rule_identity_field"`
	BundleItem           string   `json:"bundle_item"`
	ExactJoinFields      []string `json:"exact_join_fields"`
	ContentUse           string   `json:"content_use"`
	InferRuleFromContent bool     `json:"infer_rule_from_content"`
}

type mechanismRule struct {
	ID                    string                    `json:"id"`
	Version               int                       `json:"version"`
	MechanismClass        string                    `json:"mechanism_class"`
	Name                  string                    `json:"name"`
	Purpose               string                    `json:"purpose"`
	Strategy              string                    `json:"strategy"`
	AcceptedEvidenceRules []mechanismAcceptedRule   `json:"accepted_evidence_rules"`
	MinimumWitnesses      mechanismMinimumWitnesses `json:"minimum_witnesses"`
	FillableFields        []string                  `json:"fillable_fields"`
	LanguageLimits        []string                  `json:"language_limits"`
	EcosystemLimits       []string                  `json:"ecosystem_limits"`
	StatusMax             string                    `json:"status_max"`
}

type mechanismAcceptedRule struct {
	RuleID      string `json:"rule_id"`
	Role        string `json:"role"`
	WitnessKind string `json:"witness_kind"`
	StatusMax   string `json:"status_max,omitempty"`
}

type mechanismMinimumWitnesses struct {
	Partial               int `json:"partial"`
	ProvenImplementation  int `json:"proven_implementation"`
	ProvenIndependent     int `json:"proven_independent"`
	DistinctContentHashes int `json:"distinct_content_hashes"`
}

type evaluatedMechanismSource struct {
	source domain.EvidenceSource
	status domain.EvidenceStatus
	entry  mechanismAcceptedRule
	rule   mechanismRule
}

type mechanismGroup struct {
	rule     mechanismRule
	source   domain.ImmutableCandidateRef
	positive []domain.MechanismWitness
	counter  []domain.MechanismWitness
	roles    map[string]int
	ceiling  domain.MechanismEvidenceStatus
	conflict bool
}

func (DeterministicMechanismExtractor) ExtractMechanisms(ctx context.Context, input domain.MechanismExtractionInput) (domain.MechanismExtractionResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.MechanismExtractionResult{}, err
	}
	catalog, err := decodeMechanismCatalog(input.RuleCatalogJSON, input.RuleCatalogSHA256)
	if err != nil {
		return domain.MechanismExtractionResult{}, domain.ErrIntegrity
	}
	result := domain.MechanismExtractionResult{Mechanisms: []domain.ArchitecturalMechanism{}, Gaps: []domain.MechanismGap{}}
	repository := immutableRef(input.Bundle.Repository)
	if repository.Repository == "" || !immutableCommitPattern.MatchString(repository.Commit) {
		result.Gaps = append(result.Gaps, mechanismGap("IMMUTABLE_IDENTITY_REQUIRED", repository, "", "", domain.MechanismInconclusive))
		return result, nil
	}

	ruleByEvidence := make(map[string]struct {
		rule  mechanismRule
		entry mechanismAcceptedRule
	})
	for _, rule := range catalog.Rules {
		for _, entry := range rule.AcceptedEvidenceRules {
			ruleByEvidence[entry.RuleID] = struct {
				rule  mechanismRule
				entry mechanismAcceptedRule
			}{rule, entry}
		}
	}

	groups := map[string]*mechanismGroup{}
	for _, assessment := range input.Assessments {
		assessmentRef := immutableRef(assessment.Repository)
		if assessmentRef != repository {
			result.Gaps = append(result.Gaps, mechanismGap("EVIDENCE_IDENTITY_MISMATCH", assessmentRef, "", "", domain.MechanismInconclusive))
			continue
		}
		containers := []domain.CapabilityEvidence{assessment.ProblemAlignment}
		containers = append(containers, assessment.Capabilities...)
		for _, container := range containers {
			for _, source := range container.Sources {
				if !validMechanismPath(source.Path) || source.Line < 0 || source.Symbol != normalizeMechanismText(source.Symbol) {
					result.Gaps = append(result.Gaps, mechanismGap("INVALID_EVIDENCE_SOURCE", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				if immutableRef(source.Ref) != repository {
					result.Gaps = append(result.Gaps, mechanismGap("EVIDENCE_IDENTITY_MISMATCH", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				bindingCount, hashOK := mechanismBundleBinding(input.Bundle, source)
				if !hashOK {
					result.Gaps = append(result.Gaps, mechanismGap("EVIDENCE_CONTENT_HASH_MISMATCH", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				if bindingCount == 0 {
					result.Gaps = append(result.Gaps, mechanismGap("UNBOUND_EVALUATED_SOURCE", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				if bindingCount > 1 {
					result.Gaps = append(result.Gaps, mechanismGap("AMBIGUOUS_BUNDLE_BINDING", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				mapped, ok := ruleByEvidence[source.RuleID]
				if !ok {
					result.Gaps = append(result.Gaps, mechanismGap("RULE_NOT_IN_MECHANISM_CATALOG", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				if !compatibleMechanismWitness(mapped.entry, source.Kind) {
					result.Gaps = append(result.Gaps, mechanismGap("INCOMPATIBLE_WITNESS_KIND", repository, source.RuleID, source.ContentHash, domain.MechanismInconclusive))
					continue
				}
				key := repository.Repository + "@" + repository.Commit + "#" + mapped.rule.ID
				group := groups[key]
				if group == nil {
					group = &mechanismGroup{rule: mapped.rule, source: repository, positive: []domain.MechanismWitness{}, counter: []domain.MechanismWitness{}, roles: map[string]int{}, ceiling: domain.MechanismProven}
					groups[key] = group
				}
				witness := domain.MechanismWitness{Repository: repository.Repository, Commit: repository.Commit, Path: source.Path, Symbol: source.Symbol, Line: source.Line, SourceKind: source.Kind, Kind: domain.MechanismWitnessKind(mapped.entry.WitnessKind), AssessmentStatus: container.Status, RuleID: source.RuleID, ContentHash: source.ContentHash}
				applyMechanismWitness(group, witness, mapped.entry, container.Status, &result)
			}
		}
	}

	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		group := groups[key]
		mechanism, ok := projectMechanism(group, input.ConventionalBaseline)
		if !ok {
			result.Gaps = append(result.Gaps, mechanismGap("INSUFFICIENT_MECHANISM_EVIDENCE", group.source, group.rule.ID, firstWitnessHash(group), domain.MechanismInconclusive))
			continue
		}
		result.Mechanisms = append(result.Mechanisms, mechanism)
	}
	sort.Slice(result.Mechanisms, func(i, j int) bool { return result.Mechanisms[i].ID < result.Mechanisms[j].ID })
	sortMechanismGaps(result.Gaps)
	return result, nil
}

func decodeMechanismCatalog(raw []byte, claimed string) (mechanismCatalog, error) {
	digest := sha256.Sum256(raw)
	actual := hex.EncodeToString(digest[:])
	if claimed != activeMechanismRulesSHA256 || actual != activeMechanismRulesSHA256 {
		return mechanismCatalog{}, errors.New("catalog hash mismatch")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var catalog mechanismCatalog
	if decoder.Decode(&catalog) != nil || len(bytes.TrimSpace(raw[decoder.InputOffset():])) != 0 || catalog.SchemaVersion != "tekt.mechanism-rules.v3" || catalog.CatalogVersion != 3 || !catalog.ClosedWorld || catalog.UnknownClassification != "MECHANISM_CLASS_UNKNOWN" || catalog.SourceContract.InferRuleFromContent || len(catalog.Rules) != 4 {
		return mechanismCatalog{}, errors.New("invalid catalog")
	}
	classes := map[string]bool{}
	ruleIDs := map[string]bool{}
	allowedClasses := map[string]bool{"APPEND_ONLY_LOG_RECOVERY": true, "BOUNDED_QUEUE_BACKPRESSURE": true, "CONTENT_ADDRESSED_STORAGE": true, "PORT_ADAPTER_ISOLATION": true}
	for _, rule := range catalog.Rules {
		if !allowedClasses[rule.MechanismClass] || classes[rule.MechanismClass] || rule.ID == "" || rule.Version != 1 || len(rule.AcceptedEvidenceRules) != 8 || rule.MinimumWitnesses != (mechanismMinimumWitnesses{Partial: 1, ProvenImplementation: 1, ProvenIndependent: 1, DistinctContentHashes: 2}) || rule.StatusMax != "PROVEN" {
			return mechanismCatalog{}, errors.New("invalid mechanism rule")
		}
		classes[rule.MechanismClass] = true
		for _, entry := range rule.AcceptedEvidenceRules {
			if entry.RuleID == "" || ruleIDs[entry.RuleID] {
				return mechanismCatalog{}, errors.New("duplicate evidence rule")
			}
			ruleIDs[entry.RuleID] = true
			switch entry.Role {
			case "IMPLEMENTATION", "INDEPENDENT", "CONFLICT":
				if entry.StatusMax != "" {
					return mechanismCatalog{}, errors.New("unexpected status max")
				}
			case "DOCUMENTARY":
				if entry.WitnessKind != string(domain.WitnessDocumentation) || entry.StatusMax != string(domain.MechanismPartial) {
					return mechanismCatalog{}, errors.New("invalid documentary rule")
				}
			default:
				return mechanismCatalog{}, errors.New("invalid role")
			}
		}
	}
	if len(ruleIDs) != 32 {
		return mechanismCatalog{}, errors.New("invalid catalog rule count")
	}
	return catalog, nil
}

func immutableRef(ref domain.CommitRef) domain.ImmutableCandidateRef {
	if ref.Owner == "" || ref.Repo == "" || strings.ContainsAny(ref.Owner+ref.Repo, "@/\\") {
		return domain.ImmutableCandidateRef{}
	}
	return domain.ImmutableCandidateRef{Repository: ref.Owner + "/" + ref.Repo, Commit: ref.Commit}
}

func validMechanismPath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "\\") {
		return false
	}
	for _, segment := range strings.Split(path, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return normalizeMechanismText(path) == path
}

func mechanismBundleBinding(bundle domain.EvidenceBundle, source domain.EvidenceSource) (int, bool) {
	count := 0
	hashOK := true
	for _, item := range bundle.Items {
		digest := sha256.Sum256([]byte(item.Content))
		computed := "sha256:" + hex.EncodeToString(digest[:])
		if item.ContentHash != computed {
			if item.Kind == source.Kind && item.Ref == source.Ref && item.Path == source.Path && item.ContentHash == source.ContentHash {
				hashOK = false
			}
			continue
		}
		if item.Kind == source.Kind && item.Ref == source.Ref && item.Path == source.Path && item.ContentHash == source.ContentHash {
			count++
		}
	}
	return count, hashOK && contentHashPattern.MatchString(source.ContentHash)
}

func compatibleMechanismWitness(entry mechanismAcceptedRule, kind domain.EvidenceKind) bool {
	switch entry.Role {
	case "IMPLEMENTATION":
		return kind == domain.EvidenceFile || kind == domain.EvidenceSymbol
	case "DOCUMENTARY":
		return entry.WitnessKind == string(domain.WitnessDocumentation) && entry.StatusMax == string(domain.MechanismPartial) && kind == domain.EvidenceDocumentation
	case "CONFLICT":
		return true
	}
	switch entry.WitnessKind {
	case string(domain.WitnessTest):
		return kind == domain.EvidenceTest
	case string(domain.WitnessExample):
		return kind == domain.EvidenceExample
	case string(domain.WitnessManifest):
		return kind == domain.EvidenceManifest
	case string(domain.WitnessSchema), string(domain.WitnessConfiguration):
		return kind == domain.EvidenceFile || kind == domain.EvidenceManifest
	}
	return false
}

func applyMechanismWitness(group *mechanismGroup, witness domain.MechanismWitness, entry mechanismAcceptedRule, status domain.EvidenceStatus, result *domain.MechanismExtractionResult) {
	if entry.Role == "CONFLICT" {
		if status == domain.EvidenceConflict {
			group.counter = appendUniqueWitness(group.counter, witness)
			group.conflict = true
		} else {
			result.Gaps = append(result.Gaps, mechanismGap("INCONSISTENT_EVIDENCE_STATUS", group.source, witness.RuleID, witness.ContentHash, domain.MechanismInconclusive))
		}
		return
	}
	if status != domain.EvidenceProven && status != domain.EvidencePartial {
		result.Gaps = append(result.Gaps, mechanismGap("INCONSISTENT_EVIDENCE_STATUS", group.source, witness.RuleID, witness.ContentHash, domain.MechanismInconclusive))
		return
	}
	group.positive = appendUniqueWitness(group.positive, witness)
	group.roles[entry.Role]++
	if status == domain.EvidencePartial || entry.Role == "DOCUMENTARY" {
		group.ceiling = domain.MechanismPartial
	}
}

func appendUniqueWitness(values []domain.MechanismWitness, witness domain.MechanismWitness) []domain.MechanismWitness {
	for _, value := range values {
		if value == witness {
			return values
		}
	}
	return append(values, witness)
}

func projectMechanism(group *mechanismGroup, baseline *domain.ConventionalBaseline) (domain.ArchitecturalMechanism, bool) {
	sortWitnesses(group.positive)
	sortWitnesses(group.counter)
	status := domain.MechanismPartial
	if group.conflict {
		status = domain.MechanismInconclusive
	} else if group.ceiling == domain.MechanismProven && group.roles["IMPLEMENTATION"] >= group.rule.MinimumWitnesses.ProvenImplementation && group.roles["INDEPENDENT"] >= group.rule.MinimumWitnesses.ProvenIndependent && distinctWitnessHashes(group.positive) >= group.rule.MinimumWitnesses.DistinctContentHashes {
		status = domain.MechanismProven
	} else if len(group.positive) < group.rule.MinimumWitnesses.Partial {
		return domain.ArchitecturalMechanism{}, false
	}
	ruleIDs := make([]string, 0, len(group.positive)+len(group.counter))
	hashes := make([]string, 0, len(group.positive)+len(group.counter))
	for _, witness := range append(append([]domain.MechanismWitness{}, group.positive...), group.counter...) {
		ruleIDs = append(ruleIDs, witness.RuleID)
		hashes = append(hashes, witness.ContentHash)
	}
	ruleIDs = canonicalMechanismSet(ruleIDs)
	hashes = canonicalMechanismSet(hashes)
	purpose := normalizeMechanismText(group.rule.Purpose)
	strategy := normalizeMechanismText(group.rule.Strategy)
	mechanism := domain.ArchitecturalMechanism{
		MechanismClass: group.rule.MechanismClass, Name: normalizeMechanismText(group.rule.Name), Purpose: purpose,
		Trigger: purpose, Strategy: strategy, Components: []domain.MechanismComponentRef{}, DataFlow: []string{strategy},
		StateModel: []string{}, FailureModel: []string{}, Boundaries: canonicalMechanismSet(group.rule.EcosystemLimits),
		Preconditions: canonicalMechanismSet(group.rule.LanguageLimits), Tradeoffs: []string{}, Source: group.source,
		Evidence: group.positive, CounterEvidence: group.counter, Status: status, Novelty: mechanismNovelty(group.rule.MechanismClass, status, baseline), ExtractionRuleIDs: ruleIDs,
	}
	mechanism.ID = mechanismID(mechanism, group.rule.ID, hashes)
	return mechanism, true
}

func mechanismID(mechanism domain.ArchitecturalMechanism, mechanismRuleID string, hashes []string) string {
	preimage := struct {
		SchemaVersion  string                       `json:"schema_version"`
		Source         domain.ImmutableCandidateRef `json:"source"`
		RuleID         string                       `json:"rule_id"`
		MechanismClass string                       `json:"mechanism_class"`
		Purpose        string                       `json:"purpose"`
		Strategy       string                       `json:"strategy"`
		WitnessHashes  []string                     `json:"witness_hashes"`
	}{"1", mechanism.Source, mechanismRuleID, mechanism.MechanismClass, mechanism.Purpose, mechanism.Strategy, hashes}
	encoded := canonicalJSON(preimage)
	digest := sha256.Sum256(encoded)
	return "mech_" + hex.EncodeToString(digest[:])
}

func canonicalJSON(value any) []byte {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
	return bytes.TrimSuffix(buffer.Bytes(), []byte{'\n'})
}

func normalizeMechanismText(value string) string {
	if !utf8.ValidString(value) {
		return ""
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f || r >= 0x80 && r <= 0x9f {
			return ""
		}
	}
	value = norm.NFC.String(value)
	return strings.Join(strings.FieldsFunc(value, unicode.IsSpace), " ")
}

func canonicalMechanismSet(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		value = normalizeMechanismText(value)
		if value != "" {
			set[value] = true
		}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func mechanismNovelty(class string, status domain.MechanismEvidenceStatus, baseline *domain.ConventionalBaseline) domain.MechanismNovelty {
	if baseline == nil || baseline.BaselineID == "" || baseline.IdeaHash == "" {
		return domain.NoveltyUnknown
	}
	for _, value := range baseline.MechanismClasses {
		if value == class {
			return domain.NoveltyConventional
		}
	}
	if status == domain.MechanismProven {
		return domain.NoveltyDivergent
	}
	return domain.NoveltyUnknown
}

func sortWitnesses(values []domain.MechanismWitness) {
	sort.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j]
		return witnessSortKey(a) < witnessSortKey(b)
	})
}

func witnessSortKey(value domain.MechanismWitness) string {
	return strings.Join([]string{value.Repository, value.Commit, value.Path, value.Symbol, strconv.Itoa(value.Line), string(value.SourceKind), string(value.Kind), string(value.AssessmentStatus), value.RuleID, value.ContentHash}, "\x00")
}

func distinctWitnessHashes(values []domain.MechanismWitness) int {
	set := map[string]bool{}
	for _, value := range values {
		set[value.ContentHash] = true
	}
	return len(set)
}

func firstWitnessHash(group *mechanismGroup) string {
	if len(group.positive) > 0 {
		return group.positive[0].ContentHash
	}
	if len(group.counter) > 0 {
		return group.counter[0].ContentHash
	}
	return ""
}

func mechanismGap(code string, source domain.ImmutableCandidateRef, ruleID, witnessHash string, ceiling domain.MechanismEvidenceStatus) domain.MechanismGap {
	return domain.MechanismGap{Code: code, Source: source, RuleID: ruleID, WitnessHash: witnessHash, StatusCeiling: ceiling}
}

func sortMechanismGaps(values []domain.MechanismGap) {
	sort.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j]
		ka := strings.Join([]string{a.Code, a.Source.Repository, a.Source.Commit, a.RuleID, a.WitnessHash}, "\x00")
		kb := strings.Join([]string{b.Code, b.Source.Repository, b.Source.Commit, b.RuleID, b.WitnessHash}, "\x00")
		return ka < kb
	})
}
