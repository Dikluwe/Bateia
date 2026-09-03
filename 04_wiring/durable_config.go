/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/durable_config.md
 * @layer L4
 * @updated 2026-09-01
 */

package main

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

	"bateia/03_infra/runstore"
)

const durableCursorSecretBytes = 32

type durableLocalConfig struct {
	root     string
	secret   []byte
	limits   runstore.RunStoreLimits
	workers  int
	capacity int
}

func loadDurableLocalConfig() (durableLocalConfig, error) {
	root := os.Getenv("BATEIA_RUN_ROOT")
	if root == "" {
		root = filepath.Join(".bateia", "runs")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return durableLocalConfig{}, fmt.Errorf("invalid durable run root")
	}
	if err := prepareDurableRoot(root); err != nil {
		return durableLocalConfig{}, err
	}
	if err := auditDurableConfigRoot(root); err != nil {
		return durableLocalConfig{}, err
	}
	secret, err := loadOrCreateCursorSecret(filepath.Join(root, "cursor.key"))
	if err != nil {
		return durableLocalConfig{}, err
	}
	return durableLocalConfig{
		root: root, secret: secret,
		limits:  runstore.RunStoreLimits{MaxArtifactBytes: 4 << 20, MaxRunBytes: 16 << 20, MaxTotalBytes: 256 << 20},
		workers: 2, capacity: 16,
	}, nil
}

func prepareDurableRoot(root string) error {
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

func auditDurableConfigRoot(root string) error {
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
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err == nil {
		file := os.NewFile(uintptr(fd), path)
		if file == nil {
			_ = syscall.Close(fd)
			return nil, fmt.Errorf("durable cursor secret is invalid")
		}
		defer file.Close()
		info, statErr := file.Stat()
		if statErr != nil || !info.Mode().IsRegular() {
			return nil, fmt.Errorf("durable cursor secret is invalid")
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); !ok || stat.Nlink != 1 {
			return nil, fmt.Errorf("durable cursor secret is invalid")
		}
		encoded, readErr := io.ReadAll(io.LimitReader(file, durableCursorSecretBytes+1))
		if readErr != nil || len(encoded) != durableCursorSecretBytes {
			return nil, fmt.Errorf("durable cursor secret is invalid")
		}
		if err := file.Chmod(0o600); err != nil {
			return nil, fmt.Errorf("durable run store unavailable")
		}
		return encoded, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("durable run store unavailable")
	}
	secret := make([]byte, durableCursorSecretBytes)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("durable cursor secret unavailable")
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return loadOrCreateCursorSecret(path)
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
