/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_recheck_test.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type technicalRecheckProvider struct {
	tree      []domain.DirectoryNode
	extracted *domain.ExtractedRepository
	treeErr   error
	readErr   error
	calls     int
	candidate domain.RepositoryCandidate
	paths     []string
}

func (p *technicalRecheckProvider) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	return nil, errors.New("unexpected manifest")
}
func (p *technicalRecheckProvider) FetchRepoMetadata(context.Context, domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	return domain.RepositoryCandidate{}, errors.New("unexpected metadata")
}
func (p *technicalRecheckProvider) FetchTree(_ context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls++
	p.candidate = c
	return p.tree, p.treeErr
}
func (p *technicalRecheckProvider) ExtractRepoFiles(_ context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.calls++
	p.candidate = c
	p.paths = append([]string(nil), paths...)
	return p.extracted, p.readErr
}

func recheckFixture() (domain.TechnicalSourceRecheckRequest, *technicalRecheckProvider) {
	batch, bundle, ref, _ := verificationFixture()
	ref.Commit = strings.Repeat("a", 40)
	batch.Findings[0].Candidate = ref.Key()
	batch.Findings[0].Source.Revision = ref.Commit
	candidate := domain.RepositoryCandidate{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit}
	return domain.TechnicalSourceRecheckRequest{Batch: batch, Candidates: []domain.RepositoryCandidate{candidate}}, &technicalRecheckProvider{
		tree:      []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}},
		extracted: &domain.ExtractedRepository{Ref: ref, Files: []domain.ExtractedFile{{Path: "LICENSE", Content: bundle.Items[0].Content}}},
	}
}

func TestRecheckTechnicalSourcesReadsPinnedProviderBytes(t *testing.T) {
	req, provider := recheckFixture()
	_, _, _, asOf := verificationFixture()
	report, err := RecheckTechnicalSources(context.Background(), req, provider, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if report.Basis != "PROVIDER_RECHECK" || report.Review.Verifications[0].Status != domain.TechnicalVerified || !report.Bundles[0].Complete {
		t.Fatalf("expected provider verification: %+v", report)
	}
	if provider.candidate.CommitID != req.Candidates[0].CommitID || provider.candidate.URL != "https://github.com/o/r" || !reflect.DeepEqual(provider.paths, []string{"LICENSE"}) {
		t.Fatalf("unpinned provider request: %+v", provider)
	}
	if report.Bundles[0].Items[0].ContentHash != contentHash(provider.extracted.Files[0].Content) {
		t.Fatal("hash not computed from provider bytes")
	}
}

func TestRecheckTechnicalSourcesDoesNotTrustClaimedBytesOrIncompleteReads(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*technicalRecheckProvider)
	}{
		{"different bytes", func(p *technicalRecheckProvider) { p.extracted.Files[0].Content += " changed" }},
		{"wrong revision", func(p *technicalRecheckProvider) { p.extracted.Ref.Commit = strings.Repeat("b", 40) }},
		{"stale read", func(p *technicalRecheckProvider) { p.extracted.Stale = true }},
		{"missing file", func(p *technicalRecheckProvider) { p.extracted.Files = nil }},
		{"duplicate file", func(p *technicalRecheckProvider) { p.extracted.Files = append(p.extracted.Files, p.extracted.Files[0]) }},
		{"unexpected file", func(p *technicalRecheckProvider) {
			p.extracted.Files = append(p.extracted.Files, domain.ExtractedFile{Path: "other"})
		}},
		{"tree unavailable", func(p *technicalRecheckProvider) { p.treeErr = errors.New("offline") }},
		{"read unavailable", func(p *technicalRecheckProvider) { p.readErr = errors.New("offline") }},
		{"nil extraction", func(p *technicalRecheckProvider) { p.extracted = nil }},
		{"admission byte budget", func(p *technicalRecheckProvider) {
			p.extracted.Files[0].Content = strings.Repeat("x", technicalRecheckMaxBytes+1)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, p := recheckFixture()
			tc.mutate(p)
			_, _, _, asOf := verificationFixture()
			report, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
			if err != nil {
				t.Fatal(err)
			}
			if report.Review.Verifications[0].Status == domain.TechnicalVerified {
				t.Fatalf("unsafe promotion: %+v", report)
			}
			if tc.name != "different bytes" && (report.Bundles[0].Complete || len(report.Failures) == 0) {
				t.Fatalf("read failure not recorded: %+v", report)
			}
		})
	}
}

func TestRecheckTechnicalSourcesCollectsOtherRootLicenses(t *testing.T) {
	req, p := recheckFixture()
	p.tree = append(p.tree, domain.DirectoryNode{Path: "COPYING", Type: "file"})
	p.extracted.Files = append(p.extracted.Files, domain.ExtractedFile{Path: "COPYING", Content: "another license"})
	_, _, _, asOf := verificationFixture()
	report, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.paths, []string{"COPYING", "LICENSE"}) || report.Review.Verifications[0].Status == domain.TechnicalVerified {
		t.Fatalf("root ambiguity omitted: %+v %+v", p.paths, report)
	}
}

func TestRecheckTechnicalSourcesSkipsInactiveFindings(t *testing.T) {
	for _, mode := range []string{"stale", "conflict", "unsupported"} {
		t.Run(mode, func(t *testing.T) {
			req, p := recheckFixture()
			switch mode {
			case "stale":
				req.Batch.Findings[0].ValidUntil = "2026-09-21T00:00:00Z"
			case "conflict":
				f := req.Batch.Findings[0]
				f.ID = "refute"
				f.Claim = "outra redação para o mesmo predicado"
				f.Result = domain.TechnicalRefutes
				req.Batch.Findings = append(req.Batch.Findings, f)
			case "unsupported":
				req.Batch.Findings[0].Predicate = ""
			}
			_, _, _, asOf := verificationFixture()
			report, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
			if err != nil {
				t.Fatal(err)
			}
			if p.calls != 0 || len(report.Bundles) != 0 {
				t.Fatal("inactive finding fetched")
			}
		})
	}
}

func TestRecheckTechnicalSourcesRejectsUnpinnedTargets(t *testing.T) {
	cases := []func(*domain.TechnicalSourceRecheckRequest){
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates[0].CommitID = "HEAD" },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates[0].CommitID = "abc123" },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates[0].URL = "file:///tmp/repo" },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates[0].URL = "https://other.test/o/r" },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates[0].Owner = "../o" },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Candidates = append(r.Candidates, r.Candidates[0]) },
		func(r *domain.TechnicalSourceRecheckRequest) { r.Batch.Findings[0].Source.Ref = "../LICENSE" },
	}
	for i, mutate := range cases {
		req, p := recheckFixture()
		mutate(&req)
		_, _, _, asOf := verificationFixture()
		if _, err := RecheckTechnicalSources(context.Background(), req, p, asOf); err == nil {
			t.Fatalf("invalid target %d accepted", i)
		}
		if p.calls != 0 {
			t.Fatal("invalid target reached provider")
		}
	}
	req, p := recheckFixture()
	_, _, _, asOf := verificationFixture()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RecheckTechnicalSources(ctx, req, p, asOf); !errors.Is(err, context.Canceled) || p.calls != 0 {
		t.Fatalf("cancellation ignored: %v", err)
	}
}

func TestRecheckTechnicalSourcesChecksSBOMDeclaration(t *testing.T) {
	req, p := recheckFixture()
	content := `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"widget","version":"1.2.3"}]}`
	f := &req.Batch.Findings[0]
	f.Technique, f.Predicate, f.Value = domain.TechniqueSBOM, domain.TechnicalSBOMComponent, "widget@1.2.3"
	f.Claim = "SBOM declara widget"
	f.Source.Ref, f.Source.Digest = "sbom.json", "sha256:"+contentHash(content)
	p.tree = []domain.DirectoryNode{{Path: "sbom.json", Type: "file"}}
	p.extracted.Files = []domain.ExtractedFile{{Path: "sbom.json", Content: content}}
	_, _, _, asOf := verificationFixture()
	report, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
	if err != nil {
		t.Fatal(err)
	}
	if report.Review.Verifications[0].Status != domain.TechnicalVerified || !reflect.DeepEqual(p.paths, []string{"sbom.json"}) {
		t.Fatalf("SBOM not rechecked: %+v", report)
	}
	if _, changed := ApplyVerifiedTechnicalLicenses([]domain.CandidateAssessment{{Repository: p.extracted.Ref}}, report.Review); changed {
		t.Fatal("SBOM filled license gate")
	}
}
