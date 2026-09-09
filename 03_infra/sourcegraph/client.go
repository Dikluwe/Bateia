/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/client.md
 * @layer L3
 * @updated 2026-09-08
 */

package sourcegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"bateia/01_core/domain"
)

// Client é o adaptador para a API do Sourcegraph.
type Client struct {
	endpoint   string
	apiToken   string
	httpClient *http.Client
}

func NewClient(apiToken string) *Client {
	return &Client{
		endpoint: "https://sourcegraph.com/.api/graphql",
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// NewClientWithEndpoint cria um Client com endpoint customizado (útil para testes com httptest).
func NewClientWithEndpoint(apiToken string, endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type sgRepo struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type sgSearchResultItem struct {
	TypeName    string `json:"__typename"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	File        struct {
		Path string `json:"path"`
	} `json:"file"`
	LineMatches []struct {
		Preview string `json:"preview"`
	} `json:"lineMatches"`
	Repository sgRepo `json:"repository"`
}

type sgSearchResponse struct {
	Data struct {
		Search struct {
			Results struct {
				MatchCount *int                 `json:"matchCount"`
				Results    []sgSearchResultItem `json:"results"`
			} `json:"results"`
		} `json:"search"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// SearchCandidates executa a consulta de descoberta no Sourcegraph satisfazendo domain.SearchProvider.
// sgLangName converte o nome da linguagem do Bateia para o nome exato aceito pelo Sourcegraph.
var sgLangName = map[string]string{
	"go":         "Go",
	"typescript": "TypeScript",
	"ts":         "TypeScript",
	"javascript": "JavaScript",
	"js":         "JavaScript",
	"node":       "JavaScript",
	"python":     "Python",
	"py":         "Python",
	"rust":       "Rust",
	"rs":         "Rust",
	"java":       "Java",
	"kotlin":     "Kotlin",
	"kt":         "Kotlin",
	"php":        "PHP",
	"c":          "C",
	"cpp":        "C++",
	"c++":        "C++",
	"csharp":     "C#",
	"cs":         "C#",
	"dotnet":     "C#",
	"ruby":       "Ruby",
	"swift":      "Swift",
	"scala":      "Scala",
	"haskell":    "Haskell",
	"dart":       "Dart",
	"flutter":    "Dart",
	"elixir":     "Elixir",
	"clojure":    "Clojure",
	"clj":        "Clojure",
	"erlang":     "Erlang",
	"ocaml":      "OCaml",
	"zig":        "Zig",
	"lua":        "Lua",
	"perl":       "Perl",
	"r":          "R",
	"julia":      "Julia",
	"groovy":     "Groovy",
	"nim":        "Nim",
	"crystal":    "Crystal",
	"fsharp":     "F#",
	"fs":         "F#",
	"fortran":    "Fortran",
	"cobol":      "COBOL",
	"pascal":     "Pascal",
	"delphi":     "Delphi",
	"ada":        "Ada",
	"haxe":       "Haxe",
	"racket":     "Racket",
}

func (c *Client) SearchCandidates(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, error) {
	candidates, _, err := c.SearchCandidatesWithTrace(ctx, filter)
	return candidates, err
}

func (c *Client) SearchCandidatesWithTrace(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, []domain.SearchAttempt, error) {
	result, err := c.Discover(ctx, domain.DiscoveryRequest{Query: filter.Query, Language: filter.Language, MaxResults: filter.MaxResults})
	return result.Candidates, discoverySearchAttempts(filter.Query, result.Attempts, result.Matches), err
}

func (c *Client) executeSGQuery(ctx context.Context, queryStr string, maxResults int) ([]domain.RepositoryCandidate, error) {
	candidates, _, err := c.executeSGQueryDetailed(ctx, queryStr, maxResults)
	return candidates, err
}

func (c *Client) executeSGQueryDetailed(ctx context.Context, queryStr string, maxResults int) ([]domain.RepositoryCandidate, []domain.ProviderCandidateMatch, error) {
	sgResp, err := c.fetchSGQuery(ctx, queryStr)
	if err != nil {
		return nil, nil, err
	}

	if maxResults <= 0 {
		maxResults = 1
	}
	candidatesMap := make(map[string]domain.RepositoryCandidate)
	candidates := make([]domain.RepositoryCandidate, 0)
	matches := make([]domain.ProviderCandidateMatch, 0)
	matchIndexes := make(map[string]int)

	for _, res := range sgResp.Data.Search.Results.Results {
		var repo sgRepo
		if res.TypeName == "Repository" && res.Name != "" {
			repo = sgRepo{
				Name:        res.Name,
				URL:         res.URL,
				Description: res.Description,
			}
		} else if res.TypeName == "FileMatch" && res.Repository.Name != "" {
			repo = res.Repository
		}

		if repo.Name == "" {
			continue
		}

		owner, name := parseRepoName(repo.Name)
		key := strings.ToLower(owner + "/" + name)
		// Keep evidence for already retained repositories even after reaching the cap.
		if _, exists := candidatesMap[repo.Name]; !exists && len(candidates) >= maxResults {
			continue
		}
		evidences := []domain.MatchEvidence{{MatchType: "repository"}}
		if res.TypeName == "FileMatch" {
			evidences = []domain.MatchEvidence{{MatchType: "file", Path: res.File.Path}}
			if len(res.LineMatches) > 0 {
				evidences = nil
				for _, line := range res.LineMatches {
					evidences = append(evidences, domain.MatchEvidence{MatchType: "file", Path: res.File.Path, Snippet: truncateMatchSnippet(line.Preview)})
				}
			}
		}
		if _, exists := matchIndexes[key]; !exists {
			matchIndexes[key] = len(matches)
			matches = append(matches, domain.ProviderCandidateMatch{Repository: owner + "/" + name})
		}
		matchIndex := matchIndexes[key]
		for _, evidence := range evidences {
			matches[matchIndex].Evidence = appendUniqueMatchEvidence(matches[matchIndex].Evidence, evidence)
		}
		if _, exists := candidatesMap[repo.Name]; exists {
			continue
		}

		cand := domain.RepositoryCandidate{
			Owner:        owner,
			Repo:         name,
			URL:          "https://" + repo.Name,
			Description:  repo.Description,
			Stars:        0,
			CommitID:     "HEAD",
			ConsultadoEm: time.Now().UTC(),
			Stale:        false,
		}

		candidatesMap[repo.Name] = cand
		candidates = append(candidates, cand)

	}

	return candidates, matches, nil
}

func truncateMatchSnippet(value string) string {
	normalized := strings.Join(strings.Fields(value), " ")
	runes := []rune(normalized)
	if len(runes) > 240 {
		return string(runes[:240])
	}
	return normalized
}

func appendUniqueMatchEvidence(values []domain.MatchEvidence, addition domain.MatchEvidence) []domain.MatchEvidence {
	for _, value := range values {
		if value.MatchType == addition.MatchType && value.Path == addition.Path && value.Snippet == addition.Snippet {
			return values
		}
	}
	return append(values, addition)
}

func parseRepoName(fullName string) (string, string) {
	var owner, repo string
	var parts []string
	for _, p := range bytes.Split([]byte(fullName), []byte("/")) {
		parts = append(parts, string(p))
	}

	if len(parts) >= 3 {
		owner = parts[1]
		repo = parts[2]
	} else if len(parts) == 2 {
		owner = parts[0]
		repo = parts[1]
	} else {
		repo = fullName
	}

	return owner, repo
}

func (c *Client) fetchSGQuery(ctx context.Context, queryStr string) (*sgSearchResponse, error) {
	gqlQuery := `
	query SearchCandidates($query: String!) {
		search(query: $query, version: V2) {
			results {
				matchCount
				results {
					__typename
					... on Repository {
						name
						url
						description
					}
					... on FileMatch {
						file { path }
						lineMatches { preview }
						repository {
							name
							url
							description
						}
					}
				}
			}
		}
	}`

	reqBody := graphQLRequest{
		Query: gqlQuery,
		Variables: map[string]interface{}{
			"query": queryStr,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiToken != "" {
		req.Header.Set("Authorization", "token "+c.apiToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro de conexão com Sourcegraph: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status inesperado do Sourcegraph: %d: %s", resp.StatusCode, string(respData))
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sgResp sgSearchResponse
	if err := json.Unmarshal(respData, &sgResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do Sourcegraph: %w", err)
	}

	if len(sgResp.Errors) > 0 {
		return nil, fmt.Errorf("erro na API do Sourcegraph: %s", sgResp.Errors[0].Message)
	}

	// Missing or null envelopes are malformed, not valid empty searches.
	var envelope struct {
		Data *struct {
			Search *struct {
				Results *struct {
					Results json.RawMessage `json:"results"`
				} `json:"results"`
			} `json:"search"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respData, &envelope); err != nil || envelope.Data == nil || envelope.Data.Search == nil || envelope.Data.Search.Results == nil {
		return nil, fmt.Errorf("resposta Sourcegraph sem envelope de resultados válido")
	}
	rawResults := bytes.TrimSpace(envelope.Data.Search.Results.Results)
	if len(rawResults) == 0 || rawResults[0] != '[' {
		return nil, fmt.Errorf("resposta Sourcegraph sem array de resultados válido")
	}

	return &sgResp, nil
}
