/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_queries.md
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
)

func FunctionalLexiconHash(l domain.FunctionalLexicon) string {
	l.Hash = ""
	b, _ := json.Marshal(l)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func PrepareFunctionalLexicon(l domain.FunctionalLexicon) (domain.FunctionalLexicon, error) {
	if l.SchemaVersion != "1" || l.Version == "" || len(l.Frames) == 0 {
		return l, fmt.Errorf("léxico funcional inválido")
	}
	h := FunctionalLexiconHash(l)
	if l.Hash != "" && l.Hash != h {
		return l, fmt.Errorf("hash funcional inválido")
	}
	l.Hash = h
	return l, nil
}
func ExpandFunctionalQueries(ontology domain.OntologyExpansionResult, lexicon domain.FunctionalLexicon, templates []domain.FunctionalQueryTemplate) (domain.FunctionalExpansion, error) {
	l, err := PrepareFunctionalLexicon(lexicon)
	if err != nil {
		return domain.FunctionalExpansion{}, err
	}
	out := domain.FunctionalExpansion{SchemaVersion: "1", LexiconHash: l.Hash, Frames: []domain.FunctionalFrame{}, Queries: []domain.FunctionalQuery{}, Unknown: []domain.FunctionalFrame{}}
	sources := map[domain.DiscoveryConceptID][]string{}
	for _, m := range ontology.Matches {
		sources[m.ConceptID] = appendUnique(sources[m.ConceptID], m.SourceID)
	}
	ids := make([]string, 0, len(sources))
	for id := range sources {
		ids = append(ids, string(id))
	}
	sort.Strings(ids)
	sort.Slice(templates, func(i, j int) bool { return templates[i].ID < templates[j].ID })
	for _, raw := range ids {
		id := domain.DiscoveryConceptID(raw)
		base, ok := l.Frames[id]
		if !ok {
			out.Unknown = append(out.Unknown, domain.FunctionalFrame{ConceptID: id, SourceIDs: sources[id], UnknownReason: "conceito sem frame funcional"})
			continue
		}
		f := base
		f.SourceIDs = uniqueSorted(sources[id])
		if len(f.Actions) == 0 || len(f.Objects) == 0 || (len(f.Inputs) == 0 && len(f.Outputs) == 0) {
			f.UnknownReason = "slots obrigatórios ausentes"
			out.Unknown = append(out.Unknown, f)
			continue
		}
		out.Frames = append(out.Frames, f)
		selected := []domain.FunctionalQueryTemplate{}
		for _, t := range templates {
			if t.ID == "T3" && len(f.Timing) == 0 {
				continue
			}
			selected = append(selected, t)
			if len(selected) == 2 {
				break
			}
		}
		for _, t := range selected {
			text := renderFunctionalTemplate(t.Pattern, f)
			if text == "" {
				continue
			}
			q := domain.FunctionalQuery{FrameID: f.ID, ConceptID: f.ConceptID, TemplateID: t.ID, Text: text, Inputs: first(f.Inputs), Actions: first(f.Actions), Objects: first(f.Objects), Outputs: first(f.Outputs), Timing: first(f.Timing), SourceIDs: append([]string(nil), f.SourceIDs...), LexiconHash: l.Hash}
			h := sha256.Sum256([]byte(string(q.FrameID) + "|" + q.TemplateID + "|" + q.Text))
			q.ID = "fq-" + hex.EncodeToString(h[:6])
			out.Queries = append(out.Queries, q)
		}
	}
	return out, nil
}
func renderFunctionalTemplate(pattern string, f domain.FunctionalFrame) string {
	slots := map[string]string{"{input}": head(f.Inputs), "{action}": head(f.Actions), "{object}": head(f.Objects), "{output}": head(f.Outputs), "{timing}": head(f.Timing)}
	out := pattern
	for k, v := range slots {
		out = strings.ReplaceAll(out, k, v)
	}
	if strings.Contains(out, "{") {
		return ""
	}
	return strings.Join(strings.Fields(out), " ")
}
func AuditFunctionalQueries(exp domain.FunctionalExpansion, l domain.FunctionalLexicon) []domain.FunctionalQueryFinding {
	out := []domain.FunctionalQueryFinding{}
	for _, q := range exp.Queries {
		f := domain.FunctionalQueryFinding{QueryID: q.ID, Verdict: "ACCEPTED", Reasons: []string{}}
		if len(q.Actions) == 0 {
			f.MissingSlots = append(f.MissingSlots, "action")
		}
		if len(q.Objects) == 0 {
			f.MissingSlots = append(f.MissingSlots, "object")
		}
		if len(q.Inputs) == 0 && len(q.Outputs) == 0 {
			f.MissingSlots = append(f.MissingSlots, "input_or_output")
		}
		n := normalizeOntologyText(q.Text)
		for _, x := range l.Forbidden {
			if ontologyPhraseMatch(n, normalizeOntologyText(x)) {
				f.ForbiddenHits = append(f.ForbiddenHits, x)
			}
		}
		if len(f.MissingSlots) > 0 || len(f.ForbiddenHits) > 0 {
			f.Verdict = "REJECTED"
		}
		f.MissingSlots = uniqueSorted(f.MissingSlots)
		f.ForbiddenHits = uniqueSorted(f.ForbiddenHits)
		out = append(out, f)
	}
	return out
}
func first(v []string) []string {
	if len(v) == 0 {
		return []string{}
	}
	return []string{v[0]}
}
func head(v []string) string {
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
