/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/authenticated_web_experiment.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func runWebAuthenticatedExperiment(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	return runWebExperimentMode(args, stdout, provider, true)
}

func validateStep0020Authorization(value domain.WebExperimentAuthorization) error {
	return usecase.ValidateWebExperimentAuthorization(value)
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

func authenticatedExperimentVerdict(values []domain.WebSearchSnapshot) string {
	return usecase.AuthenticatedWebExperimentVerdict(values)
}
