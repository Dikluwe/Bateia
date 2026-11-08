/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_insight_contract_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func algorithmContractFixture() domain.AlgorithmInsightRequest {
	mk := func(owner string) domain.EvidenceBundle {
		ref := domain.CommitRef{Owner: owner, Repo: "scheduler", Commit: strings.Repeat("a", 40)}
		lines := []string{}
		for i := 0; i < 24; i++ {
			lines = append(lines, fmt.Sprintf("%s candidate %d", owner, i))
		}
		content := strings.Join(lines, "\n")
		digest := sha256.Sum256([]byte(content))
		return domain.EvidenceBundle{SchemaVersion: "1", Repository: ref, Items: []domain.CollectedEvidence{{Kind: domain.EvidenceDocumentation, Ref: ref, Path: "schedule.txt", Content: content, ContentHash: fmt.Sprintf("sha256:%x", digest)}}, RequirementCoverage: []domain.RequirementCollectionCoverage{{RequirementKind: "capability", Requirement: "select ready work", Status: domain.RequirementTargetedExamined, CandidatePaths: []string{"schedule.txt"}, ExaminedPaths: []string{"schedule.txt"}}}}
	}
	source, target := mk("reference"), mk("target")
	return domain.AlgorithmInsightRequest{SchemaVersion: "1", Question: "How to schedule ready work?", Sources: []domain.EvidenceBundle{source}, Target: target, Ideas: []domain.AlgorithmIdea{{Title: "Ready candidates", Principle: "Maintain candidates when readiness changes.", Steps: []string{"Observe readiness changes.", "Choose a ready candidate."}, Evidence: []domain.AlgorithmIdeaEvidence{algorithmContractReference(source, 1)}, ApplicationConditions: []string{"Readiness changes are observable."}, FailureConditions: []string{"Updates can be missed."}, Decision: domain.AlgorithmDecision{Current: "Scan all work.", Proposed: "Maintain ready work.", Benefit: "May reduce selection work.", Cost: "Update bookkeeping."}, TargetRelation: "DIFFERENT_APPROACH", TargetEvidence: []domain.AlgorithmIdeaEvidence{algorithmContractReference(target, 1)}, ValidationObservable: "Selection work and missed candidates.", RejectIf: "Any ready work is lost."}}}
}

func algorithmContractReference(b domain.EvidenceBundle, line int) domain.AlgorithmIdeaEvidence {
	item := b.Items[0]
	return domain.AlgorithmIdeaEvidence{Source: domain.ImmutableCandidateRef{Repository: b.Repository.Owner + "/" + b.Repository.Repo, Commit: b.Repository.Commit}, LocatedInsightReference: domain.LocatedInsightReference{InsightObservationReference: domain.InsightObservationReference{Path: item.Path, ContentHash: item.ContentHash, StartLine: line, EndLine: line}, Excerpt: strings.Split(item.Content, "\n")[line-1]}}
}

func algorithmContractJSON(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func checkAlgorithmContract(t *testing.T, r domain.AlgorithmInsightRequest, valid bool) domain.AlgorithmInsightReport {
	t.Helper()
	before := algorithmContractJSON(t, r)
	got, err := FormulateAlgorithmInsights(context.Background(), r)
	if string(before) != string(algorithmContractJSON(t, r)) {
		t.Fatal("caller input mutated")
	}
	if !valid {
		if err == nil || got.ID != "" {
			t.Fatalf("invalid input produced usable report: %v", err)
		}
	} else if err != nil || got.ID == "" || got.EffectState != "UNVERIFIED" || got.InterpretationState != "DECLARED_NOT_REASSESSED" || got.TargetRelationState != "DECLARED_NOT_REASSESSED" {
		t.Fatalf("valid input rejected or promoted: %v", err)
	}
	return got
}

func TestAlgorithmContractReferencePermutationsAndByteCausality(t *testing.T) {
	r := algorithmContractFixture()
	other := r.Sources[0]
	other.Repository.Owner = "alternative"
	other.Items = append([]domain.CollectedEvidence(nil), other.Items...)
	other.Items[0].Ref = other.Repository
	r.Sources = append(r.Sources, other)
	r.Ideas[0].Evidence = append(r.Ideas[0].Evidence, algorithmContractReference(other, 2), algorithmContractReference(r.Sources[0], 3))
	r.Ideas[0].TargetEvidence = append(r.Ideas[0].TargetEvidence, algorithmContractReference(r.Target, 2), algorithmContractReference(r.Target, 3))
	baseline := checkAlgorithmContract(t, r, true)
	for _, order := range [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}} {
		for _, reverseBundles := range []bool{false, true} {
			for _, targetOrder := range [][]int{{0, 1, 2}, {2, 0, 1}, {2, 1, 0}} {
				copy := r
				copy.Sources = append([]domain.EvidenceBundle(nil), r.Sources...)
				copy.Ideas = append([]domain.AlgorithmIdea(nil), r.Ideas...)
				copy.Ideas[0].Evidence = append([]domain.AlgorithmIdeaEvidence(nil), r.Ideas[0].Evidence...)
				copy.Ideas[0].TargetEvidence = append([]domain.AlgorithmIdeaEvidence(nil), r.Ideas[0].TargetEvidence...)
				for i, j := range order {
					copy.Ideas[0].Evidence[i] = r.Ideas[0].Evidence[j]
				}
				for i, j := range targetOrder {
					copy.Ideas[0].TargetEvidence[i] = r.Ideas[0].TargetEvidence[j]
				}
				if reverseBundles {
					copy.Sources[0], copy.Sources[1] = copy.Sources[1], copy.Sources[0]
				}
				if got := checkAlgorithmContract(t, copy, true); !reflect.DeepEqual(baseline, got) {
					t.Fatalf("reference/bundle permutation changed report: order=%v target=%v bundles=%v", order, targetOrder, reverseBundles)
				}
			}
		}
	}
	for _, target := range []bool{false, true} {
		changed := algorithmContractFixture()
		old := checkAlgorithmContract(t, changed, true)
		bundle := &changed.Sources[0]
		refs := changed.Ideas[0].Evidence
		if target {
			bundle = &changed.Target
			refs = changed.Ideas[0].TargetEvidence
		}
		bundle.Items[0].Content = strings.Replace(bundle.Items[0].Content, "candidate 0", "candidate changed", 1)
		hash := sha256.Sum256([]byte(bundle.Items[0].Content))
		bundle.Items[0].ContentHash = fmt.Sprintf("sha256:%x", hash)
		refs[0] = algorithmContractReference(*bundle, 1)
		next := checkAlgorithmContract(t, changed, true)
		if old.ID == next.ID || old.InputSHA256 == next.InputSHA256 || old.Ideas[0].ID == next.Ideas[0].ID {
			t.Fatal("cited source/target bytes are not causal")
		}
	}
	for _, target := range []bool{false, true} {
		duplicate := algorithmContractFixture()
		refs := &duplicate.Ideas[0].Evidence
		if target {
			refs = &duplicate.Ideas[0].TargetEvidence
		}
		*refs = append(*refs, (*refs)[0])
		checkAlgorithmContract(t, duplicate, false)
	}
}

func TestAlgorithmContractTextBoundaries(t *testing.T) {
	fields := []struct {
		name  string
		max   int
		field func(*domain.AlgorithmInsightRequest) *string
	}{
		{"question", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Question }},
		{"title", 200, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Title }},
		{"principle", 1500, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Principle }},
		{"current", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Decision.Current }},
		{"proposed", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Decision.Proposed }},
		{"benefit", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Decision.Benefit }},
		{"cost", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Decision.Cost }},
		{"observable", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].ValidationObservable }},
		{"reject", 1000, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].RejectIf }},
		{"step", 600, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].Steps[0] }},
		{"application", 600, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].ApplicationConditions[0] }},
		{"failure", 600, func(r *domain.AlgorithmInsightRequest) *string { return &r.Ideas[0].FailureConditions[0] }},
	}
	for _, f := range fields {
		t.Run(f.name, func(t *testing.T) {
			for _, value := range []string{strings.Repeat("x", f.max), strings.Repeat("x", f.max+1), "", " ", "bad\x00text", "bad\u0085text"} {
				r := algorithmContractFixture()
				*f.field(&r) = value
				checkAlgorithmContract(t, r, len(value) == f.max)
			}
			r := algorithmContractFixture()
			*f.field(&r) = strings.Repeat("é", f.max/2)
			checkAlgorithmContract(t, r, true)
			*f.field(&r) += "x"
			checkAlgorithmContract(t, r, false)
		})
	}
	for _, relation := range []string{"ALREADY_PRESENT", "DIFFERENT_APPROACH", "NOT_ESTABLISHED"} {
		r := algorithmContractFixture()
		r.Ideas[0].TargetRelation = relation
		checkAlgorithmContract(t, r, true)
	}
	r := algorithmContractFixture()
	r.SchemaVersion = "2"
	checkAlgorithmContract(t, r, false)
}

func TestAlgorithmContractCollectionBoundaries(t *testing.T) {
	for _, count := range []int{0, 1, 8, 9} {
		t.Run(fmt.Sprintf("sources/%d", count), func(t *testing.T) {
			r := algorithmContractFixture()
			source := r.Sources[0]
			r.Sources = nil
			for i := 0; i < count; i++ {
				b := source
				b.Repository.Owner = fmt.Sprintf("source%d", i)
				b.Items = append([]domain.CollectedEvidence(nil), source.Items...)
				b.Items[0].Ref = b.Repository
				r.Sources = append(r.Sources, b)
			}
			if count > 0 {
				r.Ideas[0].Evidence[0] = algorithmContractReference(r.Sources[0], 1)
			}
			checkAlgorithmContract(t, r, count >= 1 && count <= 8)
		})
		t.Run(fmt.Sprintf("ideas/%d", count), func(t *testing.T) {
			r := algorithmContractFixture()
			idea := r.Ideas[0]
			r.Ideas = nil
			for i := 0; i < count; i++ {
				v := idea
				v.Title = fmt.Sprintf("idea %d", i)
				r.Ideas = append(r.Ideas, v)
			}
			checkAlgorithmContract(t, r, count >= 1 && count <= 8)
		})
		for _, target := range []bool{false, true} {
			t.Run(fmt.Sprintf("references/target=%v/%d", target, count), func(t *testing.T) {
				r := algorithmContractFixture()
				b := r.Sources[0]
				refs := &r.Ideas[0].Evidence
				if target {
					b = r.Target
					refs = &r.Ideas[0].TargetEvidence
				}
				*refs = nil
				for i := 0; i < count; i++ {
					*refs = append(*refs, algorithmContractReference(b, i+1))
				}
				checkAlgorithmContract(t, r, count >= 1 && count <= 8)
			})
		}
	}
	for _, field := range []string{"steps", "application", "failure"} {
		for _, count := range []int{0, 1, 2, 12, 13} {
			t.Run(fmt.Sprintf("%s/%d", field, count), func(t *testing.T) {
				r := algorithmContractFixture()
				values := []string{}
				for i := 0; i < count; i++ {
					values = append(values, fmt.Sprintf("condition or step %d", i))
				}
				minimum := 1
				switch field {
				case "steps":
					r.Ideas[0].Steps = values
					minimum = 2
				case "application":
					r.Ideas[0].ApplicationConditions = values
				case "failure":
					r.Ideas[0].FailureConditions = values
				}
				checkAlgorithmContract(t, r, count >= minimum && count <= 12)
			})
		}
	}
	for _, count := range []int{512, 513} {
		t.Run(fmt.Sprintf("items/%d", count), func(t *testing.T) {
			r := algorithmContractFixture()
			b := &r.Sources[0]
			first := b.Items[0]
			for i := 1; i < count; i++ {
				item := first
				item.Path = fmt.Sprintf("note-%d.txt", i)
				b.Items = append(b.Items, item)
			}
			checkAlgorithmContract(t, r, count == 512)
		})
	}
}

func TestAlgorithmContractCitationBoundaries(t *testing.T) {
	for _, target := range []bool{false, true} {
		for _, tc := range []struct {
			name, content, excerpt string
			start, end             int
			valid                  bool
		}{
			{"20 lines", strings.Repeat("x\n", 20), "x", 1, 20, true},
			{"21 lines", strings.Repeat("x\n", 21), "x", 1, 21, false},
			{"4096 span bytes", strings.Repeat("x", 4096), "x", 1, 1, true},
			{"4097 span bytes", strings.Repeat("x", 4097), "x", 1, 1, false},
			{"2048 excerpt bytes", strings.Repeat("x", 2049), strings.Repeat("x", 2048), 1, 1, true},
			{"2049 excerpt bytes", strings.Repeat("x", 2049), strings.Repeat("x", 2049), 1, 1, false},
			{"empty excerpt", "text", "", 1, 1, false},
			{"outside literal", "first\nsecond", "second", 1, 1, false},
			{"zero line", "text", "text", 0, 1, false},
			{"backward range", "first\nsecond", "first", 2, 1, false},
		} {
			t.Run(fmt.Sprintf("target=%v/%s", target, tc.name), func(t *testing.T) {
				r := algorithmContractFixture()
				b := &r.Sources[0]
				refs := r.Ideas[0].Evidence
				if target {
					b = &r.Target
					refs = r.Ideas[0].TargetEvidence
				}
				b.Items[0].Content = tc.content
				hash := sha256.Sum256([]byte(tc.content))
				b.Items[0].ContentHash = fmt.Sprintf("sha256:%x", hash)
				refs[0].ContentHash = b.Items[0].ContentHash
				refs[0].StartLine = tc.start
				refs[0].EndLine = tc.end
				refs[0].Excerpt = tc.excerpt
				checkAlgorithmContract(t, r, tc.valid)
			})
		}
	}
}

func setAlgorithmContractPadding(r *domain.AlgorithmInsightRequest, target bool, size int) {
	b := &r.Sources[0]
	refs := r.Ideas[0].Evidence
	if target {
		b = &r.Target
		refs = r.Ideas[0].TargetEvidence
	}
	prefix := strings.Split(b.Items[0].Content, "\n")[0] + "\n"
	b.Items[0].Content = prefix + strings.Repeat("x", size-len(prefix))
	hash := sha256.Sum256([]byte(b.Items[0].Content))
	b.Items[0].ContentHash = fmt.Sprintf("sha256:%x", hash)
	refs[0].ContentHash = b.Items[0].ContentHash
}

func TestAlgorithmContractSerializedByteBoundaries(t *testing.T) {
	for _, target := range []bool{false, true} {
		r := algorithmContractFixture()
		group := func() domain.InsightExplorationRequest {
			bundles := r.Sources
			if target {
				bundles = []domain.EvidenceBundle{r.Target}
			}
			return domain.InsightExplorationRequest{SchemaVersion: "1", Question: r.Question, Bundles: bundles}
		}
		initial := len(algorithmContractJSON(t, group()))
		b := r.Sources[0]
		if target {
			b = r.Target
		}
		size := len(b.Items[0].Content) + (2 << 20) - initial
		setAlgorithmContractPadding(&r, target, size)
		// Padding removes escaped newlines; measure the actual wire representation.
		size += (2 << 20) - len(algorithmContractJSON(t, group()))
		setAlgorithmContractPadding(&r, target, size)
		if got := len(algorithmContractJSON(t, group())); got != 2<<20 {
			t.Fatalf("wrong group size: %d", got)
		}
		checkAlgorithmContract(t, r, true)
		setAlgorithmContractPadding(&r, target, size+1)
		checkAlgorithmContract(t, r, false)
	}
	r := algorithmContractFixture()
	for i := 0; i < 12; i++ {
		value := fmt.Sprintf("%03d", i) + strings.Repeat("x", 597)
		if i < 2 {
			r.Ideas[0].Steps[i] = value
		} else {
			r.Ideas[0].Steps = append(r.Ideas[0].Steps, value)
		}
	}
	r.Ideas[0].ApplicationConditions = nil
	r.Ideas[0].FailureConditions = nil
	for i := 0; i < 12; i++ {
		value := fmt.Sprintf("%03d", i) + strings.Repeat("x", 597)
		r.Ideas[0].ApplicationConditions = append(r.Ideas[0].ApplicationConditions, value)
		r.Ideas[0].FailureConditions = append(r.Ideas[0].FailureConditions, value)
	}
	setAlgorithmContractPadding(&r, false, (2<<20)-8192)
	setAlgorithmContractPadding(&r, true, (2<<20)-8192)
	excess := len(algorithmContractJSON(t, r)) - (4 << 20)
	if excess <= 0 {
		t.Fatal("fixture must straddle request ceiling")
	}
	targetSize := len(r.Target.Items[0].Content) - excess
	setAlgorithmContractPadding(&r, true, targetSize)
	if got := len(algorithmContractJSON(t, r)); got != 4<<20 {
		t.Fatalf("wrong request size: %d", got)
	}
	checkAlgorithmContract(t, r, true)
	setAlgorithmContractPadding(&r, true, targetSize+1)
	checkAlgorithmContract(t, r, false)
}
