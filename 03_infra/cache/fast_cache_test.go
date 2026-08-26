/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/fast_cache_test.md
 * @layer L3
 * @updated 2026-08-08
 */

package cache_test

import (
	"sync"
	"testing"

	"bateia/01_core/domain"
	"bateia/03_infra/cache"
)

func TestFastCacheStore_SetGetRoundTrip(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "test", Repo: "repo", Commit: "abc123"}

	data := []byte(`{"module": "github.com/test/repo"}`)
	if err := store.Set(ref, "go.mod", data); err != nil {
		t.Fatalf("Set falhou: %v", err)
	}

	got, found, err := store.Get(ref, "go.mod")
	if err != nil {
		t.Fatalf("Get falhou: %v", err)
	}
	if !found {
		t.Fatal("Esperava encontrar item no cache, mas não encontrou")
	}
	if string(got) != string(data) {
		t.Errorf("Conteúdo diverge: got=%q, want=%q", got, data)
	}
}

func TestFastCacheStore_GetMiss(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "miss", Repo: "nope"}

	_, found, err := store.Get(ref, "inexistente")
	if err != nil {
		t.Fatalf("Get inesperado: %v", err)
	}
	if found {
		t.Error("Esperava miss, mas obteve hit")
	}
}

func TestFastCacheStore_Has(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "x", Repo: "y", Commit: "z"}

	exists, _ := store.Has(ref, "k")
	if exists {
		t.Error("Has deveria retornar false antes do Set")
	}

	_ = store.Set(ref, "k", []byte("v"))
	exists, _ = store.Has(ref, "k")
	if !exists {
		t.Error("Has deveria retornar true após o Set")
	}
}

func TestFastCacheStore_Size(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "s", Repo: "r"}

	size, _ := store.Size()
	if size != 0 {
		t.Errorf("Size vazio esperado 0, recebeu: %d", size)
	}

	_ = store.Set(ref, "a", []byte("12345"))
	_ = store.Set(ref, "b", []byte("67890"))
	size, _ = store.Size()
	if size != 10 {
		t.Errorf("Size esperado 10, recebeu: %d", size)
	}
}

func TestFastCacheStore_IsolationCopy(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "iso", Repo: "test"}

	original := []byte("original")
	_ = store.Set(ref, "data", original)

	// Mutação do slice original NÃO deve afetar o cache
	original[0] = 'X'

	got, _, _ := store.Get(ref, "data")
	if got[0] == 'X' {
		t.Error("Cache não isolou a cópia no Set — mutação externa propagou")
	}

	// Mutação do slice retornado NÃO deve afetar o cache
	got[0] = 'Y'
	got2, _, _ := store.Get(ref, "data")
	if got2[0] == 'Y' {
		t.Error("Cache não isolou a cópia no Get — mutação externa propagou")
	}
}

func TestFastCacheStore_ConcurrentAccess(t *testing.T) {
	store := cache.NewFastCacheStore()
	ref := domain.CommitRef{Owner: "conc", Repo: "test"}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := "key"
			_ = store.Set(ref, key, []byte("value"))
			_, _, _ = store.Get(ref, key)
			_, _ = store.Has(ref, key)
		}(i)
	}
	wg.Wait()
}

func TestFastCacheStore_CloseAndPurge(t *testing.T) {
	store := cache.NewFastCacheStore()

	if err := store.Close(); err != nil {
		t.Errorf("Close não deveria falhar: %v", err)
	}

	purged, err := store.PurgeOld(0)
	if err != nil {
		t.Errorf("PurgeOld não deveria falhar: %v", err)
	}
	if purged != 0 {
		t.Errorf("PurgeOld deveria retornar 0: %d", purged)
	}
}
