/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/software_alternative_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSoftwareAlternativeFitSerializesAllSealedStates(t *testing.T) {
	states := []SoftwareAlternativeFit{
		AlternativeRecommended,
		AlternativePartial,
		AlternativeReference,
		AlternativeBlocked,
		AlternativeInconclusive,
		AlternativeNotAlternative,
	}
	want := []string{"RECOMMENDED", "PARTIAL", "REFERENCE", "BLOCKED", "INCONCLUSIVE", "NOT_ALTERNATIVE"}
	encoded, err := json.Marshal(states)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fit states = %v, want %v", got, want)
	}
}

func TestSoftwareAlternativeReportJSONPreservesSealedSurface(t *testing.T) {
	payload := []byte(`{
		"schema_version":"1",
		"target":{"name":"Known Tool","repository":{"owner":"acme","repo":"known"}},
		"filters":{"allowed_licenses":["Apache-2.0","MIT"],"require_local":true,"allow_archived":false,"recommended_min_coverage":0.8,"partial_min_coverage":0.5},
		"top_choice":{"owner":"org","repo":"best","commit":"abc123"},
		"alternatives":[{
			"repository":{"owner":"org","repo":"best","commit":"abc123"},
			"url":"https://github.com/org/best","description":"proven replacement","stars":42,
			"fit":"RECOMMENDED","coverage":1,"problem_alignment":"PROVEN",
			"capabilities":[{"capability":"dependency graph","status":"PROVEN"}],"gaps":[],
			"license":{"status":"PROVEN","identifier":"MIT"},
			"maintenance":{"status":"PROVEN","archived":false},
			"operation":{"local_status":"PROVEN","requirements":["self-hosted"]},
			"reasons":["all required capabilities proven"],"blockers":[],
			"evidence":[{"kind":"file","ref":{"owner":"org","repo":"best","commit":"abc123"},"path":"graph.go","rule_id":"CAP-001","strength":1}]
		}],
		"excluded":[{
			"repository":{"owner":"org","repo":"blocked","commit":"def456"},
			"url":"https://github.com/org/blocked","description":"incompatible license","stars":100,
			"fit":"BLOCKED","coverage":1,"problem_alignment":"PROVEN",
			"capabilities":[{"capability":"dependency graph","status":"PROVEN"}],"gaps":[],
			"license":{"status":"PROVEN","identifier":"GPL-3.0"},
			"maintenance":{"status":"PROVEN","archived":false},
			"operation":{"local_status":"PROVEN"},
			"reasons":[],"blockers":["license GPL-3.0 is not allowed"],"evidence":[]
		}]
	}`)
	var report SoftwareAlternativeReport
	if err := json.Unmarshal(payload, &report); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	assertJSONPathEquals(t, got, "schema_version", "1")
	assertJSONPathEquals(t, got, "target.name", "Known Tool")
	assertJSONPathEquals(t, got, "target.repository.owner", "acme")
	assertJSONPathEquals(t, got, "filters.require_local", true)
	assertJSONPathEquals(t, got, "filters.recommended_min_coverage", float64(0.8))
	assertJSONPathEquals(t, got, "top_choice.commit", "abc123")
	assertJSONPathEquals(t, got, "alternatives.0.fit", "RECOMMENDED")
	assertJSONPathEquals(t, got, "alternatives.0.capabilities.0.capability", "dependency graph")
	assertJSONPathEquals(t, got, "alternatives.0.capabilities.0.status", "PROVEN")
	assertJSONPathEquals(t, got, "alternatives.0.evidence.0.ref.commit", "abc123")
	assertJSONPathEquals(t, got, "excluded.0.fit", "BLOCKED")
	assertJSONPathEquals(t, got, "excluded.0.blockers.0", "license GPL-3.0 is not allowed")
}

func assertJSONPathEquals(t *testing.T, root map[string]any, path string, want any) {
	t.Helper()
	var current any = root
	parts := splitJSONTestPath(path)
	for _, part := range parts {
		switch node := current.(type) {
		case map[string]any:
			current = node[part]
		case []any:
			index := int(part[0] - '0')
			if len(part) != 1 || index < 0 || index >= len(node) {
				t.Fatalf("invalid test path %q", path)
			}
			current = node[index]
		default:
			t.Fatalf("path %q stopped at %q (%T)", path, part, current)
		}
	}
	if !reflect.DeepEqual(current, want) {
		t.Fatalf("%s = %#v, want %#v", path, current, want)
	}
}

func splitJSONTestPath(path string) []string {
	parts := []string{}
	start := 0
	for index := range path {
		if path[index] == '.' {
			parts = append(parts, path[start:index])
			start = index + 1
		}
	}
	return append(parts, path[start:])
}
