/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/report/link_safety.md
 * @layer L2
 * @updated 2026-09-26
 */

package report

import (
	"html"
	"net/url"
	"strings"
	"unicode"
)

// safeReportURL allows web links and, only when explicitly requested, local
// absolute file links. It does not authorize reading or executing the target.
func safeReportURL(raw string, allowLocal bool) string {
	if strings.IndexFunc(raw, unicode.IsControl) >= 0 || strings.Contains(raw, "\\") {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Opaque != "" || u.User != nil {
		return ""
	}
	decoded, err := url.PathUnescape(raw)
	if err != nil || strings.IndexFunc(decoded, unicode.IsControl) >= 0 || strings.Contains(decoded, "\\") {
		return ""
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		if u.Hostname() == "" {
			return ""
		}
	case "file":
		if !allowLocal || (u.Host != "" && u.Host != "localhost") || !strings.HasPrefix(u.Path, "/") || u.RawQuery != "" || u.Fragment != "" {
			return ""
		}
	default:
		return ""
	}
	return u.String()
}

// reportHTMLLink keeps both href and label in their own HTML contexts. An
// invalid or unsafe URL is represented as plain text, never an active link.
func reportHTMLLink(raw, label string, allowLocal bool) string {
	text := html.EscapeString(label)
	if safe := safeReportURL(raw, allowLocal); safe != "" {
		return `<a href="` + html.EscapeString(safe) + `" target="_blank" rel="noopener noreferrer">` + text + `</a>`
	}
	return text
}
