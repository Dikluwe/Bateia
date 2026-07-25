/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/main.md
 * @layer L4
 * @updated 2026-07-25
 */

package main

import (
	"log"
	"os"

	"bateia/01_core/usecase"
	"bateia/02_shell/cli"
	"bateia/03_infra/cache"
	"bateia/03_infra/github"
	"bateia/03_infra/sourcegraph"
)

func main() {
	// L4: Ler tokens de ambiente se disponíveis
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
