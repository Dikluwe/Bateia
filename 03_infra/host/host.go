/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/host.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"context"
	"crypto/rand"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"bateia/01_core/domain"
)

type Host struct{}

func New() *Host { return &Host{} }

var _ domain.CommandHost = (*Host)(nil)

func (*Host) Now() time.Time { return time.Now().UTC() }
func (*Host) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}
func (*Host) Open(p string) (domain.CommandFile, error) { return os.Open(p) }
func (*Host) OpenFile(p string, f domain.FileOpenFlags, m fs.FileMode) (domain.CommandFile, error) {
	return os.OpenFile(p, nativeFlags(f), m)
}
func (*Host) CreateTemp(d, p string) (domain.CommandFile, error) { return os.CreateTemp(d, p) }
func (*Host) ReadFile(p string) ([]byte, error)                  { return os.ReadFile(p) }
func (*Host) WriteFile(p string, b []byte, m fs.FileMode) error  { return os.WriteFile(p, b, m) }
func (*Host) Mkdir(p string, m fs.FileMode) error                { return os.Mkdir(p, m) }
func (*Host) MkdirAll(p string, m fs.FileMode) error             { return os.MkdirAll(p, m) }
func (*Host) MkdirTemp(d, p string) (string, error)              { return os.MkdirTemp(d, p) }
func (*Host) Remove(p string) error                              { return os.Remove(p) }
func (*Host) RemoveAll(p string) error                           { return os.RemoveAll(p) }
func (*Host) Rename(a, b string) error                           { return os.Rename(a, b) }
func (*Host) Stat(p string) (fs.FileInfo, error)                 { return os.Stat(p) }
func (*Host) Lstat(p string) (fs.FileInfo, error)                { return os.Lstat(p) }
func (*Host) SameFile(a, b fs.FileInfo) bool                     { return os.SameFile(a, b) }
func (*Host) LinkCount(i fs.FileInfo) (uint64, bool) {
	s, ok := i.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, false
	}
	return uint64(s.Nlink), true
}
func (*Host) Abs(p string) (string, error)          { return filepath.Abs(p) }
func (*Host) Glob(p string) ([]string, error)       { return filepath.Glob(p) }
func (*Host) EvalSymlinks(p string) (string, error) { return filepath.EvalSymlinks(p) }
func (*Host) Walk(p string, f func(string, fs.FileInfo, error) error) error {
	return filepath.Walk(p, f)
}
func (*Host) Getenv(k string) string           { return os.Getenv(k) }
func (*Host) ReadRandom(b []byte) (int, error) { return rand.Read(b) }
func (*Host) NotifyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
}
func (*Host) OpenRoot(p string) (domain.CommandRoot, error) {
	r, e := os.OpenRoot(p)
	if e != nil {
		return nil, e
	}
	return &root{Root: r}, nil
}

type root struct{ *os.Root }

func (r *root) OpenFile(p string, f domain.FileOpenFlags, m fs.FileMode) (domain.CommandFile, error) {
	return r.Root.OpenFile(p, nativeFlags(f), m)
}
func nativeFlags(f domain.FileOpenFlags) int {
	n := os.O_RDONLY
	for _, flag := range []struct {
		portable domain.FileOpenFlags
		native   int
	}{
		{domain.OpenWriteOnly, os.O_WRONLY}, {domain.OpenReadWrite, os.O_RDWR}, {domain.OpenAppend, os.O_APPEND}, {domain.OpenCreate, os.O_CREATE}, {domain.OpenExclusive, os.O_EXCL}, {domain.OpenTruncate, os.O_TRUNC}, {domain.OpenNonblock, syscall.O_NONBLOCK}, {domain.OpenNoFollow, syscall.O_NOFOLLOW},
	} {
		if f&flag.portable != 0 {
			n |= flag.native
		}
	}
	return n
}
