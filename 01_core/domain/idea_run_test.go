/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_run_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIdeaRunManifestKeepsLineageWithoutCredentialFields(t *testing.T) {
	m := IdeaRunManifest{SchemaVersion: "1", Mode: IdeaRunReplay, Status: IdeaRunInconclusive, InputHashes: map[string]string{"idea": "abc"}, RuleVersions: map[string]string{"assessment": "1"}, Stages: []IdeaRunStage{}, Failures: []ExecutionFailure{}}
	b, e := json.Marshal(m)
	if e != nil {
		t.Fatal(e)
	}
	s := string(b)
	for _, x := range []string{"input_hashes", "rule_versions", "INCONCLUSIVE"} {
		if !strings.Contains(s, x) {
			t.Fatalf("missing %s", x)
		}
	}
	for _, x := range []string{"token", "authorization", "credential"} {
		if strings.Contains(strings.ToLower(s), x) {
			t.Fatalf("credential field %s", x)
		}
	}
}
