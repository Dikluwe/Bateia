/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/cache_store_test.md
 * @layer L3
 * @updated 2026-07-25
 */

package cache

import (
	"os"
	"testing"
	"time"

	"bateia/01_core/domain"
)

func TestDiskCacheStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bateia_cache_test_*")
	if err != nil {
		t.Fatalf("Falha ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewDiskCacheStore(tempDir)
	if err != nil {
		t.Fatalf("Falha ao instanciar DiskCacheStore: %v", err)
	}

	// 1. Testar Discovery Cache com TTL
	queryKey := "go_auth_example"
	candidates := []domain.RepositoryCandidate{
		{Owner: "foo", Repo: "bar", Stars: 100, CommitID: "c1", ConsultadoEm: time.Now().UTC()},
	}

	err = store.SaveDiscovery(queryKey, candidates)
	if err != nil {
		t.Fatalf("SaveDiscovery falhou: %v", err)
	}

	cachedCandidates, hit, err := store.GetDiscovery(queryKey, 1*time.Hour)
	if err != nil || !hit {
		t.Fatalf("GetDiscovery deveria dar hit. Hit: %v, Err: %v", hit, err)
	}

	if len(cachedCandidates) != 1 || cachedCandidates[0].Owner != "foo" {
		t.Errorf("Resultado do cache de descoberta incorreto: %+v", cachedCandidates)
	}

	// 2. Testar Manifest Cache (repo@commit imutável)
	ref := domain.CommitRef{Owner: "foo", Repo: "bar", Commit: "c1"}
	manifest := &domain.Manifest{
		Ref:          ref,
		ModuleName:   "github.com/foo/bar",
		GoVersion:    "1.22",
		ConsultadoEm: time.Now().UTC(),
	}

	err = store.SaveManifest(ref, manifest)
	if err != nil {
		t.Fatalf("SaveManifest falhou: %v", err)
	}

	cachedManifest, hitM, err := store.GetManifest(ref)
	if err != nil || !hitM {
		t.Fatalf("GetManifest deveria dar hit. Hit: %v, Err: %v", hitM, err)
	}

	if cachedManifest.ModuleName != "github.com/foo/bar" {
		t.Errorf("Manifesto cacheados incorreto: %+v", cachedManifest)
	}
}
