/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/end_to_end_benchmark_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBenchmarkContractsPreserveUnknownAndNullableMetrics(t *testing.T) {
	r := EndToEndBenchmarkReport{SchemaVersion: "1", Cases: []BenchmarkCaseResult{{Correspondence: CorrespondenceUnknown}}, Metrics: EndToEndBenchmarkMetrics{EvidencePrecision: nil}}
	b, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	for _, x := range []string{"UNKNOWN", "evidence_precision\":null"} {
		if !strings.Contains(string(b), x) {
			t.Fatalf("missing %s", x)
		}
	}
}
