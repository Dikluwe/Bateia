/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/insight_reference_resolver.md
 * @layer L3
 * @updated 2026-10-28
 */
package github

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ResolveRepository pins a public GitHub reference before a bounded README read.
func (c *Client) ResolveRepository(ctx context.Context, identity domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	failure := func(code string) (domain.RepositoryCandidate, error) {
		return domain.RepositoryCandidate{}, domain.InsightReadmeFailure{Code: code}
	}
	parts := strings.Split(identity.ExternalID, "/")
	u, err := url.Parse(identity.CanonicalURL)
	if identity.ConnectorID != domain.ConnectorGitHubCatalog || len(parts) != 2 || err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.RawPath != "" || u.Path != "/"+identity.ExternalID {
		return failure("INTEGRITY_ERROR")
	}
	repo := domain.RepositoryCandidate{Owner: parts[0], Repo: parts[1], URL: identity.CanonicalURL}
	if validateRemoteCandidate(repo) != nil {
		return failure("INTEGRITY_ERROR")
	}

	var metadata struct {
		FullName      string `json:"full_name"`
		HTMLURL       string `json:"html_url"`
		Private       *bool  `json:"private"`
		DefaultBranch string `json:"default_branch"`
	}
	base := "/repos/" + repo.Owner + "/" + repo.Repo
	if err := c.getInsightReferenceJSON(ctx, base, 256<<10, &metadata); err != nil {
		return domain.RepositoryCandidate{}, err
	}
	if metadata.Private == nil || *metadata.Private || metadata.FullName != identity.ExternalID || metadata.HTMLURL != identity.CanonicalURL || (metadata.DefaultBranch == "" || len(metadata.DefaultBranch) > 256) {
		return failure("INTEGRITY_ERROR")
	}
	branch := repo
	branch.CommitID = metadata.DefaultBranch
	if validateRemoteCandidate(branch) != nil {
		return failure("INTEGRITY_ERROR")
	}
	var commit struct {
		SHA string `json:"sha"`
	}
	if err := c.getInsightReferenceJSON(ctx, base+"/commits/"+url.PathEscape(metadata.DefaultBranch), 256<<10, &commit); err != nil {
		return domain.RepositoryCandidate{}, err
	}
	if !implementationSHA(commit.SHA) {
		return failure("INTEGRITY_ERROR")
	}
	repo.CommitID = commit.SHA
	return repo, nil
}

func (c *Client) getInsightReferenceJSON(ctx context.Context, path string, limit int64, target any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return domain.InsightReadmeFailure{Code: "READ_FAILED"}
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		request.Header.Set("Authorization", "token "+c.token)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return domain.InsightReadmeFailure{Code: "READ_FAILED"}
	}
	defer response.Body.Close()
	code := ""
	switch response.StatusCode {
	case 200:
	case 401, 403:
		code = "AUTH_REQUIRED"
	case 404:
		code = "NOT_FOUND"
	case 429:
		code = "RATE_LIMIT"
	default:
		code = "READ_FAILED"
	}
	if code != "" {
		return domain.InsightReadmeFailure{Code: code}
	}
	if decodeBounded(response, limit, target) != nil {
		return domain.InsightReadmeFailure{Code: "INTEGRITY_ERROR"}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

// InsightLinkedReader keeps followup reads bounded without legacy raw fallbacks.
type InsightLinkedReader struct{ client *Client }

func NewInsightLinkedReader(client *Client) *InsightLinkedReader {
	if client == nil {
		client = NewClient("")
	}
	return &InsightLinkedReader{client}
}
func (r *InsightLinkedReader) ResolveRepository(ctx context.Context, id domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	return r.client.ResolveRepository(ctx, id)
}
func validInsightLinkedRepo(repo domain.RepositoryCandidate) bool {
	return validateRemoteCandidate(repo) == nil && implementationSHA(repo.CommitID) && repo.URL == "https://github.com/"+repo.Owner+"/"+repo.Repo
}
func (r *InsightLinkedReader) FetchTree(ctx context.Context, repo domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	if !validInsightLinkedRepo(repo) {
		return nil, domain.InsightReadmeFailure{Code: "INTEGRITY_ERROR"}
	}
	var tree githubTreeResponse
	if err := r.client.getInsightReferenceJSON(ctx, "/repos/"+repo.Owner+"/"+repo.Repo+"/git/trees/"+repo.CommitID, 256<<10, &tree); err != nil {
		return nil, err
	}
	if tree.Truncated || len(tree.Tree) >= 100 {
		return nil, domain.InsightReadmeFailure{Code: "TREE_LIMIT"}
	}
	if tree.Tree == nil {
		return nil, domain.InsightReadmeFailure{Code: "INTEGRITY_ERROR"}
	}
	nodes := []domain.DirectoryNode{}
	for _, item := range tree.Tree {
		if item.Type == "blob" && (item.Mode == "100644" || item.Mode == "100755") && !strings.ContainsAny(item.Path, "/\\") {
			nodes = append(nodes, domain.DirectoryNode{Path: item.Path, Type: "file", Size: item.Size})
		}
	}
	return nodes, nil
}
func (r *InsightLinkedReader) ExtractRepoFiles(ctx context.Context, repo domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	failure := func(code string) (*domain.ExtractedRepository, error) {
		return nil, domain.InsightReadmeFailure{Code: code}
	}
	if !validInsightLinkedRepo(repo) || len(paths) != 1 {
		return failure("INTEGRITY_ERROR")
	}
	path := paths[0]
	switch strings.ToLower(path) {
	case "readme.md", "readme.rst", "readme.txt", "readme":
	default:
		return failure("INTEGRITY_ERROR")
	}
	var file struct {
		Type     string `json:"type"`
		Path     string `json:"path"`
		Encoding string `json:"encoding"`
		Content  string `json:"content"`
		Size     *int   `json:"size"`
		SHA      string `json:"sha"`
	}
	if err := r.client.getInsightReferenceJSON(ctx, "/repos/"+repo.Owner+"/"+repo.Repo+"/contents/"+url.PathEscape(path)+"?ref="+repo.CommitID, 64<<10, &file); err != nil {
		return nil, err
	}
	if file.Type != "file" || file.Path != path || file.Encoding != "base64" || file.Size == nil || *file.Size < 0 || !implementationSHA(file.SHA) {
		return failure("INTEGRITY_ERROR")
	}
	if *file.Size > 32<<10 {
		return failure("CONTENT_LIMIT")
	}
	content, err := base64.StdEncoding.Strict().DecodeString(file.Content)
	if err != nil || len(content) != *file.Size {
		return failure("INTEGRITY_ERROR")
	}
	object := append([]byte(fmt.Sprintf("blob %d\x00", len(content))), content...)
	actual := ""
	if len(file.SHA) == 40 {
		actual = fmt.Sprintf("%x", sha1.Sum(object))
	} else {
		actual = fmt.Sprintf("%x", sha256.Sum256(object))
	}
	if file.SHA != actual {
		return failure("INTEGRITY_ERROR")
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}, Files: []domain.ExtractedFile{{Path: path, Content: string(content)}}}, nil
}
