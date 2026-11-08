/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_followup_context.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
)

func prepareInsightFollowupContext(report *domain.InsightSiteReport) error {
	// Reuse the same per-origin passage policy, including GitLab namespaces.
	shadow := domain.InsightSiteReport{Plan: report.Plan, Exploration: report.Exploration, Gaps: []string{}}
	occupied := len(report.Exploration.Cards)
	for _, s := range report.SiteContexts {
		occupied += len(s.Cards)
	}
	shadow.Exploration.Cards = make([]domain.InsightReadingCard, occupied)
	for _, f := range report.Followups {
		if f.Bundle != nil {
			shadow.SiteContexts = append(shadow.SiteContexts, domain.InsightSiteContext{CandidateID: f.EntryID, Identity: f.Identity, Bundle: *f.Bundle})
		}
	}
	if err := prepareInsightSiteContexts(&shadow); err != nil {
		return err
	}
	for _, scoped := range shadow.SiteContexts {
		for i := range report.Followups {
			f := &report.Followups[i]
			if f.EntryID == scoped.CandidateID {
				for _, card := range scoped.Cards {
					card.ID = "followup_card_" + insightHash(struct{ Entry, Card string }{f.EntryID, card.ID})
					f.Cards = append(f.Cards, card)
				}
			}
		}
	}
	report.Gaps = append(report.Gaps, shadow.Gaps...)

	return nil
}
