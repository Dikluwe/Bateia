/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/reuse_decision.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type decisionCandidate struct {
	a         domain.CandidateAssessment
	key       string
	coverage  float64
	covered   map[string]bool
	evidence  []domain.EvidenceSource
	blockers  []domain.DecisionReason
	ambiguous bool
	useful    bool
}

func DecideReuse(req domain.ReuseDecisionRequest) (result domain.ReuseDecisionReport, err error) {
	return decideReuseProjected(req, nil)
}

// decisionAssessmentProjection is used only by bounded counterfactual analysis.
// It runs after evidence admission, so hypothetical states never become sources.
type decisionAssessmentProjection func(domain.CandidateAssessment) domain.CandidateAssessment

func decideReuseProjected(req domain.ReuseDecisionRequest, projection decisionAssessmentProjection) (result domain.ReuseDecisionReport, err error) {
	defer func() {
		if err == nil && projection == nil {
			plan := BuildInvestigationPlan(req.Idea, domain.DiscoverySnapshot{}, req.Assessments)
			result.Investigation = &plan
		}
	}()
	if err := validateDecisionRequest(req); err != nil {
		return domain.ReuseDecisionReport{}, err
	}
	policyJSON, err := json.Marshal(req.Policy)
	if err != nil {
		return domain.ReuseDecisionReport{}, fmt.Errorf("não foi possível serializar a política: %w", err)
	}
	h := sha256.Sum256(policyJSON)
	report := domain.ReuseDecisionReport{SchemaVersion: "1", Policy: req.Policy, PolicyHash: hex.EncodeToString(h[:]), Alternatives: []domain.DecisionAlternative{}}
	capabilities := normalizedSorted(req.Idea.Capabilities)
	if len(req.Idea.Qualities) > 0 || len(req.Idea.OperatingConditions) > 0 {
		report.Decision = domain.DecisionUnknown
		report.Matrix = buildDecisionMatrix(capabilities, nil)
		report.Gaps = append(normalizeIdeaList(req.Idea.Qualities), normalizeIdeaList(req.Idea.OperatingConditions)...)
		report.Reasons = []domain.DecisionReason{{RuleID: "DEC-REQUIREMENTS-0087", Effect: "verify", Message: "qualidades e condições operacionais declaradas ainda não são avaliadas por candidato; verificar os critérios compilados antes de decidir"}}
		report.Effort = assessEffort(report.Decision, 0, len(report.Gaps), true)
		return report, nil
	}
	if len(req.Assessments) == 0 {
		report.Decision = domain.DecisionUnknown
		report.Matrix = buildDecisionMatrix(capabilities, nil)
		report.Gaps = append([]string(nil), capabilities...)
		report.Reasons = []domain.DecisionReason{{RuleID: "DEC-INCONCLUSIVE-002", Effect: "block", Message: "nenhum assessment foi fornecido; o corpus não sustenta decisão"}}
		report.Effort = assessEffort(report.Decision, 0, len(report.Gaps), true)
		return report, nil
	}
	capabilities = attributableDecisionCapabilities(capabilities, req.Assessments)
	if len(capabilities) == 0 {
		report.Decision = domain.DecisionUnknown
		report.Matrix = []domain.CapabilityComparison{}
		report.Gaps = []string{}
		report.Reasons = []domain.DecisionReason{{RuleID: "DEC-CATALOG-0071", Effect: "verify", Message: "o catálogo não possui regra declarada segura para atribuir presença ou ausência das capacidades aos candidatos"}}
		report.Effort = assessEffort(report.Decision, 0, 0, true)
		return report, nil
	}
	assessments := append([]domain.CandidateAssessment(nil), req.Assessments...)
	sort.Slice(assessments, func(i, j int) bool { return assessments[i].Repository.Key() < assessments[j].Repository.Key() })
	candidates := make([]decisionCandidate, 0, len(assessments))
	for _, a := range assessments {
		a = SanitizeAssessmentEvidence(a)
		if projection != nil {
			a = projection(a)
		}
		candidates = append(candidates, inspectPreparedCandidate(a, capabilities, req.Policy))
	}
	report.Matrix = buildDecisionMatrix(capabilities, candidates)
	composition, err := analyzeCompositions(candidates, capabilities, req.Policy, req.Composition)
	if err != nil {
		return domain.ReuseDecisionReport{}, err
	}
	report.Composition = &composition
	for _, c := range candidates {
		report.Alternatives = append(report.Alternatives, domain.DecisionAlternative{Decision: domain.DecisionReference, Candidate: c.key, Coverage: c.coverage, Message: fmt.Sprintf("%d%% das capacidades obrigatórias comprovadas", int(c.coverage*100))})
	}
	for _, candidate := range composition.Candidates {
		report.Alternatives = append(report.Alternatives, domain.DecisionAlternative{Decision: domain.DecisionCompose, Candidate: candidate.ID, Coverage: float64(len(candidate.Coverage)) / float64(len(capabilities)), Message: string(candidate.Status)})
	}
	if c := firstCandidate(candidates, func(c decisionCandidate) bool {
		return c.coverage == 1 && len(c.blockers) == 0 && !c.ambiguous && c.a.Operation.LocalStatus == domain.EvidenceProven
	}); c != nil {
		return finishDecision(report, capabilities, *c, domain.DecisionFork, "DEC-FORK-001", "candidato executável cobre todas as capacidades e passa pelos gates"), nil
	}
	if c := firstCandidate(candidates, func(c decisionCandidate) bool {
		return c.coverage >= req.Policy.ExtendMinCoverage && c.coverage < 1 && c.a.ProblemAlignment.Status == domain.EvidenceProven && len(validParts(c.a.ReusableParts)) > 0 && len(c.blockers) == 0 && !c.ambiguous
	}); c != nil {
		return finishDecision(report, capabilities, *c, domain.DecisionExtend, "DEC-EXTEND-001", "núcleo reutilizável cobre a capacidade central e o limiar da política"), nil
	}
	for _, candidate := range composition.Candidates {
		if candidate.Status == domain.CompositionVerified {
			report.Decision = domain.DecisionCompose
			for _, member := range candidate.Members {
				report.Components = append(report.Components, member.Component)
			}
			report.Reasons = []domain.DecisionReason{{RuleID: "DEC-COMPOSE-0088", Effect: "select", Candidate: candidate.ID, Evidence: candidate.Evidence, Message: "composição cobre as capacidades, interfaces conectadas e ensaio do conjunto declarado como executado e aprovado"}}
			report.Gaps = candidate.Gaps
			report.Effort = assessEffort(report.Decision, len(candidate.Members), 0, false)
			return report, nil
		}
	}
	for _, candidate := range composition.Candidates {
		if candidate.Status == domain.CompositionProposed {
			report.Decision = domain.DecisionUnknown
			report.Gaps = candidate.Gaps
			report.Reasons = []domain.DecisionReason{{RuleID: "DEC-COMPOSE-ABSTAIN-0088", Effect: "verify", Candidate: candidate.ID, Evidence: candidate.Evidence, Message: "cobertura complementar é hipótese; verificar interfaces, dependências e ensaio de integração antes de COMPOSE"}}
			report.Effort = assessEffort(report.Decision, len(candidate.Members), 0, true)
			return report, nil
		}
	}
	if composition.Truncated {
		report.Decision = domain.DecisionUnknown
		report.Gaps = append([]string(nil), capabilities...)
		report.Reasons = []domain.DecisionReason{{RuleID: "DEC-COMPOSE-LIMIT-0088", Effect: "verify", Message: "busca limitada de composições; ausência de solução não foi comprovada"}}
		report.Effort = assessEffort(report.Decision, 0, len(report.Gaps), true)
		return report, nil
	}
	if c := firstCandidate(candidates, func(c decisionCandidate) bool {
		return c.coverage >= req.Policy.IntegrationMinValue && len(validParts(c.a.ReusableParts)) > 0 && len(c.blockers) == 0 && !c.ambiguous
	}); c != nil {
		return finishDecision(report, capabilities, *c, domain.DecisionIntegrate, "DEC-INTEGRATE-001", "componente verificável resolve uma capacidade delimitada importante"), nil
	}
	anyAmbiguous := false
	for _, c := range candidates {
		anyAmbiguous = anyAmbiguous || c.ambiguous
		report.Blockers = append(report.Blockers, c.blockers...)
	}
	if anyAmbiguous {
		report.Decision = domain.DecisionUnknown
		report.Reasons = []domain.DecisionReason{{RuleID: "DEC-INCONCLUSIVE-001", Effect: "block", Message: "evidência obrigatória ou gate eliminatório permanece inconclusivo"}}
		covered := map[string]bool{}
		catalogLimited := false
		for _, candidate := range candidates {
			for capability := range candidate.covered {
				covered[capability] = true
			}
			catalogLimited = catalogLimited || assessmentHasCatalogLimitation(candidate.a)
			report.Reasons = append(report.Reasons, unresolvedCandidateReasons(candidate, req.Policy)...)
		}
		if catalogLimited {
			report.Reasons = append(report.Reasons, domain.DecisionReason{RuleID: "DEC-CATALOG-0071", Effect: "verify", Message: "a cobertura do catálogo é ausente ou apenas hipotética para parte dos requisitos; isso não descreve lacuna de nenhum candidato"})
		}
		report.Gaps = uncoveredGaps(capabilities, covered)
		report.Effort = assessEffort(report.Decision, 0, len(report.Gaps), true)
		return report, nil
	}
	if c := firstCandidate(candidates, func(c decisionCandidate) bool { return c.useful }); c != nil {
		r := finishDecision(report, capabilities, *c, domain.DecisionReference, "DEC-REFERENCE-001", "há evidência útil, mas gates ou ausência de componente impedem reúso direto")
		r.Blockers = report.Blockers
		return r, nil
	}
	report.Decision = domain.DecisionBuild
	report.Reasons = []domain.DecisionReason{{RuleID: "DEC-BUILD-001", Effect: "select", Message: "avaliação conclusiva não encontrou base reutilizável para a capacidade central"}}
	report.Gaps = append([]string(nil), capabilities...)
	report.Effort = assessEffort(report.Decision, 0, len(report.Gaps), false)
	return report, nil
}

func attributableDecisionCapabilities(capabilities []string, assessments []domain.CandidateAssessment) []string {
	out := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		found := false
		catalogLimited := true
		for _, assessment := range assessments {
			for _, evidence := range assessment.Capabilities {
				if normalizeDecisionText(evidence.Capability) != capability {
					continue
				}
				found = true
				if evidence.ReasonCode != "CATALOG_GAP" && evidence.ReasonCode != "HYPOTHESIS_ONLY" {
					catalogLimited = false
				}
			}
		}
		if !found || !catalogLimited {
			out = append(out, capability)
		}
	}
	return out
}

// Unknown adoption gates are evidence gaps, not missing project capabilities.
func unresolvedCandidateReasons(c decisionCandidate, policy domain.ReuseDecisionPolicy) []domain.DecisionReason {
	var out []domain.DecisionReason
	add := func(message string, evidence []domain.EvidenceSource) {
		out = append(out, domain.DecisionReason{RuleID: "DEC-EVIDENCE-0040", Effect: "verify", Candidate: c.key, Message: message, Evidence: evidence})
	}
	if c.a.License.Status == domain.EvidenceInconclusive {
		add("licença não verificada: localizar documento de licença com autoridade no projeto", c.a.License.Sources)
	}
	if c.a.Maintenance.Status == domain.EvidenceInconclusive {
		add("manutenção não verificada: obter metadados datados de atividade e arquivamento", c.a.Maintenance.Sources)
	}
	if policy.RequireLocal && (c.a.Operation.LocalStatus == domain.EvidenceInconclusive || c.a.Operation.LocalStatus == domain.EvidencePartial) {
		add("operação local não verificada: examinar requisitos e instruções de execução", c.a.Operation.Sources)
	}
	for _, capability := range c.a.Capabilities {
		if catalogLimitedEvidence(capability) {
			continue
		}
		if capability.Status == domain.EvidencePartial || capability.Status == domain.EvidenceInconclusive {
			add("capacidade não comprovada: "+capability.Capability+"; "+capability.Explanation, capability.Sources)
		}
	}
	for _, constraint := range c.a.Constraints {
		if catalogLimitedEvidence(constraint.Evidence) {
			continue
		}
		if constraint.Evidence.Status == domain.EvidenceInconclusive {
			add("restrição não verificada: "+constraint.Constraint+"; "+constraint.Evidence.Explanation, constraint.Evidence.Sources)
		}
	}
	if c.a.ProblemAlignment.Status == domain.EvidenceInconclusive && !catalogLimitedEvidence(c.a.ProblemAlignment) {
		add("relação com a intenção não verificada: "+c.a.ProblemAlignment.Explanation, c.a.ProblemAlignment.Sources)
	}
	for _, failure := range c.a.Failures {
		add("resolver falha de coleta em "+failure.Stage+": "+failure.Message, nil)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Message < out[j].Message })
	return out
}

func assessmentHasCatalogLimitation(assessment domain.CandidateAssessment) bool {
	if catalogLimitedEvidence(assessment.ProblemAlignment) {
		return true
	}
	for _, capability := range assessment.Capabilities {
		if catalogLimitedEvidence(capability) {
			return true
		}
	}
	for _, constraint := range assessment.Constraints {
		if catalogLimitedEvidence(constraint.Evidence) {
			return true
		}
	}
	return false
}

func validateDecisionRequest(req domain.ReuseDecisionRequest) error {
	if len(req.Idea.Capabilities) == 0 {
		return fmt.Errorf("a ideia exige ao menos uma capacidade obrigatória")
	}
	if req.Policy.SchemaVersion == "" {
		return fmt.Errorf("policy.schema_version é obrigatório")
	}
	if req.Policy.MaxComponents < 2 {
		return fmt.Errorf("policy.max_components deve ser ao menos 2")
	}
	if !(req.Policy.ExtendMinCoverage > 0 && req.Policy.ExtendMinCoverage <= 1) || !(req.Policy.IntegrationMinValue > 0 && req.Policy.IntegrationMinValue <= 1) {
		return fmt.Errorf("limiares da política devem estar em (0,1]")
	}
	return nil
}

func inspectCandidate(a domain.CandidateAssessment, caps []string, p domain.ReuseDecisionPolicy) decisionCandidate {
	return inspectPreparedCandidate(SanitizeAssessmentEvidence(a), caps, p)
}

func inspectPreparedCandidate(a domain.CandidateAssessment, caps []string, p domain.ReuseDecisionPolicy) decisionCandidate {
	c := decisionCandidate{a: a, key: a.Repository.Key(), covered: map[string]bool{}}
	if assessmentHasCatalogLimitation(a) {
		c.ambiguous = true
	}
	byCap := map[string]domain.CapabilityEvidence{}
	for _, e := range a.Capabilities {
		byCap[normalizeDecisionText(e.Capability)] = e
	}
	for _, cap := range caps {
		e, ok := byCap[cap]
		if !ok {
			c.ambiguous = true
			continue
		}
		c.evidence = append(c.evidence, e.Sources...)
		if e.Status == domain.EvidenceProven {
			c.covered[cap] = true
			c.useful = true
		}
		if e.Status == domain.EvidencePartial {
			c.useful = true
		}
		if e.Status == domain.EvidenceInconclusive {
			c.ambiguous = true
		}
		if e.Status == domain.EvidenceConflict {
			c.blockers = append(c.blockers, domain.DecisionReason{RuleID: "GATE-CAPABILITY-001", Effect: "block", Candidate: c.key, Evidence: e.Sources, Message: "capacidade obrigatória em conflito: " + cap})
		}
	}
	c.coverage = float64(len(c.covered)) / float64(len(caps))
	if catalogLimitedEvidence(a.ProblemAlignment) {
		c.ambiguous = true
	}
	if a.ProblemAlignment.Status == domain.EvidenceInconclusive {
		c.ambiguous = true
	}
	if len(a.Failures) > 0 {
		c.ambiguous = true
	}
	if a.Maintenance.Status == domain.EvidenceInconclusive {
		c.ambiguous = true
	}
	if a.License.Status == domain.EvidenceInconclusive {
		c.ambiguous = true
	}
	if a.License.Status != domain.EvidenceInconclusive && (a.License.Status != domain.EvidenceProven || !containsFold(p.AllowedLicenses, a.License.Identifier)) {
		c.blockers = append(c.blockers, domain.DecisionReason{RuleID: "GATE-LICENSE-001", Effect: "block", Candidate: c.key, Evidence: a.License.Sources, Message: "licença não permitida pela política"})
	}
	if a.Maintenance.Archived && !p.AllowArchived {
		c.blockers = append(c.blockers, domain.DecisionReason{RuleID: "GATE-MAINTENANCE-001", Effect: "block", Candidate: c.key, Evidence: a.Maintenance.Sources, Message: "repositório arquivado não permitido"})
	}
	if p.RequireLocal && a.Operation.LocalStatus != domain.EvidenceProven {
		if a.Operation.LocalStatus == domain.EvidenceInconclusive || a.Operation.LocalStatus == domain.EvidencePartial {
			c.ambiguous = true
		} else {
			c.blockers = append(c.blockers, domain.DecisionReason{RuleID: "GATE-OPERATION-001", Effect: "block", Candidate: c.key, Evidence: a.Operation.Sources, Message: "operação local obrigatória não foi comprovada"})
		}
	}
	for _, x := range a.Constraints {
		if catalogLimitedEvidence(x.Evidence) {
			c.ambiguous = true
			continue
		}
		if x.Evidence.Status == domain.EvidenceInconclusive {
			c.ambiguous = true
		}
		if x.Evidence.Status == domain.EvidenceConflict {
			c.blockers = append(c.blockers, domain.DecisionReason{RuleID: "GATE-CONSTRAINT-001", Effect: "block", Candidate: c.key, Evidence: x.Evidence.Sources, Message: "conflito eliminatório: " + x.Constraint})
		}
	}
	return c
}

func finishDecision(r domain.ReuseDecisionReport, capabilities []string, c decisionCandidate, decision domain.ReuseDecision, rule, msg string) domain.ReuseDecisionReport {
	r.Decision = decision
	r.Primary = &domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: c.a.Repository.Owner, Repo: c.a.Repository.Repo}, Commit: c.a.Repository.Commit}
	r.Components = validParts(c.a.ReusableParts)
	r.Gaps = uncoveredGaps(capabilities, c.covered)
	r.Reasons = []domain.DecisionReason{{RuleID: rule, Effect: "select", Candidate: c.key, Evidence: c.evidence, Message: msg}}
	r.Blockers = append(r.Blockers, c.blockers...)
	r.Effort = assessEffort(decision, len(r.Components), len(r.Gaps), false)
	return r
}

func buildDecisionMatrix(caps []string, cs []decisionCandidate) []domain.CapabilityComparison {
	out := make([]domain.CapabilityComparison, 0, len(caps))
	for _, cap := range caps {
		row := domain.CapabilityComparison{Capability: cap}
		for _, c := range cs {
			status := domain.EvidenceInconclusive
			var sources []domain.EvidenceSource
			for _, e := range c.a.Capabilities {
				if normalizeDecisionText(e.Capability) == cap {
					status = e.Status
					sources = e.Sources
					break
				}
			}
			row.Candidates = append(row.Candidates, domain.CandidateCapabilityStatus{Candidate: c.key, Status: status, Evidence: sources})
		}
		out = append(out, row)
	}
	return out
}

func assessEffort(d domain.ReuseDecision, n, gaps int, amb bool) domain.EffortAssessment {
	if amb {
		return domain.EffortAssessment{Method: "heuristic_component_gap_count", Limitations: "Não mede custo nem viabilidade; requer ensaio de integração", Level: domain.EffortUnknown, Factors: []domain.EffortFactor{{Factor: "evidence", Effect: "unknown", Evidence: "coleta ou gate inconclusivo"}}}
	}
	level := domain.EffortLow
	if n > 1 || gaps > 0 {
		level = domain.EffortMedium
	}
	if d == domain.DecisionBuild || gaps > 2 {
		level = domain.EffortHigh
	}
	return domain.EffortAssessment{Method: "heuristic_component_gap_count", Limitations: "Não mede custo nem viabilidade; requer ensaio de integração", Level: level, Factors: []domain.EffortFactor{{Factor: "components", Effect: fmt.Sprintf("%d", n), Evidence: "componentes verificáveis selecionados"}, {Factor: "functional_gaps", Effect: fmt.Sprintf("%d", gaps), Evidence: "matriz de capacidades"}}}
}
func firstCandidate(cs []decisionCandidate, p func(decisionCandidate) bool) *decisionCandidate {
	for i := range cs {
		if p(cs[i]) {
			return &cs[i]
		}
	}
	return nil
}
func uncoveredGaps(caps []string, covered map[string]bool) []string {
	var out []string
	for _, c := range caps {
		if !covered[c] {
			out = append(out, c)
		}
	}
	return out
}
func normalizedSorted(xs []string) []string {
	out := make([]string, 0, len(xs))
	seen := map[string]bool{}
	for _, x := range xs {
		n := normalizeDecisionText(x)
		if n != "" && !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}
func normalizeDecisionText(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}
func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
func containsFold(xs []string, x string) bool {
	for _, v := range xs {
		if strings.EqualFold(strings.TrimSpace(v), strings.TrimSpace(x)) {
			return true
		}
	}
	return false
}

func validParts(parts []domain.ReusableComponent) []domain.ReusableComponent {
	out := make([]domain.ReusableComponent, 0, len(parts))
	for _, part := range parts {
		if part.Name != "" && part.Kind != "" && (part.Path != "" || len(part.Symbols) > 0) && len(part.Capabilities) > 0 && len(part.Evidence) > 0 {
			out = append(out, part)
		}
	}
	return out
}
