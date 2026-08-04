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

	"bateia/01_core/usecase"
	"bateia/02_shell/cli"
	"bateia/03_infra/cache"
	"bateia/03_infra/github"
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
	// L4: Ler tokens do arquivo .env se disponíveis
	loadEnvFile(".env")

	githubToken := os.Getenv("GITHUB_TOKEN")
	sourcegraphToken := os.Getenv("SOURCEGRAPH_TOKEN")

	// L4: Instanciar dependências concretas de L3 (Infraestrutura)
	store, err := cache.NewDiskCacheStore("")
	if err != nil {
		log.Fatalf("Erro ao inicializar cache em disco L3: %v", err)
	}

	sgClient := sourcegraph.NewClient(sourcegraphToken)
	ghClient := github.NewClient(githubToken)

	// L4: Injetar adaptadores L3 nas abstrações/interfaces de L1 (Núcleo)
	funnelUsecase := usecase.NewSearchFunnel(sgClient, ghClient, store)

	// L4: Instanciar a Casca L2 (CLI) passando a fiação pronta
	app := cli.NewCLIApp(funnelUsecase)

	// Executar a aplicação com os argumentos passados pelo sistema operacional
	if err := app.Run(os.Args[1:]); err != nil {
		log.Fatalf("Erro na execução da CLI Bateia: %v", err)
	}
}
