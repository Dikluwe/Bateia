// Offline fixture replay of already collected public README bytes; no network.
package main

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type runtime struct{}

func (runtime) Now() time.Time { return time.Now() }
func (runtime) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

type catalog struct{ values []domain.ProductCandidate }

func (c catalog) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGitHubCatalog, Version: "frozen-readme-replay-0141", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"frozen_fixture"}, PublicOnly: true, Replayable: true}
}
func (c catalog) Discover(_ context.Context, _ domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: c.values}, nil
}

type reader struct{ bundles []domain.EvidenceBundle }

func (r *reader) FetchTree(_ context.Context, repo domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	for _, b := range r.bundles {
		if b.Repository.Owner == repo.Owner && b.Repository.Repo == repo.Repo && b.Repository.Commit == repo.CommitID {
			return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
		}
	}
	return nil, fmt.Errorf("fixture missing")
}
func (r *reader) ExtractRepoFiles(_ context.Context, repo domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	for _, b := range r.bundles {
		if b.Repository.Owner == repo.Owner && b.Repository.Repo == repo.Repo && b.Repository.Commit == repo.CommitID {
			return &domain.ExtractedRepository{Ref: b.Repository, Files: []domain.ExtractedFile{{Path: "README.md", Content: b.Items[0].Content}}}, nil
		}
	}
	return nil, fmt.Errorf("fixture missing")
}
func main() {
	raw, err := os.ReadFile("lab/bancada/insights/comparacao-ia-2026-10-28/bateia-result.json")
	if err != nil {
		panic(err)
	}
	var old domain.InsightSiteReport
	if err := json.Unmarshal(raw, &old); err != nil {
		panic(err)
	}
	values := []domain.ProductCandidate{}
	for _, c := range old.Candidates {
		if c.Selected && c.Product.Identity.ConnectorID == domain.ConnectorGitHubCatalog {
			values = append(values, c.Product)
		}
	}
	request := old.Plan.Request
	request.Sources = []string{"github"}
	request.MaxCalls = 5
	request.Discovery.MaxReferences = 2
	report, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), runtime{}), request, []domain.InsightSiteBinding{{Connector: catalog{values}, Reader: &reader{old.Bundles}}})
	if err != nil {
		panic(err)
	}
	found := false
	for _, lead := range report.ReferenceLeads {
		if lead.Target.ExternalID == "chnm/serendipomatic" {
			found = true
		}
	}
	if !found || len(report.ReferenceLeads) != 12 || len(report.Bundles) != 2 || len(report.Tasks[0].ReferenceLeadRefs) != 12 {
		panic("frozen collection reference or scope missing")
	}
	json.NewEncoder(os.Stdout).Encode(map[string]any{"mode": "OFFLINE_FROZEN_FIXTURE_REPLAY", "fresh_network": false, "candidate_code_executed": false, "report": report})
}
