/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/bounded_test.md
 * @layer L3
 * @updated 2026-10-09
 */
package clang

import (
	"io"
	"strings"
	"testing"
)

func TestCopyCannotBypassOutputBound(t *testing.T) {
	for _, size := range []int{1024, 1025} {
		b := &bounded{limit: 1024}
		_, err := io.Copy(b, struct{ io.Reader }{strings.NewReader(strings.Repeat("x", size))})
		if size == 1024 {
			if err != nil || b.Len() != size || b.exceeded {
				t.Fatalf("valid output: %d %v", b.Len(), err)
			}
		} else if err == nil || !b.exceeded || b.Len() > b.limit {
			t.Fatalf("copy bypassed limit: %d %v %v", b.Len(), b.exceeded, err)
		}
	}
}
