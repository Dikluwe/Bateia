/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_replay.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"context"
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func WebSnapshotHash(s domain.WebSearchSnapshot) string {
	s.SnapshotHash = ""
	b, _ := json.Marshal(s)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func CanonicalizeURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return rawURL
	}
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.TrimSuffix(u.Path, "/")
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func ReplayWebSnapshot(snapshot domain.WebSearchSnapshot) (domain.WebReplayResult, error) {
	if snapshot.SchemaVersion == "" {
		return domain.WebReplayResult{}, fmt.Errorf("snapshot inválido")
	}

	seenURLs := make(map[string]bool)
	var deduped []domain.RawWebResult
	var missing []string

	for _, resp := range snapshot.Responses {
		for _, res := range resp.Results {
			if res.URL == "" {
				missing = append(missing, "url")
				continue
			}
			canon := res.CanonicalURL
			if canon == "" {
				canon = CanonicalizeURL(res.URL)
			}
			if !seenURLs[canon] {
				seenURLs[canon] = true
				res.CanonicalURL = canon
				deduped = append(deduped, res)
			}
		}
	}

	hash := WebSnapshotHash(snapshot)
	replayable := true
	if snapshot.SnapshotHash != "" && snapshot.SnapshotHash != hash {
		replayable = false
	}

	return domain.WebReplayResult{
		SchemaVersion: "1",
		SnapshotHash:  hash,
		TotalQueries:  len(snapshot.Responses),
		UniqueHits:    len(deduped),
		Deduplicated:  deduped,
		MissingFields: missing,
		Replayable:    replayable,
	}, nil
}

func CompareGrammarExperiments(nominal, f1, f2 domain.WebSearchSnapshot, judgments domain.BlindedJudgmentsFile) (domain.GrammarExperimentReport, error) {
	jMap := make(map[string]domain.BlindedJudgment)
	for _, j := range judgments.Judgments {
		jMap[j.CanonicalURL] = j
	}

	calcMetrics := func(s domain.WebSearchSnapshot) (float64, float64) {
		r, err := ReplayWebSnapshot(s)
		if err != nil || len(r.Deduplicated) == 0 {
			return 0, 0
		}
		relCount := 0
		repoCount := 0
		n := len(r.Deduplicated)
		if n > 10 {
			n = 10
		}
		for i := 0; i < n; i++ {
			item := r.Deduplicated[i]
			if j, ok := jMap[item.CanonicalURL]; ok {
				if j.Judgment == "RELEVANT" {
					relCount++
				}
				if j.IsRepoBacked {
					repoCount++
				}
			}
		}
		precision := float64(relCount) / float64(n)
		recall := float64(repoCount) / 4.0 // 4 seeds
		return precision, recall
	}

	pNom, rNom := calcMetrics(nominal)
	pF1, rF1 := calcMetrics(f1)
	pF2, rF2 := calcMetrics(f2)

	gain := (rF2 - rNom) * 100.0

	verdict := "MAINTAINED"
	if gain >= 30.0 && pF2 >= 0.60 {
		verdict = "PROMOTED"
	} else if gain < 0 {
		verdict = "REJECTED"
	}

	return domain.GrammarExperimentReport{
		SchemaVersion:             "1",
		JudgedPrecisionAt10Nominal: pNom,
		JudgedPrecisionAt10F1:      pF1,
		JudgedPrecisionAt10F2:      pF2,
		RepoBackedRecallNominal:   rNom,
		RepoBackedRecallF1:        rF1,
		RepoBackedRecallF2:        rF2,
		RecallGainOverNominal:     gain,
		Verdict:                   verdict,
	}, nil
}

type MockObservableWebSearch struct {
	ProviderName string
}

func NewMockObservableWebSearch(provider string) *MockObservableWebSearch {
	if provider == "" {
		provider = "web-mock"
	}
	return &MockObservableWebSearch{ProviderName: provider}
}

func (c *MockObservableWebSearch) Capabilities(ctx context.Context) domain.WebSearchCapabilities {
	return domain.WebSearchCapabilities{
		Provider:        c.ProviderName,
		SupportsLocale:  true,
		PreservesRank:   true,
		MaxResultsLimit: 10,
	}
}

func (c *MockObservableWebSearch) SearchOne(ctx context.Context, req domain.ObservableWebRequest) (domain.RawWebResponse, error) {
	if req.Query == "" {
		return domain.RawWebResponse{}, fmt.Errorf("query vazia")
	}

	results := []domain.RawWebResult{
		{
			Rank:         1,
			URL:          fmt.Sprintf("https://github.com/example/%s-tool", req.QueryID),
			CanonicalURL: fmt.Sprintf("https://github.com/example/%s-tool", req.QueryID),
			Title:        fmt.Sprintf("Tool for %s", req.Query),
			Snippet:      fmt.Sprintf("Open source tool implementing %s", req.Query),
			IsAd:         false,
		},
	}

	return domain.RawWebResponse{
		RequestID:   req.RequestID,
		QueryID:     req.QueryID,
		Provider:    c.ProviderName,
		RetrievedAt: "2026-08-26T21:00:00Z",
		Status:      "SUCCESS",
		HTTPStatus:  200,
		Results:     results,
		Receipt: domain.WebReceipt{
			SanitizedURL: "https://api.search.provider/v1/search",
			ResponseTime: 42,
		},
	}, nil
}
