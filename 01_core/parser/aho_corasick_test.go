/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/aho_corasick.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestAhoCorasick_MatchAll(t *testing.T) {
	patterns := []string{"gorm", "sqlx", "diesel", "prisma"}
	ac := NewAhoCorasick(patterns)

	matches := ac.MatchAll("Importing github.com/jinzhu/gorm and sqlx driver")
	if len(matches) != 2 {
		t.Fatalf("Esperava 2 correspondências (gorm, sqlx), obteve %d", len(matches))
	}

	foundMap := make(map[string]bool)
	for _, m := range matches {
		foundMap[m] = true
	}

	if !foundMap["gorm"] || !foundMap["sqlx"] {
		t.Errorf("Deveria ter encontrado gorm e sqlx")
	}
}
