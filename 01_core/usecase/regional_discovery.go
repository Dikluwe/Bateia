/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/regional_discovery.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

var canonicalSearchLocales = []string{
	"pt-BR", "en-US", "es-ES", "fr-FR", "de-DE",
	"ja-JP", "ko-KR", "ru-RU", "ar-SA", "hi-IN",
}

func DefaultRegionalQueryCatalog() domain.RegionalQueryCatalog {
	return domain.RegionalQueryCatalog{
		SchemaVersion: "1",
		Version:       "regional-recipes-v1",
		Concepts: []domain.RegionalCatalogConcept{{
			Concept:            "recipe-manager",
			RuleID:             "recipe-manager",
			SourceQuery:        "gerenciador de receitas",
			DerivationStrategy: "controlled-catalog-exact-v1",
			Queries: []domain.RegionalCatalogQuery{
				{Locale: "pt-BR", Query: "gerenciador de receitas"},
				{Locale: "en-US", Query: "recipe manager"},
				{Locale: "es-ES", Query: "gestor de recetas"},
				{Locale: "fr-FR", Query: "gestionnaire de recettes"},
				{Locale: "de-DE", Query: "Rezeptmanager"},
				{Locale: "ja-JP", Query: "レシピ管理"},
				{Locale: "ko-KR", Query: "레시피 관리자"},
				{Locale: "ru-RU", Query: "менеджер рецептов"},
				{Locale: "ar-SA", Query: "مدير الوصفات"},
				{Locale: "hi-IN", Query: "रेसिपी प्रबंधक"},
			},
		}},
	}
}

func canonicalizeSearchLocales(values []string) ([]string, error) {
	allowed := make(map[string]int, len(canonicalSearchLocales))
	for index, locale := range canonicalSearchLocales {
		allowed[locale] = index
	}
	seen := map[string]bool{}
	for _, raw := range values {
		value := strings.TrimSpace(raw)
		if value == "all" {
			for _, locale := range canonicalSearchLocales {
				seen[locale] = true
			}
			continue
		}
		if _, ok := allowed[value]; !ok {
			return nil, fmt.Errorf("search locale %q não suportado", raw)
		}
		seen[value] = true
	}
	out := make([]string, 0, len(seen))
	for _, locale := range canonicalSearchLocales {
		if seen[locale] {
			out = append(out, locale)
		}
	}
	if len(out) > len(canonicalSearchLocales) {
		return nil, fmt.Errorf("no máximo %d search locales são permitidos", len(canonicalSearchLocales))
	}
	return out, nil
}

func validateRegionalCatalog(catalog domain.RegionalQueryCatalog) error {
	if catalog.SchemaVersion != "1" || strings.TrimSpace(catalog.Version) == "" || len(catalog.Concepts) == 0 {
		return fmt.Errorf("catálogo regional inválido")
	}
	allowed := map[string]bool{}
	for _, locale := range canonicalSearchLocales {
		allowed[locale] = true
	}
	concepts := map[string]bool{}
	for _, concept := range catalog.Concepts {
		key := strings.TrimSpace(concept.Concept)
		if key == "" || key != concept.Concept || concepts[key] || strings.TrimSpace(concept.RuleID) == "" || strings.TrimSpace(concept.SourceQuery) == "" || strings.TrimSpace(concept.DerivationStrategy) == "" {
			return fmt.Errorf("conceito regional inconsistente %q", concept.Concept)
		}
		concepts[key] = true
		locales := map[string]bool{}
		for _, query := range concept.Queries {
			if !allowed[query.Locale] || locales[query.Locale] || strings.TrimSpace(query.Query) == "" {
				return fmt.Errorf("consulta regional inconsistente para %s/%s", key, query.Locale)
			}
			locales[query.Locale] = true
		}
	}
	return nil
}

func regionalProposalID(locale, concept, query string) string {
	hash := sha256.Sum256([]byte(locale + "\x00" + concept + "\x00" + query))
	return fmt.Sprintf("exp-regional-%s-%x", locale, hash[:6])
}

func PlanRegionalDiscovery(catalog domain.RegionalQueryCatalog, concept string, locales []string) (domain.RegionalDiscoveryPlan, error) {
	if err := validateRegionalCatalog(catalog); err != nil {
		return domain.RegionalDiscoveryPlan{}, err
	}
	requested, err := canonicalizeSearchLocales(locales)
	if err != nil {
		return domain.RegionalDiscoveryPlan{}, err
	}
	out := domain.RegionalDiscoveryPlan{Requested: requested, Realized: []string{}, Gaps: []domain.RegionalDiscoveryGap{}, Proposals: []domain.SearchExpansion{}}
	if len(requested) == 0 {
		return out, nil
	}
	var selected *domain.RegionalCatalogConcept
	for index := range catalog.Concepts {
		if catalog.Concepts[index].Concept == strings.TrimSpace(concept) {
			selected = &catalog.Concepts[index]
			break
		}
	}
	if selected == nil {
		for _, locale := range requested {
			out.Gaps = append(out.Gaps, domain.RegionalDiscoveryGap{Locale: locale, Concept: strings.TrimSpace(concept), Status: "Unknown", Reason: "conceito sem tradução no catálogo regional"})
		}
		return out, nil
	}
	byLocale := map[string]string{}
	for _, query := range selected.Queries {
		byLocale[query.Locale] = query.Query
	}
	for _, locale := range requested {
		query := byLocale[locale]
		if query == "" {
			out.Gaps = append(out.Gaps, domain.RegionalDiscoveryGap{Locale: locale, Concept: selected.Concept, Status: "Unknown", Reason: "locale sem tradução controlada para o conceito"})
			continue
		}
		id := regionalProposalID(locale, selected.Concept, query)
		provenance := &domain.RegionalQueryProvenance{
			Locale: locale, CatalogVersion: catalog.Version, Concept: selected.Concept,
			RuleID: selected.RuleID, SourceQuery: selected.SourceQuery,
			DerivationStrategy: selected.DerivationStrategy, ProposalID: id,
		}
		out.Realized = append(out.Realized, locale)
		out.Proposals = append(out.Proposals, domain.SearchExpansion{
			ID: id, Facet: domain.FacetProblem, Query: query,
			Rationale: "descoberta regional controlada para " + selected.Concept,
			Covers:    []string{"problem"}, Regional: provenance,
			Derivation: &domain.ExpansionDerivation{
				Strategy: selected.DerivationStrategy,
				RuleIDs:  []string{selected.RuleID}, SourceIDs: []string{selected.SourceQuery},
			},
		})
	}
	return out, nil
}

func regionalConceptForPlan(plan domain.ExpansionPlan) string {
	triggers := []string{
		"organizar receitas culinárias", "receitas culinárias", "gerenciador de receitas", "recipe manager",
		"desperdício alimentar", "food waste", "inventário", "inventory", "previsão de consumo", "consumption forecast",
	}
	for _, expansion := range plan.Expansions {
		for _, trigger := range triggers {
			if containsPhrase(expansion.Query, trigger) || containsPhrase(trigger, expansion.Query) {
				return "recipe-manager"
			}
		}
	}
	return "unmapped-idea"
}

func addRegionalDiscovery(plan domain.ExpansionPlan, locales []string) (domain.ExpansionPlan, error) {
	canonical, err := canonicalizeSearchLocales(locales)
	if err != nil {
		return domain.ExpansionPlan{}, err
	}
	if len(canonical) == 0 {
		return plan, nil
	}
	regional, err := PlanRegionalDiscovery(DefaultRegionalQueryCatalog(), regionalConceptForPlan(plan), canonical)
	if err != nil {
		return domain.ExpansionPlan{}, err
	}
	plan.Regional = &regional
	plan.Expansions = append(plan.Expansions, regional.Proposals...)
	return plan, nil
}

func equalSearchLocales(left, right []string) bool {
	a, errA := canonicalizeSearchLocales(left)
	b, errB := canonicalizeSearchLocales(right)
	if errA != nil || errB != nil || len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}

func canonicalRegionalOrigins(origins []domain.CandidateOrigin) []domain.CandidateOrigin {
	out := append([]domain.CandidateOrigin(nil), origins...)
	sort.Slice(out, func(i, j int) bool {
		left := regionalOriginSortKey(out[i])
		right := regionalOriginSortKey(out[j])
		return left < right
	})
	deduped := make([]domain.CandidateOrigin, 0, len(out))
	seen := map[string]bool{}
	for _, origin := range out {
		encoded, _ := json.Marshal(origin)
		key := string(encoded)
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, origin)
		}
	}
	return deduped
}

func regionalOriginSortKey(origin domain.CandidateOrigin) string {
	rank := 0
	if origin.SearchLocale != "" {
		rank = len(canonicalSearchLocales) + 1
		for index, locale := range canonicalSearchLocales {
			if origin.SearchLocale == locale {
				rank = index + 1
				break
			}
		}
	}
	encoded, _ := json.Marshal(origin)
	return fmt.Sprintf("%02d\x00%s\x00%s\x00%s", rank, origin.RegionalProposalID, origin.PlannedQuery, encoded)
}
