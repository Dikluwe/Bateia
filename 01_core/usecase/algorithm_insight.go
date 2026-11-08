/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_insight.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"sort"
	"strings"
)

// FormulateAlgorithmInsights binds declared reasoning to source and target bytes.
// It checks integrity and completeness, not algorithm semantics or measured effects.
func FormulateAlgorithmInsights(ctx context.Context, request domain.AlgorithmInsightRequest) (domain.AlgorithmInsightReport, error) {
	invalid := func() (domain.AlgorithmInsightReport, error) {
		return domain.AlgorithmInsightReport{}, domain.ErrInvalidRunInput
	}
	if ctx == nil {
		return invalid()
	}
	if err := ctx.Err(); err != nil {
		return domain.AlgorithmInsightReport{}, err
	}
	raw, err := json.Marshal(request)
	if err != nil || len(raw) > 4<<20 || request.SchemaVersion != "1" || len(request.Ideas) < 1 || len(request.Ideas) > 8 || len(request.Sources) < 1 {
		return invalid()
	}
	// Own nested slices before canonicalizing or sorting caller-supplied evidence.
	var owned domain.AlgorithmInsightRequest
	if err = json.Unmarshal(raw, &owned); err != nil {
		return invalid()
	}
	request = owned
	sources, err := normalizeInsightExploration(domain.InsightExplorationRequest{SchemaVersion: "1", Question: request.Question, Bundles: request.Sources})
	if err != nil {
		return domain.AlgorithmInsightReport{}, err
	}
	target, err := normalizeInsightExploration(domain.InsightExplorationRequest{SchemaVersion: "1", Question: request.Question, Bundles: []domain.EvidenceBundle{request.Target}})
	if err != nil {
		return domain.AlgorithmInsightReport{}, err
	}
	request.Question = sources.Question
	request.Sources = sources.Bundles
	request.Target = target.Bundles[0]
	for i := range request.Ideas {
		if err := ctx.Err(); err != nil {
			return domain.AlgorithmInsightReport{}, err
		}
		idea := &request.Ideas[i]
		for _, entry := range []struct {
			value   *string
			maximum int
		}{{&idea.Title, 200}, {&idea.Principle, 1500}, {&idea.Decision.Current, 1000}, {&idea.Decision.Proposed, 1000}, {&idea.Decision.Benefit, 1000}, {&idea.Decision.Cost, 1000}, {&idea.ValidationObservable, 1000}, {&idea.RejectIf, 1000}} {
			if !validSynthesisText(*entry.value, entry.maximum) || normalizeMechanismText(*entry.value) == "" {
				return invalid()
			}
			*entry.value = normalizeMechanismText(*entry.value)
		}
		if len(idea.Steps) < 2 || !validSynthesisTextSet(idea.Steps, 12, 600) || len(idea.ApplicationConditions) < 1 || !validSynthesisTextSet(idea.ApplicationConditions, 12, 600) || len(idea.FailureConditions) < 1 || !validSynthesisTextSet(idea.FailureConditions, 12, 600) {
			return invalid()
		}
		for j := range idea.Steps {
			idea.Steps[j] = normalizeMechanismText(idea.Steps[j])
			if idea.Steps[j] == "" {
				return invalid()
			}
		}
		idea.ApplicationConditions = canonicalMechanismSet(idea.ApplicationConditions)
		idea.FailureConditions = canonicalMechanismSet(idea.FailureConditions)
		if len(idea.ApplicationConditions) == 0 || len(idea.FailureConditions) == 0 {
			return invalid()
		}
		switch idea.TargetRelation {
		case "ALREADY_PRESENT", "DIFFERENT_APPROACH", "NOT_ESTABLISHED":
		default:
			return invalid()
		}
		idea.Evidence, err = bindAlgorithmEvidence(ctx, idea.Evidence, request.Sources)
		if err != nil {
			return domain.AlgorithmInsightReport{}, err
		}
		idea.TargetEvidence, err = bindAlgorithmEvidence(ctx, idea.TargetEvidence, []domain.EvidenceBundle{request.Target})
		if err != nil {
			return domain.AlgorithmInsightReport{}, err
		}
	}
	sort.Slice(request.Ideas, func(i, j int) bool {
		return string(canonicalJSON(request.Ideas[i])) < string(canonicalJSON(request.Ideas[j]))
	})
	for i := 1; i < len(request.Ideas); i++ {
		if string(canonicalJSON(request.Ideas[i])) == string(canonicalJSON(request.Ideas[i-1])) {
			return invalid()
		}
	}
	hash := insightHash(request)
	report := domain.AlgorithmInsightReport{SchemaVersion: "insight-algorithms-v1", ID: "algorithms_" + insightHash(struct{ Contract, Input string }{"insight-algorithms-v1", hash}), InputSHA256: hash, Question: request.Question, Target: immutableRef(request.Target.Repository), Scope: "SOURCE_BOUND_DECLARED_REASONING", InterpretationState: "DECLARED_NOT_REASSESSED", TargetRelationState: "DECLARED_NOT_REASSESSED", EffectState: "UNVERIFIED", Ideas: []domain.LocatedAlgorithmIdea{}, Limitations: []string{"Identidade, hashes, faixas e literais foram ligados aos bytes fornecidos; a origem externa não foi autenticada pelo formulário.", "Princípios, condições, decisões e relação com o alvo são interpretações declaradas, sem reavaliação semântica automática.", "Recortes do alvo não comprovam ausência global de uma capacidade ou novidade; testes citados não foram executados por este comando.", "Nenhum efeito, ganho de desempenho, adoção ou superioridade foi verificado."}}
	for _, idea := range request.Ideas {
		report.Ideas = append(report.Ideas, domain.LocatedAlgorithmIdea{ID: "algorithm_" + insightHash(struct {
			Input string
			Idea  domain.AlgorithmIdea
		}{hash, idea}), AlgorithmIdea: idea})
	}
	return report, nil
}

func bindAlgorithmEvidence(ctx context.Context, references []domain.AlgorithmIdeaEvidence, bundles []domain.EvidenceBundle) ([]domain.AlgorithmIdeaEvidence, error) {
	if len(references) < 1 || len(references) > 8 {
		return nil, domain.ErrInvalidRunInput
	}
	type evidenceKey struct {
		Source domain.ImmutableCandidateRef
		Path   string
	}
	items := map[evidenceKey]domain.CollectedEvidence{}
	for _, b := range bundles {
		for _, item := range b.Items {
			items[evidenceKey{immutableRef(b.Repository), item.Path}] = item
		}
	}
	for i := range references {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		ref := &references[i]
		item, ok := items[evidenceKey{ref.Source, ref.Path}]
		if !ok {
			return nil, domain.ErrIntegrity
		}
		digest, valid := evidenceDigest(ref.ContentHash)
		if !valid || digest != contentHash(item.Content) {
			return nil, domain.ErrIntegrity
		}
		ref.ContentHash = "sha256:" + digest
		if !validObservationRange(item.Content, ref.InsightObservationReference) || len(ref.Excerpt) > 2048 || strings.TrimSpace(ref.Excerpt) == "" {
			return nil, domain.ErrInvalidRunInput
		}
		lines := observationLines(item.Content)
		span := strings.Join(lines[ref.StartLine-1:ref.EndLine], "\n")
		if !strings.Contains(span, ref.Excerpt) {
			return nil, domain.ErrIntegrity
		}
	}
	sort.Slice(references, func(i, j int) bool {
		return string(canonicalJSON(references[i])) < string(canonicalJSON(references[j]))
	})
	for i := 1; i < len(references); i++ {
		if references[i] == references[i-1] {
			return nil, domain.ErrInvalidRunInput
		}
	}
	return references, nil
}
