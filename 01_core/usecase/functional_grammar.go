/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/functional_grammar.md
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

func FunctionalGrammarHash(g domain.FunctionalGrammar) string {
	g.Hash = ""
	sort.Slice(g.Rules, func(i, j int) bool { return g.Rules[i].ID < g.Rules[j].ID })
	sort.Slice(g.Morphology, func(i, j int) bool { return g.Morphology[i].ID < g.Morphology[j].ID })
	b, _ := json.Marshal(g)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func RealizeFunctionalGrammar(expansion domain.FunctionalExpansion, grammar domain.FunctionalGrammar) (domain.FunctionalGrammarExpansion, error) {
	if grammar.SchemaVersion == "" || grammar.Version == "" || len(grammar.Rules) == 0 {
		return domain.FunctionalGrammarExpansion{}, fmt.Errorf("gramática inválida")
	}
	hash := FunctionalGrammarHash(grammar)
	rules := append([]domain.FunctionalGrammarRule(nil), grammar.Rules...)
	sort.Slice(rules, func(i, j int) bool {
		if rules[i].Intent != rules[j].Intent {
			return rules[i].Intent < rules[j].Intent
		}
		if rules[i].Priority != rules[j].Priority {
			return rules[i].Priority < rules[j].Priority
		}
		return rules[i].ID < rules[j].ID
	})
	frames := append([]domain.FunctionalFrame(nil), expansion.Frames...)
	sort.Slice(frames, func(i, j int) bool { return frames[i].ID < frames[j].ID })

	result := domain.FunctionalGrammarExpansion{SchemaVersion: "2", GrammarHash: hash}
	for _, frame := range frames {
		matched := 0
		for _, rule := range rules {
			if rule.Intent != string(frame.ConceptID) {
				continue
			}
			query, finding := realizeFrameRule(frame, grammar, rule, hash)
			if finding.Verdict != "ACCEPTED" {
				result.Unknown = append(result.Unknown, finding)
				continue
			}
			result.Queries = append(result.Queries, query)
			matched++
			if matched >= 2 {
				break
			}
		}
		if matched == 0 {
			result.Unknown = append(result.Unknown, domain.GrammarFinding{QueryID: string(frame.ID), GrammarVersion: grammar.Version, Verdict: "UNKNOWN", Reasons: []string{"no applicable grammar rule"}})
		}
	}
	sort.Slice(result.Queries, func(i, j int) bool { return result.Queries[i].ID < result.Queries[j].ID })
	return result, nil
}

func realizeFrameRule(frame domain.FunctionalFrame, grammar domain.FunctionalGrammar, rule domain.FunctionalGrammarRule, hash string) (domain.RealizedFunctionalQuery, domain.GrammarFinding) {
	finding := domain.GrammarFinding{QueryID: string(frame.ID), GrammarVersion: grammar.Version, RuleID: rule.ID, Verdict: "ACCEPTED"}
	slots := map[string]string{"input": firstValue(frame.Inputs), "action": firstValue(frame.Actions), "object": firstValue(frame.Objects), "output": firstValue(frame.Outputs), "timing": firstValue(frame.Timing)}
	for _, required := range rule.RequiredSlots {
		if slots[required] == "" {
			finding.MissingSlots = append(finding.MissingSlots, required)
		}
	}
	if len(finding.MissingSlots) > 0 {
		finding.Verdict = "UNKNOWN"
		finding.Reasons = []string{"required slots unavailable"}
		return domain.RealizedFunctionalQuery{}, finding
	}
	action, morphID, ok := realizeAction(slots["action"], rule.SubjectKind, grammar.Morphology)
	if !ok {
		finding.Verdict = "UNKNOWN"
		finding.UnresolvedVerbs = []string{slots["action"]}
		finding.Reasons = []string{"morphology rule unavailable"}
		return domain.RealizedFunctionalQuery{}, finding
	}
	slots["action"] = action
	text := rule.Pattern
	for _, name := range []string{"input", "action", "object", "output", "timing"} {
		text = strings.ReplaceAll(text, "{"+name+"}", slots[name])
	}
	text = strings.Join(strings.Fields(text), " ")
	if strings.Contains(text, "{") {
		finding.Verdict = "UNKNOWN"
		finding.Reasons = []string{"unresolved template slot"}
		return domain.RealizedFunctionalQuery{}, finding
	}
	id := stableGrammarQueryID(frame.ID, rule.ID, text, hash)
	query := domain.RealizedFunctionalQuery{FunctionalQuery: domain.FunctionalQuery{ID: id, FrameID: frame.ID, ConceptID: frame.ConceptID, TemplateID: rule.ID, Text: text, Inputs: frame.Inputs, Actions: frame.Actions, Objects: frame.Objects, Outputs: frame.Outputs, Timing: frame.Timing, SourceIDs: frame.SourceIDs, LexiconHash: hash}, GrammarVersion: grammar.Version, GrammarHash: hash, GrammarRuleID: rule.ID}
	if morphID != "" {
		query.MorphologyIDs = []string{morphID}
	}
	finding.QueryID = id
	return query, finding
}

func realizeAction(action, subject string, rules []domain.MorphologyRule) (string, string, bool) {
	if action == "" {
		return "", "", false
	}
	if subject == "imperative" || subject == "" {
		return action, "", true
	}
	for _, rule := range rules {
		if rule.SubjectKind == subject && strings.EqualFold(rule.Source, action) {
			return rule.Realized, rule.ID, true
		}
	}
	return "", "", false
}

func AuditFunctionalGrammar(expansion domain.FunctionalGrammarExpansion, grammar domain.FunctionalGrammar, forbidden []string) []domain.GrammarFinding {
	ruleIDs := map[string]bool{}
	for _, r := range grammar.Rules {
		ruleIDs[r.ID] = true
	}
	findings := make([]domain.GrammarFinding, 0, len(expansion.Queries))
	for _, q := range expansion.Queries {
		f := domain.GrammarFinding{QueryID: q.ID, GrammarVersion: q.GrammarVersion, RuleID: q.GrammarRuleID, Verdict: "ACCEPTED"}
		if q.GrammarHash != expansion.GrammarHash || q.GrammarHash != FunctionalGrammarHash(grammar) {
			f.Reasons = append(f.Reasons, "grammar hash mismatch")
		}
		if !ruleIDs[q.GrammarRuleID] {
			f.Reasons = append(f.Reasons, "unknown grammar rule")
		}
		if len(q.Actions) == 0 {
			f.MissingSlots = append(f.MissingSlots, "action")
		}
		if len(q.Objects) == 0 {
			f.MissingSlots = append(f.MissingSlots, "object")
		}
		if len(q.Inputs) == 0 && len(q.Outputs) == 0 {
			f.MissingSlots = append(f.MissingSlots, "input_or_output")
		}
		lower := strings.ToLower(q.Text)
		for _, term := range forbidden {
			if term != "" && strings.Contains(lower, strings.ToLower(term)) {
				f.ForbiddenHits = append(f.ForbiddenHits, term)
			}
		}
		if len(f.Reasons)+len(f.MissingSlots)+len(f.ForbiddenHits) > 0 {
			f.Verdict = "REJECTED"
		}
		findings = append(findings, f)
	}
	return findings
}

func firstValue(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
func stableGrammarQueryID(frame domain.FunctionalFrameID, rule, text, hash string) string {
	h := sha256.Sum256([]byte(string(frame) + "\x00" + rule + "\x00" + text + "\x00" + hash))
	return "fg-" + hex.EncodeToString(h[:6])
}
