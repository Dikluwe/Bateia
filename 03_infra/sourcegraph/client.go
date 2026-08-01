/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/sourcegraph/client.md
 * @layer L3
 * @updated 2026-07-25
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
	Repository  sgRepo `json:"repository"`
}

type sgSearchResponse struct {
	Data struct {
		Search struct {
			Results struct {
				MatchCount int                  `json:"matchCount"`
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
	"go":           "Go",
	"typescript":   "TypeScript",
	"ts":           "TypeScript",
	"javascript":   "JavaScript",
	"js":           "JavaScript",
	"node":         "JavaScript",
	"python":       "Python",
	"py":           "Python",
	"rust":         "Rust",
	"rs":           "Rust",
	"java":         "Java",
	"kotlin":       "Kotlin",
	"kt":           "Kotlin",
	"php":          "PHP",
	"c":            "C",
	"cpp":          "C++",
	"c++":          "C++",
	"csharp":       "C#",
	"cs":           "C#",
	"dotnet":       "C#",
	"ruby":         "Ruby",
	"swift":        "Swift",
	"scala":        "Scala",
	"haskell":      "Haskell",
	"dart":         "Dart",
	"flutter":      "Dart",
	"elixir":       "Elixir",
	"clojure":      "Clojure",
	"clj":          "Clojure",
	"erlang":       "Erlang",
	"ocaml":        "OCaml",
	"zig":          "Zig",
	"lua":          "Lua",
	"perl":         "Perl",
	"r":            "R",
	"julia":        "Julia",
	"groovy":       "Groovy",
	"nim":          "Nim",
	"crystal":      "Crystal",
	"fsharp":       "F#",
	"fs":           "F#",
	"fortran":      "Fortran",
	"cobol":        "COBOL",
	"pascal":       "Pascal",
	"delphi":       "Delphi",
	"ada":          "Ada",
	"haxe":         "Haxe",
	"racket":       "Racket",
}

// shortQuery extrai o(s) primeiro(s) termo(s) relevantes de uma query longa para o fallback sem lang.
func shortQuery(query string) string {
	words := strings.Fields(query)
	stopwords := map[string]bool{
		"a": true, "o": true, "e": true, "de": true, "do": true, "da": true,
		"the": true, "an": true, "and": true, "or": true, "with": true,
		"using": true, "for": true, "in": true, "on": true, "at": true, "to": true,
		"of": true, "is": true, "are": true, "api": true, "rest": true,
	}
	result := make([]string, 0, 3)
	for _, w := range words {
		if !stopwords[strings.ToLower(w)] {
			result = append(result, w)
		}
		if len(result) >= 2 {
			break
		}
	}
	if len(result) == 0 {
		return words[0]
	}
	return strings.Join(result, " ")
}

func (c *Client) SearchCandidates(ctx context.Context, filter domain.QueryFilter) ([]domain.RepositoryCandidate, error) {
	// Traduzir nome da linguagem para o formato aceito pelo Sourcegraph
	sgLang := filter.Language
	if mapped, ok := sgLangName[strings.ToLower(filter.Language)]; ok {
		sgLang = mapped
	} else if len(filter.Language) > 0 {
		sgLang = strings.ToUpper(filter.Language[:1]) + filter.Language[1:]
	}

	short := shortQuery(filter.Query)

	queriesToTry := []string{
		// 1ª tentativa: query completa com linguagem (sem repo: para API pública funcionar)
		fmt.Sprintf("lang:%s %s", sgLang, filter.Query),
		// 2ª tentativa: termos curtos com linguagem
		fmt.Sprintf("lang:%s %s", sgLang, short),
		// 3ª tentativa: sem filtro de linguagem, apenas termos curtos
		short,
	}

	for _, qStr := range queriesToTry {
		candidates, err := c.executeSGQuery(ctx, qStr, filter.MaxResults)
		if err == nil && len(candidates) > 0 {
			return candidates, nil
		}
	}

	return []domain.RepositoryCandidate{}, nil
}

func (c *Client) executeSGQuery(ctx context.Context, queryStr string, maxResults int) ([]domain.RepositoryCandidate, error) {
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
		return nil, fmt.Errorf("Sourcegraph retornou status %d: %s", resp.StatusCode, string(respData))
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

	candidatesMap := make(map[string]domain.RepositoryCandidate)
	candidates := make([]domain.RepositoryCandidate, 0)

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

		if _, exists := candidatesMap[repo.Name]; exists {
			continue
		}

		owner, name := parseRepoName(repo.Name)

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

		if len(candidates) >= maxResults {
			break
		}
	}

	return candidates, nil
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
