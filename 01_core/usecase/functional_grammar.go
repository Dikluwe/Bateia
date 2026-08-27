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
	"strings"
)

func FunctionalGrammarHash(g domain.FunctionalGrammar) string {
	g.Hash = ""
	b, _ := json.Marshal(g)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func RealizeFunctionalGrammar(expansion domain.FunctionalExpansion, grammar domain.FunctionalGrammar) (domain.FunctionalExpansion, error) {
	if grammar.SchemaVersion == "" {
		return expansion, fmt.Errorf("gramática inválida")
	}
	grammar.Hash = FunctionalGrammarHash(grammar)

	morphMap := make(map[string]string)
	for _, m := range grammar.Morphology {
		morphMap[strings.ToLower(m.Source)] = m.Realized
	}

	realizedQueries := make([]domain.FunctionalQuery, len(expansion.Queries))
	for i, q := range expansion.Queries {
		text := q.Text
		for src, dst := range morphMap {
			text = strings.ReplaceAll(text, "tool that " + src, "tool that " + dst)
		}
		q.Text = text
		realizedQueries[i] = q
	}

	res := expansion
	res.Queries = realizedQueries
	res.LexiconHash = grammar.Hash
	return res, nil
}

func AuditFunctionalGrammar(expansion domain.FunctionalExpansion, grammar domain.FunctionalGrammar) ([]domain.GrammarFinding, error) {
	var findings []domain.GrammarFinding
	for _, q := range expansion.Queries {
		f := domain.GrammarFinding{
			QueryID:        q.ID,
			GrammarVersion: grammar.Version,
			Verdict:        "PASS",
		}

		if strings.Contains(q.Text, "tool that evaluate ") || strings.HasSuffix(q.Text, "tool that evaluate") || strings.Contains(q.Text, "tool that find ") || strings.HasSuffix(q.Text, "tool that find") {
			f.Verdict = "FAIL"
			f.UnresolvedVerbs = append(f.UnresolvedVerbs, "unrealized 3rd person verb")
			f.Reasons = append(f.Reasons, "found ungrammatical pattern tool that verb")
		}

		findings = append(findings, f)
	}

	return findings, nil
}
