/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence_step0073_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package domain_test

import (
	"encoding/json"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestStep0073EvidenceSourceExposesExactContentMatch(t *testing.T) {
	if domain.EvidenceMatchStrongContent != "STRONG_CONTENT" ||
		domain.EvidenceMatchSupportingContent != "SUPPORTING_CONTENT" ||
		domain.EvidenceMatchDocumentation != "DOCUMENTATION_CONTENT" ||
		domain.EvidenceMatchPathOnly != "PATH_ONLY" {
		t.Fatalf("classes públicas inesperadas: %q %q %q %q",
			domain.EvidenceMatchStrongContent,
			domain.EvidenceMatchSupportingContent,
			domain.EvidenceMatchDocumentation,
			domain.EvidenceMatchPathOnly)
	}

	source := domain.EvidenceSource{
		Kind: domain.EvidenceSymbol, Path: "bits/basic_string.h", Symbol: "_M_limit",
		Line: 407, Snippet: "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {",
		RuleID: "cpp-limit", MatchedPattern: "_M_limit", MatchClass: domain.EvidenceMatchStrongContent,
		Strength: .6,
	}
	encoded, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"line":407`, `"snippet":"constexpr size_type _M_limit`, `"matched_pattern":"_M_limit"`, `"match_class":"STRONG_CONTENT"`} {
		if !strings.Contains(string(encoded), field) {
			t.Fatalf("observável público ausente %s em %s", field, encoded)
		}
	}
	var roundTrip domain.EvidenceSource
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if roundTrip.Line != source.Line || roundTrip.Snippet != source.Snippet ||
		roundTrip.MatchedPattern != source.MatchedPattern || roundTrip.MatchClass != source.MatchClass {
		t.Fatalf("round-trip perdeu localização: want=%+v got=%+v", source, roundTrip)
	}
}
