/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/host.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"hash"
	"io/fs"
	"time"
)

var configuredCommandHost domain.CommandHost
var errCommandHost = errors.New("command host is not configured; compose it before running CLI commands")

// ConfigureCommandHost is a startup-only composition seam.
func ConfigureCommandHost(host domain.CommandHost) { configuredCommandHost = host }
func commandHost() domain.CommandHost {
	if configuredCommandHost != nil {
		return configuredCommandHost
	}
	return unavailableCommandHost{}
}
func commandContext(ctx context.Context) context.Context {
	return domain.WithExecutionRuntime(ctx, configuredCommandHost)
}
func isCommandNotExist(err error) bool { return errors.Is(err, fs.ErrNotExist) }

type unavailableCommandHost struct{}

func (unavailableCommandHost) LocalGitRepositoryIdentities(string, int64) []domain.RepositoryIdentity {
	return nil
}
func (unavailableCommandHost) ReadLocalGitConfig(string, int64) ([]byte, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) ReadRootGitMetadata(domain.CommandRoot, string, int) ([]byte, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) CaptureLocal(context.Context, []string, string, bool) (domain.CapturedWorkspace, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) NormalizedLocalRoots(context.Context, []string, string) ([]string, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) SnapshotProvider(domain.ContentProvider, []domain.SnapshotBinding) (domain.ContentProvider, error) {
	return nil, errCommandHost
}

func (unavailableCommandHost) TreeHash(context.Context, string, string, *domain.TreeHashBudget) (string, error) {
	return "", errCommandHost
}
func (unavailableCommandHost) ValidateHashLayout(string, string) error { return errCommandHost }
func (unavailableCommandHost) HashRootedTree(context.Context, domain.CommandRoot, string, *domain.TreeHashBudget, domain.CommandRoot) (string, error) {
	return "", errCommandHost
}
func (unavailableCommandHost) HashRootedFile(context.Context, domain.CommandRoot, string, hash.Hash, []byte, *domain.TreeHashBudget, domain.CommandRoot) error {
	return errCommandHost
}

func (unavailableCommandHost) SameOutputPath(string, string) bool { return true }

func (unavailableCommandHost) ValidateTraceOutput(string, ...string) error    { return errCommandHost }
func (unavailableCommandHost) ResolveProtectedPath(string) (string, error)    { return "", errCommandHost }
func (unavailableCommandHost) WriteTraceFile([]byte, string, ...string) error { return errCommandHost }

func (unavailableCommandHost) OpenAlgorithmHistoryInput(string, string) (domain.CommandFile, domain.FileInfo, error) {
	return nil, nil, errCommandHost
}
func (unavailableCommandHost) WriteAlgorithmHistoryAtomic(string, []byte, domain.FileInfo) error {
	return errCommandHost
}
func (unavailableCommandHost) OpenIdeaRadarExportInput(string, string) (domain.CommandFile, domain.FileInfo, error) {
	return nil, nil, errCommandHost
}
func (unavailableCommandHost) RejectIdeaRadarExportAlias(string, string, domain.FileInfo) error {
	return errCommandHost
}
func (unavailableCommandHost) WriteIdeaRadarExportAtomic(string, []byte, domain.FileInfo) error {
	return errCommandHost
}

func (unavailableCommandHost) ReadDescriptor(domain.CommandFile, int64) ([]byte, error) {
	return nil, errCommandHost
}

func (unavailableCommandHost) ReadBounded(string, int64) ([]byte, error)     { return nil, errCommandHost }
func (unavailableCommandHost) PrepareDossier(string, bool) error             { return errCommandHost }
func (unavailableCommandHost) WriteAtomic(string, []byte, fs.FileMode) error { return errCommandHost }

func (unavailableCommandHost) Now() time.Time { return time.Time{} }
func (unavailableCommandHost) WithTimeout(ctx context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	next, cancel := context.WithCancel(ctx)
	cancel()
	return next, cancel
}
func (unavailableCommandHost) NotifyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	next, cancel := context.WithCancel(ctx)
	cancel()
	return next, cancel
}
func (unavailableCommandHost) Open(string) (domain.CommandFile, error) { return nil, errCommandHost }
func (unavailableCommandHost) OpenFile(string, domain.FileOpenFlags, fs.FileMode) (domain.CommandFile, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) CreateTemp(string, string) (domain.CommandFile, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) OpenRoot(string) (domain.CommandRoot, error) {
	return nil, errCommandHost
}
func (unavailableCommandHost) ReadFile(string) ([]byte, error)             { return nil, errCommandHost }
func (unavailableCommandHost) WriteFile(string, []byte, fs.FileMode) error { return errCommandHost }
func (unavailableCommandHost) Mkdir(string, fs.FileMode) error             { return errCommandHost }
func (unavailableCommandHost) MkdirAll(string, fs.FileMode) error          { return errCommandHost }
func (unavailableCommandHost) MkdirTemp(string, string) (string, error)    { return "", errCommandHost }
func (unavailableCommandHost) Remove(string) error                         { return errCommandHost }
func (unavailableCommandHost) RemoveAll(string) error                      { return errCommandHost }
func (unavailableCommandHost) Rename(string, string) error                 { return errCommandHost }
func (unavailableCommandHost) Stat(string) (fs.FileInfo, error)            { return nil, errCommandHost }
func (unavailableCommandHost) Lstat(string) (fs.FileInfo, error)           { return nil, errCommandHost }
func (unavailableCommandHost) SameFile(fs.FileInfo, fs.FileInfo) bool      { return false }
func (unavailableCommandHost) LinkCount(fs.FileInfo) (uint64, bool)        { return 0, false }
func (unavailableCommandHost) Abs(string) (string, error)                  { return "", errCommandHost }
func (unavailableCommandHost) Glob(string) ([]string, error)               { return nil, errCommandHost }
func (unavailableCommandHost) EvalSymlinks(string) (string, error)         { return "", errCommandHost }
func (unavailableCommandHost) Walk(string, func(string, fs.FileInfo, error) error) error {
	return errCommandHost
}
func (unavailableCommandHost) Getenv(string) string           { return "" }
func (unavailableCommandHost) ReadRandom([]byte) (int, error) { return 0, errCommandHost }
