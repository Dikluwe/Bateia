/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/baseline_tools_step0028_oracle_test.md
 * @layer L2
 * @updated 2026-09-03
 */

package mcp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestStep0028OracleFrozenBasalNineCapture(t *testing.T) {
	connector := &step0027ConnectorStub{}
	starter := &step0027StarterStub{}
	response := step0027RunMCP(t, true, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}\n", connector, starter)[0]
	tools := step0027Tools(t, response)
	encoded, err := json.Marshal(tools)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(encoded)
	var frozen struct {
		SHA256 string   `json:"sha256"`
		Names  []string `json:"names"`
	}
	fixture, err := os.ReadFile(filepath.Join("..", "..", "00_nucleo", "benchmark", "0031", "basal-tools-successor.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(fixture, &frozen); err != nil {
		t.Fatal(err)
	}
	if got := hex.EncodeToString(digest[:]); got != frozen.SHA256 {
		t.Fatalf("nine-tool predecessor capture changed: %s", got)
	}
	names := make([]string, len(tools))
	for index := range tools {
		names[index], _ = tools[index]["name"].(string)
	}
	if !reflect.DeepEqual(names, frozen.Names) {
		t.Fatalf("nine-tool predecessor order changed: %v", names)
	}
}
