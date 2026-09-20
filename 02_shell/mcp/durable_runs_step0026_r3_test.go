/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/durable_runs_step0026_r3_test.md
 * @layer L2
 * @updated 2026-09-03
 */
package mcp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"bateia/01_core/domain"
)

const step0026R3LegacyBenchmarkGolden = `{"jsonrpc":"2.0","id":41,"error":{"code":-32602,"message":"Invalid params: depth deve estar entre 1 e 5"}}`
const step0026R3LegacyRecommendationGolden = `{"schema_version":"1","status":"COMPLETE","summary":{"expansion_count":0,"candidate_count":0,"assessed_count":0},"decision":{"schema_version":"1","policy":{"schema_version":"","allowed_licenses":null,"require_local":false,"allow_archived":false,"max_components":0,"extend_min_coverage":0,"integration_min_value":0},"policy_hash":"","decision":"BUILD","matrix":[],"gaps":[],"reasons":[],"alternatives":[],"effort":{"level":"","factors":null}},"alternatives":{"schema_version":"1","target":{"name":""},"filters":{"allowed_licenses":null,"require_local":false,"allow_archived":false,"recommended_min_coverage":0,"partial_min_coverage":0},"alternatives":[],"excluded":[]},"failures":[],"warnings":[]}`

func TestStep0026R3LegacyBenchmarkAndRecommendationUseIndependentByteGoldens(t *testing.T) {
	benchmark := step0025ServeLegacy(t, `{"jsonrpc":"2.0","id":41,"method":"tools/call","params":{"name":"bateia_benchmark","arguments":{"repos":"acme/example"}}}`)

	fixture := domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunComplete,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionBuild, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{}, Warnings: []string{},
	}
	spy := &step0025RecommenderSpy{result: fixture}
	responses := step0025Serve(t, spy, `{"jsonrpc":"2.0","id":42,"method":"tools/call","params":{"name":"bateia_recommend","arguments":`+step0025MCPRequest()+`}}`)
	if len(responses) != 1 || responses[0].Error != nil || spy.calls != 1 {
		t.Fatalf("legacy recommendation did not execute once: %#v calls=%d", responses, spy.calls)
	}
	result := decodeStep0025Result(t, responses[0].Result)
	content := result["content"].([]interface{})
	text := content[0].(map[string]interface{})["text"].(string)
	if text != step0026R3LegacyRecommendationGolden {
		t.Fatalf("freeze recommendation golden:\n%s", text)
	}
	if len(benchmark) != 1 || benchmark[0].Error == nil {
		t.Fatalf("legacy benchmark route did not execute: %#v", benchmark)
	}
	benchmarkBytes := mustStep0025JSON(t, benchmark[0])
	if string(benchmarkBytes) != step0026R3LegacyBenchmarkGolden {
		t.Fatalf("freeze benchmark golden:\n%s", benchmarkBytes)
	}
	var closed map[string]interface{}
	if err := json.Unmarshal([]byte(text), &closed); err != nil {
		t.Fatal(err)
	}
	if closed["schema_version"] != "1" || closed["status"] != "COMPLETE" {
		t.Fatalf("unexpected legacy recommendation projection: %#v", closed)
	}
}

func TestStep0026R3PredecessorServerHashesRemainCertified(t *testing.T) {
	for path, want := range map[string]string{
		"server.go": "e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42",
		"../../01_core/usecase/recommendation.go":                                   "ad6badcdbf33019ffb0c81274961e68524da933a832c7cbbc6415b5c35491e1b",
		"../../01_core/domain/recommendation.go":                                    "2fa9bf8e8beba55f22279b33cff20084834e510dc6175ad2ee70b6791f149af6",
		"../../lab/bancada/history/benchmarks/0029/historical-hash-invalidation-o2-to-o3.json": "98c36700e72c434a2443ce4aa7e64c5b5c9b599f0eb0b5451cf07c9af46639ca",
		"../../lab/bancada/history/benchmarks/0029/historical-hash-successor-o3.json":          "cf5da34d390b710425f71a5545d46978245db3581d5f2caacf762980238c4156",
		"../../lab/bancada/history/benchmarks/0031/historical-contract-supersession.json":      "65598792dd329561e7d8a21bd1cf0bbd22c19d00876470361ac2ce6e4d6aae67",
	} {
		if got := step0026R3SHA256File(t, path); got != want {
			t.Fatalf("certified predecessor/successor lineage changed: %s got=%s want=%s", path, got, want)
		}
	}
}

func step0026R3SHA256File(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
