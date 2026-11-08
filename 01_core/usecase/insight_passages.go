/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_passages.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"strings"
	"unicode"
	"unicode/utf8"
)

func explorationTerms(question string) []string {
	stop := " a an and are as at be by can com como could da das de do dos em for from how in is it no nos of on or os para por que se sem should the to um uma what when with would "
	out, seen := []string{}, map[string]bool{}
	for _, term := range explorationWords(question) {
		if utf8.RuneCountInString(term) < 3 || strings.Contains(stop, " "+term+" ") || seen[term] {
			continue
		}
		seen[term] = true
		out = append(out, term)
		if len(out) == 24 {
			break
		}
	}
	return out
}

func explorationWords(value string) []string {
	return strings.FieldsFunc(strings.ToLower(value), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}

func selectInsightPassages(ctx context.Context, request domain.InsightExplorationRequest, terms []string) ([]domain.InsightReadingCard, domain.InsightExplorationDiagnostics, error) {
	diagnostics := domain.InsightExplorationDiagnostics{SuppliedSources: len(request.Bundles)}
	candidates := []domain.InsightReadingCard{}
	for _, bundle := range request.Bundles {
		if !bundle.Complete || len(bundle.Failures) > 0 {
			diagnostics.IncompleteBundles++
		}
		for _, item := range bundle.Items {
			if err := ctx.Err(); err != nil {
				return nil, diagnostics, err
			}
			if !explorationReadable(item) {
				diagnostics.IgnoredItems++
				continue
			}
			lines := observationLines(item.Content)
			if len(lines) > 1024 {
				diagnostics.PrefixLimitedFiles++
				lines = lines[:1024]
			}
			best := domain.InsightReadingCard{}
			for start := 0; start < len(lines); start += 8 {
				left, right := start, min(start+16, len(lines))
				for left < right && strings.TrimSpace(lines[left]) == "" {
					left++
				}
				for right > left && (strings.TrimSpace(lines[right-1]) == "" || len(strings.Join(lines[left:right], "\n")) > 2048) {
					right--
				}
				if left == right {
					continue
				}
				excerpt := strings.Join(lines[left:right], "\n")
				card := domain.InsightReadingCard{Source: immutableRef(bundle.Repository), Kind: item.Kind, Origin: "UNINTERPRETED_EXCERPT",
					Reference:            domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: item.Path, ContentHash: item.ContentHash, StartLine: left + 1, EndLine: right}, Excerpt: excerpt},
					MatchedQuestionTerms: []string{}, FamiliarSubstringMatches: []string{}, SelectionReasons: []string{},
				}
				words := map[string]bool{}
				for _, word := range explorationWords(excerpt) {
					words[word] = true
				}
				for word := range words {
					if utf8.RuneCountInString(word) >= 3 {
						card.DistinctExcerptWords++
					}
				}
				for _, term := range terms {
					if words[term] {
						card.MatchedQuestionTerms = append(card.MatchedQuestionTerms, term)
					}
				}
				folded := strings.ToLower(strings.Join(strings.Fields(excerpt), " "))
				for _, familiar := range request.FamiliarApproaches {
					if strings.Contains(folded, strings.ToLower(familiar)) {
						card.FamiliarSubstringMatches = append(card.FamiliarSubstringMatches, familiar)
					}
				}
				if best.ID == "" || explorationScoreGreater([5]int{len(card.MatchedQuestionTerms), -len(card.FamiliarSubstringMatches), card.DistinctExcerptWords}, [5]int{len(best.MatchedQuestionTerms), -len(best.FamiliarSubstringMatches), best.DistinctExcerptWords}) {
					card.ID = "card_" + insightHash(struct {
						Source    domain.ImmutableCandidateRef
						Kind      domain.EvidenceKind
						Reference domain.LocatedInsightReference
					}{card.Source, card.Kind, card.Reference})
					best = card
				}
			}
			if best.ID == "" {
				diagnostics.IgnoredItems++
			} else {
				candidates = append(candidates, best)
			}
		}
	}
	diagnostics.EligibleFiles = len(candidates)
	unique := map[string]bool{}
	for _, card := range candidates {
		unique[contentHash(card.Reference.Excerpt)] = true
	}
	diagnostics.DuplicatePassages = len(candidates) - len(unique)
	diagnostics.SelectionLimited = len(unique) > request.MaxPassages
	out := []domain.InsightReadingCard{}
	repositories, kinds, excerpts := map[string]bool{}, map[domain.EvidenceKind]bool{}, map[string]bool{}
	for len(out) < request.MaxPassages {
		if err := ctx.Err(); err != nil {
			return nil, diagnostics, err
		}
		bestIndex, bestScore := -1, [5]int{}
		for i, card := range candidates {
			if excerpts[contentHash(card.Reference.Excerpt)] {
				continue
			}
			score := [5]int{0, 0, len(card.MatchedQuestionTerms), -len(card.FamiliarSubstringMatches), card.DistinctExcerptWords}
			if !repositories[strings.ToLower(card.Source.Repository)] {
				score[0] = 1
			}
			if !kinds[card.Kind] {
				score[1] = 1
			}
			if bestIndex < 0 || explorationScoreGreater(score, bestScore) {
				bestIndex, bestScore = i, score
			}
		}
		if bestIndex < 0 {
			break
		}
		card := candidates[bestIndex]
		if bestScore[0] > 0 {
			card.SelectionReasons = append(card.SelectionReasons, "NEW_REPOSITORY")
		}
		if bestScore[1] > 0 {
			card.SelectionReasons = append(card.SelectionReasons, "NEW_EVIDENCE_KIND")
		}
		if len(card.MatchedQuestionTerms) == 0 {
			card.SelectionReasons = append(card.SelectionReasons, "NO_QUESTION_TERM_MATCH")
		} else {
			card.SelectionReasons = append(card.SelectionReasons, "LITERAL_QUESTION_TERM_MATCH")
		}
		repositories[strings.ToLower(card.Source.Repository)], kinds[card.Kind], excerpts[contentHash(card.Reference.Excerpt)] = true, true, true
		out = append(out, card)
	}
	if request.ReviewAlternatives {
		for _, bundle := range request.Bundles {
			for _, item := range bundle.Items {
				if !explorationReadable(item) {
					continue
				}
				lines := observationLines(item.Content)
				lines = lines[:min(len(lines), 1024)]
				for count := 1; count < 3; count++ {
					best := domain.InsightReadingCard{}
					for start := 0; start < len(lines); start += 8 {
						if err := ctx.Err(); err != nil {
							return nil, diagnostics, err
						}
						end := min(start+16, len(lines))
						overlap := false
						for _, used := range out {
							if used.Source.Repository == bundle.Repository.Owner+"/"+bundle.Repository.Repo && used.Source.Commit == bundle.Repository.Commit && used.Reference.Path == item.Path && start+1 <= used.Reference.EndLine && end >= used.Reference.StartLine {
								overlap = true
								break
							}
						}
						if overlap {
							continue
						}
						excerpt := strings.Join(lines[start:end], "\n")
						if strings.TrimSpace(excerpt) == "" || len(excerpt) > 2048 || excerpts[contentHash(excerpt)] {
							continue
						}
						card := domain.InsightReadingCard{Source: immutableRef(bundle.Repository), Kind: item.Kind, Origin: "UNINTERPRETED_EXCERPT", Reference: domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: item.Path, ContentHash: item.ContentHash, StartLine: start + 1, EndLine: end}, Excerpt: excerpt}, MatchedQuestionTerms: []string{}, FamiliarSubstringMatches: []string{}, SelectionReasons: []string{"ADDITIONAL_FILE_PASSAGE"}}
						words := map[string]bool{}
						for _, word := range explorationWords(excerpt) {
							words[word] = true
						}
						for word := range words {
							if utf8.RuneCountInString(word) >= 3 {
								card.DistinctExcerptWords++
							}
						}
						for _, term := range terms {
							if words[term] {
								card.MatchedQuestionTerms = append(card.MatchedQuestionTerms, term)
							}
						}
						for _, familiar := range request.FamiliarApproaches {
							if strings.Contains(strings.ToLower(strings.Join(strings.Fields(excerpt), " ")), strings.ToLower(familiar)) {
								card.FamiliarSubstringMatches = append(card.FamiliarSubstringMatches, familiar)
							}
						}
						if best.ID != "" && !explorationScoreGreater([5]int{len(card.MatchedQuestionTerms), -len(card.FamiliarSubstringMatches), card.DistinctExcerptWords}, [5]int{len(best.MatchedQuestionTerms), -len(best.FamiliarSubstringMatches), best.DistinctExcerptWords}) {
							continue
						}
						card.ID = "card_" + insightHash(struct {
							Source    domain.ImmutableCandidateRef
							Kind      domain.EvidenceKind
							Reference domain.LocatedInsightReference
						}{card.Source, card.Kind, card.Reference})
						best = card
					}
					if best.ID == "" {
						break
					}
					// Additional windows only expand files already represented in the first pass.
					represented := false
					for _, used := range out {
						if used.Source == best.Source && used.Reference.Path == item.Path {
							represented = true
							break
						}
					}
					if !represented {
						break
					}
					out = append(out, best)
					excerpts[contentHash(best.Reference.Excerpt)] = true
				}
			}
		}
		if len(out) > request.MaxPassages {
			diagnostics.SelectionLimited = true
			out = out[:request.MaxPassages]
		}
	}
	diagnostics.SelectedSources = len(repositories)
	return out, diagnostics, nil
}

func explorationScoreGreater(left, right [5]int) bool {
	for i := range left {
		if left[i] != right[i] {
			return left[i] > right[i]
		}
	}
	return false
}

func explorationReadable(item domain.CollectedEvidence) bool {
	switch item.Kind {
	case domain.EvidenceFile, domain.EvidenceSymbol, domain.EvidenceDocumentation, domain.EvidenceTest, domain.EvidenceTestImplementation, domain.EvidenceExample:
	default:
		return false
	}
	if !utf8.ValidString(item.Content) {
		return false
	}
	for _, r := range item.Content {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}
