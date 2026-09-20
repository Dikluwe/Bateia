/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/mechanism_rules.md
 * @layer L1
 * @updated 2026-09-19
 */

package usecase

import _ "embed"

const activeMechanismRulesSHA256 = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"

//go:embed mechanism-rules.json
var embeddedMechanismRules []byte

func ActiveMechanismRuleCatalog() (catalog []byte, sha256Hex string) {
	return append([]byte(nil), embeddedMechanismRules...), activeMechanismRulesSHA256
}
