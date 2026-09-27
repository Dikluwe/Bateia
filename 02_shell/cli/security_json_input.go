/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_json_input.md
 * @layer L2
 * @updated 2026-09-26
 */

package cli

import "fmt"

// JSON requests have both byte and structural budgets, before typed decoding.
// Text Git metadata and byte hashes deliberately use readIdeaInput instead.
func readIdeaJSONInput(path string) ([]byte, error) {
	data, err := readIdeaInput(path)
	if err != nil {
		return nil, err
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return nil, fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	return data, nil
}
