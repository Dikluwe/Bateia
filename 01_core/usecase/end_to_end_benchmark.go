/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/end_to_end_benchmark.md
 * @layer L1
 * @updated 2026-08-26
 */
package usecase

import (
	"bateia/01_core/domain"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func BenchmarkCaseHash(c domain.EndToEndBenchmarkCase) string {
	canonicalizeBenchmarkCase(&c)
	return recoveryHash(c)
}
func BenchmarkManifestHash(m domain.BenchmarkManifest) string { return recoveryHash(m) }

func ValidateBenchmarkCase(manifest domain.BenchmarkManifest, c domain.EndToEndBenchmarkCase) domain.CorrespondenceResult {
	r := domain.CorrespondenceResult{CaseID: c.ID, Verdict: domain.CorrespondencePreserved, Witnesses: []string{}}
	violate := func(s string) { r.Verdict = domain.CorrespondenceViolated; r.Witnesses = append(r.Witnesses, s) }
	unknown := func(s string) {
		if r.Verdict != domain.CorrespondenceViolated {
			r.Verdict = domain.CorrespondenceUnknown
		}
		r.Witnesses = append(r.Witnesses, s)
	}
	if c.SchemaVersion != "1" || c.ID == "" {
		violate("case sem identidade/schema")
	}
	if expected := manifest.CaseHashes[c.ID]; expected == "" {
		violate("case ausente do manifesto")
	} else if expected != BenchmarkCaseHash(c) {
		violate("hash do case diverge do manifesto")
	}
	if !reflect.DeepEqual(c.Idea, c.Discovery.Idea) {
		violate("DiscoverySnapshot.Idea diverge da ideia congelada")
	}
	if c.EvidenceRules.Version == "" || c.EvidenceRules.Hash == "" {
		unknown("catálogo de regras sem versão/hash")
	} else if ruleCatalogHash(c.EvidenceRules) != c.EvidenceRules.Hash {
		violate("hash do catálogo de regras inválido")
	}
	if c.ScreeningModel.SchemaVersion == "" || c.ReusePolicy.SchemaVersion == "" {
		unknown("modelo ou política sem versão")
	}
	if c.Provenance.OriginalRunHash == "" || c.Provenance.CollectorVersion == "" {
		unknown("proveniência incompleta")
	}
	discovered := map[string]domain.RepositoryCandidate{}
	for _, x := range c.Discovery.Candidates {
		discovered[x.Repository.Owner+"/"+x.Repository.Repo] = x.Repository
	}
	for _, b := range c.Bundles {
		key := b.Repository.Owner + "/" + b.Repository.Repo
		candidate, ok := discovered[key]
		if !ok {
			violate("bundle de candidato ausente na descoberta: " + key)
			continue
		}
		if b.Repository.Commit == "" || strings.EqualFold(b.Repository.Commit, "HEAD") {
			unknown("bundle em revisão opaca: " + key)
		}
		if candidate.CommitID != "" && !strings.EqualFold(candidate.CommitID, "HEAD") && candidate.CommitID != b.Repository.Commit {
			violate("commit do bundle diverge da descoberta: " + key)
		}
		if expected := c.Provenance.CandidateCommits[key]; expected == "" {
			unknown("commit ausente da proveniência: " + key)
		} else if expected != b.Repository.Commit {
			violate("commit diverge da proveniência: " + key)
		}
		for _, item := range b.Items {
			if item.Ref != b.Repository {
				violate("item pertence a outro commit: " + item.Path)
			} else if item.ContentHash == "" || item.ContentHash != contentHash(item.Content) {
				violate("hash inválido: " + item.Path)
			}
		}
	}
	capabilities := map[string]bool{}
	for _, x := range c.Idea.Capabilities {
		capabilities[normalizeEvidenceText(x)] = true
	}
	capabilities[normalizeEvidenceText(c.Idea.Problem)] = true
	for _, j := range c.ExpectedEvidence {
		if !capabilities[normalizeEvidenceText(j.Capability)] && !strings.HasPrefix(j.Capability, "gate:") {
			violate("julgamento fora da ideia/gates: " + j.Capability)
		}
	}
	if containsDecision(c.ExpectedDecision.Forbidden, c.ExpectedDecision.Expected) {
		violate("decisão esperada também proibida")
	}
	opaque := false
	for _, b := range c.Bundles {
		opaque = opaque || !b.Complete || len(b.Failures) > 0 || strings.EqualFold(b.Repository.Commit, "HEAD")
	}
	if opaque && c.ExpectedDecision.Expected != domain.DecisionUnknown {
		violate("case opaco exige INCONCLUSIVE")
	}
	sort.Strings(r.Witnesses)
	return r
}

func ValidateBenchmark(manifest domain.BenchmarkManifest, corpus domain.BenchmarkCorpus) domain.BenchmarkValidationReport {
	out := domain.BenchmarkValidationReport{SchemaVersion: "1", ManifestHash: BenchmarkManifestHash(manifest), Results: []domain.CorrespondenceResult{}}
	seen := map[domain.BenchmarkCaseID]bool{}
	for _, c := range corpus.Cases {
		r := ValidateBenchmarkCase(manifest, c)
		if seen[c.ID] {
			r.Verdict = domain.CorrespondenceViolated
			r.Witnesses = append(r.Witnesses, "case ID duplicado")
		}
		seen[c.ID] = true
		out.Results = append(out.Results, r)
		switch r.Verdict {
		case domain.CorrespondencePreserved:
			out.Preserved++
		case domain.CorrespondenceViolated:
			out.Violated++
		case domain.CorrespondenceUnknown:
			out.Unknown++
		}
	}
	sort.Slice(out.Results, func(i, j int) bool { return out.Results[i].CaseID < out.Results[j].CaseID })
	out.Valid = len(corpus.Cases) > 0 && out.Violated == 0 && out.Unknown == 0 && out.Preserved == len(corpus.Cases)
	return out
}

func EvaluateBenchmarkMutations(manifest domain.BenchmarkManifest, mutations []domain.BenchmarkMutation) domain.BenchmarkMutationReport {
	r := domain.BenchmarkMutationReport{SchemaVersion: "1", Total: len(mutations), Survivors: []string{}}
	for _, m := range mutations {
		v := ValidateBenchmarkCase(manifest, m.Case)
		if v.Verdict != domain.CorrespondencePreserved {
			r.Rejected++
		} else {
			r.Survivors = append(r.Survivors, m.ID)
		}
	}
	if r.Total > 0 {
		x := float64(r.Rejected) / float64(r.Total)
		r.MutationScore = &x
	}
	sort.Strings(r.Survivors)
	r.Sealed = r.Total > 0 && r.Rejected == r.Total
	return r
}

func RunEndToEndBenchmark(manifest domain.BenchmarkManifest, corpus domain.BenchmarkCorpus, split domain.BenchmarkSplit) (domain.EndToEndBenchmarkReport, error) {
	validation := ValidateBenchmark(manifest, corpus)
	byID := map[domain.BenchmarkCaseID]domain.CorrespondenceResult{}
	for _, r := range validation.Results {
		byID[r.CaseID] = r
	}
	report := domain.EndToEndBenchmarkReport{SchemaVersion: "1", Split: split, Cases: []domain.BenchmarkCaseResult{}}
	for _, c := range corpus.Cases {
		if split != "" && c.Split != split {
			continue
		}
		vr := byID[c.ID]
		result := domain.BenchmarkCaseResult{CaseID: c.ID, Correspondence: vr.Verdict, Witnesses: append([]string(nil), vr.Witnesses...), Utility: domain.UtilityUnjudgeable}
		if vr.Verdict != domain.CorrespondencePreserved {
			report.Cases = append(report.Cases, result)
			continue
		}
		screening, err := ScreenDiscovery(c.Discovery, c.ScreeningModel)
		if err != nil {
			return report, fmt.Errorf("case %s screening: %w", c.ID, err)
		}
		assessment, err := AssessEvidence(c.Idea, screening, c.Bundles, c.EvidenceRules)
		if err != nil {
			return report, fmt.Errorf("case %s assessment: %w", c.ID, err)
		}
		result.EvidenceMatches, result.EvidenceDivergences = compareBenchmarkEvidence(assessment, c.ExpectedEvidence)
		decision, err := DecideReuse(domain.ReuseDecisionRequest{Idea: c.Idea, Assessments: assessment.Assessments, Policy: c.ReusePolicy})
		if err != nil {
			return report, fmt.Errorf("case %s decision: %w", c.ID, err)
		}
		result.DecisionObserved = decision.Decision
		result.DecisionAccepted = decision.Decision == c.ExpectedDecision.Expected || containsDecision(c.ExpectedDecision.Acceptable, decision.Decision)
		result.ForbiddenDecision = containsDecision(c.ExpectedDecision.Forbidden, decision.Decision)
		if result.DecisionAccepted {
			result.Utility = c.UtilityJudgment.Label
		} else {
			result.Utility = domain.UtilityNotUseful
		}
		report.Cases = append(report.Cases, result)
	}
	sort.Slice(report.Cases, func(i, j int) bool { return report.Cases[i].CaseID < report.Cases[j].CaseID })
	report.Metrics = benchmarkMetrics(report.Cases)
	report.Passed = benchmarkGates(report.Metrics)
	return report, nil
}

func compareBenchmarkEvidence(a domain.AssessmentReport, expected []domain.EvidenceJudgment) (int, int) {
	observed := map[string]domain.EvidenceStatus{}
	for _, x := range a.Assessments {
		key := x.Repository.Owner + "/" + x.Repository.Repo
		for _, e := range x.Capabilities {
			observed[key+"|"+normalizeEvidenceText(e.Capability)] = e.Status
		}
	}
	matches, div := 0, 0
	for _, e := range expected {
		if observed[repositoryIdentityKey(e.Repository)+"|"+normalizeEvidenceText(e.Capability)] == e.Expected {
			matches++
		} else {
			div++
		}
	}
	return matches, div
}
func benchmarkMetrics(cases []domain.BenchmarkCaseResult) domain.EndToEndBenchmarkMetrics {
	m := domain.EndToEndBenchmarkMetrics{}
	if len(cases) == 0 {
		return m
	}
	preserved, unknown, accepted, forbidden, conclusive, correctUnknown, useful, judgeable := 0, 0, 0, 0, 0, 0, 0, 0
	eMatches, eTotal := 0, 0
	for _, c := range cases {
		if c.Correspondence == domain.CorrespondencePreserved {
			preserved++
		}
		if c.Correspondence == domain.CorrespondenceUnknown {
			unknown++
		}
		if c.Correspondence != domain.CorrespondencePreserved {
			continue
		}
		if c.DecisionAccepted {
			accepted++
		}
		if c.ForbiddenDecision {
			forbidden++
		}
		if c.DecisionObserved != domain.DecisionUnknown {
			conclusive++
		} else if c.DecisionAccepted {
			correctUnknown++
		}
		eMatches += c.EvidenceMatches
		eTotal += c.EvidenceMatches + c.EvidenceDivergences
		if c.Utility != domain.UtilityUnjudgeable {
			judgeable++
			if c.Utility == domain.UtilityUseful {
				useful++
			}
		}
	}
	n := float64(len(cases))
	m.CorrespondencePreservedRate = float64(preserved) / n
	m.CorrespondenceUnknownRate = float64(unknown) / n
	m.ForbiddenDecisionRate = float64(forbidden) / n
	m.ConclusiveCoverage = float64(conclusive) / n
	if preserved > 0 {
		x := float64(accepted) / float64(preserved)
		m.DecisionAccuracy = &x
	}
	if eTotal > 0 {
		x := float64(eMatches) / float64(eTotal)
		m.EvidencePrecision = &x
		m.EvidenceRecall = &x
	}
	unknownObserved := len(cases) - conclusive
	if unknownObserved > 0 {
		x := float64(correctUnknown) / float64(unknownObserved)
		m.CorrectInconclusiveRate = &x
	}
	if judgeable > 0 {
		x := float64(useful) / float64(judgeable)
		m.UtilityRate = &x
	}
	return m
}
func benchmarkGates(m domain.EndToEndBenchmarkMetrics) bool {
	return m.CorrespondencePreservedRate == 1 && m.CorrespondenceUnknownRate == 0 && m.EvidencePrecision != nil && *m.EvidencePrecision >= .9 && m.EvidenceRecall != nil && *m.EvidenceRecall >= .85 && m.DecisionAccuracy != nil && *m.DecisionAccuracy >= .8 && m.ForbiddenDecisionRate == 0 && m.UtilityRate != nil && *m.UtilityRate >= .8
}
func containsDecision(xs []domain.ReuseDecision, x domain.ReuseDecision) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
func CanonicalBenchmarkJSON(v any) []byte { b, _ := json.Marshal(v); return b }
func canonicalizeBenchmarkCase(c *domain.EndToEndBenchmarkCase) {
	sort.Slice(c.Discovery.Candidates, func(i, j int) bool {
		return c.Discovery.Candidates[i].Repository.Owner+"/"+c.Discovery.Candidates[i].Repository.Repo < c.Discovery.Candidates[j].Repository.Owner+"/"+c.Discovery.Candidates[j].Repository.Repo
	})
	sort.Slice(c.Bundles, func(i, j int) bool { return c.Bundles[i].Repository.Key() < c.Bundles[j].Repository.Key() })
	for i := range c.Bundles {
		sort.Slice(c.Bundles[i].Items, func(a, b int) bool { return c.Bundles[i].Items[a].Path < c.Bundles[i].Items[b].Path })
	}
	sort.Slice(c.ExpectedEvidence, func(i, j int) bool {
		return c.ExpectedEvidence[i].Repository+"|"+c.ExpectedEvidence[i].Capability < c.ExpectedEvidence[j].Repository+"|"+c.ExpectedEvidence[j].Capability
	})
}
