/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_catalog_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type cliCatalog struct{ calls int }

func (c *cliCatalog) Capabilities() domain.ProviderCapabilities {
	return domain.ProviderCapabilities{Provider: "github", Roles: []domain.DiscoveryProviderRole{domain.RoleProjectCatalog}}
}
func (c *cliCatalog) SearchCatalog(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
	c.calls++
	return domain.CatalogDiscovery{Provider: "github", Role: domain.RoleProjectCatalog, SearchUnit: "repository", Attempt: domain.CatalogAttempt{Query: r.Query, Status: "SUCCEEDED"}}, nil
}
func TestCatalogCLIRequiresOutputBeforeNetwork(t *testing.T) {
	p := &cliCatalog{}
	err := RunIdeaCatalogCommand([]string{"discover", "--plan", "missing"}, &bytes.Buffer{}, p)
	if err == nil || p.calls != 0 {
		t.Fatalf("unsafe CLI: %v %d", err, p.calls)
	}
}
func TestCatalogScreenIsOffline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.json")
	b, _ := json.Marshal(domain.CatalogSnapshot{Candidates: []domain.CatalogCandidate{{Repository: domain.RepositoryCandidate{Owner: "o", Repo: "r"}}}})
	_ = os.WriteFile(path, b, 0600)
	var out bytes.Buffer
	if err := RunIdeaCatalogCommand([]string{"screen", "--snapshot", path}, &out, nil); err != nil {
		t.Fatal(err)
	}
	var report domain.CatalogScreeningReport
	if json.Unmarshal(out.Bytes(), &report) != nil || len(report.Decisions) != 1 {
		t.Fatal("invalid report")
	}
}
