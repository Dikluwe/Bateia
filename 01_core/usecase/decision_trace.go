/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const traceMaxNodes = 1000
const traceMaxEdges = 4000
const traceMaxChanges = 100
const traceMaxJSON = 2 << 20

var traceLimitations = []string{
	"Graph and source declarations are not authenticated; fingerprints do not prove truth or source identity.",
	"Only recorded dependencies are traversed; no recorded impact does not prove absence of real impact.",
	"Review tasks do not recalculate or endorse a decision; the recorded historical decision is preserved.",
}

// traceBoundedJSON checks typed input before allocating its complete JSON.
// The conservative estimate includes even omitted fields and bounds empty vectors.
func traceBoundedJSON(value any) ([]byte, error) {
	count, estimate, text := 0, 0, 0
	var visit func(reflect.Value, int) error
	visit = func(v reflect.Value, depth int) error {
		count++
		if count > 100000 || depth > 64 {
			return fmt.Errorf("decision trace input exceeds structural budget")
		}
		estimate += 2
		if estimate > traceMaxJSON {
			return fmt.Errorf("decision trace input exceeds JSON budget")
		}
		if !v.IsValid() {
			return nil
		}
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if !v.IsNil() {
				return visit(v.Elem(), depth+1)
			}
		case reflect.String:
			s := v.String()
			if len(s) > 8192 || !utf8.ValidString(s) {
				return fmt.Errorf("decision trace string is too long or invalid UTF-8")
			}
			text += len(s)
			if text > traceMaxJSON {
				return fmt.Errorf("decision trace text exceeds budget")
			}
			b, _ := json.Marshal(s)
			estimate += len(b)
		case reflect.Slice, reflect.Array:
			if v.Len() > traceMaxEdges {
				return fmt.Errorf("decision trace vector exceeds budget")
			}
			for i := 0; i < v.Len(); i++ {
				if err := visit(v.Index(i), depth+1); err != nil {
					return err
				}
			}
		case reflect.Map:
			if v.Len() > traceMaxNodes {
				return fmt.Errorf("decision trace map exceeds budget")
			}
			iter := v.MapRange()
			for iter.Next() {
				if err := visit(iter.Key(), depth+1); err != nil {
					return err
				}
				if err := visit(iter.Value(), depth+1); err != nil {
					return err
				}
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				estimate += len(v.Type().Field(i).Name) + 4
				if err := visit(v.Field(i), depth+1); err != nil {
					return err
				}
			}
		case reflect.Float32, reflect.Float64:
			if math.IsNaN(v.Float()) || math.IsInf(v.Float(), 0) {
				return fmt.Errorf("decision trace contains non-finite number")
			}
			estimate += 32
		default:
			estimate += 32
		}
		if estimate > traceMaxJSON {
			return fmt.Errorf("decision trace input exceeds JSON budget")
		}
		return nil
	}
	if err := visit(reflect.ValueOf(value), 0); err != nil {
		return nil, err
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(b) > traceMaxJSON {
		return nil, fmt.Errorf("decision trace input exceeds JSON budget")
	}
	return b, nil
}

func traceFingerprint(value any) (string, error) {
	b, err := traceBoundedJSON(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func traceValidDigest(s string) bool {
	if len(s) != 64 || strings.ToLower(s) != s {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func traceValidID(s string) bool {
	if s == "" || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func traceValidDecision(d domain.ReuseDecision) bool {
	switch d {
	case domain.DecisionFork, domain.DecisionExtend, domain.DecisionCompose, domain.DecisionIntegrate, domain.DecisionReference, domain.DecisionBuild, domain.DecisionUnknown:
		return true
	}
	return false
}

// FreezeDecisionTrace validates and canonically seals declarations, not proof.
func FreezeDecisionTrace(graph domain.DecisionTraceGraph) (domain.DecisionTraceGraph, error) {
	fail := func(err error) (domain.DecisionTraceGraph, error) { return domain.DecisionTraceGraph{}, err }
	if len(graph.Nodes) > traceMaxNodes || len(graph.Edges) > traceMaxEdges {
		return fail(fmt.Errorf("decision trace graph exceeds budget"))
	}
	if _, err := traceBoundedJSON(graph); err != nil {
		return fail(err)
	}
	if graph.SchemaVersion != domain.DecisionTraceSchemaVersion || graph.Scope != domain.DecisionTraceScope || !traceValidDecision(graph.Decision) {
		return fail(fmt.Errorf("invalid decision trace schema, scope or decision"))
	}
	providedHash := graph.GraphHash
	if providedHash != "" && !traceValidDigest(providedHash) {
		return fail(fmt.Errorf("invalid decision trace graph hash"))
	}
	graph.Nodes = append([]domain.DecisionTraceNode{}, graph.Nodes...)
	graph.Edges = append([]domain.DecisionTraceEdge{}, graph.Edges...)
	graph.Limitations = append([]string{}, graph.Limitations...)
	for _, required := range traceLimitations {
		found := false
		for _, existing := range graph.Limitations {
			if existing == required {
				found = true
				break
			}
		}
		if !found {
			graph.Limitations = append(graph.Limitations, required)
		}
	}
	ids := make(map[string]domain.DecisionTraceNode, len(graph.Nodes))
	sourceIDs := make(map[string]string)
	decisions := 0
	for i := range graph.Nodes {
		n := &graph.Nodes[i]
		if !traceValidID(n.ID) || !traceValidDigest(n.Fingerprint) {
			return fail(fmt.Errorf("invalid node ID or fingerprint"))
		}
		if _, ok := ids[n.ID]; ok {
			return fail(fmt.Errorf("duplicate decision trace node %q", n.ID))
		}
		switch n.Kind {
		case domain.TraceInput, domain.TraceClaim:
		case domain.TraceSource:
		case domain.TraceDecision:
			decisions++
			fp, err := traceFingerprint(graph.Decision)
			if err != nil || fp != n.Fingerprint {
				return fail(fmt.Errorf("decision declaration fingerprint mismatch"))
			}
		default:
			return fail(fmt.Errorf("unknown decision trace node kind"))
		}
		if n.Source != nil {
			if n.Kind != domain.TraceSource {
				return fail(fmt.Errorf("source declaration attached to non-source node"))
			}
			copy := *n.Source
			n.Source = &copy
			fp, err := traceFingerprint(copy)
			if err != nil || fp != n.Fingerprint {
				return fail(fmt.Errorf("source declaration fingerprint mismatch"))
			}
			key, err := traceFingerprint(traceSourceKey(copy))
			if err != nil {
				return fail(err)
			}
			if _, exists := sourceIDs[key]; exists {
				return fail(fmt.Errorf("duplicate or ambiguous source identity"))
			}
			sourceIDs[key] = n.ID
		}
		ids[n.ID] = *n
	}
	if decisions != 1 {
		return fail(fmt.Errorf("decision trace must contain exactly one decision"))
	}
	adj := make(map[string][]string)
	indegree := make(map[string]int, len(ids))
	for id := range ids {
		indegree[id] = 0
	}
	seen := make(map[domain.DecisionTraceEdge]bool)
	for _, e := range graph.Edges {
		if e.From == e.To {
			return fail(fmt.Errorf("self edge in decision trace"))
		}
		if _, ok := ids[e.From]; !ok {
			return fail(fmt.Errorf("dangling dependency"))
		}
		if _, ok := ids[e.To]; !ok {
			return fail(fmt.Errorf("dangling consumer"))
		}
		switch e.Relation {
		case domain.TraceDependsOn, domain.TraceSupports, domain.TraceRefutes, domain.TraceInforms:
		default:
			return fail(fmt.Errorf("unknown decision trace relation"))
		}
		if seen[e] {
			return fail(fmt.Errorf("duplicate decision trace edge"))
		}
		seen[e] = true
		adj[e.From] = append(adj[e.From], e.To)
		indegree[e.To]++
	}
	queue := []string{}
	for id, degree := range indegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	visited := 0
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		visited++
		for _, child := range adj[id] {
			indegree[child]--
			if indegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}
	if visited != len(ids) {
		return fail(fmt.Errorf("cycle in decision trace"))
	}
	sort.Slice(graph.Nodes, func(i, j int) bool { return graph.Nodes[i].ID < graph.Nodes[j].ID })
	sort.Slice(graph.Edges, func(i, j int) bool {
		a, b := graph.Edges[i], graph.Edges[j]
		if a.From != b.From {
			return a.From < b.From
		}
		if a.To != b.To {
			return a.To < b.To
		}
		return a.Relation < b.Relation
	})
	sort.Strings(graph.Limitations)
	graph.GraphHash = ""
	hash, err := traceFingerprint(graph)
	if err != nil {
		return fail(err)
	}
	if providedHash != "" && providedHash != hash {
		return fail(fmt.Errorf("decision trace graph hash mismatch"))
	}
	graph.GraphHash = hash
	return graph, nil
}

type traceBuilder struct {
	graph    domain.DecisionTraceGraph
	nodes    map[string]domain.DecisionTraceNode
	edges    map[domain.DecisionTraceEdge]bool
	inputs   []string
	decision string
}

func (b *traceBuilder) node(kind domain.DecisionTraceNodeKind, key any, label string, record any, source *domain.EvidenceSource) (string, error) {
	id, err := traceFingerprint(struct {
		Kind domain.DecisionTraceNodeKind
		Key  any
	}{kind, key})
	if err != nil {
		return "", err
	}
	fp, err := traceFingerprint(record)
	if err != nil {
		return "", err
	}
	n := domain.DecisionTraceNode{ID: id, Kind: kind, Label: label, Fingerprint: fp, Source: source}
	if old, ok := b.nodes[id]; ok {
		if old.Fingerprint != fp {
			return "", fmt.Errorf("ambiguous declarations for trace identity %s", id)
		}
		return id, nil
	}
	if len(b.nodes) >= traceMaxNodes {
		return "", fmt.Errorf("decision trace node budget exceeded")
	}
	b.nodes[id] = n
	b.graph.Nodes = append(b.graph.Nodes, n)
	return id, nil
}

func (b *traceBuilder) edge(from, to string, relation domain.DecisionTraceRelation) error {
	e := domain.DecisionTraceEdge{From: from, To: to, Relation: relation}
	if b.edges[e] {
		return nil
	}
	if len(b.edges) >= traceMaxEdges {
		return fmt.Errorf("decision trace edge budget exceeded")
	}
	b.edges[e] = true
	b.graph.Edges = append(b.graph.Edges, e)
	return nil
}

func traceSourceKey(source domain.EvidenceSource) any {
	return struct {
		Kind            domain.EvidenceKind
		Ref             domain.CommitRef
		Path, Symbol    string
		Line            int
		RuleID, Pattern string
		MatchClass      domain.EvidenceMatchClass
	}{source.Kind, source.Ref, source.Path, source.Symbol, source.Line, source.RuleID, source.MatchedPattern, source.MatchClass}
}

func (b *traceBuilder) claim(key any, label string, record any, sources []domain.EvidenceSource) error {
	id, err := b.node(domain.TraceClaim, key, label, record, nil)
	if err != nil {
		return err
	}
	if err = b.edge(id, b.decision, domain.TraceDependsOn); err != nil {
		return err
	}
	for _, input := range b.inputs {
		if err = b.edge(input, id, domain.TraceDependsOn); err != nil {
			return err
		}
	}
	for _, source := range sources {
		// Snippet, strength and content hash are observations, not identity fields.
		key := traceSourceKey(source)
		copy := source
		sid, err := b.node(domain.TraceSource, key, "Declared source: "+source.Path, copy, &copy)
		if err != nil {
			return err
		}
		if err = b.edge(sid, id, domain.TraceInforms); err != nil {
			return err
		}
	}
	return nil
}

// BuildDecisionTrace records an existing report without running any decision rule.
func BuildDecisionTrace(input domain.DecisionTraceBuildInput) (domain.DecisionTraceGraph, error) {
	if input.SchemaVersion != "1" {
		return domain.DecisionTraceGraph{}, fmt.Errorf("unsupported trace build schema")
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return domain.DecisionTraceGraph{}, err
	}
	b := traceBuilder{graph: domain.DecisionTraceGraph{SchemaVersion: domain.DecisionTraceSchemaVersion, Decision: input.Report.Decision, Scope: domain.DecisionTraceScope, Limitations: append([]string{}, traceLimitations...)}, nodes: make(map[string]domain.DecisionTraceNode), edges: make(map[domain.DecisionTraceEdge]bool)}
	var err error
	b.decision, err = b.node(domain.TraceDecision, "recorded-decision", "Historical decision (not recalculated)", input.Report.Decision, nil)
	if err != nil {
		return domain.DecisionTraceGraph{}, err
	}
	addInput := func(key string, value any) error {
		id, err := b.node(domain.TraceInput, key, key, value, nil)
		if err != nil {
			return err
		}
		b.inputs = append(b.inputs, id)
		return b.edge(id, b.decision, domain.TraceDependsOn)
	}
	if err = addInput("report", input.Report); err != nil {
		return domain.DecisionTraceGraph{}, err
	}
	if err = addInput("policy", input.Report.Policy); err != nil {
		return domain.DecisionTraceGraph{}, err
	}
	if input.Idea != nil {
		if err = addInput("idea", input.Idea); err != nil {
			return domain.DecisionTraceGraph{}, err
		}
	}
	if input.RuleVersions != nil {
		if err = addInput("rule_versions", input.RuleVersions); err != nil {
			return domain.DecisionTraceGraph{}, err
		}
	}
	for _, row := range input.Report.Matrix {
		for _, candidate := range row.Candidates {
			key := struct{ Category, Capability, Candidate string }{"capability", row.Capability, candidate.Candidate}
			if err = b.claim(key, "Capability: "+row.Capability+" / "+candidate.Candidate, candidate, candidate.Evidence); err != nil {
				return domain.DecisionTraceGraph{}, err
			}
		}
	}
	for _, group := range []struct {
		name    string
		reasons []domain.DecisionReason
	}{{"reason", input.Report.Reasons}, {"blocker", input.Report.Blockers}} {
		for _, reason := range group.reasons {
			key := struct{ Category, RuleID, Effect, Candidate, Message string }{group.name, reason.RuleID, reason.Effect, reason.Candidate, reason.Message}
			if err = b.claim(key, group.name+": "+reason.RuleID, reason, reason.Evidence); err != nil {
				return domain.DecisionTraceGraph{}, err
			}
		}
	}
	for _, component := range input.Report.Components {
		key := struct{ Category, Name, Kind, Path string }{"component", component.Name, component.Kind, component.Path}
		if err = b.claim(key, "Component: "+component.Name, component, component.Evidence); err != nil {
			return domain.DecisionTraceGraph{}, err
		}
	}
	if input.Report.Composition != nil {
		if err = b.claim("composition-report", "Composition declaration", input.Report.Composition, nil); err != nil {
			return domain.DecisionTraceGraph{}, err
		}
		for _, composition := range input.Report.Composition.Candidates {
			key := struct{ Category, ID string }{"composition", composition.ID}
			sources := append([]domain.EvidenceSource{}, composition.Evidence...)
			for _, member := range composition.Members {
				sources = append(sources, member.Component.Evidence...)
			}
			if err = b.claim(key, "Composition: "+composition.ID, composition, sources); err != nil {
				return domain.DecisionTraceGraph{}, err
			}
		}
	}
	if input.Report.TechnicalEvidence != nil {
		if err = b.claim("technical-report", "Consultative technical review", input.Report.TechnicalEvidence, nil); err != nil {
			return domain.DecisionTraceGraph{}, err
		}
		for _, claim := range input.Report.TechnicalEvidence.Claims {
			key := struct{ Category, Candidate, Claim string }{"technical-claim", claim.Candidate, claim.Claim}
			if err = b.claim(key, "Technical declaration: "+claim.Claim, claim, nil); err != nil {
				return domain.DecisionTraceGraph{}, err
			}
		}
		for _, verification := range input.Report.TechnicalEvidence.Verifications {
			key := struct{ Category, Finding, Candidate, Predicate, Value string }{"technical-verification", verification.FindingID, verification.Candidate, string(verification.Predicate), verification.Value}
			sources := []domain.EvidenceSource{}
			if verification.Source != nil {
				sources = append(sources, *verification.Source)
			}
			if err = b.claim(key, "Technical verifier declaration: "+verification.FindingID, verification, sources); err != nil {
				return domain.DecisionTraceGraph{}, err
			}
		}
	}
	return FreezeDecisionTrace(b.graph)
}

// PlanDecisionTraceImpact traverses all registered relations, including REFUTES.
// It intentionally preserves the historical verdict and does not mutate the graph.
func PlanDecisionTraceImpact(input domain.DecisionTraceImpactInput) (domain.DecisionTraceImpactPlan, error) {
	fail := func(err error) (domain.DecisionTraceImpactPlan, error) { return domain.DecisionTraceImpactPlan{}, err }
	if input.SchemaVersion != "1" || len(input.Changes) > traceMaxChanges {
		return fail(fmt.Errorf("invalid impact schema or change budget"))
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return fail(err)
	}
	graph, err := FreezeDecisionTrace(input.Graph)
	if err != nil {
		return fail(err)
	}
	if input.Graph.GraphHash == "" {
		return fail(fmt.Errorf("impact requires a frozen graph hash"))
	}
	nodes := make(map[string]domain.DecisionTraceNode, len(graph.Nodes))
	adj := make(map[string][]string)
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	for _, edge := range graph.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}
	changes := make(map[string]string)
	seen := make(map[string]bool)
	for _, change := range input.Changes {
		node, ok := nodes[change.NodeID]
		if !ok || seen[change.NodeID] {
			return fail(fmt.Errorf("unknown or duplicate changed node"))
		}
		seen[change.NodeID] = true
		switch change.Reason {
		case "CONTENT_CHANGED", "EXPIRED", "REVOKED", "POLICY_CHANGED", "RULE_CHANGED", "INPUT_CHANGED", "JUDGMENT_CHANGED":
		default:
			return fail(fmt.Errorf("unknown change reason"))
		}
		if change.NewFingerprint != "" && !traceValidDigest(change.NewFingerprint) {
			return fail(fmt.Errorf("invalid new fingerprint"))
		}
		if change.Reason == "CONTENT_CHANGED" && change.NewFingerprint == node.Fingerprint {
			continue
		}
		changes[change.NodeID] = change.Reason
	}
	plan := domain.DecisionTraceImpactPlan{SchemaVersion: "decision-trace-impact-v1", GraphHash: graph.GraphHash, Status: "NO_RECORDED_IMPACT", Scope: domain.DecisionTraceScope, DecisionBefore: graph.Decision, DecisionAfter: graph.Decision, DecisionImpact: "NONE", ChangedNodeIDs: []string{}, ImpactedNodeIDs: []string{}, UnaffectedNodeIDs: []string{}, Tasks: []domain.DecisionTraceReviewTask{}, Limitations: append([]string{}, traceLimitations...)}
	for id := range changes {
		plan.ChangedNodeIDs = append(plan.ChangedNodeIDs, id)
	}
	sort.Strings(plan.ChangedNodeIDs)
	ancestors := make(map[string][]string)
	// Per-change visited sets bound traversal to changes*(nodes+edges).
	for _, changed := range plan.ChangedNodeIDs {
		visited := make(map[string]bool)
		queue := []string{changed}
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			if visited[id] {
				continue
			}
			visited[id] = true
			ancestors[id] = append(ancestors[id], changed)
			for _, child := range adj[id] {
				if !visited[child] {
					queue = append(queue, child)
				}
			}
		}
	}
	for _, node := range graph.Nodes {
		roots := ancestors[node.ID]
		if len(roots) == 0 {
			plan.UnaffectedNodeIDs = append(plan.UnaffectedNodeIDs, node.ID)
			continue
		}
		plan.ImpactedNodeIDs = append(plan.ImpactedNodeIDs, node.ID)
		reason := "RECORDED_DEPENDENCY_CHANGED"
		if own, ok := changes[node.ID]; ok {
			reason = own
		}
		plan.Tasks = append(plan.Tasks, domain.DecisionTraceReviewTask{NodeID: node.ID, Kind: node.Kind, ChangedAncestorIDs: roots, Reason: reason})
	}
	if len(plan.ImpactedNodeIDs) > 0 {
		plan.Status = "REVIEW_REQUIRED"
	}
	if _, err := traceBoundedJSON(plan); err != nil {
		return fail(fmt.Errorf("decision trace impact output exceeds budget: %w", err))
	}
	return plan, nil
}
