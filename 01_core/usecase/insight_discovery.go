/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_discovery.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"errors"
)

func DiscoverInsightReferences(ctx context.Context, request domain.InsightDiscoveryRequest, catalog domain.ProjectCatalog, reader domain.InsightReferenceReader) (domain.InsightDiscoveryReport, error) {
	if ctx == nil || catalog == nil || reader == nil {
		return domain.InsightDiscoveryReport{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.InsightDiscoveryReport{}, err
	}
	plan, err := PlanInsightDiscovery(request)
	if err != nil {
		return domain.InsightDiscoveryReport{}, err
	}
	capability := catalog.Capabilities()
	if capability.Provider == "" || !domain.SupportsRole(capability, domain.RoleProjectCatalog) {
		return domain.InsightDiscoveryReport{}, domain.ErrInvalidRunInput
	}
	ctx, cancel, err := domain.ExecutionTimeout(ctx, plan.Request.TotalTimeout)
	if err != nil {
		return domain.InsightDiscoveryReport{}, err
	}
	defer cancel()
	report := domain.InsightDiscoveryReport{SchemaVersion: "insight-discovery-v1", Plan: plan, Provider: capability, Attempts: []domain.InsightDiscoveryAttempt{}, Candidates: []domain.InsightDiscoveredReference{}, Readmes: []domain.InsightReadmeReceipt{}, Bundles: []domain.EvidenceBundle{}, Gaps: []string{}, NoveltyState: "NOT_ASSESSED", EffectState: "UNVERIFIED", Limitations: []string{
		"Consultas e seleção usam heurísticas lexicais e cobertura de queries/autores, sem classificar relevância semântica, novidade ou domínio.",
		"Este recorte consulta o catálogo GitHub. Presença em catálogo e README não comprovam capacidade, comportamento executado ou efeito no alvo.",
		"As consultas não são traduzidas automaticamente; falha ou vazio não prova ausência de abordagens.",
		"A coleta cobre somente um README de raiz por referência em commit declarado, até 32 KiB admitidos; os limites de transporte do provedor também se aplicam.",
		"Budgets contam operações lógicas; o provedor pode realizar requisições internas para resolver metadados e revisões sob o deadline.",
	}}
	validSearches, failed := 0, false
	for _, query := range plan.Queries {
		attempt := domain.InsightDiscoveryAttempt{QueryID: query.ID, Status: "NOT_STARTED", Attempt: domain.CatalogAttempt{ExpansionID: query.ID, Query: query.Query, Status: "NOT_STARTED"}}
		if ctx.Err() != nil {
			report.Attempts = append(report.Attempts, attempt)
			continue
		}
		operation, stop, timeoutErr := domain.ExecutionTimeout(ctx, plan.Request.AttemptTimeout)
		if timeoutErr != nil {
			return domain.InsightDiscoveryReport{}, timeoutErr
		}
		result, searchErr := catalog.SearchCatalog(operation, domain.CatalogSearchRequest{MetadataOnly: plan.Request.QueryPolicy == "context-v2", DeferRevisionResolution: plan.Request.QueryPolicy == "context-v2", ExpansionID: query.ID, Query: query.Query, Discriminators: append([]string{}, query.AnchorTerms...), MaxResults: plan.Request.MaxPerQuery, AttemptTimeout: plan.Request.AttemptTimeout})
		operationErr := operation.Err()
		stop()
		attempt.ReportedCount = len(result.Candidates)
		attempt.Attempt.Duration = result.Attempt.Duration
		if searchErr != nil || operationErr != nil {
			failed = true
			attempt.Status = "FAILED"
			attempt.ErrorCode = "PROVIDER"
			if result.Attempt.Status == "RATE_LIMITED" {
				attempt.Status = "RATE_LIMITED"
				attempt.ErrorCode = "RATE_LIMIT"
			}
			if errors.Is(searchErr, context.DeadlineExceeded) || errors.Is(operationErr, context.DeadlineExceeded) {
				attempt.Status = "TIMED_OUT"
				attempt.ErrorCode = "TIMEOUT"
			}
			if errors.Is(searchErr, context.Canceled) || errors.Is(operationErr, context.Canceled) {
				attempt.Status = "CANCELLED"
				attempt.ErrorCode = "CANCELLED"
			}
		} else if result.Provider != capability.Provider || result.Role != domain.RoleProjectCatalog || result.Attempt.Query != query.Query || result.Attempt.ExpansionID != query.ID || result.Attempt.Status != "SUCCEEDED" || !validInsightCatalogCandidates(result.Candidates, plan.Request.MaxPerQuery) {
			failed = true
			attempt.Status = "INTEGRITY_ERROR"
			attempt.ErrorCode = "INTEGRITY"
		} else {
			validSearches++
			attempt.Status = "SUCCEEDED"
			if len(result.Candidates) == 0 {
				attempt.Status = "EMPTY"
			}
			attempt.AdmittedCount = len(result.Candidates)
			report.Candidates = mergeInsightCatalogCandidates(report.Candidates, result.Candidates, query.ID, plan.Request.FamiliarApproaches)
		}
		attempt.Attempt.Status = attempt.Status
		attempt.Attempt.ErrorClass = attempt.ErrorCode
		attempt.Attempt.ResultCount = attempt.AdmittedCount
		report.Attempts = append(report.Attempts, attempt)
	}
	question := ""
	if plan.Request.QueryPolicy == "context-v2" {
		question = plan.Request.Question
		report.Limitations = append(report.Limitations, "Política contextual: correspondência literal em metadados precede diversidade; não é avaliação semântica. Revisões são resolvidas nas referências escolhidas.")
	}
	report.Candidates = selectInsightDiscoveryReferences(report.Candidates, plan.Request.MaxReferences, question)
	covered := map[string]bool{}
	for _, ref := range report.Candidates {
		if ref.Selected {
			for _, q := range ref.QueryIDs {
				covered[q] = true
			}
		}
	}
	for _, query := range plan.Queries {
		if !covered[query.ID] {
			report.Gaps = append(report.Gaps, "NO_SELECTED_REFERENCE_FOR_QUERY:"+query.ID)
		}
	}
	for _, reference := range report.Candidates {
		if !reference.Selected {
			continue
		}
		if ctx.Err() != nil {
			report.Readmes = append(report.Readmes, domain.InsightReadmeReceipt{ReferenceID: reference.ID, Source: domain.CommitRef{Owner: reference.Repository.Owner, Repo: reference.Repository.Repo, Commit: reference.Repository.CommitID}, Status: "NOT_STARTED", Scope: "ROOT_README_ONLY"})
			continue
		}
		operation, stop, timeoutErr := domain.ExecutionTimeout(ctx, plan.Request.AttemptTimeout)
		if timeoutErr != nil {
			return domain.InsightDiscoveryReport{}, timeoutErr
		}
		receipt, bundle := collectInsightReadme(operation, reference, reader)
		operationErr := operation.Err()
		stop()
		if operationErr != nil {
			receipt.Status = "CANCELLED"
			bundle = nil
			if errors.Is(operationErr, context.DeadlineExceeded) {
				receipt.Status = "TIMED_OUT"
			}
		}
		report.Readmes = append(report.Readmes, receipt)
		if bundle != nil {
			report.Bundles = append(report.Bundles, *bundle)
		} else {
			failed = true
		}
	}
	// Preserve completed observations even when the operational deadline has expired.
	report.Exploration, err = ExploreInsights(context.Background(), domain.InsightExplorationRequest{SchemaVersion: "1", AlgorithmIdeas: plan.Request.AlgorithmIdeas, Question: plan.Request.Question, FamiliarApproaches: plan.Request.FamiliarApproaches, MaxPassages: plan.Request.MaxPassages, Bundles: report.Bundles})
	if err != nil {
		return domain.InsightDiscoveryReport{}, err
	}
	report.Status = "REFERENCES_COLLECTED"
	if len(report.Candidates) == 0 {
		report.Status = "NO_CANDIDATES"
	}
	if failed {
		report.Status = "PARTIAL"
	}
	if validSearches == 0 {
		report.Status = "FAILED"
		err = domain.ErrInsightDiscoveryUnavailable
	}
	if ctx.Err() != nil {
		err = ctx.Err()
		report.Status = "CANCELLED"
		if errors.Is(err, context.DeadlineExceeded) {
			report.Status = "TIMED_OUT"
		}
	}
	report.ID = "discovery_" + insightHash(report)
	return report, err
}
