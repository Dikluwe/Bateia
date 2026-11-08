/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_candidates.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

func validInsightSiteCandidates(values []domain.ProductCandidate, id string, limit int) bool {
	if len(values) > limit {
		return false
	}
	for _, value := range values {
		if value.Identity.ConnectorID != id || !validSynthesisText(value.Identity.ExternalID, 300) || !validSynthesisText(value.Name, 500) || len(canonicalJSON(value)) > 4096 {
			return false
		}
		if raw := value.Identity.CanonicalURL; raw != "" {
			u, err := url.Parse(raw)
			if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
				return false
			}
		}
	}
	return true
}

func mergeInsightSiteCandidates(prior []domain.InsightSiteCandidate, values []domain.ProductCandidate, cap domain.DiscoveryConnectorCapabilities, queryID, locale, responseHash string, familiar []string, markets ...string) []domain.InsightSiteCandidate {
	byKey := map[string]domain.InsightSiteCandidate{}
	for _, candidate := range prior {
		byKey[candidate.ID] = candidate
	}
	for _, value := range values {
		// Own all nested slices/pointers before canonicalizing and sealing provenance.
		encoded, _ := json.Marshal(value)
		var product domain.ProductCandidate
		_ = json.Unmarshal(encoded, &product)
		product = canonicalProduct(product)
		market := ""
		if len(markets) > 0 {
			market = markets[0]
		}
		product.Origins = []domain.MarketCandidateOrigin{{ConnectorID: cap.ID, ConnectorVersion: cap.Version, QueryID: queryID, QueryLocale: locale, ResponseHash: responseHash, ObservedMarket: market}}
		id := "catalog_ref_" + insightHash(struct{ Connector, External string }{cap.ID, product.Identity.ExternalID})
		candidate, found := byKey[id]
		if !found {
			candidate = domain.InsightSiteCandidate{ID: id, Product: product, QueryIDs: []string{}, Origin: "CATALOG_METADATA_UNVERIFIED", FamiliarSubstringMatches: []string{}, SelectionReasons: []string{}}
		} else {
			origins := canonicalOrigins(append(append([]domain.MarketCandidateOrigin{}, candidate.Product.Origins...), product.Origins...))
			links := canonicalRepoLinks(append(append([]domain.ProductRepositoryLink{}, candidate.Product.RepositoryLinks...), product.RepositoryLinks...))
			left, right := candidate.Product, product
			left.Origins = nil
			right.Origins = nil
			left.RepositoryLinks = nil
			right.RepositoryLinks = nil
			if string(canonicalJSON(right)) < string(canonicalJSON(left)) {
				candidate.Product = product
			}
			candidate.Product.Origins = origins
			candidate.Product.RepositoryLinks = links
		}
		candidate.QueryIDs = canonicalMechanismSet(append(append([]string{}, candidate.QueryIDs...), queryID))
		haystack := strings.ToLower(strings.Join(strings.Fields(candidate.Product.Name+" "+candidate.Product.Summary), " "))
		candidate.FamiliarSubstringMatches = []string{}
		for _, label := range familiar {
			if strings.Contains(haystack, strings.ToLower(label)) {
				candidate.FamiliarSubstringMatches = append(candidate.FamiliarSubstringMatches, label)
			}
		}
		byKey[id] = candidate
	}
	out := make([]domain.InsightSiteCandidate, 0, len(byKey))
	for _, value := range byKey {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		return string(canonicalJSON(out[i].Product.Identity)) < string(canonicalJSON(out[j].Product.Identity))
	})
	return out
}

func selectInsightSiteCandidates(values []domain.InsightSiteCandidate, maximum int, question ...string) []domain.InsightSiteCandidate {
	out := append([]domain.InsightSiteCandidate{}, values...)
	sites, queries, publishers := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for order := 1; order <= maximum; order++ {
		best, score := -1, [5]int{}
		for i, c := range out {
			if c.Selected {
				continue
			}
			gain := 0
			for _, q := range c.QueryIDs {
				if !queries[q] {
					gain++
				}
			}
			current := [5]int{0, gain, 0, -len(c.FamiliarSubstringMatches)}
			if !sites[c.Product.Identity.ConnectorID] {
				current[0] = 1
			}
			if c.Product.Publisher != "" && !publishers[c.Product.Identity.ConnectorID+"\x00"+c.Product.Publisher] {
				current[2] = 1
			}
			if len(question) > 0 && question[0] != "" {
				matched := insightMetadataMatchCount(question[0], c.Product.Name+" "+c.Product.Summary)
				if matched == 0 {
					continue
				}
				current = [5]int{matched, current[0], gain, current[2], -len(c.FamiliarSubstringMatches)}
			}
			if best < 0 || explorationScoreGreater(current, score) {
				best, score = i, current
			}
		}
		if best < 0 {
			break
		}
		c := &out[best]
		c.Selected = true
		c.SelectionOrder = order
		if len(question) > 0 && question[0] != "" {
			score = [5]int{score[1], score[2], score[3], score[4]}
		}
		c.SelectionReasons = []string{}
		if len(question) > 0 && question[0] != "" {
			c.SelectionReasons = append(c.SelectionReasons, "LITERAL_METADATA_QUESTION_MATCH")
		}
		if score[0] > 0 {
			c.SelectionReasons = append(c.SelectionReasons, "NEW_SITE")
		}
		if score[1] > 0 {
			c.SelectionReasons = append(c.SelectionReasons, "NEW_QUERY_COVERAGE")
		}
		if score[2] > 0 {
			c.SelectionReasons = append(c.SelectionReasons, "NEW_PUBLISHER")
		}
		sites[c.Product.Identity.ConnectorID] = true
		publishers[c.Product.Identity.ConnectorID+"\x00"+c.Product.Publisher] = true
		for _, q := range c.QueryIDs {
			queries[q] = true
		}
	}
	return out
}
