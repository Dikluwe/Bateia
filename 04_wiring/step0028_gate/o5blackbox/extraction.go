/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/extraction.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const extractionCatalogSHA256 = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"

type extractionRequest struct {
	Bundle struct {
		RunID string                 `json:"run_id"`
		Items []extractionBundleItem `json:"items"`
	} `json:"bundle"`
	Assessments          []extractionAssessment `json:"assessments"`
	ConventionalBaseline *extractionBaseline    `json:"conventional_baseline"`
	RuleCatalogJSON      string                 `json:"rule_catalog_json"`
	RuleCatalogSHA256    string                 `json:"rule_catalog_sha256"`
}

type extractionBundleItem struct {
	Kind        string `json:"kind"`
	Ref         string `json:"ref"`
	Path        string `json:"path"`
	ContentHash string `json:"content_hash"`
	ContentUTF8 string `json:"content_utf8"`
}

type extractionAssessment struct {
	Repository string             `json:"repository"`
	Status     string             `json:"status"`
	Sources    []extractionSource `json:"sources"`
}

type extractionSource struct {
	Kind        string `json:"kind"`
	Ref         string `json:"ref"`
	Path        string `json:"path"`
	ContentHash string `json:"content_hash"`
	RuleID      string `json:"rule_id"`
	Symbol      string `json:"symbol"`
	Line        int    `json:"line"`
}

type extractionBaseline struct {
	BaselineID       string   `json:"baseline_id"`
	IdeaHash         string   `json:"idea_hash"`
	MechanismClasses []string `json:"mechanism_classes"`
}

type extractionCatalog struct {
	SchemaVersion         string                       `json:"schema_version"`
	CatalogID             string                       `json:"catalog_id"`
	CatalogVersion        int                          `json:"catalog_version"`
	ClosedWorld           bool                         `json:"closed_world"`
	UnknownClassification string                       `json:"unknown_classification"`
	Canonicalization      extractionCatalogCanonical   `json:"canonicalization"`
	SourceContract        extractionCatalogSource      `json:"source_contract"`
	Rules                 []extractionCatalogMechanism `json:"rules"`
}

type extractionCatalogCanonical struct {
	Text            string `json:"text"`
	Sets            string `json:"sets"`
	WitnessIdentity string `json:"witness_identity"`
	Hash            string `json:"hash"`
}

type extractionCatalogSource struct {
	EvaluatedContainer   string   `json:"evaluated_container"`
	EvaluatedSource      string   `json:"evaluated_source"`
	RuleIdentityField    string   `json:"rule_identity_field"`
	BundleItem           string   `json:"bundle_item"`
	ExactJoinFields      []string `json:"exact_join_fields"`
	ContentUse           string   `json:"content_use"`
	InferRuleFromContent bool     `json:"infer_rule_from_content"`
}

type extractionCatalogMechanism struct {
	ID                    string                      `json:"id"`
	Version               int                         `json:"version"`
	MechanismClass        string                      `json:"mechanism_class"`
	Name                  string                      `json:"name"`
	Purpose               string                      `json:"purpose"`
	Strategy              string                      `json:"strategy"`
	AcceptedEvidenceRules []extractionCatalogEvidence `json:"accepted_evidence_rules"`
	MinimumWitnesses      extractionMinimumWitnesses  `json:"minimum_witnesses"`
	FillableFields        []string                    `json:"fillable_fields"`
	LanguageLimits        []string                    `json:"language_limits"`
	EcosystemLimits       []string                    `json:"ecosystem_limits"`
	StatusMax             string                      `json:"status_max"`
}

type extractionCatalogEvidence struct {
	RuleID      string `json:"rule_id"`
	Role        string `json:"role"`
	WitnessKind string `json:"witness_kind"`
	StatusMax   string `json:"status_max,omitempty"`
}

type extractionMinimumWitnesses struct {
	Partial               int `json:"partial"`
	ProvenImplementation  int `json:"proven_implementation"`
	ProvenIndependent     int `json:"proven_independent"`
	DistinctContentHashes int `json:"distinct_content_hashes"`
}

type extractionResult struct {
	Mechanisms []extractionMechanism `json:"mechanisms"`
	Gaps       []extractionGap       `json:"gaps"`
}

type extractionMechanism struct {
	ID                string                `json:"id"`
	MechanismClass    string                `json:"mechanism_class"`
	Name              string                `json:"name"`
	Purpose           string                `json:"purpose"`
	Trigger           string                `json:"trigger"`
	Strategy          string                `json:"strategy"`
	Components        []extractionComponent `json:"components"`
	DataFlow          []string              `json:"data_flow"`
	StateModel        []string              `json:"state_model"`
	FailureModel      []string              `json:"failure_model"`
	Boundaries        []string              `json:"boundaries"`
	Preconditions     []string              `json:"preconditions"`
	Tradeoffs         []string              `json:"tradeoffs"`
	Source            extractionImmutable   `json:"source"`
	Evidence          []extractionWitness   `json:"evidence"`
	CounterEvidence   []extractionWitness   `json:"counter_evidence"`
	Status            string                `json:"status"`
	Novelty           string                `json:"novelty"`
	ExtractionRuleIDs []string              `json:"extraction_rule_ids"`
}

type extractionComponent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type extractionImmutable struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
}

type extractionWitness struct {
	Repository       string `json:"repository"`
	Commit           string `json:"commit"`
	Path             string `json:"path"`
	Symbol           string `json:"symbol"`
	Line             int    `json:"line"`
	SourceKind       string `json:"source_kind"`
	Kind             string `json:"kind"`
	AssessmentStatus string `json:"assessment_status"`
	RuleID           string `json:"rule_id"`
	ContentHash      string `json:"content_hash"`
}

type extractionGap struct {
	Code          string              `json:"code"`
	Source        extractionImmutable `json:"source"`
	RuleID        string              `json:"rule_id"`
	WitnessHash   string              `json:"witness_hash"`
	StatusCeiling string              `json:"status_ceiling"`
}

type extractionEffects struct {
	NetworkCalls            int      `json:"network_calls"`
	FilesystemReads         int      `json:"filesystem_reads"`
	FilesystemWrites        int      `json:"filesystem_writes"`
	ProcessSpawns           int      `json:"process_spawns"`
	ExternalMutations       int      `json:"external_mutations"`
	ResolveCalls            int      `json:"resolve_calls"`
	ReserveCalls            int      `json:"reserve_calls"`
	IDCalls                 int      `json:"id_calls"`
	CreateCalls             int      `json:"create_calls"`
	StoreInstanceIDs        []string `json:"store_instance_ids"`
	QueueInstanceIDs        []string `json:"queue_instance_ids"`
	RootInstanceIDs         []string `json:"root_instance_ids"`
	CursorSecretInstanceIDs []string `json:"cursor_secret_instance_ids"`
	WorkerInstanceIDs       []string `json:"worker_instance_ids"`
}

type extractionPublication struct {
	ArtifactNames    []string          `json:"artifact_names"`
	ArtifactHashes   map[string]string `json:"artifact_hashes"`
	IndexHashes      map[string]string `json:"index_hashes"`
	LeakedFieldNames []string          `json:"leaked_field_names"`
	LayerFindings    []string          `json:"layer_findings"`
}

type extractionControlledError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type extractionRuleBinding struct {
	Mechanism extractionCatalogMechanism
	Evidence  extractionCatalogEvidence
}

type extractionQualified struct {
	GroupKey string
	Rule     extractionCatalogMechanism
	Binding  extractionCatalogEvidence
	Witness  extractionWitness
}

type extractionExpectedGroup struct {
	Rule     extractionCatalogMechanism
	Source   extractionImmutable
	Positive []extractionWitness
	Conflict []extractionWitness
	Status   string
}

type extractionIssue struct {
	Code    string
	Source  extractionImmutable
	RuleID  string
	Hash    string
	Ceiling string
	Loose   bool
}

func evaluateMechanismExtract(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input extractionRequest
	if err := decodeStrict(request, &input); err != nil {
		state.fail("MECHANISM_EXTRACT request is not strict protocol JSON")
		return
	}

	extractionRequirePureObservation(state, output)

	catalogBytes, catalog, catalogOK := extractionDecodeCatalog(input)
	if !catalogOK {
		extractionRequireIntegrityRejection(state, output)
		return
	}
	state.require(input.RuleCatalogSHA256 == extractionCatalogSHA256, "extractor did not receive the frozen catalog hash")
	state.require(extractionDigest(catalogBytes) == extractionCatalogSHA256, "decoded catalog bytes do not match the frozen catalog hash")

	if output.Error != nil {
		state.fail("valid MECHANISM_EXTRACT request returned a controlled error")
		return
	}
	if rawIsNull(output.ReturnValue) {
		state.fail("valid MECHANISM_EXTRACT request returned null")
		return
	}
	var result extractionResult
	if err := decodeReturn(output.ReturnValue, &result); err != nil {
		state.fail("MECHANISM_EXTRACT return_value is not a strict extraction_result")
		return
	}
	state.require(result.Mechanisms != nil, "mechanisms must be [] rather than null")
	state.require(result.Gaps != nil, "gaps must be [] rather than null")

	bindings := extractionBindings(catalog)
	groups, issues := extractionQualify(input, bindings)
	extractionEvaluateGroups(groups)
	extractionCheckMechanisms(state, input, result, groups, issues)
	extractionCheckGaps(state, result.Gaps, issues, groups)
	extractionCheckSafety(state, input, result)
	if len(result.Mechanisms) == 0 && extractionHasConcreteOpacity(issues, groups) {
		state.unknown("mechanism evidence remains contractually unclassifiable without overclaim")
	}
}

func extractionHasConcreteOpacity(issues []extractionIssue, groups map[string]*extractionExpectedGroup) bool {
	if len(issues) > 0 {
		return true
	}
	for _, group := range groups {
		if group.Status == "" {
			return true
		}
	}
	return false
}

func extractionRequirePureObservation(state *evaluation, output ContractOutput) {
	var effects extractionEffects
	b, err := json.Marshal(output.EffectObservation)
	state.require(err == nil && decodeStrict(b, &effects) == nil, "effect observation is not strict")
	state.require(effects.NetworkCalls == 0, "mechanism extraction performed network I/O")
	state.require(effects.FilesystemReads == 0, "mechanism extraction performed filesystem reads")
	state.require(effects.FilesystemWrites == 0, "mechanism extraction performed filesystem writes")
	state.require(effects.ProcessSpawns == 0, "mechanism extraction spawned a process")
	state.require(effects.ExternalMutations == 0, "mechanism extraction mutated an external target")
	state.require(effects.ResolveCalls == 0 && effects.ReserveCalls == 0 && effects.IDCalls == 0 && effects.CreateCalls == 0, "mechanism extraction touched durable-run admission")
	state.require(effects.StoreInstanceIDs != nil && effects.QueueInstanceIDs != nil && effects.RootInstanceIDs != nil && effects.CursorSecretInstanceIDs != nil && effects.WorkerInstanceIDs != nil, "effect observation arrays must be [] rather than null")
	state.require(len(effects.StoreInstanceIDs) == 0 && len(effects.QueueInstanceIDs) == 0 && len(effects.RootInstanceIDs) == 0 && len(effects.CursorSecretInstanceIDs) == 0 && len(effects.WorkerInstanceIDs) == 0, "pure extraction exposed runtime instances")

	var publication extractionPublication
	b, err = json.Marshal(output.PublicationObservation)
	state.require(err == nil && decodeStrict(b, &publication) == nil, "publication observation is not strict")
	state.require(publication.ArtifactNames != nil && publication.ArtifactHashes != nil && publication.IndexHashes != nil && publication.LeakedFieldNames != nil && publication.LayerFindings != nil, "publication observation collections must not be null")
	state.require(len(publication.ArtifactNames) == 0 && len(publication.ArtifactHashes) == 0 && len(publication.IndexHashes) == 0, "pure extraction published artifacts")
	state.require(len(publication.LeakedFieldNames) == 0, "extraction publication reported leaked fields")
}

func extractionDecodeCatalog(input extractionRequest) ([]byte, extractionCatalog, bool) {
	if input.RuleCatalogSHA256 != extractionCatalogSHA256 {
		return nil, extractionCatalog{}, false
	}
	b, err := base64.StdEncoding.Strict().DecodeString(input.RuleCatalogJSON)
	if err != nil || extractionDigest(b) != extractionCatalogSHA256 {
		return nil, extractionCatalog{}, false
	}
	var catalog extractionCatalog
	if decodeStrict(b, &catalog) != nil || !extractionCatalogIsExact(catalog) {
		return nil, extractionCatalog{}, false
	}
	return b, catalog, true
}

func extractionRequireIntegrityRejection(state *evaluation, output ContractOutput) {
	state.require(rawIsNull(output.ReturnValue), "invalid catalog must not produce an extraction result")
	if output.Error == nil {
		state.fail("invalid catalog/hash was not rejected")
		return
	}
	b, err := json.Marshal(output.Error)
	var controlled extractionControlledError
	state.require(err == nil && decodeStrict(b, &controlled) == nil, "catalog rejection is not a controlled error")
	state.require(controlled.Code == "INTEGRITY_ERROR" && controlled.Message == "INTEGRITY_ERROR", "catalog/hash mismatch must be INTEGRITY_ERROR")
}

func extractionCatalogIsExact(c extractionCatalog) bool {
	if c.SchemaVersion != "tekt.mechanism-rules.v3" || c.CatalogID != "bateia-step0028-initial" || c.CatalogVersion != 3 || !c.ClosedWorld || c.UnknownClassification != "MECHANISM_CLASS_UNKNOWN" {
		return false
	}
	if c.SourceContract.EvaluatedContainer != "CandidateAssessment" || c.SourceContract.EvaluatedSource != "EvidenceSource nested under CapabilityEvidence" || c.SourceContract.RuleIdentityField != "EvidenceSource.RuleID" || c.SourceContract.BundleItem != "CollectedEvidence" || c.SourceContract.ContentUse != "hash verification only" || c.SourceContract.InferRuleFromContent {
		return false
	}
	if !extractionEqualStrings(c.SourceContract.ExactJoinFields, []string{"kind", "ref", "path", "content_hash"}) || len(c.Rules) != 4 {
		return false
	}
	documentary := map[string]string{
		"MR-APPEND-ONLY-LOG-V1":  "append-only-architecture-documentation",
		"MR-BOUNDED-QUEUE-V1":    "bounded-queue-architecture-documentation",
		"MR-CONTENT-ADDRESS-V1":  "content-address-architecture-documentation",
		"MR-ADAPTER-BOUNDARY-V1": "adapter-boundary-architecture-documentation",
	}
	seenRules := map[string]bool{}
	roleCounts := map[string]int{}
	for _, r := range c.Rules {
		if documentary[r.ID] == "" || seenRules[r.ID] || r.Version != 1 || r.StatusMax != "PROVEN" || r.AcceptedEvidenceRules == nil || len(r.AcceptedEvidenceRules) != 8 || r.Name == "" || r.Purpose == "" || r.Strategy == "" {
			return false
		}
		seenRules[r.ID] = true
		seenEvidence := map[string]bool{}
		foundDocumentary := false
		for _, evidence := range r.AcceptedEvidenceRules {
			if evidence.RuleID == "" || seenEvidence[evidence.RuleID] {
				return false
			}
			seenEvidence[evidence.RuleID] = true
			roleCounts[evidence.Role]++
			switch evidence.Role {
			case "IMPLEMENTATION":
				if evidence.WitnessKind != "IMPLEMENTATION" || evidence.StatusMax != "" {
					return false
				}
			case "INDEPENDENT":
				if !errorAllowed(evidence.WitnessKind, "TEST", "EXAMPLE", "MANIFEST", "SCHEMA", "CONFIGURATION") || evidence.StatusMax != "" {
					return false
				}
			case "DOCUMENTARY":
				if evidence.RuleID != documentary[r.ID] || evidence.WitnessKind != "DOCUMENTATION" || evidence.StatusMax != "PARTIAL" {
					return false
				}
				foundDocumentary = true
			case "CONFLICT":
				if evidence.WitnessKind != "COUNTEREVIDENCE" || evidence.StatusMax != "" {
					return false
				}
			default:
				return false
			}
		}
		if !foundDocumentary {
			return false
		}
		if r.MinimumWitnesses != (extractionMinimumWitnesses{1, 1, 1, 2}) || !extractionEqualStrings(r.FillableFields, []string{"purpose", "trigger", "strategy", "data_flow", "state_model", "failure_model", "boundaries", "preconditions", "tradeoffs"}) || !extractionEqualStrings(r.LanguageLimits, []string{"go", "java", "python", "rust", "typescript"}) || !extractionEqualStrings(r.EcosystemLimits, []string{"repository-backed source bundle"}) {
			return false
		}
	}
	return len(seenRules) == 4 && roleCounts["IMPLEMENTATION"] == 8 && roleCounts["INDEPENDENT"] == 12 && roleCounts["DOCUMENTARY"] == 4 && roleCounts["CONFLICT"] == 8
}

func extractionBindings(c extractionCatalog) map[string]extractionRuleBinding {
	out := make(map[string]extractionRuleBinding, 28)
	for _, mechanism := range c.Rules {
		for _, evidence := range mechanism.AcceptedEvidenceRules {
			out[evidence.RuleID] = extractionRuleBinding{Mechanism: mechanism, Evidence: evidence}
		}
	}
	return out
}

func extractionQualify(input extractionRequest, bindings map[string]extractionRuleBinding) (map[string]*extractionExpectedGroup, []extractionIssue) {
	groups := make(map[string]*extractionExpectedGroup)
	issues := make([]extractionIssue, 0)
	for _, assessment := range input.Assessments {
		for _, source := range assessment.Sources {
			repository, commit, identityOK := extractionParseRef(source.Ref)
			identity := extractionImmutable{Repository: repository, Commit: commit}
			if !identityOK || assessment.Repository != repository {
				issues = append(issues, extractionIssue{Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE", Loose: true})
				continue
			}
			binding, known := bindings[source.RuleID]
			if !known {
				issues = append(issues, extractionIssue{Code: "RULE_NOT_IN_MECHANISM_CATALOG", Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE"})
				continue
			}
			matches := make([]extractionBundleItem, 0, 2)
			for _, item := range input.Bundle.Items {
				if item.Kind == source.Kind && item.Ref == source.Ref && item.Path == source.Path && item.ContentHash == source.ContentHash {
					matches = append(matches, item)
				}
			}
			if len(matches) == 0 {
				issues = append(issues, extractionIssue{Code: "UNBOUND_EVALUATED_SOURCE", Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE"})
				continue
			}
			if len(matches) > 1 {
				issues = append(issues, extractionIssue{Code: "AMBIGUOUS_BUNDLE_BINDING", Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE"})
				continue
			}
			itemRepository, itemCommit, itemIdentityOK := extractionParseRef(matches[0].Ref)
			calculated := "sha256:" + extractionDigest([]byte(matches[0].ContentUTF8))
			if !itemIdentityOK || itemRepository != repository || itemCommit != commit || calculated != matches[0].ContentHash || calculated != source.ContentHash {
				issues = append(issues, extractionIssue{Code: "EVIDENCE_CONTENT_HASH_MISMATCH", Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE"})
				continue
			}
			if !extractionKindCompatible(binding.Evidence.WitnessKind, source.Kind) || !extractionValidPath(source.Path) || source.Line < 0 {
				issues = append(issues, extractionIssue{Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE", Loose: true})
				continue
			}
			witness := extractionWitness{Repository: repository, Commit: commit, Path: source.Path, Symbol: source.Symbol, Line: source.Line, SourceKind: source.Kind, Kind: binding.Evidence.WitnessKind, AssessmentStatus: assessment.Status, RuleID: source.RuleID, ContentHash: source.ContentHash}
			key := repository + "\x00" + commit + "\x00" + binding.Mechanism.ID
			group := groups[key]
			if group == nil {
				group = &extractionExpectedGroup{Rule: binding.Mechanism, Source: identity, Positive: []extractionWitness{}, Conflict: []extractionWitness{}}
				groups[key] = group
			}
			if binding.Evidence.Role == "CONFLICT" {
				if assessment.Status == "CONFLICT" {
					group.Conflict = append(group.Conflict, witness)
				} else {
					issues = append(issues, extractionIssue{Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE", Loose: true})
				}
				continue
			}
			if assessment.Status != "PROVEN" && assessment.Status != "PARTIAL" {
				issues = append(issues, extractionIssue{Source: identity, RuleID: source.RuleID, Hash: source.ContentHash, Ceiling: "INCONCLUSIVE", Loose: true})
				continue
			}
			group.Positive = append(group.Positive, witness)
		}
	}
	return groups, issues
}

func extractionEvaluateGroups(groups map[string]*extractionExpectedGroup) {
	for _, group := range groups {
		group.Positive = extractionUniqueWitnesses(group.Positive)
		group.Conflict = extractionUniqueWitnesses(group.Conflict)
		if len(group.Conflict) > 0 {
			group.Status = "INCONCLUSIVE"
			continue
		}
		impl, independent := false, false
		hashes := make(map[string]struct{})
		partialCeiling := false
		for _, witness := range group.Positive {
			hashes[witness.ContentHash] = struct{}{}
			if witness.AssessmentStatus == "PARTIAL" {
				partialCeiling = true
			}
			binding := extractionFindEvidence(group.Rule, witness.RuleID)
			if binding.Role == "IMPLEMENTATION" && witness.AssessmentStatus == "PROVEN" {
				impl = true
			}
			if binding.Role == "INDEPENDENT" && witness.Kind != "DOCUMENTATION" && witness.AssessmentStatus == "PROVEN" {
				independent = true
			}
		}
		if len(group.Positive) == 0 {
			group.Status = ""
		} else if impl && independent && len(hashes) >= 2 && !partialCeiling {
			group.Status = "PROVEN"
		} else {
			group.Status = "PARTIAL"
		}
	}
}

func extractionCheckMechanisms(state *evaluation, input extractionRequest, result extractionResult, groups map[string]*extractionExpectedGroup, issues []extractionIssue) {
	seen := make(map[string]bool)
	previousID := ""
	for i := range result.Mechanisms {
		mechanism := result.Mechanisms[i]
		state.require(i == 0 || previousID < mechanism.ID, "mechanisms are not strictly ordered by id")
		previousID = mechanism.ID
		state.require(!seen[mechanism.ID], "duplicate mechanism id")
		seen[mechanism.ID] = true
		state.require(strings.HasPrefix(mechanism.ID, "mech_") && isLowerHex(strings.TrimPrefix(mechanism.ID, "mech_"), 64), "mechanism id is not canonical")
		state.require(mechanism.Components != nil && mechanism.DataFlow != nil && mechanism.StateModel != nil && mechanism.FailureModel != nil && mechanism.Boundaries != nil && mechanism.Preconditions != nil && mechanism.Tradeoffs != nil && mechanism.Evidence != nil && mechanism.CounterEvidence != nil && mechanism.ExtractionRuleIDs != nil, "public mechanism slices must not be null")

		rule, found := extractionRuleByClass(groups, mechanism.Source, mechanism.MechanismClass)
		state.require(found, "mechanism has no qualified same-source catalog group")
		if !found {
			continue
		}
		key := mechanism.Source.Repository + "\x00" + mechanism.Source.Commit + "\x00" + rule.ID
		group := groups[key]
		state.require(group != nil && group.Status != "", "mechanism was materialized without qualifying evidence")
		if group == nil || group.Status == "" {
			continue
		}
		state.require(mechanism.MechanismClass == rule.MechanismClass && mechanism.Name == rule.Name && mechanism.Purpose == rule.Purpose && mechanism.Strategy == rule.Strategy, "mechanism catalog fields diverge from the frozen rule")
		state.require(mechanism.Source == group.Source, "mechanism source diverges from its immutable evidence source")
		state.require(mechanism.Status == group.Status, "mechanism status violates PROVEN/PARTIAL/INCONCLUSIVE precedence")
		state.require(strings.TrimSpace(mechanism.Trigger) != "" && len(mechanism.Boundaries) > 0, "mechanism central fields are incomplete")
		state.require(extractionWitnessSetsEqual(mechanism.Evidence, group.Positive), "mechanism evidence is not the exact bound, deduplicated positive witness set")
		state.require(extractionWitnessSetsEqual(mechanism.CounterEvidence, group.Conflict), "mechanism counter_evidence is not the exact catalog-derived conflict set")
		expectedRuleIDs := extractionRuleIDs(group.Positive, group.Conflict)
		state.require(extractionEqualStrings(mechanism.ExtractionRuleIDs, expectedRuleIDs), "extraction_rule_ids do not preserve the exact catalog rule provenance")
		expectedID := extractionMechanismID(mechanism, rule.ID)
		state.require(mechanism.ID == expectedID, "mechanism id does not match its canonical causal preimage")
		state.require(mechanism.Novelty == extractionExpectedNovelty(input.ConventionalBaseline, mechanism.MechanismClass, mechanism.Status), "novelty does not follow the frozen conventional baseline and status")
		extractionCheckCanonicalMechanism(state, mechanism)
	}
	for _, group := range groups {
		if group.Status == "" {
			continue
		}
		found := false
		for _, mechanism := range result.Mechanisms {
			if mechanism.Source == group.Source && mechanism.MechanismClass == group.Rule.MechanismClass {
				found = true
				break
			}
		}
		state.require(found, "qualifying same-source mechanism group was omitted")
	}
	_ = issues
}

func extractionCheckGaps(state *evaluation, gaps []extractionGap, issues []extractionIssue, groups map[string]*extractionExpectedGroup) {
	previous := ""
	for i, gap := range gaps {
		key := extractionGapOrderKey(gap)
		state.require(i == 0 || previous < key, "mechanism gaps are not strictly ordered by the contractual tuple")
		previous = key
		state.require(extractionValidRepository(gap.Source.Repository) && extractionValidCommit(gap.Source.Commit), "gap source is not immutable")
		state.require(isSHA256(gap.WitnessHash), "gap witness_hash is not sha256:<64 lower hex>")
		state.require(gap.StatusCeiling == "PROVEN" || gap.StatusCeiling == "PARTIAL" || gap.StatusCeiling == "INCONCLUSIVE", "gap has an invalid status ceiling")
		state.require(gap.Code != "" && len(gap.Code) <= 128 && len(gap.RuleID) <= 128 && extractionSafeText(gap.Code) && extractionSafeText(gap.RuleID), "gap contains unsafe public text")
		supported := false
		for _, issue := range issues {
			if extractionGapMatchesIssue(gap, issue) {
				supported = true
				break
			}
		}
		if !supported {
			for _, group := range groups {
				if gap.Source == group.Source && group.Status != "PROVEN" && extractionGapRelatesToGroup(gap, group) {
					supported = true
					break
				}
			}
		}
		state.require(supported, "mechanism gap has no real request-derived ambiguity or evidence limitation")
	}
	for _, issue := range issues {
		matched := false
		for _, gap := range gaps {
			if extractionGapMatchesIssue(gap, issue) {
				matched = true
				break
			}
		}
		state.require(matched, "a real extraction ambiguity/failure was not preserved as a gap")
	}
	for _, group := range groups {
		if group.Status != "" {
			continue
		}
		matched := false
		for _, gap := range gaps {
			if gap.Source == group.Source && gap.StatusCeiling == "INCONCLUSIVE" {
				matched = true
				break
			}
		}
		state.require(matched, "insufficient witness set disappeared instead of producing an Unknown-bearing gap")
	}
}

func extractionGapMatchesIssue(gap extractionGap, issue extractionIssue) bool {
	if issue.Code != "" && gap.Code != issue.Code {
		return false
	}
	if issue.Source.Repository != "" && gap.Source.Repository != issue.Source.Repository {
		return false
	}
	if issue.Source.Commit != "" && gap.Source.Commit != issue.Source.Commit {
		return false
	}
	if issue.RuleID != "" && gap.RuleID != issue.RuleID {
		return false
	}
	if !issue.Loose && issue.Hash != "" && gap.WitnessHash != issue.Hash {
		return false
	}
	return issue.Ceiling == "" || gap.StatusCeiling == issue.Ceiling
}

func extractionGapRelatesToGroup(gap extractionGap, group *extractionExpectedGroup) bool {
	if gap.RuleID == group.Rule.ID {
		return true
	}
	for _, witness := range append(append([]extractionWitness{}, group.Positive...), group.Conflict...) {
		if gap.RuleID == witness.RuleID && gap.WitnessHash == witness.ContentHash {
			return true
		}
	}
	return false
}

func extractionCheckSafety(state *evaluation, input extractionRequest, result extractionResult) {
	b, err := json.Marshal(result)
	state.require(err == nil && utf8.Valid(b), "extraction result is not valid UTF-8 JSON")
	low := strings.ToLower(string(b))
	for _, forbidden := range []string{"snippet", "content_utf8", "checkout_path", "idempotency_key", "authorization", "cookie", "private_prompt"} {
		state.require(!strings.Contains(low, forbidden), "extraction output exposes a forbidden raw/private field")
	}
	for _, item := range input.Bundle.Items {
		for _, token := range extractionLeakTokens(item.ContentUTF8) {
			state.require(!strings.Contains(string(b), token), "extraction output leaked a token from raw evidence content")
		}
	}
}

func extractionCheckCanonicalMechanism(state *evaluation, mechanism extractionMechanism) {
	for _, value := range []string{mechanism.Name, mechanism.Purpose, mechanism.Trigger, mechanism.Strategy} {
		state.require(extractionSafeText(value), "mechanism contains non-canonical or controlled text")
	}
	for _, values := range [][]string{mechanism.DataFlow, mechanism.StateModel, mechanism.FailureModel, mechanism.Boundaries, mechanism.Preconditions, mechanism.Tradeoffs, mechanism.ExtractionRuleIDs} {
		state.require(extractionCanonicalStrings(values), "mechanism set is not normalized, deduplicated and byte-sorted")
	}
	state.require(extractionCanonicalComponents(mechanism.Components), "components are not a canonical set")
	state.require(extractionCanonicalWitnesses(mechanism.Evidence) && extractionCanonicalWitnesses(mechanism.CounterEvidence), "witnesses are not deduplicated and canonically ordered")
	for _, witness := range append(append([]extractionWitness{}, mechanism.Evidence...), mechanism.CounterEvidence...) {
		state.require(witness.Repository == mechanism.Source.Repository && witness.Commit == mechanism.Source.Commit, "mechanism combines cross-repository or cross-commit evidence")
		state.require(extractionValidPath(witness.Path) && witness.Line >= 0 && isSHA256(witness.ContentHash), "witness path/line/hash is invalid")
	}
}

func extractionRuleByClass(groups map[string]*extractionExpectedGroup, source extractionImmutable, class string) (extractionCatalogMechanism, bool) {
	for _, group := range groups {
		if group.Source == source && group.Rule.MechanismClass == class {
			return group.Rule, true
		}
	}
	return extractionCatalogMechanism{}, false
}

func extractionFindEvidence(rule extractionCatalogMechanism, id string) extractionCatalogEvidence {
	for _, evidence := range rule.AcceptedEvidenceRules {
		if evidence.RuleID == id {
			return evidence
		}
	}
	return extractionCatalogEvidence{}
}

func extractionExpectedNovelty(baseline *extractionBaseline, class, status string) string {
	if !extractionBaselineValid(baseline) {
		return "UNKNOWN"
	}
	for _, candidate := range baseline.MechanismClasses {
		if candidate == class {
			return "CONVENTIONAL"
		}
	}
	if status == "PROVEN" {
		return "DIVERGENT"
	}
	return "UNKNOWN"
}

func extractionBaselineValid(b *extractionBaseline) bool {
	if b == nil || strings.TrimSpace(b.BaselineID) == "" || !isLowerHex(b.IdeaHash, 64) || b.MechanismClasses == nil {
		return false
	}
	for _, class := range b.MechanismClasses {
		if !extractionSafeText(class) {
			return false
		}
	}
	return true
}

func extractionMechanismID(mechanism extractionMechanism, ruleID string) string {
	hashes := make([]string, 0, len(mechanism.Evidence)+len(mechanism.CounterEvidence))
	for _, witness := range mechanism.Evidence {
		hashes = append(hashes, witness.ContentHash)
	}
	for _, witness := range mechanism.CounterEvidence {
		hashes = append(hashes, witness.ContentHash)
	}
	for i := range hashes {
		hashes[i] = extractionNormalizeCanonicalText(hashes[i])
	}
	hashes = extractionSortedUnique(hashes)
	preimage := struct {
		SchemaVersion  string              `json:"schema_version"`
		Source         extractionImmutable `json:"source"`
		RuleID         string              `json:"rule_id"`
		MechanismClass string              `json:"mechanism_class"`
		Purpose        string              `json:"purpose"`
		Strategy       string              `json:"strategy"`
		WitnessHashes  []string            `json:"witness_hashes"`
	}{
		SchemaVersion: "1",
		Source: extractionImmutable{
			Repository: extractionNormalizeCanonicalText(mechanism.Source.Repository),
			Commit:     extractionNormalizeCanonicalText(mechanism.Source.Commit),
		},
		RuleID:         extractionNormalizeCanonicalText(ruleID),
		MechanismClass: extractionNormalizeCanonicalText(mechanism.MechanismClass),
		Purpose:        extractionNormalizeCanonicalText(mechanism.Purpose),
		Strategy:       extractionNormalizeCanonicalText(mechanism.Strategy),
		WitnessHashes:  hashes,
	}
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(preimage) // fixed string-only struct cannot fail JSON encoding
	b := bytes.TrimSuffix(encoded.Bytes(), []byte{'\n'})
	return "mech_" + extractionDigest(b)
}

// extractionNormalizeCanonicalText is the local, I/O-free normalization used
// by causal preimages. It collapses Unicode whitespace and composes the
// canonical Latin and Hangul sequences that can otherwise have two UTF-8 byte
// representations. Unsupported combining sequences are preserved rather than
// guessed; public representation checks still reject every non-canonical
// sequence that this oracle can prove locally.
func extractionNormalizeCanonicalText(value string) string {
	value = strings.Join(strings.FieldsFunc(value, unicode.IsSpace), " ")
	runes := []rune(value)
	if len(runes) < 2 {
		return value
	}
	out := make([]rune, 0, len(runes))
	for _, current := range runes {
		if len(out) > 0 {
			prior := out[len(out)-1]
			if composed, ok := extractionComposeHangul(prior, current); ok {
				out[len(out)-1] = composed
				continue
			}
			if composed, ok := extractionLatinComposition[[2]rune{prior, current}]; ok {
				out[len(out)-1] = composed
				continue
			}
		}
		out = append(out, current)
	}
	return string(out)
}

func extractionComposeHangul(first, second rune) (rune, bool) {
	const (
		hangulSBase  = 0xAC00
		hangulLBase  = 0x1100
		hangulVBase  = 0x1161
		hangulTBase  = 0x11A7
		hangulLCount = 19
		hangulVCount = 21
		hangulTCount = 28
		hangulNCount = hangulVCount * hangulTCount
		hangulSCount = hangulLCount * hangulNCount
	)
	if first >= hangulLBase && first < hangulLBase+hangulLCount && second >= hangulVBase && second < hangulVBase+hangulVCount {
		return hangulSBase + (first-hangulLBase)*hangulNCount + (second-hangulVBase)*hangulTCount, true
	}
	if first >= hangulSBase && first < hangulSBase+hangulSCount && (first-hangulSBase)%hangulTCount == 0 && second > hangulTBase && second < hangulTBase+hangulTCount {
		return first + second - hangulTBase, true
	}
	return 0, false
}

// The table covers canonical Latin compositions used by public engineering
// text. Keeping it local avoids filesystem, environment and dependency lookup
// in the oracle authority.
var extractionLatinComposition = map[[2]rune]rune{
	{'A', '\u0300'}: 'À', {'A', '\u0301'}: 'Á', {'A', '\u0302'}: 'Â', {'A', '\u0303'}: 'Ã', {'A', '\u0308'}: 'Ä', {'A', '\u030a'}: 'Å',
	{'C', '\u0327'}: 'Ç', {'E', '\u0300'}: 'È', {'E', '\u0301'}: 'É', {'E', '\u0302'}: 'Ê', {'E', '\u0308'}: 'Ë',
	{'I', '\u0300'}: 'Ì', {'I', '\u0301'}: 'Í', {'I', '\u0302'}: 'Î', {'I', '\u0308'}: 'Ï', {'N', '\u0303'}: 'Ñ',
	{'O', '\u0300'}: 'Ò', {'O', '\u0301'}: 'Ó', {'O', '\u0302'}: 'Ô', {'O', '\u0303'}: 'Õ', {'O', '\u0308'}: 'Ö',
	{'U', '\u0300'}: 'Ù', {'U', '\u0301'}: 'Ú', {'U', '\u0302'}: 'Û', {'U', '\u0308'}: 'Ü', {'Y', '\u0301'}: 'Ý',
	{'a', '\u0300'}: 'à', {'a', '\u0301'}: 'á', {'a', '\u0302'}: 'â', {'a', '\u0303'}: 'ã', {'a', '\u0308'}: 'ä', {'a', '\u030a'}: 'å',
	{'c', '\u0327'}: 'ç', {'e', '\u0300'}: 'è', {'e', '\u0301'}: 'é', {'e', '\u0302'}: 'ê', {'e', '\u0308'}: 'ë',
	{'i', '\u0300'}: 'ì', {'i', '\u0301'}: 'í', {'i', '\u0302'}: 'î', {'i', '\u0308'}: 'ï', {'n', '\u0303'}: 'ñ',
	{'o', '\u0300'}: 'ò', {'o', '\u0301'}: 'ó', {'o', '\u0302'}: 'ô', {'o', '\u0303'}: 'õ', {'o', '\u0308'}: 'ö',
	{'u', '\u0300'}: 'ù', {'u', '\u0301'}: 'ú', {'u', '\u0302'}: 'û', {'u', '\u0308'}: 'ü', {'y', '\u0301'}: 'ý', {'y', '\u0308'}: 'ÿ',
}

func extractionWitnessSetsEqual(actual, expected []extractionWitness) bool {
	a := extractionUniqueWitnesses(actual)
	e := extractionUniqueWitnesses(expected)
	if len(actual) != len(a) || len(expected) != len(e) || len(a) != len(e) {
		return false
	}
	for i := range a {
		if a[i] != e[i] {
			return false
		}
	}
	return true
}

func extractionUniqueWitnesses(in []extractionWitness) []extractionWitness {
	out := append([]extractionWitness(nil), in...)
	sort.Slice(out, func(i, j int) bool { return extractionWitnessKey(out[i]) < extractionWitnessKey(out[j]) })
	dedup := out[:0]
	for _, witness := range out {
		if len(dedup) == 0 || extractionWitnessKey(dedup[len(dedup)-1]) != extractionWitnessKey(witness) {
			dedup = append(dedup, witness)
		}
	}
	return dedup
}

func extractionWitnessKey(w extractionWitness) string {
	b, _ := json.Marshal(w)
	return string(b)
}

func extractionRuleIDs(parts ...[]extractionWitness) []string {
	var ids []string
	for _, witnesses := range parts {
		for _, witness := range witnesses {
			ids = append(ids, witness.RuleID)
		}
	}
	return extractionSortedUnique(ids)
}

func extractionKindCompatible(witnessKind, sourceKind string) bool {
	switch witnessKind {
	case "IMPLEMENTATION":
		return sourceKind == "file" || sourceKind == "symbol"
	case "TEST":
		return sourceKind == "test"
	case "EXAMPLE":
		return sourceKind == "example"
	case "MANIFEST":
		return sourceKind == "manifest"
	case "DOCUMENTATION":
		return sourceKind == "documentation"
	case "SCHEMA", "CONFIGURATION":
		return sourceKind == "file" || sourceKind == "manifest"
	case "COUNTEREVIDENCE":
		return true
	default:
		return false
	}
}

func extractionParseRef(ref string) (string, string, bool) {
	cut := strings.LastIndexByte(ref, '@')
	if cut <= 0 || cut == len(ref)-1 || strings.Contains(ref[:cut], "@") {
		return "", "", false
	}
	repository, commit := ref[:cut], ref[cut+1:]
	return repository, commit, extractionValidRepository(repository) && extractionValidCommit(commit)
}

func extractionValidRepository(repository string) bool {
	if strings.ContainsAny(repository, " \t\r\n") || strings.Count(repository, "/") != 1 {
		return false
	}
	parts := strings.Split(repository, "/")
	return parts[0] != "" && parts[1] != "" && !strings.Contains(repository, "://")
}

func extractionValidCommit(commit string) bool {
	return isLowerHex(commit, 40) || isLowerHex(commit, 64)
}

func extractionValidPath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "\\") {
		return false
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return extractionSafeText(path)
}

func extractionSafeText(value string) bool {
	if !utf8.ValidString(value) || strings.TrimSpace(value) != value {
		return false
	}
	if extractionNormalizeCanonicalText(value) != value {
		return false
	}
	space := false
	for _, r := range value {
		if r <= 0x1f || (r >= 0x7f && r <= 0x9f) {
			return false
		}
		if unicode.IsSpace(r) {
			if space || r != ' ' {
				return false
			}
			space = true
		} else {
			space = false
		}
	}
	return true
}

func extractionCanonicalStrings(values []string) bool {
	for i, value := range values {
		if !extractionSafeText(value) || (i > 0 && values[i-1] >= value) {
			return false
		}
	}
	return true
}

func extractionCanonicalComponents(values []extractionComponent) bool {
	previous := ""
	for i, value := range values {
		if value.ID == "" || !extractionSafeText(value.Name) || !extractionSafeText(value.Kind) {
			return false
		}
		b, _ := json.Marshal(value)
		if i > 0 && previous >= string(b) {
			return false
		}
		previous = string(b)
	}
	return true
}

func extractionCanonicalWitnesses(values []extractionWitness) bool {
	for i, witness := range values {
		if i > 0 && extractionWitnessKey(values[i-1]) >= extractionWitnessKey(witness) {
			return false
		}
	}
	return true
}

func extractionGapOrderKey(g extractionGap) string {
	return g.Code + "\x00" + g.Source.Repository + "\x00" + g.Source.Commit + "\x00" + g.RuleID + "\x00" + g.WitnessHash
}

func extractionLeakTokens(content string) []string {
	var out []string
	for _, field := range strings.FieldsFunc(content, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("_-.:/@", r))
	}) {
		lower := strings.ToLower(field)
		if len(field) >= 16 && (strings.Contains(lower, "canary") || strings.Contains(lower, "secret") || strings.Contains(lower, "token") || strings.Contains(lower, "private")) {
			out = append(out, field)
		}
	}
	return extractionSortedUnique(out)
}

func extractionSortedUnique(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	dedup := out[:0]
	for _, value := range out {
		if len(dedup) == 0 || dedup[len(dedup)-1] != value {
			dedup = append(dedup, value)
		}
	}
	return dedup
}

func extractionEqualStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func extractionDigest(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
