/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/product_validation_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"reflect"
	"testing"
)

func productInput() domain.ProductValidationInput {
	input := domain.ProductValidationInput{SchemaVersion: "1", Domain: "software-search-and-evaluation", DiscoveryArms: []domain.DiscoveryArmCompletion{{Arm: "N", Planned: 12, Completed: 4}, {Arm: "F1", Planned: 12, Completed: 6}, {Arm: "F2", Planned: 12, Completed: 6}}}
	for i := 0; i < 10; i++ {
		input.CaseIDs = append(input.CaseIDs, fmt.Sprintf("case-%02d", i))
	}
	return input
}

func TestProductValidationFailsDiscoveryAndDoesNotEvaluateDownstream(t *testing.T) {
	report := ValidateSoftwareSearchProduct(productInput())
	if report.Verdict != "FAIL_DISCOVERY" || report.FailedPhase != "discovery" || report.Evidence != "NOT_EVALUATED" || report.Decision != "NOT_EVALUATED" || report.Utility != "NOT_EVALUATED" {
		t.Fatalf("short-circuit incorreto: %#v", report)
	}
}

func TestProductValidationDoesNotHideFailedArmWithAggregate(t *testing.T) {
	input := productInput()
	input.DiscoveryArms = []domain.DiscoveryArmCompletion{{Arm: "A", Planned: 10, Completed: 7}, {Arm: "B", Planned: 90, Completed: 90}}
	if got := ValidateSoftwareSearchProduct(input); got.Verdict != "FAIL_DISCOVERY" {
		t.Fatalf("média agregada escondeu braço: %#v", got)
	}
}

func TestProductValidationRejectsMalformedAndIsOrderInvariant(t *testing.T) {
	input := productInput()
	left := ValidateSoftwareSearchProduct(input)
	input.DiscoveryArms[0], input.DiscoveryArms[2] = input.DiscoveryArms[2], input.DiscoveryArms[0]
	if right := ValidateSoftwareSearchProduct(input); !reflect.DeepEqual(left, right) {
		t.Fatalf("ordem alterou saída: %#v %#v", left, right)
	}
	input = productInput()
	input.CaseIDs[1] = input.CaseIDs[0]
	if got := ValidateSoftwareSearchProduct(input); got.Verdict != "INVALID_BENCHMARK" {
		t.Fatalf("ID duplicado aceito: %#v", got)
	}
}
