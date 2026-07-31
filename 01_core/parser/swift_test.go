/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/swift.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePackageSwift_Basic(t *testing.T) {
	content := `// swift-tools-version:5.5
import PackageDescription

let package = Package(
    name: "VaporApp",
    dependencies: [
        .package(url: "https://github.com/vapor/vapor.git", from: "4.0.0"),
        .package(url: "https://github.com/vapor/fluent.git", from: "4.0.0")
    ]
)`

	ref := domain.CommitRef{Owner: "vapor", Repo: "VaporApp"}
	manifest, err := ParsePackageSwift(ref, content)
	if err != nil {
		t.Fatalf("ParsePackageSwift falhou: %v", err)
	}

	if manifest.ModuleName != "VaporApp" {
		t.Errorf("Esperado ModuleName 'VaporApp', obtido '%s'", manifest.ModuleName)
	}

	if manifest.Language != "swift" {
		t.Errorf("Esperado Language 'swift', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 dependências (vapor, fluent), obtido %d", len(manifest.Dependencies))
	}
}
