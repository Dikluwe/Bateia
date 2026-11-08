/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/insight_discovery_test.md
 * @layer L2
 * @updated 2026-11-08
 */
package cli

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type discoveryCLIPorts struct {
	calls   int
	fail    bool
	content string
}

func (*discoveryCLIPorts) Capabilities() domain.ProviderCapabilities {
	return domain.ProviderCapabilities{Provider: "fixture", Roles: []domain.DiscoveryProviderRole{domain.RoleProjectCatalog}}
}
func (p *discoveryCLIPorts) SearchCatalog(_ context.Context, r domain.CatalogSearchRequest) (domain.CatalogDiscovery, error) {
	p.calls++
	out := domain.CatalogDiscovery{Provider: "fixture", Role: domain.RoleProjectCatalog, Attempt: domain.CatalogAttempt{ExpansionID: r.ExpansionID, Query: r.Query, Status: "SUCCEEDED"}, Candidates: []domain.CatalogCandidate{{Repository: domain.RepositoryCandidate{Owner: "org", Repo: "queue", CommitID: strings.Repeat("a", 40)}}}}
	if p.fail {
		out.Attempt.Status = "RATE_LIMITED"
		return out, errors.New("private-provider-text")
	}
	return out, nil
}
func (*discoveryCLIPorts) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{{Path: "README.md", Type: "file"}}, nil
}
func (p *discoveryCLIPorts) ExtractRepoFiles(_ context.Context, r domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	content := p.content
	if content == "" {
		content = "Queue cancellation reference"
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: r.Owner, Repo: r.Repo, Commit: r.CommitID}, Files: []domain.ExtractedFile{{Path: paths[0], Content: content}}}, nil
}

func TestInsightDiscoveryCLIRequiresValidOutputAndBudgetsBeforeNetwork(t *testing.T) {
	for _, args := range [][]string{
		{"--question", "queue"},
		{"--question", "queue", "--max-per-query", "6"},
		{"--question", "queue", "--timeout", "1ms"},
		{"--question", "queue", "--question", "another"},
		{"--question", "queue", "extra"},
		{"--question", "queue", "--unknown", "yes"},
	} {
		ports := &discoveryCLIPorts{}
		var stdout bytes.Buffer
		output := filepath.Join(t.TempDir(), "out.json")
		if len(args) != 2 {
			args = append(args, "--output", output)
		}
		if err := RunInsightDiscoveryCommand(args, &stdout, ports, ports); err == nil || ports.calls != 0 || stdout.Len() != 0 {
			t.Fatal("invalid command accessed provider or published")
		}
		if _, err := os.Stat(output); !os.IsNotExist(err) {
			t.Fatal("invalid command wrote result")
		}
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"insights", "discover", "--question", "queue"}, &stdout); err == nil {
		t.Fatal("unconfigured library created a remote fallback")
	}
}

func TestInsightDiscoveryCLIPublishesContextAndReviewableFailure(t *testing.T) {
	for _, failed := range []bool{false, true} {
		ports := &discoveryCLIPorts{fail: failed}
		output := filepath.Join(t.TempDir(), "result.json")
		err := RunInsightDiscoveryCommand([]string{"--question", "queue cancellation", "--familiar", "queue", "--output", output}, insightShortWriter{}, ports, ports)
		if (err != nil) != failed {
			t.Fatalf("wrong error for failed=%v: %v", failed, err)
		}
		data, readErr := os.ReadFile(output)
		if readErr != nil {
			t.Fatal(readErr)
		}
		var report domain.InsightDiscoveryReport
		if err := json.Unmarshal(data, &report); err != nil {
			t.Fatal(err)
		}
		if ports.calls != 5 || strings.Contains(string(data), "private-provider-text") || report.ID == "" {
			t.Fatal("routes/protected output missing")
		}
		if failed {
			if report.Status != "FAILED" {
				t.Fatal("total failure hidden")
			}
		} else if len(report.Bundles) != 1 || report.Exploration.Status != "CONTEXT_PREPARED" {
			t.Fatal("automatic readme context missing")
		}
	}
}

func TestInsightDiscoveryCLIBlocksSecretsWithoutReplacingOutput(t *testing.T) {
	const secret = "test-0137-not-a-real-credential"
	t.Setenv("GITHUB_TOKEN", secret)
	ports := &discoveryCLIPorts{content: "queue " + secret}
	output := filepath.Join(t.TempDir(), "out.json")
	if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err := RunInsightDiscoveryCommand([]string{"--question", "queue", "--output", output}, &stdout, ports, ports)
	saved, _ := os.ReadFile(output)
	if err == nil || string(saved) != "preserve" || stdout.Len() != 0 || strings.Contains(err.Error(), secret) {
		t.Fatal("secret leaked or protected output replaced")
	}
}

func TestInsightDiscoveryCLIBooleanDuplicatesNeverReachProvidersOrFactory(t *testing.T) {
	for _, flags := range [][]string{
		{"--algorithm-ideas", "--algorithm-ideas"},
		{"--algorithm-ideas=true", "--algorithm-ideas=true"},
		{"--algorithm-ideas=false", "--algorithm-ideas=false"},
		{"--algorithm-ideas", "-algorithm-ideas=false"},
		{"-algorithm-ideas=false", "--algorithm-ideas"},
	} {
		for _, existing := range []bool{false, true} {
			t.Run(strings.Join(flags, "/")+"/existing="+strconv.FormatBool(existing), func(t *testing.T) {
				output := filepath.Join(t.TempDir(), "result.json")
				if existing {
					if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				ports := &discoveryCLIPorts{}
				factoryCalls := 0
				factory := func() ([]domain.InsightSiteBinding, error) { factoryCalls++; return nil, errors.New("factory reached") }
				var out bytes.Buffer
				args := []string{"--question", "queue", "--sources", "all", "--output", output}
				err := RunInsightDiscoveryCommand(append(args, flags...), &out, ports, ports, factory)
				if err == nil || !strings.Contains(err.Error(), "duplicada") || ports.calls != 0 || factoryCalls != 0 || out.Len() != 0 {
					t.Fatalf("duplicate reached effects: %v", err)
				}
				data, readErr := os.ReadFile(output)
				if existing && (readErr != nil || string(data) != "preserve") || !existing && !os.IsNotExist(readErr) {
					t.Fatal("duplicate changed output")
				}
			})
		}
	}
}

func TestInsightDiscoveryCLIBooleanFormsPreserveModes(t *testing.T) {
	for _, value := range []string{"absent", "implicit", "true", "false", "1", "0", "bad"} {
		t.Run(value, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "result.json")
			args := []string{"--question", "queue", "--query-policy", "lexical-v1", "--output", output}
			if value == "implicit" {
				args = append(args, "--algorithm-ideas")
			} else if value != "absent" {
				args = append(args, "--algorithm-ideas="+value)
			}
			ports := &discoveryCLIPorts{}
			var out bytes.Buffer
			err := RunInsightDiscoveryCommand(args, &out, ports, ports)
			if value == "bad" {
				if err == nil || ports.calls != 0 {
					t.Fatal("invalid boolean reached provider")
				}
				if _, err := os.Stat(output); !os.IsNotExist(err) {
					t.Fatal("invalid boolean published")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			var report domain.InsightDiscoveryReport
			if err := json.Unmarshal(data, &report); err != nil {
				t.Fatal(err)
			}
			want := value == "implicit" || value == "true" || value == "1"
			if report.Plan.Request.AlgorithmIdeas != want || (len(report.Exploration.AlgorithmAgenda) > 0) != want || ports.calls != 5 {
				t.Fatal("boolean mode changed")
			}
		})
	}
}
