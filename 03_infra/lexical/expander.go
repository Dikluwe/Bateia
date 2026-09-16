/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/lexical/expander.md
 * @layer L3
 * @updated 2026-09-11
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
	ID             string   `json:"id"`
	WhenAny        []string `json:"when_any"`
	ExpandTo       []string `json:"expand_to"`
	Strategy       string   `json:"strategy,omitempty"`
	Discriminators []string `json:"discriminators,omitempty"`
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
		rule.Strategy = strings.TrimSpace(rule.Strategy)
		if rule.Strategy == "" {
			rule.Strategy = "lexical"
		}
		switch rule.Strategy {
		case "lexical":
			if len(rule.Discriminators) != 0 {
				return nil, fmt.Errorf("regra lexical %s não aceita discriminadores", rule.ID)
			}
		case "technical-lexicon-v1":
			if len(rule.ExpandTo) != 1 || len(rule.Discriminators) != 2 {
				return nil, fmt.Errorf("regra técnica %s exige uma expansão e dois discriminadores", rule.ID)
			}
			if len(rule.ExpandTo[0]) > 256 {
				return nil, fmt.Errorf("regra técnica %s excede 256 bytes por consulta", rule.ID)
			}
			for discriminatorIndex := range rule.Discriminators {
				rule.Discriminators[discriminatorIndex] = strings.TrimSpace(rule.Discriminators[discriminatorIndex])
				if rule.Discriminators[discriminatorIndex] == "" || len(rule.Discriminators[discriminatorIndex]) > 64 {
					return nil, fmt.Errorf("regra técnica %s possui discriminador vazio ou acima de 64 bytes", rule.ID)
				}
			}
			if normalize(rule.Discriminators[0]) == normalize(rule.Discriminators[1]) {
				return nil, fmt.Errorf("regra técnica %s possui discriminadores duplicados", rule.ID)
			}
		default:
			return nil, fmt.Errorf("regra lexical %s possui estratégia desconhecida %q", rule.ID, rule.Strategy)
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
					Rationale:      "ampliar vocabulário por regra lexical " + rule.ID,
					Covers:         append([]string(nil), source.Covers...),
					Discriminators: append([]string(nil), rule.Discriminators...),
					Derivation: &domain.ExpansionDerivation{
						Strategy: rule.Strategy, RuleIDs: []string{rule.ID}, SourceIDs: []string{source.ID},
					},
				})
				count++
			}
			if count >= e.maxPerSourceID {
				break
			}
		}
	}
	if len(plan.Expansions)+len(proposals) > 64 {
		return nil, fmt.Errorf("plano enriquecido excede o limite de 64 expansões")
	}
	return proposals, nil
}
