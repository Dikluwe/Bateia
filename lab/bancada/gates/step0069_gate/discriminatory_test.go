/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0069_gate/discriminatory_test.md
 * @layer L4
 * @updated 2026-09-12
 */

package step0069_gate

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

type classification string

const (
	preserved classification = "Preserved"
	violated  classification = "Violated"
	unknown   classification = "Unknown"
)

type fixtureFile struct {
	SchemaVersion          string        `json:"schema_version"`
	RequiredMutationScore  float64       `json:"required_mutation_score"`
	PermutationRepetitions int           `json:"permutation_repetitions"`
	Cases                  []fixtureCase `json:"cases"`
}

type fixtureCase struct {
	ID        string `json:"id"`
	Family    string `json:"family"`
	Operation string `json:"operation"`
	Expected  string `json:"expected"`
	Witness   string `json:"witness"`
}

type semanticInput struct {
	Projections  []projection
	Observations []observation
	Targets      []targetPlan
}

type projection struct {
	ID, PerspectiveID, OriginStatus, OriginID string
}

type observation struct {
	ID, PerspectiveID, TargetID, Stance string
}

type targetPlan struct {
	ID, Outcome, Reason string
	Probes              []string
}

type node struct {
	ID, Role, Identity, Label string
	Stage, Row                int
}

type flow struct {
	ID, Kind, Identity, Source, Target, Stance string
}

type card struct {
	TargetID, Outcome, Reason string
	Probes                    []string
}

type document struct {
	SchemaVersion int      `json:"schema_version"`
	DiagramType   string   `json:"diagram_type"`
	Title         string   `json:"title"`
	Stages        []string `json:"stages"`
	Nodes         []node   `json:"nodes"`
	Flows         []flow   `json:"flows"`
	Cards         []card   `json:"cards"`
}

type receipt struct {
	Input             semanticInput
	RadarStatus       string
	PlannerStatus     string
	OpaqueSchema      bool
	BudgetExceeded    bool
	InputDigestBefore string
	InputDigestAfter  string
	Document          *document
	SerializedBytes   int
	PermutationBytes  []string
	OutputFile        bool
	Stdout            string
	AtomicRename      bool
	TruncatedEarly    bool
}

var archifyID = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func typedID(prefix, identity string) string {
	return prefix + hex.EncodeToString([]byte(identity))
}

func canonicalReceipt(repetitions int) receipt {
	in := semanticInput{
		Projections: []projection{
			{ID: "projection-recorded", PerspectiveID: "perspective:a:b", OriginStatus: "RECORDED", OriginID: "origin:alpha"},
			{ID: "projection-unknown-a", PerspectiveID: "perspective:a/b", OriginStatus: "UNKNOWN_ORIGIN"},
			{ID: "projection-unknown-b", PerspectiveID: "perspective:space edge", OriginStatus: "UNKNOWN_ORIGIN"},
		},
		Observations: []observation{
			{ID: "observation-01", PerspectiveID: "perspective:a:b", TargetID: "target:one", Stance: "support"},
			{ID: "observation-02", PerspectiveID: "perspective:a/b", TargetID: "target:one", Stance: "conflict"},
			{ID: "observation-03", PerspectiveID: "perspective:space edge", TargetID: "target:two", Stance: "ambiguous"},
		},
		Targets: []targetPlan{
			{ID: "target:one", Outcome: "ADOPT", Reason: "supported", Probes: []string{"BIND_EVIDENCE: confirm source", "CHECK_CONFLICT: resolve"}},
			{ID: "target:two", Outcome: "PROBE", Reason: "ambiguous", Probes: []string{"RECORD_ORIGIN: identify origin"}},
		},
	}
	doc := buildDocument(in)
	bytes := make([]string, repetitions)
	for i := range bytes {
		bytes[i] = "canonical-dataflow-v1-bytes"
	}
	return receipt{
		Input: in, RadarStatus: "COMPLETE", PlannerStatus: "COMPLETE",
		InputDigestBefore: "sha256:input", InputDigestAfter: "sha256:input",
		Document: doc, SerializedBytes: 4096, PermutationBytes: bytes,
		OutputFile: true, AtomicRename: true,
	}
}

func buildDocument(in semanticInput) *document {
	d := &document{SchemaVersion: 1, DiagramType: "dataflow", Title: "Bateia Idea Radar", Stages: []string{"Origins", "Perspectives", "Targets"}}
	seenOrigin := map[string]bool{}
	seenPerspective := map[string]bool{}
	seenTarget := map[string]bool{}
	for _, p := range in.Projections {
		if p.OriginStatus == "RECORDED" && !seenOrigin[p.OriginID] {
			d.Nodes = append(d.Nodes, node{ID: typedID("origin_", p.OriginID), Role: "origin", Identity: p.OriginID, Label: p.OriginID, Stage: 0})
			seenOrigin[p.OriginID] = true
		}
		if p.OriginStatus == "UNKNOWN_ORIGIN" {
			d.Nodes = append(d.Nodes, node{ID: typedID("origin_unknown_", p.ID), Role: "origin_unknown", Identity: p.ID, Label: "Unknown origin for " + p.ID, Stage: 0})
		}
		if !seenPerspective[p.PerspectiveID] {
			d.Nodes = append(d.Nodes, node{ID: typedID("perspective_", p.PerspectiveID), Role: "perspective", Identity: p.PerspectiveID, Label: p.PerspectiveID, Stage: 1})
			seenPerspective[p.PerspectiveID] = true
		}
	}
	for _, target := range in.Targets {
		if !seenTarget[target.ID] {
			d.Nodes = append(d.Nodes, node{ID: typedID("target_", target.ID), Role: "target", Identity: target.ID, Label: target.ID, Stage: 2})
			seenTarget[target.ID] = true
		}
		d.Cards = append(d.Cards, card{TargetID: target.ID, Outcome: target.Outcome, Reason: target.Reason, Probes: append([]string(nil), target.Probes...)})
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
	for _, p := range in.Projections {
		source := typedID("origin_", p.OriginID)
		if p.OriginStatus == "UNKNOWN_ORIGIN" {
			source = typedID("origin_unknown_", p.ID)
		}
		d.Flows = append(d.Flows, flow{ID: typedID("flow_projection_", p.ID), Kind: "projection", Identity: p.ID, Source: source, Target: typedID("perspective_", p.PerspectiveID)})
	}
	sort.Slice(d.Flows, func(i, j int) bool { return d.Flows[i].Identity < d.Flows[j].Identity })
	start := len(d.Flows)
	for _, o := range in.Observations {
		d.Flows = append(d.Flows, flow{ID: typedID("flow_observation_", o.ID), Kind: "observation", Identity: o.ID, Source: typedID("perspective_", o.PerspectiveID), Target: typedID("target_", o.TargetID), Stance: o.Stance})
	}
	sort.Slice(d.Flows[start:], func(i, j int) bool { return d.Flows[start+i].Identity < d.Flows[start+j].Identity })
	sort.Slice(d.Cards, func(i, j int) bool { return d.Cards[i].TargetID < d.Cards[j].TargetID })
	return d
}

func mutate(base receipt, operation string) receipt {
	r := cloneReceipt(base)
	switch operation {
	case "identity", "nil_empty_equivalent":
	case "radar_unknown_no_output":
		r.RadarStatus, r.Document, r.SerializedBytes = "UNKNOWN", nil, 0
	case "opaque_schema_no_output":
		r.OpaqueSchema, r.Document, r.SerializedBytes = true, nil, 0
	case "budget_exceeded_no_output":
		r.BudgetExceeded, r.Document, r.SerializedBytes = true, nil, 0
	case "add_target_relation":
		r.Document.Flows = append(r.Document.Flows, flow{ID: typedID("flow_observation_", "invented"), Kind: "observation", Identity: "invented", Source: typedID("target_", "target:one"), Target: typedID("target_", "target:two")})
	case "add_probe_flow":
		r.Document.Flows = append(r.Document.Flows, flow{ID: typedID("flow_observation_", "probe"), Kind: "probe", Identity: "probe", Source: typedID("target_", "target:one"), Target: typedID("target_", "target:two")})
	case "coalesce_unknown_origins":
		first := typedID("origin_unknown_", "projection-unknown-a")
		second := typedID("origin_unknown_", "projection-unknown-b")
		for i := range r.Document.Flows {
			if r.Document.Flows[i].Source == second {
				r.Document.Flows[i].Source = first
			}
		}
	case "promote_unknown_origin":
		for i := range r.Document.Nodes {
			if r.Document.Nodes[i].Role == "origin_unknown" {
				r.Document.Nodes[i].Role = "origin"
				r.Document.Nodes[i].Label = "physical origin"
				break
			}
		}
	case "collide_sanitized_ids":
		var first string
		for i := range r.Document.Nodes {
			if r.Document.Nodes[i].Role == "perspective" {
				if first == "" {
					first = r.Document.Nodes[i].ID
				} else {
					r.Document.Nodes[i].ID = first
					break
				}
			}
		}
	case "order_dependent_bytes":
		if len(r.PermutationBytes) > 1 {
			r.PermutationBytes[len(r.PermutationBytes)-1] = "different-order-bytes"
		}
	case "retain_duplicate_flow":
		r.Document.Flows = append(r.Document.Flows, r.Document.Flows[0])
	case "radar_unknown_partial":
		r.RadarStatus = "UNKNOWN"
	case "invalid_archify_id":
		r.Document.Nodes[0].ID = "9 invalid id"
	case "swap_observation_stance":
		for i := range r.Document.Flows {
			if r.Document.Flows[i].Identity == "observation-02" {
				r.Document.Flows[i].Stance = "support"
			}
		}
	case "drop_observation_flow":
		for i := range r.Document.Flows {
			if r.Document.Flows[i].Kind == "observation" {
				r.Document.Flows = append(r.Document.Flows[:i], r.Document.Flows[i+1:]...)
				break
			}
		}
	case "redirect_observation_target":
		for i := range r.Document.Flows {
			if r.Document.Flows[i].Identity == "observation-01" {
				r.Document.Flows[i].Target = typedID("target_", "target:two")
			}
		}
	case "drop_target_card":
		r.Document.Cards = r.Document.Cards[:len(r.Document.Cards)-1]
	case "drop_probe_text":
		r.Document.Cards[0].Probes = r.Document.Cards[0].Probes[:1]
	case "mutate_input_digest":
		r.InputDigestAfter = "sha256:mutated-input"
	case "overflow_with_partial":
		r.SerializedBytes = 4*1024*1024 + 1
	case "non_atomic_receipt":
		r.AtomicRename, r.TruncatedEarly = false, true
	default:
		panic("unknown fixture operation: " + operation)
	}
	return r
}

func cloneReceipt(in receipt) receipt {
	out := in
	out.Input.Projections = append([]projection(nil), in.Input.Projections...)
	out.Input.Observations = append([]observation(nil), in.Input.Observations...)
	out.Input.Targets = append([]targetPlan(nil), in.Input.Targets...)
	for i := range out.Input.Targets {
		out.Input.Targets[i].Probes = append([]string(nil), in.Input.Targets[i].Probes...)
	}
	out.PermutationBytes = append([]string(nil), in.PermutationBytes...)
	if in.Document != nil {
		d := *in.Document
		d.Stages = append([]string(nil), in.Document.Stages...)
		d.Nodes = append([]node(nil), in.Document.Nodes...)
		d.Flows = append([]flow(nil), in.Document.Flows...)
		d.Cards = append([]card(nil), in.Document.Cards...)
		for i := range d.Cards {
			d.Cards[i].Probes = append([]string(nil), in.Document.Cards[i].Probes...)
		}
		out.Document = &d
	}
	return out
}

func classify(r receipt, repetitions int) (classification, string) {
	if r.RadarStatus == "UNKNOWN" || r.PlannerStatus == "UNKNOWN" || r.OpaqueSchema || r.BudgetExceeded {
		if r.Document != nil || r.SerializedBytes != 0 || r.Stdout != "" {
			return violated, "opaque or UNKNOWN state emitted partial output"
		}
		return unknown, "explicitly opaque state failed closed"
	}
	if r.RadarStatus != "COMPLETE" || r.PlannerStatus != "COMPLETE" {
		return unknown, "unsupported authority status"
	}
	if r.Document == nil {
		return violated, "complete authorities produced no document"
	}
	if r.InputDigestBefore != r.InputDigestAfter {
		return violated, "input digest changed"
	}
	if r.SerializedBytes > 4*1024*1024 {
		return violated, "serialized byte limit exceeded with output"
	}
	if r.OutputFile && (!r.AtomicRename || r.TruncatedEarly || r.Stdout != "") {
		return violated, "output-file receipt is not atomic and exclusive"
	}
	if len(r.PermutationBytes) != repetitions {
		return violated, "missing permutation repetitions"
	}
	for _, b := range r.PermutationBytes {
		if b != r.PermutationBytes[0] {
			return violated, "equivalent permutations changed bytes"
		}
	}
	if reason := validateDocument(r.Input, r.Document); reason != "" {
		return violated, reason
	}
	return preserved, "all observable obligations preserved"
}

func validateDocument(in semanticInput, got *document) string {
	if got.SchemaVersion != 1 || got.DiagramType != "dataflow" || got.Title != "Bateia Idea Radar" || !reflect.DeepEqual(got.Stages, []string{"Origins", "Perspectives", "Targets"}) {
		return "invalid Dataflow v1 envelope"
	}
	want := buildDocument(in)
	seenIDs := map[string]bool{}
	for _, n := range got.Nodes {
		if !archifyID.MatchString(n.ID) || seenIDs[n.ID] {
			return "invalid or duplicate node ID"
		}
		seenIDs[n.ID] = true
		prefix := map[string]string{"origin": "origin_", "origin_unknown": "origin_unknown_", "perspective": "perspective_", "target": "target_"}[n.Role]
		if prefix == "" || n.ID != typedID(prefix, n.Identity) {
			return "node identity is not typed injective UTF-8 hex"
		}
		if n.Role == "origin_unknown" && !strings.Contains(strings.ToLower(n.Label), "unknown") {
			return "unknown origin lacks descriptive label"
		}
	}
	if !reflect.DeepEqual(got.Nodes, want.Nodes) {
		return "nodes are missing, invented, promoted, coalesced, or non-canonical"
	}
	seenIDs = map[string]bool{}
	for _, f := range got.Flows {
		if !archifyID.MatchString(f.ID) || seenIDs[f.ID] {
			return "invalid or duplicate flow ID"
		}
		seenIDs[f.ID] = true
		prefix := map[string]string{"projection": "flow_projection_", "observation": "flow_observation_"}[f.Kind]
		if prefix == "" || f.ID != typedID(prefix, f.Identity) {
			return "invented flow kind or non-canonical flow ID"
		}
	}
	if !reflect.DeepEqual(got.Flows, want.Flows) {
		return "flows are missing, invented, redirected, reordered, or carry wrong stance"
	}
	if !reflect.DeepEqual(got.Cards, want.Cards) {
		return "cards or ordered probe text diverge from plan"
	}
	return ""
}

func loadFixtures(t *testing.T) fixtureFile {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate gate source")
	}
	path := filepath.Join(filepath.Dir(here), "..", "..", "..", "..", "lab", "bancada", "history", "benchmarks", "0069", "adversarial-fixtures.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var fixtures fixtureFile
	if err := json.Unmarshal(b, &fixtures); err != nil {
		t.Fatalf("decode fixtures: %v", err)
	}
	if fixtures.SchemaVersion != "idea-radar-archify-adversarial-v1" {
		t.Fatalf("unexpected fixture schema %q", fixtures.SchemaVersion)
	}
	return fixtures
}

func classifyCorpus(t *testing.T, fixtures fixtureFile, order []int) map[string]classification {
	t.Helper()
	base := canonicalReceipt(fixtures.PermutationRepetitions)
	results := make(map[string]classification, len(order))
	for _, index := range order {
		c := fixtures.Cases[index]
		first, reason := classify(mutate(base, c.Operation), fixtures.PermutationRepetitions)
		second, _ := classify(mutate(base, c.Operation), fixtures.PermutationRepetitions)
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

func TestStep0069DiscriminatoryGate(t *testing.T) {
	fixtures := loadFixtures(t)
	if fixtures.PermutationRepetitions != 20 {
		t.Fatalf("contract requires 20 permutation repetitions, got %d", fixtures.PermutationRepetitions)
	}
	if fixtures.RequiredMutationScore != 1.0 {
		t.Fatalf("contract requires mutation score 1.0, got %.3f", fixtures.RequiredMutationScore)
	}
	order := make([]int, len(fixtures.Cases))
	mutants := 0
	for i, c := range fixtures.Cases {
		order[i] = i
		if strings.HasPrefix(c.ID, "M") {
			mutants++
		}
	}
	if mutants < 16 {
		t.Fatalf("need at least M01-M16, got %d mutants", mutants)
	}
	baseline := classifyCorpus(t, fixtures, order)
	for repetition := 0; repetition < fixtures.PermutationRepetitions; repetition++ {
		permuted := append([]int(nil), order...)
		shift := repetition % len(permuted)
		permuted = append(permuted[shift:], permuted[:shift]...)
		if repetition%2 == 1 {
			for left, right := 0, len(permuted)-1; left < right; left, right = left+1, right-1 {
				permuted[left], permuted[right] = permuted[right], permuted[left]
			}
		}
		if got := classifyCorpus(t, fixtures, permuted); !reflect.DeepEqual(got, baseline) {
			t.Fatalf("corpus order %d changed classifications", repetition)
		}
	}
	rejected := 0
	for id, result := range baseline {
		if strings.HasPrefix(id, "M") && result == violated {
			rejected++
		}
	}
	score := float64(rejected) / float64(mutants)
	t.Logf("mutation_score=%.3f rejected=%d valid_mutations=%d required=%.3f", score, rejected, mutants, fixtures.RequiredMutationScore)
	if score != fixtures.RequiredMutationScore {
		t.Fatalf("mutation score %.3f does not meet %.3f", score, fixtures.RequiredMutationScore)
	}
	for _, c := range fixtures.Cases {
		if strings.HasPrefix(c.ID, "U") && baseline[c.ID] != unknown {
			t.Fatalf("opaque case %s was not Unknown", c.ID)
		}
		if strings.HasPrefix(c.ID, "P") && baseline[c.ID] != preserved {
			t.Fatalf("control %s was not Preserved", c.ID)
		}
	}
}

func TestFixtureIDsAreUniqueAndContiguous(t *testing.T) {
	fixtures := loadFixtures(t)
	seen := map[string]bool{}
	for _, c := range fixtures.Cases {
		if seen[c.ID] {
			t.Fatalf("duplicate fixture ID %s", c.ID)
		}
		seen[c.ID] = true
		if c.Expected != string(preserved) && c.Expected != string(violated) && c.Expected != string(unknown) {
			t.Fatalf("%s has invalid expected classification %q", c.ID, c.Expected)
		}
	}
	for i := 1; i <= 16; i++ {
		id := fmt.Sprintf("M%02d", i)
		if !seen[id] {
			t.Errorf("required mutation %s is missing", id)
		}
	}
}
