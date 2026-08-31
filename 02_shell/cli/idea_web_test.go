/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type fakeObservableSearch struct{}

func (fakeObservableSearch) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "fake", PreservesRank: true, MaxResultsLimit: 5}
}

type recordingSearch struct{ order []string }

func (r *recordingSearch) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "recording", PreservesRank: true, MaxResultsLimit: 5}
}
func (r *recordingSearch) SearchOne(_ context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	r.order = append(r.order, req.RequestID)
	resp := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "recording", RetrievedAt: "2026-08-26T00:00:00Z", Status: "SUCCESS", Results: []domain.RawWebResult{{Rank: 1, URL: "https://github.com/acme/" + req.QueryID, Title: req.QueryID, Snippet: "result"}}}
	return usecase.SealWebResponse(req, resp), nil
}
func (fakeObservableSearch) SearchOne(_ context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	r := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "fake", RetrievedAt: "2026-08-26T00:00:00Z", Status: "SUCCESS", Results: []domain.RawWebResult{{Rank: 1, URL: "https://github.com/acme/tool", Title: "acme/tool", Snippet: "tool"}}}
	return usecase.SealWebResponse(req, r), nil
}

func TestIdeaWebCommands(t *testing.T) {
	dir := t.TempDir()
	q := filepath.Join(dir, "queries.json")
	s := filepath.Join(dir, "snapshot.json")
	content := `{"schema_version":"2","grammar_hash":"g","queries":[{"id":"q1","text":"find tools"}]}`
	if e := os.WriteFile(q, []byte(content), 0644); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := runIdeaWeb([]string{"discover-observable", "--queries", q, "--output", s}, &out, fakeObservableSearch{}); e != nil {
		t.Fatal(e)
	}
	if e := runIdeaWeb([]string{"replay", "--snapshot", s}, &out, nil); e != nil {
		t.Fatal(e)
	}
}
func TestDiscoverRequiresInjectedProvider(t *testing.T) {
	var out bytes.Buffer
	e := runIdeaWeb([]string{"discover-observable", "--queries", "missing"}, &out, nil)
	if e == nil {
		t.Fatal("missing provider accepted")
	}
}

func TestExperimentRequiresNoTokenAndInterleavesArms(t *testing.T) {
	dir := t.TempDir()
	paths := map[string]string{}
	for _, arm := range []string{"N", "F1", "F2"} {
		exp := domain.FunctionalGrammarExpansion{SchemaVersion: "2"}
		for i := 0; i < 12; i++ {
			exp.Queries = append(exp.Queries, domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: fmt.Sprintf("%s-%02d", arm, i), Text: "query"}})
		}
		b, _ := json.Marshal(exp)
		p := filepath.Join(dir, arm+".json")
		if err := os.WriteFile(p, b, 0644); err != nil {
			t.Fatal(err)
		}
		paths[arm] = p
	}
	provider := &recordingSearch{}
	args := []string{"experiment", "--nominal", paths["N"], "--functional-v1", paths["F1"], "--functional-v2", paths["F2"], "--output-dir", filepath.Join(dir, "out"), "--no-token"}
	if err := runIdeaWeb(args, &bytes.Buffer{}, provider); err != nil {
		t.Fatal(err)
	}
	want := []string{"0019-01-N", "0019-01-F1", "0019-01-F2", "0019-02-N"}
	for i := range want {
		if provider.order[i] != want[i] {
			t.Fatalf("order[%d]=%s", i, provider.order[i])
		}
	}
	args = args[:len(args)-1]
	if err := runIdeaWeb(args, &bytes.Buffer{}, provider); err == nil {
		t.Fatal("experiment accepted token-capable mode")
	}
}
