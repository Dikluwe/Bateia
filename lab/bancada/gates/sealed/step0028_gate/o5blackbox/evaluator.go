/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/evaluator.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type evaluation struct {
	finding Finding
}

func evaluate(contractInput, contractOutput []byte) Finding {
	state := &evaluation{finding: Finding{Failures: []string{}, Unknowns: []string{}}}
	if !exactTopLevelKeys(contractInput, "port", "request") {
		state.fail("contract_input does not have the exact protocol envelope")
		return state.finish()
	}
	var input ContractInput
	if err := decodeStrict(contractInput, &input); err != nil {
		state.fail("contract_input is not strict protocol JSON")
		return state.finish()
	}
	state.finding.Port = input.Port
	if !exactTopLevelKeys(contractOutput, "return_value", "error", "effect_observation", "publication_observation") {
		state.fail("mutant_contract_output does not have the exact protocol envelope")
		return state.finish()
	}
	var output ContractOutput
	if err := decodeStrict(contractOutput, &output); err != nil {
		state.fail("mutant_contract_output is not strict protocol JSON")
		return state.finish()
	}
	state.checkGlobal(output)
	if len(state.finding.Failures) == 0 {
		switch input.Port {
		case MechanismExtract:
			evaluateMechanismExtractO6(state, input.Request, output)
		case SynthesisStart:
			o6EvaluateSynthesisStart(state, input.Request, output)
		case SynthesisExecute:
			evaluateSynthesisExecute(state, input.Request, output)
		case RunRead:
			o6EvaluateRunRead(state, input.Request, output)
		case ArtifactRead:
			evaluateArtifactReadO6(state, input.Request, output)
		case RunCancel:
			o6EvaluateRunCancel(state, input.Request, output)
		case ToolsList:
			o6EvaluateToolsList(state, input.Request, output)
		case Replay:
			o6EvaluateReplay(state, input.Request, output)
		case RefinementAudit:
			o6EvaluateRefinement(state, input.Request, output)
		default:
			state.fail("contract_input port is outside the closed protocol")
		}
	}
	return state.finish()
}

func (state *evaluation) checkGlobal(output ContractOutput) {
	state.require(output.EffectObservation.NetworkCalls == 0, "network effects are forbidden")
	state.require(output.EffectObservation.FilesystemWrites == 0, "filesystem writes are forbidden")
	state.require(output.EffectObservation.ProcessSpawns == 0, "process spawning is forbidden")
	state.require(output.EffectObservation.ExternalMutations == 0, "external mutation is forbidden")
	state.require(output.EffectObservation.StoreInstanceIDs != nil, "store instance observations must be an array")
	state.require(output.EffectObservation.QueueInstanceIDs != nil, "queue instance observations must be an array")
	state.require(output.EffectObservation.RootInstanceIDs != nil, "root instance observations must be an array")
	state.require(output.EffectObservation.CursorSecretInstanceIDs != nil, "cursor-secret observations must be an array")
	state.require(output.EffectObservation.WorkerInstanceIDs != nil, "worker instance observations must be an array")
	state.require(singleObservedIdentity(output.EffectObservation.StoreInstanceIDs), "multiple store identities observed")
	state.require(singleObservedIdentity(output.EffectObservation.QueueInstanceIDs), "multiple queue identities observed")
	state.require(singleObservedIdentity(output.EffectObservation.RootInstanceIDs), "multiple root identities observed")
	state.require(singleObservedIdentity(output.EffectObservation.CursorSecretInstanceIDs), "multiple cursor-secret identities observed")
	state.require(singleObservedIdentity(output.EffectObservation.WorkerInstanceIDs), "multiple worker identities observed")
	state.require(output.PublicationObservation.ArtifactNames != nil, "artifact_names must be an array")
	state.require(output.PublicationObservation.ArtifactHashes != nil, "artifact_hashes must be an object")
	state.require(output.PublicationObservation.IndexHashes != nil, "index_hashes must be an object")
	state.require(output.PublicationObservation.LeakedFieldNames != nil, "leaked_field_names must be an array")
	state.require(output.PublicationObservation.LayerFindings != nil, "layer_findings must be an array")
	state.require(len(output.PublicationObservation.LeakedFieldNames) == 0, "unsafe fields crossed the publication boundary")
	state.require(output.EffectObservation.ResolveCalls >= 0 && output.EffectObservation.ReserveCalls >= 0 && output.EffectObservation.IDCalls >= 0 && output.EffectObservation.CreateCalls >= 0, "negative call count")

	returned := !rawIsNull(output.ReturnValue)
	failed := output.Error != nil
	state.require(!(returned && failed), "return_value and error are mutually exclusive")
	if failed {
		state.require(validControlledError(*output.Error), "error code/message is not a controlled contract pair")
	}
	for name, hash := range output.PublicationObservation.ArtifactHashes {
		state.require(name != "" && isLowerHex(hash, 64), "publication artifact hash is invalid")
	}
	for name, hash := range output.PublicationObservation.IndexHashes {
		state.require(name != "" && isLowerHex(hash, 64), "publication index hash is invalid")
	}
}

func evaluateSequence(exchanges []Exchange) []Finding {
	findings := make([]Finding, len(exchanges))
	canonicalOutputs := map[string][]byte{}
	canonicalIndices := map[string][]int{}
	identities := map[string]string{}
	for index, exchange := range exchanges {
		findings[index] = evaluate(exchange.ContractInput, exchange.ContractOutput)
		input, inputErr := canonicalJSON(exchange.ContractInput)
		output, outputErr := canonicalJSON(exchange.ContractOutput)
		if inputErr == nil && outputErr == nil {
			key := string(input)
			if prior, exists := canonicalOutputs[key]; exists && !bytes.Equal(prior, output) {
				for _, affected := range append(canonicalIndices[key], index) {
					appendSequenceFailure(&findings[affected], "same contract_input produced a different contract_output")
				}
			} else if !exists {
				canonicalOutputs[key] = output
			}
			canonicalIndices[key] = append(canonicalIndices[key], index)
		}
		var decoded ContractOutput
		if decodeStrict(exchange.ContractOutput, &decoded) == nil {
			observeSequenceIdentity(findings, index, identities, "store", decoded.EffectObservation.StoreInstanceIDs)
			observeSequenceIdentity(findings, index, identities, "queue", decoded.EffectObservation.QueueInstanceIDs)
			observeSequenceIdentity(findings, index, identities, "root", decoded.EffectObservation.RootInstanceIDs)
			observeSequenceIdentity(findings, index, identities, "cursor_secret", decoded.EffectObservation.CursorSecretInstanceIDs)
			observeSequenceIdentity(findings, index, identities, "worker", decoded.EffectObservation.WorkerInstanceIDs)
		}
	}
	return findings
}

func observeSequenceIdentity(findings []Finding, index int, identities map[string]string, kind string, observed []string) {
	if len(observed) == 0 {
		return
	}
	if prior, exists := identities[kind]; exists && prior != observed[0] {
		appendSequenceFailure(&findings[index], fmt.Sprintf("%s runtime identity changed", kind))
		return
	}
	identities[kind] = observed[0]
}

func appendSequenceFailure(finding *Finding, reason string) {
	for _, existing := range finding.Failures {
		if existing == reason {
			return
		}
	}
	finding.Failures = append(finding.Failures, reason)
	finding.Classification = Violated
}

func (state *evaluation) require(condition bool, reason string) {
	state.finding.Predicates++
	if !condition {
		state.fail(reason)
	}
}

func (state *evaluation) fail(reason string) {
	state.finding.Failures = append(state.finding.Failures, reason)
}

func (state *evaluation) unknown(reason string) {
	state.finding.Unknowns = append(state.finding.Unknowns, reason)
}

func (state *evaluation) finish() Finding {
	switch {
	case len(state.finding.Failures) > 0:
		state.finding.Classification = Violated
	case len(state.finding.Unknowns) > 0:
		state.finding.Classification = Unknown
	default:
		state.finding.Classification = Preserved
	}
	return state.finding
}

func singleObservedIdentity(values []string) bool {
	if len(values) == 0 {
		return true
	}
	if values[0] == "" {
		return false
	}
	for _, value := range values[1:] {
		if value != values[0] {
			return false
		}
	}
	return true
}

func validControlledError(contractError ControlledError) bool {
	pairs := map[string]string{
		"INVALID_PARAMS": "Invalid params", "RUN_NOT_FOUND": "RUN_NOT_FOUND",
		"SOURCE_BASIS_UNAVAILABLE": "SOURCE_BASIS_UNAVAILABLE", "IDEMPOTENCY_CONFLICT": "IDEMPOTENCY_CONFLICT",
		"QUEUE_FULL": "QUEUE_FULL", "QUOTA_EXCEEDED": "QUOTA_EXCEEDED", "INTEGRITY_ERROR": "INTEGRITY_ERROR",
		"RUN_STORE_UNAVAILABLE": "RUN_STORE_UNAVAILABLE", "ARTIFACT_NOT_READY": "ARTIFACT_NOT_READY",
	}
	return pairs[contractError.Code] == contractError.Message
}

func decodeReturn(raw json.RawMessage, target any) error {
	if rawIsNull(raw) {
		return fmt.Errorf("return_value is null")
	}
	return decodeStrict(raw, target)
}
