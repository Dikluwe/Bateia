/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0117/model.md
 * @layer lab
 * @updated 2026-10-06
 */
package phaseexperiment

import "sync/atomic"

// Model is an authored abstraction; it is not an implementation of conc.
type Model struct {
	initialized atomic.Bool
	limit       atomic.Int64
}

func (m *Model) Configure(n int64, afterCheck func()) {
	if m.initialized.Load() {
		panic("already initialized")
	}
	if afterCheck != nil {
		afterCheck()
	}
	m.limit.Store(n)
}
func (m *Model) Initialize() { m.initialized.Store(true) }

// Snapshot is intended only after a rendezvous or join in this experiment.
func (m *Model) Snapshot() (bool, int64) { return m.initialized.Load(), m.limit.Load() }
