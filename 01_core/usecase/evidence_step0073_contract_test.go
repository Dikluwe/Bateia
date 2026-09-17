/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0073_contract_test.md
 * @layer L1
 * @updated 2026-09-16
 */

package usecase_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func step0073SourceFile(lineCount int, replacements map[int]string) string {
	lines := make([]string, lineCount)
	for index := range lines {
		lines[index] = fmt.Sprintf("// neutral filler %d", index+1)
	}
	for line, value := range replacements {
		lines[line-1] = value
	}
	return strings.Join(lines, "\n")
}

func step0073Rule(id, capability, target string, strongKind domain.EvidenceKind, strong, supporting []string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0073-contract",
		TargetHints: []string{target}, CapabilityTerms: []string{capability},
		StrongEvidence:     []domain.CapabilityPattern{{Kind: strongKind, Patterns: append([]string(nil), strong...)}},
		SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: append([]string(nil), supporting...)}},
	}
}

func step0073Catalog(rules ...domain.CapabilityRule) domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "cpp-0073", Rules: rules}
}

func step0073Bundle(items ...domain.CollectedEvidence) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "local", Repo: "libstdcxx", Commit: "immutable-13"}
	for index := range items {
		items[index].Ref = ref
		sum := sha256.Sum256([]byte(items[index].Content))
		items[index].ContentHash = hex.EncodeToString(sum[:])
	}
	return domain.EvidenceBundle{
		SchemaVersion: "1", Repository: ref, Complete: true, TreeComplete: true,
		Items: items, Targets: []domain.CollectionTargetResult{}, RequirementCoverage: []domain.RequirementCollectionCoverage{},
	}
}

func step0073Assessment(t *testing.T, idea domain.IdeaSpec, bundle domain.EvidenceBundle, catalog domain.CapabilityRuleCatalog) domain.CandidateAssessment {
	t.Helper()
	report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{bundle}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Assessments) != 1 {
		t.Fatalf("assessment ausente: %+v", report)
	}
	return report.Assessments[0]
}

func step0073CapabilitySource(t *testing.T, assessment domain.CandidateAssessment, capability, pattern string) domain.EvidenceSource {
	t.Helper()
	for _, evidence := range assessment.Capabilities {
		if evidence.Capability != capability {
			continue
		}
		for _, source := range evidence.Sources {
			if source.MatchedPattern == pattern {
				return source
			}
		}
		t.Fatalf("padrão %q não citado para %q: %+v", pattern, capability, evidence)
	}
	t.Fatalf("capacidade %q ausente: %+v", capability, assessment.Capabilities)
	return domain.EvidenceSource{}
}

func TestStep0073CitesSubstantialCppOccurrenceInsteadOfFileComment(t *testing.T) {
	stringHeader := step0073SourceFile(410, map[int]string{
		25:  "/** @file basic_string.h */",
		100: "using size_type = unsigned long;",
		407: "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {",
		408: "return __off > size() - __pos ? size() - __pos : __off;",
		409: "}",
	})
	spanHeader := step0073SourceFile(421, map[int]string{
		25:  "/** @file span */",
		100: "inline constexpr size_t dynamic_extent = size_t(-1);",
		388: "constexpr auto subspan(size_type __offset, size_type __count = dynamic_extent) const noexcept {",
		389: "return span(data() + __offset, __count);",
		390: "}",
		420: "constexpr auto subspan() const noexcept { return span(data(), size()); }",
	})
	idea := domain.IdeaSpec{Problem: "cpp intrafile evidence", Capabilities: []string{
		"bounded string offset", "contiguous subspan", "remaining string size",
	}}
	catalog := step0073Catalog(
		step0073Rule("cpp-limit", "bounded string offset", "basic_string.h", domain.EvidenceSymbol, []string{"_M_limit"}, []string{"size_type"}),
		step0073Rule("cpp-subspan", "contiguous subspan", "span", domain.EvidenceSymbol, []string{"subspan"}, []string{"dynamic_extent"}),
		step0073Rule("cpp-remaining", "remaining string size", "basic_string.h", domain.EvidenceFile, []string{"size() - __pos"}, []string{"size_type"}),
	)
	bundle := step0073Bundle(
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: stringHeader},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "span", Content: spanHeader},
	)
	assessment := step0073Assessment(t, idea, bundle, catalog)

	cases := []struct {
		capability, pattern, path, snippet string
		line                               int
	}{
		{"bounded string offset", "_M_limit", "bits/basic_string.h", "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {", 407},
		{"contiguous subspan", "subspan", "span", "constexpr auto subspan(size_type __offset, size_type __count = dynamic_extent) const noexcept {", 388},
		{"remaining string size", "size() - __pos", "bits/basic_string.h", "return __off > size() - __pos ? size() - __pos : __off;", 408},
	}
	for _, tc := range cases {
		source := step0073CapabilitySource(t, assessment, tc.capability, tc.pattern)
		if source.MatchClass != domain.EvidenceMatchStrongContent || source.Path != tc.path ||
			source.Line != tc.line || source.Snippet != tc.snippet {
			t.Errorf("localização inconsistente para %q: %+v", tc.capability, source)
		}
		if strings.Contains(source.Snippet, "@file") || source.MatchClass == domain.EvidenceMatchPathOnly {
			t.Errorf("recuperação foi promovida a substância: %+v", source)
		}
	}

	for _, part := range assessment.ReusableParts {
		for _, source := range part.Evidence {
			if source.MatchClass == domain.EvidenceMatchPathOnly || strings.Contains(source.Snippet, "@file") {
				t.Fatalf("componente reutilizável contém path/doc falso: %+v", part)
			}
		}
	}
}

func TestStep0073RetrievalTermsDoNotBecomeProofPatterns(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp intrafile evidence", Capabilities: []string{"bounded string offset"}}
	rule := step0073Rule("cpp-limit", "bounded string offset", "basic_string.h", domain.EvidenceSymbol,
		[]string{"_M_limit", "size() - __pos"}, []string{"size_type"})
	screening := domain.ScreeningReport{Decisions: []domain.ScreeningDecision{{Repository: "local/libstdcxx@immutable-13", Decision: "ADVANCE"}}}
	plans, err := usecase.PlanEvidenceForIdea(idea, screening, step0073Catalog(rule))
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || len(plans[0].RequirementQueries) != 1 {
		t.Fatalf("plano dirigido ausente: %+v", plans)
	}
	if got, want := plans[0].RequirementQueries[0].Terms, []string{"basic_string.h"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("padrões probatórios vazaram para recuperação: got=%v want=%v", got, want)
	}
}

func TestStep0073CommentsStringsSubstringsAndBasenameAreNotStrongEvidence(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp intrafile evidence", Capabilities: []string{"bounded string offset"}}
	rule := step0073Rule("cpp-limit", "bounded string offset", "basic_string.h", domain.EvidenceSymbol,
		[]string{"_M_limit"}, []string{"size_type"})
	content := strings.Join([]string{
		"/** @file basic_string.h */",
		"// constexpr size_type _M_limit(size_type, size_type) const;",
		"const char* diagnostic = \"_M_limit(size_type, size_type)\";",
		"constexpr size_type _M_limited(size_type value) const noexcept;",
	}, "\n")
	assessment := step0073Assessment(t, idea,
		step0073Bundle(domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: content}),
		step0073Catalog(rule))
	got := assessment.Capabilities[0]
	if got.Status != domain.EvidenceInconclusive || got.ReasonCode != "LEXICAL_MATCH_OPAQUE" || len(got.Sources) != 0 {
		t.Fatalf("ocorrência não probatória ganhou autoridade ou virou ausência: %+v", got)
	}
	if len(assessment.ReusableParts) != 0 {
		t.Fatalf("basename/comentário gerou componente: %+v", assessment.ReusableParts)
	}
}

func TestStep0073AssessmentIsCanonicalAndInputsRemainImmutable(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp intrafile evidence", Capabilities: []string{"bounded string offset", "contiguous subspan"}}
	rules := []domain.CapabilityRule{
		step0073Rule("cpp-limit", "bounded string offset", "basic_string.h", domain.EvidenceSymbol, []string{"size() - __pos", "_M_limit"}, []string{"size_type"}),
		step0073Rule("cpp-subspan", "contiguous subspan", "span", domain.EvidenceSymbol, []string{"subspan"}, []string{"dynamic_extent"}),
	}
	items := []domain.CollectedEvidence{
		{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {\nreturn size() - __pos;\n}"},
		{Kind: domain.EvidenceFile, Path: "span", Content: "constexpr auto subspan(size_type offset) const noexcept;\nconstexpr auto subspan() const noexcept;"},
	}
	originalIdea, _ := json.Marshal(idea)
	originalRules, _ := json.Marshal(rules)
	originalItems, _ := json.Marshal(items)
	var canonical []byte
	for iteration := 0; iteration < 20; iteration++ {
		encodedRules, _ := json.Marshal(rules)
		var permutedRules []domain.CapabilityRule
		if err := json.Unmarshal(encodedRules, &permutedRules); err != nil {
			t.Fatal(err)
		}
		permutedItems := append([]domain.CollectedEvidence(nil), items...)
		if iteration%2 == 1 {
			permutedRules[0], permutedRules[1] = permutedRules[1], permutedRules[0]
			permutedItems[0], permutedItems[1] = permutedItems[1], permutedItems[0]
			patterns := permutedRules[1].StrongEvidence[0].Patterns
			patterns[0], patterns[1] = patterns[1], patterns[0]
		}
		report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{step0073Bundle(permutedItems...)}, step0073Catalog(permutedRules...))
		if err != nil {
			t.Fatal(err)
		}
		observed, _ := json.Marshal(report)
		if iteration == 0 {
			canonical = observed
		} else if string(observed) != string(canonical) {
			t.Fatalf("permutação %d alterou assessment\nwant=%s\ngot=%s", iteration, canonical, observed)
		}
	}
	afterIdea, _ := json.Marshal(idea)
	afterRules, _ := json.Marshal(rules)
	afterItems, _ := json.Marshal(items)
	if string(afterIdea) != string(originalIdea) || string(afterRules) != string(originalRules) || string(afterItems) != string(originalItems) {
		t.Fatal("avaliação modificou ideia, catálogo ou itens")
	}
}
