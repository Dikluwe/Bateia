/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/clone_detector.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

// CalculateSimilarityScore compara dois fingerprints de 64-bit e calcula a taxa de similaridade (1.0 = idêntico, 0.0 = totalmente diferente).
func CalculateSimilarityScore(fp1, fp2 uint64) float64 {
	if fp1 == fp2 {
		return 1.0
	}

	// Calcula a distância de Hamming no nível de bits (XOR)
	diffBits := fp1 ^ fp2
	setBits := 0
	for diffBits > 0 {
		setBits += int(diffBits & 1)
		diffBits >>= 1
	}

	// Retorna percentual de similaridade de bit-level
	return 1.0 - (float64(setBits) / 64.0)
}
