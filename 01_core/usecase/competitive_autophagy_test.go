/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/competitive_autophagy_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"reflect"
	"testing"
)

func autophagyBundle(owner, repo, commit, symbol string, complete bool, license string) domain.EvidenceBundle {
	ref := domain.CommitRef{Owner: owner, Repo: repo, Commit: commit}
	items := []domain.CollectedEvidence{}
	add := func(kind domain.EvidenceKind, path, content string) {
		items = append(items, domain.CollectedEvidence{Kind: kind, Ref: ref, Path: path, Content: content, ContentHash: contentHash(content)})
	}
	if symbol != "" {
		add(domain.EvidenceFile, "feature.go", "type "+symbol+" struct{}")
		add(domain.EvidenceTest, "feature_test.go", "func Test"+symbol+"() {}")
	}
	if license != "" {
		add(domain.EvidenceLicense, "LICENSE", license)
	}
	return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: items, Complete: complete, TreeComplete: complete}
}

func autophagyRequest(internal, competitors []domain.EvidenceBundle) domain.CompetitiveAutophagyRequest {
	return domain.CompetitiveAutophagyRequest{
		SchemaVersion: "1",
		Idea:          domain.IdeaSpec{Problem: "analisar código", Capabilities: []string{"dependency graph", "blast radius"}},
		Target:        domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: "local", Repo: "bateia"}, Commit: "target1"},
		Internal:      internal, Competitors: competitors, Rules: DefaultEvidenceRules(),
		ReusePolicy:    domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, AllowArchived: true},
		KatamariPolicy: domain.KatamariPolicy{SchemaVersion: "1", MaxProjects: 3, MaxComponents: 3, AllowedLicenses: []string{"MIT"}, RequireImmutableRef: true},
	}
}

func TestAutophagyKeepsLexicalMatchesPartial(t *testing.T) {
	internal := autophagyBundle("local", "bateia", "target1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	external := autophagyBundle("oss", "impact", "c1", "BlastRadius", true, "MIT License\nPermission is hereby granted, free of charge")
	result, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{internal}, []domain.EvidenceBundle{external}))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result.Autophagy.Partial, []string{"dependency graph"}) || !reflect.DeepEqual(result.Autophagy.Gaps, []string{"blast radius"}) {
		t.Fatalf("unexpected autofagia: %#v", result.Autophagy)
	}
	if len(result.Katamari.Components) != 0 || len(result.Katamari.CoveredGaps) != 0 {
		t.Fatalf("unexpected Katamari: %#v", result.Katamari)
	}

}

func TestIncompleteInternalEvidenceIsUnknownNotGap(t *testing.T) {
	internal := autophagyBundle("local", "bateia", "target1", "", false, "MIT License\nPermission is hereby granted, free of charge")
	result, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{internal}, nil))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Autophagy.Gaps) != 0 || len(result.Autophagy.Unknown) != 2 {
		t.Fatalf("opacity became gap: %#v", result.Autophagy)
	}
	if result.Katamari.Feasible {
		t.Fatal("Katamari must not absorb dependencies for unknown internal state")
	}
}

func TestDuplicateLexicalDeclarationsDoNotRequireConsolidation(t *testing.T) {
	a := autophagyBundle("local", "one", "a1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	b := autophagyBundle("local", "two", "b1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	r, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{a, b}, nil))
	if err != nil {
		t.Fatal(err)
	}
	if r.Autophagy.Findings[1].Capability == "dependency graph" {
		// Findings are sorted alphabetically; locate without depending on index.
	}
	found := false
	for _, f := range r.Autophagy.Findings {
		found = found || (f.Capability == "dependency graph" && f.Status == domain.EvidencePartial)
	}
	if !found {
		t.Fatalf("duplicate not detected: %#v", r.Autophagy.Findings)
	}
}

func TestKatamariRejectsUnverifiedComponentsAndForbiddenLicense(t *testing.T) {
	internal := autophagyBundle("local", "bateia", "target1", "", true, "MIT License\nPermission is hereby granted, free of charge")
	graph := autophagyBundle("oss", "graph", "g1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	impact := autophagyBundle("oss", "impact", "i1", "BlastRadius", true, "MIT License\nPermission is hereby granted, free of charge")
	both := autophagyBundle("oss", "both", "b1", "DependencyGraph BlastRadius", true, "MIT License\nPermission is hereby granted, free of charge")
	forbidden := autophagyBundle("oss", "forbidden", "f1", "DependencyGraph BlastRadius", true, "GNU Affero General Public License")
	r, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{internal}, []domain.EvidenceBundle{impact, graph, forbidden, both}))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Katamari.Components) != 0 {
		t.Fatalf("did not select minimum project set: %#v", r.Katamari.Components)
	}
	rejected := false
	for _, x := range r.Katamari.Rejected {
		rejected = rejected || x.Candidate == "oss/forbidden@f1"
	}
	if !rejected {
		t.Fatal("forbidden license was not witnessed")
	}
}

func TestKatamariIsInvariantToCompetitorOrder(t *testing.T) {
	internal := autophagyBundle("local", "bateia", "target1", "", true, "MIT License\nPermission is hereby granted, free of charge")
	a := autophagyBundle("oss", "a", "a1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	b := autophagyBundle("oss", "b", "b1", "BlastRadius", true, "MIT License\nPermission is hereby granted, free of charge")
	left, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{internal}, []domain.EvidenceBundle{a, b}))
	if err != nil {
		t.Fatal(err)
	}
	right, err := RunCompetitiveAutophagy(autophagyRequest([]domain.EvidenceBundle{internal}, []domain.EvidenceBundle{b, a}))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(left.Katamari.Components, right.Katamari.Components) {
		t.Fatal("input order changed proposal")
	}
}

func TestHeadAndInvalidPolicyAreRejected(t *testing.T) {
	req := autophagyRequest(nil, nil)
	req.Target.Commit = "HEAD"
	if _, err := RunCompetitiveAutophagy(req); err == nil {
		t.Fatal("HEAD target survived")
	}
	req = autophagyRequest(nil, nil)
	req.KatamariPolicy.MaxProjects = 0
	if _, err := RunCompetitiveAutophagy(req); err == nil {
		t.Fatal("invalid policy survived")
	}
}

func TestCompetitiveMutationGateRejectsEverySemanticMutation(t *testing.T) {
	validBundle := autophagyBundle("local", "bateia", "target1", "DependencyGraph", true, "MIT License\nPermission is hereby granted, free of charge")
	base := autophagyRequest([]domain.EvidenceBundle{validBundle}, nil)
	mutations := []func(domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest{
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.Target.Commit = "HEAD"
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.KatamariPolicy.MaxProjects = 0
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.KatamariPolicy.MaxComponents = 0
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.SchemaVersion = "2"
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.Internal[0].Repository.Commit = "HEAD"
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.Internal[0].Items[0].ContentHash = "tampered"
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.Internal[0].Items[0].Ref.Repo = "other"
			return r
		},
		func(r domain.CompetitiveAutophagyRequest) domain.CompetitiveAutophagyRequest {
			r.Competitors = []domain.EvidenceBundle{autophagyBundle("oss", "opaque", "HEAD", "BlastRadius", true, "MIT License\nPermission is hereby granted, free of charge")}
			return r
		},
	}
	rejected := 0
	for _, mutate := range mutations {
		encoded, _ := json.Marshal(base)
		var isolated domain.CompetitiveAutophagyRequest
		_ = json.Unmarshal(encoded, &isolated)
		if _, err := RunCompetitiveAutophagy(mutate(isolated)); err != nil {
			rejected++
		}
	}
	score := float64(rejected) / float64(len(mutations))
	if score != 1 {
		t.Fatalf("mutation_score=%v (%d/%d)", score, rejected, len(mutations))
	}
}
