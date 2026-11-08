/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_integrity_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type evidenceIntegrityPort struct {
	calls       []string
	responseRef *domain.CommitRef
	nilResponse bool
	metadata    *domain.RepositoryCandidate
	metadataErr error
	extracted   domain.RepositoryCandidate
}

func (p *evidenceIntegrityPort) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	p.calls = append(p.calls, "manifest")
	return nil, errors.New("unused")
}
func (p *evidenceIntegrityPort) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.calls = append(p.calls, "metadata")
	if p.metadata != nil {
		c = *p.metadata
	}
	return c, p.metadataErr
}
func (p *evidenceIntegrityPort) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls = append(p.calls, "tree")
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
}
func (p *evidenceIntegrityPort) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	p.calls = append(p.calls, "extract")
	p.extracted = c
	if p.nilResponse {
		return nil, nil
	}
	ref := domain.CommitRef{Owner: c.Owner, Repo: c.Repo, Commit: c.CommitID}
	if p.responseRef != nil {
		ref = *p.responseRef
	}
	return &domain.ExtractedRepository{Ref: ref, Files: []domain.ExtractedFile{{Path: "README.md", Content: "ready candidates with bounded scheduling"}}}, nil
}
func evidenceIntegrityPlan() domain.EvidenceExtractionPlan {
	ref := domain.CommitRef{Owner: "example", Repo: "scheduler", Commit: strings.Repeat("a", 40)}
	return domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit}, Budget: defaultEvidenceBudget, Targets: []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation}}}
}
func TestEvidenceCollectionRejectsMismatchedAndNilResponse(t *testing.T) {
	plan := evidenceIntegrityPlan()
	for _, field := range []string{"owner", "repo", "commit", "empty", "nil"} {
		for _, coordinated := range []bool{false, true} {
			t.Run(field+"/coordinated="+map[bool]string{false: "false", true: "true"}[coordinated], func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("provider response caused panic: %v", p)
					}
				}()
				ref := plan.Repository
				switch field {
				case "owner":
					ref.Owner = "other"
				case "repo":
					ref.Repo = "other"
				case "commit":
					ref.Commit = strings.Repeat("b", 40)
				case "empty":
					ref = domain.CommitRef{}
				}
				p := &evidenceIntegrityPort{responseRef: &ref, nilResponse: field == "nil"}
				if coordinated {
					got, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "ready candidates"})
					if !errors.Is(err, domain.ErrIntegrity) || len(got.Bundles) != 0 || got.Exploration.ID != "" {
						t.Fatalf("integrity failure produced context: %v", err)
					}
				} else {
					got, err := CollectEvidence(context.Background(), plan, p)
					if !errors.Is(err, domain.ErrIntegrity) || len(got.Items) != 0 {
						t.Fatalf("foreign bytes admitted: %v", err)
					}
				}
			})
		}
	}
}
func TestEvidenceCollectionMetadataCannotChangeIdentityOrPinnedRevision(t *testing.T) {
	plan := evidenceIntegrityPlan()
	for _, field := range []string{"owner", "repo", "newer commit"} {
		t.Run(field, func(t *testing.T) {
			metadata := plan.Candidate
			switch field {
			case "owner":
				metadata.Owner = "other"
			case "repo":
				metadata.Repo = "other"
			case "newer commit":
				metadata.CommitID = strings.Repeat("b", 40)
			}
			p := &evidenceIntegrityPort{metadata: &metadata}
			before := algorithmContractJSON(t, plan)
			got, err := CollectEvidence(context.Background(), plan, p)
			if field != "newer commit" {
				if !errors.Is(err, domain.ErrIntegrity) || !reflect.DeepEqual(p.calls, []string{"metadata"}) || len(got.Items) != 0 {
					t.Fatalf("metadata crossed identity: %v %v", err, p.calls)
				}
			} else if err != nil || p.extracted.CommitID != plan.Repository.Commit || got.Repository != plan.Repository {
				t.Fatalf("pinned revision changed: %v %+v", err, p.extracted)
			}
			if string(before) != string(algorithmContractJSON(t, plan)) {
				t.Fatal("plan mutated")
			}
		})
	}
	plan.Repository.Commit = "HEAD"
	plan.Candidate.CommitID = "HEAD"
	resolved := plan.Candidate
	resolved.CommitID = strings.Repeat("b", 40)
	p := &evidenceIntegrityPort{metadata: &resolved}
	got, err := CollectEvidence(context.Background(), plan, p)
	if err != nil || got.Repository.Commit != resolved.CommitID || p.extracted.CommitID != resolved.CommitID {
		t.Fatal("explicit resolution failed", err)
	}
	p = &evidenceIntegrityPort{metadataErr: errors.New("unavailable")}
	report, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{plan}, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "ready candidates"})
	if err != nil || len(report.UnexploredRepositories) != 1 || report.Exploration.Status != "SEARCH_NEEDED" {
		t.Fatal("operational failure no longer explicit", err)
	}
}
func TestEvidenceCollectionAdmitsAllSchemasBeforeAnyEffects(t *testing.T) {
	for _, schema := range []string{"", "2", "invalid"} {
		for _, position := range []int{0, 1} {
			plan := evidenceIntegrityPlan()
			plans := []domain.EvidenceExtractionPlan{plan, plan}
			plans[position].SchemaVersion = schema
			p := &evidenceIntegrityPort{}
			got, err := CollectInsightEvidence(context.Background(), plans, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "ready candidates"})
			if !errors.Is(err, domain.ErrInvalidRunInput) || len(p.calls) != 0 || len(got.Bundles) != 0 {
				t.Fatalf("schema %q in slot %d reached effects: %v %v", schema, position, err, p.calls)
			}
		}
	}
	plan := evidenceIntegrityPlan()
	plan.SchemaVersion = "invalid"
	p := &evidenceIntegrityPort{}
	if _, err := CollectEvidence(context.Background(), plan, p); !errors.Is(err, domain.ErrInvalidRunInput) || len(p.calls) != 0 {
		t.Fatal("unknown generic schema admitted", err)
	}
	plan.SchemaVersion = ""
	p = &evidenceIntegrityPort{}
	if _, err := CollectEvidence(context.Background(), plan, p); err != nil {
		t.Fatal("generic legacy plan lost compatibility", err)
	}
}
func TestEvidenceFormulationIdentityStableAcrossEmptyCoverageRoundtrips(t *testing.T) {
	for _, empty := range []struct {
		name     string
		coverage []domain.RequirementCollectionCoverage
	}{{"nil", nil}, {"empty", []domain.RequirementCollectionCoverage{}}} {
		t.Run(empty.name, func(t *testing.T) {
			r := algorithmFixture()
			r.Sources[0].RequirementCoverage = empty.coverage
			r.Target.RequirementCoverage = empty.coverage
			first, err := FormulateAlgorithmInsights(context.Background(), r)
			if err != nil {
				t.Fatal(err)
			}
			for round := 0; round < 5; round++ {
				raw := algorithmContractJSON(t, r)
				var next domain.AlgorithmInsightRequest
				if err := json.Unmarshal(raw, &next); err != nil {
					t.Fatal(err)
				}
				got, err := FormulateAlgorithmInsights(context.Background(), next)
				if err != nil || !reflect.DeepEqual(first, got) {
					t.Fatalf("roundtrip changed formulation identity at %d: %v", round, err)
				}
				r = next
			}
		})
	}
}
func TestEvidenceFormulationPopulatedCompatibility(t *testing.T) {
	r := algorithmContractFixture()
	got, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	const want = "algorithms_16a7f584dcdd7468cb65d8502de60f2fad2b35e0a8b62d9ab286fa7aa7506e77"
	if got.ID != want {
		t.Fatalf("populated formulation ID changed: %s", got.ID)
	}
}

func TestEvidenceCollectionRejectsInvalidPlanInWholeBatch(t *testing.T) {
	for name, mutate := range map[string]func(*domain.EvidenceExtractionPlan){
		"missing identity":  func(p *domain.EvidenceExtractionPlan) { p.Repository.Owner = "" },
		"foreign candidate": func(p *domain.EvidenceExtractionPlan) { p.Candidate.Repo = "other" },
		"zero bytes":        func(p *domain.EvidenceExtractionPlan) { p.Budget.TotalBytes = 0 },
		"negative quota":    func(p *domain.EvidenceExtractionPlan) { p.Budget.Documents = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			valid := evidenceIntegrityPlan()
			invalid := valid
			mutate(&invalid)
			p := &evidenceIntegrityPort{}
			got, err := CollectInsightEvidence(context.Background(), []domain.EvidenceExtractionPlan{valid, invalid}, p, domain.InsightExplorationRequest{SchemaVersion: "1", Question: "ready candidates"})
			if !errors.Is(err, domain.ErrInvalidRunInput) || len(p.calls) != 0 || got.Exploration.ID != "" {
				t.Fatalf("invalid batch produced effects/context: %v %v", err, p.calls)
			}
		})
	}
}
