/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_experiment_run_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
	"testing"
)

type experimentProbe struct{ requests []string }

func (*experimentProbe) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "memory"}
}
func (p *experimentProbe) SearchOne(_ context.Context, r domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	p.requests = append(p.requests, r.RequestID)
	return domain.RawWebResponse{Status: "EMPTY", Provider: "memory"}, nil
}
func TestWebExperimentRunsInterleavedWithoutStorage(t *testing.T) {
	in := WebExperimentInput{MaxPerQuery: 5}
	for _, name := range []string{"N", "F1", "F2"} {
		arm := WebExperimentArm{Name: name}
		for i := 12; i >= 1; i-- {
			arm.Expansion.Queries = append(arm.Expansion.Queries, domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: fmt.Sprintf("q-%d", i), Text: "graph"}})
		}
		in.Arms = append(in.Arms, arm)
	}
	p := &experimentProbe{}
	out, err := ExecuteWebExperiment(context.Background(), in, p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.requests) != 36 || p.requests[0] != "0019-01-N" || p.requests[1] != "0019-01-F1" || p.requests[2] != "0019-01-F2" || out.Receipt.Completed != 36 {
		t.Fatalf("schedule: %v %+v", p.requests, out.Receipt)
	}
	if len(out.Snapshots) != 3 || out.Snapshots[0].SnapshotHash == "" {
		t.Fatal("unsealed snapshots")
	}
}
func TestWebExperimentVerdictRemainsBounded(t *testing.T) {
	snaps := make([]domain.WebSearchSnapshot, 3)
	for i := range snaps {
		for n := 0; n < 12; n++ {
			snaps[i].Responses = append(snaps[i].Responses, domain.RawWebResponse{Status: "EMPTY"})
		}
	}
	if AuthenticatedWebExperimentVerdict(snaps) != "READY_FOR_BLINDED_JUDGMENT" {
		t.Fatal("complete experiment rejected")
	}
	for n := 0; n < 4; n++ {
		snaps[0].Responses[n].Status = "UNKNOWN_RATE_LIMIT"
	}
	if AuthenticatedWebExperimentVerdict(snaps) != "INCONCLUSIVE_RATE_LIMIT" {
		t.Fatal("rate limitation promoted")
	}
	snaps[0].Responses[0].Status = "UNKNOWN_PROVIDER_RESPONSE"
	for n := 1; n < 4; n++ {
		snaps[0].Responses[n].Status = "UNKNOWN_PROVIDER_RESPONSE"
	}
	if AuthenticatedWebExperimentVerdict(snaps) != "UNKNOWN_PROVIDER_EXPOSURE" {
		t.Fatal("unknown exposure promoted")
	}
	if ValidateWebExperimentAuthorization(domain.WebExperimentAuthorization{}) == nil {
		t.Fatal("opaque authorization accepted")
	}
}

func TestWebExperimentVerdictThresholdsAndAdmission(t *testing.T) {
	for _, test := range []struct {
		completed [3]int
		want      string
	}{{[3]int{10, 11, 10}, "READY_FOR_BLINDED_JUDGMENT"}, {[3]int{10, 12, 10}, "UNKNOWN_PROVIDER_EXPOSURE"}, {[3]int{9, 9, 9}, "UNKNOWN_PROVIDER_EXPOSURE"}} {
		snapshots := make([]domain.WebSearchSnapshot, 3)
		for i, count := range test.completed {
			for n := 0; n < 12; n++ {
				status := "UNKNOWN_PROVIDER_RESPONSE"
				if n < count {
					status = "EMPTY"
				}
				snapshots[i].Responses = append(snapshots[i].Responses, domain.RawWebResponse{Status: status})
			}
		}
		if got := AuthenticatedWebExperimentVerdict(snapshots); got != test.want {
			t.Fatalf("threshold %v: %s", test.completed, got)
		}
	}
	valid := WebExperimentAdmission{Authenticated: true, AuthenticatedFlag: true, HasAuthorization: true, SafeRemaining: 10, MaxPerQuery: 10}
	if err := ValidateWebExperimentAdmission(valid); err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"token", "limit", "margin", "authorization"} {
		bad := valid
		switch kind {
		case "token":
			bad.NoToken = true
		case "limit":
			bad.MaxPerQuery = 11
		case "margin":
			bad.SafeRemaining = 0
		case "authorization":
			bad.HasAuthorization = false
		}
		if ValidateWebExperimentAdmission(bad) == nil {
			t.Fatalf("invalid %s admitted", kind)
		}
	}
}
