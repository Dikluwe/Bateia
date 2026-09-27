/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/security.md
 * @layer L3
 * @updated 2026-09-26
 */
package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"bateia/01_core/domain"
)

const (
	maxFileBytes             = 2 << 20
	maxContentEnvelopeBytes  = 4 << 20
	maxTreeEnvelopeBytes     = 8 << 20
	maxMetadataEnvelopeBytes = 1 << 20
	maxExtractionBytes       = 8 << 20
	maxExtractionFiles       = 100
)

var errUnsafeRemote = errors.New("GitHub: resposta ou entrada fora dos limites seguros")

// Security failures are terminal: they must not become a weaker fallback read.
func unsafeRemote(message string) error { return fmt.Errorf("%w: %s", errUnsafeRemote, message) }

func readBounded(resp *http.Response, limit int64) ([]byte, error) {
	if resp.ContentLength > limit {
		return nil, unsafeRemote("Content-Length excede orçamento")
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > limit {
		return nil, unsafeRemote("resposta excede orçamento de bytes")
	}
	return body, nil
}

func decodeBounded(resp *http.Response, limit int64, target any) error {
	body, err := readBounded(resp, limit)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return unsafeRemote("JSON inválido: " + err.Error())
	}
	return nil
}

func validateRemoteCandidate(c domain.RepositoryCandidate) error {
	if len(c.Owner) > 100 || len(c.Repo) > 100 || len(c.CommitID) > 1024 {
		return unsafeRemote("identidade excessivamente longa")
	}
	if !implementationRepositoryPart(c.Owner, true) || !implementationRepositoryPart(c.Repo, false) {
		return unsafeRemote("owner/repo inválidos")
	}
	if c.CommitID != "" && !implementationSafePath(c.CommitID) {
		return unsafeRemote("revisão inválida")
	}
	for _, r := range c.CommitID {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("._-/", r)) {
			return unsafeRemote("revisão inválida")
		}
	}
	return nil
}

func remotePath(p string) (string, error) {
	if len(p) > 4096 || !implementationSafePath(p) {
		return "", unsafeRemote("caminho inválido")
	}
	parts := strings.Split(p, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/"), nil
}

func decodeRegularContent(item githubContentItem, expectedPath string) (string, error) {
	if item.Type != "file" || item.Target != "" || item.Path != expectedPath || !implementationSafePath(item.Path) {
		return "", unsafeRemote("arquivo não regular ou caminho divergente")
	}
	if item.Size < 0 || item.Size > maxFileBytes {
		return "", unsafeRemote("arquivo excede orçamento")
	}
	var data []byte
	switch item.Encoding {
	case "base64":
		var err error
		data, err = base64.StdEncoding.DecodeString(strings.ReplaceAll(item.Content, "\n", ""))
		if err != nil {
			return "", unsafeRemote("base64 inválido")
		}
	case "utf-8", "":
		data = []byte(item.Content)
	default:
		return "", unsafeRemote("encoding indisponível ou não suportado")
	}
	if len(data) > maxFileBytes {
		return "", unsafeRemote("arquivo decodificado excede orçamento")
	}
	return string(data), nil
}

func denyRedirect(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

func queryRef(ref string) string { return url.QueryEscape(ref) }

// Contents API may dereference symlinks that target another file in the repo.
// Therefore extraction first requires an explicit regular mode in the Git tree.
func (c *Client) requireRegularRemotePaths(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) error {
	revision := candidate.CommitID
	if revision == "" {
		revision = "HEAD"
	}
	endpoint := "/repos/" + candidate.Owner + "/" + candidate.Repo + "/git/trees/" + url.PathEscape(revision) + "?recursive=1"
	body, _, err := c.implementationGET(ctx, endpoint, maxTreeEnvelopeBytes)
	if err != nil {
		return err
	}
	var tree githubTreeResponse
	if err := json.Unmarshal(body, &tree); err != nil {
		return err
	}
	if tree.Truncated || len(tree.Tree) > 20_000 {
		return unsafeRemote("árvore truncada ou acima do orçamento")
	}
	regular := map[string]bool{}
	sizes := map[string]int64{}
	for _, item := range tree.Tree {
		if !implementationSafePath(item.Path) || item.Size < 0 {
			return unsafeRemote("entrada de árvore inválida")
		}
		if _, duplicate := regular[item.Path]; duplicate {
			return unsafeRemote("caminho duplicado na árvore")
		}
		sizes[item.Path] = item.Size
		regular[item.Path] = item.Type == "blob" && (item.Mode == "100644" || item.Mode == "100755") && item.Size <= maxFileBytes
	}
	var total int64
	for _, p := range paths {
		total += sizes[p]
		if total > maxExtractionBytes {
			return unsafeRemote("tamanhos da árvore excedem orçamento total de extração")
		}
		if !regular[p] {
			return unsafeRemote("arquivo ausente, não regular ou acima do orçamento: " + p)
		}
	}
	return nil
}
