/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/discovery_coverage_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package report_test

import (
	"bateia/01_core/domain"
	"bateia/02_shell/report"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestDiscovery0056MarkdownCoverageBoundsAndJSONCompleteness(t *testing.T) {
	matches := 307
	coverage := &domain.DiscoveryCoverage{RequestedMatches: 200, ReceivedItems: 173, ObservedRepositories: 149, ExaminedRepositories: 137, DeferredRepositories: 12, EligibleCandidates: 29, SelectedCandidates: 7, ProviderMatches: &matches, Completeness: "unknown"}
	a := domain.SearchAttempt{Query: "synthetic derived", OriginalQuery: "serializar dados", QueryPlan: &domain.DiscoveryQueryPlan{Version: "0055-v2", Status: "hypothesis"}, Coverage: coverage, RawCandidates: 137, RejectedCandidates: 108, Candidates: 7}
	for i := 0; i < 12; i++ {
		d := domain.DiscoveryCandidateDiagnostic{Repository: fmt.Sprintf("control/diagnostic_%02d_END", i), Status: "insufficient_evidence", Reason: "split_across_files", ActionFound: true, ObjectFound: true, TruncatedSnippets: 41, OmittedSnippets: 43}
		for j := 0; j < 5; j++ {
			d.Evidence = append(d.Evidence, domain.MatchEvidence{MatchType: "FileMatch", Path: fmt.Sprintf("path_%02d_%02d_END.go", i, j), Snippet: fmt.Sprintf("preview_%02d_%02d_END", i, j)})
		}
		a.CandidateDiagnostics = append(a.CandidateDiagnostics, d)
	}
	r := &domain.SearchReport{SearchAttempts: []domain.SearchAttempt{a}}
	before, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	md := report.GenerateMarkdownReport(r)
	for _, want := range []string{"200", "173", "149", "137", "12", "29", "7", "307", "41", "43", "split_across_files"} {
		if !strings.Contains(md, want) {
			t.Errorf("missing coverage/reason/retention %q in %s", want, md)
		}
	}
	lower := strings.ToLower(md)
	if !strings.Contains(lower, "unknown") && !strings.Contains(lower, "desconhecid") {
		t.Errorf("unknown completeness absent: %s", md)
	}
	if !strings.Contains(lower, "count") || !strings.Contains(lower, "resultad") || !strings.Contains(lower, "provedor") {
		t.Errorf("count provider-results unit absent: %s", md)
	}
	if !strings.Contains(lower, "insufici") || !strings.Contains(lower, "trech") || !strings.Contains(lower, "capacidade") {
		t.Errorf("snippet-insufficiency epistemic statement absent: %s", md)
	}
	if !strings.Contains(lower, "omit") || !strings.Contains(lower, "2") {
		t.Errorf("missing omitted display counts: %s", md)
	}
	for i := 0; i < 12; i++ {
		name := fmt.Sprintf("control/diagnostic_%02d_END", i)
		if strings.Contains(md, name) != (i < 10) {
			t.Errorf("diagnostic display bound %s", name)
		}
		for j := 0; j < 5; j++ {
			for _, value := range []string{fmt.Sprintf("path_%02d_%02d_END.go", i, j), fmt.Sprintf("preview_%02d_%02d_END", i, j)} {
				if strings.Contains(md, value) != (i < 10 && j < 3) {
					t.Errorf("evidence display bound %s", value)
				}
			}
		}
	}
	after, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("Markdown truncated underlying JSON")
	}
	var decoded domain.SearchReport
	if err = json.Unmarshal(after, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.SearchAttempts) != 1 || !reflect.DeepEqual(decoded.SearchAttempts[0].CandidateDiagnostics, a.CandidateDiagnostics) || !reflect.DeepEqual(decoded.SearchAttempts[0].Coverage, coverage) {
		t.Fatalf("JSON lost complete diagnostics %s", after)
	}
}
