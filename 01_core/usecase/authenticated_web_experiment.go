/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/authenticated_web_experiment.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

var validWebStates = map[string]bool{
	"SUCCESS": true, "EMPTY": true, "UNKNOWN_RATE_LIMIT": true,
	"UNKNOWN_TIMEOUT": true, "UNKNOWN_PROVIDER_RESPONSE": true,
	"REJECTED_PRIVATE": true, "SKIPPED_BUDGET": true,
	"UNKNOWN": true, "FAILED": true,
}
