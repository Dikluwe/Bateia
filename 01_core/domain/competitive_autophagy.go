/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/competitive_autophagy.md
 * @layer L1
 * @updated 2026-08-26
 */
package domain

type InternalCapabilityDisposition string

const (
	InternalKeep        InternalCapabilityDisposition = "KEEP"
	InternalConsolidate InternalCapabilityDisposition = "CONSOLIDATE"
	InternalGap         InternalCapabilityDisposition = "GAP"
	InternalUnknown     InternalCapabilityDisposition = "UNKNOWN"
)

type InternalCapabilityFinding struct {
	Capability  string                        `json:"capability"`
	Status      EvidenceStatus                `json:"status"`
	Disposition InternalCapabilityDisposition `json:"disposition"`
	Sources     []EvidenceSource              `json:"sources"`
	RuleIDs     []string                      `json:"rule_ids"`
	Rationale   string                        `json:"rationale"`
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
	Repository   CommitRef        `json:"repository"`
	Capabilities []string         `json:"capabilities"`
	Sources      []EvidenceSource `json:"sources"`
	Decision     ReuseDecision    `json:"decision"`
	Cost         int              `json:"cost"`
	Risks        []string         `json:"risks"`
}

type KatamariProposal struct {
	SchemaVersion string                `json:"schema_version"`
	Components    []KatamariComponent   `json:"components"`
	CoveredGaps   []string              `json:"covered_gaps"`
	UncoveredGaps []string              `json:"uncovered_gaps"`
	Rejected      []DecisionAlternative `json:"rejected"`
	Minimal       bool                  `json:"minimal"`
	Feasible      bool                  `json:"feasible"`
	InputHash     string                `json:"input_hash"`
}

type CompetitiveAutophagyResult struct {
	SchemaVersion string           `json:"schema_version"`
	Autophagy     AutophagyReport  `json:"autophagy"`
	Katamari      KatamariProposal `json:"katamari"`
}
