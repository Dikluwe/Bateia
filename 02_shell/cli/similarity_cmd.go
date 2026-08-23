/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/subcommands.md
 * @layer L2
 * @updated 2026-08-08
 */

package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/report"
)

// RunSimilarityCmd executa o subcomando de Ontologia e Score de Similaridade de Ideias.
func (a *CLIApp) RunSimilarityCmd(args []string) error {
	fs := flag.NewFlagSet("similarity", flag.ContinueOnError)

	targetFlag := fs.String("target", "", "Repositório alvo de comparação (owner/repo). Ignora -target-idea se informado.")
	targetIdeaFlag := fs.String("target-idea", "", "Descrição em texto da ideia/problema alvo para buscar projetos parecidos")
	queryFlag := fs.String("query", "", "Query de busca por repositórios candidatos")
	langFlag := fs.String("lang", "go", "Linguagem alvo (ex: go, typescript, python, rust ou 'go,typescript,python,rust')")
	excludeLangFlag := fs.String("exclude-lang", "", "Linguagens para excluir da varredura poliglota (ex: 'go,typescript')")
	sweepFlag := fs.Bool("sweep", false, "Executa a varredura poliglota em múltiplos ecossistemas")
	maxFlag := fs.Int("max", 15, "Número máximo de candidatos no topo do funil")
	outputFlag := fs.String("output", "relatorios/similaridade.md", "Arquivo de saída para o relatório de similaridade ontológica")

	if err := fs.Parse(args); err != nil {
		return err
	}

	selectedLang := *langFlag
	if *sweepFlag {
		selectedLang = "sweep"
	}

	if *excludeLangFlag != "" {
		filteredLangs := usecase.FilterLanguages(usecase.ParseLanguageList(selectedLang), *excludeLangFlag)
		selectedLang = strings.Join(filteredLangs, ",")
	}

	searchQuery := *queryFlag
	if searchQuery == "" {
		if *targetIdeaFlag != "" {
			searchQuery = *targetIdeaFlag
		} else if *targetFlag != "" {
			searchQuery = *targetFlag
		} else {
			return fmt.Errorf("informe ao menos --query, --target ou --target-idea para o comando similarity")
		}
	}

	var directRepos []string
	if *targetFlag != "" {
		directRepos = append(directRepos, *targetFlag)
	}

	filter := domain.QueryFilter{
		Query:       searchQuery,
		Language:    selectedLang,
		MaxResults:  *maxFlag,
		Depth:       1,
		DirectRepos: directRepos,
	}
	if err := usecase.ValidateQueryFilter(filter); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	log.Printf("[SIMILARITY ENGINE] Iniciando análise ontológica para query: '%s' (Alvo: %s)...", searchQuery, *targetFlag)
	searchReport, err := a.funnel.Execute(ctx, filter)
	if err != nil {
		return fmt.Errorf("falha ao executar análise de similaridade: %w", err)
	}

	reportContent := report.GenerateMarkdownReport(searchReport)

	_ = os.MkdirAll(filepath.Dir(*outputFlag), 0755)
	if err := os.WriteFile(*outputFlag, []byte(reportContent), 0644); err != nil {
		return fmt.Errorf("falha ao salvar relatório em %s: %w", *outputFlag, err)
	}

	log.Printf("[SIMILARITY ENGINE] Relatório de similaridade e recomendações salvo em: %s", *outputFlag)
	return nil
}
