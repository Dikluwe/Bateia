/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestEvaluateSecurityAdvisories(t *testing.T) {
	manifest := &domain.Manifest{
		ModuleName: "test-app",
		Dependencies: []domain.Dependency{
			{Name: "github.com/dgrijalva/jwt-go", Version: "v3.2.0"},
		},
	}

	advisories := EvaluateSecurityAdvisories(manifest)
	if len(advisories) != 1 {
		t.Fatalf("Esperava 1 alerta de segurança, obteve %d", len(advisories))
	}
	if advisories[0].Severity != "HIGH" {
		t.Errorf("Esperava severidade HIGH, obteve %s", advisories[0].Severity)
	}
}
