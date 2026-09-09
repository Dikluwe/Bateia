/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_retrieval_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
)

func TestFunctional0045ReferenceValidation(t *testing.T) {
	caller := "package q\nfunc QueueTask() {\n alpha()\n}\n"
	helper := "package q\nfunc alpha() {}\n"
	callerHash := fmt.Sprintf("%x", sha256.Sum256([]byte(caller)))
	helperHash := fmt.Sprintf("%x", sha256.Sum256([]byte(helper)))
	cases := []struct {
		name, mode, hash   string
		line               int
		missing            bool
		wantHits, wantRefs int
	}{
		{"valid", "direct_call", callerHash, 3, false, 1, 1},
		{"stale", "direct_call", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 3, false, 0, 0},
		{"zero-line", "direct_call", callerHash, 0, false, 0, 0},
		{"outside-caller", "direct_call", callerHash, 1, false, 0, 0},
		{"beyond-file", "direct_call", callerHash, 99, false, 0, 0},
		{"missing-caller", "direct_call", callerHash, 3, true, 0, 0},
		{"lexical-survives", "lexical", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 3, false, 1, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, _ := json.Marshal(map[string]any{"hits": []any{map[string]any{"path": "q/helper.go", "symbol": "alpha", "kind": "function", "line": 2, "end_line": 2, "sha256": helperHash, "retrieval": c.mode, "related_to": []any{map[string]any{"path": "q/entry.go", "symbol": "QueueTask", "line": c.line, "sha256": c.hash}}}}})
			var report domain.ImplementationSearchReport
			if err := json.Unmarshal(raw, &report); err != nil {
				t.Fatal(err)
			}
			files := []domain.ExtractedFile{{Path: "q/helper.go", Content: helper}}
			if !c.missing {
				files = append(files, domain.ExtractedFile{Path: "q/entry.go", Content: caller})
			}
			result := usecase.ValidateImplementationHits(report, files)
			if len(result.Hits) != c.wantHits {
				t.Fatalf("hits=%d want%d", len(result.Hits), c.wantHits)
			}
			raw, _ = json.Marshal(result)
			var wire struct {
				Hits []struct {
					Related []any `json:"related_to"`
				}
			}
			json.Unmarshal(raw, &wire)
			if len(wire.Hits) > 0 && len(wire.Hits[0].Related) != c.wantRefs {
				t.Errorf("references=%d want%d", len(wire.Hits[0].Related), c.wantRefs)
			}
		})
	}
}
