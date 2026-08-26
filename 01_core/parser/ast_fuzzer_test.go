/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ast_fuzzer_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func FuzzExtractASTSymbols(f *testing.F) {
	// Seed corpus com entradas de código Go e Rust válidas e corrompidas
	f.Add("main.go", []byte("package main\nfunc main() {}"))
	f.Add("lib.rs", []byte("pub fn run() -> Result<(), ()> { Ok(()) }"))
	f.Add("invalid.go", []byte("package ;;; ((( invalid syntax"))

	f.Fuzz(func(t *testing.T, filename string, content []byte) {
		// O parser de AST NUNCA deve provocar pânico, mesmo com código malformado
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("ExtractASTSymbols causou pânico com o arquivo %s: %v", filename, r)
			}
		}()

		_ = ExtractASTSymbols(filename, content)
	})
}
