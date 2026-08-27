/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/catalog_surfaces_test.md
 * @layer L3
 * @updated 2026-08-26
 */
package github

import (
	"bateia/01_core/domain"
	"testing"
)

func TestReadmePathRestriction(t *testing.T) {
	for _, p := range []string{"README.md", "docs/index.md"} {
		if !isRootReadme(p) {
			t.Fatal(p)
		}
	}
	for _, p := range []string{"src/readme.go", "vendor/README.md", "README.test.md"} {
		if isRootReadme(p) {
			t.Fatalf("non-product doc accepted: %s", p)
		}
	}
}
func TestSurfaceCapabilitiesAreDistinct(t *testing.T) {
	c := NewClient("")
	if NewReadmeSurface(c).SurfaceCapabilities().ID != domain.SurfaceReadme || NewTopicsSurface(c).SurfaceCapabilities().ID != domain.SurfaceTopics {
		t.Fatal("surface substitution")
	}
}
