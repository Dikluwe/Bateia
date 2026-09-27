/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_value_cost_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func TestDecisionValueCostUsesParetoDominanceWithoutHiddenExchangeRate(t *testing.T) {
	document := domain.InvestigationCost{Documents: 1}
	provider := domain.InvestigationCost{ProviderCalls: 1}
	if costDominates(document, provider) || costDominates(provider, document) {
		t.Fatal("custos incomparáveis ganharam taxa de conversão implícita")
	}
	if !costDominates(domain.InvestigationCost{Documents: 1}, domain.InvestigationCost{Documents: 2}) {
		t.Fatal("custo estritamente menor não dominou")
	}
}
