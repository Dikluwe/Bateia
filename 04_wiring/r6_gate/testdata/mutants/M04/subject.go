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
	"sync"
)

type Subject struct {
	*reference.Subject
	mu    sync.Mutex
	first map[string]conformance.Receipt
}

func New(root string) conformance.Surface {
	return &Subject{Subject: reference.NewSubject(root), first: map[string]conformance.Receipt{}}
}
func (s *Subject) Start(ctx context.Context, key string, req conformance.Request) (conformance.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.first[key]; ok {
		return r, nil
	}
	r, e := s.Subject.Start(ctx, key, req)
	if e == nil {
		s.first[key] = r
	}
	return r, e
}
