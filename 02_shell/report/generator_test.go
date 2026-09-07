/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/generator_test.md
 * @layer L2
 * @updated 2026-09-07
 */

package report

import (
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestGenerateMarkdownReportIncludesOperationalDiagnostics(t *testing.T) {
	report := &domain.SearchReport{
		Query:    domain.QueryFilter{Query: "test", Language: "go"},
		Failures: []domain.ExecutionFailure{{Scope: "rust", Stage: "sweep", Message: "provider unavailable"}},
		Warnings: []domain.OperationalWarning{{Stage: "cache", Operation: "get", Target: "query", Message: "disk unavailable"}},
	}
	output := GenerateMarkdownReport(report)
	if !strings.Contains(output, "Falha parcial") || !strings.Contains(output, "Aviso recuperável") {
		t.Fatalf("diagnósticos ausentes no relatório: %s", output)
	}
}

func TestGenerateMarkdownReport(t *testing.T) {
	report := &domain.SearchReport{
		Query:          domain.QueryFilter{Query: "hexagonal", Language: "go"},
		GeneratedAt:    time.Now().UTC(),
		TopCandidates:  5,
		TriagedCount:   3,
		ExtractedCount: 2,
		Results: []domain.ExtractedRepository{
			{
				Ref:       domain.CommitRef{Owner: "foo", Repo: "bar", Commit: "c123"},
				Candidate: domain.RepositoryCandidate{Owner: "foo", Repo: "bar", Stars: 100},
				Manifest:  &domain.Manifest{ModuleName: "github.com/foo/bar", GoVersion: "1.22"},
				Files: []domain.ExtractedFile{
					{Path: "go.mod", Content: "module github.com/foo/bar", Size: 25},
				},
			},
		},
	}

	md := GenerateMarkdownReport(report)

	if !strings.Contains(md, "# Relatório de Análise Arquitetural — Bateia") {
		t.Errorf("Markdown não contém o título esperado")
	}

	if !strings.Contains(md, "foo/bar@c123") {
		t.Errorf("Markdown não contém o commit imutável foo/bar@c123")
	}

	if !strings.Contains(md, "github.com/foo/bar") {
		t.Errorf("Markdown não contém o nome do módulo")
	}
}

func TestEvidenceLanguageDoesNotPromoteSearch(t *testing.T) {
	empty := GenerateMarkdownReport(&domain.SearchReport{})
	for _, word := range []string{"Soluções Exatas Encontradas", "Descobrimento de Necessidades Ocultas", "Presença frequente"} {
		if strings.Contains(empty, word) {
			t.Fatalf("unsupported empty claim %s", word)
		}
	}
	r := domain.SearchReport{Results: []domain.ExtractedRepository{{Ref: domain.CommitRef{Repo: "USalign"}, Candidate: domain.RepositoryCandidate{Repo: "USalign"}, RelevanceScore: 85, HealthScore: 100}}, SimilarityResults: []domain.IdeaSimilarityResult{{CandidateRef: domain.CommitRef{Repo: "USalign"}, SimilarityLabel: "Divergent"}}, Telemetry: &domain.FunnelTelemetry{TotalDurationMs: 3}}
	md := GenerateMarkdownReport(&r)
	for _, word := range []string{"Não avaliado", "Divergência", "não é benchmark do projeto", "solução e comportamento não verificados", "fato observado", "inferência", "hipótese", "não verificado"} {
		if !strings.Contains(md, word) {
			t.Fatalf("missing %q", word)
		}
	}
	if strings.Contains(md, "`100/100`") {
		t.Fatal("health inflated")
	}
}

func TestReportPresentsCompilerEvidenceAsStatic0035(t *testing.T) {
	unit := domain.SemanticUnit{Path: "src/render.cpp", Status: "complete", Symbols: []domain.SemanticSymbol{{Name: "render", File: "src/render.cpp", Line: 10, EndLine: 15, Definition: true}}, Calls: []domain.SemanticCall{{Caller: "render", Callee: "blend", File: "src/render.cpp", Line: 12}}, Inputs: []domain.SemanticInput{{Path: "renderer.h", SHA256: "header-sha"}}}
	semantic := &domain.SemanticReport{Status: "partial", Version: "clang fixture", Units: []domain.SemanticUnit{unit}}
	md := GenerateMarkdownReport(&domain.SearchReport{Query: domain.QueryFilter{Query: "render"}, Results: []domain.ExtractedRepository{{Semantic: semantic}}})
	for _, want := range []string{"Análise estática C/C++", "partial", "render → blend", "src/render.cpp:10–15", "header-sha", "sem prova de comportamento"} {
		if !strings.Contains(md, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if strings.Contains(md, "Soluções Exatas Encontradas") {
		t.Fatal("static evidence promoted")
	}
}
