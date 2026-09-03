/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/main_step0025_contract_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package main

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type step0025WiringDiscovery struct{ calls int }

func (p *step0025WiringDiscovery) Discover(_ context.Context, request domain.DiscoveryRequest) (domain.ProviderDiscovery, error) {
	p.calls++
	return domain.ProviderDiscovery{Provider: "l4-offline", Attempts: []domain.ProviderAttempt{{Query: request.Query, Status: "SUCCEEDED"}}, Candidates: []domain.RepositoryCandidate{}}, nil
}

type step0025WiringContent struct{ calls int }

func (p *step0025WiringContent) FetchManifest(context.Context, domain.RepositoryCandidate) (*domain.Manifest, error) {
	p.calls++
	return nil, errors.New("not expected")
}
func (p *step0025WiringContent) FetchTree(context.Context, domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	p.calls++
	return nil, nil
}
func (p *step0025WiringContent) ExtractRepoFiles(context.Context, domain.RepositoryCandidate, []string) (*domain.ExtractedRepository, error) {
	p.calls++
	return nil, errors.New("not expected")
}
func (p *step0025WiringContent) FetchRepoMetadata(_ context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	p.calls++
	return candidate, nil
}

func TestStep0025RealL4CompositionListsAndCallsRecommendationOffline(t *testing.T) {
	discovery := &step0025WiringDiscovery{}
	content := &step0025WiringContent{}
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_recommend","arguments":{"idea":{"problem":"food waste","capabilities":["inventory"]},"policy":{"schema_version":"1","allowed_licenses":["MIT"],"require_local":false,"allow_archived":false,"max_components":2,"extend_min_coverage":0.7,"integration_min_value":0.3},"target":{"name":""},"search":{"strategy":"literal","language":"","locales":["pt-BR"],"max_per_query":1,"max_candidates":1}}}}`,
	}, "\n") + "\n"
	writer := &bytes.Buffer{}
	server := newMCPServer((*usecase.SearchFunnel)(nil), nil, discovery, content, strings.NewReader(input), writer)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(writer.Bytes()))
	var listed, called struct {
		Result map[string]interface{} `json:"result"`
		Error  interface{}            `json:"error"`
	}
	if err := decoder.Decode(&listed); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&called); err != nil {
		t.Fatal(err)
	}
	tools, ok := listed.Result["tools"].([]interface{})
	if !ok || len(tools) != 3 || called.Error != nil || discovery.calls == 0 {
		t.Fatalf("composição L4 não roteou offline: tools=%#v call=%#v discovery=%d content=%d", listed, called, discovery.calls, content.calls)
	}
}
