/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/mechanism_extraction_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package usecase

import (
	"reflect"
	"testing"
)

func TestMechanismCanonicalSetIsStableAndRejectsControls(t *testing.T) {
	want := []string{"alpha", "beta"}
	if got := canonicalMechanismSet([]string{"beta", "alpha", "beta"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("canonical set=%v want=%v", got, want)
	}
	if got := normalizeMechanismText("unsafe\ntext"); got != "" {
		t.Fatalf("control character survived: %q", got)
	}
}
