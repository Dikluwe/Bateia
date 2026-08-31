/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_replay.md
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
	"net/url"
	"sort"
	"strings"
)

func hashJSON(value any) string {
	b, _ := json.Marshal(value)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func WebRequestHash(req domain.ObservableWebRequest) string { return hashJSON(req) }
func WebResponseHash(resp domain.RawWebResponse) string {
	resp.ResponseHash = ""
	return hashJSON(resp)
}
func WebSnapshotHash(s domain.WebSearchSnapshot) string { s.SnapshotHash = ""; return hashJSON(s) }

func CanonicalizeURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.TrimSuffix(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func SealWebResponse(req domain.ObservableWebRequest, resp domain.RawWebResponse) domain.RawWebResponse {
	resp.RequestHash = WebRequestHash(req)
	resp.ResponseHash = WebResponseHash(resp)
	return resp
}

func ReplayWebSnapshot(snapshot domain.WebSearchSnapshot) (domain.WebReplayResult, error) {
	if snapshot.SchemaVersion == "" {
		return domain.WebReplayResult{}, fmt.Errorf("snapshot inválido")
	}
	result := domain.WebReplayResult{SchemaVersion: "1", SnapshotHash: WebSnapshotHash(snapshot), TotalQueries: len(snapshot.Responses), Replayable: true}
	if snapshot.SnapshotHash == "" || snapshot.SnapshotHash != result.SnapshotHash {
		result.Replayable = false
		result.MissingFields = append(result.MissingFields, "snapshot_hash")
	}
	seenReq, seenURL := map[string]bool{}, map[string]bool{}
	for _, resp := range snapshot.Responses {
		if resp.RequestID == "" || resp.QueryID == "" || resp.Provider == "" || resp.RetrievedAt == "" || resp.Status == "" {
			result.Replayable = false
			result.MissingFields = append(result.MissingFields, "response_identity")
			continue
		}
		if seenReq[resp.RequestID] {
			result.Replayable = false
			result.MissingFields = append(result.MissingFields, "duplicate_request_id")
			continue
		}
		seenReq[resp.RequestID] = true
		if resp.RequestHash == "" || resp.ResponseHash == "" || WebResponseHash(resp) != resp.ResponseHash {
			result.Replayable = false
			result.MissingFields = append(result.MissingFields, "response_hash")
		}
		for _, hit := range resp.Results {
			canon := CanonicalizeURL(hit.CanonicalURL)
			if canon == "" {
				canon = CanonicalizeURL(hit.URL)
			}
			if hit.Rank < 1 || canon == "" || hit.Title == "" || hit.Snippet == "" {
				result.MissingFields = append(result.MissingFields, "hit_material_field")
				continue
			}
			hit.CanonicalURL = canon
			if hit.IsAd {
				result.Sponsored = append(result.Sponsored, hit)
				continue
			}
			if !seenURL[canon] {
				seenURL[canon] = true
				result.Organic = append(result.Organic, hit)
			}
		}
	}
	result.Deduplicated = append([]domain.RawWebResult(nil), result.Organic...)
	result.UniqueHits = len(result.Organic)
	return result, nil
}

func CompareGrammarExperiments(nominal, f1, f2 domain.WebSearchSnapshot, judgments domain.BlindedJudgmentsFile) (domain.GrammarExperimentReport, error) {
	jMap, expected := map[string]domain.BlindedJudgment{}, map[string]bool{}
	for _, j := range judgments.Judgments {
		canon := CanonicalizeURL(j.CanonicalURL)
		if canon == "" {
			continue
		}
		jMap[canon] = j
		if j.ExpectedID != "" {
			expected[j.ExpectedID] = true
		}
	}
	calc := func(s domain.WebSearchSnapshot) (float64, *float64, error) {
		r, err := ReplayWebSnapshot(s)
		if err != nil {
			return 0, nil, err
		}
		if !r.Replayable {
			return 0, nil, fmt.Errorf("snapshot %s não reexecutável", s.Arm)
		}
		n := len(r.Organic)
		if n > 10 {
			n = 10
		}
		if n == 0 {
			return 0, nil, nil
		}
		relevant, found := 0, map[string]bool{}
		for i := 0; i < n; i++ {
			if j, ok := jMap[r.Organic[i].CanonicalURL]; ok {
				if j.Judgment == "RELEVANT" {
					relevant++
				}
				if j.IsRepoBacked && j.ExpectedID != "" {
					found[j.ExpectedID] = true
				}
			}
		}
		precision := float64(relevant) / float64(n)
		if len(expected) == 0 {
			return precision, nil, nil
		}
		recall := float64(len(found)) / float64(len(expected))
		return precision, &recall, nil
	}
	pN, rN, e := calc(nominal)
	if e != nil {
		return domain.GrammarExperimentReport{}, e
	}
	p1, r1, e := calc(f1)
	if e != nil {
		return domain.GrammarExperimentReport{}, e
	}
	p2, r2, e := calc(f2)
	if e != nil {
		return domain.GrammarExperimentReport{}, e
	}
	report := domain.GrammarExperimentReport{SchemaVersion: "2", JudgedPrecisionAt10Nominal: pN, JudgedPrecisionAt10F1: p1, JudgedPrecisionAt10F2: p2, RepoBackedRecallNominal: rN, RepoBackedRecallF1: r1, RepoBackedRecallF2: r2, Verdict: "INCONCLUSIVE"}
	if rN != nil && r1 != nil && r2 != nil {
		report.RecallGainOverNominal = (*r2 - *rN) * 100
		if report.RecallGainOverNominal >= 30 && p2 >= .60 && *r2 >= .70 && p2 >= p1-.05 {
			report.Verdict = "PROMOTED"
		} else if *r2 < *r1 || p2 < p1-.05 {
			report.Verdict = "REJECTED"
		} else {
			report.Verdict = "VALIDATING"
		}
	}
	return report, nil
}

func StableSnapshotResponses(responses []domain.RawWebResponse) {
	sort.Slice(responses, func(i, j int) bool { return responses[i].QueryID < responses[j].QueryID })
}
