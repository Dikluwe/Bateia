/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/mcp/request_dispatch_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"bateia/01_core/domain"
)

type securityImprovementStub struct{}

func (securityImprovementStub) StartImprovementSynthesis(context.Context, string, domain.ImprovementSynthesisRequest) (domain.ImprovementSynthesisStartReceipt, error) {
	return domain.ImprovementSynthesisStartReceipt{}, errors.New("unexpected start")
}

func securityCall(id int, name, arguments string) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"tools/call","params":{"name":%q,"arguments":%s}}`, id, name, arguments)
}

func securityResponses(t *testing.T, output *bytes.Buffer) []JSONRPCResponse {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
	var result []JSONRPCResponse
	for {
		var response JSONRPCResponse
		err := decoder.Decode(&response)
		if err == io.EOF {
			return result
		}
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, response)
	}
}

func securityConfigure(s *MCPServer, mask int, reverse bool) {
	registrations := []func(){
		func() {
			if mask&1 != 0 {
				s.WithAdvancedIntake()
			}
		},
		func() {
			if mask&2 != 0 {
				s.WithDurableRuns(step0027DurableStub{})
			}
		},
		func() {
			if mask&4 != 0 {
				s.WithMarketDiscovery([]domain.MarketDiscoveryConnector{&step0027ConnectorStub{}}, &step0027StarterStub{})
			}
		},
		func() {
			if mask&8 != 0 {
				s.WithImprovementSynthesis(securityImprovementStub{})
			}
		},
	}
	if reverse {
		for i := len(registrations) - 1; i >= 0; i-- {
			registrations[i]()
		}
	} else {
		for _, register := range registrations {
			register()
		}
	}
}

func TestSecurityEveryToolSharesOneBudgetForEveryExtensionCombination(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		for _, reverse := range []bool{false, true} {
			t.Run(fmt.Sprintf("mask%d_reverse%v", mask, reverse), func(t *testing.T) {
				names := []string{"bateia_recommend"}
				if mask&1 != 0 {
					names = append(names, "bateia_prepare_recommendation")
				}
				if mask&2 != 0 {
					names = append(names, "bateia_get_run")
				}
				if mask&4 != 0 {
					names = append(names, "bateia_list_connectors")
				}
				if mask&8 != 0 {
					names = append(names, "bateia_start_improvement_synthesis")
				}
				for _, first := range names {
					calls := []string{securityCall(1, first, `{"input":{}}`), `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`}
					for i, name := range names {
						calls = append(calls, securityCall(i+3, name, `{}`))
					}
					var output bytes.Buffer
					server := NewMCPServer(nil, strings.NewReader(strings.Join(calls, "\n")), &output)
					securityConfigure(server, mask, reverse)
					limit := 1
					budget, _ := NewToolCallBudget(&limit)
					server.WithToolCallBudget(budget)
					if err := server.Serve(context.Background()); err != nil {
						t.Fatal(err)
					}
					responses := securityResponses(t, &output)
					if len(responses) != len(calls) {
						t.Fatalf("responses=%d calls=%d", len(responses), len(calls))
					}
					if responses[0].Error != nil && responses[0].Error.Code == -32010 {
						t.Fatal("first call incorrectly exhausted")
					}
					if responses[1].Error != nil {
						t.Fatalf("tools/list failed: %#v", responses[1])
					}
					for _, response := range responses[2:] {
						if response.Error == nil || response.Error.Code != -32010 {
							t.Fatalf("tool escaped budget after %s: %#v", first, response)
						}
					}
				}
			})
		}
	}
}

func TestSecurityRepeatedPrepareCannotBypassBudget(t *testing.T) {
	call := securityCall(1, "bateia_prepare_recommendation", `{"input":{}}`)
	var output bytes.Buffer
	server := NewMCPServer(nil, strings.NewReader(call+"\n"+call+"\n"+call), &output).WithAdvancedIntake()
	limit := 1
	budget, _ := NewToolCallBudget(&limit)
	server.WithToolCallBudget(budget)
	if err := server.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	responses := securityResponses(t, &output)
	if len(responses) != 3 || responses[0].Error != nil || responses[1].Error == nil || responses[2].Error == nil || responses[1].Error.Code != -32010 || responses[2].Error.Code != -32010 {
		t.Fatalf("responses=%#v", responses)
	}
}

type securityCountingReader struct {
	source     io.Reader
	bytes      int
	failIfRead bool
}

func (r *securityCountingReader) Read(target []byte) (int, error) {
	if r.failIfRead {
		return 0, errors.New("cancelled context must not read")
	}
	n, err := r.source.Read(target)
	r.bytes += n
	return n, err
}

func TestSecurityFramesBoundedBeforeAllExtensionsDecode(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		for _, suffix := range []string{"\n", ""} {
			t.Run(fmt.Sprintf("mask%d_lf%v", mask, suffix != ""), func(t *testing.T) {
				source := &securityCountingReader{source: strings.NewReader(strings.Repeat("x", 128<<10) + suffix)}
				var output bytes.Buffer
				server := NewMCPServer(nil, source, &output)
				securityConfigure(server, mask, false)
				if err := server.Serve(context.Background()); err == nil {
					t.Fatal("oversized frame accepted")
				}
				if source.bytes > maxMCPMessageBytes+1 {
					t.Fatalf("unbounded read: %d", source.bytes)
				}
				if output.Len() != 0 {
					t.Fatalf("oversized message dispatched: %s", output.String())
				}
			})
		}
	}
}

func TestSecurityStructureLimitsRejectBeforeHandlerAndRecover(t *testing.T) {
	for _, value := range []string{strings.Repeat("[", 40) + "0" + strings.Repeat("]", 40), "[" + strings.Repeat("0,", 10000) + "0]"} {
		bad := securityCall(1, "bateia_prepare_recommendation", value)
		good := securityCall(2, "bateia_prepare_recommendation", `{"input":{}}`)
		var output bytes.Buffer
		server := NewMCPServer(nil, strings.NewReader(bad+"\n"+good), &output).WithAdvancedIntake()
		if err := server.Serve(context.Background()); err != nil {
			t.Fatal(err)
		}
		responses := securityResponses(t, &output)
		if len(responses) != 2 || responses[0].Error == nil || responses[0].Error.Code != -32700 || responses[1].Error != nil {
			t.Fatalf("responses=%#v", responses)
		}
	}
}

func TestSecurityCancelledSessionDoesNotReadOrDispatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for mask := 0; mask < 16; mask++ {
		var output bytes.Buffer
		server := NewMCPServer(nil, &securityCountingReader{failIfRead: true}, &output)
		securityConfigure(server, mask, false)
		if err := server.Serve(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("mask%d err=%v", mask, err)
		}
		if output.Len() != 0 {
			t.Fatalf("cancelled dispatched: %s", output.String())
		}
	}
}

type securityContextDurable struct {
	step0027DurableStub
	got context.Context
}

func (s *securityContextDurable) GetRun(ctx context.Context, _ domain.RunID) (domain.RecommendationRunView, error) {
	s.got = ctx
	return domain.RecommendationRunView{}, domain.ErrRunNotFound
}

type securityContextMarket struct{ got context.Context }

func (s *securityContextMarket) StartMarketDiscovery(ctx context.Context, _ string, _ domain.MarketDiscoveryStartRequest) (domain.MarketDiscoveryStartReceipt, error) {
	s.got = ctx
	return domain.MarketDiscoveryStartReceipt{}, domain.ErrInvalidRunInput
}

func TestSecurityDurableReceivesSessionContextInEveryRegistrationOrder(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		var output bytes.Buffer
		spy := &securityContextDurable{}
		server := NewMCPServer(nil, strings.NewReader(securityCall(1, "bateia_get_run", `{"run_id":"run_0123456789abcdef0123456789abcdef"}`)), &output)
		if reverse {
			securityConfigure(server, 13, true)
			server.WithDurableRuns(spy)
		} else {
			server.WithDurableRuns(spy)
			securityConfigure(server, 13, false)
		}
		type contextKey string
		ctx := context.WithValue(context.Background(), contextKey("marker"), "session")
		if err := server.Serve(ctx); err != nil {
			t.Fatal(err)
		}
		if spy.got != ctx {
			t.Fatal("session context replaced by Background")
		}
	}
}

func TestSecurityExtensionRegistrationOrderPreservesToolInventory(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		var output bytes.Buffer
		server := NewMCPServer(nil, strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`), &output)
		securityConfigure(server, 15, reverse)
		securityConfigure(server, 15, !reverse)
		if err := server.Serve(context.Background()); err != nil {
			t.Fatal(err)
		}
		responses := securityResponses(t, &output)
		tools := responseMap(t, responses[0].Result)["tools"].([]interface{})
		if len(tools) != 11 {
			t.Fatalf("missing or duplicate tools: %d", len(tools))
		}
		seen := map[string]bool{}
		for _, tool := range tools {
			name := tool.(map[string]interface{})["name"].(string)
			if seen[name] {
				t.Fatalf("duplicate %s", name)
			}
			seen[name] = true
		}
	}
}

func TestSecurityUnknownAndMalformedCallsConsumeSharedBudget(t *testing.T) {
	for _, first := range []string{
		securityCall(1, "unknown_tool", `{}`),
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":"bad"}`,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"bateia_start_market_discovery","arguments":`,
	} {
		var output bytes.Buffer
		server := NewMCPServer(nil, strings.NewReader(first+"\n"+securityCall(2, "bateia_prepare_recommendation", `{"input":{}}`)), &output)
		securityConfigure(server, 15, false)
		limit := 1
		budget, _ := NewToolCallBudget(&limit)
		server.WithToolCallBudget(budget)
		if err := server.Serve(context.Background()); err != nil {
			t.Fatal(err)
		}
		responses := securityResponses(t, &output)
		if len(responses) != 2 || responses[0].Error == nil || responses[0].Error.Code == -32010 || responses[1].Error == nil || responses[1].Error.Code != -32010 {
			t.Fatalf("first=%s responses=%#v", first, responses)
		}
	}
}
