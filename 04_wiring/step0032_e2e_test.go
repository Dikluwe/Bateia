/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0032_e2e_test.md
 * @layer L4
 * @updated 2026-09-03
 */
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bateia/01_core/usecase"
	"bateia/02_shell/cli"
	"bateia/02_shell/mcp"
	"bateia/03_infra/cache"
	"bateia/03_infra/content"
	"bateia/03_infra/lexical"
	localrepo "bateia/03_infra/local"
	"bateia/03_infra/runstore"
)

func step0032RPCCall(t *testing.T, writer io.Writer, scanner *bufio.Scanner, request string) mcp.JSONRPCResponse {
	t.Helper()
	if _, err := io.WriteString(writer, request+"\n"); err != nil {
		t.Fatalf("falha ao enviar request MCP: %v", err)
	}
	if !scanner.Scan() {
		t.Fatalf("nenhuma resposta MCP recebida: %v", scanner.Err())
	}
	var response mcp.JSONRPCResponse
	if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
		t.Fatalf("resposta MCP com JSON inválido: %v (raw: %s)", err, scanner.Bytes())
	}
	return response
}

func step0032ExtractTextResult(t *testing.T, response mcp.JSONRPCResponse) map[string]interface{} {
	t.Helper()
	if response.Error != nil {
		t.Fatalf("resposta MCP retornou erro: %#v", response.Error)
	}
	encoded, err := json.Marshal(response.Result)
	if err != nil {
		t.Fatalf("falha ao serializar result: %v", err)
	}
	var envelope struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil || len(envelope.Content) < 1 || envelope.Content[0].Type != "text" {
		t.Fatalf("envelope MCP inválido: %v (result: %#v)", err, response.Result)
	}
	var object map[string]interface{}
	if err := json.Unmarshal([]byte(envelope.Content[0].Text), &object); err != nil {
		t.Fatalf("JSON textual MCP inválido: %v (text: %s)", err, envelope.Content[0].Text)
	}
	return object
}

func TestStep0032E2ECLILocalFunnelWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "relatorio_e2e_teste.md")

	diskStore, err := cache.NewDiskCacheStore(filepath.Join(tempDir, "cache"))
	if err != nil {
		t.Fatalf("falha ao criar DiskCacheStore: %v", err)
	}

	contentProvider := content.NewRouter(nil, localrepo.NewClient())
	funnelUsecase := usecase.NewSearchFunnel(nil, contentProvider, diskStore)
	app := cli.NewCLIApp(funnelUsecase)

	// Executa a CLI apontando para a worktree local atual (offline, sem dependência de rede externa)
	args := []string{"-path", ".", "-query", "funnel", "-output", outputPath}
	if err := app.Run(args); err != nil {
		t.Fatalf("execução E2E da CLI falhou: %v", err)
	}

	contentBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("relatório de saída não foi gerado no caminho esperado: %v", err)
	}

	reportContent := string(contentBytes)
	if len(reportContent) == 0 {
		t.Fatalf("relatório E2E gerado está vazio")
	}

	if !strings.Contains(reportContent, "Bateia") {
		t.Fatalf("relatório E2E não contém cabeçalho do Bateia: %s", reportContent)
	}
}

func TestStep0032E2EMCPSyncAndAsyncWorkflow(t *testing.T) {
	root := t.TempDir()
	secret := []byte("0123456789abcdef0123456789abcdef")
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}

	discovery := &step0025WiringDiscovery{}
	contentMock := &step0025WiringContent{}
	lexicalExpander, err := lexical.NewDefault(2)
	if err != nil {
		t.Fatalf("falha ao instanciar lexicalExpander: %v", err)
	}
	recommendation := usecase.NewRecommendationService(lexicalExpander, discovery, contentMock)

	reqReader, reqWriter := io.Pipe()
	respReader, respWriter := io.Pipe()

	server, shutdown, err := NewDurableMCPServer((*usecase.SearchFunnel)(nil), recommendation, root, secret, limits, 1, 4, reqReader, respWriter)
	if err != nil {
		t.Fatalf("falha ao instanciar NewDurableMCPServer: %v", err)
	}

	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(context.Background()) }()
	scanner := bufio.NewScanner(respReader)

	// 1. Handshake MCP initialize
	initResp := step0032RPCCall(t, reqWriter, scanner, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"e2e-client","version":"1.0"}}}`)
	if initResp.Error != nil {
		t.Fatalf("MCP initialize retornou erro: %#v", initResp.Error)
	}

	// 2. tools/list - valida as ferramentas duráveis do Bateia
	listResp := step0032RPCCall(t, reqWriter, scanner, `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	listBytes, _ := json.Marshal(listResp.Result)
	listStr := string(listBytes)
	requiredTools := []string{"bateia_recommend", "bateia_start_recommendation", "bateia_get_run", "bateia_get_artifact", "bateia_cancel_run"}
	for _, tool := range requiredTools {
		if !strings.Contains(listStr, `"`+tool+`"`) {
			t.Fatalf("ferramenta esperada %s não encontrada em tools/list: %s", tool, listStr)
		}
	}

	// 3. bateia_recommend com payload mínimo (validação E2E do autopreenchimento V1 do Passo 0031)
	minimalSyncPayload := `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bateia_recommend","arguments":{"idea":{"problem":"caching mechanism","capabilities":["lru_eviction"]}}}}`
	syncCallResp := step0032RPCCall(t, reqWriter, scanner, minimalSyncPayload)
	syncResult := step0032ExtractTextResult(t, syncCallResp)
	if syncResult == nil {
		t.Fatalf("bateia_recommend com defaults V1 retornou resultado vazio")
	}
	if status, ok := syncResult["status"].(string); !ok || status == "" {
		t.Fatalf("bateia_recommend não retornou status válido: %#v", syncResult)
	}
	// Valida se os defaults V1 foram materializados na decisão/alternativas
	if decision, ok := syncResult["decision"].(map[string]interface{}); !ok {
		t.Fatalf("bateia_recommend não retornou objeto de decisão: %#v", syncResult)
	} else if policy, ok := decision["policy"].(map[string]interface{}); !ok {
		t.Fatalf("bateia_recommend não materializou política padrão: %#v", decision)
	} else if _, hasLicenses := policy["allowed_licenses"]; !hasLicenses {
		t.Fatalf("bateia_recommend não preencheu allowed_licenses padrão: %#v", policy)
	}

	// 4. bateia_start_recommendation com idempotency_key e payload mínimo
	startPayload := `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"bateia_start_recommendation","arguments":{"idempotency_key":"e2e-step0032-test","idea":{"problem":"caching mechanism","capabilities":["lru_eviction"]}}}}`
	startCallResp := step0032RPCCall(t, reqWriter, scanner, startPayload)
	startResult := step0032ExtractTextResult(t, startCallResp)

	runID, ok := startResult["run_id"].(string)
	if !ok || runID == "" {
		t.Fatalf("bateia_start_recommendation não retornou run_id válido: %#v", startResult)
	}
	if startResult["lifecycle"] != "QUEUED" {
		t.Fatalf("lifecycle inicial esperado QUEUED, recebido: %v", startResult["lifecycle"])
	}

	// 5. bateia_get_run por run_id até resolução do lifecycle
	var terminalResult map[string]interface{}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		getCallResp := step0032RPCCall(t, reqWriter, scanner, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"bateia_get_run","arguments":{"run_id":"`+runID+`"}}}`)
		terminalResult = step0032ExtractTextResult(t, getCallResp)
		if lc, _ := terminalResult["lifecycle"].(string); lc == "SUCCEEDED" || lc == "FAILED" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if terminalResult["lifecycle"] != "SUCCEEDED" {
		t.Fatalf("execução durável não concluiu com SUCCEEDED: %#v", terminalResult)
	}

	// 6. bateia_get_artifact para leitura segura de artefato persistido em disco
	artifactResp := step0032RPCCall(t, reqWriter, scanner, `{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"bateia_get_artifact","arguments":{"run_id":"`+runID+`","artifact":"recommendation","cursor":"","limit":1}}}`)
	artifactResult := step0032ExtractTextResult(t, artifactResp)
	if artifactResult["artifact"] != "recommendation" {
		t.Fatalf("artefato retornado incorreto: %#v", artifactResult)
	}
	if hash, _ := artifactResult["artifact_hash"].(string); !strings.HasPrefix(hash, "sha256:") {
		t.Fatalf("hash de artefato ausente ou inválido: %v", hash)
	}

	// 7. bateia_cancel_run
	cancelResp := step0032RPCCall(t, reqWriter, scanner, `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"bateia_cancel_run","arguments":{"run_id":"`+runID+`"}}}`)
	cancelResult := step0032ExtractTextResult(t, cancelResp)
	if cancelResult["lifecycle"] != "SUCCEEDED" {
		t.Fatalf("cancel_run não preservou estado terminal: %#v", cancelResult)
	}

	// Finalização graciosa
	_ = reqWriter.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = shutdown(ctx)
	_ = <-serveDone
}
