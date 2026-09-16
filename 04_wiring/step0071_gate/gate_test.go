/**
 * Crystalline Lineage
 * @prompt 00_nucleo/passos/0071-expor-cobertura-do-catalogo-sem-atribuir-lacunas-ao-candidato.md
 * @layer L4
 * @updated 2026-09-16
 */

package step0071_gate_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func gateRule(id string, basis domain.EvidenceRuleBasis, terms ...string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, Basis: basis, Family: "software-capability", Provenance: "step-0071-gate",
		TargetHints: []string{"source symbol"}, CapabilityTerms: terms,
		StrongEvidence:     []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{"FluxEngine"}}},
		SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceDocumentation, Patterns: []string{"quantum flux"}}},
	}
}

func gateCatalog(version string, rules ...domain.CapabilityRule) domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: version, Rules: rules}
}

func gateBundle(content string, complete bool) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "org", Repo: "candidate", Commit: "abc123"}
	items := []domain.CollectedEvidence{}
	if content != "" {
		sum := sha256.Sum256([]byte(content))
		items = append(items, domain.CollectedEvidence{Kind: domain.EvidenceFile, Ref: ref, Path: "flux.go", Content: content, ContentHash: hex.EncodeToString(sum[:])})
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: items, Targets: []domain.CollectionTargetResult{}, Complete: complete, TreeComplete: complete}
}

func gateEntry(t *testing.T, coverage domain.EvidenceCatalogCoverage, requirement string) domain.CatalogCoverageEntry {
	t.Helper()
	for _, entry := range coverage.Entries {
		if entry.Requirement == requirement {
			return entry
		}
	}
	t.Fatalf("entrada ausente para %q: %+v", requirement, coverage)
	return domain.CatalogCoverageEntry{}
}

func gateCapability(t *testing.T, report domain.AssessmentReport, capability string) domain.CapabilityEvidence {
	t.Helper()
	if len(report.Assessments) != 1 {
		t.Fatalf("assessment inesperado: %+v", report)
	}
	for _, evidence := range report.Assessments[0].Capabilities {
		if evidence.Capability == capability {
			return evidence
		}
	}
	t.Fatalf("capacidade ausente: %s", capability)
	return domain.CapabilityEvidence{}
}

func TestStep0071GateP01PublicSchemaAndArrays(t *testing.T) {
	if domain.EvidenceRuleDeclared != "DECLARED" || domain.EvidenceRuleHypothesis != "HYPOTHESIS" ||
		domain.CatalogCoverageCovered != "COVERED" || domain.CatalogCoverageGap != "CATALOG_GAP" || domain.CatalogCoverageHypothesis != "HYPOTHESIS" {
		t.Fatal("valores públicos instáveis")
	}
	report := domain.EvidenceCatalogCoverage{SchemaVersion: "1", RulesVersion: "gate", RulesHash: strings.Repeat("a", 64), Entries: []domain.CatalogCoverageEntry{}}
	b, err := json.Marshal(report)
	if err != nil || !strings.Contains(string(b), `"entries":[]`) {
		t.Fatalf("array vazio deve ser []: err=%v json=%s", err, b)
	}
	var round domain.EvidenceCatalogCoverage
	if err := json.Unmarshal(b, &round); err != nil || round.Entries == nil {
		t.Fatalf("round-trip inválido: err=%v %+v", err, round)
	}
	hyp, _ := json.Marshal(gateRule("hyp", domain.EvidenceRuleHypothesis, "quantum flux"))
	if !strings.Contains(string(hyp), `"basis":"HYPOTHESIS"`) {
		t.Fatalf("basis de hipótese perdida: %s", hyp)
	}
}

func TestStep0071GateP02ExactDeclaredAliasIsCovered(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "routing", Capabilities: []string{"Quantum-Flux"}}
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, gateCatalog("provided", gateRule("flux", domain.EvidenceRuleDeclared, "quantum flux", "flux routing")))
	if err != nil {
		t.Fatal(err)
	}
	entry := gateEntry(t, coverage, "Quantum-Flux")
	if entry.Status != domain.CatalogCoverageCovered || entry.RuleID != "flux" || entry.RuleBasis != domain.EvidenceRuleDeclared {
		t.Fatalf("alias exato não coberto: %+v", entry)
	}
}

func TestStep0071GateO01CatalogGapNeverDescribesCandidate(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "routing", Capabilities: []string{"quantum flux"}}
	catalog := gateCatalog("provided", gateRule("other", domain.EvidenceRuleDeclared, "other capability"))
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		t.Fatal(err)
	}
	entry := gateEntry(t, coverage, "quantum flux")
	if entry.Status != domain.CatalogCoverageGap || entry.ReasonCode != "CATALOG_GAP" {
		t.Fatalf("gap não preservado: %+v", entry)
	}
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{gateBundle("type FluxEngine struct{}", true)}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	got := gateCapability(t, report, "quantum flux")
	if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "CATALOG_GAP" || len(got.Sources) != 0 || got.Confidence != 0 {
		t.Fatalf("gap atribuído ao candidato: %+v", got)
	}
	if len(report.Assessments[0].Gaps) != 0 || len(report.Assessments[0].ReusableParts) != 0 {
		t.Fatalf("gap/componente indevido: %+v", report.Assessments[0])
	}
	plan := usecase.BuildInvestigationPlan(idea, domain.DiscoverySnapshot{}, report.Assessments)
	found := false
	for _, action := range plan.Actions {
		if action.Candidate != "" && action.Classification == "CATALOG_GAP" {
			found = true
		}
		if action.Classification == "NOT_VERIFIED" && action.Candidate != "" {
			t.Fatalf("CATALOG_GAP degradado: %+v", action)
		}
	}
	if !found {
		t.Fatalf("classificação CATALOG_GAP ausente: %+v", plan.Actions)
	}
	coverage2, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		t.Fatal(err)
	}
	b1, _ := json.Marshal(coverage)
	b2, _ := json.Marshal(coverage2)
	if string(b1) != string(b2) {
		t.Fatal("conteúdo do candidato influenciou cobertura")
	}
}

func TestStep0071GateO02UnsafeAndAmbiguousMatchesStayUnknown(t *testing.T) {
	fragment := domain.IdeaSpec{Capabilities: []string{"graph"}}
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(fragment, gateCatalog("fragment", gateRule("dependency", domain.EvidenceRuleDeclared, "dependency graph")))
	if err != nil || gateEntry(t, coverage, "graph").Status != domain.CatalogCoverageGap {
		t.Fatalf("fragmento reverso não permaneceu Unknown: err=%v %+v", err, coverage)
	}
	ambiguous := gateCatalog("ambiguous",
		gateRule("first", domain.EvidenceRuleDeclared, "quantum-flux"),
		gateRule("second", domain.EvidenceRuleDeclared, "quantum flux"),
	)
	coverage, err = usecase.AnalyzeEvidenceCatalogCoverage(domain.IdeaSpec{Capabilities: []string{"quantum flux"}}, ambiguous)
	if err != nil || gateEntry(t, coverage, "quantum flux").Status != domain.CatalogCoverageGap {
		t.Fatalf("primeira regra venceu ambiguidade: err=%v %+v", err, coverage)
	}
}

func TestStep0071GateN02HypothesisHasNoDeclaredAuthority(t *testing.T) {
	idea := domain.IdeaSpec{Capabilities: []string{"quantum flux"}}
	catalog := gateCatalog("hyp", gateRule("hyp", domain.EvidenceRuleHypothesis, "quantum flux"))
	coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, catalog)
	if err != nil {
		t.Fatal(err)
	}
	entry := gateEntry(t, coverage, "quantum flux")
	if entry.Status != domain.CatalogCoverageHypothesis || entry.RuleBasis != domain.EvidenceRuleHypothesis {
		t.Fatalf("hipótese mascarada: %+v", entry)
	}
	for _, bundle := range []domain.EvidenceBundle{gateBundle("type FluxEngine struct{}", true), gateBundle("", true)} {
		report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, catalog)
		if err != nil {
			t.Fatal(err)
		}
		got := gateCapability(t, report, "quantum flux")
		if got.Status == domain.EvidenceProven || got.Status == domain.EvidenceAbsent || got.ReasonCode != "HYPOTHESIS_ONLY" || len(report.Assessments[0].ReusableParts) != 0 {
			t.Fatalf("hipótese ganhou autoridade: %+v", got)
		}
	}
}

func TestStep0071GateO03IncompleteCollectionKeepsItsOwnUnknownCause(t *testing.T) {
	idea := domain.IdeaSpec{Capabilities: []string{"quantum flux"}}
	catalog := gateCatalog("declared", gateRule("flux", domain.EvidenceRuleDeclared, "quantum flux"))
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{gateBundle("", false)}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	got := gateCapability(t, report, "quantum flux")
	if got.Status != domain.EvidenceInconclusive || got.ReasonCode == "CATALOG_GAP" || got.ReasonCode == "HYPOTHESIS_ONLY" {
		t.Fatalf("incompletude mascarada: %+v", got)
	}
}

func TestStep0071GateN01ProvidedCatalogCrossesPipeline(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "routing", Capabilities: []string{"quantum flux"}}
	catalog := gateCatalog("provided-0071", gateRule("flux", domain.EvidenceRuleDeclared, "quantum flux"))
	snapshot := domain.DiscoverySnapshot{SchemaVersion: "1", Idea: idea, Strategy: "literal", Candidates: []domain.DiscoveredCandidate{{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "candidate", CommitID: "abc123", URL: "https://github.com/org/candidate"}}}}
	req := domain.IdeaRunRequest{Idea: idea, Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 5, ExtendMinCoverage: .5, IntegrationMinValue: .5}, Mode: domain.IdeaRunReplay, Strategy: "literal", MaxPerQuery: 5, MaxCandidates: 5}
	result, err := usecase.RunIdeaPipeline(context.Background(), req, usecase.IdeaRunDependencies{ReplaySnapshot: &snapshot, ReplayBundles: []domain.EvidenceBundle{gateBundle("", true)}, EvidenceRules: catalog})
	if err != nil {
		t.Fatal(err)
	}
	if result.CatalogCoverage.RulesVersion != "provided-0071" || result.Assessment.RulesVersion != "provided-0071" || result.Assessment.CatalogCoverage.RulesHash != result.CatalogCoverage.RulesHash {
		t.Fatalf("pipeline recuou ao default: run=%+v assessment=%+v", result.CatalogCoverage, result.Assessment.CatalogCoverage)
	}
	service, err := usecase.NewRecommendationServiceWithEvidenceRules(nil, nil, nil, catalog)
	if err != nil || service == nil {
		t.Fatalf("construtor integrado rejeitou catálogo válido: service=%v err=%v", service, err)
	}
}

func TestStep0071GateD01TwentyPermutationsAreByteStable(t *testing.T) {
	baseRules := []domain.CapabilityRule{
		gateRule("alpha", domain.EvidenceRuleDeclared, "alpha capability", "alpha alias"),
		gateRule("beta", domain.EvidenceRuleDeclared, "beta capability", "beta alias"),
	}
	var wantCoverage, wantAssessment []byte
	for i := 0; i < 20; i++ {
		idea := domain.IdeaSpec{Problem: "alpha capability", Capabilities: []string{"alpha capability", "beta capability"}}
		rules := append([]domain.CapabilityRule(nil), baseRules...)
		for j := range rules {
			rules[j].CapabilityTerms = append([]string(nil), rules[j].CapabilityTerms...)
		}
		bundle := gateBundle("type FluxEngine struct{}", true)
		if i%2 == 1 {
			idea.Capabilities[0], idea.Capabilities[1] = idea.Capabilities[1], idea.Capabilities[0]
			rules[0], rules[1] = rules[1], rules[0]
			for j := range rules {
				rules[j].CapabilityTerms[0], rules[j].CapabilityTerms[1] = rules[j].CapabilityTerms[1], rules[j].CapabilityTerms[0]
			}
		}
		catalog := gateCatalog("stable", rules...)
		beforeIdea, _ := json.Marshal(idea)
		beforeCatalog, _ := json.Marshal(catalog)
		beforeBundle, _ := json.Marshal(bundle)
		coverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, catalog)
		if err != nil {
			t.Fatal(err)
		}
		report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, catalog)
		if err != nil {
			t.Fatal(err)
		}
		gotCoverage, _ := json.Marshal(coverage)
		gotAssessment, _ := json.Marshal(report)
		afterIdea, _ := json.Marshal(idea)
		afterCatalog, _ := json.Marshal(catalog)
		afterBundle, _ := json.Marshal(bundle)
		if string(beforeIdea) != string(afterIdea) || string(beforeCatalog) != string(afterCatalog) || string(beforeBundle) != string(afterBundle) {
			t.Fatalf("input mutado na permutação %d", i)
		}
		if i == 0 {
			wantCoverage, wantAssessment = gotCoverage, gotAssessment
		} else if string(gotCoverage) != string(wantCoverage) || string(gotAssessment) != string(wantAssessment) {
			t.Fatalf("saída não determinística na permutação %d\ncoverage=%s\nassessment=%s", i, gotCoverage, gotAssessment)
		}
	}
}

func TestStep0071GateN03CLICatalogArtifactIsCausal(t *testing.T) {
	// A prova E2E detalhada permanece no contrato CLI congelado. Este gate exige
	// adicionalmente que a fiação declare o flag, o hash causal, o artifact e o resume.
	source, err := os.ReadFile(filepath.Join("..", "..", "02_shell", "cli", "idea_recommend.go"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, witness := range []string{`String("rules"`, `evidence_rules`, `catalog-coverage.json`, `NewRecommendationServiceWithEvidenceRules`, `reusableEvidenceCheckpoint`} {
		if !strings.Contains(text, witness) {
			t.Fatalf("witness causal CLI ausente: %s", witness)
		}
	}
}

func TestStep0071GateN04CoveragePrecedesCollectionStatically(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("..", "..", "01_core", "usecase", "idea_run.go"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	analyze := strings.Index(text, "AnalyzeEvidenceCatalogCoverage(")
	collect := strings.Index(text, "CollectEvidence(")
	if analyze < 0 || collect < 0 || analyze > collect {
		t.Fatalf("cobertura deve ser anexada antes de CollectEvidence: analyze=%d collect=%d", analyze, collect)
	}
}
