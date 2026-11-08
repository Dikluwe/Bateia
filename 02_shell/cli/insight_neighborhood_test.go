/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_neighborhood_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type neighborhoodCLIReader struct {
	reads []string
	extra string
	fail  bool
}

func (r *neighborhoodCLIReader) ResolveRepository(_ context.Context, id domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	if r.fail {
		return domain.RepositoryCandidate{}, domain.InsightReadmeFailure{Code: "AUTH_REQUIRED"}
	}
	parts := strings.Split(id.ExternalID, "/")
	return domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: id.CanonicalURL, CommitID: strings.Repeat("a", 40)}, nil
}
func (*neighborhoodCLIReader) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
}
func (r *neighborhoodCLIReader) ExtractRepoFiles(_ context.Context, repo domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	r.reads = append(r.reads, repo.Repo)
	text := "# Leaf\nNo further references.\n"
	if repo.Repo == "repo" {
		text = "# Root\n[Child](https://gitee.com/org/child)\n"
	}
	if repo.Repo == "child" {
		text = "# Child\n[Grandchild](https://gitee.com/org/grandchild)\n" + r.extra
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}, Files: []domain.ExtractedFile{{Path: "README.md", Content: text}}}, nil
}

func neighborhoodCLIInput(t *testing.T) (string, *neighborhoodCLIReader, *insightSiteCLIFixture, domain.InsightSiteFactory) {
	t.Helper()
	reader, connector := &neighborhoodCLIReader{}, &insightSiteCLIFixture{}
	factory := func() ([]domain.InsightSiteBinding, error) {
		return []domain.InsightSiteBinding{{Connector: connector, RepositoryReader: reader}}, nil
	}
	input := filepath.Join(t.TempDir(), "initial.json")
	if err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--sources", "gitee", "--max-calls", "1", "--follow-links", "0", "--output", input}, &bytes.Buffer{}, nil, nil, factory); err != nil {
		t.Fatal(err)
	}
	return input, reader, connector, factory
}

func TestInsightNeighborhoodCLIChainsReportsWithoutQuestionOrCatalog(t *testing.T) {
	input, reader, connector, factory := neighborhoodCLIInput(t)
	before, _ := os.ReadFile(input)
	second, third := filepath.Join(t.TempDir(), "second.json"), filepath.Join(t.TempDir(), "third.json")
	args := []string{"--input", input, "--focus", "https://gitee.com/org/child", "--follow-links", "1", "--output", second}
	if err := RunInsightNeighborhoodCommand(args, &bytes.Buffer{}, factory); err != nil {
		t.Fatal(err)
	}
	if err := RunInsightNeighborhoodCommand([]string{"--input", second, "--focus", "https://gitee.com/org/child", "--output", third}, &bytes.Buffer{}, factory); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(input)
	data, _ := os.ReadFile(third)
	var report domain.InsightSiteReport
	_ = json.Unmarshal(data, &report)
	if string(before) != string(after) || connector.calls != 1 || strings.Join(reader.reads, ",") != "repo,child,grandchild" || report.Calls != 0 || report.Plan.Request.Discovery.Question != "queue" || report.Neighborhood == nil {
		t.Fatal("CLI required a new investigation or lost input/history")
	}
	if _, err := readInsightMemory(third); err != nil {
		t.Fatal("expanded report cannot supply next memory", err)
	}
}

func TestInsightNeighborhoodCLIRejectsInputsBeforeFactory(t *testing.T) {
	input, _, _, _ := neighborhoodCLIInput(t)
	invalid := filepath.Join(t.TempDir(), "invalid.json")
	_ = os.WriteFile(invalid, []byte(`{"schema_version":"insight-site-discovery-v1","unknown":true}`), 0600)
	for _, extra := range [][]string{
		{"--input", invalid}, {"--focus", "https://gitee.com/unknown/repo"},
		{"--question", "new question"}, {"--sources", "all"}, {"--max-passages", "0"},
		{"--follow-links", "5"}, {"--timeout", "0s"}, {"--focus", ""},
		{"--output", input},
	} {
		called := false
		factory := func() ([]domain.InsightSiteBinding, error) { called = true; return nil, nil }
		args := []string{"--input", input, "--output", filepath.Join(t.TempDir(), "out.json")}
		// Replace input/output once so invalid JSON and aliases reach their respective checks.
		if extra[0] == "--input" {
			args[1], extra = extra[1], nil
		} else if extra[0] == "--output" {
			args[3], extra = extra[1], nil
		}
		if err := RunInsightNeighborhoodCommand(append(args, extra...), &bytes.Buffer{}, factory); err == nil || called {
			t.Fatal("invalid input reached factory", extra, err)
		}
	}
	alias := filepath.Join(t.TempDir(), "alias.json")
	if err := os.Link(input, alias); err != nil {
		t.Fatal(err)
	}
	if err := RunInsightNeighborhoodCommand([]string{"--input", input, "--output", alias}, &bytes.Buffer{}, nil); err == nil {
		t.Fatal("hardlink input overwritten")
	}
}

func TestInsightNeighborhoodCLIFeedbackFailuresAndSecretProtection(t *testing.T) {
	input, reader, _, factory := neighborhoodCLIInput(t)
	feedback, output := filepath.Join(t.TempDir(), "feedback.json"), filepath.Join(t.TempDir(), "result.json")
	if err := runInsightMemory([]string{"--input", input, "--target", "https://gitee.com/org/child", "--state", "discarded", "--reason", "Already explored", "--output", feedback}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if err := RunInsightNeighborhoodCommand([]string{"--input", input, "--memory", feedback, "--follow-links", "1", "--output", output}, &bytes.Buffer{}, factory); err != nil || len(reader.reads) != 1 {
		t.Fatal("feedback was not applied", err)
	}
	reader.fail = true
	if err := RunInsightNeighborhoodCommand([]string{"--input", input, "--follow-links", "1", "--output", output}, &bytes.Buffer{}, factory); err == nil {
		t.Fatal("failed read appeared successful")
	}
	data, _ := os.ReadFile(output)
	var report domain.InsightSiteReport
	if json.Unmarshal(data, &report) != nil || report.ID == "" || report.Status != "FAILED" || report.Followups[0].Status != "AUTH_REQUIRED" {
		t.Fatal("failed expansion did not publish pending history")
	}
	reader.fail = false
	secret := "fixture-neighborhood-token<>&"
	t.Setenv("GITHUB_TOKEN", secret)
	reader.extra = secret
	_ = os.WriteFile(output, []byte("preserve"), 0600)
	err := RunInsightNeighborhoodCommand([]string{"--input", input, "--follow-links", "1", "--output", output}, &bytes.Buffer{}, factory)
	data, _ = os.ReadFile(output)
	if err == nil || string(data) != "preserve" || strings.Contains(err.Error(), secret) {
		t.Fatal("secret escaped publication protection")
	}
}

func TestNeighborhoodCLIInheritsFocusAndSupportsExplicitGeneralMode(t *testing.T) {
	input, reader, connector, factory := neighborhoodCLIInput(t)
	dir := t.TempDir()
	second, third, general := filepath.Join(dir, "second.json"), filepath.Join(dir, "third.json"), filepath.Join(dir, "general.json")
	if err := RunInsightNeighborhoodCommand([]string{"--input", input, "--focus", "https://gitee.com/org/child", "--follow-links", "1", "--output", second}, &bytes.Buffer{}, factory); err != nil {
		t.Fatal(err)
	}
	if err := RunInsightNeighborhoodCommand([]string{"--input", second, "--output", third}, &bytes.Buffer{}, factory); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(third)
	var report domain.InsightSiteReport
	if json.Unmarshal(raw, &report) != nil || report.Neighborhood.Focus == nil || report.Neighborhood.Focus.ExternalID != "org/child" || strings.Join(reader.reads, ",") != "repo,child,grandchild" || connector.calls != 1 {
		t.Fatal("CLI did not inherit focus without a new search")
	}
	if err := RunInsightNeighborhoodCommand([]string{"--input", second, "--all-trails", "--follow-links", "0", "--output", general}, &bytes.Buffer{}, factory); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile(general)
	report = domain.InsightSiteReport{}
	if json.Unmarshal(raw, &report) != nil || report.Neighborhood.Focus != nil {
		t.Fatal("explicit general mode retained focus")
	}
	for _, n := range report.Neighborhood.Neighbors {
		if n.Distance != nil {
			t.Fatal("general mode fabricated focus distance")
		}
	}
	for _, flags := range [][]string{{"--focus", "https://gitee.com/org/child", "--all-trails"}, {"--all-trails", "--all-trails"}} {
		called := false
		unused := func() ([]domain.InsightSiteBinding, error) { called = true; return nil, nil }
		args := append([]string{"--input", second, "--output", filepath.Join(dir, "reject.json")}, flags...)
		if err := RunInsightNeighborhoodCommand(args, &bytes.Buffer{}, unused); err == nil || called {
			t.Fatal("conflicting/duplicate scope flags reached factory", err)
		}
	}
}
