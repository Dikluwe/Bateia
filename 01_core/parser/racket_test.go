/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/racket.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseInfoRkt_Basic(t *testing.T) {
	content := `#lang info
(define collection "spiffy-app")
(define deps '("base" "rackunit-lib" "web-server-lib"))
`

	ref := domain.CommitRef{Owner: "racket", Repo: "spiffy-app"}
	manifest, err := ParseInfoRkt(ref, content)
	if err != nil {
		t.Fatalf("ParseInfoRkt falhou: %v", err)
	}

	if manifest.Language != "racket" {
		t.Errorf("Esperado Language 'racket', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (rackunit-lib, web-server-lib), obtido %d", len(manifest.Dependencies))
	}
}
