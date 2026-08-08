/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/vfs/memory_vfs.md
 * @layer L3
 * @updated 2026-08-08
 */

package vfs

import (
	"fmt"
	"sync"

	"bateia/01_core/domain"
)

// MemoryVFS abstrai um sistema de arquivos virtual completamente mantido em memória em L3.
type MemoryVFS struct {
	mu    sync.RWMutex
	files map[string][]byte
}

// NewMemoryVFS instancia uma nova VFS em memória.
func NewMemoryVFS() *MemoryVFS {
	return &MemoryVFS{
		files: make(map[string][]byte),
	}
}

// WriteFile armazena um arquivo virtual.
func (vfs *MemoryVFS) WriteFile(path string, data []byte) error {
	vfs.mu.Lock()
	defer vfs.mu.Unlock()
	buf := make([]byte, len(data))
	copy(buf, data)
	vfs.files[path] = buf
	return nil
}

// ReadFile lê o conteúdo de um arquivo virtual.
func (vfs *MemoryVFS) ReadFile(path string) ([]byte, error) {
	vfs.mu.RLock()
	defer vfs.mu.RUnlock()
	data, ok := vfs.files[path]
	if !ok {
		return nil, fmt.Errorf("arquivo virtual %s não encontrado", path)
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	return buf, nil
}

// ListNodes retorna a árvore de diretórios simulada do VFS.
func (vfs *MemoryVFS) ListNodes() []domain.DirectoryNode {
	vfs.mu.RLock()
	defer vfs.mu.RUnlock()

	var nodes []domain.DirectoryNode
	for p := range vfs.files {
		nodes = append(nodes, domain.DirectoryNode{
			Path: p,
			Type: "file",
			Size: int64(len(vfs.files[p])),
		})
	}
	return nodes
}
