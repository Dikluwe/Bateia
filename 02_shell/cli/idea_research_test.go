/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_research_test.md
 * @layer L2
 * @updated 2026-10-06
 */

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

type researchConnectorStub struct {
	id      string
	failure bool
}

func (s researchConnectorStub) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: s.id, Version: "test"}
}
func (s researchConnectorStub) Discover(_ context.Context, r domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	if s.failure {
		return domain.ConnectorDiscovery{}, errors.New("provider offline")
	}
	if r.Query != "symbolic algebra" {
		panic("query was changed")
	}
	return domain.ConnectorDiscovery{Status: domain.AttemptSucceeded, ResponseHash: "fixture-hash", Candidates: []domain.ProductCandidate{{Identity: domain.ProductIdentity{ConnectorID: s.id, ExternalID: "123", CanonicalURL: "https://zbmath.org/software/123"}, Name: "Example", ScholarlyEvidence: &domain.ScholarlyEvidence{Catalogue: "zbmath", RecordURL: "https://zbmath.org/software/123", Assessment: "INDEXED"}}}}, nil
}
func TestResearchCLIWritesPartialSnapshotWithOrigin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "research.json")
	err := RunIdeaResearchCommand([]string{"discover", "--query", "symbolic algebra", "--output", path}, nil, []domain.MarketDiscoveryConnector{researchConnectorStub{id: domain.ConnectorZBMath}, researchConnectorStub{id: domain.ConnectorJOSS, failure: true}})
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var snapshot researchSnapshot
	if err := json.Unmarshal(body, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Attempts) != 2 || snapshot.Attempts[1].Error == "" || len(snapshot.Candidates) != 1 || snapshot.Candidates[0].Origins[0].ResponseHash != "fixture-hash" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
}
