/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/local_workspace_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCaptureProjectionDoesNotExposePhysicalPaths(t *testing.T) {
	value := CapturedWorkspaceView{ContainerPath: "private-container", Bindings: []SnapshotBinding{{Directory: "private-binding"}}}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "private-") {
		t.Fatal("private capture path exposed")
	}
}
