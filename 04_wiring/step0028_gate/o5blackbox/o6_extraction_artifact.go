/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o6_extraction_artifact.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// o6Check prevents an opacity classification from masking a violation found in
// the same contractual observation.
type o6Check struct {
	state *evaluation
	valid bool
}

func newO6Check(state *evaluation) *o6Check {
	return &o6Check{state: state, valid: true}
}

func (c *o6Check) require(condition bool, reason string) {
	c.state.require(condition, reason)
	if !condition {
		c.valid = false
	}
}

func (c *o6Check) fail(reason string) {
	c.valid = false
	c.state.fail(reason)
}

func (c *o6Check) unknownIf(condition bool, reason string) {
	if c.valid && condition {
		c.state.unknown(reason)
	}
}

// evaluateMechanismExtractO6 is the O6 entry point for MECHANISM_EXTRACT.  It
// intentionally treats only contractually declared sets as sets. Data-flow,
// state, failure, component and witness presentation order is not promoted into
// an observable unless R3 explicitly makes it one.
func evaluateMechanismExtractO6(state *evaluation, request json.RawMessage, output ContractOutput) {
	check := newO6Check(state)
	var input extractionRequest
	if err := decodeStrict(request, &input); err != nil {
		check.fail("MECHANISM_EXTRACT request is not strict protocol JSON")
		return
	}

	o6CheckExtractionEffects(check, output)
	catalogBytes, catalog, catalogOK := extractionDecodeCatalog(input)
	if !catalogOK {
		o6CheckIntegrityRejection(check, output)
		return
	}
	check.require(input.RuleCatalogSHA256 == extractionCatalogSHA256, "extractor did not receive the frozen catalog hash")
	check.require(extractionDigest(catalogBytes) == extractionCatalogSHA256, "decoded catalog bytes do not match the frozen catalog hash")
	if output.Error != nil {
		check.fail("valid MECHANISM_EXTRACT request returned a controlled error")
		return
	}
	if rawIsNull(output.ReturnValue) {
		check.fail("valid MECHANISM_EXTRACT request returned null")
		return
	}

	var result extractionResult
	if err := decodeReturn(output.ReturnValue, &result); err != nil {
		check.fail("MECHANISM_EXTRACT return_value is not a strict extraction_result")
		return
	}
	check.require(result.Mechanisms != nil, "mechanisms must be [] rather than null")
	check.require(result.Gaps != nil, "gaps must be [] rather than null")

	groups, issues := extractionQualify(input, extractionBindings(catalog))
	extractionEvaluateGroups(groups)
	o6CheckExtractionMechanisms(check, input, result.Mechanisms, groups)
	matchedIssues := o6CheckExtractionGaps(check, result.Gaps, issues, groups)
	o6CheckExtractionSafety(check, input, result)

	hasInconclusive := false
	for _, mechanism := range result.Mechanisms {
		switch mechanism.Status {
		case "INCONCLUSIVE":
			hasInconclusive = true
		}
	}
	// A resolved mechanism does not erase an independent unresolved source,
	// identity or qualification ambiguity in the same extraction result.
	// A missing/invalid conventional baseline is observable as a real limit
	// only when it is what prevents classifying a PROVEN mechanism's novelty.
	// PARTIAL mechanisms remain UNKNOWN independently of baseline availability.
	baselineLimitsNovelty := false
	if !extractionBaselineValid(input.ConventionalBaseline) {
		for _, mechanism := range result.Mechanisms {
			if mechanism.Status == "PROVEN" && mechanism.Novelty == "UNKNOWN" {
				baselineLimitsNovelty = true
				break
			}
		}
	}
	minimumWitnessInsufficient := false
	for _, group := range groups {
		if o15GroupHasSemanticWitnessLimit(group) {
			minimumWitnessInsufficient = true
			break
		}
	}
	for _, gap := range result.Gaps {
		upperCode := strings.ToUpper(gap.Code)
		if strings.Contains(upperCode, "INSUFFICIENT") && (gap.StatusCeiling == "PARTIAL" || gap.StatusCeiling == "INCONCLUSIVE") && o6GapHasProvenance(gap, issues, groups) {
			minimumWitnessInsufficient = true
			break
		}
	}
	realOpacity := (len(issues) > 0 && matchedIssues == len(issues)) || hasInconclusive || baselineLimitsNovelty || minimumWitnessInsufficient
	check.unknownIf(realOpacity, "mechanism extraction preserves a real R4 ambiguity without a simultaneous violation")
}

func o6CheckExtractionEffects(check *o6Check, output ContractOutput) {
	var effects extractionEffects
	b, err := json.Marshal(output.EffectObservation)
	check.require(err == nil && decodeStrict(b, &effects) == nil, "effect observation is not strict")
	check.require(effects.NetworkCalls == 0, "mechanism extraction performed network I/O")
	check.require(effects.FilesystemReads == 0, "mechanism extraction performed filesystem reads")
	check.require(effects.FilesystemWrites == 0, "mechanism extraction performed filesystem writes")
	check.require(effects.ProcessSpawns == 0, "mechanism extraction spawned a process")
	check.require(effects.ExternalMutations == 0, "mechanism extraction mutated an external target")
	check.require(effects.ResolveCalls == 0 && effects.ReserveCalls == 0 && effects.IDCalls == 0 && effects.CreateCalls == 0, "mechanism extraction touched durable-run admission")
	check.require(effects.StoreInstanceIDs != nil && effects.QueueInstanceIDs != nil && effects.RootInstanceIDs != nil && effects.CursorSecretInstanceIDs != nil && effects.WorkerInstanceIDs != nil, "effect observation arrays must be [] rather than null")
	check.require(len(effects.StoreInstanceIDs) == 0 && len(effects.QueueInstanceIDs) == 0 && len(effects.RootInstanceIDs) == 0 && len(effects.CursorSecretInstanceIDs) == 0 && len(effects.WorkerInstanceIDs) == 0, "pure extraction exposed runtime instances")

	var publication extractionPublication
	b, err = json.Marshal(output.PublicationObservation)
	check.require(err == nil && decodeStrict(b, &publication) == nil, "publication observation is not strict")
	check.require(publication.ArtifactNames != nil && publication.ArtifactHashes != nil && publication.IndexHashes != nil && publication.LeakedFieldNames != nil && publication.LayerFindings != nil, "publication observation collections must not be null")
	check.require(len(publication.ArtifactNames) == 0 && len(publication.ArtifactHashes) == 0 && len(publication.IndexHashes) == 0, "pure extraction published artifacts")
	check.require(len(publication.LeakedFieldNames) == 0, "extraction publication reported leaked fields")
}

func o6CheckIntegrityRejection(check *o6Check, output ContractOutput) {
	check.require(rawIsNull(output.ReturnValue), "invalid catalog must not produce an extraction result")
	if output.Error == nil {
		check.fail("invalid catalog/hash was not rejected")
		return
	}
	b, err := json.Marshal(output.Error)
	var controlled extractionControlledError
	check.require(err == nil && decodeStrict(b, &controlled) == nil, "catalog rejection is not a controlled error")
	check.require(controlled.Code == "INTEGRITY_ERROR" && controlled.Message == "INTEGRITY_ERROR", "catalog/hash mismatch must be INTEGRITY_ERROR")
}

func o6CheckExtractionMechanisms(check *o6Check, input extractionRequest, mechanisms []extractionMechanism, groups map[string]*extractionExpectedGroup) {
	seenIDs := make(map[string]struct{}, len(mechanisms))
	previousID := ""
	materialized := make(map[string]struct{}, len(mechanisms))
	for i, mechanism := range mechanisms {
		check.require(i == 0 || previousID < mechanism.ID, "mechanisms are not strictly ordered by id")
		previousID = mechanism.ID
		_, duplicate := seenIDs[mechanism.ID]
		check.require(!duplicate, "duplicate mechanism id")
		seenIDs[mechanism.ID] = struct{}{}
		check.require(strings.HasPrefix(mechanism.ID, "mech_") && isLowerHex(strings.TrimPrefix(mechanism.ID, "mech_"), 64), "mechanism id is not canonical")
		check.require(mechanism.Components != nil && mechanism.DataFlow != nil && mechanism.StateModel != nil && mechanism.FailureModel != nil && mechanism.Boundaries != nil && mechanism.Preconditions != nil && mechanism.Tradeoffs != nil && mechanism.Evidence != nil && mechanism.CounterEvidence != nil && mechanism.ExtractionRuleIDs != nil, "public mechanism slices must not be null")
		check.require(len(mechanism.Components) <= 32 && len(mechanism.DataFlow) <= 32 && len(mechanism.StateModel) <= 32 && len(mechanism.FailureModel) <= 32 && len(mechanism.Boundaries) <= 32 && len(mechanism.Preconditions) <= 32 && len(mechanism.Tradeoffs) <= 32, "mechanism structural field exceeds its limit")
		check.require(len(mechanism.Evidence) <= 64 && len(mechanism.CounterEvidence) <= 64 && len(mechanism.ExtractionRuleIDs) <= 64, "mechanism provenance exceeds its limit")

		rule, found := extractionRuleByClass(groups, mechanism.Source, mechanism.MechanismClass)
		check.require(found, "mechanism has no qualified same-source catalog group")
		if !found {
			continue
		}
		key := mechanism.Source.Repository + "\x00" + mechanism.Source.Commit + "\x00" + rule.ID
		group := groups[key]
		check.require(group != nil && group.Status != "", "mechanism was materialized without qualifying evidence")
		if group == nil || group.Status == "" {
			continue
		}
		materialized[key] = struct{}{}
		check.require(mechanism.MechanismClass == rule.MechanismClass && mechanism.Name == rule.Name && mechanism.Purpose == rule.Purpose && mechanism.Strategy == rule.Strategy, "mechanism catalog fields diverge from the frozen rule")
		check.require(mechanism.Source == group.Source, "mechanism source diverges from its immutable evidence source")
		check.require(mechanism.Status == group.Status, "mechanism status violates PROVEN/PARTIAL/INCONCLUSIVE precedence")
		check.require(strings.TrimSpace(mechanism.Trigger) != "" && len(mechanism.Boundaries) > 0, "mechanism central fields are incomplete")
		check.require(extractionWitnessSetsEqual(mechanism.Evidence, group.Positive), "mechanism evidence is not the exact bound, deduplicated positive witness set")
		check.require(extractionWitnessSetsEqual(mechanism.CounterEvidence, group.Conflict), "mechanism counter_evidence is not the exact catalog-derived conflict set")
		check.require(extractionEqualStrings(mechanism.ExtractionRuleIDs, extractionRuleIDs(group.Positive, group.Conflict)), "extraction_rule_ids do not preserve canonical rule provenance")
		check.require(mechanism.ID == extractionMechanismID(mechanism, rule.ID), "mechanism id does not match its canonical causal preimage")
		check.require(mechanism.Novelty == extractionExpectedNovelty(input.ConventionalBaseline, mechanism.MechanismClass, mechanism.Status), "novelty does not follow baseline and status")
		o6CheckMechanismRepresentation(check, mechanism)
	}
	for key, group := range groups {
		if group.Status == "" {
			continue
		}
		_, found := materialized[key]
		check.require(found, "qualifying same-source mechanism group was omitted")
	}
}

func o6CheckMechanismRepresentation(check *o6Check, mechanism extractionMechanism) {
	for _, value := range []string{mechanism.Name, mechanism.Purpose, mechanism.Trigger, mechanism.Strategy} {
		check.require(value != "" && extractionSafeText(value), "mechanism contains unsafe public text")
	}
	for _, values := range [][]string{mechanism.DataFlow, mechanism.StateModel, mechanism.FailureModel, mechanism.Boundaries, mechanism.Preconditions, mechanism.Tradeoffs} {
		for _, value := range values {
			check.require(value != "" && extractionSafeText(value), "mechanism structural sequence contains unsafe text")
		}
	}
	for _, component := range mechanism.Components {
		check.require(component.ID != "" && extractionSafeText(component.ID) && component.Name != "" && extractionSafeText(component.Name) && component.Kind != "" && extractionSafeText(component.Kind), "component contains unsafe or empty fields")
	}
	// ExtractionRuleIDs is explicitly a canonical set. The remaining public
	// slices retain their domain order; equality of witness sets is checked above.
	check.require(extractionCanonicalStrings(mechanism.ExtractionRuleIDs), "extraction_rule_ids is not a canonical set")
	for _, witness := range append(append([]extractionWitness{}, mechanism.Evidence...), mechanism.CounterEvidence...) {
		check.require(witness.Repository == mechanism.Source.Repository && witness.Commit == mechanism.Source.Commit, "mechanism combines cross-source evidence")
		check.require(extractionValidPath(witness.Path) && witness.Line >= 0 && isSHA256(witness.ContentHash), "witness path/line/hash is invalid")
		check.require(extractionSafeText(witness.Symbol), "witness symbol contains unsafe text")
	}
}

func o6CheckExtractionGaps(check *o6Check, gaps []extractionGap, issues []extractionIssue, groups map[string]*extractionExpectedGroup) int {
	previous := ""
	for i, gap := range gaps {
		key := extractionGapOrderKey(gap)
		check.require(i == 0 || previous < key, "mechanism gaps are not strictly ordered by the R3 tuple")
		previous = key
		check.require(gap.Code != "" && len(gap.Code) <= 128 && extractionSafeText(gap.Code), "gap code is invalid")
		check.require(len(gap.RuleID) <= 128 && extractionSafeText(gap.RuleID), "gap rule_id is invalid")
		check.require(extractionValidRepository(gap.Source.Repository) && extractionValidCommit(gap.Source.Commit), "gap source is not immutable")
		check.require(isSHA256(gap.WitnessHash), "gap witness_hash is invalid")
		check.require(gap.StatusCeiling == "PROVEN" || gap.StatusCeiling == "PARTIAL" || gap.StatusCeiling == "INCONCLUSIVE", "gap status ceiling is invalid")
		check.require(o6GapHasProvenance(gap, issues, groups), "gap has no request-derived provenance")
	}
	matched := 0
	for _, issue := range issues {
		found := false
		for _, gap := range gaps {
			if o6GapMatchesIssue(gap, issue) {
				found = true
				break
			}
		}
		check.require(found, "a real extraction ambiguity/failure was not preserved as a gap")
		if found {
			matched++
		}
	}
	for _, group := range groups {
		if group.Status != "" {
			continue
		}
		found := false
		for _, gap := range gaps {
			if gap.Source == group.Source && gap.StatusCeiling == "INCONCLUSIVE" {
				found = true
				break
			}
		}
		check.require(found, "insufficient witnesses disappeared instead of producing an Unknown-bearing gap")
	}
	return matched
}

func o6GapMatchesIssue(gap extractionGap, issue extractionIssue) bool {
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
	if !issue.Loose && issue.Hash != "" && isSHA256(issue.Hash) && gap.WitnessHash != issue.Hash {
		return false
	}
	return issue.Ceiling == "" || gap.StatusCeiling == issue.Ceiling
}

func o6GapHasProvenance(gap extractionGap, issues []extractionIssue, groups map[string]*extractionExpectedGroup) bool {
	for _, issue := range issues {
		if o6GapMatchesIssue(gap, issue) {
			return true
		}
	}
	for _, group := range groups {
		if gap.Source != group.Source {
			continue
		}
		if gap.RuleID == group.Rule.ID {
			return true
		}
		for _, witness := range append(append([]extractionWitness{}, group.Positive...), group.Conflict...) {
			if gap.RuleID == witness.RuleID && gap.WitnessHash == witness.ContentHash {
				return true
			}
		}
	}
	return false
}

func o6CheckExtractionSafety(check *o6Check, input extractionRequest, result extractionResult) {
	b, err := json.Marshal(result)
	check.require(err == nil && utf8.Valid(b), "extraction result is not valid UTF-8 JSON")
	low := strings.ToLower(string(b))
	for _, forbidden := range []string{"snippet", "content_utf8", "checkout_path", "idempotency_key", "authorization", "cookie", "private_prompt"} {
		check.require(!strings.Contains(low, forbidden), "extraction output exposes a forbidden raw/private field")
	}
	for _, item := range input.Bundle.Items {
		for _, token := range extractionLeakTokens(item.ContentUTF8) {
			check.require(!strings.Contains(string(b), token), "extraction output leaked raw evidence content")
		}
	}
}

type o6ArtifactRequest struct {
	RunID    string  `json:"run_id"`
	Artifact string  `json:"artifact"`
	Cursor   *string `json:"cursor"`
}

type o6ArtifactResult struct {
	Artifact      string  `json:"artifact"`
	SHA256        string  `json:"sha256"`
	Size          int64   `json:"size"`
	CanonicalJSON string  `json:"canonical_json"`
	NextCursor    *string `json:"next_cursor"`
}

type o6ArtifactDocument struct {
	SchemaVersion string
	Artifact      string
	EntriesRaw    json.RawMessage
	Entries       []any
	SemanticJSON  []byte
}

// evaluateArtifactReadO6 is the O6 entry point for ARTIFACT_READ. Artifact
// bytes use their declared fixed struct-field order; TEKT_CANONICAL_JSON_V1 is
// computed separately and is never substituted for those publication bytes.
func evaluateArtifactReadO6(state *evaluation, request json.RawMessage, output ContractOutput) {
	check := newO6Check(state)
	var input o6ArtifactRequest
	if err := decodeStrict(request, &input); err != nil {
		check.fail("ARTIFACT_READ request is not strict protocol JSON")
		return
	}
	check.require(strings.HasPrefix(input.RunID, "run_") && isLowerHex(strings.TrimPrefix(input.RunID, "run_"), 32), "artifact request run_id is invalid")
	check.require(o6AllowedArtifact(input.Artifact), "artifact request is outside the closed allowlist")
	if input.Cursor != nil {
		check.require(*input.Cursor != "" && len(*input.Cursor) <= 512 && extractionSafeText(*input.Cursor), "artifact cursor is not an opaque bounded token")
	}
	o6CheckArtifactEffects(check, output)

	if output.Error != nil {
		o6EvaluateArtifactError(check, input, output)
		return
	}
	if rawIsNull(output.ReturnValue) {
		check.fail("successful ARTIFACT_READ returned null")
		return
	}
	var result o6ArtifactResult
	if err := decodeReturn(output.ReturnValue, &result); err != nil {
		check.fail("ARTIFACT_READ return_value is not a strict artifact_result")
		return
	}
	check.require(result.Artifact == input.Artifact, "artifact response changed the requested artifact identity")
	check.require(isLowerHex(result.SHA256, 64), "artifact response sha256 is not 64 lower hex")
	check.require(result.Size >= 0, "artifact size is negative")
	if result.NextCursor != nil {
		check.require(*result.NextCursor != "" && len(*result.NextCursor) <= 512 && extractionSafeText(*result.NextCursor), "next_cursor is not an opaque bounded token")
		if input.Cursor != nil {
			check.require(*result.NextCursor != *input.Cursor, "pagination returned the same cursor and cannot advance")
		}
	}

	artifactBytes, err := base64.StdEncoding.Strict().DecodeString(result.CanonicalJSON)
	check.require(err == nil, "canonical_json is not strict base64")
	if err != nil {
		return
	}
	document, err := o6DecodeFixedArtifactDocument(artifactBytes)
	check.require(err == nil, "artifact bytes are not fixed-field-order canonical JSON")
	if err != nil {
		return
	}
	check.require(document.SchemaVersion == "1", "artifact document schema_version is not 1")
	check.require(document.Artifact == input.Artifact && document.Artifact == result.Artifact, "artifact document identity is inconsistent")
	check.require(result.Size >= int64(len(artifactBytes)), "artifact page is larger than the declared artifact size")
	if input.Cursor == nil && result.NextCursor == nil {
		check.require(result.Size == int64(len(artifactBytes)), "unpaginated artifact size does not match its bytes")
	}
	if result.Size == int64(len(artifactBytes)) {
		check.require(result.SHA256 == extractionDigest(artifactBytes), "artifact sha256 is not over the fixed-order publication bytes")
		if !bytes.Equal(artifactBytes, document.SemanticJSON) {
			check.require(result.SHA256 != extractionDigest(document.SemanticJSON), "artifact sha256 incorrectly uses TEKT semantic canonical bytes")
		}
	}
	check.require(o6ArtifactEntriesOrdered(input.Artifact, document.Entries), "artifact entries violate their declared order")
	check.require(!o6ArtifactContainsForbiddenFields(document.Entries), "artifact contains forbidden raw/private fields")

	realOpacity := o23OpaqueArtifactRule(document, result, output.PublicationObservation) || o6ArtifactHasRealOpacity(input.Artifact, document.Entries)
	check.unknownIf(realOpacity, "artifact preserves a real contractual Unknown without a simultaneous violation")
}

func o6EvaluateArtifactError(check *o6Check, input o6ArtifactRequest, output ContractOutput) {
	check.require(rawIsNull(output.ReturnValue), "ARTIFACT_READ error must return null")
	b, err := json.Marshal(output.Error)
	var controlled extractionControlledError
	check.require(err == nil && decodeStrict(b, &controlled) == nil, "artifact error is not controlled")
	expectedMessage := controlled.Code
	if controlled.Code == "INVALID_PARAMS" {
		expectedMessage = "Invalid params"
	}
	check.require(controlled.Message == expectedMessage, "artifact controlled error code/message diverge")
	opaque := controlled.Code == "RUN_NOT_FOUND" || controlled.Code == "ARTIFACT_NOT_READY" || controlled.Code == "RUN_STORE_UNAVAILABLE" || controlled.Code == "INTEGRITY_ERROR" || (controlled.Code == "INVALID_PARAMS" && input.Cursor != nil)
	check.require(opaque, "schema-valid ARTIFACT_READ returned an inapplicable controlled error")
	check.unknownIf(opaque, "artifact availability, integrity, compatibility or cursor authenticity is opaque to this request")
}

func o6CheckArtifactEffects(check *o6Check, output ContractOutput) {
	var effects extractionEffects
	b, err := json.Marshal(output.EffectObservation)
	check.require(err == nil && decodeStrict(b, &effects) == nil, "artifact effect observation is not strict")
	// A configured store may perform reads. A read operation may not cause any
	// externally visible write, admission, process or network effect.
	check.require(effects.NetworkCalls == 0, "artifact read performed network I/O")
	check.require(effects.FilesystemWrites == 0 && effects.ProcessSpawns == 0 && effects.ExternalMutations == 0, "artifact read caused an unsafe effect")
	check.require(effects.ResolveCalls == 0 && effects.ReserveCalls == 0 && effects.IDCalls == 0 && effects.CreateCalls == 0, "artifact read touched run admission")
	check.require(effects.StoreInstanceIDs != nil && effects.QueueInstanceIDs != nil && effects.RootInstanceIDs != nil && effects.CursorSecretInstanceIDs != nil && effects.WorkerInstanceIDs != nil, "artifact effect arrays must not be null")

	var publication extractionPublication
	b, err = json.Marshal(output.PublicationObservation)
	check.require(err == nil && decodeStrict(b, &publication) == nil, "artifact publication observation is not strict")
	check.require(publication.ArtifactNames != nil && publication.ArtifactHashes != nil && publication.IndexHashes != nil && publication.LeakedFieldNames != nil && publication.LayerFindings != nil, "artifact publication collections must not be null")
	check.require(len(publication.LeakedFieldNames) == 0, "artifact read exposed forbidden fields")
	for _, hash := range publication.ArtifactHashes {
		check.require(isLowerHex(hash, 64), "publication artifact hash is invalid")
	}
	for _, hash := range publication.IndexHashes {
		check.require(isLowerHex(hash, 64), "publication index hash is invalid")
	}
}

func o6DecodeFixedArtifactDocument(raw []byte) (o6ArtifactDocument, error) {
	var document o6ArtifactDocument
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	token, err := dec.Token()
	if err != nil || token != json.Delim('{') {
		return document, strconv.ErrSyntax
	}
	expected := []string{"schema_version", "artifact", "entries"}
	for i, key := range expected {
		if !dec.More() {
			return document, strconv.ErrSyntax
		}
		name, err := dec.Token()
		if err != nil || name != key {
			return document, strconv.ErrSyntax
		}
		switch i {
		case 0:
			if err := dec.Decode(&document.SchemaVersion); err != nil {
				return document, err
			}
		case 1:
			if err := dec.Decode(&document.Artifact); err != nil {
				return document, err
			}
		case 2:
			if err := dec.Decode(&document.EntriesRaw); err != nil {
				return document, err
			}
		}
	}
	if dec.More() {
		return document, strconv.ErrSyntax
	}
	if token, err = dec.Token(); err != nil || token != json.Delim('}') {
		return document, strconv.ErrSyntax
	}
	if _, err := dec.Token(); err != io.EOF {
		return document, strconv.ErrSyntax
	}

	entriesValue, err := o6ParseUniqueJSON(document.EntriesRaw)
	if err != nil {
		return document, err
	}
	entries, ok := entriesValue.([]any)
	if !ok || entries == nil {
		return document, strconv.ErrSyntax
	}
	document.Entries = entries

	fixed, err := o6MarshalFixedArtifactRoot(document.SchemaVersion, document.Artifact, document.EntriesRaw)
	if err != nil || !bytes.Equal(raw, fixed) {
		return document, strconv.ErrSyntax
	}
	root := map[string]any{"schema_version": document.SchemaVersion, "artifact": document.Artifact, "entries": entriesValue}
	document.SemanticJSON, err = o6MarshalTEKTSemantic(root)
	if err != nil {
		return document, err
	}
	return document, nil
}

func o6MarshalFixedArtifactRoot(schemaVersion, artifact string, entries json.RawMessage) ([]byte, error) {
	schema, err := o6MarshalJSONString(schemaVersion)
	if err != nil {
		return nil, err
	}
	name, err := o6MarshalJSONString(artifact)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(entries)+64)
	out = append(out, `{"schema_version":`...)
	out = append(out, schema...)
	out = append(out, `,"artifact":`...)
	out = append(out, name...)
	out = append(out, `,"entries":`...)
	out = append(out, entries...)
	out = append(out, '}')
	return out, nil
}

func o6ParseUniqueJSON(raw []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	value, err := o6ParseUniqueValue(dec)
	if err != nil {
		return nil, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, strconv.ErrSyntax
	}
	return value, nil
}

func o6ParseUniqueValue(dec *json.Decoder) (any, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch token := token.(type) {
	case json.Delim:
		switch token {
		case '{':
			object := make(map[string]any)
			for dec.More() {
				nameToken, err := dec.Token()
				if err != nil {
					return nil, err
				}
				name, ok := nameToken.(string)
				if !ok {
					return nil, strconv.ErrSyntax
				}
				if _, duplicate := object[name]; duplicate {
					return nil, strconv.ErrSyntax
				}
				value, err := o6ParseUniqueValue(dec)
				if err != nil {
					return nil, err
				}
				object[name] = value
			}
			if close, err := dec.Token(); err != nil || close != json.Delim('}') {
				return nil, strconv.ErrSyntax
			}
			return object, nil
		case '[':
			array := make([]any, 0)
			for dec.More() {
				value, err := o6ParseUniqueValue(dec)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			if close, err := dec.Token(); err != nil || close != json.Delim(']') {
				return nil, strconv.ErrSyntax
			}
			return array, nil
		default:
			return nil, strconv.ErrSyntax
		}
	case nil, bool, string, json.Number:
		return token, nil
	default:
		return nil, strconv.ErrSyntax
	}
}

// o6MarshalTEKTSemantic implements TEKT_CANONICAL_JSON_V1 for already parsed,
// duplicate-free finite JSON. It is deliberately separate from artifact bytes.
func o6MarshalTEKTSemantic(value any) ([]byte, error) {
	var out bytes.Buffer
	if err := o6WriteTEKTValue(&out, value); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func o6WriteTEKTValue(out *bytes.Buffer, value any) error {
	switch value := value.(type) {
	case nil:
		out.WriteString("null")
	case bool:
		if value {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case string:
		encoded, err := o6MarshalJSONString(value)
		if err != nil {
			return err
		}
		out.Write(encoded)
	case json.Number:
		if _, err := strconv.ParseFloat(value.String(), 64); err != nil {
			return err
		}
		out.WriteString(value.String())
	case []any:
		out.WriteByte('[')
		for i, item := range value {
			if i > 0 {
				out.WriteByte(',')
			}
			if err := o6WriteTEKTValue(out, item); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		out.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				out.WriteByte(',')
			}
			encoded, err := o6MarshalJSONString(key)
			if err != nil {
				return err
			}
			out.Write(encoded)
			out.WriteByte(':')
			if err := o6WriteTEKTValue(out, value[key]); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return strconv.ErrSyntax
	}
	return nil
}

func o6MarshalJSONString(value string) ([]byte, error) {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(out.Bytes(), []byte{'\n'}), nil
}

func o6ArtifactEntriesOrdered(artifact string, entries []any) bool {
	previous := ""
	for i, entry := range entries {
		object, ok := entry.(map[string]any)
		if !ok {
			return false
		}
		key, constrained := o6ArtifactEntryOrderKey(artifact, object)
		if constrained && i > 0 && previous >= key {
			return false
		}
		if constrained {
			previous = key
		}
	}
	return true
}

func o6ArtifactEntryOrderKey(artifact string, object map[string]any) (string, bool) {
	stringField := func(name string) string {
		value, _ := object[name].(string)
		return value
	}
	switch artifact {
	case "mechanism_catalog":
		return stringField("id"), stringField("id") != ""
	case "mechanism_gaps":
		source, _ := object["source"].(map[string]any)
		return stringField("code") + "\x00" + o6MapString(source, "repository") + "\x00" + o6MapString(source, "commit") + "\x00" + stringField("rule_id") + "\x00" + stringField("witness_hash"), true
	case "architecture_comparison":
		return stringField("target_property") + "\x00" + stringField("mechanism_id") + "\x00" + stringField("source_run_id"), true
	case "improvement_hypotheses":
		return stringField("id"), stringField("id") != ""
	case "hypothesis_gaps":
		return stringField("hypothesis_id") + "\x00" + stringField("code") + "\x00" + stringField("source_ref"), true
	case "validation_plans":
		index := ""
		if number, ok := object["step_index"].(json.Number); ok {
			index = number.String()
		}
		if len(index) == 0 || len(index) > 20 {
			return "", false
		}
		return stringField("hypothesis_id") + "\x00" + strings.Repeat("0", 20-len(index)) + index, true
	default:
		return "", false
	}
}

func o6MapString(object map[string]any, key string) string {
	if object == nil {
		return ""
	}
	value, _ := object[key].(string)
	return value
}

func o6ArtifactContainsForbiddenFields(value any) bool {
	forbidden := map[string]struct{}{
		"snippet": {}, "content_utf8": {}, "raw_content": {}, "checkout_path": {},
		"idempotency_key": {}, "authorization": {}, "cookie": {}, "private_prompt": {},
	}
	var walk func(any) bool
	walk = func(value any) bool {
		switch value := value.(type) {
		case []any:
			for _, item := range value {
				if walk(item) {
					return true
				}
			}
		case map[string]any:
			for key, item := range value {
				if _, found := forbidden[strings.ToLower(key)]; found || walk(item) {
					return true
				}
			}
		}
		return false
	}
	return walk(value)
}

func o6ArtifactHasRealOpacity(artifact string, entries []any) bool {
	if len(entries) == 0 {
		return false
	}
	for _, entry := range entries {
		if o6JSONHasUncertainty(artifact, entry) {
			return true
		}
	}
	return false
}

func o6JSONHasUncertainty(artifact string, value any) bool {
	switch value := value.(type) {
	case []any:
		for _, item := range value {
			if o6JSONHasUncertainty(artifact, item) {
				return true
			}
		}
	case map[string]any:
		if o15ArtifactHasRelationalOpacity(artifact, value) {
			return true
		}
		for key, item := range value {
			if text, ok := item.(string); ok {
				switch strings.ToLower(key) {
				case "code":
					switch text {
					case "MECHANISM_CLASS_UNKNOWN", "SOURCE_BASIS_UNAVAILABLE":
						return true
					}
				}
			}
			if o6JSONHasUncertainty(artifact, item) {
				return true
			}
		}
	}
	return false
}

// o15GroupHasSemanticWitnessLimit separates the accepted documentary-only
// PARTIAL case from a technical claim whose frozen minimum witness relation
// cannot reach PROVEN. It uses catalog roles, never case identities.
func o15GroupHasSemanticWitnessLimit(group *extractionExpectedGroup) bool {
	if group == nil || group.Status != "PARTIAL" {
		return false
	}
	hasTechnical := false
	partialTechnicalAssessment := false
	for _, witness := range group.Positive {
		binding := extractionFindEvidence(group.Rule, witness.RuleID)
		if binding.Role == "DOCUMENTARY" {
			continue
		}
		hasTechnical = true
		if witness.AssessmentStatus == "PARTIAL" {
			partialTechnicalAssessment = true
		}
	}
	return hasTechnical && partialTechnicalAssessment
}

// o15ArtifactHasRelationalOpacity recognizes only frozen protocol structure:
// full hypothesis relations and controlled gap codes inside their declared
// artifact. Free-form titles, descriptions and expected effects are ignored.
func o15ArtifactHasRelationalOpacity(artifact string, entry map[string]any) bool {
	// O22 gives hypothesis_gaps with canonical hypothesis identities a closed,
	// structural rule. Those entries may become Unknown only through the
	// document-level predicate, which also checks field order, cardinality and
	// publication integrity. This guard prevents predecessor fuzzy predicates
	// from weakening the new code/source_ref relation.
	if artifact == "hypothesis_gaps" && o23ValidHypothesisID(o6MapString(entry, "hypothesis_id")) {
		return false
	}
	if o20PublicE2EHypothesisOpacity(artifact, entry) {
		return true
	}
	if o19ExactPublicGapOpacity(artifact, entry) {
		return true
	}
	if artifact == "hypothesis_gaps" && o19PublicGapCode(o6MapString(entry, "code")) {
		return false
	}
	if o18ArtifactHasPublicOpacity(artifact, entry) {
		return true
	}
	if o14ArtifactHasRelationalOpacity(entry) {
		return true
	}
	if artifact == "mechanism_gaps" || artifact == "hypothesis_gaps" {
		code := strings.ToUpper(o6MapString(entry, "code"))
		if strings.Contains(code, "WITNESS") && strings.Contains(code, "INSUFFICIENT") {
			return true
		}
	}
	// O03 is a predecessor predicate: UNKNOWN becomes opaque only under the
	// incapable hypothesis status. UNKNOWN+PROPOSED remains decidable.
	return o6MapString(entry, "target_state") == "UNKNOWN" && (o6MapString(entry, "status") == "INCONCLUSIVE" || o14BaselineAbsent(entry))
}

func o20PublicE2EHypothesisOpacity(artifact string, entry map[string]any) bool {
	if artifact != "improvement_hypotheses" || o6MapString(entry, "status") != "PROPOSED" {
		return false
	}
	if plan, present := entry["validation_plan"].([]any); present && len(plan) > 0 {
		return false
	}
	if o6MapString(entry, "target_state") == "UNKNOWN" {
		if gaps, ok := entry["target_gap_refs"].([]any); ok && len(gaps) == 0 {
			return true
		}
	}
	mechanisms, mechanismsOK := entry["mechanism_refs"].([]any)
	unknowns, unknownsOK := entry["integration_unknowns"].([]any)
	if mechanismsOK && len(mechanisms) > 1 && unknownsOK && len(unknowns) > 0 {
		return true
	}
	if o6MapString(entry, "license") == "UNRECOGNIZED" && o6MapString(entry, "transfer_mode") == "REFERENCE_ONLY" {
		return true
	}
	search, ok := entry["counterevidence_search"].(map[string]any)
	if !ok || o6MapString(search, "status") != "NONE_FOUND_IN_SCOPE" {
		return false
	}
	scope, scopeOK := search["scope_refs"].([]any)
	results, resultsOK := search["result_refs"].([]any)
	return scopeOK && len(scope) > 0 && resultsOK && len(results) == 0
}

func o19ExactPublicGapOpacity(artifact string, entry map[string]any) bool {
	if artifact != "hypothesis_gaps" {
		return false
	}
	hypothesisID, hypothesisOK := entry["hypothesis_id"].(string)
	code, codeOK := entry["code"].(string)
	sourceRef, sourceOK := entry["source_ref"].(string)
	if !hypothesisOK || !codeOK || !sourceOK || hypothesisID == "" || sourceRef == "" {
		return false
	}
	// Canonical O22 identities are governed exclusively by the exact O23
	// document rule. Predecessor public projections used non-canonical neutral
	// identities and retain their historical classification here.
	if o23ValidHypothesisID(hypothesisID) {
		return false
	}
	return o19PublicGapCode(code)
}

func o23OpaqueArtifactRule(document o6ArtifactDocument, result o6ArtifactResult, publication PublicationObservation) bool {
	if document.SchemaVersion != "1" || document.Artifact != "hypothesis_gaps" || len(document.Entries) != 1 {
		return false
	}
	entry, ok := document.Entries[0].(map[string]any)
	if !ok || len(entry) != 3 {
		return false
	}
	hypothesisID := o6MapString(entry, "hypothesis_id")
	code := o6MapString(entry, "code")
	sourceRef := o6MapString(entry, "source_ref")
	if !o23ValidHypothesisID(hypothesisID) || !o23ControlledPair(code, sourceRef) {
		return false
	}
	idJSON, errID := o6MarshalJSONString(hypothesisID)
	codeJSON, errCode := o6MarshalJSONString(code)
	sourceJSON, errSource := o6MarshalJSONString(sourceRef)
	if errID != nil || errCode != nil || errSource != nil {
		return false
	}
	expected := make([]byte, 0, len(hypothesisID)+len(code)+len(sourceRef)+48)
	expected = append(expected, `[{"hypothesis_id":`...)
	expected = append(expected, idJSON...)
	expected = append(expected, `,"code":`...)
	expected = append(expected, codeJSON...)
	expected = append(expected, `,"source_ref":`...)
	expected = append(expected, sourceJSON...)
	expected = append(expected, '}', ']')
	if !bytes.Equal(document.EntriesRaw, expected) {
		return false
	}
	if publication.ArtifactHashes == nil || publication.ArtifactHashes["hypothesis_gaps"] != result.SHA256 {
		return false
	}
	for _, name := range publication.ArtifactNames {
		if name == "hypothesis_gaps" {
			return true
		}
	}
	return false
}

func o23ValidHypothesisID(value string) bool {
	return strings.HasPrefix(value, "hyp_") && isLowerHex(strings.TrimPrefix(value, "hyp_"), 64)
}

func o23ControlledCode(code string) bool {
	switch code {
	case "COMPOSITION_COMPATIBILITY_UNMEASURED",
		"LICENSE_UNRECOGNIZED",
		"COUNTEREVIDENCE_BUDGET_EXHAUSTED",
		"IMMUTABLE_IDENTITY_AMBIGUOUS",
		"PROVIDER_CAPTURE_ABSENT",
		"UTILITY_COST_BASELINE_ABSENT",
		"SECURITY_BASELINE_ABSENT":
		return true
	default:
		return false
	}
}

func o23ControlledPair(code, sourceRef string) bool {
	switch code {
	case "COMPOSITION_COMPATIBILITY_UNMEASURED":
		return sourceRef == "hypothesis:composition"
	case "LICENSE_UNRECOGNIZED":
		return sourceRef == "hypothesis:license"
	case "COUNTEREVIDENCE_BUDGET_EXHAUSTED":
		return sourceRef == "counterevidence:search"
	case "IMMUTABLE_IDENTITY_AMBIGUOUS":
		return sourceRef == "source:identity"
	case "PROVIDER_CAPTURE_ABSENT":
		return sourceRef == "provider:capture"
	case "UTILITY_COST_BASELINE_ABSENT":
		return sourceRef == "baseline:utility_cost"
	case "SECURITY_BASELINE_ABSENT":
		return sourceRef == "baseline:security"
	default:
		return false
	}
}

func o19PublicGapCode(code string) bool {
	switch code {
	case "COMPOSITION_COMPATIBILITY_UNMEASURED",
		"LICENSE_NOT_RECOGNIZED",
		"COUNTEREVIDENCE_BUDGET_EXHAUSTED",
		"IMMUTABLE_IDENTITY_AMBIGUOUS",
		"PROVIDER_CAPTURE_UNAVAILABLE",
		"UTILITY_COST_BASELINE_ABSENT",
		"SECURITY_BASELINE_ABSENT":
		return true
	default:
		return false
	}
}

func o18ArtifactHasPublicOpacity(artifact string, entry map[string]any) bool {
	status := o6MapString(entry, "status")
	plan, planOK := entry["validation_plan"].([]any)
	mechanisms, mechanismsOK := entry["mechanism_refs"].([]any)
	unknowns, unknownsOK := entry["integration_unknowns"].([]any)
	if artifact == "improvement_hypotheses" && mechanismsOK && len(mechanisms) > 1 && unknownsOK && len(unknowns) > 0 && status == "INCONCLUSIVE" && planOK && len(plan) > 0 {
		return true
	}
	search, searchOK := entry["counterevidence_search"].(map[string]any)
	if artifact == "improvement_hypotheses" && searchOK && o6MapString(search, "status") == "BUDGET_EXHAUSTED" && status == "INCONCLUSIVE" && planOK && len(plan) > 0 {
		return true
	}
	if artifact != "mechanism_gaps" && artifact != "hypothesis_gaps" {
		return false
	}
	code := strings.ToUpper(o6MapString(entry, "code"))
	hasAny := func(tokens ...string) bool {
		for _, token := range tokens {
			if strings.Contains(code, token) {
				return true
			}
		}
		return false
	}
	if artifact == "hypothesis_gaps" && code == "LICENSE_NOT_RECOGNIZED" {
		return true
	}
	if (hasAny("COMPOSITION", "COMPATIBILITY") && hasAny("UNMEASURED", "UNKNOWN", "INSUFFICIENT")) ||
		(hasAny("COUNTEREVIDENCE", "BUDGET") && strings.Contains(code, "EXHAUSTED")) ||
		(hasAny("IDENTITY", "BINDING") && hasAny("MISSING", "AMBIGUOUS", "UNBOUND", "UNRESOLVED")) ||
		(hasAny("PROVIDER", "CAPTURE") && hasAny("UNKNOWN", "UNAVAILABLE", "NOT_EVIDENCED", "MISSING")) {
		return true
	}
	if artifact != "hypothesis_gaps" {
		return false
	}
	baselineMissing := strings.Contains(code, "BASELINE") && hasAny("MISSING", "ABSENT", "WITHOUT")
	return baselineMissing && (hasAny("UTILITY", "COST") || strings.Contains(code, "SECURITY"))
}

// o17ArtifactHasTypedOpacity is the consolidated R4/O8 decision table. It
// consumes only enum-valued fields and controlled gap codes in their declared
// artifact; arbitrary human text is deliberately outside the predicate.
func o17ArtifactHasTypedOpacity(artifact string, entry map[string]any) bool {
	status := o6MapString(entry, "status")
	if status == "INCONCLUSIVE" && o6MapString(entry, "target_state") == "UNKNOWN" {
		return true
	}
	mechanisms, mechanismsOK := entry["mechanism_refs"].([]any)
	unknowns, unknownsOK := entry["integration_unknowns"].([]any)
	if status == "INCONCLUSIVE" && mechanismsOK && len(mechanisms) >= 2 && unknownsOK && len(unknowns) > 0 {
		return true
	}
	if o6MapString(entry, "license") == "UNRECOGNIZED" && (status == "INCONCLUSIVE" || o6MapString(entry, "transfer_mode") == "UNKNOWN") {
		return true
	}
	if o14CounterEvidenceBudgetExhausted(entry) {
		return true
	}
	if artifact != "mechanism_gaps" && artifact != "hypothesis_gaps" {
		return false
	}
	code := strings.ToUpper(o6MapString(entry, "code"))
	if code == "" {
		return false
	}
	if artifact == "hypothesis_gaps" && code == "LICENSE_NOT_RECOGNIZED" {
		return true
	}
	limited := strings.Contains(code, "MISSING") || strings.Contains(code, "ABSENT") || strings.Contains(code, "UNKNOWN") || strings.Contains(code, "UNAVAILABLE") || strings.Contains(code, "UNMEASURED") || strings.Contains(code, "INSUFFICIENT") || strings.Contains(code, "AMBIGUOUS") || strings.Contains(code, "UNBOUND") || strings.Contains(code, "UNRESOLVED") || strings.Contains(code, "WITHOUT") || strings.Contains(code, "NOT_EVIDENCED") || strings.Contains(code, "NOT_MEASURED") || strings.Contains(code, "NOT_ESTABLISHED") || strings.Contains(code, "NOT_AVAILABLE")
	if strings.Contains(code, "BUDGET") && strings.Contains(code, "EXHAUSTED") {
		return true
	}
	if !limited {
		return false
	}
	return strings.Contains(code, "BASELINE") || strings.Contains(code, "COMPATIBILITY") || strings.Contains(code, "WITNESS") ||
		strings.Contains(code, "IDENTITY") || strings.Contains(code, "IMMUTABLE") || strings.Contains(code, "BINDING") ||
		strings.Contains(code, "PROVIDER") || strings.Contains(code, "CAPTURE") ||
		strings.Contains(code, "UTILITY") || strings.Contains(code, "SECURITY")
}

func o14ArtifactHasRelationalOpacity(entry map[string]any) bool {
	status := o6MapString(entry, "status")
	transfer := o6MapString(entry, "transfer_mode")
	incapable := status == "INCONCLUSIVE" || o6MapString(entry, "status_ceiling") == "INCONCLUSIVE" || transfer == "UNKNOWN"
	context := o14HasEpistemicContext(entry)

	if o6MapString(entry, "target_state") == "UNKNOWN" && o14BaselineAbsent(entry) {
		return true
	}
	if o14CompositionUnmeasured(entry) && incapable && context {
		return true
	}
	if o6MapString(entry, "license") == "UNRECOGNIZED" && incapable && context {
		return true
	}
	if o14CounterEvidenceBudgetExhausted(entry) && status == "INCONCLUSIVE" && context {
		return true
	}
	if o14ImmutableIdentityOpaque(entry) && incapable {
		return true
	}
	if o14ProviderAvailabilityOpaque(entry) && context {
		return true
	}
	return o14UtilityOrSecurityWithoutBaseline(entry) && o14BaselineAbsent(entry) && context
}

func o14BaselineAbsent(entry map[string]any) bool {
	if absent, ok := entry["baseline_absent"].(bool); ok && absent {
		return true
	}
	if present, ok := entry["baseline_present"].(bool); ok && !present {
		return true
	}
	if baseline, exists := entry["baseline_ref"]; exists && baseline == nil {
		return true
	}
	if baseline, exists := entry["baseline_value"]; exists && baseline == nil {
		return true
	}
	for _, name := range []string{"baseline_state", "baseline_status"} {
		switch o6MapString(entry, name) {
		case "ABSENT", "MISSING", "UNKNOWN", "UNAVAILABLE", "NOT_EVIDENCED":
			return true
		}
	}
	for _, name := range []string{"target_gap_refs", "baseline_gap_refs"} {
		if refs, ok := entry[name].([]any); ok && len(refs) > 0 {
			return true
		}
	}
	return false
}

func o14HasEpistemicContext(entry map[string]any) bool {
	for _, name := range []string{"validation_plan", "validation_steps", "missing_evidence", "target_gap_refs", "baseline_gap_refs"} {
		switch value := entry[name].(type) {
		case []any:
			if len(value) > 0 {
				return true
			}
		case string:
			if strings.TrimSpace(value) != "" {
				return true
			}
		}
	}
	return o6MapString(entry, "target_property") != "" || o6MapString(entry, "source_ref") != "" || o6MapString(entry, "source_run_id") != ""
}

func o14CompositionUnmeasured(entry map[string]any) bool {
	if measured, ok := entry["composition_compatibility_measured"].(bool); ok && !measured {
		return true
	}
	mechanisms, mechanismsOK := entry["mechanism_refs"].([]any)
	unknowns, unknownsOK := entry["integration_unknowns"].([]any)
	return mechanismsOK && len(mechanisms) >= 2 && unknownsOK && len(unknowns) > 0
}

func o14CounterEvidenceBudgetExhausted(entry map[string]any) bool {
	if exhausted, ok := entry["counterevidence_budget_exhausted"].(bool); ok && exhausted {
		return true
	}
	if o6MapString(entry, "counterevidence_search_status") == "BUDGET_EXHAUSTED" {
		return true
	}
	for _, name := range []string{"counter_evidence_search", "counterevidence_search"} {
		if search, ok := entry[name].(map[string]any); ok && o6MapString(search, "status") == "BUDGET_EXHAUSTED" {
			return true
		}
	}
	return false
}

func o14ImmutableIdentityOpaque(entry map[string]any) bool {
	code := strings.ToUpper(o6MapString(entry, "code"))
	if code == "UNBOUND_EVALUATED_SOURCE" || code == "AMBIGUOUS_BUNDLE_BINDING" {
		return true
	}
	return strings.Contains(code, "IDENTITY") && (strings.Contains(code, "MISSING") || strings.Contains(code, "AMBIGUOUS"))
}

func o14ProviderAvailabilityOpaque(entry map[string]any) bool {
	availability := o6MapString(entry, "provider_availability")
	unavailable := availability == "UNKNOWN" || availability == "NOT_EVIDENCED"
	if evidenced, ok := entry["provider_availability_evidenced"].(bool); ok && !evidenced {
		unavailable = true
	}
	if !unavailable {
		return false
	}
	if frozen, ok := entry["frozen_capture_present"].(bool); ok {
		return !frozen
	}
	if capture, exists := entry["frozen_capture_ref"]; exists {
		return capture == nil || o6MapString(entry, "frozen_capture_ref") == ""
	}
	return true
}

func o14UtilityOrSecurityWithoutBaseline(entry map[string]any) bool {
	kind := strings.ToUpper(o6MapString(entry, "effect_kind"))
	if kind == "UTILITY_COST" || kind == "SECURITY_EFFECT" {
		return true
	}
	for _, key := range []string{"utility_cost", "security_effect", "utility_effect"} {
		if _, exists := entry[key]; exists {
			return true
		}
	}
	return false
}

func o11TargetBaselineUnresolved(entry map[string]any) bool {
	for _, name := range []string{"target_gap_refs", "baseline_gap_refs", "missing_evidence"} {
		if refs, ok := entry[name].([]any); ok && len(refs) > 0 {
			return true
		}
	}
	if present, ok := entry["baseline_present"].(bool); ok && !present {
		return true
	}
	if baseline, exists := entry["baseline_ref"]; exists && baseline == nil {
		return true
	}
	return false
}

func o6AllowedArtifact(artifact string) bool {
	switch artifact {
	case "request", "recommendation", "decision", "alternatives", "failures", "provenance",
		"connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps",
		"mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses",
		"hypothesis_gaps", "validation_plans", "synthesis_provenance":
		return true
	default:
		return false
	}
}
