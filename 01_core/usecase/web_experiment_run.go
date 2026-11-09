/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/web_experiment_run.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type WebExperimentArm struct {
	Name      string
	Expansion domain.FunctionalGrammarExpansion
}
type WebExperimentInput struct {
	Arms                       []WebExperimentArm
	Authenticated              bool
	Authorization              *domain.WebExperimentAuthorization
	SafeRemaining, MaxPerQuery int
	Resume                     bool
}
type WebExperimentOutput struct {
	Snapshots []domain.WebSearchSnapshot
	Receipt   domain.WebExperimentReceipt
}
type WebExperimentResumeReader interface {
	PreviousSnapshot(string) (domain.WebSearchSnapshot, error)
}

func ExecuteWebExperiment(ctx context.Context, input WebExperimentInput, provider domain.ObservableWebSearch, resume WebExperimentResumeReader) (WebExperimentOutput, error) {
	if provider == nil || len(input.Arms) != 3 {
		return WebExperimentOutput{}, fmt.Errorf("experiment exige três braços e provedor injetado")
	}
	if input.Authenticated {
		if err := ValidateWebExperimentAuthorizationValue(input.Authorization); err != nil {
			return WebExperimentOutput{}, err
		}
	}
	if input.Resume && resume == nil {
		return WebExperimentOutput{}, fmt.Errorf("retomada exige leitor de snapshots")
	}

	arms := make([]struct {
		name string
		exp  domain.FunctionalGrammarExpansion
	}, len(input.Arms))
	seenArms := map[string]bool{}
	for i, arm := range input.Arms {
		if (arm.Name != "N" && arm.Name != "F1" && arm.Name != "F2") || seenArms[arm.Name] {
			return WebExperimentOutput{}, fmt.Errorf("braço inválido ou duplicado: %s", arm.Name)
		}
		seenArms[arm.Name] = true
		arms[i].name = arm.Name
		arms[i].exp = arm.Expansion
		arms[i].exp.Queries = append([]domain.RealizedFunctionalQuery(nil), arm.Expansion.Queries...)
	}
	for i := range arms {
		if len(arms[i].exp.Queries) != 12 {
			return WebExperimentOutput{}, fmt.Errorf("braço %s possui %d queries; esperado 12", arms[i].name, len(arms[i].exp.Queries))
		}
		seenQueries := map[string]bool{}
		for _, query := range arms[i].exp.Queries {
			if strings.TrimSpace(query.ID) == "" || seenQueries[query.ID] || strings.TrimSpace(query.Text) == "" {
				return WebExperimentOutput{}, fmt.Errorf("query inválida ou duplicada no braço %s", arms[i].name)
			}
			seenQueries[query.ID] = true
		}
		sort.SliceStable(arms[i].exp.Queries, func(a, b int) bool {
			return naturalQueryPosition(arms[i].exp.Queries[a].ID) < naturalQueryPosition(arms[i].exp.Queries[b].ID)
		})
	}
	capabilities := provider.Capabilities(ctx)
	if input.Authenticated && (!capabilities.Authenticated || !capabilities.PublicOnly || capabilities.Provider != "github-repository-search") {
		return WebExperimentOutput{}, fmt.Errorf("provedor autenticado público do GitHub não disponível")
	}
	snapshots := make([]domain.WebSearchSnapshot, 3)
	runID := "run-0019"
	requestPrefix := "0019"
	var budget *domain.WebExperimentBudget
	if input.Authenticated {
		runID, requestPrefix = "run-0020", "0020"
		budget = &domain.WebExperimentBudget{QueriesPerArm: 12, PrimaryRequests: 36, MaxResultsPerQuery: 10, AttemptTimeoutSec: 20, TotalTimeoutSec: 600, MaxConcurrency: 2, MaxRetries: 1, SafeRemaining: input.SafeRemaining}
	}
	for i, arm := range arms {
		schema := "1"
		if input.Authenticated {
			schema = "2"
		}
		snapshots[i] = domain.WebSearchSnapshot{SchemaVersion: schema, RunID: runID, Arm: arm.name, GrammarHash: arm.exp.GrammarHash, QuerySetHash: WebExperimentQuerySetHash(arm.exp.Queries), Provider: capabilities.Provider, Authorization: input.Authorization, Budget: budget}
	}
	limit := input.MaxPerQuery
	if limit > 10 {
		limit = 10
	}
	resumed := map[string]domain.RawWebResponse{}
	if input.Authenticated && input.Resume {
		for i, arm := range arms {
			var previous domain.WebSearchSnapshot
			var readErr error
			previous, readErr = resume.PreviousSnapshot(arm.name)
			if err := readErr; err != nil {
				return WebExperimentOutput{}, fmt.Errorf("retomada sem integridade/reprodutibilidade para %s: %w", arm.name, err)
			}
			if err := validateWebExperimentResume(previous, snapshots[i], arm.exp.Queries, i, limit); err != nil {
				return WebExperimentOutput{}, fmt.Errorf("retomada sem integridade/reprodutibilidade para %s: %w", arm.name, err)
			}
			for _, response := range previous.Responses {
				resumed[arm.name+"|"+response.QueryID] = response
			}
		}
	}
	totalCtx := ctx
	var totalCancel context.CancelFunc
	if input.Authenticated {
		var timeoutErr error
		totalCtx, totalCancel, timeoutErr = domain.ExecutionTimeout(ctx, 10*time.Minute)
		if timeoutErr != nil {
			return WebExperimentOutput{}, timeoutErr
		}
		defer totalCancel()
	}
	receipt := domain.WebExperimentReceipt{SchemaVersion: "1", RunID: runID, Planned: 36}
	stopForBudget := false
	for pos := 0; pos < 12; pos++ {
		for armIndex := 0; armIndex < 3; armIndex++ {
			q := arms[armIndex].exp.Queries[pos]
			req := domain.ObservableWebRequest{RequestID: fmt.Sprintf("%s-%02d-%s", requestPrefix, pos+1, arms[armIndex].name), QueryID: q.ID, Query: q.Text, MaxResults: limit, Locale: "en-US"}
			if previous, ok := resumed[arms[armIndex].name+"|"+q.ID]; ok && (previous.Status == "SUCCESS" || previous.Status == "EMPTY") {
				snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, previous)
				receipt.Completed++
				continue
			}
			if stopForBudget || totalCtx.Err() != nil {
				resp := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: capabilities.Provider, RetrievedAt: domain.RuntimeNow(ctx).UTC().Format(time.RFC3339), Status: "SKIPPED_BUDGET", Arm: arms[armIndex].name, Position: pos*3 + armIndex + 1, AttemptNumber: 1}
				resp = SealWebResponse(req, resp)
				snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, resp)
				receipt.SkippedBudget++
				continue
			}
			attemptCtx := totalCtx
			var cancel context.CancelFunc
			if input.Authenticated {
				var timeoutErr error
				attemptCtx, cancel, timeoutErr = domain.ExecutionTimeout(totalCtx, 20*time.Second)
				if timeoutErr != nil {
					return WebExperimentOutput{}, timeoutErr
				}
			}
			receipt.Started++
			if _, ok := resumed[arms[armIndex].name+"|"+q.ID]; ok {
				receipt.Retried++
			}
			resp, err := provider.SearchOne(attemptCtx, req)
			if cancel != nil {
				cancel()
			}
			if err != nil && resp.Status == "" {
				status := "UNKNOWN_PROVIDER_RESPONSE"
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					status = "UNKNOWN_TIMEOUT"
				}
				resp = domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: capabilities.Provider, RetrievedAt: domain.RuntimeNow(ctx).UTC().Format(time.RFC3339), Status: status}
				resp = SealWebResponse(req, resp)
			}
			resp.Arm, resp.Position, resp.AttemptNumber = arms[armIndex].name, pos*3+armIndex+1, 1
			resp = SealWebResponse(req, resp)
			snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, resp)
			if resp.Status == "SUCCESS" || resp.Status == "EMPTY" {
				receipt.Completed++
			} else {
				receipt.Unknown++
			}
			if input.Authenticated && resp.Receipt.RateLimit.Remaining != nil && *resp.Receipt.RateLimit.Remaining <= int64(input.SafeRemaining) {
				stopForBudget = true
				receipt.StopReason = "SAFE_RATE_LIMIT_MARGIN"
			}
		}
	}
	for i := range snapshots {
		StableSnapshotResponses(snapshots[i].Responses)
		snapshots[i].SnapshotHash = WebSnapshotHash(snapshots[i])
	}
	if input.Authenticated {
		receipt.Verdict = AuthenticatedWebExperimentVerdict(snapshots)
	}
	return WebExperimentOutput{Snapshots: snapshots, Receipt: receipt}, nil
}

// Admission checks the frozen snapshot against the current plan without resealing it.
func validateWebExperimentResume(previous, planned domain.WebSearchSnapshot, queries []domain.RealizedFunctionalQuery, armIndex, limit int) error {
	if previous.SchemaVersion != "1" && previous.SchemaVersion != "2" {
		return domain.ErrIntegrity
	}
	replay, err := ReplayWebSnapshot(previous)
	if err != nil || !replay.Replayable || len(replay.MissingFields) != 0 {
		return domain.ErrIntegrity
	}
	if previous.RunID != planned.RunID || previous.Arm != planned.Arm || previous.Provider != planned.Provider || previous.GrammarHash != planned.GrammarHash || previous.QuerySetHash != planned.QuerySetHash {
		return domain.ErrIntegrity
	}
	if (previous.SchemaVersion == "2" || previous.Authorization != nil) && !reflect.DeepEqual(previous.Authorization, planned.Authorization) {
		return domain.ErrIntegrity
	}
	if (previous.SchemaVersion == "2" || previous.Budget != nil) && !reflect.DeepEqual(previous.Budget, planned.Budget) {
		return domain.ErrIntegrity
	}
	positions := map[string]int{}
	for pos, query := range queries {
		positions[query.ID] = pos
	}
	seen := map[string]bool{}
	for _, response := range previous.Responses {
		pos, known := positions[response.QueryID]
		if !known || seen[response.QueryID] {
			return domain.ErrIntegrity
		}
		seen[response.QueryID] = true
		query := queries[pos]
		expected := domain.ObservableWebRequest{RequestID: fmt.Sprintf("0020-%02d-%s", pos+1, planned.Arm), QueryID: query.ID, Query: query.Text, MaxResults: limit, Locale: "en-US"}
		if response.Request != expected || response.RequestID != expected.RequestID || response.Provider != planned.Provider || response.Arm != planned.Arm || response.Position != pos*3+armIndex+1 || response.AttemptNumber != 1 {
			return domain.ErrIntegrity
		}
	}
	return nil
}

func ValidateWebExperimentAuthorization(a domain.WebExperimentAuthorization) error {
	if a.Provider != "github" || a.Purpose != "public_repository_search" || a.Authentication != "read_only_token" || a.PrivateRepositoriesAllowed || a.PlannedQueries != 36 || !a.Authorized {
		return fmt.Errorf("autorização do Passo 0020 inválida")
	}
	return nil
}

func naturalQueryPosition(id string) int {
	i := len(id)
	for i > 0 && id[i-1] >= '0' && id[i-1] <= '9' {
		i--
	}
	if i == len(id) {
		return 0
	}
	n, _ := strconv.Atoi(id[i:])
	return n
}

func AuthenticatedWebExperimentVerdict(snapshots []domain.WebSearchSnapshot) string {
	minRate, maxRate, rateLimited := 1.0, 0.0, false
	for _, snapshot := range snapshots {
		completed := 0
		for _, response := range snapshot.Responses {
			if response.Status == "SUCCESS" || response.Status == "EMPTY" {
				completed++
			}
			if response.Status == "UNKNOWN_RATE_LIMIT" {
				rateLimited = true
			}
			if snapshot.Budget != nil && response.Receipt.RateLimit.Remaining != nil && *response.Receipt.RateLimit.Remaining <= int64(snapshot.Budget.SafeRemaining) {
				rateLimited = true
			}
		}
		rate := float64(completed) / 12
		if rate < minRate {
			minRate = rate
		}
		if rate > maxRate {
			maxRate = rate
		}
	}
	if minRate >= .8 && maxRate-minRate <= .1 {
		return "READY_FOR_BLINDED_JUDGMENT"
	}
	if rateLimited {
		return "INCONCLUSIVE_RATE_LIMIT"
	}
	return "UNKNOWN_PROVIDER_EXPOSURE"
}
func WebExperimentQuerySetHash(queries []domain.RealizedFunctionalQuery) string {
	copyQueries := append([]domain.RealizedFunctionalQuery(nil), queries...)
	sort.Slice(copyQueries, func(i, j int) bool { return copyQueries[i].ID < copyQueries[j].ID })
	b, _ := json.Marshal(copyQueries)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func ValidateWebExperimentAuthorizationValue(value *domain.WebExperimentAuthorization) error {
	if value == nil {
		return fmt.Errorf("autorização do Passo 0020 inválida")
	}
	return ValidateWebExperimentAuthorization(*value)
}

func DiscoverWebQueries(ctx context.Context, exp domain.FunctionalGrammarExpansion, arm string, maxPerQuery int, provider domain.ObservableWebSearch) domain.WebSearchSnapshot {
	capabilities := provider.Capabilities(ctx)

	var responses []domain.RawWebResponse
	for _, q := range exp.Queries {
		req := domain.ObservableWebRequest{
			RequestID:  fmt.Sprintf("req-%s", q.ID),
			QueryID:    q.ID,
			Query:      q.Text,
			MaxResults: maxPerQuery,
			Locale:     "en-US",
		}
		resp, err := provider.SearchOne(ctx, req)
		if err != nil {
			resp = domain.RawWebResponse{
				RequestID:   req.RequestID,
				QueryID:     req.QueryID,
				Provider:    capabilities.Provider,
				RetrievedAt: domain.RuntimeNow(ctx).UTC().Format(time.RFC3339),
				Status:      "FAILED",
			}
		}
		responses = append(responses, resp)
	}

	snapshot := domain.WebSearchSnapshot{
		SchemaVersion: "1",
		RunID:         "run-0018",
		Arm:           arm,
		GrammarHash:   exp.GrammarHash,
		QuerySetHash:  WebExperimentQuerySetHash(exp.Queries),
		Provider:      capabilities.Provider,
		Responses:     responses,
	}
	snapshot.SnapshotHash = WebSnapshotHash(snapshot)

	return snapshot
}

type WebExperimentAdmission struct {
	Authenticated, AuthenticatedFlag, NoToken, HasAuthorization bool
	SafeRemaining, MaxPerQuery                                  int
}

func ValidateWebExperimentAdmission(admission WebExperimentAdmission) error {
	if admission.Authenticated {
		if !admission.AuthenticatedFlag || admission.NoToken {
			return fmt.Errorf("authenticated-experiment exige --authenticated e proíbe --no-token")
		}
		if !admission.HasAuthorization || admission.SafeRemaining < 1 {
			return fmt.Errorf("authenticated-experiment exige authorization e safe-remaining positivo")
		}
		if admission.MaxPerQuery != 10 {
			return fmt.Errorf("authenticated-experiment fixa max-per-query em 10")
		}
	} else if !admission.NoToken {
		return fmt.Errorf("experiment exige --no-token")
	}
	return nil
}
