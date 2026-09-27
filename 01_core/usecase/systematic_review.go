/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/systematic_review.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

func reviewHash(value any) string {
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// FreezeReviewProtocol materializes all rules that must exist before candidate selection.
func FreezeReviewProtocol(draft domain.ReviewProtocolDraft) (domain.ReviewProtocol, error) {
	if draft.Intake.Status != domain.IntakeReady {
		return domain.ReviewProtocol{}, fmt.Errorf("intake deve estar READY antes de congelar o protocolo")
	}
	if strings.TrimSpace(draft.Intake.InputHash) == "" {
		return domain.ReviewProtocol{}, fmt.Errorf("intake.input_hash é obrigatório")
	}
	if strings.TrimSpace(draft.EvidenceRulesVersion) == "" || strings.TrimSpace(draft.EvidenceRulesHash) == "" {
		return domain.ReviewProtocol{}, fmt.Errorf("versão e hash das regras de evidência são obrigatórios")
	}
	if len(draft.Intake.Request.Idea.Capabilities) == 0 {
		return domain.ReviewProtocol{}, fmt.Errorf("ao menos uma capacidade é obrigatória")
	}
	modelHash := reviewHash(draft.ScreeningModel)
	policy := draft.UpdatePolicy
	if policy.Mode == "" {
		policy = domain.LivingUpdatePolicy{Mode: "LIVING", Trigger: "NEW_EVIDENCE", RequireSameProtocol: true, PreservePriorVersion: true}
	}
	protocol := domain.ReviewProtocol{
		SchemaVersion: domain.SystematicReviewSchemaVersion,
		IntakeHash:    draft.Intake.InputHash,
		Question:      draft.Intake.Request.Idea,
		Target:        draft.Intake.Request.Target,
		Search: domain.ReviewSearchProtocol{
			Strategy: draft.Intake.Request.Search.Strategy, Language: draft.Intake.Request.Search.Language,
			Locales: normalizedSorted(draft.Intake.Request.Search.Locales), MaxPerQuery: draft.Intake.Request.Search.MaxPerQuery,
			MaxCandidates: draft.Intake.Request.Search.MaxCandidates, Provider: draft.Intake.Execution.Provider,
		},
		Eligibility: reviewEligibility(draft.Intake.Request),
		Screening: domain.ReviewScreeningProtocol{
			ModelReference: draft.ScreeningModelReference, ModelHash: modelHash,
			AdvanceThreshold: draft.ScreeningModel.AdvanceThreshold, DropThreshold: draft.ScreeningModel.DropThreshold,
			AuditSampleRate: draft.ScreeningModel.AuditSampleRate, Calibration: draft.ScreeningModel.CalibrationStatus,
		},
		Evidence:     domain.ReviewEvidenceProtocol{RulesVersion: draft.EvidenceRulesVersion, RulesHash: draft.EvidenceRulesHash},
		Decision:     draft.Intake.Request.Policy,
		Execution:    draft.Intake.Execution,
		UpdatePolicy: policy,
	}
	protocol.ProtocolHash = protocolDigest(protocol)
	return protocol, nil
}

func reviewEligibility(request domain.RecommendationRequest) []domain.EligibilityCriterion {
	criteria := []domain.EligibilityCriterion{
		{ID: "ELIG-IDENTITY-001", Kind: "INCLUDE", Rule: "repositório e revisão imutável identificáveis", SourceRef: "candidate.commit_id"},
		{ID: "ELIG-ALIGNMENT-001", Kind: "INCLUDE", Rule: "alinhamento com o problema deve ser comprovado", SourceRef: "assessment.problem_alignment"},
		{ID: "ELIG-CAPABILITY-001", Kind: "INCLUDE", Rule: "ao menos uma capacidade requerida deve ter evidência admissível", SourceRef: "assessment.capabilities"},
		{ID: "ELIG-CONFLICT-001", Kind: "EXCLUDE", Rule: "conflito de evidência impede recomendação", SourceRef: "assessment"},
	}
	for _, license := range normalizedLicenses(request.Policy.AllowedLicenses) {
		criteria = append(criteria, domain.EligibilityCriterion{ID: "ELIG-LICENSE-" + strings.ToUpper(strings.ReplaceAll(license, " ", "-")), Kind: "INCLUDE", Rule: "licença permitida: " + license, SourceRef: "policy.allowed_licenses"})
	}
	if !request.Policy.AllowArchived {
		criteria = append(criteria, domain.EligibilityCriterion{ID: "ELIG-ARCHIVED-001", Kind: "EXCLUDE", Rule: "repositórios arquivados são excluídos", SourceRef: "policy.allow_archived"})
	}
	if request.Policy.RequireLocal {
		criteria = append(criteria, domain.EligibilityCriterion{ID: "ELIG-LOCAL-001", Kind: "INCLUDE", Rule: "operação local deve ser comprovada", SourceRef: "policy.require_local"})
	}
	if request.Target.Repository != nil {
		criteria = append(criteria, domain.EligibilityCriterion{ID: "ELIG-TARGET-001", Kind: "EXCLUDE", Rule: "o próprio alvo não é alternativa a si mesmo", SourceRef: "target.repository"})
	}
	sort.Slice(criteria, func(i, j int) bool { return criteria[i].ID < criteria[j].ID })
	return criteria
}

func protocolDigest(protocol domain.ReviewProtocol) string {
	copy := protocol
	copy.ProtocolHash = ""
	return reviewHash(copy)
}

func validateProtocol(protocol domain.ReviewProtocol) error {
	if protocol.SchemaVersion != domain.SystematicReviewSchemaVersion {
		return fmt.Errorf("schema de protocolo não suportado: %q", protocol.SchemaVersion)
	}
	if protocol.ProtocolHash == "" || protocol.ProtocolHash != protocolDigest(protocol) {
		return fmt.Errorf("protocol_hash não corresponde ao protocolo congelado")
	}
	return nil
}

// SynthesizeSystematicReview derives the ledger, flow and certainty from existing Bateia artifacts.
func SynthesizeSystematicReview(input domain.ReviewSynthesisInput) (domain.SystematicReviewPackage, error) {
	if err := validateProtocol(input.Protocol); err != nil {
		return domain.SystematicReviewPackage{}, err
	}
	if strings.TrimSpace(input.SearchedThrough) == "" {
		return domain.SystematicReviewPackage{}, fmt.Errorf("searched_through é obrigatório")
	}
	ledger := buildDispositionLedger(input)
	flow, err := deriveFunnelFlow(input, ledger)
	if err != nil {
		return domain.SystematicReviewPackage{}, err
	}
	result := domain.SystematicReviewPackage{
		SchemaVersion: domain.SystematicReviewSchemaVersion, Protocol: input.Protocol,
		SearchedThrough: input.SearchedThrough, Ledger: ledger, Flow: flow,
		Certainty: buildCertaintyReport(input), Decision: input.Decision, Alternatives: input.Alternatives,
	}
	result.PackageHash = packageDigest(result)
	return result, nil
}

func reviewRepositoryKey(owner, repo string) string {
	return strings.ToLower(strings.TrimSpace(owner)) + "/" + strings.ToLower(strings.TrimSpace(repo))
}

func repositoryStringKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if index := strings.Index(value, "@"); index >= 0 {
		value = value[:index]
	}
	return value
}

func buildDispositionLedger(input domain.ReviewSynthesisInput) domain.CandidateDispositionLedger {
	records := []domain.ReviewSourceRecord{}
	canonical := map[string]string{}
	candidateRecords := map[string][]string{}
	seenSnapshotCandidate := map[string]domain.RepositoryCandidate{}
	for _, discovered := range input.Discovery.Candidates {
		key := reviewRepositoryKey(discovered.Repository.Owner, discovered.Repository.Repo)
		seenSnapshotCandidate[key] = discovered.Repository
	}
	for expansionIndex, expansion := range input.Discovery.PerExpansion {
		for candidateIndex, candidate := range expansion.Provider.Candidates {
			key := reviewRepositoryKey(candidate.Owner, candidate.Repo)
			id := reviewHash([]any{expansionIndex, candidateIndex, expansion.Provider.Provider, expansion.Expansion.ID, expansion.Expansion.Query, key})[:20]
			disposition := domain.SourceRecordCanonical
			canonicalID := ""
			if previous := canonical[key]; previous != "" {
				disposition, canonicalID = domain.SourceRecordDuplicate, previous
			} else {
				canonical[key] = id
			}
			records = append(records, domain.ReviewSourceRecord{RecordID: id, CandidateID: key, Repository: candidate, Provider: expansion.Provider.Provider, ExpansionID: expansion.Expansion.ID, PlannedQuery: expansion.Expansion.Query, Disposition: disposition, CanonicalRecordID: canonicalID})
			candidateRecords[key] = append(candidateRecords[key], id)
			seenSnapshotCandidate[key] = candidate
		}
	}
	keys := make([]string, 0, len(seenSnapshotCandidate))
	for key := range seenSnapshotCandidate {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if len(candidateRecords[key]) > 0 {
			continue
		}
		candidate := seenSnapshotCandidate[key]
		id := reviewHash([]string{"snapshot", key})[:20]
		records = append(records, domain.ReviewSourceRecord{RecordID: id, CandidateID: key, Repository: candidate, Provider: "DISCOVERY_SNAPSHOT", Disposition: domain.SourceRecordCanonical})
		candidateRecords[key] = []string{id}
		canonical[key] = id
	}
	sort.Slice(records, func(i, j int) bool { return records[i].RecordID < records[j].RecordID })

	screening := map[string]domain.ScreeningDecision{}
	for _, decision := range input.Screening.Decisions {
		screening[repositoryStringKey(decision.Repository)] = decision
	}
	assessments := map[string]domain.CandidateAssessment{}
	for _, assessment := range input.Assessment.Assessments {
		assessments[reviewRepositoryKey(assessment.Repository.Owner, assessment.Repository.Repo)] = assessment
	}
	included := map[string]domain.SoftwareAlternative{}
	for _, alternative := range input.Alternatives.Alternatives {
		included[reviewRepositoryKey(alternative.Repository.Owner, alternative.Repository.Repo)] = alternative
	}
	excluded := map[string]domain.SoftwareAlternative{}
	for _, alternative := range input.Alternatives.Excluded {
		excluded[reviewRepositoryKey(alternative.Repository.Owner, alternative.Repository.Repo)] = alternative
	}

	candidates := make([]domain.CandidateDisposition, 0, len(keys))
	for _, key := range keys {
		metadata := seenSnapshotCandidate[key]
		ref := domain.CommitRef{Owner: strings.ToLower(strings.TrimSpace(metadata.Owner)), Repo: strings.ToLower(strings.TrimSpace(metadata.Repo)), Commit: metadata.CommitID}
		events := []domain.CandidateDispositionEvent{{Stage: "IDENTIFICATION", Outcome: "UNIQUE", SourceRef: canonical[key]}}
		if decision, ok := screening[key]; ok {
			events = append(events, domain.CandidateDispositionEvent{Stage: "SCREENING", Outcome: decision.Decision, SourceRef: decision.Repository})
		}
		if assessment, ok := assessments[key]; ok {
			outcome := "ASSESSED"
			if len(assessment.Failures) > 0 {
				outcome = "PARTIAL"
			}
			events = append(events, domain.CandidateDispositionEvent{Stage: "EVIDENCE", Outcome: outcome, SourceRef: assessment.Repository.Key()})
		}
		terminal, reason := terminalDisposition(key, input.Protocol, screening[key], assessments[key], included[key], excluded[key])
		events = append(events, domain.CandidateDispositionEvent{Stage: "ELIGIBILITY", Outcome: string(terminal), ReasonCode: reason})
		candidates = append(candidates, domain.CandidateDisposition{CandidateID: key, Repository: ref, SourceRecordIDs: append([]string{}, candidateRecords[key]...), Events: events, Terminal: terminal, ReasonCode: reason})
	}
	return domain.CandidateDispositionLedger{SchemaVersion: domain.SystematicReviewSchemaVersion, ProtocolHash: input.Protocol.ProtocolHash, SourceRecords: records, Candidates: candidates}
}

func terminalDisposition(key string, protocol domain.ReviewProtocol, screening domain.ScreeningDecision, assessment domain.CandidateAssessment, included, excluded domain.SoftwareAlternative) (domain.CandidateTerminalDecision, domain.CandidateDispositionReason) {
	if included.Repository.Owner != "" {
		switch included.Fit {
		case domain.AlternativeRecommended:
			return domain.CandidateIncluded, domain.ReasonIncludedRecommended
		case domain.AlternativePartial:
			return domain.CandidateIncluded, domain.ReasonIncludedPartial
		default:
			return domain.CandidateIncluded, domain.ReasonIncludedReference
		}
	}
	if excluded.Repository.Owner != "" {
		if protocol.Target.Repository != nil && key == reviewRepositoryKey(protocol.Target.Repository.Owner, protocol.Target.Repository.Repo) {
			return domain.CandidateExcluded, domain.ReasonTargetRepository
		}
		if excluded.Maintenance.Archived && !protocol.Decision.AllowArchived {
			return domain.CandidateExcluded, domain.ReasonArchivedBlocked
		}
		if protocol.Decision.RequireLocal && excluded.Operation.LocalStatus != domain.EvidenceProven {
			return domain.CandidateExcluded, domain.ReasonLocalRequired
		}
		if licenseBlocked(excluded.License, protocol.Decision.AllowedLicenses) {
			return domain.CandidateExcluded, domain.ReasonLicenseBlocked
		}
		if hasAssessmentConflict(assessment) {
			return domain.CandidateExcluded, domain.ReasonEvidenceConflict
		}
		switch excluded.Fit {
		case domain.AlternativeNotAlternative:
			return domain.CandidateExcluded, domain.ReasonNotAlternative
		case domain.AlternativeBlocked:
			return domain.CandidateExcluded, domain.ReasonPolicyBlocked
		default:
			return domain.CandidateInconclusive, domain.ReasonEvidenceInconclusive
		}
	}
	if strings.EqualFold(screening.Decision, "DROP") {
		return domain.CandidateExcluded, domain.ReasonScreeningDrop
	}
	if assessment.Repository.Owner == "" {
		return domain.CandidateInconclusive, domain.ReasonEvidenceNotAssessed
	}
	if len(assessment.Failures) > 0 {
		return domain.CandidateInconclusive, domain.ReasonEvidenceRetrieval
	}
	if hasAssessmentConflict(assessment) {
		return domain.CandidateExcluded, domain.ReasonEvidenceConflict
	}
	return domain.CandidateInconclusive, domain.ReasonEvidenceInconclusive
}

func licenseBlocked(license domain.LicenseAssessment, allowed []string) bool {
	if license.Status != domain.EvidenceProven {
		return true
	}
	for _, value := range allowed {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(license.Identifier)) {
			return false
		}
	}
	return len(allowed) > 0
}

func hasAssessmentConflict(assessment domain.CandidateAssessment) bool {
	if assessment.ProblemAlignment.Status == domain.EvidenceConflict || assessment.License.Status == domain.EvidenceConflict || assessment.Maintenance.Status == domain.EvidenceConflict || assessment.Operation.LocalStatus == domain.EvidenceConflict {
		return true
	}
	for _, capability := range assessment.Capabilities {
		if capability.Status == domain.EvidenceConflict {
			return true
		}
	}
	for _, constraint := range assessment.Constraints {
		if constraint.Evidence.Status == domain.EvidenceConflict {
			return true
		}
	}
	return false
}

func deriveFunnelFlow(input domain.ReviewSynthesisInput, ledger domain.CandidateDispositionLedger) (domain.FunnelFlowReport, error) {
	flow := domain.FunnelFlowReport{SchemaVersion: domain.SystematicReviewSchemaVersion, ProtocolHash: ledger.ProtocolHash, IdentifiedRecords: len(ledger.SourceRecords), UniqueCandidates: len(ledger.Candidates)}
	for _, record := range ledger.SourceRecords {
		if record.Disposition == domain.SourceRecordDuplicate {
			flow.DuplicateRecords++
		}
	}
	for _, expansion := range input.Discovery.PerExpansion {
		if count := len(expansion.Provider.Attempts); count > 0 {
			flow.ProviderRejected += expansion.Provider.Attempts[count-1].RejectedCandidates
		}
	}
	known := map[string]bool{}
	for _, decision := range input.Screening.Decisions {
		key := repositoryStringKey(decision.Repository)
		if known[key] {
			continue
		}
		known[key] = true
		flow.Screened++
		switch strings.ToUpper(decision.Decision) {
		case "ADVANCE":
			flow.ScreeningAdvance++
		case "AUDIT":
			flow.ScreeningAudit++
		case "DROP":
			flow.ScreeningDrop++
		}
	}
	flow.Unscreened = flow.UniqueCandidates - flow.Screened
	if flow.Unscreened < 0 {
		return domain.FunnelFlowReport{}, fmt.Errorf("screening contém mais candidatos únicos que o ledger")
	}
	assessed := map[string]bool{}
	failed := map[string]bool{}
	for _, assessment := range input.Assessment.Assessments {
		key := reviewRepositoryKey(assessment.Repository.Owner, assessment.Repository.Repo)
		assessed[key] = true
		if len(assessment.Failures) > 0 {
			failed[key] = true
		}
	}
	for _, bundle := range input.EvidenceBundles {
		key := reviewRepositoryKey(bundle.Repository.Owner, bundle.Repository.Repo)
		if len(bundle.Failures) > 0 || !bundle.Complete {
			failed[key] = true
		}
	}
	flow.Assessed, flow.EvidenceFailures = len(assessed), len(failed)
	reasons := map[domain.CandidateDispositionReason]int{}
	for _, candidate := range ledger.Candidates {
		switch candidate.Terminal {
		case domain.CandidateIncluded:
			flow.Included++
		case domain.CandidateExcluded:
			flow.Excluded++
		case domain.CandidateInconclusive:
			flow.Inconclusive++
		default:
			return domain.FunnelFlowReport{}, fmt.Errorf("candidato %s sem disposição terminal válida", candidate.CandidateID)
		}
		if candidate.ReasonCode == "" {
			return domain.FunnelFlowReport{}, fmt.Errorf("candidato %s sem reason_code terminal", candidate.CandidateID)
		}
		reasons[candidate.ReasonCode]++
	}
	for reason, count := range reasons {
		flow.TerminalReasons = append(flow.TerminalReasons, domain.ReasonCount{ReasonCode: reason, Count: count})
	}
	sort.Slice(flow.TerminalReasons, func(i, j int) bool { return flow.TerminalReasons[i].ReasonCode < flow.TerminalReasons[j].ReasonCode })
	canonical := flow.IdentifiedRecords - flow.DuplicateRecords
	if canonical != flow.UniqueCandidates {
		return domain.FunnelFlowReport{}, fmt.Errorf("fluxo não conserva identificação: %d canônicos para %d candidatos", canonical, flow.UniqueCandidates)
	}
	if flow.Included+flow.Excluded+flow.Inconclusive != flow.UniqueCandidates {
		return domain.FunnelFlowReport{}, fmt.Errorf("fluxo não conserva disposições terminais")
	}
	flow.ConservationChecked = true
	return flow, nil
}

func buildCertaintyReport(input domain.ReviewSynthesisInput) domain.ReviewCertaintyReport {
	bundles := map[string]domain.EvidenceBundle{}
	for _, bundle := range input.EvidenceBundles {
		bundles[reviewRepositoryKey(bundle.Repository.Owner, bundle.Repository.Repo)] = bundle
	}
	claims := []domain.ClaimCertainty{}
	for _, assessment := range input.Assessment.Assessments {
		key := reviewRepositoryKey(assessment.Repository.Owner, assessment.Repository.Repo)
		bundle, hasBundle := bundles[key]
		claims = append(claims, certaintyForEvidence(key, "problem_alignment", assessment.ProblemAlignment, bundle, hasBundle, input.Discovery))
		for _, capability := range assessment.Capabilities {
			claims = append(claims, certaintyForEvidence(key, "capability:"+normalizeIdeaText(capability.Capability), capability, bundle, hasBundle, input.Discovery))
		}
		for _, constraint := range assessment.Constraints {
			claims = append(claims, certaintyForEvidence(key, "constraint:"+normalizeIdeaText(constraint.Constraint), constraint.Evidence, bundle, hasBundle, input.Discovery))
		}
		claims = append(claims, certaintyForSources(key, "license", assessment.License.Status, assessment.License.Sources, bundle, hasBundle, input.Discovery))
		claims = append(claims, certaintyForSources(key, "maintenance", assessment.Maintenance.Status, assessment.Maintenance.Sources, bundle, hasBundle, input.Discovery))
		claims = append(claims, certaintyForSources(key, "local_operation", assessment.Operation.LocalStatus, assessment.Operation.Sources, bundle, hasBundle, input.Discovery))
	}
	sort.Slice(claims, func(i, j int) bool {
		if claims[i].CandidateID == claims[j].CandidateID {
			return claims[i].ClaimID < claims[j].ClaimID
		}
		return claims[i].CandidateID < claims[j].CandidateID
	})
	return domain.ReviewCertaintyReport{SchemaVersion: domain.SystematicReviewSchemaVersion, Claims: claims, Biases: assessReviewBiases(input)}
}

func certaintyForEvidence(candidate, claim string, evidence domain.CapabilityEvidence, bundle domain.EvidenceBundle, hasBundle bool, discovery domain.DiscoverySnapshot) domain.ClaimCertainty {
	return certaintyForSources(candidate, claim, evidence.Status, evidence.Sources, bundle, hasBundle, discovery)
}

func certaintyForSources(candidate, claim string, status domain.EvidenceStatus, sources []domain.EvidenceSource, bundle domain.EvidenceBundle, hasBundle bool, discovery domain.DiscoverySnapshot) domain.ClaimCertainty {
	concerns := []domain.ClaimConcern{}
	strong, documentation, pathOnly := false, false, false
	for _, source := range sources {
		if source.MatchClass == domain.EvidenceMatchPathOnly {
			pathOnly = true
		}
		if source.Kind == domain.EvidenceDocumentation {
			documentation = true
		}
		if source.MatchClass == domain.EvidenceMatchStrongContent || source.Kind == domain.EvidenceTest || source.Kind == domain.EvidenceTestImplementation || source.Kind == domain.EvidenceLicense || source.Kind == domain.EvidenceRepository {
			strong = true
		}
	}
	concerns = append(concerns, concern("source_authority", func() domain.CertaintyConcern {
		if len(sources) == 0 {
			return domain.ConcernUnknown
		}
		if pathOnly && !strong {
			return domain.ConcernSerious
		}
		if documentation && !strong {
			return domain.ConcernSome
		}
		return domain.ConcernNone
	}(), "força e autoridade dos artefatos citados"))
	consistency := domain.ConcernNone
	if status == domain.EvidenceConflict {
		consistency = domain.ConcernSerious
	}
	concerns = append(concerns, concern("consistency", consistency, "concordância entre fontes observadas"))
	directness := domain.ConcernNone
	if !strong && documentation {
		directness = domain.ConcernSerious
	} else if !strong {
		directness = domain.ConcernSome
	}
	concerns = append(concerns, concern("directness", directness, "proximidade entre fonte e claim"))
	completeness := domain.ConcernUnknown
	if hasBundle {
		completeness = domain.ConcernNone
		if !bundle.Complete || !bundle.TreeComplete || len(bundle.Failures) > 0 {
			completeness = domain.ConcernSerious
		} else {
			for _, coverage := range bundle.RequirementCoverage {
				if !coverage.SearchComplete || len(coverage.SkippedPaths) > 0 {
					completeness = domain.ConcernSome
					break
				}
			}
		}
	}
	concerns = append(concerns, concern("collection_completeness", completeness, "cobertura e falhas da coleta"))
	availability := domain.ConcernNone
	if len(discovery.Failures) > 0 {
		availability = domain.ConcernSerious
	} else if len(discovery.PerExpansion) < 2 {
		availability = domain.ConcernSome
	}
	concerns = append(concerns, concern("source_availability", availability, "amplitude das fontes de descoberta"))
	return domain.ClaimCertainty{CandidateID: candidate, ClaimID: claim, ObservedStatus: status, Level: certaintyLevel(status, concerns), Concerns: concerns, SourceCount: len(sources)}
}

func concern(dimension string, value domain.CertaintyConcern, reason string) domain.ClaimConcern {
	return domain.ClaimConcern{Dimension: dimension, Concern: value, Reason: reason}
}

func certaintyLevel(status domain.EvidenceStatus, concerns []domain.ClaimConcern) domain.CertaintyLevel {
	if status == domain.EvidenceInconclusive || len(concerns) == 0 {
		return domain.CertaintyUnknown
	}
	worst := domain.ConcernNone
	rank := map[domain.CertaintyConcern]int{domain.ConcernNone: 0, domain.ConcernSome: 1, domain.ConcernSerious: 2, domain.ConcernUnknown: 3}
	for _, item := range concerns {
		if rank[item.Concern] > rank[worst] {
			worst = item.Concern
		}
	}
	if worst == domain.ConcernUnknown {
		return domain.CertaintyUnknown
	}
	if status == domain.EvidenceConflict || worst == domain.ConcernSerious {
		return domain.CertaintyFragile
	}
	if status == domain.EvidenceProven && worst == domain.ConcernNone {
		return domain.CertaintyRobust
	}
	return domain.CertaintySupported
}

func assessReviewBiases(input domain.ReviewSynthesisInput) []domain.BiasAssessment {
	biases := []domain.BiasAssessment{}
	add := func(name string, risk domain.BiasRisk, reason, mitigation string) {
		biases = append(biases, domain.BiasAssessment{Bias: name, Risk: risk, Reason: reason, Mitigation: mitigation})
	}
	coverageRisk := domain.BiasLow
	if len(input.Discovery.Failures) > 0 || len(input.Discovery.PerExpansion) == 0 {
		coverageRisk = domain.BiasPresent
	}
	add("SOURCE_COVERAGE", coverageRisk, "falhas e amplitude dos provedores observados", "reexecutar fontes ausentes e comparar o delta")
	languageRisk := domain.BiasLow
	if len(input.Discovery.SearchLocales) <= 1 {
		languageRisk = domain.BiasPresent
	}
	add("LANGUAGE", languageRisk, "quantidade de locales pesquisados", "adicionar expansões e fontes em outros idiomas")
	add("POPULARITY", domain.BiasPresent, "metadados de popularidade podem influenciar ordenação", "não usar estrelas como prova de capacidade")
	automationRisk := domain.BiasPresent
	if input.Screening.Strategy == "" {
		automationRisk = domain.BiasUnknown
	}
	add("AUTOMATION", automationRisk, "screening automatizado antecede revisão humana", "auditar drops e manter active learning sem poder de exclusão")
	evidenceRisk := domain.BiasLow
	if len(input.EvidenceBundles) < len(input.Assessment.Assessments) {
		evidenceRisk = domain.BiasPresent
	}
	add("EVIDENCE_AVAILABILITY", evidenceRisk, "nem todo assessment possui bundle observável", "marcar inconclusão e tentar recuperação")
	add("SURVIVORSHIP", domain.BiasPresent, "apenas repositórios encontráveis entram no funil", "registrar falhas e termos sem resultados")
	add("RECENCY", domain.BiasUnknown, "o pacote não infere atualidade sem searched_through externo", "versionar searched_through e atualizar incrementalmente")
	dupRisk := domain.BiasLow
	if len(input.Discovery.Candidates) == 0 {
		dupRisk = domain.BiasUnknown
	}
	add("DUPLICATION", dupRisk, "deduplicação usa identidade normalizada owner/repo", "reter todos os source records no ledger")
	sort.Slice(biases, func(i, j int) bool { return biases[i].Bias < biases[j].Bias })
	return biases
}

func packageDigest(pack domain.SystematicReviewPackage) string {
	copy := pack
	copy.PackageHash = ""
	return reviewHash(copy)
}

// UpdateLivingReview appends a version to the matching protocol lineage.
func UpdateLivingReview(input domain.LivingReviewUpdateInput) (domain.LivingReviewIndex, error) {
	current := input.CurrentPackage
	if err := validateProtocol(current.Protocol); err != nil {
		return domain.LivingReviewIndex{}, err
	}
	if current.PackageHash == "" || current.PackageHash != packageDigest(current) {
		return domain.LivingReviewIndex{}, fmt.Errorf("current_package possui hash inválido")
	}
	index := input.Index
	if index.SchemaVersion == "" {
		index = domain.LivingReviewIndex{SchemaVersion: domain.SystematicReviewSchemaVersion, Lineages: []domain.LivingReviewLineage{}}
	}
	if index.SchemaVersion != domain.SystematicReviewSchemaVersion {
		return domain.LivingReviewIndex{}, fmt.Errorf("schema de índice vivo não suportado")
	}
	lineageIndex := -1
	for i := range index.Lineages {
		if index.Lineages[i].ProtocolHash == current.Protocol.ProtocolHash {
			lineageIndex = i
			break
		}
	}
	if lineageIndex < 0 {
		index.Lineages = append(index.Lineages, domain.LivingReviewLineage{LineageID: current.Protocol.ProtocolHash[:20], ProtocolHash: current.Protocol.ProtocolHash, Versions: []domain.LivingReviewVersion{}})
		lineageIndex = len(index.Lineages) - 1
	}
	lineage := &index.Lineages[lineageIndex]
	for _, version := range lineage.Versions {
		if version.PackageHash == current.PackageHash {
			return index, nil
		}
	}
	if input.PreviousPackage != nil && input.PreviousPackage.Protocol.ProtocolHash != current.Protocol.ProtocolHash {
		return domain.LivingReviewIndex{}, fmt.Errorf("previous_package pertence a outro protocolo; nova linhagem não aceita delta cruzado")
	}
	delta := livingDelta(input.PreviousPackage, current)
	lineage.Versions = append(lineage.Versions, domain.LivingReviewVersion{Version: len(lineage.Versions) + 1, PackageHash: current.PackageHash, SearchedThrough: current.SearchedThrough, Delta: delta})
	sort.Slice(index.Lineages, func(i, j int) bool { return index.Lineages[i].LineageID < index.Lineages[j].LineageID })
	return index, nil
}

func livingDelta(previous *domain.SystematicReviewPackage, current domain.SystematicReviewPackage) domain.LivingReviewDelta {
	delta := domain.LivingReviewDelta{ToPackageHash: current.PackageHash, AddedCandidates: []string{}, RemovedCandidates: []string{}, DispositionChanges: []domain.LivingDispositionChange{}, CertaintyChanges: []domain.LivingCertaintyChange{}}
	currentCandidates := dispositionMap(current.Ledger)
	if previous == nil {
		for key := range currentCandidates {
			delta.AddedCandidates = append(delta.AddedCandidates, key)
		}
		sort.Strings(delta.AddedCandidates)
		return delta
	}
	delta.FromPackageHash = previous.PackageHash
	priorCandidates := dispositionMap(previous.Ledger)
	for key, value := range currentCandidates {
		prior, ok := priorCandidates[key]
		if !ok {
			delta.AddedCandidates = append(delta.AddedCandidates, key)
		} else if prior.Terminal != value.Terminal || prior.ReasonCode != value.ReasonCode {
			delta.DispositionChanges = append(delta.DispositionChanges, domain.LivingDispositionChange{CandidateID: key, Before: prior.Terminal, After: value.Terminal, BeforeReason: prior.ReasonCode, AfterReason: value.ReasonCode})
		}
	}
	for key := range priorCandidates {
		if _, ok := currentCandidates[key]; !ok {
			delta.RemovedCandidates = append(delta.RemovedCandidates, key)
		}
	}
	priorCertainty := certaintyMap(previous.Certainty)
	for key, value := range certaintyMap(current.Certainty) {
		if prior, ok := priorCertainty[key]; ok && prior.Level != value.Level {
			delta.CertaintyChanges = append(delta.CertaintyChanges, domain.LivingCertaintyChange{CandidateID: value.CandidateID, ClaimID: value.ClaimID, Before: prior.Level, After: value.Level})
		}
	}
	delta.DecisionChanged = previous.Decision.Decision != current.Decision.Decision || reviewHash(previous.Decision.Primary) != reviewHash(current.Decision.Primary)
	sort.Strings(delta.AddedCandidates)
	sort.Strings(delta.RemovedCandidates)
	sort.Slice(delta.DispositionChanges, func(i, j int) bool {
		return delta.DispositionChanges[i].CandidateID < delta.DispositionChanges[j].CandidateID
	})
	sort.Slice(delta.CertaintyChanges, func(i, j int) bool {
		left, right := delta.CertaintyChanges[i], delta.CertaintyChanges[j]
		return left.CandidateID+"\x00"+left.ClaimID < right.CandidateID+"\x00"+right.ClaimID
	})
	return delta
}

func dispositionMap(ledger domain.CandidateDispositionLedger) map[string]domain.CandidateDisposition {
	result := map[string]domain.CandidateDisposition{}
	for _, item := range ledger.Candidates {
		result[item.CandidateID] = item
	}
	return result
}

func certaintyMap(report domain.ReviewCertaintyReport) map[string]domain.ClaimCertainty {
	result := map[string]domain.ClaimCertainty{}
	for _, item := range report.Claims {
		result[item.CandidateID+"\x00"+item.ClaimID] = item
	}
	return result
}

// PrioritizeReviewCandidates ranks unreviewed items and has no authority to exclude them.
func PrioritizeReviewCandidates(request domain.ActiveLearningRequest) (domain.ActiveLearningPlan, error) {
	if request.BatchSize <= 0 {
		return domain.ActiveLearningPlan{}, fmt.Errorf("batch_size deve ser positivo")
	}
	reviewed := map[string]domain.RelevanceLabel{}
	decisionByKey := map[string]domain.ScreeningDecision{}
	for _, decision := range request.Screening.Decisions {
		decisionByKey[repositoryStringKey(decision.Repository)] = decision
	}
	for _, judgment := range request.Reviewed {
		key := repositoryStringKey(judgment.Repository)
		if _, ok := decisionByKey[key]; !ok {
			return domain.ActiveLearningPlan{}, fmt.Errorf("julgamento refere candidato ausente: %s", judgment.Repository)
		}
		if judgment.Label != domain.RelevanceRelevant && judgment.Label != domain.RelevanceIrrelevant && judgment.Label != domain.RelevanceUncertain {
			return domain.ActiveLearningPlan{}, fmt.Errorf("rótulo inválido para %s", judgment.Repository)
		}
		reviewed[key] = judgment.Label
	}
	positive, negative := [][]float64{}, [][]float64{}
	knownExpansions := map[string]bool{}
	for key, label := range reviewed {
		decision := decisionByKey[key]
		if label == domain.RelevanceRelevant {
			positive = append(positive, screeningVector(decision))
		} else if label == domain.RelevanceIrrelevant {
			negative = append(negative, screeningVector(decision))
		}
		for _, origin := range decision.Origins {
			knownExpansions[origin.ExpansionID] = true
		}
	}
	priorities := []domain.ActiveLearningPriority{}
	for _, decision := range request.Screening.Decisions {
		key := repositoryStringKey(decision.Repository)
		if _, ok := reviewed[key]; ok {
			continue
		}
		base := clamp(decision.RankingScore)
		if decision.Probability != nil {
			base = clamp(*decision.Probability)
		}
		vector := screeningVector(decision)
		estimated := estimateRelevance(base, vector, positive, negative)
		uncertainty := 1 - math.Abs(2*estimated-1)
		novel := 0.0
		for _, origin := range decision.Origins {
			if !knownExpansions[origin.ExpansionID] {
				novel = 1
				break
			}
		}
		priority := clamp(0.60*estimated + 0.30*uncertainty + 0.10*novel)
		priorities = append(priorities, domain.ActiveLearningPriority{Repository: decision.Repository, Priority: priority, EstimatedRelevance: estimated, Uncertainty: uncertainty, Novelty: novel, Reason: "relevância estimada + incerteza + novidade de origem"})
	}
	sort.Slice(priorities, func(i, j int) bool {
		if priorities[i].Priority == priorities[j].Priority {
			return priorities[i].Repository < priorities[j].Repository
		}
		return priorities[i].Priority > priorities[j].Priority
	})
	if len(priorities) > request.BatchSize {
		priorities = priorities[:request.BatchSize]
	}
	return domain.ActiveLearningPlan{SchemaVersion: domain.SystematicReviewSchemaVersion, Authority: "PRIORITIZATION_ONLY", CanExclude: false, Priorities: priorities}, nil
}

func screeningVector(decision domain.ScreeningDecision) []float64 {
	s := decision.Signals
	fallback := 1 / float64(1+maxInt(0, s.FallbackDepth))
	return []float64{clamp(decision.RankingScore), clamp(s.ExpansionPrior), clamp(s.QueryRetention), fallback, clamp(s.MatchSpecificity), clamp(s.DescriptionAlignment), clamp(float64(s.OriginConvergence) / 3), clamp(s.CapabilityCoverage), clamp(1 - s.NonGoalConflict), clamp(s.EvidenceCompleteness)}
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func estimateRelevance(base float64, vector []float64, positive, negative [][]float64) float64 {
	if len(positive) == 0 && len(negative) == 0 {
		return base
	}
	positiveDistance, negativeDistance := nearestDistance(vector, positive), nearestDistance(vector, negative)
	switch {
	case len(positive) > 0 && len(negative) > 0:
		learned := (negativeDistance + 1e-9) / (positiveDistance + negativeDistance + 2e-9)
		return clamp(0.35*base + 0.65*learned)
	case len(positive) > 0:
		return clamp(0.70*base + 0.30/(1+positiveDistance))
	default:
		return clamp(0.70*base + 0.30*(negativeDistance/(1+negativeDistance)))
	}
}

func nearestDistance(vector []float64, examples [][]float64) float64 {
	if len(examples) == 0 {
		return 0
	}
	best := math.Inf(1)
	for _, example := range examples {
		total := 0.0
		for index := range vector {
			delta := vector[index] - example[index]
			total += delta * delta
		}
		best = math.Min(best, math.Sqrt(total/float64(len(vector))))
	}
	return best
}

// EvaluateActiveLearning reveals gold labels only after each item is selected.
func EvaluateActiveLearning(corpus domain.ActiveLearningCorpus) (domain.ActiveLearningEvaluation, error) {
	if corpus.SchemaVersion != domain.SystematicReviewSchemaVersion {
		return domain.ActiveLearningEvaluation{}, fmt.Errorf("schema de corpus não suportado: %q", corpus.SchemaVersion)
	}
	result := domain.ActiveLearningEvaluation{SchemaVersion: domain.SystematicReviewSchemaVersion, Cases: []domain.ActiveLearningCaseMetrics{}, StoppingAuthority: false}
	for _, testCase := range corpus.Cases {
		metrics, err := evaluateActiveLearningCase(testCase)
		if err != nil {
			return domain.ActiveLearningEvaluation{}, fmt.Errorf("caso %s: %w", testCase.ID, err)
		}
		result.Cases = append(result.Cases, metrics)
		result.AverageWorkSaved95 += metrics.WorkSavedAt95
		result.AverageRecallHalf += metrics.RecallAtHalf
	}
	if len(result.Cases) > 0 {
		result.AverageWorkSaved95 /= float64(len(result.Cases))
		result.AverageRecallHalf /= float64(len(result.Cases))
	}
	return result, nil
}

func evaluateActiveLearningCase(testCase domain.ActiveLearningCase) (domain.ActiveLearningCaseMetrics, error) {
	gold := map[string]domain.RelevanceLabel{}
	for _, judgment := range testCase.Gold {
		key := repositoryStringKey(judgment.Repository)
		if _, exists := gold[key]; exists {
			return domain.ActiveLearningCaseMetrics{}, fmt.Errorf("gold duplicado para %s", key)
		}
		gold[key] = judgment.Label
	}
	if len(gold) != len(testCase.Screening.Decisions) {
		return domain.ActiveLearningCaseMetrics{}, fmt.Errorf("gold deve cobrir exatamente todos os candidatos")
	}
	relevant := 0
	for _, label := range gold {
		if label == domain.RelevanceRelevant {
			relevant++
		}
	}
	reviewed := append([]domain.CandidateJudgment{}, testCase.Seeds...)
	found := 0
	for _, seed := range reviewed {
		label, ok := gold[repositoryStringKey(seed.Repository)]
		if !ok || label != seed.Label {
			return domain.ActiveLearningCaseMetrics{}, fmt.Errorf("seed não corresponde ao gold: %s", seed.Repository)
		}
		if label == domain.RelevanceRelevant {
			found++
		}
	}
	metrics := domain.ActiveLearningCaseMetrics{CaseID: testCase.ID, Candidates: len(gold), Relevant: relevant}
	halfBudget := int(math.Ceil(float64(len(gold)) / 2))
	for len(reviewed) < len(gold) {
		plan, err := PrioritizeReviewCandidates(domain.ActiveLearningRequest{Screening: testCase.Screening, Reviewed: reviewed, BatchSize: 1})
		if err != nil {
			return metrics, err
		}
		if len(plan.Priorities) != 1 {
			return metrics, fmt.Errorf("priorizador não retornou item não revisado")
		}
		selected := plan.Priorities[0].Repository
		label, ok := gold[repositoryStringKey(selected)]
		if !ok {
			return metrics, fmt.Errorf("candidato selecionado sem gold: %s", selected)
		}
		reviewed = append(reviewed, domain.CandidateJudgment{Repository: selected, Label: label, Reason: "revealed_after_selection"})
		if label == domain.RelevanceRelevant {
			found++
		}
		if len(reviewed) == halfBudget {
			metrics.RecallAtHalf = recall(found, relevant)
		}
		if metrics.ScreenedFor95Recall == 0 && recall(found, relevant) >= 0.95 {
			metrics.ScreenedFor95Recall = len(reviewed)
		}
	}
	if halfBudget <= len(testCase.Seeds) {
		seedFound := 0
		for _, seed := range testCase.Seeds[:minInt(halfBudget, len(testCase.Seeds))] {
			if seed.Label == domain.RelevanceRelevant {
				seedFound++
			}
		}
		metrics.RecallAtHalf = recall(seedFound, relevant)
	}
	if relevant == 0 {
		metrics.ScreenedFor95Recall = 0
		metrics.WorkSavedAt95 = 0
	} else {
		metrics.WorkSavedAt95 = float64(len(gold)-metrics.ScreenedFor95Recall) / float64(len(gold))
	}
	return metrics, nil
}

func recall(found, relevant int) float64 {
	if relevant == 0 {
		return 1
	}
	return float64(found) / float64(relevant)
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
