/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/improvement_sources_test.md
 * @layer L1
 * @updated 2026-10-28
 */

package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type synthesisSourceTestStore struct {
	runs    map[domain.RunID]domain.RecommendationRun
	entries map[domain.RunArtifactName][][]byte
	getErr  error
	readErr error
	page    func(domain.ArtifactReadRequest) domain.ArtifactPage
	reads   []domain.ArtifactReadRequest
}

func (s *synthesisSourceTestStore) Get(_ context.Context, id domain.RunID) (domain.RecommendationRun, error) {
	if s.getErr != nil {
		return domain.RecommendationRun{}, s.getErr
	}
	return s.runs[id], nil
}
func (s *synthesisSourceTestStore) ReadArtifact(_ context.Context, request domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	s.reads = append(s.reads, request)
	if s.readErr != nil {
		return domain.ArtifactPage{}, s.readErr
	}
	if s.page != nil {
		return s.page(request), nil
	}
	return domain.ArtifactPage{RunID: request.RunID, Artifact: request.Artifact, ArtifactHash: "sha256:" + strings.Repeat("a", 64), Entries: s.entries[request.Artifact]}, nil
}

func synthesisSourceFixture() (*synthesisSourceTestStore, []domain.SourceArtifactRef) {
	recommendation := domain.RunID("run_" + strings.Repeat("1", 32))
	market := domain.RunID("run_" + strings.Repeat("2", 32))
	hash := strings.Repeat("a", 64)
	mechanismID := "mech_" + strings.Repeat("3", 64)
	encode := func(value any) []byte {
		b, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		return b
	}
	store := &synthesisSourceTestStore{runs: map[domain.RunID]domain.RecommendationRun{}, entries: map[domain.RunArtifactName][][]byte{
		domain.ArtifactMechanismCatalog: {encode(domain.ArchitecturalMechanism{ID: mechanismID, Name: "Bounded queue", Purpose: "Bound work", Strategy: "Backpressure", Status: domain.MechanismProven})},
		domain.ArtifactMechanismGaps:    {encode(domain.MechanismGap{Code: "TARGET_UNKNOWN", RuleID: "queue", WitnessHash: "witness"})},
		domain.ArtifactMarketCandidates: {[]byte(`{"candidate_id":"candidate-1","identity":{"connector_id":"catalog","external_id":"project-1"}}`)},
		domain.ArtifactMarketSnapshot:   {[]byte(`{"observation_id":"observation-1","query_id":"query-1","attempt_id":"attempt-1"}`)},
	}}
	sources := []domain.SourceArtifactRef{}
	for _, record := range []struct {
		id    domain.RunID
		kind  domain.RunKind
		names []domain.RunArtifactName
	}{
		{recommendation, domain.RunKindRecommendation, []domain.RunArtifactName{domain.ArtifactMechanismCatalog, domain.ArtifactMechanismGaps}},
		{market, domain.RunKindMarketDiscovery, []domain.RunArtifactName{domain.ArtifactMarketCandidates, domain.ArtifactMarketSnapshot}},
	} {
		run := domain.RecommendationRun{RunID: record.id, RunKind: record.kind, Lifecycle: domain.RunSucceeded}
		for _, name := range record.names {
			run.Artifacts = append(run.Artifacts, domain.RunArtifactIndex{Name: name, State: domain.ArtifactPublished, Hash: "sha256:" + hash})
			sources = append(sources, domain.SourceArtifactRef{RunID: record.id, RunKind: record.kind, Artifact: name, SHA256: hash})
		}
		store.runs[record.id] = run
	}
	return store, sources
}

// The same fixtures and frozen digests exercised the original service methods.
func resolveSourcesForTest(ctx context.Context, store *synthesisSourceTestStore, sources []domain.SourceArtifactRef) (resolvedSynthesisSources, error) {
	return resolveSynthesisSources(ctx, store, sources)
}

func TestSynthesisSourcesCharacterization(t *testing.T) {
	for _, degraded := range []bool{false, true} {
		t.Run(fmt.Sprint(degraded), func(t *testing.T) {
			store, sources := synthesisSourceFixture()
			if degraded {
				run := store.runs[sources[0].RunID]
				run.Lifecycle = domain.RunFailed
				store.runs[run.RunID] = run
			}
			resolved, err := resolveSourcesForTest(context.Background(), store, sources)
			if err != nil {
				t.Fatal(err)
			}
			if resolved.degraded != degraded || len(resolved.mechanisms) != 1 || len(resolved.mechanismGaps) != 1 || len(resolved.observations) != 7 {
				t.Fatalf("unexpected basis: %+v", resolved)
			}
			for _, read := range store.reads {
				if read.Limit != 50 || read.Cursor != "" {
					t.Fatalf("pagination changed: %+v", read)
				}
			}
			request := domain.ImprovementSynthesisRequest{Idea: domain.IdeaSpec{Problem: "Bound concurrent work", Capabilities: []string{"queue"}}, SourceArtifacts: sources, Policy: domain.ImprovementSynthesisPolicy{MaxHypotheses: 2, MaxMechanismsPerHypothesis: 1, RequireCounterEvidenceSearch: true}}
			analysis := synthesizeImprovement(request, improvementAnalysisBasis{mechanisms: resolved.mechanisms, degraded: resolved.degraded})
			projection := struct {
				Mechanisms     map[string]domain.ArchitecturalMechanism
				Gaps           map[string]domain.MechanismGap
				Observations   map[string]bool
				Degraded       bool
				Comparisons    []domain.ArchitectureComparison
				Hypotheses     []domain.ImprovementHypothesis
				HypothesisGaps []domain.HypothesisGap
				Plans          []domain.ValidationPlanEntry
				Provenance     []domain.SynthesisProvenance
				Status         domain.ImprovementSynthesisStatus
			}{resolved.mechanisms, resolved.mechanismGaps, resolved.observations, resolved.degraded, analysis.comparisons, analysis.hypotheses, analysis.gaps, analysis.plans, analysis.provenance, analysis.status}
			encoded, err := json.Marshal(projection)
			if err != nil {
				t.Fatal(err)
			}
			got := fmt.Sprintf("%x", sha256.Sum256(encoded))
			expected := map[bool]string{false: "384d47fc8e04ce1abea8d011ef224ce84969371e6557b040c14074eee68603f1", true: "2784813f67e98c4c5927259e561e54b6b66641d14fd1f4d72a77c4829a63ad1d"}[degraded]
			if got != expected {
				t.Fatalf("projection digest = %s, want %s", got, expected)
			}
			if err := validateProposalReferences([]domain.ImprovementProposal{{MechanismRefs: []string{"mech_" + strings.Repeat("3", 64)}, ProductObservationRefs: []string{"observation-1"}}}, resolved); err != nil {
				t.Fatal(err)
			}
			if !errors.Is(validateProposalReferences([]domain.ImprovementProposal{{MechanismRefs: []string{"missing"}}}, resolved), domain.ErrInvalidRunInput) {
				t.Fatal("unknown mechanism accepted")
			}
			if !errors.Is(validateProposalReferences([]domain.ImprovementProposal{{ProductObservationRefs: []string{"missing"}}}, resolved), domain.ErrInvalidRunInput) {
				t.Fatal("unknown observation accepted")
			}
		})
	}
}

func TestSynthesisSourcesPreserveFailures(t *testing.T) {
	sentinel := errors.New("read failed")
	cases := []struct {
		name   string
		change func(*synthesisSourceTestStore, []domain.SourceArtifactRef)
		want   error
	}{
		{"get", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) { s.getErr = sentinel }, sentinel},
		{"read", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) { s.readErr = sentinel }, sentinel},
		{"run kind", func(s *synthesisSourceTestStore, r []domain.SourceArtifactRef) {
			v := s.runs[r[0].RunID]
			v.RunKind = domain.RunKindMarketDiscovery
			s.runs[v.RunID] = v
		}, domain.ErrInvalidRunInput},
		{"running source", func(s *synthesisSourceTestStore, r []domain.SourceArtifactRef) {
			v := s.runs[r[0].RunID]
			v.Lifecycle = domain.RunRunning
			s.runs[v.RunID] = v
		}, domain.ErrSourceBasisUnavailable},
		{"missing artifact", func(s *synthesisSourceTestStore, r []domain.SourceArtifactRef) {
			v := s.runs[r[0].RunID]
			v.Artifacts = nil
			s.runs[v.RunID] = v
		}, domain.ErrSourceBasisUnavailable},
		{"index hash", func(s *synthesisSourceTestStore, r []domain.SourceArtifactRef) {
			v := s.runs[r[0].RunID]
			v.Artifacts[0].Hash = "wrong"
			s.runs[v.RunID] = v
		}, domain.ErrIntegrity},
		{"duplicate mechanism", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) {
			s.entries[domain.ArtifactMechanismCatalog] = append(s.entries[domain.ArtifactMechanismCatalog], s.entries[domain.ArtifactMechanismCatalog][0])
		}, domain.ErrIntegrity},
		{"malformed JSON", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) {
			s.entries[domain.ArtifactMechanismCatalog] = [][]byte{[]byte("{")}
		}, domain.ErrIntegrity},
		{"invalid mechanism ID", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) {
			s.entries[domain.ArtifactMechanismCatalog] = [][]byte{[]byte(`{"id":"bad"}`)}
		}, domain.ErrIntegrity},
		{"invalid gap", func(s *synthesisSourceTestStore, _ []domain.SourceArtifactRef) {
			s.entries[domain.ArtifactMechanismGaps] = [][]byte{[]byte(`{}`)}
		}, domain.ErrIntegrity},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store, sources := synthesisSourceFixture()
			tc.change(store, sources)
			result, err := resolveSourcesForTest(context.Background(), store, sources)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
			if !reflect.DeepEqual(result, resolvedSynthesisSources{}) {
				t.Fatalf("partial basis escaped failure: %+v", result)
			}
		})
	}
}

func TestSynthesisSourcesPaginationIntegrityAndHashFallback(t *testing.T) {
	for _, fault := range []string{"none", "run", "artifact", "hash", "empty cursor", "repeated cursor", "page bound"} {
		t.Run(fault, func(t *testing.T) {
			store, sources := synthesisSourceFixture()
			source := sources[0]
			run := store.runs[source.RunID]
			run.Artifacts[0].ArtifactHash = run.Artifacts[0].Hash
			run.Artifacts[0].Hash = ""
			store.runs[source.RunID] = run
			pageCount := 0
			store.page = func(r domain.ArtifactReadRequest) domain.ArtifactPage {
				pageCount++
				page := domain.ArtifactPage{RunID: r.RunID, Artifact: r.Artifact, ArtifactHash: "sha256:" + source.SHA256, Entries: [][]byte{}}
				if pageCount == 1 {
					page.Entries = store.entries[source.Artifact]
					next := "page-2"
					page.NextCursor = &next
				}
				switch fault {
				case "run":
					page.RunID = "other"
				case "artifact":
					page.Artifact = domain.ArtifactMarketCandidates
				case "hash":
					page.ArtifactHash = "other"
				case "empty cursor":
					next := ""
					page.NextCursor = &next
				case "repeated cursor":
					next := "page-2"
					page.NextCursor = &next
				case "page bound":
					next := fmt.Sprint(pageCount)
					page.NextCursor = &next
				}
				return page
			}
			result, err := resolveSourcesForTest(context.Background(), store, []domain.SourceArtifactRef{source})
			if fault == "none" {
				if err != nil || len(result.mechanisms) != 1 || pageCount != 2 || store.reads[1].Cursor != "page-2" {
					t.Fatalf("bad pagination: %+v %v", result, err)
				}
			} else if !errors.Is(err, domain.ErrIntegrity) {
				t.Fatalf("%s accepted: %v", fault, err)
			}
			if fault == "page bound" && pageCount != 1000 {
				t.Fatalf("unbounded pagination: %d", pageCount)
			}
		})
	}
}

func TestSynthesisSourceEntryBytesAreCopied(t *testing.T) {
	store, sources := synthesisSourceFixture()
	entries, err := readAllSynthesisArtifactEntries(context.Background(), store, sources[0], "sha256:"+sources[0].SHA256)
	if err != nil {
		t.Fatal(err)
	}
	original := append([]byte(nil), entries[0]...)
	store.entries[sources[0].Artifact][0][0] = '!'
	if !reflect.DeepEqual(entries[0], original) {
		t.Fatal("reader retains mutable page bytes")
	}
}
