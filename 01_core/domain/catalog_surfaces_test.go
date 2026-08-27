/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/catalog_surfaces_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestSurfaceIdentityIsNominal(t *testing.T) {
	if SameSurface(SurfaceReadme, SurfaceTopics) || !SameSurface(SurfaceReadme, SurfaceReadme) {
		t.Fatal("surface mutation survived")
	}
}
func TestOpaqueSurfaceMatchCarriesUnknownInsteadOfFabricatedRepository(t *testing.T) {
	m := SurfaceMatch{Surface: SurfacePackage, Package: "opaque", UnknownReason: "source unavailable"}
	if m.Repository != "" || m.UnknownReason == "" {
		t.Fatal("opaque source promoted")
	}
}
func TestRepeatedQueriesDoNotBecomeIndependentSurfaces(t *testing.T) {
	a, b := SurfaceReadme, SurfaceReadme
	if !SameSurface(a, b) {
		t.Fatal("same surface counted twice")
	}
}
