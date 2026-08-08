/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/clone_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestCalculateSimilarityScore(t *testing.T) {
	fp1 := ComputeFingerprint("func Hello() { println(\"world\") }")
	fp2 := ComputeFingerprint("func Hello() { println(\"world\") }")

	score := CalculateSimilarityScore(fp1, fp2)
	if score != 1.0 {
		t.Errorf("Esperava taxa de similaridade 1.0 para fingerprints idênticos, obteve %f", score)
	}

	fp3 := ComputeFingerprint("func Goodbye() { os.Exit(0) }")
	scoreDiff := CalculateSimilarityScore(fp1, fp3)
	if scoreDiff >= 1.0 {
		t.Errorf("Esperava taxa de similaridade < 1.0 para funções diferentes, obteve %f", scoreDiff)
	}
}
