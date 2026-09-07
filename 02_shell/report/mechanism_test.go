/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/mechanism_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package report

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"strings"
	"testing"
)

func TestUnknownAggregateAndMechanismFirst(t *testing.T) {
	repo := domain.ExtractedRepository{Ref: domain.CommitRef{Owner: "o", Repo: "r"}, Manifest: &domain.Manifest{Language: "cpp"}, QueryEvidence: &domain.QueryEvidence{Context: "x.cpp", Definitions: []domain.ImplementationEvidence{{Symbol: domain.SemanticSymbol{ID: "a", Name: "entry", File: "x.cpp", Line: 1, EndLine: 3}}}}}
	dist := usecase.CalculateLanguageDistribution([]domain.ExtractedRepository{repo}, nil)
	if dist[0].TopSimilarity != -1 {
		t.Fatal("missing score became zero")
	}
	md := GenerateMarkdownReport(&domain.SearchReport{Results: []domain.ExtractedRepository{repo}, LanguageDistribution: dist, SimilarityResults: []domain.IdeaSimilarityResult{{OverallScore: -1}}})
	if strings.Index(md, "Resposta à consulta") > strings.Index(md, "Resumo do Funil") {
		t.Fatal("answer buried")
	}
	for _, bad := range []string{"Corpo examinado integralmente", "**Commit-base local**: ``", "`0%`"} {
		if strings.Contains(md, bad) {
			t.Fatalf("misleading %q", bad)
		}
	}
	if !strings.Contains(md, "indisponível") || !strings.Contains(md, "Não avaliado") {
		t.Fatal("unknown hidden")
	}
}
