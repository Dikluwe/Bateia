/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/competitive_autophagy.md
 * @layer L1
 * @updated 2026-09-07
 */
package domain

type InternalCapabilityDisposition string

const (
	InternalKeep          InternalCapabilityDisposition = "KEEP"
	InternalConsolidate   InternalCapabilityDisposition = "CONSOLIDATE"
	InternalGap           InternalCapabilityDisposition = "GAP"
	InternalUnknown       InternalCapabilityDisposition = "UNKNOWN"
	InternalReviewOverlap InternalCapabilityDisposition = "REVIEW_OVERLAP"
)

type InternalCapabilityFinding struct {
	Capability       string                        `json:"capability"`
	Status           EvidenceStatus                `json:"status"`
	Disposition      InternalCapabilityDisposition `json:"disposition"`
	Sources          []EvidenceSource              `json:"sources"`
	RuleIDs          []string                      `json:"rule_ids"`
	Rationale        string                        `json:"rationale"`
	NextChecks       []string                      `json:"next_checks"`
	CoverageComplete bool                          `json:"coverage_complete"`
}

type AutophagyReport struct {
	SchemaVersion string                      `json:"schema_version"`
	Target        ResolvedRepository          `json:"target"`
	Findings      []InternalCapabilityFinding `json:"findings"`
	Existing      []string                    `json:"existing"`
	Partial       []string                    `json:"partial"`
	Gaps          []string                    `json:"gaps"`
	Unknown       []string                    `json:"unknown"`
	InputHash     string                      `json:"input_hash"`
	Scope         []CommitRef                 `json:"scope"`
	Options       []EvolutionOption           `json:"options"`
}

type KatamariPolicy struct {
	SchemaVersion       string   `json:"schema_version"`
	MaxProjects         int      `json:"max_projects"`
	MaxComponents       int      `json:"max_components"`
	AllowedLicenses     []string `json:"allowed_licenses"`
	RequiredOperation   string   `json:"required_operation,omitempty"`
	RequireImmutableRef bool     `json:"require_immutable_ref"`
	AllowUnknown        bool     `json:"allow_unknown"`
}

type CompetitiveAutophagyRequest struct {
	SchemaVersion  string                `json:"schema_version"`
	Idea           IdeaSpec              `json:"idea"`
	Target         ResolvedRepository    `json:"target"`
	Internal       []EvidenceBundle      `json:"internal"`
	Competitors    []EvidenceBundle      `json:"competitors"`
	Rules          CapabilityRuleCatalog `json:"rules"`
	ReusePolicy    ReuseDecisionPolicy   `json:"reuse_policy"`
	KatamariPolicy KatamariPolicy        `json:"katamari_policy"`
}

type KatamariComponent struct {
	Repository        CommitRef        `json:"repository"`
	Capabilities      []string         `json:"capabilities"`
	Sources           []EvidenceSource `json:"sources"`
	Decision          ReuseDecision    `json:"decision"`
	Cost              int              `json:"cost,omitempty"`
	CostStatus        string           `json:"cost_status"`
	NextChecks        []string         `json:"next_checks"`
	IntegrationStatus string           `json:"integration_status"`
	Risks             []string         `json:"risks"`
}

type KatamariProposal struct {
	SchemaVersion      string                `json:"schema_version"`
	Components         []KatamariComponent   `json:"components"`
	CoveredGaps        []string              `json:"covered_gaps"`
	UncoveredGaps      []string              `json:"uncovered_gaps"`
	Rejected           []DecisionAlternative `json:"rejected"`
	Minimal            bool                  `json:"minimal"`
	Feasible           bool                  `json:"feasible"`
	InputHash          string                `json:"input_hash"`
	Status             string                `json:"status"`
	CoverageComplete   bool                  `json:"coverage_complete"`
	MinimalityScope    string                `json:"minimality_scope"`
	EligibleCandidates int                   `json:"eligible_candidates"`
	Options            []EvolutionOption     `json:"options"`
	Warnings           []string              `json:"warnings"`
}

type CompetitiveAutophagyResult struct {
	SchemaVersion string           `json:"schema_version"`
	Autophagy     AutophagyReport  `json:"autophagy"`
	Katamari      KatamariProposal `json:"katamari"`
}

// EvolutionOption is a bounded proposal, never authorization to change a project.
type EvolutionOption struct {
	Action         string           `json:"action"`
	Capability     string           `json:"capability"`
	Candidate      string           `json:"candidate,omitempty"`
	EvidenceStatus EvidenceStatus   `json:"evidence_status"`
	Sources        []EvidenceSource `json:"sources"`
	Rationale      string           `json:"rationale"`
	NextChecks     []string         `json:"next_checks"`
}
