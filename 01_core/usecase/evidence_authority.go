/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/evidence_authority.md
 * @layer L1
 * @updated 2026-09-07
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/hex"
	"strings"
)

func evidenceDigest(value string) (string, bool) {
	value = strings.TrimPrefix(value, "sha256:")
	decoded, err := hex.DecodeString(value)
	return strings.ToLower(value), err == nil && len(decoded) == 32
}

// EvidenceSourcesAdmissible validates declared source pointers, not behavior.
// Content-bound consumers must additionally verify their collected bytes.
func EvidenceSourcesAdmissible(ref domain.CommitRef, sources []domain.EvidenceSource, implementation bool) bool {
	if !competitiveRefValid(ref) || len(sources) == 0 {
		return false
	}
	foundImplementation := false
	for _, s := range sources {
		_, hashOK := evidenceDigest(s.ContentHash)
		if !hashOK || s.Ref != ref || !validMechanismPath(s.Path) || strings.TrimSpace(s.RuleID) == "" || s.Line < 0 {
			return false
		}
		if s.Kind != domain.EvidenceRepository && (s.Line <= 0 || strings.TrimSpace(s.Snippet) == "") {
			return false
		}
		switch s.Kind {
		case domain.EvidenceSymbol, domain.EvidenceFile, domain.EvidenceTestImplementation:
			if !isCodePath(s.Path) {
				return false
			}
			foundImplementation = true
		case domain.EvidenceTest, domain.EvidenceExample:
			if !isCodePath(s.Path) {
				return false
			}
		case domain.EvidenceDocumentation, domain.EvidenceManifest, domain.EvidenceLicense, domain.EvidenceRepository:
		default:
			return false
		}
	}
	return !implementation || foundImplementation
}

func authoritySourcesOfKind(ref domain.CommitRef, sources []domain.EvidenceSource, kind domain.EvidenceKind) bool {
	if !EvidenceSourcesAdmissible(ref, sources, false) {
		return false
	}
	for _, s := range sources {
		if s.Kind == kind {
			return true
		}
	}
	return false
}

// Known built-in rule IDs retain the same requirement mapping as the producer.
// Custom IDs remain supplied assessments; pointer admission is not content validation.
func knownEvidenceRulesMatchCapability(e domain.CapabilityEvidence) bool {
	catalog := DefaultEvidenceRules()
	var matched *domain.CapabilityRule
	if strings.TrimSpace(e.Capability) != "" {
		matched = findRule(e.Capability, catalog)
	}
	for _, source := range e.Sources {
		for _, rule := range catalog.Rules {
			if source.RuleID == rule.ID && (matched == nil || matched.ID != rule.ID) {
				return false
			}
		}
	}
	return true
}

func sanitizeCapability(ref domain.CommitRef, e domain.CapabilityEvidence) domain.CapabilityEvidence {
	e.Sources = append([]domain.EvidenceSource(nil), e.Sources...)
	if e.Status == "" {
		e.Status = domain.EvidenceInconclusive
	}
	if e.Status == domain.EvidenceAbsent && !EvidenceSourcesAdmissible(ref, e.Sources, false) {
		e.Status = domain.EvidenceInconclusive
		e.Confidence = 0
		e.Explanation = "ausência sem testemunho admissível do escopo examinado; " + e.Explanation
	}
	if e.Status == domain.EvidenceProven && !knownEvidenceRulesMatchCapability(e) {
		e.Status = domain.EvidenceInconclusive
		e.Confidence = 0
		e.Explanation = "regra conhecida da fonte não corresponde à capacidade avaliada; " + e.Explanation
	}
	if e.Status == domain.EvidenceProven && EvidenceSourcesAdmissible(ref, e.Sources, false) {
		for _, source := range e.Sources {
			if negativeLocalOperation(normalizeEvidenceText(source.Snippet)) && matchesAny(normalizeEvidenceText(e.Capability), []string{"local", "offline"}) {
				e.Status = domain.EvidenceConflict
				e.Explanation = "fonte declara negação de operação local"
			}
		}
	}
	if e.Status == domain.EvidenceProven && !EvidenceSourcesAdmissible(ref, e.Sources, true) {
		e.Status = domain.EvidenceInconclusive
		if EvidenceSourcesAdmissible(ref, e.Sources, false) {
			e.Status = domain.EvidencePartial
		}
		e.Explanation = "capacidade não comprovada: fonte ausente, incompatível ou apenas declarativa; " + e.Explanation
		e.Confidence = 0
	}
	if (e.Status == domain.EvidencePartial || e.Status == domain.EvidenceConflict) && !EvidenceSourcesAdmissible(ref, e.Sources, false) {
		e.Status = domain.EvidenceInconclusive
		e.Explanation = "testemunho sem identidade admissível; " + e.Explanation
		e.Confidence = 0
	}
	return e
}

// SanitizeAssessmentEvidence is the common admission boundary for supplied
// assessments. It preserves caller-owned slices and never manufactures sources.
func SanitizeAssessmentEvidence(a domain.CandidateAssessment) domain.CandidateAssessment {
	a.Capabilities = append([]domain.CapabilityEvidence(nil), a.Capabilities...)
	a.Constraints = append([]domain.ConstraintEvidence(nil), a.Constraints...)
	a.ProblemAlignment = sanitizeCapability(a.Repository, a.ProblemAlignment)
	proven := map[string]domain.CapabilityEvidence{}
	for i, e := range a.Capabilities {
		a.Capabilities[i] = sanitizeCapability(a.Repository, e)
		if a.Capabilities[i].Status == domain.EvidenceProven {
			proven[normalizeEvidenceText(e.Capability)] = a.Capabilities[i]
		}
	}
	for i, c := range a.Constraints {
		a.Constraints[i].Evidence = sanitizeCapability(a.Repository, c.Evidence)
	}
	if a.License.Status == "" || (a.License.Status == domain.EvidenceProven && !authoritySourcesOfKind(a.Repository, a.License.Sources, domain.EvidenceLicense)) {
		a.License.Status = domain.EvidenceInconclusive
	}
	if a.License.Status == domain.EvidenceProven {
		valid := false
		for _, s := range a.License.Sources {
			valid = valid || s.Kind == domain.EvidenceLicense && authoritativeLicensePath(s.Path)
		}
		if !valid {
			a.License.Status = domain.EvidenceInconclusive
		}
	}
	if a.Maintenance.Status == "" || (a.Maintenance.Status == domain.EvidenceProven && !authoritySourcesOfKind(a.Repository, a.Maintenance.Sources, domain.EvidenceRepository)) {
		a.Maintenance.Status = domain.EvidenceInconclusive
	}
	if EvidenceSourcesAdmissible(a.Repository, a.Operation.Sources, false) {
		for _, source := range a.Operation.Sources {
			if negativeLocalOperation(normalizeEvidenceText(source.Snippet)) {
				a.Operation.LocalStatus = domain.EvidenceConflict
			}
		}
	}
	if a.Operation.LocalStatus == domain.EvidenceConflict && !EvidenceSourcesAdmissible(a.Repository, a.Operation.Sources, false) {
		a.Operation.LocalStatus = domain.EvidenceInconclusive
	}
	if a.Operation.LocalStatus == "" || (a.Operation.LocalStatus == domain.EvidenceProven && !EvidenceSourcesAdmissible(a.Repository, a.Operation.Sources, true)) {
		a.Operation.LocalStatus = domain.EvidenceInconclusive
	}
	parts := make([]domain.ReusableComponent, 0, len(a.ReusableParts))
	for _, p := range a.ReusableParts {
		if !EvidenceSourcesAdmissible(a.Repository, p.Evidence, true) {
			continue
		}
		p.Capabilities = append([]string(nil), p.Capabilities...)
		var caps []string
		for _, c := range p.Capabilities {
			e, ok := proven[normalizeEvidenceText(c)]
			if !ok {
				continue
			}
			linked := false
			for _, s := range p.Evidence {
				for _, t := range e.Sources {
					sd, _ := evidenceDigest(s.ContentHash)
					td, _ := evidenceDigest(t.ContentHash)
					linked = linked || s.Ref == t.Ref && s.Path == t.Path && sd == td && s.RuleID == t.RuleID
				}
			}
			if linked {
				caps = append(caps, c)
			}
		}
		p.Capabilities = caps
		if len(caps) > 0 {
			parts = append(parts, p)
		}
	}
	a.ReusableParts = parts
	return a
}

func negativeLocalOperation(line string) bool {
	return matchesAny(line, []string{"cannot run locally", "can not run locally", "does not run locally", "never runs locally", "cannot work offline", "does not work offline", "no local execution", "local execution is not supported", "nao executa localmente", "nao funciona offline"})
}

func negatedExternalTransfer(line string) bool {
	return matchesAny(line, []string{"never uploads", "never upload", "does not upload", "do not upload", "not upload", "never sends code", "never send code", "does not send code", "do not send code", "nao envia", "nao enviamos"})
}
