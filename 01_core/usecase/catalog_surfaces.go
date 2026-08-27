/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/catalog_surfaces.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"fmt"
	"sort"
	"strings"
)

func PlanCatalogSurfaces(plan domain.ExpansionPlan, surfaces []domain.CatalogSurfaceCapabilities) []domain.SurfaceQueryPlan {
	out := []domain.SurfaceQueryPlan{}
	for _, expansion := range plan.Expansions {
		for _, surface := range surfaces {
			p := domain.SurfaceQueryPlan{ExpansionID: expansion.ID, Surface: surface.ID, Original: expansion.Query, Reason: "discriminator ausente"}
			if len(expansion.Discriminators) > 0 {
				p.Effective = expansion.Discriminators[0]
				p.PreservedDiscriminators = []string{expansion.Discriminators[0]}
				p.Valid = true
				p.Reason = "discriminator completo preservado"
			}
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return string(out[i].Surface)+"|"+out[i].ExpansionID < string(out[j].Surface)+"|"+out[j].ExpansionID
	})
	return out
}

func DiscoverCatalogSurfaces(ctx context.Context, plans []domain.SurfaceQueryPlan, providers []domain.CatalogSurface, maxResults int) (domain.MultiViewSnapshot, error) {
	out := domain.MultiViewSnapshot{SchemaVersion: "1", QueryPlans: append([]domain.SurfaceQueryPlan(nil), plans...), Discoveries: []domain.SurfaceDiscovery{}, Candidates: []domain.MultiViewCandidate{}, Unknown: []domain.SurfaceDiscovery{}}
	bySurface := map[domain.CatalogSurfaceID]domain.CatalogSurface{}
	for _, p := range providers {
		c := p.SurfaceCapabilities()
		bySurface[c.ID] = p
	}
	byRepo := map[string]int{}
	for _, plan := range plans {
		if !plan.Valid {
			continue
		}
		provider := bySurface[plan.Surface]
		if provider == nil {
			return out, fmt.Errorf("superfície %s sem provedor", plan.Surface)
		}
		caps := provider.SurfaceCapabilities()
		if !domain.SameSurface(caps.ID, plan.Surface) {
			return out, fmt.Errorf("provedor incompatível com %s", plan.Surface)
		}
		discovery, err := provider.SearchSurface(ctx, domain.SurfaceSearchRequest{ExpansionID: plan.ExpansionID, Surface: plan.Surface, Query: plan.Effective, Discriminators: plan.PreservedDiscriminators, MaxResults: maxResults})
		discovery.Plan = plan
		out.Discoveries = append(out.Discoveries, discovery)
		if err != nil || discovery.Attempt.Status == "RATE_LIMITED" || discovery.Attempt.Status == "TIMED_OUT" {
			out.Unknown = append(out.Unknown, discovery)
			continue
		}
		for _, m := range discovery.Matches {
			if m.Repository == "" || m.CommitID == "" || m.UnknownReason != "" {
				u := discovery
				u.Matches = []domain.SurfaceMatch{m}
				out.Unknown = append(out.Unknown, u)
				continue
			}
			key := strings.ToLower(m.Repository)
			idx, ok := byRepo[key]
			if !ok {
				parts := strings.SplitN(m.Repository, "/", 2)
				if len(parts) != 2 {
					continue
				}
				idx = len(out.Candidates)
				byRepo[key] = idx
				out.Candidates = append(out.Candidates, domain.MultiViewCandidate{Repository: domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], CommitID: m.CommitID}, Matches: []domain.SurfaceMatch{}, Surfaces: []domain.CatalogSurfaceID{}})
			}
			c := &out.Candidates[idx]
			c.Matches = append(c.Matches, m)
			seen := false
			for _, s := range c.Surfaces {
				if s == m.Surface {
					seen = true
				}
			}
			if !seen {
				c.Surfaces = append(c.Surfaces, m.Surface)
			}
		}
	}
	for i := range out.Candidates {
		sort.Slice(out.Candidates[i].Surfaces, func(a, b int) bool { return out.Candidates[i].Surfaces[a] < out.Candidates[i].Surfaces[b] })
		sort.Slice(out.Candidates[i].Matches, func(a, b int) bool {
			return string(out.Candidates[i].Matches[a].Surface)+out.Candidates[i].Matches[a].ExpansionID < string(out.Candidates[i].Matches[b].Surface)+out.Candidates[i].Matches[b].ExpansionID
		})
	}
	sort.Slice(out.Candidates, func(i, j int) bool {
		return out.Candidates[i].Repository.Owner+"/"+out.Candidates[i].Repository.Repo < out.Candidates[j].Repository.Owner+"/"+out.Candidates[j].Repository.Repo
	})
	return out, nil
}

func ScreenMultiView(snapshot domain.MultiViewSnapshot) domain.CatalogScreeningReport {
	out := domain.CatalogScreeningReport{SchemaVersion: "1", Decisions: []domain.CatalogDecision{}}
	for _, c := range snapshot.Candidates {
		score := .35
		signals := []domain.CatalogSignal{{ID: "IDENTITY_RESOLVED", Contribution: .35, Reason: "commit imutável"}}
		verdict := domain.CatalogReview
		if len(c.Surfaces) >= 2 {
			score += .5
			signals = append(signals, domain.CatalogSignal{ID: "INDEPENDENT_SURFACES", Contribution: .5, Reason: "discriminator em duas superfícies"})
			verdict = domain.CatalogAdvance
		}
		out.Decisions = append(out.Decisions, domain.CatalogDecision{Repository: c.Repository, Verdict: verdict, Score: score, Signals: signals})
	}
	sort.Slice(out.Decisions, func(i, j int) bool { return out.Decisions[i].Score > out.Decisions[j].Score })
	return out
}
