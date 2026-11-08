/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_recovery_followup_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type followupRecoveryPort struct {
	budgetRevisionPort
	resolve         bool
	metadataFailure bool
	treeFailure     bool
}

func (p *followupRecoveryPort) FetchRepoMetadata(_ context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.call("metadata", c)
	if p.metadataFailure {
		return c, errors.New("resolution unavailable")
	}
	if p.resolve {
		c.CommitID = strings.Repeat("a", 40)
	}
	return c, nil
}
func (p *followupRecoveryPort) FetchTree(ctx context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	tree, err := p.budgetRevisionPort.FetchTree(ctx, c)
	if p.treeFailure {
		return nil, errors.New("tree unavailable")
	}
	return tree, err
}
func followupExecute(t *testing.T, plan domain.EvidenceRecoveryPlan, original []domain.EvidenceBundle, p *followupRecoveryPort) domain.RecoveryRoundResult {
	t.Helper()
	got, err := ExecuteEvidenceRecovery(context.Background(), domain.IdeaSpec{Problem: "dependency graph", Capabilities: []string{"dependency graph"}}, domain.ScreeningReport{}, original, domain.ReuseDecisionReport{Decision: domain.DecisionUnknown}, plan, p, DefaultEvidenceRules(), domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .7, IntegrationMinValue: .5})
	if err != nil {
		t.Fatal(err)
	}
	return got
}
func TestRecoveryRetiresOnlySuccessfullyResolvedPlaceholder(t *testing.T) {
	for _, commit := range []string{"HEAD", ""} {
		for _, resolved := range []bool{true, false} {
			t.Run(fmt.Sprintf("%q/%v", commit, resolved), func(t *testing.T) {
				ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: commit}
				pinned := ref
				pinned.Commit = strings.Repeat("b", 40)
				oldText := "historical bytes must not acquire another revision"
				digest := sha256.Sum256([]byte(oldText))
				original := []domain.EvidenceBundle{{SchemaVersion: "1", Repository: ref, Failures: []domain.ExecutionFailure{{Stage: "resolution", Message: "old failure"}}, Items: []domain.CollectedEvidence{{Ref: ref, Path: "old.go", Content: oldText, ContentHash: fmt.Sprintf("%x", digest)}}}, {SchemaVersion: "1", Repository: pinned}}
				before, _ := json.Marshal(original)
				plan := budgetPlan(ref)
				plan.Targets[0].Path = "@repository-metadata"
				plan.Targets[0].Kind = domain.EvidenceRepository
				p := &followupRecoveryPort{resolve: resolved, metadataFailure: !resolved}
				got := followupExecute(t, plan, original, p)
				unresolved := false
				for _, f := range got.Diagnosis.Findings {
					unresolved = unresolved || f.Cause == domain.CauseUnresolvedRevision
				}
				if unresolved == resolved {
					t.Fatalf("resolution status stale: resolved=%v diagnosis=%+v", resolved, got.Diagnosis)
				}
				foundPinned, foundOpaque := false, false
				for _, b := range got.Bundles {
					foundPinned = foundPinned || b.Repository == pinned
					if mutableIdeaRunCommit(b.Repository.Commit) {
						foundOpaque = true
					}
					if resolved {
						for _, item := range b.Items {
							if item.Content == oldText || item.Ref != b.Repository {
								t.Fatal("historical content reassigned")
							}
						}
					}
				}
				if !foundPinned || foundOpaque == resolved {
					t.Fatalf("active origins wrong: %+v", got.Bundles)
				}
				after, _ := json.Marshal(original)
				if string(before) != string(after) {
					t.Fatal("original mutated")
				}
			})
		}
	}
	// Reading a fixed revision is not an explicit resolution of the old placeholder.
	opaque := domain.CommitRef{Owner: "o", Repo: "r", Commit: "HEAD"}
	fixed := opaque
	fixed.Commit = strings.Repeat("a", 40)
	got := followupExecute(t, budgetPlan(fixed), []domain.EvidenceBundle{{SchemaVersion: "1", Repository: opaque}}, &followupRecoveryPort{})
	kept := false
	for _, b := range got.Bundles {
		kept = kept || b.Repository == opaque
	}
	if !kept {
		t.Fatal("unattempted placeholder retired")
	}
}
func TestOriginalRecoveryHistoryRequiresAttemptWitness(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	target := domain.ExtractionTarget{Path: "LICENSE", Kind: domain.EvidenceLicense}
	for _, status := range []domain.CollectionTargetStatus{domain.TargetCollected, domain.TargetNotFound, domain.TargetFailed, domain.TargetSkippedBudget, domain.TargetUnsupported, domain.TargetPlanned, ""} {
		t.Run(string(status), func(t *testing.T) {
			b := domain.EvidenceBundle{Repository: ref, Plan: domain.EvidenceExtractionPlan{Targets: []domain.ExtractionTarget{target}}}
			if status != "" {
				b.Targets = []domain.CollectionTargetResult{{Target: target, Status: status}}
			}
			previous := AttemptedOriginalRecoveryTargets([]domain.EvidenceBundle{b})
			d := domain.InconclusiveDiagnosis{SchemaVersion: "1", Findings: []domain.InconclusiveFinding{{ID: "license", Candidate: ref.Key(), Cause: domain.CauseUnknownLicense, RuleID: "GATE-LICENSE-001", Recoverable: true}}}
			plan, err := PlanEvidenceRecovery(d, DefaultRecoveryPolicy(), 1, previous)
			if err != nil {
				t.Fatal(err)
			}
			attempted := status == domain.TargetCollected || status == domain.TargetNotFound || status == domain.TargetFailed
			if len(plan.Targets) != 1 || ((plan.Targets[0].Path == "LICENSE") == attempted) {
				t.Fatalf("incorrect attempted history: status=%q previous=%+v plan=%+v", status, previous, plan)
			}
		})
	}
	b := domain.EvidenceBundle{Repository: ref, Plan: domain.EvidenceExtractionPlan{Targets: []domain.ExtractionTarget{target}}, Items: []domain.CollectedEvidence{{Path: "LICENSE", Kind: domain.EvidenceLicense}}}
	if len(AttemptedOriginalRecoveryTargets([]domain.EvidenceBundle{b})) != 1 {
		t.Fatal("legacy item witness lost")
	}
}
func TestRecoveryCountsActualProviderCalls(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	for _, tc := range []struct {
		name                   string
		derived, fail, limited bool
		want                   int
	}{{"explicit", false, false, false, 3}, {"derived", true, false, false, 4}, {"tree failure", false, true, false, 3}, {"budget", false, false, true, 3}} {
		t.Run(tc.name, func(t *testing.T) {
			plan := budgetPlan(ref)
			if tc.derived {
				plan.Targets[0].Path = ""
				plan.Targets[0].Patterns = []string{"license"}
			}
			if tc.limited {
				other := ref
				other.Repo = "z"
				plan = budgetPlan(ref, other)
				plan.Policy.AdditionalBytes = 400
			}
			p := &followupRecoveryPort{treeFailure: tc.fail}
			p.content = strings.Repeat("x", 512)
			got := followupExecute(t, plan, nil, p)
			if len(p.calls) != tc.want || got.Metrics.ProviderCalls != tc.want {
				t.Fatalf("calls: actual=%d reported=%d want=%d", len(p.calls), got.Metrics.ProviderCalls, tc.want)
			}
		})
	}
}

func TestExecutedRecoveryHistoryUsesOutcomeAndKeepsPlannedIdentity(t *testing.T) {
	ref := evidenceIntegrityPlan().Repository
	for _, derived := range []bool{false, true} {
		for _, status := range []domain.CollectionTargetStatus{domain.TargetCollected, domain.TargetNotFound, domain.TargetFailed, domain.TargetSkippedBudget, domain.TargetUnsupported} {
			t.Run(fmt.Sprintf("%v/%s", derived, status), func(t *testing.T) {
				plan := budgetPlan(ref)
				if derived {
					plan.Targets[0].Path = ""
					plan.Targets[0].Patterns = []string{"license"}
				}
				outcome := domain.EvidenceBundle{Repository: ref, Targets: []domain.CollectionTargetResult{{Target: domain.ExtractionTarget{Path: "LICENSE", Kind: domain.EvidenceLicense, Reason: "recovery 0"}, Status: domain.TargetCollected}, {Target: domain.ExtractionTarget{Path: plan.Targets[0].Path, Kind: domain.EvidenceLicense, Reason: "recovery 0"}, Status: status}}}
				before, _ := json.Marshal(plan)
				got := ExecutedRecoveryTargets(plan, []domain.EvidenceBundle{outcome})
				attempted := status == domain.TargetCollected || status == domain.TargetNotFound || status == domain.TargetFailed
				if attempted {
					if len(got) != 1 {
						t.Fatal("attempted target lost")
					}
					encoded, _ := json.Marshal(got[0])
					want, _ := json.Marshal(plan.Targets[0])
					if string(encoded) != string(want) {
						t.Fatal("planned identity changed")
					}
				} else if len(got) != 0 {
					t.Fatal("omission treated as execution", got)
				}
				after, _ := json.Marshal(plan)
				if string(before) != string(after) {
					t.Fatal("plan mutated")
				}
			})
		}
	}
	if got := ExecutedRecoveryTargets(budgetPlan(ref), nil); len(got) != 0 {
		t.Fatal("missing outcome fabricated history")
	}
}

func TestRecoveryBudgetOmissionDoesNotResolvePlaceholder(t *testing.T) {
	first := domain.CommitRef{Owner: "a", Repo: "first", Commit: strings.Repeat("a", 40)}
	opaque := domain.CommitRef{Owner: "z", Repo: "last", Commit: "HEAD"}
	plan := budgetPlan(first, opaque)
	plan.Policy.AdditionalBytes = 400
	plan.Targets[1].Path = "@repository-metadata"
	plan.Targets[1].Kind = domain.EvidenceRepository
	p := &followupRecoveryPort{resolve: true}
	p.content = strings.Repeat("x", 512)
	got := followupExecute(t, plan, []domain.EvidenceBundle{{SchemaVersion: "1", Repository: opaque}}, p)
	unresolved := false
	for _, f := range got.Diagnosis.Findings {
		unresolved = unresolved || f.Cause == domain.CauseUnresolvedRevision
	}
	if !unresolved || got.Metrics.ProviderCalls != 3 {
		t.Fatalf("omitted target resolved: %+v", got)
	}
	for _, seen := range p.refs {
		if seen.Owner == opaque.Owner {
			t.Fatal("omitted origin reached provider")
		}
	}
}
