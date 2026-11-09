/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/integrity_0171_composition_test.md
 * @layer L1
 * @updated 2026-11-09
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func integrity0171MCDA(r domain.ReuseDecisionRequest) domain.MulticriteriaDecisionRequest {
	return domain.MulticriteriaDecisionRequest{Idea: r.Idea, Assessment: domain.AssessmentReport{SchemaVersion: "1", Assessments: r.Assessments}, Policy: r.Policy, Composition: r.Composition, Criteria: []domain.MulticriteriaCriterion{{ID: "coverage", Direction: domain.CriterionMaximize, ScaleMin: 0, ScaleMax: 1}}}
}

func TestIntegrity0171CompositionRejectsDivergentMemberDeclarations(t *testing.T) {
	for name, mutate := range map[string]func(*domain.ReusableComponent){
		"export-only":         func(p *domain.ReusableComponent) { p.Capabilities = []string{"export"} },
		"unlinked-capability": func(p *domain.ReusableComponent) { p.Capabilities = append(p.Capabilities, "unlinked") },
		"capability-order": func(p *domain.ReusableComponent) {
			p.Capabilities[0], p.Capabilities[1] = p.Capabilities[1], p.Capabilities[0]
		},
		"evidence-order": func(p *domain.ReusableComponent) { p.Evidence[0], p.Evidence[1] = p.Evidence[1], p.Evidence[0] },
		"extra-evidence": func(p *domain.ReusableComponent) { p.Evidence = append(p.Evidence, decisionSource("b", "additional")) },
		"caveat":         func(p *domain.ReusableComponent) { p.Caveats = []string{"requires adaptation"} },
	} {
		for _, reverse := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/forward", true: "/reverse"}[reverse], func(t *testing.T) {
				req, a, b := compositionFixture()
				req.Composition = verifiedComposition(a, b)
				original := integrity0170Clone(req.Assessments[1].ReusableParts[0])
				other := integrity0170Clone(original)
				mutate(&other)
				req.Assessments[1].ReusableParts = []domain.ReusableComponent{original, other}
				if reverse {
					req.Assessments[1].ReusableParts = []domain.ReusableComponent{other, original}
				}
				before, _ := json.Marshal(req)
				got, err := DecideReuse(req)
				if err == nil || !reflect.DeepEqual(got, domain.ReuseDecisionReport{}) {
					t.Errorf("divergent member yielded report: err=%v decision=%s", err, got.Decision)
				}
				mcda, err := AnalyzeMulticriteriaDecision(integrity0171MCDA(req))
				if err == nil || !reflect.DeepEqual(mcda, domain.MulticriteriaDecisionReport{}) {
					t.Errorf("MCDA admitted divergent member: err=%v frontier=%v", err, mcda.Frontier)
				}
				after, _ := json.Marshal(req)
				if string(before) != string(after) {
					t.Fatal("input changed")
				}
			})
		}
	}
}

func TestIntegrity0171CompositionIdenticalMembersAreIdempotent(t *testing.T) {
	req, a, b := compositionFixture()
	req.Composition = verifiedComposition(a, b)
	control, err := DecideReuse(req)
	if err != nil || control.Decision != domain.DecisionCompose {
		t.Fatal(err, control)
	}
	mc, err := AnalyzeMulticriteriaDecision(integrity0171MCDA(req))
	if err != nil {
		t.Fatal(err)
	}
	for i := range req.Assessments {
		req.Assessments[i].ReusableParts = append(req.Assessments[i].ReusableParts, integrity0170Clone(req.Assessments[i].ReusableParts[0]))
	}
	got, err := DecideReuse(req)
	if err != nil || !reflect.DeepEqual(control, got) {
		t.Fatal("identical members changed composition", err, got)
	}
	gotMC, err := AnalyzeMulticriteriaDecision(integrity0171MCDA(req))
	if err != nil || !reflect.DeepEqual(mc, gotMC) {
		t.Fatal("identical members changed MCDA", err, gotMC)
	}
}

func TestIntegrity0171CompositionInadmissiblePartDoesNotConflict(t *testing.T) {
	req, a, b := compositionFixture()
	req.Composition = verifiedComposition(a, b)
	control, err := DecideReuse(req)
	if err != nil {
		t.Fatal(err)
	}
	other := integrity0170Clone(req.Assessments[1].ReusableParts[0])
	other.Evidence = []domain.EvidenceSource{decisionSource("outsider", "unrelated")}
	other.Caveats = []string{"inadmissible"}
	req.Assessments[1].ReusableParts = append(req.Assessments[1].ReusableParts, other)
	got, err := DecideReuse(req)
	if err != nil || !reflect.DeepEqual(control, got) {
		t.Fatal("inadmissible part changed composition", err, got)
	}
}
