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
	"flag"
	"fmt"
	"io"
	"time"
)

func runIdeaWeb(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: idea web <discover-observable|replay|compare>")
	}
	switch args[0] {
	case "discover-observable":
		return runWebDiscoverObservable(args[1:], stdout)
	case "replay":
		return runWebReplay(args[1:], stdout)
	case "compare":
		return runWebCompare(args[1:], stdout)
	default:
		return fmt.Errorf("operação web desconhecida")
	}
}

func runWebDiscoverObservable(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("web discover-observable", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	queriesPath := fs.String("queries", "", "FunctionalExpansion")
	providerName := fs.String("provider", "web-mock", "provedor")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo por query")
	output := fs.String("output", "", "saída")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *queriesPath == "" {
		return fmt.Errorf("discover-observable exige queries")
	}
	var exp domain.FunctionalExpansion
	if err := readOntologyJSON(*queriesPath, &exp); err != nil {
		return err
	}

	client := usecase.NewMockObservableWebSearch(*providerName)
	ctx := context.Background()

	var responses []domain.RawWebResponse
	for _, q := range exp.Queries {
		req := domain.ObservableWebRequest{
			RequestID:  fmt.Sprintf("req-%s", q.ID),
			QueryID:    q.ID,
			Query:      q.Text,
			MaxResults: *maxPerQuery,
			Locale:     "en-US",
		}
		resp, err := client.SearchOne(ctx, req)
		if err != nil {
			resp = domain.RawWebResponse{
				RequestID:   req.RequestID,
				QueryID:     req.QueryID,
				Provider:    *providerName,
				RetrievedAt: time.Now().UTC().Format(time.RFC3339),
				Status:      "FAILED",
			}
		}
		responses = append(responses, resp)
	}

	snapshot := domain.WebSearchSnapshot{
		SchemaVersion: "1",
		RunID:         "run-0018",
		Arm:           "F2",
		GrammarHash:   exp.LexiconHash,
		QuerySetHash:  exp.LexiconHash,
		Provider:      *providerName,
		Responses:     responses,
	}
	snapshot.SnapshotHash = usecase.WebSnapshotHash(snapshot)

	return writeIdeaJSON(snapshot, *output, stdout)
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
