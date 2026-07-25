/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/funnel.md
 * @layer L1
 * @updated 2026-07-25
 */

package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// SearchFunnel UseCase que implementa a orquestração do funil de redução de 3 estágios em L1.
type SearchFunnel struct {
	searcher domain.SearchProvider
	fetcher  domain.ContentProvider
	store    domain.CacheStore
}

func NewSearchFunnel(searcher domain.SearchProvider, fetcher domain.ContentProvider, store domain.CacheStore) *SearchFunnel {
	return &SearchFunnel{
		searcher: searcher,
		fetcher:  fetcher,
		store:    store,
	}
}

// Execute executa a busca através dos 3 estágios do funil com cache e invariante de redução.
func (f *SearchFunnel) Execute(ctx context.Context, filter domain.QueryFilter) (*domain.SearchReport, error) {
	report := &domain.SearchReport{
		Query:       filter,
		GeneratedAt: time.Now().UTC(),
		Results:     make([]domain.ExtractedRepository, 0),
	}

	queryKey := fmt.Sprintf("%s_%s_%d", filter.Query, filter.Language, filter.MaxResults)

	// ==========================================
	// Estágio 1: Descoberta (Topo do Funil)
	// ==========================================
	log.Printf("[ESTÁGIO 1] Buscando candidatos no topo do funil...")
	candidates, hit, err := f.store.GetDiscovery(queryKey, 2*time.Hour)
	if err != nil {
		log.Printf("[WARNING] Erro ao ler cache de descoberta: %v", err)
	}

	if !hit || len(candidates) == 0 {
		log.Printf("[ESTÁGIO 1] Cache miss/vazio. Consultando provedor de busca...")
		candidates, err = f.searcher.SearchCandidates(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("falha no estágio 1 (Descoberta): %w", err)
		}
		if len(candidates) > 0 {
			_ = f.store.SaveDiscovery(queryKey, candidates)
		}
	} else {
		log.Printf("[ESTÁGIO 1] Cache HIT para a busca '%s' (%d candidatos).", filter.Query, len(candidates))
	}

	report.TopCandidates = len(candidates)
	if len(candidates) == 0 {
		return report, nil
	}

	// ==========================================
	// Estágio 2: Triagem (Meio do Funil)
	// ==========================================
	log.Printf("[ESTÁGIO 2] Triando manifestos e estrutura (%d candidatos)...", len(candidates))
	triaged := make([]domain.TriagedRepository, 0)
	discardedCount := 0

	for _, cand := range candidates {
		ref := domain.CommitRef{
			Owner:  cand.Owner,
			Repo:   cand.Repo,
			Commit: cand.CommitID,
		}

		// 1. Tentar manifesto via Cache
		manifest, hitM, _ := f.store.GetManifest(ref)
		if !hitM {
			manifest, err = f.fetcher.FetchManifest(ctx, cand)
			if err != nil {
				if err.Error() == "MANIFEST_NOT_FOUND (404/410)" {
					cand.Stale = true
					discardedCount++
					continue
				}
				log.Printf("[WARNING] Falha ao buscar manifesto para %s: %v", ref.Key(), err)
			} else {
				_ = f.store.SaveManifest(ref, manifest)
			}
		}

		// 2. Tentar estrutura de diretórios se filtros foram especificados
		tree, err := f.fetcher.FetchTree(ctx, cand)
		score := parser.StructuralScore{Approved: true, Reason: "Sem filtro estrutural"}
		if err == nil && len(filter.MustHave) > 0 {
			score = parser.EvaluateStructure(tree, filter.MustHave)
		}

		if score.Approved {
			triaged = append(triaged, domain.TriagedRepository{
				Candidate:   cand,
				Manifest:    manifest,
				Tree:        tree,
				Approved:    true,
				ScoreReason: score.Reason,
			})
		} else {
			discardedCount++
		}
	}

	report.TriagedCount = len(triaged)
	report.DiscardedCount = discardedCount
	report.DiscardReason = "Filtro estrutural ou repositório indisponível (404/stale)"

	log.Printf("[ESTÁGIO 2] Triagem concluída: %d aprovados de %d candidatos.", len(triaged), len(candidates))

	// ==========================================
	// Estágio 3: Extração (Fundo do Funil)
	// Limita o fundo do funil a no máximo 10 repositórios validados
	// ==========================================
	maxExtract := 10
	if len(triaged) < maxExtract {
		maxExtract = len(triaged)
	}

	log.Printf("[ESTÁGIO 3] Extraindo até %d repositórios finais...", maxExtract)

	for i := 0; i < maxExtract; i++ {
		tRepo := triaged[i]
		ref := domain.CommitRef{
			Owner:  tRepo.Candidate.Owner,
			Repo:   tRepo.Candidate.Repo,
			Commit: tRepo.Candidate.CommitID,
		}

		// Tentar extração via cache imutável (repo@commit)
		extracted, hitE, _ := f.store.GetExtracted(ref)
		if !hitE {
			log.Printf("[ESTÁGIO 3] Baixando conteúdo para %s...", ref.Key())
			extracted, err = f.fetcher.ExtractRepoFiles(ctx, tRepo.Candidate, []string{"go.mod", "README.md"})
			if err != nil {
				log.Printf("[WARNING] Falha na extração de %s: %v", ref.Key(), err)
				continue
			}
			extracted.Manifest = tRepo.Manifest
			_ = f.store.SaveExtracted(ref, extracted)
		} else {
			log.Printf("[ESTÁGIO 3] Cache HIT na extração para %s.", ref.Key())
		}

		report.Results = append(report.Results, *extracted)
	}

	report.ExtractedCount = len(report.Results)
	log.Printf("[CONCLUÍDO] Funil executado com sucesso! %d repositórios entregues no relatório.", report.ExtractedCount)

	return report, nil
}
