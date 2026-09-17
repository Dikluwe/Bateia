/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0073_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-16
 */

package step0073_gate_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func oracle0073File(lineCount int, replacements map[int]string) string {
	lines := make([]string, lineCount)
	for index := range lines {
		lines[index] = fmt.Sprintf("// filler %03d", index+1)
	}
	for line, value := range replacements {
		lines[line-1] = value
	}
	return strings.Join(lines, "\n")
}

func oracle0073Rule(id, capability, target string, kind domain.EvidenceKind, strong, supporting []string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: id, Basis: domain.EvidenceRuleDeclared, Family: "cpp-api", Provenance: "step-0073-adversary",
		TargetHints: []string{target}, CapabilityTerms: []string{capability},
		StrongEvidence:     []domain.CapabilityPattern{{Kind: kind, Patterns: append([]string(nil), strong...)}},
		SupportingEvidence: []domain.CapabilityPattern{{Kind: domain.EvidenceFile, Patterns: append([]string(nil), supporting...)}},
	}
}

func oracle0073Catalog(rules ...domain.CapabilityRule) domain.CapabilityRuleCatalog {
	return domain.CapabilityRuleCatalog{SchemaVersion: "2", Version: "step-0073-adversary", Rules: rules}
}

func oracle0073Bundle(complete bool, items ...domain.CollectedEvidence) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "local", Repo: "libstdcxx", Commit: "immutable-13"}
	for index := range items {
		items[index].Ref = ref
		sum := sha256.Sum256([]byte(items[index].Content))
		items[index].ContentHash = hex.EncodeToString(sum[:])
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Complete: complete, TreeComplete: complete, Items: items}
}

func oracle0073Assess(t *testing.T, idea domain.IdeaSpec, bundle domain.EvidenceBundle, catalog domain.CapabilityRuleCatalog) domain.CandidateAssessment {
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

func oracle0073Capability(t *testing.T, assessment domain.CandidateAssessment, capability string) domain.CapabilityEvidence {
	t.Helper()
	for _, observed := range assessment.Capabilities {
		if observed.Capability == capability {
			return observed
		}
	}
	t.Fatalf("capacidade %q ausente: %+v", capability, assessment.Capabilities)
	return domain.CapabilityEvidence{}
}

func oracle0073PatternSources(observed domain.CapabilityEvidence, pattern string) []domain.EvidenceSource {
	var sources []domain.EvidenceSource
	for _, source := range observed.Sources {
		if source.MatchedPattern == pattern {
			sources = append(sources, source)
		}
	}
	return sources
}

// M01, M02, M04, M05, M06, M08 and M09: the cited occurrence, not the first
// retrievable line, is the indivisible public evidence unit.
func TestStep0073GateCitesExactCppDeclarationAndExpression(t *testing.T) {
	stringHeader := oracle0073File(410, map[int]string{
		25:  "/** @file basic_string.h */",
		100: "// _M_limit(size_type __pos, size_type __off) is documented here",
		101: "const char* diagnostic = \"_M_limit(size_type, size_type)\";",
		102: "constexpr size_type _M_limited(size_type value) const noexcept;",
		407: "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {",
		408: "return __off > size() - __pos ? size() - __pos : __off;",
		409: "}",
	})
	spanHeader := oracle0073File(421, map[int]string{
		25:  "/** @file span */",
		80:  "// subspan(offset) is not the implementation",
		81:  "const char* diagnostic = \"subspan(offset)\";",
		82:  "constexpr auto subspanner(size_type offset) const noexcept;",
		388: "constexpr auto subspan(size_type __offset, size_type __count = dynamic_extent) const noexcept {",
		389: "return span(data() + __offset, __count);",
		420: "constexpr auto subspan() const noexcept { return span(data(), size()); }",
	})
	idea := domain.IdeaSpec{Problem: "cpp exact evidence", Capabilities: []string{"bounded offset", "contiguous subview", "remaining size"}}
	catalog := oracle0073Catalog(
		oracle0073Rule("limit", "bounded offset", "basic_string.h", domain.EvidenceSymbol, []string{"_M_limit"}, []string{"size_type"}),
		oracle0073Rule("subspan", "contiguous subview", "span", domain.EvidenceSymbol, []string{"subspan"}, []string{"dynamic_extent"}),
		oracle0073Rule("remaining", "remaining size", "basic_string.h", domain.EvidenceFile, []string{"size() - __pos"}, []string{"size_type"}),
	)
	assessment := oracle0073Assess(t, idea, oracle0073Bundle(true,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: stringHeader},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "span", Content: spanHeader}), catalog)

	cases := []struct {
		capability, pattern, path, snippet string
		line                               int
	}{
		{"bounded offset", "_M_limit", "bits/basic_string.h", "constexpr size_type _M_limit(size_type __pos, size_type __off) const noexcept {", 407},
		{"contiguous subview", "subspan", "span", "constexpr auto subspan(size_type __offset, size_type __count = dynamic_extent) const noexcept {", 388},
		{"remaining size", "size() - __pos", "bits/basic_string.h", "return __off > size() - __pos ? size() - __pos : __off;", 408},
	}
	for _, tc := range cases {
		observed := oracle0073Capability(t, assessment, tc.capability)
		sources := oracle0073PatternSources(observed, tc.pattern)
		if len(sources) == 0 {
			t.Fatalf("%s: padrão %q não citado: %+v", tc.capability, tc.pattern, observed)
		}
		source := sources[0]
		if source.Path != tc.path || source.Line != tc.line || source.Snippet != tc.snippet ||
			source.MatchClass != domain.EvidenceMatchStrongContent || source.MatchedPattern != tc.pattern {
			t.Errorf("%s: ocorrência inconsistente: %+v", tc.capability, source)
		}
		if strings.Contains(source.Snippet, "@file") || source.MatchClass == domain.EvidenceMatchPathOnly {
			t.Errorf("%s: recuperação ganhou autoridade: %+v", tc.capability, source)
		}
	}
}

// M03, M07, M08, M09 and M13: retrieval, opaque syntax and lexical decoys
// cannot manufacture a source or an absence claim.
func TestStep0073GateRejectsPathCommentsStringsSubstringsAndOpaqueSyntax(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp opaque syntax", Capabilities: []string{"bounded offset"}}
	rule := oracle0073Rule("limit", "bounded offset", "basic_string.h", domain.EvidenceSymbol, []string{"_M_limit"}, []string{"size_type"})
	content := strings.Join([]string{
		"/** @file basic_string.h */",
		"// constexpr size_type _M_limit(size_type, size_type) const;",
		"/* _M_limit(size_type, size_type) */",
		"const char* diagnostic = \"_M_limit(size_type, size_type)\";",
		"char marker = '_'; // _M_limit(1, 2)",
		"constexpr size_type _M_limited(size_type value) const noexcept;",
		"#define LIMIT(name) _M_##name",
		"constexpr size_type _M_",
		"limit(size_type value) const noexcept;",
	}, "\n")
	assessment := oracle0073Assess(t, idea,
		oracle0073Bundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: content}),
		oracle0073Catalog(rule))
	observed := oracle0073Capability(t, assessment, "bounded offset")
	if observed.Status != domain.EvidenceInconclusive || observed.ReasonCode != "LEXICAL_MATCH_OPAQUE" || len(observed.Sources) != 0 {
		t.Fatalf("sintaxe opaca/decoys ganharam autoridade ou ausência: %+v", observed)
	}
	if len(assessment.ReusableParts) != 0 {
		t.Fatalf("path/comentário/literal gerou componente: %+v", assessment.ReusableParts)
	}
	for _, source := range observed.Sources {
		if source.MatchClass == domain.EvidenceMatchPathOnly {
			t.Fatalf("PATH_ONLY emitido em sources: %+v", source)
		}
	}
}

// M04 plus the call/declaration edge: both supported lexical forms are located;
// ordering remains the public class/path/line/pattern order, without semantic invention.
func TestStep0073GateRecognizesCppMemberAndCallWithoutLanguageKeyword(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "cpp lexical slice", Capabilities: []string{"member declaration", "member call"}}
	catalog := oracle0073Catalog(
		oracle0073Rule("decl", "member declaration", "widget.hpp", domain.EvidenceSymbol, []string{"reserve_exact"}, nil),
		oracle0073Rule("call", "member call", "widget.hpp", domain.EvidenceSymbol, []string{"commit_exact"}, nil),
	)
	content := "constexpr size_type reserve_exact(size_type n) const noexcept;\nstate.commit_exact(value);"
	assessment := oracle0073Assess(t, idea,
		oracle0073Bundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "widget.hpp", Content: content}), catalog)
	cases := []struct {
		capability, pattern, snippet string
		line                         int
	}{
		{"member declaration", "reserve_exact", "constexpr size_type reserve_exact(size_type n) const noexcept;", 1},
		{"member call", "commit_exact", "state.commit_exact(value);", 2},
	}
	for _, tc := range cases {
		sources := oracle0073PatternSources(oracle0073Capability(t, assessment, tc.capability), tc.pattern)
		if len(sources) != 1 || sources[0].Line != tc.line || sources[0].Snippet != tc.snippet || sources[0].MatchClass != domain.EvidenceMatchStrongContent {
			t.Fatalf("forma C++ não localizada para %s: %+v", tc.capability, sources)
		}
	}
}

// M10 and M12: strong content wins over earlier support and distinct occurrences
// remain distinct, sorted by line, while exact duplicates are absent.
func TestStep0073GatePrecedenceAndDistinctOccurrencesAreCanonical(t *testing.T) {
	idea := domain.IdeaSpec{Problem: "precedence", Capabilities: []string{"contiguous subview"}}
	rule := oracle0073Rule("subspan", "contiguous subview", "span", domain.EvidenceSymbol,
		[]string{"subspan"}, []string{"dynamic_extent"})
	content := strings.Join([]string{
		"inline constexpr size_t dynamic_extent = size_t(-1);",
		"constexpr auto subspan(size_type offset) const noexcept;",
		"constexpr auto subspan() const noexcept;",
	}, "\n")
	assessment := oracle0073Assess(t, idea,
		oracle0073Bundle(true, domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "span", Content: content}),
		oracle0073Catalog(rule))
	observed := oracle0073Capability(t, assessment, "contiguous subview")
	if len(observed.Sources) < 2 || observed.Sources[0].MatchClass != domain.EvidenceMatchStrongContent || observed.Sources[0].Line != 2 {
		t.Fatalf("forte não precede apoio/linha principal incorreta: %+v", observed.Sources)
	}
	strong := oracle0073PatternSources(observed, "subspan")
	if len(strong) != 2 || strong[0].Line != 2 || strong[1].Line != 3 {
		t.Fatalf("ocorrências legítimas foram perdidas ou desordenadas: %+v", strong)
	}
	for index, source := range observed.Sources {
		if source.MatchClass == domain.EvidenceMatchPathOnly {
			t.Fatalf("PATH_ONLY emitido na posição %d: %+v", index, source)
		}
		for previous := 0; previous < index; previous++ {
			if reflect.DeepEqual(observed.Sources[previous], source) {
				t.Fatalf("duplicata idêntica preservada: %+v", source)
			}
		}
	}
	for _, part := range assessment.ReusableParts {
		if len(part.Evidence) == 0 || part.Evidence[0].MatchClass != domain.EvidenceMatchStrongContent || part.Evidence[0].Line != 2 {
			t.Fatalf("componente existente não usa a primeira ocorrência forte: %+v", part)
		}
	}
}

// Opaque controls: the locator must not erase catalog/hypothesis/collection causes.
func TestStep0073GatePreservesEarlierUnknownCauses(t *testing.T) {
	declared := oracle0073Rule("decl", "declared capability", "decl.hpp", domain.EvidenceSymbol, []string{"declared_symbol"}, nil)
	hypothesis := oracle0073Rule("hyp", "hypothesis capability", "hyp.hpp", domain.EvidenceSymbol, []string{"hypothesis_symbol"}, nil)
	hypothesis.Basis = domain.EvidenceRuleHypothesis
	idea := domain.IdeaSpec{Problem: "unknown causes", Capabilities: []string{"declared capability", "hypothesis capability", "catalog gap"}}
	bundle := oracle0073Bundle(false,
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "decl.hpp", Content: "// collection ended before symbol"},
		domain.CollectedEvidence{Kind: domain.EvidenceFile, Path: "hyp.hpp", Content: "constexpr auto hypothesis_symbol() noexcept;"})
	assessment := oracle0073Assess(t, idea, bundle, oracle0073Catalog(declared, hypothesis))
	wants := map[string]string{
		"declared capability":   "COLLECTION_INCOMPLETE",
		"hypothesis capability": "HYPOTHESIS_ONLY",
		"catalog gap":           "CATALOG_GAP",
	}
	for capability, reason := range wants {
		observed := oracle0073Capability(t, assessment, capability)
		if observed.ReasonCode != reason {
			t.Errorf("%s: causa Unknown perdida: %+v", capability, observed)
		}
		if observed.Status == domain.EvidenceProven || observed.Status == domain.EvidenceAbsent ||
			(capability == "catalog gap" && observed.Status == domain.EvidencePartial) {
			t.Errorf("%s: status proibido sob Unknown: %+v", capability, observed)
		}
		if capability == "catalog gap" && len(observed.Sources) != 0 {
			t.Errorf("%s: fonte indevida sem catálogo: %+v", capability, observed.Sources)
		}
	}
	for _, part := range assessment.ReusableParts {
		for _, capability := range part.Capabilities {
			if capability == "hypothesis capability" || capability == "catalog gap" {
				t.Errorf("Unknown gerou componente reutilizável: %+v", part)
			}
		}
	}
}

// M11 and M14: twenty permutations are byte-identical and caller-owned values
// remain byte-identical.
func TestStep0073GateAssessmentIsDeterministicAndInputsImmutable(t *testing.T) {
	baseIdea := domain.IdeaSpec{Problem: "deterministic cpp locator", Capabilities: []string{"bounded offset", "contiguous subview"}}
	baseRules := []domain.CapabilityRule{
		oracle0073Rule("limit", "bounded offset", "basic_string.h", domain.EvidenceSymbol, []string{"size() - __pos", "_M_limit"}, []string{"size_type"}),
		oracle0073Rule("subspan", "contiguous subview", "span", domain.EvidenceSymbol, []string{"subspan"}, []string{"dynamic_extent"}),
	}
	baseItems := []domain.CollectedEvidence{
		{Kind: domain.EvidenceFile, Path: "bits/basic_string.h", Content: "constexpr size_type _M_limit(size_type p, size_type o) const noexcept;\nreturn size() - __pos;"},
		{Kind: domain.EvidenceFile, Path: "span", Content: "constexpr auto subspan(size_type offset) const noexcept;\nconstexpr auto subspan() const noexcept;"},
	}
	ideaBefore, _ := json.Marshal(baseIdea)
	rulesBefore, _ := json.Marshal(baseRules)
	itemsBefore, _ := json.Marshal(baseItems)
	var canonical []byte
	for iteration := 0; iteration < 20; iteration++ {
		ideaJSON, _ := json.Marshal(baseIdea)
		var idea domain.IdeaSpec
		_ = json.Unmarshal(ideaJSON, &idea)
		rulesJSON, _ := json.Marshal(baseRules)
		var rules []domain.CapabilityRule
		_ = json.Unmarshal(rulesJSON, &rules)
		itemsJSON, _ := json.Marshal(baseItems)
		var items []domain.CollectedEvidence
		_ = json.Unmarshal(itemsJSON, &items)
		if iteration%2 == 1 {
			sort.Sort(sort.Reverse(sort.StringSlice(idea.Capabilities)))
			rules[0], rules[1] = rules[1], rules[0]
			items[0], items[1] = items[1], items[0]
			for index := range rules {
				for left, right := 0, len(rules[index].StrongEvidence[0].Patterns)-1; left < right; left, right = left+1, right-1 {
					rules[index].StrongEvidence[0].Patterns[left], rules[index].StrongEvidence[0].Patterns[right] = rules[index].StrongEvidence[0].Patterns[right], rules[index].StrongEvidence[0].Patterns[left]
				}
			}
		}
		report, err := usecase.AssessEvidence(idea, domain.ScreeningReport{}, []domain.EvidenceBundle{oracle0073Bundle(true, items...)}, oracle0073Catalog(rules...))
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
	ideaAfter, _ := json.Marshal(baseIdea)
	rulesAfter, _ := json.Marshal(baseRules)
	itemsAfter, _ := json.Marshal(baseItems)
	if string(ideaBefore) != string(ideaAfter) || string(rulesBefore) != string(rulesAfter) || string(itemsBefore) != string(itemsAfter) {
		t.Fatal("avaliação modificou ideia, catálogo ou bundle do chamador")
	}
}
