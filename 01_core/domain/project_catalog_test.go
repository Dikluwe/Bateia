/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/project_catalog_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

import "testing"

func TestProviderRolesAreNominalAndNotInterchangeable(t *testing.T) {
	github := ProviderCapabilities{Provider: "github", Roles: []DiscoveryProviderRole{RoleProjectCatalog}, SearchUnits: []string{"repository"}}
	sourcegraph := ProviderCapabilities{Provider: "sourcegraph", Roles: []DiscoveryProviderRole{RoleCodeProof}, SearchUnits: []string{"code_occurrence"}}
	if !SupportsRole(github, RoleProjectCatalog) || SupportsRole(github, RoleCodeProof) {
		t.Fatal("GitHub role mutation survived")
	}
	if !SupportsRole(sourcegraph, RoleCodeProof) || SupportsRole(sourcegraph, RoleProjectCatalog) {
		t.Fatal("Sourcegraph role mutation survived")
	}
}

func TestUnknownIdentityCannotBeRepresentedAsCatalogSuccess(t *testing.T) {
	d := CatalogDecision{Repository: RepositoryCandidate{Owner: "org", Repo: "opaque"}, Verdict: CatalogReview}
	if d.Repository.CommitID != "" || d.Verdict != CatalogReview {
		t.Fatal("opaque identity was promoted")
	}
}
