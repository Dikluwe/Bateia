/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/durableconfig/config_test.md
 * @layer L3
 * @updated 2026-09-18
 */

package durableconfig_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"bateia/03_infra/durableconfig"
)

const (
	runRootEnvironment = "BATEIA_RUN_ROOT"
	cursorKeyName      = "cursor.key"
	cursorKeyBytes     = 32
)

func TestLoadLocalCreatesAndReloadsExplicitRoot(t *testing.T) {
	workingDirectory := t.TempDir()
	t.Chdir(workingDirectory)
	t.Setenv(runRootEnvironment, filepath.Join("state", "runs"))

	expectedRoot, err := filepath.Abs(filepath.Join("state", "runs"))
	if err != nil {
		t.Fatalf("resolve expected root: %v", err)
	}

	first, err := durableconfig.LoadLocal()
	if err != nil {
		t.Fatalf("first LoadLocal: %v", err)
	}
	assertConfigDefaults(t, first)
	if first.Root != expectedRoot {
		t.Fatalf("root = %q, want %q", first.Root, expectedRoot)
	}
	if len(first.Secret) != cursorKeyBytes {
		t.Fatalf("secret length = %d, want %d", len(first.Secret), cursorKeyBytes)
	}

	for _, directory := range []string{
		filepath.Join(workingDirectory, "state"),
		expectedRoot,
	} {
		info, statErr := os.Stat(directory)
		if statErr != nil {
			t.Fatalf("stat created directory %q: %v", directory, statErr)
		}
		if !info.IsDir() {
			t.Fatalf("created component %q is not a directory", directory)
		}
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0o700 {
			t.Errorf("directory %q mode = %04o, want 0700", directory, info.Mode().Perm())
		}
	}

	keyPath := filepath.Join(expectedRoot, cursorKeyName)
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("read created cursor key: %v", err)
	}
	if !bytes.Equal(keyBytes, first.Secret) {
		t.Fatal("returned secret differs from the persisted cursor key")
	}
	keyInfo, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("stat created cursor key: %v", err)
	}
	if !keyInfo.Mode().IsRegular() {
		t.Fatalf("cursor key mode = %v, want a regular file", keyInfo.Mode())
	}
	if runtime.GOOS != "windows" && keyInfo.Mode().Perm() != 0o600 {
		t.Errorf("cursor key mode = %04o, want 0600", keyInfo.Mode().Perm())
	}

	second, err := durableconfig.LoadLocal()
	if err != nil {
		t.Fatalf("second LoadLocal: %v", err)
	}
	assertConfigDefaults(t, second)
	if second.Root != expectedRoot {
		t.Fatalf("reloaded root = %q, want %q", second.Root, expectedRoot)
	}
	if !bytes.Equal(second.Secret, first.Secret) {
		t.Fatal("reloading changed the durable cursor key")
	}
}

func TestLoadLocalUsesDefaultRootForAbsentOrEmptyEnvironment(t *testing.T) {
	testCases := []struct {
		name  string
		empty bool
	}{
		{name: "absent"},
		{name: "empty", empty: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			workingDirectory := t.TempDir()
			t.Chdir(workingDirectory)
			if testCase.empty {
				t.Setenv(runRootEnvironment, "")
			} else {
				unsetEnvironment(t, runRootEnvironment)
			}

			expectedRoot, err := filepath.Abs(filepath.Join(".bateia", "runs"))
			if err != nil {
				t.Fatalf("resolve expected default root: %v", err)
			}
			config, err := durableconfig.LoadLocal()
			if err != nil {
				t.Fatalf("LoadLocal: %v", err)
			}
			if config.Root != expectedRoot {
				t.Fatalf("root = %q, want default %q", config.Root, expectedRoot)
			}
			assertConfigDefaults(t, config)
		})
	}
}

func TestLoadLocalRejectsSymlinkRoot(t *testing.T) {
	requirePOSIXFilesystem(t, "symlinks")
	parent := t.TempDir()
	realRoot := filepath.Join(parent, "real-root")
	if err := os.Mkdir(realRoot, 0o700); err != nil {
		t.Fatalf("create real root: %v", err)
	}
	symlinkRoot := filepath.Join(parent, "symlink-root")
	if err := os.Symlink(realRoot, symlinkRoot); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	requireLoadErrorAt(t, symlinkRoot)
}

func TestLoadLocalRejectsSymlinkIntermediateComponent(t *testing.T) {
	requirePOSIXFilesystem(t, "symlinks")
	parent := t.TempDir()
	realIntermediate := filepath.Join(parent, "real-intermediate")
	if err := os.Mkdir(realIntermediate, 0o700); err != nil {
		t.Fatalf("create real intermediate directory: %v", err)
	}
	symlinkIntermediate := filepath.Join(parent, "symlink-intermediate")
	if err := os.Symlink(realIntermediate, symlinkIntermediate); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	requireLoadErrorAt(t, filepath.Join(symlinkIntermediate, "runs"))
}

func TestLoadLocalRejectsSymlinkCursorKey(t *testing.T) {
	requirePOSIXFilesystem(t, "symlinks")
	parent := t.TempDir()
	root := filepath.Join(parent, "runs")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	target := filepath.Join(parent, "external-key")
	writeFileWithMode(t, target, bytes.Repeat([]byte{0x51}, cursorKeyBytes), 0o600)
	if err := os.Symlink(target, filepath.Join(root, cursorKeyName)); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	requireLoadErrorAt(t, root)
}

func TestLoadLocalRejectsRootThatIsAFile(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runs")
	writeFileWithMode(t, root, []byte("not a directory"), 0o600)

	requireLoadErrorAt(t, root)
}

func TestLoadLocalRejectsInvalidCursorKeyKindsAndSizes(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, keyPath string)
	}{
		{
			name: "short",
			setup: func(t *testing.T, keyPath string) {
				writeFileWithMode(t, keyPath, bytes.Repeat([]byte{0x31}, cursorKeyBytes-1), 0o600)
			},
		},
		{
			name: "long",
			setup: func(t *testing.T, keyPath string) {
				writeFileWithMode(t, keyPath, bytes.Repeat([]byte{0x33}, cursorKeyBytes+1), 0o600)
			},
		},
		{
			name: "non-regular-directory",
			setup: func(t *testing.T, keyPath string) {
				if err := os.Mkdir(keyPath, 0o700); err != nil {
					t.Fatalf("create directory at cursor key path: %v", err)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "runs")
			if err := os.Mkdir(root, 0o700); err != nil {
				t.Fatalf("create root: %v", err)
			}
			testCase.setup(t, filepath.Join(root, cursorKeyName))

			requireLoadErrorAt(t, root)
		})
	}
}

func TestLoadLocalRejectsHardlinkedCursorWithoutChangingExternalMode(t *testing.T) {
	requirePOSIXFilesystem(t, "hardlinks and POSIX modes")
	parent := t.TempDir()
	root := filepath.Join(parent, "runs")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	keyPath := filepath.Join(root, cursorKeyName)
	keyBytes := bytes.Repeat([]byte{0x64}, cursorKeyBytes)
	writeFileWithMode(t, keyPath, keyBytes, 0o644)
	externalLink := filepath.Join(parent, "external-link")
	if err := os.Link(keyPath, externalLink); err != nil {
		t.Skipf("hardlinks unavailable: %v", err)
	}

	requireLoadErrorAt(t, root)
	assertFileModeAndContents(t, externalLink, 0o644, keyBytes)
}

func TestLoadLocalRejectsCursorLinkedFromAnotherFile(t *testing.T) {
	requirePOSIXFilesystem(t, "hardlinks and POSIX modes")
	parent := t.TempDir()
	root := filepath.Join(parent, "runs")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	externalFile := filepath.Join(parent, "external-key")
	keyBytes := bytes.Repeat([]byte{0x65}, cursorKeyBytes)
	writeFileWithMode(t, externalFile, keyBytes, 0o640)
	if err := os.Link(externalFile, filepath.Join(root, cursorKeyName)); err != nil {
		t.Skipf("hardlinks unavailable: %v", err)
	}

	requireLoadErrorAt(t, root)
	assertFileModeAndContents(t, externalFile, 0o640, keyBytes)
}

func TestLoadLocalNormalizesRegularCursorKeyMode(t *testing.T) {
	requirePOSIXFilesystem(t, "POSIX modes")
	root := filepath.Join(t.TempDir(), "runs")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}
	keyPath := filepath.Join(root, cursorKeyName)
	keyBytes := bytes.Repeat([]byte{0x72}, cursorKeyBytes)
	writeFileWithMode(t, keyPath, keyBytes, 0o644)
	t.Setenv(runRootEnvironment, root)

	config, err := durableconfig.LoadLocal()
	if err != nil {
		t.Fatalf("LoadLocal: %v", err)
	}
	if !bytes.Equal(config.Secret, keyBytes) {
		t.Fatal("mode normalization changed the cursor key bytes")
	}
	assertFileModeAndContents(t, keyPath, 0o600, keyBytes)
}

func TestLoadLocalConcurrentCreationUsesOneWinningKey(t *testing.T) {
	root := filepath.Join(t.TempDir(), "new", "runs")
	t.Setenv(runRootEnvironment, root)

	const callers = 24
	type result struct {
		config durableconfig.Config
		err    error
	}
	results := make(chan result, callers)
	start := make(chan struct{})
	var workers sync.WaitGroup
	workers.Add(callers)
	for range callers {
		go func() {
			defer workers.Done()
			<-start
			config, err := durableconfig.LoadLocal()
			results <- result{config: config, err: err}
		}()
	}
	close(start)
	workers.Wait()
	close(results)

	var winningSecret []byte
	for current := range results {
		if current.err != nil {
			t.Errorf("concurrent LoadLocal: %v", current.err)
			continue
		}
		assertConfigDefaults(t, current.config)
		if current.config.Root != root {
			t.Errorf("concurrent root = %q, want %q", current.config.Root, root)
		}
		if len(current.config.Secret) != cursorKeyBytes {
			t.Errorf("concurrent secret length = %d, want %d", len(current.config.Secret), cursorKeyBytes)
			continue
		}
		if winningSecret == nil {
			winningSecret = append([]byte(nil), current.config.Secret...)
			continue
		}
		if !bytes.Equal(current.config.Secret, winningSecret) {
			t.Error("concurrent callers observed different cursor keys")
		}
	}

	persisted, err := os.ReadFile(filepath.Join(root, cursorKeyName))
	if err != nil {
		t.Fatalf("read winning cursor key: %v", err)
	}
	if !bytes.Equal(persisted, winningSecret) {
		t.Fatal("persisted cursor key differs from the key observed by callers")
	}
}

func assertConfigDefaults(t *testing.T, config durableconfig.Config) {
	t.Helper()
	if config.Limits.MaxArtifactBytes != 4<<20 {
		t.Errorf("MaxArtifactBytes = %d, want %d", config.Limits.MaxArtifactBytes, int64(4<<20))
	}
	if config.Limits.MaxRunBytes != 16<<20 {
		t.Errorf("MaxRunBytes = %d, want %d", config.Limits.MaxRunBytes, int64(16<<20))
	}
	if config.Limits.MaxTotalBytes != 256<<20 {
		t.Errorf("MaxTotalBytes = %d, want %d", config.Limits.MaxTotalBytes, int64(256<<20))
	}
	if config.Workers != 2 {
		t.Errorf("Workers = %d, want 2", config.Workers)
	}
	if config.Capacity != 16 {
		t.Errorf("Capacity = %d, want 16", config.Capacity)
	}
}

func requireLoadErrorAt(t *testing.T, root string) {
	t.Helper()
	t.Setenv(runRootEnvironment, root)
	if _, err := durableconfig.LoadLocal(); err == nil {
		t.Fatal("LoadLocal succeeded, want an error")
	}
}

func requirePOSIXFilesystem(t *testing.T, capability string) {
	t.Helper()
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		t.Skipf("%s unavailable or not POSIX-compatible on %s", capability, runtime.GOOS)
	}
}

func unsetEnvironment(t *testing.T, name string) {
	t.Helper()
	value, present := os.LookupEnv(name)
	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("unset %s: %v", name, err)
	}
	t.Cleanup(func() {
		var err error
		if present {
			err = os.Setenv(name, value)
		} else {
			err = os.Unsetenv(name)
		}
		if err != nil {
			t.Errorf("restore %s: %v", name, err)
		}
	})
}

func writeFileWithMode(t *testing.T, path string, contents []byte, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, contents, mode); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, mode); err != nil {
			t.Fatalf("chmod %q: %v", path, err)
		}
	}
}

func assertFileModeAndContents(t *testing.T, path string, expectedMode os.FileMode, expectedContents []byte) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %q: %v", path, err)
	}
	if info.Mode().Perm() != expectedMode {
		t.Errorf("%q mode = %04o, want %04o", path, info.Mode().Perm(), expectedMode)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if !bytes.Equal(contents, expectedContents) {
		t.Errorf("%q contents changed", path)
	}
}
