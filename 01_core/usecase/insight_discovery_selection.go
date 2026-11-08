/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery_selection.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"sort"
	"strings"
)

func validInsightCatalogCandidates(candidates []domain.CatalogCandidate, limit int) bool {
	if len(candidates) > limit {
		return false
	}
	for _, candidate := range candidates {
		r := candidate.Repository
		if !validSynthesisText(r.Owner, 100) || !validSynthesisText(r.Repo, 100) || r.Owner != normalizeMechanismText(r.Owner) || r.Repo != normalizeMechanismText(r.Repo) || strings.ContainsAny(r.Owner+r.Repo, "/\\@ :?#") || r.Owner == "." || r.Owner == ".." || r.Repo == "." || r.Repo == ".." || len(r.Description) > 4000 || len(r.URL) > 2048 || len(candidate.Topics) > 32 {
			return false
		}
		for _, topic := range candidate.Topics {
			if !validSynthesisText(topic, 128) {
				return false
			}
		}
		if r.CommitID != "" && r.CommitID != "HEAD" && !immutableCommitPattern.MatchString(r.CommitID) {
			return false
		}
	}
	return true
}

func mergeInsightCatalogCandidates(existing []domain.InsightDiscoveredReference, candidates []domain.CatalogCandidate, queryID string, familiar []string) []domain.InsightDiscoveredReference {
	byKey := map[string]domain.InsightDiscoveredReference{}
	for _, ref := range existing {
		byKey[repositoryKey(ref.Repository)] = ref
	}
	for _, candidate := range candidates {
		repo := candidate.Repository
		repo.SearchAttempts = nil
		key := repositoryKey(repo)
		ref, found := byKey[key]
		if !found {
			ref = domain.InsightDiscoveredReference{ID: "reference_" + insightHash(key), Repository: repo, Topics: []string{}, QueryIDs: []string{}, FamiliarSubstringMatches: []string{}, ObservedCommits: []string{}, SelectionReasons: []string{}}
		}
		if immutableCommitPattern.MatchString(repo.CommitID) {
			ref.ObservedCommits = canonicalMechanismSet(append(append([]string{}, ref.ObservedCommits...), repo.CommitID))
		}
		ref.RevisionConflict = len(ref.ObservedCommits) > 1
		if string(canonicalJSON(repo)) < string(canonicalJSON(ref.Repository)) {
			ref.Repository = repo
		}
		if len(ref.ObservedCommits) == 1 {
			ref.Repository.CommitID = ref.ObservedCommits[0]
		}
		ref.Topics = canonicalMechanismSet(append(append([]string{}, ref.Topics...), candidate.Topics...))
		ref.QueryIDs = canonicalMechanismSet(append(append([]string{}, ref.QueryIDs...), queryID))
		byKey[key] = ref
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := []domain.InsightDiscoveredReference{}
	for _, key := range keys {
		ref := byKey[key]
		haystack := strings.ToLower(strings.Join(strings.Fields(ref.Repository.Repo+" "+ref.Repository.Description+" "+strings.Join(ref.Topics, " ")), " "))
		ref.FamiliarSubstringMatches = []string{}
		for _, label := range familiar {
			if strings.Contains(haystack, strings.ToLower(label)) {
				ref.FamiliarSubstringMatches = append(ref.FamiliarSubstringMatches, label)
			}
		}
		out = append(out, ref)
	}
	return out
}

func selectInsightDiscoveryReferences(candidates []domain.InsightDiscoveredReference, maximum int, question ...string) []domain.InsightDiscoveredReference {
	out := append([]domain.InsightDiscoveredReference{}, candidates...)
	queries, owners := map[string]bool{}, map[string]bool{}
	for count := 0; count < maximum; count++ {
		best, score := -1, [5]int{}
		for i, candidate := range out {
			if candidate.Selected {
				continue
			}
			gain := 0
			for _, q := range candidate.QueryIDs {
				if !queries[q] {
					gain++
				}
			}
			current := [5]int{gain, 0, -len(candidate.FamiliarSubstringMatches)}
			if len(question) > 0 && question[0] != "" {
				matched := insightMetadataMatchCount(question[0], candidate.Repository.Repo+" "+candidate.Repository.Description+" "+strings.Join(candidate.Topics, " "))
				if matched == 0 {
					continue
				}
				current = [5]int{matched, gain, 0, -len(candidate.FamiliarSubstringMatches)}
				if !owners[strings.ToLower(candidate.Repository.Owner)] {
					current[2] = 1
				}
			}
			if (len(question) == 0 || question[0] == "") && !owners[strings.ToLower(candidate.Repository.Owner)] {
				current[1] = 1
			}
			if best < 0 || explorationScoreGreater(current, score) {
				best, score = i, current
			}
		}
		if best < 0 {
			break
		}
		ref := &out[best]
		ref.Selected = true
		ref.SelectionReasons = append([]string{}, ref.SelectionReasons...)
		if len(question) > 0 && question[0] != "" {
			ref.SelectionReasons = append(ref.SelectionReasons, "LITERAL_METADATA_QUESTION_MATCH")
		}
		queryGain, ownerGain := score[0] > 0, score[1] > 0
		if len(question) > 0 && question[0] != "" {
			queryGain = score[1] > 0
			ownerGain = score[2] > 0
		}
		if queryGain {
			ref.SelectionReasons = append(ref.SelectionReasons, "NEW_QUERY_COVERAGE")
		}
		if ownerGain {
			ref.SelectionReasons = append(ref.SelectionReasons, "NEW_OWNER")
		}
		for _, q := range ref.QueryIDs {
			queries[q] = true
		}
		owners[strings.ToLower(ref.Repository.Owner)] = true
	}
	if len(question) > 0 && question[0] != "" {
		for i := range out {
			if !out[i].Selected && insightMetadataMatchCount(question[0], out[i].Repository.Repo+" "+out[i].Repository.Description+" "+strings.Join(out[i].Topics, " ")) == 0 {
				out[i].SelectionReasons = append(append([]string{}, out[i].SelectionReasons...), "NO_METADATA_QUESTION_MATCH")
			}
		}
	}
	return out
}

func insightMetadataMatchCount(question, metadata string) int {
	words := map[string]bool{}
	for _, word := range explorationWords(metadata) {
		words[word] = true
	}
	count := 0
	for _, term := range explorationTerms(question) {
		if words[term] {
			count++
		}
	}
	return count
}
