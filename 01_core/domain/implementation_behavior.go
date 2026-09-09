/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/implementation_behavior.md
 * @layer L1
 * @updated 2026-09-08
 */
package domain

// ImplementationBehavior records syntax, never execution or requirement satisfaction.
type ImplementationBehavior struct {
	Source      ImplementationReference `json:"source"`
	EndLine     int                     `json:"end_line"`
	Status      string                  `json:"status"`
	Reason      string                  `json:"reason"`
	Facts       []BehaviorFact          `json:"facts"`
	Limitations []string                `json:"limitations"`
}

// BehaviorFact is an instruction observed in the referenced function's source bytes.
type BehaviorFact struct {
	Kind    string `json:"kind"`
	Line    int    `json:"line"`
	EndLine int    `json:"end_line"`
	Target  string `json:"target"`
	Snippet string `json:"snippet"`
}
