/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0069_gate/discriminatory_r1_test.md
 * @layer L4
 * @updated 2026-09-12
 */

package step0069_gate

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

type r1FixtureFile struct {
	SchemaVersion          string          `json:"schema_version"`
	RequiredMutationScore  float64         `json:"required_mutation_score"`
	PermutationRepetitions int             `json:"permutation_repetitions"`
	Cases                  []r1FixtureCase `json:"cases"`
}

type r1FixtureCase struct {
	ID, Family, Operation, Expected, Witness string
}

type r1Perspective struct {
	ID, Kind, ProjectionID, OriginID string
}

type r1Target struct {
	ID, ObservationID, PerspectiveID string
}

type r1Input struct {
	Origins      []string
	Perspectives []r1Perspective
	Targets      []r1Target
}

type r1Node struct {
	Stage, Row          int
	Role, Kind          string
	Identity, ID, Label string
}

type r1Flow struct {
	Kind, ProjectionID, OriginID, ObservationID, Label string
}

type r1Card struct {
	TargetID, Title string
	Items           []string
}

type r1Document struct {
	Nodes []r1Node
	Flows []r1Flow
	Cards []r1Card
}

type r1Receipt struct {
	Input                r1Input
	Document             *r1Document
	RequestedCounts      [3]int
	RenderBudgetExceeded bool
	ReasonCode           string
	SerializedBytes      int
	PermutationBytes     []string
}

func r1TypedID(prefix, identity string) string {
	return prefix + hex.EncodeToString([]byte(identity))
}

func normativeShort(identity string) string {
	b := []byte(identity)
	for i := range b {
		if !((b[i] >= 'A' && b[i] <= 'Z') || (b[i] >= 'a' && b[i] <= 'z') || (b[i] >= '0' && b[i] <= '9')) {
			b[i] = '_'
		}
	}
	if len(b) <= 8 {
		return string(b)
	}
	return string(b[:7]) + "~"
}

func r1Control(repetitions int) r1Receipt {
	in := r1Input{
		Origins: []string{"origin:alpha/long", "origin beta-long", "origin.gamma-long", "origin#delta-long", "origin:épsilon-long"},
		Perspectives: []r1Perspective{
			{ID: "perspective:problem-long", Kind: "PROB", ProjectionID: "projection:problem/complete", OriginID: "origin:alpha/long"},
			{ID: "perspective:behavior-long", Kind: "BEHV", ProjectionID: "projection:behavior/complete", OriginID: "origin beta-long"},
			{ID: "perspective:architecture-long", Kind: "ARCH", ProjectionID: "projection:architecture/complete", OriginID: "origin.gamma-long"},
			{ID: "perspective:constraints-long", Kind: "CONS", ProjectionID: "projection:constraints/complete", OriginID: "origin#delta-long"},
			{ID: "perspective:operations-long", Kind: "OPER", ProjectionID: "projection:operations/complete", OriginID: "origin:épsilon-long"},
		},
		Targets: []r1Target{
			{ID: "target:one/long", ObservationID: "observation:one/complete", PerspectiveID: "perspective:problem-long"},
			{ID: "target two-long", ObservationID: "observation:two/complete", PerspectiveID: "perspective:behavior-long"},
			{ID: "target.three-long", ObservationID: "observation:three/complete", PerspectiveID: "perspective:architecture-long"},
			{ID: "target#four-long", ObservationID: "observation:four/complete", PerspectiveID: "perspective:constraints-long"},
			{ID: "target:cinco-longo", ObservationID: "observation:five/complete", PerspectiveID: "perspective:operations-long"},
		},
	}
	doc := buildR1Document(in)
	permutations := make([]string, repetitions)
	for i := range permutations {
		permutations[i] = "canonical-r1-bytes"
	}
	return r1Receipt{Input: in, Document: doc, RequestedCounts: [3]int{5, 5, 5}, SerializedBytes: 8192, PermutationBytes: permutations}
}

func buildR1Document(in r1Input) *r1Document {
	d := &r1Document{}
	for _, identity := range in.Origins {
		d.Nodes = append(d.Nodes, r1Node{Stage: 0, Role: "origin", Identity: identity, ID: r1TypedID("origin_", identity), Label: "O " + normativeShort(identity)})
	}
	for _, p := range in.Perspectives {
		d.Nodes = append(d.Nodes, r1Node{Stage: 1, Role: "perspective", Kind: p.Kind, Identity: p.ID, ID: r1TypedID("perspective_", p.ID), Label: "P " + p.Kind})
		d.Flows = append(d.Flows, r1Flow{Kind: "projection", ProjectionID: p.ProjectionID, OriginID: p.OriginID, Label: "projection=" + p.ProjectionID + " origin=" + p.OriginID})
	}
	for _, target := range in.Targets {
		d.Nodes = append(d.Nodes, r1Node{Stage: 2, Role: "target", Identity: target.ID, ID: r1TypedID("target_", target.ID), Label: "T " + normativeShort(target.ID)})
		d.Flows = append(d.Flows, r1Flow{Kind: "observation", ObservationID: target.ObservationID, Label: "observation=" + target.ObservationID})
		d.Cards = append(d.Cards, r1Card{TargetID: target.ID, Title: "Target " + target.ID, Items: []string{"provenance=" + target.ObservationID}})
	}
	sort.Slice(d.Nodes, func(i, j int) bool {
		if d.Nodes[i].Stage != d.Nodes[j].Stage {
			return d.Nodes[i].Stage < d.Nodes[j].Stage
		}
		return d.Nodes[i].Identity < d.Nodes[j].Identity
	})
	rows := [3]int{}
	for i := range d.Nodes {
		d.Nodes[i].Row = rows[d.Nodes[i].Stage]
		rows[d.Nodes[i].Stage]++
	}
	return d
}

func cloneR1Receipt(in r1Receipt) r1Receipt {
	out := in
	out.Input.Origins = append([]string(nil), in.Input.Origins...)
	out.Input.Perspectives = append([]r1Perspective(nil), in.Input.Perspectives...)
	out.Input.Targets = append([]r1Target(nil), in.Input.Targets...)
	out.PermutationBytes = append([]string(nil), in.PermutationBytes...)
	if in.Document != nil {
		d := *in.Document
		d.Nodes = append([]r1Node(nil), in.Document.Nodes...)
		d.Flows = append([]r1Flow(nil), in.Document.Flows...)
		d.Cards = append([]r1Card(nil), in.Document.Cards...)
		for i := range d.Cards {
			d.Cards[i].Items = append([]string(nil), in.Document.Cards[i].Items...)
		}
		out.Document = &d
	}
	return out
}

func rejectR1Budget(r *r1Receipt, stage int) {
	r.RequestedCounts[stage]++
	r.RenderBudgetExceeded = true
	r.ReasonCode = "render_budget_exceeded"
	r.Document = nil
	r.SerializedBytes = 0
}

func addSixthNode(r *r1Receipt, stage int) {
	role, prefix, identity, label := "origin", "origin_", "origin:six/complete", "O "+normativeShort("origin:six/complete")
	if stage == 1 {
		role, prefix, identity, label = "perspective", "perspective_", "perspective:six/complete", "P PROB"
	}
	if stage == 2 {
		role, prefix, identity, label = "target", "target_", "target:six/complete", "T "+normativeShort("target:six/complete")
	}
	r.Document.Nodes = append(r.Document.Nodes, r1Node{Stage: stage, Row: 5, Role: role, Kind: "PROB", Identity: identity, ID: r1TypedID(prefix, identity), Label: label})
	r.RequestedCounts[stage]++
	r.RenderBudgetExceeded = true
	r.ReasonCode = "render_budget_exceeded"
}

func mutateR1(base r1Receipt, operation string) r1Receipt {
	r := cloneR1Receipt(base)
	switch operation {
	case "five_by_five_by_five", "long_punctuated_identities":
	case "sixth_origin_rejected":
		rejectR1Budget(&r, 0)
	case "sixth_perspective_rejected":
		rejectR1Budget(&r, 1)
	case "sixth_target_rejected":
		rejectR1Budget(&r, 2)
	case "emit_sixth_origin":
		addSixthNode(&r, 0)
	case "emit_sixth_perspective":
		addSixthNode(&r, 1)
	case "emit_sixth_target":
		addSixthNode(&r, 2)
	case "row_outside_renderer":
		for i := range r.Document.Nodes {
			if r.Document.Nodes[i].Stage == 1 && r.Document.Nodes[i].Row == 4 {
				r.Document.Nodes[i].Row = 5
				break
			}
		}
	case "label_over_ten_bytes":
		r.Document.Nodes[0].Label = "O 123456789"
	case "wrong_short_label":
		r.Document.Nodes[0].Label = "O wrong"
	case "id_uses_short_identity":
		r.Document.Nodes[0].ID = r1TypedID("origin_", normativeShort(r.Document.Nodes[0].Identity))
	case "flow_loses_full_identity":
		r.Document.Flows[0].Label = "projection=project~ origin=origin~"
	case "card_loses_full_identity":
		r.Document.Cards[0].Title = "Target target~"
		r.Document.Cards[0].Items = []string{"provenance=observ~"}
	case "multibyte_label_over_ten_bytes":
		r.Document.Nodes[0].Label = "T ééééé"
	default:
		panic("unknown r1 fixture operation: " + operation)
	}
	return r
}

func classifyR1(r r1Receipt, repetitions int) (classification, string) {
	if r.RenderBudgetExceeded {
		if r.ReasonCode != "render_budget_exceeded" {
			return violated, "render budget used a non-canonical reason"
		}
		over := 0
		for _, count := range r.RequestedCounts {
			if count == 6 {
				over++
			} else if count != 5 {
				return violated, "render budget receipt has non-observable stage counts"
			}
		}
		if over != 1 {
			return violated, "render budget receipt does not identify one sixth stage entity"
		}
		if r.Document != nil || r.SerializedBytes != 0 {
			return violated, "render budget emitted a partial document"
		}
		return unknown, "render budget failed closed"
	}
	if r.Document == nil {
		return violated, "renderable input produced no document"
	}
	if r.RequestedCounts != [3]int{5, 5, 5} {
		return violated, "renderable receipt has inconsistent requested stage counts"
	}
	if len(r.PermutationBytes) != repetitions {
		return violated, "missing r1 permutation repetitions"
	}
	for _, value := range r.PermutationBytes {
		if value != r.PermutationBytes[0] {
			return violated, "r1 equivalent permutations changed bytes"
		}
	}
	if reason := validateR1Document(r.Input, r.Document); reason != "" {
		return violated, reason
	}
	return preserved, "r1 renderer boundary preserved"
}

func validateR1Document(in r1Input, got *r1Document) string {
	want := buildR1Document(in)
	counts := [3]int{}
	for _, n := range got.Nodes {
		if n.Stage < 0 || n.Stage > 2 {
			return "invalid stage"
		}
		counts[n.Stage]++
		if n.Row < 0 || n.Row > 4 {
			return "row outside 0..4"
		}
		if len([]byte(n.Label)) > 10 || !isASCII(n.Label) {
			return "node label exceeds ten ASCII bytes"
		}
		prefix := map[string]string{"origin": "origin_", "perspective": "perspective_", "target": "target_"}[n.Role]
		if prefix == "" || n.ID != r1TypedID(prefix, n.Identity) {
			return "node ID does not preserve full identity"
		}
		expectedLabel := ""
		switch n.Role {
		case "origin":
			expectedLabel = "O " + normativeShort(n.Identity)
		case "perspective":
			if !containsKind(n.Kind) {
				return "invalid perspective kind"
			}
			expectedLabel = "P " + n.Kind
		case "target":
			expectedLabel = "T " + normativeShort(n.Identity)
		}
		if n.Label != expectedLabel {
			return "node label violates normative short function"
		}
	}
	for _, count := range counts {
		if count > 5 {
			return "more than five nodes in a stage"
		}
	}
	if !reflect.DeepEqual(got.Nodes, want.Nodes) {
		return "nodes or canonical rows diverge"
	}
	if len(got.Flows) != len(want.Flows) {
		return "flow set diverges"
	}
	for _, f := range got.Flows {
		if f.Kind == "projection" && (!strings.Contains(f.Label, f.ProjectionID) || !strings.Contains(f.Label, f.OriginID)) {
			return "projection flow label loses full identity"
		}
		if f.Kind == "observation" && !strings.Contains(f.Label, f.ObservationID) {
			return "observation flow label loses full identity"
		}
	}
	if len(got.Cards) != len(in.Targets) {
		return "card set diverges"
	}
	byTarget := map[string]r1Target{}
	for _, target := range in.Targets {
		byTarget[target.ID] = target
	}
	for _, c := range got.Cards {
		target, ok := byTarget[c.TargetID]
		if !ok || !strings.Contains(c.Title, target.ID) || !strings.Contains(strings.Join(c.Items, "\n"), target.ObservationID) {
			return "card loses full TargetID or provenance"
		}
	}
	return ""
}

func isASCII(value string) bool {
	for _, b := range []byte(value) {
		if b > 0x7f {
			return false
		}
	}
	return true
}

func containsKind(kind string) bool {
	for _, allowed := range []string{"PROB", "BEHV", "ARCH", "CONS", "OPER"} {
		if kind == allowed {
			return true
		}
	}
	return false
}

func loadR1Fixtures(t *testing.T) r1FixtureFile {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate r1 gate source")
	}
	path := filepath.Join(filepath.Dir(here), "..", "..", "..", "..", "lab", "bancada", "history", "benchmarks", "0069", "adversarial-fixtures-r1.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read r1 fixtures: %v", err)
	}
	var fixtures r1FixtureFile
	if err := json.Unmarshal(b, &fixtures); err != nil {
		t.Fatalf("decode r1 fixtures: %v", err)
	}
	if fixtures.SchemaVersion != "idea-radar-archify-adversarial-r1" {
		t.Fatalf("unexpected r1 fixture schema %q", fixtures.SchemaVersion)
	}
	return fixtures
}

func classifyR1Corpus(t *testing.T, fixtures r1FixtureFile, order []int) map[string]classification {
	t.Helper()
	base := r1Control(fixtures.PermutationRepetitions)
	results := make(map[string]classification, len(order))
	for _, index := range order {
		c := fixtures.Cases[index]
		first, reason := classifyR1(mutateR1(base, c.Operation), fixtures.PermutationRepetitions)
		second, _ := classifyR1(mutateR1(base, c.Operation), fixtures.PermutationRepetitions)
		if first != second {
			t.Fatalf("%s is not repeatable: %s != %s", c.ID, first, second)
		}
		if string(first) != c.Expected {
			t.Errorf("%s (%s): got %s, want %s; classifier witness: %s; fixture witness: %s", c.ID, c.Family, first, c.Expected, reason, c.Witness)
		}
		results[c.ID] = first
	}
	return results
}

func TestStep0069DiscriminatoryGateR1(t *testing.T) {
	fixtures := loadR1Fixtures(t)
	if fixtures.PermutationRepetitions != 20 || fixtures.RequiredMutationScore != 1.0 {
		t.Fatalf("invalid r1 budget: repetitions=%d score=%.3f", fixtures.PermutationRepetitions, fixtures.RequiredMutationScore)
	}
	order := make([]int, len(fixtures.Cases))
	mutants := 0
	for i, c := range fixtures.Cases {
		order[i] = i
		if strings.HasPrefix(c.ID, "M") {
			mutants++
		}
	}
	baseline := classifyR1Corpus(t, fixtures, order)
	for repetition := 0; repetition < fixtures.PermutationRepetitions; repetition++ {
		permuted := append([]int(nil), order...)
		shift := repetition % len(permuted)
		permuted = append(permuted[shift:], permuted[:shift]...)
		if repetition%2 == 1 {
			for left, right := 0, len(permuted)-1; left < right; left, right = left+1, right-1 {
				permuted[left], permuted[right] = permuted[right], permuted[left]
			}
		}
		if got := classifyR1Corpus(t, fixtures, permuted); !reflect.DeepEqual(got, baseline) {
			t.Fatalf("r1 corpus order %d changed classifications", repetition)
		}
	}
	rejected := 0
	for id, result := range baseline {
		if strings.HasPrefix(id, "M") && result == violated {
			rejected++
		}
	}
	score := float64(rejected) / float64(mutants)
	t.Logf("r1_mutation_score=%.3f rejected=%d valid_mutations=%d required=%.3f", score, rejected, mutants, fixtures.RequiredMutationScore)
	if score != fixtures.RequiredMutationScore {
		t.Fatalf("r1 mutation score %.3f does not meet %.3f", score, fixtures.RequiredMutationScore)
	}
	for _, c := range fixtures.Cases {
		if strings.HasPrefix(c.ID, "P") && baseline[c.ID] != preserved {
			t.Fatalf("r1 control %s was not Preserved", c.ID)
		}
		if strings.HasPrefix(c.ID, "U") && baseline[c.ID] != unknown {
			t.Fatalf("r1 opaque budget case %s was not Unknown", c.ID)
		}
	}
}
