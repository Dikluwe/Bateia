/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/app.md
 * @layer L2
 * @updated 2026-09-24
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
	funnel         *usecase.SearchFunnel
	recommendation mcp.RecommendationRecommender
	durableRuns    mcp.DurableRecommendationRuns
}

func NewCLIApp(funnel *usecase.SearchFunnel, recommenders ...mcp.RecommendationRecommender) *CLIApp {
	app := &CLIApp{funnel: funnel}
	if len(recommenders) > 0 {
		app.recommendation = recommenders[0]
	}
	return app
}

func (a *CLIApp) Run(args []string) error {
	if len(args) > 0 {
		cmd := strings.ToLower(args[0])
		switch cmd {
		case "similarity", "compare":
			return a.RunSimilarityCmd(args[1:])
		case "search", "funnel":
			return a.runSearchFlags(args[1:])
		case "analyze", "local":
			return a.runSearchFlags(append([]string{"-path", "."}, args[1:]...))
		case "idea":
			return RunIdeaCommand(args[1:], os.Stdout)
		case "algorithm":
			return RunAlgorithmHistoryCommand(args[1:], os.Stdout)
		case "benchmark":
			return RunBenchmarkCommand(args[1:], os.Stdout)
		case "mcp":
			mcpServer := mcp.NewMCPServer(a.funnel, os.Stdin, os.Stdout, a.recommendation).WithDurableRuns(a.durableRuns).WithAdvancedIntake()
			return mcpServer.Serve(context.Background())
		case "tui":
			return RunTUI(a.funnel, "", "go")
		case "help":
			a.PrintUsage()
			return nil
		}
	}
	return a.runSearchFlags(args)
}

func (a *CLIApp) PrintUsage() {
	fmt.Println("  security scan --binary /caminho/binario-go [--output scan.json] Consulta govulncheck instalado; não executa o candidato")
	fmt.Println("  effects verify Verifica offline sort-preserve-v1 em casos limitados; effects --help para limites")
	fmt.Println("🔍 Bateia — Motor de Busca de Código, Análise Arquitetural & Inteligência Ontológica")
	fmt.Println("\nSubcomandos Disponíveis:")
	fmt.Println("  similarity    Busca projetos parecidos, calcula Score de Similaridade e sugere melhorias")
	fmt.Println("  search        Executa o funil tradicional de busca por código, manifestos e necessidades ocultas")
	fmt.Println("  analyze       Analisa uma worktree local, incluindo alterações ainda não enviadas ao GitHub")
	fmt.Println("  idea plan     Materializa offline um plano rastreável de expansão da ideia")
	fmt.Println("  idea expand   Enriquece offline o plano com vocabulário alternativo")
	fmt.Println("  idea evaluate Compara estratégias de expansão num corpus versionado")
	fmt.Println("  idea discover Executa expansões num provedor e grava snapshot rastreável")
	fmt.Println("  idea ontology Expande, avalia e audita a ontologia lexical offline")
	fmt.Println("  idea catalog Descobre projetos no GitHub e executa a peneira offline")
	fmt.Println("  idea discovery evaluate Avalia offline um snapshot com julgamentos")
	fmt.Println("  idea screen   Aplica a peneira explicável a um snapshot")
	fmt.Println("  idea screening calibrate|audit Calibra ou audita a peneira offline")
	fmt.Println("  idea evidence plan|collect|assess|evaluate|diagnose|reconcile|verify|recheck Comprova e audita evidências, sem LLM")
	fmt.Println("  idea decide   Decide offline; `decide frontier` expõe Pareto e `decide evaluate` mede risco-cobertura")
	fmt.Println("  idea recommend Orquestra o funil e ranqueia alternativas com evidência")
	fmt.Println("  idea intake prepare|apply|diagnose Estrutura entradas e abstenções avançadas sem executar efeitos")
	fmt.Println("  idea review protocol|synthesize|update|prioritize|evaluate Mantém revisão auditável e viva sem alterar a busca simples")
	fmt.Println("  idea recover diagnose|plan|value|execute Recupera evidências e prioriza ações por valor decisório")
	fmt.Println("  idea autofagia|katamari Audita capacidades internas e compõe somente lacunas comprovadas")
	fmt.Println("  idea radar export Exporta o radar offline como Archify Dataflow v1")
	fmt.Println("  idea landscape analyze Classifica substitutos, complementos e analogias; compara snapshots")
	fmt.Println("  algorithm history analyze Reconstrói alegações históricas com proveniência")
	fmt.Println("  benchmark validate|mutate|run Audita offline o benchmark ponta a ponta")
	fmt.Println("  mcp           Inicia o Bateia como Servidor MCP (Model Context Protocol) via stdio")
	fmt.Println("  tui           Inicia a interface de terminal interativa TUI")
	fmt.Println("\nExemplo de Uso:")
	fmt.Println("  bateia similarity -target \"gofiber/fiber\" -query \"web framework\" -lang go")
	fmt.Println("  bateia search -query \"quicksort\" -lang go -max 20")
	fmt.Println("  bateia analyze -path . -query \"arquitetura do projeto\"")
	fmt.Println("  bateia idea plan --input ideia.json --output plano.json")
	fmt.Println("  bateia idea expand --input ideia.json --strategy lexical")
	fmt.Println("  bateia idea evaluate --corpus corpus.json --strategy lexical")
	fmt.Println("  bateia idea discover --input ideia.json --strategy lexical --output snapshot.json")
	fmt.Println("  bateia idea ontology expand --input ideia.json --ontology discovery-ontology.json --plan-output plano.json")
	fmt.Println("  bateia idea catalog discover --plan plano.json --provider github --output catalogo.json")
	fmt.Println("  bateia idea decide --idea ideia.json --assessment assessment.json --policy reuse-policy.json")
	fmt.Println("  bateia idea evidence reconcile --input achados-tecnicos.json --as-of 2026-09-24T00:00:00Z")
	fmt.Println("  bateia idea evidence verify --input achados-tecnicos.json --bundles bundle.json --as-of 2026-09-24T00:00:00Z")
	fmt.Println("  bateia idea evidence recheck --input releitura.json --output releitura-report.json --bundles-output bundles-relidos.json")
	fmt.Println("  bateia idea decide frontier --input multicriteria-request.json --output frontier.json")
	fmt.Println("  bateia idea decide evaluate --input selective-decisions.json --output risk-coverage.json")
	fmt.Println("  bateia idea recommend --idea ideia.json --policy reuse-policy.json --target-name Software --target-repo owner/repo --output dossie")
	fmt.Println("  bateia idea intake prepare --input intake-draft.json --output intake.json")
	fmt.Println("  bateia idea review protocol --input review-protocol-draft.json --output review-protocol.json")
	fmt.Println("  bateia idea recover value --input decision-value-request.json --output decision-value-plan.json")
	fmt.Println("  bateia idea radar export --input radar-input.json --output radar.archify.json")
	fmt.Println("  bateia idea landscape analyze --input landscape.json --output landscape-report.json")
	fmt.Println("  bateia algorithm history analyze --input history.json --output history-report.json")
}

func (a *CLIApp) runSearchFlags(args []string) error {
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
	outputFlag := fs.String("output", "relatorio_bateia.md", "Caminho do arquivo de relatório de saída")
	mustHaveFlag := fs.String("must-have", "", "Pastas/estruturas requeridas separadas por vírgula (ex: cmd,internal)")
	localPathFlag := fs.String("path", "", "Caminho de uma worktree local para análise sem GitHub")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *manFlag {
		PrintManPage()
		return nil
	}

	if *mcpFlag {
		mcpServer := mcp.NewMCPServer(a.funnel, os.Stdin, os.Stdout, a.recommendation).WithDurableRuns(a.durableRuns).WithAdvancedIntake()
		return mcpServer.Serve(context.Background())
	}

	if *localPathFlag != "" {
		absolutePath, err := filepath.Abs(*localPathFlag)
		if err != nil {
			return fmt.Errorf("caminho local inválido: %w", err)
		}
		info, err := os.Stat(absolutePath)
		if err != nil {
			return fmt.Errorf("não foi possível acessar a worktree local: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("o caminho local deve ser um diretório: %s", absolutePath)
		}
		*localPathFlag = absolutePath
		if strings.TrimSpace(*queryFlag) == "" {
			*queryFlag = "Análise arquitetural local de " + filepath.Base(absolutePath)
		}
	}

	if *tuiFlag || (*queryFlag == "" && !*mcpFlag && len(args) == 0) {
		return RunTUI(a.funnel, *queryFlag, *langFlag)
	}

	if *queryFlag == "" {
		a.PrintUsage()
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
		LocalPath:   *localPathFlag,
	}
	if err := usecase.ValidateQueryFilter(filter); err != nil {
		return fmt.Errorf("parâmetros inválidos: %w", err)
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
		// Keep full compiler receipts available without repeating them in the main answer.
		audit, auditErr := report.ExportJSONReport(searchReport)
		if auditErr != nil {
			return auditErr
		}
		auditPath := *outputFlag + ".evidence.json"
		if err := os.MkdirAll(filepath.Dir(auditPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(auditPath, audit, 0644); err != nil {
			return err
		}
		auditAbs, err := filepath.Abs(auditPath)
		if err != nil {
			return err
		}
		reportContent = report.GenerateMarkdownReport(searchReport) + "\n[Dados integrais e recibos de compilação](<" + auditAbs + ">)\n"
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
	if len(searchReport.Failures) > 0 || len(searchReport.Warnings) > 0 {
		fmt.Printf("Diagnóstico: %d falhas parciais, %d avisos recuperáveis (detalhes no relatório)\n", len(searchReport.Failures), len(searchReport.Warnings))
	}

	return nil
}

type LogProgressReporter struct{}

func (LogProgressReporter) ReportProgress(event domain.ProgressEvent) {
	log.Printf("[%s] %s", strings.ToUpper(event.Stage), event.Message)
}
