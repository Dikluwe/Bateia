/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/fast_cache.md
 * @layer L3
 * @updated 2026-08-08
 */

package cache

import (
	"sync"
	"time"

	"bateia/01_core/domain"
)

// FastCacheStore implementa um cache em memória arena de alta velocidade para L3.
type FastCacheStore struct {
	mu    sync.RWMutex
	items map[string][]byte
}

// NewFastCacheStore cria uma nova instância de FastCacheStore.
func NewFastCacheStore() *FastCacheStore {
	return &FastCacheStore{
		items: make(map[string][]byte),
	}
}

func (s *FastCacheStore) Set(ref domain.CommitRef, key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fullKey := ref.Key() + ":" + key
	buf := make([]byte, len(data))
	copy(buf, data)
	s.items[fullKey] = buf
	return nil
}

func (s *FastCacheStore) Get(ref domain.CommitRef, key string) ([]byte, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fullKey := ref.Key() + ":" + key
	data, ok := s.items[fullKey]
	if !ok {
		return nil, false, nil
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	return buf, true, nil
}

func (s *FastCacheStore) Has(ref domain.CommitRef, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fullKey := ref.Key() + ":" + key
	_, ok := s.items[fullKey]
	return ok, nil
}

func (s *FastCacheStore) SetSearchCache(query string, lang string, candidates []domain.RepositoryCandidate) error {
	return nil
}

func (s *FastCacheStore) GetSearchCache(query string, lang string) ([]domain.RepositoryCandidate, bool, error) {
	return nil, false, nil
}

func (s *FastCacheStore) Close() error {
	return nil
}

func (s *FastCacheStore) Size() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var total int64
	for _, v := range s.items {
		total += int64(len(v))
	}
	return total, nil
}

func (s *FastCacheStore) PurgeOld(maxAge time.Duration) (int, error) {
	return 0, nil
}
