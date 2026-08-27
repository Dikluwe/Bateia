/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_test.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func testBundle(complete bool, items ...domain.CollectedEvidence) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	for i := range items {
		items[i].Ref = ref
		if items[i].ContentHash == "" {
			items[i].ContentHash = contentHash(items[i].Content)
		}
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: complete, TreeComplete: complete, Items: items}
}

func TestDocumentationAloneDoesNotProveCapability(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Provides dependency graph analysis."})
	got := assessCapability("grafo de dependências", b, findRule("grafo de dependências", DefaultEvidenceRules()))
	if got.Status == domain.EvidenceProven {
		t.Fatalf("documentação isolada produziu PROVEN: %+v", got)
	}
}

func TestWeakSignalNeverBecomesAbsence(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceRepository, Path: "@repository-metadata", Content: "static analysis"})
	got := assessCapability("static analysis", b, findRule("static analysis", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceInconclusive {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestImplementedSymbolAndRelatedTestCanProve(t *testing.T) {
	b := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "graph.go", Content: "type DependencyGraph struct { Nodes []string }"},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "graph_test.go", Content: "func TestDependencyGraph(t *testing.T) {}"},
	)
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceProven {
		t.Fatalf("status = %s, esperado PROVEN", got.Status)
	}
	if len(got.Sources) == 0 || got.Sources[0].Ref.Commit == "" {
		t.Fatal("fonte não rastreável")
	}
}

func TestFailureNeverBecomesAbsent(t *testing.T) {
	b := testBundle(false)
	b.Failures = []domain.ExecutionFailure{{Scope: "org/repo@abc123", Stage: "content", Message: "timeout"}}
	got := assessCapability("static analysis", b, findRule("static analysis", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceInconclusive {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestVendorNoiseCannotProve(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "vendor/k8s/graph_test.go", Content: "func TestDependencyGraph(t *testing.T) {}"})
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status == domain.EvidenceProven {
		t.Fatalf("vendor produziu PROVEN: %+v", got)
	}
}

func TestAssessmentIsDeterministicAndDoesNotMutateBundle(t *testing.T) {
	itemA := domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "z_test.go", Content: "func TestStaticAnalyzer(t *testing.T) {}"}
	itemB := domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "static analysis"}
	a := testBundle(true, itemA, itemB)
	before := a
	b := testBundle(true, itemB, itemA)
	rule := findRule("static analysis", DefaultEvidenceRules())
	left := assessCapability("static analysis", a, rule)
	right := assessCapability("static analysis", b, rule)
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("ordem alterou resultado\n%+v\n%+v", left, right)
	}
	if !reflect.DeepEqual(a, before) {
		t.Fatal("entrada foi modificada")
	}
}

func TestPlanAcceptsResolutionPendingAndSelectsDeterministically(t *testing.T) {
	report := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{
		{Repository: "z/drop@1", Decision: "DROP", RankingScore: .1}, {Repository: "a/drop@1", Decision: "DROP", RankingScore: .2},
		{Repository: "org/a@abc", Decision: "ADVANCE", RankingScore: .9},
	}}
	plans, err := PlanEvidence(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 3 || plans[1].Repository.Owner != "a" {
		t.Fatalf("seleção inesperada: %+v", plans)
	}
	report.Decisions = []domain.ScreeningDecision{{Repository: "org/mutable", Decision: "ADVANCE"}}
	pending, err := PlanEvidence(report)
	if err != nil {
		t.Fatal(err)
	}
	if pending[0].Repository.Commit != "HEAD" {
		t.Fatalf("esperado HEAD pendente, recebido %s", pending[0].Repository.Commit)
	}
}

func TestPlanIncludesExplicitCalibrationControls(t *testing.T) {
	report := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/main", Decision: "ADVANCE"}, {Repository: "org/control", Decision: "DROP"}}}
	plans, err := PlanEvidenceForIdeaAndRepositories(domain.IdeaSpec{Capabilities: []string{"static analysis"}}, report, DefaultEvidenceRules(), []string{"org/control", "extra/noise"})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, plan := range plans {
		seen[plan.Repository.Owner+"/"+plan.Repository.Repo] = true
	}
	if !seen["org/control"] || !seen["extra/noise"] {
		t.Fatalf("controles ausentes: %+v", plans)
	}
}

func TestPlanDeclaresProblemAsEvidenceTarget(t *testing.T) {
	report := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo", Decision: "ADVANCE"}}}
	idea := domain.IdeaSpec{Problem: "analisar código sem executar o programa", Capabilities: []string{"static analysis"}}
	plans, err := PlanEvidenceForIdea(idea, report, DefaultEvidenceRules())
	if err != nil {
		t.Fatal(err)
	}
	got := plans[0].Targets[0].Capabilities
	if !containsNormalized(got, idea.Problem) || !containsNormalized(got, idea.Capabilities[0]) {
		t.Fatalf("problema e capacidade não declarados no alvo: %v", got)
	}
}

type evidenceProvider struct {
	files   []domain.ExtractedFile
	treeErr error
}

func (p evidenceProvider) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("unused")
}
func (p evidenceProvider) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	if p.treeErr != nil {
		return nil, p.treeErr
	}
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "graph_test.go", Type: "file"}}, nil
}
func (p evidenceProvider) ExtractRepoFiles(_ context.Context, _ domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{Files: p.files}, nil
}
func (p evidenceProvider) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	candidate.CommitID = "resolved123"
	return candidate, nil
}

func TestCollectDeduplicatesByHashAndPreservesPartialFailure(t *testing.T) {
	plan := domain.EvidenceExtractionPlan{Repository: domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}, Candidate: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc"}, Budget: defaultEvidenceBudget}
	p := evidenceProvider{treeErr: errors.New("tree unavailable"), files: []domain.ExtractedFile{{Path: "README.md", Content: "same"}, {Path: "COPY.md", Content: "same"}}}
	b, err := CollectEvidence(context.Background(), plan, p)
	if err != nil {
		t.Fatal(err)
	}
	if b.Complete || len(b.Failures) != 1 || len(b.Items) != 2 {
		t.Fatalf("bundle parcial inesperado: %+v", b)
	}
}

func TestMissingRelevantTargetMakesAbsenceInconclusive(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	plan := domain.EvidenceExtractionPlan{Repository: ref, Candidate: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "abc"}, Budget: defaultEvidenceBudget, Targets: []domain.ExtractionTarget{{Path: "README.md", Kind: domain.EvidenceDocumentation, Capabilities: []string{"dependency graph"}}}}
	b, err := CollectEvidence(context.Background(), plan, evidenceProvider{files: []domain.ExtractedFile{{Path: "go.mod", Content: "module x"}}})
	if err != nil {
		t.Fatal(err)
	}
	if b.Complete {
		t.Fatal("alvo relevante ausente manteve bundle completo")
	}
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	if got.Status != domain.EvidenceInconclusive {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestCommentDistantFromDeclarationDoesNotProduceSymbol(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "graph.go", Content: "// dependency graph is planned\nfunc unrelated() {}"})
	got := assessCapability("dependency graph", b, findRule("dependency graph", DefaultEvidenceRules()))
	for _, source := range got.Sources {
		if source.Kind == domain.EvidenceSymbol {
			t.Fatalf("comentário virou símbolo: %+v", source)
		}
	}
}

func TestNegativeEvidenceMustBeContextual(t *testing.T) {
	rule := findRule("static analysis", DefaultEvidenceRules())
	distant := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Static analysis is available.\nA separate deployment feature is not implemented."})
	if got := assessCapability("static analysis", distant, rule); got.Status == domain.EvidenceConflict {
		t.Fatal("negativa distante criou conflito")
	}
	explicit := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Static analysis is not implemented."})
	if got := assessCapability("static analysis", explicit, rule); got.Status != domain.EvidenceConflict {
		t.Fatalf("status = %s", got.Status)
	}
}

func TestSymbolPreservesNameLineAndDrivesReusableComponent(t *testing.T) {
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc123"}
	b := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "dependency graph"},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "src/graph.go", Content: "package graph\n\ntype DependencyGraph struct{}"},
	)
	report, err := AssessEvidence(domain.IdeaSpec{Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "org/repo", Decision: "ADVANCE"}}}, []domain.EvidenceBundle{b}, DefaultEvidenceRules())
	if err != nil {
		t.Fatal(err)
	}
	a := report.Assessments[0]
	if a.ScreeningDecision.Decision != "ADVANCE" {
		t.Fatal("decisão da peneira perdida após SHA")
	}
	if len(a.ReusableParts) != 1 || a.ReusableParts[0].Path != "src/graph.go" || a.ReusableParts[0].Symbols[0] != "DependencyGraph" {
		t.Fatalf("componente inesperado: %+v", a.ReusableParts)
	}
	s := a.ReusableParts[0].Evidence[0]
	if s.Line != 3 || s.Symbol != "DependencyGraph" {
		t.Fatalf("localização inesperada: %+v", s)
	}
	if a.Repository != ref {
		t.Fatalf("ref inesperada: %+v", a.Repository)
	}
}

func TestLicenseDetectionAvoidsMITSubstring(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceLicense, Path: "LICENSE", Content: "Limitations and submissions are described here."})
	if got := assessLicense(b); got.Status == domain.EvidenceProven {
		t.Fatalf("substring acidental identificou licença: %+v", got)
	}
}

func TestExplicitLocalOperationCanBeProven(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "The CLI can run locally and works offline."})
	if got := assessOperation(b); got.LocalStatus != domain.EvidenceProven {
		t.Fatalf("status = %s", got.LocalStatus)
	}
}

func TestUnrelatedSendAndAPIDoNotCreateOperationConflict(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Send an email from your application.\nThe API can also run locally."})
	if got := assessOperation(b); got.LocalStatus == domain.EvidenceConflict {
		t.Fatalf("conflito operacional falso: %+v", got)
	}
}

func TestExplicitRepositoryUploadCreatesOperationConflict(t *testing.T) {
	b := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "The CLI uploads your repository to the cloud API."})
	if got := assessOperation(b); got.LocalStatus != domain.EvidenceConflict {
		t.Fatalf("status = %s", got.LocalStatus)
	}
}

func TestExplicitMaintenanceEndOverridesRecentMetadata(t *testing.T) {
	b := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceRepository, Path: "@repository-metadata", Content: `{"owner":"org","repo":"repo","updated_at":"2026-08-01T00:00:00Z"}`},
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Attention: this tool is no longer maintained."},
	)
	if got := assessMaintenance(b); got.Status != domain.EvidenceConflict {
		t.Fatalf("status = %s; warnings=%v", got.Status, got.Warnings)
	}
}

func TestStaticAnalysisRuleCoversBenchmarkProblemWording(t *testing.T) {
	if got := findRule("analisar código sem executar o programa", DefaultEvidenceRules()); got == nil || got.ID != "static-analysis" {
		t.Fatalf("regra = %+v", got)
	}
}

func TestCalibratedStructuralPatternsRequireRelatedDeclaration(t *testing.T) {
	rule := findRule("blast radius", DefaultEvidenceRules())
	positive := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Reports blast radius for a code change."},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "conformance_test.go", Content: "func TestNoSpilloverWhenTheChangeStaysInScope(t *testing.T) {}"},
	)
	if got := assessCapability("blast radius", positive, rule); got.Status != domain.EvidenceProven {
		t.Fatalf("positivo = %s", got.Status)
	}
	negative := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "network_test.go", Content: "func TestNetworkSpilloverBucket(t *testing.T) {}"})
	if got := assessCapability("blast radius", negative, rule); got.Status == domain.EvidenceProven {
		t.Fatalf("spillover não estrutural produziu PROVEN: %+v", got)
	}
}

func TestCalibratedRecipeAndEventPatternsRejectMentions(t *testing.T) {
	recipeRule := findRule("recipe import", DefaultEvidenceRules())
	recipe := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Imports recipes from Tandoor."},
		domain.CollectedEvidence{Kind: domain.EvidenceTest, Path: "tandoor_test.go", Content: "func TestTandoorImport(t *testing.T) {}"},
	)
	if got := assessCapability("recipe import", recipe, recipeRule); got.Status != domain.EvidenceProven {
		t.Fatalf("recipe = %s", got.Status)
	}
	mention := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "widgets.go", Content: `var label = "Tandoor import"`})
	if got := assessCapability("recipe import", mention, recipeRule); got.Status == domain.EvidenceProven {
		t.Fatalf("menção produziu PROVEN: %+v", got)
	}
	eventRule := findRule("asynchronous event processing", DefaultEvidenceRules())
	event := testBundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Path: "README.md", Content: "Event-driven architecture with asynchronous events."},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "event_dispatcher.go", Content: "type EventDispatcher struct{}"},
	)
	if got := assessCapability("asynchronous event processing", event, eventRule); got.Status != domain.EvidenceProven {
		t.Fatalf("event = %s", got.Status)
	}
}

func TestRecipeOrganizationSymbolsAreSpecific(t *testing.T) {
	rule := findRule("recipe organization", DefaultEvidenceRules())
	for _, declaration := range []string{"class MelaRecipes(Integration):", "class RepositoryRecipes(BaseRepository):", "class RecipeController:"} {
		bundle := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "recipes.py", Content: declaration})
		if got := assessCapability("recipe organization", bundle, rule); got.Status != domain.EvidenceProven {
			t.Fatalf("%q = %s", declaration, got.Status)
		}
	}
	for _, declaration := range []string{"class Repository:", "class Controller:", "class Recipe:"} {
		bundle := testBundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "models.py", Content: declaration})
		if got := assessCapability("recipe organization", bundle, rule); got.Status == domain.EvidenceProven {
			t.Fatalf("genérico %q produziu PROVEN", declaration)
		}
	}
}

func TestEvidenceTreePrioritizesCapabilityPaths(t *testing.T) {
	tree := []domain.DirectoryNode{{Path: "aaa/command.go", Type: "file"}, {Path: "msg/event_dispatcher.go", Type: "file"}, {Path: "zzz.go", Type: "file"}}
	targets := []domain.ExtractionTarget{{Path: "README.md", Capabilities: []string{"asynchronous event processing"}}}
	got := prioritizeEvidenceTree(tree, targets)
	if got[0].Path != "msg/event_dispatcher.go" {
		t.Fatalf("prioridade inesperada: %+v", got)
	}
}

func TestEvidenceMetricsAreReproducible(t *testing.T) {
	valid := true
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}
	source := domain.EvidenceSource{Kind: domain.EvidenceSymbol, Ref: ref, Path: "graph.go", Symbol: "DependencyGraph", Line: 1, ContentHash: "hash", RuleID: "dependency-graph", Strength: .9}
	report := domain.AssessmentReport{Assessments: []domain.CandidateAssessment{{Repository: ref, Capabilities: []domain.CapabilityEvidence{{Capability: "dependency graph", Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{source}}}}}}
	bundles := []domain.EvidenceBundle{{Targets: []domain.CollectionTargetResult{{Status: domain.TargetCollected}}}}
	judgments := []domain.EvidenceJudgment{{Repository: "org/repo", Capability: "dependency graph", Expected: domain.EvidenceProven, ComponentValid: &valid}}
	a := EvaluateEvidenceMetrics(report, bundles, judgments)
	b := EvaluateEvidenceMetrics(report, bundles, judgments)
	if !reflect.DeepEqual(a, b) || a.EvidencePrecision == nil || *a.EvidencePrecision != 1 || a.SourceTraceRate != 1 {
		t.Fatalf("métricas inesperadas: %+v", a)
	}
}

func TestComponentPrecisionOnlyReviewsEmittedComponents(t *testing.T) {
	valid := true
	ref := domain.CommitRef{Owner: "org", Repo: "repo", Commit: "abc"}
	evidence := domain.CapabilityEvidence{Capability: "dependency graph", Status: domain.EvidenceProven}
	judgments := []domain.EvidenceJudgment{{Repository: "org/repo", Capability: "dependency graph", Expected: domain.EvidenceProven, ComponentValid: &valid}}
	without := EvaluateEvidenceMetrics(domain.AssessmentReport{Assessments: []domain.CandidateAssessment{{Repository: ref, Capabilities: []domain.CapabilityEvidence{evidence}}}}, nil, judgments)
	if without.ComponentPrecision != nil {
		t.Fatalf("componente ausente entrou no denominador: %+v", without)
	}
	assessment := domain.CandidateAssessment{Repository: ref, Capabilities: []domain.CapabilityEvidence{evidence}, ReusableParts: []domain.ReusableComponent{{Capabilities: []string{"dependency graph"}}}}
	with := EvaluateEvidenceMetrics(domain.AssessmentReport{Assessments: []domain.CandidateAssessment{assessment}}, nil, judgments)
	if with.ComponentPrecision == nil || *with.ComponentPrecision != 1 {
		t.Fatalf("precisão = %+v", with.ComponentPrecision)
	}
}
