/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/durableconfig/config.md
 * @layer L3
 * @updated 2026-09-18
 */

package durableconfig

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"bateia/03_infra/runstore"
)

const cursorSecretBytes = 32

type Config struct {
	Root     string
	Secret   []byte
	Limits   runstore.RunStoreLimits
	Workers  int
	Capacity int
}

func LoadLocal() (Config, error) {
	root := os.Getenv("BATEIA_RUN_ROOT")
	if root == "" {
		root = filepath.Join(".bateia", "runs")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return Config{}, fmt.Errorf("invalid durable run root")
	}
	if err := prepareRoot(root); err != nil {
		return Config{}, err
	}
	if err := auditRoot(root); err != nil {
		return Config{}, err
	}
	secret, err := loadOrCreateCursorSecret(filepath.Join(root, "cursor.key"))
	if err != nil {
		return Config{}, err
	}
	return Config{
		Root: root, Secret: secret,
		Limits:   runstore.RunStoreLimits{MaxArtifactBytes: 4 << 20, MaxRunBytes: 16 << 20, MaxTotalBytes: 256 << 20},
		Workers:  2,
		Capacity: 16,
	}, nil
}

func prepareRoot(root string) error {
	volume := filepath.VolumeName(root)
	current := string(os.PathSeparator)
	if volume != "" {
		current = volume + string(os.PathSeparator)
	}
	relative, err := filepath.Rel(current, root)
	if err != nil || relative == "." {
		return fmt.Errorf("invalid durable run root")
	}
	for _, component := range strings.Split(relative, string(os.PathSeparator)) {
		if component == "" || component == "." || component == ".." {
			return fmt.Errorf("invalid durable run root")
		}
		current = filepath.Join(current, component)
		info, statErr := os.Lstat(current)
		if errors.Is(statErr, fs.ErrNotExist) {
			if err := os.Mkdir(current, 0o700); err != nil {
				if errors.Is(err, fs.ErrExist) {
					info, statErr = os.Lstat(current)
					if statErr == nil && info.Mode()&os.ModeSymlink == 0 && info.IsDir() {
						continue
					}
				}
				return fmt.Errorf("durable run store unavailable")
			}
			continue
		}
		if statErr != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("durable run root is invalid")
		}
	}
	return nil
}

func auditRoot(root string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("durable run root is invalid")
		}
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("durable run root is invalid")
		}
		if entry.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("durable run root is invalid")
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); !ok || stat.Nlink != 1 {
			return fmt.Errorf("durable run root is invalid")
		}
		return nil
	})
}

func loadOrCreateCursorSecret(path string) ([]byte, error) {
	secret, state, err := readCursorSecret(path)
	if err != nil {
		return nil, err
	}
	if state == cursorSecretReady {
		return secret, nil
	}
	if state == cursorSecretIncomplete {
		return waitForCursorSecret(path)
	}

	secret = make([]byte, cursorSecretBytes)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("durable cursor secret unavailable")
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return waitForCursorSecret(path)
	}
	if err != nil {
		return nil, fmt.Errorf("durable run store unavailable")
	}
	written, writeErr := file.Write(secret)
	syncErr := file.Sync()
	closeErr := file.Close()
	if writeErr != nil || written != len(secret) || syncErr != nil || closeErr != nil {
		return nil, fmt.Errorf("durable run store unavailable")
	}
	return secret, nil
}

type cursorSecretState uint8

const (
	cursorSecretMissing cursorSecretState = iota
	cursorSecretIncomplete
	cursorSecretReady
)

func readCursorSecret(path string) ([]byte, cursorSecretState, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, cursorSecretMissing, nil
		}
		return nil, cursorSecretMissing, fmt.Errorf("durable run store unavailable")
	}
	file := os.NewFile(uintptr(fd), path)
	if file == nil {
		_ = syscall.Close(fd)
		return nil, cursorSecretMissing, fmt.Errorf("durable cursor secret is invalid")
	}
	defer file.Close()
	info, statErr := file.Stat()
	if statErr != nil || !info.Mode().IsRegular() {
		return nil, cursorSecretMissing, fmt.Errorf("durable cursor secret is invalid")
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); !ok || stat.Nlink != 1 {
		return nil, cursorSecretMissing, fmt.Errorf("durable cursor secret is invalid")
	}
	encoded, readErr := io.ReadAll(io.LimitReader(file, cursorSecretBytes+1))
	if readErr != nil {
		return nil, cursorSecretMissing, fmt.Errorf("durable cursor secret is invalid")
	}
	if len(encoded) != cursorSecretBytes {
		return nil, cursorSecretIncomplete, nil
	}
	if err := file.Chmod(0o600); err != nil {
		return nil, cursorSecretMissing, fmt.Errorf("durable run store unavailable")
	}
	return encoded, cursorSecretReady, nil
}

func waitForCursorSecret(path string) ([]byte, error) {
	const attempts = 100
	for attempt := 0; attempt < attempts; attempt++ {
		secret, state, err := readCursorSecret(path)
		if err != nil {
			return nil, err
		}
		if state == cursorSecretReady {
			return secret, nil
		}
		time.Sleep(time.Millisecond)
	}
	return nil, fmt.Errorf("durable cursor secret is invalid")
}
