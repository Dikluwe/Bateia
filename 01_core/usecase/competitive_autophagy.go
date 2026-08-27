/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/competitive_autophagy.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"fmt"
	"sort"
	"strings"
)

// RunCompetitiveAutophagy assesses frozen internal and external bytes and then
// enumerates the smallest admissible composition. It performs no I/O.
func RunCompetitiveAutophagy(req domain.CompetitiveAutophagyRequest) (domain.CompetitiveAutophagyResult, error) {
	if req.SchemaVersion != "1" || req.Target.Identity.Owner == "" || req.Target.Identity.Repo == "" || req.Target.Commit == "" || strings.EqualFold(req.Target.Commit, "HEAD") {
		return domain.CompetitiveAutophagyResult{}, fmt.Errorf("autofagia exige schema 1 e alvo em commit imutável")
	}
	if req.KatamariPolicy.SchemaVersion != "1" || req.KatamariPolicy.MaxProjects <= 0 || req.KatamariPolicy.MaxComponents <= 0 {
		return domain.CompetitiveAutophagyResult{}, fmt.Errorf("política Katamari inválida")
	}
	internal, err := AssessEvidence(req.Idea, domain.ScreeningReport{}, req.Internal, req.Rules)
	if err != nil {
		return domain.CompetitiveAutophagyResult{}, fmt.Errorf("evidência interna: %w", err)
	}
	autophagy := buildAutophagyReport(req, internal)
	external, err := AssessEvidence(req.Idea, domain.ScreeningReport{}, req.Competitors, req.Rules)
	if err != nil {
		return domain.CompetitiveAutophagyResult{}, fmt.Errorf("evidência competitiva: %w", err)
	}
	katamari, err := buildKatamari(req, autophagy, external)
	if err != nil {
		return domain.CompetitiveAutophagyResult{}, err
	}
	return domain.CompetitiveAutophagyResult{SchemaVersion: "1", Autophagy: autophagy, Katamari: katamari}, nil
}

func AnalyzeAutophagy(idea domain.IdeaSpec, target domain.ResolvedRepository, bundles []domain.EvidenceBundle, rules domain.CapabilityRuleCatalog) (domain.AutophagyReport, error) {
	req := domain.CompetitiveAutophagyRequest{SchemaVersion: "1", Idea: idea, Target: target, Internal: bundles, Rules: rules}
	if target.Identity.Owner == "" || target.Identity.Repo == "" || target.Commit == "" || strings.EqualFold(target.Commit, "HEAD") {
		return domain.AutophagyReport{}, fmt.Errorf("autofagia exige alvo em commit imutável")
	}
	assessment, err := AssessEvidence(idea, domain.ScreeningReport{}, bundles, rules)
	if err != nil {
		return domain.AutophagyReport{}, err
	}
	return buildAutophagyReport(req, assessment), nil
}

func ProposeKatamari(autophagy domain.AutophagyReport, assessments domain.AssessmentReport, policy domain.KatamariPolicy, reusePolicy domain.ReuseDecisionPolicy) (domain.KatamariProposal, error) {
	if policy.SchemaVersion != "1" || policy.MaxProjects <= 0 || policy.MaxComponents <= 0 {
		return domain.KatamariProposal{}, fmt.Errorf("política Katamari inválida")
	}
	req := domain.CompetitiveAutophagyRequest{KatamariPolicy: policy, ReusePolicy: reusePolicy}
	return buildKatamari(req, autophagy, assessments)
}

func buildAutophagyReport(req domain.CompetitiveAutophagyRequest, report domain.AssessmentReport) domain.AutophagyReport {
	out := domain.AutophagyReport{SchemaVersion: "1", Target: req.Target, Findings: []domain.InternalCapabilityFinding{}, Existing: []string{}, Partial: []string{}, Gaps: []string{}, Unknown: []string{}}
	for _, capability := range req.Idea.Capabilities {
		finding := aggregateInternalCapability(capability, report.Assessments)
		out.Findings = append(out.Findings, finding)
		switch finding.Disposition {
		case domain.InternalKeep, domain.InternalConsolidate:
			out.Existing = append(out.Existing, capability)
		case domain.InternalGap:
			out.Gaps = append(out.Gaps, capability)
		default:
			if finding.Status == domain.EvidencePartial {
				out.Partial = append(out.Partial, capability)
			} else {
				out.Unknown = append(out.Unknown, capability)
			}
		}
	}
	sortAutophagyReport(&out)
	out.InputHash = recoveryHash(struct {
		Idea     domain.IdeaSpec
		Target   domain.ResolvedRepository
		Internal []domain.EvidenceBundle
		Rules    domain.CapabilityRuleCatalog
	}{req.Idea, req.Target, req.Internal, req.Rules})
	return out
}

func aggregateInternalCapability(capability string, assessments []domain.CandidateAssessment) domain.InternalCapabilityFinding {
	f := domain.InternalCapabilityFinding{Capability: capability, Sources: []domain.EvidenceSource{}, RuleIDs: []string{}}
	proven, absent, opaque := 0, 0, 0
	for _, a := range assessments {
		for _, e := range a.Capabilities {
			if normalizeEvidenceText(e.Capability) != normalizeEvidenceText(capability) {
				continue
			}
			f.Sources = append(f.Sources, e.Sources...)
			for _, s := range e.Sources {
				if s.RuleID != "" {
					f.RuleIDs = append(f.RuleIDs, s.RuleID)
				}
			}
			switch e.Status {
			case domain.EvidenceProven:
				proven++
			case domain.EvidenceAbsent:
				absent++
			default:
				opaque++
				if f.Status != domain.EvidenceConflict {
					f.Status = e.Status
				}
			}
		}
	}
	f.RuleIDs = uniqueSorted(f.RuleIDs)
	sort.Slice(f.Sources, func(i, j int) bool { return evidenceSourceKey(f.Sources[i]) < evidenceSourceKey(f.Sources[j]) })
	switch {
	case proven > 1:
		f.Status, f.Disposition, f.Rationale = domain.EvidenceProven, domain.InternalConsolidate, "capacidade comprovada em mais de uma implementação interna"
	case proven == 1:
		f.Status, f.Disposition, f.Rationale = domain.EvidenceProven, domain.InternalKeep, "capacidade interna comprovada"
	case opaque > 0 || len(assessments) == 0:
		if f.Status == "" {
			f.Status = domain.EvidenceInconclusive
		}
		f.Disposition, f.Rationale = domain.InternalUnknown, "evidência interna parcial, conflitante ou opaca"
	case absent > 0:
		f.Status, f.Disposition, f.Rationale = domain.EvidenceAbsent, domain.InternalGap, "ausência sustentada pela coleta interna completa"
	default:
		f.Status, f.Disposition, f.Rationale = domain.EvidenceInconclusive, domain.InternalUnknown, "capacidade sem avaliação interna correspondente"
	}
	return f
}

type katamariCandidate struct {
	component domain.KatamariComponent
	covered   map[string]bool
}

func buildKatamari(req domain.CompetitiveAutophagyRequest, auto domain.AutophagyReport, report domain.AssessmentReport) (domain.KatamariProposal, error) {
	out := domain.KatamariProposal{SchemaVersion: "1", Components: []domain.KatamariComponent{}, CoveredGaps: []string{}, UncoveredGaps: append([]string(nil), auto.Gaps...), Rejected: []domain.DecisionAlternative{}, Minimal: true}
	candidates := []katamariCandidate{}
	for _, a := range report.Assessments {
		candidate, reason := admissibleKatamariCandidate(req, auto.Gaps, a)
		if reason != "" {
			out.Rejected = append(out.Rejected, domain.DecisionAlternative{Decision: domain.DecisionUnknown, Candidate: a.Repository.Key(), Message: reason})
			continue
		}
		if len(candidate.covered) == 0 {
			out.Rejected = append(out.Rejected, domain.DecisionAlternative{Decision: domain.DecisionReference, Candidate: a.Repository.Key(), Message: "nenhuma lacuna interna comprovada é coberta"})
			continue
		}
		candidates = append(candidates, candidate)
	}
	if len(candidates) > 20 {
		return out, fmt.Errorf("enumeração Katamari limitada a 20 candidatos, recebeu %d", len(candidates))
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].component.Repository.Key() < candidates[j].component.Repository.Key()
	})
	best := selectKatamariSubset(candidates, auto.Gaps, req.KatamariPolicy)
	for _, index := range best {
		out.Components = append(out.Components, candidates[index].component)
	}
	covered := map[string]bool{}
	for _, i := range best {
		for gap := range candidates[i].covered {
			covered[gap] = true
		}
	}
	out.CoveredGaps, out.UncoveredGaps = partitionGaps(auto.Gaps, covered)
	out.Feasible = len(out.Components) > 0
	sort.Slice(out.Rejected, func(i, j int) bool { return out.Rejected[i].Candidate < out.Rejected[j].Candidate })
	out.InputHash = recoveryHash(struct {
		Autophagy domain.AutophagyReport
		External  domain.AssessmentReport
		Policy    domain.KatamariPolicy
		Reuse     domain.ReuseDecisionPolicy
	}{auto, report, req.KatamariPolicy, req.ReusePolicy})
	return out, nil
}

func admissibleKatamariCandidate(req domain.CompetitiveAutophagyRequest, gaps []string, a domain.CandidateAssessment) (katamariCandidate, string) {
	c := katamariCandidate{component: domain.KatamariComponent{Repository: a.Repository, Decision: domain.DecisionIntegrate, Sources: []domain.EvidenceSource{}, Risks: []string{}}, covered: map[string]bool{}}
	if req.KatamariPolicy.RequireImmutableRef && (a.Repository.Commit == "" || strings.EqualFold(a.Repository.Commit, "HEAD")) {
		return c, "revisão não imutável"
	}
	if len(a.Failures) > 0 && !req.KatamariPolicy.AllowUnknown {
		return c, "coleta competitiva possui falhas"
	}
	if !licenseAllowed(a.License, req.KatamariPolicy.AllowedLicenses) {
		return c, "licença ausente, opaca ou não permitida"
	}
	if req.ReusePolicy.RequireLocal || strings.EqualFold(req.KatamariPolicy.RequiredOperation, "local") {
		if a.Operation.LocalStatus != domain.EvidenceProven {
			return c, "operação local não comprovada"
		}
	}
	if a.Maintenance.Archived && !req.ReusePolicy.AllowArchived {
		return c, "repositório arquivado proibido pela política"
	}
	gapSet := normalizedSet(gaps)
	for _, e := range a.Capabilities {
		key := normalizeEvidenceText(e.Capability)
		if e.Status == domain.EvidenceProven && gapSet[key] {
			c.covered[key] = true
			c.component.Capabilities = append(c.component.Capabilities, e.Capability)
			c.component.Sources = append(c.component.Sources, e.Sources...)
		}
	}
	c.component.Capabilities = uniqueSorted(c.component.Capabilities)
	sort.Slice(c.component.Sources, func(i, j int) bool {
		return evidenceSourceKey(c.component.Sources[i]) < evidenceSourceKey(c.component.Sources[j])
	})
	c.component.Cost = 100 + 10*len(c.component.Capabilities)
	if a.Maintenance.Archived {
		c.component.Cost += 50
		c.component.Risks = append(c.component.Risks, "archived")
	}
	return c, ""
}

func selectKatamariSubset(candidates []katamariCandidate, gaps []string, policy domain.KatamariPolicy) []int {
	best := []int{}
	bestCoverage, bestCost, bestRedundancy, bestKey := -1, 0, 0, ""
	for mask := 1; mask < (1 << len(candidates)); mask++ {
		indices, cost, occurrences := []int{}, 0, map[string]int{}
		for i := range candidates {
			if mask&(1<<i) == 0 {
				continue
			}
			indices = append(indices, i)
			cost += candidates[i].component.Cost
			for g := range candidates[i].covered {
				occurrences[g]++
			}
		}
		if len(indices) > policy.MaxProjects || len(indices) > policy.MaxComponents {
			continue
		}
		coverage, redundancy := 0, 0
		for _, n := range occurrences {
			coverage++
			if n > 1 {
				redundancy += n - 1
			}
		}
		keyParts := []string{}
		for _, i := range indices {
			keyParts = append(keyParts, candidates[i].component.Repository.Key())
		}
		key := strings.Join(keyParts, ",")
		better := coverage > bestCoverage || (coverage == bestCoverage && (len(best) == 0 || len(indices) < len(best))) || (coverage == bestCoverage && len(indices) == len(best) && (cost < bestCost || (cost == bestCost && (redundancy < bestRedundancy || (redundancy == bestRedundancy && key < bestKey)))))
		if better {
			best, bestCoverage, bestCost, bestRedundancy, bestKey = append([]int(nil), indices...), coverage, cost, redundancy, key
		}
	}
	return best
}

func licenseAllowed(l domain.LicenseAssessment, allowed []string) bool {
	if l.Status != domain.EvidenceProven {
		return false
	}
	for _, x := range allowed {
		if strings.EqualFold(strings.TrimSpace(x), strings.TrimSpace(l.Identifier)) {
			return true
		}
	}
	return false
}
func normalizedSet(xs []string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[normalizeEvidenceText(x)] = true
	}
	return m
}
func partitionGaps(gaps []string, covered map[string]bool) ([]string, []string) {
	yes, no := []string{}, []string{}
	for _, gap := range gaps {
		if covered[normalizeEvidenceText(gap)] {
			yes = append(yes, gap)
		} else {
			no = append(no, gap)
		}
	}
	sort.Strings(yes)
	sort.Strings(no)
	return yes, no
}
func sortAutophagyReport(r *domain.AutophagyReport) {
	sort.Slice(r.Findings, func(i, j int) bool { return r.Findings[i].Capability < r.Findings[j].Capability })
	sort.Strings(r.Existing)
	sort.Strings(r.Partial)
	sort.Strings(r.Gaps)
	sort.Strings(r.Unknown)
}
func uniqueSorted(xs []string) []string {
	seen, out := map[string]bool{}, []string{}
	for _, x := range xs {
		if x != "" && !seen[x] {
			seen[x], out = true, append(out, x)
		}
	}
	sort.Strings(out)
	return out
}
func evidenceSourceKey(s domain.EvidenceSource) string {
	return s.Ref.Key() + "|" + s.Path + "|" + s.Symbol + "|" + s.RuleID
}
