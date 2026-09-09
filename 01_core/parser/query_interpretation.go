/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_interpretation.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"go/token"
	"sort"
	"strings"
	"unicode"

	"bateia/01_core/domain"
)

// InterpretImplementationQuery preserves constraints as unverified wording only.
func InterpretImplementationQuery(query string) []domain.ImplementationQuery {
	lines := map[string]bool{}
	for _, line := range strings.Split(query, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines[line] = true
		}
	}
	result := []domain.ImplementationQuery{}
	for _, line := range sortedQueryValues(lines) {
		interpretation := domain.ImplementationQuery{Original: line, Mode: "natural"}
		terms, symbols, markers := map[string]bool{}, map[string]bool{}, map[string]bool{}
		expansionSources := map[string]bool{}
		natural := false
		addTerms := func(text string) {
			for term := range implementationTerms(text) {
				terms[term] = true
			}
		}
		addSymbol := func(symbol string) {
			symbols[symbol] = true
			addTerms(symbol)
			for _, term := range implementationWords(symbol) {
				terms[term] = true
			}
		}
		addNatural := func(text string, expand bool) {
			words := strings.FieldsFunc(text, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
			})
			for _, word := range words {
				marker := strings.ToLower(word)
				switch marker {
				case "sem", "não", "nao", "nunca", "not", "without", "never":
					markers[marker] = true
					natural = true
					continue
				}
				explicit := strings.Contains(word, "_")
				for i, r := range word {
					if i > 0 && unicode.IsUpper(r) {
						explicit = true
					}
				}
				if explicit && token.IsIdentifier(word) {
					addSymbol(word)
				} else {
					natural = true
					addTerms(word)
					if expand {
						expansionSources[marker] = true
					}
				}
			}
		}
		if token.IsIdentifier(line) {
			addSymbol(line)
		} else {
			// Only single, balanced delimiters around a Go identifier grant explicitness.
			rest := line
			for rest != "" {
				open := strings.IndexByte(rest, '`')
				if open < 0 {
					addNatural(rest, true)
					break
				}
				addNatural(rest[:open], true)
				rest = rest[open:]
				opening := len(rest) - len(strings.TrimLeft(rest, "`"))
				close := strings.IndexByte(rest[opening:], '`')
				if close < 0 {
					addNatural(rest, false)
					break
				}
				close += opening
				closing := len(rest[close:]) - len(strings.TrimLeft(rest[close:], "`"))
				content := rest[opening:close]
				if opening == 1 && closing == 1 && token.IsIdentifier(content) {
					addSymbol(content)
				} else {
					addNatural(content, false)
				}
				rest = rest[close+closing:]
			}
		}
		if len(symbols) > 0 {
			interpretation.Mode = "symbol"
			if natural {
				interpretation.Mode = "mixed"
			}
		}
		if len(markers) == 0 {
			for _, source := range sortedQueryValues(expansionSources) {
				for _, expansion := range implementationQueryExpansions(source) {
					interpretation.Expansions = append(interpretation.Expansions, expansion)
					for _, term := range expansion.Terms {
						terms[term] = true
					}
				}
			}
		}
		if interpretation.Mode == "natural" {
			plan := PlanDiscoveryQuery(line)
			if plan.Status == "hypothesis" {
				for i, group := range [][]string{plan.ActionTerms, plan.ObjectTerms} {
					interpretation.Expansions = append(interpretation.Expansions, domain.ImplementationQueryExpansion{
						Source: line, Terms: append([]string(nil), group...),
						RuleID: "0058-v1:" + plan.Version + ":" + plan.RuleIDs[i],
						Kind:   "translation", Status: "hypothesis",
					})
					for _, term := range group {
						terms[term] = true
					}
				}
			}
		}
		if interpretation.Mode == "natural" && len(markers) == 0 {
			for _, source := range sortedQueryValues(expansionSources) {
				if source != "caminho" && source != "caminhos" {
					continue
				}
				interpretation.Expansions = append(interpretation.Expansions, domain.ImplementationQueryExpansion{
					Source: source, Terms: []string{"path"}, RuleID: "0060-v1:path:translation",
					Kind: "translation", Status: "hypothesis",
				})
				terms["path"] = true
			}
		}
		interpretation.Terms = sortedQueryValues(terms)
		interpretation.Symbols = sortedQueryValues(symbols)
		for _, marker := range sortedQueryValues(markers) {
			interpretation.Constraints = append(interpretation.Constraints, domain.ImplementationQueryConstraint{
				Marker: marker, Text: line, Status: "not_verified",
			})
		}
		result = append(result, interpretation)
	}
	return result
}

// implementationQueryExpansions applies the finite 0052 catalog only to original words.
func implementationQueryExpansions(source string) []domain.ImplementationQueryExpansion {
	var family string
	var translation, mechanism []string
	switch source {
	case "descomprimir", "descompressão", "descompressao":
		family, translation, mechanism = "decompression", []string{"decompress", "decompression"}, []string{"decode", "decoded", "decompresses"}
	case "comprimir", "compressão", "compressao":
		family, translation, mechanism = "compression", []string{"compress", "compression"}, []string{"encode", "encoded", "compresses"}
	case "serializar", "serialização", "serializacao":
		family, translation, mechanism = "serialization", []string{"serialize", "serialization"}, []string{"marshal"}
	case "desserializar", "desserialização", "desserializacao":
		family, translation, mechanism = "deserialization", []string{"deserialize", "deserialization"}, []string{"unmarshal"}
	case "ordenar", "ordenação", "ordenacao":
		family, translation = "sorting", []string{"sort", "sorting"}
	case "comparar", "comparação", "comparacao":
		family, translation = "comparison", []string{"compare", "comparison"}
	case "bloco", "blocos":
		family, translation = "block", []string{"block", "blocks"}
	case "corrompido", "corrompidos", "corrompida", "corrompidas":
		family, translation = "corrupt", []string{"corrupt", "corrupted"}
	case "dados":
		family, translation = "data", []string{"data"}
	case "destino":
		family, translation, mechanism = "destination", []string{"destination"}, []string{"dst"}
	case "reutilizar":
		family, translation = "reuse", []string{"reuse"}
	case "rejeitar":
		family, translation = "reject", []string{"reject"}
	default:
		return nil
	}
	var result []domain.ImplementationQueryExpansion
	for i, derived := range [][]string{translation, mechanism} {
		if len(derived) == 0 {
			continue
		}
		kind := []string{"translation", "mechanism"}[i]
		sort.Strings(derived)
		result = append(result, domain.ImplementationQueryExpansion{
			Source: source, Terms: derived, RuleID: "0052-v1:" + family + ":" + kind,
			Kind: kind, Status: "hypothesis",
		})
	}
	return result
}

// implementationQueryActionAdmission checks one query line against one isolated field.
// The second return value reports whether action-family filtering was applicable.
func implementationQueryActionAdmission(query domain.ImplementationQuery, field map[string]bool) (bool, bool) {
	context := false
	for _, term := range query.Terms {
		if (term == "block" || term == "blocks" || term == "corrupt" || term == "corrupted") && field[term] {
			context = true
		}
	}
	families := map[string]bool{}
	for _, expansion := range query.Expansions {
		parts := strings.Split(expansion.RuleID, ":")
		if len(parts) != 3 || parts[0] != "0052-v1" {
			continue
		}
		family := parts[1]
		switch family {
		case "compression", "decompression", "serialization", "deserialization", "sorting", "comparison":
		default:
			continue
		}
		if _, present := families[family]; !present {
			families[family] = false
		}
		if expansion.Kind != "translation" && expansion.Kind != "mechanism" {
			continue
		}
		if expansion.Kind == "mechanism" && (family == "compression" || family == "decompression") && !context {
			continue
		}
		for _, term := range expansion.Terms {
			if field[term] {
				families[family] = true
			}
		}
	}
	for _, satisfied := range families {
		if !satisfied {
			return false, true
		}
	}
	return true, len(families) > 0
}

func sortedQueryValues(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
