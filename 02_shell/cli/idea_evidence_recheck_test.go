/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_recheck_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type cliTechnicalRecheckProvider struct {
	domain.ContentProvider
	ref     domain.CommitRef
	content string
}

func (p cliTechnicalRecheckProvider) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "LICENSE", Type: "file"}}, nil
}
func (p cliTechnicalRecheckProvider) ExtractRepoFiles(context.Context, domain.RepositoryCandidate, []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{Ref: p.ref, Files: []domain.ExtractedFile{{Path: "LICENSE", Content: p.content}}}, nil
}

func TestEvidenceRecheckExportsFreshBundlesAndReport(t *testing.T) {
	dir := t.TempDir()
	ref := domain.CommitRef{Owner: "o", Repo: "r", Commit: strings.Repeat("a", 40)}
	content := securityMITLicense
	req := domain.TechnicalSourceRecheckRequest{
		Candidates: []domain.RepositoryCandidate{{Owner: ref.Owner, Repo: ref.Repo, CommitID: ref.Commit}},
		Batch:      domain.TechnicalEvidenceBatch{SchemaVersion: "1", Findings: []domain.TechnicalEvidence{{ID: "lic", Technique: domain.TechniqueLicense, Candidate: ref.Key(), Claim: "license MIT", Predicate: domain.TechnicalLicenseID, Value: "MIT", Result: domain.TechnicalSupports, Source: domain.TechnicalEvidenceSource{Ref: "LICENSE", Revision: ref.Commit, Digest: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(content)))}, Method: "read license", ObservedAt: "2026-09-20T00:00:00Z"}}},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	input, bundlesPath := filepath.Join(dir, "request.json"), filepath.Join(dir, "bundles.json")
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	args := []string{"--input", input, "--as-of", "2026-09-24T00:00:00Z", "--bundles-output", bundlesPath}
	if err := RunIdeaEvidenceRecheckCommand(args, &out, cliTechnicalRecheckProvider{ref: ref, content: content}); err != nil {
		t.Fatal(err)
	}
	var report domain.TechnicalSourceRecheckReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Basis != "PROVIDER_RECHECK" || report.Review.Verifications[0].Status != domain.TechnicalVerified {
		t.Fatalf("not rechecked: %s", out.String())
	}
	bundlesData, err := os.ReadFile(bundlesPath)
	if err != nil {
		t.Fatal(err)
	}
	var bundles []domain.EvidenceBundle
	if err := json.Unmarshal(bundlesData, &bundles); err != nil || len(bundles) != 1 || !bundles[0].Complete {
		t.Fatalf("bundles not exported: %v %s", err, bundlesData)
	}
	for _, extra := range [][]string{
		{"--output", input},
		{"--output", bundlesPath},
		{"--output", filepath.Join(dir, "child", "..", "bundles.json")},
		{"--timeout", "0s"}, {"--timeout", "6m"}, {"--as-of", "invalid"},
	} {
		if err := RunIdeaEvidenceRecheckCommand(append(append([]string(nil), args...), extra...), &out, nil); err == nil {
			t.Fatalf("invalid command accepted: %v", extra)
		}
	}
	after, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(data, after) {
		t.Fatal("input overwritten")
	}
	if err := RunIdeaEvidenceRecheckCommand(nil, &out, nil); err == nil {
		t.Fatal("missing input accepted")
	}
}
