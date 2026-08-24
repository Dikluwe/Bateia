/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/idea_discovery.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDiscoveryEvaluationSerializesUnknownPrecisionAsNull(t *testing.T) {
	data, err := json.Marshal(DiscoveryEvaluation{SchemaVersion: "1"})
	if err != nil { t.Fatal(err) }
	if !strings.Contains(string(data), `"judged_precision":null`) { t.Fatalf("precisão desconhecida não é null: %s", data) }
}
