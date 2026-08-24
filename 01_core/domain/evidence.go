/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence.md
 * @layer L1
 * @updated 2026-08-24
 */

package domain

type EvidenceStatus string

const (
	EvidenceProven       EvidenceStatus = "PROVEN"
	EvidencePartial      EvidenceStatus = "PARTIAL"
	EvidenceAbsent       EvidenceStatus = "ABSENT"
	EvidenceInconclusive EvidenceStatus = "INCONCLUSIVE"
	EvidenceConflict     EvidenceStatus = "CONFLICT"
)

type EvidenceKind string

const (
	EvidenceDocumentation EvidenceKind = "documentation"
	EvidenceManifest      EvidenceKind = "manifest"
	EvidenceFile          EvidenceKind = "file"
	EvidenceSymbol        EvidenceKind = "symbol"
	EvidenceTest          EvidenceKind = "test"
	EvidenceExample       EvidenceKind = "example"
	EvidenceLicense       EvidenceKind = "license"
	EvidenceRepository    EvidenceKind = "repository_metadata"
)

type EvidenceSource struct {
	Kind        EvidenceKind `json:"kind"`
	Ref         CommitRef    `json:"ref"`
	Path        string       `json:"path,omitempty"`
	Symbol      string       `json:"symbol,omitempty"`
	Line        int          `json:"line,omitempty"`
	Snippet     string       `json:"snippet,omitempty"`
	RuleID      string       `json:"rule_id"`
	ContentHash string       `json:"content_hash,omitempty"`
	Strength    float64      `json:"strength"`
}

type CapabilityEvidence struct {
	Capability  string           `json:"capability"`
	Status      EvidenceStatus   `json:"status"`
	Sources     []EvidenceSource `json:"sources"`
	Confidence  float64          `json:"confidence"`
	Explanation string           `json:"explanation"`
}

type ConstraintEvidence struct {
	Constraint string             `json:"constraint"`
	Evidence   CapabilityEvidence `json:"evidence"`
}

type LicenseAssessment struct {
	Status     EvidenceStatus   `json:"status"`
	Identifier string           `json:"identifier,omitempty"`
	Sources    []EvidenceSource `json:"sources,omitempty"`
	Caveat     string           `json:"caveat,omitempty"`
}

type MaintenanceAssessment struct {
	Status   EvidenceStatus   `json:"status"`
	Archived bool             `json:"archived"`
	Sources  []EvidenceSource `json:"sources,omitempty"`
	Warnings []string         `json:"warnings,omitempty"`
}

type OperationAssessment struct {
	LocalStatus  EvidenceStatus   `json:"local_status"`
	Sources      []EvidenceSource `json:"sources,omitempty"`
	Requirements []string         `json:"requirements,omitempty"`
}

type ReusableComponent struct {
	Name         string           `json:"name"`
	Kind         string           `json:"kind"`
	Path         string           `json:"path,omitempty"`
	Symbols      []string         `json:"symbols,omitempty"`
	Capabilities []string         `json:"capabilities"`
	Evidence     []EvidenceSource `json:"evidence"`
	Caveats      []string         `json:"caveats,omitempty"`
}

type CandidateAssessment struct {
	Repository        CommitRef             `json:"repository"`
	ScreeningDecision ScreeningDecision     `json:"screening_decision"`
	ProblemAlignment  CapabilityEvidence    `json:"problem_alignment"`
	Capabilities      []CapabilityEvidence  `json:"capabilities"`
	Constraints       []ConstraintEvidence  `json:"constraints"`
	License           LicenseAssessment     `json:"license"`
	Maintenance       MaintenanceAssessment `json:"maintenance"`
	Operation         OperationAssessment   `json:"operation"`
	ReusableParts     []ReusableComponent   `json:"reusable_parts"`
	Gaps              []string              `json:"gaps"`
	Failures          []ExecutionFailure    `json:"failures,omitempty"`
}

type ExtractionBudget struct {
	Documents  int   `json:"documents"`
	Manifests  int   `json:"manifests"`
	CodeFiles  int   `json:"code_files"`
	Tests      int   `json:"tests_examples"`
	TotalBytes int64 `json:"total_bytes"`
	TreeDepth  int   `json:"tree_depth"`
}

type ExtractionTarget struct {
	Path         string       `json:"path"`
	Kind         EvidenceKind `json:"kind"`
	Capabilities []string     `json:"capabilities,omitempty"`
	Reason       string       `json:"reason"`
}

type EvidenceExtractionPlan struct {
	SchemaVersion string              `json:"schema_version"`
	Repository    CommitRef           `json:"repository"`
	Candidate     RepositoryCandidate `json:"candidate"`
	Decision      ScreeningDecision   `json:"screening_decision"`
	Targets       []ExtractionTarget  `json:"targets"`
	Budget        ExtractionBudget    `json:"budget"`
}

type CollectedEvidence struct {
	Kind        EvidenceKind `json:"kind"`
	Ref         CommitRef    `json:"ref"`
	Path        string       `json:"path"`
	Content     string       `json:"content"`
	ContentHash string       `json:"content_hash"`
}

type EvidenceBundle struct {
	SchemaVersion string                 `json:"schema_version"`
	Repository    CommitRef              `json:"repository"`
	Plan          EvidenceExtractionPlan `json:"plan"`
	Items         []CollectedEvidence    `json:"items"`
	Failures      []ExecutionFailure     `json:"failures,omitempty"`
	Complete      bool                   `json:"complete"`
	BytesRead     int64                  `json:"bytes_read"`
}

type CapabilityPattern struct {
	Kind     EvidenceKind `json:"kind"`
	Patterns []string     `json:"patterns"`
}

type CapabilityRule struct {
	ID                 string              `json:"id"`
	CapabilityTerms    []string            `json:"capability_terms"`
	StrongEvidence     []CapabilityPattern `json:"strong_evidence"`
	SupportingEvidence []CapabilityPattern `json:"supporting_evidence"`
	NegativeEvidence   []string            `json:"negative_evidence"`
}

type CapabilityRuleCatalog struct {
	SchemaVersion string           `json:"schema_version"`
	Version       string           `json:"version"`
	Rules         []CapabilityRule `json:"rules"`
	Hash          string           `json:"hash,omitempty"`
}

type AssessmentReport struct {
	SchemaVersion string                `json:"schema_version"`
	RulesVersion  string                `json:"rules_version"`
	RulesHash     string                `json:"rules_hash"`
	Assessments   []CandidateAssessment `json:"assessments"`
}
