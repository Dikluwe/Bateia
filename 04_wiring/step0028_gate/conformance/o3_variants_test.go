//go:build step0028_o3_invalidated

/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o3_variants_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"testing"
)

type o3CoverageDiagnostic struct {
	SchemaVersion        string          `json:"schema_version"`
	Visibility           string          `json:"visibility"`
	Step                 string          `json:"step"`
	SourceOracleRevision string          `json:"source_oracle_revision"`
	TargetOracleRevision string          `json:"target_oracle_revision"`
	ContractRevision     string          `json:"contract_revision"`
	GatePolicy           json.RawMessage `json:"gate_policy"`
	PrivacyContract      json.RawMessage `json:"privacy_contract"`
	Entries              []struct {
		ID                              string   `json:"id"`
		PublicContractBinding           string   `json:"public_contract_binding"`
		RequiredObservableFields        []string `json:"required_observable_fields"`
		RequiredMutantSubjectCapability string   `json:"required_mutant_subject_capability"`
		CurrentO2Gap                    string   `json:"current_o2_gap"`
		WhetherExecutable               bool     `json:"whether_executable"`
	} `json:"entries"`
}

type o3PublicCoverageMatrix struct {
	SchemaVersion           string `json:"schema_version"`
	Step                    string `json:"step"`
	OracleRevision          string `json:"oracle_revision"`
	ContractRevision        string `json:"contract_revision"`
	SourceCoverageGapSHA256 string `json:"source_coverage_gap_sha256"`
	Entries                 []struct {
		ID              string   `json:"id"`
		ControlID       string   `json:"control_id"`
		ObservableField string   `json:"observable_field"`
		Capability      string   `json:"capability"`
		RequiredFields  []string `json:"required_fields"`
	} `json:"entries"`
}

func TestStep0028O3EveryNegativeIDHasAnIndependentRejectedVariant(t *testing.T) {
	diagnosticDocument := readO3File(t, repositoryPath(t, "00_nucleo", "benchmark", "0028", "coverage-gap-o2.json"))
	digest := sha256.Sum256(diagnosticDocument)
	if got := hex.EncodeToString(digest[:]); got != "92da449660fdd1f0fba1906cfe9077168820d1e838c671633c33fd0c83c69f3b" {
		t.Fatalf("O2 diagnostic hash changed: %s", got)
	}
	var diagnostic o3CoverageDiagnostic
	decodeO3Strict(t, diagnosticDocument, &diagnostic)
	if diagnostic.SchemaVersion != "tekt.step0028.coverage-gap-public.v2" || diagnostic.Visibility != "PUBLIC_ALLOWLISTED_ORACLE_FEEDBACK" || diagnostic.SourceOracleRevision != "O2" || diagnostic.TargetOracleRevision != "O3" || diagnostic.ContractRevision != "R3" {
		t.Fatal("unexpected O2->O3 diagnostic identity")
	}

	variantDocument := readO3File(t, repositoryPath(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "fixtures", "variants-o3.json"))
	variants, err := DecodeO3Variants(variantDocument)
	if err != nil {
		t.Fatal(err)
	}
	matrixDocument := readO3File(t, repositoryPath(t, "00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o3.json"))
	matrixDigest := sha256.Sum256(matrixDocument)
	if got := hex.EncodeToString(matrixDigest[:]); got != "ef3ff51ca1ca80812eae792ad6b7724bea8af73779ed666096830b2f211a3c00" {
		t.Fatalf("O3 coverage matrix hash changed: %s", got)
	}
	var matrix o3PublicCoverageMatrix
	decodeO3Strict(t, matrixDocument, &matrix)
	if matrix.SchemaVersion != "tekt.step0028.coverage-matrix-public.v3" || matrix.Step != "0028" || matrix.OracleRevision != "O3" || matrix.ContractRevision != "R3" || matrix.SourceCoverageGapSHA256 != "92da449660fdd1f0fba1906cfe9077168820d1e838c671633c33fd0c83c69f3b" || len(matrix.Entries) != 160 {
		t.Fatal("unexpected O3 coverage matrix identity")
	}
	matrixByID := map[string]struct {
		control string
		field   string
	}{}
	for _, entry := range matrix.Entries {
		if entry.Capability == "" || !containsString(entry.RequiredFields, entry.ObservableField) {
			t.Fatalf("matrix entry %s is not bound to its observable", entry.ID)
		}
		matrixByID[entry.ID] = struct{ control, field string }{entry.ControlID, entry.ObservableField}
	}
	requiredByID := map[string]map[string]bool{}
	for _, entry := range diagnostic.Entries {
		if len(entry.ID) != 4 || entry.ID[0] != 'M' {
			continue
		}
		fields := map[string]bool{}
		for _, field := range entry.RequiredObservableFields {
			fields[field] = true
		}
		requiredByID[entry.ID] = fields
	}
	if len(requiredByID) != 160 {
		t.Fatalf("sanitized negative domain = %d, want 160", len(requiredByID))
	}

	controls := map[string]Case{}
	for _, testCase := range O3PositiveCases() {
		controls[testCase.ID] = testCase
	}
	orders := variantOrders(variants)
	for orderIndex, ordered := range orders {
		t.Run(fmt.Sprintf("order-%d", orderIndex+1), func(t *testing.T) {
			seen := map[string]bool{}
			for _, variant := range ordered {
				if seen[variant.ID] || !requiredByID[variant.ID][variant.ObservableField] {
					t.Fatalf("variant %s is duplicate or changes undeclared field %q", variant.ID, variant.ObservableField)
				}
				binding, ok := matrixByID[variant.ID]
				if !ok || binding.control != variant.ControlID || binding.field != variant.ObservableField {
					t.Fatalf("variant %s diverges from public O3 matrix", variant.ID)
				}
				seen[variant.ID] = true
				control, ok := controls[variant.ControlID]
				if !ok {
					t.Fatalf("variant %s names absent control %s", variant.ID, variant.ControlID)
				}
				baseline, err := ReferenceSubject{}.Observe(context.Background(), control)
				if err != nil || !reflect.DeepEqual(baseline, control.Expected) {
					t.Fatalf("baseline %s for %s is not preserved: %v", control.ID, variant.ID, err)
				}
				subject := VariantSubject{Base: ReferenceSubject{}, Variant: variant}
				mutated, err := subject.Observe(context.Background(), control)
				if err != nil {
					t.Fatalf("variant %s did not execute: %v", variant.ID, err)
				}
				differences := differingObservationFields(t, baseline, mutated)
				if !reflect.DeepEqual(differences, []string{variant.ObservableField}) {
					t.Fatalf("variant %s changed %v, want only %s", variant.ID, differences, variant.ObservableField)
				}
				findings := Evaluate(context.Background(), subject, []Case{control})
				if len(findings) != 1 || findings[0].Outcome != "Violated" {
					t.Fatalf("variant %s survived: %+v", variant.ID, findings)
				}
			}
			if len(seen) != 160 {
				t.Fatalf("executed variants = %d, want 160", len(seen))
			}
		})
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func readO3File(t *testing.T, path string) []byte {
	t.Helper()
	document, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return document
}

func decodeO3Strict(t *testing.T, document []byte, target any) {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatal("trailing diagnostic JSON")
	}
}

func differingObservationFields(t *testing.T, baseline, mutated Observation) []string {
	t.Helper()
	encode := func(value Observation) map[string]json.RawMessage {
		document, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		fields := map[string]json.RawMessage{}
		if err := json.Unmarshal(document, &fields); err != nil {
			t.Fatal(err)
		}
		return fields
	}
	want, got := encode(baseline), encode(mutated)
	differences := []string{}
	for field, value := range want {
		if !bytes.Equal(value, got[field]) {
			differences = append(differences, field)
		}
	}
	sort.Strings(differences)
	return differences
}

func variantOrders(variants []VariantSpec) [][]VariantSpec {
	forward := append([]VariantSpec(nil), variants...)
	reverse := append([]VariantSpec(nil), variants...)
	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}
	rotated := append([]VariantSpec(nil), variants...)
	if len(rotated) > 0 {
		pivot := 53 % len(rotated)
		rotated = append(append([]VariantSpec(nil), rotated[pivot:]...), rotated[:pivot]...)
	}
	return [][]VariantSpec{forward, reverse, rotated}
}
