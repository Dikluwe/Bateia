/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/o23_refinement_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

type o23Pair struct {
	Code      string `json:"code"`
	SourceRef string `json:"source_ref"`
}

func TestO23ProtocolPairsAndOneRequirementNeighbors(t *testing.T) {
	var protocol struct {
		BindingFreeEvaluator struct {
			OpaqueArtifactRule struct {
				CanonicalDocument struct {
					Pairs []o23Pair `json:"controlled_code_source_ref_pairs"`
				} `json:"canonical_document"`
			} `json:"opaque_artifact_rule"`
		} `json:"binding_free_evaluator"`
	}
	raw := readPublicO5(t, "00_nucleo", "benchmark", "0028", "executable-mutant-protocol-o22.json")
	if digestBytes(raw) != "7436cb8007ce0b9ef96a31ea22321ba3ba260ab06656ca414798f2715f215daf" {
		t.Fatal("O22 protocol hash changed")
	}
	if _, err := o6ParseUniqueJSON(raw); err != nil || json.Unmarshal(raw, &protocol) != nil {
		t.Fatal("invalid O22 protocol")
	}
	pairs := protocol.BindingFreeEvaluator.OpaqueArtifactRule.CanonicalDocument.Pairs
	if len(pairs) != 7 {
		t.Fatalf("controlled pair count %d", len(pairs))
	}
	for i, pair := range pairs {
		t.Run(pair.Code, func(t *testing.T) {
			if !o23ControlledPair(pair.Code, pair.SourceRef) {
				t.Fatal("protocol pair absent from evaluator table")
			}
			control := o23GapExchange(t, pair.Code, pair.SourceRef)
			if finding := Evaluate(control.ContractInput, control.ContractOutput); finding.Classification != Unknown {
				t.Fatalf("pair %d not Unknown: %+v", i, finding)
			}

			wrongSource := o23GapExchange(t, pair.Code, "source:neighbor")
			if finding := Evaluate(wrongSource.ContractInput, wrongSource.ContractOutput); finding.Classification != Preserved {
				t.Fatalf("source_ref neighbor not Preserved: %+v", finding)
			}
			wrongCode := o23GapExchange(t, pair.Code+"_NEIGHBOR", pair.SourceRef)
			if finding := Evaluate(wrongCode.ContractInput, wrongCode.ContractOutput); finding.Classification != Preserved {
				t.Fatalf("code neighbor not Preserved: %+v", finding)
			}
		})
	}
}

func TestO23ClosedShapeAndViolationPrecedence(t *testing.T) {
	control := o23GapExchange(t, "LICENSE_UNRECOGNIZED", "hypothesis:license")
	var output ContractOutput
	var result o6ArtifactResult
	if decodeStrict(control.ContractOutput, &output) != nil || decodeStrict(output.ReturnValue, &result) != nil {
		t.Fatal("decode")
	}
	result.SHA256 = strings.Repeat("0", 64)
	output.ReturnValue = mustJSON(t, result)
	if finding := Evaluate(control.ContractInput, mustJSON(t, output)); finding.Classification != Violated {
		t.Fatalf("opaque relation masked integrity violation: %+v", finding)
	}

	entry := `{"code":"LICENSE_UNRECOGNIZED","hypothesis_id":"hyp_` + strings.Repeat("a", 64) + `","source_ref":"hypothesis:license"}`
	wrongOrder := o23RawGapExchange(t, entry)
	if finding := Evaluate(wrongOrder.ContractInput, wrongOrder.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("wrong field order not Preserved: %+v", finding)
	}
	entry = `{"hypothesis_id":"hyp_` + strings.Repeat("a", 64) + `","code":"LICENSE_UNRECOGNIZED","source_ref":"hypothesis:license","extra":true}`
	extra := o23RawGapExchange(t, entry)
	if finding := Evaluate(extra.ContractInput, extra.ContractOutput); finding.Classification != Preserved {
		t.Fatalf("additional property not Preserved: %+v", finding)
	}
}

func TestO23PublicArtifactsAndManifest(t *testing.T) {
	paths := [][]string{{"00_nucleo", "benchmark", "0028", "oracle", "coverage-matrix-o23.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-invalidation-o20-to-o23.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "oracle-authorship-o23.json"}, {"00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o23.json"}}
	for _, path := range paths {
		if _, err := o6ParseUniqueJSON(readPublicO5(t, path...)); err != nil {
			t.Fatalf("strict JSON %s: %v", strings.Join(path, "/"), err)
		}
	}
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		Visibility     string `json:"visibility"`
		OracleRevision string `json:"oracle_revision"`
		Files          []struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"files"`
	}
	if json.Unmarshal(readPublicO5(t, "00_nucleo", "benchmark", "0028", "oracle", "protected", "manifest-o23.json"), &manifest) != nil || manifest.SchemaVersion != "tekt.step0028.oracle-protected-manifest.v23" || manifest.Visibility != "PROTECTED_UNTIL_SEAL" || manifest.OracleRevision != "O23" || len(manifest.Files) != 42 {
		t.Fatal("manifest identity")
	}
	seen := map[string]bool{}
	for _, file := range manifest.Files {
		if seen[file.Path] || filepath.IsAbs(file.Path) || strings.Contains(file.Path, "..") || !isLowerHex(file.SHA256, 64) {
			t.Fatalf("unsafe %q", file.Path)
		}
		seen[file.Path] = true
		if digestBytes(readPublicO5(t, strings.Split(file.Path, "/")...)) != file.SHA256 {
			t.Fatalf("manifest mismatch %s", file.Path)
		}
	}
}

func o23GapExchange(t *testing.T, code, sourceRef string) Exchange {
	t.Helper()
	id, errID := o6MarshalJSONString("hyp_" + strings.Repeat("a", 64))
	codeJSON, errCode := o6MarshalJSONString(code)
	sourceJSON, errSource := o6MarshalJSONString(sourceRef)
	if errID != nil || errCode != nil || errSource != nil {
		t.Fatal("marshal gap")
	}
	entry := `{"hypothesis_id":` + string(id) + `,"code":` + string(codeJSON) + `,"source_ref":` + string(sourceJSON) + `}`
	return o23RawGapExchange(t, entry)
}

func o23RawGapExchange(t *testing.T, entry string) Exchange {
	t.Helper()
	contents := []byte(`{"schema_version":"1","artifact":"hypothesis_gaps","entries":[` + entry + `]}`)
	hash := digestBytes(contents)
	request := map[string]any{"run_id": runID('3'), "artifact": "hypothesis_gaps", "cursor": nil}
	returned := map[string]any{"artifact": "hypothesis_gaps", "sha256": hash, "size": len(contents), "canonical_json": base64.StdEncoding.EncodeToString(contents), "next_cursor": nil}
	publication := emptyPublication()
	publication.ArtifactNames = []string{"hypothesis_gaps"}
	publication.ArtifactHashes = map[string]string{"hypothesis_gaps": hash}
	return exchange(t, ArtifactRead, request, returned, nil, calls(0, 0, 0, 0), publication)
}
