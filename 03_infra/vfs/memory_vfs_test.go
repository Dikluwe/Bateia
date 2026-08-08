/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/vfs/memory_vfs.md
 * @layer L3
 * @updated 2026-08-08
 */

package vfs_test

import (
	"sync"
	"testing"

	"bateia/03_infra/vfs"
)

func TestMemoryVFS_WriteReadRoundTrip(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	data := []byte("module github.com/test/repo\n\ngo 1.21")

	if err := fs.WriteFile("go.mod", data); err != nil {
		t.Fatalf("WriteFile falhou: %v", err)
	}

	got, err := fs.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("ReadFile falhou: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("Conteúdo diverge: got=%q, want=%q", got, data)
	}
}

func TestMemoryVFS_ReadNotFound(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	_, err := fs.ReadFile("inexistente.txt")
	if err == nil {
		t.Fatal("Esperava erro ao ler arquivo inexistente, mas não recebeu")
	}
}

func TestMemoryVFS_OverwriteFile(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	_ = fs.WriteFile("config.yaml", []byte("v1"))
	_ = fs.WriteFile("config.yaml", []byte("v2"))

	got, _ := fs.ReadFile("config.yaml")
	if string(got) != "v2" {
		t.Errorf("Esperava 'v2' após overwrite, recebeu: %q", got)
	}
}

func TestMemoryVFS_IsolationCopy(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	original := []byte("imutavel")
	_ = fs.WriteFile("f.txt", original)

	// Mutação do slice original NÃO deve afetar o VFS
	original[0] = 'X'
	got, _ := fs.ReadFile("f.txt")
	if got[0] == 'X' {
		t.Error("VFS não isolou a cópia no Write — mutação externa propagou")
	}

	// Mutação do slice retornado NÃO deve afetar o VFS
	got[0] = 'Y'
	got2, _ := fs.ReadFile("f.txt")
	if got2[0] == 'Y' {
		t.Error("VFS não isolou a cópia no Read — mutação externa propagou")
	}
}

func TestMemoryVFS_ListNodes(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	_ = fs.WriteFile("go.mod", []byte("module test"))
	_ = fs.WriteFile("main.go", []byte("package main"))
	_ = fs.WriteFile("README.md", []byte("# Test"))

	nodes := fs.ListNodes()
	if len(nodes) != 3 {
		t.Fatalf("Esperava 3 nós, recebeu: %d", len(nodes))
	}

	pathSet := make(map[string]bool)
	for _, n := range nodes {
		pathSet[n.Path] = true
		if n.Type != "file" {
			t.Errorf("Type esperado 'file', recebeu: %q para %s", n.Type, n.Path)
		}
		if n.Size <= 0 {
			t.Errorf("Size esperado > 0 para %s, recebeu: %d", n.Path, n.Size)
		}
	}

	for _, expected := range []string{"go.mod", "main.go", "README.md"} {
		if !pathSet[expected] {
			t.Errorf("Arquivo %s ausente no ListNodes", expected)
		}
	}
}

func TestMemoryVFS_ListNodesEmpty(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	nodes := fs.ListNodes()
	if len(nodes) != 0 {
		t.Errorf("Esperava 0 nós em VFS vazio, recebeu: %d", len(nodes))
	}
}

func TestMemoryVFS_ConcurrentAccess(t *testing.T) {
	fs := vfs.NewMemoryVFS()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = fs.WriteFile("shared.txt", []byte("data"))
			_, _ = fs.ReadFile("shared.txt")
			_ = fs.ListNodes()
		}(i)
	}
	wg.Wait()
}
