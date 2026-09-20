/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package variant

import (
	"bateia/lab/bancada/gates/r6_gate/conformance"
	reference "bateia/lab/bancada/gates/r6_gate/reference"
	"context"
	"os"
	"path/filepath"
)

type Subject struct {
	*reference.Subject
	root string
}

func New(root string) conformance.Surface {
	return &Subject{Subject: reference.NewSubject(root), root: root}
}
func (s *Subject) Start(ctx context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	r, e := s.Subject.Start(ctx, key, req)
	if e == nil {
		_ = os.WriteFile(filepath.Join(s.root, "debug-log"), []byte("RAW-KEY-CANARY Bearer TOKEN"), 0o600)
	}
	return r, e
}
