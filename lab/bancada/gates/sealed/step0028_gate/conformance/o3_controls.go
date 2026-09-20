/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_controls.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import "sort"

func O3PositiveCases() []Case {
	tools := []Tool{
		{Name: "bateia_search", Definition: "sha256:t1"},
		{Name: "bateia_benchmark", Definition: "sha256:t2"},
		{Name: "bateia_recommend", Definition: "sha256:t3"},
		{Name: "bateia_start_recommendation", Definition: "sha256:t4"},
		{Name: "bateia_get_run", Definition: "sha256:t5"},
		{Name: "bateia_get_artifact", Definition: "sha256:t6"},
		{Name: "bateia_cancel_run", Definition: "sha256:t7"},
		{Name: "bateia_list_connectors", Definition: "sha256:t8"},
		{Name: "bateia_start_market_discovery", Definition: "sha256:t9"},
	}
	refinements := map[string]string{"0025": "receipt-0025", "0026": "receipt-0026", "0027": "receipt-0027"}
	cases := []Case{
		{ID: "P01", Obligation: P01ExtractionBinding, Operation: OperationExtract, Input: Input{
			ScenarioVersion: "O3", Implementation: true, Independent: true, SameImmutableSource: true, ExactBundleBinding: true, ContentHashValid: true, ParentStatus: "PROVEN", ClassInBaseline: true,
			Owner: "acme", Repo: "append-log", Commit: "0123456789abcdef", Path: "core/log.go", Language: "Go", Ecosystem: "go", SourceHashes: []string{"sha256:impl", "sha256:test"}, SourceKinds: []string{"IMPLEMENTATION", "TEST"}, RuleIDs: []string{"append-record", "recovery-replay"}, WitnessRoles: []string{"IMPLEMENTATION", "INDEPENDENT"}, SourceRefs: []string{"acme/append-log@0123456789abcdef:core/log.go", "acme/append-log@0123456789abcdef:core/log_test.go"}, MechanismIDs: []string{"mechanism-append-log"}, AppliedLimits: map[string]int{"evidence": 2},
		}},
		{ID: "P02", Obligation: P02EpistemicSeparation, Operation: OperationExtract, Input: Input{
			ScenarioVersion: "O3", DocumentationOnly: true, SameImmutableSource: true, ExactBundleBinding: true, ContentHashValid: true, ParentStatus: "PARTIAL", Owner: "acme", Repo: "docs-only", Commit: "abcdef0123456789", Path: "README.md", Language: "Markdown", Ecosystem: "docs", SourceHashes: []string{"sha256:docs"}, SourceKinds: []string{"DOCUMENTATION"}, RuleIDs: []string{"documented-claim"}, WitnessRoles: []string{"DOCUMENTATION"}, SourceRefs: []string{"acme/docs-only@abcdef0123456789:README.md"},
		}},
		{ID: "P03", Obligation: P03ProductOnly, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", ProductOnly: true, TargetState: "UNKNOWN", RequestedTransfer: "PATTERN_ADAPTATION", License: "COMPATIBLE", MechanismCount: 1, ClaimedNovelty: "UNKNOWN", SourceRefs: []string{"product:observed-feature", "source-run:sha256:source"}, ArtifactNames: []string{"comparative-synthesis.json", "improvement-hypotheses.json", "transfer-plans.json"}, ArtifactEntries: []string{"comparison:observation", "hypothesis:proposed", "plan:falsifiable"}, ArtifactsPublished: true, ArtifactsCanonical: true, ValidationPlan: []string{"measure baseline", "apply bounded experiment", "compare invariant"}, Assumptions: []string{"product observation is not mechanism evidence"}, Preconditions: []string{"target remains unknown"}, Risks: []string{"false analogy"}, Tradeoffs: []string{"coverage versus certainty"}, HypothesisIDs: []string{"hyp-product-observation"},
		}},
		{ID: "P04", Obligation: P04UnknownTarget, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", TargetState: "UNKNOWN", MechanismStatus: "PROVEN", ClaimedNovelty: "CONVENTIONAL", RequestedTransfer: "PATTERN_ADAPTATION", License: "COMPATIBLE", MechanismCount: 1, SourceRefs: []string{"mechanism:sha256:m1", "target:sha256:t1"}, TargetGapRefs: []string{"gap:target-not-observed"}, ValidationPlan: []string{"instrument target", "collect immutable snapshot", "reclassify only with evidence"}, HypothesisIDs: []string{"hyp-unknown-target"},
		}},
		{ID: "P05", Obligation: P05License, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", TargetState: "PRESENT", MechanismStatus: "PROVEN", ClaimedNovelty: "CONVENTIONAL", License: "UNRECOGNIZED", RequestedTransfer: "COMPONENT_REUSE", MechanismCount: 1, SourceRefs: []string{"mechanism:sha256:licensed", "license:sha256:unknown"}, ValidationPlan: []string{"resolve license provenance", "permit reuse only after compatibility"}, HypothesisIDs: []string{"hyp-license-ceiling"},
		}},
		{ID: "P06", Obligation: P06Composition, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", TargetState: "PRESENT", MechanismStatus: "PROVEN", ClaimedNovelty: "DIVERGENT", License: "COMPATIBLE", RequestedTransfer: "PATTERN_ADAPTATION", MechanismCount: 2, CompatibilityWitnessed: false, MechanismRefs: []string{"mechanism:queue", "mechanism:ledger"}, SourceRefs: []string{"source:queue@commit", "source:ledger@commit"}, ValidationPlan: []string{"test interface model", "test state and failure models", "measure performance"}, HypothesisIDs: []string{"hyp-composition"},
		}},
		{ID: "P07", Obligation: P07CounterEvidence, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", TargetState: "PRESENT", MechanismStatus: "PROVEN", ClaimedNovelty: "CONVENTIONAL", License: "COMPATIBLE", RequestedTransfer: "PATTERN_ADAPTATION", MechanismCount: 1, CounterSearch: "CENTRAL_CONFLICT", CounterScope: []string{"central invariant", "recovery boundary"}, CounterEvidenceRefs: []string{"counter:sha256:central"}, SourceRefs: []string{"mechanism:sha256:m7"}, ValidationPlan: []string{"reproduce counterexample", "reject unsupported causal link"}, HypothesisIDs: []string{"hyp-counterevidence"},
		}},
		{ID: "P08", Obligation: P08PopularityNovelty, Operation: OperationHypothesis, Input: Input{
			ScenarioVersion: "O3", TargetState: "PRESENT", MechanismStatus: "PROVEN", ClaimedNovelty: "DIVERGENT", License: "COMPATIBLE", Fit: "PARTIAL", RequestedTransfer: "PATTERN_ADAPTATION", MechanismCount: 1, Popularity: 999999999, AdoptionRate: 1, Ordering: []string{"mechanism-a", "mechanism-b"}, CanonicalItems: []string{"mechanism-b", "mechanism-a"}, ArtifactNames: []string{"comparative-synthesis.json"}, ArtifactEntries: []string{"mechanism-a", "mechanism-b"}, ArtifactsPublished: true, ArtifactsCanonical: true, SourceRefs: []string{"source:a", "source:b"}, ValidationPlan: []string{"compare by causal fit", "ignore popularity for support"},
		}},
		{ID: "P09", Obligation: P09Idempotency, Operation: OperationStart, Input: Input{
			ScenarioVersion: "O3", SourceExists: true, SourceTerminal: true, SourcePublished: true, SourceHashMatches: true, CanonicalDocument: true, RequestHash: "sha256:request", CanonicalItems: []string{"source-b", "source-a"}, ConcurrentStarts: 4, SourceRefs: []string{"run:source"}, AppliedLimits: map[string]int{"queue_depth": 1}, ArtifactIDs: []string{"artifact-source"},
		}},
		{ID: "P10", Obligation: P10CancelRecovery, Operation: OperationLifecycle, Input: Input{
			ScenarioVersion: "O3", Lifecycle: "RUNNING", ImprovementStatus: "PARTIAL", TargetState: "UNKNOWN", MechanismStatus: "PROVEN", CounterSearch: "NONE_FOUND_IN_SCOPE", Cancelled: true, RecoveryRequested: true, ArtifactNames: []string{"comparative-synthesis.json"}, ArtifactState: "PUBLISHED", ArtifactHashes: map[string]string{"comparative-synthesis.json": "sha256:preserved"}, AppliedLimits: map[string]int{"completed_stage": 1},
		}},
		{ID: "P11", Obligation: P11OfflineReplay, Operation: OperationReplay, Input: Input{
			ScenarioVersion: "O3", ReplayItems: []string{"capture-target", "capture-ledger", "capture-market"}, CatalogHashMatches: true, CaptureHashMatches: true, NetworkAvailable: false, ArtifactNames: []string{"comparative-synthesis.json", "improvement-hypotheses.json", "transfer-plans.json", "mechanisms.json", "gaps.json", "counterevidence.json", "provenance.json"}, ArtifactHashes: map[string]string{"comparative-synthesis.json": "sha256:a", "improvement-hypotheses.json": "sha256:b", "transfer-plans.json": "sha256:c", "mechanisms.json": "sha256:d", "gaps.json": "sha256:e", "counterevidence.json": "sha256:f", "provenance.json": "sha256:g"}, IndexHashes: map[string]string{"artifact-index": "sha256:index", "run-index": "sha256:run-index"}, SnapshotSHA256: "sha256:snapshot", SourceRefs: []string{"capture:market", "capture:ledger", "capture:target"},
		}},
		{ID: "P12", Obligation: P12RefinementRuntime, Operation: OperationFull, Input: Input{
			ScenarioVersion: "O3", StoreIDs: []string{"store-1", "store-1"}, QueueIDs: []string{"queue-1", "queue-1"}, RootIDs: []string{"root-1", "root-1"}, SecretIDs: []string{"secret-1", "secret-1"}, WorkerIDs: []string{"worker-1", "worker-1"}, BasalTools: tools, BasalToolsSHA256: "ee630cff5eada0dee4afe82ad9513dea7b760f0beb7905e915d2714a9114a596", ActualBasalToolsSHA256: "ee630cff5eada0dee4afe82ad9513dea7b760f0beb7905e915d2714a9114a596", ToolDefinitions: toolDefinitionMap(tools), SchemaDefinitions: map[string]string{"bateia_start_improvement_synthesis": "sha256:tool10-schema"}, RefinementReceipts: refinements, RequiredRefinementHashes: refinements, ArtifactHashes: map[string]string{"refinement-0025": "receipt-0025", "refinement-0026": "receipt-0026", "refinement-0027": "receipt-0027"}, LayerFindings: []string{"L1_PURE", "L2_NO_L3", "L4_WIRING_ONLY"},
		}},
	}
	for index := range cases {
		cases[index].Expected = frozenO3Expected(cases[index])
	}
	return cases
}

func frozenO3Expected(testCase Case) Observation {
	input := testCase.Input
	result := frozenProjection(input)
	switch testCase.ID {
	case "P01":
		result.Decision, result.MechanismStatus, result.Novelty = "MECHANISM_PRESERVED", "PROVEN", "CONVENTIONAL"
		result.BindingVerified, result.ProvenanceComplete = true, true
		result.CanonicalSHA256 = canonicalDigest([]string{input.Owner, input.Repo, input.Commit, input.Path, input.Language, input.Ecosystem})
	case "P02":
		result.Decision, result.MechanismStatus, result.Novelty = "MECHANISM_PRESERVED", "PARTIAL", "UNKNOWN"
		result.BindingVerified, result.ProvenanceComplete = true, true
		result.CanonicalSHA256 = canonicalDigest([]string{input.Owner, input.Repo, input.Commit, input.Path, input.Language, input.Ecosystem})
	case "P03":
		result.Decision, result.MechanismStatus, result.Novelty = "OBSERVATION_ONLY", "NONE", "UNKNOWN"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
		result.ArtifactSHA256 = canonicalDigest(sortedCopy(input.ArtifactEntries))
	case "P04":
		result.MechanismStatus, result.Novelty = "PROVEN", "CONVENTIONAL"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
	case "P05":
		result.MechanismStatus, result.Novelty = "PROVEN", "CONVENTIONAL"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
		result.WarningCodes = []string{"LICENSE_NOT_RECOGNIZED"}
	case "P06":
		result.MechanismStatus, result.Novelty = "PROVEN", "DIVERGENT"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
		result.IntegrationUnknowns = []string{"failure model compatibility", "interface compatibility", "license compatibility", "performance compatibility", "state model compatibility"}
	case "P07":
		result.MechanismStatus, result.Novelty = "INCONCLUSIVE", "CONVENTIONAL"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "INCONCLUSIVE", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
	case "P08":
		result.MechanismStatus, result.Novelty = "PROVEN", "DIVERGENT"
		result.HypothesisStatus, result.TransferMode, result.AdoptionDisposition = "PROPOSED", "PATTERN_ADAPTATION", "NOT_AN_ADOPTION_DECISION"
		result.ArtifactSHA256 = canonicalDigest(sortedCopy(input.ArtifactEntries))
		result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
	case "P09":
		result.ResolveCalls, result.ReserveCalls, result.IDCalls, result.CreateCalls = 1, 1, 1, 1
		result.RunID, result.Lifecycle = "run_reference", "QUEUED"
		result.RunIDs = []string{"run_reference", "run_reference", "run_reference", "run_reference"}
		result.CanonicalSHA256 = canonicalDigest(sortedCopy(input.CanonicalItems))
	case "P10":
		result.Lifecycle, result.HypothesisStatus = "CANCELLED", "PROPOSED"
		result.RecoveryState, result.AdoptionDisposition = "INTERRUPTED", "NOT_AN_ADOPTION_DECISION"
	case "P11":
		result.ReplaySHA256 = canonicalDigest(sortedCopy(input.ReplayItems))
	case "P12":
		result.RuntimeUnique, result.RefinementsPreserved, result.BasalHashVerified = true, true, true
		result.StoreIdentity, result.QueueIdentity, result.RootIdentity = "store-1", "queue-1", "root-1"
		result.SecretIdentity, result.WorkerIdentity = "secret-1", "worker-1"
		result.Tools = append([]Tool(nil), input.BasalTools...)
		result.ArtifactSHA256 = canonicalMapDigest(input.ArtifactHashes)
	}
	return result
}

func frozenProjection(input Input) Observation {
	result := emptyObservation()
	projectInputFields(&result, input)
	return result
}

func sortedCopy(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func toolDefinitionMap(tools []Tool) map[string]string {
	result := make(map[string]string, len(tools))
	for _, tool := range tools {
		result[tool.Name] = tool.Definition
	}
	return result
}
