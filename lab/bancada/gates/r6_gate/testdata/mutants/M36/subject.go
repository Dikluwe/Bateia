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
func (s *Subject) GetRun(ctx context.Context, id string) (conformance.Run, error) {
	run, e := s.Subject.GetRun(ctx, id)
	if e == conformance.ErrIntegrity {
		run, ok := s.RawRun(id)
		if ok {
			return run, nil
		}
	}
	return run, e
}
