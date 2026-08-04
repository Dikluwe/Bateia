/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/security.md
 * @layer L1
 * @updated 2026-08-04
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
			{Name: "golang.org/x/net", Version: "v0.16.0"},
		},
	}

	advisories := EvaluateSecurityAdvisories(manifest)
	if len(advisories) != 2 {
		t.Fatalf("Esperava 2 alertas de segurança, obteve %d", len(advisories))
	}

	severities := make(map[string]bool)
	for _, a := range advisories {
		severities[a.Severity] = true
	}

	if !severities["HIGH"] || !severities["CRITICAL"] {
		t.Errorf("Esperava severidades HIGH e CRITICAL no teste de segurança")
	}
}
