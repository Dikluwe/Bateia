/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/symbol_table_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestSymbolTable_Operations(t *testing.T) {
	st := NewSymbolTable()
	st.Insert("parser", domain.ASTSymbol{Name: "ParseCargoToml", Kind: "function", Line: 25})
	st.Insert("parser", domain.ASTSymbol{Name: "Manifest", Kind: "struct", Line: 10})

	funcs := st.LookupKind("function")
	if len(funcs) != 1 {
		t.Fatalf("Esperava 1 função, obteve %d", len(funcs))
	}

	if funcs[0].Name != "ParseCargoToml" {
		t.Errorf("Esperava ParseCargoToml, obteve %s", funcs[0].Name)
	}

	scopeSyms := st.LookupScope("parser")
	if len(scopeSyms) != 2 {
		t.Fatalf("Esperava 2 símbolos no escopo 'parser', obteve %d", len(scopeSyms))
	}
}
