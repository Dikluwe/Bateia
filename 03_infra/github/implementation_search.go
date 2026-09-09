/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/implementation_search.md
 * @layer L3
 * @updated 2026-09-08
 */
package github

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"sort"
	"strings"
	"unicode"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

type implementationTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size *int64 `json:"size"`
}

func remoteImplementationBudget(b domain.ImplementationSearchBudget) domain.ImplementationSearchBudget {
	if b.MaxFiles <= 0 || b.MaxFiles > 64 {
		b.MaxFiles = 64
	}
	if b.MaxBytes <= 0 || b.MaxBytes > 4*1024*1024 {
		b.MaxBytes = 4 * 1024 * 1024
	}
	if b.MaxFileBytes <= 0 || b.MaxFileBytes > 512*1024 {
		b.MaxFileBytes = 512 * 1024
	}
	if b.MaxResults <= 0 || b.MaxResults > 24 {
		b.MaxResults = 24
	}
	return b
}

func implementationSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func implementationRepositoryPart(s string, owner bool) bool {
	if s == "" || s == "." || s == ".." {
		return false
	}
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' {
			continue
		}
		if !owner && (r == '.' || r == '_') {
			continue
		}
		return false
	}
	return !owner || !strings.HasPrefix(s, "-") && !strings.HasSuffix(s, "-")
}

func implementationSafePath(s string) bool {
	if s == "" || s == "." || s == ".." || strings.HasPrefix(s, "/") || strings.HasPrefix(s, "../") || path.Clean(s) != s || strings.ContainsAny(s, "\\:") {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

// implementationGET never follows redirects, so the snapshot request and request budget remain exact.
func (c *Client) implementationGET(ctx context.Context, endpoint string, cap int64) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.baseURL, "/")+endpoint, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	client := *c.httpClient
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("GitHub HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, cap+1))
	if err != nil {
		return body, resp.StatusCode, err
	}
	if int64(len(body)) > cap {
		return body, resp.StatusCode, fmt.Errorf("resposta GitHub excede limite de bytes")
	}
	return body, resp.StatusCode, nil
}

// implementationPartialContent accounts for source already received even when the
// transport or a bounded JSON envelope ends in the middle of the response.
// This prefix is accounting only and is never admitted to the parser.
func implementationPartialContent(body []byte) []byte {
	decoder := json.NewDecoder(bytes.NewReader(body))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil
	}
	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			return nil
		}
		if token == "content" {
			raw := bytes.TrimSpace(body[decoder.InputOffset():])
			if len(raw) == 0 || raw[0] != ':' {
				return nil
			}
			raw = bytes.TrimSpace(raw[1:])
			if len(raw) == 0 || raw[0] != '"' {
				return nil
			}
			end := len(raw)
			escaped := false
			for i := 1; i < len(raw); i++ {
				if escaped {
					escaped = false
					continue
				}
				if raw[i] == '\\' {
					escaped = true
					continue
				}
				if raw[i] == '"' {
					end = i + 1
					break
				}
			}
			var content string
			if json.Unmarshal(raw[:end], &content) != nil {
				// A final escape may be incomplete; only remove that bounded suffix.
				for trim := 0; trim <= 6 && end-trim > 0; trim++ {
					prefix := append(append([]byte(nil), raw[:end-trim]...), '"')
					if json.Unmarshal(prefix, &content) == nil {
						break
					}
				}
			}
			data, _ := base64.StdEncoding.DecodeString(content)
			return data
		}
		var skip json.RawMessage
		if decoder.Decode(&skip) != nil {
			return nil
		}
	}
	return nil
}

// SearchImplementations collects verified Git blobs at the candidate's immutable commit.
func (c *Client) SearchImplementations(ctx context.Context, candidate domain.RepositoryCandidate, query string, budget domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	budget = remoteImplementationBudget(budget)
	report := domain.ImplementationSearchReport{Status: "partial", Method: "go_ast_query_v1", Budget: budget, Hits: []domain.ImplementationHit{}, Examined: []domain.ImplementationInput{}}
	fail := func(err error) (domain.ImplementationSearchReport, error) {
		report.Issues = append(report.Issues, err.Error())
		return report, err
	}
	if !implementationRepositoryPart(candidate.Owner, true) || !implementationRepositoryPart(candidate.Repo, false) || !implementationSHA(candidate.CommitID) {
		return fail(fmt.Errorf("identidade GitHub inválida: owner/repo e commit SHA completo obrigatórios"))
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	base := "/repos/" + candidate.Owner + "/" + candidate.Repo + "/git/"
	body, _, err := c.implementationGET(ctx, base+"trees/"+candidate.CommitID+"?recursive=1", 8*1024*1024)
	if err != nil {
		report.Limited = len(body) > 8*1024*1024
		return fail(err)
	}
	var tree struct {
		Truncated *bool                      `json:"truncated"`
		Tree      *[]implementationTreeEntry `json:"tree"`
	}
	if err := json.Unmarshal(body, &tree); err != nil {
		return fail(fmt.Errorf("árvore inválida: %w", err))
	}
	if tree.Tree == nil || tree.Truncated == nil {
		return fail(fmt.Errorf("árvore GitHub sem campos obrigatórios"))
	}
	var issues []string
	limited := false
	if *tree.Truncated {
		issues = append(issues, "árvore GitHub truncada")
		limited = true
	}
	entries := map[string]implementationTreeEntry{}
	ambiguous := map[string]bool{}
	for _, entry := range *tree.Tree {
		if !implementationSafePath(entry.Path) {
			issues = append(issues, "caminho inseguro na árvore: "+entry.Path)
			continue
		}
		if previous, ok := entries[entry.Path]; ok {
			sameSize := previous.Size == nil && entry.Size == nil || previous.Size != nil && entry.Size != nil && *previous.Size == *entry.Size
			if previous.Type != entry.Type || previous.Mode != entry.Mode || !strings.EqualFold(previous.SHA, entry.SHA) || !sameSize {
				ambiguous[entry.Path] = true
			}
			continue
		}
		entries[entry.Path] = entry
	}
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	files := []domain.ExtractedFile{}
	var bytesRead int64
	attempts := 0
	var terminal error
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			terminal = err
			break
		}
		entry := entries[name]
		if !strings.HasSuffix(name, ".go") || !parser.ImplementationSearchPathAllowed(name, query) {
			continue
		}
		if ambiguous[name] {
			issues = append(issues, "caminho ambíguo: "+name)
			continue
		}
		if entry.Type != "blob" || entry.Mode != "100644" && entry.Mode != "100755" {
			issues = append(issues, "arquivo não regular: "+name)
			continue
		}
		if entry.Size == nil || *entry.Size < 0 || !implementationSHA(entry.SHA) {
			issues = append(issues, "metadados de blob inválidos: "+name)
			continue
		}
		if attempts >= budget.MaxFiles || bytesRead >= budget.MaxBytes {
			limited = true
			issues = append(issues, "limite de coleta atingido")
			break
		}
		cap := min(budget.MaxFileBytes, budget.MaxBytes-bytesRead)
		if *entry.Size > cap {
			limited = true
			issues = append(issues, "limite de bytes impede leitura: "+name)
			continue
		}
		attempts++
		// Allow base64 line wrapping plus a bounded metadata envelope.
		body, status, fetchErr := c.implementationGET(ctx, base+"blobs/"+entry.SHA, ((cap+2)/3)*4+64*1024)
		var blob struct {
			SHA      string  `json:"sha"`
			Size     *int64  `json:"size"`
			Encoding string  `json:"encoding"`
			Content  *string `json:"content"`
		}
		decodeErr := json.Unmarshal(body, &blob)
		var data []byte
		var contentErr error
		if decodeErr == nil && blob.Content != nil && blob.Encoding == "base64" {
			data, contentErr = base64.StdEncoding.DecodeString(*blob.Content)
			bytesRead += int64(len(data))
		} else if decodeErr != nil {
			bytesRead += int64(len(implementationPartialContent(body)))
		}
		if fetchErr != nil {
			if int64(len(body)) > ((cap+2)/3)*4+64*1024 {
				limited = true
			}
			issues = append(issues, fmt.Sprintf("leitura %s: %v", name, fetchErr))
			if status == http.StatusForbidden || status == http.StatusTooManyRequests {
				terminal = fetchErr
				break
			}
			if ctx.Err() != nil {
				terminal = ctx.Err()
				break
			}
			continue
		}
		if int64(len(data)) > cap {
			limited = true
			issues = append(issues, "conteúdo excede orçamento: "+name)
			continue
		}
		if decodeErr != nil || contentErr != nil || blob.Content == nil || blob.Encoding != "base64" || blob.Size == nil || *blob.Size != int64(len(data)) || *entry.Size != int64(len(data)) || !strings.EqualFold(blob.SHA, entry.SHA) {
			issues = append(issues, "blob incompleto ou inválido: "+name)
			continue
		}
		hash := sha1.New()
		fmt.Fprintf(hash, "blob %d%c", len(data), 0)
		hash.Write(data)
		if !strings.EqualFold(hex.EncodeToString(hash.Sum(nil)), entry.SHA) {
			issues = append(issues, "identidade de blob divergente: "+name)
			continue
		}
		files = append(files, domain.ExtractedFile{Path: name, Content: string(data), Size: int64(len(data))})
	}
	if terminal == nil {
		terminal = ctx.Err()
	}
	report = parser.RankGoImplementations(query, files, budget.MaxResults)
	report.Budget = budget
	report.BytesRead = bytesRead
	report.Limited = report.Limited || limited
	report.Issues = append(report.Issues, issues...)
	if terminal != nil {
		report.Issues = append(report.Issues, terminal.Error())
	}
	sort.Strings(report.Issues)
	if report.Limited || len(report.Issues) > 0 {
		report.Status = "partial"
	}
	return report, terminal
}
