/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/competitive_planning.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
)

// This validates a declared pointer, not the behavior claimed by its author.
func competitiveSourcesValid(ref domain.CommitRef, sources []domain.EvidenceSource, implementation bool) bool {
	return EvidenceSourcesAdmissible(ref, sources, implementation)
}

func competitiveRefValid(ref domain.CommitRef) bool {
	c := strings.ToLower(strings.TrimSpace(ref.Commit))
	return strings.TrimSpace(ref.Owner) != "" && strings.TrimSpace(ref.Repo) != "" && c != "" && c != "head" && c != "worktree" && c != "main" && c != "master" && c != "latest" && !strings.HasPrefix(c, "refs/") && !strings.ContainsAny(c, " ~^:{}")
}

func competitiveGates(a domain.CandidateAssessment, p domain.KatamariPolicy, reuse domain.ReuseDecisionPolicy) (blocks, checks []string) {
	a = SanitizeAssessmentEvidence(a)
	if !competitiveRefValid(a.Repository) {
		blocks = append(blocks, "revisão do candidato não congelada")
	}
	if len(a.Failures) > 0 {
		checks = append(checks, "resolver falhas da coleta antes de qualificar o candidato")
	}
	for _, c := range a.Constraints {
		switch c.Evidence.Status {
		case domain.EvidenceConflict, domain.EvidenceAbsent:
			blocks = append(blocks, "restrição não atendida: "+c.Constraint)
		case domain.EvidenceProven:
			if !competitiveSourcesValid(a.Repository, c.Evidence.Sources, false) {
				checks = append(checks, "rastrear comprovação da restrição: "+c.Constraint)
			}
		default:
			checks = append(checks, "verificar restrição: "+c.Constraint)
		}
	}
	capStates := map[string]domain.EvidenceStatus{}
	for _, c := range a.Capabilities {
		key := normalizeEvidenceText(c.Capability)
		if prev, ok := capStates[key]; ok && prev != c.Status {
			checks = append(checks, "reconciliar avaliações da capacidade: "+c.Capability)
		}
		capStates[key] = c.Status
		if c.Status == domain.EvidenceConflict {
			blocks = append(blocks, "capacidade em conflito: "+c.Capability)
		}
	}
	if a.ProblemAlignment.Status == domain.EvidenceConflict || a.ProblemAlignment.Status == domain.EvidenceAbsent {
		blocks = append(blocks, "candidato não atende à intenção")
	}
	if a.ProblemAlignment.Status != domain.EvidenceProven || !competitiveSourcesValid(a.Repository, a.ProblemAlignment.Sources, false) {
		checks = append(checks, "verificar relação da implementação com a intenção")
	}
	if a.License.Status == domain.EvidenceProven {
		if !licenseAllowed(a.License, p.AllowedLicenses) || (len(reuse.AllowedLicenses) > 0 && !licenseAllowed(a.License, reuse.AllowedLicenses)) {
			blocks = append(blocks, "licença não permitida pelas políticas")
		}
		valid := competitiveSourcesValid(a.Repository, a.License.Sources, false)
		for _, s := range a.License.Sources {
			valid = valid && s.Kind == domain.EvidenceLicense && authoritativeLicensePath(s.Path)
		}
		if !valid {
			checks = append(checks, "localizar licença com autoridade e fonte rastreável")
		}
	} else if a.License.Status == domain.EvidenceConflict {
		blocks = append(blocks, "licença em conflito")
	} else {
		checks = append(checks, "verificar licença do projeto")
	}
	if a.Maintenance.Archived && !reuse.AllowArchived {
		blocks = append(blocks, "repositório arquivado proibido pela política")
	}
	if a.Maintenance.Status == domain.EvidenceConflict || a.Maintenance.Status == domain.EvidenceAbsent {
		blocks = append(blocks, "manutenção em conflito ou encerrada")
	}
	if a.Maintenance.Status != domain.EvidenceProven || !competitiveSourcesValid(a.Repository, a.Maintenance.Sources, false) {
		checks = append(checks, "obter evidência datada de manutenção")
	}
	if reuse.RequireLocal || strings.EqualFold(p.RequiredOperation, "local") {
		switch a.Operation.LocalStatus {
		case domain.EvidenceConflict, domain.EvidenceAbsent:
			blocks = append(blocks, "operação local obrigatória não atendida")
		case domain.EvidenceProven:
			if !competitiveSourcesValid(a.Repository, a.Operation.Sources, false) {
				checks = append(checks, "rastrear comprovação de operação local")
			}
		default:
			checks = append(checks, "verificar operação local obrigatória")
		}
	}
	return uniqueSorted(blocks), uniqueSorted(checks)
}

func integrationChecks() []string {
	return []string{"delimitar implementação, interface e dependências que serão integradas", "verificar compatibilidade com o alvo e com os demais componentes", "avaliar termos aplicáveis ao modo de reúso proposto", "estimar adaptação, migração, operação e possibilidade de reversão", "definir e executar teste de aceitação isolado antes da adoção"}
}

func internalEvolutionOption(f domain.InternalCapabilityFinding) domain.EvolutionOption {
	o := domain.EvolutionOption{Capability: f.Capability, EvidenceStatus: f.Status, Sources: f.Sources, NextChecks: f.NextChecks, Rationale: f.Rationale}
	switch f.Disposition {
	case domain.InternalKeep:
		o.Action = "KEEP"
	case domain.InternalReviewOverlap:
		o.Action = "REVIEW_OVERLAP"
	case domain.InternalGap:
		o.Action = "COMPARE_BUILD_REUSE"
		o.NextChecks = []string{"comparar uma implementação interna delimitada com opções externas antes de escolher"}
	default:
		o.Action = "INVESTIGATE_INTERNAL"
		if f.Status == domain.EvidencePartial {
			o.Action = "ASSESS_INTERNAL_EXTENSION"
		}
	}
	return o
}

func relevantEvolutionOptions(auto domain.AutophagyReport, a domain.CandidateAssessment, checks []string) []domain.EvolutionOption {
	wanted := normalizedSet(append(append(append(append([]string{}, auto.Gaps...), auto.Partial...), auto.Unknown...), auto.Existing...))
	var out []domain.EvolutionOption
	for _, c := range a.Capabilities {
		if !wanted[normalizeEvidenceText(c.Capability)] || !competitiveSourcesValid(a.Repository, c.Sources, false) || c.Status == domain.EvidenceAbsent || c.Status == domain.EvidenceConflict {
			continue
		}
		action := "INVESTIGATE_EXTERNAL"
		if normalizedSet(auto.Existing)[normalizeEvidenceText(c.Capability)] {
			action = "COMPARE_REPLACEMENT"
		}
		next := append([]string{}, checks...)
		next = append(next, "verificar comportamento requerido para "+c.Capability)
		next = append(next, "comparar manter/completar internamente com integrar ou substituir; ganho ainda não medido")
		out = append(out, domain.EvolutionOption{Action: action, Capability: c.Capability, Candidate: a.Repository.Key(), EvidenceStatus: c.Status, Sources: c.Sources, Rationale: "oportunidade de investigação; não conta como cobertura nem autoriza adoção", NextChecks: uniqueSorted(next)})
	}
	return out
}

// These inputs represent sets of evidence, candidates and capabilities. This is
// deliberately scoped to competitive analysis, not a generic JSON canonicalizer.
func competitiveHash(v any) string {
	data, _ := json.Marshal(v)
	var value any
	_ = json.Unmarshal(data, &value)
	return recoveryHash(competitiveSetValue(value))
}
func competitiveSetValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, a := range x {
			x[k] = competitiveSetValue(a)
		}
		return x
	case []any:
		byJSON := map[string]any{}
		for _, a := range x {
			a = competitiveSetValue(a)
			b, _ := json.Marshal(a)
			byJSON[string(b)] = a
		}
		keys := make([]string, 0, len(byJSON))
		for k := range byJSON {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		out := make([]any, 0, len(keys))
		for _, k := range keys {
			out = append(out, byJSON[k])
		}
		return out
	default:
		return v
	}
}

func sortEvolutionOptions(options []domain.EvolutionOption) []domain.EvolutionOption {
	byKey := map[string]domain.EvolutionOption{}
	for _, o := range options {
		o.Sources = canonicalEvidenceSources(o.Sources)
		o.NextChecks = uniqueSorted(o.NextChecks)
		byKey[competitiveHash(o)] = o
	}
	out := make([]domain.EvolutionOption, 0, len(byKey))
	for _, o := range byKey {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		return a.Capability+"|"+a.Action+"|"+a.Candidate+"|"+competitiveHash(a) < b.Capability+"|"+b.Action+"|"+b.Candidate+"|"+competitiveHash(b)
	})
	return out
}

func uniqueCompetitiveRefs(refs []domain.CommitRef) []domain.CommitRef {
	byKey := map[string]domain.CommitRef{}
	for _, r := range refs {
		byKey[r.Key()] = r
	}
	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]domain.CommitRef, 0, len(keys))
	for _, k := range keys {
		out = append(out, byKey[k])
	}
	return out
}

// Older list-only inventories remain readable, but cannot establish missing
// behavior or justify new dependencies without their scoped findings.
func validatedAutophagy(input domain.AutophagyReport) (domain.AutophagyReport, []string) {
	out := domain.AutophagyReport{SchemaVersion: "2", Target: input.Target, Scope: uniqueCompetitiveRefs(input.Scope), Existing: []string{}, Partial: []string{}, Gaps: []string{}, Unknown: []string{}, InputHash: input.InputHash}
	digest, err := hex.DecodeString(input.InputHash)
	scoped := input.SchemaVersion == "2" && err == nil && len(digest) == 32 && len(input.Scope) > 0
	for _, r := range input.Scope {
		scoped = scoped && competitiveRefValid(r)
	}
	byCap := map[string][]domain.InternalCapabilityFinding{}
	names := map[string]string{}
	for _, name := range append(append(append(append([]string{}, input.Gaps...), input.Existing...), input.Partial...), input.Unknown...) {
		key := normalizeEvidenceText(name)
		if old, ok := names[key]; !ok || name < old {
			names[key] = name
		}
	}
	for _, f := range input.Findings {
		k := normalizeEvidenceText(f.Capability)
		if old, ok := names[k]; !ok || f.Capability < old {
			names[k] = f.Capability
		}
		byCap[k] = append(byCap[k], f)
	}
	keys := make([]string, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	warnings := []string{}
	for _, key := range keys {
		if key == "" {
			continue
		}
		f := domain.InternalCapabilityFinding{Capability: names[key], Status: domain.EvidenceInconclusive, Disposition: domain.InternalUnknown, Rationale: "inventário sem achado único e escopo verificável; reexecutar autofagia", NextChecks: []string{"reexecutar autofagia com bundles congelados e cobertura explícita"}}
		unique := map[string]domain.InternalCapabilityFinding{}
		for _, x := range byCap[key] {
			unique[competitiveHash(x)] = x
		}
		if scoped && len(unique) == 1 {
			for _, x := range unique {
				f = x
			}
			switch {
			case f.Status == domain.EvidenceAbsent && f.Disposition == domain.InternalGap && f.CoverageComplete:
				out.Gaps = append(out.Gaps, f.Capability)
			case f.Status == domain.EvidenceProven && scopedFindingSources(f.Sources, input.Scope, true) && (f.Disposition == domain.InternalKeep || f.Disposition == domain.InternalReviewOverlap):
				out.Existing = append(out.Existing, f.Capability)
			case f.Status == domain.EvidencePartial && scopedFindingSources(f.Sources, input.Scope, false):
				out.Partial = append(out.Partial, f.Capability)
			default:
				f.Status = domain.EvidenceInconclusive
				f.Disposition = domain.InternalUnknown
				out.Unknown = append(out.Unknown, f.Capability)
			}
		} else {
			out.Unknown = append(out.Unknown, f.Capability)
			warnings = append(warnings, "reavaliar inventário de "+f.Capability+": listas isoladas não comprovam ausência ou presença")
		}
		out.Findings = append(out.Findings, f)
		out.Options = append(out.Options, internalEvolutionOption(f))
	}
	sortAutophagyReport(&out)
	out.Options = sortEvolutionOptions(out.Options)
	return out, uniqueSorted(warnings)
}

func scopedFindingSources(sources []domain.EvidenceSource, scope []domain.CommitRef, implementation bool) bool {
	if len(sources) == 0 {
		return false
	}
	refs := map[string]domain.CommitRef{}
	for _, r := range scope {
		refs[r.Key()] = r
	}
	byRef := map[string][]domain.EvidenceSource{}
	for _, s := range sources {
		if _, ok := refs[s.Ref.Key()]; !ok {
			return false
		}
		byRef[s.Ref.Key()] = append(byRef[s.Ref.Key()], s)
	}
	for k, s := range byRef {
		if !competitiveSourcesValid(refs[k], s, implementation) {
			return false
		}
	}
	return true
}

func competitiveOperationSupported(operation string) bool {
	return operation == "" || strings.EqualFold(operation, "local")
}

func competitiveCapabilities(values []string) []string {
	names := map[string]string{}
	for _, name := range values {
		key := normalizeEvidenceText(name)
		if key == "" {
			continue
		}
		if old, ok := names[key]; !ok || name < old {
			names[key] = name
		}
	}
	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
