/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/ports.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

const (
	mechanismRulesDigest = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"
	receipt0025Digest    = "bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d"
	receipt0026Digest    = "a04d6b07a939c1ce5fb7fb184e5cb56cc3cc34ed0b2b54c8087415a18a4a99ed"
	receipt0027Digest    = "acf05df728b05ccdc5bc15f78a18e49f19e743b2016ab1141e553ca2417c7e43"
)

var synthesisArtifactOrder = []string{
	"architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance",
}

var toolOrder = []string{
	"bateia_search", "bateia_benchmark", "bateia_recommend", "bateia_start_recommendation", "bateia_get_run",
	"bateia_get_artifact", "bateia_cancel_run", "bateia_list_connectors", "bateia_start_market_discovery", "bateia_start_improvement_synthesis",
}

type ideaInput struct {
	Problem        string   `json:"problem"`
	Users          []string `json:"users"`
	Capabilities   []string `json:"capabilities"`
	Constraints    []string `json:"constraints"`
	NonGoals       []string `json:"non_goals"`
	PreferredStack []string `json:"preferred_stack"`
}

type sourceArtifactInput struct {
	RunID, RunKind, Artifact, SHA256 string
}

func (source *sourceArtifactInput) UnmarshalJSON(data []byte) error {
	type wire struct {
		RunID    string `json:"run_id"`
		RunKind  string `json:"run_kind"`
		Artifact string `json:"artifact"`
		SHA256   string `json:"sha256"`
	}
	var value wire
	if err := decodeStrict(data, &value); err != nil {
		return err
	}
	source.RunID, source.RunKind, source.Artifact, source.SHA256 = value.RunID, value.RunKind, value.Artifact, value.SHA256
	return nil
}

type proposalInput struct {
	ProposalID             string   `json:"proposal_id"`
	Title                  string   `json:"title"`
	TargetProperty         string   `json:"target_property"`
	ExpectedEffect         string   `json:"expected_effect"`
	MechanismRefs          []string `json:"mechanism_refs"`
	ProductObservationRefs []string `json:"product_observation_refs"`
	Assumptions            []string `json:"assumptions"`
	Derivation             string   `json:"derivation"`
	Rationale              string   `json:"rationale"`
}

type synthesisPolicyInput struct {
	MaxHypotheses                int  `json:"max_hypotheses"`
	MaxMechanismsPerHypothesis   int  `json:"max_mechanisms_per_hypothesis"`
	RequireCounterevidenceSearch bool `json:"require_counterevidence_search"`
}

type synthesisStartInput struct {
	IdempotencyKey string                `json:"idempotency_key"`
	Idea           ideaInput             `json:"idea"`
	Sources        []sourceArtifactInput `json:"source_artifacts"`
	Proposals      []proposalInput       `json:"proposals"`
	Policy         synthesisPolicyInput  `json:"policy"`
}

type runInput struct {
	RunID string `json:"run_id"`
}
type cancelInput struct {
	RunID           string `json:"run_id"`
	ExpectedVersion uint64 `json:"expected_version"`
}
type artifactInput struct {
	RunID    string  `json:"run_id"`
	Artifact string  `json:"artifact"`
	Cursor   *string `json:"cursor"`
}
type startReceipt struct {
	SchemaVersion string `json:"schema_version"`
	RunID         string `json:"run_id"`
	RunKind       string `json:"run_kind"`
	Lifecycle     string `json:"lifecycle"`
	Version       uint64 `json:"version"`
}
type progressView struct {
	Planned   int    `json:"planned"`
	Completed int    `json:"completed"`
	Stage     string `json:"stage"`
}
type runView struct {
	RunID                 string          `json:"run_id"`
	RunKind               string          `json:"run_kind"`
	Lifecycle             string          `json:"lifecycle"`
	Version               uint64          `json:"version"`
	RecommendationStatus  json.RawMessage `json:"recommendation_status"`
	MarketDiscoveryStatus json.RawMessage `json:"market_discovery_status"`
	ImprovementStatus     json.RawMessage `json:"improvement_status"`
	Progress              progressView    `json:"progress"`
}
type artifactResult struct {
	Artifact      string  `json:"artifact"`
	SHA256        string  `json:"sha256"`
	Size          int     `json:"size"`
	CanonicalJSON string  `json:"canonical_json"`
	NextCursor    *string `json:"next_cursor"`
}
type toolResult struct {
	Tools []toolDescriptor `json:"tools"`
}
type toolDescriptor struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}
type replayInput struct {
	CatalogSHA256   string            `json:"catalog_sha256"`
	CaptureSHA256   string            `json:"capture_sha256"`
	ArtifactHashes  map[string]string `json:"artifact_hashes"`
	CanonicalInputs []string          `json:"canonical_inputs"`
}
type replayResult struct {
	ArtifactHashes map[string]string `json:"artifact_hashes"`
	IndexHashes    map[string]string `json:"index_hashes"`
	NetworkCalls   int               `json:"network_calls"`
}
type refinementInput struct {
	ReceiptHashes struct {
		R0025 string `json:"0025"`
		R0026 string `json:"0026"`
		R0027 string `json:"0027"`
	} `json:"receipt_hashes"`
	RuntimeInstanceIDs struct {
		Store        string `json:"store"`
		Queue        string `json:"queue"`
		Root         string `json:"root"`
		CursorSecret string `json:"cursor_secret"`
		Worker       string `json:"worker"`
	} `json:"runtime_instance_ids"`
	LayerFindings []string `json:"layer_findings"`
}
type refinementResult struct {
	ReceiptsVerified struct {
		R0025 bool `json:"0025"`
		R0026 bool `json:"0026"`
		R0027 bool `json:"0027"`
	} `json:"receipts_verified"`
	RuntimeUnique bool     `json:"runtime_unique"`
	LayerFindings []string `json:"layer_findings"`
}

func evaluateSynthesisStart(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input synthesisStartInput
	if decodeStrict(request, &input) != nil {
		state.fail("synthesis start request is not strict protocol JSON")
		return
	}
	if !validSynthesisStartInput(input) {
		state.require(isError(output, "INVALID_PARAMS"), "locally invalid synthesis request must fail with Invalid params")
		state.require(rawIsNull(output.ReturnValue), "invalid synthesis request must not return a receipt")
		state.require(noAdmissionCalls(output), "invalid synthesis request reached source resolution or admission")
		state.require(noPublication(output), "invalid synthesis request published artifacts")
		return
	}
	state.require(output.EffectObservation.ResolveCalls >= 0 && output.EffectObservation.ResolveCalls <= len(input.Sources), "source resolution count is outside the declared source set")
	state.require(output.EffectObservation.ReserveCalls <= 1 && output.EffectObservation.IDCalls <= 1 && output.EffectObservation.CreateCalls <= 1, "start operation repeated an admission side effect")
	state.require(output.EffectObservation.CreateCalls <= output.EffectObservation.IDCalls, "run creation occurred without run-ID generation")
	state.require(output.EffectObservation.IDCalls <= output.EffectObservation.ReserveCalls, "run-ID generation occurred before capacity reservation")
	state.require(noPublication(output), "start operation must not publish synthesis artifacts inline")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed synthesis start returned a receipt")
		state.require(startErrorAllowed(output.Error.Code), "synthesis start returned an error outside its contract")
		switch output.Error.Code {
		case "RUN_NOT_FOUND", "SOURCE_BASIS_UNAVAILABLE", "IDEMPOTENCY_CONFLICT", "INTEGRITY_ERROR", "INVALID_PARAMS":
			state.require(output.EffectObservation.ReserveCalls == 0 && output.EffectObservation.IDCalls == 0 && output.EffectObservation.CreateCalls == 0, "pre-admission failure mutated admission state")
		case "QUEUE_FULL", "QUOTA_EXCEEDED":
			state.require(output.EffectObservation.IDCalls == 0 && output.EffectObservation.CreateCalls == 0, "rejected admission created a run")
		}
		if output.Error.Code == "SOURCE_BASIS_UNAVAILABLE" {
			state.unknown("declared immutable source basis is unavailable without reread")
		}
		return
	}
	var receipt startReceipt
	state.require(decodeReturn(output.ReturnValue, &receipt) == nil, "successful synthesis start did not return a strict receipt")
	state.require(receipt.SchemaVersion == "1", "synthesis start receipt schema version changed")
	state.require(isRunID(receipt.RunID), "synthesis start receipt has an invalid run_id")
	state.require(receipt.RunKind == "IMPROVEMENT_SYNTHESIS", "synthesis start receipt has the wrong run kind")
	state.require(receipt.Lifecycle == "QUEUED", "synthesis start must return the queued state")
	state.require(receipt.Version >= 1, "synthesis start receipt has a zero version")
	state.require(output.EffectObservation.ResolveCalls == len(input.Sources), "successful synthesis start did not resolve every declared source")
	state.require(output.EffectObservation.CreateCalls == 0 || (output.EffectObservation.ReserveCalls == 1 && output.EffectObservation.IDCalls == 1), "fresh run creation did not preserve admission order")
}

func evaluateSynthesisExecute(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input runInput
	if decodeStrict(request, &input) != nil || !isRunID(input.RunID) {
		state.fail("synthesis execute request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "worker execution used start/admission operations")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed synthesis execution returned a value")
		state.require(errorAllowed(output.Error.Code, "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE"), "synthesis execution returned an error outside its contract")
		state.require(noPublication(output), "failed synthesis execution exposed a partial publication")
		return
	}
	state.require(rawIsNull(output.ReturnValue), "synthesis executor must return no value")
	state.require(equalStrings(output.PublicationObservation.ArtifactNames, synthesisArtifactOrder), "synthesis artifacts were not published in exact contract order")
	state.require(exactHashKeys(output.PublicationObservation.ArtifactHashes, synthesisArtifactOrder), "synthesis publication artifact hashes are incomplete or extraneous")
	state.require(exactHashKeys(output.PublicationObservation.IndexHashes, synthesisArtifactOrder), "synthesis publication index hashes are incomplete or extraneous")
}

func evaluateRunRead(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input runInput
	if decodeStrict(request, &input) != nil || !isRunID(input.RunID) {
		state.fail("run read request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "run read performed admission work")
	state.require(noPublication(output), "run read mutated publication state")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed run read returned a run")
		state.require(errorAllowed(output.Error.Code, "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE"), "run read returned an error outside its contract")
		return
	}
	var view runView
	state.require(decodeReturn(output.ReturnValue, &view) == nil, "run read did not return a strict run view")
	evaluateRunView(state, input.RunID, view)
}

func evaluateArtifactRead(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input artifactInput
	if decodeStrict(request, &input) != nil || !isRunID(input.RunID) || !artifactGloballyAllowed(input.Artifact) {
		state.fail("artifact read request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "artifact read performed admission work")
	state.require(noPublication(output), "artifact read mutated publication state")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed artifact read returned an artifact")
		state.require(errorAllowed(output.Error.Code, "INVALID_PARAMS", "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE", "ARTIFACT_NOT_READY"), "artifact read returned an error outside its contract")
		return
	}
	var result artifactResult
	if decodeReturn(output.ReturnValue, &result) != nil {
		state.fail("artifact read did not return a strict artifact result")
		return
	}
	state.require(result.Artifact == input.Artifact, "artifact response does not match the requested artifact")
	state.require(isLowerHex(result.SHA256, 64), "artifact response hash is invalid")
	contents, err := base64.StdEncoding.DecodeString(result.CanonicalJSON)
	state.require(err == nil, "canonical_json is not canonical base64")
	if err != nil {
		return
	}
	state.require(base64.StdEncoding.EncodeToString(contents) == result.CanonicalJSON, "canonical_json uses a non-canonical base64 encoding")
	state.require(result.Size == len(contents), "artifact size does not match canonical bytes")
	state.require(digestBytes(contents) == result.SHA256, "artifact hash does not match canonical bytes")
	var compact bytes.Buffer
	compactErr := json.Compact(&compact, contents)
	state.require(compactErr == nil && bytes.Equal(compact.Bytes(), contents), "artifact document is not compact canonical JSON")
	if compactErr != nil {
		return
	}
	var document any
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.UseNumber()
	if decoder.Decode(&document) != nil {
		state.fail("artifact document cannot be decoded")
		return
	}
	object, ok := document.(map[string]any)
	state.require(ok, "artifact document must be an object")
	if !ok {
		return
	}
	state.require(exactObjectKeys(object, "schema_version", "artifact", "entries"), "artifact document top-level shape changed")
	state.require(object["schema_version"] == "1", "artifact document schema version changed")
	state.require(object["artifact"] == input.Artifact, "artifact document name does not match its index")
	_, entriesAreArray := object["entries"].([]any)
	state.require(entriesAreArray, "artifact entries must be an array")
	state.require(!containsUnsafePublishedField(document), "artifact contains a forbidden raw or control field")
	if containsContractualUnknown(document) {
		state.unknown("artifact preserves an explicit contractual unknown")
	}
}

func evaluateRunCancel(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input cancelInput
	if decodeStrict(request, &input) != nil || !isRunID(input.RunID) || input.ExpectedVersion == 0 {
		state.fail("run cancel request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "run cancellation performed admission work")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed cancellation returned a run")
		state.require(errorAllowed(output.Error.Code, "INVALID_PARAMS", "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE"), "run cancellation returned an error outside its contract")
		state.require(noPublication(output), "failed cancellation published artifacts")
		return
	}
	var view runView
	if decodeReturn(output.ReturnValue, &view) != nil {
		state.fail("run cancellation did not return a strict run view")
		return
	}
	evaluateRunView(state, input.RunID, view)
	state.require(view.Version >= input.ExpectedVersion, "cancel response regressed the expected version")
	state.require(errorAllowed(view.Lifecycle, "CANCELLING", "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED"), "cancel response did not enter cancellation or a terminal race")
	if view.Lifecycle != "SUCCEEDED" {
		state.require(noPublication(output), "cancelled or degraded run published a new artifact set")
	}
}

func evaluateToolsList(state *evaluation, request json.RawMessage, output ContractOutput) {
	var empty map[string]json.RawMessage
	if decodeStrict(request, &empty) != nil || len(empty) != 0 {
		state.fail("tools/list request must be the empty object")
		return
	}
	state.require(output.Error == nil, "tools/list returned an error")
	state.require(noAdmissionCalls(output), "tools/list performed runtime admission work")
	state.require(noPublication(output), "tools/list mutated publication state")
	var result toolResult
	if decodeReturn(output.ReturnValue, &result) != nil {
		state.fail("tools/list did not return a strict tool registry")
		return
	}
	state.require(len(result.Tools) == len(toolOrder), "tools/list must expose exactly ten tools")
	if len(result.Tools) != len(toolOrder) {
		return
	}
	for index, expectedName := range toolOrder {
		state.require(result.Tools[index].Name == expectedName, "tool registry order or name changed")
		state.require(len(result.Tools[index].InputSchema) != 0 && !rawIsNull(result.Tools[index].InputSchema), "tool input schema is missing")
	}
	state.require(result.Tools[9].Description == "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis", "tool 10 description changed")
	expected, expectedErr := canonicalJSON([]byte(tool10Schema))
	actual, actualErr := canonicalJSON(result.Tools[9].InputSchema)
	state.require(expectedErr == nil && actualErr == nil && bytes.Equal(expected, actual), "tool 10 input schema is not the exact R3 schema")
	state.require(toolSixHasExactArtifactEnum(result.Tools[5].InputSchema), "artifact getter does not expose the exact additive artifact enum")
}

func evaluateReplay(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input replayInput
	if decodeStrict(request, &input) != nil {
		state.fail("replay request is not strict protocol JSON")
		return
	}
	validEncoding := len(input.ArtifactHashes) > 0 && len(input.CanonicalInputs) > 0 && isLowerHex(input.CaptureSHA256, 64)
	for name, hash := range input.ArtifactHashes {
		validEncoding = validEncoding && name != "" && isLowerHex(hash, 64)
	}
	for _, encoded := range input.CanonicalInputs {
		contents, err := base64.StdEncoding.DecodeString(encoded)
		var compact bytes.Buffer
		compactErr := json.Compact(&compact, contents)
		validEncoding = validEncoding && err == nil && compactErr == nil && bytes.Equal(contents, compact.Bytes()) && base64.StdEncoding.EncodeToString(contents) == encoded
	}
	if input.CatalogSHA256 != mechanismRulesDigest || !validEncoding {
		state.require(isError(output, "INTEGRITY_ERROR"), "replay with divergent frozen input must fail with INTEGRITY_ERROR")
		state.require(rawIsNull(output.ReturnValue), "failed replay returned hashes")
		state.require(noAdmissionCalls(output) && noPublication(output), "failed replay produced runtime or publication effects")
		return
	}
	state.require(noAdmissionCalls(output), "replay performed runtime admission work")
	state.require(noPublication(output), "replay mutated published state")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed replay returned hashes")
		state.require(output.Error.Code == "INTEGRITY_ERROR", "offline replay may fail only on frozen-input integrity")
		return
	}
	var result replayResult
	if decodeReturn(output.ReturnValue, &result) != nil {
		state.fail("replay did not return strict hash projections")
		return
	}
	state.require(result.NetworkCalls == 0, "replay reported a network call")
	state.require(len(result.ArtifactHashes) > 0 && len(result.IndexHashes) > 0, "replay did not project artifact and index hashes")
	state.require(hashMapValid(result.ArtifactHashes) && hashMapValid(result.IndexHashes), "replay projected an invalid hash map")
	state.require(reflect.DeepEqual(result.ArtifactHashes, input.ArtifactHashes), "replay artifact bytes differ from the frozen expected hashes")
}

func evaluateRefinement(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input refinementInput
	if decodeStrict(request, &input) != nil {
		state.fail("refinement request is not strict protocol JSON")
		return
	}
	state.require(noAdmissionCalls(output), "refinement audit performed admission work")
	state.require(noPublication(output), "refinement audit mutated publication state")
	state.require(output.Error == nil, "refinement audit returned an error instead of evidence")
	var result refinementResult
	if decodeReturn(output.ReturnValue, &result) != nil {
		state.fail("refinement audit did not return a strict result")
		return
	}
	state.require(result.ReceiptsVerified.R0025 == (input.ReceiptHashes.R0025 == receipt0025Digest), "0025 receipt verification does not follow the frozen predecessor hash")
	state.require(result.ReceiptsVerified.R0026 == (input.ReceiptHashes.R0026 == receipt0026Digest), "0026 receipt verification does not follow the frozen predecessor hash")
	state.require(result.ReceiptsVerified.R0027 == (input.ReceiptHashes.R0027 == receipt0027Digest), "0027 receipt verification does not follow the frozen predecessor hash")
	requiredLayers := []string{"L1_PURE", "L2_NO_L3_OR_L4", "L3_NO_L2_OR_L4", "L4_WIRING_ONLY"}
	state.require(equalStringSets(input.LayerFindings, requiredLayers), "refinement request does not cover every required layer boundary")
	state.require(equalStringSets(result.LayerFindings, input.LayerFindings), "refinement result changed the audited layer findings")
	state.require(result.RuntimeUnique, "refinement did not attest the unique runtime")
	state.require(observedIdentityMatches(output.EffectObservation.StoreInstanceIDs, input.RuntimeInstanceIDs.Store), "observed store identity differs from the audited runtime")
	state.require(observedIdentityMatches(output.EffectObservation.QueueInstanceIDs, input.RuntimeInstanceIDs.Queue), "observed queue identity differs from the audited runtime")
	state.require(observedIdentityMatches(output.EffectObservation.RootInstanceIDs, input.RuntimeInstanceIDs.Root), "observed root identity differs from the audited runtime")
	state.require(observedIdentityMatches(output.EffectObservation.CursorSecretInstanceIDs, input.RuntimeInstanceIDs.CursorSecret), "observed cursor-secret identity differs from the audited runtime")
	state.require(observedIdentityMatches(output.EffectObservation.WorkerInstanceIDs, input.RuntimeInstanceIDs.Worker), "observed worker identity differs from the audited runtime")
}

func evaluateRunView(state *evaluation, requestedRunID string, view runView) {
	state.require(view.RunID == requestedRunID && isRunID(view.RunID), "run view identity does not match the request")
	state.require(errorAllowed(view.RunKind, "RECOMMENDATION", "MARKET_DISCOVERY", "IMPROVEMENT_SYNTHESIS"), "run view has an unknown run kind")
	state.require(errorAllowed(view.Lifecycle, "QUEUED", "RUNNING", "CANCELLING", "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED"), "run view has an unknown lifecycle")
	state.require(view.Version >= 1, "run view has a zero version")
	state.require(view.Progress.Planned >= 0 && view.Progress.Completed >= 0 && view.Progress.Completed <= view.Progress.Planned, "run progress is not bounded or monotonic")
	state.require(view.Progress.Stage != "", "run progress stage is empty")
	terminal := errorAllowed(view.Lifecycle, "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED")
	statuses := []json.RawMessage{view.RecommendationStatus, view.MarketDiscoveryStatus, view.ImprovementStatus}
	corresponding := 0
	switch view.RunKind {
	case "RECOMMENDATION":
		corresponding = 0
	case "MARKET_DISCOVERY":
		corresponding = 1
	case "IMPROVEMENT_SYNTHESIS":
		corresponding = 2
		state.require(view.Progress.Planned == 5, "improvement synthesis must retain five planned stages")
	}
	for index, status := range statuses {
		if terminal && index == corresponding {
			state.require(!rawIsNull(status), "terminal run is missing its kind-specific status")
		} else {
			state.require(rawIsNull(status), "run exposed a status for the wrong kind or before terminal")
		}
	}
	if view.Lifecycle == "QUEUED" {
		state.require(view.Progress.Completed == 0 && view.Progress.Stage == "queued", "queued progress changed")
	}
	if terminal {
		state.require(view.Progress.Stage == "terminal", "terminal run must expose terminal progress stage")
	}
	if view.RunKind == "IMPROVEMENT_SYNTHESIS" && terminal {
		var status string
		if decodeStrict(view.ImprovementStatus, &status) != nil || !errorAllowed(status, "COMPLETE", "PARTIAL", "INCONCLUSIVE") {
			state.fail("improvement_status is outside the closed enum")
		} else if status == "INCONCLUSIVE" {
			state.unknown("improvement synthesis concluded with explicit epistemic insufficiency")
		}
		if view.Lifecycle == "SUCCEEDED" {
			state.require(view.Progress.Completed == 5, "successful synthesis did not complete all five stages")
		}
	}
}

func validSynthesisStartInput(input synthesisStartInput) bool {
	if !validIdempotencyKey(input.IdempotencyKey) || len(input.Sources) < 1 || len(input.Sources) > 12 || len(input.Proposals) > 20 {
		return false
	}
	if input.Policy.MaxHypotheses < 1 || input.Policy.MaxHypotheses > 12 || input.Policy.MaxMechanismsPerHypothesis < 1 || input.Policy.MaxMechanismsPerHypothesis > 4 || !input.Policy.RequireCounterevidenceSearch {
		return false
	}
	seenSources := map[string]bool{}
	for _, source := range input.Sources {
		if !isRunID(source.RunID) || !isLowerHex(source.SHA256, 64) || !sourcePairAllowed(source.RunKind, source.Artifact) {
			return false
		}
		identity := source.RunID + "\x00" + source.RunKind + "\x00" + source.Artifact
		if seenSources[identity] {
			return false
		}
		seenSources[identity] = true
	}
	seenProposals := map[string]bool{}
	for _, proposal := range input.Proposals {
		id := normalizeSetText(proposal.ProposalID)
		if id == "" || seenProposals[id] || proposal.Derivation != "EXTERNAL_AGENT_PROPOSAL" || len(proposal.MechanismRefs) > 4 {
			return false
		}
		seenProposals[id] = true
		if hasDuplicateNormalized(proposal.MechanismRefs) || hasDuplicateNormalized(proposal.ProductObservationRefs) || hasDuplicateNormalized(proposal.Assumptions) {
			return false
		}
	}
	return !hasDuplicateNormalized(input.Idea.Users) && !hasDuplicateNormalized(input.Idea.Capabilities) && !hasDuplicateNormalized(input.Idea.Constraints) && !hasDuplicateNormalized(input.Idea.NonGoals) && !hasDuplicateNormalized(input.Idea.PreferredStack)
}

func validIdempotencyKey(value string) bool {
	if len(value) < 1 || len(value) > 128 || strings.Contains(value, "..") {
		return false
	}
	for _, character := range []byte(value) {
		if !((character >= 0x20 && character <= 0x2e) || (character >= 0x30 && character <= 0x5b) || (character >= 0x5d && character <= 0x7e)) {
			return false
		}
	}
	return true
}

func normalizeSetText(value string) string {
	return strings.Join(strings.FieldsFunc(value, unicode.IsSpace), " ")
}
func hasDuplicateNormalized(values []string) bool {
	seen := map[string]bool{}
	for _, value := range values {
		normalized := normalizeSetText(value)
		if seen[normalized] {
			return true
		}
		seen[normalized] = true
	}
	return false
}
func sourcePairAllowed(runKind, artifact string) bool {
	return (runKind == "RECOMMENDATION" && (artifact == "mechanism_catalog" || artifact == "mechanism_gaps")) || (runKind == "MARKET_DISCOVERY" && (artifact == "market_candidates" || artifact == "market_snapshot"))
}
func isRunID(value string) bool {
	return len(value) == 36 && strings.HasPrefix(value, "run_") && isLowerHex(value[4:], 32)
}
func isError(output ContractOutput, code string) bool {
	return output.Error != nil && output.Error.Code == code
}
func errorAllowed(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
func startErrorAllowed(code string) bool {
	return errorAllowed(code, "INVALID_PARAMS", "RUN_NOT_FOUND", "SOURCE_BASIS_UNAVAILABLE", "IDEMPOTENCY_CONFLICT", "QUEUE_FULL", "QUOTA_EXCEEDED", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE")
}
func noAdmissionCalls(output ContractOutput) bool {
	return output.EffectObservation.ResolveCalls == 0 && output.EffectObservation.ReserveCalls == 0 && output.EffectObservation.IDCalls == 0 && output.EffectObservation.CreateCalls == 0
}
func noPublication(output ContractOutput) bool {
	return len(output.PublicationObservation.ArtifactNames) == 0 && len(output.PublicationObservation.ArtifactHashes) == 0 && len(output.PublicationObservation.IndexHashes) == 0
}
func equalStrings(left, right []string) bool { return reflect.DeepEqual(left, right) }
func equalStringSets(left, right []string) bool {
	leftCopy := append([]string(nil), left...)
	rightCopy := append([]string(nil), right...)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	return reflect.DeepEqual(leftCopy, rightCopy)
}
func exactHashKeys(values map[string]string, expected []string) bool {
	if len(values) != len(expected) {
		return false
	}
	for _, name := range expected {
		if !isLowerHex(values[name], 64) {
			return false
		}
	}
	return true
}
func hashMapValid(values map[string]string) bool {
	for name, hash := range values {
		if name == "" || !isLowerHex(hash, 64) {
			return false
		}
	}
	return true
}
func exactObjectKeys(object map[string]any, expected ...string) bool {
	if len(object) != len(expected) {
		return false
	}
	for _, key := range expected {
		if _, exists := object[key]; !exists {
			return false
		}
	}
	return true
}

func containsUnsafePublishedField(value any) bool {
	forbidden := map[string]bool{"body": true, "checkout_path": true, "content": true, "content_utf8": true, "cookie": true, "header": true, "idempotency_key": true, "prompt": true, "query": true, "rationale": true, "raw": true, "snippet": true, "token": true, "url": true}
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if forbidden[strings.ToLower(key)] || containsUnsafePublishedField(child) {
				return true
			}
		}
	case []any:
		for _, child := range typed {
			if containsUnsafePublishedField(child) {
				return true
			}
		}
	}
	return false
}

func containsContractualUnknown(value any) bool {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if text, scalar := child.(string); scalar {
				switch key {
				case "target_state", "novelty", "transfer_mode":
					if text == "UNKNOWN" {
						return true
					}
				case "license":
					if text == "UNRECOGNIZED" {
						return true
					}
				case "status", "status_ceiling":
					if text == "INCONCLUSIVE" {
						return true
					}
				case "code":
					if text == "MECHANISM_CLASS_UNKNOWN" {
						return true
					}
				}
			}
			if containsContractualUnknown(child) {
				return true
			}
		}
	case []any:
		for _, child := range typed {
			if containsContractualUnknown(child) {
				return true
			}
		}
	}
	return false
}

func artifactGloballyAllowed(name string) bool {
	return errorAllowed(name, "request", "recommendation", "decision", "alternatives", "failures", "provenance", "connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps", "mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance")
}
func observedIdentityMatches(observed []string, expected string) bool {
	return len(observed) > 0 && singleObservedIdentity(observed) && observed[0] == expected
}

func toolSixHasExactArtifactEnum(raw json.RawMessage) bool {
	var schema map[string]any
	if decodeStrict(raw, &schema) != nil {
		return false
	}
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return false
	}
	artifact, ok := properties["artifact"].(map[string]any)
	if !ok {
		return false
	}
	values, ok := artifact["enum"].([]any)
	if !ok || len(values) != 18 {
		return false
	}
	for index, expected := range []string{"request", "recommendation", "decision", "alternatives", "failures", "provenance", "connector_capabilities", "market_plan", "market_snapshot", "market_candidates", "market_gaps", "mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance"} {
		if values[index] != expected {
			return false
		}
	}
	return true
}

const tool10Schema = `{
"type":"object","additionalProperties":false,"properties":{
"idempotency_key":{"type":"string","minLength":1,"maxLength":128,"pattern":"^[\\u0020-\\u002E\\u0030-\\u005B\\u005D-\\u007E]+$"},
"idea":{"type":"object","additionalProperties":false,"properties":{"problem":{"type":"string","minLength":1,"maxLength":1000,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},"users":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"capabilities":{"type":"array","minItems":1,"maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"constraints":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"non_goals":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"preferred_stack":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}}},"required":["problem","capabilities"]},
"source_artifacts":{"type":"array","minItems":1,"maxItems":12,"items":{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"},"run_kind":{"type":"string","enum":["RECOMMENDATION","MARKET_DISCOVERY"]},"artifact":{"type":"string","enum":["mechanism_catalog","mechanism_gaps","market_candidates","market_snapshot"]},"sha256":{"type":"string","pattern":"^[0-9a-f]{64}$"}},"required":["run_id","run_kind","artifact","sha256"]}},
"proposals":{"type":"array","maxItems":20,"items":{"type":"object","additionalProperties":false,"properties":{"proposal_id":{"type":"string","minLength":1,"maxLength":128,"pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$"},"title":{"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},"target_property":{"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},"expected_effect":{"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},"mechanism_refs":{"type":"array","maxItems":4,"items":{"type":"string","pattern":"^mech_[0-9a-f]{64}$"}},"product_observation_refs":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":256,"pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,255}$"}},"assumptions":{"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"derivation":{"type":"string","const":"EXTERNAL_AGENT_PROPOSAL"},"rationale":{"type":"string","minLength":1,"maxLength":1000,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},"required":["proposal_id","title","target_property","expected_effect","mechanism_refs","product_observation_refs","assumptions","derivation","rationale"]}},
"policy":{"type":"object","additionalProperties":false,"properties":{"max_hypotheses":{"type":"integer","minimum":1,"maximum":12},"max_mechanisms_per_hypothesis":{"type":"integer","minimum":1,"maximum":4},"require_counterevidence_search":{"type":"boolean","const":true}},"required":["max_hypotheses","max_mechanisms_per_hypothesis","require_counterevidence_search"]}},
"required":["idempotency_key","idea","source_artifacts","proposals","policy"]}`
