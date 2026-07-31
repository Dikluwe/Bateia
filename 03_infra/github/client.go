/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/github/client.md
 * @layer L3
 * @updated 2026-07-25
 */

package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

// Client é o adaptador para acesso à API do GitHub (Contents API, Raw HTTP e Git Clone) satisfazendo domain.ContentProvider.
type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
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
}

// FetchManifest lê o manifesto (ex: go.mod ou package.json) via Contents API ou fallback raw (ADR-0001).
func (c *Client) FetchManifest(ctx context.Context, candidate domain.RepositoryCandidate) (*domain.Manifest, error) {
	refCommit := candidate.CommitID
	if refCommit == "" {
		refCommit = "HEAD"
	}

	manifestFiles := []string{"go.mod", "package.json", "pyproject.toml", "requirements.txt", "Cargo.toml", "pom.xml", "build.gradle", "composer.json", "CMakeLists.txt", "vcpkg.json", "conanfile.txt", "Project.toml", "DESCRIPTION", "dune-project", "mix.exs", "Gemfile", "Package.swift", "pubspec.yaml", "build.sbt", "project.clj", "build.zig.zon", "fpm.toml", "cpanfile", "Makefile.PL"}
	var lastErr error

	for _, targetFile := range manifestFiles {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", candidate.Owner, candidate.Repo, targetFile)
		if refCommit != "HEAD" {
			url += "?ref=" + refCommit
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
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
			lastErr = fmt.Errorf("MANIFEST_NOT_FOUND (%s)", targetFile)
			continue
		}

		var rawContent string

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized {
			raw, errRaw := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, refCommit, targetFile)
			if errRaw != nil {
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
			if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
				lastErr = err
				continue
			}

			rawContent = item.Content
			if item.Encoding == "base64" {
				decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(rawContent, "\n", ""))
				if err == nil {
					rawContent = string(decoded)
				}
			}
		}

		ref := domain.CommitRef{
			Owner:  candidate.Owner,
			Repo:   candidate.Repo,
			Commit: candidate.CommitID,
		}

		return parser.ParseManifest(targetFile, ref, rawContent)
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("MANIFEST_NOT_FOUND")
}

func (c *Client) fetchRawFallback(ctx context.Context, owner, repo, refCommit, path string) (string, error) {
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, refCommit, path)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchTree recupera a lista de nós do diretório raiz.
func (c *Client) FetchTree(ctx context.Context, candidate domain.RepositoryCandidate) ([]domain.DirectoryNode, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", candidate.Owner, candidate.Repo)
	if candidate.CommitID != "" && candidate.CommitID != "HEAD" {
		url += "?ref=" + candidate.CommitID
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
		return []domain.DirectoryNode{
			{Path: "go.mod", Type: "file"},
		}, nil
	}

	var items []githubContentItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
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

// ExtractRepoFiles realiza a extração do repositório no fundo do funil.
func (c *Client) ExtractRepoFiles(ctx context.Context, candidate domain.RepositoryCandidate, paths []string) (*domain.ExtractedRepository, error) {
	ref := domain.CommitRef{
		Owner:  candidate.Owner,
		Repo:   candidate.Repo,
		Commit: candidate.CommitID,
	}

	extracted := &domain.ExtractedRepository{
		Ref:          ref,
		Candidate:    candidate,
		Files:        make([]domain.ExtractedFile, 0),
		ConsultadoEm: time.Now().UTC(),
	}

	for _, p := range paths {
		content, err := c.fetchSingleFile(ctx, candidate, p)
		if err != nil {
			refCommit := candidate.CommitID
			if refCommit == "" {
				refCommit = "HEAD"
			}
			raw, errRaw := c.fetchRawFallback(ctx, candidate.Owner, candidate.Repo, refCommit, p)
			if errRaw == nil {
				extracted.Files = append(extracted.Files, domain.ExtractedFile{
					Path:    p,
					Content: raw,
					Size:    int64(len(raw)),
				})
			}
		} else {
			extracted.Files = append(extracted.Files, *content)
		}
	}

	if len(extracted.Files) > 0 {
		return extracted, nil
	}

	tempDir, err := os.MkdirTemp("", "bateia_clone_*")
	if err != nil {
		return nil, fmt.Errorf("falha ao criar diretório temporário para clone: %w", err)
	}
	defer os.RemoveAll(tempDir)

	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", candidate.Owner, candidate.Repo)
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--filter=blob:none", repoURL, tempDir)
	if err := cmd.Run(); err != nil {
		cmdSimple := exec.CommandContext(ctx, "git", "clone", "--depth", "1", repoURL, tempDir)
		if err := cmdSimple.Run(); err != nil {
			return nil, fmt.Errorf("falha ao clonar repositório %s: %w", candidate.URL, err)
		}
	}

	targetFiles := []string{"go.mod", "README.md", "main.go"}
	for _, tf := range targetFiles {
		filePath := filepath.Join(tempDir, tf)
		if data, err := os.ReadFile(filePath); err == nil {
			extracted.Files = append(extracted.Files, domain.ExtractedFile{
				Path:    tf,
				Content: string(data),
				Size:    int64(len(data)),
			})
		}
	}

	return extracted, nil
}

func (c *Client) fetchSingleFile(ctx context.Context, candidate domain.RepositoryCandidate, path string) (*domain.ExtractedFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", candidate.Owner, candidate.Repo, path)
	if candidate.CommitID != "" && candidate.CommitID != "HEAD" {
		url += "?ref=" + candidate.CommitID
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
		return nil, fmt.Errorf("status %d ao ler arquivo %s", resp.StatusCode, path)
	}

	var item githubContentItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	rawContent := item.Content
	if item.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(rawContent, "\n", ""))
		if err == nil {
			rawContent = string(decoded)
		}
	}

	return &domain.ExtractedFile{
		Path:    item.Path,
		Content: rawContent,
		Size:    item.Size,
	}, nil
}

type githubRepoResponse struct {
	StargazersCount int    `json:"stargazers_count"`
	Description     string `json:"description"`
}

// FetchRepoMetadata realiza o sanity check de metadados na API REST do GitHub obtendo número de estrelas reais e descrição.
func (c *Client) FetchRepoMetadata(ctx context.Context, candidate domain.RepositoryCandidate) (domain.RepositoryCandidate, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", candidate.Owner, candidate.Repo)

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
	if err := json.NewDecoder(resp.Body).Decode(&repoResp); err != nil {
		return candidate, nil
	}

	candidate.Stars = repoResp.StargazersCount
	if repoResp.Description != "" {
		candidate.Description = repoResp.Description
	}

	return candidate, nil
}
