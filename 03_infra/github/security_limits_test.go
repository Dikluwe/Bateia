/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/security_limits_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package github

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestRemoteEndpointsRejectOversizedBodiesBeforeDecoding(t *testing.T) {
	for _, endpoint := range []string{"tree", "metadata", "manifest", "raw", "catalog", "surface"} {
		t.Run(endpoint, func(t *testing.T) {
			c := NewClient("")
			calls := 0
			c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if calls > 1 {
					t.Fatal("oversized body caused fallback or extra request")
				}
				// Content-Length is enough to reject without reading/allocating the body.
				resp := securityResponse("{}", 200)
				resp.ContentLength = maxTreeEnvelopeBytes + 1
				return resp, nil
			})
			var err error
			switch endpoint {
			case "tree":
				_, err = c.FetchTree(context.Background(), securityCandidate())
			case "metadata":
				_, err = c.FetchRepoMetadata(context.Background(), securityCandidate())
			case "manifest":
				_, err = c.FetchManifest(context.Background(), securityCandidate())
			case "raw":
				_, err = c.fetchRawFallback(context.Background(), "org", "repo", strings.Repeat("a", 40), "LICENSE")
			case "catalog":
				_, err = c.SearchCatalog(context.Background(), domain.CatalogSearchRequest{Query: "q", MaxResults: 1})
			case "surface":
				_, err = NewTopicsSurface(c).SearchSurface(context.Background(), domain.SurfaceSearchRequest{Surface: domain.SurfaceTopics, Query: "q", MaxResults: 1})
			}
			if !errors.Is(err, errUnsafeRemote) {
				t.Fatalf("endpoint %s exceeded body budget: %v", endpoint, err)
			}
		})
	}
}

func TestRemoteTreeRejectsInvalidPathsAndPreservesNonRegularTypes(t *testing.T) {
	for _, tc := range []struct{ mode, path, want string }{
		{"120000", "LICENSE", "symlink"},
		{"160000", "component", "submodule"},
		{"100644", "../outside", ""},
	} {
		c := NewClient("")
		c.httpClient.Transport = securityTransport(func(*http.Request) (*http.Response, error) {
			return securityResponse(securityTree([]string{tc.path}, tc.mode, 1), 200), nil
		})
		nodes, err := c.FetchTree(context.Background(), securityCandidate())
		if tc.want == "" {
			if !errors.Is(err, errUnsafeRemote) {
				t.Fatal("unsafe tree path accepted")
			}
		} else if err != nil || len(nodes) != 1 || nodes[0].Type != tc.want {
			t.Fatalf("non-regular mode disguised as file: %+v %v", nodes, err)
		}
	}
}

func TestExtractionHEADRemainsExplicitlyMutable(t *testing.T) {
	candidate := domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "HEAD"}
	c := NewClient("")
	c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
		if strings.Contains(r.URL.Path, "/git/trees/HEAD") {
			return securityResponse(securityTree([]string{"dir/file"}, "100644", 3), 200), nil
		}
		return securityResponse(securityContent("dir/file", "abc"), 200), nil
	})
	got, err := c.ExtractRepoFiles(context.Background(), candidate, []string{"dir/file"})
	if err != nil || got.Ref.Commit != "HEAD" || got.Files[0].Content != "abc" {
		t.Fatalf("HEAD falsely promoted or rejected: %+v %v", got, err)
	}
}

func TestExtractionRawLimitIsTerminal(t *testing.T) {
	c := NewClient("")
	calls := 0
	c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if strings.Contains(r.URL.Path, "/git/trees/") {
			return securityResponse(securityTree([]string{"LICENSE"}, "100644", 1), 200), nil
		}
		if r.URL.Host == "raw.githubusercontent.com" {
			return securityResponse(strings.Repeat("x", maxFileBytes+1), 200), nil
		}
		return securityResponse("unavailable", 503), nil
	})
	if got, err := c.ExtractRepoFiles(context.Background(), securityCandidate(), []string{"LICENSE"}); got != nil || !errors.Is(err, errUnsafeRemote) {
		t.Fatalf("raw oversized accepted: %+v %v", got, err)
	}
	if calls != 3 {
		t.Fatalf("unexpected fallback attempts: %d", calls)
	}
}
