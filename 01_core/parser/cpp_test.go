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

func TestParseCMakeListsTxt_Basic(t *testing.T) {
	content := `
cmake_minimum_required(VERSION 3.14)
project(my_engine)

find_package(Boost REQUIRED COMPONENTS system)
find_package(OpenCV REQUIRED)
`

	ref := domain.CommitRef{Owner: "cpp", Repo: "my_engine"}
	manifest, err := ParseCMakeListsTxt(ref, content)
	if err != nil {
		t.Fatalf("ParseCMakeListsTxt falhou: %v", err)
	}

	if manifest.ModuleName != "my_engine" {
		t.Errorf("Esperado ModuleName 'my_engine', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "cpp" {
		t.Errorf("Esperado Language 'cpp', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}

func TestParseVcpkgJson_Basic(t *testing.T) {
	content := `{
		"name": "my-game-app",
		"dependencies": [
			"fmt",
			{"name": "spdlog", "version>=": "1.10.0"}
		]
	}`

	ref := domain.CommitRef{Owner: "game", Repo: "my-game-app"}
	manifest, err := ParseVcpkgJson(ref, content)
	if err != nil {
		t.Fatalf("ParseVcpkgJson falhou: %v", err)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências, obtido %d", len(manifest.Dependencies))
	}
}
