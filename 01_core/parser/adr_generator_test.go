/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/adr_generator_test.md
 * @layer L1
 * @updated 2026-09-26
 */

package parser

import (
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestGenerateADRDocument(t *testing.T) {
	rep := &domain.SearchReport{
		Query:          domain.QueryFilter{Query: "time tracking", Language: "go"},
		GeneratedAt:    time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC),
		TopCandidates:  12,
		TriagedCount:   4,
		ExtractedCount: 2,
		Results: []domain.ExtractedRepository{{
			Candidate:       domain.RepositoryCandidate{Owner: "solidtime", Repo: "solidtime", URL: "https://github.com/solidtime/solidtime"},
			Ref:             domain.CommitRef{Owner: "solidtime", Repo: "solidtime", Commit: "abc123"},
			RelevanceScore:  85,
			RelevanceReason: "Excellent match",
			Manifest: &domain.Manifest{Dependencies: []domain.Dependency{
				{Name: "direct-one", Version: "v1"},
				{Name: "indirect-one", Version: "v2", Indirect: true},
				{Name: "direct-two", Version: "v2"},
				{Name: "direct-three", Version: "v3"},
				{Name: "direct-four", Version: "v4"},
				{Name: "direct-five", Version: "v5"},
				{Name: "direct-six", Version: "v6"},
			}},
		}},
	}
	adr := GenerateADRDocument(rep)
	assertADRProposal(t, adr)
	for _, want := range []string{
		"# ADR-001", "2026-09-26", "time tracking", "Linguagem: go",
		"analisou 12 candidatos", "selecionou 4 na triagem", "extraiu 2 repositórios",
		"[solidtime/solidtime](https://github.com/solidtime/solidtime)",
		"Score de relevância da pesquisa: 85/100", "não é garantia de adequação",
		"Excellent match", "solidtime/solidtime@abc123", "Propõe-se avaliar",
		"Hipótese de benefício", "a pesquisa não comprova esse efeito",
		"compatibilidade funcional", "licença", "custo de integração", "riscos de segurança",
		"Dependências diretas observadas no candidato (não aprovadas)",
		"direct-one", "direct-two", "direct-three", "direct-four", "direct-five",
	} {
		if !strings.Contains(adr, want) {
			t.Errorf("ADR deve preservar %q; documento: %s", want, adr)
		}
	}
	for _, unwanted := range []string{"indirect-one", "direct-six", "Dependências de Referência Recomendadas"} {
		if strings.Contains(adr, unwanted) {
			t.Errorf("ADR não deve recomendar ou incluir %q", unwanted)
		}
	}
	if adr != GenerateADRDocument(rep) {
		t.Fatal("A proposta deve ser determinística para o mesmo relatório")
	}
}

func TestGenerateADRDocumentWithoutResults(t *testing.T) {
	adr := GenerateADRDocument(&domain.SearchReport{Query: domain.QueryFilter{Query: "unresolved need"}})
	assertADRProposal(t, adr)
	for _, want := range []string{"não contém candidatos extraídos", "não demonstra a inexistência de alternativas adequadas", "Próximo passo", "Revisar a pergunta"} {
		if !strings.Contains(adr, want) {
			t.Errorf("ADR vazio deve explicar %q", want)
		}
	}
	for _, unwanted := range []string{"Propõe-se avaliar", "Hipótese de benefício", "Score de relevância", "Dependências diretas observadas", "atendeu aos critérios estritos"} {
		if strings.Contains(adr, unwanted) {
			t.Errorf("Sem resultados, não inferir %q", unwanted)
		}
	}
}

func TestGenerateADRDocumentWithoutDirectDependencies(t *testing.T) {
	for _, manifest := range []*domain.Manifest{nil, {}, {Dependencies: []domain.Dependency{{Name: "only-indirect", Indirect: true}}}} {
		adr := GenerateADRDocument(&domain.SearchReport{Results: []domain.ExtractedRepository{{Manifest: manifest}}})
		assertADRProposal(t, adr)
		if strings.Contains(adr, "Dependências diretas observadas") || strings.Contains(adr, "only-indirect") {
			t.Errorf("Não criar seção vazia nem recomendar dependência indireta: %s", adr)
		}
	}
}

func assertADRProposal(t *testing.T, adr string) {
	t.Helper()
	if !strings.Contains(adr, "**Status**: PROPOSTO") {
		t.Errorf("Pesquisa deve gerar proposta: %s", adr)
	}
	for _, unsupported := range []string{"Aceito", "ACEITO", "testados em produção", "Redução de dívida técnica", "aprovou", "Recomenda-se a adoção"} {
		if strings.Contains(adr, unsupported) {
			t.Errorf("A pesquisa não sustenta %q", unsupported)
		}
	}
}

func TestGenerateADRDocumentUnpinnedReference(t *testing.T) {
	for _, ref := range []domain.CommitRef{{}, {Owner: "example", Repo: "candidate"}, {Owner: "example", Repo: "candidate", Commit: "HEAD"}} {
		adr := GenerateADRDocument(&domain.SearchReport{Results: []domain.ExtractedRepository{{Ref: ref}}})
		assertADRProposal(t, adr)
		if strings.Contains(adr, "**Rastreabilidade**: Commit ") && !strings.Contains(adr, "Commit não informado") {
			t.Errorf("Referência sem revisão não pode ser rotulada como commit: %s", adr)
		}
		if ref.Owner != "" && (!strings.Contains(adr, "example/candidate") || !strings.Contains(adr, "commit não fixado")) {
			t.Errorf("Preservar repositório e limitação do commit: %s", adr)
		}
		if ref.Owner == "" && !strings.Contains(adr, "Commit não informado") {
			t.Errorf("Explicitar ausência de referência: %s", adr)
		}
	}
}
