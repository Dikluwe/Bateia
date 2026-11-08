/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_sites_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type insightSiteCLIFixture struct {
	calls int
	name  string
}

func (*insightSiteCLIFixture) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGiteeCatalog, Version: "fixture", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"catalog"}, PublicOnly: true, Replayable: true}
}
func (p *insightSiteCLIFixture) Discover(_ context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	p.calls++
	name := p.name
	if name == "" {
		name = "org/repo"
	}
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorGiteeCatalog, ExternalID: "org/repo", CanonicalURL: "https://gitee.com/org/repo"}, Name: name, Summary: "queue cancellation reference", ProductKind: "repository"}}}, nil
}
func TestInsightSiteCLIOptInAndLegacyFactoryIsolation(t *testing.T) {
	ports := &discoveryCLIPorts{}
	site := &insightSiteCLIFixture{}
	factoryCalls := 0
	factory := func() ([]domain.InsightSiteBinding, error) {
		factoryCalls++
		return []domain.InsightSiteBinding{{Connector: site}}, nil
	}
	out := filepath.Join(t.TempDir(), "out.json")
	var stdout bytes.Buffer
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--max-calls", "2", "--output", out}, &stdout, ports, ports, factory); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(out)
	var report domain.InsightSiteReport
	if err := json.Unmarshal(data, &report); err != nil || report.SchemaVersion != "insight-site-discovery-v1" || report.Calls != 2 || factoryCalls != 1 || ports.calls != 0 || stdout.Len() != 0 {
		t.Fatal("multisite command not isolated")
	}
	unused := func() ([]domain.InsightSiteBinding, error) {
		t.Error("legacy invoked site factory")
		return nil, errors.New("not used")
	}
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--output", out}, &stdout, ports, ports, unused); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(out)
	if !bytes.Contains(data, []byte("insight-discovery-v1")) {
		t.Fatal("legacy schema changed")
	}
}
func TestInsightSiteCLIValidatesBeforeFactoryAndBlocksNewTokens(t *testing.T) {
	for _, flags := range [][]string{{"--max-calls", "0"}, {"--sources", "all", "--max-calls", "61"}, {"--sources", "all", "--locale", "bad_locale"}, {"--sources", ""}} {
		called := false
		factory := func() ([]domain.InsightSiteBinding, error) { called = true; return nil, nil }
		args := append([]string{"--question", "queue", "--output", filepath.Join(t.TempDir(), "out.json")}, flags...)
		if err := RunInsightDiscoveryCommand(args, &bytes.Buffer{}, nil, nil, factory); err == nil || called {
			t.Fatal("invalid flags called factory")
		}
	}
	for _, key := range []string{"BATEIA_GITEE_TOKEN", "BATEIA_ATOMGIT_TOKEN"} {
		t.Run(key, func(t *testing.T) {
			secret := "test<token>&\"not-real-credential"
			t.Setenv(key, secret)
			site := &insightSiteCLIFixture{name: "name " + secret}
			out := filepath.Join(t.TempDir(), "out.json")
			if err := os.WriteFile(out, []byte("preserve"), 0600); err != nil {
				t.Fatal(err)
			}
			factory := func() ([]domain.InsightSiteBinding, error) {
				return []domain.InsightSiteBinding{{Connector: site}}, nil
			}
			err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--output", out}, &bytes.Buffer{}, nil, nil, factory)
			saved, _ := os.ReadFile(out)
			if err == nil || string(saved) != "preserve" || strings.Contains(err.Error(), secret) {
				t.Fatal("new token escaped protection")
			}
		})
	}
}
