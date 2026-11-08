/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/insight_discovery_wiring_test.md
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

func TestInsightDiscoveryMainInjectsOnlinePortsBeforeOfflineDispatch(t *testing.T) {
	if os.Getenv("BATEIA_TEST_INSIGHT_DISCOVERY_CHILD") == "1" {
		os.Args = []string{"bateia", "idea", "insights", "discover", "--question", "backpressure"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestInsightDiscoveryMainInjectsOnlinePortsBeforeOfflineDispatch$")
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"BATEIA_TEST_INSIGHT_DISCOVERY_CHILD=1"}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "exige --output antes de acessar a rede") || strings.Contains(stderr.String(), "configurados pelo executável") {
		t.Fatalf("online ports were not injected before validation: %v %s", err, stderr.String())
	}
}
