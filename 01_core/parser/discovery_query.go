/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/discovery_query.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser

import (
	"bateia/01_core/domain"
	"strings"
	"unicode"
)

type discoveryFamily struct {
	id              string
	triggers, terms []string
}

var discoveryActions = []discoveryFamily{
	{"action.watch", []string{"monitorar", "monitoramento", "vigiar", "watch"}, []string{"watch", "notify"}},
	{"action.interpolate", []string{"interpolar", "interpolacao", "interpolation", "blend"}, []string{"interpolat", "blend"}},
	{"action.intersection", []string{"intersecao", "intersection", "intersect", "presentes em ambos"}, []string{"intersect", "intersection"}},
	{"action.union", []string{"uniao", "union", "presentes em qualquer"}, []string{"union"}},
	{"action.compress", []string{"comprimir", "compressao", "compress"}, []string{"compress"}},
	{"action.serialize", []string{"serializar", "serializacao", "serialize"}, []string{"serializ", "marshal"}},
}
var discoveryObjects = []discoveryFamily{
	{"object.file", []string{"arquivo", "arquivos", "file", "files", "diretorio", "diretorios", "directory"}, []string{"file", "directory"}},
	{"object.color", []string{"cor", "cores", "color", "colours"}, []string{"color", "colour"}},
	{"object.bit", []string{"bit", "bits", "bitset"}, []string{"bit", "bitset"}},
	{"object.data", []string{"dado", "dados", "data"}, []string{"data"}},
}

func discoveryNormalized(query string) string {
	query = strings.ToLower(query)
	query = strings.Map(func(r rune) rune {
		switch r {
		case 'á', 'à', 'â', 'ã', 'ä':
			return 'a'
		case 'é', 'è', 'ê', 'ë':
			return 'e'
		case 'í', 'ì', 'î', 'ï':
			return 'i'
		case 'ó', 'ò', 'ô', 'õ', 'ö':
			return 'o'
		case 'ú', 'ù', 'û', 'ü':
			return 'u'
		case 'ç':
			return 'c'
		}
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, query)
	return " " + strings.Join(strings.FieldsFunc(query, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' }), " ") + " "
}

func discoveryRecognized(text string, families []discoveryFamily) []discoveryFamily {
	var found []discoveryFamily
	for _, family := range families {
		for _, trigger := range family.triggers {
			if strings.Contains(text, " "+trigger+" ") {
				found = append(found, family)
				break
			}
		}
	}
	return found
}

// PlanDiscoveryQuery applies only the finite 0055-v2 catalog.
func PlanDiscoveryQuery(query string) domain.DiscoveryQueryPlan {
	plan := domain.DiscoveryQueryPlan{Original: query, Version: "0055-v2", Status: "literal_only", Reason: "ação e objeto reconhecidos insuficientes", Unverified: []string{query}}
	if strings.ContainsAny(query, "\r\n\u0085\u2028\u2029") {
		plan.Reason = "consulta multilinha"
		return plan
	}
	// Reject code/query syntax conservatively instead of guessing its meaning.
	if strings.ContainsAny(query, "\"'`:~=<>!|&+*/\\^$()[]{};%-") {
		plan.Reason = "sintaxe explícita"
		return plan
	}
	normalized := discoveryNormalized(query)
	for _, negation := range []string{"sem", "nao", "nunca", "not", "without", "never"} {
		if strings.Contains(normalized, " "+negation+" ") {
			plan.Reason = "negação não interpretada"
			return plan
		}
	}
	for _, operator := range []string{"AND", "OR", "NOT"} {
		for _, word := range strings.Fields(query) {
			if word == operator {
				plan.Reason = "operador de consulta"
				return plan
			}
		}
	}
	actions, objects := discoveryRecognized(normalized, discoveryActions), discoveryRecognized(normalized, discoveryObjects)
	if len(actions) > 1 || len(objects) > 1 {
		plan.Reason = "famílias ambíguas"
		return plan
	}
	if len(actions) != 1 || len(objects) != 1 {
		return plan
	}
	plan.Status = "hypothesis"
	plan.Reason = "hipótese lexical de ação e objeto; implementação e restrições não verificadas"
	plan.ActionTerms = append([]string(nil), actions[0].terms...)
	plan.ObjectTerms = append([]string(nil), objects[0].terms...)
	plan.RuleIDs = []string{actions[0].id, objects[0].id}
	return plan
}

// DiscoveryEvidenceMatches checks only supplied retained text, not implementation semantics.
func DiscoveryEvidenceMatches(plan domain.DiscoveryQueryPlan, text string) bool {
	if plan.Status != "hypothesis" {
		return false
	}
	text = strings.ToLower(text)
	contains := func(terms []string) bool {
		for _, term := range terms {
			if term != "" && strings.Contains(text, strings.ToLower(term)) {
				return true
			}
		}
		return false
	}
	return contains(plan.ActionTerms) && contains(plan.ObjectTerms)
}
