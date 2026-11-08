/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_reading.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
	"net/url"
	"strings"
)

func insightReadmeFailureStatus(err error) string {
	var failure domain.InsightReadmeFailure
	if errors.As(err, &failure) {
		switch failure.Code {
		case "AUTH_REQUIRED", "RATE_LIMIT", "NOT_FOUND", "INTEGRITY_ERROR", "CONTENT_LIMIT", "TREE_LIMIT":
			return failure.Code
		}
	}
	return "READ_FAILED"
}

func validInsightSiteRepository(identity domain.ProductIdentity) bool {
	host := ""
	switch identity.ConnectorID {
	case domain.ConnectorGitHubCatalog:
		host = "github.com"
	case domain.ConnectorGiteeCatalog:
		host = "gitee.com"
	case domain.ConnectorAtomGitCatalog:
		host = "atomgit.com"
	case domain.ConnectorCodebergCatalog:
		host = "codeberg.org"
	case domain.ConnectorGitLabCatalog:
		host = "gitlab.com"
	default:
		return false
	}
	u, err := url.Parse(identity.CanonicalURL)
	if err != nil || u.Scheme != "https" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.RawPath != "" || (u.Host != host && !(identity.ConnectorID == domain.ConnectorAtomGitCatalog && u.Host == "gitcode.com")) || u.Path != "/"+identity.ExternalID {
		return false
	}
	parts := strings.Split(identity.ExternalID, "/")
	if len(parts) < 2 || (identity.ConnectorID != domain.ConnectorGitLabCatalog && len(parts) != 2) || (identity.ConnectorID == domain.ConnectorGitLabCatalog && len(identity.ExternalID) > 200) {
		return false
	}
	for _, p := range parts {
		if p == "" || p == "." || p == ".." || len(p) > 100 {
			return false
		}
		for _, r := range p {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
				return false
			}
		}
	}
	return true
}

func collectInsightSiteReadme(ctx context.Context, candidate domain.InsightSiteCandidate, reader domain.InsightSiteRepositoryReader) (domain.InsightReadmeReceipt, *domain.EvidenceBundle) {
	receipt := domain.InsightReadmeReceipt{ReferenceID: candidate.ID, Status: "INTEGRITY_ERROR", Scope: "ROOT_README_ONLY"}
	identity := candidate.Product.Identity
	if !validInsightSiteRepository(identity) {
		return receipt, nil
	}
	if ctx.Err() != nil {
		receipt.Status = "CANCELLED"
		return receipt, nil
	}
	repo, err := reader.ResolveRepository(ctx, identity)
	if err != nil {
		receipt.Status = insightReadmeFailureStatus(err)
		return receipt, nil
	}
	parts := strings.Split(identity.ExternalID, "/")
	if repo.Owner != strings.Join(parts[:len(parts)-1], "/") || repo.Repo != parts[len(parts)-1] || repo.URL != identity.CanonicalURL || !immutableCommitPattern.MatchString(repo.CommitID) {
		return receipt, nil
	}
	return collectInsightReadme(ctx, domain.InsightDiscoveredReference{ID: candidate.ID, Repository: repo}, reader)
}

func prepareInsightSiteContexts(report *domain.InsightSiteReport) error {
	remaining := report.Plan.Request.Discovery.MaxPassages - len(report.Exploration.Cards)
	for i := range report.SiteContexts {
		scoped := &report.SiteContexts[i]
		scoped.Cards = []domain.InsightReadingCard{}
		if remaining <= 0 {
			report.Gaps = append(report.Gaps, "READING_CARD_BUDGET:"+scoped.CandidateID)
			continue
		}
		request := domain.InsightExplorationRequest{SchemaVersion: "1", Question: report.Plan.Request.Discovery.Question, FamiliarApproaches: report.Plan.Request.Discovery.FamiliarApproaches, MaxPassages: 1, Bundles: []domain.EvidenceBundle{scoped.Bundle}}
		var cards []domain.InsightReadingCard
		if scoped.Identity.ConnectorID == domain.ConnectorGitLabCatalog && strings.Contains(scoped.Bundle.Repository.Owner, "/") {
			if !validInsightSiteRepository(scoped.Identity) || scoped.Bundle.Repository.Owner+"/"+scoped.Bundle.Repository.Repo != scoped.Identity.ExternalID || !immutableCommitPattern.MatchString(scoped.Bundle.Repository.Commit) || len(scoped.Bundle.Items) != 1 || !validMechanismPath(scoped.Bundle.Items[0].Path) || len(scoped.Bundle.Items[0].Content) > 32<<10 || validateEvidenceBundle(scoped.Bundle) != nil {
				return domain.ErrIntegrity
			}
			var err error
			cards, _, err = selectInsightPassages(context.Background(), request, explorationTerms(request.Question))
			if err != nil {
				return err
			}
			for i := range cards {
				card := &cards[i]
				card.Source = domain.ImmutableCandidateRef{Repository: scoped.Identity.ExternalID, Commit: scoped.Bundle.Repository.Commit}
				card.ID = "card_" + insightHash(struct {
					Source    domain.ImmutableCandidateRef
					Kind      domain.EvidenceKind
					Reference domain.LocatedInsightReference
				}{card.Source, card.Kind, card.Reference})
			}
		} else {
			explored, err := ExploreInsights(context.Background(), request)
			if err != nil {
				return err
			}
			cards = explored.Cards
		}
		for _, card := range cards {
			card.ID = "site_card_" + insightHash(struct {
				Identity domain.ProductIdentity
				CardID   string
			}{scoped.Identity, card.ID})
			scoped.Cards = append(scoped.Cards, card)
			remaining--
		}
	}
	return nil
}
