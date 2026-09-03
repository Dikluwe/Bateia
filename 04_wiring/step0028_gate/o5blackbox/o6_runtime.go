/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o6_runtime.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"
)

// The O6 entry points are deliberately separate from O5. The integrator may
// replace only the corresponding dispatcher calls after sealing this revision.
func o6EvaluateSynthesisStart(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input synthesisStartInput
	localValid := decodeStrict(request, &input) == nil && o6SynthesisStartRequestValid(request, input)
	state.require(noPublication(output), "synthesis start published inline output")
	state.require(len(output.PublicationObservation.LayerFindings) == 0, "synthesis start invented evidence in publication findings")

	if !localValid {
		state.require(isError(output, "INVALID_PARAMS"), "locally invalid synthesis input did not return Invalid params")
		state.require(rawIsNull(output.ReturnValue), "locally invalid synthesis input returned a receipt")
		state.require(noAdmissionCalls(output), "local validation occurred after source resolution or admission")
		return
	}

	state.require(o6AdmissionPrefix(output), "admission calls do not form a reserve-id-create prefix")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed synthesis start returned a receipt")
		state.require(startErrorAllowed(output.Error.Code), "synthesis start returned an error outside its contract")
		o6EvaluateStartFailureEffects(state, input, output)
		if output.Error.Code == "SOURCE_BASIS_UNAVAILABLE" {
			state.unknown("immutable historical source basis is unavailable without reread")
		}
		return
	}

	state.require(output.EffectObservation.ResolveCalls == len(input.Sources), "successful synthesis start did not resolve every source")
	state.require(o6AdmissionTriple(output) == "reuse" || o6AdmissionTriple(output) == "fresh", "successful synthesis start stopped partway through admission")
	var receipt startReceipt
	if decodeReturn(output.ReturnValue, &receipt) != nil || !o6ExactRawObject(output.ReturnValue, "schema_version", "run_id", "run_kind", "lifecycle", "version") {
		state.fail("successful synthesis start returned an invalid receipt shape")
		return
	}
	state.require(receipt.SchemaVersion == "1", "synthesis receipt schema version changed")
	state.require(isRunID(receipt.RunID), "synthesis receipt run identity is invalid")
	state.require(receipt.RunKind == "IMPROVEMENT_SYNTHESIS", "synthesis receipt run kind changed")
	state.require(receipt.Lifecycle == "QUEUED", "synthesis start did not return QUEUED")
	state.require(receipt.Version >= 1, "synthesis receipt version is invalid")
	for _, source := range input.Sources {
		state.require(receipt.RunID != source.RunID, "synthesis run identity collides with an immutable source run")
	}
	if o6AdmissionTriple(output) == "fresh" {
		state.require(receipt.Version == 1, "freshly created synthesis run did not start at version one")
	}
}

func o6EvaluateRunRead(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input runInput
	if decodeStrict(request, &input) != nil || !o6ExactRawObject(request, "run_id") || !isRunID(input.RunID) {
		state.fail("run read request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "run read performed admission work")
	state.require(noPublication(output), "run read changed publication state")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed run read returned a run")
		state.require(errorAllowed(output.Error.Code, "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE"), "run read returned an error outside its contract")
		return
	}
	var view runView
	if decodeReturn(output.ReturnValue, &view) != nil || !o6RunViewShapeValid(output.ReturnValue) {
		state.fail("run read returned an invalid run representation")
		return
	}
	o6EvaluateRunView(state, input.RunID, view)
}

func o6EvaluateRunCancel(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input cancelInput
	if decodeStrict(request, &input) != nil || !o6ExactRawObject(request, "run_id", "expected_version") || !isRunID(input.RunID) || input.ExpectedVersion == 0 {
		state.fail("run cancel request is invalid")
		return
	}
	state.require(noAdmissionCalls(output), "cancellation started new resolution or admission work")
	state.require(noPublication(output), "cancellation published a new artifact set")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed cancellation returned a run")
		state.require(errorAllowed(output.Error.Code, "INVALID_PARAMS", "RUN_NOT_FOUND", "INTEGRITY_ERROR", "RUN_STORE_UNAVAILABLE"), "cancellation returned an error outside its contract")
		return
	}
	var view runView
	if decodeReturn(output.ReturnValue, &view) != nil || !o6RunViewShapeValid(output.ReturnValue) {
		state.fail("cancellation returned an invalid run representation")
		return
	}
	o6EvaluateRunView(state, input.RunID, view)
	state.require(errorAllowed(view.Lifecycle, "CANCELLING", "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED"), "cancellation returned a non-cancelling nonterminal state")
	switch view.Lifecycle {
	case "CANCELLING":
		state.require(view.Version == input.ExpectedVersion+1, "cancellation CAS must advance exactly one version")
	case "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED":
		// A terminal result may be the state already observed by the caller or
		// the single successful CAS transition caused by this request. Further
		// advancement would expose a second, non-idempotent transition.
		state.require(view.Version == input.ExpectedVersion || view.Version == input.ExpectedVersion+1, "terminal cancellation is neither idempotent nor a single CAS transition")
	}
}

func o6EvaluateReplay(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input replayInput
	if decodeStrict(request, &input) != nil || !o6ExactRawObject(request, "catalog_sha256", "capture_sha256", "artifact_hashes", "canonical_inputs") {
		state.fail("replay request is not strict protocol JSON")
		return
	}
	requestValid, hasUnknown := o6ReplayInputState(input)
	if input.CatalogSHA256 != mechanismRulesDigest || !requestValid {
		state.require(isError(output, "INTEGRITY_ERROR"), "replay accepted a divergent or malformed frozen input")
		state.require(rawIsNull(output.ReturnValue), "failed replay returned produced hashes")
		state.require(noAdmissionCalls(output) && noPublication(output), "failed replay caused runtime effects")
		return
	}
	state.require(noAdmissionCalls(output), "replay performed admission work")
	state.require(noPublication(output), "replay mutated live publication state")
	if output.Error != nil {
		state.require(rawIsNull(output.ReturnValue), "failed replay returned produced hashes")
		state.require(output.Error.Code == "INTEGRITY_ERROR", "offline replay returned a non-integrity error")
		if output.Error.Code == "INTEGRITY_ERROR" {
			state.unknown("valid frozen replay input exposes an integrity failure whose cause is not observable at this boundary")
		}
		return
	}
	var result replayResult
	if decodeReturn(output.ReturnValue, &result) != nil || !o6ExactRawObject(output.ReturnValue, "artifact_hashes", "index_hashes", "network_calls") {
		state.fail("replay returned an invalid result shape")
		return
	}
	state.require(result.NetworkCalls == 0, "replay reported network access")
	state.require(o6ProducedReplayHashes(result.ArtifactHashes), "replay artifact hashes do not cover exactly the seven produced artifacts")
	state.require(o6ProducedReplayHashes(result.IndexHashes), "replay index hashes do not cover exactly the seven produced artifacts")
	if hasUnknown {
		state.unknown("frozen replay input preserves a contractual unknown")
	}
}

func o6EvaluateRefinement(state *evaluation, request json.RawMessage, output ContractOutput) {
	var input refinementInput
	if decodeStrict(request, &input) != nil || !o6RefinementRequestShapeValid(request) {
		state.fail("refinement request is not strict protocol JSON")
		return
	}
	state.require(noAdmissionCalls(output), "refinement audit performed admission work")
	state.require(noPublication(output), "refinement audit changed publication state")
	state.require(output.Error == nil, "refinement audit returned an operational error")
	var result refinementResult
	if decodeReturn(output.ReturnValue, &result) != nil || !o6RefinementResultShapeValid(output.ReturnValue) {
		state.fail("refinement audit returned an invalid result shape")
		return
	}

	receiptTruth := []bool{
		input.ReceiptHashes.R0025 == receipt0025Digest,
		input.ReceiptHashes.R0026 == receipt0026Digest,
		input.ReceiptHashes.R0027 == receipt0027Digest,
	}
	resultTruth := []bool{result.ReceiptsVerified.R0025, result.ReceiptsVerified.R0026, result.ReceiptsVerified.R0027}
	for index := range receiptTruth {
		state.require(resultTruth[index] == receiptTruth[index], "refinement receipt result overstates or contradicts its frozen input")
	}

	state.require(o6LayerSetValid(input.LayerFindings), "refinement request contains an unknown or duplicate layer finding")
	state.require(o6SameStringSet(result.LayerFindings, input.LayerFindings), "refinement result changed the supplied layer evidence")
	completeLayers := o6SameStringSet(input.LayerFindings, o6RequiredLayers)

	observed, identitiesMatch := o6RuntimeEvidence(input, output)
	state.require(identitiesMatch, "refinement runtime identity evidence contradicts the audited runtime")
	if observed {
		state.require(result.RuntimeUnique, "complete runtime identity evidence was not represented as unique")
	} else {
		state.require(!result.RuntimeUnique, "runtime uniqueness was claimed without complete identity observations")
	}

	if !receiptTruth[0] || !receiptTruth[1] || !receiptTruth[2] || !completeLayers || !observed || !result.RuntimeUnique {
		state.unknown("refinement evidence is incomplete without a positive overclaim")
	}
}

func o6EvaluateToolsList(state *evaluation, request json.RawMessage, output ContractOutput) {
	if !o6ExactRawObject(request) {
		state.fail("tools/list request must be the empty object")
		return
	}
	state.require(output.Error == nil, "tools/list returned an error")
	state.require(noAdmissionCalls(output), "tools/list performed runtime work")
	state.require(noPublication(output), "tools/list changed publication state")
	var result toolResult
	if decodeReturn(output.ReturnValue, &result) != nil {
		state.fail("tools/list returned an invalid registry")
		return
	}
	state.require(len(result.Tools) == len(toolOrder), "tools/list must expose exactly ten tools")
	if len(result.Tools) != len(toolOrder) {
		return
	}
	for index, name := range toolOrder {
		state.require(result.Tools[index].Name == name, "tool registry name or order changed")
		state.require(o6SchemaObject(result.Tools[index].InputSchema) != nil, "tool input schema is not an object")
	}
	state.require(result.Tools[9].Description == "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis", "tool 10 description changed")
	state.require(o6Tool10SchemaValid(result.Tools[9].InputSchema), "tool 10 schema omits or weakens a required R3 constraint")
	state.require(toolSixHasExactArtifactEnum(result.Tools[5].InputSchema), "artifact getter enum is not the exact additive enum")
}

func o6EvaluateStartFailureEffects(state *evaluation, input synthesisStartInput, output ContractOutput) {
	resolve := output.EffectObservation.ResolveCalls
	state.require(resolve >= 0 && resolve <= len(input.Sources), "source resolution count exceeds the declared sources")
	switch output.Error.Code {
	case "RUN_NOT_FOUND", "SOURCE_BASIS_UNAVAILABLE", "INTEGRITY_ERROR":
		state.require(resolve >= 1, "source-dependent failure occurred without a source lookup")
		state.require(o6AdmissionTriple(output) == "reuse", "source-dependent failure reached admission")
	case "IDEMPOTENCY_CONFLICT":
		state.require(resolve == len(input.Sources), "idempotency was evaluated before all sources were frozen")
		state.require(o6AdmissionTriple(output) == "reuse", "idempotency conflict mutated admission state")
	case "QUEUE_FULL":
		state.require(resolve == len(input.Sources), "queue admission occurred before all sources were frozen")
		state.require(output.EffectObservation.ReserveCalls == 1 && output.EffectObservation.IDCalls == 0 && output.EffectObservation.CreateCalls == 0, "queue rejection does not match reserve-before-ID order")
	case "QUOTA_EXCEEDED":
		state.require(resolve == len(input.Sources), "quota admission occurred before all sources were frozen")
		state.require(output.EffectObservation.IDCalls == 0 && output.EffectObservation.CreateCalls == 0, "quota rejection created a run identity")
	case "INVALID_PARAMS":
		state.require(resolve == len(input.Sources), "resolved-reference validation did not inspect every declared source")
		state.require(o6AdmissionTriple(output) == "reuse", "resolved-reference validation reached admission")
	case "RUN_STORE_UNAVAILABLE":
		// Store failure may occur during source resolution, idempotency lookup, or
		// creation. The observable calls must still be a valid operation prefix.
		state.require(o6AdmissionPrefix(output), "store failure does not preserve admission order")
	}
}

func o6AdmissionPrefix(output ContractOutput) bool {
	reserve, id, create := output.EffectObservation.ReserveCalls, output.EffectObservation.IDCalls, output.EffectObservation.CreateCalls
	if reserve < 0 || reserve > 1 || id < 0 || id > 1 || create < 0 || create > 1 {
		return false
	}
	return (reserve == 0 && id == 0 && create == 0) ||
		(reserve == 1 && id == 0 && create == 0) ||
		(reserve == 1 && id == 1 && create == 0) ||
		(reserve == 1 && id == 1 && create == 1)
}

func o6AdmissionTriple(output ContractOutput) string {
	reserve, id, create := output.EffectObservation.ReserveCalls, output.EffectObservation.IDCalls, output.EffectObservation.CreateCalls
	switch {
	case reserve == 0 && id == 0 && create == 0:
		return "reuse"
	case reserve == 1 && id == 1 && create == 1:
		return "fresh"
	default:
		return "partial"
	}
}

func o6SynthesisStartRequestValid(raw json.RawMessage, input synthesisStartInput) bool {
	if !o6ExactRawObject(raw, "idempotency_key", "idea", "source_artifacts", "proposals", "policy") || !validIdempotencyKey(input.IdempotencyKey) {
		return false
	}
	var root map[string]json.RawMessage
	if decodeStrict(raw, &root) != nil || !o6IdeaValid(root["idea"], input.Idea) || !o6PolicyValid(root["policy"], input.Policy) {
		return false
	}
	if input.Sources == nil || len(input.Sources) < 1 || len(input.Sources) > 12 || input.Proposals == nil || len(input.Proposals) > 20 {
		return false
	}
	seenSources := map[string]bool{}
	runKinds := map[string]string{}
	artifactByHash := map[string]string{}
	hasMechanismSource := false
	hasProductSource := false
	for _, source := range input.Sources {
		identity := source.RunID + "\x00" + source.RunKind + "\x00" + source.Artifact
		if !isRunID(source.RunID) || !sourcePairAllowed(source.RunKind, source.Artifact) || !isLowerHex(source.SHA256, 64) || seenSources[identity] {
			return false
		}
		if prior, exists := runKinds[source.RunID]; exists && prior != source.RunKind {
			return false
		}
		if prior, exists := artifactByHash[source.SHA256]; exists && prior != source.Artifact {
			return false
		}
		runKinds[source.RunID] = source.RunKind
		artifactByHash[source.SHA256] = source.Artifact
		hasMechanismSource = hasMechanismSource || source.Artifact == "mechanism_catalog"
		hasProductSource = hasProductSource || source.Artifact == "market_candidates" || source.Artifact == "market_snapshot"
		seenSources[identity] = true
	}
	seenProposals := map[string]bool{}
	for _, proposal := range input.Proposals {
		id := normalizeSetText(proposal.ProposalID)
		if !o6Identifier(id, 128) || seenProposals[id] || proposal.Derivation != "EXTERNAL_AGENT_PROPOSAL" ||
			!o6Text(proposal.Title, 512) || !o6Text(proposal.TargetProperty, 512) || !o6Text(proposal.ExpectedEffect, 512) || !o6Text(proposal.Rationale, 1000) ||
			proposal.MechanismRefs == nil || len(proposal.MechanismRefs) > 4 || proposal.ProductObservationRefs == nil || len(proposal.ProductObservationRefs) > 20 || proposal.Assumptions == nil || len(proposal.Assumptions) > 20 ||
			hasDuplicateNormalized(proposal.MechanismRefs) || hasDuplicateNormalized(proposal.ProductObservationRefs) || hasDuplicateNormalized(proposal.Assumptions) {
			return false
		}
		seenProposals[id] = true
		for _, ref := range proposal.MechanismRefs {
			if len(ref) != 69 || !strings.HasPrefix(ref, "mech_") || !isLowerHex(ref[5:], 64) {
				return false
			}
		}
		for _, ref := range proposal.ProductObservationRefs {
			if !o6Identifier(ref, 256) {
				return false
			}
			if strings.HasPrefix(ref, "mech_") && len(ref) == 69 && isLowerHex(ref[5:], 64) {
				return false
			}
		}
		for _, assumption := range proposal.Assumptions {
			if !o6Text(assumption, 512) {
				return false
			}
		}
		if len(proposal.MechanismRefs) > 0 && !hasMechanismSource {
			return false
		}
		if len(proposal.ProductObservationRefs) > 0 && !hasProductSource {
			return false
		}
	}
	return true
}

func o6IdeaValid(raw json.RawMessage, idea ideaInput) bool {
	if !o6ObjectHasAllowedRequired(raw,
		[]string{"problem", "users", "capabilities", "constraints", "non_goals", "preferred_stack"},
		[]string{"problem", "capabilities"}) || !o6Text(idea.Problem, 1000) || idea.Capabilities == nil || len(idea.Capabilities) < 1 {
		return false
	}
	sets := [][]string{idea.Users, idea.Capabilities, idea.Constraints, idea.NonGoals, idea.PreferredStack}
	for index, set := range sets {
		if (index == 1 && set == nil) || len(set) > 20 || hasDuplicateNormalized(set) {
			return false
		}
		for _, item := range set {
			if !o6Text(item, 200) {
				return false
			}
		}
	}
	return true
}

func o6PolicyValid(raw json.RawMessage, policy synthesisPolicyInput) bool {
	return o6ExactRawObject(raw, "max_hypotheses", "max_mechanisms_per_hypothesis", "require_counterevidence_search") &&
		policy.MaxHypotheses >= 1 && policy.MaxHypotheses <= 12 && policy.MaxMechanismsPerHypothesis >= 1 && policy.MaxMechanismsPerHypothesis <= 4 && policy.RequireCounterevidenceSearch
}

func o6Text(value string, maximumBytes int) bool {
	if value == "" || normalizeSetText(value) == "" || !utf8.ValidString(value) || len([]byte(value)) > maximumBytes {
		return false
	}
	for _, character := range value {
		if character <= 0x1f || character == 0x7f || (character >= 0x80 && character <= 0x9f) {
			return false
		}
	}
	return true
}

func o6Identifier(value string, maximumBytes int) bool {
	if len(value) < 1 || len(value) > maximumBytes || !o6ASCIIAlphaNumeric(value[0]) {
		return false
	}
	for index := 1; index < len(value); index++ {
		character := value[index]
		if !o6ASCIIAlphaNumeric(character) && character != '.' && character != '_' && character != ':' && character != '-' {
			return false
		}
	}
	return true
}

func o6ASCIIAlphaNumeric(character byte) bool {
	return (character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9')
}

func o6RunViewShapeValid(raw json.RawMessage) bool {
	if !o6ExactRawObject(raw, "run_id", "run_kind", "lifecycle", "version", "recommendation_status", "market_discovery_status", "improvement_status", "progress") {
		return false
	}
	var root map[string]json.RawMessage
	return decodeStrict(raw, &root) == nil && o6ExactRawObject(root["progress"], "planned", "completed", "stage")
}

func o6EvaluateRunView(state *evaluation, requestedRunID string, view runView) {
	state.require(view.RunID == requestedRunID && isRunID(view.RunID), "run representation does not match the requested identity")
	state.require(errorAllowed(view.RunKind, "RECOMMENDATION", "MARKET_DISCOVERY", "IMPROVEMENT_SYNTHESIS"), "run representation has an unknown kind")
	state.require(errorAllowed(view.Lifecycle, "QUEUED", "RUNNING", "CANCELLING", "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED"), "run representation has an unknown lifecycle")
	state.require(view.Version >= 1, "run representation has an invalid version")
	state.require(view.Progress.Planned >= 0 && view.Progress.Completed >= 0 && view.Progress.Completed <= view.Progress.Planned, "run progress is outside its planned bound")
	state.require(o6Text(view.Progress.Stage, 512), "run progress stage is invalid")

	terminal := errorAllowed(view.Lifecycle, "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED")
	statuses := []json.RawMessage{view.RecommendationStatus, view.MarketDiscoveryStatus, view.ImprovementStatus}
	corresponding := map[string]int{"RECOMMENDATION": 0, "MARKET_DISCOVERY": 1, "IMPROVEMENT_SYNTHESIS": 2}[view.RunKind]
	for index, status := range statuses {
		if index != corresponding || !terminal {
			state.require(rawIsNull(status), "run exposed a status for the wrong kind or before terminal")
		}
	}
	if view.Lifecycle == "SUCCEEDED" {
		state.require(!rawIsNull(statuses[corresponding]), "successful run omitted its kind-specific result status")
	}
	if view.Lifecycle == "QUEUED" {
		state.require(view.Progress.Completed == 0 && view.Progress.Stage == "queued", "queued progress representation changed")
	}
	if terminal {
		state.require(view.Progress.Stage == "terminal", "terminal lifecycle is not represented by terminal stage")
	}
	if view.RunKind != "IMPROVEMENT_SYNTHESIS" {
		return
	}
	state.require(view.Progress.Planned == 5, "improvement synthesis does not retain five planned stages")
	state.require(o6SynthesisStageValid(view.Lifecycle, view.Progress.Stage), "improvement synthesis stage is outside the ordered lifecycle vocabulary")
	if !rawIsNull(view.ImprovementStatus) {
		var status string
		if decodeStrict(view.ImprovementStatus, &status) != nil || !errorAllowed(status, "COMPLETE", "PARTIAL", "INCONCLUSIVE") {
			state.fail("improvement status is outside the closed enum")
		} else if view.Lifecycle != "SUCCEEDED" && status == "COMPLETE" {
			state.fail("degraded synthesis lifecycle cannot claim COMPLETE")
		} else if status == "INCONCLUSIVE" {
			state.unknown("improvement synthesis preserves explicit epistemic insufficiency")
		}
	}
	if view.Lifecycle == "SUCCEEDED" {
		state.require(view.Progress.Completed == 5, "successful synthesis did not complete five stages")
	} else if !terminal {
		state.require(view.Progress.Completed < 5, "nonterminal synthesis claims every stage completed")
	}
}

func o6SynthesisStageValid(lifecycle, stage string) bool {
	if lifecycle == "QUEUED" {
		return stage == "queued"
	}
	if errorAllowed(lifecycle, "CANCELLED", "FAILED", "SUCCEEDED", "INTERRUPTED") {
		return stage == "terminal"
	}
	return errorAllowed(stage, "queued", "synthesis_compare", "synthesis_derive", "synthesis_counterevidence", "synthesis_validate", "synthesis_publish")
}

var o6ProducedArtifacts = []string{
	"mechanism_catalog", "mechanism_gaps", "architecture_comparison", "improvement_hypotheses", "hypothesis_gaps", "validation_plans", "synthesis_provenance",
}

func o6ReplayInputState(input replayInput) (bool, bool) {
	if !isLowerHex(input.CatalogSHA256, 64) || !isLowerHex(input.CaptureSHA256, 64) || len(input.ArtifactHashes) < 1 || len(input.CanonicalInputs) < 1 {
		return false, false
	}
	for name, hash := range input.ArtifactHashes {
		if name == "" || !isLowerHex(hash, 64) {
			return false, false
		}
	}
	hasUnknown := false
	for _, encoded := range input.CanonicalInputs {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil || base64.StdEncoding.EncodeToString(decoded) != encoded {
			return false, false
		}
		var value any
		decoder := json.NewDecoder(bytes.NewReader(decoded))
		decoder.UseNumber()
		if decoder.Decode(&value) == nil && o6ContainsEpistemicUnknown(value) {
			hasUnknown = true
		}
	}
	return true, hasUnknown
}

func o6ProducedReplayHashes(values map[string]string) bool {
	if len(values) != len(o6ProducedArtifacts) {
		return false
	}
	for _, artifact := range o6ProducedArtifacts {
		if !isLowerHex(values[artifact], 64) {
			return false
		}
	}
	return true
}

func o6ContainsEpistemicUnknown(value any) bool {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			lowerKey := strings.ToLower(key)
			if child == nil && strings.Contains(lowerKey, "baseline") {
				return true
			}
			if flag, ok := child.(bool); ok {
				if (!flag && errorAllowed(lowerKey, "baseline_present", "provider_availability_evidenced", "composition_compatibility_measured")) ||
					(flag && errorAllowed(lowerKey, "counterevidence_budget_exhausted", "baseline_absent")) {
					return true
				}
			}
			if text, ok := child.(string); ok {
				if (errorAllowed(lowerKey, "target_state", "novelty", "transfer_mode", "provider_availability") && text == "UNKNOWN") ||
					(lowerKey == "license" && text == "UNRECOGNIZED") ||
					(lowerKey == "status" && errorAllowed(text, "INCONCLUSIVE", "BUDGET_EXHAUSTED")) ||
					(lowerKey == "status_ceiling" && text == "INCONCLUSIVE") {
					return true
				}
			}
			if lowerKey == "integration_unknowns" {
				if unknowns, ok := child.([]any); ok && len(unknowns) > 0 {
					return true
				}
			}
			if o6ContainsEpistemicUnknown(child) {
				return true
			}
		}
	case []any:
		for _, child := range typed {
			if o6ContainsEpistemicUnknown(child) {
				return true
			}
		}
	}
	return false
}

var o6RequiredLayers = []string{"L1_PURE", "L2_NO_L3_OR_L4", "L3_NO_L2_OR_L4", "L4_WIRING_ONLY"}

func o6LayerSetValid(values []string) bool {
	if values == nil || len(values) > len(o6RequiredLayers) {
		return false
	}
	seen := map[string]bool{}
	for _, value := range values {
		if seen[value] || !errorAllowed(value, o6RequiredLayers...) {
			return false
		}
		seen[value] = true
	}
	return true
}

func o6RuntimeEvidence(input refinementInput, output ContractOutput) (bool, bool) {
	pairs := []struct {
		observed []string
		declared string
	}{
		{output.EffectObservation.StoreInstanceIDs, input.RuntimeInstanceIDs.Store},
		{output.EffectObservation.QueueInstanceIDs, input.RuntimeInstanceIDs.Queue},
		{output.EffectObservation.RootInstanceIDs, input.RuntimeInstanceIDs.Root},
		{output.EffectObservation.CursorSecretInstanceIDs, input.RuntimeInstanceIDs.CursorSecret},
		{output.EffectObservation.WorkerInstanceIDs, input.RuntimeInstanceIDs.Worker},
	}
	complete := true
	for _, pair := range pairs {
		if len(pair.observed) == 0 {
			complete = false
			continue
		}
		if !singleObservedIdentity(pair.observed) || pair.observed[0] != pair.declared {
			return false, false
		}
	}
	return complete, true
}

func o6RefinementRequestShapeValid(raw json.RawMessage) bool {
	if !o6ExactRawObject(raw, "receipt_hashes", "runtime_instance_ids", "layer_findings") {
		return false
	}
	var root map[string]json.RawMessage
	return decodeStrict(raw, &root) == nil &&
		o6ExactRawObject(root["receipt_hashes"], "0025", "0026", "0027") &&
		o6ExactRawObject(root["runtime_instance_ids"], "store", "queue", "root", "cursor_secret", "worker")
}

func o6RefinementResultShapeValid(raw json.RawMessage) bool {
	if !o6ExactRawObject(raw, "receipts_verified", "runtime_unique", "layer_findings") {
		return false
	}
	var root map[string]json.RawMessage
	return decodeStrict(raw, &root) == nil && o6ExactRawObject(root["receipts_verified"], "0025", "0026", "0027")
}

func o6SameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}

func o6Tool10SchemaValid(raw json.RawMessage) bool {
	schema := o6SchemaObject(raw)
	if schema == nil || schema["type"] != "object" || schema["additionalProperties"] != false || !o6SchemaRequired(schema, "idempotency_key", "idea", "source_artifacts", "proposals", "policy") {
		return false
	}
	properties, ok := schema["properties"].(map[string]any)
	if !ok || !o6MapHasExactly(properties, "idempotency_key", "idea", "source_artifacts", "proposals", "policy") {
		return false
	}
	idempotency, ok := properties["idempotency_key"].(map[string]any)
	if !ok || !o6JSONNumberEquals(idempotency["minLength"], 1) || !o6JSONNumberEquals(idempotency["maxLength"], 128) {
		return false
	}
	idea, ok := properties["idea"].(map[string]any)
	sources, sourcesOK := properties["source_artifacts"].(map[string]any)
	proposals, proposalsOK := properties["proposals"].(map[string]any)
	policy, policyOK := properties["policy"].(map[string]any)
	if !ok || !sourcesOK || !proposalsOK || !policyOK || idea["additionalProperties"] != false || policy["additionalProperties"] != false {
		return false
	}
	ideaProperties, ideaPropertiesOK := idea["properties"].(map[string]any)
	if !ideaPropertiesOK || !o6MapHasExactly(ideaProperties, "problem", "users", "capabilities", "constraints", "non_goals", "preferred_stack") {
		return false
	}
	if !o6SchemaRequired(idea, "problem", "capabilities") || !o6ArrayBounds(sources, 1, 12) || !o6ArrayMaximum(proposals, 20) {
		return false
	}
	sourceItems, ok := sources["items"].(map[string]any)
	if !ok || sourceItems["additionalProperties"] != false || !o6SchemaRequired(sourceItems, "run_id", "run_kind", "artifact", "sha256") {
		return false
	}
	sourceProperties, ok := sourceItems["properties"].(map[string]any)
	if !ok || !o6MapHasExactly(sourceProperties, "run_id", "run_kind", "artifact", "sha256") || !o6SchemaEnum(sourceProperties["run_kind"], "RECOMMENDATION", "MARKET_DISCOVERY") || !o6SchemaEnum(sourceProperties["artifact"], "mechanism_catalog", "mechanism_gaps", "market_candidates", "market_snapshot") {
		return false
	}
	proposalItems, ok := proposals["items"].(map[string]any)
	if !ok || proposalItems["additionalProperties"] != false || !o6SchemaRequired(proposalItems, "proposal_id", "title", "target_property", "expected_effect", "mechanism_refs", "product_observation_refs", "assumptions", "derivation", "rationale") {
		return false
	}
	proposalProperties, ok := proposalItems["properties"].(map[string]any)
	if !ok || !o6MapHasExactly(proposalProperties, "proposal_id", "title", "target_property", "expected_effect", "mechanism_refs", "product_observation_refs", "assumptions", "derivation", "rationale") {
		return false
	}
	derivation, ok := proposalProperties["derivation"].(map[string]any)
	if !ok || derivation["const"] != "EXTERNAL_AGENT_PROPOSAL" {
		return false
	}
	policyProperties, ok := policy["properties"].(map[string]any)
	if !ok || !o6MapHasExactly(policyProperties, "max_hypotheses", "max_mechanisms_per_hypothesis", "require_counterevidence_search") || !o6SchemaRequired(policy, "max_hypotheses", "max_mechanisms_per_hypothesis", "require_counterevidence_search") {
		return false
	}
	requireSearch, ok := policyProperties["require_counterevidence_search"].(map[string]any)
	return ok && requireSearch["const"] == true && o6NumericBounds(policyProperties["max_hypotheses"], 1, 12) && o6NumericBounds(policyProperties["max_mechanisms_per_hypothesis"], 1, 4)
}

func o6SchemaObject(raw json.RawMessage) map[string]any {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var schema map[string]any
	if decoder.Decode(&schema) != nil || schema == nil {
		return nil
	}
	return schema
}

func o6SchemaRequired(schema map[string]any, required ...string) bool {
	values, ok := schema["required"].([]any)
	if !ok || len(values) != len(required) {
		return false
	}
	seen := map[string]bool{}
	for _, value := range values {
		text, ok := value.(string)
		if !ok || seen[text] {
			return false
		}
		seen[text] = true
	}
	for _, name := range required {
		if !seen[name] {
			return false
		}
	}
	return true
}

func o6SchemaEnum(value any, required ...string) bool {
	schema, ok := value.(map[string]any)
	if !ok {
		return false
	}
	values, ok := schema["enum"].([]any)
	if !ok || len(values) != len(required) {
		return false
	}
	for index, item := range values {
		if item != required[index] {
			return false
		}
	}
	return true
}

func o6ArrayBounds(schema map[string]any, minimum, maximum int) bool {
	return schema["type"] == "array" && o6JSONNumberEquals(schema["minItems"], minimum) && o6JSONNumberEquals(schema["maxItems"], maximum)
}

func o6ArrayMaximum(schema map[string]any, maximum int) bool {
	return schema["type"] == "array" && o6JSONNumberEquals(schema["maxItems"], maximum)
}

func o6NumericBounds(value any, minimum, maximum int) bool {
	schema, ok := value.(map[string]any)
	return ok && o6JSONNumberEquals(schema["minimum"], minimum) && o6JSONNumberEquals(schema["maximum"], maximum)
}

func o6JSONNumberEquals(value any, expected int) bool {
	number, ok := value.(json.Number)
	return ok && number.String() == strconv.Itoa(expected)
}

func o6ExactRawObject(raw json.RawMessage, keys ...string) bool {
	var object map[string]json.RawMessage
	if decodeStrict(raw, &object) != nil || len(object) != len(keys) {
		return false
	}
	for _, key := range keys {
		if _, exists := object[key]; !exists {
			return false
		}
	}
	return true
}

func o6ObjectHasAllowedRequired(raw json.RawMessage, allowed, required []string) bool {
	var object map[string]json.RawMessage
	if decodeStrict(raw, &object) != nil {
		return false
	}
	allowedSet := map[string]bool{}
	for _, key := range allowed {
		allowedSet[key] = true
	}
	for key := range object {
		if !allowedSet[key] {
			return false
		}
	}
	for _, key := range required {
		if _, exists := object[key]; !exists {
			return false
		}
	}
	return true
}

func o6MapHasExactly(object map[string]any, keys ...string) bool {
	if len(object) != len(keys) {
		return false
	}
	for _, key := range keys {
		if _, exists := object[key]; !exists {
			return false
		}
	}
	return true
}
