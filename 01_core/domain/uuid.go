/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/uuid.md
 * @layer L1
 * @updated 2026-08-08
 */

package domain

import (
	"crypto/sha1"
	"fmt"
)

// GenerateRepoUUID gera um UUID v5 determinístico e idempotente em L1 a partir da chave do repositório e commit SHA.
func GenerateRepoUUID(repoKey, commitSHA string) string {
	input := fmt.Sprintf("bateia:repo:%s@%s", repoKey, commitSHA)
	hash := sha1.Sum([]byte(input))

	// Formata como UUID v5 (variante RFC 4122 e versão 5)
	hash[6] = (hash[6] & 0x0f) | 0x50
	hash[8] = (hash[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
}
