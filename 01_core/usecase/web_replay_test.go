/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_replay_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func sealedSnapshot(arm string, urls ...string) domain.WebSearchSnapshot {
	req := domain.ObservableWebRequest{RequestID: "req-" + arm, QueryID: "q-" + arm, Query: "query", MaxResults: 5}
	resp := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "test", RetrievedAt: "2026-08-26T00:00:00Z", Status: "SUCCESS"}
	for i, u := range urls {
		resp.Results = append(resp.Results, domain.RawWebResult{Rank: i + 1, URL: u, CanonicalURL: u, Title: "title", Snippet: "snippet"})
	}
	resp = SealWebResponse(req, resp)
	s := domain.WebSearchSnapshot{SchemaVersion: "1", RunID: "run", Arm: arm, Provider: "test", Responses: []domain.RawWebResponse{resp}}
	s.SnapshotHash = WebSnapshotHash(s)
	return s
}
func TestReplayVerifiesAndDeduplicates(t *testing.T) {
	s := sealedSnapshot("F2", "https://example.com/repo/", "https://example.com/repo")
	r, e := ReplayWebSnapshot(s)
	if e != nil || !r.Replayable || r.UniqueHits != 1 {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestReplayRejectsTamperedSnapshot(t *testing.T) {
	s := sealedSnapshot("F2", "https://example.com/repo")
	s.Responses[0].Results[0].Title = "changed"
	r, _ := ReplayWebSnapshot(s)
	if r.Replayable {
		t.Fatal("tampering survived")
	}
}
func TestReplaySeparatesSponsored(t *testing.T) {
	s := sealedSnapshot("F2", "https://example.com/ad")
	s.Responses[0].Results[0].IsAd = true
	s.Responses[0] = SealWebResponse(domain.ObservableWebRequest{RequestID: "req-F2", QueryID: "q-F2", Query: "query", MaxResults: 5}, s.Responses[0])
	s.SnapshotHash = WebSnapshotHash(s)
	r, _ := ReplayWebSnapshot(s)
	if len(r.Organic) != 0 || len(r.Sponsored) != 1 {
		t.Fatal("ad counted organic")
	}
}
func TestReplayIsDeterministic(t *testing.T) {
	s := sealedSnapshot("F2", "https://example.com/a")
	a, _ := ReplayWebSnapshot(s)
	b, _ := ReplayWebSnapshot(s)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("replay changed")
	}
}
func TestCompareUsesExplicitExpectedIDs(t *testing.T) {
	n := sealedSnapshot("N", "https://example.com/noise")
	f1 := sealedSnapshot("F1", "https://github.com/acme/one")
	f2 := sealedSnapshot("F2", "https://github.com/acme/one", "https://github.com/acme/two", "https://example.com/noise")
	j := domain.BlindedJudgmentsFile{Judgments: []domain.BlindedJudgment{{CanonicalURL: "https://github.com/acme/one", Judgment: "RELEVANT", IsRepoBacked: true, ExpectedID: "seed-1"}, {CanonicalURL: "https://github.com/acme/two", Judgment: "RELEVANT", IsRepoBacked: true, ExpectedID: "seed-2"}, {CanonicalURL: "https://example.com/noise", Judgment: "IRRELEVANT"}}}
	r, e := CompareGrammarExperiments(n, f1, f2, j)
	if e != nil {
		t.Fatal(e)
	}
	if r.RepoBackedRecallF2 == nil || *r.RepoBackedRecallF2 != 1 {
		t.Fatalf("%+v", r)
	}
}
func TestCompareEmptyExpectedDenominatorIsNull(t *testing.T) {
	s := sealedSnapshot("N", "https://example.com/x")
	r, e := CompareGrammarExperiments(s, s, s, domain.BlindedJudgmentsFile{})
	if e != nil {
		t.Fatal(e)
	}
	if r.RepoBackedRecallF2 != nil || r.Verdict != "INCONCLUSIVE" {
		t.Fatalf("%+v", r)
	}
}
