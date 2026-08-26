/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/lua_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseLuaRockspec_Basic(t *testing.T) {
	content := `package = "lapis"
version = "1.9.0-1"
source = { url = "git://..." }
dependencies = {
   "lua >= 5.1",
   "ansicolors",
   "luafilesystem"
}
`

	ref := domain.CommitRef{Owner: "leafo", Repo: "lapis"}
	manifest, err := ParseLuaRockspec(ref, content)
	if err != nil {
		t.Fatalf("ParseLuaRockspec falhou: %v", err)
	}

	if manifest.ModuleName != "lapis" {
		t.Errorf("Esperado ModuleName 'lapis', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "lua" {
		t.Errorf("Esperado Language 'lua', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (ansicolors, luafilesystem), obtido %d", len(manifest.Dependencies))
	}
}
