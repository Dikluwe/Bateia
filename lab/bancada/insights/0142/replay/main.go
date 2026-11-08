// Offline workflow demo: frozen parent READMEs and explicitly synthetic child documents.
package main

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type clock struct{}

func (clock) Now() time.Time { return time.Now() }
func (clock) WithTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

type catalog struct{ values []domain.ProductCandidate }

func (c catalog) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorGitHubCatalog, Version: "offline-synthetic-followup-0142", Kind: domain.ConnectorRepositoryCatalog, QueryLocales: []string{"en-US"}, IdentityLevels: []string{"owner_repo"}, EvidenceSurfaces: []string{"fixture"}, PublicOnly: true, Replayable: true}
}
func (c catalog) Discover(context.Context, domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, Candidates: c.values}, nil
}

type reader struct {
	parents  map[string]domain.EvidenceBundle
	resolves []string
	reads    []string
}

func (r *reader) ResolveRepository(_ context.Context, id domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	r.resolves = append(r.resolves, id.ExternalID)
	parts := strings.Split(id.ExternalID, "/")
	return domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: id.CanonicalURL, CommitID: strings.Repeat("c", 40)}, nil
}
func (r *reader) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
}
func (r *reader) ExtractRepoFiles(_ context.Context, repo domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	id := repo.Owner + "/" + repo.Repo
	r.reads = append(r.reads, id)
	ref := domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}
	text := "# Synthetic destination fixture\nThis is authored test data for " + id + ", not remotely observed repository contents.\nIt exercises memory and bounded reading only.\n"
	if parent, ok := r.parents[id]; ok {
		ref = parent.Repository
		text = parent.Items[0].Content
	}
	return &domain.ExtractedRepository{Ref: ref, Files: []domain.ExtractedFile{{Path: "README.md", Content: text}}}, nil
}
func main() {
	raw, err := os.ReadFile("lab/bancada/insights/comparacao-ia-2026-10-28/bateia-result.json")
	if err != nil {
		panic(err)
	}
	var frozen domain.InsightSiteReport
	if json.Unmarshal(raw, &frozen) != nil {
		panic("invalid frozen source")
	}
	values := []domain.ProductCandidate{}
	r := &reader{parents: map[string]domain.EvidenceBundle{}}
	for _, c := range frozen.Candidates {
		if c.Selected && c.Product.Identity.ConnectorID == domain.ConnectorGitHubCatalog {
			values = append(values, c.Product)
		}
	}
	for _, b := range frozen.Bundles {
		r.parents[b.Repository.Owner+"/"+b.Repository.Repo] = b
	}
	request := frozen.Plan.Request
	request.Sources = []string{"github"}
	request.MaxCalls = 5
	request.Discovery.MaxReferences = 2
	request.Discovery.MaxPassages = 6
	bindings := []domain.InsightSiteBinding{{Connector: catalog{values}, Reader: r, FollowReader: r}}
	first, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), clock{}), request, bindings)
	if err != nil {
		panic(err)
	}
	if len(r.resolves) != 2 || r.resolves[0] != "chnm/serendipomatic" {
		panic("first ranked trail not deepened")
	}
	firstDestinations := append([]string{}, r.resolves...)
	request.Memory = first.ExplorationMemory
	r.resolves = nil
	r.reads = nil
	second, err := usecase.DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), clock{}), request, bindings)
	if err != nil {
		panic(err)
	}
	if len(r.resolves) != 2 || len(r.reads) != 2 {
		panic("memory did not advance")
	}
	for _, id := range r.reads {
		for _, old := range firstDestinations {
			if id == old {
				panic("repeated examined target")
			}
		}
		if _, ok := r.parents[id]; ok {
			panic("repeated parent")
		}
	}
	discarded := ""
	var memory domain.InsightExplorationMemory
	for _, entry := range second.ExplorationMemory.Entries {
		if entry.State == "PENDING" {
			discarded = entry.Identity.CanonicalURL
			memory, err = usecase.UpdateInsightExplorationMemory(*second.ExplorationMemory, discarded, "discarded", "Synthetic workflow decision; pursue another path")
			break
		}
	}
	if err != nil || discarded == "" {
		panic("feedback fixture failed")
	}
	output := map[string]any{"mode": "OFFLINE_WORKFLOW_FIXTURE", "network": false, "child_documents": "SYNTHETIC_NOT_PUBLIC_OBSERVATIONS", "candidate_code_executed": false, "first": first, "second": second, "first_destinations": firstDestinations, "second_destinations": r.resolves, "feedback_memory": memory, "discarded_target": discarded}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		panic(fmt.Sprint(err))
	}
}
