/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/web/authenticated_client.md
 * @layer L3
 * @updated 2026-08-31
 */
package web

import (
	"bateia/01_core/domain"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func readRateLimit(header http.Header) domain.WebRateLimit {
	parse := func(name string) *int64 {
		value := strings.TrimSpace(header.Get(name))
		if value == "" {
			return nil
		}
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil
		}
		return &n
	}
	return domain.WebRateLimit{Limit: parse("X-RateLimit-Limit"), Remaining: parse("X-RateLimit-Remaining"), Used: parse("X-RateLimit-Used"), Reset: parse("X-RateLimit-Reset")}
}

func containsSensitiveWebValue(value any, token string) bool {
	b, _ := json.Marshal(value)
	text := strings.ToLower(string(b))
	for _, forbidden := range []string{"authorization: bearer", "authorization=", "bearer ", "set-cookie", "cookie="} {
		if strings.Contains(text, forbidden) {
			return true
		}
	}
	return token != "" && strings.Contains(string(b), token)
}
