/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/fingerprint_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestComputeFingerprint(t *testing.T) {
	code1 := "package main\nfunc Hello() { println(\"World\") }"
	code2 := "package main\nfunc Hello() { println(\"World\") }"
	code3 := "package main\nfunc Goodbye() { println(\"Bye\") }"

	fp1 := ComputeFingerprint(code1)
	fp2 := ComputeFingerprint(code2)
	fp3 := ComputeFingerprint(code3)

	if !AreSimilar(fp1, fp2) {
		t.Errorf("Esperava fingerprints idênticos para código igual")
	}

	if AreSimilar(fp1, fp3) {
		t.Errorf("Esperava fingerprints diferentes para código distinto")
	}
}
