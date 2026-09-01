/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/assisted_judgment.md
 * @layer L1
 * @updated 2026-08-31
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"sort"
	"strings"
)

const AssistedRubricVersion = "step0021-v1"

func PrepareBlindCandidateUnion(snapshots []domain.WebSearchSnapshot) (domain.BlindCandidatePackage, error) {
	byID := map[string]domain.BlindCandidate{}
	for _, snapshot := range snapshots {
		for _, response := range snapshot.Responses {
			for _, hit := range response.Results {
				canonical := CanonicalizeURL(hit.CanonicalURL)
				if canonical == "" {
					canonical = CanonicalizeURL(hit.URL)
				}
				if canonical == "" {
					continue
				}
				h := sha256.Sum256([]byte(canonical))
				id := "candidate-" + hex.EncodeToString(h[:12])
				candidate := domain.BlindCandidate{CandidateID: id, Fragments: []domain.BlindEvidenceFragment{{FragmentID: "title", Content: strings.TrimSpace(hit.Title)}, {FragmentID: "description", Content: strings.TrimSpace(hit.Snippet)}}}
				byID[id] = candidate
			}
		}
	}
	out := domain.BlindCandidatePackage{SchemaVersion: "1", RubricVersion: AssistedRubricVersion}
	for _, candidate := range byID {
		out.Candidates = append(out.Candidates, candidate)
	}
	sort.Slice(out.Candidates, func(i, j int) bool { return out.Candidates[i].CandidateID < out.Candidates[j].CandidateID })
	out.PackageHash = blindPackageHash(out)
	return out, nil
}

func blindPackageHash(pkg domain.BlindCandidatePackage) string {
	pkg.PackageHash = ""
	b, _ := json.Marshal(pkg)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func DetermineHumanEscalations(pkg domain.BlindCandidatePackage, judgments []domain.AssistedJudgment) []string {
	candidateByID := map[string]domain.BlindCandidate{}
	for _, candidate := range pkg.Candidates {
		candidateByID[candidate.CandidateID] = candidate
	}
	mandatory := map[string]bool{}
	judgmentByID := map[string]domain.AssistedJudgment{}
	for _, judgment := range judgments {
		judgmentByID[judgment.CandidateID] = judgment
		if judgment.Confidence < .8 || math.IsNaN(judgment.Confidence) || math.IsInf(judgment.Confidence, 0) || judgment.Confidence > 1 || judgment.Verdict == "UNCERTAIN" || judgment.Verdict == "CONFLICT" || !judgmentEvidenceValid(candidateByID[judgment.CandidateID], judgment) {
			mandatory[judgment.CandidateID] = true
		}
	}
	var control []string
	for _, candidate := range pkg.Candidates {
		if !mandatory[candidate.CandidateID] {
			control = append(control, candidate.CandidateID)
		}
	}
	sort.Strings(control)
	count := int(math.Ceil(float64(len(control)) * .1))
	for i := 0; i < count; i++ {
		mandatory[control[i]] = true
	}
	out := make([]string, 0, len(mandatory))
	for id := range mandatory {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func judgmentEvidenceValid(candidate domain.BlindCandidate, judgment domain.AssistedJudgment) bool {
	if candidate.CandidateID == "" || len(judgment.Evidence) == 0 || judgment.RubricVersion != AssistedRubricVersion {
		return false
	}
	fragments := map[string]bool{}
	for _, fragment := range candidate.Fragments {
		fragments[fragment.FragmentID] = strings.TrimSpace(fragment.Content) != ""
	}
	for _, pointer := range judgment.Evidence {
		if pointer.CandidateID != candidate.CandidateID || !fragments[pointer.FragmentID] {
			return false
		}
	}
	return validAssistedVerdict(judgment.Verdict) && judgment.Confidence >= 0 && judgment.Confidence <= 1 && !math.IsNaN(judgment.Confidence) && !math.IsInf(judgment.Confidence, 0)
}

func validAssistedVerdict(value string) bool {
	switch value {
	case "RELEVANT", "PARTIAL", "IRRELEVANT", "UNCERTAIN", "CONFLICT":
		return true
	}
	return false
}

func ValidateAssistedJudgments(pkg domain.BlindCandidatePackage, judgments []domain.AssistedJudgment, reviews []domain.HumanReview) domain.AssistedOracleReport {
	result := domain.AssistedOracleReport{SchemaVersion: "1", PackageHash: pkg.PackageHash}
	if len(pkg.Candidates) == 0 {
		result.Verdict = "INCONCLUSIVE_NO_CANDIDATES"
		return result
	}
	candidateByID := map[string]domain.BlindCandidate{}
	for _, candidate := range pkg.Candidates {
		candidateByID[candidate.CandidateID] = candidate
	}
	judgmentByID := map[string]domain.AssistedJudgment{}
	for _, judgment := range judgments {
		candidate := candidateByID[judgment.CandidateID]
		if !judgmentEvidenceValid(candidate, judgment) {
			result.Verdict = "FAIL_EVIDENCE_GROUNDING"
		}
		judgmentByID[judgment.CandidateID] = judgment
	}
	escalated := map[string]bool{}
	for _, id := range DetermineHumanEscalations(pkg, judgments) {
		escalated[id] = true
	}
	reviewByID := map[string]domain.HumanReview{}
	for _, review := range reviews {
		reviewByID[review.CandidateID] = review
	}
	if result.Verdict == "" {
		result.Verdict = "ORACLE_ACCEPTED"
	}
	for _, candidate := range pkg.Candidates {
		judgment, hasJudgment := judgmentByID[candidate.CandidateID]
		record := domain.AssistedOracleRecord{CandidateID: candidate.CandidateID, Escalated: escalated[candidate.CandidateID]}
		if hasJudgment {
			record.AIHistory = []domain.AssistedJudgment{judgment}
		}
		if !hasJudgment {
			record.Effective = "UNKNOWN"
			result.Verdict = "FAIL_EVIDENCE_GROUNDING"
			result.Records = append(result.Records, record)
			continue
		}
		if review, ok := reviewByID[candidate.CandidateID]; ok {
			record.HumanHistory = []domain.HumanReview{review}
			switch review.Action {
			case "CONFIRMED":
				record.Effective = judgment.Verdict
			case "OVERRIDDEN", "RESOLVED":
				if validAssistedVerdict(review.Verdict) && strings.TrimSpace(review.Reason) != "" {
					record.Effective = review.Verdict
				} else {
					record.Effective = "UNKNOWN"
					result.Verdict = "INCONCLUSIVE_HUMAN_REVIEW"
				}
			default:
				record.Effective = "UNKNOWN"
				result.Verdict = "INCONCLUSIVE_HUMAN_REVIEW"
			}
		} else if record.Escalated {
			record.Effective = "UNKNOWN"
			if result.Verdict == "ORACLE_ACCEPTED" {
				result.Verdict = "INCONCLUSIVE_HUMAN_REVIEW"
			}
		} else {
			record.Effective = judgment.Verdict
		}
		result.Records = append(result.Records, record)
	}
	sort.Slice(result.Records, func(i, j int) bool { return result.Records[i].CandidateID < result.Records[j].CandidateID })
	result.Metrics = assistedOracleMetrics(result.Records)
	if result.Verdict == "ORACLE_ACCEPTED" && result.Metrics.DefinitiveReviewed > 0 && (result.Metrics.RawAgreementRate < .8 || result.Metrics.OverrideRate > .2) {
		result.Verdict = "FAIL_AI_RELIABILITY"
	}
	return result
}

func assistedOracleMetrics(records []domain.AssistedOracleRecord) domain.AssistedOracleMetrics {
	metrics := domain.AssistedOracleMetrics{}
	grounded := 0
	for _, record := range records {
		if len(record.AIHistory) > 0 {
			grounded++
		}
		if !record.Escalated {
			continue
		}
		metrics.Escalated++
		if len(record.HumanHistory) == 0 {
			continue
		}
		metrics.Reviewed++
		ai := record.AIHistory[0].Verdict
		review := record.HumanHistory[0]
		if ai == "UNCERTAIN" || ai == "CONFLICT" {
			if record.Effective != "UNKNOWN" {
				metrics.ResolvedUncertain++
			}
			continue
		}
		metrics.DefinitiveReviewed++
		if record.Effective == ai {
			metrics.DefinitiveAgreements++
		} else if review.Action == "OVERRIDDEN" || review.Action == "RESOLVED" {
			metrics.DefinitiveOverrides++
		}
	}
	if len(records) > 0 {
		metrics.EvidenceGroundingRate = float64(grounded) / float64(len(records))
	}
	if metrics.DefinitiveReviewed > 0 {
		metrics.RawAgreementRate = float64(metrics.DefinitiveAgreements) / float64(metrics.DefinitiveReviewed)
		metrics.OverrideRate = float64(metrics.DefinitiveOverrides) / float64(metrics.DefinitiveReviewed)
	}
	return metrics
}

func ReplayAssistedOracle(pkg domain.BlindCandidatePackage, judgments []domain.AssistedJudgment, reviews []domain.HumanReview) domain.AssistedOracleReport {
	return ValidateAssistedJudgments(pkg, judgments, reviews)
}
