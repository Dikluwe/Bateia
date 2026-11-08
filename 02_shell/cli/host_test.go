/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/host_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"errors"
	"testing"
)

func TestUnconfiguredHostRefusesFileEffects(t *testing.T) {
	h := unavailableCommandHost{}
	if _, e := h.ReadFile("must-not-open"); !errors.Is(e, errCommandHost) {
		t.Fatal(e)
	}
	if e := h.WriteFile("must-not-create", nil, 0600); !errors.Is(e, errCommandHost) {
		t.Fatal(e)
	}
	if _, e := h.OpenRoot("must-not-open"); !errors.Is(e, errCommandHost) {
		t.Fatal(e)
	}
	if !h.Now().IsZero() {
		t.Fatal("missing host invented a date")
	}
}
