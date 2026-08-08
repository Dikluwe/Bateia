/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/polyglot_sweep.md
 * @layer L1
 * @updated 2026-08-08
 */

package usecase

import (
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// DefaultPolyglotSuite lista completa de todos os ecossistemas suportados pelo Bateia.
var DefaultPolyglotSuite = []string{
	"go", "typescript", "python", "rust", "java", "kotlin", "php", "c", "cpp", "csharp",
	"julia", "r", "lua", "ocaml", "elixir", "ruby", "swift", "haskell", "dart", "scala",
	"clojure", "zig", "pascal", "fortran", "basic", "fsharp", "perl", "groovy", "cobol",
	"nim", "crystal", "v", "erlang", "d", "gleam", "coldfusion", "actionscript", "ada",
	"mojo", "objc", "haxe", "raku", "racket", "chapel", "pony", "odin",
}

// ParseLanguageList analisa a string de linguagem da query/CLI e retorna uma lista desmembrada.
func ParseLanguageList(langInput string) []string {
	trimmed := strings.TrimSpace(strings.ToLower(langInput))
	if trimmed == "" || trimmed == "sweep" || trimmed == "all" || trimmed == "polyglot" {
		return DefaultPolyglotSuite
	}

	parts := strings.Split(trimmed, ",")
	var result []string
	seen := make(map[string]bool)

	for _, p := range parts {
		item := strings.TrimSpace(p)
		if item != "" && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return DefaultPolyglotSuite
	}
	return result
}

// FilterLanguages remove as linguagens especificadas na string de exclusao (ex: "go,typescript").
func FilterLanguages(languages []string, excludeInput string) []string {
	if strings.TrimSpace(excludeInput) == "" {
		return languages
	}

	excluded := make(map[string]bool)
	for _, item := range strings.Split(excludeInput, ",") {
		e := strings.TrimSpace(strings.ToLower(item))
		if e != "" {
			excluded[e] = true
			switch e {
			case "go":
				excluded["golang"] = true
			case "typescript":
				excluded["ts"] = true
			case "javascript":
				excluded["js"] = true
				excluded["node"] = true
			case "python":
				excluded["py"] = true
			case "rust":
				excluded["rs"] = true
			case "csharp":
				excluded["cs"] = true
				excluded["dotnet"] = true
			case "cpp", "c++":
				excluded["c"] = true
				excluded["cpp"] = true
				excluded["c++"] = true
			}
		}
	}

	var result []string
	for _, l := range languages {
		lLower := strings.ToLower(l)
		if !excluded[lLower] {
			result = append(result, l)
		}
	}

	if len(result) == 0 {
		return DefaultPolyglotSuite
	}
	return result
}


// CalculateLanguageDistribution compila a distribuição estatística de repositórios por ecossistema.
func CalculateLanguageDistribution(results []domain.ExtractedRepository, simResults []domain.IdeaSimilarityResult) []domain.LanguageDistribution {
	if len(results) == 0 {
		return nil
	}

	simMap := make(map[string]int)
	for _, sim := range simResults {
		simMap[sim.CandidateRef.Key()] = sim.OverallScore
	}

	type stats struct {
		lang          string
		count         int
		totalStars    int
		topSimilarity int
		packagesMap   map[string]int
	}

	langStats := make(map[string]*stats)
	totalResults := float64(len(results))

	for _, repo := range results {
		lang := "desconhecida"
		if repo.Manifest != nil && repo.Manifest.Language != "" {
			lang = strings.ToLower(repo.Manifest.Language)
		} else if repo.Candidate.Description != "" {
			lang = inferLanguageFromRepo(repo)
		}

		if _, exists := langStats[lang]; !exists {
			langStats[lang] = &stats{
				lang:        lang,
				packagesMap: make(map[string]int),
			}
		}

		st := langStats[lang]
		st.count++
		st.totalStars += repo.Candidate.Stars

		if score, ok := simMap[repo.Ref.Key()]; ok {
			if score > st.topSimilarity {
				st.topSimilarity = score
			}
		}

		if repo.Manifest != nil {
			for _, dep := range repo.Manifest.Dependencies {
				st.packagesMap[dep.Name]++
			}
		}
	}

	var distributions []domain.LanguageDistribution

	langIcons := map[string]string{
		"go": "Go 🐹", "typescript": "TypeScript 🟦", "ts": "TypeScript 🟦",
		"python": "Python 🐍", "py": "Python 🐍", "rust": "Rust 🦀", "rs": "Rust 🦀",
		"java": "Java ☕", "kotlin": "Kotlin ☕", "kt": "Kotlin ☕", "php": "PHP 🐘",
		"c": "C ⚙️", "cpp": "C++ ⚙️", "c++": "C++ ⚙️", "csharp": "C# 💜", "cs": "C# 💜", "dotnet": "C# 💜",
		"julia": "Julia 📊", "r": "R 📊", "lua": "Lua 🌙", "ocaml": "OCaml 🐫",
		"elixir": "Elixir 💧", "ruby": "Ruby 💎", "swift": "Swift 🐦", "haskell": "Haskell 📜",
		"dart": "Dart / Flutter 🎯", "flutter": "Dart / Flutter 🎯", "scala": "Scala 🔴",
		"clojure": "Clojure 🍃", "clj": "Clojure 🍃", "zig": "Zig ⚡", "pascal": "Pascal / Delphi 🏛️",
		"delphi": "Pascal / Delphi 🏛️", "fortran": "Fortran 🧮", "basic": "BASIC / VB 💻",
		"vb": "BASIC / VB 💻", "fsharp": "F# 🔷", "fs": "F# 🔷", "perl": "Perl 🐪",
		"groovy": "Groovy 🍃", "cobol": "COBOL 📜", "nim": "Nim 👑", "crystal": "Crystal 💎",
		"v": "V (Vlang) ⚡", "vlang": "V (Vlang) ⚡", "erlang": "Erlang 🔴", "d": "D 🎯",
		"gleam": "Gleam ⭐️", "coldfusion": "ColdFusion ❄️", "cf": "ColdFusion ❄️",
		"actionscript": "ActionScript ⚡", "flex": "ActionScript ⚡", "ada": "Ada 🛡️",
		"mojo": "Mojo 🚀", "objc": "Objective-C 🍎", "objective-c": "Objective-C 🍎",
		"haxe": "Haxe 🎯", "raku": "Raku 🦋", "racket": "Racket 🌐", "chapel": "Chapel 🧱",
		"pony": "Pony 🐴", "odin": "Odin ⚔️",
	}

	for lang, st := range langStats {
		pct := math.Round((float64(st.count) / totalResults) * 100)
		avgStars := 0
		if st.count > 0 {
			avgStars = st.totalStars / st.count
		}

		type pkgCount struct {
			name  string
			count int
		}
		var pkgs []pkgCount
		for name, c := range st.packagesMap {
			pkgs = append(pkgs, pkgCount{name: name, count: c})
		}
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].count > pkgs[j].count
		})

		var dominant []string
		for i := 0; i < len(pkgs) && i < 3; i++ {
			dominant = append(dominant, pkgs[i].name)
		}

		displayName := strings.ToUpper(lang[:1]) + lang[1:]
		if mapped, ok := langIcons[strings.ToLower(lang)]; ok {
			displayName = mapped
		}

		distributions = append(distributions, domain.LanguageDistribution{
			Language:      displayName,
			Count:         st.count,
			Percentage:    pct,
			TopSimilarity: st.topSimilarity,
			AvgStars:      avgStars,
			DominantTech:  dominant,
		})
	}

	sort.Slice(distributions, func(i, j int) bool {
		return distributions[i].Count > distributions[j].Count
	})

	return distributions
}

func inferLanguageFromRepo(repo domain.ExtractedRepository) string {
	for _, f := range repo.Files {
		p := strings.ToLower(f.Path)
		if strings.HasSuffix(p, ".go") || p == "go.mod" {
			return "go"
		}
		if strings.HasSuffix(p, ".rs") || p == "cargo.toml" {
			return "rust"
		}
		if strings.HasSuffix(p, ".py") || p == "pyproject.toml" {
			return "python"
		}
		if strings.HasSuffix(p, ".ts") || strings.HasSuffix(p, ".js") || p == "package.json" {
			return "typescript"
		}
	}
	return "poliglota"
}
