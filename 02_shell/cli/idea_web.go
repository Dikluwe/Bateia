/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func runIdeaWeb(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: idea web <discover-observable|experiment|replay|compare>")
	}
	switch args[0] {
	case "discover-observable":
		return runWebDiscoverObservable(args[1:], stdout, provider)
	case "experiment":
		return runWebExperiment(args[1:], stdout, provider)
	case "authenticated-experiment":
		return runWebAuthenticatedExperiment(args[1:], stdout, provider)
	case "replay":
		return runWebReplay(args[1:], stdout)
	case "compare":
		return runWebCompare(args[1:], stdout)
	default:
		return fmt.Errorf("operação web desconhecida")
	}
}

func runWebExperiment(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	return runWebExperimentMode(args, stdout, provider, false)
}

func runWebExperimentMode(args []string, stdout io.Writer, provider domain.ObservableWebSearch, authenticated bool) error {
	fs := flag.NewFlagSet("web experiment", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	nominalPath := fs.String("nominal", "", "queries N")
	f1Path := fs.String("functional-v1", "", "queries F1")
	f2Path := fs.String("functional-v2", "", "queries F2")
	outputDir := fs.String("output-dir", "", "diretório de snapshots")
	defaultMax := 5
	if authenticated {
		defaultMax = 10
	}
	maxPerQuery := fs.Int("max-per-query", defaultMax, "máximo por query")
	noToken := fs.Bool("no-token", false, "proíbe autenticação")
	authenticatedFlag := fs.Bool("authenticated", false, "autoriza autenticação por variável de ambiente")
	authorizationPath := fs.String("authorization", "", "autorização redigida do experimento")
	safeRemaining := fs.Int("safe-remaining", 10, "margem segura da cota")
	resumeDir := fs.String("resume-dir", "", "snapshots autenticados anteriores para retomada")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if authenticated {
		if !*authenticatedFlag || *noToken {
			return fmt.Errorf("authenticated-experiment exige --authenticated e proíbe --no-token")
		}
		if *authorizationPath == "" || *safeRemaining < 1 {
			return fmt.Errorf("authenticated-experiment exige authorization e safe-remaining positivo")
		}
		if *maxPerQuery != 10 {
			return fmt.Errorf("authenticated-experiment fixa max-per-query em 10")
		}
	} else if !*noToken {
		return fmt.Errorf("experiment exige --no-token")
	}
	if provider == nil {
		return fmt.Errorf("experiment exige provedor injetado")
	}
	if *nominalPath == "" || *f1Path == "" || *f2Path == "" || *outputDir == "" {
		return fmt.Errorf("experiment exige nominal, functional-v1, functional-v2 e output-dir")
	}
	var authorization *domain.WebExperimentAuthorization
	if authenticated {
		var value domain.WebExperimentAuthorization
		if err := readOntologyJSON(*authorizationPath, &value); err != nil {
			return err
		}
		if err := validateStep0020Authorization(value); err != nil {
			return err
		}
		authorization = &value
	}
	arms := []struct {
		name, path string
		exp        domain.FunctionalGrammarExpansion
	}{{name: "N", path: *nominalPath}, {name: "F1", path: *f1Path}, {name: "F2", path: *f2Path}}
	for i := range arms {
		if err := readOntologyJSON(arms[i].path, &arms[i].exp); err != nil {
			return err
		}
		if len(arms[i].exp.Queries) != 12 {
			return fmt.Errorf("braço %s possui %d queries; esperado 12", arms[i].name, len(arms[i].exp.Queries))
		}
		sort.SliceStable(arms[i].exp.Queries, func(a, b int) bool {
			return naturalQueryPosition(arms[i].exp.Queries[a].ID) < naturalQueryPosition(arms[i].exp.Queries[b].ID)
		})
	}
	ctx := context.Background()
	capabilities := provider.Capabilities(ctx)
	if authenticated && (!capabilities.Authenticated || !capabilities.PublicOnly || capabilities.Provider != "github-repository-search") {
		return fmt.Errorf("provedor autenticado público do GitHub não disponível")
	}
	snapshots := make([]domain.WebSearchSnapshot, 3)
	runID := "run-0019"
	requestPrefix := "0019"
	var budget *domain.WebExperimentBudget
	if authenticated {
		runID, requestPrefix = "run-0020", "0020"
		budget = &domain.WebExperimentBudget{QueriesPerArm: 12, PrimaryRequests: 36, MaxResultsPerQuery: 10, AttemptTimeoutSec: 20, TotalTimeoutSec: 600, MaxConcurrency: 2, MaxRetries: 1, SafeRemaining: *safeRemaining}
	}
	for i, arm := range arms {
		schema := "1"
		if authenticated {
			schema = "2"
		}
		snapshots[i] = domain.WebSearchSnapshot{SchemaVersion: schema, RunID: runID, Arm: arm.name, GrammarHash: arm.exp.GrammarHash, QuerySetHash: querySetHash(arm.exp.Queries), Provider: capabilities.Provider, Authorization: authorization, Budget: budget}
	}
	resumed := map[string]domain.RawWebResponse{}
	if authenticated && *resumeDir != "" {
		for i, arm := range arms {
			var previous domain.WebSearchSnapshot
			path := filepath.Join(*resumeDir, "snapshot-"+strings.ToLower(arm.name)+".json")
			if err := readOntologyJSON(path, &previous); err != nil {
				return fmt.Errorf("retomada inválida para %s: %w", arm.name, err)
			}
			if previous.Arm != arm.name || previous.QuerySetHash != snapshots[i].QuerySetHash {
				return fmt.Errorf("retomada diverge do braço/query set %s", arm.name)
			}
			for _, response := range previous.Responses {
				resumed[arm.name+"|"+response.QueryID] = response
			}
		}
	}
	limit := *maxPerQuery
	if limit > 10 {
		limit = 10
	}
	totalCtx := ctx
	var totalCancel context.CancelFunc
	if authenticated {
		totalCtx, totalCancel = context.WithTimeout(ctx, 10*time.Minute)
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
				resp := domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: capabilities.Provider, RetrievedAt: time.Now().UTC().Format(time.RFC3339), Status: "SKIPPED_BUDGET", Arm: arms[armIndex].name, Position: pos*3 + armIndex + 1, AttemptNumber: 1}
				resp = usecase.SealWebResponse(req, resp)
				snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, resp)
				receipt.SkippedBudget++
				continue
			}
			attemptCtx := totalCtx
			var cancel context.CancelFunc
			if authenticated {
				attemptCtx, cancel = context.WithTimeout(totalCtx, 20*time.Second)
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
				resp = domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: capabilities.Provider, RetrievedAt: time.Now().UTC().Format(time.RFC3339), Status: status}
				resp = usecase.SealWebResponse(req, resp)
			}
			resp.Arm, resp.Position, resp.AttemptNumber = arms[armIndex].name, pos*3+armIndex+1, 1
			resp = usecase.SealWebResponse(req, resp)
			snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, resp)
			if resp.Status == "SUCCESS" || resp.Status == "EMPTY" {
				receipt.Completed++
			} else {
				receipt.Unknown++
			}
			if authenticated && resp.Receipt.RateLimit.Remaining != nil && *resp.Receipt.RateLimit.Remaining <= int64(*safeRemaining) {
				stopForBudget = true
				receipt.StopReason = "SAFE_RATE_LIMIT_MARGIN"
			}
		}
	}
	for i := range snapshots {
		usecase.StableSnapshotResponses(snapshots[i].Responses)
		snapshots[i].SnapshotHash = usecase.WebSnapshotHash(snapshots[i])
		if err := rejectSensitiveWebArtifact(snapshots[i]); err != nil {
			return err
		}
	}
	if authenticated {
		receipt.Verdict = authenticatedExperimentVerdict(snapshots)
	}
	if err := rejectSensitiveWebArtifact(receipt); err != nil {
		return err
	}
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return err
	}
	for i := range snapshots {
		path := filepath.Join(*outputDir, "snapshot-"+strings.ToLower(snapshots[i].Arm)+".json")
		if err := writeIdeaJSON(snapshots[i], path, stdout); err != nil {
			return err
		}
	}
	if authenticated {
		if err := writeIdeaJSON(receipt, filepath.Join(*outputDir, "receipt.json"), stdout); err != nil {
			return err
		}
	}
	return nil
}

func runWebDiscoverObservable(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	fs := flag.NewFlagSet("web discover-observable", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	queriesPath := fs.String("queries", "", "FunctionalExpansion")
	arm := fs.String("arm", "F2", "braço N, F1 ou F2")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo por query")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *queriesPath == "" {
		return fmt.Errorf("discover-observable exige queries")
	}
	if provider == nil {
		return fmt.Errorf("discover-observable exige provedor injetado por L4")
	}
	var exp domain.FunctionalGrammarExpansion
	if err := readOntologyJSON(*queriesPath, &exp); err != nil {
		return err
	}

	ctx := context.Background()
	capabilities := provider.Capabilities(ctx)

	var responses []domain.RawWebResponse
	for _, q := range exp.Queries {
		req := domain.ObservableWebRequest{
			RequestID:  fmt.Sprintf("req-%s", q.ID),
			QueryID:    q.ID,
			Query:      q.Text,
			MaxResults: *maxPerQuery,
			Locale:     "en-US",
		}
		resp, err := provider.SearchOne(ctx, req)
		if err != nil {
			resp = domain.RawWebResponse{
				RequestID:   req.RequestID,
				QueryID:     req.QueryID,
				Provider:    capabilities.Provider,
				RetrievedAt: time.Now().UTC().Format(time.RFC3339),
				Status:      "FAILED",
			}
		}
		responses = append(responses, resp)
	}

	snapshot := domain.WebSearchSnapshot{
		SchemaVersion: "1",
		RunID:         "run-0018",
		Arm:           *arm,
		GrammarHash:   exp.GrammarHash,
		QuerySetHash:  querySetHash(exp.Queries),
		Provider:      capabilities.Provider,
		Responses:     responses,
	}
	snapshot.SnapshotHash = usecase.WebSnapshotHash(snapshot)

	return writeIdeaJSON(snapshot, *output, stdout)
}

func querySetHash(queries []domain.RealizedFunctionalQuery) string {
	copyQueries := append([]domain.RealizedFunctionalQuery(nil), queries...)
	sort.Slice(copyQueries, func(i, j int) bool { return copyQueries[i].ID < copyQueries[j].ID })
	b, _ := json.Marshal(copyQueries)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func runWebReplay(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("web replay", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	snapshotPath := fs.String("snapshot", "", "WebSearchSnapshot")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *snapshotPath == "" {
		return fmt.Errorf("replay exige snapshot")
	}
	var snap domain.WebSearchSnapshot
	if err := readOntologyJSON(*snapshotPath, &snap); err != nil {
		return err
	}
	res, err := usecase.ReplayWebSnapshot(snap)
	if err != nil {
		return err
	}
	return writeIdeaJSON(res, *output, stdout)
}

func runWebCompare(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("web compare", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	nominalPath := fs.String("nominal", "", "Nominal Snapshot")
	f1Path := fs.String("functional-v1", "", "F1 Snapshot")
	f2Path := fs.String("functional-v2", "", "F2 Snapshot")
	judgmentsPath := fs.String("judgments", "", "Judgments JSON")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *nominalPath == "" || *f1Path == "" || *f2Path == "" || *judgmentsPath == "" {
		return fmt.Errorf("compare exige nominal, functional-v1, functional-v2 e judgments")
	}

	var nom, f1, f2 domain.WebSearchSnapshot
	var j domain.BlindedJudgmentsFile
	if err := readOntologyJSON(*nominalPath, &nom); err != nil {
		return err
	}
	if err := readOntologyJSON(*f1Path, &f1); err != nil {
		return err
	}
	if err := readOntologyJSON(*f2Path, &f2); err != nil {
		return err
	}
	if err := readOntologyJSON(*judgmentsPath, &j); err != nil {
		return err
	}

	report, err := usecase.CompareGrammarExperiments(nom, f1, f2, j)
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}
