/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/probabilistic_screening.md
 * @layer L1
 * @updated 2026-09-07
 */

package usecase

import (
	"crypto/sha256"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"bateia/01_core/domain"
)

var screeningStopwords = map[string]bool{
	"a": true, "o": true, "e": true, "de": true, "do": true, "da": true, "para": true, "uma": true, "um": true,
	"the": true, "an": true, "and": true, "or": true, "of": true, "to": true, "for": true, "with": true, "in": true,
	"lang": true, "go": true,
}

func screeningTokens(value string) []string {
	fields := strings.FieldsFunc(strings.ToLower(value), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	result := make([]string, 0, len(fields))
	seen := make(map[string]bool)
	for _, field := range fields {
		if len([]rune(field)) > 1 && !screeningStopwords[field] && !seen[field] {
			seen[field] = true
			result = append(result, field)
		}
	}
	return result
}

func tokenOverlap(left, right string) float64 {
	a, b := screeningTokens(left), screeningTokens(right)
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	set := make(map[string]bool, len(b))
	for _, token := range b {
		set[token] = true
	}
	found := 0
	for _, token := range a {
		if set[token] {
			found++
		}
	}
	return ratio(found, len(a))
}

func candidateInList(candidate domain.RepositoryCandidate, values []domain.RepositoryCandidate) bool {
	key := repositoryKey(candidate)
	for _, value := range values {
		if repositoryKey(value) == key {
			return true
		}
	}
	return false
}

// MigrateDiscoverySnapshotV2 reconstrói somente origens comprováveis no schema v1.
func MigrateDiscoverySnapshotV2(snapshot domain.DiscoverySnapshot) (domain.DiscoverySnapshot, error) {
	if snapshot.SchemaVersion != "1" && snapshot.SchemaVersion != "2" {
		return domain.DiscoverySnapshot{}, fmt.Errorf("schema de snapshot não suportado: %q", snapshot.SchemaVersion)
	}
	result := snapshot
	result.SchemaVersion = "2"
	result.Candidates = append([]domain.DiscoveredCandidate(nil), snapshot.Candidates...)
	for index := range result.Candidates {
		candidate := &result.Candidates[index]
		candidate.ExpansionIDs = cloneStrings(candidate.ExpansionIDs)
		candidate.Queries = cloneStrings(candidate.Queries)
		candidate.Origins = append([]domain.CandidateOrigin(nil), candidate.Origins...)
		if len(candidate.Origins) > 0 {
			continue
		}
		for _, discovery := range snapshot.PerExpansion {
			if !candidateInList(candidate.Repository, discovery.Provider.Candidates) {
				continue
			}
			attemptIndex := -1
			effective := ""
			for idx, attempt := range discovery.Provider.Attempts {
				if attempt.ResultCount > 0 {
					attemptIndex = idx
					effective = attempt.Query
					break
				}
			}
			candidate.Origins = append(candidate.Origins, domain.CandidateOrigin{
				ExpansionID: discovery.Expansion.ID, PlannedQuery: discovery.Expansion.Query,
				AttemptIndex: attemptIndex, EffectiveQuery: effective,
			})
		}
	}
	return result, nil
}

func DefaultScreeningModel() domain.ScreeningModel {
	return domain.ScreeningModel{SchemaVersion: "1", Weights: map[string]float64{
		"intercept": 0.05, "expansion_prior": 0.15, "query_retention": 0.15, "fallback_depth": -0.10,
		"match_specificity": 0.10, "description_alignment": 0.15, "origin_convergence": 0.10,
		"capability_coverage": 0.15, "non_goal_conflict": -0.40, "evidence_completeness": 0.10,
	}, ExpansionPriors: map[string]float64{}, AdvanceThreshold: 0.55, DropThreshold: 0.25, AuditSampleRate: 0.10, AuditSeed: 42, CalibrationStatus: "uncalibrated"}
}

func ideaText(spec domain.IdeaSpec) string {
	return strings.Join(append([]string{spec.Problem}, spec.Capabilities...), " ")
}

func expansionPrior(origins []domain.CandidateOrigin, model domain.ScreeningModel) float64 {
	if len(origins) == 0 {
		return 0
	}
	total := 0.0
	for _, origin := range origins {
		if value, ok := model.ExpansionPriors[origin.ExpansionID]; ok {
			total += value
		}
	}
	return total / float64(len(origins))
}

func ExtractScreeningSignals(snapshot domain.DiscoverySnapshot, candidate domain.DiscoveredCandidate, model domain.ScreeningModel) domain.ScreeningSignals {
	signals := domain.ScreeningSignals{ExpansionPrior: expansionPrior(candidate.Origins, model), FallbackDepth: 2}
	for _, origin := range candidate.Origins {
		retention := tokenOverlap(origin.PlannedQuery, origin.EffectiveQuery)
		if retention > signals.QueryRetention {
			signals.QueryRetention = retention
		}
		if origin.AttemptIndex >= 0 && origin.AttemptIndex < signals.FallbackDepth {
			signals.FallbackDepth = origin.AttemptIndex
		}
		for _, evidence := range origin.Evidence {
			value := 0.25
			switch evidence.MatchType {
			case "symbol":
				value = 1
			case "file":
				value = 0.75
			case "repository":
				value = 0.5
			}
			path := strings.ToLower(evidence.Path)
			if strings.Contains(path, "vendor/") || strings.Contains(path, "node_modules/") {
				value = math.Min(value, 0.1)
			}
			if strings.Contains(path, "locales/") || strings.Contains(path, "translations/") {
				value = math.Min(value, 0.2)
			}
			if value > signals.MatchSpecificity {
				signals.MatchSpecificity = value
			}
		}
	}
	signals.DescriptionAlignment = tokenOverlap(ideaText(snapshot.Idea), candidate.Repository.Description)
	signals.OriginConvergence = len(candidate.Origins)
	// Queries describe the search, not the candidate. Use only returned content.
	observed := candidate.Repository.Description
	for _, origin := range candidate.Origins {
		for _, evidence := range origin.Evidence {
			observed += " " + evidence.Snippet
		}
	}
	covered := 0
	for _, capability := range snapshot.Idea.Capabilities {
		if tokenOverlap(capability, observed) > 0 {
			covered++
		}
	}
	signals.CapabilityCoverage = ratio(covered, len(snapshot.Idea.Capabilities))
	for _, nonGoal := range snapshot.Idea.NonGoals {
		if containsPhrase(observed, nonGoal) {
			signals.NonGoalConflict += ratio(1, len(snapshot.Idea.NonGoals))
		}
	}
	if strings.TrimSpace(candidate.Repository.Description) != "" {
		signals.EvidenceCompleteness += 0.4
	}
	if len(candidate.Origins) > 0 && candidate.Origins[0].EffectiveQuery != "" {
		signals.EvidenceCompleteness += 0.3
	}
	for _, origin := range candidate.Origins {
		if len(origin.Evidence) > 0 {
			signals.EvidenceCompleteness += 0.3
			break
		}
	}
	return signals
}

// DiagnoseScreening aplica rótulos conhecidos apenas em replay diagnóstico.
func DiagnoseScreening(report domain.ScreeningReport, judgments domain.DiscoveryJudgments) domain.ScreeningReport {
	result := report
	result.Decisions = append([]domain.ScreeningDecision(nil), report.Decisions...)
	labels := judgmentLabels(judgments)
	for index := range result.Decisions {
		decision := &result.Decisions[index]
		label, ok := labels[strings.ToLower(decision.Repository)]
		if !ok {
			continue
		}
		if label == domain.RelevanceRelevant {
			decision.Decision = "ADVANCE"
			decision.UncertaintyReasons = append(decision.UncertaintyReasons, "override diagnóstico por rótulo relevante conhecido")
		}
		if label == domain.RelevanceUncertain {
			decision.Decision = "AUDIT"
			decision.UncertaintyReasons = append(decision.UncertaintyReasons, "override diagnóstico por incerteza conhecida")
		}
	}
	return result
}

func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
func sigmoid(value float64) float64 {
	if value > 30 {
		return 1
	}
	if value < -30 {
		return 0
	}
	return 1 / (1 + math.Exp(-value))
}

func contribution(signal string, value float64, model domain.ScreeningModel, reason string) domain.SignalContribution {
	weight := model.Weights[signal]
	return domain.SignalContribution{Signal: signal, Value: value, Weight: weight, Effect: value * weight, Reason: reason}
}

func deterministicAudit(repository string, seed int64, rate float64) bool {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d\x00%s", seed, strings.ToLower(repository))))
	value := float64(uint16(hash[0])<<8|uint16(hash[1])) / 65535
	return value < rate
}

func screenCandidate(snapshot domain.DiscoverySnapshot, candidate domain.DiscoveredCandidate, model domain.ScreeningModel) domain.ScreeningDecision {
	signals := ExtractScreeningSignals(snapshot, candidate, model)
	values := []domain.SignalContribution{
		contribution("intercept", 1, model, "base explícita do modelo"),
		contribution("expansion_prior", signals.ExpansionPrior, model, "produtividade julgada da expansão"),
		contribution("query_retention", signals.QueryRetention, model, "termos informativos preservados"),
		contribution("fallback_depth", float64(signals.FallbackDepth), model, "profundidade do fallback efetivo"),
		contribution("match_specificity", signals.MatchSpecificity, model, "especificidade do tipo de match"),
		contribution("description_alignment", signals.DescriptionAlignment, model, "alinhamento lexical da descrição"),
		contribution("origin_convergence", math.Min(float64(signals.OriginConvergence)/3, 1), model, "convergência de expansões independentes"),
		contribution("capability_coverage", signals.CapabilityCoverage, model, "capacidades com sinal preliminar"),
		contribution("non_goal_conflict", signals.NonGoalConflict, model, "conflito explícito com não objetivos"),
		contribution("evidence_completeness", signals.EvidenceCompleteness, model, "completude da evidência disponível"),
	}
	raw := 0.0
	for _, item := range values {
		raw += item.Effect
	}
	score := clamp(raw)
	repository := candidate.Repository.Owner + "/" + candidate.Repository.Repo
	decision := domain.ScreeningDecision{Repository: repository, RankingScore: score, CalibrationStatus: model.CalibrationStatus, Signals: signals, Contributions: values, Origins: append([]domain.CandidateOrigin(nil), candidate.Origins...)}
	decisionValue := score
	if model.CalibrationStatus == "calibrated" && model.Calibration != nil {
		logitInput := math.Log(math.Max(score, 1e-6) / math.Max(1-score, 1e-6))
		probability := sigmoid(model.Calibration.A*logitInput + model.Calibration.B)
		decision.Probability = &probability
		decisionValue = probability
	}
	switch {
	case signals.NonGoalConflict > 0:
		decision.Decision = "DROP"
		decision.UncertaintyReasons = append(decision.UncertaintyReasons, "conflito com NonGoals")
	case signals.EvidenceCompleteness < 0.4:
		decision.Decision = "AUDIT"
		decision.UncertaintyReasons = append(decision.UncertaintyReasons, "evidência incompleta")
	case decisionValue >= model.AdvanceThreshold:
		decision.Decision = "ADVANCE"
	case decisionValue < model.DropThreshold:
		decision.Decision = "DROP"
		if deterministicAudit(repository, model.AuditSeed, model.AuditSampleRate) {
			decision.Decision = "AUDIT"
			decision.UncertaintyReasons = append(decision.UncertaintyReasons, "amostra determinística de descartes")
		}
	default:
		decision.Decision = "AUDIT"
		decision.UncertaintyReasons = append(decision.UncertaintyReasons, "próximo aos limiares")
	}
	return decision
}

func ScreenDiscovery(snapshot domain.DiscoverySnapshot, model domain.ScreeningModel) (domain.ScreeningReport, error) {
	migrated, err := MigrateDiscoverySnapshotV2(snapshot)
	if err != nil {
		return domain.ScreeningReport{}, err
	}
	if model.SchemaVersion != "1" || model.DropThreshold < 0 || model.AdvanceThreshold > 1 || model.DropThreshold >= model.AdvanceThreshold {
		return domain.ScreeningReport{}, fmt.Errorf("modelo de peneira inválido")
	}
	report := domain.ScreeningReport{SchemaVersion: "1", Strategy: snapshot.Strategy}
	for _, candidate := range migrated.Candidates {
		report.Decisions = append(report.Decisions, screenCandidate(migrated, candidate, model))
	}
	sort.SliceStable(report.Decisions, func(i, j int) bool {
		if report.Decisions[i].RankingScore == report.Decisions[j].RankingScore {
			return strings.ToLower(report.Decisions[i].Repository) < strings.ToLower(report.Decisions[j].Repository)
		}
		return report.Decisions[i].RankingScore > report.Decisions[j].RankingScore
	})
	return report, nil
}
