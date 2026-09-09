/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/query_expansion_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package report_test

import (
	"bateia/01_core/domain"
	"bateia/02_shell/report"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestExpansion0052MarkdownBoundsAndJSONCompleteness(t *testing.T) {
	q := domain.ImplementationQuery{Original: "consulta original", Mode: "natural"}
	for i := 0; i < 7; i++ {
		e := domain.ImplementationQueryExpansion{Source: fmt.Sprintf("origem%d", i), RuleID: fmt.Sprintf("0052-v1-family%d", i), Kind: "translation", Status: "hypothesis"}
		for j := 0; j < 10; j++ {
			e.Terms = append(e.Terms, fmt.Sprintf("termo%d_%d", i, j))
		}
		q.Expansions = append(q.Expansions, e)
	}
	var b strings.Builder
	report.WriteQueryInterpretation(&b, []domain.ImplementationQuery{q})
	s := b.String()
	for _, part := range []string{"origem0", "0052-v1-family0", "translation", "hypothesis", "termo0_7"} {
		if !strings.Contains(s, part) {
			t.Errorf("missing visible provenance %q in %s", part, s)
		}
	}
	for _, part := range []string{"origem6", "termo0_8", "termo0_9"} {
		if strings.Contains(s, part) {
			t.Errorf("display limit leaked %q", part)
		}
	}
	lower := strings.ToLower(s)
	if !strings.Contains(lower, "prova") || (!strings.Contains(lower, "capacidade") && !strings.Contains(lower, "intenção")) {
		t.Errorf("missing hypothesis warning: %s", s)
	}
	if !strings.Contains(s, "2") || !strings.Contains(s, "1") {
		t.Errorf("missing remaining counts: %s", s)
	}
	raw, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "origem6") || !strings.Contains(string(raw), "termo6_9") {
		t.Fatal("JSON lost full receipts")
	}
}
