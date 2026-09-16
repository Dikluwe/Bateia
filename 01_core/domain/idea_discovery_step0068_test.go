/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/idea_discovery_step0068_test.md
 * @layer L1
 * @updated 2026-09-11
 */

package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStep0068DiscoveryRequestDerivationJSONContract(t *testing.T) {
	derivation := &ExpansionDerivation{
		Strategy:  "technical-lexicon-v1",
		RuleIDs:   []string{"robust-predicates.symbol-pair"},
		SourceIDs: []string{"exp-capability-source"},
	}
	request := DiscoveryRequest{ExpansionID: "exp-technical", Query: "orient2d incircle", Derivation: derivation}

	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	for _, witness := range []string{`"derivation":`, `"strategy":"technical-lexicon-v1"`, `"rule_ids":["robust-predicates.symbol-pair"]`, `"source_ids":["exp-capability-source"]`} {
		if !strings.Contains(string(encoded), witness) {
			t.Fatalf("DiscoveryRequest perdeu derivação %q: %s", witness, encoded)
		}
	}

	var decoded DiscoveryRequest
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Derivation == nil || decoded.Derivation.Strategy != derivation.Strategy || len(decoded.Derivation.RuleIDs) != 1 || len(decoded.Derivation.SourceIDs) != 1 {
		t.Fatalf("round-trip perdeu derivação: %+v", decoded)
	}
}

func TestStep0068DiscoveryRequestKeepsLegacyJSONWhenDerivationIsAbsent(t *testing.T) {
	encoded, err := json.Marshal(DiscoveryRequest{ExpansionID: "legacy", Query: "code impact analysis", MaxResults: 5})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "derivation") {
		t.Fatalf("campo opcional alterou JSON legado: %s", encoded)
	}
}
