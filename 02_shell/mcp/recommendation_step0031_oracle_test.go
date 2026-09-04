/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/recommendation_step0031_oracle_test.md
 * @layer L2
 * @updated 2026-09-03
 */
package mcp

import (
	"bateia/01_core/domain"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var step0031Licenses = []string{"0BSD", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "BSL-1.0", "ISC", "MIT", "Unlicense", "Zlib"}

type step0031SyncSpy struct {
	requests []domain.RecommendationRequest
	mutate   bool
	result   domain.RecommendationResult
}

func (s *step0031SyncSpy) Recommend(_ context.Context, request domain.RecommendationRequest) (domain.RecommendationResult, error) {
	s.requests = append(s.requests, step0031Clone(request))
	if s.mutate {
		step0031CorruptSlices(request)
	}
	return s.result, nil
}

type step0031DurableSpy struct {
	keys     []string
	requests []domain.RecommendationRequest
	mutate   bool
}

func (s *step0031DurableSpy) Start(_ context.Context, key string, request domain.RecommendationRequest) (domain.RecommendationStartReceipt, error) {
	s.keys = append(s.keys, key)
	s.requests = append(s.requests, step0031Clone(request))
	if s.mutate {
		step0031CorruptSlices(request)
	}
	return domain.RecommendationStartReceipt{}, nil
}
func (*step0031DurableSpy) GetRun(context.Context, domain.RunID) (domain.RecommendationRunView, error) {
	return domain.RecommendationRunView{}, nil
}
func (*step0031DurableSpy) GetArtifact(context.Context, domain.ArtifactReadRequest) (domain.ArtifactPage, error) {
	return domain.ArtifactPage{}, nil
}
func (*step0031DurableSpy) Cancel(context.Context, domain.RunID) (domain.RecommendationRunView, error) {
	return domain.RecommendationRunView{}, nil
}

func step0031Strings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string{}, values...)
}

func step0031Clone(request domain.RecommendationRequest) domain.RecommendationRequest {
	clone := request
	clone.Idea.Users = step0031Strings(request.Idea.Users)
	clone.Idea.Capabilities = step0031Strings(request.Idea.Capabilities)
	clone.Idea.Constraints = step0031Strings(request.Idea.Constraints)
	clone.Idea.NonGoals = step0031Strings(request.Idea.NonGoals)
	clone.Idea.PreferredStack = step0031Strings(request.Idea.PreferredStack)
	clone.Policy.AllowedLicenses = step0031Strings(request.Policy.AllowedLicenses)
	clone.Search.Locales = step0031Strings(request.Search.Locales)
	return clone
}

func step0031CorruptSlices(request domain.RecommendationRequest) {
	if len(request.Policy.AllowedLicenses) > 0 {
		request.Policy.AllowedLicenses[0] = "CORRUPTED"
	}
	if len(request.Search.Locales) > 0 {
		request.Search.Locales[0] = "CORRUPTED"
	}
}

func step0031Serve(t *testing.T, sync *step0031SyncSpy, durable *step0031DurableSpy, lines ...string) []JSONRPCResponse {
	t.Helper()
	writer := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(strings.Join(lines, "\n")+"\n"), writer, sync).WithDurableRuns(durable)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var responses []JSONRPCResponse
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			return responses
		} else if err != nil {
			t.Fatalf("invalid JSON-RPC response: %v\n%s", err, writer.Bytes())
		}
		responses = append(responses, response)
	}
}

func step0031Call(id int, tool, arguments string) string {
	return `{"jsonrpc":"2.0","id":` + strconv.Itoa(id) + `,"method":"tools/call","params":{"name":"` + tool + `","arguments":` + arguments + `}}`
}

func step0031Async(arguments string) string {
	if arguments == `{}` {
		return `{"idempotency_key":"step0031-key"}`
	}
	return `{"idempotency_key":"step0031-key",` + strings.TrimPrefix(arguments, "{")
}

func step0031JSONEqual(left, right interface{}) bool {
	leftBytes, _ := json.Marshal(left)
	rightBytes, _ := json.Marshal(right)
	var leftValue, rightValue interface{}
	_ = json.Unmarshal(leftBytes, &leftValue)
	_ = json.Unmarshal(rightBytes, &rightValue)
	return reflect.DeepEqual(leftValue, rightValue)
}

func step0031Required(value interface{}) []string {
	items, _ := value.([]interface{})
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.(string))
	}
	return result
}

func TestStep0031SchemasDeclareOnlyBusinessIntentAndExactDefaults(t *testing.T) {
	responses := step0031Serve(t, &step0031SyncSpy{}, &step0031DurableSpy{},
		`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	)
	if len(responses) != 2 || responses[0].Error != nil || responses[1].Error != nil || !step0031JSONEqual(responses[0].Result, responses[1].Result) {
		t.Fatalf("tools/list failed: %#v", responses)
	}
	encoded, _ := json.Marshal(responses[0].Result)
	var result map[string]interface{}
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatal(err)
	}
	tools := result["tools"].([]interface{})
	byName := map[string]map[string]interface{}{}
	for _, raw := range tools {
		tool := raw.(map[string]interface{})
		byName[tool["name"].(string)] = tool
	}
	recommend := byName["bateia_recommend"]["inputSchema"].(map[string]interface{})
	start := byName["bateia_start_recommendation"]["inputSchema"].(map[string]interface{})
	if got := step0031Required(recommend["required"]); !reflect.DeepEqual(got, []string{"idea"}) {
		t.Fatalf("recommend required = %#v", got)
	}
	if got := step0031Required(start["required"]); !reflect.DeepEqual(got, []string{"idempotency_key", "idea"}) {
		t.Fatalf("start required = %#v", got)
	}
	recommendProperties := recommend["properties"].(map[string]interface{})
	startProperties := start["properties"].(map[string]interface{})
	for _, name := range []string{"idea", "policy", "target", "search"} {
		if !reflect.DeepEqual(recommendProperties[name], startProperties[name]) {
			t.Fatalf("sync/async schema differs at %s", name)
		}
		schema := recommendProperties[name].(map[string]interface{})
		if schema["additionalProperties"] != false {
			t.Fatalf("%s additionalProperties = %#v", name, schema["additionalProperties"])
		}
	}
	idea := recommendProperties["idea"].(map[string]interface{})
	if got := step0031Required(idea["required"]); !reflect.DeepEqual(got, []string{"problem", "capabilities"}) {
		t.Fatalf("idea required = %#v", got)
	}
	for _, name := range []string{"policy", "target", "search"} {
		if got := step0031Required(recommendProperties[name].(map[string]interface{})["required"]); len(got) != 0 {
			t.Fatalf("%s still requires operational members: %#v", name, got)
		}
	}
	policy := recommendProperties["policy"].(map[string]interface{})["properties"].(map[string]interface{})
	search := recommendProperties["search"].(map[string]interface{})["properties"].(map[string]interface{})
	defaults := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"policy.schema_version", policy["schema_version"].(map[string]interface{})["default"], "1"},
		{"policy.allowed_licenses", policy["allowed_licenses"].(map[string]interface{})["default"], step0031Licenses},
		{"policy.require_local", policy["require_local"].(map[string]interface{})["default"], false},
		{"policy.allow_archived", policy["allow_archived"].(map[string]interface{})["default"], false},
		{"policy.max_components", policy["max_components"].(map[string]interface{})["default"], 5},
		{"policy.extend_min_coverage", policy["extend_min_coverage"].(map[string]interface{})["default"], .5},
		{"policy.integration_min_value", policy["integration_min_value"].(map[string]interface{})["default"], .5},
		{"search.strategy", search["strategy"].(map[string]interface{})["default"], "lexical"},
		{"search.language", search["language"].(map[string]interface{})["default"], ""},
		{"search.locales", search["locales"].(map[string]interface{})["default"], []string{"pt-BR", "en-US"}},
		{"search.max_per_query", search["max_per_query"].(map[string]interface{})["default"], 5},
		{"search.max_candidates", search["max_candidates"].(map[string]interface{})["default"], 8},
	}
	for _, item := range defaults {
		if !step0031JSONEqual(item.got, item.want) {
			t.Errorf("%s default = %#v, want %#v", item.name, item.got, item.want)
		}
	}
	targetName := recommendProperties["target"].(map[string]interface{})["properties"].(map[string]interface{})["name"]
	targetJSON, _ := json.Marshal(targetName)
	annotation := strings.ToLower(string(targetJSON))
	if !strings.Contains(annotation, "idea.problem") || (!strings.Contains(annotation, "omit") && !strings.Contains(annotation, "default")) {
		t.Fatalf("target.name does not identify its dynamic default: %s", targetJSON)
	}
}

func step0031DefaultRequest() domain.RecommendationRequest {
	return domain.RecommendationRequest{
		Idea:   domain.IdeaSpec{Problem: "food waste", Capabilities: []string{"inventory"}},
		Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: step0031Strings(step0031Licenses), MaxComponents: 5, ExtendMinCoverage: .5, IntegrationMinValue: .5},
		Target: domain.SoftwareAlternativeTarget{Name: "food waste"},
		Search: domain.RecommendationSearch{Strategy: "lexical", Language: "", Locales: []string{"pt-BR", "en-US"}, MaxPerQuery: 5, MaxCandidates: 8},
	}
}

func TestStep0031MinimalEmptyObjectsAndSyncAsyncMaterializeExactIndependentCopies(t *testing.T) {
	cases := []string{
		`{"idea":{"problem":"food waste","capabilities":["inventory"]}}`,
		`{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{},"target":{},"search":{}}`,
	}
	orders := []struct {
		name    string
		indices []int
		reverse bool
	}{{"forward", []int{0, 1}, false}, {"reverse", []int{1, 0}, true}}
	for _, order := range orders {
		t.Run(order.name, func(t *testing.T) {
			sync := &step0031SyncSpy{mutate: true}
			durable := &step0031DurableSpy{mutate: true}
			var lines []string
			id := 1
			for _, index := range order.indices {
				paths := []struct{ tool, arguments string }{{"bateia_recommend", cases[index]}, {"bateia_start_recommendation", step0031Async(cases[index])}}
				if order.reverse {
					paths[0], paths[1] = paths[1], paths[0]
				}
				for _, path := range paths {
					lines = append(lines, step0031Call(id, path.tool, path.arguments))
					id++
				}
			}
			responses := step0031Serve(t, sync, durable, lines...)
			if len(responses) != 4 || len(sync.requests) != 2 || len(durable.requests) != 2 {
				t.Fatalf("routing mismatch: responses=%d sync=%d async=%d", len(responses), len(sync.requests), len(durable.requests))
			}
			want := step0031DefaultRequest()
			for _, request := range append(sync.requests, durable.requests...) {
				if !reflect.DeepEqual(request, want) {
					t.Errorf("resolved request differs:\ngot  %#v\nwant %#v", request, want)
				}
			}
			if !reflect.DeepEqual(sync.requests, durable.requests) || !reflect.DeepEqual(durable.keys, []string{"step0031-key", "step0031-key"}) {
				t.Fatalf("sync/async parity failed: sync=%#v async=%#v keys=%#v", sync.requests, durable.requests, durable.keys)
			}
			for _, response := range responses {
				if response.Error != nil {
					t.Errorf("positive request failed: %#v", response.Error)
				}
			}
		})
	}
}

func step0031ExplicitCases() []struct {
	arguments string
	want      domain.RecommendationRequest
} {
	wantExplicit := step0031DefaultRequest()
	wantExplicit.Policy.AllowedLicenses = []string{}
	wantExplicit.Policy.AllowArchived = true
	wantExplicit.Target.Name = ""
	wantExplicit.Search.Locales = []string{}
	wantComplete := domain.RecommendationRequest{
		Idea:   domain.IdeaSpec{Problem: "food waste", Users: []string{"families"}, Capabilities: []string{"inventory"}, Constraints: []string{"offline"}, NonGoals: []string{"sales"}, PreferredStack: []string{"Go"}},
		Policy: domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, RequireLocal: true, AllowArchived: true, MaxComponents: 3, ExtendMinCoverage: .7, IntegrationMinValue: .3},
		Target: domain.SoftwareAlternativeTarget{Name: "pantry"},
		Search: domain.RecommendationSearch{Strategy: "literal", Language: "go", Locales: []string{"es-ES"}, MaxPerQuery: 2, MaxCandidates: 4},
	}
	return []struct {
		arguments string
		want      domain.RecommendationRequest
	}{
		{`{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"allowed_licenses":[],"require_local":false,"allow_archived":true},"target":{"name":""},"search":{"language":"","locales":[]}}`, wantExplicit},
		{`{"idea":{"problem":"food waste","users":["families"],"capabilities":["inventory"],"constraints":["offline"],"non_goals":["sales"],"preferred_stack":["Go"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":true,"allow_archived":true,"max_components":3,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"pantry"},"search":{"strategy":"literal","language":"go","locales":["es-ES"],"max_per_query":2,"max_candidates":4}}`, wantComplete},
	}
}

func TestStep0031ExplicitEmptyFalseAndCompleteRequestsArePreservedOnBothPaths(t *testing.T) {
	cases := step0031ExplicitCases()
	for _, order := range []struct {
		name    string
		indices []int
		reverse bool
	}{{"forward", []int{0, 1}, false}, {"reverse", []int{1, 0}, true}} {
		t.Run(order.name, func(t *testing.T) {
			sync := &step0031SyncSpy{}
			durable := &step0031DurableSpy{}
			var lines []string
			id := 1
			for _, index := range order.indices {
				paths := []struct{ tool, arguments string }{{"bateia_recommend", cases[index].arguments}, {"bateia_start_recommendation", step0031Async(cases[index].arguments)}}
				if order.reverse {
					paths[0], paths[1] = paths[1], paths[0]
				}
				for _, path := range paths {
					lines = append(lines, step0031Call(id, path.tool, path.arguments))
					id++
				}
			}
			responses := step0031Serve(t, sync, durable, lines...)
			if len(responses) != 4 || len(sync.requests) != 2 || len(durable.requests) != 2 {
				t.Fatalf("calls failed: responses=%#v sync=%#v async=%#v", responses, sync.requests, durable.requests)
			}
			for position, caseIndex := range order.indices {
				want := cases[caseIndex].want
				if !reflect.DeepEqual(sync.requests[position], want) || !reflect.DeepEqual(durable.requests[position], want) {
					t.Errorf("request[%d] not preserved: sync=%#v async=%#v want=%#v", caseIndex, sync.requests[position], durable.requests[position], want)
				}
				if caseIndex == 0 && (sync.requests[position].Policy.AllowedLicenses == nil || sync.requests[position].Search.Locales == nil) {
					t.Error("explicit empty slices collapsed to nil")
				}
			}
		})
	}
}

func TestStep0031V1AllowlistNeverAssignsLicenseWithoutEvidence(t *testing.T) {
	request := step0031DefaultRequest()
	arguments, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	policy := request.Policy
	result := domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunComplete,
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Policy: policy, Decision: domain.DecisionUnknown, Matrix: []domain.CapabilityComparison{}, Gaps: []string{"inventory"}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Target: request.Target, Filters: domain.SoftwareAlternativeFilters{AllowedLicenses: step0031Strings(step0031Licenses)}, Alternatives: []domain.SoftwareAlternative{{Repository: domain.CommitRef{Owner: "acme", Repo: "opaque", Commit: "abc"}, Fit: domain.AlternativeInconclusive, Capabilities: []domain.SoftwareAlternativeCapability{}, Gaps: []string{"inventory"}, License: domain.LicenseAssessment{Status: domain.EvidenceInconclusive, Caveat: "license evidence missing"}, Reasons: []string{"license evidence remains inconclusive"}, Blockers: []string{}, Evidence: []domain.EvidenceSource{}}}, Excluded: []domain.SoftwareAlternative{}},
		Failures:     []domain.ExecutionFailure{}, Warnings: []string{},
	}
	spy := &step0031SyncSpy{result: result}
	responses := step0031Serve(t, spy, &step0031DurableSpy{}, step0031Call(1, "bateia_recommend", string(arguments)))
	if len(responses) != 1 || responses[0].Error != nil || len(spy.requests) != 1 {
		t.Fatalf("baseline projection failed: responses=%#v requests=%#v", responses, spy.requests)
	}
	envelopeBytes, _ := json.Marshal(responses[0].Result)
	var envelope struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(envelopeBytes, &envelope); err != nil || len(envelope.Content) != 1 {
		t.Fatalf("invalid projection envelope: %v %s", err, envelopeBytes)
	}
	var projected domain.RecommendationResult
	if err := json.Unmarshal([]byte(envelope.Content[0].Text), &projected); err != nil {
		t.Fatal(err)
	}
	if len(projected.Alternatives.Alternatives) != 1 {
		t.Fatalf("opaque candidate missing: %#v", projected.Alternatives)
	}
	license := projected.Alternatives.Alternatives[0].License
	if license.Status != domain.EvidenceInconclusive || license.Identifier != "" || len(license.Sources) != 0 {
		t.Fatalf("V1 acceptance allowlist became candidate evidence: %#v", license)
	}
	if !reflect.DeepEqual(projected.Decision.Policy.AllowedLicenses, step0031Licenses) {
		t.Fatalf("control allowlist missing from same projection: %#v", projected.Decision.Policy)
	}
}

func TestStep0031ExplicitZeroNullInvalidAndUnknownFailClosedWithParity(t *testing.T) {
	cases := []struct{ name, arguments string }{
		{"missing-idea", `{}`},
		{"missing-capability", `{"idea":{"problem":"food waste"}}`},
		{"zero-max-components", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"max_components":0}}`},
		{"zero-extend-threshold", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"extend_min_coverage":0}}`},
		{"zero-integration-threshold", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"integration_min_value":0}}`},
		{"zero-max-per-query", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"search":{"max_per_query":0}}`},
		{"zero-max-candidates", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"search":{"max_candidates":0}}`},
		{"null-policy", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":null}`},
		{"null-license-list", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"allowed_licenses":null}}`},
		{"null-target-name", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"target":{"name":null}}`},
		{"null-language", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"search":{"language":null}}`},
		{"empty-version", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"schema_version":""}}`},
		{"empty-strategy", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"search":{"strategy":""}}`},
		{"invalid-threshold-relation", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"extend_min_coverage":0.4,"integration_min_value":0.6}}`},
		{"unknown-root", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"provider":"hidden"}`},
		{"unknown-nested", `{"idea":{"problem":"food waste","capabilities":["inventory"]},"search":{"mode":"hidden"}}`},
	}
	indices := make([]int, len(cases))
	for index := range cases {
		indices[index] = index
	}
	reversed := append([]int{}, indices...)
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	for _, order := range []struct {
		name       string
		indices    []int
		asyncFirst bool
	}{{"forward", indices, false}, {"reverse", reversed, true}} {
		t.Run(order.name, func(t *testing.T) {
			sync, durable := &step0031SyncSpy{}, &step0031DurableSpy{}
			var lines []string
			id := 1
			for _, index := range order.indices {
				item := cases[index]
				paths := []struct{ tool, arguments string }{{"bateia_recommend", item.arguments}, {"bateia_start_recommendation", step0031Async(item.arguments)}}
				if order.asyncFirst {
					paths[0], paths[1] = paths[1], paths[0]
				}
				for _, path := range paths {
					lines = append(lines, step0031Call(id, path.tool, path.arguments))
					id++
				}
			}
			responses := step0031Serve(t, sync, durable, lines...)
			if len(responses) != len(lines) {
				t.Fatalf("responses=%d want=%d", len(responses), len(lines))
			}
			for _, response := range responses {
				if response.Error == nil || response.Error.Code != -32602 || response.Error.Message != "Invalid params" {
					t.Errorf("invalid request response = %#v", response)
				}
			}
			if len(sync.requests) != 0 || len(durable.requests) != 0 {
				t.Fatalf("invalid request reached a port: sync=%#v async=%#v", sync.requests, durable.requests)
			}
		})
	}
}
