/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/observable_client_step0020_contract_test.md
 * @layer L3
 * @updated 2026-08-31
 */
package web

import (
	"bateia/01_core/domain"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type step0020RoundTripper func(*http.Request) (*http.Response, error)

func (f step0020RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func step0020Response(status int, body string, headers map[string]string) *http.Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(body))}
}

func step0020Search(t *testing.T, token string, rt step0020RoundTripper) (domain.RawWebResponse, error) {
	t.Helper()
	c := NewObservableWebClient(&http.Client{Transport: rt}, token)
	return c.SearchOne(context.Background(), domain.ObservableWebRequest{RequestID: "request-0020", QueryID: "query-0020", Query: "public tools", MaxResults: 10})
}

func TestObservableWebClientStep0020Contract(t *testing.T) {
	const canary = "step0020-TOKEN-CANARY-never-persist"

	t.Run("bearer is sent but canary is redacted from response error and receipt", func(t *testing.T) {
		seen := ""
		got, err := step0020Search(t, canary, func(r *http.Request) (*http.Response, error) {
			seen = r.Header.Get("Authorization")
			return step0020Response(http.StatusForbidden, `{"message":"denied"}`, nil), nil
		})
		if seen != "Bearer "+canary {
			t.Fatalf("Authorization = %q", seen)
		}
		encoded, marshalErr := json.Marshal(got)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		observable := string(encoded)
		if err != nil {
			observable += err.Error()
		}
		if strings.Contains(observable, canary) || strings.Contains(observable, "Authorization") {
			t.Fatalf("credential leaked in observable artifact: %s", observable)
		}
	})

	t.Run("private result is rejected without persisted hit", func(t *testing.T) {
		got, _ := step0020Search(t, canary, func(*http.Request) (*http.Response, error) {
			return step0020Response(http.StatusOK, `{"items":[{"full_name":"secret/repo","html_url":"https://github.com/secret/repo","description":"private","private":true}]}`, nil), nil
		})
		if got.Status != "REJECTED_PRIVATE" || len(got.Results) != 0 {
			t.Fatalf("status=%q results=%+v", got.Status, got.Results)
		}
	})

	for _, status := range []int{http.StatusForbidden, http.StatusTooManyRequests} {
		status := status
		t.Run(http.StatusText(status), func(t *testing.T) {
			got, _ := step0020Search(t, canary, func(*http.Request) (*http.Response, error) {
				return step0020Response(status, `{}`, nil), nil
			})
			if got.Status != "UNKNOWN_RATE_LIMIT" || len(got.Results) != 0 {
				t.Fatalf("HTTP %d: status=%q results=%+v", status, got.Status, got.Results)
			}
		})
	}

	t.Run("context timeout is classified", func(t *testing.T) {
		c := NewObservableWebClient(&http.Client{Transport: step0020RoundTripper(func(r *http.Request) (*http.Response, error) {
			<-r.Context().Done()
			return nil, r.Context().Err()
		})}, canary)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		got, err := c.SearchOne(ctx, domain.ObservableWebRequest{RequestID: "request-timeout", QueryID: "query-timeout", Query: "x"})
		if !errors.Is(err, context.Canceled) || got.Status != "UNKNOWN_TIMEOUT" || got.RequestHash == "" || got.ResponseHash == "" {
			t.Fatalf("status=%q request_hash=%q response_hash=%q err=%v", got.Status, got.RequestHash, got.ResponseHash, err)
		}
	})

	t.Run("invalid JSON is unknown provider response and sealed", func(t *testing.T) {
		got, err := step0020Search(t, canary, func(*http.Request) (*http.Response, error) {
			return step0020Response(http.StatusOK, `{"items":`, nil), nil
		})
		if err == nil || got.Status != "UNKNOWN_PROVIDER_RESPONSE" || got.RequestHash == "" || got.ResponseHash == "" || len(got.Results) != 0 {
			t.Fatalf("status=%q request_hash=%q response_hash=%q results=%+v err=%v", got.Status, got.RequestHash, got.ResponseHash, got.Results, err)
		}
	})

	t.Run("empty 200 is EMPTY", func(t *testing.T) {
		got, err := step0020Search(t, canary, func(*http.Request) (*http.Response, error) {
			return step0020Response(http.StatusOK, `{"items":[]}`, nil), nil
		})
		if err != nil || got.Status != "EMPTY" || len(got.Results) != 0 {
			t.Fatalf("status=%q results=%+v err=%v", got.Status, got.Results, err)
		}
	})

	t.Run("rate limit headers preserve integer counters only", func(t *testing.T) {
		got, err := step0020Search(t, canary, func(*http.Request) (*http.Response, error) {
			return step0020Response(http.StatusOK, `{"items":[]}`, map[string]string{
				"X-RateLimit-Limit": "5000", "X-RateLimit-Remaining": "4999", "X-RateLimit-Reset": "1700000000", "X-RateLimit-Resource": "search",
			}), nil
		})
		if err != nil {
			t.Fatal(err)
		}
		encoded, marshalErr := json.Marshal(got.Receipt)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		text := string(encoded)
		for _, integer := range []string{"5000", "4999", "1700000000"} {
			if !strings.Contains(text, integer) {
				t.Errorf("receipt does not preserve integer %s: %s", integer, text)
			}
		}
		rateLimitEncoded, marshalErr := json.Marshal(got.Receipt.RateLimit)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		if strings.Contains(string(rateLimitEncoded), "search") || strings.Contains(text, `"5000"`) || strings.Contains(text, `"4999"`) || strings.Contains(text, `"1700000000"`) {
			t.Errorf("receipt preserved a non-integer rate-limit value: %s", text)
		}
	})
}
