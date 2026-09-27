/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_sbom.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"encoding/json"
	"strings"
)

// This verifies only a literal declaration in the top-level component list.
// It is not a full CycloneDX schema validator, authenticity check or dependency
// graph resolver. Unsupported versions never become positive evidence.
func verifyTechnicalSBOMDeclaration(content, target string) (domain.TechnicalVerificationStatus, string) {
	if err := parser.ValidateBoundedJSON([]byte(content)); err != nil {
		return domain.TechnicalUnverified, "SBOM JSON inválido, ambíguo ou acima do orçamento"
	}
	var bom struct {
		Format      string `json:"bomFormat"`
		SpecVersion string `json:"specVersion"`
		Components  []struct {
			Ref     string `json:"bom-ref"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"components"`
	}
	if err := json.Unmarshal([]byte(content), &bom); err != nil || bom.Format != "CycloneDX" {
		return domain.TechnicalUnverified, "SBOM não é declaração JSON CycloneDX reconhecível"
	}
	if bom.SpecVersion != "1.6" {
		return domain.TechnicalUnsupported, "versão SBOM sem verificador: somente declarações 1.6 são suportadas"
	}
	if len(bom.Components) > 2000 {
		return domain.TechnicalUnverified, "SBOM excede limite de 2000 componentes"
	}
	refs := map[string]bool{}
	matched := false
	for _, component := range bom.Components {
		if len(component.Ref) > 8192 || len(component.Name) > 8192 || len(component.Version) > 8192 {
			return domain.TechnicalUnverified, "campo de componente SBOM excede orçamento"
		}
		if component.Ref != "" {
			if refs[component.Ref] {
				return domain.TechnicalUnverified, "SBOM contém bom-ref duplicado"
			}
			refs[component.Ref] = true
		}
		if strings.TrimSpace(component.Name) != "" && strings.TrimSpace(component.Version) != "" && component.Name+"@"+component.Version == target {
			matched = true
		}
	}
	if !matched {
		return domain.TechnicalUnverified, "componente versionado não declarado no SBOM coletado"
	}
	return domain.TechnicalVerified, "componente declarado no SBOM; presença no artefato executado não inferida"
}
