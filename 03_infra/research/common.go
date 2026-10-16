/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/research/common.md
 * @layer L3
 * @updated 2026-10-06
 */

package research

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bateia/01_core/domain"
)

const maxResponseBytes = 4 << 20

var locales = []string{"ar-SA", "de-DE", "en-US", "es-ES", "fr-FR", "hi-IN", "ja-JP", "ko-KR", "pt-BR", "ru-RU"}

type httpClient struct{ client *http.Client }

func newHTTPClient(client *http.Client) httpClient {
	if client == nil {
		client = &http.Client{}
	}
	copy := *client
	if copy.Timeout <= 0 || copy.Timeout > 30*time.Second {
		copy.Timeout = 20 * time.Second
	}
	copy.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("redirect forbidden") }
	return httpClient{client: &copy}
}

func (h httpClient) get(ctx context.Context, endpoint string) ([]byte, error) {
	status, body, err := h.getStatus(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("research catalogue returned HTTP %d", status)
	}
	return body, nil
}

func (h httpClient) getStatus(ctx context.Context, endpoint string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Bateia-Research-Discovery/1.0")
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil || len(body) > maxResponseBytes {
		return 0, nil, fmt.Errorf("research catalogue response exceeds limit or cannot be read")
	}
	return resp.StatusCode, body, nil
}

func validPublicURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Hostname() == "" || parsed.User != nil {
		return ""
	}
	return parsed.String()
}

func safeText(value string, max int) string {
	value = strings.TrimSpace(strings.ToValidUTF8(value, ""))
	if len(value) > max {
		value = value[:max]
	}
	return value
}

func responseHash(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

func result(status domain.MarketAttemptStatus, candidates []domain.ProductCandidate, hashes ...[]byte) domain.ConnectorDiscovery {
	if candidates == nil {
		candidates = []domain.ProductCandidate{}
	}
	combined := []byte{}
	for _, body := range hashes {
		combined = append(combined, body...)
	}
	return domain.ConnectorDiscovery{Status: status, Candidates: candidates, ResultCount: len(candidates), AcceptedCount: len(candidates), ResponseHash: responseHash(combined), ObservedAt: time.Now().UTC().Format(time.RFC3339)}
}

func validateRequest(request domain.MarketDiscoveryRequest, max int) (string, int, error) {
	query := strings.TrimSpace(request.Query)
	if query == "" || len(query) > 256 {
		return "", 0, domain.ErrInvalidRunInput
	}
	limit := request.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > max {
		limit = max
	}
	return query, limit, nil
}
