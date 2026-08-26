/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/generator_test.md
 * @layer L2
 * @updated 2026-07-25
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
