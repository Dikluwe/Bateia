//go:build step0028_candidate

/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/architectural_mechanism_step0028_oracle_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package domain_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type step0028ExtractorSignatureProbe struct{}

func (step0028ExtractorSignatureProbe) ExtractMechanisms(context.Context, domain.MechanismExtractionInput) (domain.MechanismExtractionResult, error) {
	return domain.MechanismExtractionResult{}, nil
}

type step0028StarterSignatureProbe struct{}

func (step0028StarterSignatureProbe) StartImprovementSynthesis(context.Context, string, domain.ImprovementSynthesisRequest) (domain.ImprovementSynthesisStartReceipt, error) {
	return domain.ImprovementSynthesisStartReceipt{}, nil
}

type step0028ExecutorSignatureProbe struct{}

func (step0028ExecutorSignatureProbe) ExecuteImprovementSynthesis(context.Context, domain.RunID) error {
	return nil
}

var _ domain.MechanismExtractor = step0028ExtractorSignatureProbe{}
var _ domain.ImprovementSynthesisStarter = step0028StarterSignatureProbe{}
var _ domain.ImprovementSynthesisExecutor = step0028ExecutorSignatureProbe{}

func TestStep0028OraclePublicEnumsAndBasalLicenseTypeDoNotCollide(t *testing.T) {
	want := map[any]string{
		domain.MechanismProven: "PROVEN", domain.MechanismPartial: "PARTIAL", domain.MechanismInconclusive: "INCONCLUSIVE",
		domain.NoveltyConventional: "CONVENTIONAL", domain.NoveltyDivergent: "DIVERGENT", domain.NoveltyUnknown: "UNKNOWN",
		domain.WitnessImplementation: "IMPLEMENTATION", domain.WitnessTest: "TEST", domain.WitnessExample: "EXAMPLE",
		domain.WitnessManifest: "MANIFEST", domain.WitnessSchema: "SCHEMA", domain.WitnessConfiguration: "CONFIGURATION",
		domain.WitnessDocumentation: "DOCUMENTATION", domain.WitnessCounter: "COUNTEREVIDENCE",
		domain.HypothesisProposed: "PROPOSED", domain.HypothesisSupported: "SUPPORTED", domain.HypothesisConflicting: "CONFLICTING", domain.HypothesisInconclusive: "INCONCLUSIVE",
		domain.ImprovementComplete: "COMPLETE", domain.ImprovementPartial: "PARTIAL", domain.ImprovementInconclusive: "INCONCLUSIVE",
		domain.TransferComponentReuse: "COMPONENT_REUSE", domain.TransferPatternAdapt: "PATTERN_ADAPTATION", domain.TransferReferenceOnly: "REFERENCE_ONLY", domain.TransferUnknown: "UNKNOWN",
		domain.DerivationRuleDerived: "RULE_DERIVED", domain.DerivationExternal: "EXTERNAL_AGENT_PROPOSAL",
		domain.TargetPresent: "PRESENT", domain.TargetAbsent: "ABSENT", domain.TargetPartial: "PARTIAL", domain.TargetUnknown: "UNKNOWN",
		domain.ImprovementLicenseCompatible: "COMPATIBLE", domain.ImprovementLicenseIncompatible: "INCOMPATIBLE", domain.ImprovementLicenseUnrecognized: "UNRECOGNIZED",
		domain.CounterEvidenceFound: "FOUND", domain.CounterEvidenceNoneInScope: "NONE_FOUND_IN_SCOPE", domain.CounterEvidenceBudgetExhausted: "BUDGET_EXHAUSTED",
	}
	for value, expected := range want {
		if got := reflect.ValueOf(value).String(); got != expected {
			t.Fatalf("enum %T = %q, want %q", value, got, expected)
		}
	}

	// R3 must preserve the already-certified candidate-license struct.
	basal := domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT"}
	if basal.Identifier != "MIT" {
		t.Fatal("basal LicenseAssessment was shadowed by the synthesis enum")
	}
}

func TestStep0028OracleSafeMechanismJSONHasExactNonNullSurface(t *testing.T) {
	mechanism := domain.ArchitecturalMechanism{
		ID:             "mech_1111111111111111111111111111111111111111111111111111111111111111",
		MechanismClass: "APPEND_ONLY_LOG_RECOVERY", Name: "append-only log with explicit recovery",
		Purpose: "preserve history", Trigger: "interruption", Strategy: "append and replay",
		Components: []domain.MechanismComponentRef{}, DataFlow: []string{}, StateModel: []string{}, FailureModel: []string{},
		Boundaries: []string{"immutable repository"}, Preconditions: []string{}, Tradeoffs: []string{},
		Source:          domain.ImmutableCandidateRef{Repository: "atlas/ledger", Commit: strings.Repeat("1", 40)},
		Evidence:        []domain.MechanismWitness{{Repository: "atlas/ledger", Commit: strings.Repeat("1", 40), Path: "ledger/log.go", Symbol: "Append", Line: 41, SourceKind: domain.EvidenceSymbol, Kind: domain.WitnessImplementation, AssessmentStatus: domain.EvidenceProven, RuleID: "append-only-write", ContentHash: "sha256:" + strings.Repeat("a", 64)}},
		CounterEvidence: []domain.MechanismWitness{}, Status: domain.MechanismProven, Novelty: domain.NoveltyConventional,
		ExtractionRuleIDs: []string{"MR-APPEND-ONLY-LOG-V1"},
	}
	encoded, err := json.Marshal(mechanism)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(encoded, &object); err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{"boundaries", "components", "counter_evidence", "data_flow", "evidence", "extraction_rule_ids", "failure_model", "id", "mechanism_class", "name", "novelty", "preconditions", "purpose", "source", "state_model", "status", "strategy", "tradeoffs", "trigger"}
	for _, key := range wantKeys {
		if _, ok := object[key]; !ok {
			t.Fatalf("mechanism JSON missing %q: %s", key, encoded)
		}
	}
	for _, key := range []string{"components", "data_flow", "state_model", "failure_model", "preconditions", "tradeoffs", "counter_evidence"} {
		if object[key] == nil {
			t.Fatalf("mechanism JSON emitted null %q: %s", key, encoded)
		}
	}
	for _, forbidden := range []string{"snippet", "content", "readme", "checkout", "url"} {
		if strings.Contains(strings.ToLower(string(encoded)), `"`+forbidden+`"`) {
			t.Fatalf("unsafe field %q in mechanism JSON: %s", forbidden, encoded)
		}
	}
}

func TestStep0028OracleHypothesisJSONPreservesUnknownsAndNoAdoptionDecision(t *testing.T) {
	hypothesis := domain.ImprovementHypothesis{
		ID:    "hyp_2222222222222222222222222222222222222222222222222222222222222222",
		Title: "combine bounded admission with recovery", TargetProperty: "durable bounded ingestion", TargetState: domain.TargetUnknown,
		TargetGapRefs: []string{}, ExpectedEffect: "bound work while retaining recoverability",
		MechanismRefs: []string{"mech_1", "mech_2"}, ProductObservationRefs: []string{}, TransferMode: domain.TransferPatternAdapt,
		License: domain.ImprovementLicenseUnrecognized, Assumptions: []string{"the target accepts a log"}, Preconditions: []string{},
		IntegrationUnknowns: []string{"failure model compatibility", "interface compatibility", "license compatibility", "performance compatibility", "state model compatibility"},
		Risks:               []string{"replay latency"}, Tradeoffs: []string{}, EvidenceRefs: []string{"artifact:mechanisms"}, CounterEvidenceRefs: []string{},
		CounterEvidenceSearch: domain.CounterEvidenceSearch{Status: domain.CounterEvidenceNoneInScope, ScopeRefs: []string{"artifact:mechanisms"}, ResultRefs: []string{}},
		ValidationPlan:        []domain.ValidationStep{{Question: "does replay remain bounded?", Observable: "replay duration", SafeProcedure: "offline fixture replay", PassCriterion: "within frozen bound", RejectCriterion: "exceeds frozen bound", RelativeCost: "LOW", MissingEvidence: "target baseline"}},
		Derivation:            domain.DerivationRuleDerived, Status: domain.HypothesisProposed, AdoptionDisposition: "NOT_AN_ADOPTION_DECISION",
	}
	encoded, err := json.Marshal(hypothesis)
	if err != nil {
		t.Fatal(err)
	}
	for _, exact := range []string{`"target_state":"UNKNOWN"`, `"license":"UNRECOGNIZED"`, `"status":"PROPOSED"`, `"adoption_disposition":"NOT_AN_ADOPTION_DECISION"`, `"integration_unknowns":[`} {
		if !strings.Contains(string(encoded), exact) {
			t.Fatalf("hypothesis JSON lacks %s: %s", exact, encoded)
		}
	}
	for _, nullArray := range []string{`"target_gap_refs":null`, `"product_observation_refs":null`, `"preconditions":null`, `"tradeoffs":null`, `"counter_evidence_refs":null`} {
		if strings.Contains(string(encoded), nullArray) {
			t.Fatalf("hypothesis emitted null array %s: %s", nullArray, encoded)
		}
	}
}

func TestStep0028OracleSucceededLifecycleDoesNotImplySupportedOrComplete(t *testing.T) {
	status := domain.ImprovementPartial
	view := domain.RecommendationRunView{
		SchemaVersion: "1", RunID: "run_0123456789abcdef0123456789abcdef", RunKind: domain.RunKindImprovementSynthesis,
		Lifecycle: domain.RunSucceeded, Version: 7, Stage: domain.RunStageTerminal,
		Progress:             domain.RecommendationRunProgress{Planned: 5, Completed: 5},
		RecommendationStatus: nil, MarketDiscoveryStatus: nil, ImprovementStatus: &status,
		Artifacts: []domain.RunArtifactIndex{}, Failure: nil,
	}
	encoded, err := json.Marshal(view)
	if err != nil {
		t.Fatal(err)
	}
	for _, exact := range []string{`"run_kind":"IMPROVEMENT_SYNTHESIS"`, `"lifecycle":"SUCCEEDED"`, `"recommendation_status":null`, `"market_discovery_status":null`, `"improvement_status":"PARTIAL"`} {
		if !strings.Contains(string(encoded), exact) {
			t.Fatalf("run view lacks %s: %s", exact, encoded)
		}
	}
	if strings.Contains(string(encoded), "SUPPORTED") || strings.Contains(string(encoded), `"improvement_status":"COMPLETE"`) {
		t.Fatalf("lifecycle promoted epistemic status: %s", encoded)
	}
}
