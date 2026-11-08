/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/alternative_target_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"encoding/json"
	"reflect"
	"testing"
)

type targetHintReader struct {
	hints map[string][]domain.RepositoryIdentity
	calls []string
}

func (r *targetHintReader) RepositoryIdentities(url string) []domain.RepositoryIdentity {
	r.calls = append(r.calls, url)
	return r.hints[url]
}

func TestAlternativeTargetRulesPreserveSourcesAndInputs(t *testing.T) {
	ref := domain.CommitRef{Owner: "local", Repo: "captured", Commit: "immutable-bytes"}
	source := domain.EvidenceSource{Ref: ref, Path: "graph.go", Line: 9, ContentHash: "frozen", Strength: .5}
	assessment := domain.CandidateAssessment{Repository: ref, ScreeningDecision: domain.ScreeningDecision{Repository: ref.Key()}, ProblemAlignment: domain.CapabilityEvidence{Sources: []domain.EvidenceSource{source}}, Capabilities: []domain.CapabilityEvidence{{Capability: "graph", Sources: []domain.EvidenceSource{source}}}, License: domain.LicenseAssessment{Sources: []domain.EvidenceSource{source}}, Maintenance: domain.MaintenanceAssessment{Sources: []domain.EvidenceSource{source}}, Operation: domain.OperationAssessment{Sources: []domain.EvidenceSource{source}}}
	discovery := domain.DiscoverySnapshot{Candidates: []domain.DiscoveredCandidate{
		{Repository: domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit, URL: "file://frozen"}},
		{Repository: domain.RepositoryCandidate{Owner: "local", Repo: "other", URL: "file://other"}},
		{Repository: domain.RepositoryCandidate{Owner: "local", Repo: "unknown", URL: "file://unknown"}},
	}}
	assessments := []domain.CandidateAssessment{assessment}
	beforeDiscovery, _ := json.Marshal(discovery)
	beforeAssessments, _ := json.Marshal(assessments)
	reader := &targetHintReader{hints: map[string][]domain.RepositoryIdentity{"file://frozen": {{Owner: "unrelated", Repo: "remote"}, {Owner: "ORG", Repo: "TARGET"}}, "file://other": {{Owner: "other", Repo: "tool"}}}}
	target := &domain.RepositoryIdentity{Owner: " Org ", Repo: " Target "}
	got, values := usecase.AlternativeInputsForTarget(domain.IdeaRunLocal, discovery, assessments, target, reader)
	if !reflect.DeepEqual(reader.calls, []string{"file://frozen", "file://other", "file://unknown"}) {
		t.Fatalf("lookup order: %v", reader.calls)
	}
	if got.Candidates[0].Repository.Owner != "org" || got.Candidates[0].Repository.Repo != "target" || got.Candidates[0].Repository.CommitID != ref.Commit {
		t.Fatal("target binding or byte identity changed")
	}
	if got.Candidates[1].Repository.Repo != "other" || got.Candidates[1].Repository.Stale || !got.Candidates[2].Repository.Stale {
		t.Fatal("hints changed unrelated/unknown classification")
	}
	a := values[0]
	if a.Repository.Owner != "org" || a.Repository.Repo != "target" || a.Repository.Commit != ref.Commit || a.ScreeningDecision.Repository != a.Repository.Key() {
		t.Fatal("assessment binding diverged")
	}
	for _, sources := range [][]domain.EvidenceSource{a.ProblemAlignment.Sources, a.Capabilities[0].Sources, a.License.Sources, a.Maintenance.Sources, a.Operation.Sources} {
		want := source
		want.Ref.Owner = "org"
		want.Ref.Repo = "target"
		if !reflect.DeepEqual(sources, []domain.EvidenceSource{want}) {
			t.Fatalf("source provenance changed: %+v", sources)
		}
	}
	afterDiscovery, _ := json.Marshal(discovery)
	afterAssessments, _ := json.Marshal(assessments)
	if string(beforeDiscovery) != string(afterDiscovery) || string(beforeAssessments) != string(afterAssessments) {
		t.Fatal("input was mutated")
	}
	// The existing ranking owns the exclusion; a consulted hint is not new evidence.
	report, err := usecase.RankSoftwareAlternatives(domain.SoftwareAlternativeTarget{Name: "Target", Repository: target}, domain.IdeaSpec{Capabilities: []string{"graph"}}, got, values, usecase.DefaultRecommendationPolicyV1())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, excluded := range report.Excluded {
		if excluded.Repository.Owner == "org" && excluded.Repository.Repo == "target" {
			found = true
			if excluded.Fit != domain.AlternativeNotAlternative {
				t.Fatal("target became eligible")
			}
		}
	}
	if !found {
		t.Fatal("target exclusion missing")
	}
}

func TestAlternativeTargetDoesNotReadHintsOutsideLocalTarget(t *testing.T) {
	discovery := domain.DiscoverySnapshot{Candidates: []domain.DiscoveredCandidate{{Repository: domain.RepositoryCandidate{Owner: "unchanged", Repo: "repo"}}}}
	assessments := []domain.CandidateAssessment{{Repository: domain.CommitRef{Owner: "unchanged", Repo: "repo"}}}
	for _, mode := range []domain.IdeaRunMode{domain.IdeaRunLive, domain.IdeaRunReplay, domain.IdeaRunLocal} {
		target := &domain.RepositoryIdentity{Owner: "org", Repo: "target"}
		if mode == domain.IdeaRunLocal {
			target = nil
		}
		got, values := usecase.AlternativeInputsForTarget(mode, discovery, assessments, target, nil)
		if !reflect.DeepEqual(got, discovery) || !reflect.DeepEqual(values, assessments) {
			t.Fatal("unrelated mode changed inputs")
		}
	}
}
