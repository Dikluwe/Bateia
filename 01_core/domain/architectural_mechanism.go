/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/architectural_mechanism.md
 * @layer L1
 * @updated 2026-09-02
 */

package domain

import "context"

type MechanismEvidenceStatus string

const (
	MechanismProven       MechanismEvidenceStatus = "PROVEN"
	MechanismPartial      MechanismEvidenceStatus = "PARTIAL"
	MechanismInconclusive MechanismEvidenceStatus = "INCONCLUSIVE"
)

type MechanismNovelty string

const (
	NoveltyConventional MechanismNovelty = "CONVENTIONAL"
	NoveltyDivergent    MechanismNovelty = "DIVERGENT"
	NoveltyUnknown      MechanismNovelty = "UNKNOWN"
)

type MechanismWitnessKind string

const (
	WitnessImplementation MechanismWitnessKind = "IMPLEMENTATION"
	WitnessTest           MechanismWitnessKind = "TEST"
	WitnessExample        MechanismWitnessKind = "EXAMPLE"
	WitnessManifest       MechanismWitnessKind = "MANIFEST"
	WitnessSchema         MechanismWitnessKind = "SCHEMA"
	WitnessConfiguration  MechanismWitnessKind = "CONFIGURATION"
	WitnessDocumentation  MechanismWitnessKind = "DOCUMENTATION"
	WitnessCounter        MechanismWitnessKind = "COUNTEREVIDENCE"
)

type ImmutableCandidateRef struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
}

type MechanismWitness struct {
	Repository       string               `json:"repository"`
	Commit           string               `json:"commit"`
	Path             string               `json:"path"`
	Symbol           string               `json:"symbol"`
	Line             int                  `json:"line"`
	SourceKind       EvidenceKind         `json:"source_kind"`
	Kind             MechanismWitnessKind `json:"kind"`
	AssessmentStatus EvidenceStatus       `json:"assessment_status"`
	RuleID           string               `json:"rule_id"`
	ContentHash      string               `json:"content_hash"`
}

type MechanismComponentRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ArchitecturalMechanism struct {
	ID                string                  `json:"id"`
	MechanismClass    string                  `json:"mechanism_class"`
	Name              string                  `json:"name"`
	Purpose           string                  `json:"purpose"`
	Trigger           string                  `json:"trigger"`
	Strategy          string                  `json:"strategy"`
	Components        []MechanismComponentRef `json:"components"`
	DataFlow          []string                `json:"data_flow"`
	StateModel        []string                `json:"state_model"`
	FailureModel      []string                `json:"failure_model"`
	Boundaries        []string                `json:"boundaries"`
	Preconditions     []string                `json:"preconditions"`
	Tradeoffs         []string                `json:"tradeoffs"`
	Source            ImmutableCandidateRef   `json:"source"`
	Evidence          []MechanismWitness      `json:"evidence"`
	CounterEvidence   []MechanismWitness      `json:"counter_evidence"`
	Status            MechanismEvidenceStatus `json:"status"`
	Novelty           MechanismNovelty        `json:"novelty"`
	ExtractionRuleIDs []string                `json:"extraction_rule_ids"`
}

type ConventionalBaseline struct {
	BaselineID       string   `json:"baseline_id"`
	IdeaHash         string   `json:"idea_hash"`
	MechanismClasses []string `json:"mechanism_classes"`
}

type MechanismGap struct {
	Code          string                  `json:"code"`
	Source        ImmutableCandidateRef   `json:"source"`
	RuleID        string                  `json:"rule_id"`
	WitnessHash   string                  `json:"witness_hash"`
	StatusCeiling MechanismEvidenceStatus `json:"status_ceiling"`
}

type MechanismExtractionInput struct {
	Bundle               EvidenceBundle
	Assessments          []CandidateAssessment
	ConventionalBaseline *ConventionalBaseline
	RuleCatalogJSON      []byte
	RuleCatalogSHA256    string
}

type MechanismExtractionResult struct {
	Mechanisms []ArchitecturalMechanism `json:"mechanisms"`
	Gaps       []MechanismGap           `json:"gaps"`
}

type MechanismExtractor interface {
	ExtractMechanisms(context.Context, MechanismExtractionInput) (MechanismExtractionResult, error)
}
