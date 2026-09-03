/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/conformance/o4_conformance_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package conformance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

const o4CommitmentDocumentSHA256 = "0408fa27c802444c3131c138c7e09500b802f2109ba52a5826c3e1bd834d14e5"

func TestStep0028O4CommitmentBindingAndSemanticExecutions(t *testing.T) {
	document := readO4File(t, repositoryPath(t, "00_nucleo", "benchmark", "0028", "mutation-commitments-o4.json"))
	digest := sha256.Sum256(document)
	if got := hex.EncodeToString(digest[:]); got != o4CommitmentDocumentSHA256 {
		t.Fatalf("O4 commitment document hash changed: %s", got)
	}
	commitments, err := DecodeO4Commitments(document)
	if err != nil {
		t.Fatal(err)
	}
	if root, err := O4EntryCommitmentRoot(commitments.Entries); err != nil || root != "a486a67b330946daea5d9bfb41f7764ce35a1e300f94bb77210341d3893a4a3a" {
		t.Fatalf("entry commitment root mismatch: %s (%v)", root, err)
	}

	negative, positive, opaque := partitionO4Commitments(t, commitments.Entries)
	for orderIndex, order := range o4NegativeOrders(negative) {
		t.Run(fmt.Sprintf("negative-order-%d", orderIndex+1), func(t *testing.T) {
			seenInputs := map[string]string{}
			for _, commitment := range order {
				testCase, err := O4NegativeCase(commitment)
				if err != nil {
					t.Fatal(err)
				}
				assertO4ObservableDomain(t, commitment)
				baseline, err := (O4ReferenceSubject{}).Execute(context.Background(), testCase, commitment)
				if err != nil {
					t.Fatalf("%s reference failed: %v", commitment.ID, err)
				}
				variant, err := (O4SemanticVariantSubject{}).Execute(context.Background(), testCase, commitment)
				if err != nil {
					t.Fatalf("%s semantic execution failed: %v", commitment.ID, err)
				}
				finding := EvaluateO4Negative(commitment, testCase.Expected, baseline, variant)
				if finding.Outcome != "Violated" || len(finding.ChangedFields) != 1 || !finding.EntryHashVerified || !finding.InputMutationProven {
					t.Fatalf("%s survived normative predicate: %+v", commitment.ID, finding)
				}
				if prior, duplicate := seenInputs[variant.InputSHA256]; duplicate {
					t.Fatalf("%s reuses semantic input of %s", commitment.ID, prior)
				}
				seenInputs[variant.InputSHA256] = commitment.ID
				encoded, _ := json.Marshal(variant.Observation)
				if strings.Contains(string(encoded), commitment.ID) || strings.Contains(string(encoded), commitment.CanonicalAttackEntrySHA256) || strings.Contains(string(encoded), commitment.SemanticVariantID) {
					t.Fatalf("%s leaked an ID/hash marker into its observation", commitment.ID)
				}
			}
			if len(seenInputs) != 160 {
				t.Fatalf("executed semantic variants = %d, want 160", len(seenInputs))
			}
		})
	}

	for _, commitment := range positive {
		testCase, err := O4PositiveCase(commitment)
		if err != nil {
			t.Fatal(err)
		}
		execution, err := (O4ReferenceSubject{}).Execute(context.Background(), testCase, commitment)
		if err != nil {
			t.Fatalf("%s reference control failed: %v", commitment.ID, err)
		}
		if finding := EvaluateO4Positive(commitment, testCase.Expected, execution); finding.Outcome != "Preserved" || !finding.EntryHashVerified {
			t.Fatalf("%s positive control was not preserved: %+v", commitment.ID, finding)
		}
	}

	for _, commitment := range opaque {
		testCase, err := O4OpaqueCase(commitment)
		if err != nil {
			t.Fatal(err)
		}
		execution, err := (O4OpaqueSubject{}).Execute(context.Background(), testCase, commitment)
		if err != nil {
			t.Fatalf("%s opaque branch failed: %v", commitment.ID, err)
		}
		if finding := EvaluateO4Opaque(commitment, execution); finding.Outcome != "Unknown" || !finding.EntryHashVerified {
			t.Fatalf("%s opacity was coerced: %+v", commitment.ID, finding)
		}
	}
}

func TestStep0028O4SemanticExecutorIsPreObservationAndNonTautological(t *testing.T) {
	path := repositoryPath(t, "04_wiring", "step0028_gate", "conformance", "o4_semantic_variants.go")
	source := string(readO4File(t, path))
	for _, forbidden := range []string{"ReferenceSubject", ".Expected", "alterObservable", "ALTERED_BY_", "switch commitment.ID", "case \"M001\""} {
		if strings.Contains(source, forbidden) {
			t.Fatalf("O4 semantic executor contains forbidden tautology token %q", forbidden)
		}
	}
	for _, required := range []string{"mutated.Input.SemanticFault", "observeO4SemanticInput(mutated)", "INPUT_BEFORE_OBSERVATION"} {
		if !strings.Contains(source, required) {
			t.Fatalf("O4 semantic executor lacks pre-observation evidence %q", required)
		}
	}
}

func TestStep0028O4ProtectedManifestBindsEveryO4HarnessFile(t *testing.T) {
	manifestDocument := readO4File(t, repositoryPath(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o4.json"))
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		Visibility     string `json:"visibility"`
		OracleRevision string `json:"oracle_revision"`
		Files          []struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"files"`
	}
	if err := json.Unmarshal(manifestDocument, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v4" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O4" || len(manifest.Files) != 14 {
		t.Fatalf("unexpected O4 protected manifest identity/count")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || strings.Contains(file.Path, "..") || len(file.SHA256) != 64 {
			t.Fatalf("unsafe or duplicate manifest path %q", file.Path)
		}
		seen[file.Path] = true
		contents := readO4File(t, repositoryPath(t, strings.Split(file.Path, "/")...))
		digest := sha256.Sum256(contents)
		if got := hex.EncodeToString(digest[:]); got != file.SHA256 {
			t.Fatalf("manifest hash mismatch for %s: %s", file.Path, got)
		}
	}
}

func partitionO4Commitments(t *testing.T, entries []O4Commitment) (negative, positive, opaque []O4Commitment) {
	t.Helper()
	for _, entry := range entries {
		switch entry.ID[0] {
		case 'M':
			negative = append(negative, entry)
		case 'P':
			positive = append(positive, entry)
		case 'O':
			opaque = append(opaque, entry)
		default:
			t.Fatalf("unknown commitment class %s", entry.ID)
		}
	}
	if len(negative) != 160 || len(positive) != 12 || len(opaque) != 10 {
		t.Fatalf("partition = %d/%d/%d, want 160/12/10", len(negative), len(positive), len(opaque))
	}
	return negative, positive, opaque
}

func o4NegativeOrders(entries []O4Commitment) [][]O4Commitment {
	forward := append([]O4Commitment(nil), entries...)
	reverse := append([]O4Commitment(nil), entries...)
	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}
	rotated := append([]O4Commitment(nil), entries...)
	pivot := 61 % len(rotated)
	rotated = append(append([]O4Commitment(nil), rotated[pivot:]...), rotated[:pivot]...)
	return [][]O4Commitment{forward, reverse, rotated}
}

func assertO4ObservableDomain(t *testing.T, commitment O4Commitment) {
	t.Helper()
	fields := observationJSONFields(Observation{})
	for _, field := range commitment.RequiredObservableFields {
		if _, represented := fields[field]; !represented {
			t.Fatalf("%s requires unrepresented observation field %s", commitment.ID, field)
		}
	}
	fault, err := o4SemanticFault(commitment.SemanticVariantID)
	if err != nil {
		t.Fatal(err)
	}
	if _, represented := fields[strings.ToLower(fault)]; !represented {
		t.Fatalf("%s semantic fault %s is not an observation field", commitment.ID, fault)
	}
}

func readO4File(t *testing.T, path string) []byte {
	t.Helper()
	document, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return document
}
