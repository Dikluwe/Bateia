/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/symbol_table.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"sync"

	"bateia/01_core/domain"
)

// SymbolTable armazena e organiza símbolos de AST por escopo e tipo em L1.
type SymbolTable struct {
	mu      sync.RWMutex
	symbols map[string][]domain.ASTSymbol
}

// NewSymbolTable constrói uma nova Tabela de Símbolos.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string][]domain.ASTSymbol),
	}
}

// Insert adiciona um símbolo à tabela sob uma chave de escopo (ex: nome do pacote ou arquivo).
func (st *SymbolTable) Insert(scope string, sym domain.ASTSymbol) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.symbols[scope] = append(st.symbols[scope], sym)
}

// LookupScope busca todos os símbolos pertencentes a um escopo específico.
func (st *SymbolTable) LookupScope(scope string) []domain.ASTSymbol {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.symbols[scope]
}

// LookupKind busca todos os símbolos de um tipo específico (ex: "function", "struct", "trait") em todos os escopos.
func (st *SymbolTable) LookupKind(kind string) []domain.ASTSymbol {
	st.mu.RLock()
	defer st.mu.RUnlock()

	var matched []domain.ASTSymbol
	for _, list := range st.symbols {
		for _, sym := range list {
			if sym.Kind == kind {
				matched = append(matched, sym)
			}
		}
	}
	return matched
}
