/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/forge/public_insight_reader.md
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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

// PublicInsightReader reads anonymous public repositories on a fixed service.
type PublicInsightReader struct{ http *InsightReader }

func NewCodebergInsightReader(client *http.Client) *PublicInsightReader {
	return &PublicInsightReader{http: &InsightReader{forge: newBaseForge(domain.ConnectorCodebergCatalog, "readme-v1", "https://codeberg.org/api/v1", client)}}
}
func NewGitLabInsightReader(client *http.Client) *PublicInsightReader {
	return &PublicInsightReader{http: &InsightReader{forge: newBaseForge(domain.ConnectorGitLabCatalog, "readme-v1", "https://gitlab.com/api/v4", client)}}
}
func (r *PublicInsightReader) gitlab() bool { return r.http.forge.id == domain.ConnectorGitLabCatalog }

func (r *PublicInsightReader) endpoint(identity domain.ProductIdentity) (string, error) {
	if identity.ConnectorID != r.http.forge.id || !validForgeIdentity(identity.ConnectorID, identity.ExternalID, identity.CanonicalURL) {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	parts := strings.Split(identity.ExternalID, "/")
	if !r.gitlab() && len(parts) != 2 {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	for _, part := range parts {
		if len(part) > 100 {
			return "", readmeFailure("INTEGRITY_ERROR")
		}
		for _, c := range part {
			if c > unicode.MaxASCII {
				return "", readmeFailure("INTEGRITY_ERROR")
			}
		}
	}
	u, _ := url.Parse(identity.CanonicalURL)
	if u.RawPath != "" {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	if r.gitlab() {
		return r.http.forge.baseURL + "/projects/" + url.PathEscape(identity.ExternalID), nil
	}
	return r.http.forge.baseURL + "/repos/" + identity.ExternalID, nil
}

func (r *PublicInsightReader) ResolveRepository(ctx context.Context, identity domain.ProductIdentity) (domain.RepositoryCandidate, error) {
	endpoint, err := r.endpoint(identity)
	if err != nil {
		return domain.RepositoryCandidate{}, err
	}
	var metadata struct {
		genericForgeItem
		DefaultBranch string `json:"default_branch"`
	}
	if err := r.http.get(ctx, endpoint, nil, 256<<10, &metadata); err != nil {
		return domain.RepositoryCandidate{}, err
	}
	public := flagValue(metadata.Private, false) && publicForgeItem(r.http.forge.id, metadata.genericForgeItem)
	name, rawURL := metadata.FullName, metadata.HTMLURL
	if r.gitlab() {
		public = metadata.Visibility == "public" && publicForgeItem(r.http.forge.id, metadata.genericForgeItem)
		name, rawURL = metadata.PathWithNS, metadata.WebURL
	}
	branch := metadata.DefaultBranch
	if !public || name != identity.ExternalID || rawURL != identity.CanonicalURL || !validPublicReadmeBranch(branch) {
		return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
	}
	sha := ""
	if r.gitlab() {
		var commit struct {
			ID string `json:"id"`
		}
		if err := r.http.get(ctx, endpoint+"/repository/commits/"+url.PathEscape(branch), nil, 256<<10, &commit); err != nil {
			return domain.RepositoryCandidate{}, err
		}
		sha = commit.ID
	} else {
		var branchInfo struct {
			Name   string `json:"name"`
			Commit struct {
				ID string `json:"id"`
			} `json:"commit"`
		}
		if err := r.http.get(ctx, endpoint+"/branches/"+url.PathEscape(branch), nil, 256<<10, &branchInfo); err != nil {
			return domain.RepositoryCandidate{}, err
		}
		if branchInfo.Name != branch {
			return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
		}
		sha = branchInfo.Commit.ID
	}
	if !immutableReadmeSHA(sha) {
		return domain.RepositoryCandidate{}, readmeFailure("INTEGRITY_ERROR")
	}
	split := strings.LastIndex(identity.ExternalID, "/")
	return domain.RepositoryCandidate{Owner: identity.ExternalID[:split], Repo: identity.ExternalID[split+1:], URL: identity.CanonicalURL, CommitID: sha}, nil
}

func validPublicReadmeBranch(branch string) bool {
	if branch == "" || len(branch) > 256 || strings.TrimSpace(branch) != branch {
		return false
	}
	for _, part := range strings.Split(branch, "/") {
		if part == "" || part == "." || part == ".." || strings.HasSuffix(part, ".") || strings.HasSuffix(part, ".lock") {
			return false
		}
	}
	for _, c := range branch {
		if unicode.IsControl(c) || unicode.IsSpace(c) || strings.ContainsRune("~^:?*[\\", c) {
			return false
		}
	}
	return true
}
func (r *PublicInsightReader) repositoryEndpoint(repo domain.RepositoryCandidate) (string, error) {
	if !immutableReadmeSHA(repo.CommitID) {
		return "", readmeFailure("INTEGRITY_ERROR")
	}
	return r.endpoint(domain.ProductIdentity{ConnectorID: r.http.forge.id, ExternalID: repo.Owner + "/" + repo.Repo, CanonicalURL: repo.URL})
}

type publicReadmeTreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Mode string `json:"mode"`
}

func (r *PublicInsightReader) FetchTree(ctx context.Context, repo domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	endpoint, err := r.repositoryEndpoint(repo)
	if err != nil {
		return nil, err
	}
	entries := []publicReadmeTreeEntry(nil)
	if r.gitlab() {
		if err := r.http.get(ctx, endpoint+"/repository/tree", url.Values{"ref": {repo.CommitID}, "page": {"1"}, "per_page": {"100"}}, 256<<10, &entries); err != nil {
			return nil, err
		}
	} else {
		var tree struct {
			Tree       []publicReadmeTreeEntry `json:"tree"`
			Truncated  bool                    `json:"truncated"`
			TotalCount int                     `json:"total_count"`
		}
		if err := r.http.get(ctx, endpoint+"/git/trees/"+repo.CommitID, url.Values{"page": {"1"}, "per_page": {"100"}}, 256<<10, &tree); err != nil {
			return nil, err
		}
		if tree.Truncated || tree.TotalCount > len(tree.Tree) {
			return nil, readmeFailure("TREE_LIMIT")
		}
		entries = tree.Tree
	}
	if entries == nil {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	if len(entries) >= 100 {
		return nil, readmeFailure("TREE_LIMIT")
	}
	nodes := []domain.DirectoryNode{}
	for _, entry := range entries {
		if entry.Type == "blob" && (entry.Mode == "100644" || entry.Mode == "100755") && !strings.ContainsAny(entry.Path, "/\\") {
			nodes = append(nodes, domain.DirectoryNode{Path: entry.Path, Type: "file"})
		}
	}
	return nodes, nil
}
func (r *PublicInsightReader) ExtractRepoFiles(ctx context.Context, repo domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
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
		Type          string `json:"type"`
		Encoding      string `json:"encoding"`
		Path          string `json:"path"`
		FilePath      string `json:"file_path"`
		Size          *int   `json:"size"`
		SHA           string `json:"sha"`
		BlobID        string `json:"blob_id"`
		CommitID      string `json:"commit_id"`
		ContentSHA256 string `json:"content_sha256"`
		Content       string `json:"content"`
	}
	route := "/contents/"
	if r.gitlab() {
		route = "/repository/files/"
	}
	if err := r.http.get(ctx, endpoint+route+url.PathEscape(path), url.Values{"ref": {repo.CommitID}}, 64<<10, &file); err != nil {
		return nil, err
	}
	sha := file.SHA
	if r.gitlab() {
		if file.FilePath != path || file.CommitID != repo.CommitID {
			return nil, readmeFailure("INTEGRITY_ERROR")
		}
		sha = file.BlobID
	} else if file.Path != path || file.Type != "file" {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	if file.Encoding != "base64" || file.Size == nil || *file.Size < 0 || !immutableReadmeSHA(sha) {
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
	if len(sha) == 40 {
		actual = fmt.Sprintf("%x", sha1.Sum(object))
	} else {
		actual = fmt.Sprintf("%x", sha256.Sum256(object))
	}
	if sha != actual {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	if r.gitlab() && file.ContentSHA256 != "" && file.ContentSHA256 != fmt.Sprintf("%x", sha256.Sum256(content)) {
		return nil, readmeFailure("INTEGRITY_ERROR")
	}
	return &domain.ExtractedRepository{Ref: domain.CommitRef{Owner: repo.Owner, Repo: repo.Repo, Commit: repo.CommitID}, Files: []domain.ExtractedFile{{Path: path, Content: string(content)}}}, nil
}
