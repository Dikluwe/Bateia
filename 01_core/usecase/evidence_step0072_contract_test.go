/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0072_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type step0072Content struct {
	tree  []domain.DirectoryNode
	files map[string]string
}

func (p step0072Content) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, domain.ErrManifestNotFound
}
func (p step0072Content) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return append([]domain.DirectoryNode(nil), p.tree...), nil
}
func (p step0072Content) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	files := make([]domain.ExtractedFile, 0, len(paths))
	for _, path := range paths {
		if content, ok := p.files[path]; ok {
			files = append(files, domain.ExtractedFile{Path: path, Content: content, Size: int64(len(content))})
		}
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}, Files: files}, nil
}
func (p step0072Content) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return candidate, nil
}

func step0072Catalog() domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "cpp-0072", Rules: []domain.CapabilityRule{
		{ID: "cpp-string", Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0072-contract", TargetHints: []string{"basic_string"}, CapabilityTerms: []string{"owning string abstraction"}, StrongEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"basic_string"}}}, SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: []string{"basic_string.h"}}}},
		{ID: "cpp-span", Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0072-contract", TargetHints: []string{"span"}, CapabilityTerms: []string{"contiguous span view"}, StrongEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"span<"}}}, SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: []string{"span"}}}},
	}}
}

func step0072Queries() []domain.EvidenceRequirementQuery {
	return []domain.EvidenceRequirementQuery{
		{RequirementKind: "capability", Requirement: "owning string abstraction", RuleID: "cpp-string", Terms: []string{"basic_string", "basic_string.h"}},
		{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}},
	}
}

func step0072Plan(queries []domain.EvidenceRequirementQuery) domain.EvidenceExtractionPlan {
	ref := domain.CommitRef{Owner: "local", Repo: "libstdcxx", Commit: "immutable-13"}
	return domain.EvidenceExtractionPlan{
		SchemaVersion: "1", Repository: ref,
		Candidate:          domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit, URL: "file:///fixture"},
		RequirementQueries: queries,
		Budget:             domain.ExtractionBudget{CodeFiles: 2, RequirementFiles: 1, Documents: 0, Manifests: 0, Tests: 0, TotalBytes: 64 * 1024, TreeDepth: 8},
	}
}

func step0072Paths(bundle domain.EvidenceBundle) []string {
	paths := []string{}
	for _, item := range bundle.Items {
		if item.Path != "@repository-metadata" {
			paths = append(paths, item.Path)
		}
	}
	sort.Strings(paths)
	return paths
}

func step0072Coverage(bundle domain.EvidenceBundle, requirement string) domain.RequirementCollectionCoverage {
	for _, coverage := range bundle.RequirementCoverage {
		if coverage.Requirement == requirement {
			return coverage
		}
	}
	return domain.RequirementCollectionCoverage{}
}

func TestStep0072LargeTreeReservesOneTargetPerRequirement(t *testing.T) {
	tree := make([]domain.DirectoryNode, 0, 405)
	files := map[string]string{
		"span":                "template<class T> class span {};",
		"bits/basic_string.h": "template<class C> class basic_string {};",
		"spanner":             "not a C++ header",
		"NOTICE":              "extensionless but unrelated",
	}
	for i := 0; i < 400; i++ {
		path := fmt.Sprintf("algorithm/distractor_%03d.hpp", i)
		tree = append(tree, domain.DirectoryNode{Path: path, Type: "file"})
		files[path] = "template<class T> void unrelated(T);"
	}
	tree = append(tree,
		domain.DirectoryNode{Path: "spanner", Type: "file"},
		domain.DirectoryNode{Path: "NOTICE", Type: "file"},
		domain.DirectoryNode{Path: "span", Type: "file"},
		domain.DirectoryNode{Path: "bits/basic_string.h", Type: "file"},
	)
	provider := step0072Content{tree: tree, files: files}
	originalTree, _ := json.Marshal(tree)
	originalQueries, _ := json.Marshal(step0072Queries())

	var canonical []byte
	for iteration := 0; iteration < 20; iteration++ {
		permuted := append([]domain.DirectoryNode(nil), tree...)
		queries := append([]domain.EvidenceRequirementQuery(nil), step0072Queries()...)
		if iteration%2 == 1 {
			for left, right := 0, len(permuted)-1; left < right; left, right = left+1, right-1 {
				permuted[left], permuted[right] = permuted[right], permuted[left]
			}
			queries[0], queries[1] = queries[1], queries[0]
		}
		provider.tree = permuted
		bundle, err := usecase.CollectEvidence(context.Background(), step0072Plan(queries), provider)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := step0072Paths(bundle), []string{"bits/basic_string.h", "span"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("iteração %d: reserva dirigida não venceu distratores: got=%v want=%v", iteration, got, want)
		}
		for _, requirement := range []string{"owning string abstraction", "contiguous span view"} {
			coverage := step0072Coverage(bundle, requirement)
			if coverage.Status != domain.RequirementTargetedExamined || len(coverage.ExaminedPaths) != 1 || len(coverage.SkippedPaths) != 0 {
				t.Fatalf("cobertura dirigida inválida para %q: %+v", requirement, coverage)
			}
		}
		observed, _ := json.Marshal(struct {
			Paths    []string                               `json:"paths"`
			Coverage []domain.RequirementCollectionCoverage `json:"coverage"`
		}{step0072Paths(bundle), bundle.RequirementCoverage})
		if iteration == 0 {
			canonical = observed
		} else if string(observed) != string(canonical) {
			t.Fatalf("reordenação alterou seleção/cobertura\nwant=%s\ngot=%s", canonical, observed)
		}
	}
	afterTree, _ := json.Marshal(tree)
	afterQueries, _ := json.Marshal(step0072Queries())
	if string(afterTree) != string(originalTree) || string(afterQueries) != string(originalQueries) {
		t.Fatal("coleta modificou entradas")
	}
}

func TestStep0072PathMatchDoesNotProveAndSkippedPlausiblePathStaysUnknown(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp api", Capabilities: []string{"contiguous span view"}}
	plan := step0072Plan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	plan.Budget.TotalBytes = 1 // metadado consome o budget antes do caminho plausível
	provider := step0072Content{tree: []domain.DirectoryNode{{Path: "span", Type: "file"}}, files: map[string]string{"span": "unrelated content"}}
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	coverage := step0072Coverage(bundle, "contiguous span view")
	if coverage.Status != domain.RequirementCollectionIncomplete || coverage.ReasonCode != "COLLECTION_INCOMPLETE" || len(coverage.SkippedPaths) != 1 {
		t.Fatalf("budget não preservou superfície omitida: %+v", coverage)
	}
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, step0072Catalog())
	if err != nil {
		t.Fatal(err)
	}
	got := report.Assessments[0].Capabilities[0]
	if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "COLLECTION_INCOMPLETE" || got.Confidence != 0 || len(report.Assessments[0].Gaps) != 0 {
		t.Fatalf("omissão de coleta atribuída ao candidato: evidence=%+v gaps=%v", got, report.Assessments[0].Gaps)
	}
}

func TestStep0072SmallFocusedTreeKeepsLegacyCollection(t *testing.T) {
	provider := step0072Content{
		tree:  []domain.DirectoryNode{{Path: "graph.go", Type: "file"}},
		files: map[string]string{"graph.go": "package graph\ntype DependencyGraph struct{}"},
	}
	plan := step0072Plan(nil)
	plan.Budget.CodeFiles = 20
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	if got := step0072Paths(bundle); !reflect.DeepEqual(got, []string{"graph.go"}) {
		t.Fatalf("regressão na árvore pequena: %v", got)
	}
}

func TestStep0072PlanningCompilesDeclaredRuleQueriesOnly(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp api", Capabilities: []string{"contiguous span view", "unknown facility"}}
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "local/libstdcxx@immutable-13", Decision: "ADVANCE"}}}
	plans, err := usecase.PlanEvidenceForIdea(idea, screening, step0072Catalog())
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || len(plans[0].RequirementQueries) != 1 {
		t.Fatalf("consultas devem existir apenas para regra declarada coberta: %+v", plans)
	}
	query := plans[0].RequirementQueries[0]
	if query.RuleID != "cpp-span" || query.Requirement != "contiguous span view" || !reflect.DeepEqual(query.Terms, []string{"span"}) {
		t.Fatalf("consulta compilada inesperada: %+v", query)
	}
}
