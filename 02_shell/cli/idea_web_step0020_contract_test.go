/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web_step0020_contract_test.md
 * @layer L2
 * @updated 2026-08-31
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
	"reflect"
	"strings"
	"testing"
)

func step0020QueryFiles(t *testing.T, dir string) map[string]string {
	t.Helper()
	paths := map[string]string{}
	for _, arm := range []string{"N", "F1", "F2"} {
		exp := domain.FunctionalGrammarExpansion{SchemaVersion: "2", GrammarHash: "frozen-" + arm}
		for i := 1; i <= 12; i++ {
			exp.Queries = append(exp.Queries, domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: fmt.Sprintf("%s%d", arm, i), Text: "public repository query"}})
		}
		b, _ := json.Marshal(exp)
		path := filepath.Join(dir, strings.ToLower(arm)+".json")
		if err := os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
		paths[arm] = path
	}
	return paths
}

func step0020ExperimentArgs(paths map[string]string, out string, extra ...string) []string {
	args := []string{"experiment", "--nominal", paths["N"], "--functional-v1", paths["F1"], "--functional-v2", paths["F2"], "--output-dir", out, "--no-token"}
	return append(args, extra...)
}

type step0020BudgetProvider struct{ requests []domain.ObservableWebRequest }

func (p *step0020BudgetProvider) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "github-repository-search", PreservesRank: true, MaxResultsLimit: 10}
}
func (p *step0020BudgetProvider) SearchOne(_ context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	p.requests = append(p.requests, req)
	r := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "github-repository-search", RetrievedAt: "2026-08-31T12:00:00Z", Status: "EMPTY"}
	return usecase.SealWebResponse(req, r), nil
}

func TestStep0020ContractFrozenBudgetAndInterleaving(t *testing.T) {
	dir := t.TempDir()
	paths := step0020QueryFiles(t, dir)
	p := &step0020BudgetProvider{}
	err := runIdeaWeb(step0020ExperimentArgs(paths, filepath.Join(dir, "out"), "--max-per-query", "10"), &bytes.Buffer{}, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.requests) != 36 {
		t.Fatalf("requests=%d, esperado orçamento congelado de 36", len(p.requests))
	}
	for i, req := range p.requests {
		arm := []string{"N", "F1", "F2"}[i%3]
		position := i/3 + 1
		if req.QueryID != fmt.Sprintf("%s%d", arm, position) {
			t.Fatalf("ordem[%d]=%s", i, req.QueryID)
		}
	}
}

func TestStep0020ContractRejectsMaxResultsAboveTen(t *testing.T) {
	dir := t.TempDir()
	paths := step0020QueryFiles(t, dir)
	p := &step0020BudgetProvider{}
	err := runIdeaWeb(step0020ExperimentArgs(paths, filepath.Join(dir, "out"), "--max-per-query", "99"), &bytes.Buffer{}, p)
	if err != nil {
		t.Fatal(err)
	}
	for _, req := range p.requests {
		if req.MaxResults > 10 {
			t.Fatalf("request %s excedeu 10 resultados: %d", req.RequestID, req.MaxResults)
		}
	}
}

func TestStep0020ContractRejectsExtraArmRequestBeforeNetwork(t *testing.T) {
	dir := t.TempDir()
	paths := step0020QueryFiles(t, dir)
	var nominal domain.FunctionalGrammarExpansion
	b, err := os.ReadFile(paths["N"])
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &nominal); err != nil {
		t.Fatal(err)
	}
	nominal.Queries = append(nominal.Queries, domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: "N13", Text: "unauthorized extra query"}})
	b, _ = json.Marshal(nominal)
	if err := os.WriteFile(paths["N"], b, 0600); err != nil {
		t.Fatal(err)
	}

	p := &step0020BudgetProvider{}
	err = runIdeaWeb(step0020ExperimentArgs(paths, filepath.Join(dir, "out")), &bytes.Buffer{}, p)
	if err == nil {
		t.Fatal("braço com 13ª query foi aceito")
	}
	if len(p.requests) != 0 {
		t.Fatalf("ataque extra-arm alcançou a rede: requests=%d", len(p.requests))
	}
}

type step0020CanaryProvider struct{ canary string }

func (p step0020CanaryProvider) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "github-repository-search", MaxResultsLimit: 10}
}
func (p step0020CanaryProvider) SearchOne(_ context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	r := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "github-repository-search", RetrievedAt: "2026-08-31T12:00:00Z", Status: "SUCCESS", Receipt: domain.WebReceipt{SanitizedURL: "https://example.invalid/?Authorization=Bearer+" + p.canary}, Results: []domain.RawWebResult{{Rank: 1, URL: "https://github.com/public/repo", Title: p.canary, Snippet: "public"}}}
	return usecase.SealWebResponse(req, r), nil
}

func TestStep0020ContractCanaryAbortsPersistence(t *testing.T) {
	dir := t.TempDir()
	paths := step0020QueryFiles(t, dir)
	out := filepath.Join(dir, "out")
	canary := "TEKT_STEP0020_CANARY_DO_NOT_PERSIST"
	err := runIdeaWeb(step0020ExperimentArgs(paths, out), &bytes.Buffer{}, step0020CanaryProvider{canary: canary})
	if err == nil {
		t.Fatal("execução persistiu artefato após detectar token-canário/Authorization")
	}
	entries, readErr := os.ReadDir(out)
	if readErr == nil && len(entries) != 0 {
		t.Fatalf("%d snapshots persistidos após falha de segurança", len(entries))
	}
}

func TestStep0020ContractInvalidSnapshotParserAndOfflineReplay(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	out := filepath.Join(dir, "replay.json")
	if err := os.WriteFile(bad, []byte(`{"schema_version":`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := runIdeaWeb([]string{"replay", "--snapshot", bad, "--output", out}, &bytes.Buffer{}, nil); err == nil {
		t.Fatal("parser inválido aceito")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatal("replay inválido produziu artefato")
	}
}

func TestStep0020ContractSnapshotsAreDeterministic(t *testing.T) {
	dir := t.TempDir()
	paths := step0020QueryFiles(t, dir)
	for _, name := range []string{"a", "b"} {
		if err := runIdeaWeb(step0020ExperimentArgs(paths, filepath.Join(dir, name)), &bytes.Buffer{}, &step0020BudgetProvider{}); err != nil {
			t.Fatal(err)
		}
	}
	for _, arm := range []string{"n", "f1", "f2"} {
		a, _ := os.ReadFile(filepath.Join(dir, "a", "snapshot-"+arm+".json"))
		b, _ := os.ReadFile(filepath.Join(dir, "b", "snapshot-"+arm+".json"))
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("snapshot %s não determinístico", arm)
		}
	}
}
