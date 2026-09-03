/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/improvement_synthesis_step0028_oracle_test.md
 * @layer L4
 * @updated 2026-09-02
 */

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStep0028OracleWiringReusesOneRuntimeAndContainsNoBusinessPolicy(t *testing.T) {
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var production string
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatal(readErr)
		}
		production += string(content)
	}
	if strings.Count(production, "runstore.NewFileStore(") != 1 || strings.Count(production, "runworker.NewLocalQueue(") != 1 {
		t.Fatal("Step 0028 created a parallel durable store or queue")
	}
	for _, forbidden := range []string{"AdoptionRate", "MechanismProven", "HypothesisSupported", "json.Unmarshal", "sha256.Sum", "http.Get", "http.NewRequest", "os.ReadFile("} {
		if strings.Contains(production, forbidden) {
			t.Fatalf("business/parsing/I/O policy leaked into L4: %s", forbidden)
		}
	}
}
