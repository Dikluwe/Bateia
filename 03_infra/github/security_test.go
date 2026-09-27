/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/security_test.md
 * @layer L3
 * @updated 2026-09-26
 */
package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

import "bateia/01_core/domain"

type securityTransport func(*http.Request) (*http.Response, error)

func (f securityTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func securityResponse(body string, code int) *http.Response {
	return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header), ContentLength: -1}
}

func securityCandidate() domain.RepositoryCandidate {
	return domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: strings.Repeat("a", 40)}
}

func securityTree(paths []string, mode string, size int64) string {
	items := []map[string]any{}
	for _, path := range paths {
		items = append(items, map[string]any{"path": path, "type": "blob", "mode": mode, "size": size})
	}
	body, _ := json.Marshal(map[string]any{"truncated": false, "tree": items})
	return string(body)
}

func securityContent(path, text string) string {
	body, _ := json.Marshal(map[string]any{"path": path, "type": "file", "size": len(text), "encoding": "base64", "content": base64.StdEncoding.EncodeToString([]byte(text))})
	return string(body)
}

func TestSecureExtractionPositiveAndPinnedRawFallback(t *testing.T) {
	candidate := securityCandidate()
	for _, rawFallback := range []bool{false, true} {
		t.Run(fmt.Sprint(rawFallback), func(t *testing.T) {
			c := NewClient("secret")
			calls := 0
			c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if strings.Contains(r.URL.Path, "/git/trees/") {
					if !strings.Contains(r.URL.Path, candidate.CommitID) {
						t.Fatal("tree not pinned")
					}
					return securityResponse(securityTree([]string{"LICENSE"}, "100644", 3), 200), nil
				}
				if r.URL.Host == "raw.githubusercontent.com" {
					if r.Header.Get("Authorization") != "" || !strings.Contains(r.URL.Path, candidate.CommitID) {
						t.Fatal("raw credentials or mutable revision")
					}
					return securityResponse("abc", 200), nil
				}
				if r.URL.Query().Get("ref") != candidate.CommitID {
					t.Fatal("contents not pinned")
				}
				if rawFallback {
					return securityResponse("unavailable", 503), nil
				}
				return securityResponse(securityContent("LICENSE", "abc"), 200), nil
			})
			extracted, err := c.ExtractRepoFiles(context.Background(), candidate, []string{"LICENSE"})
			if err != nil || len(extracted.Files) != 1 || extracted.Files[0].Content != "abc" || extracted.Ref.Commit != candidate.CommitID {
				t.Fatalf("positive: %+v %v", extracted, err)
			}
			if calls != 2 && !rawFallback || calls != 3 && rawFallback {
				t.Fatalf("unexpected calls: %d", calls)
			}
		})
	}
}

func TestSecureExtractionRejectsGitSymlinkAndSubmoduleBeforeContent(t *testing.T) {
	for _, mode := range []string{"120000", "160000", "", "100644"} {
		t.Run(mode, func(t *testing.T) {
			c := NewClient("")
			calls := 0
			c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if calls > 1 {
					t.Fatal("unsafe entry reached content or raw fallback")
				}
				size := int64(0)
				if mode == "100644" {
					size = maxFileBytes + 1
				}
				return securityResponse(securityTree([]string{"LICENSE"}, mode, size), 200), nil
			})
			if result, err := c.ExtractRepoFiles(context.Background(), securityCandidate(), []string{"LICENSE"}); !errors.Is(err, errUnsafeRemote) || result != nil {
				t.Fatalf("unsafe entry accepted: %+v %v", result, err)
			}
		})
	}
}

func TestSecureExtractionRejectsInvalidInputWithoutNetwork(t *testing.T) {
	for _, tc := range []struct {
		candidate domain.RepositoryCandidate
		paths     []string
	}{
		{domain.RepositoryCandidate{Owner: "../org", Repo: "repo"}, []string{"LICENSE"}},
		{domain.RepositoryCandidate{Owner: "org", Repo: "repo?x=1"}, []string{"LICENSE"}},
		{domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: "main?x=1"}, []string{"LICENSE"}},
		{securityCandidate(), []string{"../secret"}},
		{securityCandidate(), []string{"/secret"}},
		{securityCandidate(), []string{"dir/../secret"}},
		{securityCandidate(), []string{"C:\\secret"}},
		{securityCandidate(), []string{"LICENSE", "LICENSE"}},
		{securityCandidate(), make([]string, maxExtractionFiles+1)},
	} {
		c := NewClient("")
		c.httpClient.Transport = securityTransport(func(*http.Request) (*http.Response, error) {
			t.Fatal("invalid input performed network request")
			return nil, nil
		})
		if _, err := c.ExtractRepoFiles(context.Background(), tc.candidate, tc.paths); !errors.Is(err, errUnsafeRemote) {
			t.Fatalf("input accepted: %+v %v", tc, err)
		}
	}
}

func TestSecureExtractionMalformedOrOversizedContentsCannotFallback(t *testing.T) {
	cases := []string{
		strings.Repeat(" ", maxContentEnvelopeBytes+1),
		`{"path":"LICENSE","type":"symlink","target":"/etc/passwd","content":"abc"}`,
		`{"path":"other","type":"file","content":"abc"}`,
		`{"path":"LICENSE","type":"file","encoding":"base64","content":"!"}`,
		securityContent("LICENSE", strings.Repeat("x", maxFileBytes+1)),
	}
	for i, body := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			c := NewClient("")
			calls := 0
			c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if calls == 1 {
					return securityResponse(securityTree([]string{"LICENSE"}, "100644", 1), 200), nil
				}
				if calls > 2 {
					t.Fatal("unsafe envelope silently used raw fallback")
				}
				return securityResponse(body, 200), nil
			})
			if result, err := c.ExtractRepoFiles(context.Background(), securityCandidate(), []string{"LICENSE"}); !errors.Is(err, errUnsafeRemote) || result != nil {
				t.Fatalf("accepted: %+v %v", result, err)
			}
		})
	}
}

func TestSecureExtractionTotalBudget(t *testing.T) {
	paths := []string{"a", "b", "c", "d", "e"}
	c := NewClient("")
	c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
		if strings.Contains(r.URL.Path, "/git/trees/") {
			return securityResponse(securityTree(paths, "100644", maxFileBytes), 200), nil
		}
		return securityResponse(securityContent(filepath.Base(r.URL.Path), strings.Repeat("x", maxFileBytes)), 200), nil
	})
	if result, err := c.ExtractRepoFiles(context.Background(), securityCandidate(), paths); !errors.Is(err, errUnsafeRemote) || result != nil {
		t.Fatalf("total budget not enforced: %+v %v", result, err)
	}
}

func TestBoundedRemoteReadLimitsKnownAndChunkedSizes(t *testing.T) {
	for _, known := range []bool{false, true} {
		resp := securityResponse(strings.Repeat("x", 101), 200)
		if known {
			resp.ContentLength = 101
		}
		if _, err := readBounded(resp, 100); !errors.Is(err, errUnsafeRemote) {
			t.Fatalf("limit not enforced: %v", err)
		}
	}
}

func TestGitHubRedirectNeverForwardsCredentials(t *testing.T) {
	redirected := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirected++
		t.Error("redirect followed")
		w.WriteHeader(200)
	}))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 302) }))
	defer source.Close()
	c := NewClientWithBaseURL("secret", source.URL)
	if _, err := c.FetchTree(context.Background(), securityCandidate()); err == nil {
		t.Fatal("redirect considered successful")
	}
	if redirected != 0 {
		t.Fatal("request escaped host")
	}
}

func TestFailedExtractionNeverClonesOrReadsLocalSymlink(t *testing.T) {
	// Real local Git fixture demonstrates the former attack path without allowing
	// the provider to spawn any Git process, even with hostile inherited config.
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git unavailable for fixture")
	}
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("synthetic secret"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repo, "LICENSE")); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "--quiet", repo}, {"-C", repo, "add", "LICENSE"}, {"-C", repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@invalid", "commit", "--quiet", "-m", "fixture"}} {
		cmd := exec.Command(git, args...)
		cmd.Env = []string{"PATH=/usr/bin:/bin", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("fixture git: %s %v", output, err)
		}
	}
	cmd := exec.Command(git, "clone", "--quiet", "--no-hardlinks", repo, filepath.Join(t.TempDir(), "clone"))
	cmd.Env = []string{"PATH=/usr/bin:/bin", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("fixture clone: %s %v", output, err)
	}
	// An executable PATH shim would fail the test if the old clone fallback ran.
	shimDir := t.TempDir()
	marker := filepath.Join(shimDir, "called")
	shim := "#!/bin/sh\nprintf called > '" + marker + "'\nexit 1\n"
	if err := os.WriteFile(filepath.Join(shimDir, "git"), []byte(shim), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", shimDir)
	t.Setenv("GITHUB_TOKEN", "synthetic-token")
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "url.file://"+repo+".insteadOf")
	t.Setenv("GIT_CONFIG_VALUE_0", "https://github.com/org/repo.git")
	c := NewClient("")
	c.httpClient.Transport = securityTransport(func(r *http.Request) (*http.Response, error) {
		if strings.Contains(r.URL.Path, "/git/trees/") {
			return securityResponse(securityTree([]string{"LICENSE"}, "100644", 20), 200), nil
		}
		return securityResponse("unavailable", 503), nil
	})
	result, err := c.ExtractRepoFiles(context.Background(), securityCandidate(), []string{"LICENSE"})
	if err == nil || result != nil {
		t.Fatalf("failure became evidence: %+v %v", result, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("Git subprocess was invoked")
	}
}
