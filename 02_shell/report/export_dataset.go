/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/export_dataset.md
 * @layer L2
 * @updated 2026-08-04
 */

package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

// ExportJSONReport serializa o relatório para o formato JSON formatado.
func ExportJSONReport(report *domain.SearchReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// ExportSQLiteSchema gera o script SQL DDL e DML para importar os dados no SQLite.
func ExportSQLiteSchema(report *domain.SearchReport) string {
	var sb strings.Builder

	sb.WriteString("-- Bateia Architectural Intelligence Dataset (SQLite Export)\n")
	sb.WriteString("CREATE TABLE IF NOT EXISTS repositories (\n")
	sb.WriteString("  id INTEGER PRIMARY KEY AUTOINCREMENT,\n")
	sb.WriteString("  owner TEXT NOT NULL,\n")
	sb.WriteString("  repo TEXT NOT NULL,\n")
	sb.WriteString("  commit_hash TEXT,\n")
	sb.WriteString("  score INTEGER,\n")
	sb.WriteString("  reason TEXT\n")
	sb.WriteString(");\n\n")

	sb.WriteString("CREATE TABLE IF NOT EXISTS symbols (\n")
	sb.WriteString("  id INTEGER PRIMARY KEY AUTOINCREMENT,\n")
	sb.WriteString("  repo_id INTEGER,\n")
	sb.WriteString("  name TEXT NOT NULL,\n")
	sb.WriteString("  kind TEXT NOT NULL,\n")
	sb.WriteString("  file TEXT,\n")
	sb.WriteString("  line INTEGER\n")
	sb.WriteString(");\n\n")

	for i, res := range report.Results {
		repoID := i + 1
		sb.WriteString(fmt.Sprintf("INSERT INTO repositories (id, owner, repo, commit_hash, score, reason) VALUES (%d, '%s', '%s', '%s', %d, '%s');\n",
			repoID,
			strings.ReplaceAll(res.Candidate.Owner, "'", "''"),
			strings.ReplaceAll(res.Candidate.Repo, "'", "''"),
			strings.ReplaceAll(res.Ref.Commit, "'", "''"),
			res.RelevanceScore,
			strings.ReplaceAll(res.RelevanceReason, "'", "''"),
		))

		for _, sym := range res.Symbols {
			sb.WriteString(fmt.Sprintf("INSERT INTO symbols (repo_id, name, kind, file, line) VALUES (%d, '%s', '%s', '%s', %d);\n",
				repoID,
				strings.ReplaceAll(sym.Name, "'", "''"),
				strings.ReplaceAll(sym.Kind, "'", "''"),
				strings.ReplaceAll(sym.File, "'", "''"),
				sym.Line,
			))
		}
	}

	return sb.String()
}
