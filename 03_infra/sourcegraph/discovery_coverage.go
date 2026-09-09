/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/discovery_coverage.md
 * @layer L3
 * @updated 2026-09-08
 */
package sourcegraph

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"context"
	"strings"
	"time"
)

func discoveryRequestedMatches(limit int) int {
	if limit >= 20 {
		return 200
	}
	if limit <= 2 {
		return 20
	}
	return 10 * limit
}

func discoveryContains(text string, terms []string) bool {
	for _, term := range terms {
		if term != "" && strings.Contains(strings.ToLower(text), strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func diagnoseDiscovery(plan domain.DiscoveryQueryPlan, diagnostic *domain.DiscoveryCandidateDiagnostic) []domain.MatchEvidence {
	texts := map[string][]string{}
	for _, item := range diagnostic.Evidence {
		texts[item.Path] = append(texts[item.Path], item.Snippet)
	}
	accepted := map[string]bool{}
	for path, snippets := range texts {
		text := strings.Join(snippets, "\n")
		diagnostic.ActionFound = diagnostic.ActionFound || discoveryContains(text, plan.ActionTerms)
		diagnostic.ObjectFound = diagnostic.ObjectFound || discoveryContains(text, plan.ObjectTerms)
		accepted[path] = parser.DiscoveryEvidenceMatches(plan, text)
	}
	var evidence []domain.MatchEvidence
	for _, item := range diagnostic.Evidence {
		if accepted[item.Path] {
			evidence = append(evidence, item)
		}
	}
	diagnostic.Status = "insufficient_evidence"
	switch {
	case len(evidence) > 0:
		diagnostic.Status, diagnostic.Reason = "selected", "matched_same_file"
	case len(diagnostic.Evidence) == 0:
		diagnostic.Reason = "no_file_evidence"
	case diagnostic.ActionFound && diagnostic.ObjectFound:
		diagnostic.Reason = "split_across_files"
	case diagnostic.ActionFound:
		diagnostic.Reason = "missing_object"
	case diagnostic.ObjectFound:
		diagnostic.Reason = "missing_action"
	default:
		diagnostic.Reason = "missing_both"
	}
	return evidence
}

func (c *Client) executeDerivedDiscovery(ctx context.Context, query string, limit int, plan domain.DiscoveryQueryPlan) ([]domain.RepositoryCandidate, []domain.ProviderCandidateMatch, *domain.DiscoveryCoverage, []domain.DiscoveryCandidateDiagnostic, error) {
	response, err := c.fetchSGQuery(ctx, query)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if limit <= 0 {
		limit = 1
	}
	coverage := &domain.DiscoveryCoverage{RequestedMatches: discoveryRequestedMatches(limit), ReceivedItems: len(response.Data.Search.Results.Results), ProviderMatches: response.Data.Search.Results.MatchCount, Completeness: "unknown"}
	observed := map[string]bool{}
	indexes := map[string]int{}
	var pool []domain.RepositoryCandidate
	var diagnostics []domain.DiscoveryCandidateDiagnostic
	retainedCounts := map[string]int{}
	for _, item := range response.Data.Search.Results.Results {
		if err := ctx.Err(); err != nil {
			return nil, nil, nil, nil, err
		}
		var repo sgRepo
		switch item.TypeName {
		case "Repository":
			repo = sgRepo{Name: item.Name, URL: item.URL, Description: item.Description}
		case "FileMatch":
			repo = item.Repository
		default:
			continue
		}
		owner, name := parseRepoName(repo.Name)
		if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
			continue
		}
		key := strings.ToLower(owner + "/" + name)
		observed[key] = true
		index, exists := indexes[key]
		if !exists {
			if len(pool) >= coverage.RequestedMatches {
				continue
			}
			index = len(pool)
			indexes[key] = index
			pool = append(pool, domain.RepositoryCandidate{Owner: owner, Repo: name, URL: "https://" + repo.Name, Description: repo.Description, CommitID: "HEAD", ConsultadoEm: time.Now().UTC()})
			diagnostics = append(diagnostics, domain.DiscoveryCandidateDiagnostic{Repository: owner + "/" + name, Evidence: []domain.MatchEvidence{}})
		}
		if item.TypeName != "FileMatch" {
			continue
		}
		for _, line := range item.LineMatches {
			diagnostic := &diagnostics[index]
			if retainedCounts[key] >= 32 {
				diagnostic.OmittedSnippets++
				continue
			}
			retainedCounts[key]++
			normalized := strings.Join(strings.Fields(line.Preview), " ")
			snippet := truncateMatchSnippet(normalized)
			if snippet != normalized {
				diagnostic.TruncatedSnippets++
			}
			// Empty snippets and missing paths cannot participate in admission.
			if snippet == "" || strings.TrimSpace(item.File.Path) == "" {
				continue
			}
			diagnostic.Evidence = append(diagnostic.Evidence, domain.MatchEvidence{MatchType: "file", Path: item.File.Path, Snippet: snippet})
		}
	}
	coverage.ObservedRepositories = len(observed)
	coverage.ExaminedRepositories = len(pool)
	coverage.DeferredRepositories = len(observed) - len(pool)
	candidates := make([]domain.RepositoryCandidate, 0)
	matches := make([]domain.ProviderCandidateMatch, 0)
	for i, candidate := range pool {
		diagnostic := &diagnostics[i]
		evidence := diagnoseDiscovery(plan, diagnostic)
		if diagnostic.Status != "selected" {
			continue
		}
		coverage.EligibleCandidates++
		if len(candidates) >= limit {
			diagnostic.Status, diagnostic.Reason = "candidate_limit", "candidate_limit"
			continue
		}
		candidates = append(candidates, candidate)
		matches = append(matches, domain.ProviderCandidateMatch{Repository: diagnostic.Repository, Evidence: evidence})
	}
	coverage.SelectedCandidates = len(candidates)
	return candidates, matches, coverage, diagnostics, nil
}
