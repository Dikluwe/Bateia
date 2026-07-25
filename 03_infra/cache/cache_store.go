/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/cache_store.md
 * @layer L3
 * @updated 2026-07-25
 */

package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"bateia/01_core/domain"
)

type cacheEntry[T any] struct {
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}

type DiskCacheStore struct {
	baseDir string
	mu      sync.RWMutex
}

func NewDiskCacheStore(baseDir string) (*DiskCacheStore, error) {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			baseDir = ".bateia_cache"
		} else {
			baseDir = filepath.Join(home, ".cache", "bateia")
		}
	}

	subdirs := []string{"discovery", "manifests", "extracted"}
	for _, sub := range subdirs {
		dir := filepath.Join(baseDir, sub)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretório de cache %s: %w", dir, err)
		}
	}

	return &DiskCacheStore{baseDir: baseDir}, nil
}

func sanitizeFilename(name string) string {
	var result []rune
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result = append(result, r)
		} else {
			result = append(result, '_')
		}
	}
	return string(result)
}

// GetDiscovery recupera candidatos de busca do cache se válidos dentro do TTL.
func (d *DiskCacheStore) GetDiscovery(queryKey string, ttl time.Duration) ([]domain.RepositoryCandidate, bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fileName := sanitizeFilename(queryKey) + ".json"
	path := filepath.Join(d.baseDir, "discovery", fileName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry cacheEntry[[]domain.RepositoryCandidate]
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, nil
	}

	if time.Since(entry.Timestamp) > ttl {
		return nil, false, nil
	}

	return entry.Data, true, nil
}

// SaveDiscovery armazena os candidatos de descoberta no cache.
func (d *DiskCacheStore) SaveDiscovery(queryKey string, candidates []domain.RepositoryCandidate) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileName := sanitizeFilename(queryKey) + ".json"
	path := filepath.Join(d.baseDir, "discovery", fileName)

	entry := cacheEntry[[]domain.RepositoryCandidate]{
		Timestamp: time.Now().UTC(),
		Data:      candidates,
	}

	bytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

// GetManifest recupera o manifesto cacheados por repo@commit (imutável).
func (d *DiskCacheStore) GetManifest(ref domain.CommitRef) (*domain.Manifest, bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fileName := sanitizeFilename(ref.Key()) + ".json"
	path := filepath.Join(d.baseDir, "manifests", fileName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry cacheEntry[domain.Manifest]
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, nil
	}

	return &entry.Data, true, nil
}

// SaveManifest salva o manifesto por repo@commit.
func (d *DiskCacheStore) SaveManifest(ref domain.CommitRef, manifest *domain.Manifest) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileName := sanitizeFilename(ref.Key()) + ".json"
	path := filepath.Join(d.baseDir, "manifests", fileName)

	entry := cacheEntry[domain.Manifest]{
		Timestamp: time.Now().UTC(),
		Data:      *manifest,
	}

	bytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

// GetExtracted recupera o resultado da extração por repo@commit (imutável).
func (d *DiskCacheStore) GetExtracted(ref domain.CommitRef) (*domain.ExtractedRepository, bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	fileName := sanitizeFilename(ref.Key()) + ".json"
	path := filepath.Join(d.baseDir, "extracted", fileName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry cacheEntry[domain.ExtractedRepository]
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, nil
	}

	return &entry.Data, true, nil
}

// SaveExtracted salva o repositório extraído por repo@commit.
func (d *DiskCacheStore) SaveExtracted(ref domain.CommitRef, extracted *domain.ExtractedRepository) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	fileName := sanitizeFilename(ref.Key()) + ".json"
	path := filepath.Join(d.baseDir, "extracted", fileName)

	entry := cacheEntry[domain.ExtractedRepository]{
		Timestamp: time.Now().UTC(),
		Data:      *extracted,
	}

	bytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
