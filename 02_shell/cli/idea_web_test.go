/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_web_test.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestIdeaWebCommands(t *testing.T) {
	dir := t.TempDir()
	queriesFile := filepath.Join(dir, "queries.json")
	snapshotFile := filepath.Join(dir, "snapshot.json")

	queriesContent := `{"schema_version":"1","queries":[{"id":"q1","text":"tool that evaluates code"}]}`
	if err := os.WriteFile(queriesFile, []byte(queriesContent), 0644); err != nil {
		t.Fatalf("write file error: %v", err)
	}

	var stdout bytes.Buffer
	err := runIdeaWeb([]string{"discover-observable", "--queries", queriesFile, "--output", snapshotFile}, &stdout)
	if err != nil {
		t.Fatalf("discover-observable error: %v", err)
	}

	stdout.Reset()
	err = runIdeaWeb([]string{"replay", "--snapshot", snapshotFile}, &stdout)
	if err != nil {
		t.Fatalf("replay error: %v", err)
	}
}
