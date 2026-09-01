/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/authenticated_web_experiment_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthenticatedWebExperimentDomainIsRedacted(t *testing.T) {
	value := WebExperimentAuthorization{Provider: "github", Purpose: "public_repository_search", Authentication: "read_only_token", PlannedQueries: 36, Authorized: true}
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(string(b)), "bearer ") || strings.Contains(strings.ToLower(string(b)), "authorization\"") {
		t.Fatalf("autorização redigida contém credencial: %s", b)
	}
}
