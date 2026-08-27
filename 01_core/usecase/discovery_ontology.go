/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/discovery_ontology.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func normalizeOntologyText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer("á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a", "é", "e", "ê", "e", "ë", "e", "í", "i", "ï", "i", "ó", "o", "õ", "o", "ô", "o", "ö", "o", "ú", "u", "ü", "u", "ç", "c")
	value = replacer.Replace(value)
	return strings.Join(strings.FieldsFunc(value, func(r rune) bool { return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '+') }), " ")
}

func OntologyHash(ontology domain.DiscoveryOntology) string {
	ontology.Hash = ""
	sortOntology(&ontology)
	b, _ := json.Marshal(ontology)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func PrepareDiscoveryOntology(ontology domain.DiscoveryOntology) (domain.DiscoveryOntology, error) {
	if ontology.SchemaVersion != "1" || strings.TrimSpace(ontology.Version) == "" || len(ontology.Concepts) == 0 {
		return ontology, fmt.Errorf("ontologia exige schema 1, version e concepts")
	}
	sortOntology(&ontology)
	seen := map[domain.DiscoveryConceptID]bool{}
	for i := range ontology.Concepts {
		c := &ontology.Concepts[i]
		if c.ID == "" || seen[c.ID] || len(c.Triggers) == 0 || len(c.Queries) == 0 {
			return ontology, fmt.Errorf("conceito inválido ou duplicado: %q", c.ID)
		}
		seen[c.ID] = true
		if len(c.Queries) > 2 {
			return ontology, fmt.Errorf("conceito %s excede duas queries", c.ID)
		}
		for _, q := range c.Queries {
			if normalizeOntologyText(q.Text) == "" || len(q.Discriminators) == 0 {
				return ontology, fmt.Errorf("query inválida em %s", c.ID)
			}
		}
	}
	calculated := OntologyHash(ontology)
	if ontology.Hash != "" && ontology.Hash != calculated {
		return ontology, fmt.Errorf("hash da ontologia inválido")
	}
	ontology.Hash = calculated
	return ontology, nil
}

func ExpandWithOntology(plan domain.ExpansionPlan, ontology domain.DiscoveryOntology) (domain.OntologyExpansionResult, error) {
	prepared, err := PrepareDiscoveryOntology(ontology)
	if err != nil {
		return domain.OntologyExpansionResult{}, err
	}
	out := domain.OntologyExpansionResult{SchemaVersion: "1", OntologyVersion: prepared.Version, OntologyHash: prepared.Hash, Matches: []domain.ConceptMatch{}, Expansions: []domain.SearchExpansion{}, Uncovered: []string{}}
	matchedSources := map[string]bool{}
	seenQuery := map[string]int{}
	for _, source := range plan.Expansions {
		text := normalizeOntologyText(source.Query)
		for _, concept := range prepared.Concepts {
			trigger := ""
			for _, candidate := range concept.Triggers {
				if ontologyPhraseMatch(text, normalizeOntologyText(candidate)) {
					trigger = candidate
					break
				}
			}
			if trigger == "" {
				continue
			}
			matchedSources[source.ID] = true
			out.Matches = append(out.Matches, domain.ConceptMatch{ConceptID: concept.ID, SourceID: source.ID, Trigger: trigger, Strength: 1})
			for _, query := range concept.Queries {
				key := normalizeOntologyText(query.Text)
				if index, ok := seenQuery[key]; ok {
					out.Expansions[index].Covers = appendUnique(out.Expansions[index].Covers, source.Covers...)
					out.Expansions[index].Derivation.SourceIDs = appendUnique(out.Expansions[index].Derivation.SourceIDs, source.ID)
					continue
				}
				expansion := domain.SearchExpansion{Facet: source.Facet, Query: query.Text, Discriminators: append([]string(nil), query.Discriminators...), Rationale: "conceito de descoberta " + string(concept.ID), Covers: append([]string(nil), source.Covers...), Derivation: &domain.ExpansionDerivation{Strategy: "ontology:" + prepared.Version, RuleIDs: []string{string(concept.ID)}, SourceIDs: []string{source.ID}}}
				expansion.ID = semanticExpansionID(expansion.Derivation.Strategy, expansion.Facet, expansion.Query)
				seenQuery[key] = len(out.Expansions)
				out.Expansions = append(out.Expansions, expansion)
			}
		}
	}
	for _, source := range plan.Expansions {
		if !matchedSources[source.ID] {
			out.Uncovered = append(out.Uncovered, source.ID)
		}
	}
	sort.Slice(out.Matches, func(i, j int) bool {
		a, b := out.Matches[i], out.Matches[j]
		return string(a.ConceptID)+"|"+a.SourceID < string(b.ConceptID)+"|"+b.SourceID
	})
	sort.Slice(out.Expansions, func(i, j int) bool {
		return normalizeOntologyText(out.Expansions[i].Query) < normalizeOntologyText(out.Expansions[j].Query)
	})
	sort.Strings(out.Uncovered)
	return out, nil
}

func AuditOntologyQueries(expansions []domain.SearchExpansion, ontology domain.DiscoveryOntology) (domain.QueryQualityReport, error) {
	prepared, err := PrepareDiscoveryOntology(ontology)
	if err != nil {
		return domain.QueryQualityReport{}, err
	}
	byID := map[string]domain.DiscoveryConcept{}
	for _, c := range prepared.Concepts {
		byID[string(c.ID)] = c
	}
	out := domain.QueryQualityReport{SchemaVersion: "1", OntologyVersion: prepared.Version, OntologyHash: prepared.Hash, Findings: []domain.QueryQualityFinding{}, Accepted: []domain.SearchExpansion{}}
	forbidden := []string{}
	for _, c := range prepared.Concepts {
		forbidden = append(forbidden, c.Forbidden...)
	}
	for _, expansion := range expansions {
		f := domain.QueryQualityFinding{ExpansionID: expansion.ID, Query: expansion.Query, Verdict: domain.QueryAccepted, Reasons: []string{}}
		var concept domain.DiscoveryConcept
		if expansion.Derivation == nil || len(expansion.Derivation.RuleIDs) != 1 {
			f.Verdict = domain.QueryRejected
			f.Reasons = append(f.Reasons, "proveniência ontológica ausente")
		} else {
			concept = byID[expansion.Derivation.RuleIDs[0]]
		}
		for _, q := range concept.Queries {
			if normalizeOntologyText(q.Text) == normalizeOntologyText(expansion.Query) {
				f.Discriminators = append([]string(nil), q.Discriminators...)
			}
		}
		tokens := strings.Fields(normalizeOntologyText(expansion.Query))
		if len(tokens) < 2 || len(tokens) > 6 {
			f.Verdict = domain.QueryRejected
			f.Reasons = append(f.Reasons, "query deve ter entre 2 e 6 tokens")
		}
		for _, d := range f.Discriminators {
			if !ontologyPhraseMatch(normalizeOntologyText(expansion.Query), normalizeOntologyText(d)) {
				f.Missing = append(f.Missing, d)
			}
		}
		if len(f.Discriminators) == 0 || len(f.Missing) > 0 {
			f.Verdict = domain.QueryRejected
			f.Reasons = append(f.Reasons, "discriminator ausente")
		}
		for _, x := range forbidden {
			if ontologyPhraseMatch(normalizeOntologyText(expansion.Query), normalizeOntologyText(x)) {
				f.Verdict = domain.QueryRejected
				f.Reasons = append(f.Reasons, "termo proibido: "+x)
			}
		}
		if isGenericOntologyQuery(expansion.Query) {
			f.Verdict = domain.QueryRejected
			f.Reasons = append(f.Reasons, "fallback genérico")
		}
		f.Discriminators, f.Missing, f.Reasons = uniqueSorted(f.Discriminators), uniqueSorted(f.Missing), uniqueSorted(f.Reasons)
		if f.Verdict == domain.QueryAccepted {
			out.Accepted = append(out.Accepted, expansion)
		} else {
			out.Rejected++
		}
		out.Findings = append(out.Findings, f)
	}
	sort.Slice(out.Findings, func(i, j int) bool { return out.Findings[i].ExpansionID < out.Findings[j].ExpansionID })
	sort.Slice(out.Accepted, func(i, j int) bool { return out.Accepted[i].ID < out.Accepted[j].ID })
	return out, nil
}

func EvaluateDiscoveryOntology(corpus domain.OntologyEvaluationCorpus, ontology domain.DiscoveryOntology) (domain.OntologyEvaluationReport, error) {
	prepared, err := PrepareDiscoveryOntology(ontology)
	if err != nil {
		return domain.OntologyEvaluationReport{}, err
	}
	if len(corpus.Cases) == 0 {
		return domain.OntologyEvaluationReport{}, fmt.Errorf("corpus ontológico vazio")
	}
	out := domain.OntologyEvaluationReport{SchemaVersion: "1", OntologyVersion: prepared.Version, OntologyHash: prepared.Hash, Cases: []domain.OntologyCaseResult{}}
	tp, foundTotal, expectedTotal, forbiddenFound, forbiddenTotal, accepted, queryTotal, requiredDisc, retainedDisc, uncovered, sources, controlFound := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	for _, c := range corpus.Cases {
		plan, e := ExpandIdea(c.Idea)
		if e != nil {
			return out, fmt.Errorf("case %s: %w", c.ID, e)
		}
		expanded, e := ExpandWithOntology(*plan, prepared)
		if e != nil {
			return out, e
		}
		quality, e := AuditOntologyQueries(expanded.Expansions, prepared)
		if e != nil {
			return out, e
		}
		foundSet := map[domain.DiscoveryConceptID]bool{}
		for _, m := range expanded.Matches {
			foundSet[m.ConceptID] = true
		}
		expectedSet := map[domain.DiscoveryConceptID]bool{}
		for _, x := range c.ExpectedConcepts {
			expectedSet[x] = true
		}
		result := domain.OntologyCaseResult{ID: c.ID, Found: []domain.DiscoveryConceptID{}, Missing: []domain.DiscoveryConceptID{}, Forbidden: []domain.DiscoveryConceptID{}, Quality: quality}
		for x := range foundSet {
			result.Found = append(result.Found, x)
			foundTotal++
			if expectedSet[x] {
				tp++
			}
			if strings.HasPrefix(c.Group, "control") {
				controlFound++
			}
		}
		for _, x := range c.ExpectedConcepts {
			expectedTotal++
			if !foundSet[x] {
				result.Missing = append(result.Missing, x)
			}
		}
		for _, x := range c.ForbiddenConcepts {
			forbiddenTotal++
			if foundSet[x] {
				forbiddenFound++
				result.Forbidden = append(result.Forbidden, x)
			}
		}
		for _, d := range c.RequiredDiscriminators {
			requiredDisc++
			for _, q := range quality.Accepted {
				if ontologyPhraseMatch(normalizeOntologyText(q.Query), normalizeOntologyText(d)) {
					retainedDisc++
					break
				}
			}
		}
		accepted += len(quality.Accepted)
		queryTotal += len(expanded.Expansions)
		uncovered += len(expanded.Uncovered)
		sources += len(plan.Expansions)
		sort.Slice(result.Found, func(i, j int) bool { return result.Found[i] < result.Found[j] })
		sort.Slice(result.Missing, func(i, j int) bool { return result.Missing[i] < result.Missing[j] })
		out.Cases = append(out.Cases, result)
	}
	if foundTotal > 0 {
		x := float64(tp) / float64(foundTotal)
		out.Metrics.ConceptPrecision = &x
	}
	if expectedTotal > 0 {
		x := float64(tp) / float64(expectedTotal)
		out.Metrics.ConceptRecall = &x
	}
	out.Metrics.QueryAcceptanceRate = ratio(accepted, queryTotal)
	out.Metrics.DiscriminatorRetention = ratio(retainedDisc, requiredDisc)
	out.Metrics.ForbiddenLeakRate = ratio(forbiddenFound, forbiddenTotal)
	out.Metrics.ControlContamination = ratio(controlFound, foundTotal)
	out.Metrics.UncoveredRate = ratio(uncovered, sources)
	out.Passed = out.Metrics.ConceptPrecision != nil && *out.Metrics.ConceptPrecision >= .9 && out.Metrics.ConceptRecall != nil && *out.Metrics.ConceptRecall >= .85 && out.Metrics.QueryAcceptanceRate >= .9 && out.Metrics.DiscriminatorRetention == 1 && out.Metrics.ForbiddenLeakRate == 0 && out.Metrics.ControlContamination == 0
	return out, nil
}

func ontologyPhraseMatch(text, phrase string) bool {
	return phrase != "" && strings.Contains(" "+text+" ", " "+phrase+" ")
}
func isGenericOntologyQuery(q string) bool {
	n := normalizeOntologyText(q)
	for _, x := range []string{"decidir entre", "encontrar projeto", "comparar software"} {
		if n == x {
			return true
		}
	}
	return false
}
func sortOntology(o *domain.DiscoveryOntology) {
	sort.Slice(o.Concepts, func(i, j int) bool { return o.Concepts[i].ID < o.Concepts[j].ID })
	for i := range o.Concepts {
		c := &o.Concepts[i]
		sort.Strings(c.Triggers)
		sort.Slice(c.Queries, func(a, b int) bool {
			return normalizeOntologyText(c.Queries[a].Text) < normalizeOntologyText(c.Queries[b].Text)
		})
		for q := range c.Queries {
			sort.Strings(c.Queries[q].Discriminators)
			sort.Strings(c.Queries[q].NegativeTerms)
		}
		sort.Slice(c.Related, func(a, b int) bool { return c.Related[a] < c.Related[b] })
		sort.Strings(c.Forbidden)
	}
}
