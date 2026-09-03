/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/coverage_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type sanitizedCoverageFeedback struct {
	SchemaVersion        string `json:"schema_version"`
	Visibility           string `json:"visibility"`
	Step                 string `json:"step"`
	SourceOracleRevision string `json:"source_oracle_revision"`
	TargetOracleRevision string `json:"target_oracle_revision"`
	ContractRevision     string `json:"contract_revision"`
	PrivacyContract      struct {
		EntryFields  []string `json:"entry_fields"`
		Excluded     []string `json:"excluded"`
		IDsAreOpaque bool     `json:"ids_are_opaque"`
	} `json:"privacy_contract"`
	Entries []struct {
		ID                        string `json:"id"`
		PublicContractBinding     string `json:"public_contract_binding"`
		RequiredExecutableSubject string `json:"required_executable_subject"`
		O1CoverageGap             string `json:"o1_coverage_gap"`
	} `json:"entries"`
}

type oracleCoverageMatrix struct {
	SchemaVersion    string `json:"schema_version"`
	Step             string `json:"step"`
	OracleRevision   string `json:"oracle_revision"`
	ContractRevision string `json:"contract_revision"`
	Phase            string `json:"phase"`
	SourceFeedback   struct {
		CoverageGapO1SHA256           string `json:"coverage_gap_o1_sha256"`
		DiscriminatoryFailureO1SHA256 string `json:"discriminatory_failure_o1_sha256"`
	} `json:"source_feedback"`
	CoveragePolicy struct {
		NegativeDomain                    string `json:"negative_domain"`
		PositiveDomain                    string `json:"positive_domain"`
		RangesAreInclusive                bool   `json:"ranges_are_inclusive"`
		EveryRangeRequiresExecutableCases bool   `json:"every_range_requires_executable_case_ids"`
		EveryPositiveRequiresOwnCaseID    bool   `json:"every_positive_requires_its_own_case_id"`
		OpaqueDomainIsSeparate            bool   `json:"opaque_domain_is_separate"`
		OpaqueCasesCountTowardScore       bool   `json:"opaque_cases_count_toward_mutation_score"`
	} `json:"coverage_policy"`
	MutationGroups   []coverageGroup   `json:"mutation_groups"`
	PositiveControls []positiveControl `json:"positive_controls"`
}

type coverageGroup struct {
	Group             string      `json:"group"`
	First             string      `json:"first"`
	Last              string      `json:"last"`
	Operations        []Operation `json:"operations"`
	Capability        string      `json:"capability"`
	ExecutableCaseIDs []string    `json:"executable_case_ids"`
}

type positiveControl struct {
	ID                string      `json:"id"`
	Operations        []Operation `json:"operations"`
	Capability        string      `json:"capability"`
	ExecutableCaseIDs []string    `json:"executable_case_ids"`
}

func TestO2CoverageMatrixExactlyCoversSanitizedFeedback(t *testing.T) {
	feedbackPath := repositoryPath(t, "00_nucleo", "benchmark", "0028", "coverage-gap-o1.json")
	matrixPath := repositoryPath(t, "00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o2.json")

	var feedback sanitizedCoverageFeedback
	decodeStrictFile(t, feedbackPath, &feedback)
	var matrix oracleCoverageMatrix
	decodeStrictFile(t, matrixPath, &matrix)
	frozenCasesPath := repositoryPath(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "fixtures", "conformance-cases.json")
	frozenDocument, err := os.ReadFile(frozenCasesPath)
	if err != nil {
		t.Fatal(err)
	}
	frozenCases, err := DecodeCases(frozenDocument)
	if err != nil {
		t.Fatal(err)
	}
	casesByID := make(map[string]Case, len(frozenCases))
	for _, testCase := range frozenCases {
		casesByID[testCase.ID] = testCase
	}

	if feedback.SchemaVersion != "tekt.step0028.coverage-gap-public.v1" {
		t.Fatalf("unexpected sanitized feedback schema %q", feedback.SchemaVersion)
	}
	if feedback.Visibility != "PUBLIC_ALLOWLISTED_ORACLE_FEEDBACK" || feedback.Step != "0028" || feedback.SourceOracleRevision != "O1" || feedback.TargetOracleRevision != "O2" || feedback.ContractRevision != "R3" || !feedback.PrivacyContract.IDsAreOpaque {
		t.Fatal("unexpected sanitized feedback identity or privacy boundary")
	}
	if matrix.SchemaVersion != "tekt.step0028.oracle-coverage-matrix.v2" || matrix.Step != "0028" || matrix.OracleRevision != "O2" || matrix.ContractRevision != "R3" || matrix.Phase != "PRE_CANDIDATE" {
		t.Fatalf("unexpected matrix identity: %+v", matrix)
	}
	if matrix.SourceFeedback.CoverageGapO1SHA256 != "7eca960c91c865e2a24c697c2568c3a99c11f90a597be16848d9c28f8f7d00e1" || matrix.SourceFeedback.DiscriminatoryFailureO1SHA256 != "d85a18daa25d523f2dca9d5225e2789929a375181032a089fcca302a578a67ec" {
		t.Fatal("matrix is not bound to the authorized O1 feedback")
	}
	if matrix.CoveragePolicy.NegativeDomain != "M001-M160" || matrix.CoveragePolicy.PositiveDomain != "P01-P12" || !matrix.CoveragePolicy.RangesAreInclusive || !matrix.CoveragePolicy.EveryRangeRequiresExecutableCases || !matrix.CoveragePolicy.EveryPositiveRequiresOwnCaseID || !matrix.CoveragePolicy.OpaqueDomainIsSeparate || matrix.CoveragePolicy.OpaqueCasesCountTowardScore {
		t.Fatal("matrix coverage policy does not preserve the O2 gate boundary")
	}

	wantMutations, wantPositives := splitFeedbackIDs(t, feedback)
	gotMutations := expandMutationGroups(t, matrix.MutationGroups, casesByID)
	gotPositives := validatePositiveControls(t, matrix.PositiveControls, casesByID)
	assertSameIDs(t, "mutations", wantMutations, gotMutations)
	assertSameIDs(t, "positive controls", wantPositives, gotPositives)
}

func splitFeedbackIDs(t *testing.T, feedback sanitizedCoverageFeedback) ([]string, []string) {
	t.Helper()
	mutations := make([]string, 0, 160)
	positives := make([]string, 0, 12)
	seen := map[string]bool{}
	for _, entry := range feedback.Entries {
		if entry.ID == "" || seen[entry.ID] || entry.PublicContractBinding == "" || entry.RequiredExecutableSubject == "" || entry.O1CoverageGap == "" {
			t.Fatalf("invalid or duplicate sanitized feedback entry %q", entry.ID)
		}
		seen[entry.ID] = true
		switch {
		case len(entry.ID) == 4 && entry.ID[0] == 'M':
			mutations = append(mutations, entry.ID)
		case len(entry.ID) == 3 && entry.ID[0] == 'P':
			positives = append(positives, entry.ID)
		default:
			t.Fatalf("unexpected sanitized feedback id %q", entry.ID)
		}
	}
	return mutations, positives
}

func expandMutationGroups(t *testing.T, groups []coverageGroup, casesByID map[string]Case) []string {
	t.Helper()
	allowedOperations := executableOperations()
	seen := map[string]bool{}
	ids := make([]string, 0, 160)
	for _, group := range groups {
		if group.Group == "" || group.Capability == "" || len(group.Operations) == 0 || len(group.ExecutableCaseIDs) == 0 {
			t.Fatalf("incomplete mutation coverage group %q", group.Group)
		}
		validateOperations(t, group.Operations, allowedOperations)
		validateCaseIDs(t, group.ExecutableCaseIDs)
		validateCaseBindings(t, group.ExecutableCaseIDs, group.Operations, casesByID)
		first := parseMutationID(t, group.First)
		last := parseMutationID(t, group.Last)
		if first > last {
			t.Fatalf("descending mutation range %s-%s", group.First, group.Last)
		}
		for number := first; number <= last; number++ {
			id := fmt.Sprintf("M%03d", number)
			if seen[id] {
				t.Fatalf("mutation %s covered more than once", id)
			}
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

func validatePositiveControls(t *testing.T, controls []positiveControl, casesByID map[string]Case) []string {
	t.Helper()
	allowedOperations := executableOperations()
	seen := map[string]bool{}
	ids := make([]string, 0, len(controls))
	for _, control := range controls {
		if control.ID == "" || seen[control.ID] || control.Capability == "" || len(control.Operations) == 0 || len(control.ExecutableCaseIDs) == 0 {
			t.Fatalf("invalid positive control %q", control.ID)
		}
		seen[control.ID] = true
		validateOperations(t, control.Operations, allowedOperations)
		validateCaseIDs(t, control.ExecutableCaseIDs)
		validateCaseBindings(t, control.ExecutableCaseIDs, control.Operations, casesByID)
		ownCase := false
		for _, caseID := range control.ExecutableCaseIDs {
			ownCase = ownCase || strings.HasPrefix(caseID, control.ID+"-")
		}
		if !ownCase {
			t.Fatalf("positive control %s is not bound to its declared executable case", control.ID)
		}
		ids = append(ids, control.ID)
	}
	return ids
}

func validateCaseBindings(t *testing.T, caseIDs []string, operations []Operation, casesByID map[string]Case) {
	t.Helper()
	allowed := map[Operation]bool{}
	for _, operation := range operations {
		allowed[operation] = true
	}
	for _, caseID := range caseIDs {
		testCase, ok := casesByID[caseID]
		if !ok {
			t.Fatalf("coverage matrix references absent executable case %q", caseID)
		}
		if !allowed[testCase.Operation] {
			t.Fatalf("case %s operation %s is outside declared operations %v", caseID, testCase.Operation, operations)
		}
	}
}

func executableOperations() map[Operation]bool {
	return map[Operation]bool{
		OperationExtract: true, OperationStart: true, OperationLifecycle: true,
		OperationReplay: true, OperationRuntime: true, OperationTools: true,
		OperationRefinement: true, OperationHypothesis: true, OperationSafety: true,
		OperationArtifact: true, OperationPagination: true, OperationSchema: true,
		OperationFull: true,
	}
}

func validateOperations(t *testing.T, operations []Operation, allowed map[Operation]bool) {
	t.Helper()
	seen := map[Operation]bool{}
	for _, operation := range operations {
		if !allowed[operation] || seen[operation] {
			t.Fatalf("invalid or duplicate executable operation %q", operation)
		}
		seen[operation] = true
	}
}

func validateCaseIDs(t *testing.T, caseIDs []string) {
	t.Helper()
	seen := map[string]bool{}
	for _, caseID := range caseIDs {
		if caseID == "" || seen[caseID] {
			t.Fatalf("empty or duplicate executable case id %q", caseID)
		}
		seen[caseID] = true
	}
}

func parseMutationID(t *testing.T, id string) int {
	t.Helper()
	if len(id) != 4 || id[0] != 'M' {
		t.Fatalf("invalid mutation id %q", id)
	}
	number, err := strconv.Atoi(id[1:])
	if err != nil || number < 1 || number > 160 {
		t.Fatalf("mutation id out of domain %q", id)
	}
	return number
}

func assertSameIDs(t *testing.T, label string, want, got []string) {
	t.Helper()
	sort.Strings(want)
	sort.Strings(got)
	if fmt.Sprint(want) != fmt.Sprint(got) {
		t.Fatalf("%s coverage mismatch: want=%v got=%v", label, want, got)
	}
}

func repositoryPath(t *testing.T, elements ...string) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate coverage test source")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
	return filepath.Join(append([]string{root}, elements...)...)
}

func decodeStrictFile(t *testing.T, path string, target any) {
	t.Helper()
	document, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("trailing JSON in %s", path)
	}
}
