/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/ontology_matcher_test.md
 * @layer L1
 * @updated 2026-09-07
 */

package usecase_test

import (
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func TestExtractOntologyFromRepo(t *testing.T) {
	repo := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "gofiber", Repo: "fiber"},
		Candidate: domain.RepositoryCandidate{
			Owner:       "gofiber",
			Repo:        "fiber",
			Description: "Express inspired web framework written in Go",
		},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "github.com/valyala/fasthttp"},
				{Name: "github.com/stretchr/testify"},
			},
		},
		Files: []domain.ExtractedFile{
			{Path: "cmd/main.go"},
			{Path: "internal/router.go"},
			{Path: ".github/workflows/ci.yml"},
		},
	}

	ont := usecase.ExtractOntologyFromRepo(repo)

	if ont.Ref.Owner != "gofiber" || ont.Ref.Repo != "fiber" {
		t.Errorf("Esperado ref gofiber/fiber, recebido: %v", ont.Ref)
	}

	if len(ont.IntentKeywords) == 0 {
		t.Error("Esperado palavras-chave extraídas da descrição")
	}

	foundTesting := false
	for _, c := range ont.Capabilities {
		if c == "testing" {
			foundTesting = true
		}
	}
	if !foundTesting {
		t.Error("Esperado capacidade 'testing' detectada via testify")
	}
}

func TestCalculateIdeaSimilarity(t *testing.T) {
	target := usecase.ExtractOntologyFromIdea("Fast web framework with router and middleware")

	candidateRepo := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "gin-gonic", Repo: "gin"},
		Candidate: domain.RepositoryCandidate{
			Owner:       "gin-gonic",
			Repo:        "gin",
			Description: "Gin is a HTTP web framework written in Go",
		},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "github.com/stretchr/testify"},
			},
		},
	}
	candidate := usecase.ExtractOntologyFromRepo(candidateRepo)

	res := usecase.CalculateIdeaSimilarity(target, candidate)

	if res.OverallScore != -1 || res.SimilarityLabel != "Indeterminate (insufficient evidence)" {
		t.Errorf("Dimensões ausentes devem ser indeterminadas: %+v", res)
	}

	if res.Vectors.DomainIntentScore <= 0 {
		t.Errorf("Esperado DomainIntentScore > 0, obtido: %d", res.Vectors.DomainIntentScore)
	}
}

func TestMineOntologicalImprovements(t *testing.T) {
	target := usecase.ExtractOntologyFromIdea("Code search engine")

	repo1 := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "org", Repo: "repo1"},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "go.uber.org/zap"},
			},
		},
	}
	repo2 := domain.ExtractedRepository{
		Ref: domain.CommitRef{Owner: "org", Repo: "repo2"},
		Manifest: &domain.Manifest{
			Dependencies: []domain.Dependency{
				{Name: "go.uber.org/zap"},
			},
		},
	}

	repos := []domain.ExtractedRepository{repo1, repo2}
	results := []domain.IdeaSimilarityResult{
		{CandidateRef: repo1.Ref, OverallScore: 80},
		{CandidateRef: repo2.Ref, OverallScore: 75},
	}

	recs := usecase.MineOntologicalImprovements(target, results, repos)

	if len(recs) == 0 {
		t.Fatal("Esperado ao menos 1 recomendação para dependência 'zap' compartilhada")
	}

	if recs[0].AdoptionRate != 100 {
		t.Errorf("Esperado taxa de adoção 100%%, recebido: %d%%", recs[0].AdoptionRate)
	}
}
