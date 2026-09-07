/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/recommendation_step0025_contract_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package mcp

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type step0025RecommenderSpy struct {
	calls    int
	requests []domain.RecommendationRequest
	result   domain.RecommendationResult
	err      error
}

func (s *step0025RecommenderSpy) Recommend(_ context.Context, request domain.RecommendationRequest) (domain.RecommendationResult, error) {
	s.calls++
	s.requests = append(s.requests, request)
	return s.result, s.err
}

func step0025MCPRequest() string {
	return `{"idea":{"problem":"reduzir desperdício alimentar doméstico","users":["famílias"],"capabilities":["inventário","previsão de consumo"],"constraints":["privacidade"]},"policy":{"schema_version":"1","allowed_licenses":["MIT","Apache-2.0"],"require_local":false,"allow_archived":false,"max_components":3,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":"gestor de alimentos doméstico"},"search":{"strategy":"lexical","language":"","locales":["pt-BR","en-US","ja-JP"],"max_per_query":5,"max_candidates":5}}`
}

func step0025Serve(t *testing.T, recommender RecommendationRecommender, lines ...string) []JSONRPCResponse {
	t.Helper()
	input := strings.Join(lines, "\n") + "\n"
	writer := &bytes.Buffer{}
	server := NewMCPServer(nil, strings.NewReader(input), writer, recommender)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var responses []JSONRPCResponse
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("resposta JSON-RPC inválida: %v\n%s", err, writer.Bytes())
		}
		responses = append(responses, response)
	}
	return responses
}

func decodeStep0025Result(t *testing.T, value interface{}) map[string]interface{} {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]interface{}
	if err := json.Unmarshal(encoded, &object); err != nil {
		t.Fatal(err)
	}
	return object
}

const step0025ExactToolJSON = `{
  "name":"bateia_recommend",
  "description":"Transforma uma ideia em decisão de reúso e alternativas de software com evidências rastreáveis",
  "inputSchema":{
    "type":"object","additionalProperties":false,
    "properties":{
      "idea":{"type":"object","additionalProperties":false,"properties":{"problem":{"type":"string","minLength":1},"users":{"type":"array","items":{"type":"string"}},"capabilities":{"type":"array","minItems":1,"items":{"type":"string","minLength":1}},"constraints":{"type":"array","items":{"type":"string"}},"non_goals":{"type":"array","items":{"type":"string"}},"preferred_stack":{"type":"array","items":{"type":"string"}}},"required":["problem","capabilities"]},
      "policy":{"type":"object","additionalProperties":false,"properties":{"schema_version":{"type":"string","const":"1","default":"1"},"allowed_licenses":{"type":"array","items":{"type":"string","minLength":1},"default":["0BSD","Apache-2.0","BSD-2-Clause","BSD-3-Clause","BSL-1.0","ISC","MIT","Unlicense","Zlib"]},"require_local":{"type":"boolean","default":false},"allow_archived":{"type":"boolean","default":false},"max_components":{"type":"integer","minimum":2,"default":5},"extend_min_coverage":{"type":"number","exclusiveMinimum":0,"maximum":1,"default":0.5},"integration_min_value":{"type":"number","exclusiveMinimum":0,"maximum":1,"default":0.5}},"required":[]},
      "target":{"type":"object","additionalProperties":false,"properties":{"name":{"type":"string","description":"Quando omitido, usa idea.problem como default"}},"required":[]},
      "search":{"type":"object","additionalProperties":false,"properties":{"strategy":{"type":"string","enum":["literal","lexical"],"default":"lexical"},"language":{"type":"string","default":""},"locales":{"type":"array","uniqueItems":true,"items":{"type":"string","enum":["pt-BR","en-US","es-ES","fr-FR","de-DE","ja-JP","ko-KR","ru-RU","ar-SA","hi-IN"]},"default":["pt-BR","en-US"]},"max_per_query":{"type":"integer","minimum":1,"maximum":25,"default":5},"max_candidates":{"type":"integer","minimum":1,"maximum":20,"default":8}},"required":[]}
    },"required":["idea"]
  }
}`

func TestStep0025ToolsListIsExactAndLegacySchemasRemainFrozen(t *testing.T) {
	responses := step0025Serve(t, nil, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if len(responses) != 1 || responses[0].Error != nil {
		t.Fatalf("tools/list falhou: %#v", responses)
	}
	result := decodeStep0025Result(t, responses[0].Result)
	tools, ok := result["tools"].([]interface{})
	if !ok || len(tools) != 3 {
		t.Fatalf("conjunto exato de ferramentas violado: %#v", result["tools"])
	}
	wantNames := []string{"bateia_search", "bateia_benchmark", "bateia_recommend"}
	for index, want := range wantNames {
		tool := tools[index].(map[string]interface{})
		if tool["name"] != want {
			t.Fatalf("tools[%d]=%v, esperado %s", index, tool["name"], want)
		}
	}
	var want interface{}
	if err := json.Unmarshal([]byte(step0025ExactToolJSON), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tools[2], want) {
		got, _ := json.Marshal(tools[2])
		wantJSON, _ := json.Marshal(want)
		t.Fatalf("schema bateia_recommend divergente:\ngot  %s\nwant %s", got, wantJSON)
	}
	legacy, _ := json.Marshal(tools[:2])
	for _, frozen := range []string{
		`"name":"bateia_search"`, `"required":["query"]`, `"default":"go"`,
		`"name":"bateia_benchmark"`, `"required":["repos"]`, `"default":"Benchmark de Repositórios"`,
	} {
		if !strings.Contains(string(legacy), frozen) {
			t.Fatalf("schema legado mudou, testemunha %s ausente: %s", frozen, legacy)
		}
	}
}

func TestStep0025RecommendUsesInjectedServiceOnceAndPreservesItsProjection(t *testing.T) {
	low := domain.CommitRef{Owner: "acme", Repo: "service-choice", Commit: "aaa"}
	high := domain.CommitRef{Owner: "acme", Repo: "popular", Commit: "bbb"}
	fixture := domain.RecommendationResult{
		SchemaVersion: "1", Status: domain.IdeaRunComplete,
		Summary:      domain.RecommendationSummary{ExpansionCount: 4, CandidateCount: 2, AssessedCount: 2},
		Decision:     domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionExtend, Matrix: []domain.CapabilityComparison{{Capability: "inventário", Candidates: []domain.CandidateCapabilityStatus{{Candidate: "acme/service-choice", Status: domain.EvidenceProven, Evidence: []domain.EvidenceSource{{Kind: domain.EvidenceFile, Ref: low, Path: "inventory.go", Symbol: "Inventory", Line: 17, RuleID: "CAP-1", ContentHash: "hash", Strength: .875, Snippet: "REMOVE-ME"}}}}}}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}},
		Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", TopChoice: &low, Alternatives: []domain.SoftwareAlternative{{Repository: low, Stars: 1, Fit: domain.AlternativeRecommended, Capabilities: []domain.SoftwareAlternativeCapability{}, Gaps: []string{}, Reasons: []string{}, Blockers: []string{}, Evidence: []domain.EvidenceSource{{Kind: domain.EvidenceFile, Ref: low, Path: "inventory.go", Symbol: "Inventory", Line: 17, RuleID: "CAP-1", ContentHash: "hash", Strength: .875, Snippet: "REMOVE-ME"}}, DiscoveryLocales: []string{"pt-BR", "en-US", "ja-JP"}}, {Repository: high, Stars: 999999, Fit: domain.AlternativeRecommended, Capabilities: []domain.SoftwareAlternativeCapability{}, Gaps: []string{}, Reasons: []string{}, Blockers: []string{}, Evidence: []domain.EvidenceSource{}}}, Excluded: []domain.SoftwareAlternative{{Repository: domain.CommitRef{Owner: "acme", Repo: "blocked", Commit: "ccc"}, Fit: domain.AlternativeBlocked, Capabilities: []domain.SoftwareAlternativeCapability{}, Gaps: []string{}, Reasons: []string{"incompatible"}, Blockers: []string{"license"}, Evidence: []domain.EvidenceSource{}}}},
		Failures:     []domain.ExecutionFailure{}, Warnings: []string{},
	}
	spy := &step0025RecommenderSpy{result: fixture}
	call := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_recommend","arguments":` + step0025MCPRequest() + `}}`
	responses := step0025Serve(t, spy, call)
	if spy.calls != 1 || len(spy.requests) != 1 {
		t.Fatalf("serviço chamado %d vezes", spy.calls)
	}
	if !reflect.DeepEqual(spy.requests[0].Search.Locales, []string{"pt-BR", "en-US", "ja-JP"}) {
		t.Fatalf("locales alterados antes do serviço: %#v", spy.requests[0].Search.Locales)
	}
	if len(responses) != 1 || responses[0].Error != nil {
		t.Fatalf("call falhou: %#v", responses)
	}
	result := decodeStep0025Result(t, responses[0].Result)
	content := result["content"].([]interface{})
	item := content[0].(map[string]interface{})
	text, ok := item["text"].(string)
	if !ok || item["type"] != "text" || strings.HasSuffix(text, "\n") || strings.HasPrefix(text, "```") {
		t.Fatalf("envelope textual inválido: %#v", item)
	}
	want, err := usecase.MarshalRecommendationResult(fixture)
	if err != nil {
		t.Fatal(err)
	}
	if text != string(want) {
		t.Fatalf("MCP reclassificou/truncou projeção L1:\ngot  %s\nwant %s", text, want)
	}
	if !strings.Contains(text, `"top_choice":{"owner":"acme","repo":"service-choice","commit":"aaa"}`) {
		t.Fatalf("L2 refez ranking por popularidade: %s", text)
	}
	for _, required := range []string{`"excluded":[{"repository":{"owner":"acme","repo":"blocked","commit":"ccc"}`, `"top_choice":{"owner":"acme","repo":"service-choice","commit":"aaa"}`, `"discovery_locales":["pt-BR","en-US","ja-JP"]`} {
		if !strings.Contains(text, required) {
			t.Fatalf("projeção completa/proveniência truncada; falta %s em %s", required, text)
		}
	}
	for _, evidenceField := range []string{`"kind":"file"`, `"ref":{"owner":"acme","repo":"service-choice","commit":"aaa"}`, `"path":"inventory.go"`, `"symbol":"Inventory"`, `"line":17`, `"rule_id":"CAP-1"`, `"content_hash":"hash"`, `"strength":0.875`} {
		if strings.Count(text, evidenceField) < 2 {
			t.Fatalf("campo EvidenceSource não foi preservado em decision e alternatives: %s em %s", evidenceField, text)
		}
	}
	if strings.Contains(text, `"snippet"`) || strings.Contains(text, "REMOVE-ME") {
		t.Fatalf("snippet não foi removido: %s", text)
	}
}

func TestStep0025InvalidArgumentsNeverReachServiceAndCannotAcquireBehavior(t *testing.T) {
	base := step0025MCPRequest()
	cases := map[string]string{
		"missing-problem":      strings.Replace(base, `"problem":"reduzir desperdício alimentar doméstico",`, "", 1),
		"missing-capabilities": strings.Replace(base, `,"capabilities":["inventário","previsão de consumo"]`, "", 1),
		"bad-limit":            strings.Replace(base, `"max_per_query":5`, `"max_per_query":26`, 1),
		"unsupported-locale":   strings.Replace(base, `"ja-JP"`, `"pt-PT"`, 1),
		"unknown-strategy":     strings.Replace(base, `"lexical"`, `"neural"`, 1),
	}
	for _, field := range []string{"output_directory", "replay_directory", "local_path", "local_paths", "path", "directory", "command", "url", "token", "header", "headers", "cookie", "mode", "resume", "run_id", "provider"} {
		cases["prohibited-"+field] = strings.TrimSuffix(base, "}") + `,"` + field + `":"STEP0025_PROHIBITED"}`
	}
	for name, arguments := range cases {
		t.Run(name, func(t *testing.T) {
			spy := &step0025RecommenderSpy{}
			call := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_recommend","arguments":` + arguments + `}}`
			responses := step0025Serve(t, spy, call)
			if spy.calls != 0 {
				t.Fatalf("entrada inválida chegou ao serviço: %#v", spy.requests)
			}
			if len(responses) != 1 || responses[0].Error == nil || responses[0].Error.Code != -32602 || !strings.HasPrefix(responses[0].Error.Message, "Invalid params") {
				t.Fatalf("erro inválido: %#v", responses)
			}
		})
	}
}

func TestStep0025FailureIsSanitizedAndRepeatedCallsAreDeterministic(t *testing.T) {
	canaries := []string{"STEP0025_TOKEN", "STEP0025_AUTHORIZATION", "STEP0025_HEADER", "STEP0025_COOKIE", "STEP0025_ENV", "STEP0025_PROVIDER_CONFIG", "STEP0025_ADAPTER_METADATA"}
	spy := &step0025RecommenderSpy{err: errors.New("provider unavailable Authorization: Bearer " + strings.Join(canaries, " Cookie= header="))}
	call := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_recommend","arguments":` + step0025MCPRequest() + `}}`
	responses := step0025Serve(t, spy, call, call)
	if spy.calls != 2 || len(responses) != 2 {
		t.Fatalf("repetição não executou duas chamadas independentes: calls=%d responses=%d", spy.calls, len(responses))
	}
	first, _ := json.Marshal(responses[0])
	second, _ := json.Marshal(responses[1])
	if string(first) != string(second) {
		t.Fatalf("mesma chamada mudou bytes:\n%s\n%s", first, second)
	}
	for _, response := range responses {
		if response.Error == nil || response.Error.Code != -32603 || strings.Contains(strings.ToLower(response.Error.Message), "authorization") {
			t.Fatalf("segredo ou erro interno exposto: %#v", response.Error)
		}
		for _, canary := range canaries {
			if strings.Contains(response.Error.Message, canary) {
				t.Fatalf("canário %q exposto: %#v", canary, response.Error)
			}
		}
	}
}

func TestStep0025MCPHasNoCLIInfraOrFilesystemExecutionPath(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		source, err := os.ReadFile(entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		file, err := parser.ParseFile(token.NewFileSet(), entry.Name(), source, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imported := range file.Imports {
			path := strings.Trim(imported.Path.Value, `"`)
			if strings.Contains(path, "/02_shell/cli") || strings.Contains(path, "/03_infra/") || path == "os" || path == "io/ioutil" || path == "path/filepath" {
				t.Fatalf("MCP adquiriu dependência proibida %q em %s", path, entry.Name())
			}
		}
	}
	source, err := os.ReadFile("server.go")
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "server.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			name := step0025Selector(call.Fun)
			filesystem := strings.HasPrefix(name, "Create") || name == "OpenFile" || strings.HasPrefix(name, "Temp") || strings.HasPrefix(name, "Mkdir") || name == "ReadFile" || name == "WriteFile"
			if name == "RunIdeaRecommendCommand" || filesystem {
				t.Fatalf("entrypoint MCP chama CLI/filesystem: %s", name)
			}
		}
		return true
	})
}

type step0025LegacySearch struct{}

func (step0025LegacySearch) SearchCandidates(context.Context, domain.QueryFilter) ([]domain.RepositoryCandidate, error) {
	return []domain.RepositoryCandidate{}, nil
}

type step0025LegacyContent struct{}

func (step0025LegacyContent) FetchManifest(_ context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	return &domain.Manifest{Ref: domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID}, ModuleName: "github.com/" + candidate.Owner + "/" + candidate.Repo, Language: "go", GoVersion: "1.22", Dependencies: []domain.Dependency{}}, nil
}
func (step0025LegacyContent) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	return []domain.DirectoryNode{}, nil
}
func (step0025LegacyContent) ExtractRepoFiles(_ context.Context, candidate domain.RepositoryCandidate, _ []string) (*domain.ExtractedRepository, error) {
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: "baseline-commit"}, Candidate: candidate, Files: []domain.ExtractedFile{}, ConsultadoEm: time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)}, nil
}
func (step0025LegacyContent) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	candidate.CommitID = "baseline-commit"
	return candidate, nil
}

type step0025LegacyCache struct{}

func (step0025LegacyCache) GetDiscovery(string, time.Duration) ([]domain.RepositoryCandidate, bool, error) {
	return nil, false, nil
}
func (step0025LegacyCache) SaveDiscovery(string, []domain.RepositoryCandidate) error { return nil }
func (step0025LegacyCache) GetManifest(domain.CommitRef) (*domain.Manifest, bool, error) {
	return nil, false, nil
}
func (step0025LegacyCache) SaveManifest(domain.CommitRef, *domain.Manifest) error { return nil }
func (step0025LegacyCache) GetExtracted(domain.CommitRef) (*domain.ExtractedRepository, bool, error) {
	return nil, false, nil
}
func (step0025LegacyCache) SaveExtracted(domain.CommitRef, *domain.ExtractedRepository) error {
	return nil
}

func step0025LegacyFunnel() *usecase.SearchFunnel {
	return usecase.NewSearchFunnel(step0025LegacySearch{}, step0025LegacyContent{}, step0025LegacyCache{})
}

func step0025ServeLegacy(t *testing.T, lines ...string) []JSONRPCResponse {
	t.Helper()
	writer := &bytes.Buffer{}
	server := NewMCPServer(step0025LegacyFunnel(), strings.NewReader(strings.Join(lines, "\n")+"\n"), writer)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var responses []JSONRPCResponse
	for {
		var response JSONRPCResponse
		if err := decoder.Decode(&response); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		responses = append(responses, response)
	}
	return responses
}

func step0025CanonicalizeLegacyGeneratedAt(t *testing.T, encoded []byte) string {
	t.Helper()
	const marker = `**Gerado em**: `
	const suffix = ` UTC  \n\n`
	const baselineValue = `2026-09-01 17:04:00`
	value := string(encoded)
	if strings.Count(value, marker) != 1 {
		t.Fatalf("resposta legada não contém exatamente um GeneratedAt: %s", value)
	}
	start := strings.Index(value, marker) + len(marker)
	end := strings.Index(value[start:], suffix)
	if end < 0 {
		t.Fatalf("GeneratedAt legado não tem envelope congelado: %s", value)
	}
	return value[:start] + baselineValue + value[start+end:]
}

func TestStep0025LegacySuccessfulSearchUsesEvidenceContract(t *testing.T) {
	searchCall := `{"jsonrpc":"2.0","id":40,"method":"tools/call","params":{"name":"bateia_search","arguments":{"query":"legacy success"}}}`
	searchActual := step0025ServeLegacy(t, searchCall)
	if len(searchActual) != 1 {
		t.Fatalf("resposta search ausente: %#v", searchActual)
	}
	gotSearch := step0025CanonicalizeLegacyGeneratedAt(t, mustStep0025JSON(t, searchActual[0]))
	if strings.Contains(gotSearch, "Soluções Exatas Encontradas") || strings.Contains(gotSearch, "Presença frequente") || !strings.Contains(gotSearch, "Candidatos Encontrados") || !strings.Contains(gotSearch, "Não avaliado") {
		t.Fatalf("contrato de evidência ausente: %s", gotSearch)
	}
}

func mustStep0025JSON(t *testing.T, value interface{}) []byte {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func TestStep0025LegacyValidationErrorsAreByteCompatible(t *testing.T) {
	responses := step0025Serve(t, nil,
		`{"jsonrpc":"2.0","id":20,"method":"tools/call","params":{"name":"bateia_search","arguments":{"query":""}}}`,
		`{"jsonrpc":"2.0","id":21,"method":"tools/call","params":{"name":"bateia_benchmark","arguments":{"repos":"invalid"}}}`,
		`{"jsonrpc":"2.0","id":22,"method":"tools/call","params":{"name":"unknown","arguments":{}}}`,
	)
	want := []struct {
		code    int
		message string
	}{{-32602, "Invalid params: query não pode ser vazia"}, {-32602, "Invalid params: depth deve estar entre 1 e 5"}, {-32601, "Unknown tool: unknown"}}
	if len(responses) != len(want) {
		t.Fatalf("respostas=%d", len(responses))
	}
	for i, expected := range want {
		if responses[i].Error == nil || responses[i].Error.Code != expected.code || responses[i].Error.Message != expected.message {
			t.Fatalf("legado[%d] mudou: %#v", i, responses[i].Error)
		}
	}
}

func TestStep0025SuccessfulRepeatAndFreshServerAreByteIdentical(t *testing.T) {
	fixture := domain.RecommendationResult{SchemaVersion: "1", Status: domain.IdeaRunComplete, Summary: domain.RecommendationSummary{}, Decision: domain.ReuseDecisionReport{SchemaVersion: "1", Decision: domain.DecisionBuild, Matrix: []domain.CapabilityComparison{}, Gaps: []string{}, Reasons: []domain.DecisionReason{}, Alternatives: []domain.DecisionAlternative{}}, Alternatives: domain.SoftwareAlternativeReport{SchemaVersion: "1", Alternatives: []domain.SoftwareAlternative{}, Excluded: []domain.SoftwareAlternative{}}, Failures: []domain.ExecutionFailure{}, Warnings: []string{}}
	call := `{"jsonrpc":"2.0","id":30,"method":"tools/call","params":{"name":"bateia_recommend","arguments":` + step0025MCPRequest() + `}}`
	inProcess := &step0025RecommenderSpy{result: fixture}
	responses := step0025Serve(t, inProcess, call, call)
	fresh := &step0025RecommenderSpy{result: fixture}
	freshResponse := step0025Serve(t, fresh, call)
	if inProcess.calls != 2 || fresh.calls != 1 || len(responses) != 2 || len(freshResponse) != 1 {
		t.Fatalf("traço de chamadas divergente: same=%d fresh=%d", inProcess.calls, fresh.calls)
	}
	bytes0, _ := json.Marshal(responses[0])
	bytes1, _ := json.Marshal(responses[1])
	bytesFresh, _ := json.Marshal(freshResponse[0])
	if string(bytes0) != string(bytes1) || string(bytes0) != string(bytesFresh) || !reflect.DeepEqual(inProcess.requests[0], inProcess.requests[1]) || !reflect.DeepEqual(inProcess.requests[0], fresh.requests[0]) {
		t.Fatalf("sucesso repetido/fresco divergiu:\n%s\n%s\n%s", bytes0, bytes1, bytesFresh)
	}
	for _, forbidden := range []string{`"run_id"`, `"resume"`, `"replay"`} {
		if strings.Contains(string(bytes0), forbidden) {
			t.Fatalf("estado persistente apareceu no sucesso: %s", bytes0)
		}
	}
}

func TestStep0025L4InjectsRecommendationServiceIntoMCP(t *testing.T) {
	path := "../../04_wiring/main.go"
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
	if err != nil {
		t.Fatal(err)
	}
	serviceVars := map[string]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		assign, ok := node.(*ast.AssignStmt)
		if !ok || len(assign.Rhs) != 1 {
			return true
		}
		call, ok := assign.Rhs[0].(*ast.CallExpr)
		if !ok || step0025Selector(call.Fun) != "NewRecommendationService" {
			return true
		}
		for _, lhs := range assign.Lhs {
			if id, ok := lhs.(*ast.Ident); ok {
				serviceVars[id.Name] = true
			}
		}
		return true
	})
	if len(serviceVars) == 0 {
		t.Fatal("L4 não constrói RecommendationService")
	}
	injected := false
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || step0025Selector(call.Fun) != "NewMCPServer" {
			return true
		}
		for _, arg := range call.Args {
			if id, ok := arg.(*ast.Ident); ok && serviceVars[id.Name] {
				injected = true
			}
		}
		return true
	})
	if !injected {
		t.Fatalf("NewMCPServer não recebe o serviço construído por L4: %v", serviceVars)
	}
}

func step0025Selector(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.SelectorExpr:
		return value.Sel.Name
	case *ast.Ident:
		return value.Name
	default:
		return ""
	}
}
