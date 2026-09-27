/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/reuse_decision_product_test.md
 * @layer L1
 * @updated 2026-09-24
 */

package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func productCandidate(repo string, caps []string, proven map[string]bool, parts []domain.ReusableComponent) domain.CandidateAssessment {
	src := decisionSource(repo, "product-review")
	a := domain.CandidateAssessment{Repository: src.Ref, ProblemAlignment: domain.CapabilityEvidence{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{src}}, License: domain.LicenseAssessment{Status: domain.EvidenceProven, Identifier: "MIT", Sources: []domain.EvidenceSource{decisionSource(repo, "license")}}, Maintenance: domain.MaintenanceAssessment{Status: domain.EvidenceProven, Sources: []domain.EvidenceSource{decisionSource(repo, "maintenance")}}, Operation: domain.OperationAssessment{LocalStatus: domain.EvidenceProven, Sources: []domain.EvidenceSource{src}}, ReusableParts: parts}
	for _, cap := range caps {
		status := domain.EvidenceAbsent
		if proven[cap] {
			status = domain.EvidenceProven
		}
		a.Capabilities = append(a.Capabilities, domain.CapabilityEvidence{Capability: cap, Status: status, Sources: []domain.EvidenceSource{decisionSource(repo, "product-"+cap)}})
	}
	return a
}
func productPart(repo, name string, caps ...string) domain.ReusableComponent {
	sources := []domain.EvidenceSource{}
	for _, capability := range caps {
		sources = append(sources, decisionSource(repo, "product-"+capability))
	}
	return domain.ReusableComponent{Name: name, Kind: "package", Path: name, Capabilities: caps, Evidence: sources}
}

func TestReuseDecisionProductReplayFiveDistinctIdeas(t *testing.T) {
	type scenario struct {
		name        string
		idea        domain.IdeaSpec
		assessments []domain.CandidateAssessment
		composition *domain.CompositionInput
		want        domain.ReuseDecision
	}
	all := func(xs ...string) map[string]bool {
		m := map[string]bool{}
		for _, x := range xs {
			m[x] = true
		}
		return m
	}
	scenarios := []scenario{
		{name: "busca local", idea: domain.IdeaSpec{Problem: "buscar projetos locais", Capabilities: []string{"indexar", "ranquear"}}, want: domain.DecisionFork},
		{name: "analisador estrutural", idea: domain.IdeaSpec{Problem: "analisar impacto", Capabilities: []string{"parsear", "impacto", "sarif"}}, want: domain.DecisionExtend},
		{name: "pipeline de eventos", idea: domain.IdeaSpec{Problem: "processar eventos", Capabilities: []string{"ingerir", "transformar", "exportar"}}, want: domain.DecisionCompose},
		{name: "editor de receitas", idea: domain.IdeaSpec{Problem: "editar receitas", Capabilities: []string{"parsear", "editar", "sincronizar"}}, want: domain.DecisionIntegrate},
		{name: "simulador inédito", idea: domain.IdeaSpec{Problem: "simular domínio inédito", Capabilities: []string{"simular", "validar"}}, want: domain.DecisionBuild},
	}
	scenarios[0].assessments = []domain.CandidateAssessment{productCandidate("local-search", scenarios[0].idea.Capabilities, all("indexar", "ranquear"), nil)}
	scenarios[1].assessments = []domain.CandidateAssessment{productCandidate("analyzer", scenarios[1].idea.Capabilities, all("parsear", "impacto"), []domain.ReusableComponent{productPart("analyzer", "engine", "parsear", "impacto")})}
	a := productCandidate("ingest", scenarios[2].idea.Capabilities, all("ingerir"), []domain.ReusableComponent{productPart("ingest", "source", "ingerir")})
	a.ProblemAlignment.Status = domain.EvidencePartial
	b := productCandidate("flow", scenarios[2].idea.Capabilities, all("transformar", "exportar"), []domain.ReusableComponent{productPart("flow", "pipeline", "transformar", "exportar")})
	b.ProblemAlignment.Status = domain.EvidencePartial
	scenarios[2].assessments = []domain.CandidateAssessment{a, b}
	scenarios[2].composition = verifiedComposition(compositionMemberID(a.Repository, a.ReusableParts[0]), compositionMemberID(b.Repository, b.ReusableParts[0]), "ingest")
	c := productCandidate("recipe-parser", scenarios[3].idea.Capabilities, all("parsear"), []domain.ReusableComponent{productPart("recipe-parser", "parser", "parsear")})
	c.ProblemAlignment.Status = domain.EvidencePartial
	scenarios[3].assessments = []domain.CandidateAssessment{c}
	d := productCandidate("unrelated", scenarios[4].idea.Capabilities, map[string]bool{}, nil)
	d.ProblemAlignment.Status = domain.EvidenceAbsent
	scenarios[4].assessments = []domain.CandidateAssessment{d}
	useful := 0
	for _, s := range scenarios {
		r, err := DecideReuse(domain.ReuseDecisionRequest{Idea: s.idea, Assessments: s.assessments, Policy: decisionPolicy(), Composition: s.composition})
		if err != nil {
			t.Fatalf("%s: %v", s.name, err)
		}
		if r.Decision == s.want {
			useful++
		} else {
			t.Errorf("%s: got %s want %s", s.name, r.Decision, s.want)
		}
	}
	if float64(useful)/float64(len(scenarios)) < .8 {
		t.Fatalf("useful decisions %d/%d", useful, len(scenarios))
	}
}
