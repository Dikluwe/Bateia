/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/fingerprint.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "hash/fnv"

// ComputeFingerprint calcula um hash de 64-bit FNV-1a em O(1) para comparação rápida de ASTs/código.
func ComputeFingerprint(content string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(content))
	return h.Sum64()
}

// AreSimilar verifica se dois fingerprints de 64-bit correspondem ao mesmo hash sintático.
func AreSimilar(fp1, fp2 uint64) bool {
	return fp1 == fp2
}
