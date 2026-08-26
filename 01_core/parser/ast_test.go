/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ast_test.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"testing"
)

func TestExtractASTSymbols_Go(t *testing.T) {
	code := `package main

type User struct {
	ID string
}

type Processor interface {
	Process() error
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetID() string {
	return u.ID
}
`

	symbols := ExtractASTSymbols("main.go", []byte(code))
	if len(symbols) != 4 {
		t.Fatalf("Esperava 4 símbolos, obteve %d", len(symbols))
	}

	expected := map[string]string{
		"User":      "struct",
		"Processor": "interface",
		"NewUser":   "function",
		"GetID":     "method",
	}

	for _, s := range symbols {
		kind, ok := expected[s.Name]
		if !ok {
			t.Errorf("Símbolo inesperado: %s", s.Name)
		} else if s.Kind != kind {
			t.Errorf("Para símbolo %s esperava kind %s, obteve %s", s.Name, kind, s.Kind)
		}
	}
}
