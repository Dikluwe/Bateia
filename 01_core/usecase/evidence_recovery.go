/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_recovery.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func DefaultRecoveryPolicy() domain.RecoveryPolicy {
	return domain.RecoveryPolicy{SchemaVersion: "1", MaxRounds: 2, MaxCandidates: 3, AdditionalDocuments: 4, AdditionalCodeFiles: 8, AdditionalTests: 6, AdditionalBytes: 2 * 1024 * 1024, AdditionalTreeDepth: 2, RetryProviderFailure: true}
}

func DiagnoseInconclusive(runID, decisionHash string, decision domain.ReuseDecisionReport, assessment domain.AssessmentReport, bundles []domain.EvidenceBundle, catalog domain.CapabilityRuleCatalog) domain.InconclusiveDiagnosis {
	d := domain.InconclusiveDiagnosis{SchemaVersion: "1", RunID: runID, DecisionHash: decisionHash, Findings: []domain.InconclusiveFinding{}}
	if decision.Decision != domain.DecisionUnknown {
		return d
	}
	byRepo := map[string]domain.EvidenceBundle{}
	for _, b := range bundles {
		byRepo[b.Repository.Key()] = b
	}
	if len(assessment.Assessments) == 0 {
		d.Findings = append(d.Findings, newFinding(domain.CauseInsufficientCorpus, "", "", "DEC-INCONCLUSIVE-002", nil, nil, []string{"candidate assessment"}, false, "corpus não contém assessments"))
	}
	globalCatalogFindings := map[string]bool{}
	for _, a := range assessment.Assessments {
		key := a.Repository.Key()
		candidateKey := a.Repository.Key()
		b := byRepo[key]
		if a.Repository.Commit == "" || strings.EqualFold(a.Repository.Commit, "HEAD") {
			d.Findings = append(d.Findings, newFinding(domain.CauseUnresolvedRevision, candidateKey, "", "GATE-REVISION-001", nil, a.Failures, []string{"immutable commit"}, true, "revisão do candidato não foi resolvida"))
		}
		providerFailure := hasFailure(a.Failures) || hasFailure(b.Failures)
		budget := bundleHasStatus(b, domain.TargetSkippedBudget)
		for _, e := range a.Capabilities {
			if e.Status != domain.EvidenceInconclusive && !catalogLimitedEvidence(e) {
				continue
			}
			cause, rule, missing, recoverable, message := classifyCapabilityFinding(e, catalog, providerFailure, budget)
			findingCandidate := candidateKey
			findingEvidence := e.Sources
			findingFailures := append(a.Failures, b.Failures...)
			if catalogLimitedEvidence(e) {
				findingCandidate = ""
				findingEvidence = nil
				findingFailures = nil
				key := e.ReasonCode + "|" + normalizeEvidenceText(e.Capability)
				if globalCatalogFindings[key] {
					continue
				}
				globalCatalogFindings[key] = true
			}
			d.Findings = append(d.Findings, newFinding(cause, findingCandidate, e.Capability, rule, findingEvidence, findingFailures, missing, recoverable, message))
		}
		if a.License.Status == domain.EvidenceInconclusive {
			d.Findings = append(d.Findings, newFinding(domain.CauseUnknownLicense, candidateKey, "", "GATE-LICENSE-001", a.License.Sources, a.Failures, []string{"license artifact", "repository license metadata"}, true, "licença permanece desconhecida"))
		}
		if a.Operation.LocalStatus == domain.EvidenceInconclusive {
			d.Findings = append(d.Findings, newFinding(domain.CauseUnknownOperation, candidateKey, "", "GATE-OPERATION-001", a.Operation.Sources, a.Failures, []string{"runtime manifest", "local operation evidence"}, true, "modo de operação permanece desconhecido"))
		}
		for _, c := range a.Constraints {
			if c.Evidence.Status == domain.EvidenceConflict {
				d.Findings = append(d.Findings, newFinding(domain.CauseConflictingEvidence, candidateKey, c.Constraint, "GATE-CONSTRAINT-001", c.Evidence.Sources, a.Failures, []string{}, false, "conflito eliminatório exige revisão humana"))
			}
		}
	}
	sort.Slice(d.Findings, func(i, j int) bool { return findingSortKey(d.Findings[i]) < findingSortKey(d.Findings[j]) })
	for i := range d.Findings {
		d.Findings[i].ID = recoveryHash(struct{ Run, Cause, Candidate, Capability, Rule string }{runID, string(d.Findings[i].Cause), d.Findings[i].Candidate, d.Findings[i].Capability, d.Findings[i].RuleID})[:16]
		if d.Findings[i].Recoverable {
			d.Recoverable++
		} else {
			d.Terminal++
		}
	}
	return d
}

func classifyCapabilityFinding(e domain.CapabilityEvidence, catalog domain.CapabilityRuleCatalog, providerFailure, budget bool) (domain.InconclusiveCause, string, []string, bool, string) {
	switch e.ReasonCode {
	case "CATALOG_GAP":
		return domain.CauseRuleCoverageGap, "catalog-gap:" + normalizeEvidenceText(e.Capability), []string{"declared evidence rule"}, false, "o catálogo não possui regra declarada segura; isso não descreve o candidato"
	case "HYPOTHESIS_ONLY":
		ruleID := "hypothesis:" + normalizeEvidenceText(e.Capability)
		if rule := findRule(e.Capability, catalog); rule != nil {
			ruleID = rule.ID
		}
		return domain.CauseRuleCoverageGap, ruleID, []string{"declared rule validation"}, false, "a regra disponível é apenas hipótese sem autoridade; isso não descreve o candidato"
	}
	rule := findRule(e.Capability, catalog)
	if rule == nil {
		return domain.CauseRuleCoverageGap, "unmapped:" + normalizeEvidenceText(e.Capability), []string{"versioned evidence rule"}, false, "nenhuma regra versionada cobre a capacidade"
	}
	if budget {
		return domain.CauseBudgetExhausted, rule.ID, []string{"uncollected planned target"}, true, "budget encerrou coleta relevante"
	}
	if providerFailure {
		return domain.CauseProviderFailure, rule.ID, []string{"successful provider response"}, true, "provedor falhou durante coleta"
	}
	hasDoc, hasImpl, hasTest := false, false, false
	for _, s := range e.Sources {
		switch s.Kind {
		case domain.EvidenceDocumentation:
			hasDoc = true
		case domain.EvidenceFile, domain.EvidenceSymbol, domain.EvidenceExample:
			hasImpl = true
		case domain.EvidenceTest, domain.EvidenceTestImplementation:
			hasTest = true
		}
	}
	if hasImpl && !hasTest {
		return domain.CauseMissingTest, rule.ID, []string{"related test"}, true, "implementação observada sem teste independente"
	}
	if hasDoc && !hasImpl {
		return domain.CauseMissingImplementation, rule.ID, []string{"implementation symbol", "implementation file"}, true, "documentação observada sem implementação"
	}
	return domain.CauseIncompleteCollection, rule.ID, []string{"independent implementation and test evidence"}, true, "coleta não atingiu evidência independente suficiente"
}
func newFinding(c domain.InconclusiveCause, candidate, capability, rule string, e []domain.EvidenceSource, f []domain.ExecutionFailure, missing []string, recoverable bool, message string) domain.InconclusiveFinding {
	return domain.InconclusiveFinding{Cause: c, Candidate: candidate, Capability: capability, RuleID: rule, Evidence: append([]domain.EvidenceSource(nil), e...), Failures: append([]domain.ExecutionFailure(nil), f...), Missing: append([]string(nil), missing...), Recoverable: recoverable, Message: message}
}
func hasFailure(f []domain.ExecutionFailure) bool { return len(f) > 0 }
func bundleHasStatus(b domain.EvidenceBundle, s domain.CollectionTargetStatus) bool {
	for _, t := range b.Targets {
		if t.Status == s {
			return true
		}
	}
	return false
}
func findingSortKey(f domain.InconclusiveFinding) string {
	return f.Candidate + "|" + f.Capability + "|" + string(f.Cause) + "|" + f.RuleID
}

func PlanEvidenceRecovery(d domain.InconclusiveDiagnosis, policy domain.RecoveryPolicy, round int, previous []domain.RecoveryTarget) (domain.EvidenceRecoveryPlan, error) {
	if policy.SchemaVersion == "" || policy.MaxRounds < 1 || policy.MaxCandidates < 1 || policy.AdditionalBytes < 1 {
		return domain.EvidenceRecoveryPlan{}, fmt.Errorf("política de recuperação inválida")
	}
	if round < 1 || round > policy.MaxRounds {
		return domain.EvidenceRecoveryPlan{}, fmt.Errorf("rodada %d excede max_rounds", round)
	}
	plan := domain.EvidenceRecoveryPlan{SchemaVersion: "1", Round: round, DiagnosisHash: canonicalDiagnosisHash(d), Policy: policy, PolicyHash: recoveryHash(policy), Targets: []domain.RecoveryTarget{}, Skipped: []domain.InconclusiveFinding{}}
	seen := map[string]bool{}
	for _, t := range previous {
		seen[targetIdentity(t)] = true
	}
	candidates := map[string]bool{}
	findings := append([]domain.InconclusiveFinding(nil), d.Findings...)
	sort.Slice(findings, func(i, j int) bool { return findingSortKey(findings[i]) < findingSortKey(findings[j]) })
	for _, f := range findings {
		if !f.Recoverable {
			plan.Skipped = append(plan.Skipped, f)
			continue
		}
		ref, err := parseRepositoryReference(f.Candidate)
		if err != nil {
			plan.Skipped = append(plan.Skipped, f)
			continue
		}
		if !candidates[f.Candidate] && len(candidates) >= policy.MaxCandidates {
			plan.Skipped = append(plan.Skipped, f)
			continue
		}
		options := targetsForFinding(f, ref, policy)
		if len(options) == 0 {
			plan.Skipped = append(plan.Skipped, f)
			continue
		}
		var target domain.RecoveryTarget
		found := false
		for _, option := range options {
			if !seen[targetIdentity(option)] {
				target = option
				found = true
				break
			}
		}
		if !found {
			plan.Skipped = append(plan.Skipped, f)
			continue
		}
		seen[targetIdentity(target)] = true
		candidates[f.Candidate] = true
		plan.Targets = append(plan.Targets, target)
	}
	sort.Slice(plan.Targets, func(i, j int) bool { return targetIdentity(plan.Targets[i]) < targetIdentity(plan.Targets[j]) })
	sort.Slice(plan.Skipped, func(i, j int) bool { return findingSortKey(plan.Skipped[i]) < findingSortKey(plan.Skipped[j]) })
	return plan, nil
}

func targetsForFinding(f domain.InconclusiveFinding, ref domain.CommitRef, p domain.RecoveryPolicy) []domain.RecoveryTarget {
	t := domain.RecoveryTarget{FindingID: f.ID, Candidate: ref, Capability: f.Capability, RuleID: f.RuleID, Expected: []domain.EvidenceStatus{domain.EvidenceProven, domain.EvidencePartial, domain.EvidenceAbsent, domain.EvidenceConflict}, Reason: f.Message}
	terms := strings.Fields(normalizeEvidenceText(f.Capability))
	for _, x := range terms {
		if len(x) >= 4 {
			t.Patterns = append(t.Patterns, x)
		}
	}
	switch f.Cause {
	case domain.CauseUnresolvedRevision:
		t.Kind = domain.EvidenceRepository
		t.Path = "@repository-metadata"
		return []domain.RecoveryTarget{t}
	case domain.CauseUnknownLicense:
		t.Kind = domain.EvidenceLicense
		out := []domain.RecoveryTarget{}
		for _, path := range []string{"LICENSE", "COPYING", "NOTICE", "@repository-metadata"} {
			copy := t
			copy.Path = path
			if path == "@repository-metadata" {
				copy.Kind = domain.EvidenceRepository
			}
			out = append(out, copy)
		}
		return out
	case domain.CauseUnknownOperation:
		t.Kind = domain.EvidenceManifest
		out := []domain.RecoveryTarget{}
		for _, path := range []string{"go.mod", "package.json", "Cargo.toml", "pyproject.toml", "README.md"} {
			copy := t
			copy.Path = path
			if path == "README.md" {
				copy.Kind = domain.EvidenceDocumentation
			}
			out = append(out, copy)
		}
		return out
	case domain.CauseMissingTest:
		t.Kind = domain.EvidenceTest
	case domain.CauseMissingImplementation, domain.CauseIncompleteCollection, domain.CauseBudgetExhausted:
		t.Kind = domain.EvidenceFile
	case domain.CauseProviderFailure:
		if !p.RetryProviderFailure {
			return nil
		}
		t.Kind = domain.EvidenceRepository
		t.Path = "@repository-metadata"
	default:
		return nil
	}
	return []domain.RecoveryTarget{t}
}
func targetIdentity(t domain.RecoveryTarget) string {
	return strings.ToLower(t.Candidate.Key() + "|" + string(t.Kind) + "|" + t.Path + "|" + strings.Join(t.Patterns, ",") + "|" + t.Capability + "|" + t.RuleID)
}
func recoveryHash(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func canonicalDiagnosisHash(d domain.InconclusiveDiagnosis) string {
	copy := d
	copy.Findings = append([]domain.InconclusiveFinding(nil), d.Findings...)
	sort.Slice(copy.Findings, func(i, j int) bool { return findingSortKey(copy.Findings[i]) < findingSortKey(copy.Findings[j]) })
	return recoveryHash(copy)
}

func ExecuteEvidenceRecovery(ctx context.Context, idea domain.IdeaSpec, screening domain.ScreeningReport, original []domain.EvidenceBundle, before domain.ReuseDecisionReport, plan domain.EvidenceRecoveryPlan, provider domain.ContentProvider, rules domain.CapabilityRuleCatalog, decisionPolicy domain.ReuseDecisionPolicy) (domain.RecoveryRoundResult, error) {
	if ctx == nil {
		return domain.RecoveryRoundResult{}, domain.ErrInvalidRunInput
	}
	if err := ctx.Err(); err != nil {
		return domain.RecoveryRoundResult{}, err
	}
	if plan.Policy.AdditionalBytes <= 0 {
		return domain.RecoveryRoundResult{}, domain.ErrInvalidRunInput
	}
	if len(plan.Targets) == 0 {
		return domain.RecoveryRoundResult{}, fmt.Errorf("plano de recuperação não possui targets")
	}
	plans := recoveryExtractionPlans(plan, original)
	counter := &recoveryProviderCounter{provider: provider}
	observed := provider
	if provider != nil {
		observed = counter
	}
	resolved := map[string]bool{}
	collected := []domain.EvidenceBundle{}
	omitted := []domain.EvidenceBundle{}
	remaining := plan.Policy.AdditionalBytes
	executedTargets := 0
	for _, p := range plans {
		if err := ctx.Err(); err != nil {
			return domain.RecoveryRoundResult{}, err
		}
		if remaining == 0 {
			b := domain.EvidenceBundle{SchemaVersion: "1", Repository: p.Repository, Plan: p}
			for _, target := range plan.Targets {
				if target.Candidate == p.Repository {
					b.Targets = append(b.Targets, domain.CollectionTargetResult{Target: domain.ExtractionTarget{Path: target.Path, Kind: target.Kind, Capabilities: []string{target.Capability}, Reason: "recovery " + target.FindingID}, Status: domain.TargetSkippedBudget, Message: "orçamento de bytes da rodada esgotado"})
				}
			}
			omitted = append(omitted, b)
			continue
		}
		prepared, _, err := prepareRecoveryExtraction(ctx, p, plan.Targets, observed)
		if err != nil {
			return domain.RecoveryRoundResult{}, err
		}
		p = prepared
		p.Budget.TotalBytes = remaining
		groupPlan := domain.EvidenceRecoveryPlan{}
		for _, target := range plan.Targets {
			if target.Candidate == p.Repository {
				groupPlan.Targets = append(groupPlan.Targets, target)
			}
		}
		b, e := collectEvidence(ctx, p, observed, true)
		if err := ctx.Err(); err != nil {
			return domain.RecoveryRoundResult{}, err
		}
		if errors.Is(e, domain.ErrIntegrity) {
			return domain.RecoveryRoundResult{}, domain.ErrIntegrity
		}
		if e != nil {
			b = domain.EvidenceBundle{SchemaVersion: "1", Repository: p.Repository, Plan: p, Failures: []domain.ExecutionFailure{{Scope: p.Repository.Key(), Stage: "recovery", Message: e.Error()}}}
		}
		if e == nil && mutableIdeaRunCommit(p.Repository.Commit) && immutableCommitPattern.MatchString(b.Repository.Commit) {
			resolved[p.Repository.Owner+"/"+p.Repository.Repo] = true
		}
		for _, target := range groupPlan.Targets {
			if recoveryTargetAttempted(target, b) || recoveryTargetPartiallyRead(target, b) {
				executedTargets++
			}
		}
		remaining -= b.BytesRead
		collected = append(collected, b)
	}
	activeOriginal := []domain.EvidenceBundle{}
	for _, b := range original {
		if mutableIdeaRunCommit(b.Repository.Commit) && resolved[b.Repository.Owner+"/"+b.Repository.Repo] {
			continue
		}
		activeOriginal = append(activeOriginal, b)
	}
	merged, duplicates := mergeRecoveryBundles(activeOriginal, append(collected, omitted...))
	valid := []domain.EvidenceBundle{}
	failed := []domain.EvidenceBundle{}
	for _, b := range merged {
		if b.Repository.Commit == "" || strings.EqualFold(b.Repository.Commit, "HEAD") {
			failed = append(failed, b)
		} else {
			valid = append(valid, b)
		}
	}
	assessment, err := AssessEvidence(idea, screening, valid, rules)
	if err != nil {
		return domain.RecoveryRoundResult{}, err
	}
	for _, b := range failed {
		assessment.Assessments = append(assessment.Assessments, inconclusiveRunAssessment(idea, b, rules, assessment.CatalogCoverage))
	}
	decision, err := DecideReuse(domain.ReuseDecisionRequest{Idea: idea, Assessments: assessment.Assessments, Policy: decisionPolicy})
	if err != nil {
		return domain.RecoveryRoundResult{}, err
	}
	diagnosis := DiagnoseInconclusive("recovery", recoveryHash(decision), decision, assessment, merged, rules)
	metrics := recoveryMetrics(plan, before, decision, original, collected, duplicates, diagnosis, counter.calls)
	metrics.TargetsExecuted = executedTargets
	return domain.RecoveryRoundResult{Bundles: merged, Assessment: assessment, Decision: decision, Diagnosis: diagnosis, Metrics: metrics}, nil
}

func recoveryExtractionPlans(plan domain.EvidenceRecoveryPlan, original []domain.EvidenceBundle) []domain.EvidenceExtractionPlan {
	byRepo := map[string]domain.EvidenceBundle{}
	for _, b := range original {
		byRepo[b.Repository.Key()] = b
	}
	grouped := map[string][]domain.RecoveryTarget{}
	for _, t := range plan.Targets {
		grouped[t.Candidate.Key()] = append(grouped[t.Candidate.Key()], t)
	}
	keys := []string{}
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []domain.EvidenceExtractionPlan{}
	for _, key := range keys {
		targets := grouped[key]
		base := byRepo[key]
		candidate := base.Plan.Candidate
		if candidate.Owner == "" {
			candidate = domain.RepositoryCandidate{Owner: targets[0].Candidate.Owner, Repo: targets[0].Candidate.Repo, CommitID: targets[0].Candidate.Commit, URL: "https://github.com/" + targets[0].Candidate.Owner + "/" + targets[0].Candidate.Repo}
		}
		ref := targets[0].Candidate
		candidate.CommitID = ref.Commit
		extraction := []domain.ExtractionTarget{}
		for _, t := range targets {
			if t.Path == "@repository-metadata" {
				continue
			}
			if t.Path != "" {
				extraction = append(extraction, domain.ExtractionTarget{Path: t.Path, Kind: t.Kind, Capabilities: []string{t.Capability}, Reason: "recovery " + t.FindingID})
			}
		}
		extraction = dedupeExtractionTargets(extraction)
		out = append(out, domain.EvidenceExtractionPlan{SchemaVersion: "1", Repository: ref, Candidate: candidate, Decision: base.Plan.Decision, Targets: extraction, Budget: domain.ExtractionBudget{Documents: plan.Policy.AdditionalDocuments, CodeFiles: plan.Policy.AdditionalCodeFiles, Tests: plan.Policy.AdditionalTests, Manifests: 2, TotalBytes: plan.Policy.AdditionalBytes, TreeDepth: plan.Policy.AdditionalTreeDepth}})
	}
	return out
}

func prepareRecoveryExtraction(ctx context.Context, p domain.EvidenceExtractionPlan, targets []domain.RecoveryTarget, provider domain.ContentProvider) (domain.EvidenceExtractionPlan, int, error) {
	derived := []domain.RecoveryTarget{}
	for _, target := range targets {
		if target.Candidate == p.Repository && target.Path == "" {
			derived = append(derived, target)
		}
	}
	if len(derived) == 0 {
		return p, 0, nil
	}
	if err := ctx.Err(); err != nil {
		return p, 0, err
	}
	if provider == nil {
		return p, 0, domain.ErrInvalidRunInput
	}
	tree, err := provider.FetchTree(ctx, p.Candidate)
	if contextErr := ctx.Err(); contextErr != nil {
		return p, 1, contextErr
	}
	if errors.Is(err, domain.ErrIntegrity) {
		return p, 1, domain.ErrIntegrity
	}
	if err == nil {
		for _, node := range tree {
			if node.Type != "file" {
				continue
			}
			for _, target := range derived {
				if recoveryPathMatches(node.Path, target) {
					p.Targets = append(p.Targets, domain.ExtractionTarget{Path: node.Path, Kind: target.Kind, Capabilities: []string{target.Capability}, Reason: "recovery " + target.FindingID})
				}
			}
		}
	}
	p.Targets = dedupeExtractionTargets(p.Targets)
	return p, 1, nil
}
func recoveryPathMatches(path string, t domain.RecoveryTarget) bool {
	kind := classifyEvidencePath(path)
	if t.Kind == domain.EvidenceTest && kind != domain.EvidenceTest {
		return false
	}
	if t.Kind == domain.EvidenceFile && !isCodePath(path) {
		return false
	}
	low := strings.ToLower(slashPath(path))
	if len(t.Patterns) == 0 {
		return true
	}
	for _, p := range t.Patterns {
		if strings.Contains(low, p) {
			return true
		}
	}
	return t.Kind == domain.EvidenceTest
}
func dedupeExtractionTargets(in []domain.ExtractionTarget) []domain.ExtractionTarget {
	seen := map[string]bool{}
	out := []domain.ExtractionTarget{}
	sort.Slice(in, func(i, j int) bool { return in[i].Path < in[j].Path })
	for _, t := range in {
		key := string(t.Kind) + "|" + t.Path
		if !seen[key] {
			seen[key] = true
			out = append(out, t)
		}
	}
	return out
}
func mergeRecoveryBundles(original, newBundles []domain.EvidenceBundle) ([]domain.EvidenceBundle, int) {
	byRepo := map[string]domain.EvidenceBundle{}
	for _, b := range original {
		b.Items = append([]domain.CollectedEvidence(nil), b.Items...)
		b.Targets = append([]domain.CollectionTargetResult(nil), b.Targets...)
		b.Failures = append([]domain.ExecutionFailure(nil), b.Failures...)
		byRepo[b.Repository.Key()] = b
	}
	duplicates := 0
	for _, n := range newBundles {
		key := n.Repository.Key()
		b := byRepo[key]
		if b.Repository.Owner == "" {
			n.Items = append([]domain.CollectedEvidence(nil), n.Items...)
			byRepo[key] = n
			continue
		}
		if !strings.EqualFold(n.Repository.Commit, "HEAD") && n.Repository.Commit != "" {
			b.Repository = n.Repository
			b.Plan.Repository = n.Repository
			b.Plan.Candidate = n.Plan.Candidate
		}
		seen := map[string]bool{}
		for _, i := range b.Items {
			seen[i.ContentHash] = true
		}
		for _, i := range n.Items {
			if seen[i.ContentHash] {
				duplicates++
				continue
			}
			seen[i.ContentHash] = true
			b.Items = append(b.Items, i)
			b.BytesRead += int64(len(i.Content))
		}
		b.Targets = append(b.Targets, n.Targets...)
		b.Failures = append(b.Failures, n.Failures...)
		// Uma tentativa adicional incompleta não invalida a cobertura conclusiva já
		// registrada pelo run original; ela apenas deixa de ampliar essa cobertura.
		b.Complete = b.Complete || n.Complete
		b.TreeComplete = b.TreeComplete || n.TreeComplete
		byRepo[key] = b
	}
	keys := []string{}
	for k := range byRepo {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []domain.EvidenceBundle{}
	for _, k := range keys {
		b := byRepo[k]
		sort.Slice(b.Items, func(i, j int) bool { return b.Items[i].Path < b.Items[j].Path })
		out = append(out, b)
	}
	return out, duplicates
}
func recoveryMetrics(plan domain.EvidenceRecoveryPlan, before, after domain.ReuseDecisionReport, original, collected []domain.EvidenceBundle, duplicates int, d domain.InconclusiveDiagnosis, calls int) domain.RecoveryMetrics {
	bytes := int64(0)
	for _, b := range collected {
		bytes += b.BytesRead
	}
	m := domain.RecoveryMetrics{Round: plan.Round, FindingsBefore: len(plan.Targets) + len(plan.Skipped), TargetsPlanned: len(plan.Targets), TargetsExecuted: len(plan.Targets), TargetsNew: len(plan.Targets), DuplicateContent: duplicates, AdditionalBytes: bytes, ProviderCalls: calls, RemainingFindings: len(d.Findings), DecisionBefore: string(before.Decision), DecisionAfter: string(after.Decision)}
	m.ResolvedFindings = m.FindingsBefore - m.RemainingFindings
	if m.ResolvedFindings < 0 {
		m.ResolvedFindings = 0
	}
	if m.TargetsNew > 0 {
		x := float64(m.ResolvedFindings) / float64(m.TargetsNew)
		m.UsefulEvidenceRate = &x
	}
	if m.FindingsBefore > 0 {
		x := float64(m.ResolvedFindings) / float64(m.FindingsBefore)
		m.InconclusiveResolution = &x
	}
	if before.Decision != after.Decision {
		m.StatusTransitions = 1
	}
	return m
}

func OriginalRecoveryTargets(bundles []domain.EvidenceBundle) []domain.RecoveryTarget {
	out := []domain.RecoveryTarget{}
	for _, b := range bundles {
		for _, target := range b.Plan.Targets {
			out = append(out, domain.RecoveryTarget{Candidate: b.Repository, Capability: firstString(target.Capabilities), Path: target.Path, Kind: target.Kind, RuleID: recoveryHistoryRule(target.Kind)})
		}
		for _, target := range b.Targets {
			out = append(out, domain.RecoveryTarget{Candidate: b.Repository, Capability: firstString(target.Target.Capabilities), Path: target.Target.Path, Kind: target.Target.Kind, RuleID: recoveryHistoryRule(target.Target.Kind)})
		}
	}
	return out
}
func recoveryHistoryRule(kind domain.EvidenceKind) string {
	switch kind {
	case domain.EvidenceLicense:
		return "GATE-LICENSE-001"
	case domain.EvidenceManifest:
		return "GATE-OPERATION-001"
	}
	return "original"
}
func firstString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// ExecutedRecoveryTargets retains planned identities only when outcomes witness an attempt.
func ExecutedRecoveryTargets(plan domain.EvidenceRecoveryPlan, bundles []domain.EvidenceBundle) []domain.RecoveryTarget {
	out := []domain.RecoveryTarget{}
	for _, target := range plan.Targets {
		matches := []domain.EvidenceBundle{}
		for _, b := range bundles {
			if b.Repository.Key() == target.Candidate.Key() {
				matches = append(matches, b)
			}
		}
		if len(matches) == 0 && mutableIdeaRunCommit(target.Candidate.Commit) {
			for _, b := range bundles {
				if b.Repository.Owner == target.Candidate.Owner && b.Repository.Repo == target.Candidate.Repo && immutableCommitPattern.MatchString(b.Repository.Commit) {
					matches = append(matches, b)
				}
			}
		}
		if len(matches) == 1 && recoveryTargetAttempted(target, matches[0]) {
			out = append(out, target)
		}
	}
	return out
}
func recoveryAttemptStatus(status domain.CollectionTargetStatus) bool {
	return status == domain.TargetCollected || status == domain.TargetNotFound || status == domain.TargetFailed
}
func recoveryTargetAttempted(target domain.RecoveryTarget, bundle domain.EvidenceBundle) bool {
	for i := len(bundle.Targets) - 1; i >= 0; i-- {
		result := bundle.Targets[i]
		if recoveryTargetResultMatches(target, result) {
			return recoveryAttemptStatus(result.Status)
		}
	}
	for _, item := range bundle.Items {
		if item.Path == target.Path || (target.Path == "" && (item.Path == "@repository-metadata" || recoveryPathMatches(item.Path, target))) {
			return true
		}
	}
	for _, failure := range bundle.Failures {
		if failure.Stage == "recovery" {
			return true
		}
	}
	return false
}

// The counter observes port invocations; it performs no concrete effects itself.
type recoveryProviderCounter struct {
	provider domain.ContentProvider
	calls    int
}

func (p *recoveryProviderCounter) FetchManifest(ctx context.Context, c domain.RepositoryCandidate) (*domain.Manifest, error) {
	p.calls++
	return p.provider.FetchManifest(ctx, c)
}
func (p *recoveryProviderCounter) FetchRepoMetadata(ctx context.Context, c domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.calls++
	return p.provider.FetchRepoMetadata(ctx, c)
}
func (p *recoveryProviderCounter) FetchTree(ctx context.Context, c domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls++
	return p.provider.FetchTree(ctx, c)
}
func (p *recoveryProviderCounter) ExtractRepoFiles(ctx context.Context, c domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	p.calls++
	return p.provider.ExtractRepoFiles(ctx, c, paths)
}

// AttemptedOriginalRecoveryTargets is the active history; the legacy projection stays intact.
func AttemptedOriginalRecoveryTargets(bundles []domain.EvidenceBundle) []domain.RecoveryTarget {
	observed := []domain.EvidenceBundle{}
	for _, b := range bundles {
		filtered := b
		filtered.Plan.Targets = nil
		filtered.Targets = nil
		for _, target := range b.Plan.Targets {
			if recoveryTargetAttempted(domain.RecoveryTarget{Path: target.Path, Kind: target.Kind}, b) {
				filtered.Plan.Targets = append(filtered.Plan.Targets, target)
			}
		}
		for _, result := range b.Targets {
			if recoveryAttemptStatus(result.Status) {
				filtered.Targets = append(filtered.Targets, result)
			}
		}
		observed = append(observed, filtered)
	}
	return OriginalRecoveryTargets(observed)
}

func recoveryTargetResultMatches(target domain.RecoveryTarget, result domain.CollectionTargetResult) bool {
	return (target.FindingID != "" && result.Target.Reason == "recovery "+target.FindingID) || (target.Path != "" && result.Target.Path == target.Path) || (target.Path == "" && result.Target.Kind == target.Kind && recoveryPathMatches(result.Target.Path, target))
}
func recoveryTargetPartiallyRead(target domain.RecoveryTarget, bundle domain.EvidenceBundle) bool {
	for _, result := range bundle.Targets {
		if result.Status == domain.TargetSkippedBudget && result.Bytes > 0 && recoveryTargetResultMatches(target, result) {
			return true
		}
	}
	return false
}
