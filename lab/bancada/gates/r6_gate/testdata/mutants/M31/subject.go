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
)

type Subject struct{ *reference.Subject }

func New(root string) conformance.Surface { return &Subject{reference.NewSubject(root)} }
func (s *Subject) Read(ctx context.Context, id, name, cursor string, limit int) (conformance.Page, error) {
	page, err := s.Subject.Read(ctx, id, name, cursor, limit)
	if err == conformance.ErrNotReady {
		return s.ReadUnsigned(id, name, cursor, limit)
	}
	return page, err
}
