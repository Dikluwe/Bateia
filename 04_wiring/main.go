/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/main.md
 * @layer L4
 * @updated 2026-08-04
 */

package main

import (
	"log"
	"os"
	"strings"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/cli"
	"bateia/03_infra/cache"
	"bateia/03_infra/content"
	"bateia/03_infra/github"
	"bateia/03_infra/lexical"
	localrepo "bateia/03_infra/local"
	"bateia/03_infra/sourcegraph"
)

func loadEnvFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			if os.Getenv(key) == "" && val != "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

func main() {
	// Comandos offline são despachados antes de .env, cache e clientes remotos.
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "idea") {
		lexicalExpander, err := lexical.NewDefault(2)
		if err != nil {
			log.Fatalf("Erro ao carregar vocabulário lexical: %v", err)
		}
		if err := cli.RunIdeaCommand(os.Args[2:], os.Stdout, lexicalExpander); err != nil {
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

	// L4: Injetar adaptadores L3 nas abstrações/interfaces de L1 (Núcleo)
	funnelUsecase := usecase.NewSearchFunnel(sgClient, contentProvider, cacheStore).WithProgressReporter(cli.LogProgressReporter{})

	// L4: Instanciar a Casca L2 (CLI) passando a fiação pronta
	app := cli.NewCLIApp(funnelUsecase)

	// Executar a aplicação com os argumentos passados pelo sistema operacional
	if err := app.Run(os.Args[1:]); err != nil {
		log.Fatalf("Erro na execução da CLI Bateia: %v", err)
	}
}
