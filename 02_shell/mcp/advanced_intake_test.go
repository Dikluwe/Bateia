/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/advanced_intake_test.md
 * @layer L2
 * @updated 2026-09-24
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

func serveAdvancedIntake(t *testing.T, advanced bool, lines ...string) []JSONRPCResponse {
	t.Helper()
	output := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(strings.Join(lines, "\n")+"\n"), output)
	if advanced {
		server.WithAdvancedIntake()
	}
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
	responses := []JSONRPCResponse{}
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			return responses
		} else if err != nil {
			t.Fatalf("resposta inválida: %v\n%s", err, output.Bytes())
		}
		responses = append(responses, response)
	}
}

func responseMap(t *testing.T, value interface{}) map[string]interface{} {
	t.Helper()
	encoded, _ := json.Marshal(value)
	var result map[string]interface{}
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestAdvancedIntakeIsOptInAndPreservesLegacyToolsDeeply(t *testing.T) {
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	legacy := responseMap(t, serveAdvancedIntake(t, false, request)[0].Result)["tools"].([]interface{})
	advanced := responseMap(t, serveAdvancedIntake(t, true, request)[0].Result)["tools"].([]interface{})
	if len(legacy) != 3 || len(advanced) != 4 || !reflect.DeepEqual(legacy, advanced[:3]) {
		t.Fatalf("extensão não aditiva: legacy=%#v advanced=%#v", legacy, advanced)
	}
	tool := advanced[3].(map[string]interface{})
	if tool["name"] != "bateia_prepare_recommendation" {
		t.Fatalf("ferramenta = %#v", tool)
	}
}

func TestAdvancedIntakePreparesQuestionsWithoutExecutingRecommendation(t *testing.T) {
	call := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_prepare_recommendation","arguments":{"input":{}}}}`
	responses := serveAdvancedIntake(t, true, call)
	if len(responses) != 1 || responses[0].Error != nil {
		t.Fatalf("responses = %#v", responses)
	}
	result := responseMap(t, responses[0].Result)
	content := result["content"].([]interface{})[0].(map[string]interface{})
	var intake domain.RecommendationIntake
	if err := json.Unmarshal([]byte(content["text"].(string)), &intake); err != nil {
		t.Fatal(err)
	}
	if intake.Status != domain.IntakeNeedsInput || len(intake.Questions) != 2 {
		t.Fatalf("intake = %#v", intake)
	}
}

func TestAdvancedIntakeRejectsNullAndUnknownFields(t *testing.T) {
	for _, arguments := range []string{`{"input":null}`, `{"input":{},"invented":true}`} {
		call := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_prepare_recommendation","arguments":` + arguments + `}}`
		responses := serveAdvancedIntake(t, true, call)
		if len(responses) != 1 || responses[0].Error == nil || responses[0].Error.Code != -32602 {
			t.Fatalf("arguments %s: responses = %#v", arguments, responses)
		}
	}
}
