/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_replay_step0020_contract_test.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func step0020Snapshot(status string) domain.WebSearchSnapshot {
	req := domain.ObservableWebRequest{RequestID: "step0020-r1", QueryID: "N1", Query: "public repositories", MaxResults: 10}
	resp := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: "github-repository-search", RetrievedAt: "2026-08-31T12:00:00Z", Status: status}
	resp = SealWebResponse(req, resp)
	s := domain.WebSearchSnapshot{SchemaVersion: "1", RunID: "step0020", Arm: "N", Provider: "github-repository-search", Responses: []domain.RawWebResponse{resp}}
	s.SnapshotHash = WebSnapshotHash(s)
	return s
}

func TestStep0020ContractRejectsTamperedRequestHash(t *testing.T) {
	s := step0020Snapshot("EMPTY")
	s.Responses[0].RequestHash = strings.Repeat("a", 64)
	s.Responses[0].ResponseHash = WebResponseHash(s.Responses[0])
	s.SnapshotHash = WebSnapshotHash(s)
	r, err := ReplayWebSnapshot(s)
	if err == nil && r.Replayable {
		t.Fatal("request_hash adulterado sobreviveu após o atacante recalcular os hashes externos")
	}
}

func TestStep0020ContractEmptyIsNotUnknown(t *testing.T) {
	empty, err := ReplayWebSnapshot(step0020Snapshot("EMPTY"))
	if err != nil {
		t.Fatal(err)
	}
	unknown, err := ReplayWebSnapshot(step0020Snapshot("UNKNOWN_PROVIDER_RESPONSE"))
	if err != nil {
		t.Fatal(err)
	}
	eb, _ := json.Marshal(empty)
	ub, _ := json.Marshal(unknown)
	if !strings.Contains(string(eb), "EMPTY") || !strings.Contains(string(ub), "UNKNOWN_PROVIDER_RESPONSE") {
		t.Fatalf("replay apagou a causa observável: EMPTY=%s UNKNOWN=%s", eb, ub)
	}
}

func TestStep0020ContractRejectsInvalidNormalizedState(t *testing.T) {
	r, err := ReplayWebSnapshot(step0020Snapshot("NOT_A_STEP0020_STATE"))
	if err == nil && r.Replayable {
		t.Fatal("estado externo fora do vocabulário obrigatório foi aceito")
	}
}

func TestStep0020ContractReplayOfflineAndDeterministic(t *testing.T) {
	s := step0020Snapshot("EMPTY")
	a, err := ReplayWebSnapshot(s)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ReplayWebSnapshot(s)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("duas reexecuções puramente locais divergiram")
	}

	one := []domain.RawWebResponse{{QueryID: "N2"}, {QueryID: "N1"}}
	two := []domain.RawWebResponse{{QueryID: "N1"}, {QueryID: "N2"}}
	StableSnapshotResponses(one)
	StableSnapshotResponses(two)
	if !reflect.DeepEqual(one, two) {
		t.Fatal("normalização dependeu da ordem de chegada")
	}
}
