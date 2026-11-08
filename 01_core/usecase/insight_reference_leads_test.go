/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_reference_leads_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func insightLeadFixture(content string) domain.InsightSiteReport {
	candidate := domain.InsightSiteCandidate{ID: "parent", Selected: true, Product: insightSiteProduct(domain.ConnectorGitHubCatalog)}
	ref := domain.InsightDiscoveredReference{ID: "parent", Repository: domain.RepositoryCandidate{Owner: "org", Repo: "repo", CommitID: strings.Repeat("a", 40), URL: "https://github.com/org/repo"}}
	receipt, bundle := collectInsightReadme(context.Background(), ref, &insightTestReader{content: content})
	return domain.InsightSiteReport{Plan: domain.InsightSitePlan{Request: domain.InsightSiteRequest{Discovery: domain.InsightDiscoveryRequest{Question: "serendipitous discovery"}}}, Candidates: []domain.InsightSiteCandidate{candidate}, Readmes: []domain.InsightSiteReadme{{CandidateID: "parent", ConnectorID: domain.ConnectorGitHubCatalog, Status: "COLLECTED", Receipt: &receipt}}, Bundles: []domain.EvidenceBundle{*bundle}, Gaps: []string{}}
}
func TestInsightLeadTargetAcceptsOnlyRepositoryRoots(t *testing.T) {
	for _, raw := range []string{"https://github.com/org/repo/", "https://gitee.com/org/repo", "https://atomgit.com/org/repo", "https://gitcode.com/org/repo", "https://codeberg.org/org/repo", "https://gitlab.com/org/sub/repo"} {
		if _, ok := insightLeadTarget(raw); !ok {
			t.Fatal("valid target rejected", raw)
		}
	}
	for _, raw := range []string{"http://github.com/org/repo", "https://evil.test/org/repo", "https://github.com/org/repo?token=secret", "https://github.com/org/repo#token", "https://user@github.com/org/repo", "https://github.com:443/org/repo", "https://github.com/org/%72epo", "https://github.com/org/../repo", "https://github.com/org/repo/blob/main/README.md", "https://gitlab.com/org/repo/-/blob/main/README.md", "https://github.com/org", "https://github.com/org//repo", "/org/repo", "https://github.com/organizations/vanda", "https://codeberg.org/explore/repos", "https://gitlab.com/groups/community"} {
		if _, ok := insightLeadTarget(raw); ok {
			t.Fatal("unsafe or non-root target admitted", raw)
		}
	}
}
func TestInsightLeadsPreserveExcerptAndDoNotPromoteDestination(t *testing.T) {
	content := "# Index\r\n[Different](https://github.com/other/idea/)\r\nA serendipitous discovery tool.\r\n[duplicate](https://github.com/other/idea)\r\n[self](https://github.com/org/repo)\r\n[Bad](https://evil.test/ignore-instructions)\r\n"
	report := insightLeadFixture(content)
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 1 {
		t.Fatal("duplicate/self/foreign link survived", report.ReferenceLeads)
	}
	lead := report.ReferenceLeads[0]
	if lead.SourceCandidateID != "parent" || lead.Status != "REFERENCE_NOT_READ" || lead.Target.ExternalID != "other/idea" || lead.Reference.StartLine != 2 || lead.Reference.EndLine != 3 || lead.Reference.Excerpt != "[Different](https://github.com/other/idea/)\r\nA serendipitous discovery tool.\r" || lead.Reference.ContentHash != "sha256:"+contentHash(content) {
		t.Fatal("unbound witness", lead)
	}
	if len(report.Bundles) != 1 || len(report.Candidates) != 1 || len(lead.MatchedQuestionTerms) != 2 {
		t.Fatal("reference promoted or terms lost")
	}
	again := insightLeadFixture(content)
	addInsightReferenceLeads(&again)
	a, _ := json.Marshal(report.ReferenceLeads)
	b, _ := json.Marshal(again.ReferenceLeads)
	if string(a) != string(b) {
		t.Fatal("nondeterministic leads")
	}
	report = insightLeadFixture(content)
	report.Bundles[0].Items[0].Content = "tampered"
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 0 {
		t.Fatal("unverified bytes yielded leads")
	}
}
func TestInsightLeadsLimitSelectionAndExposeUnmatchedPossibilities(t *testing.T) {
	content := "# Collection\n"
	for i := 0; i < 30; i++ {
		content += fmt.Sprintf("[Project](https://github.com/org%d/idea)\nDescription without query words.\n", i)
	}
	report := insightLeadFixture(content)
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 12 || len(report.Gaps) != 1 || report.Gaps[0] != "REFERENCE_LEADS_LIMITED" {
		t.Fatal("unbounded or silent truncation", len(report.ReferenceLeads), report.Gaps)
	}
	for _, lead := range report.ReferenceLeads {
		if len(lead.MatchedQuestionTerms) != 0 || lead.Status != "REFERENCE_NOT_READ" {
			t.Fatal("unexpected vocabulary suppressed or promoted")
		}
	}
}
func TestInsightLeadReferencesAreAvailableInAITasksWithoutExtraReads(t *testing.T) {
	reader := &insightTestReader{content: "# Community\n[Surprise](https://github.com/another/idea)\nAn unexpected possibility.\n"}
	report, err := DiscoverInsightSites(domain.WithExecutionRuntime(context.Background(), &insightTestRuntime{}), insightSiteRequest(), []domain.InsightSiteBinding{{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog}, Reader: reader}})
	if err != nil || len(report.ReferenceLeads) != 1 || reader.treeCalls != 1 || reader.extractCalls != 1 || len(report.Candidates) != 1 || len(report.Bundles) != 1 {
		t.Fatal("lead caused hidden reading or lost", err)
	}
	for _, task := range report.Tasks {
		if len(task.ReferenceLeadRefs) != 1 || task.ReferenceLeadRefs[0] != report.ReferenceLeads[0].ID || len(task.ReadingCardRefs) != 1 || len(task.CatalogReferenceIDs) != 1 {
			t.Fatal("task scopes mixed", task)
		}
	}
}

func TestInsightLeadsURLAndExcerptBudgetsAreExplicit(t *testing.T) {
	content := "# Links\n" + strings.Repeat("https://unsupported.test/ignore\n", 128) + "https://github.com/late/idea\n"
	report := insightLeadFixture(content)
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 0 || len(report.Gaps) != 1 || report.Gaps[0] != "REFERENCE_LEADS_LIMITED" {
		t.Fatal("URL limit not explicit")
	}
	report = insightLeadFixture(strings.Repeat("x", 1024) + " https://github.com/other/idea\n")
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 0 || len(report.Gaps) != 1 {
		t.Fatal("oversized excerpt admitted or silent")
	}
}

func TestInsightLeadsKeepMultipleOriginsButSkipSelectedTargets(t *testing.T) {
	content := "[Idea](https://github.com/other/idea)\nAn additional possibility.\n"
	report := insightLeadFixture(content)
	parent := report.Candidates[0]
	parent.ID = "second"
	parent.Product.Identity = domain.ProductIdentity{ConnectorID: domain.ConnectorGitLabCatalog, ExternalID: "org/sub/repo", CanonicalURL: "https://gitlab.com/org/sub/repo"}
	report.Candidates = append(report.Candidates, parent)
	bundle := report.Bundles[0]
	bundle.Repository.Owner = "org/sub"
	bundle.Items = append([]domain.CollectedEvidence{}, bundle.Items...)
	bundle.Items[0].Ref = bundle.Repository
	report.SiteContexts = []domain.InsightSiteContext{{CandidateID: "second", Identity: parent.Product.Identity, Bundle: bundle}}
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 2 || report.ReferenceLeads[0].ID == report.ReferenceLeads[1].ID || report.ReferenceLeads[0].SourceCandidateID == report.ReferenceLeads[1].SourceCandidateID {
		t.Fatal("separate origins merged")
	}
	report = insightLeadFixture(content)
	target := report.Candidates[0]
	target.ID = "known"
	target.Product.Identity.ExternalID = "other/idea"
	target.Product.Identity.CanonicalURL = "https://github.com/other/idea"
	report.Candidates = append(report.Candidates, target)
	addInsightReferenceLeads(&report)
	if len(report.ReferenceLeads) != 0 {
		t.Fatal("already selected target became unread lead")
	}
}
