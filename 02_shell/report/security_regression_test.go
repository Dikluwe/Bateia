/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/security_regression_test.md
 * @layer L2
 * @updated 2026-09-26
 */

package report

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestReportURLPolicy(t *testing.T) {
	for _, tc := range []struct {
		url   string
		local bool
		valid bool
	}{
		{"https://github.com/org/repo", false, true},
		{"http://example.com/repo", false, true},
		{"file:///tmp/repository", true, true},
		{"file://localhost/tmp/repository", true, true},
		{"file:///tmp/repository", false, false},
		{"file://remote-server/share/repository", true, false},
		{"file:relative", true, false},
		{"file:///tmp/repository?query=1", true, false},
		{"javascript:alert(1)", true, false},
		{"JaVaScRiPt:alert(1)", false, false},
		{"data:text/html,<script>danger()</script>", true, false},
		{"//example.com/repo", true, false},
		{"/relative", true, false},
		{"https:example.com/repo", false, false},
		{"https://user:password@example.com/repo", false, false},
		{"https://example.com/%0aevil", false, false},
		{"https://example.com/\\evil", false, false},
		{"https://example.com/%5cevil", false, false},
		{"https://example.com/repo\n", false, false},
	} {
		t.Run(tc.url, func(t *testing.T) {
			if got := safeReportURL(tc.url, tc.local); (got != "") != tc.valid {
				t.Fatalf("safeReportURL(%q, %t) = %q; valid=%t", tc.url, tc.local, got, tc.valid)
			}
		})
	}
}

func TestHTMLLinksEscapeAttributesAndText(t *testing.T) {
	got := reportHTMLLink(`https://example.com/?q=" onmouseover="danger()`, `<img src=x onerror=danger()>`, false)
	if strings.Contains(got, ` onmouseover="`) || strings.Contains(got, "<img") {
		t.Fatalf("unescaped attribute or markup: %s", got)
	}
	if !strings.Contains(got, "&#34;") || !strings.Contains(got, `rel="noopener noreferrer"`) {
		t.Fatalf("expected escaped query and opener protection: %s", got)
	}
	if got := reportHTMLLink("javascript:danger()", "blocked", true); got != "blocked" {
		t.Fatalf("unsafe link must become text: %s", got)
	}
}

func TestHTMLReportUntrustedDataCannotIntroduceMarkup(t *testing.T) {
	payload := "</pre><script>danger()</script><img src=x onerror=danger()>"
	rep := &domain.SearchReport{
		Query:       domain.QueryFilter{Query: payload, Language: payload},
		GeneratedAt: time.Date(2026, 9, 26, 0, 0, 0, 0, time.FixedZone(payload, 0)),
		Results: []domain.ExtractedRepository{{
			Candidate:          domain.RepositoryCandidate{Owner: payload, Repo: payload, URL: "javascript:danger()"},
			Manifest:           &domain.Manifest{Dependencies: []domain.Dependency{{Name: payload}}},
			Files:              []domain.ExtractedFile{{Path: `source" onmouseover="danger().go`, Content: payload}},
			SecurityAdvisories: []domain.SecurityAdvisory{{Source: "data:text/html," + payload, Message: payload}},
			SecurityScan:       &domain.SecurityScanReport{Status: domain.SecurityScannedNoMatch, Scanner: payload, Scope: payload, Limitations: []string{payload}},
		}},
	}
	got := GenerateHTMLReport(rep)
	for _, bad := range []string{"<script", "<img", `href="javascript:`, `href="data:`, ` onmouseover="`, "cdn.jsdelivr.net", "mermaid.initialize"} {
		if strings.Contains(got, bad) {
			t.Errorf("HTML contains active payload %q", bad)
		}
	}
	if !strings.Contains(got, "&lt;script&gt;") || !strings.Contains(got, "script-src 'none'") || !strings.Contains(got, "mermaid-source") {
		t.Fatal("report must retain inert content, defensive CSP and offline graph source")
	}
	for _, tag := range regexp.MustCompile(`<a\s[^>]*>`).FindAllString(got, -1) {
		if !strings.Contains(tag, `rel="noopener noreferrer"`) {
			t.Errorf("unprotected anchor: %s", tag)
		}
	}
}

func TestMermaidLabelsCannotBreakGraphOrMarkdownFence(t *testing.T) {
	payload := "\"]\nclick Root \"javascript:danger()\"\n```\n</pre><script>danger()</script> & #quot;"
	rep := &domain.SearchReport{
		Query: domain.QueryFilter{Query: payload},
		Results: []domain.ExtractedRepository{{
			Candidate: domain.RepositoryCandidate{Owner: payload, Repo: "repo"},
			Manifest:  &domain.Manifest{Dependencies: []domain.Dependency{{Name: payload}, {Name: "a-b"}, {Name: "a_b"}}},
		}},
	}
	got := GenerateMermaidGraph(rep)
	if strings.Count(got, "```") != 2 || strings.Count(got, "\n") != 8 {
		t.Fatalf("payload broke line or fence boundaries: %q", got)
	}
	for _, bad := range []string{"\nclick ", "<script", "</pre>", "#quot;"} {
		if strings.Contains(got, bad) {
			t.Errorf("unescaped Mermaid syntax %q", bad)
		}
	}
	if !strings.Contains(got, "#34;") || !strings.Contains(got, "#35;quot#59;") || !strings.Contains(got, "Dep_1_2") || !strings.Contains(got, "Dep_1_3") {
		t.Fatal("labels require entity escaping and dependency IDs independent of input")
	}
}

func TestSourceFileURLConfinesLocalLinksAndEscapesPath(t *testing.T) {
	repo := domain.ExtractedRepository{Candidate: domain.RepositoryCandidate{URL: "file:///tmp/repo"}}
	for _, path := range []string{"../outside", "/outside", "a/../../outside", "a\\..\\outside"} {
		if got := sourceFileURL(repo, path); got != "" {
			t.Errorf("unsafe path %q => %q", path, got)
		}
	}
	if got := sourceFileURL(repo, `a b".go`); got != "file:///tmp/repo/a%20b%22.go" {
		t.Fatalf("local path must be URL-escaped, got %q", got)
	}
	repo.Candidate.URL = "file://remote-server/share"
	if got := sourceFileURL(repo, "ok.go"); got != "" {
		t.Fatalf("remote file authority allowed: %q", got)
	}
}

func TestSecurityScanRenderedEvenWithoutAdvisories(t *testing.T) {
	for _, tc := range []struct {
		scan *domain.SecurityScanReport
		want string
	}{
		{nil, "Não avaliado"},
		{&domain.SecurityScanReport{Status: domain.SecurityNotScanned}, "NOT_SCANNED"},
		{&domain.SecurityScanReport{Status: domain.SecurityScannedNoMatch}, "não demonstra ausência de vulnerabilidades"},
		{&domain.SecurityScanReport{Status: domain.SecurityVulnerabilityFound}, "alcance e impacto precisam de validação"},
	} {
		rep := &domain.SearchReport{Results: []domain.ExtractedRepository{{SecurityScan: tc.scan}}}
		for _, rendered := range []string{GenerateHTMLReport(rep), GenerateMarkdownReport(rep)} {
			if !strings.Contains(rendered, tc.want) {
				t.Fatalf("missing honest scanner state %q", tc.want)
			}
		}
	}
}
