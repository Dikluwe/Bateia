/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/web_input.md
 * @layer L2
 * @updated 2026-10-29
 */
package cli

// WebTokenRequested preserves the early literal probe before full command parsing.
func WebTokenRequested(args []string) bool {
	for _, arg := range args {
		if arg == "--no-token" {
			return false
		}
	}
	return true
}
