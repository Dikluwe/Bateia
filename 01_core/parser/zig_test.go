/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/zig_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseBuildZigZon_Basic(t *testing.T) {
	content := `.{
    .name = "httpz_app",
    .version = "0.0.0",
    .dependencies = .{
        .httpz = .{
            .url = "https://github.com/karlseguin/httpz/archive/master.tar.gz",
        },
    },
}`

	ref := domain.CommitRef{Owner: "zig-lang", Repo: "httpz_app"}
	manifest, err := ParseBuildZigZon(ref, content)
	if err != nil {
		t.Fatalf("ParseBuildZigZon falhou: %v", err)
	}

	if manifest.ModuleName != "httpz_app" {
		t.Errorf("Esperado ModuleName 'httpz_app', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "zig" {
		t.Errorf("Esperado Language 'zig', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 1 {
		t.Fatalf("Esperado 1 dependência (httpz), obtido %d", len(manifest.Dependencies))
	}
}
