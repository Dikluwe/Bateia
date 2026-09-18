/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_step0075_test.md
 * @layer L1
 * @updated 2026-09-17
 */

package usecase

import (
	"reflect"
	"testing"

	"bateia/01_core/domain"
)

func step0075Rule(pattern string) domain.CapabilityRule {
	return domain.CapabilityRule{
		ID: "self-assignment-guard-hashtable", Basis: domain.EvidenceRuleDeclared,
		Family: "cpp-api-design", Provenance: "step-0075-regression",
		TargetHints:     []string{"hashtable.h", "robin_hood.h"},
		CapabilityTerms: []string{"guard copy assignment against self-assignment"},
		StrongEvidence:  []domain.CapabilityPattern{{Kind: domain.EvidenceSymbol, Patterns: []string{pattern}}},
	}
}

func step0075Bundle(content string) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: "local", Repo: "cpp", Commit: "immutable"}
	return domain.EvidenceBundle{
		Repository: ref, Complete: true, TreeComplete: true,
		Items: []domain.CollectedEvidence{{Kind: domain.EvidenceFile, Ref: ref, Path: "bits/hashtable.h", Content: content, ContentHash: contentHash(content)}},
	}
}

func TestStep0075LocatesOperatorRichCodeExpressions(t *testing.T) {
	cases := []struct {
		name, pattern, content, snippet string
	}{
		{"self guard", "&__ht == this", "void assign(const table& __ht) {\n  if (  &__ht   ==   this ) return;\n}", "if ( &__ht == this ) return;"},
		{"short self guard", "&o == this", "if (&o == this) { return *this; }", "if (&o == this) { return *this; }"},
		{"member relink", "__p->_M_nxt =", "__p->_M_nxt = __next;", "__p->_M_nxt = __next;"},
		{"modulo", "__num % __den", "return __num % __den;", "return __num % __den;"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := assessCapability("guard copy assignment against self-assignment", step0075Bundle(tc.content), ptrRule(step0075Rule(tc.pattern)))
			if got.Status != domain.EvidencePartial || len(got.Sources) != 1 {
				t.Fatalf("expressão composta não localizada: %+v", got)
			}
			source := got.Sources[0]
			if source.MatchedPattern != tc.pattern || source.MatchClass != domain.EvidenceMatchStrongContent || source.Snippet != tc.snippet || source.Line < 1 {
				t.Fatalf("fonte não aponta para a expressão exata: %+v", source)
			}
		})
	}
}

func TestStep0075RejectsNonCodeAndAmbiguousCompoundText(t *testing.T) {
	cases := []struct {
		name, pattern, content string
	}{
		{"line comment", "&__ht == this", "// if (&__ht == this) return;"},
		{"block comment", "&__ht == this", "/* if (&__ht == this) return; */"},
		{"string", "&__ht == this", `const char* note = "&__ht == this";`},
		{"identifier substring", "&__ht == this", "if (&__htable == this) return;"},
		{"natural words", "self assignment guard", "self assignment guard"},
		{"multiline", "&__ht == this", "if (&__ht ==\n    this) return;"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := assessCapability("guard copy assignment against self-assignment", step0075Bundle(tc.content), ptrRule(step0075Rule(tc.pattern)))
			if len(got.Sources) != 0 || got.Status == domain.EvidencePartial || got.Status == domain.EvidenceProven {
				t.Fatalf("texto não admissível ganhou autoridade: %+v", got)
			}
		})
	}
}

func TestStep0075ProofExpressionsDoNotLeakIntoPathQueries(t *testing.T) {
	rule := step0075Rule("&__ht == this")
	coverage := domain.EvidenceCatalogCoverage{Entries: []domain.CatalogCoverageEntry{{
		RequirementKind: "capability", Requirement: rule.CapabilityTerms[0], RuleID: rule.ID,
		RuleBasis: domain.EvidenceRuleDeclared, Status: domain.CatalogCoverageCovered,
	}}}
	got := compileEvidenceRequirementQueries(coverage, domain.CapabilityRuleCatalog{SchemaVersion: "2", Rules: []domain.CapabilityRule{rule}})
	if len(got) != 1 || !reflect.DeepEqual(got[0].Terms, []string{"hashtable.h", "robin_hood.h"}) {
		t.Fatalf("padrão probatório vazou para recuperação: %+v", got)
	}
}

func ptrRule(rule domain.CapabilityRule) *domain.CapabilityRule { return &rule }
