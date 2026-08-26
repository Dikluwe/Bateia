/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/validation.md
 * @layer L1
 * @updated 2026-08-23
 */

package usecase

import (
	"fmt"
	"strings"

	"bateia/01_core/domain"
)

const (
	MaxDiscoveryResults = 100
	MaxSearchDepth      = 5
)

var supportedLanguageAliases = map[string]bool{
	"go": true, "golang": true, "typescript": true, "ts": true, "javascript": true, "js": true, "node": true,
	"python": true, "py": true, "rust": true, "rs": true, "java": true, "kotlin": true, "kt": true,
	"php": true, "c": true, "cpp": true, "c++": true, "csharp": true, "cs": true, "dotnet": true,
	"julia": true, "r": true, "lua": true, "ocaml": true, "elixir": true, "ruby": true, "swift": true,
	"haskell": true, "dart": true, "flutter": true, "scala": true, "clojure": true, "clj": true,
	"zig": true, "pascal": true, "delphi": true, "fortran": true, "basic": true, "vb": true,
	"fsharp": true, "fs": true, "perl": true, "groovy": true, "cobol": true, "nim": true,
	"crystal": true, "v": true, "vlang": true, "erlang": true, "d": true, "dlang": true,
	"gleam": true, "coldfusion": true, "cf": true, "actionscript": true, "flex": true, "ada": true,
	"mojo": true, "objective-c": true, "objc": true, "haxe": true, "raku": true, "racket": true,
	"chapel": true, "pony": true, "odin": true,
	"sweep": true, "all": true, "polyglot": true,
}

// ValidateQueryFilter valida entradas compartilhadas por CLI, MCP e chamadas diretas ao caso de uso.
func ValidateQueryFilter(filter domain.QueryFilter) error {
	if strings.TrimSpace(filter.Query) == "" {
		return fmt.Errorf("query não pode ser vazia")
	}
	if filter.MaxResults < 1 || filter.MaxResults > MaxDiscoveryResults {
		return fmt.Errorf("max_results deve estar entre 1 e %d", MaxDiscoveryResults)
	}
	if filter.Depth < 1 || filter.Depth > MaxSearchDepth {
		return fmt.Errorf("depth deve estar entre 1 e %d", MaxSearchDepth)
	}
	for _, language := range strings.Split(strings.ToLower(strings.TrimSpace(filter.Language)), ",") {
		language = strings.TrimSpace(language)
		if language == "" || !supportedLanguageAliases[language] {
			return fmt.Errorf("linguagem não suportada: %q", language)
		}
	}
	for _, repo := range filter.DirectRepos {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" || strings.ContainsAny(repo, " \\?#") {
			return fmt.Errorf("repositório inválido %q; use owner/repo", repo)
		}
	}
	if filter.LocalPath != "" && len(filter.DirectRepos) > 0 {
		return fmt.Errorf("local_path e direct_repos são mutuamente exclusivos")
	}
	return nil
}
