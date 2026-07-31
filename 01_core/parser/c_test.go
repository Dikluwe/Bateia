/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/c.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseMakefile_Basic(t *testing.T) {
	content := `
CC = gcc
TARGET = myapp
CFLAGS = -Wall -O2
LDLIBS = -lpthread -lm
`
	ref := domain.CommitRef{Owner: "org", Repo: "myapp"}
	m, err := ParseMakefile(ref, content)
	if err != nil {
		t.Fatalf("ParseMakefile falhou: %v", err)
	}
	if m.ModuleName != "myapp" {
		t.Errorf("ModuleName: esperado 'myapp', obtido '%s'", m.ModuleName)
	}
	if m.Language != "c" {
		t.Errorf("Language: esperado 'c', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências, obtido %d", len(m.Dependencies))
	}
}

func TestParseMakefile_PROG(t *testing.T) {
	content := `
CC = gcc
PROG = server
LDFLAGS = -lssl -lcrypto
`
	ref := domain.CommitRef{Owner: "org", Repo: "fallback"}
	m, err := ParseMakefile(ref, content)
	if err != nil {
		t.Fatalf("ParseMakefile falhou: %v", err)
	}
	if m.ModuleName != "server" {
		t.Errorf("ModuleName via PROG: esperado 'server', obtido '%s'", m.ModuleName)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências (-lssl, -lcrypto), obtido %d", len(m.Dependencies))
	}
}

func TestParseMakefile_FallbackRepo(t *testing.T) {
	content := `
CC = gcc
LDFLAGS = -lsqlite3
`
	ref := domain.CommitRef{Owner: "org", Repo: "myrepo"}
	m, err := ParseMakefile(ref, content)
	if err != nil {
		t.Fatalf("ParseMakefile falhou: %v", err)
	}
	if m.ModuleName != "myrepo" {
		t.Errorf("Fallback para ref.Repo: esperado 'myrepo', obtido '%s'", m.ModuleName)
	}
}

func TestParseMakefile_Empty(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "empty"}
	m, err := ParseMakefile(ref, "")
	if err != nil {
		t.Fatalf("ParseMakefile falhou em conteúdo vazio: %v", err)
	}
	if m.Language != "c" {
		t.Errorf("Language: esperado 'c', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 0 {
		t.Errorf("esperado 0 dependências, obtido %d", len(m.Dependencies))
	}
}

func TestParseMesonBuild_Basic(t *testing.T) {
	content := `
project('libuv', 'c',
  version : '1.48.0',
  license : 'MIT')

threads_dep = dependency('threads')
dl_dep = dependency('dl', required : false)
`
	ref := domain.CommitRef{Owner: "libuv", Repo: "libuv"}
	m, err := ParseMesonBuild(ref, content)
	if err != nil {
		t.Fatalf("ParseMesonBuild falhou: %v", err)
	}
	if m.ModuleName != "libuv" {
		t.Errorf("ModuleName: esperado 'libuv', obtido '%s'", m.ModuleName)
	}
	if m.Language != "c" {
		t.Errorf("Language: esperado 'c', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências, obtido %d", len(m.Dependencies))
	}
}

func TestParseMesonBuild_NoDeps(t *testing.T) {
	content := `project('minimal', 'c')`
	ref := domain.CommitRef{Owner: "org", Repo: "minimal"}
	m, err := ParseMesonBuild(ref, content)
	if err != nil {
		t.Fatalf("ParseMesonBuild falhou: %v", err)
	}
	if m.ModuleName != "minimal" {
		t.Errorf("ModuleName: esperado 'minimal', obtido '%s'", m.ModuleName)
	}
	if len(m.Dependencies) != 0 {
		t.Errorf("esperado 0 dependências, obtido %d", len(m.Dependencies))
	}
}

func TestParseMesonBuild_FallbackRepo(t *testing.T) {
	content := `dependency('glib-2.0')`
	ref := domain.CommitRef{Owner: "org", Repo: "myrepo"}
	m, err := ParseMesonBuild(ref, content)
	if err != nil {
		t.Fatalf("ParseMesonBuild falhou: %v", err)
	}
	if m.ModuleName != "myrepo" {
		t.Errorf("Fallback para ref.Repo: esperado 'myrepo', obtido '%s'", m.ModuleName)
	}
}
