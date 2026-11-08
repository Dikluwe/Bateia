/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/command_files.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"context"
	"hash"
	"io"
	"io/fs"
)

type FileInfo = fs.FileInfo
type FileMode = fs.FileMode
type DirEntry = fs.DirEntry
type FileOpenFlags uint32

const (
	OpenReadOnly  FileOpenFlags = 0
	OpenWriteOnly FileOpenFlags = 1 << iota
	OpenReadWrite
	OpenAppend
	OpenCreate
	OpenExclusive
	OpenTruncate
	OpenNonblock
	OpenNoFollow
)

type CommandFile interface {
	io.Reader
	io.Writer
	io.Closer
	Stat() (fs.FileInfo, error)
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	ReadDir(int) ([]fs.DirEntry, error)
}
type CommandRoot interface {
	Close() error
	OpenFile(string, FileOpenFlags, fs.FileMode) (CommandFile, error)
	Stat(string) (fs.FileInfo, error)
	Lstat(string) (fs.FileInfo, error)
	MkdirAll(string, fs.FileMode) error
	Remove(string) error
	Rename(string, string) error
}

type CommandStorage interface {
	ReadBounded(string, int64) ([]byte, error)
	ReadDescriptor(CommandFile, int64) ([]byte, error)
	PrepareDossier(string, bool) error
	WriteAtomic(string, []byte, fs.FileMode) error
	WriteFile(string, []byte, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
	Stat(string) (fs.FileInfo, error)
	Lstat(string) (fs.FileInfo, error)
	Abs(string) (string, error)
	Glob(string) ([]string, error)
}
type CommandExports interface {
	OpenAlgorithmHistoryInput(string, string) (CommandFile, FileInfo, error)
	WriteAlgorithmHistoryAtomic(string, []byte, FileInfo) error
	OpenIdeaRadarExportInput(string, string) (CommandFile, FileInfo, error)
	RejectIdeaRadarExportAlias(string, string, FileInfo) error
	WriteIdeaRadarExportAtomic(string, []byte, FileInfo) error
}
type CommandTraceFiles interface {
	ValidateTraceOutput(string, ...string) error
	ResolveProtectedPath(string) (string, error)
	WriteTraceFile([]byte, string, ...string) error
	SameOutputPath(string, string) bool
}
type CommandWorkspace interface {
	OpenRoot(string) (CommandRoot, error)
	TreeHash(context.Context, string, string, *TreeHashBudget) (string, error)
	ValidateHashLayout(string, string) error
	HashRootedTree(context.Context, CommandRoot, string, *TreeHashBudget, CommandRoot) (string, error)
	HashRootedFile(context.Context, CommandRoot, string, hash.Hash, []byte, *TreeHashBudget, CommandRoot) error
	LocalGitRepositoryIdentities(string, int64) []RepositoryIdentity
	ReadLocalGitConfig(string, int64) ([]byte, error)
	ReadRootGitMetadata(CommandRoot, string, int) ([]byte, error)
	CaptureLocal(context.Context, []string, string, bool) (CapturedWorkspace, error)
	NormalizedLocalRoots(context.Context, []string, string) ([]string, error)
	SnapshotProvider(ContentProvider, []SnapshotBinding) (ContentProvider, error)
}
type CommandEnvironment interface {
	Getenv(string) string
	NotifyContext(context.Context) (context.Context, context.CancelFunc)
}
type CommandHost interface {
	ExecutionRuntime
	CommandStorage
	CommandExports
	CommandTraceFiles
	CommandWorkspace
	CommandEnvironment
}
