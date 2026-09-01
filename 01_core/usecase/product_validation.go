/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/product_validation.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"sort"
	"strings"
)

func ValidateSoftwareSearchProduct(input domain.ProductValidationInput) domain.ProductValidationReport {
	report := domain.ProductValidationReport{SchemaVersion: "1", Domain: input.Domain, Evidence: "NOT_EVALUATED", Decision: "NOT_EVALUATED", Utility: "NOT_EVALUATED", Diagnostics: []string{}}
	invalid := func(message string) domain.ProductValidationReport {
		report.Verdict = "INVALID_BENCHMARK"
		report.Diagnostics = append(report.Diagnostics, message)
		sort.Strings(report.Diagnostics)
		return report
	}
	if input.SchemaVersion != "1" || input.Domain != "software-search-and-evaluation" {
		return invalid("schema ou domínio fora do contrato")
	}
	seenCases := map[string]bool{}
	for _, id := range input.CaseIDs {
		id = strings.TrimSpace(id)
		if id == "" || seenCases[id] {
			return invalid("case ID vazio ou duplicado")
		}
		seenCases[id] = true
	}
	if len(seenCases) != 10 {
		return invalid(fmt.Sprintf("corpus exige 10 casos; recebeu %d", len(seenCases)))
	}
	seenArms := map[string]bool{}
	for _, arm := range input.DiscoveryArms {
		if strings.TrimSpace(arm.Arm) == "" || seenArms[arm.Arm] || arm.Planned <= 0 || arm.Completed < 0 || arm.Completed > arm.Planned {
			return invalid("braço ausente, duplicado ou com contagem inválida")
		}
		seenArms[arm.Arm] = true
		rate := float64(arm.Completed) / float64(arm.Planned)
		report.Discovery = append(report.Discovery, domain.ProductArmMetric{Arm: arm.Arm, Planned: arm.Planned, Completed: arm.Completed, CompletionRate: rate, Passed: rate >= .8})
	}
	if len(report.Discovery) == 0 {
		return invalid("nenhum braço de descoberta")
	}
	sort.Slice(report.Discovery, func(i, j int) bool { return report.Discovery[i].Arm < report.Discovery[j].Arm })
	for _, arm := range report.Discovery {
		if !arm.Passed {
			report.Verdict = "FAIL_DISCOVERY"
			report.FailedPhase = "discovery"
			report.Diagnostics = append(report.Diagnostics, fmt.Sprintf("braço %s concluiu %.2f%%; mínimo 80%%", arm.Arm, arm.CompletionRate*100))
		}
	}
	if report.Verdict == "FAIL_DISCOVERY" {
		sort.Strings(report.Diagnostics)
		return report
	}
	report.Verdict = "INVALID_BENCHMARK"
	report.Diagnostics = append(report.Diagnostics, "entrada não contém métricas downstream seladas")
	return report
}
