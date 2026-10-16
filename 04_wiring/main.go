/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/main.md
 * @layer L4
 * @updated 2026-10-06
 */

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/cli"
	"bateia/02_shell/mcp"
	"bateia/03_infra/cache"
	"bateia/03_infra/clang"
	"bateia/03_infra/content"
	"bateia/03_infra/durableconfig"
	"bateia/03_infra/effects"
	"bateia/03_infra/envfile"
	"bateia/03_infra/github"
	"bateia/03_infra/lexical"
	localrepo "bateia/03_infra/local"
	"bateia/03_infra/research"
	"bateia/03_infra/sourcegraph"
	"bateia/03_infra/vulnerability"
	webinfra "bateia/03_infra/web"
)

func newMCPServer(funnel *usecase.SearchFunnel, semantic usecase.SemanticExpander, discovery domain.ObservableDiscoveryProvider, content domain.ContentProvider, reader io.Reader, writer io.Writer) *mcp.MCPServer {
	recommendationService := usecase.NewRecommendationService(semantic, discovery, content)
	return mcp.NewMCPServer(funnel, reader, writer, recommendationService)
}

func parseMCPMaxToolCalls(raw string, present bool) (*int, error) {
	if !present {
		return nil, nil
	}
	if raw == "" || raw[0] < '1' || raw[0] > '9' {
		return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
	}
	for index := 1; index < len(raw); index++ {
		if raw[index] < '0' || raw[index] > '9' {
			return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
		}
	}
	parsed, err := strconv.ParseInt(raw, 10, strconv.IntSize)
	if err != nil || parsed <= 0 {
		return nil, fmt.Errorf("BATEIA_MCP_MAX_TOOL_CALLS must be a canonical positive decimal integer")
	}
	limit := int(parsed)
	return &limit, nil
}

func loadEnvFile(filepath string) {
	envfile.Load(filepath)
}

func main() {
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "security") {
		if len(os.Args) < 3 || !strings.EqualFold(os.Args[2], "scan") {
			log.Fatal("Uso: bateia security scan --binary CAMINHO_ABSOLUTO [--output arquivo.json] [--timeout 1m]")
		}
		if err := cli.RunSecurityScanCommand(os.Args[3:], os.Stdout, vulnerability.NewRunner()); err != nil {
			log.Fatalf("Erro no scanner de segurança: %v", err)
		}
		return
	}
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "effects") {
		if err := cli.RunEffectsCommand(os.Args[2:], os.Stdout, effects.NewRunner()); err != nil {
			log.Fatalf("Erro na execução effects: %v", err)
		}
		return
	}
	// Comandos offline são despachados antes de .env, cache e clientes remotos.
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "benchmark") {
		if err := cli.RunBenchmarkCommand(os.Args[2:], os.Stdout); err != nil {
			log.Fatalf("Erro na execução da CLI Bateia: %v", err)
		}
		return
	}
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "algorithm") {
		if err := cli.RunAlgorithmHistoryCommand(os.Args[2:], os.Stdout); err != nil {
			log.Fatalf("Erro na execução da CLI Bateia: %v", err)
		}
		return
	}
	if len(os.Args) > 2 && strings.EqualFold(os.Args[1], "idea") && strings.EqualFold(os.Args[2], "research") {
		if err := cli.RunIdeaResearchCommand(os.Args[3:], os.Stdout, []domain.MarketDiscoveryConnector{research.NewZBMath(nil), research.NewJOSS(nil)}); err != nil {
			log.Fatalf("Erro na pesquisa acadêmica: %v", err)
		}
		return
	}
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "idea") {
		// Trace is entirely offline: no .env, vocabulary, cache or remote clients.
		if len(os.Args) > 2 && strings.EqualFold(os.Args[2], "insights") {
			if err := cli.RunIdeaCommand(os.Args[2:], os.Stdout); err != nil {
				log.Fatalf("Erro na investigação offline de insights: %v", err)
			}
			return
		}
		if len(os.Args) > 2 && strings.EqualFold(os.Args[2], "trace") {
			if err := cli.RunIdeaCommand(os.Args[2:], os.Stdout); err != nil {
				log.Fatalf("Erro na rastreabilidade consultiva: %v", err)
			}
			return
		}
		lexicalExpander, err := lexical.NewDefault(3)
		if err != nil {
			log.Fatalf("Erro ao carregar vocabulário lexical: %v", err)
		}
		if len(os.Args) > 2 && strings.EqualFold(os.Args[2], "discover") {
			sgClient := sourcegraph.NewClient(os.Getenv("SOURCEGRAPH_TOKEN"))
			if err := cli.RunIdeaDiscoverCommand(os.Args[3:], os.Stdout, os.Stderr, lexicalExpander, sgClient); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		if len(os.Args) > 2 && strings.EqualFold(os.Args[2], "catalog") {
			loadEnvFile(".env")
			ghClient := github.NewClient(os.Getenv("GITHUB_TOKEN"))
			if err := cli.RunIdeaCatalogCommand(os.Args[3:], os.Stdout, ghClient, github.NewReadmeSurface(ghClient), github.NewTopicsSurface(ghClient)); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		if len(os.Args) > 3 && strings.EqualFold(os.Args[2], "evidence") && strings.EqualFold(os.Args[3], "recheck") {
			loadEnvFile(".env")
			provider := github.NewClient(os.Getenv("GITHUB_TOKEN"))
			if err := cli.RunIdeaEvidenceRecheckCommand(os.Args[4:], os.Stdout, provider); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		if len(os.Args) > 3 && strings.EqualFold(os.Args[2], "evidence") && strings.EqualFold(os.Args[3], "collect") {
			loadEnvFile(".env")
			provider := content.NewRouter(github.NewClient(os.Getenv("GITHUB_TOKEN")), localrepo.NewClient())
			if err := cli.RunIdeaEvidenceCollectCommand(os.Args[4:], os.Stdout, provider); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		if len(os.Args) > 2 && strings.EqualFold(os.Args[2], "recommend") {
			loadEnvFile(".env")
			cli.SetSnapshotProviderFactory(func(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
				return localrepo.NewSnapshotProvider(fallback, bindings)
			})
			sgClient := sourcegraph.NewClient(os.Getenv("SOURCEGRAPH_TOKEN"))
			provider := content.NewRouter(github.NewClient(os.Getenv("GITHUB_TOKEN")), localrepo.NewClient())
			recommendationService := usecase.NewRecommendationService(lexicalExpander, sgClient, provider)
			if err := cli.RunIdeaRecommendCommand(os.Args[3:], os.Stdout, lexicalExpander, sgClient, provider, recommendationService); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		if len(os.Args) > 3 && strings.EqualFold(os.Args[2], "recover") && strings.EqualFold(os.Args[3], "execute") {
			loadEnvFile(".env")
			provider := content.NewRouter(github.NewClient(os.Getenv("GITHUB_TOKEN")), localrepo.NewClient())
			if err := cli.RunIdeaRecoverCommand(os.Args[3:], os.Stdout, provider); err != nil {
				log.Fatalf("Erro na execução da CLI Bateia: %v", err)
			}
			return
		}
		noToken := false
		for _, arg := range os.Args[2:] {
			if arg == "--no-token" {
				noToken = true
			}
		}
		if !noToken {
			loadEnvFile(".env")
		}
		token := ""
		if !noToken {
			token = os.Getenv("GITHUB_TOKEN")
		}
		webProvider := webinfra.NewObservableWebClient(nil, token)
		if err := cli.RunIdeaCommandWithWeb(os.Args[2:], os.Stdout, lexicalExpander, webProvider); err != nil {
			log.Fatalf("Erro na execução da CLI Bateia: %v", err)
		}
		return
	}

	// L4: Ler tokens do arquivo .env se disponíveis
	loadEnvFile(".env")

	githubToken := os.Getenv("GITHUB_TOKEN")
	sourcegraphToken := os.Getenv("SOURCEGRAPH_TOKEN")

	// L4: Instanciar dependências concretas de L3 (Infraestrutura com suporte Bleve local)
	store, err := cache.NewBleveStore("")
	var cacheStore domain.CacheStore
	if err != nil {
		log.Printf("Aviso: Falha ao inicializar BleveStore (%v). Usando DiskCacheStore tradicional.", err)
		diskStore, diskErr := cache.NewDiskCacheStore("")
		if diskErr != nil {
			log.Fatalf("Erro ao inicializar cache em disco L3: %v", diskErr)
		}
		cacheStore = diskStore
	} else {
		defer store.Close()
		cacheStore = store
	}

	sgClient := sourcegraph.NewClient(sourcegraphToken)
	ghClient := github.NewClient(githubToken)
	contentProvider := content.NewRouter(ghClient, localrepo.NewClient())
	lexicalExpander, err := lexical.NewDefault(3)
	if err != nil {
		log.Fatalf("Erro ao carregar vocabulário lexical: %v", err)
	}

	// L4: Injetar adaptadores L3 nas abstrações/interfaces de L1 (Núcleo)
	funnelUsecase := usecase.NewSearchFunnel(sgClient, contentProvider, cacheStore).WithProgressReporter(cli.LogProgressReporter{}).WithSemanticProvider(clang.NewClient("clang").WithFlowTool(os.Getenv("BATEIA_CLANG_FLOW_TOOL")))
	recommendationService := usecase.NewRecommendationService(lexicalExpander, sgClient, contentProvider)
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "mcp") {
		rawToolCallLimit, toolCallLimitPresent := os.LookupEnv("BATEIA_MCP_MAX_TOOL_CALLS")
		toolCallLimit, toolCallLimitErr := parseMCPMaxToolCalls(rawToolCallLimit, toolCallLimitPresent)
		if toolCallLimitErr != nil {
			log.Fatalf("Erro ao inicializar orçamento MCP: %v", toolCallLimitErr)
		}
		toolCallBudget, toolCallBudgetErr := mcp.NewToolCallBudget(toolCallLimit)
		if toolCallBudgetErr != nil {
			log.Fatalf("Erro ao inicializar orçamento MCP: %v", toolCallBudgetErr)
		}
		config, configErr := durableconfig.LoadLocal()
		if configErr != nil {
			log.Fatalf("Erro ao inicializar configuração MCP durável: %v", configErr)
		}
		connectors, connectorErr := defaultMarketDiscoveryConnectors()
		if connectorErr != nil {
			log.Fatalf("Erro ao inicializar conectores públicos: %v", connectorErr)
		}
		mechanismCatalog, mechanismCatalogSHA256 := usecase.ActiveMechanismRuleCatalog()
		server, shutdown, serverErr := NewImprovementSynthesisMCPServer(funnelUsecase, recommendationService, connectors, mechanismCatalog, mechanismCatalogSHA256, nil, config.Root, config.Secret, config.Limits, config.Workers, config.Capacity, os.Stdin, os.Stdout)
		if serverErr != nil {
			log.Fatalf("Erro ao inicializar descoberta de mercado MCP: %v", serverErr)
		}
		server.WithAdvancedIntake()
		server.WithToolCallBudget(toolCallBudget)
		serveErr := server.Serve(context.Background())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		shutdownErr := shutdown(ctx)
		cancel()
		if serveErr != nil {
			log.Fatalf("Erro na execução MCP: %v", serveErr)
		}
		if shutdownErr != nil {
			log.Fatalf("Erro ao encerrar execuções MCP duráveis: %v", shutdownErr)
		}
		return
	}

	// L4: Instanciar a Casca L2 (CLI) passando a fiação pronta.
	app := cli.NewCLIApp(funnelUsecase, recommendationService)
	var shutdownDurable func(context.Context) error
	if durableMCPRequested(os.Args[1:]) {
		config, configErr := durableconfig.LoadLocal()
		if configErr != nil {
			log.Fatalf("Erro ao inicializar configuração MCP durável: %v", configErr)
		}
		durableRuns, shutdown, durableErr := newDurableRecommendationRuns(
			recommendationService, config.Root, config.Secret, config.Limits, config.Workers, config.Capacity,
		)
		if durableErr != nil {
			log.Fatalf("Erro ao inicializar execuções MCP duráveis: %v", durableErr)
		}
		app.WithDurableRuns(durableRuns)
		shutdownDurable = shutdown
	}

	// Executar a aplicação com os argumentos passados pelo sistema operacional
	runErr := app.Run(os.Args[1:])
	var shutdownErr error
	if shutdownDurable != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		shutdownErr = shutdownDurable(ctx)
		cancel()
	}
	if runErr != nil {
		log.Fatalf("Erro na execução da CLI Bateia: %v", runErr)
	}
	if shutdownErr != nil {
		log.Fatalf("Erro ao encerrar execuções MCP duráveis: %v", shutdownErr)
	}
}

func durableMCPRequested(args []string) bool {
	if len(args) > 0 && strings.EqualFold(args[0], "mcp") {
		return true
	}
	for _, arg := range args {
		if strings.EqualFold(arg, "-mcp") || strings.EqualFold(arg, "--mcp") {
			return true
		}
	}
	return false
}
