/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/cyclomatic.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestCalculateCyclomaticComplexity(t *testing.T) {
	simpleCode := `
func Hello() {
	println("Hello")
}
`
	if CalculateCyclomaticComplexity(simpleCode) != 1 {
		t.Errorf("Esperava complexidade 1 para função simples, obteve %d", CalculateCyclomaticComplexity(simpleCode))
	}

	complexCode := `
func Process(x int) {
	if x > 0 && x < 100 {
		for i := 0; i < x; i++ {
			switch i {
			case 1:
				println("1")
			case 2:
				println("2")
			}
		}
	}
}
`
	cc := CalculateCyclomaticComplexity(complexCode)
	if cc < 5 {
		t.Errorf("Esperava complexidade >= 5 para função com bifurcações, obteve %d", cc)
	}
}
