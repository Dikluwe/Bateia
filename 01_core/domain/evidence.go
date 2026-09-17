/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/evidence.md
 * @layer L1
 * @updated 2026-09-16
 */

package domain

import "encoding/json"

type EvidenceStatus string

const (
	EvidenceProven       EvidenceStatus = "PROVEN"
	EvidencePartial      EvidenceStatus = "PARTIAL"
	EvidenceAbsent       EvidenceStatus = "ABSENT"
	EvidenceInconclusive EvidenceStatus = "INCONCLUSIVE"
	EvidenceConflict     EvidenceStatus = "CONFLICT"
)

type RepositoryIdentity struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func (r RepositoryIdentity) Key() string { return r.Owner + "/" + r.Repo }

type ResolvedRepository struct {
	Identity RepositoryIdentity `json:"identity"`
	Commit   string             `json:"commit"`
}

type EvidenceKind string

const (
	EvidenceDocumentation      EvidenceKind = "documentation"
	EvidenceManifest           EvidenceKind = "manifest"
	EvidenceFile               EvidenceKind = "file"
	EvidenceSymbol             EvidenceKind = "symbol"
	EvidenceTest               EvidenceKind = "test"
	EvidenceTestImplementation EvidenceKind = "test_implementation"
	EvidenceExample            EvidenceKind = "example"
	EvidenceLicense            EvidenceKind = "license"
	EvidenceRepository         EvidenceKind = "repository_metadata"
)

// EvidenceMatchClass records where a rule pattern matched. Path-only matches
// are useful to explain retrieval, but are never admissible evidence sources.
type EvidenceMatchClass string

const (
	EvidenceMatchStrongContent     EvidenceMatchClass = "STRONG_CONTENT"
	EvidenceMatchSupportingContent EvidenceMatchClass = "SUPPORTING_CONTENT"
	EvidenceMatchDocumentation     EvidenceMatchClass = "DOCUMENTATION_CONTENT"
	EvidenceMatchPathOnly          EvidenceMatchClass = "PATH_ONLY"
)

type EvidenceSource struct {
	Kind           EvidenceKind       `json:"kind"`
	Ref            CommitRef          `json:"ref"`
	Path           string             `json:"path,omitempty"`
	Symbol         string             `json:"symbol,omitempty"`
	Line           int                `json:"line,omitempty"`
	Snippet        string             `json:"snippet,omitempty"`
	RuleID         string             `json:"rule_id"`
	MatchedPattern string             `json:"matched_pattern,omitempty"`
	MatchClass     EvidenceMatchClass `json:"match_class,omitempty"`
	ContentHash    string             `json:"content_hash,omitempty"`
	Strength       float64            `json:"strength"`
}

type CapabilityEvidence struct {
	Capability  string           `json:"capability"`
	Status      EvidenceStatus   `json:"status"`
	Sources     []EvidenceSource `json:"sources"`
	Confidence  float64          `json:"confidence"`
	Explanation string           `json:"explanation"`
	ReasonCode  string           `json:"reason_code,omitempty"`
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
	Documents        int   `json:"documents"`
	Manifests        int   `json:"manifests"`
	CodeFiles        int   `json:"code_files"`
	RequirementFiles int   `json:"requirement_files"`
	Tests            int   `json:"tests_examples"`
	TotalBytes       int64 `json:"total_bytes"`
	TreeDepth        int   `json:"tree_depth"`
}

type RequirementCollectionStatus string

const (
	RequirementTargetedExamined     RequirementCollectionStatus = "TARGETED_EXAMINED"
	RequirementNoPlausiblePath      RequirementCollectionStatus = "NO_PLAUSIBLE_PATH"
	RequirementCollectionIncomplete RequirementCollectionStatus = "COLLECTION_INCOMPLETE"
)

type EvidenceRequirementQuery struct {
	RequirementKind string   `json:"requirement_kind"`
	Requirement     string   `json:"requirement"`
	RuleID          string   `json:"rule_id"`
	Terms           []string `json:"terms"`
}

type RequirementCollectionCoverage struct {
	RequirementKind string                      `json:"requirement_kind"`
	Requirement     string                      `json:"requirement"`
	RuleID          string                      `json:"rule_id"`
	Status          RequirementCollectionStatus `json:"status"`
	SearchComplete  bool                        `json:"search_complete"`
	CandidatePaths  []string                    `json:"candidate_paths"`
	ExaminedPaths   []string                    `json:"examined_paths"`
	SkippedPaths    []string                    `json:"skipped_paths"`
	ReasonCode      string                      `json:"reason_code"`
}

type ExtractionTarget struct {
	Path         string       `json:"path"`
	Kind         EvidenceKind `json:"kind"`
	Capabilities []string     `json:"capabilities,omitempty"`
	Reason       string       `json:"reason"`
}

type CollectionTargetStatus string

const (
	TargetPlanned       CollectionTargetStatus = "PLANNED"
	TargetCollected     CollectionTargetStatus = "COLLECTED"
	TargetNotFound      CollectionTargetStatus = "NOT_FOUND"
	TargetFailed        CollectionTargetStatus = "FAILED"
	TargetSkippedBudget CollectionTargetStatus = "SKIPPED_BUDGET"
	TargetUnsupported   CollectionTargetStatus = "UNSUPPORTED"
)

type CollectionTargetResult struct {
	Target  ExtractionTarget       `json:"target"`
	Status  CollectionTargetStatus `json:"status"`
	Bytes   int64                  `json:"bytes,omitempty"`
	Message string                 `json:"message,omitempty"`
}

type EvidenceExtractionPlan struct {
	SchemaVersion      string                     `json:"schema_version"`
	Repository         CommitRef                  `json:"repository"`
	Candidate          RepositoryCandidate        `json:"candidate"`
	Decision           ScreeningDecision          `json:"screening_decision"`
	Targets            []ExtractionTarget         `json:"targets"`
	Budget             ExtractionBudget           `json:"budget"`
	RequirementQueries []EvidenceRequirementQuery `json:"requirement_queries"`
}

type CollectedEvidence struct {
	Kind        EvidenceKind `json:"kind"`
	Ref         CommitRef    `json:"ref"`
	Path        string       `json:"path"`
	Content     string       `json:"content"`
	ContentHash string       `json:"content_hash"`
}

type EvidenceBundle struct {
	ImplementationSearch *ImplementationSearchReport     `json:"implementation_search,omitempty"`
	SchemaVersion        string                          `json:"schema_version"`
	Repository           CommitRef                       `json:"repository"`
	Plan                 EvidenceExtractionPlan          `json:"plan"`
	Items                []CollectedEvidence             `json:"items"`
	Failures             []ExecutionFailure              `json:"failures,omitempty"`
	Complete             bool                            `json:"complete"`
	BytesRead            int64                           `json:"bytes_read"`
	Targets              []CollectionTargetResult        `json:"target_results"`
	TreeComplete         bool                            `json:"tree_complete"`
	DepthReached         int                             `json:"depth_reached"`
	RequirementCoverage  []RequirementCollectionCoverage `json:"requirement_coverage"`
}

func (plan EvidenceExtractionPlan) MarshalJSON() ([]byte, error) {
	type alias EvidenceExtractionPlan
	copy := alias(plan)
	if copy.Targets == nil {
		copy.Targets = []ExtractionTarget{}
	}
	if copy.RequirementQueries == nil {
		copy.RequirementQueries = []EvidenceRequirementQuery{}
	}
	return json.Marshal(copy)
}

func (bundle EvidenceBundle) MarshalJSON() ([]byte, error) {
	type alias EvidenceBundle
	copy := alias(bundle)
	if copy.Items == nil {
		copy.Items = []CollectedEvidence{}
	}
	if copy.Targets == nil {
		copy.Targets = []CollectionTargetResult{}
	}
	if copy.RequirementCoverage == nil {
		copy.RequirementCoverage = []RequirementCollectionCoverage{}
	} else {
		copy.RequirementCoverage = append([]RequirementCollectionCoverage(nil), copy.RequirementCoverage...)
	}
	for i := range copy.RequirementCoverage {
		if copy.RequirementCoverage[i].CandidatePaths == nil {
			copy.RequirementCoverage[i].CandidatePaths = []string{}
		}
		if copy.RequirementCoverage[i].ExaminedPaths == nil {
			copy.RequirementCoverage[i].ExaminedPaths = []string{}
		}
		if copy.RequirementCoverage[i].SkippedPaths == nil {
			copy.RequirementCoverage[i].SkippedPaths = []string{}
		}
	}
	return json.Marshal(copy)
}

type CapabilityPattern struct {
	Kind     EvidenceKind `json:"kind"`
	Patterns []string     `json:"patterns"`
}

type EvidenceRuleBasis string

const (
	EvidenceRuleDeclared   EvidenceRuleBasis = "DECLARED"
	EvidenceRuleHypothesis EvidenceRuleBasis = "HYPOTHESIS"
)

type CatalogCoverageStatus string

const (
	CatalogCoverageCovered    CatalogCoverageStatus = "COVERED"
	CatalogCoverageGap        CatalogCoverageStatus = "CATALOG_GAP"
	CatalogCoverageHypothesis CatalogCoverageStatus = "HYPOTHESIS"
)

type CatalogCoverageEntry struct {
	RequirementKind string                `json:"requirement_kind"`
	Requirement     string                `json:"requirement"`
	Status          CatalogCoverageStatus `json:"status"`
	RuleID          string                `json:"rule_id,omitempty"`
	RuleBasis       EvidenceRuleBasis     `json:"rule_basis,omitempty"`
	ReasonCode      string                `json:"reason_code"`
}

type EvidenceCatalogCoverage struct {
	SchemaVersion string                 `json:"schema_version"`
	RulesVersion  string                 `json:"rules_version"`
	RulesHash     string                 `json:"rules_hash"`
	Entries       []CatalogCoverageEntry `json:"entries"`
}

func (coverage EvidenceCatalogCoverage) MarshalJSON() ([]byte, error) {
	type alias EvidenceCatalogCoverage
	copy := alias(coverage)
	if copy.Entries == nil {
		copy.Entries = []CatalogCoverageEntry{}
	}
	return json.Marshal(copy)
}

func (coverage *EvidenceCatalogCoverage) UnmarshalJSON(data []byte) error {
	type alias EvidenceCatalogCoverage
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if decoded.Entries == nil {
		decoded.Entries = []CatalogCoverageEntry{}
	}
	*coverage = EvidenceCatalogCoverage(decoded)
	return nil
}

type CapabilityRule struct {
	ID                  string              `json:"id"`
	Basis               EvidenceRuleBasis   `json:"basis,omitempty"`
	Family              string              `json:"family,omitempty"`
	Provenance          string              `json:"provenance,omitempty"`
	TargetHints         []string            `json:"target_hints,omitempty"`
	RequiresIndependent bool                `json:"requires_independent,omitempty"`
	CapabilityTerms     []string            `json:"capability_terms"`
	StrongEvidence      []CapabilityPattern `json:"strong_evidence"`
	SupportingEvidence  []CapabilityPattern `json:"supporting_evidence"`
	NegativeEvidence    []string            `json:"negative_evidence"`
}

type CapabilityRuleCatalog struct {
	SchemaVersion string           `json:"schema_version"`
	Version       string           `json:"version"`
	Rules         []CapabilityRule `json:"rules"`
	Hash          string           `json:"hash,omitempty"`
}

type AssessmentReport struct {
	SchemaVersion   string                  `json:"schema_version"`
	RulesVersion    string                  `json:"rules_version"`
	RulesHash       string                  `json:"rules_hash"`
	CatalogCoverage EvidenceCatalogCoverage `json:"catalog_coverage"`
	Assessments     []CandidateAssessment   `json:"assessments"`
}

type EvidenceJudgment struct {
	Repository     string         `json:"repository"`
	Capability     string         `json:"capability"`
	Expected       EvidenceStatus `json:"expected"`
	ComponentValid *bool          `json:"component_valid,omitempty"`
}

type EvidenceMetrics struct {
	EvidencePrecision  *float64 `json:"evidence_precision"`
	FalseAbsenceRate   *float64 `json:"false_absence_rate"`
	InconclusiveRate   float64  `json:"inconclusive_rate"`
	SourceTraceRate    float64  `json:"source_trace_rate"`
	BudgetCompletion   float64  `json:"budget_completion"`
	ComponentPrecision *float64 `json:"component_precision"`
	ProvenCount        int      `json:"proven_count"`
	SampleSufficient   bool     `json:"sample_sufficient"`
}
