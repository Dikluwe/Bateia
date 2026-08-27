/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/discovery_ontology_test.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"encoding/json"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func loadOntologyFixture(t *testing.T) domain.DiscoveryOntology {
	t.Helper()
	concept := func(id domain.DiscoveryConceptID, trigger, query, discriminator string) domain.DiscoveryConcept {
		return domain.DiscoveryConcept{ID: id, Triggers: []string{trigger}, Queries: []domain.ConceptQuery{{Text: query, Discriminators: []string{discriminator}}}, Rationale: "teste"}
	}
	o := domain.DiscoveryOntology{SchemaVersion: "1", Version: "test", Concepts: []domain.DiscoveryConcept{
		concept("software-reuse-discovery", "não começar do zero", "software reuse discovery", "software reuse"),
		concept("repository-recommendation", "buscar repositórios", "repository recommendation", "repository recommendation"),
		concept("adoption-decision", "decisão de adoção", "dependency adoption decision", "adoption decision"),
		concept("software-composition", "composição de projetos", "software component composition", "component composition"),
		concept("internal-reuse", "reúso interno", "internal code reuse", "code reuse"),
	}}
	o.Concepts[0].Forbidden = []string{"katamari"}
	return o
}

func TestOntologyRecognizesIntentAndRejectsControls(t *testing.T) {
	o := loadOntologyFixture(t)
	for _, tc := range []struct {
		text string
		want domain.DiscoveryConceptID
	}{
		{"Não começar do zero", "software-reuse-discovery"}, {"buscar repositórios", "repository-recommendation"},
		{"decisão de adoção", "adoption-decision"}, {"composição de projetos", "software-composition"}, {"reúso interno", "internal-reuse"},
	} {
		p, _ := ExpandIdea(domain.IdeaSpec{Problem: tc.text})
		got, err := ExpandWithOntology(*p, o)
		if err != nil || len(got.Matches) != 1 || got.Matches[0].ConceptID != tc.want {
			t.Fatalf("%q: %+v %v", tc.text, got.Matches, err)
		}
	}
	for _, text := range []string{"jogo Katamari", "autofagia biológica", "composicionalidade"} {
		p, _ := ExpandIdea(domain.IdeaSpec{Problem: text})
		got, _ := ExpandWithOntology(*p, o)
		if len(got.Matches) != 0 {
			t.Fatalf("controle %q contaminado: %+v", text, got.Matches)
		}
	}
}

func TestOntologyHashAndOutputIgnoreOrdering(t *testing.T) {
	o := loadOntologyFixture(t)
	encoded, _ := json.Marshal(o)
	var reversed domain.DiscoveryOntology
	_ = json.Unmarshal(encoded, &reversed)
	for i, j := 0, len(reversed.Concepts)-1; i < j; i, j = i+1, j-1 {
		reversed.Concepts[i], reversed.Concepts[j] = reversed.Concepts[j], reversed.Concepts[i]
	}
	if OntologyHash(o) != OntologyHash(reversed) {
		t.Fatal("hash depende da ordem")
	}
	p, _ := ExpandIdea(domain.IdeaSpec{Problem: "não começar do zero", Capabilities: []string{"buscar repositórios"}})
	a, _ := ExpandWithOntology(*p, o)
	b, _ := ExpandWithOntology(*p, reversed)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("saída depende da ordem da ontologia")
	}
}

func TestOntologyAuditRejectsGenericAndCompetitorName(t *testing.T) {
	o := loadOntologyFixture(t)
	bad := []domain.SearchExpansion{
		{ID: "generic", Query: "decidir entre", Derivation: &domain.ExpansionDerivation{RuleIDs: []string{"adoption-decision"}}},
		{ID: "leak", Query: "katamari software reuse", Derivation: &domain.ExpansionDerivation{RuleIDs: []string{"software-reuse-discovery"}}},
	}
	r, err := AuditOntologyQueries(bad, o)
	if err != nil {
		t.Fatal(err)
	}
	if r.Rejected != 2 || len(r.Accepted) != 0 {
		t.Fatalf("auditoria permissiva: %+v", r)
	}
}
