/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/research/joss.md
 * @layer L3
 * @updated 2026-10-06
 */

package research

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"bateia/01_core/domain"
)

var jossDOI = regexp.MustCompile(`^10\.21105/joss\.[0-9]{1,8}$`)

type JOSS struct {
	http        httpClient
	crossrefURL string
	paperURL    string
}

func NewJOSS(client *http.Client) *JOSS {
	return &JOSS{http: newHTTPClient(client), crossrefURL: "https://api.crossref.org/journals/2475-9066/works", paperURL: "https://joss.theoj.org/papers/"}
}

func (c *JOSS) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorJOSS, Version: "joss-crossref-v1", Kind: domain.ConnectorProductCatalog,
		QueryLocales: append([]string(nil), locales...), Markets: []string{}, SupportsMarketFilter: false,
		IdentityLevels: []string{"doi"}, EvidenceSurfaces: []string{"peer_review", "software_repository", "software_archive"}, Replayable: true, PublicOnly: true}
}

func (c *JOSS) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	query, limit, err := validateRequest(request, 10)
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	params := url.Values{}
	params.Set("query.title", query)
	params.Set("rows", "20")
	params.Set("select", "DOI,title")
	searchBody, err := c.http.get(ctx, c.crossrefURL+"?"+params.Encode())
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	var search struct {
		Message struct {
			Items []struct {
				DOI string `json:"DOI"`
			} `json:"items"`
		} `json:"message"`
	}
	if err := json.Unmarshal(searchBody, &search); err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	candidates := make([]domain.ProductCandidate, 0, limit)
	bodies := [][]byte{searchBody}
	seen := map[string]bool{}
	for _, item := range search.Message.Items {
		if len(candidates) >= limit {
			break
		}
		doi := strings.ToLower(strings.TrimSpace(item.DOI))
		if !jossDOI.MatchString(doi) || seen[doi] {
			continue
		}
		seen[doi] = true
		body, err := c.http.get(ctx, c.paperURL+doi+".json")
		if err != nil {
			continue
		}
		bodies = append(bodies, body)
		var paper struct {
			Title              string   `json:"title"`
			State              string   `json:"state"`
			DOI                string   `json:"doi"`
			SoftwareRepository string   `json:"software_repository"`
			SoftwareArchive    string   `json:"software_archive"`
			PaperReview        string   `json:"paper_review"`
			Tags               []string `json:"tags"`
		}
		if err := json.Unmarshal(body, &paper); err != nil || paper.State != "accepted" || !strings.EqualFold(paper.DOI, doi) || strings.TrimSpace(paper.Title) == "" {
			continue
		}
		recordURL := c.paperURL + doi
		candidate := domain.ProductCandidate{
			Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorJOSS, ExternalID: doi, CanonicalURL: recordURL},
			Name:     safeText(paper.Title, 300), Summary: safeText(strings.Join(paper.Tags, "; "), 500), ProductKind: "research_software",
			ListingLocales: []string{request.QueryLocale}, ObservedMarkets: []string{}, RepositoryLinks: []domain.ProductRepositoryLink{},
			Resolution:        domain.IdentityDistinct,
			ScholarlyEvidence: &domain.ScholarlyEvidence{Catalogue: "joss", RecordURL: recordURL, DOI: doi, ReviewURL: validPublicURL(paper.PaperReview), Assessment: "PEER_REVIEWED"},
		}
		if upstream := validPublicURL(paper.SoftwareRepository); upstream != "" {
			candidate.FOSSLink = &domain.FOSSLinkage{UpstreamURL: upstream, Status: domain.RepositoryLinkUnverified}
		}
		candidates = append(candidates, candidate)
	}
	status := domain.AttemptEmptyObserved
	if len(candidates) > 0 {
		status = domain.AttemptSucceeded
	}
	return result(status, candidates, bodies...), nil
}

var _ domain.MarketDiscoveryConnector = (*JOSS)(nil)
