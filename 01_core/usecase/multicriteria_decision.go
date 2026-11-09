/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/multicriteria_decision.md
 * @layer L1
 * @updated 2026-11-08
 */

package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

const multicriteriaDecisionMethod = "GATES_THEN_INTERVAL_PARETO_WITH_DECLARED_SCENARIOS_V1"

// AnalyzeMulticriteriaDecision preserves hard gates, computes interval-aware
// Pareto relations and evaluates only preference scenarios declared by the caller.
func AnalyzeMulticriteriaDecision(request domain.MulticriteriaDecisionRequest) (domain.MulticriteriaDecisionReport, error) {
	if err := validateMulticriteriaRequest(request); err != nil {
		return domain.MulticriteriaDecisionReport{}, err
	}
	report := domain.MulticriteriaDecisionReport{
		SchemaVersion: domain.MulticriteriaDecisionSchemaVersion,
		Method:        multicriteriaDecisionMethod,
		Eligibility:   []domain.CandidateEligibility{}, Frontier: []string{}, RobustFrontier: []string{},
		Dominated: []domain.DominatedCandidate{}, Pairwise: []domain.PairwiseRelation{},
		Scenarios: []domain.PreferenceScenarioResult{}, Recommended: []string{}, Questions: []domain.MulticriteriaQuestion{},
		Sensitivity: domain.RankingSensitivity{RobustWinners: []string{}, PossibleWinners: []string{}},
	}
	criteria := canonicalCriteria(request.Criteria)
	composition, err := multicriteriaCompositions(request)
	if err != nil {
		return domain.MulticriteriaDecisionReport{}, err
	}
	report.Composition = &composition
	measurements := effectiveMeasurements(request, criteria, composition.Candidates)
	eligible, inconclusive := evaluateCandidateEligibility(request, &report, composition.Candidates)
	if len(eligible) == 0 {
		if inconclusive {
			report.StopReason = "ELIGIBILITY_INCONCLUSIVE"
		} else {
			report.StopReason = "NO_ELIGIBLE_CANDIDATES"
		}
		if composition.Truncated {
			report.StopReason = "COMPOSITION_SEARCH_LIMIT"
		}
		return report, nil
	}

	dominatedBy := map[string][]string{}
	possiblyDominated := map[string]bool{}
	for left := 0; left < len(eligible); left++ {
		for right := left + 1; right < len(eligible); right++ {
			relation, leftCertain, rightCertain, leftPossible, rightPossible := compareCandidates(eligible[left], eligible[right], criteria, measurements)
			report.Pairwise = append(report.Pairwise, relation)
			if leftCertain {
				dominatedBy[eligible[right]] = append(dominatedBy[eligible[right]], eligible[left])
			}
			if rightCertain {
				dominatedBy[eligible[left]] = append(dominatedBy[eligible[left]], eligible[right])
			}
			if leftPossible {
				possiblyDominated[eligible[right]] = true
			}
			if rightPossible {
				possiblyDominated[eligible[left]] = true
			}
		}
	}
	for _, candidate := range eligible {
		if len(dominatedBy[candidate]) == 0 {
			report.Frontier = append(report.Frontier, candidate)
		} else {
			sort.Strings(dominatedBy[candidate])
			report.Dominated = append(report.Dominated, domain.DominatedCandidate{Candidate: candidate, DominatedBy: dominatedBy[candidate]})
		}
		if !possiblyDominated[candidate] {
			report.RobustFrontier = append(report.RobustFrontier, candidate)
		}
	}

	for _, scenario := range canonicalScenarios(request.Scenarios) {
		report.Scenarios = append(report.Scenarios, evaluatePreferenceScenario(scenario, report.Frontier, criteria, measurements))
	}
	report.Sensitivity = summarizeSensitivity(report.Scenarios)
	report.Questions = buildMulticriteriaQuestions(report.Pairwise, report.Frontier)
	finishMulticriteriaReport(&report)
	if composition.Truncated {
		report.Recommended = []string{}
		report.StopReason = "COMPOSITION_SEARCH_LIMIT"
	}
	return report, nil
}

func validateMulticriteriaRequest(request domain.MulticriteriaDecisionRequest) error {
	if err := validateDecisionRequest(domain.ReuseDecisionRequest{Idea: request.Idea, Assessments: request.Assessment.Assessments, Policy: request.Policy}); err != nil {
		return err
	}
	if len(request.Assessment.Assessments) == 0 {
		return fmt.Errorf("assessment exige ao menos um candidato")
	}
	if len(request.Criteria) == 0 {
		return fmt.Errorf("criteria exige ao menos um critério")
	}
	candidates := map[string]domain.CommitRef{}
	for _, assessment := range request.Assessment.Assessments {
		key := assessment.Repository.Key()
		if _, duplicate := candidates[key]; strings.TrimSpace(key) == "/" || duplicate {
			return fmt.Errorf("candidato inválido ou duplicado: %q", key)
		}
		candidates[key] = assessment.Repository
	}
	composition, err := multicriteriaCompositions(request)
	if err != nil {
		return err
	}
	composed := map[string]domain.CompositionCandidate{}
	for _, candidate := range composition.Candidates {
		composed[candidate.ID] = candidate
		candidates[candidate.ID] = domain.CommitRef{}
	}
	criteria := map[string]domain.MulticriteriaCriterion{}
	for _, criterion := range request.Criteria {
		if strings.TrimSpace(criterion.ID) == "" || criteria[criterion.ID].ID != "" {
			return fmt.Errorf("criterion.id vazio ou duplicado: %q", criterion.ID)
		}
		if criterion.Direction != domain.CriterionMaximize && criterion.Direction != domain.CriterionMinimize {
			return fmt.Errorf("direção inválida no critério %s", criterion.ID)
		}
		if !finite(criterion.ScaleMin) || !finite(criterion.ScaleMax) || criterion.ScaleMin >= criterion.ScaleMax {
			return fmt.Errorf("escala inválida no critério %s", criterion.ID)
		}
		criteria[criterion.ID] = criterion
	}
	seenMeasurements := map[string]bool{}
	for _, measurement := range request.Measurements {
		candidateRef, candidateExists := candidates[measurement.Candidate]
		if !candidateExists {
			return fmt.Errorf("medição referencia candidato desconhecido: %s", measurement.Candidate)
		}
		criterion, ok := criteria[measurement.CriterionID]
		if !ok {
			return fmt.Errorf("medição referencia critério desconhecido: %s", measurement.CriterionID)
		}
		identity := measurement.Candidate + "\x00" + measurement.CriterionID
		if seenMeasurements[identity] {
			return fmt.Errorf("medição duplicada para %s/%s", measurement.Candidate, measurement.CriterionID)
		}
		seenMeasurements[identity] = true
		if !validEvidenceStatus(measurement.Status) {
			return fmt.Errorf("status inválido para %s/%s", measurement.Candidate, measurement.CriterionID)
		}
		if measurement.Status == domain.EvidenceProven || measurement.Status == domain.EvidencePartial {
			if !finite(measurement.Interval.Low) || !finite(measurement.Interval.High) || measurement.Interval.Low > measurement.Interval.High || measurement.Interval.Low < criterion.ScaleMin || measurement.Interval.High > criterion.ScaleMax {
				return fmt.Errorf("intervalo inválido para %s/%s", measurement.Candidate, measurement.CriterionID)
			}
			if measurement.Status == domain.EvidenceProven && measurement.Interval.Low != measurement.Interval.High {
				return fmt.Errorf("medição PROVEN deve ser pontual para %s/%s", measurement.Candidate, measurement.CriterionID)
			}
			if composedCandidate, ok := composed[measurement.Candidate]; ok && !compositionEvidenceAdmissible(compositionMemberMap(composedCandidate), measurement.Sources, false) || !ok && !EvidenceSourcesAdmissible(candidateRef, measurement.Sources, false) {
				return fmt.Errorf("medição %s exige fonte para %s/%s", measurement.Status, measurement.Candidate, measurement.CriterionID)
			}
		}
	}
	seenScenarios := map[string]bool{}
	for _, scenario := range request.Scenarios {
		if strings.TrimSpace(scenario.ID) == "" || seenScenarios[scenario.ID] {
			return fmt.Errorf("preference_scenario.id vazio ou duplicado: %q", scenario.ID)
		}
		seenScenarios[scenario.ID] = true
		weights := map[string]bool{}
		total := 0.0
		for _, weight := range scenario.Weights {
			if _, ok := criteria[weight.CriterionID]; !ok || weights[weight.CriterionID] {
				return fmt.Errorf("peso desconhecido ou duplicado em %s: %s", scenario.ID, weight.CriterionID)
			}
			if !finite(weight.Weight) || weight.Weight < 0 || weight.Weight > 1 {
				return fmt.Errorf("peso inválido em %s/%s", scenario.ID, weight.CriterionID)
			}
			weights[weight.CriterionID] = true
			total += weight.Weight
		}
		if len(weights) != len(criteria) || math.Abs(total-1) > 1e-9 {
			return fmt.Errorf("cenário %s deve declarar todos os critérios com soma 1", scenario.ID)
		}
	}
	return nil
}

func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func canonicalCriteria(input []domain.MulticriteriaCriterion) []domain.MulticriteriaCriterion {
	out := append([]domain.MulticriteriaCriterion(nil), input...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func canonicalScenarios(input []domain.PreferenceScenario) []domain.PreferenceScenario {
	out := append([]domain.PreferenceScenario(nil), input...)
	for index := range out {
		out[index].Weights = append([]domain.CriterionWeight(nil), out[index].Weights...)
		sort.Slice(out[index].Weights, func(i, j int) bool { return out[index].Weights[i].CriterionID < out[index].Weights[j].CriterionID })
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func effectiveMeasurements(request domain.MulticriteriaDecisionRequest, criteria []domain.MulticriteriaCriterion, compositions []domain.CompositionCandidate) map[string]map[string]domain.CriterionMeasurement {
	out := map[string]map[string]domain.CriterionMeasurement{}
	for _, assessment := range request.Assessment.Assessments {
		candidate := assessment.Repository.Key()
		out[candidate] = map[string]domain.CriterionMeasurement{}
		for _, criterion := range criteria {
			out[candidate][criterion.ID] = domain.CriterionMeasurement{Candidate: candidate, CriterionID: criterion.ID, Status: domain.EvidenceInconclusive, Interval: domain.CriterionInterval{Low: criterion.ScaleMin, High: criterion.ScaleMax}}
		}
	}
	for _, composition := range compositions {
		out[composition.ID] = map[string]domain.CriterionMeasurement{}
		for _, criterion := range criteria {
			out[composition.ID][criterion.ID] = domain.CriterionMeasurement{Candidate: composition.ID, CriterionID: criterion.ID, Status: domain.EvidenceInconclusive, Interval: domain.CriterionInterval{Low: criterion.ScaleMin, High: criterion.ScaleMax}}
		}
	}
	for _, measurement := range request.Measurements {
		criterion := criterionByID(criteria, measurement.CriterionID)
		copy := measurement
		copy.Sources = append([]domain.EvidenceSource(nil), measurement.Sources...)
		if measurement.Status != domain.EvidenceProven && measurement.Status != domain.EvidencePartial {
			copy.Interval = domain.CriterionInterval{Low: criterion.ScaleMin, High: criterion.ScaleMax}
		}
		out[measurement.Candidate][measurement.CriterionID] = copy
	}
	return out
}

func evaluateCandidateEligibility(request domain.MulticriteriaDecisionRequest, report *domain.MulticriteriaDecisionReport, compositions []domain.CompositionCandidate) ([]string, bool) {
	assessments := append([]domain.CandidateAssessment(nil), request.Assessment.Assessments...)
	sort.Slice(assessments, func(i, j int) bool { return assessments[i].Repository.Key() < assessments[j].Repository.Key() })
	capabilities := normalizedSorted(request.Idea.Capabilities)
	eligible := []string{}
	hasInconclusive := false
	for _, assessment := range assessments {
		prepared := SanitizeAssessmentEvidence(assessment)
		candidate := inspectPreparedCandidate(prepared, capabilities, request.Policy)
		entry := domain.CandidateEligibility{Candidate: candidate.key, Status: domain.CandidateEligible, Reasons: []domain.DecisionReason{}}
		if len(candidate.blockers) > 0 {
			entry.Status = domain.CandidateIneligible
			entry.Reasons = append(entry.Reasons, candidate.blockers...)
		} else if candidate.ambiguous {
			entry.Status = domain.CandidateEligibilityInconclusive
			entry.Reasons = append(entry.Reasons, unresolvedCandidateReasons(candidate, request.Policy)...)
			hasInconclusive = true
		} else {
			eligible = append(eligible, candidate.key)
		}
		report.Eligibility = append(report.Eligibility, entry)
	}
	for _, composition := range compositions {
		entry := domain.CandidateEligibility{Candidate: composition.ID, Reasons: []domain.DecisionReason{}}
		switch composition.Status {
		case domain.CompositionVerified:
			entry.Status = domain.CandidateEligible
			eligible = append(eligible, composition.ID)
		case domain.CompositionIncompatible:
			entry.Status = domain.CandidateIneligible
			entry.Reasons = append(entry.Reasons, domain.DecisionReason{RuleID: "MCDA-COMPOSE-CONFLICT-0088", Effect: "block", Candidate: composition.ID, Message: "composição incompatível pelas restrições ou ensaio declarado"})
		default:
			entry.Status = domain.CandidateEligibilityInconclusive
			entry.Reasons = append(entry.Reasons, domain.DecisionReason{RuleID: "MCDA-COMPOSE-UNKNOWN-0088", Effect: "verify", Candidate: composition.ID, Message: "compatibilidade ou integração não comprovada"})
			hasInconclusive = true
		}
		report.Eligibility = append(report.Eligibility, entry)
	}
	sort.Strings(eligible)
	return eligible, hasInconclusive
}

func compareCandidates(left, right string, criteria []domain.MulticriteriaCriterion, measurements map[string]map[string]domain.CriterionMeasurement) (domain.PairwiseRelation, bool, bool, bool, bool) {
	relation := domain.PairwiseRelation{Left: left, Right: right, LeftBetter: []string{}, RightBetter: []string{}, UnresolvedCriteria: []string{}, Comparisons: []domain.PairwiseCriterionComparison{}}
	leftCertain, leftStrict := true, false
	rightCertain, rightStrict := true, false
	leftPossible, leftPossibleStrict := true, false
	rightPossible, rightPossibleStrict := true, false
	allExactEqual := true
	for _, criterion := range criteria {
		leftMeasurement := measurements[left][criterion.ID]
		rightMeasurement := measurements[right][criterion.ID]
		leftInterval, rightInterval := leftMeasurement.Interval, rightMeasurement.Interval
		relation.Comparisons = append(relation.Comparisons, domain.PairwiseCriterionComparison{CriterionID: criterion.ID, Left: leftInterval, Right: rightInterval})
		leftNoWorse, leftBetter, leftCanNoWorse, leftCanBetter := intervalPreference(leftInterval, rightInterval, criterion.Direction)
		rightNoWorse, rightBetter, rightCanNoWorse, rightCanBetter := intervalPreference(rightInterval, leftInterval, criterion.Direction)
		leftCertain = leftCertain && leftNoWorse
		leftStrict = leftStrict || leftBetter
		rightCertain = rightCertain && rightNoWorse
		rightStrict = rightStrict || rightBetter
		leftPossible = leftPossible && leftCanNoWorse
		leftPossibleStrict = leftPossibleStrict || leftCanBetter
		rightPossible = rightPossible && rightCanNoWorse
		rightPossibleStrict = rightPossibleStrict || rightCanBetter
		if leftBetter {
			relation.LeftBetter = append(relation.LeftBetter, criterion.ID)
		}
		if rightBetter {
			relation.RightBetter = append(relation.RightBetter, criterion.ID)
		}
		exactEqual := leftMeasurement.Status == domain.EvidenceProven && rightMeasurement.Status == domain.EvidenceProven && leftInterval.Low == rightInterval.Low
		allExactEqual = allExactEqual && exactEqual
		if leftMeasurement.Status != domain.EvidenceProven || rightMeasurement.Status != domain.EvidenceProven || leftInterval.Low != leftInterval.High || rightInterval.Low != rightInterval.High {
			relation.UnresolvedCriteria = append(relation.UnresolvedCriteria, criterion.ID)
		}
	}
	leftDominates := leftCertain && leftStrict
	rightDominates := rightCertain && rightStrict
	leftMayDominate := leftPossible && leftPossibleStrict
	rightMayDominate := rightPossible && rightPossibleStrict
	switch {
	case leftDominates:
		relation.Relation = domain.PairwiseDominates
	case rightDominates:
		relation.Relation = domain.PairwiseDominatedBy
	case allExactEqual:
		relation.Relation = domain.PairwiseEquivalent
	case len(relation.UnresolvedCriteria) > 0:
		relation.Relation = domain.PairwiseUncertain
	default:
		relation.Relation = domain.PairwiseIncomparable
	}
	return relation, leftDominates, rightDominates, leftMayDominate, rightMayDominate
}

func intervalPreference(left, right domain.CriterionInterval, direction domain.CriterionDirection) (certainNoWorse, certainBetter, possibleNoWorse, possibleBetter bool) {
	if direction == domain.CriterionMaximize {
		return left.Low >= right.High, left.Low > right.High, left.High >= right.Low, left.High > right.Low
	}
	return left.High <= right.Low, left.High < right.Low, left.Low <= right.High, left.Low < right.High
}

func evaluatePreferenceScenario(scenario domain.PreferenceScenario, candidates []string, criteria []domain.MulticriteriaCriterion, measurements map[string]map[string]domain.CriterionMeasurement) domain.PreferenceScenarioResult {
	weights := map[string]float64{}
	for _, weight := range scenario.Weights {
		weights[weight.CriterionID] = weight.Weight
	}
	result := domain.PreferenceScenarioResult{ScenarioID: scenario.ID, Scores: []domain.ScenarioCandidateScore{}, NecessaryWinners: []string{}, PossibleWinners: []string{}}
	for _, candidate := range candidates {
		score := domain.ScenarioCandidateScore{Candidate: candidate, Contributions: []domain.CriterionContribution{}}
		for _, criterion := range criteria {
			utility := normalizeCriterionUtility(measurements[candidate][criterion.ID].Interval, criterion)
			weighted := domain.CriterionInterval{Low: utility.Low * weights[criterion.ID], High: utility.High * weights[criterion.ID]}
			score.Score.Low += weighted.Low
			score.Score.High += weighted.High
			score.Contributions = append(score.Contributions, domain.CriterionContribution{CriterionID: criterion.ID, Weight: weights[criterion.ID], Utility: utility, Weighted: weighted})
		}
		result.Scores = append(result.Scores, score)
	}
	for _, score := range result.Scores {
		necessary, possible := true, true
		for _, other := range result.Scores {
			if other.Candidate == score.Candidate {
				continue
			}
			necessary = necessary && score.Score.Low >= other.Score.High
			possible = possible && score.Score.High >= other.Score.Low
		}
		if necessary {
			result.NecessaryWinners = append(result.NecessaryWinners, score.Candidate)
		}
		if possible {
			result.PossibleWinners = append(result.PossibleWinners, score.Candidate)
		}
	}
	return result
}

func normalizeCriterionUtility(interval domain.CriterionInterval, criterion domain.MulticriteriaCriterion) domain.CriterionInterval {
	span := criterion.ScaleMax - criterion.ScaleMin
	if math.IsInf(span, 0) {
		// Finite opposite extremes can overflow their difference; rescale first.
		interval.Low, interval.High = interval.Low/2, interval.High/2
		criterion.ScaleMin, criterion.ScaleMax = criterion.ScaleMin/2, criterion.ScaleMax/2
		span = criterion.ScaleMax - criterion.ScaleMin
	}
	if criterion.Direction == domain.CriterionMaximize {
		return domain.CriterionInterval{Low: (interval.Low - criterion.ScaleMin) / span, High: (interval.High - criterion.ScaleMin) / span}
	}
	return domain.CriterionInterval{Low: (criterion.ScaleMax - interval.High) / span, High: (criterion.ScaleMax - interval.Low) / span}
}

func summarizeSensitivity(results []domain.PreferenceScenarioResult) domain.RankingSensitivity {
	summary := domain.RankingSensitivity{RobustWinners: []string{}, PossibleWinners: []string{}}
	if len(results) == 0 {
		return summary
	}
	possible := map[string]bool{}
	robust := map[string]bool{}
	for _, candidate := range results[0].NecessaryWinners {
		robust[candidate] = true
	}
	firstNecessary := strings.Join(results[0].NecessaryWinners, "\x00")
	for _, result := range results {
		current := stringSet(result.NecessaryWinners)
		for candidate := range robust {
			if !current[candidate] {
				delete(robust, candidate)
			}
		}
		if strings.Join(result.NecessaryWinners, "\x00") != firstNecessary {
			summary.PreferenceSensitive = true
		}
		for _, candidate := range result.PossibleWinners {
			possible[candidate] = true
		}
	}
	for candidate := range robust {
		summary.RobustWinners = append(summary.RobustWinners, candidate)
	}
	for candidate := range possible {
		summary.PossibleWinners = append(summary.PossibleWinners, candidate)
	}
	sort.Strings(summary.RobustWinners)
	sort.Strings(summary.PossibleWinners)
	return summary
}

func buildMulticriteriaQuestions(relations []domain.PairwiseRelation, frontier []string) []domain.MulticriteriaQuestion {
	onFrontier := stringSet(frontier)
	questions := []domain.MulticriteriaQuestion{}
	for _, relation := range relations {
		if !onFrontier[relation.Left] || !onFrontier[relation.Right] {
			continue
		}
		switch relation.Relation {
		case domain.PairwiseIncomparable:
			criteria := append(append([]string{}, relation.LeftBetter...), relation.RightBetter...)
			questions = append(questions, domain.MulticriteriaQuestion{Kind: "PREFERENCE_TRADEOFF", Candidates: []string{relation.Left, relation.Right}, Criteria: criteria, Question: fmt.Sprintf("Entre %s (melhor em %s) e %s (melhor em %s), qual trade-off você prefere?", relation.Left, strings.Join(relation.LeftBetter, ", "), relation.Right, strings.Join(relation.RightBetter, ", "))})
		case domain.PairwiseUncertain:
			questions = append(questions, domain.MulticriteriaQuestion{Kind: "EVIDENCE_GAP", Candidates: []string{relation.Left, relation.Right}, Criteria: append([]string{}, relation.UnresolvedCriteria...), Question: fmt.Sprintf("Obter medições mais precisas de %s para comparar %s e %s.", strings.Join(relation.UnresolvedCriteria, ", "), relation.Left, relation.Right)})
		}
	}
	return questions
}

func finishMulticriteriaReport(report *domain.MulticriteriaDecisionReport) {
	switch {
	case len(report.Frontier) == 1:
		report.Recommended = append(report.Recommended, report.Frontier[0])
		report.StopReason = "SINGLE_PARETO_CANDIDATE"
	case len(report.Scenarios) == 0:
		report.StopReason = "PREFERENCES_REQUIRED"
	case len(report.Sensitivity.RobustWinners) > 0:
		report.Recommended = append(report.Recommended, report.Sensitivity.RobustWinners...)
		report.StopReason = "ROBUST_WINNER"
	case report.Sensitivity.PreferenceSensitive:
		report.StopReason = "PREFERENCE_SENSITIVE"
	default:
		report.StopReason = "EVIDENCE_INCONCLUSIVE"
	}
}

func criterionByID(criteria []domain.MulticriteriaCriterion, id string) domain.MulticriteriaCriterion {
	for _, criterion := range criteria {
		if criterion.ID == id {
			return criterion
		}
	}
	return domain.MulticriteriaCriterion{}
}

func stringSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		out[value] = true
	}
	return out
}
