/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_site_plan_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestInsightSitePlanCanonicalSelectionAndRejections(t *testing.T) {
	g, a := domain.InsightSiteBinding{Connector: &insightSiteFixtureConnector{id: domain.ConnectorGiteeCatalog}}, domain.InsightSiteBinding{Connector: &insightSiteFixtureConnector{id: domain.ConnectorAtomGitCatalog}}
	r := insightSiteRequest()
	r.Sources = []string{"gitee", "atomgit"}
	before, _ := json.Marshal(r)
	first, err := PlanInsightSites(r, []domain.InsightSiteBinding{g, a})
	if err != nil {
		t.Fatal(err)
	}
	r.Sources = []string{domain.ConnectorAtomGitCatalog, domain.ConnectorGiteeCatalog}
	again, err := PlanInsightSites(r, []domain.InsightSiteBinding{a, g})
	if err != nil || !reflect.DeepEqual(first, again) {
		t.Fatal("order or aliases changed plan")
	}
	var original domain.InsightSiteRequest
	_ = json.Unmarshal(before, &original)
	if original.Sources[0] != "gitee" {
		t.Fatal("input alias mutated")
	}
	for _, sources := range [][]string{{}, {"unknown"}, {"gitee", "gitee"}, {"gitee", domain.ConnectorGiteeCatalog}, {"all", "gitee"}, {"bitbucket"}, {""}} {
		r := insightSiteRequest()
		r.Sources = sources
		if _, err := PlanInsightSites(r, []domain.InsightSiteBinding{g, a}); err == nil {
			t.Fatal("invalid selection admitted", sources)
		}
	}
	g.Reader = &insightTestReader{}
	if _, err := PlanInsightSites(insightSiteRequest(), []domain.InsightSiteBinding{g}); err == nil {
		t.Fatal("cross-host reader admitted")
	}
	for _, budget := range []int{-1, 61} {
		r := insightSiteRequest()
		r.MaxCalls = budget
		if _, err := PlanInsightSites(r, []domain.InsightSiteBinding{a}); err == nil {
			t.Fatal("invalid budget admitted")
		}
	}
	r = insightSiteRequest()
	r.Locale = "en_US"
	if _, err := PlanInsightSites(r, []domain.InsightSiteBinding{a}); err == nil {
		t.Fatal("invalid locale admitted")
	}
}

func TestInsightSiteGenerationLimitPrecedesEffects(t *testing.T) {
	for _, generation := range []int{999998, 999999} {
		for _, followLinks := range []int{0, 1} {
			t.Run(fmt.Sprintf("%d/links%d", generation, followLinks), func(t *testing.T) {
				memory := memoryFixture(t)
				memory.Generation = generation
				memory = sealInsightMemory(memory)
				request := insightSiteRequest()
				request.Memory, request.FollowLinks = &memory, &followLinks
				before, _ := json.Marshal(request)
				_, admissionErr := ValidateInsightSiteOptions(request)
				reader, bindings := followFixture()
				catalogCalls := 0
				bindings[0].Connector = &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog, hook: func() { catalogCalls++ }}
				report, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
				after, _ := json.Marshal(request)
				if string(before) != string(after) {
					t.Fatal("discovery changed caller history")
				}
				if generation == 999999 {
					if !errors.Is(admissionErr, domain.ErrInvalidRunInput) || !errors.Is(err, domain.ErrInvalidRunInput) || report.ID != "" || catalogCalls != 0 || len(reader.reads) != 0 || len(reader.resolves) != 0 {
						t.Fatal("exhausted generation reached discovery effects", admissionErr, err, catalogCalls)
					}
					return
				}
				if admissionErr != nil || err != nil || report.ExplorationMemory == nil || report.ExplorationMemory.Generation != 999999 || catalogCalls == 0 {
					t.Fatal("last valid discovery round failed", admissionErr, err)
				}
				if _, err := ValidateInsightExplorationMemory(report.ExplorationMemory); err != nil {
					t.Fatal("last discovery memory is unusable", err)
				}
			})
		}
	}
}
