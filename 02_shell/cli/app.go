/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/app.md
 * @layer L2
 * @updated 2026-07-25
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
	"bateia/01_core/parser"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/02_shell/report"
)

type CLIApp struct {
	funnel *usecase.SearchFunnel
}

func NewCLIApp(funnel *usecase.SearchFunnel) *CLIApp {
	return &CLIApp{funnel: funnel}
}

func (a *CLIApp) Run(args []string) error {
	fs := flag.NewFlagSet("bateia", flag.ContinueOnError)

	queryFlag := fs.String("query", "", "Query de busca por arquitetura ou padrão de código (obrigatório)")
	langFlag := fs.String("lang", "go", "Linguagem alvo (go, typescript/ts/node, python/py, rust/rs, java/kotlin/kt, php, c/cpp/c++, csharp/cs/dotnet, julia, r, lua, ocaml, elixir, ruby, swift, haskell, dart/flutter, scala, clojure/clj, zig, pascal/delphi, fortran, basic/vb, fsharp/fs, perl, groovy, cobol, nim, crystal, v/vlang, erlang, d, gleam, coldfusion/cf, actionscript/flex, ada, mojo, objective-c/objc, haxe, raku, racket, chapel, pony, odin)")
	maxFlag := fs.Int("max", 20, "Número de candidatos para descoberta no topo do funil")
	depthFlag := fs.Int("depth", 1, "Profundidade da busca recursiva em grafo de dependências (default: 1)")
	reposFlag := fs.String("repos", "", "Lista separada por vírgulas de repositórios (owner/repo) para Modo Benchmark direto. Ignora o estágio 1 e a flag -lang.")
	mcpFlag := fs.Bool("mcp", false, "Inicia o Bateia como um Servidor MCP (Model Context Protocol) via stdio")
	tuiFlag := fs.Bool("tui", false, "Inicia o Bateia na interface interativa TUI (Terminal User Interface)")
	labFlag := fs.Bool("lab", false, "Exporta o resultado no formato de Dataset de Algoritmos & Testes para o Lab de Projetos")
	manFlag := fs.Bool("man", false, "Exibe a manpage UNIX bateia(1)")
	formatFlag := fs.String("format", "markdown", "Formato do relatório de saída: markdown, html, json, sqlite ou lab")
	outputFlag := fs.String("output", "00_nucleo/relatorios/relatorio_bateia.md", "Caminho do arquivo de relatório de saída")
	mustHaveFlag := fs.String("must-have", "", "Pastas/estruturas requeridas separadas por vírgula (ex: cmd,internal)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *manFlag {
		PrintManPage()
		return nil
	}

	if *mcpFlag {
		mcpServer := mcp.NewMCPServer(a.funnel, os.Stdin, os.Stdout)
		return mcpServer.Serve(context.Background())
	}

	if *tuiFlag || (*queryFlag == "" && !*mcpFlag && len(args) == 0) {
		return RunTUI(a.funnel, *queryFlag, *langFlag)
	}

	if *queryFlag == "" {
		fmt.Println("Bateia — Motor de Busca de Código para Análise Arquitetural & Lab de Projetos")
		fmt.Println("Uso: bateia -query \"quicksort\" [-lang go] [-max 20] [-lab] [-output relatorios/lab_quicksort.json]")
		fmt.Println("Ou modo TUI interativo: bateia -tui")
		fmt.Println("Ou modo MCP: bateia -mcp")
		fs.PrintDefaults()
		return fmt.Errorf("parâmetro -query é obrigatório quando não estiver nos modos -tui ou -mcp")
	}

	var mustHave []string
	if *mustHaveFlag != "" {
		parts := strings.Split(*mustHaveFlag, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				mustHave = append(mustHave, trimmed)
			}
		}
	}

	var directRepos []string
	if *reposFlag != "" {
		parts := strings.Split(*reposFlag, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				directRepos = append(directRepos, trimmed)
			}
		}
	}

	filter := domain.QueryFilter{
		Query:       *queryFlag,
		Language:    *langFlag,
		MaxResults:  *maxFlag,
		MustHave:    mustHave,
		Depth:       *depthFlag,
		DirectRepos: directRepos,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	log.Printf("Iniciando busca no Bateia para: '%s' (linguagem: %s)...", filter.Query, filter.Language)
	searchReport, err := a.funnel.Execute(ctx, filter)
	if err != nil {
		return fmt.Errorf("erro na execução do funil de busca: %w", err)
	}

	// Gerar Relatório (Markdown, HTML, JSON, SQLite ou Lab)
	var reportContent string
	fmtLower := strings.ToLower(*formatFlag)
	outLower := strings.ToLower(*outputFlag)

	if *labFlag || fmtLower == "lab" || strings.HasSuffix(outLower, ".lab.json") {
		labDataset := parser.ExtractLabDataset(searchReport)
		if strings.HasSuffix(outLower, ".md") {
			reportContent = report.GenerateLabMarkdown(labDataset)
		} else {
			jsonBytes, err := report.ExportLabJSON(labDataset)
			if err != nil {
				return fmt.Errorf("erro ao exportar Lab JSON: %w", err)
			}
			reportContent = string(jsonBytes)
		}
	} else if fmtLower == "json" || strings.HasSuffix(outLower, ".json") {
		jsonBytes, err := report.ExportJSONReport(searchReport)
		if err != nil {
			return fmt.Errorf("erro ao exportar JSON: %w", err)
		}
		reportContent = string(jsonBytes)
	} else if fmtLower == "sqlite" || fmtLower == "sql" || strings.HasSuffix(outLower, ".sql") || strings.HasSuffix(outLower, ".sqlite") {
		reportContent = report.ExportSQLiteSchema(searchReport)
	} else if fmtLower == "adr" {
		reportContent = parser.GenerateADRDocument(searchReport)
	} else if fmtLower == "html" || strings.HasSuffix(outLower, ".html") {
		reportContent = report.GenerateHTMLReport(searchReport)
	} else {
		reportContent = report.GenerateMarkdownReport(searchReport)
	}

	_ = os.MkdirAll(filepath.Dir(*outputFlag), 0755)
	err = os.WriteFile(*outputFlag, []byte(reportContent), 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar o relatório em %s: %w", *outputFlag, err)
	}

	absPath, _ := filepath.Abs(*outputFlag)
	fmt.Printf("\nBusca concluída! Relatório gerado com sucesso em: %s\n", absPath)
	fmt.Printf("Funil: %d descobertos -> %d triados -> %d extraídos (%d descartados: %s)\n",
		searchReport.TopCandidates, searchReport.TriagedCount, searchReport.ExtractedCount, searchReport.DiscardedCount, searchReport.DiscardReason)

	return nil
}
