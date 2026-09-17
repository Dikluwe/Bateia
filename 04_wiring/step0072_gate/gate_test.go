/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0072_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-16
 */

package step0072_gate_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	localinfra "bateia/03_infra/local"
)

type oracleContent struct {
	tree      []domain.DirectoryNode
	files     map[string]string
	treeErr   error
	requested [][]string
}

func (p *oracleContent) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, domain.ErrManifestNotFound
}

func (p *oracleContent) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return append([]domain.DirectoryNode(nil), p.tree...), p.treeErr
}

func (p *oracleContent) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.requested = append(p.requested, append([]string(nil), paths...))
	files := make([]domain.ExtractedFile, 0, len(paths))
	for _, path := range paths {
		if content, ok := p.files[path]; ok {
			files = append(files, domain.ExtractedFile{Path: path, Content: content, Size: int64(len(content))})
		}
	}
	if len(files) == 0 {
		return nil, errors.New("nenhum arquivo solicitado foi extraído")
	}
	return &domain.ExtractedRepository{
		Ref:   domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID},
		Files: files,
	}, nil
}

func (p *oracleContent) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return candidate, nil
}

func oracleCatalog() domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "step0072-oracle", Rules: []domain.CapabilityRule{
		{
			ID: "cpp-string", Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0072-oracle",
			TargetHints: []string{"basic_string"}, CapabilityTerms: []string{"owning string abstraction"},
			StrongEvidence:     []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"basic_string"}}},
			SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: []string{"basic_string.h"}}},
		},
		{
			ID: "cpp-span", Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0072-oracle",
			TargetHints: []string{"span"}, CapabilityTerms: []string{"contiguous span view"},
			StrongEvidence:     []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"span<"}}},
			SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: []string{"span"}}},
		},
		{
			ID: "cpp-hypothesis", Basis: domain.EvidenceRuleHypothesis, Family: "cpp-api", Provenance: "step-0072-oracle",
			TargetHints: []string{"experimental_view"}, CapabilityTerms: []string{"experimental view"},
			StrongEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"experimental_view"}}},
		},
	}}
}

func oracleQueries() []domain.EvidenceRequirementQuery {
	return []domain.EvidenceRequirementQuery{
		{RequirementKind: "capability", Requirement: "owning string abstraction", RuleID: "cpp-string", Terms: []string{"basic_string", "basic_string.h"}},
		{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}},
	}
}

func oraclePlan(queries []domain.EvidenceRequirementQuery) domain.EvidenceExtractionPlan {
	ref := domain.CommitRef{Owner: "local", Repo: "libstdcxx", Commit: "immutable-13"}
	return domain.EvidenceExtractionPlan{
		SchemaVersion: "1", Repository: ref,
		Candidate:          domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit, URL: "file:///oracle"},
		RequirementQueries: queries,
		Budget: domain.ExtractionBudget{
			CodeFiles: 2, RequirementFiles: 1, Documents: 0, Manifests: 0, Tests: 0,
			TotalBytes: 128 * 1024, TreeDepth: 8,
		},
	}
}

func coverageFor(bundle domain.EvidenceBundle, requirement string) domain.RequirementCollectionCoverage {
	for _, coverage := range bundle.RequirementCoverage {
		if coverage.Requirement == requirement {
			return coverage
		}
	}
	return domain.RequirementCollectionCoverage{}
}

func collectedPaths(bundle domain.EvidenceBundle) []string {
	paths := []string{}
	for _, item := range bundle.Items {
		if item.Path != "@repository-metadata" {
			paths = append(paths, item.Path)
		}
	}
	sort.Strings(paths)
	return paths
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

// M13 and determinism: catalog gaps and hypotheses cannot synthesize declared queries.
func TestStep0072GatePlanningIsDeclaredDeterministicAndImmutable(t *testing.T) {
	baseIdea := domain.IdeaSpec{Problem: "cpp api", Capabilities: []string{
		"contiguous span view", "owning string abstraction", "unknown facility", "experimental view",
	}}
	baseCatalog := oracleCatalog()
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "local/libstdcxx@immutable-13", Decision: "ADVANCE"}}}
	ideaBefore, _ := json.Marshal(baseIdea)
	catalogBefore, _ := json.Marshal(baseCatalog)

	var canonical []byte
	for iteration := 0; iteration < 20; iteration++ {
		idea := baseIdea
		idea.Capabilities = append([]string(nil), baseIdea.Capabilities...)
		catalog := baseCatalog
		catalog.Rules = append([]domain.CapabilityRule(nil), baseCatalog.Rules...)
		if iteration%2 == 1 {
			sort.Sort(sort.Reverse(sort.StringSlice(idea.Capabilities)))
			for left, right := 0, len(catalog.Rules)-1; left < right; left, right = left+1, right-1 {
				catalog.Rules[left], catalog.Rules[right] = catalog.Rules[right], catalog.Rules[left]
			}
			for i := range catalog.Rules {
				sort.Sort(sort.Reverse(sort.StringSlice(catalog.Rules[i].TargetHints)))
			}
		}
		plans, err := usecase.PlanEvidenceForIdea(idea, screening, catalog)
		if err != nil {
			t.Fatal(err)
		}
		if len(plans) != 1 || len(plans[0].RequirementQueries) != 2 {
			t.Fatalf("iteração %d: somente duas regras declaradas devem gerar consultas: %+v", iteration, plans)
		}
		for _, query := range plans[0].RequirementQueries {
			if query.RuleID == "cpp-hypothesis" || query.Requirement == "unknown facility" {
				t.Fatalf("iteração %d: regra não declarada/coberta gerou consulta: %+v", iteration, query)
			}
		}
		encoded, _ := json.Marshal(plans[0].RequirementQueries)
		if iteration == 0 {
			canonical = encoded
		} else if string(encoded) != string(canonical) {
			t.Fatalf("iteração %d: ordem alterou consultas\nwant=%s\ngot=%s", iteration, canonical, encoded)
		}
	}
	ideaAfter, _ := json.Marshal(baseIdea)
	catalogAfter, _ := json.Marshal(baseCatalog)
	if string(ideaBefore) != string(ideaAfter) || string(catalogBefore) != string(catalogAfter) {
		t.Fatal("planejamento modificou ideia ou catálogo de entrada")
	}
}

// M02, M03, M04, M05, M06, M11: reservation wins before general lexical sampling.
func TestStep0072GateTargetedReservationBeatsDistractorsAcrossPermutations(t *testing.T) {
	tree := make([]domain.DirectoryNode, 0, 128)
	files := map[string]string{
		"span": "template<class T> class span {};", "bits/basic_string.h": "template<class C> class basic_string {};",
		"spanner": "extensionless false positive", "NOTICE": "unrelated extensionless", "span.txt": "not a header",
		"docs/basic_string.md": "documentation is not the header",
	}
	for i := 0; i < 120; i++ {
		path := fmt.Sprintf("algorithm/a_distractor_%03d.hpp", i)
		tree = append(tree, domain.DirectoryNode{Path: path, Type: "file"})
		files[path] = "template<class T> void unrelated(T);"
	}
	for _, path := range []string{"spanner", "NOTICE", "span.txt", "docs/basic_string.md", "span", "bits/basic_string.h"} {
		tree = append(tree, domain.DirectoryNode{Path: path, Type: "file"})
	}
	treeBefore, _ := json.Marshal(tree)
	queriesBefore, _ := json.Marshal(oracleQueries())

	var canonical []byte
	for iteration := 0; iteration < 20; iteration++ {
		permuted := append([]domain.DirectoryNode(nil), tree...)
		queries := oracleQueries()
		if iteration%2 == 1 {
			for left, right := 0, len(permuted)-1; left < right; left, right = left+1, right-1 {
				permuted[left], permuted[right] = permuted[right], permuted[left]
			}
			queries[0], queries[1] = queries[1], queries[0]
			for i := range queries {
				sort.Sort(sort.Reverse(sort.StringSlice(queries[i].Terms)))
			}
		}
		provider := &oracleContent{tree: permuted, files: files}
		bundle, err := usecase.CollectEvidence(context.Background(), oraclePlan(queries), provider)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"bits/basic_string.h", "span"}
		if got := collectedPaths(bundle); !reflect.DeepEqual(got, want) {
			t.Fatalf("iteração %d: seleção dirigida inválida: got=%v want=%v", iteration, got, want)
		}
		if len(provider.requested) != 1 || len(provider.requested[0]) < 2 ||
			!contains(provider.requested[0][:2], "span") || !contains(provider.requested[0][:2], "bits/basic_string.h") {
			t.Fatalf("iteração %d: amostra geral consumiu reserva: %v", iteration, provider.requested)
		}
		for _, requirement := range []string{"owning string abstraction", "contiguous span view"} {
			coverage := coverageFor(bundle, requirement)
			if coverage.Status != domain.RequirementTargetedExamined || !coverage.SearchComplete || len(coverage.ExaminedPaths) != 1 || len(coverage.SkippedPaths) != 0 {
				t.Fatalf("iteração %d: cobertura inválida para %q: %+v", iteration, requirement, coverage)
			}
		}
		spanCoverage := coverageFor(bundle, "contiguous span view")
		if contains(spanCoverage.CandidatePaths, "spanner") || contains(spanCoverage.CandidatePaths, "span.txt") || contains(spanCoverage.CandidatePaths, "NOTICE") {
			t.Fatalf("iteração %d: falso positivo extensionless/substring: %+v", iteration, spanCoverage)
		}
		stringCoverage := coverageFor(bundle, "owning string abstraction")
		if contains(stringCoverage.CandidatePaths, "docs/basic_string.md") {
			t.Fatalf("iteração %d: documento promovido como header: %+v", iteration, stringCoverage)
		}
		observed, _ := json.Marshal(struct {
			Paths    []string                               `json:"paths"`
			Coverage []domain.RequirementCollectionCoverage `json:"coverage"`
		}{collectedPaths(bundle), bundle.RequirementCoverage})
		if iteration == 0 {
			canonical = observed
		} else if string(observed) != string(canonical) {
			t.Fatalf("iteração %d: observação não canônica\nwant=%s\ngot=%s", iteration, canonical, observed)
		}
	}
	treeAfter, _ := json.Marshal(tree)
	queriesAfter, _ := json.Marshal(oracleQueries())
	if string(treeBefore) != string(treeAfter) || string(queriesBefore) != string(queriesAfter) {
		t.Fatal("coleta modificou árvore ou consultas")
	}
}

// M07: locating a plausible name is retrieval only, never semantic evidence.
func TestStep0072GatePathMatchNeverPromotesCapability(t *testing.T) {
	provider := &oracleContent{
		tree:  []domain.DirectoryNode{{Path: "span", Type: "file"}},
		files: map[string]string{"span": "plain unrelated prose without implementation"},
	}
	plan := oraclePlan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	plan.Budget.CodeFiles = 1
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	report, err := usecase.AssessEvidence(
		domain.IdeaSpec{Problem: "cpp api", Capabilities: []string{"contiguous span view"}},
		domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, oracleCatalog(),
	)
	if err != nil {
		t.Fatal(err)
	}
	evidence := report.Assessments[0].Capabilities[0]
	if evidence.Status == domain.EvidenceProven || evidence.Status == domain.EvidencePartial || len(report.Assessments[0].ReusableParts) != 0 {
		t.Fatalf("nome de caminho foi promovido a evidência sem conteúdo: evidence=%+v reusable=%+v", evidence, report.Assessments[0].ReusableParts)
	}
}

// M08, M09, M10: an omitted plausible path remains explicit candidate-independent Unknown.
func TestStep0072GateSkippedPlausiblePathIsExplicitUnknown(t *testing.T) {
	provider := &oracleContent{
		tree:  []domain.DirectoryNode{{Path: "span", Type: "file"}},
		files: map[string]string{"span": "template<class T> class span {};"},
	}
	plan := oraclePlan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	plan.Budget.TotalBytes = 1
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	coverage := coverageFor(bundle, "contiguous span view")
	if coverage.Status != domain.RequirementCollectionIncomplete || coverage.ReasonCode != "COLLECTION_INCOMPLETE" ||
		!reflect.DeepEqual(coverage.CandidatePaths, []string{"span"}) || len(coverage.ExaminedPaths) != 0 || !reflect.DeepEqual(coverage.SkippedPaths, []string{"span"}) {
		t.Fatalf("superfície omitida não ficou explícita: %+v", coverage)
	}
	report, err := usecase.AssessEvidence(
		domain.IdeaSpec{Problem: "cpp api", Capabilities: []string{"contiguous span view"}},
		domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, oracleCatalog(),
	)
	if err != nil {
		t.Fatal(err)
	}
	evidence := report.Assessments[0].Capabilities[0]
	if evidence.Status != domain.EvidenceInconclusive || evidence.ReasonCode != "COLLECTION_INCOMPLETE" || evidence.Confidence != 0 || len(report.Assessments[0].Gaps) != 0 {
		t.Fatalf("Unknown foi atribuído ao candidato: evidence=%+v gaps=%v", evidence, report.Assessments[0].Gaps)
	}
}

// M12: shared evidence is charged/read once while represented for both requirements.
func TestStep0072GateSharedTargetIsDeduplicatedWithoutLosingCoverage(t *testing.T) {
	queries := []domain.EvidenceRequirementQuery{
		{RequirementKind: "capability", Requirement: "first view", RuleID: "first", Terms: []string{"shared.hpp"}},
		{RequirementKind: "capability", Requirement: "second view", RuleID: "second", Terms: []string{"shared.hpp"}},
	}
	provider := &oracleContent{
		tree:  []domain.DirectoryNode{{Path: "include/shared.hpp", Type: "file"}},
		files: map[string]string{"include/shared.hpp": "struct shared {};"},
	}
	plan := oraclePlan(queries)
	plan.Budget.CodeFiles = 1
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	if got := collectedPaths(bundle); !reflect.DeepEqual(got, []string{"include/shared.hpp"}) {
		t.Fatalf("arquivo compartilhado lido mais de uma vez ou perdido: %v", got)
	}
	if len(provider.requested) != 1 {
		t.Fatalf("provedor chamado mais de uma vez: %v", provider.requested)
	}
	count := 0
	for _, path := range provider.requested[0] {
		if path == "include/shared.hpp" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("arquivo compartilhado solicitado %d vezes: %v", count, provider.requested[0])
	}
	for _, requirement := range []string{"first view", "second view"} {
		coverage := coverageFor(bundle, requirement)
		if coverage.Status != domain.RequirementTargetedExamined || !reflect.DeepEqual(coverage.ExaminedPaths, []string{"include/shared.hpp"}) {
			t.Fatalf("deduplicação apagou cobertura de %q: %+v", requirement, coverage)
		}
	}
	var sum int64
	for _, item := range bundle.Items {
		sum += int64(len(item.Content))
	}
	if bundle.BytesRead != sum {
		t.Fatalf("bytes cobrados mais de uma vez: bytes_read=%d soma_itens=%d", bundle.BytesRead, sum)
	}
}

// Opaque tree failures cannot become complete searches or absence.
func TestStep0072GateTreeFailurePreservesIncompleteUnknown(t *testing.T) {
	provider := &oracleContent{treeErr: errors.New("enumeração interrompida"), files: map[string]string{}}
	plan := oraclePlan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.TreeComplete || bundle.Complete {
		t.Fatalf("falha de árvore foi convertida em sucesso: %+v", bundle)
	}
	coverage := coverageFor(bundle, "contiguous span view")
	if coverage.Status != domain.RequirementCollectionIncomplete || coverage.SearchComplete || coverage.ReasonCode != "COLLECTION_INCOMPLETE" {
		t.Fatalf("falha de enumeração não preservou Unknown: %+v", coverage)
	}
}

// Complete search with no match is diagnostic only and does not invent absence evidence.
func TestStep0072GateCompleteSearchCanReportNoPlausiblePath(t *testing.T) {
	provider := &oracleContent{
		tree:  []domain.DirectoryNode{{Path: "vector.hpp", Type: "file"}},
		files: map[string]string{"vector.hpp": "template<class T> class vector {};"},
	}
	plan := oraclePlan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	plan.Budget.CodeFiles = 1
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	coverage := coverageFor(bundle, "contiguous span view")
	if coverage.Status != domain.RequirementNoPlausiblePath || !coverage.SearchComplete || len(coverage.CandidatePaths) != 0 || len(coverage.ExaminedPaths) != 0 || len(coverage.SkippedPaths) != 0 {
		t.Fatalf("busca completa sem candidato mal diagnosticada: %+v", coverage)
	}
}

func TestStep0072GateZeroRequirementQuotaUsesPositiveBoundedDefault(t *testing.T) {
	provider := &oracleContent{
		tree: []domain.DirectoryNode{
			{Path: "aaa_distractor.hpp", Type: "file"},
			{Path: "span", Type: "file"},
		},
		files: map[string]string{
			"aaa_distractor.hpp": "void unrelated();",
			"span":               "template<class T> class span {};",
		},
	}
	plan := oraclePlan([]domain.EvidenceRequirementQuery{{RequirementKind: "capability", Requirement: "contiguous span view", RuleID: "cpp-span", Terms: []string{"span"}}})
	plan.Budget.RequirementFiles = 0
	plan.Budget.CodeFiles = 1
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	if got := collectedPaths(bundle); !reflect.DeepEqual(got, []string{"span"}) {
		t.Fatalf("quota dirigida zero não aplicou default positivo: %v", got)
	}
	coverage := coverageFor(bundle, "contiguous span view")
	if coverage.Status != domain.RequirementTargetedExamined || !reflect.DeepEqual(coverage.ExaminedPaths, []string{"span"}) {
		t.Fatalf("default dirigido não materializou cobertura: %+v", coverage)
	}
}

// M14: targeted recovery is additive and preserves the focused legacy path.
func TestStep0072GateSmallGoTreeKeepsLegacyCollection(t *testing.T) {
	provider := &oracleContent{
		tree:  []domain.DirectoryNode{{Path: "graph.go", Type: "file"}},
		files: map[string]string{"graph.go": "package graph\ntype DependencyGraph struct{}"},
	}
	plan := oraclePlan(nil)
	plan.Budget.CodeFiles = 20
	bundle, err := usecase.CollectEvidence(context.Background(), plan, provider)
	if err != nil {
		t.Fatal(err)
	}
	if got := collectedPaths(bundle); !reflect.DeepEqual(got, []string{"graph.go"}) {
		t.Fatalf("regressão na coleta legada: %v", got)
	}
}

// M01: success must mean the local surface was not silently cut at the old limit.
func TestStep0072GateLocalTreeIncludesEntryAfterOldTwentyThousandCut(t *testing.T) {
	root := t.TempDir()
	const distractors = 20_001
	for i := 0; i < distractors; i++ {
		file, err := os.OpenFile(filepath.Join(root, fmt.Sprintf("distractor_%05d.hpp", i)), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "zzzz_span"), []byte("template<class T> class span {};"), 0o600); err != nil {
		t.Fatal(err)
	}
	tree, err := localinfra.NewClient().FetchTree(context.Background(), domain.RepositoryCandidate{
		Owner: "local", Repo: "large", CommitID: "immutable-13", URL: "file://" + root,
	})
	if err != nil {
		t.Fatal(err)
	}
	seen := false
	for _, node := range tree {
		if node.Path == "zzzz_span" && node.Type == "file" {
			seen = true
			break
		}
	}
	if !seen || len(tree) < distractors+1 {
		t.Fatalf("FetchTree devolveu sucesso truncado: entries=%d final=%v", len(tree), seen)
	}
}

func TestStep0072GatePublicArraysRemainNonNull(t *testing.T) {
	bundle := domain.EvidenceBundle{
		Plan: domain.EvidenceExtractionPlan{RequirementQueries: []domain.EvidenceRequirementQuery{}},
		RequirementCoverage: []domain.RequirementCollectionCoverage{{
			CandidatePaths: []string{}, ExaminedPaths: []string{}, SkippedPaths: []string{},
		}},
	}
	encoded, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"\"requirement_queries\":null", "\"requirement_coverage\":null", "\"candidate_paths\":null", "\"examined_paths\":null", "\"skipped_paths\":null"} {
		if strings.Contains(string(encoded), forbidden) {
			t.Fatalf("array público serializado como null: %s", encoded)
		}
	}
}
