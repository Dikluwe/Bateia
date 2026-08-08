/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/coupling.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

// CouplingMetrics armazena o acoplamento aferente (Ca), eferente (Ce) e o Índice de Instabilidade (I) de Robert C. Martin.
type CouplingMetrics struct {
	Afferent    int     `json:"afferent"`
	Efferent    int     `json:"efferent"`
	Instability float64 `json:"instability"`
}

// CalculateCoupling computa as métricas formais de acoplamento e instabilidade I = Ce / (Ca + Ce).
func CalculateCoupling(afferent, efferent int) CouplingMetrics {
	total := afferent + efferent
	instability := 0.0
	if total > 0 {
		instability = float64(efferent) / float64(total)
	}

	return CouplingMetrics{
		Afferent:    afferent,
		Efferent:    efferent,
		Instability: instability,
	}
}
