/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/research/research_test.md
 * @layer L3
 * @updated 2026-10-06
 */

package research

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func TestZBMathPreservesBibliographicEvidenceWithoutInventingReview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("search_string") != "computer algebra" {
			t.Errorf("query=%s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":{"execution_bool":true},"result":[{"id":6663,"name":"Giac","description":"zbMATH Open Web Interface contents unavailable due to conflicting licenses.","homepage":"https://university.example/giac","source_code":null,"keywords":["computer algebra system","CAS"],"articles_count":31,"related_software":[{"id":6663,"name":"Giac"},{"id":5895,"name":"Symja"},{"id":5895,"name":"Duplicate"},{"id":0,"name":"Bad"}]}]}`))
	}))
	defer server.Close()
	connector := NewZBMath(server.Client())
	connector.endpoint = server.URL
	got, err := connector.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "computer algebra", QueryLocale: "en-US", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.AttemptSucceeded || got.ResponseHash == "" || len(got.Candidates) != 1 {
		t.Fatalf("result=%+v", got)
	}
	candidate := got.Candidates[0]
	if candidate.Name != "Giac" || candidate.Identity.ExternalID != "6663" || candidate.ScholarlyEvidence.Assessment != "INDEXED" || candidate.ScholarlyEvidence.ReferenceCount != 31 || candidate.ScholarlyEvidence.DOI != "" {
		t.Fatalf("candidate=%+v", candidate)
	}
	if len(candidate.ScholarlyEvidence.RelatedSoftware) != 1 || candidate.ScholarlyEvidence.RelatedSoftware[0].RecordURL != "https://zbmath.org/software/5895" {
		t.Fatalf("related=%+v", candidate.ScholarlyEvidence.RelatedSoftware)
	}
	if strings.Contains(candidate.Summary, "unavailable") || candidate.FOSSLink.Status != domain.RepositoryLinkUnverified {
		t.Fatalf("candidate=%+v", candidate)
	}
}

func TestJOSSRequiresAcceptedArticleAndValidDOI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/crossref") {
			_, _ = w.Write([]byte(`{"message":{"items":[{"DOI":"10.21105/joss.03438"},{"DOI":"10.21105/joss.01118"},{"DOI":"10.9999/foreign"}]}}`))
			return
		}
		if strings.HasSuffix(r.URL.Path, "03438.json") {
			_, _ = w.Write([]byte(`{"title":"Computer algebra in R","state":"accepted","doi":"10.21105/joss.03438","software_repository":"https://example.edu/source","paper_review":"https://github.com/openjournals/joss-reviews/issues/3438","tags":["algebra"]}`))
			return
		}
		_, _ = w.Write([]byte(`{"title":"Unreviewed","state":"submitted","doi":"10.21105/joss.01118"}`))
	}))
	defer server.Close()
	connector := NewJOSS(server.Client())
	connector.crossrefURL = server.URL + "/crossref"
	connector.paperURL = server.URL + "/papers/"
	got, err := connector.Discover(context.Background(), domain.MarketDiscoveryRequest{Query: "computer algebra", QueryLocale: "en-US", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.AttemptSucceeded || len(got.Candidates) != 1 {
		t.Fatalf("result=%+v", got)
	}
	candidate := got.Candidates[0]
	if candidate.ScholarlyEvidence.Assessment != "PEER_REVIEWED" || candidate.ScholarlyEvidence.DOI != "10.21105/joss.03438" || candidate.FOSSLink.UpstreamURL != "https://example.edu/source" || candidate.FOSSLink.Status != domain.RepositoryLinkUnverified {
		t.Fatalf("candidate=%+v", candidate)
	}
}
