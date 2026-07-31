/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cpp.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

// --- CMakeLists.txt ---

func TestParseCMakeListsTxt_CppExplicit(t *testing.T) {
	content := `
cmake_minimum_required(VERSION 3.14)
project(my_engine CXX)

find_package(Boost REQUIRED COMPONENTS system)
find_package(OpenCV REQUIRED)
`
	ref := domain.CommitRef{Owner: "cpp", Repo: "my_engine"}
	m, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}
	if m.Language != "cpp" {
		t.Errorf("project(CXX): esperado 'cpp', obtido '%s'", m.Language)
	}
	if m.ModuleName != "my_engine" {
		t.Errorf("ModuleName: esperado 'my_engine', obtido '%s'", m.ModuleName)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências, obtido %d", len(m.Dependencies))
	}
}

func TestParseCMakeListsTxt_CAndCXX(t *testing.T) {
	// Quando ambas as linguagens estão declaradas, o projeto é C++ (tem CXX)
	content := `
cmake_minimum_required(VERSION 3.14)
project(mixed C CXX)

find_package(OpenSSL REQUIRED)
`
	ref := domain.CommitRef{Owner: "org", Repo: "mixed"}
	m, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}
	if m.Language != "cpp" {
		t.Errorf("project(C CXX): esperado 'cpp', obtido '%s'", m.Language)
	}
}

func TestParseCMakeListsTxt_CPureOnly(t *testing.T) {
	// project(name C) sem CXX → deve produzir Language: "c"
	content := `
cmake_minimum_required(VERSION 3.10)
project(libuv C)

find_package(Threads REQUIRED)
`
	ref := domain.CommitRef{Owner: "libuv", Repo: "libuv"}
	m, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}
	if m.Language != "c" {
		t.Errorf("project(C): esperado 'c', obtido '%s'", m.Language)
	}
	if m.ModuleName != "libuv" {
		t.Errorf("ModuleName: esperado 'libuv', obtido '%s'", m.ModuleName)
	}
}

func TestParseCMakeListsTxt_NoLanguage(t *testing.T) {
	// project(name) sem linguagem → padrão histórico "cpp"
	content := `
cmake_minimum_required(VERSION 3.14)
project(my_engine)

find_package(Boost REQUIRED COMPONENTS system)
find_package(OpenCV REQUIRED)
`
	ref := domain.CommitRef{Owner: "cpp", Repo: "my_engine"}
	m, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}
	if m.Language != "cpp" {
		t.Errorf("project() sem linguagem: esperado 'cpp' (padrão), obtido '%s'", m.Language)
	}
}

func TestParseCMakeListsTxt_LanguagesKeyword(t *testing.T) {
	// Forma moderna: project(name VERSION x LANGUAGES C)
	content := `
cmake_minimum_required(VERSION 3.18)
project(libev VERSION 4.33 LANGUAGES C)

find_package(Threads REQUIRED)
`
	ref := domain.CommitRef{Owner: "org", Repo: "libev"}
	m, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}
	if m.Language != "c" {
		t.Errorf("LANGUAGES C: esperado 'c', obtido '%s'", m.Language)
	}
}

// --- vcpkg.json ---

func TestParseVcpkgJson_Basic(t *testing.T) {
	content := `{
		"name": "my-game-app",
		"dependencies": [
			"fmt",
			{"name": "spdlog", "version>=": "1.10.0"}
		]
	}`
	ref := domain.CommitRef{Owner: "game", Repo: "my-game-app"}
	m, err := ParseVcpkgJson(ref, content)
	if err != nil {
		t.Fatalf("ParseVcpkgJson falhou: %v", err)
	}
	if m.Language != "cpp" {
		t.Errorf("vcpkg.json: esperado 'cpp', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências, obtido %d", len(m.Dependencies))
	}
}

// --- conanfile.txt ---

func TestParseConanfileTxt_CppWithLibcxx(t *testing.T) {
	content := `
[settings]
os=Linux
compiler=gcc
compiler.version=11
compiler.libcxx=libstdc++11

[requires]
boost/1.79.0
`
	ref := domain.CommitRef{Owner: "org", Repo: "mylib"}
	m, err := ParseConanfileTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseConanfileTxt falhou: %v", err)
	}
	if m.Language != "cpp" {
		t.Errorf("com compiler.libcxx: esperado 'cpp', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 1 {
		t.Errorf("esperado 1 dependência, obtido %d", len(m.Dependencies))
	}
}

func TestParseConanfileTxt_CPureNoLibcxx(t *testing.T) {
	content := `
[settings]
os=Linux
compiler=gcc
compiler.version=12

[requires]
zlib/1.2.13
openssl/3.1.0
`
	ref := domain.CommitRef{Owner: "org", Repo: "clib"}
	m, err := ParseConanfileTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseConanfileTxt falhou: %v", err)
	}
	if m.Language != "c" {
		t.Errorf("sem compiler.libcxx: esperado 'c', obtido '%s'", m.Language)
	}
	if len(m.Dependencies) != 2 {
		t.Errorf("esperado 2 dependências, obtido %d", len(m.Dependencies))
	}
}
