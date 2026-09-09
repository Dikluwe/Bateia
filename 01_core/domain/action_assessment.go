/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/action_assessment.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain

type RequirementActionAssessment struct {
	Query          string            `json:"query"`
	RequiredAction string            `json:"required_action"`
	Status         string            `json:"status"`
	Reason         string            `json:"reason"`
	Candidates     []ActionCandidate `json:"candidates"`
}
type ActionCandidate struct {
	Source   ImplementationReference `json:"source"`
	EndLine  int                     `json:"end_line"`
	Relation string                  `json:"relation"`
	Reason   string                  `json:"reason"`
	Evidence []ActionEvidence        `json:"evidence"`
}
type ActionEvidence struct {
	Kind    string `json:"kind"`
	Action  string `json:"action"`
	Term    string `json:"term"`
	Line    int    `json:"line"`
	EndLine int    `json:"end_line"`
	Snippet string `json:"snippet"`
}
