/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/manifest_keys_test.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import (
	"reflect"
	"testing"
)

func TestManifestKeysPreserveInputAndResultOwnership(t *testing.T) {
	source := map[string][]string{"z": {"2"}, "á": {"3"}, "a": {"1"}}
	original := map[string][]string{"z": {"2"}, "á": {"3"}, "a": {"1"}}
	keys := sortedManifestKeys(source)
	if !reflect.DeepEqual(keys, []string{"a", "z", "á"}) {
		t.Fatalf("noncanonical keys: %v", keys)
	}
	keys[0] = "changed"
	if !reflect.DeepEqual(source, original) {
		t.Fatalf("input mutated: %v", source)
	}
	if next := sortedManifestKeys(source); !reflect.DeepEqual(next, []string{"a", "z", "á"}) {
		t.Fatalf("returned lists share state: %v", next)
	}
	for _, empty := range []map[string]string{nil, {}} {
		if len(sortedManifestKeys(empty)) != 0 {
			t.Fatal("empty object invented dependencies")
		}
	}
}
