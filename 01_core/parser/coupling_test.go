/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/coupling_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestCalculateCoupling(t *testing.T) {
	// Módulo Estável (Ca = 10, Ce = 0 -> I = 0.0)
	stable := CalculateCoupling(10, 0)
	if stable.Instability != 0.0 {
		t.Errorf("Esperava instabilidade 0.0 para módulo estável, obteve %f", stable.Instability)
	}

	// Módulo Instável (Ca = 0, Ce = 10 -> I = 1.0)
	unstable := CalculateCoupling(0, 10)
	if unstable.Instability != 1.0 {
		t.Errorf("Esperava instabilidade 1.0 para módulo instável, obteve %f", unstable.Instability)
	}

	// Módulo Balanceado (Ca = 5, Ce = 5 -> I = 0.5)
	balanced := CalculateCoupling(5, 5)
	if balanced.Instability != 0.5 {
		t.Errorf("Esperava instabilidade 0.5 para módulo balanceado, obteve %f", balanced.Instability)
	}
}
