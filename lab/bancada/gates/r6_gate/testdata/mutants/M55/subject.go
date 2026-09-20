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
	"net/http"
)

type Subject struct {
	*reference.Subject
	client *http.Client
}

func New(root string) conformance.Surface {
	return &Subject{Subject: reference.NewSubject(root), client: http.DefaultClient}
}
