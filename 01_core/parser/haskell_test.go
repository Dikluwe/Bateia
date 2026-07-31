/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/haskell.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseHaskellCabal_Basic(t *testing.T) {
	content := `name:                servant-app
version:             0.1.0.0
build-type:          Simple
cabal-version:       >=1.10

executable servant-app
  main-is:             Main.hs
  build-depends:       base >= 4.7 && < 5,
                       servant,
                       aeson
`

	ref := domain.CommitRef{Owner: "haskell-servant", Repo: "servant-app"}
	manifest, err := ParseHaskellCabal(ref, content)
	if err != nil {
		t.Fatalf("ParseHaskellCabal falhou: %v", err)
	}

	if manifest.ModuleName != "servant-app" {
		t.Errorf("Esperado ModuleName 'servant-app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "haskell" {
		t.Errorf("Esperado Language 'haskell', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (servant, aeson), obtido %d", len(manifest.Dependencies))
	}
}
