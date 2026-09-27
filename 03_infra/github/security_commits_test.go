/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/security_commits_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package github

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"bateia/01_core/domain"
)

func TestMetadataNeverResolvesBranchOrAbbreviationAsImmutableCommit(t *testing.T) {
	for _, returned := range []string{"main", "abc123", "../secret", "invalid?x=1"} {
		c := NewClient("")
		calls := 0
		c.httpClient.Transport = securityTransport(func(*http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return securityResponse(`{"default_branch":"main"}`, 200), nil
			}
			return securityResponse(`{"sha":"`+returned+`"}`, 200), nil
		})
		candidate, err := c.FetchRepoMetadata(context.Background(), domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "HEAD"})
		if !errors.Is(err, errUnsafeRemote) || candidate.CommitID != "HEAD" {
			t.Fatalf("invalid resolution promoted: %+v %v", candidate, err)
		}
	}
}

func TestMetadataRejectsUnsafeDefaultBranchWithoutCommitRequest(t *testing.T) {
	c := NewClient("")
	calls := 0
	c.httpClient.Transport = securityTransport(func(*http.Request) (*http.Response, error) {
		calls++
		if calls > 1 {
			t.Fatal("unsafe branch used in request")
		}
		return securityResponse(`{"default_branch":"../outside"}`, 200), nil
	})
	if _, err := c.FetchRepoMetadata(context.Background(), domain.RepositoryCandidate{Owner: "org", Repo: "repo"}); !errors.Is(err, errUnsafeRemote) {
		t.Fatal("unsafe default branch accepted")
	}
}
