/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/probabilistic_screening.md
 * @layer L1
 * @updated 2026-08-24
 */

package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

type calibrationObservation struct {
	score float64
	label float64
}

type predictedObservation struct {
	probability float64
	label       float64
}

func judgmentLabels(judgments domain.DiscoveryJudgments) map[string]domain.RelevanceLabel {
	result := make(map[string]domain.RelevanceLabel)
	for _, judgment := range judgments.Judgments {
		result[strings.ToLower(strings.TrimSpace(judgment.Repository))] = judgment.Label
	}
	return result
}

func deriveExpansionPriors(groups []domain.ScreeningDatasetGroup) map[string]float64 {
	total, relevant := make(map[string]int), make(map[string]int)
	for _, group := range groups {
		labels := judgmentLabels(group.Judgments)
		snapshot, err := MigrateDiscoverySnapshotV2(group.Snapshot)
		if err != nil {
			continue
		}
		for _, candidate := range snapshot.Candidates {
			label, ok := labels[repositoryKey(candidate.Repository)]
			if !ok || label == domain.RelevanceUncertain {
				continue
			}
			for _, origin := range candidate.Origins {
				total[origin.ExpansionID]++
				if label == domain.RelevanceRelevant {
					relevant[origin.ExpansionID]++
				}
			}
		}
	}
	priors := make(map[string]float64, len(total))
	for id, count := range total {
		priors[id] = ratio(relevant[id], count)
	}
	return priors
}

func observations(groups []domain.ScreeningDatasetGroup, model domain.ScreeningModel) ([]calibrationObservation, error) {
	result := make([]calibrationObservation, 0)
	for _, group := range groups {
		report, err := ScreenDiscovery(group.Snapshot, model)
		if err != nil {
			return nil, fmt.Errorf("grupo %s: %w", group.IdeaID, err)
		}
		labels := judgmentLabels(group.Judgments)
		for _, decision := range report.Decisions {
			label, ok := labels[strings.ToLower(decision.Repository)]
			if !ok || label == domain.RelevanceUncertain {
				continue
			}
			y := 0.0
			if label == domain.RelevanceRelevant {
				y = 1
			}
			result = append(result, calibrationObservation{score: decision.RankingScore, label: y})
		}
	}
	return result, nil
}

func fitPlatt(values []calibrationObservation) domain.CalibrationParameters {
	a, b := 1.0, 0.0
	if len(values) == 0 {
		return domain.CalibrationParameters{}
	}
	for iteration := 0; iteration < 2000; iteration++ {
		gradientA, gradientB := 0.0, 0.0
		for _, value := range values {
			score := math.Min(math.Max(value.score, 1e-6), 1-1e-6)
			x := math.Log(score / (1 - score))
			predicted := sigmoid(a*x + b)
			gradientA += (predicted - value.label) * x
			gradientB += predicted - value.label
		}
		rate := 0.05 / float64(len(values))
		a -= rate * gradientA
		b -= rate * gradientB
	}
	return domain.CalibrationParameters{A: a, B: b}
}

func predictCalibration(value calibrationObservation, parameters domain.CalibrationParameters) float64 {
	score := math.Min(math.Max(value.score, 1e-6), 1-1e-6)
	return sigmoid(parameters.A*math.Log(score/(1-score)) + parameters.B)
}

func calibrationMetrics(values []calibrationObservation, parameters domain.CalibrationParameters) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	brier, loss := 0.0, 0.0
	for _, value := range values {
		p := math.Min(math.Max(predictCalibration(value, parameters), 1e-9), 1-1e-9)
		brier += (p - value.label) * (p - value.label)
		loss += -(value.label*math.Log(p) + (1-value.label)*math.Log(1-p))
	}
	return brier / float64(len(values)), loss / float64(len(values))
}

func prevalence(values []calibrationObservation) float64 {
	total := 0.0
	for _, value := range values {
		total += value.label
	}
	if len(values) == 0 {
		return 0
	}
	return total / float64(len(values))
}

func prevalenceBrier(values []calibrationObservation, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, v := range values {
		total += (p - v.label) * (p - v.label)
	}
	return total / float64(len(values))
}

func precisionAtK(values []calibrationObservation, parameters domain.CalibrationParameters, k int) float64 {
	copyValues := append([]calibrationObservation(nil), values...)
	sort.SliceStable(copyValues, func(i, j int) bool {
		return predictCalibration(copyValues[i], parameters) > predictCalibration(copyValues[j], parameters)
	})
	if len(copyValues) < k {
		k = len(copyValues)
	}
	if k == 0 {
		return 0
	}
	positives := 0.0
	for _, v := range copyValues[:k] {
		positives += v.label
	}
	return positives / float64(k)
}

func predictedPrecisionAtK(values []predictedObservation, k int) float64 {
	ordered := append([]predictedObservation(nil), values...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].probability > ordered[j].probability })
	if len(ordered) < k {
		k = len(ordered)
	}
	if k == 0 {
		return 0
	}
	positive := 0.0
	for _, value := range ordered[:k] {
		positive += value.label
	}
	return positive / float64(k)
}

func buildCalibrationBins(values []predictedObservation) []domain.CalibrationBin {
	bins := make([]domain.CalibrationBin, 5)
	probabilityTotals := make([]float64, 5)
	labelTotals := make([]float64, 5)
	for index := range bins {
		bins[index].LowerBound = float64(index) / 5
		bins[index].UpperBound = float64(index+1) / 5
	}
	for _, value := range values {
		index := int(value.probability * 5)
		if index > 4 {
			index = 4
		}
		bins[index].Count++
		probabilityTotals[index] += value.probability
		labelTotals[index] += value.label
	}
	for index := range bins {
		if bins[index].Count > 0 {
			bins[index].AverageProbability = probabilityTotals[index] / float64(bins[index].Count)
			bins[index].ObservedRate = labelTotals[index] / float64(bins[index].Count)
		}
		bins[index].Sufficient = bins[index].Count >= 5
	}
	return bins
}

// CalibrateScreening executa leave-one-idea-out e só libera probabilidade se superar prevalência.
func CalibrateScreening(groups []domain.ScreeningDatasetGroup, base domain.ScreeningModel) (domain.CalibrationReport, error) {
	if len(groups) < 2 {
		return domain.CalibrationReport{}, fmt.Errorf("calibração exige ao menos duas ideias")
	}
	report := domain.CalibrationReport{SchemaVersion: "1", Status: "uncalibrated"}
	weightedModel, weightedBaseline, weightedLoss, totalTest := 0.0, 0.0, 0.0, 0
	oof := make([]predictedObservation, 0)
	for heldOut := range groups {
		train := append([]domain.ScreeningDatasetGroup(nil), groups[:heldOut]...)
		train = append(train, groups[heldOut+1:]...)
		foldModel := base
		foldModel.ExpansionPriors = deriveExpansionPriors(train)
		foldModel.CalibrationStatus = "uncalibrated"
		foldModel.Calibration = nil
		trainValues, err := observations(train, foldModel)
		if err != nil {
			return domain.CalibrationReport{}, err
		}
		testValues, err := observations([]domain.ScreeningDatasetGroup{groups[heldOut]}, foldModel)
		if err != nil {
			return domain.CalibrationReport{}, err
		}
		parameters := fitPlatt(trainValues)
		for _, value := range testValues {
			oof = append(oof, predictedObservation{probability: predictCalibration(value, parameters), label: value.label})
		}
		brier, loss := calibrationMetrics(testValues, parameters)
		baseline := prevalenceBrier(testValues, prevalence(trainValues))
		report.Folds = append(report.Folds, domain.CalibrationFold{HeldOutIdea: groups[heldOut].IdeaID, TrainingCount: len(trainValues), TestCount: len(testValues), ModelBrier: brier, PrevalenceBrier: baseline, ModelLogLoss: loss})
		weightedModel += brier * float64(len(testValues))
		weightedBaseline += baseline * float64(len(testValues))
		weightedLoss += loss * float64(len(testValues))
		totalTest += len(testValues)
	}
	if totalTest == 0 {
		return domain.CalibrationReport{}, fmt.Errorf("não há julgamentos conclusivos")
	}
	report.ObservationCount = totalTest
	report.ModelBrier = weightedModel / float64(totalTest)
	report.PrevalenceBrier = weightedBaseline / float64(totalTest)
	report.ModelLogLoss = weightedLoss / float64(totalTest)
	report.PrecisionAtK = predictedPrecisionAtK(oof, 10)
	report.Bins = buildCalibrationBins(oof)
	finalModel := base
	finalModel.ExpansionPriors = deriveExpansionPriors(groups)
	all, err := observations(groups, finalModel)
	if err != nil {
		return domain.CalibrationReport{}, err
	}
	parameters := fitPlatt(all)
	if report.ModelBrier < report.PrevalenceBrier {
		report.Status = "calibrated"
		finalModel.CalibrationStatus = "calibrated"
		finalModel.Calibration = &parameters
	} else {
		finalModel.CalibrationStatus = "uncalibrated"
		finalModel.Calibration = nil
	}
	report.Model = finalModel
	return report, nil
}

func ValidateAuditRevisions(revisions []domain.JudgmentRevision) (domain.AuditReport, error) {
	report := domain.AuditReport{SchemaVersion: "1", Revisions: append([]domain.JudgmentRevision(nil), revisions...)}
	seen := make(map[string]bool)
	for _, revision := range revisions {
		key := strings.ToLower(strings.TrimSpace(revision.Repository))
		if key == "" || seen[key] || strings.TrimSpace(revision.Evidence) == "" || strings.TrimSpace(revision.Reason) == "" || !validRelevanceLabel(revision.PreviousLabel) || !validRelevanceLabel(revision.NewLabel) {
			return domain.AuditReport{}, fmt.Errorf("revisão inválida para %q", revision.Repository)
		}
		seen[key] = true
		report.ReviewedCount++
		if revision.PreviousLabel == domain.RelevanceIrrelevant && revision.NewLabel == domain.RelevanceRelevant {
			report.FalseNegativeCount++
		}
	}
	return report, nil
}

// SelectAuditCandidates inclui todos os relevantes/incertos e controles negativos nos extremos.
func SelectAuditCandidates(screening domain.ScreeningReport, judgments domain.DiscoveryJudgments, lowCount, boundaryCount int) (domain.AuditSelection, error) {
	if lowCount < 0 || boundaryCount < 0 {
		return domain.AuditSelection{}, fmt.Errorf("contagens de auditoria não podem ser negativas")
	}
	labels := judgmentLabels(judgments)
	decisions := make(map[string]domain.ScreeningDecision, len(screening.Decisions))
	for _, decision := range screening.Decisions {
		decisions[strings.ToLower(decision.Repository)] = decision
	}
	selection := domain.AuditSelection{SchemaVersion: "1"}
	selected := make(map[string]bool)
	irrelevant := make([]domain.ScreeningDecision, 0)
	for key, label := range labels {
		decision, ok := decisions[key]
		if !ok {
			continue
		}
		if label == domain.RelevanceRelevant || label == domain.RelevanceUncertain {
			reason := "rótulo relevante conhecido"
			if label == domain.RelevanceUncertain {
				reason = "incerteza conhecida"
			}
			selection.Candidates = append(selection.Candidates, domain.AuditCandidate{Repository: decision.Repository, CurrentLabel: label, RankingScore: decision.RankingScore, SelectionReason: reason})
			selected[key] = true
		} else if label == domain.RelevanceIrrelevant {
			irrelevant = append(irrelevant, decision)
		}
	}
	sort.SliceStable(irrelevant, func(i, j int) bool {
		if irrelevant[i].RankingScore == irrelevant[j].RankingScore {
			return irrelevant[i].Repository < irrelevant[j].Repository
		}
		return irrelevant[i].RankingScore < irrelevant[j].RankingScore
	})
	appendControls := func(values []domain.ScreeningDecision, count int, reason string) {
		for _, decision := range values {
			if count == 0 {
				break
			}
			key := strings.ToLower(decision.Repository)
			if selected[key] {
				continue
			}
			selection.Candidates = append(selection.Candidates, domain.AuditCandidate{Repository: decision.Repository, CurrentLabel: domain.RelevanceIrrelevant, RankingScore: decision.RankingScore, SelectionReason: reason})
			selected[key] = true
			count--
		}
	}
	appendControls(irrelevant, lowCount, "controle entre menores scores")
	reversed := append([]domain.ScreeningDecision(nil), irrelevant...)
	sort.SliceStable(reversed, func(i, j int) bool {
		if reversed[i].RankingScore == reversed[j].RankingScore {
			return reversed[i].Repository < reversed[j].Repository
		}
		return reversed[i].RankingScore > reversed[j].RankingScore
	})
	appendControls(reversed, boundaryCount, "controle irrelevante próximo ao avanço")
	sort.SliceStable(selection.Candidates, func(i, j int) bool {
		return strings.ToLower(selection.Candidates[i].Repository) < strings.ToLower(selection.Candidates[j].Repository)
	})
	return selection, nil
}
