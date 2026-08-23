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

func TestEvaluateSecurityAdvisoriesAvoidsFalsePositives(t *testing.T) {
	manifest := &domain.Manifest{Dependencies: []domain.Dependency{
		{Name: "requests-mock", Version: "2.30.0"},
		{Name: "log4j-wrapper", Version: "2.14.0"},
		{Name: "requests", Version: "2.31.0"},
		{Name: "spring-core", Version: "5.3.18"},
		{Name: "pyyaml", Version: "unknown"},
	}}
	if advisories := EvaluateSecurityAdvisories(manifest); len(advisories) != 0 {
		t.Fatalf("falsos positivos encontrados: %+v", advisories)
	}
}

func TestEvaluateSecurityAdvisoriesIncludesProvenance(t *testing.T) {
	manifest := &domain.Manifest{Dependencies: []domain.Dependency{{Name: "requests", Version: "2.30.0"}}}
	advisories := EvaluateSecurityAdvisories(manifest)
	if len(advisories) != 1 {
		t.Fatalf("esperava um achado, recebeu %d", len(advisories))
	}
	if advisories[0].Source == "" || advisories[0].DatabaseUpdatedAt == "" || advisories[0].Confidence == "" {
		t.Fatalf("achado sem proveniência: %+v", advisories[0])
	}
}
