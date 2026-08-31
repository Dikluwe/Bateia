/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	case "replay":
		return runWebReplay(args[1:], stdout)
	case "compare":
		return runWebCompare(args[1:], stdout)
	default:
		return fmt.Errorf("operação web desconhecida")
	}
}

func runWebExperiment(args []string, stdout io.Writer, provider domain.ObservableWebSearch) error {
	fs := flag.NewFlagSet("web experiment", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	nominalPath := fs.String("nominal", "", "queries N")
	f1Path := fs.String("functional-v1", "", "queries F1")
	f2Path := fs.String("functional-v2", "", "queries F2")
	outputDir := fs.String("output-dir", "", "diretório de snapshots")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo por query")
	noToken := fs.Bool("no-token", false, "proíbe autenticação")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*noToken {
		return fmt.Errorf("experiment exige --no-token")
	}
	if provider == nil {
		return fmt.Errorf("experiment exige provedor injetado")
	}
	if *nominalPath == "" || *f1Path == "" || *f2Path == "" || *outputDir == "" {
		return fmt.Errorf("experiment exige nominal, functional-v1, functional-v2 e output-dir")
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
		sort.Slice(arms[i].exp.Queries, func(a, b int) bool { return arms[i].exp.Queries[a].ID < arms[i].exp.Queries[b].ID })
	}
	ctx := context.Background()
	capabilities := provider.Capabilities(ctx)
	snapshots := make([]domain.WebSearchSnapshot, 3)
	for i, arm := range arms {
		snapshots[i] = domain.WebSearchSnapshot{SchemaVersion: "1", RunID: "run-0019", Arm: arm.name, GrammarHash: arm.exp.GrammarHash, QuerySetHash: querySetHash(arm.exp.Queries), Provider: capabilities.Provider}
	}
	for pos := 0; pos < 12; pos++ {
		for armIndex := 0; armIndex < 3; armIndex++ {
			q := arms[armIndex].exp.Queries[pos]
			req := domain.ObservableWebRequest{RequestID: fmt.Sprintf("0019-%02d-%s", pos+1, arms[armIndex].name), QueryID: q.ID, Query: q.Text, MaxResults: *maxPerQuery, Locale: "en-US"}
			resp, err := provider.SearchOne(ctx, req)
			if err != nil && resp.Status == "" {
				resp = domain.RawWebResponse{RequestID: req.RequestID, QueryID: req.QueryID, Provider: capabilities.Provider, RetrievedAt: time.Now().UTC().Format(time.RFC3339), Status: "UNKNOWN"}
				resp = usecase.SealWebResponse(req, resp)
			}
			snapshots[armIndex].Responses = append(snapshots[armIndex].Responses, resp)
		}
	}
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return err
	}
	for i := range snapshots {
		usecase.StableSnapshotResponses(snapshots[i].Responses)
		snapshots[i].SnapshotHash = usecase.WebSnapshotHash(snapshots[i])
		path := filepath.Join(*outputDir, "snapshot-"+strings.ToLower(snapshots[i].Arm)+".json")
		if err := writeIdeaJSON(snapshots[i], path, stdout); err != nil {
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
