/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/insight_neighborhood_sources_test.md
 * @layer L1
 * @updated 2026-10-29
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestInsightNeighborhoodRejectsTamperedHistoryBeforeEffects(t *testing.T) {
	first, reader, bindings := neighborhoodFixture(t)
	for name, change := range map[string]func(*domain.InsightNeighborhoodRequest){
		"report seal": func(r *domain.InsightNeighborhoodRequest) { r.Input.Status = "forged" },
		"content": func(r *domain.InsightNeighborhoodRequest) {
			r.Input.Bundles[0].Items[0].Content += "forged"
			r.Input.ID = ""
			r.Input.ID = "site_discovery_" + insightHash(r.Input)
		},
		"fingerprint": func(r *domain.InsightNeighborhoodRequest) {
			i := &r.Input.Bundles[0].Items[0]
			i.Content = "# Forged\n[Target](https://github.com/evil/path)\n"
			i.ContentHash = "sha256:" + contentHash(i.Content)
			r.Input.Bundles[0].BytesRead = int64(len(i.Content))
			r.Input.Readmes[0].Receipt.Bytes = len(i.Content)
			r.Input.ID = ""
			r.Input.ID = "site_discovery_" + insightHash(r.Input)
		},
		"unknown focus": func(r *domain.InsightNeighborhoodRequest) { r.Focus = "https://github.com/unknown/path" },
		"bad limit":     func(r *domain.InsightNeighborhoodRequest) { n := 5; r.FollowLinks = &n },
		"bad passages":  func(r *domain.InsightNeighborhoodRequest) { r.MaxPassages = 17 },
	} {
		t.Run(name, func(t *testing.T) {
			raw, _ := json.Marshal(first)
			var input domain.InsightSiteReport
			_ = json.Unmarshal(raw, &input)
			request := domain.InsightNeighborhoodRequest{Input: input}
			change(&request)
			if _, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), request, bindings); err == nil || len(reader.resolves) != 0 {
				t.Fatal("invalid request reached a reader", err)
			}
		})
	}
}

func frozenNeighborhoodOriginsFixture(t *testing.T, kind string, sizes []int) domain.InsightSiteReport {
	t.Helper()
	report := domain.InsightSiteReport{SchemaVersion: "insight-site-discovery-v1", Plan: domain.InsightSitePlan{Request: insightSiteRequest()}}
	memory := domain.InsightExplorationMemory{SchemaVersion: "insight-exploration-memory-v1", Generation: 1, Entries: []domain.InsightMemoryEntry{}}
	for i, size := range sizes {
		connector, host, owner := domain.ConnectorGitHubCatalog, "github.com", "org"
		if kind == "site" {
			connector, host, owner = domain.ConnectorGitLabCatalog, "gitlab.com", "org/subgroup"
		} else if kind == "followup" {
			connector, host = domain.ConnectorGiteeCatalog, "gitee.com"
		}
		repo := fmt.Sprintf("repo%d", i)
		identity := domain.ProductIdentity{ConnectorID: connector, ExternalID: owner + "/" + repo, CanonicalURL: "https://" + host + "/" + owner + "/" + repo}
		ref := domain.CommitRef{Owner: owner, Repo: repo, Commit: strings.Repeat("a", 40)}
		content := strings.Repeat("x", size)
		item := domain.CollectedEvidence{Kind: domain.EvidenceDocumentation, Ref: ref, Path: "README.md", Content: content, ContentHash: "sha256:" + contentHash(content)}
		bundle := domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: []domain.CollectedEvidence{item}, BytesRead: int64(size), Complete: true}
		entryID := insightMemoryEntryID(identity)
		observation := domain.InsightReadObservation{Source: ref, Path: item.Path, ContentHash: item.ContentHash}
		memory.Entries = append(memory.Entries, domain.InsightMemoryEntry{ID: entryID, Identity: identity, State: "EXAMINED", FirstGeneration: 1, LastQuestion: "queue", ReadAttempts: 1, LastStatus: "COLLECTED", LastRead: &observation})
		id := fmt.Sprintf("candidate%d", i)
		if kind == "followup" {
			id = entryID
		}
		receipt := domain.InsightReadmeReceipt{ReferenceID: id, Source: ref, Status: "COLLECTED", Path: item.Path, Bytes: size, Scope: "ROOT_README_ONLY"}
		switch kind {
		case "github":
			report.Candidates = append(report.Candidates, domain.InsightSiteCandidate{ID: id, Selected: true, Product: domain.ProductCandidate{Identity: identity}})
			report.Bundles = append(report.Bundles, bundle)
		case "site":
			report.SiteContexts = append(report.SiteContexts, domain.InsightSiteContext{CandidateID: id, Identity: identity, Bundle: bundle})
		case "followup":
			report.Followups = append(report.Followups, domain.InsightFollowup{EntryID: id, Identity: identity, Status: "COLLECTED", Receipt: &receipt, Bundle: &bundle})
		default:
			t.Fatal("unknown fixture kind", kind)
		}
		if kind != "followup" {
			report.Readmes = append(report.Readmes, domain.InsightSiteReadme{CandidateID: id, ConnectorID: connector, Status: "COLLECTED", Receipt: &receipt})
		}
	}
	memory = sealInsightMemory(memory)
	if _, err := ValidateInsightExplorationMemory(&memory); err != nil {
		t.Fatal("invalid fixture history", err)
	}
	report.ExplorationMemory = &memory
	return report
}

func TestNeighborhoodFrozenOriginExactLimits(t *testing.T) {
	cases := []struct {
		name  string
		sizes []int
		valid bool
	}{
		{"one file at limit", []int{32 << 10}, true},
		{"one file over limit", []int{(32 << 10) + 1}, false},
		{"eight documents", []int{1, 1, 1, 1, 1, 1, 1, 1}, true},
		{"nine documents", []int{1, 1, 1, 1, 1, 1, 1, 1, 1}, false},
		{"aggregate at limit", []int{32 << 10, 32 << 10, 32 << 10, 32 << 10}, true},
		{"aggregate over limit", []int{32 << 10, 32 << 10, 32 << 10, 32 << 10, 1}, false},
	}
	for _, kind := range []string{"github", "site", "followup"} {
		for _, tc := range cases {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				input := frozenNeighborhoodOriginsFixture(t, kind, tc.sizes)
				before, _ := json.Marshal(input)
				origins, err := insightNeighborhoodOrigins(input)
				if tc.valid {
					if err != nil || len(origins) != len(tc.sizes) {
						t.Fatal("inclusive source limit rejected", err, len(origins))
					}
				} else if err == nil {
					t.Fatal("excess frozen sources admitted")
				}
				after, _ := json.Marshal(input)
				if string(before) != string(after) {
					t.Fatal("source validation changed frozen input")
				}
				if !tc.valid {
					assertFrozenOriginRejectedBeforeRead(t, input)
				}
			})
		}
	}
}

func assertFrozenOriginRejectedBeforeRead(t *testing.T, input domain.InsightSiteReport) {
	t.Helper()
	// Keep an eligible pending destination so late validation would expose a read.
	prepared := *input.ExplorationMemory
	prepared.Entries = append([]domain.InsightMemoryEntry{}, prepared.Entries...)
	prepared = sealInsightMemory(prepared)
	memory, err := ValidateInsightExplorationMemory(&prepared)
	if err != nil {
		t.Fatal("source control must have an independently valid memory", err)
	}
	identity := domain.ProductIdentity{ConnectorID: domain.ConnectorGitHubCatalog, ExternalID: "pending/target", CanonicalURL: "https://github.com/pending/target"}
	memory.Entries = append(memory.Entries, domain.InsightMemoryEntry{ID: insightMemoryEntryID(identity), Identity: identity, State: "PENDING", FirstGeneration: 1, LastQuestion: "queue"})
	memory = sealInsightMemory(memory)
	input.ExplorationMemory = &memory
	input.ID = ""
	input.ID = "site_discovery_" + insightHash(input)
	reader, bindings := followFixture()
	one := 1
	result, err := ExpandInsightNeighborhood(domain.WithExecutionRuntime(t.Context(), &insightTestRuntime{}), domain.InsightNeighborhoodRequest{Input: input, FollowLinks: &one}, bindings)
	if err == nil || result.ID != "" || len(reader.reads) != 0 || len(reader.resolves) != 0 {
		t.Fatal("bad frozen input was published or reached a reader", err)
	}
}

func TestNeighborhoodFrozenOriginRejectsInconsistentWitnesses(t *testing.T) {
	for _, kind := range []string{"github", "site", "followup"} {
		for _, mode := range []string{"receipt source", "receipt path", "receipt bytes", "receipt scope", "receipt reference", "receipt absent", "last read", "duplicate"} {
			t.Run(kind+"/"+mode, func(t *testing.T) {
				input := frozenNeighborhoodOriginsFixture(t, kind, []int{20})
				if _, err := insightNeighborhoodOrigins(input); err != nil {
					t.Fatal("positive source control failed", err)
				}
				var r *domain.InsightReadmeReceipt
				if kind == "followup" {
					r = input.Followups[0].Receipt
				} else {
					r = input.Readmes[0].Receipt
				}
				switch mode {
				case "receipt source":
					r.Source.Commit = strings.Repeat("b", 40)
				case "receipt path":
					r.Path = "README.txt"
				case "receipt bytes":
					r.Bytes++
				case "receipt scope":
					r.Scope = "FULL_REPOSITORY"
				case "receipt reference":
					r.ReferenceID = "another-entry"
				case "receipt absent":
					if kind == "followup" {
						input.Followups[0].Receipt = nil
					} else {
						input.Readmes[0].Receipt = nil
					}
				case "last read":
					input.ExplorationMemory.Entries[0].LastRead.ContentHash = "sha256:" + strings.Repeat("b", 64)
				case "duplicate":
					if kind == "github" {
						input.Bundles = append(input.Bundles, input.Bundles[0])
					} else if kind == "site" {
						input.SiteContexts = append(input.SiteContexts, input.SiteContexts[0])
					} else {
						input.Followups = append(input.Followups, input.Followups[0])
					}
				}
				if _, err := insightNeighborhoodOrigins(input); err == nil {
					t.Fatal("inconsistent witness admitted")
				}
				assertFrozenOriginRejectedBeforeRead(t, input)
			})
		}
	}
}
