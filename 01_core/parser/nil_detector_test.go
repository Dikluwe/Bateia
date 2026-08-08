/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/nil_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestEvaluateNilPointerRisk(t *testing.T) {
	code := `
package main

func process(req *Request) {
	line1 := "init"
	line2 := "init"
	line3 := "init"
	line4 := "init"
	line5 := "init"
	line6 := "init"
	req.Execute()
}
`

	violations := EvaluateNilPointerRisk("main.go", code)
	if len(violations) != 1 {
		t.Fatalf("Esperava 1 violação de NIL_POINTER_RISK, obteve %d", len(violations))
	}

	if violations[0].Rule != "NIL_POINTER_RISK" {
		t.Errorf("Esperava regra NIL_POINTER_RISK, obteve %s", violations[0].Rule)
	}
}
