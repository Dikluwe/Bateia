/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_research.md
 * @layer L2
 * @updated 2026-10-06
 */

package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

type researchAttempt struct {
	Query        string                     `json:"query"`
	ConnectorID  string                     `json:"connector_id"`
	Status       domain.MarketAttemptStatus `json:"status"`
	ResultCount  int                        `json:"result_count"`
	ResponseHash string                     `json:"response_hash,omitempty"`
	Error        string                     `json:"error,omitempty"`
}

type researchSnapshot struct {
	SchemaVersion string                    `json:"schema_version"`
	Queries       []string                  `json:"queries"`
	Attempts      []researchAttempt         `json:"attempts"`
	Candidates    []domain.ProductCandidate `json:"candidates"`
}

// RunIdeaResearchCommand searches scholarly software catalogues and always writes a traceable snapshot.
func RunIdeaResearchCommand(args []string, stdout io.Writer, connectors []domain.MarketDiscoveryConnector) error {
	if len(args) == 0 || args[0] != "discover" {
		return fmt.Errorf("uso: bateia idea research discover --query TEXTO --output arquivo.json")
	}
	fs := flag.NewFlagSet("idea research discover", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	query := fs.String("query", "", "Consulta textual")
	input := fs.String("input", "", "IdeaSpec JSON; alternativo a --query")
	provider := fs.String("provider", "all", "all, zbmath ou joss")
	max := fs.Int("max-per-query", 10, "Máximo por fonte e consulta")
	timeout := fs.Duration("timeout", 90*time.Second, "Prazo total")
	output := fs.String("output", "", "Snapshot JSON obrigatório")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if (*query == "") == (*input == "") || *output == "" || *max < 1 || *max > 25 || *timeout <= 0 {
		return fmt.Errorf("exige exatamente um de --query ou --input, --output, max 1..25 e timeout positivo")
	}
	selected := make([]domain.MarketDiscoveryConnector, 0, 2)
	for _, connector := range connectors {
		id := connector.Capabilities().ID
		if *provider == "all" || *provider == "zbmath" && id == domain.ConnectorZBMath || *provider == "joss" && id == domain.ConnectorJOSS {
			selected = append(selected, connector)
		}
	}
	if len(selected) == 0 || (*provider != "all" && *provider != "zbmath" && *provider != "joss") {
		return fmt.Errorf("provedor de pesquisa desconhecido ou indisponível %q", *provider)
	}
	queries := []string{}
	if *query != "" {
		queries = append(queries, strings.TrimSpace(*query))
	} else {
		idea, err := readIdeaSpec(*input)
		if err != nil {
			return err
		}
		plan, err := usecase.ExpandIdea(idea)
		if err != nil {
			return err
		}
		for _, expansion := range plan.Expansions {
			if len(queries) == 4 {
				break
			}
			if q := strings.TrimSpace(expansion.Query); q != "" {
				queries = append(queries, q)
			}
		}
	}
	if len(queries) == 0 {
		return domain.ErrInvalidRunInput
	}
	for _, q := range queries {
		if len(q) > 256 {
			return domain.ErrInvalidRunInput
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	snapshot := researchSnapshot{SchemaVersion: "1", Queries: queries, Attempts: []researchAttempt{}, Candidates: []domain.ProductCandidate{}}
	seen := map[string]bool{}
	succeeded := 0
	for queryIndex, q := range queries {
		for _, connector := range selected {
			id := connector.Capabilities().ID
			if ctx.Err() != nil {
				break
			}
			discovery, err := connector.Discover(ctx, domain.MarketDiscoveryRequest{QueryID: fmt.Sprintf("research-%d", queryIndex+1), QueryLocale: "en-US", Query: q, Limit: *max})
			attempt := researchAttempt{Query: q, ConnectorID: id, Status: discovery.Status, ResultCount: len(discovery.Candidates), ResponseHash: discovery.ResponseHash}
			if err != nil {
				attempt.Status = domain.AttemptUnknownProvider
				attempt.Error = err.Error()
			} else {
				succeeded++
				for _, candidate := range discovery.Candidates {
					key := candidate.Identity.ConnectorID + "\x00" + candidate.Identity.ExternalID
					if seen[key] {
						continue
					}
					seen[key] = true
					candidate.Origins = []domain.MarketCandidateOrigin{{ConnectorID: id, ConnectorVersion: connector.Capabilities().Version, QueryID: fmt.Sprintf("research-%d", queryIndex+1), QueryLocale: "en-US", ResponseHash: discovery.ResponseHash}}
					snapshot.Candidates = append(snapshot.Candidates, candidate)
				}
			}
			snapshot.Attempts = append(snapshot.Attempts, attempt)
		}
	}
	if err := writeJSONAtomic(*output, snapshot); err != nil {
		return err
	}
	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "pesquisa acadêmica: %d candidatos, snapshot %s\n", len(snapshot.Candidates), *output)
	}
	if succeeded == 0 {
		return fmt.Errorf("nenhuma consulta acadêmica concluída")
	}
	return nil
}
