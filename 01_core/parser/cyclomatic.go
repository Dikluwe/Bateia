/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cyclomatic.md
 * @layer L1
 * @updated 2026-09-07
 */

package parser

import (
	"go/scanner"
	"go/token"
)

// GoLexicalComplexityMethod describes the scope and counting rule for this estimate.
const GoLexicalComplexityMethod = "estimativa lexical Go: 1 + tokens if/for/case/&&/||; sem CFG"

// EstimateGoLexicalComplexity counts selected Go tokens, excluding comments and
// literals. It neither constructs a control-flow graph nor validates full Go
// syntax. Results for other languages have no supported interpretation.
func EstimateGoLexicalComplexity(code string) int {
	fset := token.NewFileSet()
	file := fset.AddFile("", -1, len(code))
	var lexer scanner.Scanner
	// Malformed input must not cause diagnostic I/O in the pure parser layer.
	lexer.Init(file, []byte(code), func(token.Position, string) {}, 0)
	estimate := 1
	for {
		_, tok, _ := lexer.Scan()
		switch tok {
		case token.EOF:
			return estimate
		case token.IF, token.FOR, token.CASE, token.LAND, token.LOR:
			estimate++
		}
	}
}

// CalculateCyclomaticComplexity preserves the legacy API for a lexical Go estimate.
// Deprecated: use EstimateGoLexicalComplexity and GoLexicalComplexityMethod.
// This value is not McCabe complexity computed from a control-flow graph.
func CalculateCyclomaticComplexity(code string) int {
	return EstimateGoLexicalComplexity(code)
}
