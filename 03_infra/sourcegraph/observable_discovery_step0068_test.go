/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/observable_discovery_step0068_test.md
 * @layer L3
 * @updated 2026-09-11
 */

package sourcegraph_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/sourcegraph"
)

func step0068TechnicalRequest(query, rule string, discriminators ...string) domain.DiscoveryRequest {
	return domain.DiscoveryRequest{
		ExpansionID: "exp-technical", Query: query, Language: "cpp", RepositoryScope: "control/scope.a+b",
		Discriminators: append([]string(nil), discriminators...), MaxResults: 5, MaxAttempts: 3,
		Derivation: &domain.ExpansionDerivation{
			Strategy: "technical-lexicon-v1", RuleIDs: []string{rule}, SourceIDs: []string{"exp-source"},
		},
	}
}

func step0068Repository(name, description string) map[string]interface{} {
	return map[string]interface{}{
		"__typename": "Repository", "name": "github.com/control/" + name,
		"url": "/github.com/control/" + name, "description": description,
	}
}

func TestStep0068TechnicalAdmissionRequiresBothDiscriminatorsInSameFile(t *testing.T) {
	items := []map[string]interface{}{
		step0068Repository("Metadata", "orient2d incircle"),
		discovery0055File("Split", "a.cpp", "orient2d"),
		discovery0055File("Split", "b.cpp", "incircle"),
		discovery0055File("Single", "a.cpp", "orient2d only"),
		discovery0055File("PathOnly", "orient2d-incircle.cpp", "unrelated text"),
		discovery0055File("SameFile", "predicates.cpp", "double ORIENT2D(...); double INCIRCLE(...);"),
	}
	var calls atomic.Int32
	var effective string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		var body struct {
			Variables map[string]string `json:"variables"`
		}
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		effective = body.Variables["query"]
		_, _ = writer.Write(mockSGResponse(items))
	}))
	defer server.Close()

	client := sourcegraph.NewClientWithEndpoint("", server.URL)
	result, err := client.Discover(t.Context(), step0068TechnicalRequest("orient2d incircle", "robust-predicates.symbol-pair", "orient2d", "incircle"))
	if err != nil {
		t.Fatal(err)
	}
	wantQuery := "repo:^" + regexp.QuoteMeta("control/scope.a+b") + "$ lang:C++ orient2d incircle count:50"
	if calls.Load() != 1 || effective != wantQuery {
		t.Fatalf("consulta técnica/fallback incorreto: calls=%d got=%q want=%q", calls.Load(), effective, wantQuery)
	}
	if strings.Contains(effective, "predicados geometricos") || strings.Contains(effective, "geometric predicates") || strings.Contains(effective, "patterntype:regexp") {
		t.Fatalf("consulta técnica sofreu redução genérica: %q", effective)
	}
	if len(result.Candidates) != 1 || result.Candidates[0].Owner != "control" || result.Candidates[0].Repo != "SameFile" {
		t.Fatalf("admissão same-file incorreta: %+v", result.Candidates)
	}
	if len(result.Attempts) != 1 {
		t.Fatalf("caminho técnico executou fallback: %+v", result.Attempts)
	}
	attempt := result.Attempts[0]
	if attempt.QueryPlan == nil || attempt.QueryPlan.Version != "0068-technical-v1" || attempt.QueryPlan.Reason != "technical_discriminators" || attempt.QueryPlan.Original != "orient2d incircle" {
		t.Fatalf("plano técnico observável incorreto: %+v", attempt.QueryPlan)
	}
	if !reflect.DeepEqual(attempt.QueryPlan.ActionTerms, []string{"orient2d"}) || !reflect.DeepEqual(attempt.QueryPlan.ObjectTerms, []string{"incircle"}) || !reflect.DeepEqual(attempt.QueryPlan.RuleIDs, []string{"robust-predicates.symbol-pair"}) || !reflect.DeepEqual(attempt.QueryPlan.Unverified, []string{"orient2d incircle"}) {
		t.Fatalf("grupos/proveniência do plano incorretos: %+v", attempt.QueryPlan)
	}
	if attempt.Coverage == nil || attempt.Coverage.Completeness != "unknown" || attempt.Coverage.ReceivedItems != len(items) || attempt.Coverage.ExaminedRepositories != 5 || attempt.Coverage.EligibleCandidates != 1 || attempt.Coverage.SelectedCandidates != 1 || attempt.RawCandidates != 5 || attempt.RejectedCandidates != 4 {
		t.Fatalf("cobertura/custo técnico incorreto: %+v", attempt)
	}
	if len(attempt.CandidateDiagnostics) != 5 {
		t.Fatalf("diagnósticos ausentes: %+v", attempt.CandidateDiagnostics)
	}
	wantReasons := map[string]string{
		"control/Metadata": "no_file_evidence", "control/Split": "split_across_files",
		"control/Single": "missing_object", "control/PathOnly": "missing_both", "control/SameFile": "matched_same_file",
	}
	for _, diagnostic := range attempt.CandidateDiagnostics {
		if diagnostic.Reason != wantReasons[diagnostic.Repository] {
			t.Errorf("diagnóstico %s=%+v", diagnostic.Repository, diagnostic)
		}
		if diagnostic.Repository == "control/SameFile" && (!diagnostic.ActionFound || !diagnostic.ObjectFound || diagnostic.Status != "selected") {
			t.Errorf("same-file não registrou ambos os grupos: %+v", diagnostic)
		}
	}
	if len(result.Matches) != 1 || result.Matches[0].Repository != "control/SameFile" || result.Matches[0].AttemptIndex != 0 || result.Matches[0].EffectiveQuery != effective || len(result.Matches[0].Evidence) == 0 {
		t.Fatalf("AdmissionEvidence técnico incompleto: %+v", result.Matches)
	}
}

func TestStep0068InvalidOrOpaqueTechnicalDerivationFailsClosedBeforeIO(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		_, _ = writer.Write(mockSGResponse(nil))
	}))
	defer server.Close()

	valid := step0068TechnicalRequest("orient2d incircle", "robust-predicates.symbol-pair", "orient2d", "incircle")
	cases := map[string]func(domain.DiscoveryRequest) domain.DiscoveryRequest{
		"opaque strategy": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Derivation.Strategy = "opaque"
			return request
		},
		"missing rule": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Derivation.RuleIDs = nil
			return request
		},
		"two rules": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Derivation.RuleIDs = []string{"a", "b"}
			return request
		},
		"missing source": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Derivation.SourceIDs = nil
			return request
		},
		"two sources": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Derivation.SourceIDs = []string{"a", "b"}
			return request
		},
		"one discriminator": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Discriminators = []string{"orient2d"}
			return request
		},
		"three discriminators": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Discriminators = []string{"a", "b", "c"}
			return request
		},
		"duplicate discriminators": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Discriminators = []string{"Orient2D", " orient2d "}
			return request
		},
		"empty discriminator": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Discriminators = []string{"orient2d", " "}
			return request
		},
		"long discriminator": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Discriminators = []string{strings.Repeat("x", 65), "incircle"}
			return request
		},
		"empty query": func(request domain.DiscoveryRequest) domain.DiscoveryRequest { request.Query = " "; return request },
		"long query": func(request domain.DiscoveryRequest) domain.DiscoveryRequest {
			request.Query = strings.Repeat("q", 257)
			return request
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			request := valid
			request.Derivation = &domain.ExpansionDerivation{
				Strategy: valid.Derivation.Strategy, RuleIDs: append([]string(nil), valid.Derivation.RuleIDs...), SourceIDs: append([]string(nil), valid.Derivation.SourceIDs...),
			}
			request.Discriminators = append([]string(nil), valid.Discriminators...)
			before := calls.Load()
			result, err := sourcegraph.NewClientWithEndpoint("", server.URL).Discover(t.Context(), mutate(request))
			if err == nil {
				t.Fatalf("entrada técnica inválida virou EMPTY/sucesso: %+v", result)
			}
			if calls.Load() != before || len(result.Attempts) != 0 || len(result.Candidates) != 0 {
				t.Fatalf("entrada inválida tocou I/O ou fabricou observação: calls=%d result=%+v", calls.Load()-before, result)
			}
		})
	}
}

func TestStep0068TechnicalEmptyAndErrorAreSingleAttemptWithoutFallback(t *testing.T) {
	for _, mode := range []string{"empty", "error"} {
		t.Run(mode, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				if mode == "error" {
					http.Error(writer, "unavailable", http.StatusServiceUnavailable)
					return
				}
				_, _ = writer.Write(mockSGResponse(nil))
			}))
			defer server.Close()
			result, err := sourcegraph.NewClientWithEndpoint("", server.URL).Discover(t.Context(), step0068TechnicalRequest("Shewchuk exact predicates", "robust-predicates.reference-pair", "Shewchuk", "predicates"))
			if calls.Load() != 1 || len(result.Attempts) != 1 || len(result.Candidates) != 0 {
				t.Fatalf("modo=%s calls=%d result=%+v", mode, calls.Load(), result)
			}
			if mode == "empty" {
				if err != nil || result.Attempts[0].Status != "EMPTY" || result.Attempts[0].Coverage == nil {
					t.Fatalf("resposta vazia válida incorreta: err=%v attempt=%+v", err, result.Attempts[0])
				}
			} else if err == nil || result.Attempts[0].Status != "FAILED" || result.Attempts[0].Error == "" {
				t.Fatalf("erro virou vazio: err=%v attempt=%+v", err, result.Attempts[0])
			}
		})
	}
}

func TestStep0068TechnicalProviderErrorPreservesIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		http.Error(writer, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	_, err := sourcegraph.NewClientWithEndpoint("", server.URL).Discover(t.Context(), step0068TechnicalRequest("orient2d incircle", "robust-predicates.symbol-pair", "orient2d", "incircle"))
	if err == nil || errors.Is(err, nil) {
		t.Fatal("erro técnico foi convertido em EMPTY")
	}
}
