/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/export_dataset.md
 * @layer L2
 * @updated 2026-08-04
 */

package report

import (
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestExportJSONReport(t *testing.T) {
	rep := &domain.SearchReport{
		Query: domain.QueryFilter{Query: "test"},
	}

	data, err := ExportJSONReport(rep)
	if err != nil {
		t.Fatalf("Erro ao exportar JSON: %v", err)
	}
	if !strings.Contains(string(data), `"query": "test"`) {
		t.Errorf("Esperava query no JSON exportado")
	}
}

func TestExportSQLiteSchema(t *testing.T) {
	rep := &domain.SearchReport{
		Results: []domain.ExtractedRepository{
			{
				Candidate: domain.RepositoryCandidate{Owner: "owner", Repo: "repo"},
			},
		},
	}

	sqlOut := ExportSQLiteSchema(rep)
	if !strings.Contains(sqlOut, "CREATE TABLE IF NOT EXISTS repositories") {
		t.Errorf("Esperava DDL do SQLite no export")
	}
}
