/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_extended_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"testing"
)

func TestStep0028O3ExtendedSubjectsExposeCollisionPaginationArtifactStateAndSchemas(t *testing.T) {
	cases := []Case{
		{ID: "O3-ID-COLLISION", Obligation: P09Idempotency, Operation: OperationStart, Input: Input{ScenarioVersion: "O3", SourceExists: true, SourceTerminal: true, SourcePublished: true, SourceHashMatches: true, CanonicalDocument: true, RequestHash: "sha256:collision", ForceCollision: true, MechanismIDs: []string{"mechanism-collision"}, AppliedLimits: map[string]int{"id_attempts": 1}}},
		{ID: "O3-PAGINATION-SNAPSHOT", Obligation: P12RefinementRuntime, Operation: OperationPagination, Input: Input{ScenarioVersion: "O3", Cursor: "cursor-1", NextCursor: "cursor-2", SnapshotSHA256: "sha256:snapshot", PageItems: []string{"artifact-a", "artifact-b"}, AppliedLimits: map[string]int{"page_size": 2}}},
		{ID: "O3-ARTIFACT-PUBLISHED", Obligation: P03ProductOnly, Operation: OperationArtifact, Input: Input{ScenarioVersion: "O3", Lifecycle: "SUCCEEDED", ArtifactState: "PUBLISHED", ArtifactTransitions: []string{"STAGED", "PUBLISHED"}, ArtifactNames: []string{"hypotheses.json"}, ArtifactHashes: map[string]string{"hypotheses.json": "sha256:hypotheses"}, ArtifactIDs: []string{"artifact-hypotheses"}}},
		{ID: "O3-ARTIFACT-FAILED", Obligation: P10CancelRecovery, Operation: OperationArtifact, Input: Input{ScenarioVersion: "O3", Lifecycle: "FAILED", ArtifactState: "FAILED", ArtifactTransitions: []string{"STAGED", "FAILED"}, ArtifactNames: []string{"partial.json"}, ArtifactHashes: map[string]string{"partial.json": "sha256:partial"}}},
		{ID: "O3-SCHEMA-STRICT", Obligation: P12RefinementRuntime, Operation: OperationSchema, Input: Input{ScenarioVersion: "O3", ToolDefinitions: map[string]string{"bateia_start_improvement_synthesis": "sha256:definition"}, SchemaDefinitions: map[string]string{"input": "sha256:input-schema", "output": "sha256:output-schema"}, Language: "JSON Schema", Ecosystem: "MCP"}},
	}
	for index := range cases {
		cases[index].Expected = frozenProjection(cases[index].Input)
	}
	cases[0].Expected.ResolveCalls, cases[0].Expected.ReserveCalls, cases[0].Expected.IDCalls = 1, 1, 1
	cases[0].Expected.CreateCalls, cases[0].Expected.CollisionDetected, cases[0].Expected.ErrorCode = 0, true, "ID_COLLISION"
	cases[2].Expected.Lifecycle, cases[2].Expected.ArtifactSHA256 = "SUCCEEDED", canonicalMapDigest(cases[2].Input.ArtifactHashes)
	cases[3].Expected.Lifecycle, cases[3].Expected.ArtifactSHA256 = "FAILED", canonicalMapDigest(cases[3].Input.ArtifactHashes)

	findings := Evaluate(context.Background(), ReferenceSubject{}, cases)
	for _, finding := range findings {
		if finding.Outcome != "Preserved" {
			t.Fatalf("extended subject %s failed: %s", finding.CaseID, finding.Detail)
		}
	}
}
