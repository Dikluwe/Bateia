/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/polyglot_sweep.md
 * @layer L1
 * @updated 2026-08-08
 */

package usecase_test

import (
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func TestParseLanguageList(t *testing.T) {
	langs := usecase.ParseLanguageList("go,typescript,python,rust")
	if len(langs) != 4 {
		t.Fatalf("Esperado 4 linguagens, recebido: %d", len(langs))
	}
	if langs[0] != "go" || langs[1] != "typescript" || langs[2] != "python" || langs[3] != "rust" {
		t.Errorf("Linguagens parseadas incorretamente: %v", langs)
	}

	sweepLangs := usecase.ParseLanguageList("sweep")
	if len(sweepLangs) < 4 {
		t.Errorf("Esperado ao menos 4 linguagens padrão no modo sweep, recebido: %v", sweepLangs)
	}
}

func TestFilterLanguages(t *testing.T) {
	initial := []string{"go", "typescript", "python", "rust", "java", "cpp"}
	filtered := usecase.FilterLanguages(initial, "go,typescript")

	for _, l := range filtered {
		if l == "go" || l == "typescript" {
			t.Errorf("Linguagem %s deveria ter sido excluida", l)
		}
	}
	if len(filtered) != 4 {
		t.Errorf("Esperado 4 linguagens remanescentes, recebido: %d", len(filtered))
	}
}

func TestCalculateLanguageDistribution(t *testing.T) {
	repoGo := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "gofiber", Repo: "fiber"},
		Manifest: &domain.Manifest{
			Language: "go",
			Dependencies: []domain.Dependency{
				{Name: "fasthttp"},
			},
		},
		Candidate: domain.RepositoryCandidate{Stars: 1000},
	}

	repoRust := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "tokio-rs", Repo: "axum"},
		Manifest: &domain.Manifest{
			Language: "rust",
			Dependencies: []domain.Dependency{
				{Name: "tokio"},
			},
		},
		Candidate: domain.RepositoryCandidate{Stars: 2000},
	}

	results := []domain.ExtractedRepository{repoGo, repoRust}
	simResults := []domain.IdeaSimilarityResult{
		{CandidateRef: repoGo.Ref, OverallScore: 85},
		{CandidateRef: repoRust.Ref, OverallScore: 90},
	}

	dist := usecase.CalculateLanguageDistribution(results, simResults)

	if len(dist) != 2 {
		t.Fatalf("Esperado 2 ecossistemas na distribuição, recebido: %d", len(dist))
	}

	foundRust := false
	for _, d := range dist {
		if d.Count != 1 {
			t.Errorf("Esperado 1 repositório por ecossistema, recebido: %d para %s", d.Count, d.Language)
		}
		if d.TopSimilarity <= 0 {
			t.Errorf("Esperado TopSimilarity > 0 para %s", d.Language)
		}
		if d.Language == "Rust 🦀" {
			foundRust = true
		}
	}

	if !foundRust {
		t.Error("Esperado encontrar 'Rust 🦀' formatado na distribuição")
	}
}
