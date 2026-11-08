/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/local_snapshot_capture.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"context"
)

const localCollectionVersion = "frozen-local-v1"
const maxCapturedGitConfigBytes = 64 << 10

type capturedLocalRun struct {
	Candidates []domain.RepositoryCandidate
	Bindings   []domain.SnapshotBinding
	Identities map[string][]domain.RepositoryIdentity
	directory  string
	resource   domain.CapturedWorkspace
}

func (s *capturedLocalRun) Close() error {
	if s == nil || s.resource == nil {
		return nil
	}
	return s.resource.Close()
}
func normalizedLocalHashRoots(ctx context.Context, paths []string, excluded string) ([]string, error) {
	return commandHost().NormalizedLocalRoots(ctx, paths, excluded)
}
func captureLocalCandidates(ctx context.Context, paths []string, excluded string, metadata bool) (*capturedLocalRun, error) {
	s, err := commandHost().CaptureLocal(ctx, paths, excluded, metadata)
	if err != nil {
		return nil, err
	}
	view := s.View()
	return &capturedLocalRun{Candidates: view.Candidates, Bindings: view.Bindings, Identities: view.Identities, directory: view.ContainerPath, resource: s}, nil
}

type SnapshotProviderFactory func(domain.ContentProvider, []domain.SnapshotBinding) (domain.ContentProvider, error)

var snapshotProviderFactory SnapshotProviderFactory

func SetSnapshotProviderFactory(factory SnapshotProviderFactory) { snapshotProviderFactory = factory }
func buildSnapshotProvider(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
	if snapshotProviderFactory != nil {
		return snapshotProviderFactory(fallback, bindings)
	}
	return commandHost().SnapshotProvider(fallback, bindings)
}
