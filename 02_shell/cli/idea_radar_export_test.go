/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_radar_export_test.md
 * @layer L2
 * @updated 2026-09-12
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"bateia/01_core/domain"
)

func TestRunIdeaRadarExportCommandStrictUsageAndJSON(t *testing.T) {
	valid := marshalRadarInput0069(t, cliRadarInput0069())
	cases := []struct {
		name string
		args []string
		body []byte
	}{
		{name: "missing export", args: nil, body: valid},
		{name: "missing input flag", args: []string{"export"}, body: valid},
		{name: "missing input value", args: []string{"export", "--input"}, body: valid},
		{name: "unknown flag", args: []string{"export", "--input", "INPUT", "--wat"}, body: valid},
		{name: "positional surplus", args: []string{"export", "--input", "INPUT", "surplus"}, body: valid},
		{name: "duplicate input", args: []string{"export", "--input", "INPUT", "--input", "INPUT"}, body: valid},
		{name: "duplicate output", args: []string{"export", "--input", "INPUT", "--output", "one", "--output", "two"}, body: valid},
		{name: "invalid json", args: []string{"export", "--input", "INPUT"}, body: []byte(`{"schema_version":`)},
		{name: "unknown json field", args: []string{"export", "--input", "INPUT"}, body: append(valid[:len(valid)-1], []byte(`,"unknown":true}`)...)},
		{name: "trailing json", args: []string{"export", "--input", "INPUT"}, body: append(append([]byte{}, valid...), []byte(` {}`)...)},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			inputPath := filepath.Join(directory, "input.json")
			if err := os.WriteFile(inputPath, test.body, 0o600); err != nil {
				t.Fatal(err)
			}
			args := replaceInputPlaceholder0069(test.args, inputPath)
			var stdout bytes.Buffer
			if err := RunIdeaRadarExportCommand(args, &stdout); err == nil {
				t.Fatal("expected strict parse/usage error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("partial stdout on error: %q", stdout.Bytes())
			}
		})
	}
}

func TestRunIdeaRadarExportCommandWritesOnlyStdout(t *testing.T) {
	directory := t.TempDir()
	inputPath := filepath.Join(directory, "input.json")
	inputBytes := marshalRadarInput0069(t, cliRadarInput0069())
	if err := os.WriteFile(inputPath, inputBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	if err := RunIdeaRadarExportCommand([]string{"export", "--input", inputPath}, &stdout); err != nil {
		t.Fatal(err)
	}
	if stdout.Len() == 0 || !bytes.HasSuffix(stdout.Bytes(), []byte("\n")) {
		t.Fatalf("stdout is not the sole complete JSON document: %q", stdout.Bytes())
	}
	assertDataflowEnvelope0069(t, stdout.Bytes())
	gotInput, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotInput, inputBytes) {
		t.Fatal("input changed byte-for-byte")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "input.json" {
		t.Fatalf("stdout mode created filesystem artifacts: %v", entryNames0069(entries))
	}
}

func TestRunIdeaRadarExportCommandAtomicOutputAndSilentStdout(t *testing.T) {
	directory := t.TempDir()
	inputPath := filepath.Join(directory, "input.json")
	outputPath := filepath.Join(directory, "radar.archify.json")
	inputBytes := marshalRadarInput0069(t, cliRadarInput0069())
	if err := os.WriteFile(inputPath, inputBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, []byte("old destination"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	if err := RunIdeaRadarExportCommand([]string{"export", "--input", inputPath, "--output", outputPath}, &stdout); err != nil {
		t.Fatal(err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("output mode polluted stdout: %q", stdout.Bytes())
	}
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	assertDataflowEnvelope0069(t, output)
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("atomic output left temporary artifacts: %v", entryNames0069(entries))
	}
	gotInput, err := os.ReadFile(inputPath)
	if err != nil || !bytes.Equal(gotInput, inputBytes) {
		t.Fatalf("input changed: err=%v got=%q", err, gotInput)
	}
}

func TestRunIdeaRadarExportCommandPreservesDestinationOnEveryPreCommitError(t *testing.T) {
	cases := []struct {
		name string
		body []byte
	}{
		{name: "malformed", body: []byte(`{"schema_version":`)},
		{name: "radar unknown", body: marshalRadarInput0069(t, domain.IdeaRadarInput{SchemaVersion: "opaque"})},
		{name: "empty radar", body: marshalRadarInput0069(t, domain.IdeaRadarInput{SchemaVersion: domain.IdeaRadarInputSchemaVersion})},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			inputPath := filepath.Join(directory, "input.json")
			outputPath := filepath.Join(directory, "output.json")
			originalDestination := []byte("irreplaceable destination")
			if err := os.WriteFile(inputPath, test.body, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(outputPath, originalDestination, 0o600); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			err := RunIdeaRadarExportCommand([]string{"export", "--input", inputPath, "--output", outputPath}, &stdout)
			if err == nil || stdout.Len() != 0 {
				t.Fatalf("want silent failure, got err=%v stdout=%q", err, stdout.Bytes())
			}
			got, readErr := os.ReadFile(outputPath)
			if readErr != nil || !bytes.Equal(got, originalDestination) {
				t.Fatalf("destination changed on error: err=%v got=%q", readErr, got)
			}
		})
	}
}

func TestRunIdeaRadarExportCommandRejectsInputOutputAliases(t *testing.T) {
	directory := t.TempDir()
	inputPath := filepath.Join(directory, "input.json")
	inputBytes := marshalRadarInput0069(t, cliRadarInput0069())
	if err := os.WriteFile(inputPath, inputBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	hardlinkPath := filepath.Join(directory, "hardlink.json")
	if err := os.Link(inputPath, hardlinkPath); err != nil {
		t.Fatal(err)
	}
	symlinkPath := filepath.Join(directory, "symlink.json")
	if err := os.Symlink(inputPath, symlinkPath); err != nil {
		t.Fatal(err)
	}
	aliases := []string{inputPath, filepath.Join(directory, ".", "input.json"), hardlinkPath, symlinkPath}
	for _, outputPath := range aliases {
		t.Run(filepath.Base(outputPath), func(t *testing.T) {
			var stdout bytes.Buffer
			if err := RunIdeaRadarExportCommand([]string{"export", "--input", inputPath, "--output", outputPath}, &stdout); err == nil {
				t.Fatalf("accepted output alias %q", outputPath)
			}
			if stdout.Len() != 0 {
				t.Fatalf("partial stdout: %q", stdout.Bytes())
			}
			got, err := os.ReadFile(inputPath)
			if err != nil || !bytes.Equal(got, inputBytes) {
				t.Fatalf("input changed via alias: err=%v got=%q", err, got)
			}
		})
	}
}

func TestRunIdeaRadarExportCommandWriteFailureDoesNotTouchExistingDestination(t *testing.T) {
	directory := t.TempDir()
	inputPath := filepath.Join(directory, "input.json")
	inputBytes := marshalRadarInput0069(t, cliRadarInput0069())
	if err := os.WriteFile(inputPath, inputBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	outputDirectory := filepath.Join(directory, "destination-is-directory")
	if err := os.Mkdir(outputDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(outputDirectory, "marker")
	if err := os.WriteFile(marker, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	err := RunIdeaRadarExportCommand([]string{"export", "--input", inputPath, "--output", outputDirectory}, &stdout)
	if err == nil || stdout.Len() != 0 {
		t.Fatalf("want silent write failure, got err=%v stdout=%q", err, stdout.Bytes())
	}
	gotMarker, markerErr := os.ReadFile(marker)
	if markerErr != nil || string(gotMarker) != "keep" {
		t.Fatalf("existing destination content changed: err=%v marker=%q", markerErr, gotMarker)
	}
	gotInput, inputErr := os.ReadFile(inputPath)
	if inputErr != nil || !bytes.Equal(gotInput, inputBytes) {
		t.Fatalf("input changed on write failure: err=%v got=%q", inputErr, gotInput)
	}
}

func cliRadarInput0069() domain.IdeaRadarInput {
	return domain.IdeaRadarInput{
		SchemaVersion: domain.IdeaRadarInputSchemaVersion,
		Projections: []domain.IdeaRadarProjection{
			{ID: "problem", Kind: domain.IdeaProjectionProblem, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-a"},
			{ID: "behavior", Kind: domain.IdeaProjectionBehavior, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-b"},
			{ID: "architecture", Kind: domain.IdeaProjectionArchitecture, OriginStatus: domain.EvidenceOriginRecorded, OriginID: "origin-c"},
		},
		Observations: []domain.IdeaRadarObservation{
			{ID: "observation-a", ProjectionID: "problem", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
			{ID: "observation-b", ProjectionID: "behavior", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
			{ID: "observation-c", ProjectionID: "architecture", CandidateID: "candidate", ConceptID: "concept", Stance: domain.IdeaRadarStanceSupport},
		},
	}
}

func marshalRadarInput0069(t *testing.T, input domain.IdeaRadarInput) []byte {
	t.Helper()
	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func replaceInputPlaceholder0069(args []string, inputPath string) []string {
	result := append([]string(nil), args...)
	for index := range result {
		if result[index] == "INPUT" {
			result[index] = inputPath
		}
	}
	return result
}

func assertDataflowEnvelope0069(t *testing.T, encoded []byte) {
	t.Helper()
	if !bytes.HasSuffix(encoded, []byte("\n")) {
		t.Fatal("document lacks final newline")
	}
	var raw map[string]any
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		t.Fatalf("invalid output JSON: %v", err)
	}
	if raw["schema_version"] != float64(1) || raw["diagram_type"] != "dataflow" {
		t.Fatalf("not Dataflow v1: %#v", raw)
	}
	meta, _ := raw["meta"].(map[string]any)
	if meta["title"] != "Bateia Idea Radar" {
		t.Fatalf("wrong title: %#v", meta)
	}
}

func entryNames0069(entries []os.DirEntry) []string {
	names := make([]string, len(entries))
	for index, entry := range entries {
		names[index] = entry.Name()
	}
	sort.Strings(names)
	return names
}
