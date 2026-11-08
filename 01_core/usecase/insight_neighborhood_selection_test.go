/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_selection_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func aliasNeighborhoodFixture(t *testing.T) (domain.InsightSiteReport, *memoryFollowReader, []domain.InsightSiteBinding) {
	t.Helper()
	r, bindings := followFixture()
	r.contents["org/repo"] = "# Root\n[Alpha](https://github.com/branch/alpha)\n[Clone URL](https://github.com/branch/alpha.git)\n[Beta](https://github.com/branch/beta)\n"
	r.contents["branch/alpha"] = "# Alpha\nA leaf without further links.\n"
	request, zero := insightSiteRequest(), 0
	request.FollowLinks = &zero
	first, err := DiscoverInsightSites(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings)
	if err != nil {
		t.Fatal(err)
	}
	r.reads, r.resolves = nil, nil
	bindings[0].Connector = &insightSiteFixtureConnector{id: domain.ConnectorGitHubCatalog, hook: func() { t.Fatal("expand queried a catalog") }}
	return first, r, bindings
}

func neighborFor(t *testing.T, report domain.InsightSiteReport, name string) domain.InsightNeighborhoodNeighbor {
	t.Helper()
	for _, n := range report.Neighborhood.Neighbors {
		if n.Identity.ExternalID == name {
			return n
		}
	}
	t.Fatal("missing neighbor", name)
	return domain.InsightNeighborhoodNeighbor{}
}

func TestNeighborhoodKnownAliasDoesNotConsumeBudgetOrRewriteHistory(t *testing.T) {
	first, reader, bindings := aliasNeighborhoodFixture(t)
	before, _ := json.Marshal(first)
	one := 1
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	second, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Focus: "https://github.com/branch/alpha.git", FollowLinks: &one}, bindings)
	if err != nil || !reflect.DeepEqual(reader.resolves, []string{"branch/alpha"}) || second.Neighborhood.Focus.ExternalID != "branch/alpha" {
		t.Fatal("alias focus did not read only known canonical target", err, reader.resolves)
	}
	alias := neighborFor(t, second, "branch/alpha.git")
	canonical := neighborFor(t, second, "branch/alpha")
	if alias.SelectionStatus != "DUPLICATE_ALIAS" || alias.RelatedEntryID != canonical.EntryID || *alias.Distance != 0 || canonical.SelectionStatus != "FRONTIER" || len(second.Neighborhood.FrontierEntryIDs) != 1 {
		t.Fatal("alias suppression not explained", alias, canonical)
	}
	var previous domain.InsightMemoryEntry
	for _, e := range first.ExplorationMemory.Entries {
		if e.Identity.ExternalID == "branch/alpha.git" {
			previous = e
		}
	}
	for _, e := range second.ExplorationMemory.Entries {
		if e.ID == previous.ID && !reflect.DeepEqual(e, previous) {
			t.Fatal("suppressed alias state/counters changed")
		}
	}
	third, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: second}, bindings)
	if err != nil || third.Status != "FRONTIER_EXHAUSTED" || third.Neighborhood.Focus.ExternalID != "branch/alpha" || len(reader.resolves) != 1 || neighborFor(t, third, "branch/alpha.git").SelectionStatus != "ALIAS_ALREADY_EXAMINED" {
		t.Fatal("inherited focus leaked or reread alias", err, reader.resolves)
	}
	after, _ := json.Marshal(first)
	if string(before) != string(after) {
		t.Fatal("historical input mutated")
	}
}

func TestNeighborhoodAliasDiscardAndConflictingDecisions(t *testing.T) {
	first, reader, bindings := aliasNeighborhoodFixture(t)
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	one := 1
	closed, err := UpdateInsightExplorationMemory(*first.ExplorationMemory, "https://github.com/branch/alpha", "discarded", "Not useful for this investigation")
	if err != nil {
		t.Fatal(err)
	}
	result, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Memory: &closed, FollowLinks: &one}, bindings)
	if err != nil || !reflect.DeepEqual(reader.resolves, []string{"branch/beta"}) || neighborFor(t, result, "branch/alpha.git").SelectionStatus != "ALIAS_DISCARDED" {
		t.Fatal("discard bypassed through alias", err, reader.resolves)
	}
	if _, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: first, Memory: &closed, Focus: "https://github.com/branch/alpha.git"}); err == nil {
		t.Fatal("alias focus bypassed discarded canonical target")
	}
	conflict, err := UpdateInsightExplorationMemory(*first.ExplorationMemory, "https://github.com/branch/alpha.git", "discarded", "Explicit alias decision")
	if err != nil {
		t.Fatal(err)
	}
	reader.resolves = nil
	result, err = ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Memory: &conflict, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil || result.Status != "PARTIAL" || len(reader.resolves) != 0 || neighborFor(t, result, "branch/alpha").SelectionStatus != "ALIAS_STATE_CONFLICT" {
		t.Fatal("conflicting history silently resolved", err, result.Status)
	}
}

func TestNeighborhoodAliasWithPastExaminationRequiresExplicitCanonicalReopen(t *testing.T) {
	first, reader, bindings := aliasNeighborhoodFixture(t)
	memory, err := ValidateInsightExplorationMemory(first.ExplorationMemory)
	if err != nil {
		t.Fatal(err)
	}
	identity, _ := insightLeadTarget("https://github.com/branch/alpha.git")
	alias := insightMemoryFind(&memory, identity)
	// Declared historical fixture, not public repository evidence.
	alias.State, alias.ReadAttempts, alias.LastStatus = "EXAMINED", 1, "COLLECTED"
	alias.LastRead = &domain.InsightReadObservation{Source: domain.CommitRef{Owner: "branch", Repo: "alpha.git", Commit: strings.Repeat("a", 40)}, Path: "README.md", ContentHash: "sha256:" + contentHash("fixture")}
	memory = sealInsightMemory(memory)
	one := 1
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	result, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Memory: &memory, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil || len(reader.resolves) != 0 || neighborFor(t, result, "branch/alpha").SelectionStatus != "ALIAS_ALREADY_EXAMINED" {
		t.Fatal("prior alias examination caused duplicate read", err)
	}
	memory, err = UpdateInsightExplorationMemory(memory, "https://github.com/branch/alpha", "pending", "Explicitly examine the canonical repository again")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Memory: &memory, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil || !reflect.DeepEqual(reader.resolves, []string{"branch/alpha"}) {
		t.Fatal("explicit canonical reopening ignored", err)
	}
}

func TestNeighborhoodAliasRecognitionIsHostAndForgeScoped(t *testing.T) {
	for _, fixture := range []struct {
		base, alias string
		matches     bool
	}{
		{"https://github.com/org/repo", "https://github.com/org/repo.git", true},
		{"https://gitee.com/org/repo", "https://gitee.com/org/repo.git", true},
		{"https://gitlab.com/org/sub/repo", "https://gitlab.com/org/sub/repo.git", true},
		{"https://codeberg.org/org/repo", "https://codeberg.org/org/repo.git", true},
		{"https://atomgit.com/org/repo", "https://atomgit.com/org/repo.git", true},
		{"https://atomgit.com/org/repo", "https://gitcode.com/org/repo.git", false},
		{"https://github.com/org/repo", "https://gitee.com/org/repo.git", false},
		{"https://github.com/org/repo", "https://github.com/org/Repo.git", false},
		{"https://github.com/org/repo", "https://github.com/org/repo.gitlab", false},
		{"https://github.com/org/repo.git", "https://github.com/org/repo.git.git", false},
	} {
		base, _ := insightLeadTarget(fixture.base)
		alias, _ := insightLeadTarget(fixture.alias)
		memory := domain.InsightExplorationMemory{Entries: []domain.InsightMemoryEntry{{Identity: base}, {Identity: alias}}}
		actual := insightNeighborhoodAliasTarget(&memory, alias)
		if (actual == base) != fixture.matches {
			t.Fatal("alias crossed identity boundary", fixture)
		}
		memory.Entries = memory.Entries[1:]
		if insightNeighborhoodAliasTarget(&memory, alias) != alias {
			t.Fatal("unknown base was invented")
		}
	}
}

func TestNeighborhoodDistancePrecedesRetriesAndGeneralModeOmitsDistance(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	one, zero := 1, 0
	ctx := domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{})
	second, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: first, Focus: "https://github.com/branch/alpha", FollowLinks: &one}, bindings)
	if err != nil {
		t.Fatal(err)
	}
	third, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: second, FollowLinks: &zero}, bindings)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := UpdateInsightExplorationMemory(*third.ExplorationMemory, "https://github.com/org/repo", "pending", "Revisit the focus first")
	if err != nil {
		t.Fatal(err)
	}
	memory, err = UpdateInsightExplorationMemory(memory, "https://github.com/branch/alpha", "pending", "Revisit this immediate neighbor")
	if err != nil {
		t.Fatal(err)
	}
	for i := range memory.Entries {
		switch memory.Entries[i].Identity.ExternalID {
		case "org/repo":
			memory.Entries[i].ReadAttempts = 9
		case "branch/alpha":
			memory.Entries[i].ReadAttempts = 6
		}
	}
	memory = sealInsightMemory(memory)
	reader.resolves = nil
	result, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: third, Memory: &memory, Focus: "https://github.com/org/repo", FollowLinks: &one}, bindings)
	if err != nil || !reflect.DeepEqual(reader.resolves, []string{"org/repo"}) {
		t.Fatal("focus distance lost to retry count", err, reader.resolves)
	}
	expected := []string{"org/repo", "branch/beta", "branch/alpha", "deeper/one"}
	actual := []string{}
	for _, n := range result.Neighborhood.Neighbors {
		actual = append(actual, n.Identity.ExternalID)
	}
	if !reflect.DeepEqual(actual, expected) || *neighborFor(t, result, "deeper/one").Distance != 2 || neighborFor(t, result, "org/repo").ReadAttempts != 9 {
		t.Fatal("selection snapshot/order wrong", actual)
	}
	general, err := ExpandInsightNeighborhood(ctx, domain.InsightNeighborhoodRequest{Input: third, AllTrails: true, FollowLinks: &zero}, bindings)
	if err != nil || general.Neighborhood.Focus != nil {
		t.Fatal("general mode retained inherited focus", err)
	}
	for _, n := range general.Neighborhood.Neighbors {
		if n.Distance != nil {
			t.Fatal("distance invented without focus")
		}
	}
	if _, err := ValidateInsightNeighborhood(domain.InsightNeighborhoodRequest{Input: third, AllTrails: true, Focus: "https://github.com/org/repo"}); err == nil {
		t.Fatal("conflicting scope flags accepted")
	}
}
