/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/research/zbmath.md
 * @layer L3
 * @updated 2026-10-06
 */

package research

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"bateia/01_core/domain"
)

type ZBMath struct {
	http     httpClient
	endpoint string
}

func NewZBMath(client *http.Client) *ZBMath {
	return &ZBMath{http: newHTTPClient(client), endpoint: "https://api.zbmath.org/v1/software/_search"}
}

func (c *ZBMath) Capabilities() domain.DiscoveryConnectorCapabilities {
	return domain.DiscoveryConnectorCapabilities{ID: domain.ConnectorZBMath, Version: "zbmath-api-v1", Kind: domain.ConnectorProductCatalog,
		QueryLocales: append([]string(nil), locales...), Markets: []string{}, SupportsMarketFilter: false,
		IdentityLevels: []string{"zbmath_software_id"}, EvidenceSurfaces: []string{"bibliographic_index", "publication_count", "source_link"}, Replayable: true, PublicOnly: true}
}

func (c *ZBMath) Discover(ctx context.Context, request domain.MarketDiscoveryRequest) (domain.ConnectorDiscovery, error) {
	query, limit, err := validateRequest(request, 25)
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	params := url.Values{}
	params.Set("search_string", query)
	httpStatus, body, err := c.http.getStatus(ctx, c.endpoint+"?"+params.Encode())
	if err != nil {
		return domain.ConnectorDiscovery{}, err
	}
	var payload struct {
		Result []struct {
			ID              int      `json:"id"`
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			Homepage        string   `json:"homepage"`
			SourceCode      string   `json:"source_code"`
			License         string   `json:"license_terms"`
			Keywords        []string `json:"keywords"`
			Authors         []string `json:"authors"`
			ArticlesCount   int      `json:"articles_count"`
			URL             string   `json:"zbmath_url"`
			RelatedSoftware []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"related_software"`
		} `json:"result"`
		Status struct {
			Successful   bool   `json:"execution_bool"`
			InternalCode string `json:"internal_code"`
		} `json:"status"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return domain.ConnectorDiscovery{}, fmt.Errorf("invalid zbMATH response: %w", err)
	}
	if httpStatus == http.StatusNotFound && payload.Status.InternalCode == "successful access. No results found." {
		return result(domain.AttemptEmptyObserved, nil, body), nil
	}
	if httpStatus != http.StatusOK || !payload.Status.Successful {
		return domain.ConnectorDiscovery{}, fmt.Errorf("zbMATH returned HTTP %d", httpStatus)
	}
	candidates := make([]domain.ProductCandidate, 0, limit)
	for _, item := range payload.Result {
		if len(candidates) >= limit {
			break
		}
		if item.ID <= 0 || strings.TrimSpace(item.Name) == "" {
			continue
		}
		recordURL := "https://zbmath.org/software/" + strconv.Itoa(item.ID)
		summary := safeText(item.Description, 500)
		if strings.Contains(summary, "contents unavailable due to conflicting licenses") {
			summary = ""
		}
		if summary == "" {
			summary = safeText(strings.Join(item.Keywords, "; "), 500)
		}
		candidate := domain.ProductCandidate{
			Identity: domain.ProductIdentity{ConnectorID: domain.ConnectorZBMath, ExternalID: strconv.Itoa(item.ID), CanonicalURL: recordURL},
			Name:     safeText(item.Name, 200), Summary: summary, ProductKind: "research_software",
			ListingLocales: []string{request.QueryLocale}, ObservedMarkets: []string{}, RepositoryLinks: []domain.ProductRepositoryLink{},
			Resolution:        domain.IdentityDistinct,
			ScholarlyEvidence: &domain.ScholarlyEvidence{Catalogue: "zbmath", RecordURL: recordURL, Assessment: "INDEXED", ReferenceCount: item.ArticlesCount},
		}
		if len(item.Authors) > 0 {
			candidate.Publisher = safeText(strings.Join(item.Authors, "; "), 300)
		}
		seenRelated := map[int]bool{}
		for _, related := range item.RelatedSoftware {
			if len(candidate.ScholarlyEvidence.RelatedSoftware) == 30 {
				break
			}
			if related.ID <= 0 || related.ID == item.ID || seenRelated[related.ID] || strings.TrimSpace(related.Name) == "" {
				continue
			}
			seenRelated[related.ID] = true
			candidate.ScholarlyEvidence.RelatedSoftware = append(candidate.ScholarlyEvidence.RelatedSoftware, domain.ScholarlyRelatedSoftware{
				Name: safeText(related.Name, 200), RecordURL: "https://zbmath.org/software/" + strconv.Itoa(related.ID),
			})
		}
		upstream := validPublicURL(item.SourceCode)
		if upstream == "" {
			upstream = validPublicURL(item.Homepage)
		}
		if upstream != "" {
			candidate.FOSSLink = &domain.FOSSLinkage{UpstreamURL: upstream, Status: domain.RepositoryLinkUnverified}
		}
		candidates = append(candidates, candidate)
	}
	status := domain.AttemptEmptyObserved
	if len(candidates) > 0 {
		status = domain.AttemptSucceeded
	}
	return result(status, candidates, body), nil
}

var _ domain.MarketDiscoveryConnector = (*ZBMath)(nil)
