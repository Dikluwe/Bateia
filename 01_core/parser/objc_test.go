/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/objc_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParsePodfile_Basic(t *testing.T) {
	content := `platform :ios, '14.0'
use_frameworks!

target 'MyApp' do
  pod 'AFNetworking', '~> 4.0'
  pod 'SDWebImage', '~> 5.0'
end`

	ref := domain.CommitRef{Owner: "cocoapods", Repo: "MyApp"}
	manifest, err := ParsePodfile(ref, content)
	if err != nil {
		t.Fatalf("ParsePodfile falhou: %v", err)
	}

	if manifest.Language != "objective-c" {
		t.Errorf("Esperado Language 'objective-c', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 pods (AFNetworking, SDWebImage), obtido %d", len(manifest.Dependencies))
	}
}
