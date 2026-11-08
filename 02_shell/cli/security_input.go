/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_input.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

const maxIdeaInputBytes = 16 << 20

func readIdeaInput(path string) ([]byte, error) {
	return commandHost().ReadBounded(path, maxIdeaInputBytes)
}

func readIdeaDescriptor(file domain.CommandFile) ([]byte, error) {
	return commandHost().ReadDescriptor(file, maxIdeaInputBytes)
}
func validateIdeaJSONStructure(data []byte) error { return parser.ValidateBoundedJSON(data) }
