/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/client.md
 * @layer L3
 * @updated 2026-10-28
 */

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// Client acessa Contents API e Raw HTTP com limites; nunca clona ou executa código remoto.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: "https://api.github.com",
		httpClient: &http.Client{
			Timeout:       20 * time.Second,
			CheckRedirect: denyRedirect,
		},
	}
}

// NewClientWithBaseURL cria um Client com base URL customizada (útil para testes com httptest).
func NewClientWithBaseURL(token string, baseURL string) *Client {
	return &Client{
		token:   token,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:       10 * time.Second,
			CheckRedirect: denyRedirect,
		},
	}
}

// RateLimitInfo guarda o estado atual da quota da API do GitHub.
type RateLimitInfo struct {
	Remaining int
	ResetAt   time.Time
}

func (c *Client) parseRateLimit(resp *http.Header) RateLimitInfo {
	info := RateLimitInfo{}
	if rem := resp.Get("X-RateLimit-Remaining"); rem != "" {
		info.Remaining, _ = strconv.Atoi(rem)
	}
	if reset := resp.Get("X-RateLimit-Reset"); reset != "" {
		sec, _ := strconv.ParseInt(reset, 10, 64)
		info.ResetAt = time.Unix(sec, 0)
	}
	return info
}

type githubContentItem struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	DownloadURL string `json:"download_url"`
	Target      string `json:"target"`
}

type githubTreeResponse struct {
	Truncated bool `json:"truncated"`
	Tree      []struct {
		Path string `json:"path"`
		Type string `json:"type"`
		Mode string `json:"mode"`
		Size int64  `json:"size"`
	} `json:"tree"`
}

// FetchManifest lê o manifesto (ex: go.mod ou package.json) via Contents API ou fallback raw (ADR-0001).
func (c *Client) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return nil, err
	}
	refCommit := candidate.CommitID
	if refCommit == "" {
		refCommit = "HEAD"
	}

	manifestFiles := []string{"go.mod", "package.json", "pyproject.toml", "requirements.txt", "Cargo.toml", "pom.xml", "build.gradle", "composer.json", "CMakeLists.txt", "vcpkg.json", "conanfile.txt", "Project.toml", "DESCRIPTION", "dune-project", "mix.exs", "Gemfile", "Package.swift", "pubspec.yaml", "build.sbt", "project.clj", "build.zig.zon", "fpm.toml", "cpanfile", "Makefile.PL", "shard.yml", "v.mod", "rebar.config", "dub.json", "gleam.toml", "box.json", "flex-config.xml", "alire.toml", "mojoproject.toml", "Podfile", "haxelib.json", "META6.json", "info.rkt", "mason.toml", "corral.json", "ols.json"}
	if rootItems, err := c.fetchRootItems(ctx, candidate); err == nil {
		manifestFiles = selectManifestFiles(rootItems, manifestFiles)
		if len(manifestFiles) == 0 {
			return nil, domain.ErrManifestNotFound
		}
	} else if errors.Is(err, errUnsafeRemote) {
		return nil, err
	}
	if len(manifestFiles) > maxExtractionFiles {
		return nil, unsafeRemote("quantidade de manifestos excede orçamento de probes")
	}
	var lastErr error

	for _, targetFile := range manifestFiles {
		escaped, err := remotePath(targetFile)
		if err != nil {
			return nil, err
		}
		endpoint := fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.baseURL, candidate.Owner, candidate.Repo, escaped)
		if refCommit != "HEAD" {
			endpoint += "?ref=" + neturl.QueryEscape(refCommit)
		}

		req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if c.token != "" {
			req.Header.Set("Authorization", "token "+c.token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, readErr := readBounded(resp, maxContentEnvelopeBytes)
		resp.Body.Close()
		if readErr != nil {
			if errors.Is(readErr, errUnsafeRemote) {
				return nil, readErr
			}
			lastErr = readErr
			continue
		}

		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
			lastErr = fmt.Errorf("%w (%s)", domain.ErrManifestNotFound, targetFile)
			continue
		}

		var rawContent string

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized {
			raw, errRaw := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, refCommit, targetFile)
			if errRaw != nil {
				if errors.Is(errRaw, errUnsafeRemote) {
					return nil, errRaw
				}
				rateInfo := c.parseRateLimit(&resp.Header)
				lastErr = fmt.Errorf("cota de API atingida (%d restantes, reseta %v) e fallback raw falhou: %w", rateInfo.Remaining, rateInfo.ResetAt, errRaw)
				continue
			}
			rawContent = raw
		} else if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status inesperado no GitHub Contents API para %s: %d", targetFile, resp.StatusCode)
			continue
		} else {
			var item githubContentItem
			if err := json.Unmarshal(body, &item); err != nil {
				return nil, unsafeRemote("JSON de conteúdo inválido: " + err.Error())
			}

			rawContent, err = decodeRegularContent(item, targetFile)
			if err != nil {
				return nil, err
			}
		}

		ref := domain.CommitRef{
			Owner:  candidate.Owner,
			Repo:   candidate.Repo,
			Commit: candidate.CommitID,
		}

		manifest, err := parser.ParseManifestAt(targetFile, ref, rawContent, time.Now().UTC())
		if err != nil {
			lastErr = err
			continue
		}

		// Enriquece manifesto com conteúdo do README para README Score Booster
		if manifest.Content == nil {
			readmeContent, _ := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, refCommit, "README.md")
			if readmeContent != "" {
				manifest.Content = []byte(readmeContent)
				manifest.ContentPath = "README.md"
			}
		}

		// Se o manifesto não tem dependências, tenta go.sum como fallback (stdlib-only repos)
		if len(manifest.Dependencies) == 0 && targetFile == "go.mod" {
			goSumContent, err := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, refCommit, "go.sum")
			if err == nil && goSumContent != "" {
				sumDeps := parseGoSumDependencies(goSumContent)
				manifest.Dependencies = sumDeps
			}
		}

		return manifest, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, domain.ErrManifestNotFound
}

func selectManifestFiles(items []githubContentItem, priority []string) []string {
	available := make(map[string]string)
	var additional []string
	for _, item := range items {
		if item.Type != "file" || parser.DetectManifestType(item.Path) == "" {
			continue
		}
		available[strings.ToLower(item.Name)] = item.Path
	}
	selected := make([]string, 0, len(available))
	selectedSet := make(map[string]bool)
	for _, name := range priority {
		if path, ok := available[strings.ToLower(name)]; ok {
			selected = append(selected, path)
			selectedSet[path] = true
		}
	}
	for _, path := range available {
		if !selectedSet[path] {
			additional = append(additional, path)
		}
	}
	sort.Strings(additional)
	return append(selected, additional...)
}

func (c *Client) fetchRootItems(ctx context.Context, candidate domain.RepositoryCandidate) ([]githubContentItem, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/repos/%s/%s/contents", c.baseURL, candidate.Owner, candidate.Repo)
	if candidate.CommitID != "" && candidate.CommitID != "HEAD" {
		url += "?ref=" + queryRef(candidate.CommitID)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d ao listar raiz", resp.StatusCode)
	}
	var items []githubContentItem
	if err := decodeBounded(resp, maxTreeEnvelopeBytes, &items); err != nil {
		return nil, err
	}
	if len(items) > 20_000 {
		return nil, unsafeRemote("raiz excede limite de nós")
	}
	for _, item := range items {
		if !implementationSafePath(item.Path) || item.Size < 0 {
			return nil, unsafeRemote("entrada de árvore inválida")
		}
	}
	return items, nil
}

// parseGoSumDependencies extrai dependências do go.sum quando o go.mod não as lista.
func parseGoSumDependencies(content string) []domain.Dependency {
	seen := make(map[string]bool)
	var deps []domain.Dependency
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		// Formato: module version/go.mod hash
		modParts := strings.SplitN(parts[0], "@", 2)
		if len(modParts) < 2 {
			continue
		}
		modName := modParts[0]
		modVersion := strings.TrimSuffix(modParts[1], "/go.mod")
		if seen[modName] {
			continue
		}
		seen[modName] = true
		deps = append(deps, domain.Dependency{
			Name:     modName,
			Version:  modVersion,
			Indirect: true, // go.sum não distingue direto/indireto
		})
	}
	return deps
}

func (c *Client) fetchRawFallback(ctx context.Context, owner, repo, refCommit, path string) (string, error) {
	if err := validateRemoteCandidate(domain.RepositoryCandidate{Owner: owner, Repo: repo, CommitID: refCommit}); err != nil {
		return "", err
	}
	escaped, err := remotePath(path)
	if err != nil {
		return "", err
	}
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, neturl.PathEscape(refCommit), escaped)
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status raw: %d", resp.StatusCode)
	}

	body, err := readBounded(resp, maxFileBytes)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchTree recupera a lista de nós do diretório raiz.
func (c *Client) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return nil, err
	}
	if candidate.CommitID != "" && !strings.EqualFold(candidate.CommitID, "HEAD") {
		url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", c.baseURL, candidate.Owner, candidate.Repo, neturl.PathEscape(candidate.CommitID))
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if c.token != "" {
			req.Header.Set("Authorization", "token "+c.token)
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d ao listar árvore no commit", resp.StatusCode)
		}
		var tree githubTreeResponse
		if err := decodeBounded(resp, maxTreeEnvelopeBytes, &tree); err != nil {
			return nil, err
		}
		if tree.Truncated {
			return nil, fmt.Errorf("árvore GitHub truncada; coleta não conclusiva")
		}
		nodes := make([]domain.DirectoryNode, 0, len(tree.Tree))
		for _, item := range tree.Tree {
			if len(nodes) >= 20_000 {
				return nil, fmt.Errorf("árvore excede limite seguro de 20000 nós")
			}
			if !implementationSafePath(item.Path) || item.Size < 0 {
				return nil, unsafeRemote("entrada de árvore inválida")
			}
			typeName := "file"
			if item.Type == "tree" {
				typeName = "dir"
			} else if item.Type == "commit" || item.Mode == "160000" {
				typeName = "submodule"
			} else if item.Mode == "120000" {
				typeName = "symlink"
			} else if item.Type != "blob" || item.Mode != "" && item.Mode != "100644" && item.Mode != "100755" {
				return nil, unsafeRemote("tipo de árvore não suportado")
			}
			nodes = append(nodes, domain.DirectoryNode{Path: item.Path, Type: typeName, Size: item.Size})
		}
		return nodes, nil
	}
	items, err := c.fetchRootItems(ctx, candidate)
	if err != nil {
		return nil, err
	}

	nodes := make([]domain.DirectoryNode, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, domain.DirectoryNode{
			Path: item.Path,
			Type: item.Type,
			Size: item.Size,
		})
	}

	return nodes, nil
}

// ExtractRepoFiles reads only regular Git files. Clone fallback is deliberately disabled:
// HTTP failure must not become an unbounded clone, checkout, hook or filesystem read.
func (c *Client) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return nil, err
	}
	if len(paths) > maxExtractionFiles {
		return nil, unsafeRemote("extração excede limite de arquivos")
	}
	seen := map[string]bool{}
	for _, p := range paths {
		if _, err := remotePath(p); err != nil {
			return nil, err
		}
		if seen[p] {
			return nil, unsafeRemote("caminho duplicado")
		}
		seen[p] = true
	}
	extracted := &domain.ExtractedRepository{
		Ref:       domain.CommitRef{Owner: candidate.Owner, Repo: candidate.Repo, Commit: candidate.CommitID},
		Candidate: candidate, Files: []domain.ExtractedFile{}, ConsultadoEm: time.Now().UTC(),
	}
	if len(paths) == 0 {
		return extracted, nil
	}
	if err := c.requireRegularRemotePaths(ctx, candidate, paths); err != nil {
		return nil, err
	}
	var total int64
	for _, p := range paths {
		if total >= maxExtractionBytes {
			return nil, unsafeRemote("extração esgotou orçamento total de bytes")
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		content, err := c.fetchSingleFile(ctx, candidate, p)
		if err != nil {
			if errors.Is(err, errUnsafeRemote) || ctx.Err() != nil {
				return nil, err
			}
			revision := candidate.CommitID
			if revision == "" {
				revision = "HEAD"
			}
			raw, rawErr := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, revision, p)
			if rawErr != nil {
				return nil, fmt.Errorf("Contents API e raw indisponíveis para %s (clone desabilitado): %w", p, rawErr)
			}
			content = &domain.ExtractedFile{Path: p, Content: raw, Size: int64(len(raw))}
		}
		total += int64(len(content.Content))
		if total > maxExtractionBytes {
			return nil, unsafeRemote("extração excede orçamento total de bytes")
		}
		extracted.Files = append(extracted.Files, *content)
	}
	return extracted, nil
}

func (c *Client) fetchSingleFile(ctx context.Context, candidate domain.RepositoryCandidate, path string) (*domain.ExtractedFile, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return nil, err
	}
	escaped, err := remotePath(path)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.baseURL, candidate.Owner, candidate.Repo, escaped)
	if candidate.CommitID != "" && candidate.CommitID != "HEAD" {
		endpoint += "?ref=" + neturl.QueryEscape(candidate.CommitID)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d ao ler arquivo %s", resp.StatusCode, path)
	}

	var item githubContentItem
	if err := decodeBounded(resp, maxContentEnvelopeBytes, &item); err != nil {
		return nil, err
	}

	rawContent, err := decodeRegularContent(item, path)
	if err != nil {
		return nil, err
	}

	return &domain.ExtractedFile{
		Path:    item.Path,
		Content: rawContent,
		Size:    int64(len(rawContent)),
	}, nil
}

type githubRepoResponse struct {
	StargazersCount int       `json:"stargazers_count"`
	Description     string    `json:"description"`
	DefaultBranch   string    `json:"default_branch"`
	Archived        bool      `json:"archived"`
	PushedAt        time.Time `json:"pushed_at"`
}

type githubCommitResponse struct {
	SHA string `json:"sha"`
}

// FetchRepoMetadata realiza o sanity check de metadados na API REST do GitHub obtendo número de estrelas reais e descrição.
func (c *Client) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	if err := validateRemoteCandidate(candidate); err != nil {
		return candidate, err
	}
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, candidate.Owner, candidate.Repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return candidate, nil
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return candidate, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return candidate, nil
	}

	var repoResp githubRepoResponse
	if err := decodeBounded(resp, maxMetadataEnvelopeBytes, &repoResp); err != nil {
		return candidate, err
	}

	candidate.Stars = repoResp.StargazersCount
	candidate.Archived = repoResp.Archived
	candidate.DefaultBranch = repoResp.DefaultBranch
	candidate.UpdatedAt = repoResp.PushedAt
	candidate.ConsultadoEm = time.Now().UTC()
	if repoResp.Description != "" {
		candidate.Description = repoResp.Description
	}
	if candidate.CommitID == "" || strings.EqualFold(candidate.CommitID, "HEAD") {
		branch := repoResp.DefaultBranch
		if branch == "" {
			branch = "HEAD"
		}
		branchCandidate := candidate
		branchCandidate.CommitID = branch
		if err := validateRemoteCandidate(branchCandidate); err != nil {
			return candidate, err
		}
		commitURL := fmt.Sprintf("%s/repos/%s/%s/commits/%s", c.baseURL, candidate.Owner, candidate.Repo, neturl.PathEscape(branch))
		commitReq, requestErr := http.NewRequestWithContext(ctx, "GET", commitURL, nil)
		if requestErr == nil {
			commitReq.Header.Set("Accept", "application/vnd.github.v3+json")
			if c.token != "" {
				commitReq.Header.Set("Authorization", "token "+c.token)
			}
			if commitResp, requestErr := c.httpClient.Do(commitReq); requestErr == nil {
				defer commitResp.Body.Close()
				if commitResp.StatusCode == http.StatusOK {
					var commit githubCommitResponse
					if err := decodeBounded(commitResp, maxMetadataEnvelopeBytes, &commit); err != nil {
						return candidate, err
					}
					if commit.SHA != "" {
						if !implementationSHA(commit.SHA) {
							return candidate, unsafeRemote("API de commits não retornou SHA completo")
						}
						pinned := candidate
						pinned.CommitID = commit.SHA
						if err := validateRemoteCandidate(pinned); err != nil {
							return candidate, err
						}
						candidate.CommitID = commit.SHA
					}
				}
			}
		}
	}

	return candidate, nil
}
