/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0170_web_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type integrity0170Runtime struct{}

func (integrity0170Runtime) Now() time.Time { return time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC) }
func (integrity0170Runtime) WithTimeout(c context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	return context.WithCancel(c)
}

type integrity0170Provider struct{ requests []domain.ObservableWebRequest }

func (*integrity0170Provider) Capabilities(context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{Provider: "github-repository-search", Authenticated: true, PublicOnly: true, PreservesRank: true, MaxResultsLimit: 10}
}
func (p *integrity0170Provider) SearchOne(_ context.Context, r domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	p.requests = append(p.requests, r)
	return domain.RawWebResponse{RequestID: r.RequestID, QueryID: r.QueryID, Provider: "github-repository-search", RetrievedAt: "2026-11-08T12:00:00Z", Status: "SUCCESS", HTTPStatus: 200, Results: []domain.RawWebResult{{Rank: 1, URL: "https://github.com/public/" + r.QueryID, Title: "Public tool " + r.QueryID, Snippet: "dependency graph evidence", Visibility: "PUBLIC"}}}, nil
}

type integrity0170Reader map[string]domain.WebSearchSnapshot

func (r integrity0170Reader) PreviousSnapshot(arm string) (domain.WebSearchSnapshot, error) {
	s, ok := r[arm]
	if !ok {
		return s, errors.New("absent")
	}
	return s, nil
}
func integrity0170Input() WebExperimentInput {
	in := WebExperimentInput{Authenticated: true, SafeRemaining: 100, MaxPerQuery: 5, Authorization: &domain.WebExperimentAuthorization{Provider: "github", Purpose: "public_repository_search", Authentication: "read_only_token", PlannedQueries: 36, Authorized: true}}
	for _, name := range []string{"N", "F1", "F2"} {
		a := WebExperimentArm{Name: name, Expansion: domain.FunctionalGrammarExpansion{SchemaVersion: "2", GrammarHash: "grammar-" + name}}
		for i := 12; i >= 1; i-- {
			a.Expansion.Queries = append(a.Expansion.Queries, domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: fmt.Sprintf("%s%d", name, i), Text: fmt.Sprintf("public graph %s %d", name, i)}})
		}
		in.Arms = append(in.Arms, a)
	}
	return in
}
func integrity0170Context() context.Context {
	return domain.WithExecutionRuntime(context.Background(), integrity0170Runtime{})
}
func integrity0170Execute(t *testing.T) WebExperimentOutput {
	t.Helper()
	out, err := ExecuteWebExperiment(integrity0170Context(), integrity0170Input(), &integrity0170Provider{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Snapshots) != 3 {
		t.Fatal("missing snapshots")
	}
	return out
}
func integrity0170Clone[T any](x T) T {
	b, _ := json.Marshal(x)
	var y T
	if err := json.Unmarshal(b, &y); err != nil {
		panic(err)
	}
	return y
}
func integrity0170Reseal(s *domain.WebSearchSnapshot) {
	for i := range s.Responses {
		s.Responses[i] = SealWebResponse(s.Responses[i].Request, s.Responses[i])
	}
	s.SnapshotHash = WebSnapshotHash(*s)
}
func integrity0170ReaderFor(out WebExperimentOutput) integrity0170Reader {
	r := integrity0170Reader{}
	for _, s := range out.Snapshots {
		r[s.Arm] = s
	}
	return r
}

func TestIntegrity0170ExecutorFeedsBlindUnion(t *testing.T) {
	out := integrity0170Execute(t)
	for _, schema := range []string{"2", "1"} {
		t.Run(schema, func(t *testing.T) {
			snapshots := integrity0170Clone(out.Snapshots)
			for i := range snapshots {
				if snapshots[i].SchemaVersion != "2" {
					t.Fatalf("authenticated producer schema=%s", snapshots[i].SchemaVersion)
				}
				snapshots[i].SchemaVersion = schema
				snapshots[i].SnapshotHash = WebSnapshotHash(snapshots[i])
				replay, err := ReplayWebSnapshot(snapshots[i])
				if err != nil || !replay.Replayable || len(replay.MissingFields) != 0 {
					t.Fatalf("producer not replayable: %v %+v", err, replay)
				}
			}
			before, _ := json.Marshal(snapshots)
			pkg, err := PrepareBlindCandidateUnion(snapshots)
			if err != nil {
				t.Fatal(err)
			}
			if pkg.SchemaVersion != "1" || len(pkg.Candidates) != 36 {
				t.Fatalf("public producer hits lost: %+v", pkg)
			}
			after, _ := json.Marshal(snapshots)
			if string(before) != string(after) {
				t.Fatal("input mutated")
			}
			for _, c := range pkg.Candidates {
				if len(c.Fragments) < 2 {
					t.Fatal("material public evidence lost")
				}
			}
		})
	}
	for _, kind := range []string{"unknown-schema", "missing-schema", "tamper", "private", "missing-material"} {
		t.Run(kind, func(t *testing.T) {
			s := integrity0170Clone(out.Snapshots)
			switch kind {
			case "unknown-schema":
				s[2].SchemaVersion = "3"
			case "missing-schema":
				s[2].SchemaVersion = ""
			case "tamper":
				s[2].Responses[0].Results[0].Title = "altered"
			case "private":
				s[2].Responses[0].Results[0].Visibility = "PRIVATE"
			case "missing-material":
				s[2].Responses[0].Provider = ""
			}
			if kind != "tamper" {
				integrity0170Reseal(&s[2])
			}
			got, err := PrepareBlindCandidateUnion(s)
			if !errors.Is(err, domain.ErrIntegrity) || !reflect.DeepEqual(got, domain.BlindCandidatePackage{}) {
				t.Fatalf("bad snapshot published: %v %+v", err, got)
			}
		})
	}
}

func TestIntegrity0170ResumeRejectsBeforeAnySearch(t *testing.T) {
	baseline := integrity0170Execute(t)
	mutations := map[string]func(*domain.WebSearchSnapshot){
		"schema":                func(s *domain.WebSearchSnapshot) { s.SchemaVersion = "3" },
		"missing-schema":        func(s *domain.WebSearchSnapshot) { s.SchemaVersion = "" },
		"run":                   func(s *domain.WebSearchSnapshot) { s.RunID = "other" },
		"arm":                   func(s *domain.WebSearchSnapshot) { s.Arm = "N" },
		"provider":              func(s *domain.WebSearchSnapshot) { s.Provider = "other" },
		"grammar":               func(s *domain.WebSearchSnapshot) { s.GrammarHash = "other" },
		"query-set":             func(s *domain.WebSearchSnapshot) { s.QuerySetHash = "other" },
		"authorization":         func(s *domain.WebSearchSnapshot) { s.Authorization.Purpose = "other" },
		"missing-authorization": func(s *domain.WebSearchSnapshot) { s.Authorization = nil },
		"budget":                func(s *domain.WebSearchSnapshot) { s.Budget.MaxResultsPerQuery++ },
		"missing-budget":        func(s *domain.WebSearchSnapshot) { s.Budget = nil },
		"query":                 func(s *domain.WebSearchSnapshot) { s.Responses[0].Request.Query = "foreign plan" },
		"request-id": func(s *domain.WebSearchSnapshot) {
			s.Responses[0].Request.RequestID = "foreign"
			s.Responses[0].RequestID = "foreign"
		},
		"query-id": func(s *domain.WebSearchSnapshot) {
			s.Responses[0].Request.QueryID = "foreign"
			s.Responses[0].QueryID = "foreign"
		},
		"max-results":       func(s *domain.WebSearchSnapshot) { s.Responses[0].Request.MaxResults++ },
		"locale":            func(s *domain.WebSearchSnapshot) { s.Responses[0].Request.Locale = "pt-BR" },
		"response-provider": func(s *domain.WebSearchSnapshot) { s.Responses[0].Provider = "other" },
		"response-arm":      func(s *domain.WebSearchSnapshot) { s.Responses[0].Arm = "N" },
		"position":          func(s *domain.WebSearchSnapshot) { s.Responses[0].Position++ },
		"attempt":           func(s *domain.WebSearchSnapshot) { s.Responses[0].AttemptNumber = 2 },
		"duplicate-query":   func(s *domain.WebSearchSnapshot) { s.Responses = append(s.Responses, s.Responses[0]) },
		"missing-material":  func(s *domain.WebSearchSnapshot) { s.Responses[0].RetrievedAt = "" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			out := integrity0170Clone(baseline)
			s := &out.Snapshots[2]
			mutate(s)
			integrity0170Reseal(s)
			r := integrity0170ReaderFor(out)
			first := r["N"]
			first.Responses = first.Responses[:1]
			first.SnapshotHash = WebSnapshotHash(first)
			r["N"] = first
			in := integrity0170Input()
			in.Resume = true
			p := &integrity0170Provider{}
			got, err := ExecuteWebExperiment(integrity0170Context(), in, p, r)
			if err == nil || !reflect.DeepEqual(got, WebExperimentOutput{}) || len(p.requests) != 0 {
				t.Fatalf("invalid resumed plan searched/published: err=%v requests=%d out=%+v", err, len(p.requests), got)
			}
			if err != nil {
				msg := strings.ToLower(err.Error())
				if !strings.Contains(msg, "integr") && !strings.Contains(msg, "reprod") {
					t.Errorf("error does not identify integrity/reproducibility: %v", err)
				}
			}
		})
	}
	t.Run("tampered-seal", func(t *testing.T) {
		out := integrity0170Clone(baseline)
		out.Snapshots[2].Responses[0].Results[0].Title = "tampered"
		r := integrity0170ReaderFor(out)
		first := r["N"]
		first.Responses = first.Responses[:1]
		first.SnapshotHash = WebSnapshotHash(first)
		r["N"] = first
		in := integrity0170Input()
		in.Resume = true
		p := &integrity0170Provider{}
		got, err := ExecuteWebExperiment(integrity0170Context(), in, p, r)
		if err == nil || len(p.requests) != 0 || !reflect.DeepEqual(got, WebExperimentOutput{}) {
			t.Fatalf("tampered snapshot repaired: %v %+v", err, got)
		}
	})
}

func TestIntegrity0170ResumeValidSubsetsAndLegacy(t *testing.T) {
	baseline := integrity0170Execute(t)
	for _, schema := range []string{"1", "2"} {
		for _, partial := range []bool{false, true} {
			t.Run(fmt.Sprintf("schema%s/partial%v", schema, partial), func(t *testing.T) {
				out := integrity0170Clone(baseline)
				for i := range out.Snapshots {
					s := &out.Snapshots[i]
					s.SchemaVersion = schema
					if schema == "1" {
						s.Authorization = nil
						s.Budget = nil
					}
					if partial {
						s.Responses = s.Responses[:2]
						s.Responses[0].Status = "EMPTY"
						s.Responses[0].Results = nil
						s.Responses[1].Status = "UNKNOWN_PROVIDER_RESPONSE"
						s.Responses[1].Results = nil
						integrity0170Reseal(s)
					} else {
						s.SnapshotHash = WebSnapshotHash(*s)
					}
				}
				in := integrity0170Input()
				in.Resume = true
				p := &integrity0170Provider{}
				got, err := ExecuteWebExperiment(integrity0170Context(), in, p, integrity0170ReaderFor(out))
				if err != nil {
					t.Fatal(err)
				}
				want := 0
				if partial {
					want = 33
				}
				if len(p.requests) != want || got.Receipt.Completed != 36 {
					t.Fatalf("resume state lost: calls=%d receipt=%+v", len(p.requests), got.Receipt)
				}
				if !partial && schema == "2" && !reflect.DeepEqual(got.Snapshots, baseline.Snapshots) {
					t.Error("full authenticated resume changed snapshots or seals")
				}
				if partial {
					for _, snapshot := range got.Snapshots {
						if snapshot.Responses[0].Status != "EMPTY" {
							t.Error("EMPTY was searched or reinterpreted")
						}
					}
					for i, r := range p.requests {
						arm := []string{"N", "F1", "F2"}[i%3]
						position := i/3 + 2
						if r.RequestID != fmt.Sprintf("0020-%02d-%s", position, arm) {
							t.Fatalf("partial sequence[%d]=%+v", i, r)
						}
					}
				}
			})
		}
	}
	for _, field := range []string{"authorization", "budget"} {
		t.Run("legacy-present-mismatch/"+field, func(t *testing.T) {
			out := integrity0170Clone(baseline)
			s := &out.Snapshots[2]
			s.SchemaVersion = "1"
			if field == "authorization" {
				s.Authorization.Purpose = "other"
			} else {
				s.Budget.MaxResultsPerQuery++
			}
			s.SnapshotHash = WebSnapshotHash(*s)
			in := integrity0170Input()
			in.Resume = true
			p := &integrity0170Provider{}
			got, err := ExecuteWebExperiment(integrity0170Context(), in, p, integrity0170ReaderFor(out))
			if err == nil || len(p.requests) != 0 || !reflect.DeepEqual(got, WebExperimentOutput{}) {
				t.Fatal("legacy mismatched metadata accepted", err)
			}
		})
	}
}

func TestIntegrity0170PlanAdmissionAndInputImmutability(t *testing.T) {
	for _, kind := range []string{"duplicate-arm", "unknown-arm", "duplicate-id", "empty-id", "empty-text", "valid"} {
		t.Run(kind, func(t *testing.T) {
			in := integrity0170Input()
			switch kind {
			case "duplicate-arm":
				in.Arms[2].Name = "F1"
			case "unknown-arm":
				in.Arms[2].Name = "OTHER"
			case "duplicate-id":
				in.Arms[2].Expansion.Queries[0].ID = in.Arms[2].Expansion.Queries[1].ID
			case "empty-id":
				in.Arms[2].Expansion.Queries[0].ID = ""
			case "empty-text":
				in.Arms[2].Expansion.Queries[0].Text = ""
			}
			before, _ := json.Marshal(in)
			p := &integrity0170Provider{}
			out, err := ExecuteWebExperiment(integrity0170Context(), in, p, nil)
			after, _ := json.Marshal(in)
			if string(before) != string(after) {
				t.Error("caller query slices sorted in place")
			}
			if kind == "valid" {
				if err != nil || len(p.requests) != 36 || out.Receipt.Completed != 36 {
					t.Fatal(err, out.Receipt)
				}
				return
			}
			if err == nil || len(p.requests) != 0 || !reflect.DeepEqual(out, WebExperimentOutput{}) {
				t.Fatalf("invalid plan admitted: %v %d", err, len(p.requests))
			}
		})
	}
}
