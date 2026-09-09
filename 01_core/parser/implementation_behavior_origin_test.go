/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/implementation_behavior_origin_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestBehavior57OriginDocumentationCannotPromoteBody(t *testing.T) {
	source := "package sample\n// Adjust records calibration instructions.\nfunc Adjust(counter *int) { *counter += 4 }\n"
	hit := domain.ImplementationHit{Path: "calibration.go", Symbol: "Adjust", Kind: "definition", Line: 3, EndLine: 3, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(source))), Retrieval: "documentation"}
	files := []domain.ExtractedFile{{Path: hit.Path, Content: source, Size: int64(len(source))}}
	documentary := parser.InspectImplementationBehavior(domain.ImplementationSearchReport{Hits: []domain.ImplementationHit{hit}}, files)
	if len(documentary) != 1 || documentary[0].Status != "not_evaluated" || len(documentary[0].Facts) != 0 || !strings.Contains(strings.ToLower(documentary[0].Reason), "doc") {
		t.Fatalf("documentary selection promoted to body evidence: %+v", documentary)
	}
	lexical := hit
	lexical.Retrieval = "lexical"
	for _, hits := range [][]domain.ImplementationHit{{hit, lexical}, {lexical, hit}} {
		got := parser.InspectImplementationBehavior(domain.ImplementationSearchReport{Hits: hits}, files)
		if len(got) != 2 {
			t.Fatalf("independent channels collapsed: %+v", got)
		}
		observed, unevaluated := 0, 0
		for _, b := range got {
			switch b.Status {
			case "observed":
				observed++
				if len(b.Facts) != 1 || b.Facts[0].Kind != "assignment" || b.Facts[0].Target != "*counter" || b.Facts[0].Snippet != "*counter += 4" {
					t.Errorf("lexical syntax missing: %+v", b)
				}
			case "not_evaluated":
				unevaluated++
				if len(b.Facts) != 0 || !strings.Contains(strings.ToLower(b.Reason), "doc") {
					t.Errorf("documentary channel lacks explicit limitation: %+v", b)
				}
			default:
				t.Errorf("unexpected status: %+v", b)
			}
		}
		if observed != 1 || unevaluated != 1 {
			t.Errorf("channels contaminated each other: %+v", got)
		}
	}
}
