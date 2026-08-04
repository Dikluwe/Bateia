/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/cache/bleve_store.md
 * @layer L3
 * @updated 2026-08-04
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

	"github.com/blevesearch/bleve/v2"
)

type BleveStore struct {
	baseStore *DiskCacheStore
	index     bleve.Index
	mu        sync.RWMutex
}

type bleveDoc struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "candidate", "manifest", "extracted"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// NewBleveStore inicializa o repositório de cache e o índice Bleve embarcado em disco.
func NewBleveStore(baseDir string) (*BleveStore, error) {
	diskStore, err := NewDiskCacheStore(baseDir)
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(diskStore.baseDir, "bleve_index")
	var index bleve.Index

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf("falha ao criar índice Bleve em %s: %w", indexPath, err)
		}
	} else {
		index, err = bleve.Open(indexPath)
		if err != nil {
			// Se o índice estiver corrompido, exclui e recria
			_ = os.RemoveAll(indexPath)
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New(indexPath, mapping)
			if err != nil {
				return nil, fmt.Errorf("falha ao recriar índice Bleve: %w", err)
			}
		}
	}

	return &BleveStore{
		baseStore: diskStore,
		index:     index,
	}, nil
}

func (b *BleveStore) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.index != nil {
		return b.index.Close()
	}
	return nil
}

func (b *BleveStore) GetDiscovery(queryKey string, ttl time.Duration) ([]domain.RepositoryCandidate, bool, error) {
	return b.baseStore.GetDiscovery(queryKey, ttl)
}

func (b *BleveStore) SaveDiscovery(queryKey string, candidates []domain.RepositoryCandidate) error {
	err := b.baseStore.SaveDiscovery(queryKey, candidates)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, cand := range candidates {
		bytes, _ := json.Marshal(cand)
		docID := fmt.Sprintf("cand_%s_%s_%s", cand.Owner, cand.Repo, cand.CommitID)
		doc := bleveDoc{
			ID:        docID,
			Type:      "candidate",
			Content:   string(bytes),
			Timestamp: time.Now().UTC(),
		}
		_ = b.index.Index(doc.ID, doc)
	}

	return nil
}

func (b *BleveStore) GetManifest(ref domain.CommitRef) (*domain.Manifest, bool, error) {
	return b.baseStore.GetManifest(ref)
}

func (b *BleveStore) SaveManifest(ref domain.CommitRef, manifest *domain.Manifest) error {
	err := b.baseStore.SaveManifest(ref, manifest)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	bytes, _ := json.Marshal(manifest)
	doc := bleveDoc{
		ID:        "man_" + ref.Key(),
		Type:      "manifest",
		Content:   string(bytes),
		Timestamp: time.Now().UTC(),
	}
	_ = b.index.Index(doc.ID, doc)
	return nil
}

func (b *BleveStore) GetExtracted(ref domain.CommitRef) (*domain.ExtractedRepository, bool, error) {
	return b.baseStore.GetExtracted(ref)
}

func (b *BleveStore) SaveExtracted(ref domain.CommitRef, extracted *domain.ExtractedRepository) error {
	err := b.baseStore.SaveExtracted(ref, extracted)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	bytes, _ := json.Marshal(extracted)
	doc := bleveDoc{
		ID:        "ext_" + ref.Key(),
		Type:      "extracted",
		Content:   string(bytes),
		Timestamp: time.Now().UTC(),
	}
	_ = b.index.Index(doc.ID, doc)
	return nil
}

// SearchLocal realiza uma busca por palavra-chave diretamente no índice Bleve local offline.
func (b *BleveStore) SearchLocal(term string, limit int) ([]string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	query := bleve.NewMatchQuery(term)
	searchRequest := bleve.NewSearchRequestOptions(query, limit, 0, false)
	searchResult, err := b.index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("erro na busca Bleve: %w", err)
	}

	var results []string
	for _, hit := range searchResult.Hits {
		results = append(results, hit.ID)
	}
	return results, nil
}
