/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/pony.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseCorralJson_Basic(t *testing.T) {
	content := `{
  "deps": [
    {
      "locator": "github.com/ponylang/net_ssl.git"
    }
  ]
}`

	ref := domain.CommitRef{Owner: "ponylang", Repo: "pony_app"}
	manifest, err := ParseCorralJson(ref, content)
	if err != nil {
		t.Fatalf("ParseCorralJson falhou: %v", err)
	}

	if manifest.Language != "pony" {
		t.Errorf("Esperado Language 'pony', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 1 {
		t.Fatalf("Esperado 1 dependência (net_ssl.git), obtido %d", len(manifest.Dependencies))
	}
}
