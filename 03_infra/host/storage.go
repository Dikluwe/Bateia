/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/host/storage.md
 * @layer L3
 * @updated 2026-10-28
 */
package host

import (
	"bateia/01_core/domain"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const maxIdeaInputBytes = 16 << 20

func (*Host) ReadBounded(path string, limit int64) ([]byte, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid input byte budget")
	}
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return (&Host{}).ReadDescriptor(f, limit)
}

func (*Host) ReadDescriptor(f domain.CommandFile, limit int64) ([]byte, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("invalid input byte budget")
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("entrada deve ser arquivo regular de até %d bytes", limit)
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("entrada excedeu limite durante leitura")
	}
	return data, nil
}
func (*Host) PrepareDossier(path string, resume bool) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("diretório de saída vazio")
	}
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("diretório de saída não pode ser symlink")
		}
		if !info.IsDir() {
			return fmt.Errorf("saída exige diretório")
		}
		if !resume {
			dir, e := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
			if e != nil {
				return e
			}
			opened, e := dir.Stat()
			if e != nil || !opened.IsDir() || !os.SameFile(info, opened) {
				dir.Close()
				return fmt.Errorf("diretório de saída alterado durante validação")
			}
			entries, readErr := dir.ReadDir(1)
			closeErr := dir.Close()
			if readErr != nil && readErr != io.EOF {
				return readErr
			}
			if closeErr != nil {
				return closeErr
			}
			if len(entries) > 0 {
				return fmt.Errorf("diretório de saída já contém artefatos; use --resume")
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, 0755)
}
func (*Host) WriteAtomic(path string, data []byte, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".bateia-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(mode); err == nil {
		_, err = tmp.Write(data)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(name, path)
}
