/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/algorithm_insight.md
 * @layer L1
 * @updated 2026-11-08
 */
package domain

type AlgorithmIdeaEvidence struct {
	Source ImmutableCandidateRef `json:"source"`
	LocatedInsightReference
}
type AlgorithmDecision struct {
	Current  string `json:"current"`
	Proposed string `json:"proposed"`
	Benefit  string `json:"benefit"`
	Cost     string `json:"cost"`
}
type AlgorithmIdea struct {
	Title                 string                  `json:"title"`
	Principle             string                  `json:"principle"`
	Steps                 []string                `json:"steps"`
	Evidence              []AlgorithmIdeaEvidence `json:"evidence"`
	ApplicationConditions []string                `json:"application_conditions"`
	FailureConditions     []string                `json:"failure_conditions"`
	Decision              AlgorithmDecision       `json:"decision"`
	TargetRelation        string                  `json:"target_relation"`
	TargetEvidence        []AlgorithmIdeaEvidence `json:"target_evidence"`
	ValidationObservable  string                  `json:"validation_observable"`
	RejectIf              string                  `json:"reject_if"`
}
type AlgorithmInsightRequest struct {
	SchemaVersion string           `json:"schema_version"`
	Question      string           `json:"question"`
	Target        EvidenceBundle   `json:"target"`
	Sources       []EvidenceBundle `json:"sources"`
	Ideas         []AlgorithmIdea  `json:"ideas"`
}
type LocatedAlgorithmIdea struct {
	ID string `json:"id"`
	AlgorithmIdea
}
type AlgorithmInsightReport struct {
	SchemaVersion       string                 `json:"schema_version"`
	ID                  string                 `json:"id"`
	InputSHA256         string                 `json:"input_sha256"`
	Question            string                 `json:"question"`
	Target              ImmutableCandidateRef  `json:"target"`
	Scope               string                 `json:"scope"`
	InterpretationState string                 `json:"interpretation_state"`
	TargetRelationState string                 `json:"target_relation_state"`
	EffectState         string                 `json:"effect_state"`
	Ideas               []LocatedAlgorithmIdea `json:"ideas"`
	Limitations         []string               `json:"limitations"`
}
