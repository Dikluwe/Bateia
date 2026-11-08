/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/market_composition_test.md
 * @layer L4
 * @updated 2026-10-28
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/mcp"
	"bateia/03_infra/runstore"
)

type compositionServerFactory func([]domain.MarketDiscoveryConnector, string, []byte, io.Reader, io.Writer) (*mcp.MCPServer, func(context.Context) error, error)

func compositionFactories() map[string]compositionServerFactory {
	limits := runstore.RunStoreLimits{MaxArtifactBytes: 1 << 20, MaxRunBytes: 4 << 20, MaxTotalBytes: 32 << 20}
	recommendation := usecase.NewRecommendationService(nil, &step0025WiringDiscovery{}, &step0025WiringContent{})
	return map[string]compositionServerFactory{
		"market": func(connectors []domain.MarketDiscoveryConnector, root string, catalog []byte, reader io.Reader, writer io.Writer) (*mcp.MCPServer, func(context.Context) error, error) {
			return NewMarketDiscoveryMCPServer(nil, recommendation, connectors, root, []byte("0123456789abcdef0123456789abcdef"), limits, 1, 4, reader, writer)
		},
		"synthesis": func(connectors []domain.MarketDiscoveryConnector, root string, catalog []byte, reader io.Reader, writer io.Writer) (*mcp.MCPServer, func(context.Context) error, error) {
			_, digest := usecase.ActiveMechanismRuleCatalog()
			return NewImprovementSynthesisMCPServer(nil, recommendation, connectors, catalog, digest, nil, root, []byte("0123456789abcdef0123456789abcdef"), limits, 1, 4, reader, writer)
		},
	}
}

func TestConfiguredMarketCompositionPublishesRegisteredSourcesAndReopens(t *testing.T) {
	for _, profile := range []string{"single", "legacy", "default", "default-bitbucket", "federated", "federated-bitbucket"} {
		t.Run(profile, func(t *testing.T) {
			t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
			t.Setenv("BATEIA_GITEE_TOKEN", "")
			t.Setenv("BATEIA_ATOMGIT_TOKEN", "")
			if strings.HasSuffix(profile, "-bitbucket") {
				t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "example")
			}
			var connectors []domain.MarketDiscoveryConnector
			var err error
			if strings.HasPrefix(profile, "federated") {
				connectors, err = federatedMarketDiscoveryConnectors()
			} else {
				connectors, err = defaultMarketDiscoveryConnectors()
			}
			if err != nil {
				t.Fatal(err)
			}
			for _, id := range []string{domain.ConnectorGiteeCatalog, domain.ConnectorAtomGitCatalog} {
				found := false
				for _, c := range connectors {
					if c.Capabilities().ID == id {
						found = true
					}
				}
				if !found {
					t.Fatalf("missing configured connector %s", id)
				}
			}
			if profile == "single" {
				connectors = connectors[:1]
			}
			if profile == "legacy" {
				connectors = connectors[:3]
			}
			registry, err := usecase.NewDiscoveryConnectorRegistry(connectors)
			if err != nil {
				t.Fatal(err)
			}
			expected := []string{}
			for _, capability := range registry.Capabilities() {
				expected = append(expected, capability.ID)
			}
			catalog, _ := usecase.ActiveMechanismRuleCatalog()
			for name, factory := range compositionFactories() {
				t.Run(name, func(t *testing.T) {
					root := filepath.Join(t.TempDir(), "runs")
					output := &bytes.Buffer{}
					input := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}` + "\n" + `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"bateia_list_connectors","arguments":{}}}` + "\n")
					server, shutdown, err := factory(connectors, root, catalog, input, output)
					if err != nil {
						t.Fatalf("valid %s registry of %d sources rejected: %v", profile, len(connectors), err)
					}
					closed := false
					defer func() {
						if !closed {
							ctx, cancel := context.WithTimeout(context.Background(), time.Second)
							defer cancel()
							if err := shutdown(ctx); err != nil {
								t.Error(err)
							}
						}
					}()
					if err := server.Serve(context.Background()); err != nil {
						t.Fatal(err)
					}
					decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
					var listed, called mcp.JSONRPCResponse
					if err := decoder.Decode(&listed); err != nil {
						t.Fatal(err)
					}
					if err := decoder.Decode(&called); err != nil {
						t.Fatal(err)
					}
					if listed.Error != nil {
						t.Fatal(listed.Error)
					}
					raw, _ := json.Marshal(listed.Result)
					var tools struct {
						Tools []struct {
							Name        string `json:"name"`
							InputSchema struct {
								Properties struct {
									ConnectorIDs struct {
										Items struct {
											Enum []string `json:"enum"`
										} `json:"items"`
										MaxItems int `json:"maxItems"`
									} `json:"connector_ids"`
								} `json:"properties"`
							} `json:"inputSchema"`
						} `json:"tools"`
					}
					if err := json.Unmarshal(raw, &tools); err != nil {
						t.Fatal(err)
					}
					found, hasSynthesis := false, false
					for _, tool := range tools.Tools {
						if tool.Name == "bateia_start_improvement_synthesis" {
							hasSynthesis = true
						}
						if tool.Name == "bateia_start_market_discovery" {
							found = true
							ids := tool.InputSchema.Properties.ConnectorIDs.Items.Enum
							if !reflect.DeepEqual(ids, expected) {
								t.Errorf("advertised sources %v differ from registered %v", ids, expected)
							}
							max := tool.InputSchema.Properties.ConnectorIDs.MaxItems
							if max < 1 || max > 10 || len(expected) > 10 && max != 10 {
								t.Errorf("invalid selection limit %d", max)
							}
						}
					}
					if !found || hasSynthesis != (name == "synthesis") {
						t.Fatalf("wrong feature registration: %s", raw)
					}
					result := step0026L4Text(t, called)
					gotJSON, _ := json.Marshal(result["connectors"])
					var got []domain.DiscoveryConnectorCapabilities
					if err := json.Unmarshal(gotJSON, &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, registry.Capabilities()) {
						t.Fatal("listed connector capabilities differ from registry")
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					if err := shutdown(ctx); err != nil {
						t.Fatal(err)
					}
					closed = true
					reopened, closeAgain, err := factory(connectors, root, catalog, strings.NewReader(""), io.Discard)
					if err != nil {
						t.Fatalf("shutdown did not release resources: %v", err)
					}
					if err := reopened.Serve(context.Background()); err != nil {
						t.Error(err)
					}
					if err := closeAgain(ctx); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}

func TestConfiguredMarketCompositionRejectsInvalidInputsBeforeCreatingResources(t *testing.T) {
	t.Setenv("BATEIA_BITBUCKET_WORKSPACES", "")
	t.Setenv("BATEIA_GITEE_TOKEN", "")
	t.Setenv("BATEIA_ATOMGIT_TOKEN", "")
	valid, err := defaultMarketDiscoveryConnectors()
	if err != nil {
		t.Fatal(err)
	}
	catalog, _ := usecase.ActiveMechanismRuleCatalog()
	invalid := map[string][]domain.MarketDiscoveryConnector{
		"empty": nil, "nil": {nil}, "duplicate": {valid[0], valid[0]},
		"invalid-capabilities": {usecase.CapabilityOnlyConnector{Capability: domain.DiscoveryConnectorCapabilities{ID: "invalid"}}},
	}
	for name, factory := range compositionFactories() {
		for condition, connectors := range invalid {
			t.Run(name+"/"+condition, func(t *testing.T) {
				root := filepath.Join(t.TempDir(), "must-not-exist")
				server, shutdown, err := factory(connectors, root, catalog, strings.NewReader(""), io.Discard)
				if !errors.Is(err, domain.ErrInvalidRunInput) || server != nil || shutdown != nil {
					t.Fatalf("invalid registry accepted: %v", err)
				}
				if _, err := os.Stat(root); !os.IsNotExist(err) {
					t.Fatalf("invalid registry created resources: %v", err)
				}
			})
		}
	}
	root := filepath.Join(t.TempDir(), "invalid-catalog")
	server, shutdown, err := compositionFactories()["synthesis"](valid, root, []byte(`{}`), strings.NewReader(""), io.Discard)
	if err == nil || server != nil || shutdown != nil {
		t.Fatal("corrupt catalog accepted")
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("invalid catalog created resources: %v", err)
	}
}
