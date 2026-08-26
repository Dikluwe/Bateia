/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/gomod_test.md
 * @layer L1
 * @updated 2026-07-25
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseGoMod(t *testing.T) {
	sampleGoMod := `module github.com/example/bateia-sample

go 1.22.0

require (
	github.com/google/go-cmp v0.6.0
	github.com/pkg/errors v0.9.1 // indirect
)

require github.com/spf13/cobra v1.8.0
`

	ref := domain.CommitRef{Owner: "example", Repo: "bateia-sample", Commit: "abc1234"}
	manifest, err := ParseGoMod(ref, sampleGoMod)
	if err != nil {
		t.Fatalf("ParseGoMod falhou inesperadamente: %v", err)
	}

	if manifest.ModuleName != "github.com/example/bateia-sample" {
		t.Errorf("Esperado módulo github.com/example/bateia-sample, obtido: %s", manifest.ModuleName)
	}

	if manifest.GoVersion != "1.22.0" {
		t.Errorf("Esperado go version 1.22.0, obtido: %s", manifest.GoVersion)
	}

	if len(manifest.Dependencies) != 3 {
		t.Fatalf("Esperado 3 dependências, obtido: %d", len(manifest.Dependencies))
	}

	if manifest.Dependencies[0].Name != "github.com/google/go-cmp" || manifest.Dependencies[0].Indirect {
		t.Errorf("Dependência 0 incorreta: %+v", manifest.Dependencies[0])
	}

	if manifest.Dependencies[1].Name != "github.com/pkg/errors" || !manifest.Dependencies[1].Indirect {
		t.Errorf("Dependência 1 (indireta) incorreta: %+v", manifest.Dependencies[1])
	}

	if manifest.Dependencies[2].Name != "github.com/spf13/cobra" {
		t.Errorf("Dependência 2 incorreta: %+v", manifest.Dependencies[2])
	}
}
