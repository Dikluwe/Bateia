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
	"bateia/lab/bancada/gates/r6_gate/storefixture"
	"context"
)

type Subject struct {
	*reference.Subject
	store storefixture.Concrete
}

func New(root string) conformance.Surface {
	return &Subject{Subject: reference.NewSubject(root), store: storefixture.Concrete{Root: root}}
}
func (s *Subject) Call(ctx context.Context, tool string) ([]byte, error) {
	if tool == "legacy" {
		return []byte(`{"tools":["changed","search","benchmark"]}`), nil
	}
	return s.Subject.Call(ctx, tool)
}
