/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_neighborhood_wiring_test.md
 * @layer L4
 * @updated 2026-10-28
 */
package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestInsightNeighborhoodMainDispatchesOnlineBeforeOffline(t *testing.T) {
	if os.Getenv("BATEIA_TEST_NEIGHBORHOOD_CHILD") == "1" {
		os.Args = []string{"bateia", "idea", "insights", "expand"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestInsightNeighborhoodMainDispatchesOnlineBeforeOffline$")
	cmd.Dir, cmd.Env = t.TempDir(), []string{"BATEIA_TEST_NEIGHBORHOOD_CHILD=1"}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "exige --output antes de acessar a rede") {
		t.Fatalf("expand not wired before offline dispatch: %v %s", err, stderr.String())
	}
}
