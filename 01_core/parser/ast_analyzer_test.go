/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ast_analyzer.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser_test

import (
	"testing"

	"bateia/01_core/parser"
)

func TestExtractGoSymbols(t *testing.T) {
	code := `package sample

// Service representa um serviço de domínio.
type Service interface {
	Execute() error
}

// UserConfig contém configurações.
type UserConfig struct {
	Name string
}

// NewService cria novo serviço.
func NewService() Service {
	return nil
}

func (u *UserConfig) Validate() bool {
	return true
}
`

	symbols, err := parser.ExtractGoSymbols("sample.go", code)
	if err != nil {
		t.Fatalf("Erro inesperado no ExtractGoSymbols: %v", err)
	}

	if len(symbols) != 4 {
		t.Fatalf("Esperado 4 símbolos extraídos, obteve %d", len(symbols))
	}

	expectedMap := map[string]string{
		"Service":    "interface",
		"UserConfig": "struct",
		"NewService": "function",
		"Validate":   "method",
	}

	for _, s := range symbols {
		expectedKind, found := expectedMap[s.Name]
		if !found {
			t.Errorf("Símbolo inesperado encontrado: %s", s.Name)
			continue
		}
		if s.Kind != expectedKind {
			t.Errorf("Para símbolo %s, esperado kind %s, obteve %s", s.Name, expectedKind, s.Kind)
		}
	}
}
