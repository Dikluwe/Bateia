/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r6_gate/discriminatory_step0026_r6_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package conformance

import (
	"context"
	"encoding/json"
)

var (
	ErrInvalid   = codedError("Invalid params")
	ErrConflict  = codedError("IDEMPOTENCY_CONFLICT")
	ErrQuota     = codedError("QUOTA_EXCEEDED")
	ErrNotFound  = codedError("RUN_NOT_FOUND")
	ErrNotReady  = codedError("ARTIFACT_NOT_READY")
	ErrIntegrity = codedError("INTEGRITY_ERROR")
	ErrCursor    = codedError("CURSOR_INVALID")
	ErrCAS       = codedError("VERSION_CONFLICT")
	ErrImmutable = codedError("ARTIFACT_IMMUTABLE")
)

type codedError string

func (e codedError) Error() string { return string(e) }

type Request struct {
	Problem             string
	Capabilities, Users []string
}
type Receipt struct {
	RunID, Lifecycle, Hash string
	Version                uint64
}
type Progress struct {
	Planned, Completed int
	Stage              string
}
type Artifact struct {
	Name, State, Hash string
	Version           uint64
}
type Run struct {
	ID, Lifecycle, Status string
	Version               uint64
	Progress              Progress
	Artifacts             map[string]Artifact
	TerminalEvents        []string
}
type Page struct {
	Entries []json.RawMessage
	Next    string
}
type Effects struct{ Store, Queue, IDs, Provider int }
type CrashObservation struct {
	Boundary    string
	Confirmed   bool
	Final, Temp bool
	Reopen      string
}

type Surface interface {
	Start(context.Context, string, Request) (Receipt, error)
	GetRun(context.Context, string) (Run, error)
	Transition(context.Context, string, uint64, string, Progress) (Run, error)
	Publish(context.Context, string, uint64, string, []byte) (Run, error)
	Read(context.Context, string, string, string, int) (Page, error)
	Cancel(context.Context, string) (Run, error)
	Complete(context.Context, string) (Run, error)
	Recover(context.Context) error
	Shutdown(context.Context) error
	Call(context.Context, string) ([]byte, error)
	ValidateRoot() error
	CrashCreate() ([]CrashObservation, error)
	CrashTransition() ([]CrashObservation, error)
	CrashPublish() ([]CrashObservation, error)
	Effects() Effects
}

type Factory func(root string) Surface
