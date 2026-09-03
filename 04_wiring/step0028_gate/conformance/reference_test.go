/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/reference_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStep0028O2ReferenceDemonstratesEveryFrozenConformanceCase(t *testing.T) {
	fixture := filepath.Join("..", "..", "..", "00_nucleo", "benchmark", "0028", "oracle", "protected", "fixtures", "conformance-cases.json")
	document, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	cases, err := DecodeCases(document)
	if err != nil {
		t.Fatal(err)
	}
	findings := Evaluate(context.Background(), ReferenceSubject{}, cases)
	seen := map[Obligation]bool{}
	for _, finding := range findings {
		if finding.Outcome != "Preserved" {
			t.Fatalf("reference failed %s/%s: %s", finding.CaseID, finding.Obligation, finding.Detail)
		}
		seen[finding.Obligation] = true
	}
	if len(seen) != 12 {
		t.Fatalf("reference demonstrated %d obligations, want 12", len(seen))
	}
}
