/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/assisted_judgment.md
 * @layer L1
 * @updated 2026-08-31
 */
package domain

type BlindEvidenceFragment struct {
	FragmentID string `json:"fragment_id"`
	Content    string `json:"content"`
}

type BlindCandidate struct {
	CandidateID string                  `json:"candidate_id"`
	Fragments   []BlindEvidenceFragment `json:"fragments"`
}

type BlindCandidatePackage struct {
	SchemaVersion string           `json:"schema_version"`
	RubricVersion string           `json:"rubric_version"`
	Candidates    []BlindCandidate `json:"candidates"`
	PackageHash   string           `json:"package_hash,omitempty"`
}

type EvidencePointer struct {
	CandidateID string `json:"candidate_id"`
	FragmentID  string `json:"fragment_id"`
}

type AssistedJudgment struct {
	CandidateID   string            `json:"candidate_id"`
	Verdict       string            `json:"verdict"`
	Confidence    float64           `json:"confidence"`
	Reasons       []string          `json:"reasons"`
	Evidence      []EvidencePointer `json:"evidence"`
	RubricVersion string            `json:"rubric_version"`
}

type HumanReview struct {
	CandidateID string `json:"candidate_id"`
	Action      string `json:"action"`
	Verdict     string `json:"verdict,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

type AssistedOracleRecord struct {
	CandidateID      string             `json:"candidate_id"`
	AIHistory        []AssistedJudgment `json:"ai_history"`
	HumanHistory     []HumanReview      `json:"human_history"`
	Effective        string             `json:"effective_verdict"`
	Escalated        bool               `json:"escalated"`
	EscalationCauses []string           `json:"escalation_causes,omitempty"`
}

type AssistedOracleReport struct {
	SchemaVersion string                 `json:"schema_version"`
	Verdict       string                 `json:"verdict"`
	Records       []AssistedOracleRecord `json:"records"`
	PackageHash   string                 `json:"package_hash"`
	Metrics       AssistedOracleMetrics  `json:"metrics"`
}

type AssistedOracleMetrics struct {
	Escalated             int     `json:"escalated"`
	Reviewed              int     `json:"reviewed"`
	ResolvedUncertain     int     `json:"resolved_uncertain"`
	DefinitiveReviewed    int     `json:"definitive_reviewed"`
	DefinitiveAgreements  int     `json:"definitive_agreements"`
	DefinitiveOverrides   int     `json:"definitive_overrides"`
	RawAgreementRate      float64 `json:"raw_agreement_rate"`
	OverrideRate          float64 `json:"override_rate"`
	EvidenceGroundingRate float64 `json:"evidence_grounding_rate"`
}

func (r AssistedOracleReport) Record(candidateID string) AssistedOracleRecord {
	for _, record := range r.Records {
		if record.CandidateID == candidateID {
			return record
		}
	}
	return AssistedOracleRecord{}
}

func (r AssistedOracleReport) EffectiveVerdict(candidateID string) string {
	return r.Record(candidateID).Effective
}
