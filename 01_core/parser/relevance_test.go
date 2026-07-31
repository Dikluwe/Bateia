/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/relevance.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestEvaluateRelevance_DirectImportMatch(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "github.com/example/myapp",
		Dependencies: []domain.Dependency{
			{Name: "github.com/gofiber/fiber/v2", Version: "v2.52.0", Indirect: false},
		},
	}
	cand := domain.RepositoryCandidate{Owner: "example", Repo: "myapp"}

	result := EvaluateRelevance("fiber", manifest, cand, nil)

	if result.Score != ScoreDirectImport {
		t.Errorf("esperava score %d, obteve %d", ScoreDirectImport, result.Score)
	}
	if result.Label != LabelDirectImport {
		t.Errorf("esperava label %s, obteve %s", LabelDirectImport, result.Label)
	}
}

func TestEvaluateRelevance_IndirectImportMatch(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "github.com/example/myapp",
		Dependencies: []domain.Dependency{
			{Name: "github.com/gofiber/fiber/v2", Version: "v2.52.0", Indirect: true},
		},
	}
	cand := domain.RepositoryCandidate{Owner: "example", Repo: "myapp"}

	result := EvaluateRelevance("fiber", manifest, cand, nil)

	if result.Score != ScoreIndirectImport {
		t.Errorf("esperava score %d, obteve %d", ScoreIndirectImport, result.Score)
	}
	if result.Label != LabelIndirectImport {
		t.Errorf("esperava label %s, obteve %s", LabelIndirectImport, result.Label)
	}
}

func TestEvaluateRelevance_ModuleNameMatch(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "github.com/gofiber/fiber",
	}
	cand := domain.RepositoryCandidate{Owner: "gofiber", Repo: "fiber"}

	result := EvaluateRelevance("fiber", manifest, cand, nil)

	if result.Score != ScoreModuleName {
		t.Errorf("esperava score %d, obteve %d", ScoreModuleName, result.Score)
	}
}

func TestEvaluateRelevance_StructuralMatch(t *testing.T) {
	tree := []domain.DirectoryNode{
		{Path: "cmd/server"},
		{Path: "fiber"},
	}
	cand := domain.RepositoryCandidate{Owner: "example", Repo: "structapp"}

	result := EvaluateRelevance("fiber", nil, cand, tree)

	if result.Score != ScoreStructural {
		t.Errorf("esperava score %d, obteve %d", ScoreStructural, result.Score)
	}
}

func TestEvaluateRelevance_GenericTextMatch(t *testing.T) {
	cand := domain.RepositoryCandidate{Owner: "moby", Repo: "moby"}

	result := EvaluateRelevance("fiber", nil, cand, nil)

	if result.Score != ScoreGenericText {
		t.Errorf("esperava score %d, obteve %d", ScoreGenericText, result.Score)
	}
}

func TestNormalizeQuery(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"  github.com/gofiber/fiber/v2 ", "github.com/gofiber/fiber"},
		{"Fiber", "fiber"},
		{"github.com/gin-gonic/gin/v3", "github.com/gin-gonic/gin"},
	}

	for _, c := range cases {
		got := NormalizeQuery(c.input)
		if got != c.expected {
			t.Errorf("NormalizeQuery(%q) = %q; esperava %q", c.input, got, c.expected)
		}
	}
}
