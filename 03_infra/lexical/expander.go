/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/lexical/expander.md
 * @layer L3
 * @updated 2026-08-24
 */

package lexical

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

//go:embed vocabulary.json
var defaultVocabulary []byte

type Rule struct {
	ID       string   `json:"id"`
	WhenAny  []string `json:"when_any"`
	ExpandTo []string `json:"expand_to"`
}

type Vocabulary struct {
	Rules []Rule `json:"rules"`
}

type Expander struct {
	rules          []Rule
	maxPerSourceID int
}

func normalize(value string) string { return strings.ToLower(strings.Join(strings.Fields(value), " ")) }

func containsPhrase(text, phrase string) bool {
	text, phrase = normalize(text), normalize(phrase)
	return phrase != "" && strings.Contains(" "+text+" ", " "+phrase+" ")
}

func NewFromJSON(data []byte, maxPerSource int) (*Expander, error) {
	var vocabulary Vocabulary
	if err := json.Unmarshal(data, &vocabulary); err != nil {
		return nil, fmt.Errorf("vocabulário lexical JSON inválido: %w", err)
	}
	if len(vocabulary.Rules) == 0 {
		return nil, fmt.Errorf("vocabulário lexical deve conter ao menos uma regra")
	}
	if maxPerSource <= 0 {
		return nil, fmt.Errorf("limite por origem deve ser maior que zero")
	}
	seen := make(map[string]bool)
	for index := range vocabulary.Rules {
		rule := &vocabulary.Rules[index]
		rule.ID = strings.TrimSpace(rule.ID)
		if rule.ID == "" || seen[rule.ID] {
			return nil, fmt.Errorf("regra lexical %d possui id vazio ou duplicado", index)
		}
		seen[rule.ID] = true
		if len(rule.WhenAny) == 0 || len(rule.ExpandTo) == 0 {
			return nil, fmt.Errorf("regra lexical %s exige when_any e expand_to", rule.ID)
		}
		for _, value := range append(append([]string(nil), rule.WhenAny...), rule.ExpandTo...) {
			if normalize(value) == "" {
				return nil, fmt.Errorf("regra lexical %s possui termo vazio", rule.ID)
			}
		}
	}
	return &Expander{rules: vocabulary.Rules, maxPerSourceID: maxPerSource}, nil
}

func NewDefault(maxPerSource int) (*Expander, error) {
	return NewFromJSON(defaultVocabulary, maxPerSource)
}

func (e *Expander) Expand(plan domain.ExpansionPlan) ([]domain.SearchExpansion, error) {
	proposals := make([]domain.SearchExpansion, 0)
	for _, source := range plan.Expansions {
		count := 0
		for _, rule := range e.rules {
			matched := false
			for _, trigger := range rule.WhenAny {
				if containsPhrase(source.Query, trigger) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			for _, query := range rule.ExpandTo {
				if count >= e.maxPerSourceID {
					break
				}
				proposals = append(proposals, domain.SearchExpansion{
					Facet: source.Facet, Query: strings.Join(strings.Fields(query), " "),
					Rationale: "ampliar vocabulário por regra lexical " + rule.ID,
					Covers:    append([]string(nil), source.Covers...),
					Derivation: &domain.ExpansionDerivation{
						Strategy: "lexical", RuleIDs: []string{rule.ID}, SourceIDs: []string{source.ID},
					},
				})
				count++
			}
			if count >= e.maxPerSourceID {
				break
			}
		}
	}
	return proposals, nil
}
