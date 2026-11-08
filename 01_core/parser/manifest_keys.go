/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/manifest_keys.md
 * @layer L1
 * @updated 2026-10-28
 */

package parser

import "sort"

// sortedManifestKeys keeps object iteration independent of Go map ordering.
func sortedManifestKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
