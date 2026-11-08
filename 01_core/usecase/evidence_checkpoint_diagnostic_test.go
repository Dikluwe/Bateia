/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_checkpoint_diagnostic_test.md
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
	"strings"
	"testing"
)

func TestCheckpointDiagnosticReasonsAndReadBoundaries(t *testing.T) {
	for _, test := range []struct {
		name    string
		change  func(*domain.IdeaRunManifest, *usecase.EvidenceCheckpointCriteria, *checkpointHashReader)
		reason  usecase.EvidenceCheckpointReason
		subject string
		reads   int
	}{
		{"eligible", func(_ *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {}, usecase.CheckpointReusable, "", 2},
		{"status", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.Status = domain.IdeaRunFailed
		}, usecase.CheckpointStatusIneligible, "status", 0},
		{"schema", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.SchemaVersion = "secret-value"
		}, usecase.CheckpointSchemaUnsupported, "schema_version", 0},
		{"input", func(_ *domain.IdeaRunManifest, c *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			c.EvidenceHash = "secret-value"
		}, usecase.CheckpointEvidenceChanged, "evidence", 0},
		{"rules", func(_ *domain.IdeaRunManifest, c *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			c.EvidenceRulesHash = "secret-value"
		}, usecase.CheckpointRulesChanged, "evidence_rules", 0},
		{"model", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.InputHashes["screening_model"] = ""
		}, usecase.CheckpointModelChanged, "screening_model", 0},
		{"first-hash", func(_ *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, r *checkpointHashReader) {
			r.hashes["evidence/bundles.json"] = "secret-value"
		}, usecase.CheckpointArtifactMismatch, "evidence/bundles.json", 1},
		{"second-hash", func(_ *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, r *checkpointHashReader) {
			r.hashes["discovery-snapshot.json"] = "secret-value"
		}, usecase.CheckpointArtifactMismatch, "discovery-snapshot.json", 2},
		{"read-error", func(_ *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, r *checkpointHashReader) {
			r.failures["evidence/bundles.json"] = errors.New("secret-value /private/root")
		}, usecase.CheckpointArtifactUnreadable, "evidence/bundles.json", 1},
		{"record-missing", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.Stages[0].Outputs = []domain.IdeaRunArtifact{{Path: "evidence/bundles.json", Hash: "bundles"}}
		}, usecase.CheckpointArtifactUnrecorded, "discovery-snapshot.json", 1},
		{"both-records-missing", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.Stages = nil
		}, usecase.CheckpointArtifactUnrecorded, "discovery-snapshot.json", 0},
		{"empty-hash", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.Stages[0].Outputs = []domain.IdeaRunArtifact{{Path: "discovery-snapshot.json", Hash: "snapshot"}, {Path: "evidence/bundles.json"}}
		}, usecase.CheckpointArtifactUnrecorded, "evidence/bundles.json", 1},
		{"no-fallback-to-duplicate", func(p *domain.IdeaRunManifest, _ *usecase.EvidenceCheckpointCriteria, _ *checkpointHashReader) {
			p.Stages[0].Outputs[2].Hash = "secret-value"
			p.Stages[0].Outputs[4].Hash = "bundles"
		}, usecase.CheckpointArtifactMismatch, "evidence/bundles.json", 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			previous, criteria, reader := checkpointPolicyFixture()
			test.change(&previous, &criteria, reader)
			before, _ := json.Marshal([]any{previous, criteria})
			got := usecase.EvaluateEvidenceCheckpoint(previous, criteria, reader)
			if got.Reason != test.reason || got.Subject != test.subject || got.Reusable != (test.reason == usecase.CheckpointReusable) || len(reader.calls) != test.reads {
				t.Fatalf("assessment=%+v reads=%v", got, reader.calls)
			}
			calls := append([]string(nil), reader.calls...)
			reader.calls = nil
			if usecase.CanReuseEvidenceCheckpoint(previous, criteria, reader) != got.Reusable || !reflect.DeepEqual(calls, reader.calls) {
				t.Fatal("legacy predicate changed eligibility or read sequence")
			}
			again := usecase.EvaluateEvidenceCheckpoint(previous, criteria, reader)
			after, _ := json.Marshal([]any{previous, criteria})
			if got != again || string(before) != string(after) {
				t.Fatal("nondeterministic assessment or mutated inputs")
			}
			encoded, _ := json.Marshal(got)
			if strings.Contains(string(encoded), "secret-value") || strings.Contains(string(encoded), "/private/") {
				t.Fatalf("diagnostic exposed input or reader error: %s", encoded)
			}
		})
	}
}

func TestCheckpointDiagnosticPrecedence(t *testing.T) {
	previous, criteria, reader := checkpointPolicyFixture()
	previous.Status, previous.SchemaVersion = domain.IdeaRunRunning, "0"
	previous.InputHashes["evidence"], previous.InputHashes["evidence_rules"] = "old", "old"
	previous.InputHashes["screening_model"] = "old"
	previous.RuleVersions["expansion"], previous.RuleVersions["screening"] = "old", "old"
	checks := []struct {
		reason usecase.EvidenceCheckpointReason
		fix    func()
	}{
		{usecase.CheckpointStatusIneligible, func() { previous.Status = domain.IdeaRunComplete }},
		{usecase.CheckpointSchemaUnsupported, func() { previous.SchemaVersion = "1" }},
		{usecase.CheckpointEvidenceChanged, func() { previous.InputHashes["evidence"] = criteria.EvidenceHash }},
		{usecase.CheckpointRulesChanged, func() { previous.InputHashes["evidence_rules"] = criteria.EvidenceRulesHash }},
		{usecase.CheckpointModelChanged, func() { delete(previous.InputHashes, "screening_model") }},
		{usecase.CheckpointVersionChanged, func() { previous.RuleVersions["expansion"] = criteria.RuleVersions["expansion"] }},
		{usecase.CheckpointVersionChanged, func() { previous.RuleVersions["screening"] = criteria.RuleVersions["screening"] }},
	}
	for _, check := range checks {
		if got := usecase.EvaluateEvidenceCheckpoint(previous, criteria, reader); got.Reason != check.reason || len(reader.calls) != 0 {
			t.Fatalf("wrong first impediment: %+v reads=%v", got, reader.calls)
		}
		check.fix()
	}
	if got := usecase.EvaluateEvidenceCheckpoint(previous, criteria, nil); got.Reason != usecase.CheckpointReaderUnavailable || got.Subject != "" {
		t.Fatalf("reader unavailable: %+v", got)
	}
	for _, version := range []string{"expansion", "screening", "evidence", "analyzer", "implementation_search", "discovery_query", "local_snapshot", "local_collection"} {
		p, c, r := checkpointPolicyFixture()
		delete(p.RuleVersions, version)
		if got := usecase.EvaluateEvidenceCheckpoint(p, c, r); got.Reason != usecase.CheckpointVersionChanged || got.Subject != version || len(r.calls) != 0 {
			t.Fatalf("version %s: %+v", version, got)
		}
	}
}
