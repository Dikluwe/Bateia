/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_git_metadata.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import "bateia/01_core/domain"

func readLocalGitConfig(root string) ([]byte, error) {
	return commandHost().ReadLocalGitConfig(root, maxIdeaInputBytes)
}
func readLocalGitConfigLimit(root string, limit int64) ([]byte, error) {
	return commandHost().ReadLocalGitConfig(root, limit)
}
func readRootGitMetadata(root domain.CommandRoot, name string, limit int) ([]byte, error) {
	return commandHost().ReadRootGitMetadata(root, name, limit)
}
