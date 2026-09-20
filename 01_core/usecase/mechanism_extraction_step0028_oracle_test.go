//go:build step0028_candidate

/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/mechanism_extraction_step0028_oracle_test.md
 * @layer L1
 * @updated 2026-09-02
 */

package usecase_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type step0028FixtureRef struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Commit string `json:"commit"`
}

type step0028BundleItem struct {
	Kind        domain.EvidenceKind `json:"kind"`
	Path        string              `json:"path"`
	Content     string              `json:"content"`
	ContentHash string              `json:"content_hash"`
}

type step0028EvaluatedSource struct {
	ParentStatus domain.EvidenceStatus `json:"parent_status"`
	Owner        string                `json:"owner"`
	Repo         string                `json:"repo"`
	Commit       string                `json:"commit"`
	Kind         domain.EvidenceKind   `json:"kind"`
	Path         string                `json:"path"`
	Symbol       string                `json:"symbol"`
	Line         int                   `json:"line"`
	RuleID       string                `json:"rule_id"`
	ContentHash  string                `json:"content_hash"`
	Snippet      string                `json:"snippet"`
}

type step0028ExtractionCase struct {
	Repository       step0028FixtureRef        `json:"repository"`
	BaselineClasses  []string                  `json:"baseline_classes"`
	BundleItems      []step0028BundleItem      `json:"bundle_items"`
	EvaluatedSources []step0028EvaluatedSource `json:"evaluated_sources"`
}

type step0028ExtractionFixture struct {
	SchemaVersion     string                            `json:"schema_version"`
	RuleCatalogSHA256 string                            `json:"rule_catalog_sha256"`
	Cases             map[string]step0028ExtractionCase `json:"cases"`
}

func step0028FixturePath(parts ...string) string {
	all := append([]string{"..", "..", "lab", "bancada", "history", "benchmarks", "0028", "oracle", "protected", "fixtures"}, parts...)
	return filepath.Join(all...)
}

func step0028LoadExtractionFixture(t *testing.T) step0028ExtractionFixture {
	t.Helper()
	content, err := os.ReadFile(step0028FixturePath("extraction.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture step0028ExtractionFixture
	if err := json.Unmarshal(content, &fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestStep0028OracleProtectedExtractionFixturesHaveExactHashesAndBindings(t *testing.T) {
	fixture := step0028LoadExtractionFixture(t)
	if fixture.SchemaVersion != "tekt.step0028.extraction-fixtures.v1" || fixture.RuleCatalogSHA256 != "a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5" {
		t.Fatalf("fixture binding changed: %+v", fixture)
	}
	if len(fixture.Cases) != 5 {
		t.Fatalf("extraction case count = %d, want 5", len(fixture.Cases))
	}
	for name, testCase := range fixture.Cases {
		for _, item := range testCase.BundleItems {
			digest := sha256.Sum256([]byte(item.Content))
			if got := "sha256:" + hex.EncodeToString(digest[:]); got != item.ContentHash {
				t.Fatalf("%s %s content hash = %s, want %s", name, item.Path, got, item.ContentHash)
			}
		}
	}

	cross := fixture.Cases["cross_repository_witness"]
	if len(cross.EvaluatedSources) != 2 || cross.EvaluatedSources[0].Repo == cross.EvaluatedSources[1].Repo {
		t.Fatal("cross-repository witness control lost its independent identities")
	}
	docs := fixture.Cases["documentation_only"]
	if len(docs.BundleItems) != 1 || docs.BundleItems[0].Kind != domain.EvidenceDocumentation {
		t.Fatal("documentation-only control gained a non-document witness")
	}
}

func step0028DomainInput(testCase step0028ExtractionCase, catalog []byte, catalogHash string) domain.MechanismExtractionInput {
	defaultRef := domain.CommitRef{Owner: testCase.Repository.Owner, Repo: testCase.Repository.Repo, Commit: testCase.Repository.Commit}
	items := make([]domain.CollectedEvidence, 0, len(testCase.BundleItems))
	for _, item := range testCase.BundleItems {
		items = append(items, domain.CollectedEvidence{Kind: item.Kind, Ref: defaultRef, Path: item.Path, Content: item.Content, ContentHash: item.ContentHash})
	}
	byRef := map[string][]step0028EvaluatedSource{}
	for _, source := range testCase.EvaluatedSources {
		owner, repo, commit := source.Owner, source.Repo, source.Commit
		if owner == "" {
			owner = defaultRef.Owner
		}
		if repo == "" {
			repo = defaultRef.Repo
		}
		if commit == "" {
			commit = defaultRef.Commit
		}
		key := owner + "/" + repo + "@" + commit
		source.Owner, source.Repo, source.Commit = owner, repo, commit
		byRef[key] = append(byRef[key], source)
	}
	keys := make([]string, 0, len(byRef))
	for key := range byRef {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	assessments := make([]domain.CandidateAssessment, 0, len(keys))
	for _, key := range keys {
		sources := byRef[key]
		ref := domain.CommitRef{Owner: sources[0].Owner, Repo: sources[0].Repo, Commit: sources[0].Commit}
		capabilities := make([]domain.CapabilityEvidence, 0, len(sources))
		for index, source := range sources {
			capabilities = append(capabilities, domain.CapabilityEvidence{
				Capability: "oracle-capability-" + string(rune('a'+index)), Status: source.ParentStatus, Confidence: 1,
				Sources: []domain.EvidenceSource{{Kind: source.Kind, Ref: ref, Path: source.Path, Symbol: source.Symbol, Line: source.Line, Snippet: source.Snippet, RuleID: source.RuleID, ContentHash: source.ContentHash, Strength: 1}},
			})
		}
		assessments = append(assessments, domain.CandidateAssessment{Repository: ref, Capabilities: capabilities})
	}
	baseline := &domain.ConventionalBaseline{BaselineID: "oracle-baseline", IdeaHash: "sha256:" + strings.Repeat("7", 64), MechanismClasses: append([]string(nil), testCase.BaselineClasses...)}
	return domain.MechanismExtractionInput{Bundle: domain.EvidenceBundle{SchemaVersion: "1", Repository: defaultRef, Items: items, Complete: true, TreeComplete: true}, Assessments: assessments, ConventionalBaseline: baseline, RuleCatalogJSON: catalog, RuleCatalogSHA256: catalogHash}
}

// Step0028OracleCheckExtractor is the protected conformance entrypoint. The
// verifier binds a candidate through the public MechanismExtractor interface;
// no concrete constructor or private name is part of this oracle contract.
func Step0028OracleCheckExtractor(t *testing.T, subject domain.MechanismExtractor, catalog []byte) {
	t.Helper()
	fixture := step0028LoadExtractionFixture(t)
	for name, testCase := range fixture.Cases {
		t.Run(name, func(t *testing.T) {
			input := step0028DomainInput(testCase, catalog, fixture.RuleCatalogSHA256)
			first, err := subject.ExtractMechanisms(context.Background(), input)
			if err != nil {
				t.Fatalf("extract: %v", err)
			}
			reordered := input
			reordered.Bundle.Items = append([]domain.CollectedEvidence(nil), input.Bundle.Items...)
			sort.Slice(reordered.Bundle.Items, func(i, j int) bool { return reordered.Bundle.Items[i].Path > reordered.Bundle.Items[j].Path })
			reordered.Assessments = append([]domain.CandidateAssessment(nil), input.Assessments...)
			sort.Slice(reordered.Assessments, func(i, j int) bool {
				return reordered.Assessments[i].Repository.Key() > reordered.Assessments[j].Repository.Key()
			})
			second, err := subject.ExtractMechanisms(context.Background(), reordered)
			if err != nil || !reflect.DeepEqual(first, second) {
				t.Fatalf("reordering changed result: %v\nfirst=%+v\nsecond=%+v", err, first, second)
			}
			encoded, _ := json.Marshal(first)
			for _, forbidden := range []string{"ORACLE_TEST_CANARY", "append-record-v1", "replay-test-v1", "counter-unbounded-v1", `"snippet"`, `"content"`} {
				if strings.Contains(string(encoded), forbidden) {
					t.Fatalf("unsafe raw material %q crossed projection: %s", forbidden, encoded)
				}
			}
			switch name {
			case "append_log_proven":
				if len(first.Mechanisms) != 1 || first.Mechanisms[0].Status != domain.MechanismProven || first.Mechanisms[0].MechanismClass != "APPEND_ONLY_LOG_RECOVERY" || first.Mechanisms[0].Novelty != domain.NoveltyConventional {
					t.Fatalf("append control = %+v", first)
				}
			case "queue_proven_reordered":
				if len(first.Mechanisms) != 1 || first.Mechanisms[0].Status != domain.MechanismProven || first.Mechanisms[0].MechanismClass != "BOUNDED_QUEUE_BACKPRESSURE" {
					t.Fatalf("queue control = %+v", first)
				}
			case "documentation_only":
				if len(first.Mechanisms) > 0 && first.Mechanisms[0].Status == domain.MechanismProven {
					t.Fatalf("documentation-only became PROVEN: %+v", first)
				}
			case "cross_repository_witness":
				for _, mechanism := range first.Mechanisms {
					if mechanism.Status == domain.MechanismProven {
						t.Fatalf("cross-repo witnesses combined: %+v", first)
					}
				}
			case "central_counterevidence":
				if len(first.Mechanisms) != 1 || first.Mechanisms[0].Status != domain.MechanismInconclusive || len(first.Mechanisms[0].CounterEvidence) == 0 {
					t.Fatalf("central counterevidence disappeared: %+v", first)
				}
			}
		})
	}
}
