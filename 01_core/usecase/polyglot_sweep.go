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

// DefaultPolyglotSuite lista padrão de ecossistemas para o modo varredura (--sweep).
var DefaultPolyglotSuite = []string{"go", "typescript", "python", "rust", "java", "cpp"}

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
			// Tentar inferir da descrição ou fallback
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

	for lang, st := range langStats {
		pct := math.Round((float64(st.count) / totalResults) * 100)
		avgStars := 0
		if st.count > 0 {
			avgStars = st.totalStars / st.count
		}

		// Extrair top 3 pacotes dominantes por ecossistema
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
		switch lang {
		case "go":
			displayName = "Go 🐹"
		case "typescript", "ts":
			displayName = "TypeScript 🟦"
		case "python", "py":
			displayName = "Python 🐍"
		case "rust", "rs":
			displayName = "Rust 🦀"
		case "java":
			displayName = "Java ☕"
		case "cpp", "c++", "c":
			displayName = "C / C++ ⚙️"
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
