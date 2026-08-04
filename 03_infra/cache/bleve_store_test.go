/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/bleve_store.md
 * @layer L3
 * @updated 2026-08-04
 */

package cache_test

import (
	"os"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/03_infra/cache"
)

func TestBleveStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bateia_bleve_test_*")
	if err != nil {
		t.Fatalf("Falha ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := cache.NewBleveStore(tempDir)
	if err != nil {
		t.Fatalf("Erro ao inicializar BleveStore: %v", err)
	}
	defer store.Close()

	cand := domain.RepositoryCandidate{
		Owner:       "gofiber",
		Repo:        "fiber",
		CommitID:    "v2.52.0",
		Stars:       30000,
		Description: "Fast HTTP engine written in Go",
	}

	err = store.SaveDiscovery("fiber", []domain.RepositoryCandidate{cand})
	if err != nil {
		t.Fatalf("Erro ao salvar descoberta no BleveStore: %v", err)
	}

	retrieved, found, err := store.GetDiscovery("fiber", 1*time.Hour)
	if err != nil || !found {
		t.Fatalf("Esperado encontrar descoberta no cache, found=%v, err=%v", found, err)
	}

	if len(retrieved) != 1 || retrieved[0].Repo != "fiber" {
		t.Errorf("Candidato recuperado inconsistente: %+v", retrieved)
	}

	hits, err := store.SearchLocal("HTTP", 10)
	if err != nil {
		t.Fatalf("Erro na busca local Bleve: %v", err)
	}

	if len(hits) == 0 {
		t.Errorf("Esperado encontrar hits no índice Bleve para o termo 'HTTP'")
	}
}
