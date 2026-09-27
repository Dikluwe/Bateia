/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technology_landscape_test.md
 * @layer L1
 * @updated 2026-09-24
 */
package usecase

import (
	"testing"

	"bateia/01_core/domain"
)

func landscapeClaim(id string, status domain.EvidenceStatus) domain.LandscapeCapability {
	return domain.LandscapeCapability{ID: id, Status: status, Sources: []string{"proof:" + id}}
}

func landscapeAssertion(status domain.EvidenceStatus) domain.LandscapeAssertion {
	return domain.LandscapeAssertion{Status: status, Sources: []string{"proof:operation"}}
}

func landscapeRequest() domain.TechnologyLandscapeRequest {
	return domain.TechnologyLandscapeRequest{SchemaVersion: "1", Current: domain.LandscapeSnapshot{
		AsOf: "2026-09-24", Need: []string{"search", "export"}, Candidates: []domain.LandscapeCandidate{
			{ID: "a", Kind: domain.LandscapeRepository, Name: "A", Capabilities: []domain.LandscapeCapability{landscapeClaim("search", domain.EvidenceProven), landscapeClaim("export", domain.EvidenceProven)}, Standalone: landscapeAssertion(domain.EvidenceProven)},
			{ID: "b", Kind: domain.LandscapeRepository, Name: "B", Capabilities: []domain.LandscapeCapability{landscapeClaim("search", domain.EvidenceProven), landscapeClaim("export", domain.EvidenceAbsent)}, Integration: landscapeAssertion(domain.EvidenceProven)},
			{ID: "c", Kind: domain.LandscapeTechnology, Name: "C", Capabilities: []domain.LandscapeCapability{landscapeClaim("search", domain.EvidenceAbsent), landscapeClaim("export", domain.EvidenceAbsent)}, Analogy: landscapeAssertion(domain.EvidenceProven)},
			{ID: "d", Kind: domain.LandscapeProduct, Name: "D", Capabilities: []domain.LandscapeCapability{landscapeClaim("search", domain.EvidenceProven)}},
		},
	}}
}

func TestTechnologyLandscapeSeparatesRolesFromSimilarity(t *testing.T) {
	report, err := AnalyzeTechnologyLandscape(landscapeRequest())
	if err != nil {
		t.Fatal(err)
	}
	want := []domain.LandscapeRole{domain.LandscapeSubstitute, domain.LandscapeComplement, domain.LandscapeAnalogue, domain.LandscapeUnresolved}
	for i, role := range want {
		if report.Assessments[i].Role != role {
			t.Fatalf("candidate %s = %s, want %s", report.Assessments[i].CandidateID, report.Assessments[i].Role, role)
		}
	}
	if report.Assessments[0].Name != "A" || report.Assessments[0].Standalone.Sources[0] != "proof:operation" || len(report.Assessments[0].Capabilities) != 2 {
		t.Fatalf("relatório perdeu nome ou fontes: %+v", report.Assessments[0])
	}
	if report.Assessments[1].Standalone.Status != domain.EvidenceInconclusive {
		t.Fatalf("asserção omitida não virou INCONCLUSIVE: %+v", report.Assessments[1].Standalone)
	}
	if len(report.Assessments[3].Unknown) != 1 || report.Assessments[3].Unknown[0] != "export" {
		t.Fatalf("ausência de prova virou ausência de capacidade: %+v", report.Assessments[3])
	}
	if report.Counts[domain.LandscapeSubstitute] != 1 || report.Counts[domain.LandscapeUnresolved] != 1 {
		t.Fatalf("contagens incorretas: %+v", report.Counts)
	}
}

func TestTechnologyLandscapeReportsMovementsWithoutMaturityScore(t *testing.T) {
	r := landscapeRequest()
	r.Current.Candidates = r.Current.Candidates[:1]
	r.Current.Candidates[0].Signals = []domain.LandscapeSignal{{Kind: domain.LandscapeReleases, Value: 3, Unit: "count", WindowDays: 90, ObservedAt: "2026-09-24", Sources: []string{"release-feed:2"}}}
	previous := r.Current
	previous.AsOf = "2026-08-24"
	previous.Candidates = append([]domain.LandscapeCandidate(nil), r.Current.Candidates...)
	previous.Candidates[0].Signals = []domain.LandscapeSignal{{Kind: domain.LandscapeReleases, Value: 2, Unit: "count", WindowDays: 90, ObservedAt: "2026-08-24", Sources: []string{"release-feed:1"}}}
	r.Previous = &previous
	report, err := AnalyzeTechnologyLandscape(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Changes) != 1 || report.Changes[0].Kind != "UPDATED" || len(report.Changes[0].Signals) != 1 || report.Changes[0].Signals[0].Movement != "INCREASED" {
		t.Fatalf("delta incorreto: %+v", report.Changes)
	}
	if report.Changes[0].Signals[0].Before == nil || *report.Changes[0].Signals[0].Before != 2 {
		t.Fatalf("valor anterior perdido: %+v", report.Changes[0])
	}
}

func TestTechnologyLandscapeRejectsUngroundedAndIncomparableInput(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*domain.TechnologyLandscapeRequest)
	}{
		{"duplicate need", func(r *domain.TechnologyLandscapeRequest) { r.Current.Need = []string{"search", "search"} }},
		{"claim without source", func(r *domain.TechnologyLandscapeRequest) { r.Current.Candidates[0].Capabilities[0].Sources = nil }},
		{"future signal", func(r *domain.TechnologyLandscapeRequest) {
			r.Current.Candidates[0].Signals = []domain.LandscapeSignal{{Kind: domain.LandscapeArchived, Value: 1, Unit: "boolean", ObservedAt: "2026-09-25", Sources: []string{"repo:1"}}}
		}},
		{"duplicate candidates", func(r *domain.TechnologyLandscapeRequest) {
			r.Current.Candidates = append(r.Current.Candidates, r.Current.Candidates[0])
		}},
		{"changed need", func(r *domain.TechnologyLandscapeRequest) {
			p := r.Current
			p.AsOf = "2026-08-24"
			p.Need = []string{"other"}
			p.Candidates = nil
			r.Previous = &p
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := landscapeRequest()
			tc.mutate(&r)
			if _, err := AnalyzeTechnologyLandscape(r); err == nil {
				t.Fatalf("entrada inválida aceita: %s", tc.name)
			}
		})
	}
}
