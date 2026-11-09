/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0171_web_test.md
 * @layer L1
 * @updated 2026-11-09
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type integrity0171Probe struct {
	integrity0170Provider
	caps, reads int
}

func (p *integrity0171Probe) Capabilities(c context.Context) domain.WebSearchCapabilities {
	p.caps++
	return p.integrity0170Provider.Capabilities(c)
}
func (p *integrity0171Probe) PreviousSnapshot(string) (domain.WebSearchSnapshot, error) {
	p.reads++
	return domain.WebSearchSnapshot{}, fmt.Errorf("absent")
}
func integrity0171Input() WebExperimentInput {
	in := integrity0170Input()
	in.MaxPerQuery = 10
	return in
}

type integrity0171Runtime struct{ now time.Time }

func (r integrity0171Runtime) Now() time.Time { return r.now }
func (r integrity0171Runtime) WithTimeout(c context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	return context.WithCancel(c)
}
func integrity0171Ctx(now time.Time) context.Context {
	return domain.WithExecutionRuntime(context.Background(), integrity0171Runtime{now})
}

func TestIntegrity0171WebAdmissionPrecedesEveryPort(t *testing.T) {
	for _, auth := range []bool{false, true} {
		for _, limit := range []int{-10, -1, 0, 1, 5, 9, 10, 11, 100} {
			for _, margin := range []int{-1, 0, 1} {
				t.Run(fmt.Sprintf("auth%v/max%d/margin%d", auth, limit, margin), func(t *testing.T) {
					in := integrity0171Input()
					in.Authenticated = auth
					in.MaxPerQuery = limit
					in.SafeRemaining = margin
					valid := limit > 0
					if auth {
						valid = limit == 10 && margin >= 1
					}
					// Valid limit controls start fresh; invalid admission must precede even a failing resume reader.
					in.Resume = !valid
					before, _ := json.Marshal(in)
					p := &integrity0171Probe{}
					out, err := ExecuteWebExperiment(integrity0170Context(), in, p, p)
					if !valid {
						if err == nil || !reflect.DeepEqual(out, WebExperimentOutput{}) || p.caps != 0 || p.reads != 0 || len(p.requests) != 0 {
							t.Fatalf("invalid request touched ports: %v caps=%d reads=%d searches=%d receipt=%+v", err, p.caps, p.reads, len(p.requests), out.Receipt)
						}
					} else {
						if err != nil || len(p.requests) != 36 || out.Receipt.Completed != 36 {
							t.Fatalf("valid boundary failed: %v %+v", err, out.Receipt)
						}
						want := limit
						if want > 10 {
							want = 10
						}
						for _, r := range p.requests {
							if r.MaxResults != want {
								t.Errorf("max=%d want%d", r.MaxResults, want)
							}
						}
					}
					after, _ := json.Marshal(in)
					if string(before) != string(after) {
						t.Error("input mutated")
					}
				})
			}
		}
	}
	for _, kind := range []string{"nil", "not-authorized", "purpose"} {
		t.Run("authorization/"+kind, func(t *testing.T) {
			in := integrity0171Input()
			in.Resume = true
			switch kind {
			case "nil":
				in.Authorization = nil
			case "not-authorized":
				in.Authorization.Authorized = false
			case "purpose":
				in.Authorization.Purpose = "other"
			}
			p := &integrity0171Probe{}
			got, err := ExecuteWebExperiment(integrity0170Context(), in, p, p)
			if err == nil || !reflect.DeepEqual(got, WebExperimentOutput{}) || p.caps != 0 || p.reads != 0 || len(p.requests) != 0 {
				t.Fatal("invalid authorization touched ports", err, p, got)
			}
		})
	}
}

func integrity0171SnapshotBase(t *testing.T) WebExperimentOutput {
	t.Helper()
	out, err := ExecuteWebExperiment(integrity0170Context(), integrity0171Input(), &integrity0170Provider{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestIntegrity0171ResumeHonorsPersistedRateLimits(t *testing.T) {
	baseline := integrity0171SnapshotBase(t)
	now := integrity0170Runtime{}.Now()
	cases := []struct {
		name          string
		remaining     int64
		reset         *int64
		zero, blocked bool
	}{
		{name: "missing-reset", remaining: 100, blocked: true},
		{name: "negative-remaining", remaining: -1, blocked: true},
		{name: "future-reset", remaining: 100, reset: integrity0171Int(now.Unix() + 1), blocked: true},
		{name: "past-reset", remaining: 100, reset: integrity0171Int(now.Unix() - 1)},
		{name: "equal-reset", remaining: 100, reset: integrity0171Int(now.Unix())},
		{name: "zero-clock", remaining: 100, reset: integrity0171Int(now.Unix() - 1), zero: true, blocked: true},
		{name: "above-margin", remaining: 101},
	}
	for _, tc := range cases {
		for _, status := range []string{"SUCCESS", "EMPTY", "UNKNOWN_RATE_LIMIT"} {
			for _, arm := range []int{0, 2} {
				t.Run(fmt.Sprintf("%s/%s/arm%d", tc.name, status, arm), func(t *testing.T) {
					out := integrity0170Clone(baseline)
					for i := range out.Snapshots {
						s := &out.Snapshots[i]
						s.Responses = s.Responses[:2]
						s.Responses[1].Status = "UNKNOWN_PROVIDER_RESPONSE"
						s.Responses[1].Results = nil
						s.Responses[0].Receipt.RateLimit = domain.WebRateLimit{Remaining: integrity0171Int(500), Reset: integrity0171Int(now.Unix() - 1)}
					}
					// An expired low-margin response after a live restriction cannot clear it.
					out.Snapshots[2].Responses[0].Receipt.RateLimit = domain.WebRateLimit{Remaining: integrity0171Int(0), Reset: integrity0171Int(now.Unix() - 1)}
					// A later response/arm can constrain the first search; a permissive response cannot erase it.
					target := &out.Snapshots[arm].Responses[1]
					target.Status = status
					if status == "SUCCESS" {
						target.Results = integrity0170Clone(baseline.Snapshots[arm].Responses[1].Results)
					}
					target.Receipt.RateLimit = domain.WebRateLimit{Remaining: integrity0171Int(tc.remaining), Reset: tc.reset}
					for i := range out.Snapshots {
						integrity0170Reseal(&out.Snapshots[i])
					}
					reader := integrity0170ReaderFor(out)
					before, _ := json.Marshal(reader)
					in := integrity0171Input()
					in.Resume = true
					p := &integrity0170Provider{}
					ctx := integrity0171Ctx(now)
					if tc.zero {
						ctx = integrity0171Ctx(time.Time{})
					}
					got, err := ExecuteWebExperiment(ctx, in, p, reader)
					if err != nil {
						t.Fatal(err)
					}
					reused := 3
					if status == "SUCCESS" || status == "EMPTY" {
						reused++
					}
					if tc.blocked {
						if len(p.requests) != 0 || got.Receipt.Started != 0 || got.Receipt.Retried != 0 || got.Receipt.Completed != reused || got.Receipt.SkippedBudget != 36-reused || got.Receipt.StopReason != "SAFE_RATE_LIMIT_MARGIN" {
							t.Fatalf("persisted margin ignored: calls%d %+v", len(p.requests), got.Receipt)
						}
					} else if len(p.requests) != 36-reused || got.Receipt.Completed != 36 || got.Receipt.SkippedBudget != 0 {
						t.Fatalf("expired/absent restriction blocked: calls%d %+v", len(p.requests), got.Receipt)
					}
					for i, s := range got.Snapshots {
						replay, err := ReplayWebSnapshot(s)
						if err != nil || !replay.Replayable || len(replay.MissingFields) != 0 {
							t.Errorf("resume not replayable: %v %+v", err, replay)
						}
						for _, r := range s.Responses {
							old := out.Snapshots[i].Responses
							wasComplete := false
							for _, prev := range old {
								if prev.QueryID == r.QueryID && (prev.Status == "SUCCESS" || prev.Status == "EMPTY") {
									wasComplete = true
									if !reflect.DeepEqual(prev, r) {
										t.Error("completed response changed")
									}
								}
							}
							if tc.blocked && !wasComplete && r.Status != "SKIPPED_BUDGET" {
								t.Errorf("unknown promoted: %s", r.Status)
							}
						}
					}
					after, _ := json.Marshal(reader)
					if string(before) != string(after) {
						t.Error("reader snapshots mutated")
					}
				})
			}
		}
	}
}
func integrity0171Int(v int64) *int64 { return &v }

func TestIntegrity0171QuotaCannotMaskSnapshotAdmissionFailure(t *testing.T) {
	out := integrity0171SnapshotBase(t)
	out.Snapshots[0].Responses[0].Receipt.RateLimit = domain.WebRateLimit{Remaining: integrity0171Int(0)}
	integrity0170Reseal(&out.Snapshots[0])
	out.Snapshots[2].Responses[0].Results[0].Title = "tampered"
	in := integrity0171Input()
	in.Resume = true
	p := &integrity0170Provider{}
	got, err := ExecuteWebExperiment(integrity0170Context(), in, p, integrity0170ReaderFor(out))
	if err == nil || len(p.requests) != 0 || !reflect.DeepEqual(got, WebExperimentOutput{}) {
		t.Fatal("quota hid invalid snapshot", err, got)
	}
}
