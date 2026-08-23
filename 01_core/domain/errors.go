/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/types.md
 * @layer L1
 * @updated 2026-08-23
 */

package domain

import "errors"

// ErrManifestNotFound identifica repositórios sem um manifesto suportado.
// Adaptadores devem envolver este erro com contexto usando fmt.Errorf e %w.
var ErrManifestNotFound = errors.New("manifesto suportado não encontrado")
