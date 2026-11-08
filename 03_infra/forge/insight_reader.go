/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/insight_reader.md
 * @layer L3
 * @updated 2026-10-28
 */
package forge

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

// InsightReader has no cache or shared mutable per-repository state.
type InsightReader struct{ forge *BaseForgeConnector }

func NewGiteeInsightReader(token string, client *http.Client) *InsightReader {
	c := newBaseForge(domain.ConnectorGiteeCatalog, "readme-v1", "https://gitee.com/api/v5", client)
	c.token = strings.TrimSpace(token)
	return &InsightReader{forge: c}
}
func NewAtomGitInsightReader(token string, client *http.Client) *InsightReader {
	c := newBaseForge(domain.ConnectorAtomGitCatalog, "readme-v1", "https://api.atomgit.com/api/v5", client)
	c.token = strings.TrimSpace(token)
	return &InsightReader{forge: c}
}

func readmeFailure(code string) error { return domain.InsightReadmeFailure{Code: code} }
func immutableReadmeSHA(sha string) bool {
	if len(sha) != 40 && len(sha) != 64 {
		return false
	}
	for _, r := range sha {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}

func (r *InsightReader) endpoint(identity domain.ProductIdentity) (string, error) {
	if identity.ConnectorID != r.forge.id || !validForgeIdentity(identity.ConnectorID, identity.ExternalID, identity.CanonicalURL) || len(strings.Split(identity.ExternalID, "/")) != 2 {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	for _, c := range identity.ExternalID {
		if c > unicode.MaxASCII {
			return "", readmeFailure("INTEGRITY_ERROR")
		}
	}
	u, _ := url.Parse(identity.CanonicalURL)
	if u.RawPath != "" {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	base := r.forge.baseURL
	if identity.ConnectorID == domain.ConnectorAtomGitCatalog && u.Host == "gitcode.com" {
		base = "https://api.gitcode.com/api/v5"
	}
	return base + "/repos/" + identity.ExternalID, nil
}

func (r *InsightReader) get(ctx context.Context, endpoint string, params url.Values, limit int, output any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if r.forge.id == domain.ConnectorAtomGitCatalog && r.forge.token == "" {
		return readmeFailure("AUTH_REQUIRED")
	}
	if params == nil {
		params = url.Values{}
	}
	if r.forge.id == domain.ConnectorGiteeCatalog && r.forge.token != "" {
		params.Set("access_token", r.forge.token)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return readmeFailure("READ_FAILED")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Bateia-Insight-Reader/1.0")
	if r.forge.id == domain.ConnectorAtomGitCatalog {
		req.Header.Set("Authorization", "Bearer "+r.forge.token)
	}
	resp, err := r.forge.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return readmeFailure("READ_FAILED")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
	case 401, 403:
		return readmeFailure("AUTH_REQUIRED")
	case 404:
		return readmeFailure("NOT_FOUND")
	case 429:
		return readmeFailure("RATE_LIMIT")
	default:
		return readmeFailure("READ_FAILED")
	}
	if resp.Header.Get("X-Next-Page") != "" || strings.Contains(resp.Header.Get("Link"), `rel="next"`) {
		return readmeFailure("TREE_LIMIT")
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(limit+1)))
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil || len(body) > limit || json.Unmarshal(body, output) != nil {
		return readmeFailure("INTEGRITY_ERROR")
	}
	return nil
}

func (r *InsightReader) ResolveRepository(ctx context.Context, identity domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	endpoint, err := r.endpoint(identity)
	if err != nil {
		return domain.RepositoryCandidate{}, err
	}
	var metadata struct {
		genericForgeItem
		DefaultBranch string `json:"default_branch"`
	}
	if err := r.get(ctx, endpoint, nil, 256<<10, &metadata); err != nil {
		return domain.RepositoryCandidate{}, err
	}
	rawURL := metadata.HTMLURL
	if rawURL == "" {
		rawURL = metadata.WebURL
	}
	branch := metadata.DefaultBranch
	if metadata.FullName != identity.ExternalID || rawURL != identity.CanonicalURL || !publicForgeItem(r.forge.id, metadata.genericForgeItem) || branch == "" || len(branch) > 256 || strings.TrimSpace(branch) != branch {
		return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
	}
	for _, part := range strings.Split(branch, "/") {
		if part == "" || part == "." || part == ".." || strings.HasSuffix(part, ".") || strings.HasSuffix(part, ".lock") {
			return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
		}
	}
	for _, c := range branch {
		if unicode.IsControl(c) || unicode.IsSpace(c) || strings.ContainsRune("~^:?*[\\", c) {
			return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
		}
	}
	var commit struct {
		SHA string `json:"sha"`
	}
	if err := r.get(ctx, endpoint+"/commits/"+url.PathEscape(branch), nil, 256<<10, &commit); err != nil {
		return domain.RepositoryCandidate{}, err
	}
	if !immutableReadmeSHA(commit.SHA) {
		return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
	}
	parts := strings.Split(identity.ExternalID, "/")
	return domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: identity.CanonicalURL, CommitID: commit.SHA}, nil
}

func (r *InsightReader) repositoryEndpoint(repo domain.RepositoryCandidate) (string, error) {
	if !immutableReadmeSHA(repo.CommitID) {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	return r.endpoint(domain.ProductIdentity{ConnectorID: r.forge.id, ExternalID: repo.Owner + "/" + repo.Repo, CanonicalURL: repo.URL})
}

func (r *InsightReader) FetchTree(ctx context.Context, repo domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	endpoint, err := r.repositoryEndpoint(repo)
	if err != nil {
		return nil, err
	}
	var tree struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
			Mode string `json:"mode"`
		} `json:"tree"`
		Truncated bool `json:"truncated"`
	}
	if err := r.get(ctx, endpoint+"/git/trees/"+repo.CommitID, url.Values{"page": {"1"}, "per_page": {"100"}}, 256<<10, &tree); err != nil {
		return nil, err
	}
	if tree.Truncated || len(tree.Tree) >= 100 {
		return nil, readmeFailure("TREE_LIMIT")
	}
	if tree.Tree == nil {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	nodes := []domain.DirectoryNode{}
	for _, node := range tree.Tree {
		if node.Type == "blob" && (node.Mode == "100644" || node.Mode == "100755") && !strings.ContainsAny(node.Path, "/\\") {
			nodes = append(nodes, domain.DirectoryNode{Path: node.Path, Type: "file"})
		}
	}
	return nodes, nil
}

func (r *InsightReader) ExtractRepoFiles(ctx context.Context, repo domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	endpoint, err := r.repositoryEndpoint(repo)
	if err != nil {
		return nil, err
	}
	if len(paths) != 1 {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	path := paths[0]
	switch strings.ToLower(path) {
	case "readme.md", "readme.rst", "readme.txt", "readme":
	default:
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	var file struct {
		Type     string `json:"type"`
		Encoding string `json:"encoding"`
		Path     string `json:"path"`
		Size     *int   `json:"size"`
		SHA      string `json:"sha"`
		Content  string `json:"content"`
	}
	if err := r.get(ctx, endpoint+"/contents/"+url.PathEscape(path), url.Values{"ref": {repo.CommitID}}, 64<<10, &file); err != nil {
		return nil, err
	}
	if file.Type != "file" || file.Encoding != "base64" || file.Path != path || file.Size == nil || *file.Size < 0 || !immutableReadmeSHA(file.SHA) {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	if *file.Size > 32<<10 {
		return nil, readmeFailure("CONTENT_LIMIT")
	}
	content, err := base64.StdEncoding.Strict().DecodeString(file.Content)
	if err != nil || len(content) != *file.Size {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	object := append([]byte(fmt.Sprintf("blob %d\x00", len(content))), content...)
	actual := ""
	if len(file.SHA) == 40 {
		sum := sha1.Sum(object)
		actual = hex.EncodeToString(sum[:])
	} else {
		sum := sha256.Sum256(object)
		actual = hex.EncodeToString(sum[:])
	}
	if actual != file.SHA {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}, Files: []domain.ExtractedFile{{Path: path, Content: string(content)}}}, nil
}
