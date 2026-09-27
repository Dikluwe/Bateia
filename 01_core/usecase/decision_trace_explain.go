/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/decision_trace_explain.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"bateia/01_core/domain"
)

const traceStatementMaxBytes = 1024

// tracePointer escapes RFC 6901 tokens even though the current model uses
// fixed field names and array indexes, never user-provided keys as pointers.
func tracePointer(tokens ...string) string {
	if len(tokens) == 0 {
		return ""
	}
	escaped := make([]string, len(tokens))
	for i, token := range tokens {
		escaped[i] = strings.ReplaceAll(strings.ReplaceAll(token, "~", "~0"), "/", "~1")
	}
	return "/" + strings.Join(escaped, "/")
}

func traceStatement(statement string) (string, bool) {
	if len(statement) <= traceStatementMaxBytes {
		return statement, false
	}
	end := traceStatementMaxBytes
	for !utf8.ValidString(statement[:end]) {
		end--
	}
	return statement[:end], true
}

type traceExplanationIndex struct {
	nodes   map[string]domain.DecisionTraceNode
	entries map[string]*domain.DecisionTraceExplanationEntry
}

func (index *traceExplanationIndex) occurrence(kind domain.DecisionTraceNodeKind, key any, record any, document, pointer, statement string) error {
	id, err := traceFingerprint(struct {
		Kind domain.DecisionTraceNodeKind
		Key  any
	}{kind, key})
	if err != nil {
		return err
	}
	fp, err := traceFingerprint(record)
	if err != nil {
		return err
	}
	node, ok := index.nodes[id]
	if !ok || node.Kind != kind || node.Fingerprint != fp {
		return fmt.Errorf("explanation occurrence does not bind graph node")
	}
	entry := index.entries[id]
	if entry == nil {
		text, truncated := traceStatement(statement)
		entry = &domain.DecisionTraceExplanationEntry{NodeID: id, Kind: node.Kind, Label: node.Label, Fingerprint: fp, Statement: text, StatementTruncated: truncated, Locations: []domain.DecisionTraceLocation{}}
		index.entries[id] = entry
	}
	location := domain.DecisionTraceLocation{Document: document, Pointer: pointer}
	for _, old := range entry.Locations {
		if old == location {
			return nil
		}
	}
	entry.Locations = append(entry.Locations, location)
	return nil
}

func (index *traceExplanationIndex) sources(sources []domain.EvidenceSource, pointer string) error {
	for i, source := range sources {
		if err := index.occurrence(domain.TraceSource, traceSourceKey(source), source, "report", pointer+"/"+strconv.Itoa(i), "Declared source: "+source.Path); err != nil {
			return err
		}
	}
	return nil
}

// ExplainDecisionTrace binds a frozen historical graph to its exact declared
// build input, then indexes record locations. It neither authenticates records
// nor executes decision rules or interprets statuses as independently proven.
func ExplainDecisionTrace(input domain.DecisionTraceExplainInput) (domain.DecisionTraceExplanation, error) {
	fail := func(err error) (domain.DecisionTraceExplanation, error) {
		return domain.DecisionTraceExplanation{}, err
	}
	if input.SchemaVersion != "1" || len(input.NodeIDs) > traceMaxNodes {
		return fail(fmt.Errorf("invalid explanation schema or selection budget"))
	}
	if _, err := traceBoundedJSON(input); err != nil {
		return fail(err)
	}
	if input.Graph.GraphHash == "" {
		return fail(fmt.Errorf("explanation requires a frozen graph hash"))
	}
	graph, err := FreezeDecisionTrace(input.Graph)
	if err != nil {
		return fail(err)
	}
	built, err := BuildDecisionTrace(input.BuildInput)
	if err != nil {
		return fail(err)
	}
	if graph.GraphHash != built.GraphHash {
		return fail(fmt.Errorf("explanation graph does not match build input"))
	}
	index := traceExplanationIndex{nodes: make(map[string]domain.DecisionTraceNode, len(graph.Nodes)), entries: make(map[string]*domain.DecisionTraceExplanationEntry, len(graph.Nodes))}
	for _, node := range graph.Nodes {
		index.nodes[node.ID] = node
	}
	selected := make(map[string]bool, len(input.NodeIDs))
	for _, id := range input.NodeIDs {
		if _, exists := index.nodes[id]; !exists || selected[id] {
			return fail(fmt.Errorf("unknown or duplicate explanation node ID"))
		}
		selected[id] = true
	}
	report := input.BuildInput.Report
	if err = index.occurrence(domain.TraceInput, "report", report, "report", "", "Declared decision report"); err != nil {
		return fail(err)
	}
	if err = index.occurrence(domain.TraceInput, "policy", report.Policy, "report", "/policy", "Declared reuse policy"); err != nil {
		return fail(err)
	}
	if input.BuildInput.Idea != nil {
		if err = index.occurrence(domain.TraceInput, "idea", input.BuildInput.Idea, "idea", "", "Declared initial idea"); err != nil {
			return fail(err)
		}
	}
	if input.BuildInput.RuleVersions != nil {
		if err = index.occurrence(domain.TraceInput, "rule_versions", input.BuildInput.RuleVersions, "rule_versions", "", "Declared rule versions"); err != nil {
			return fail(err)
		}
	}
	if err = index.occurrence(domain.TraceDecision, "recorded-decision", report.Decision, "report", "/decision", "Historical decision: "+string(report.Decision)+" (not recalculated)"); err != nil {
		return fail(err)
	}
	for i, row := range report.Matrix {
		for j, candidate := range row.Candidates {
			pointer := tracePointer("matrix", strconv.Itoa(i), "candidates", strconv.Itoa(j))
			key := struct{ Category, Capability, Candidate string }{"capability", row.Capability, candidate.Candidate}
			if err = index.occurrence(domain.TraceClaim, key, candidate, "report", pointer, "Declared capability "+row.Capability+" / "+candidate.Candidate+": "+string(candidate.Status)); err != nil {
				return fail(err)
			}
			if err = index.sources(candidate.Evidence, pointer+"/evidence"); err != nil {
				return fail(err)
			}
		}
	}
	for _, group := range []struct {
		name, field string
		reasons     []domain.DecisionReason
	}{{"reason", "reasons", report.Reasons}, {"blocker", "blockers", report.Blockers}} {
		for i, reason := range group.reasons {
			pointer := tracePointer(group.field, strconv.Itoa(i))
			key := struct{ Category, RuleID, Effect, Candidate, Message string }{group.name, reason.RuleID, reason.Effect, reason.Candidate, reason.Message}
			if err = index.occurrence(domain.TraceClaim, key, reason, "report", pointer, "Recorded "+group.name+" "+reason.RuleID+" ["+reason.Effect+"]: "+reason.Message); err != nil {
				return fail(err)
			}
			if err = index.sources(reason.Evidence, pointer+"/evidence"); err != nil {
				return fail(err)
			}
		}
	}
	for i, component := range report.Components {
		pointer := tracePointer("components", strconv.Itoa(i))
		key := struct{ Category, Name, Kind, Path string }{"component", component.Name, component.Kind, component.Path}
		if err = index.occurrence(domain.TraceClaim, key, component, "report", pointer, "Declared component "+component.Name+" ["+component.Kind+"]: "+component.Path); err != nil {
			return fail(err)
		}
		if err = index.sources(component.Evidence, pointer+"/evidence"); err != nil {
			return fail(err)
		}
	}
	if report.Composition != nil {
		if err = index.occurrence(domain.TraceClaim, "composition-report", report.Composition, "report", "/composition", "Declared composition report"); err != nil {
			return fail(err)
		}
		for i, composition := range report.Composition.Candidates {
			pointer := tracePointer("composition", "candidates", strconv.Itoa(i))
			key := struct{ Category, ID string }{"composition", composition.ID}
			if err = index.occurrence(domain.TraceClaim, key, composition, "report", pointer, "Declared composition "+composition.ID+": "+string(composition.Status)); err != nil {
				return fail(err)
			}
			if err = index.sources(composition.Evidence, pointer+"/evidence"); err != nil {
				return fail(err)
			}
			for j, member := range composition.Members {
				if err = index.sources(member.Component.Evidence, pointer+"/members/"+strconv.Itoa(j)+"/component/evidence"); err != nil {
					return fail(err)
				}
			}
		}
	}
	if report.TechnicalEvidence != nil {
		if err = index.occurrence(domain.TraceClaim, "technical-report", report.TechnicalEvidence, "report", "/technical_evidence", "Declared consultative technical review"); err != nil {
			return fail(err)
		}
		for i, claim := range report.TechnicalEvidence.Claims {
			key := struct{ Category, Candidate, Claim string }{"technical-claim", claim.Candidate, claim.Claim}
			if err = index.occurrence(domain.TraceClaim, key, claim, "report", tracePointer("technical_evidence", "claims", strconv.Itoa(i)), "Technical declaration "+claim.Candidate+" / "+claim.Claim+": "+string(claim.Status)); err != nil {
				return fail(err)
			}
		}
		for i, verification := range report.TechnicalEvidence.Verifications {
			pointer := tracePointer("technical_evidence", "verifications", strconv.Itoa(i))
			key := struct{ Category, Finding, Candidate, Predicate, Value string }{"technical-verification", verification.FindingID, verification.Candidate, string(verification.Predicate), verification.Value}
			if err = index.occurrence(domain.TraceClaim, key, verification, "report", pointer, "Verifier declaration "+verification.FindingID+" ["+string(verification.Status)+"]: "+verification.Message); err != nil {
				return fail(err)
			}
			if verification.Source != nil {
				source := *verification.Source
				if err = index.occurrence(domain.TraceSource, traceSourceKey(source), source, "report", pointer+"/source", "Declared source: "+source.Path); err != nil {
					return fail(err)
				}
			}
		}
	}
	if len(index.entries) != len(index.nodes) {
		return fail(fmt.Errorf("explanation is missing bound graph nodes"))
	}
	reportFP, err := traceFingerprint(report)
	if err != nil {
		return fail(err)
	}
	result := domain.DecisionTraceExplanation{SchemaVersion: domain.DecisionTraceExplanationSchemaVersion, GraphHash: graph.GraphHash, ReportFingerprint: reportFP, Scope: domain.DecisionTraceScope, Entries: []domain.DecisionTraceExplanationEntry{}, Limitations: []string{
		"Locations bind the exact declared build input to the recorded graph; they do not authenticate sources or prove claims.",
		"Only recorded dependencies are indexed; absence of a locator does not prove absence of a real dependency.",
		"Statements are bounded display projections; locations reference the untruncated records and historical decisions are not recalculated.",
	}}
	for id, entry := range index.entries {
		if len(selected) != 0 && !selected[id] {
			continue
		}
		sort.Slice(entry.Locations, func(i, j int) bool {
			if entry.Locations[i].Document != entry.Locations[j].Document {
				return entry.Locations[i].Document < entry.Locations[j].Document
			}
			return entry.Locations[i].Pointer < entry.Locations[j].Pointer
		})
		result.Entries = append(result.Entries, *entry)
	}
	sort.Slice(result.Entries, func(i, j int) bool { return result.Entries[i].NodeID < result.Entries[j].NodeID })
	if _, err := traceBoundedJSON(result); err != nil {
		return fail(fmt.Errorf("decision trace explanation output exceeds budget: %w", err))
	}
	return result, nil
}
