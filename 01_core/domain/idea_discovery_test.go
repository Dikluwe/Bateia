/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_discovery_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDiscoveryEvaluationSerializesUnknownPrecisionAsNull(t *testing.T) {
	data, err := json.Marshal(DiscoveryEvaluation{SchemaVersion: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"judged_precision":null`) {
		t.Fatalf("precisão desconhecida não é null: %s", data)
	}
}

func TestM21RegionalDiscoveryPublicJSONContract(t *testing.T) {
	regional := RegionalQueryProvenance{
		Locale:             "ja-JP",
		CatalogVersion:     "regional-recipes-v1",
		Concept:            "recipe-manager",
		RuleID:             "recipe-manager",
		SourceQuery:        "gerenciador de receitas",
		DerivationStrategy: "controlled-catalog-exact-v1",
		ProposalID:         "regional-ja-JP-recipe-manager",
	}
	expansion := SearchExpansion{ID: regional.ProposalID, Query: "レシピ管理", Regional: &regional}
	request := DiscoveryRequest{ExpansionID: expansion.ID, Query: expansion.Query, Language: "go", SearchLocale: regional.Locale, RegionalProposalID: regional.ProposalID}
	attempt := ProviderAttempt{Query: expansion.Query, SearchLocale: regional.Locale, RegionalProposalID: regional.ProposalID}
	receipt := ProviderAttemptReceipt{ExpansionID: expansion.ID, PlannedQuery: expansion.Query, EffectiveQuery: expansion.Query, SearchLocale: regional.Locale, RegionalProposalID: regional.ProposalID}
	origin := CandidateOrigin{ExpansionID: expansion.ID, PlannedQuery: expansion.Query, EffectiveQuery: expansion.Query, SearchLocale: regional.Locale, RegionalProposalID: regional.ProposalID}
	alternative := SoftwareAlternative{Repository: CommitRef{Owner: "org", Repo: "recipes", Commit: "abc123"}, DiscoveryLocales: []string{"ja-JP"}}

	for name, value := range map[string]any{"request": request, "attempt": attempt, "receipt": receipt, "origin": origin} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		for _, token := range []string{`"search_locale":"ja-JP"`, `"regional_proposal_id":"regional-ja-JP-recipe-manager"`} {
			if !strings.Contains(string(encoded), token) {
				t.Fatalf("%s lost %s: %s", name, token, encoded)
			}
		}
	}
	expansionJSON, err := json.Marshal(expansion)
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{`"regional":{`, `"locale":"ja-JP"`, `"proposal_id":"regional-ja-JP-recipe-manager"`} {
		if !strings.Contains(string(expansionJSON), token) {
			t.Fatalf("expansion lost regional provenance %s: %s", token, expansionJSON)
		}
	}
	alternativeJSON, err := json.Marshal(alternative)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(alternativeJSON), `"discovery_locales":["ja-JP"]`) {
		t.Fatalf("alternative lost discovery locales: %s", alternativeJSON)
	}
}

func TestM21RegionalOriginRequiresAuditableSourceQuery(t *testing.T) {
	invalid := CandidateOrigin{ExpansionID: "regional-ja-JP", SearchLocale: "ja-JP", RegionalProposalID: "regional-ja-JP-recipe-manager", EffectiveQuery: "レシピ管理"}
	if err := ValidateRegionalCandidateOrigin(invalid); err == nil {
		t.Fatal("regional origin without planned/source query was accepted")
	}
	valid := invalid
	valid.PlannedQuery = "レシピ管理"
	if err := ValidateRegionalCandidateOrigin(valid); err != nil {
		t.Fatalf("auditable regional origin rejected: %v", err)
	}
}

func TestM21ZeroRegionalValuesPreserveLegacyJSON(t *testing.T) {
	request, err := json.Marshal(DiscoveryRequest{ExpansionID: "exp", Query: "query", Language: "go", MaxResults: 3})
	if err != nil {
		t.Fatal(err)
	}
	wantRequest := `{"expansion_id":"exp","query":"query","language":"go","max_results":3}`
	if string(request) != wantRequest {
		t.Fatalf("zero-locale request changed legacy JSON:\ngot  %s\nwant %s", request, wantRequest)
	}
	snapshot, err := json.Marshal(DiscoverySnapshot{SchemaVersion: "2", Idea: IdeaSpec{}, Strategy: "literal", PerExpansion: []ExpansionDiscovery{}, Candidates: []DiscoveredCandidate{}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(snapshot), "search_locale") || strings.Contains(string(snapshot), "regional") {
		t.Fatalf("zero-locale snapshot leaked M21 fields: %s", snapshot)
	}
}
