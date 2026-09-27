/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_intake_test.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
)

func TestIdeaIntakePrepareAndApplyOffline(t *testing.T) {
	dir := t.TempDir()
	draftPath := filepath.Join(dir, "draft.json")
	if err := os.WriteFile(draftPath, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	preparedPath := filepath.Join(dir, "intake.json")
	if err := RunIdeaCommand([]string{"intake", "prepare", "--input", draftPath, "--output", preparedPath}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var prepared domain.RecommendationIntake
	if err := readStrictIntakeJSON(preparedPath, &prepared); err != nil {
		t.Fatal(err)
	}
	if prepared.Status != domain.IntakeNeedsInput || len(prepared.Questions) != 2 {
		t.Fatalf("prepared = %#v", prepared)
	}

	answers := []domain.InformationAnswer{
		{NeedID: "initial:idea.problem", AnswerKind: domain.AnswerText, Value: json.RawMessage(`"Pesquisar alternativas"`), Authority: domain.AuthorityUser, Confirmed: true},
		{NeedID: "initial:idea.capabilities", AnswerKind: domain.AnswerStringList, Value: json.RawMessage(`["buscar","comparar"]`), Authority: domain.AuthorityUser, Confirmed: true},
	}
	answersJSON, _ := json.Marshal(answers)
	answersPath := filepath.Join(dir, "answers.json")
	if err := os.WriteFile(answersPath, answersJSON, 0644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := RunIdeaCommand([]string{"intake", "apply", "--input", preparedPath, "--answers", answersPath}, &stdout); err != nil {
		t.Fatal(err)
	}
	var ready domain.RecommendationIntake
	if err := json.Unmarshal(stdout.Bytes(), &ready); err != nil {
		t.Fatal(err)
	}
	if ready.Status != domain.IntakeReady || ready.Request.Idea.Problem != "Pesquisar alternativas" {
		t.Fatalf("ready = %#v", ready)
	}
}

func TestIdeaIntakeRejectsUnknownInputFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(path, []byte(`{"invented":true}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := RunIdeaCommand([]string{"intake", "prepare", "--input", path}, &bytes.Buffer{}); err == nil {
		t.Fatal("campo desconhecido foi aceito")
	}
}
