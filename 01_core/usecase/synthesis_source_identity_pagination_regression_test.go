/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/synthesis_source_identity_pagination_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var synthesisAdmissionSentinel = errors.New("valid sources reached idempotency")

type synthesisAdmissionStore struct {
	domain.RecommendationRunStore
	reader   *synthesisSourceTestStore
	resolves int
}

func (s *synthesisAdmissionStore) Get(ctx context.Context, id domain.RunID) (domain.RecommendationRun, error) {
	return s.reader.Get(ctx, id)
}
func (s *synthesisAdmissionStore) ReadArtifact(ctx context.Context, r domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	return s.reader.ReadArtifact(ctx, r)
}
func (s *synthesisAdmissionStore) ResolveIdempotency(context.Context, string) (domain.IdempotencyRecord, bool, error) {
	s.resolves++
	return domain.IdempotencyRecord{}, false, synthesisAdmissionSentinel
}

type synthesisAdmissionQueue struct {
	domain.RecommendationJobQueue
	calls int
}

func (q *synthesisAdmissionQueue) Reserve(context.Context) (domain.RecommendationQueueReservation, error) {
	q.calls++
	return nil, errors.New("unexpected queue reservation")
}

type synthesisAdmissionIDs struct{}

func (synthesisAdmissionIDs) NewRunID() (domain.RunID, error) {
	return "", errors.New("unexpected ID generation")
}

func checkSynthesisAdmission(t *testing.T, reader *synthesisSourceTestStore, sources []domain.SourceArtifactRef, want error) {
	t.Helper()
	store, queue := &synthesisAdmissionStore{reader: reader}, &synthesisAdmissionQueue{}
	_, catalog := ActiveMechanismRuleCatalog()
	service := NewRecommendationRunService(store, queue, synthesisAdmissionIDs{}, nil).ConfigureImprovementSynthesis(catalog)
	request := domain.ImprovementSynthesisRequest{Idea: domain.IdeaSpec{Problem: "Bound concurrent work", Capabilities: []string{"queue"}}, SourceArtifacts: append([]domain.SourceArtifactRef(nil), sources...), Policy: domain.ImprovementSynthesisPolicy{MaxHypotheses: 1, MaxMechanismsPerHypothesis: 1, RequireCounterEvidenceSearch: true}}
	got, err := service.StartImprovementSynthesis(context.Background(), "source-admission-control", request)
	if !errors.Is(err, want) || !reflect.DeepEqual(got, domain.ImprovementSynthesisStartReceipt{}) {
		t.Errorf("unexpected source admission: %v %+v", err, got)
	}
	wantResolves := 0
	if errors.Is(want, synthesisAdmissionSentinel) {
		wantResolves = 1
	}
	if store.resolves != wantResolves || queue.calls != 0 {
		t.Fatal("effects escaped source admission", store.resolves, queue.calls)
	}
}

func TestSynthesisRejectsMismatchedRunIdentityBeforePages(t *testing.T) {
	for _, returnedID := range []domain.RunID{"", domain.RunID("run_" + strings.Repeat("f", 32))} {
		t.Run(string(returnedID), func(t *testing.T) {
			reader, sources := synthesisSourceFixture()
			sources = sources[:1]
			run := reader.runs[sources[0].RunID]
			run.RunID = returnedID
			reader.runs[sources[0].RunID] = run
			got, err := resolveSynthesisSources(context.Background(), reader, sources)
			if !errors.Is(err, domain.ErrIntegrity) || !reflect.DeepEqual(got, resolvedSynthesisSources{}) {
				t.Errorf("wrong run admitted: %v %+v", err, got)
			}
			if len(reader.reads) != 0 {
				t.Error("read pages before checking run identity")
			}
			reader.reads = nil
			checkSynthesisAdmission(t, reader, sources, domain.ErrIntegrity)
			if len(reader.reads) != 0 {
				t.Error("public admission read wrong run pages")
			}
		})
	}
	reader, sources := synthesisSourceFixture()
	checkSynthesisAdmission(t, reader, sources, synthesisAdmissionSentinel)
}

func TestSynthesisRejectsAnyRepeatedContinuationCursor(t *testing.T) {
	for _, tc := range []struct {
		name      string
		next      []string
		wantReads int
		valid     bool
	}{
		{"normal end", nil, 1, true},
		{"normal pages", []string{"A", "B"}, 3, true},
		{"adjacent", []string{"A", "A"}, 2, false},
		{"nonadjacent", []string{"A", "B", "A"}, 3, false},
		{"longer cycle", []string{"A", "B", "C", "A"}, 4, false},
		{"empty", []string{"A", ""}, 2, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader, sources := synthesisSourceFixture()
			sources = sources[1:2] // gaps allow independent entries without duplicate mechanism IDs
			calls := 0
			reader.page = func(r domain.ArtifactReadRequest) domain.ArtifactPage {
				page := domain.ArtifactPage{RunID: r.RunID, Artifact: r.Artifact, ArtifactHash: "sha256:" + sources[0].SHA256, Entries: [][]byte{[]byte(fmt.Sprintf(`{"code":"gap%d"}`, calls))}}
				if calls < len(tc.next) {
					next := tc.next[calls]
					page.NextCursor = &next
				}
				calls++
				return page
			}
			got, err := resolveSynthesisSources(context.Background(), reader, sources)
			if tc.valid {
				if err != nil || len(got.mechanismGaps) != tc.wantReads {
					t.Fatal("valid pages changed", err, got)
				}
			} else if !errors.Is(err, domain.ErrIntegrity) || !reflect.DeepEqual(got, resolvedSynthesisSources{}) {
				t.Errorf("repeated cursor admitted partial basis: %v %+v", err, got)
			}
			if calls != tc.wantReads {
				t.Errorf("revisited cursor: calls=%d want=%d", calls, tc.wantReads)
			}
			calls = 0
			reader.reads = nil
			want := domain.ErrIntegrity
			if tc.valid {
				want = synthesisAdmissionSentinel
			}
			checkSynthesisAdmission(t, reader, sources, want)
			if calls != tc.wantReads {
				t.Errorf("public admission revisited cursor: calls=%d want=%d", calls, tc.wantReads)
			}
		})
	}
}

func TestSynthesisCursorHistoryIsScopedToEachArtifact(t *testing.T) {
	reader, sources := synthesisSourceFixture()
	sources = sources[:2]
	calls := map[domain.RunArtifactName]int{}
	reader.page = func(r domain.ArtifactReadRequest) domain.ArtifactPage {
		calls[r.Artifact]++
		page := domain.ArtifactPage{RunID: r.RunID, Artifact: r.Artifact, ArtifactHash: "sha256:" + sources[0].SHA256}
		if calls[r.Artifact] == 1 {
			next := "shared-cursor"
			page.NextCursor = &next
		}
		return page
	}
	checkSynthesisAdmission(t, reader, sources, synthesisAdmissionSentinel)
	if calls[sources[0].Artifact] != 2 || calls[sources[1].Artifact] != 2 {
		t.Fatal("cursor history leaked across artifacts", calls)
	}
}
