/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package variant

import (
	"bateia/04_wiring/r6_gate/conformance"
	reference "bateia/04_wiring/r6_gate/reference"
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type Subject struct{ *reference.Subject }

func New(root string) conformance.Surface { return &Subject{reference.NewSubject(root)} }
func (s *Subject) Start(ctx context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	d := sha256.Sum256([]byte(key))
	return s.StartWithID(ctx, key, req, "run_"+hex.EncodeToString(d[:16]))
}
