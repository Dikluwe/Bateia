/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_reference_leads.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var insightLeadURLPattern = regexp.MustCompile(`https://[^\s<>()"\[\]]+`)

type insightLeadOrigin struct {
	candidate domain.InsightSiteCandidate
	bundle    domain.EvidenceBundle
}

func insightLeadTarget(raw string) (domain.ProductIdentity, bool) {
	raw = strings.TrimSuffix(strings.TrimRight(raw, ",;"), "/")
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.RawPath != "" {
		return domain.ProductIdentity{}, false
	}
	ids := map[string]string{"github.com": domain.ConnectorGitHubCatalog, "gitee.com": domain.ConnectorGiteeCatalog, "atomgit.com": domain.ConnectorAtomGitCatalog, "gitcode.com": domain.ConnectorAtomGitCatalog, "codeberg.org": domain.ConnectorCodebergCatalog, "gitlab.com": domain.ConnectorGitLabCatalog}
	id, ok := ids[u.Host]
	if !ok {
		return domain.ProductIdentity{}, false
	}
	external := strings.TrimPrefix(u.Path, "/")
	if len(external) > 200 || strings.Contains(external, "/-/") {
		return domain.ProductIdentity{}, false
	}
	target := domain.ProductIdentity{ConnectorID: id, ExternalID: external, CanonicalURL: raw}

	parts := strings.Split(external, "/")
	if len(parts) < 2 || (id != domain.ConnectorGitLabCatalog && len(parts) != 2) {
		return domain.ProductIdentity{}, false
	}
	reserved := map[string]bool{"organizations": true, "orgs": true, "users": true, "topics": true, "collections": true, "marketplace": true, "settings": true, "search": true, "explore": true, "dashboard": true, "help": true, "groups": true, "api": true, "login": true, "logout": true, "signup": true, "user": true}
	if reserved[strings.ToLower(parts[0])] {
		return domain.ProductIdentity{}, false
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." || len(part) > 100 {
			return domain.ProductIdentity{}, false
		}
		for _, r := range part {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
				return domain.ProductIdentity{}, false
			}
		}
	}
	return target, true
}

func addInsightReferenceLeads(report *domain.InsightSiteReport) {
	selected := map[string]domain.InsightSiteCandidate{}
	known := map[string]bool{}
	key := func(identity domain.ProductIdentity) string {
		return identity.ConnectorID + "\x00" + identity.ExternalID
	}
	for _, candidate := range report.Candidates {
		if candidate.Selected {
			selected[candidate.ID] = candidate
			known[key(candidate.Product.Identity)] = true
		}
	}
	origins := []insightLeadOrigin{}
	for _, bundle := range report.Bundles {
		for _, receipt := range report.Readmes {
			if receipt.Status == "COLLECTED" && receipt.Receipt != nil && receipt.Receipt.Source == bundle.Repository && receipt.ConnectorID == domain.ConnectorGitHubCatalog {
				if candidate, ok := selected[receipt.CandidateID]; ok {
					origins = append(origins, insightLeadOrigin{candidate, bundle})
				}
			}
		}
	}
	for _, scoped := range report.SiteContexts {
		if candidate, ok := selected[scoped.CandidateID]; ok && candidate.Product.Identity == scoped.Identity {
			origins = append(origins, insightLeadOrigin{candidate, scoped.Bundle})
		}
	}
	deriveInsightReferenceLeads(report, origins, known)
}

func deriveInsightReferenceLeads(report *domain.InsightSiteReport, origins []insightLeadOrigin, known map[string]bool) {
	key := insightMemoryKey
	sort.Slice(origins, func(i, j int) bool { return origins[i].candidate.ID < origins[j].candidate.ID })
	terms := explorationTerms(report.Plan.Request.Discovery.Question)
	leads := []domain.InsightReferenceLead{}
	pairs := map[string]bool{}
	truncated := false
	for _, origin := range origins {
		if validateEvidenceBundle(origin.bundle) != nil || len(origin.bundle.Items) != 1 {
			report.Gaps = append(report.Gaps, "INVALID_BUNDLE_FOR_REFERENCE_LEADS:"+origin.candidate.ID)
			continue
		}
		item := origin.bundle.Items[0]
		lines := observationLines(item.Content)
		matches := insightLeadURLPattern.FindAllStringIndex(item.Content, 129)
		if len(matches) > 128 {
			truncated = true
			matches = matches[:128]
		}
		for _, match := range matches {
			target, ok := insightLeadTarget(item.Content[match[0]:match[1]])
			if !ok || known[key(target)] {
				continue
			}
			pair := origin.candidate.ID + "\x00" + key(target)
			if pairs[pair] {
				continue
			}
			if len(leads) >= 64 {
				truncated = true
				break
			}
			// Line coordinates use the same normalization as insight reading cards.
			prefix := item.Content[:match[0]]
			start := strings.Count(prefix, "\n")
			if start >= len(lines) {
				continue
			}
			end := start + 1
			if end < len(lines) && strings.TrimSpace(lines[end]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[end]), "#") && len(lines[start])+len(lines[end])+1 <= 1024 {
				end++
			}
			excerpt := strings.Join(lines[start:end], "\n")
			if len(excerpt) > 1024 {
				truncated = true
				continue
			}
			lead := domain.InsightReferenceLead{SourceCandidateID: origin.candidate.ID, SourceIdentity: origin.candidate.Product.Identity, Source: origin.bundle.Repository, Reference: domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: item.Path, ContentHash: item.ContentHash, StartLine: start + 1, EndLine: end}, Excerpt: excerpt}, Target: target, Status: "REFERENCE_NOT_READ", MatchedQuestionTerms: []string{}, SelectionReasons: []string{}}
			words := map[string]bool{}
			for _, word := range explorationWords(excerpt) {
				words[word] = true
			}
			for _, term := range terms {
				if words[term] {
					lead.MatchedQuestionTerms = append(lead.MatchedQuestionTerms, term)
				}
			}
			lead.ID = "reference_lead_" + insightHash(lead)
			leads = append(leads, lead)
			pairs[pair] = true
		}
	}
	sort.Slice(leads, func(i, j int) bool { return leads[i].ID < leads[j].ID })
	usedOrigins, usedTargets := map[string]bool{}, map[string]bool{}
	for len(report.ReferenceLeads) < 12 {
		best, score := -1, [5]int{}
		for i, lead := range leads {
			if lead.ID == "" {
				continue
			}
			current := [5]int{0, 0, len(lead.MatchedQuestionTerms)}
			if !usedTargets[key(lead.Target)] {
				current[0] = 1
			}
			if !usedOrigins[lead.SourceCandidateID] {
				current[1] = 1
			}
			if best < 0 || explorationScoreGreater(current, score) {
				best, score = i, current
			}
		}
		if best < 0 {
			break
		}
		lead := leads[best]
		if score[0] > 0 {
			lead.SelectionReasons = append(lead.SelectionReasons, "NEW_LINKED_TARGET")
		}
		if score[1] > 0 {
			lead.SelectionReasons = append(lead.SelectionReasons, "NEW_ORIGIN_DOCUMENT")
		}
		if score[2] > 0 {
			lead.SelectionReasons = append(lead.SelectionReasons, "LITERAL_QUESTION_TERM_MATCH")
		}
		report.ReferenceLeads = append(report.ReferenceLeads, lead)
		usedTargets[key(lead.Target)], usedOrigins[lead.SourceCandidateID] = true, true
		leads[best].ID = ""
	}
	if len(leads) > len(report.ReferenceLeads) {
		truncated = true
	}
	if truncated {
		report.Gaps = append(report.Gaps, "REFERENCE_LEADS_LIMITED")
	}
}
