/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/composition.md
 * @layer L1
 * @updated 2026-11-08
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

	"bateia/01_core/domain"
)

const compositionPartLimit = 16
const compositionSearchLimit = 30000
const compositionResultLimit = 100

func compositionMemberID(ref domain.CommitRef, part domain.ReusableComponent) string {
	symbols := append([]string(nil), part.Symbols...)
	sort.Strings(symbols)
	identity, _ := json.Marshal(struct {
		Kind, Path, Name string
		Symbols          []string
	}{part.Kind, part.Path, part.Name, symbols})
	digest := sha256.Sum256(identity)
	return ref.Key() + ":" + part.Path + ":" + part.Name + "#" + hex.EncodeToString(digest[:6])
}

func partCapabilityLinked(assessment domain.CandidateAssessment, part domain.ReusableComponent, capability string) bool {
	for _, evidence := range assessment.Capabilities {
		if normalizeDecisionText(evidence.Capability) != capability || evidence.Status != domain.EvidenceProven {
			continue
		}
		for _, partSource := range part.Evidence {
			for _, witness := range evidence.Sources {
				if partSource.Ref == witness.Ref && partSource.Kind == witness.Kind && partSource.Path == witness.Path && partSource.RuleID == witness.RuleID && partSource.ContentHash == witness.ContentHash {
					return true
				}
			}
		}
	}
	return false
}

func analyzeCompositions(cs []decisionCandidate, caps []string, policy domain.ReuseDecisionPolicy, input *domain.CompositionInput) (domain.CompositionReport, error) {
	report := domain.CompositionReport{SchemaVersion: "1", Method: "BOUNDED_SUBSET_CONSTRAINT_SEARCH_V1", Candidates: []domain.CompositionCandidate{}}
	members := make([]domain.CompositionMember, 0)
	seen := map[string]bool{}
	for _, candidate := range cs {
		if len(candidate.blockers) > 0 || candidate.ambiguous {
			continue
		}
		for _, part := range validParts(candidate.a.ReusableParts) {
			if !EvidenceSourcesAdmissible(candidate.a.Repository, part.Evidence, false) {
				continue
			}
			claimed := normalizedSorted(part.Capabilities)
			proven := make([]string, 0, len(claimed))
			for _, capability := range claimed {
				if candidate.covered[capability] && contains(caps, capability) && partCapabilityLinked(candidate.a, part, capability) {
					proven = append(proven, capability)
				}
			}
			if len(proven) == 0 {
				continue
			}
			part.Capabilities = proven
			id := compositionMemberID(candidate.a.Repository, part)
			if seen[id] {
				continue
			}
			seen[id] = true
			members = append(members, domain.CompositionMember{ID: id, Repository: candidate.a.Repository, Component: part})
		}
	}
	sort.Slice(members, func(i, j int) bool { return members[i].ID < members[j].ID })
	if err := validateCompositionInput(input, members); err != nil {
		return report, err
	}
	if len(members) > compositionPartLimit {
		members = members[:compositionPartLimit]
		report.Truncated = true
		report.StopReason = "PART_LIMIT"
	}
	if len(members) < 2 {
		return report, nil
	}
	max := policy.MaxComponents
	if max > len(members) {
		max = len(members)
	}
	if max > 6 {
		max = 6
		report.Truncated = true
		report.StopReason = "MAX_COMPONENTS_LIMIT"
	}
	visited := 0
	selected := make([]domain.CompositionMember, 0, max)
	var visit func(int)
	visit = func(start int) {
		if visited >= compositionSearchLimit || len(report.Candidates) >= compositionResultLimit {
			report.Truncated = true
			report.StopReason = "SEARCH_LIMIT"
			return
		}
		if len(selected) >= 2 {
			visited++
			composition := evaluateComposition(selected, caps, input)
			if len(composition.Gaps) == 0 {
				report.Candidates = append(report.Candidates, composition)
			}
		}
		if len(selected) == max {
			return
		}
		for i := start; i < len(members); i++ {
			if report.StopReason == "SEARCH_LIMIT" {
				return
			}
			selected = append(selected, members[i])
			visit(i + 1)
			selected = selected[:len(selected)-1]
		}
	}
	visit(0)
	sort.Slice(report.Candidates, func(i, j int) bool {
		a, b := report.Candidates[i], report.Candidates[j]
		if a.Status != b.Status {
			return compositionStatusOrder(a.Status) < compositionStatusOrder(b.Status)
		}
		if len(a.Members) != len(b.Members) {
			return len(a.Members) < len(b.Members)
		}
		return a.ID < b.ID
	})
	for i := range report.Candidates {
		report.Candidates[i].Evidence = projectEvidence(report.Candidates[i].Evidence)
		for j := range report.Candidates[i].Members {
			member := &report.Candidates[i].Members[j]
			member.Component.Evidence = projectEvidence(member.Component.Evidence)
		}
	}
	return report, nil
}

func compositionStatusOrder(status domain.CompositionStatus) int {
	switch status {
	case domain.CompositionVerified:
		return 0
	case domain.CompositionProposed:
		return 1
	default:
		return 2
	}
}

func validateCompositionInput(input *domain.CompositionInput, members []domain.CompositionMember) error {
	if input == nil {
		return nil
	}
	known := map[string]bool{}
	for _, member := range members {
		known[member.ID] = true
	}
	connections := map[[2]string]domain.CompositionConnection{}
	for _, edge := range input.Connections {
		if !known[edge.From] || !known[edge.To] || edge.From == edge.To {
			return fmt.Errorf("composition.connections referencia membros desconhecidos ou idênticos")
		}
		if !validEvidenceStatus(edge.Status) {
			return fmt.Errorf("composition.connections possui status inválido")
		}
		if edge.Cost != nil && (!finite(edge.Cost.Low) || !finite(edge.Cost.High) || edge.Cost.Low < 0 || edge.Cost.High < edge.Cost.Low || strings.TrimSpace(edge.Cost.Unit) == "") {
			return fmt.Errorf("composition.connections possui custo inválido")
		}
		key := [2]string{edge.From, edge.To}
		if previous, exists := connections[key]; exists && !reflect.DeepEqual(previous, edge) {
			return fmt.Errorf("composition.connections possui declarações divergentes para o mesmo par ordenado")
		}
		connections[key] = edge
	}
	for _, exclusion := range input.Exclusions {
		if !known[exclusion.Left] || !known[exclusion.Right] || exclusion.Left == exclusion.Right {
			return fmt.Errorf("composition.exclusions referencia membros desconhecidos ou idênticos")
		}
	}
	for _, requirement := range input.Requirements {
		if !known[requirement.Component] || !known[requirement.Requires] || requirement.Component == requirement.Requires {
			return fmt.Errorf("composition.requirements referencia membros desconhecidos ou idênticos")
		}
	}
	for _, verification := range input.Verifications {
		if !validEvidenceStatus(verification.Status) || len(verification.Members) < 2 {
			return fmt.Errorf("composition.verifications inválida")
		}
		unique := map[string]bool{}
		for _, id := range verification.Members {
			if !known[id] || unique[id] {
				return fmt.Errorf("composition.verifications referencia membro desconhecido ou duplicado")
			}
			unique[id] = true
		}
	}
	return nil
}

func evaluateComposition(members []domain.CompositionMember, caps []string, input *domain.CompositionInput) domain.CompositionCandidate {
	out := domain.CompositionCandidate{Members: append([]domain.CompositionMember(nil), members...), Status: domain.CompositionProposed, Gaps: []string{}}
	ids := make([]string, 0, len(members))
	selected := map[string]domain.CompositionMember{}
	counts := map[string]int{}
	for _, member := range members {
		ids = append(ids, member.ID)
		selected[member.ID] = member
		for _, capability := range member.Component.Capabilities {
			counts[capability]++
		}
		out.Evidence = append(out.Evidence, member.Component.Evidence...)
	}
	out.ID = "compose:" + strings.Join(ids, "+")
	for _, capability := range caps {
		if counts[capability] == 0 {
			out.Gaps = append(out.Gaps, capability)
		} else {
			out.Coverage = append(out.Coverage, capability)
			if counts[capability] > 1 {
				out.Overlap = append(out.Overlap, capability)
			}
		}
	}
	if input == nil {
		out.Unknowns = []string{"interfaces não declaradas", "ensaio de integração do conjunto não declarado", "custo de integração desconhecido"}
		return out
	}
	for _, exclusion := range input.Exclusions {
		if selected[exclusion.Left].ID != "" && selected[exclusion.Right].ID != "" {
			out.Status = domain.CompositionIncompatible
			out.Unknowns = append(out.Unknowns, "exclusão: "+exclusion.Reason)
		}
	}
	for _, requirement := range input.Requirements {
		if selected[requirement.Component].ID != "" && selected[requirement.Requires].ID == "" {
			out.Status = domain.CompositionIncompatible
			out.Unknowns = append(out.Unknowns, "dependência ausente: "+requirement.Component+" requer "+requirement.Requires)
		}
	}
	parent := map[string]string{}
	for _, id := range ids {
		parent[id] = id
	}
	var root func(string) string
	root = func(id string) string {
		for parent[id] != id {
			id = parent[id]
		}
		return id
	}
	cost := domain.IntegrationCost{}
	costKnown := true
	usedEdges := 0
	edges := append([]domain.CompositionConnection(nil), input.Connections...)
	sort.Slice(edges, func(i, j int) bool {
		a, b := edges[i], edges[j]
		return a.From+"\x00"+a.To < b.From+"\x00"+b.To
	})
	for _, edge := range edges {
		if selected[edge.From].ID == "" || selected[edge.To].ID == "" {
			continue
		}
		if edge.Status == domain.EvidenceConflict {
			out.Status = domain.CompositionIncompatible
			out.Unknowns = append(out.Unknowns, "interface conflitante: "+edge.From+" ↔ "+edge.To)
			continue
		}
		if edge.Status != domain.EvidenceProven || !compositionEvidenceAdmissible(selected, edge.Sources, true) {
			continue
		}
		left, right := root(edge.From), root(edge.To)
		if left == right {
			continue
		}
		parent[right] = left
		usedEdges++
		out.Evidence = append(out.Evidence, edge.Sources...)
		if edge.Cost == nil {
			costKnown = false
		} else if cost.Unit != "" && cost.Unit != edge.Cost.Unit {
			costKnown = false
		} else {
			cost.Unit = edge.Cost.Unit
			cost.Low += edge.Cost.Low
			cost.High += edge.Cost.High
		}
	}
	if usedEdges != len(members)-1 {
		out.Unknowns = append(out.Unknowns, "grafo de interfaces comprovadas não conecta todos os membros")
	}
	if costKnown && usedEdges == len(members)-1 && cost.Unit != "" && !math.IsInf(cost.High, 0) {
		out.IntegrationCost = &cost
	} else {
		out.Unknowns = append(out.Unknowns, "custo de integração desconhecido")
	}
	verified := false
	for _, verification := range input.Verifications {
		if !sameCompositionMembers(ids, verification.Members) {
			continue
		}
		if verification.Status == domain.EvidenceConflict || (verification.Executed && !verification.Passed) {
			out.Status = domain.CompositionIncompatible
			out.Unknowns = append(out.Unknowns, "ensaio de integração falhou")
			continue
		}
		if verification.Status == domain.EvidenceProven && verification.Executed && verification.Passed && compositionEvidenceAdmissible(selected, verification.Sources, true) {
			verified = true
			out.Evidence = append(out.Evidence, verification.Sources...)
		}
	}
	if !verified {
		out.Unknowns = append(out.Unknowns, "ensaio executado e aprovado do conjunto exato não comprovado")
	}
	if out.Status != domain.CompositionIncompatible && usedEdges == len(members)-1 && verified {
		out.Status = domain.CompositionVerified
	}
	sort.Strings(out.Unknowns)
	return out
}

func sameCompositionMembers(sorted []string, other []string) bool {
	if len(sorted) != len(other) {
		return false
	}
	copy := append([]string(nil), other...)
	sort.Strings(copy)
	for i := range sorted {
		if sorted[i] != copy[i] {
			return false
		}
	}
	return true
}

func compositionEvidenceAdmissible(members map[string]domain.CompositionMember, sources []domain.EvidenceSource, requireTest bool) bool {
	if len(sources) == 0 {
		return false
	}
	test := false
	for _, source := range sources {
		found := false
		for _, member := range members {
			if member.Repository == source.Ref {
				found = true
				break
			}
		}
		if !found || !EvidenceSourcesAdmissible(source.Ref, []domain.EvidenceSource{source}, false) {
			return false
		}
		digest := sha256.Sum256([]byte(source.Snippet))
		if source.Snippet != "" && strings.TrimPrefix(source.ContentHash, "sha256:") != hex.EncodeToString(digest[:]) {
			return false
		}
		test = test || source.Kind == domain.EvidenceTestImplementation
	}
	return !requireTest || test
}
