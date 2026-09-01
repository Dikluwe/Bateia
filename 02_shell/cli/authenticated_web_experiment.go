/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/authenticated_web_experiment.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bateia/01_core/domain"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func runWebAuthenticatedExperiment(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	return runWebExperimentMode(args, stdout, provider, true)
}

func validateStep0020Authorization(a domain.WebExperimentAuthorization) error {
	if a.Provider != "github" || a.Purpose != "public_repository_search" || a.Authentication != "read_only_token" || a.PrivateRepositoriesAllowed || a.PlannedQueries != 36 || !a.Authorized {
		return fmt.Errorf("autorização do Passo 0020 inválida")
	}
	return nil
}

func naturalQueryPosition(id string) int {
	i := len(id)
	for i > 0 && id[i-1] >= '0' && id[i-1] <= '9' {
		i--
	}
	if i == len(id) {
		return 0
	}
	n, _ := strconv.Atoi(id[i:])
	return n
}

func rejectSensitiveWebArtifact(value any) error {
	b, _ := json.Marshal(value)
	text := strings.ToLower(string(b))
	for _, forbidden := range []string{"authorization=", "authorization: bearer", "bearer ", "set-cookie", "cookie=", "tekt_step0020_canary", "step0020-token-canary"} {
		if strings.Contains(text, forbidden) {
			return fmt.Errorf("artefato rejeitado por política de segurança")
		}
	}
	return nil
}

func authenticatedExperimentVerdict(snapshots []domain.WebSearchSnapshot) string {
	minRate, maxRate, rateLimited := 1.0, 0.0, false
	for _, snapshot := range snapshots {
		completed := 0
		for _, response := range snapshot.Responses {
			if response.Status == "SUCCESS" || response.Status == "EMPTY" {
				completed++
			}
			if response.Status == "UNKNOWN_RATE_LIMIT" {
				rateLimited = true
			}
			if snapshot.Budget != nil && response.Receipt.RateLimit.Remaining != nil && *response.Receipt.RateLimit.Remaining <= int64(snapshot.Budget.SafeRemaining) {
				rateLimited = true
			}
		}
		rate := float64(completed) / 12
		if rate < minRate {
			minRate = rate
		}
		if rate > maxRate {
			maxRate = rate
		}
	}
	if minRate >= .8 && maxRate-minRate <= .1 {
		return "READY_FOR_BLINDED_JUDGMENT"
	}
	if rateLimited {
		return "INCONCLUSIVE_RATE_LIMIT"
	}
	return "UNKNOWN_PROVIDER_EXPOSURE"
}
