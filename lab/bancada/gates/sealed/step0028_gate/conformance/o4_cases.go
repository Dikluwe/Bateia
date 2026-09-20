/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o4_cases.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import "fmt"

func O4PositiveCase(commitment O4Commitment) (Case, error) {
	for _, testCase := range O3PositiveCases() {
		if testCase.ID == commitment.ID && testCase.Operation == commitment.RequiredPort {
			return testCase, nil
		}
	}
	return Case{}, fmt.Errorf("no positive control for %s/%s", commitment.ID, commitment.RequiredPort)
}

func O4NegativeCase(commitment O4Commitment) (Case, error) {
	var testCase Case
	switch commitment.RequiredPort {
	case OperationExtract:
		testCase = mustO4Control("P01")
	case OperationHypothesis:
		testCase = mustO4Control("P04")
	case OperationStart:
		testCase = mustO4Control("P09")
	case OperationLifecycle:
		testCase = mustO4Control("P10")
	case OperationReplay:
		testCase = mustO4Control("P11")
	case OperationRuntime:
		testCase = mustO4Control("P12")
		testCase.Operation = OperationRuntime
	case OperationFull:
		testCase = mustO4Control("P12")
	case OperationTools:
		testCase = mustO4Control("P12")
		testCase.Operation = OperationTools
		testCase.Input.AppendTool10 = true
	case OperationRefinement:
		testCase = mustO4Control("P12")
		testCase.Operation = OperationRefinement
		testCase.Input.RefinementHashes = cloneMap(testCase.Input.RequiredRefinementHashes)
	case OperationPagination:
		testCase = Case{Operation: OperationPagination, Input: Input{
			ScenarioVersion: "O4", AppliedLimits: map[string]int{"page_size": 2},
			SnapshotSHA256: "sha256:frozen-snapshot", Cursor: "cursor:1", NextCursor: "cursor:2", PageItems: []string{"item-a", "item-b"},
		}}
	case OperationSchema:
		testCase = Case{Operation: OperationSchema, Input: Input{
			ScenarioVersion:   "O4",
			ToolDefinitions:   map[string]string{"bateia_start_improvement_synthesis": "sha256:tool10"},
			SchemaDefinitions: map[string]string{"bateia_start_improvement_synthesis": "sha256:schema10"},
		}}
	case OperationArtifact:
		testCase = Case{Operation: OperationArtifact, Input: Input{
			ScenarioVersion: "O4", Lifecycle: "SUCCEEDED", ArtifactState: "PUBLISHED",
			ArtifactTransitions: []string{"STAGED", "PUBLISHED"},
			ArtifactNames:       []string{"comparative-synthesis.json", "improvement-hypotheses.json", "transfer-plans.json", "mechanisms.json", "gaps.json", "counterevidence.json", "provenance.json"},
			ArtifactHashes:      map[string]string{"comparative-synthesis.json": "sha256:a", "improvement-hypotheses.json": "sha256:b", "transfer-plans.json": "sha256:c", "mechanisms.json": "sha256:d", "gaps.json": "sha256:e", "counterevidence.json": "sha256:f", "provenance.json": "sha256:g"},
		}}
	case OperationSafety:
		testCase = Case{Operation: OperationSafety, Input: Input{
			ScenarioVersion: "O4", Lifecycle: "", RawFields: []string{}, ExternalEffects: []string{}, LayerFindings: []string{},
		}}
	default:
		return Case{}, fmt.Errorf("no negative fixture port for %s/%s", commitment.ID, commitment.RequiredPort)
	}
	testCase.ID = commitment.ID
	testCase.Input.ScenarioVersion = "O4"
	testCase.Expected = frozenO4Expected(testCase)
	return testCase, nil
}

func O4OpaqueCase(commitment O4Commitment) (Case, error) {
	switch commitment.RequiredPort {
	case OperationExtract, OperationHypothesis, OperationReplay, OperationFull:
	default:
		return Case{}, fmt.Errorf("unsupported opaque port %s", commitment.RequiredPort)
	}
	testCase := Case{ID: commitment.ID, Operation: commitment.RequiredPort, Input: Input{
		ScenarioVersion: "O4", Opaque: true, OpacityReason: "UNRESOLVED_PUBLIC_EVIDENCE",
		TargetState: "UNKNOWN", MechanismStatus: "UNKNOWN", ClaimedNovelty: "UNKNOWN",
		License: "UNRECOGNIZED", RequestedTransfer: "PATTERN_ADAPTATION",
		Assumptions: []string{"evidence boundary unresolved"}, Risks: []string{"classification would overclaim"},
		ValidationPlan: []string{"obtain immutable witness before classification"},
		LayerFindings:  []string{"UNKNOWN_AT_OPAQUE_BOUNDARY"},
	}}
	return testCase, nil
}

func mustO4Control(id string) Case {
	for _, testCase := range O3PositiveCases() {
		if testCase.ID == id {
			return testCase
		}
	}
	panic("missing frozen O4 control " + id)
}

func frozenO4Expected(testCase Case) Observation {
	input := testCase.Input
	result := frozenProjection(input)
	switch testCase.Operation {
	case OperationExtract:
		result.Decision, result.MechanismStatus, result.Novelty = "MECHANISM_PRESERVED", "PROVEN", "CONVENTIONAL"
		result.BindingVerified, result.ProvenanceComplete = true, true
		result.CanonicalSHA256 = canonicalDigest([]string{input.Owner, input.Repo, input.Commit, input.Path, input.Language, input.Ecosystem})
	case OperationHypothesis:
		result.MechanismStatus, result.Novelty = "PROVEN", "CONVENTIONAL"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
	case OperationStart:
		result.ResolveCalls, result.ReserveCalls, result.IDCalls, result.CreateCalls = 1, 1, 1, 1
		result.RunID, result.Lifecycle = "run_reference", "QUEUED"
		result.RunIDs = []string{"run_reference", "run_reference", "run_reference", "run_reference"}
		result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
	case OperationLifecycle:
		result.Lifecycle, result.HypothesisStatus = "CANCELLED", "PROPOSED"
		result.RecoveryState, result.AdoptionDisposition = "INTERRUPTED", "NOT_AN_ADOPTION_DECISION"
	case OperationReplay:
		result.ReplaySHA256 = canonicalDigest(sortedCopy(input.ReplayItems))
	case OperationRuntime:
		result.RuntimeUnique = true
		result.StoreIdentity, result.QueueIdentity, result.RootIdentity = "store-1", "queue-1", "root-1"
		result.SecretIdentity, result.WorkerIdentity = "secret-1", "worker-1"
	case OperationTools:
		result.BasalHashVerified = true
		result.Tools = append([]Tool(nil), input.BasalTools...)
		result.Tools = append(result.Tools, Tool{Name: "bateia_start_improvement_synthesis", Definition: "sha256:frozen-tool-10-schema"})
	case OperationRefinement:
		result.RefinementsPreserved = true
	case OperationArtifact:
		result.Lifecycle = input.Lifecycle
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
	case OperationFull:
		result = frozenO3Expected(Case{ID: "P12", Operation: OperationFull, Input: input})
	case OperationPagination, OperationSchema, OperationSafety:
		// The frozen projection is the complete normative observation.
	}
	return result
}
