/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/corpus_scoring.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"fmt"
	"math"
	"sort"

	"bateia/01_core/domain"
)

type corpusDocument struct {
	hit        domain.ImplementationHit
	name, body map[string]bool
	functional bool
	exact      bool
}

type corpusStatistics struct {
	documents   int
	frequency   map[string]int
	averageBody float64
}

func collectCorpusStatistics(documents []corpusDocument) corpusStatistics {
	stats := corpusStatistics{documents: len(documents), frequency: map[string]int{}, averageBody: 1}
	bodies, length := 0, 0
	for _, doc := range documents {
		terms := map[string]bool{}
		for term := range doc.name {
			terms[term] = true
		}
		for term := range doc.body {
			terms[term] = true
		}
		for term := range terms {
			stats.frequency[term]++
		}
		if doc.functional {
			bodies++
			length += len(doc.body)
		}
	}
	if bodies > 0 && length > bodies {
		stats.averageBody = float64(length) / float64(bodies)
	}
	return stats
}

// score uses the same corpus for global ranking and each individual requirement.
func (stats corpusStatistics) score(required map[string]bool, doc corpusDocument) (int, []string, []domain.ImplementationTermWeight) {
	terms := make([]string, 0, len(required))
	for term := range required {
		terms = append(terms, term)
	}
	sort.Strings(terms)
	matched := []string{}
	weights := []domain.ImplementationTermWeight{}
	sum := 0.0
	norm := 0.25 + 0.75*float64(len(doc.body))/stats.averageBody
	for _, term := range terms {
		df := stats.frequency[term]
		if df == 0 || (!doc.name[term] && !doc.body[term]) {
			continue
		}
		idf := math.Log1p((float64(stats.documents-df) + 0.5) / (float64(df) + 0.5))
		matched = append(matched, term)
		// Fields are emitted in lexical order as well, so float summation is stable.
		if doc.body[term] {
			weight := 3 * idf / norm
			sum += weight
			weights = append(weights, domain.ImplementationTermWeight{Term: term, Field: "body", DocumentFrequency: df, CorpusDocuments: stats.documents, Weight: weight})
		}
		if doc.name[term] {
			weight := 2 * idf
			sum += weight
			weights = append(weights, domain.ImplementationTermWeight{Term: term, Field: "name", DocumentFrequency: df, CorpusDocuments: stats.documents, Weight: weight})
		}
	}
	score := 0
	if sum > 0 {
		score = int(math.Ceil(900 * sum / (8 + sum)))
		if doc.functional {
			score += 1000
		}
	}
	return score, matched, weights
}

func (stats corpusStatistics) reason(doc corpusDocument) string {
	return fmt.Sprintf("go_ast_corpus_v1: IDF=log1p((N-df+0.5)/(df+0.5)); peso name=2*IDF, body=3*IDF/(0.25+0.75*len/avg); Score=ceil(900*soma/(8+soma))+1000 se função pertinente+100000 se exato; N=%d declarações nomeadas AST elegíveis, incluindo não matches, df união nome/corpo; len=%d termos únicos AST do corpo, avg=%.6g (média funções com corpo, piso 1); recorte lido, não ecossistema; expansão lexical PT/EN, não probabilidade nem comprovação de comportamento", stats.documents, len(doc.body), stats.averageBody)
}
