/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/market_consolidation_order_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMarketConsolidationPreservesVariantsAndProvenanceAcrossOrders(t *testing.T) {
	a := domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: "catalog", ExternalID: "product", CanonicalURL: "https://example.com/product"}, Name: "Zulu", Summary: "Zulu description", Publisher: "Zulu publisher", ProductKind: "application", ListingLocales: []string{"pt-BR"}, ObservedMarkets: []string{"BR"}, PublisherCountries: []string{"BR"}, RepositoryLinks: []domain.ProductRepositoryLink{{OwnerRepo: "org/repo", Commit: "pin", Status: domain.RepositoryLinkUnverified, Source: "z"}}, Origins: []domain.MarketCandidateOrigin{{ConnectorID: "catalog", QueryID: "q", QueryLocale: "pt-BR", ResponseHash: "z"}}, Witnesses: []domain.IdentityWitness{{Kind: "first", Value: "z", ConnectorIDs: []string{"z", "catalog"}}}}
	b := a
	b.Name, b.Summary, b.Publisher = "Alpha", "Alpha description", "Alpha publisher"
	b.ListingLocales = []string{"en-US"}
	b.ObservedMarkets = []string{"US"}
	b.PublisherCountries = []string{"US"}
	b.RepositoryLinks = []domain.ProductRepositoryLink{{OwnerRepo: "org/repo", Commit: "pin", Status: domain.RepositoryLinkConflict, Source: "a"}}
	b.Origins = []domain.MarketCandidateOrigin{{ConnectorID: "catalog", QueryID: "q", QueryLocale: "en-US", ResponseHash: "a"}}
	b.Witnesses = []domain.IdentityWitness{{Kind: "second", Value: "a", ConnectorIDs: []string{"catalog", "a"}}}
	plan := []domain.MarketDiscoveryAttempt{{AttemptID: "attempt", QueryID: "q", ConnectorID: "catalog", ConnectorVersion: "fixture", QueryLocale: "en-US", Status: domain.AttemptPlanned, PlannedCost: 1}}
	var control domain.MarketDiscoveryArtifacts
	for i, values := range [][]domain.ProductCandidate{{a, b, a}, {a, a, b}, {b, a, a}} {
		before, _ := json.Marshal(values)
		merged, e := ConsolidateProductCandidates(values)
		if e != nil || len(merged) != 1 {
			t.Fatal(e, merged)
		}
		p := merged[0]
		if p.Name != "Alpha" || p.Summary != "Alpha description" || p.Publisher != "Alpha publisher" || p.Resolution != domain.IdentitySame {
			t.Errorf("incidental representative: %+v", p)
		}
		if !reflect.DeepEqual(p.ListingLocales, []string{"en-US", "pt-BR"}) || !reflect.DeepEqual(p.ObservedMarkets, []string{"BR", "US"}) || !reflect.DeepEqual(p.PublisherCountries, []string{"BR", "US"}) || len(p.Origins) != 2 || len(p.RepositoryLinks) != 2 || len(p.Witnesses) != 2 {
			t.Errorf("lost provenance: %+v", p)
		}
		got, e := ProjectMarketDiscoveryArtifacts(nil, plan, []domain.ConnectorDiscovery{{AttemptID: "attempt", Status: domain.AttemptSucceeded, Candidates: values}})
		if e != nil {
			t.Fatal(e)
		}
		if i == 0 {
			control = got
		} else if !reflect.DeepEqual(got, control) {
			t.Error("order changed market artifacts")
		}
		after, _ := json.Marshal(values)
		if string(before) != string(after) {
			t.Fatal("input mutated")
		}
	}
	b.Identity.ConnectorID = "different"
	separate, e := ConsolidateProductCandidates([]domain.ProductCandidate{a, b})
	if e != nil || len(separate) != 2 {
		t.Fatal("connector identities merged", e, separate)
	}
}
func TestMarketCanonicalControlKeepsFrozenHashes(t *testing.T) {
	p := domain.ProductCandidate{Identity: domain.ProductIdentity{ConnectorID: "catalog", ExternalID: "product", CanonicalURL: "https://example.com/product"}, Name: "Canonical", ProductKind: "application"}
	one, e := ConsolidateProductCandidates([]domain.ProductCandidate{p})
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(one)
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != "d77a9cf7e9106a9e1d481c30a7f0ef59ae7fb978a776fd73d9f5109de7ffd431" {
		t.Fatal("canonical unique market report changed")
	}
	same, e := ConsolidateProductCandidates([]domain.ProductCandidate{p, p})
	if e != nil || len(same) != 1 {
		t.Fatal(e, same)
	}
	raw, _ = json.Marshal(same)
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != "343ce75d0945ce04d03e336c735a533187ea358e4f62e7300504806ec6d46dc5" {
		t.Fatal("canonical duplicate market report changed")
	}
}
