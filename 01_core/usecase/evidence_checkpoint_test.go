/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_checkpoint_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type checkpointHashReader struct {
	hashes   map[string]string
	failures map[string]error
	calls    []string
}

func (r *checkpointHashReader) ArtifactSHA256(path string) (string, error) {
	r.calls = append(r.calls, path)
	if err := r.failures[path]; err != nil {
		return "", err
	}
	hash, ok := r.hashes[path]
	if !ok {
		return "", errors.New("unexpected artifact")
	}
	return hash, nil
}
func checkpointPolicyFixture() (domain.IdeaRunManifest, usecase.EvidenceCheckpointCriteria, *checkpointHashReader) {
	versions := map[string]string{"expansion": "1", "screening": "1", "evidence": "e1", "analyzer": "a1", "implementation_search": "i1", "discovery_query": "q1", "local_snapshot": "s1", "local_collection": "c1"}
	previousVersions := map[string]string{}
	for key, value := range versions {
		previousVersions[key] = value
	}
	previous := domain.IdeaRunManifest{SchemaVersion: "1", Status: domain.IdeaRunComplete, InputHashes: map[string]string{"evidence": "inputs", "evidence_rules": "rules", "policy": "old-policy", "target": "old-target"}, RuleVersions: previousVersions, Stages: []domain.IdeaRunStage{{Outputs: []domain.IdeaRunArtifact{
		{Path: "../../must-not-be-read", Hash: "private"},
		{Path: "discovery-snapshot.json", Hash: ""},
		{Path: "evidence/bundles.json", Hash: "bundles"},
		{Path: "discovery-snapshot.json", Hash: "snapshot"},
		{Path: "evidence/bundles.json", Hash: "ignored-later-duplicate"},
	}}}}
	criteria := usecase.EvidenceCheckpointCriteria{EvidenceHash: "inputs", EvidenceRulesHash: "rules", RuleVersions: versions}
	reader := &checkpointHashReader{hashes: map[string]string{"discovery-snapshot.json": "snapshot", "evidence/bundles.json": "bundles"}, failures: map[string]error{}}
	return previous, criteria, reader
}
func TestEvidenceCheckpointStatesAndReadOrder(t *testing.T) {
	for _, state := range []domain.IdeaRunStatus{domain.IdeaRunComplete, domain.IdeaRunPartial, domain.IdeaRunInconclusive, domain.IdeaRunRunning, domain.IdeaRunFailed, "", "UNKNOWN"} {
		t.Run(string(state), func(t *testing.T) {
			previous, criteria, reader := checkpointPolicyFixture()
			previous.Status = state
			before, _ := json.Marshal(previous)
			got := usecase.CanReuseEvidenceCheckpoint(previous, criteria, reader)
			want := state == domain.IdeaRunComplete || state == domain.IdeaRunPartial || state == domain.IdeaRunInconclusive
			if got != want {
				t.Fatalf("reused=%v want=%v", got, want)
			}
			if want && !reflect.DeepEqual(reader.calls, []string{"evidence/bundles.json", "discovery-snapshot.json"}) {
				t.Fatalf("logical read order/filter changed: %v", reader.calls)
			}
			if !want && len(reader.calls) != 0 {
				t.Fatal("ineligible state performed reads")
			}
			after, _ := json.Marshal(previous)
			if string(before) != string(after) {
				t.Fatal("manifest was mutated")
			}
		})
	}
}
func TestEvidenceCheckpointRejectsCausalMismatchBeforeReads(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*domain.IdeaRunManifest, *usecase.EvidenceCheckpointCriteria)
	}{
		{"schema", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria) { p.SchemaVersion = "0" }},
		{"evidence-inputs", func(_ *domain.IdeaRunManifest, c *usecase.EvidenceCheckpointCriteria) { c.EvidenceHash = "changed" }},
		{"evidence-rules", func(_ *domain.IdeaRunManifest, c *usecase.EvidenceCheckpointCriteria) {
			c.EvidenceRulesHash = "changed"
		}},
	}
	for _, key := range []string{"expansion", "screening", "evidence", "analyzer", "implementation_search", "discovery_query", "local_snapshot", "local_collection"} {
		key := key
		cases = append(cases, struct {
			name   string
			mutate func(*domain.IdeaRunManifest, *usecase.EvidenceCheckpointCriteria)
		}{"changed-" + key, func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria) { p.RuleVersions[key] = "old" }})
		cases = append(cases, struct {
			name   string
			mutate func(*domain.IdeaRunManifest, *usecase.EvidenceCheckpointCriteria)
		}{"missing-" + key, func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria) { delete(p.RuleVersions, key) }})
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			p, c, r := checkpointPolicyFixture()
			test.mutate(&p, &c)
			if usecase.CanReuseEvidenceCheckpoint(p, c, r) || len(r.calls) != 0 {
				t.Fatalf("causal mismatch reused or read artifacts: %v", r.calls)
			}
		})
	}
}
func TestEvidenceCheckpointPreservesModelPresenceAndDecisionOnlyChanges(t *testing.T) {
	for _, test := range []struct {
		name         string
		oldPresent   bool
		old, current string
		want         bool
	}{
		{"legacy-absent", false, "", "", true},
		{"empty-is-present", true, "", "", false},
		{"removed", true, "m", "", false},
		{"added", false, "", "m", false},
		{"same", true, "m", "m", true},
		{"changed", true, "old", "new", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			p, c, r := checkpointPolicyFixture()
			if test.oldPresent {
				p.InputHashes["screening_model"] = test.old
			}
			c.ScreeningModelHash = test.current
			got := usecase.CanReuseEvidenceCheckpoint(p, c, r)
			if got != test.want {
				t.Fatalf("reused=%v want=%v", got, test.want)
			}
			if !test.want && len(r.calls) != 0 {
				t.Fatal("model mismatch performed reads")
			}
		})
	}
	p, c, r := checkpointPolicyFixture()
	p.InputHashes["policy"] = "changed"
	p.InputHashes["target"] = "changed"
	p.RuleVersions["decision"] = "old"
	c.RuleVersions["decision"] = "new"
	p.RuleVersions["alternatives"] = "old"
	c.RuleVersions["alternatives"] = "new"
	if !usecase.CanReuseEvidenceCheckpoint(p, c, r) {
		t.Fatal("decision-only change invalidated compatible evidence")
	}
}
func TestEvidenceCheckpointRejectsUnverifiedArtifactsWithoutFallbackToDuplicates(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*domain.IdeaRunManifest, *checkpointHashReader)
		reads  int
	}{
		{"first-required-mismatch", func(_ *domain.IdeaRunManifest, r *checkpointHashReader) {
			r.hashes["evidence/bundles.json"] = "tampered"
		}, 1},
		{"second-required-mismatch", func(_ *domain.IdeaRunManifest, r *checkpointHashReader) {
			r.hashes["discovery-snapshot.json"] = "tampered"
		}, 2},
		{"read-error", func(_ *domain.IdeaRunManifest, r *checkpointHashReader) {
			r.failures["evidence/bundles.json"] = errors.New("unavailable")
		}, 1},
		{"missing-output", func(p *domain.IdeaRunManifest, _ *checkpointHashReader) {
			p.Stages[0].Outputs = []domain.IdeaRunArtifact{{Path: "evidence/bundles.json", Hash: "bundles"}}
		}, 1},
		{"empty-outputs", func(p *domain.IdeaRunManifest, _ *checkpointHashReader) { p.Stages = nil }, 0},
		{"bad-first-duplicate", func(p *domain.IdeaRunManifest, _ *checkpointHashReader) {
			p.Stages[0].Outputs[2].Hash = "wrong-first"
			p.Stages[0].Outputs[4].Hash = "bundles"
		}, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			p, c, r := checkpointPolicyFixture()
			test.mutate(&p, r)
			if usecase.CanReuseEvidenceCheckpoint(p, c, r) || len(r.calls) != test.reads {
				t.Fatalf("reused invalid artifact or changed short-circuit: %v", r.calls)
			}
		})
	}
	p, c, _ := checkpointPolicyFixture()
	if usecase.CanReuseEvidenceCheckpoint(p, c, nil) {
		t.Fatal("missing reader accepted")
	}
}
