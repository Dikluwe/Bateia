/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/idea_radar_archify_test.md
 * @layer L2
 * @updated 2026-09-12
 */

package report

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type archifyDocument0069 struct {
	SchemaVersion int                        `json:"schema_version"`
	DiagramType   string                     `json:"diagram_type"`
	Meta          map[string]any             `json:"meta"`
	Stages        []any                      `json:"stages"`
	Nodes         []map[string]any           `json:"nodes"`
	Flows         []map[string]any           `json:"flows"`
	Cards         []map[string]any           `json:"cards"`
	Raw           map[string]json.RawMessage `json:"-"`
}

func TestExportIdeaRadarArchifyGoldenShapeAndTopology(t *testing.T) {
	input := convergentRadarInput0069()
	wantInput := cloneRadarInput0069(input)

	got, err := ExportIdeaRadarArchify(input)
	if err != nil {
		t.Fatalf("ExportIdeaRadarArchify: %v", err)
	}
	if !reflect.DeepEqual(input, wantInput) {
		t.Fatalf("input mutated\nbefore: %#v\nafter:  %#v", wantInput, input)
	}
	if !bytes.HasSuffix(got, []byte("\n")) || bytes.HasSuffix(got, []byte("\n\n")) {
		t.Fatalf("document must have exactly one final newline: %q", got)
	}
	if bytes.Contains(got, []byte("\t")) {
		t.Fatal("document must use two-space indentation, not tabs")
	}

	doc := decodeArchify0069(t, got)
	if doc.SchemaVersion != 1 || doc.DiagramType != "dataflow" {
		t.Fatalf("not Dataflow v1: schema_version=%d diagram_type=%q", doc.SchemaVersion, doc.DiagramType)
	}
	if title, _ := doc.Meta["title"].(string); title != "Bateia Idea Radar" {
		t.Fatalf("meta.title = %q", title)
	}
	assertStages0069(t, doc.Stages, []string{"Origins", "Perspectives", "Targets"})
	if len(doc.Nodes) != 7 || len(doc.Flows) != 6 || len(doc.Cards) != 1 {
		t.Fatalf("golden cardinality: nodes=%d flows=%d cards=%d", len(doc.Nodes), len(doc.Flows), len(doc.Cards))
	}

	nodes := indexByID0069(t, doc.Nodes, "node")
	flows := indexByID0069(t, doc.Flows, "flow")
	for id, node := range nodes {
		if node["type"] != "external" {
			t.Errorf("node %q type = %#v, want external", id, node["type"])
		}
		stage := integer0069(t, node["stage"], "node.stage")
		if stage < 0 || stage > 2 {
			t.Errorf("node %q has invalid stage %d", id, stage)
		}
		_ = integer0069(t, node["row"], "node.row")
		assertNodeLabelFits0069(t, id, node)
	}

	wantNodeIDs := []string{
		archifyID0069("origin_", "origin-a"),
		archifyID0069("origin_", "origin-b"),
		archifyID0069("origin_", "origin-c"),
		archifyID0069("perspective_", "p-architecture"),
		archifyID0069("perspective_", "p-behavior"),
		archifyID0069("perspective_", "p-problem"),
		archifyID0069("target_", "candidate#concept"),
	}
	assertOrderedIDs0069(t, doc.Nodes, wantNodeIDs, "golden nodes")
	for _, id := range wantNodeIDs {
		if nodes[id] == nil {
			t.Errorf("golden node %q missing", id)
		}
	}
	wantLabels := map[string]string{
		archifyID0069("origin_", "origin-a"):            "O origin_a",
		archifyID0069("origin_", "origin-b"):            "O origin_b",
		archifyID0069("origin_", "origin-c"):            "O origin_c",
		archifyID0069("perspective_", "p-architecture"): "P ARCH",
		archifyID0069("perspective_", "p-behavior"):     "P BEHV",
		archifyID0069("perspective_", "p-problem"):      "P PROB",
		archifyID0069("target_", "candidate#concept"):   "T candida~",
	}
	for id, wantLabel := range wantLabels {
		if gotLabel, _ := nodes[id]["label"].(string); gotLabel != wantLabel {
			t.Errorf("node %q label = %q, want %q", id, gotLabel, wantLabel)
		}
	}
	wantFlowIDs := make([]string, 0, 6)
	for _, identity := range []string{"p-architecture", "p-behavior", "p-problem"} {
		wantFlowIDs = append(wantFlowIDs, archifyID0069("flow_projection_", identity))
		if flows[archifyID0069("flow_projection_", identity)] == nil {
			t.Errorf("projection flow for %q missing", identity)
		}
	}
	for _, identity := range []string{"o-architecture", "o-behavior", "o-problem"} {
		wantFlowIDs = append(wantFlowIDs, archifyID0069("flow_observation_", identity))
		if flows[archifyID0069("flow_observation_", identity)] == nil {
			t.Errorf("observation flow for %q missing", identity)
		}
	}
	assertOrderedIDs0069(t, doc.Flows, wantFlowIDs, "golden flows")
	assertOrderedTitles0069(t, doc.Cards, []string{"candidate#concept"}, "golden cards")
	assertOnlyRadarTopology0069(t, nodes, doc.Flows)
	for _, identity := range []string{"o-architecture", "o-behavior", "o-problem"} {
		flow := flows[archifyID0069("flow_observation_", identity)]
		if !strings.Contains(strings.ToUpper(flattenStrings0069(flow)), "SUPPORT") {
			t.Errorf("observation flow %q omits stance: %#v", identity, flow)
		}
	}

	cardText := strings.ToUpper(flattenStrings0069(doc.Cards[0]))
	for _, fragment := range []string{"CONVERGENT", "IDEA_CONVERGENT", "BIND_EVIDENCE", "BIND_PHYSICAL_EVIDENCE"} {
		if !strings.Contains(cardText, fragment) {
			t.Errorf("card omits %q: %s", fragment, cardText)
		}
	}
	for _, flow := range doc.Flows {
		flowText := strings.ToUpper(flattenStrings0069(flow))
		if strings.Contains(flowText, "PROBE") || strings.Contains(flowText, "BIND_EVIDENCE") {
			t.Errorf("probe leaked into flow: %#v", flow)
		}
	}
}

func TestExportIdeaRadarArchifyRendererBudgetBoundary(t *testing.T) {
	boundary := radarGridInput0069(5, 5, 5)
	encoded, err := ExportIdeaRadarArchify(boundary)
	if err != nil {
		t.Fatalf("5/5/5 renderer boundary must remain representable: %v", err)
	}
	doc := decodeArchify0069(t, encoded)
	if len(doc.Nodes) != 15 || len(doc.Flows) != 30 || len(doc.Cards) != 5 {
		t.Fatalf("5/5/5 boundary cardinality: nodes=%d flows=%d cards=%d", len(doc.Nodes), len(doc.Flows), len(doc.Cards))
	}
	rowsByStage := map[int][]int{}
	for _, node := range doc.Nodes {
		stage := integer0069(t, node["stage"], "node.stage")
		rowsByStage[stage] = append(rowsByStage[stage], integer0069(t, node["row"], "node.row"))
		id, _ := node["id"].(string)
		assertNodeLabelFits0069(t, id, node)
	}
	for stage := 0; stage < 3; stage++ {
		if want := []int{0, 1, 2, 3, 4}; !reflect.DeepEqual(rowsByStage[stage], want) {
			t.Errorf("stage %d rows = %v, want %v", stage, rowsByStage[stage], want)
		}
	}

	overflows := map[string]domain.IdeaRadarInput{
		"six perspectives":   radarGridInput0069(6, 1, 5),
		"six unique origins": radarGridInput0069(6, 1, 6),
		"six targets":        radarGridInput0069(5, 6, 5),
	}
	for name, input := range overflows {
		t.Run(name, func(t *testing.T) {
			got, err := ExportIdeaRadarArchify(input)
			if err == nil || len(got) != 0 {
				t.Fatalf("renderer overflow must fail without bytes: err=%v bytes=%q", err, got)
			}
			if !strings.Contains(err.Error(), "render_budget_exceeded") {
				t.Fatalf("renderer overflow reason = %q, want render_budget_exceeded", err)
			}
		})
	}
}

func TestExportIdeaRadarArchifyShortLabelsPreserveFullIdentityElsewhere(t *testing.T) {
	input := domain.IdeaRadarInput{
		SchemaVersion: domain.IdeaRadarInputSchemaVersion,
		Projections: []domain.IdeaRadarProjection{
			{ID: "problem:projection/full", Kind: domain.IdeaProjectionProblem, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin:/problem/full"},
			{ID: "behavior:projection/full", Kind: domain.IdeaProjectionBehavior, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin:/behavior/full"},
			{ID: "architecture:projection/full", Kind: domain.IdeaProjectionArchitecture, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin:/architecture/full"},
		},
	}
	targetID := "candidate:full/value#concept/full:value"
	for index, projection := range input.Projections {
		input.Observations = append(input.Observations, domain.IdeaRadarObservation{
			ID: fmt.Sprintf("observation:%d/full/value", index), ProjectionID: projection.ID,
			CandidateID: "candidate:full/value", ConceptID: "concept/full:value", Stance: domain.IdeaRadarStanceSupport,
		})
	}

	encoded, err := ExportIdeaRadarArchify(input)
	if err != nil {
		t.Fatal(err)
	}
	doc := decodeArchify0069(t, encoded)
	nodes := indexByID0069(t, doc.Nodes, "node")
	flows := indexByID0069(t, doc.Flows, "flow")
	for _, node := range doc.Nodes {
		id, _ := node["id"].(string)
		assertNodeLabelFits0069(t, id, node)
	}
	for index, projection := range input.Projections {
		originNodeID := archifyID0069("origin_", projection.OriginID)
		if nodes[originNodeID] == nil {
			t.Errorf("full origin identity absent from ID: %q", projection.OriginID)
		}
		perspectiveNodeID := archifyID0069("perspective_", projection.ID)
		if nodes[perspectiveNodeID] == nil {
			t.Errorf("full perspective identity absent from ID: %q", projection.ID)
		}
		projectionFlow := flows[archifyID0069("flow_projection_", projection.ID)]
		projectionText := flattenStrings0069(projectionFlow)
		for _, identity := range []string{projection.ID, projection.OriginID} {
			if !strings.Contains(projectionText, identity) {
				t.Errorf("projection flow label omits full identity %q: %#v", identity, projectionFlow)
			}
		}
		observationID := fmt.Sprintf("observation:%d/full/value", index)
		observationFlow := flows[archifyID0069("flow_observation_", observationID)]
		if !strings.Contains(flattenStrings0069(observationFlow), observationID) {
			t.Errorf("observation flow label omits full identity %q: %#v", observationID, observationFlow)
		}
	}
	targetNodeID := archifyID0069("target_", targetID)
	if nodes[targetNodeID] == nil {
		t.Errorf("full target identity absent from ID: %q", targetID)
	}
	assertOrderedTitles0069(t, doc.Cards, []string{targetID}, "long identity cards")
	cardText := flattenStrings0069(doc.Cards[0])
	for index := range input.Observations {
		if !strings.Contains(cardText, input.Observations[index].ID) {
			t.Errorf("card provenance omits observation %q: %s", input.Observations[index].ID, cardText)
		}
	}
}

func TestExportIdeaRadarArchifyCanonicalDedupAndInjectiveIDs(t *testing.T) {
	input := convergentRadarInput0069()
	input.Observations = append(input.Observations,
		domain.IdeaRadarObservation{ID: "z-duplicate", ProjectionID: "p-problem", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
	)
	base, err := ExportIdeaRadarArchify(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Observations[len(input.Observations)-1].ID = "a-duplicate"
	changed, err := ExportIdeaRadarArchify(input)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(base, changed) {
		t.Fatal("canonical retained observation ID must be the lexicographically smallest semantic duplicate")
	}
	doc := decodeArchify0069(t, changed)
	if len(doc.Flows) != 6 {
		t.Fatalf("semantic duplicate created a flow: %d", len(doc.Flows))
	}
	ids := indexByID0069(t, doc.Flows, "flow")
	if ids[archifyID0069("flow_observation_", "a-duplicate")] == nil {
		t.Fatal("flow does not use canonical retained observation ID")
	}

	sharedOrigin := convergentRadarInput0069()
	sharedOrigin.Projections[1].OriginID = sharedOrigin.Projections[0].OriginID
	sharedEncoded, err := ExportIdeaRadarArchify(sharedOrigin)
	if err != nil {
		t.Fatal(err)
	}
	sharedDoc := decodeArchify0069(t, sharedEncoded)
	if len(sharedDoc.Nodes) != 6 || len(sharedDoc.Flows) != 6 {
		t.Fatalf("recorded origins were not deduplicated by exact OriginID: nodes=%d flows=%d", len(sharedDoc.Nodes), len(sharedDoc.Flows))
	}

	identityInput := identityRadarInput0069()
	encoded, err := ExportIdeaRadarArchify(identityInput)
	if err != nil {
		t.Fatal(err)
	}
	identityDoc := decodeArchify0069(t, encoded)
	targets := make([]string, 0, 2)
	for _, node := range identityDoc.Nodes {
		id, _ := node["id"].(string)
		if strings.HasPrefix(id, "target_") {
			targets = append(targets, id)
		}
	}
	sort.Strings(targets)
	wantTargets := []string{archifyID0069("target_", "a/b#concept"), archifyID0069("target_", "a:b#concept")}
	sort.Strings(wantTargets)
	if !reflect.DeepEqual(targets, wantTargets) {
		t.Fatalf("lossy IDs for a:b vs a/b: got %v want %v", targets, wantTargets)
	}
	for _, collection := range [][]map[string]any{identityDoc.Nodes, identityDoc.Flows} {
		seen := map[string]bool{}
		for _, item := range collection {
			id, _ := item["id"].(string)
			if id == "" || !validArchifyID0069(id) || seen[id] {
				t.Fatalf("invalid or duplicate Archify ID %q", id)
			}
			seen[id] = true
		}
	}
	assertOrderedTitles0069(t, identityDoc.Cards, []string{"a/b#concept", "a:b#concept"}, "identity cards")
}

func TestExportIdeaRadarArchifyUnknownOriginsAreProjectionScoped(t *testing.T) {
	input := convergentRadarInput0069()
	input.Projections[0].OriginStatus, input.Projections[0].OriginID = domain.EvidenceOriginUnknown, ""
	input.Projections[1].OriginStatus, input.Projections[1].OriginID = domain.EvidenceOriginUnknown, ""

	encoded, err := ExportIdeaRadarArchify(input)
	if err != nil {
		t.Fatal(err)
	}
	doc := decodeArchify0069(t, encoded)
	nodes := indexByID0069(t, doc.Nodes, "node")
	unknownIDs := []string{
		archifyID0069("origin_unknown_", "p-problem"),
		archifyID0069("origin_unknown_", "p-behavior"),
	}
	for _, id := range unknownIDs {
		node := nodes[id]
		if node == nil || !strings.Contains(strings.ToUpper(flattenStrings0069(node)), "UNKNOWN") {
			t.Errorf("missing descriptive projection-scoped unknown origin %q: %#v", id, node)
		}
	}
	if len(nodes) != 7 {
		t.Fatalf("unknown origins were coalesced or physical origin fabricated: nodes=%d", len(nodes))
	}
	assertOnlyRadarTopology0069(t, nodes, doc.Flows)
	cardText := strings.ToUpper(flattenStrings0069(doc.Cards[0]))
	for _, fragment := range []string{"UNKNOWN_ORIGIN", "ORIGIN_UNRECORDED", "RECORD_ORIGIN"} {
		if !strings.Contains(cardText, fragment) {
			t.Errorf("unknown-origin card omits %q: %s", fragment, cardText)
		}
	}
}

func TestExportIdeaRadarArchifyDeterministicNilEmptyAndDeepImmutable(t *testing.T) {
	canonical := convergentRadarInput0069()
	canonical.Projections[0].Terms = nil
	canonical.Projections[0].Actions = []string{}
	canonical.Projections[0].Objects = nil
	canonical.Projections[0].Required = []string{}
	wantInput := cloneRadarInput0069(canonical)
	want, err := ExportIdeaRadarArchify(canonical)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(canonical, wantInput) {
		t.Fatal("export mutated a nested input slice")
	}

	emptyVariant := cloneRadarInput0069(canonical)
	emptyVariant.Projections[0].Terms = []string{}
	emptyVariant.Projections[0].Actions = nil
	emptyVariant.Projections[0].Objects = []string{}
	emptyVariant.Projections[0].Required = nil
	got, err := ExportIdeaRadarArchify(emptyVariant)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("nil and empty non-semantic text slices changed output bytes")
	}
	canonicalVariant := cloneRadarInput0069(canonical)
	for index := range canonicalVariant.Projections {
		canonicalVariant.Projections[index].ID = " \t" + strings.ToUpper(canonicalVariant.Projections[index].ID) + "\n"
	}
	for index := range canonicalVariant.Observations {
		canonicalVariant.Observations[index].ID = " " + strings.ToUpper(canonicalVariant.Observations[index].ID) + " "
		canonicalVariant.Observations[index].ProjectionID = "\t" + strings.ToUpper(canonicalVariant.Observations[index].ProjectionID) + "\r"
		canonicalVariant.Observations[index].CandidateID = " CANDIDATE "
		canonicalVariant.Observations[index].ConceptID = " CONCEPT "
	}
	canonicalVariant.Projections[0].Terms = []string{""}
	canonicalVariantBytes, err := ExportIdeaRadarArchify(canonicalVariant)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(canonicalVariantBytes, want) {
		t.Fatal("0066-equivalent casing, edge whitespace, and empty descriptive text changed bytes")
	}

	for repetition := 0; repetition < 20; repetition++ {
		permuted := cloneRadarInput0069(canonical)
		rotateProjections0069(permuted.Projections, repetition)
		reverseObservations0069(permuted.Observations, repetition%2 == 0)
		before := cloneRadarInput0069(permuted)
		got, err := ExportIdeaRadarArchify(permuted)
		if err != nil {
			t.Fatalf("permutation %d: %v", repetition, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("permutation %d changed bytes\nwant:\n%s\ngot:\n%s", repetition, want, got)
		}
		if !reflect.DeepEqual(permuted, before) {
			t.Fatalf("permutation %d mutated deeply reachable input", repetition)
		}
	}
}

func TestExportIdeaRadarArchifyFailsClosed(t *testing.T) {
	tooMany := convergentRadarInput0069()
	tooMany.Projections = make([]domain.IdeaRadarProjection, 65)
	for index := range tooMany.Projections {
		tooMany.Projections[index] = domain.IdeaRadarProjection{
			ID: fmt.Sprintf("projection-%02d", index), Kind: domain.IdeaProjectionProblem,
			OriginStatus: domain.EvidenceOriginRecorded, OriginID: fmt.Sprintf("origin-%02d", index),
		}
	}
	cases := map[string]domain.IdeaRadarInput{
		"empty":          {},
		"invalid schema": func() domain.IdeaRadarInput { v := convergentRadarInput0069(); v.SchemaVersion = "opaque"; return v }(),
		"budget":         tooMany,
		"ambiguous": func() domain.IdeaRadarInput {
			v := convergentRadarInput0069()
			v.Observations = append(v.Observations, domain.IdeaRadarObservation{ID: "o-problem", ProjectionID: "p-problem", CandidateID: "other", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport})
			return v
		}(),
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := ExportIdeaRadarArchify(input)
			if err == nil || len(got) != 0 {
				t.Fatalf("want error and zero partial document, got err=%v bytes=%q", err, got)
			}
		})
	}
}

func convergentRadarInput0069() domain.IdeaRadarInput {
	return domain.IdeaRadarInput{
		SchemaVersion: domain.IdeaRadarInputSchemaVersion,
		Projections: []domain.IdeaRadarProjection{
			{ID: "p-problem", Kind: domain.IdeaProjectionProblem, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-a", Terms: []string{"term"}},
			{ID: "p-behavior", Kind: domain.IdeaProjectionBehavior, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-b", Actions: []string{"action"}},
			{ID: "p-architecture", Kind: domain.IdeaProjectionArchitecture, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-c", Objects: []string{"object"}},
		},
		Observations: []domain.IdeaRadarObservation{
			{ID: "o-problem", ProjectionID: "p-problem", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
			{ID: "o-behavior", ProjectionID: "p-behavior", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
			{ID: "o-architecture", ProjectionID: "p-architecture", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
		},
	}
}

func identityRadarInput0069() domain.IdeaRadarInput {
	input := convergentRadarInput0069()
	input.Observations = nil
	for _, candidate := range []string{"a:b", "a/b"} {
		for index, projection := range input.Projections {
			input.Observations = append(input.Observations, domain.IdeaRadarObservation{
				ID:           fmt.Sprintf("identity-%d-%s", index, strings.ReplaceAll(candidate, ":", "colon")),
				ProjectionID: projection.ID, CandidateID: candidate, ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport,
			})
		}
	}
	return input
}

func radarGridInput0069(perspectives, targets, uniqueOrigins int) domain.IdeaRadarInput {
	kinds := []domain.IdeaProjectionKind{
		domain.IdeaProjectionProblem,
		domain.IdeaProjectionBehavior,
		domain.IdeaProjectionArchitecture,
		domain.IdeaProjectionConstraint,
		domain.IdeaProjectionOperation,
	}
	input := domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion}
	for perspective := 0; perspective < perspectives; perspective++ {
		origin := perspective
		if uniqueOrigins > 0 {
			origin %= uniqueOrigins
		}
		input.Projections = append(input.Projections, domain.IdeaRadarProjection{
			ID: fmt.Sprintf("perspective-%02d", perspective), Kind: kinds[perspective%len(kinds)],
			OriginStatus: domain.EvidenceOriginRecorded, OriginID: fmt.Sprintf("origin-%02d", origin),
		})
	}
	for target := 0; target < targets; target++ {
		for perspective := 0; perspective < perspectives; perspective++ {
			input.Observations = append(input.Observations, domain.IdeaRadarObservation{
				ID:           fmt.Sprintf("observation-%02d-%02d", target, perspective),
				ProjectionID: fmt.Sprintf("perspective-%02d", perspective),
				CandidateID:  fmt.Sprintf("candidate-%02d", target), ConceptID: "concept",
				Stance: domain.IdeaRadarStanceSupport,
			})
		}
	}
	return input
}

func cloneRadarInput0069(input domain.IdeaRadarInput) domain.IdeaRadarInput {
	clone := input
	clone.Projections = append([]domain.IdeaRadarProjection(nil), input.Projections...)
	for index := range clone.Projections {
		clone.Projections[index].Terms = appendPreservingNil0069(input.Projections[index].Terms)
		clone.Projections[index].Actions = appendPreservingNil0069(input.Projections[index].Actions)
		clone.Projections[index].Objects = appendPreservingNil0069(input.Projections[index].Objects)
		clone.Projections[index].Required = appendPreservingNil0069(input.Projections[index].Required)
	}
	clone.Observations = append([]domain.IdeaRadarObservation(nil), input.Observations...)
	return clone
}

func appendPreservingNil0069(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

func decodeArchify0069(t *testing.T, encoded []byte) archifyDocument0069 {
	t.Helper()
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	wantKeys := []string{"cards", "diagram_type", "flows", "meta", "nodes", "schema_version", "stages"}
	gotKeys := make([]string, 0, len(raw))
	for key := range raw {
		gotKeys = append(gotKeys, key)
	}
	sort.Strings(gotKeys)
	if !reflect.DeepEqual(gotKeys, wantKeys) {
		t.Fatalf("top-level schema keys = %v, want %v", gotKeys, wantKeys)
	}
	var doc archifyDocument0069
	if err := json.Unmarshal(encoded, &doc); err != nil {
		t.Fatal(err)
	}
	for index, card := range doc.Cards {
		cardKeys := make([]string, 0, len(card))
		for key := range card {
			cardKeys = append(cardKeys, key)
		}
		sort.Strings(cardKeys)
		if want := []string{"dot", "items", "title"}; !reflect.DeepEqual(cardKeys, want) {
			t.Fatalf("card %d schema keys = %v, want %v", index, cardKeys, want)
		}
	}
	doc.Raw = raw
	return doc
}

func assertStages0069(t *testing.T, stages []any, want []string) {
	t.Helper()
	got := make([]string, len(stages))
	for index, stage := range stages {
		switch value := stage.(type) {
		case string:
			got[index] = value
		case map[string]any:
			for _, key := range []string{"name", "label", "title"} {
				if text, ok := value[key].(string); ok {
					got[index] = text
					break
				}
			}
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stages = %v, want %v", got, want)
	}
}

func indexByID0069(t *testing.T, values []map[string]any, kind string) map[string]map[string]any {
	t.Helper()
	result := make(map[string]map[string]any, len(values))
	for _, value := range values {
		id, _ := value["id"].(string)
		if id == "" || result[id] != nil {
			t.Fatalf("%s has missing/duplicate ID %q", kind, id)
		}
		result[id] = value
	}
	return result
}

func assertOrderedIDs0069(t *testing.T, values []map[string]any, want []string, label string) {
	t.Helper()
	got := make([]string, len(values))
	for index, value := range values {
		got[index], _ = value["id"].(string)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s IDs = %v, want %v", label, got, want)
	}
}

func assertOrderedTitles0069(t *testing.T, values []map[string]any, want []string, label string) {
	t.Helper()
	got := make([]string, len(values))
	seen := make(map[string]bool, len(values))
	for index, value := range values {
		if _, hasID := value["id"]; hasID {
			t.Fatalf("%s card %d has forbidden Dataflow v1 property id: %#v", label, index, value)
		}
		got[index], _ = value["title"].(string)
		if got[index] == "" || seen[got[index]] {
			t.Fatalf("%s has missing/duplicate title %q", label, got[index])
		}
		seen[got[index]] = true
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s titles = %v, want %v", label, got, want)
	}
}

func assertNodeLabelFits0069(t *testing.T, id string, node map[string]any) {
	t.Helper()
	label, ok := node["label"].(string)
	if !ok || label == "" {
		t.Fatalf("node %q lacks label: %#v", id, node)
	}
	if len(label) > 10 {
		t.Errorf("node %q label exceeds 10 bytes: %q (%d)", id, label, len(label))
	}
	for index := 0; index < len(label); index++ {
		if label[index] > 0x7f {
			t.Errorf("node %q label is not ASCII: %q", id, label)
			break
		}
	}
}

func assertOnlyRadarTopology0069(t *testing.T, nodes map[string]map[string]any, flows []map[string]any) {
	t.Helper()
	projectionFlows, observationFlows := 0, 0
	for _, flow := range flows {
		from := endpoint0069(flow, "from", "source", "source_id")
		to := endpoint0069(flow, "to", "target", "target_id")
		fromNode, toNode := nodes[from], nodes[to]
		if fromNode == nil || toNode == nil {
			t.Fatalf("flow references absent node: %#v", flow)
		}
		fromStage := integer0069(t, fromNode["stage"], "from.stage")
		toStage := integer0069(t, toNode["stage"], "to.stage")
		if !((fromStage == 0 && toStage == 1) || (fromStage == 1 && toStage == 2)) {
			t.Errorf("invented topology %d -> %d in %#v", fromStage, toStage, flow)
		} else if fromStage == 0 {
			projectionFlows++
		} else {
			observationFlows++
		}
	}
	if projectionFlows != 3 || observationFlows != 3 {
		t.Errorf("golden topology counts = %d origin->perspective, %d perspective->target", projectionFlows, observationFlows)
	}
}

func endpoint0069(value map[string]any, keys ...string) string {
	for _, key := range keys {
		if text, ok := value[key].(string); ok {
			return text
		}
	}
	return ""
}

func integer0069(t *testing.T, value any, field string) int {
	t.Helper()
	number, ok := value.(float64)
	if !ok || number != float64(int(number)) {
		t.Fatalf("%s is not an integer: %#v", field, value)
	}
	return int(number)
}

func archifyID0069(prefix, identity string) string {
	return prefix + hex.EncodeToString([]byte(identity))
}

func validArchifyID0069(value string) bool {
	if value == "" || !((value[0] >= 'a' && value[0] <= 'z') || (value[0] >= 'A' && value[0] <= 'Z')) {
		return false
	}
	for index := 1; index < len(value); index++ {
		char := value[index]
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func flattenStrings0069(value any) string {
	parts := make([]string, 0)
	var walk func(any)
	walk = func(current any) {
		switch typed := current.(type) {
		case string:
			parts = append(parts, typed)
		case []any:
			for _, item := range typed {
				walk(item)
			}
		case map[string]any:
			keys := make([]string, 0, len(typed))
			for key := range typed {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				walk(typed[key])
			}
		}
	}
	walk(value)
	return strings.Join(parts, " ")
}

func rotateProjections0069(values []domain.IdeaRadarProjection, amount int) {
	if len(values) == 0 {
		return
	}
	amount %= len(values)
	rotated := append(append([]domain.IdeaRadarProjection{}, values[amount:]...), values[:amount]...)
	copy(values, rotated)
}

func reverseObservations0069(values []domain.IdeaRadarObservation, reverse bool) {
	if !reverse {
		return
	}
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}
