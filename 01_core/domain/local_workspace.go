/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/local_workspace.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

type TreeHashBudget struct {
	EntriesLeft int
	BytesLeft   int64
}
type CapturedWorkspaceView struct {
	Candidates    []RepositoryCandidate
	Bindings      []SnapshotBinding `json:"-"`
	Identities    map[string][]RepositoryIdentity
	ContainerPath string `json:"-"`
}
type CapturedWorkspace interface {
	View() CapturedWorkspaceView
	Close() error
}
