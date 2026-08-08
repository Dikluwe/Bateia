/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/ontology_matcher.md
 * @layer L1
 * @updated 2026-08-08
 */

package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"bateia/01_core/domain"
)

// ExtractOntologyFromRepo analisa o repositório extraído e constrói o perfil ontológico.
func ExtractOntologyFromRepo(repo domain.ExtractedRepository) domain.ProjectOntology {
	keywordsMap := make(map[string]bool)
	capabilitiesMap := make(map[string]bool)
	var packages []string

	// 1. Palavras-chave da Descrição e do Nome
	textToAnalyze := strings.ToLower(repo.Candidate.Description + " " + repo.Candidate.Repo)
	words := strings.FieldsFunc(textToAnalyze, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '/' || r == '.' || r == ',' || r == ':'
	})

	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true, "for": true, "in": true,
		"on": true, "with": true, "is": true, "to": true, "of": true, "de": true, "com": true,
		"para": true, "uma": true, "um": true, "o": true, "os": true, "as": true, "em": true,
		"project": true, "tool": true, "library": true, "framework": true, "written": true,
	}

	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 2 && !stopWords[w] {
			keywordsMap[w] = true
		}
	}

	// 2. Extrair Pacotes do Manifesto
	if repo.Manifest != nil {
		for _, dep := range repo.Manifest.Dependencies {
			pkgName := strings.ToLower(dep.Name)
			packages = append(packages, pkgName)

			// Capacidades baseadas nas dependências conhecidas
			if strings.Contains(pkgName, "test") || strings.Contains(pkgName, "check") || strings.Contains(pkgName, "assert") {
				capabilitiesMap["testing"] = true
			}
			if strings.Contains(pkgName, "zap") || strings.Contains(pkgName, "log") || strings.Contains(pkgName, "otel") || strings.Contains(pkgName, "prometheus") {
				capabilitiesMap["observability"] = true
			}
			if strings.Contains(pkgName, "pgx") || strings.Contains(pkgName, "sql") || strings.Contains(pkgName, "orm") || strings.Contains(pkgName, "gorm") || strings.Contains(pkgName, "prisma") {
				capabilitiesMap["database"] = true
			}
			if strings.Contains(pkgName, "redis") || strings.Contains(pkgName, "cache") || strings.Contains(pkgName, "memcached") {
				capabilitiesMap["caching"] = true
			}
			if strings.Contains(pkgName, "grpc") || strings.Contains(pkgName, "proto") || strings.Contains(pkgName, "graphql") {
				capabilitiesMap["polyglot_rpc"] = true
			}
		}
	}

	// 3. Capacidades baseadas em arquivos
	for _, f := range repo.Files {
		p := strings.ToLower(f.Path)
		if strings.Contains(p, "docker") || strings.Contains(p, ".github/workflows") {
			capabilitiesMap["devops_ci"] = true
		}
		if strings.Contains(p, "cmd/") {
			capabilitiesMap["cli_entrypoint"] = true
		}
		if strings.Contains(p, "internal/") || strings.Contains(p, "pkg/") {
			capabilitiesMap["modular_architecture"] = true
		}
	}

	// Estilo Arquitetural preliminar
	archStyle := "Standard Layout"
	if capabilitiesMap["modular_architecture"] {
		archStyle = "Modular Clean/Hexagonal"
	}

	var keywords []string
	for k := range keywordsMap {
		keywords = append(keywords, k)
	}
	sort.Strings(keywords)

	var capabilities []string
	for c := range capabilitiesMap {
		capabilities = append(capabilities, c)
	}
	sort.Strings(capabilities)

	return domain.ProjectOntology{
		Ref:                repo.Ref,
		ProblemDomain:      deriveProblemDomain(keywords),
		IntentKeywords:     keywords,
		Capabilities:       capabilities,
		EcosystemPackages:  packages,
		ArchitecturalStyle: archStyle,
	}
}

// ExtractOntologyFromIdea cria uma ontologia simples a partir do texto/descrição de uma ideia.
func ExtractOntologyFromIdea(ideaText string) domain.ProjectOntology {
	words := strings.FieldsFunc(strings.ToLower(ideaText), func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '/' || r == '.' || r == ','
	})

	var keywords []string
	stopWords := map[string]bool{"a": true, "the": true, "and": true, "for": true, "with": true, "in": true, "to": true, "de": true, "com": true, "para": true}
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 2 && !stopWords[w] {
			keywords = append(keywords, w)
		}
	}

	return domain.ProjectOntology{
		Ref:                domain.CommitRef{Owner: "user", Repo: "idea-target"},
		ProblemDomain:      deriveProblemDomain(keywords),
		IntentKeywords:     keywords,
		Capabilities:       []string{"testing", "observability", "database", "caching", "devops_ci"},
		ArchitecturalStyle: "Target Vision",
	}
}

// CalculateIdeaSimilarity calcula o Score de Similaridade entre a ontologia Alvo e a ontologia Candidata.
func CalculateIdeaSimilarity(target, candidate domain.ProjectOntology) domain.IdeaSimilarityResult {
	// Vetor 1: Intenção & Semântica (Jaccard Index de Palavras-Chave)
	v1Score, commonKW := jaccardSimilarity(target.IntentKeywords, candidate.IntentKeywords)

	// Vetor 2: Ecossistema & Stack
	v2Score, commonStack := jaccardSimilarity(target.EcosystemPackages, candidate.EcosystemPackages)

	// Vetor 3: Topologia Estrutural
	v3Score := 50
	if target.ArchitecturalStyle == candidate.ArchitecturalStyle {
		v3Score = 100
	} else if strings.Contains(target.ArchitecturalStyle, "Modular") && strings.Contains(candidate.ArchitecturalStyle, "Modular") {
		v3Score = 80
	}

	// Vetor 4: Cobertura de Capacidades
	v4Score, _ := jaccardSimilarity(target.Capabilities, candidate.Capabilities)

	// Calcular Capacidades Faltantes no Alvo que a Candidata possui
	var missingCaps []string
	targetCapMap := make(map[string]bool)
	for _, c := range target.Capabilities {
		targetCapMap[c] = true
	}
	for _, c := range candidate.Capabilities {
		if !targetCapMap[c] {
			missingCaps = append(missingCaps, c)
		}
	}

	// Score Global Ponderado:
	// Vet1 (35%) + Vet2 (30%) + Vet4 (20%) + Vet3 (15%)
	overall := int(math.Round(float64(v1Score)*0.35 + float64(v2Score)*0.30 + float64(v4Score)*0.20 + float64(v3Score)*0.15))

	label := "Divergent"
	if overall >= 80 {
		label = "Identical / Core Twin Idea"
	} else if overall >= 60 {
		label = "High Similarity"
	} else if overall >= 40 {
		label = "Moderate Overlap"
	}

	candidateRepoName := candidate.Ref.Key()
	if candidateRepoName == "" || candidateRepoName == "/" {
		candidateRepoName = candidate.ProblemDomain
	}

	return domain.IdeaSimilarityResult{
		TargetRef:       target.Ref,
		CandidateRef:    candidate.Ref,
		CandidateRepo:   candidateRepoName,
		OverallScore:    overall,
		SimilarityLabel: label,
		Vectors: domain.SimilarityVector{
			DomainIntentScore: v1Score,
			StackMatchScore:   v2Score,
			TopologyScore:     v3Score,
			CapabilityScore:   v4Score,
		},
		CommonKeywords:      commonKW,
		CommonStack:         commonStack,
		MissingCapabilities: missingCaps,
	}
}

// MineOntologicalImprovements extrai recomendações de melhorias aprendidas com repositórios de alta similaridade.
func MineOntologicalImprovements(target domain.ProjectOntology, results []domain.IdeaSimilarityResult, repos []domain.ExtractedRepository) []domain.ImprovementRecommendation {
	var recommendations []domain.ImprovementRecommendation

	// Filtrar apenas candidatos relevantes com similaridade >= 40%
	var relevant []domain.IdeaSimilarityResult
	for _, r := range results {
		if r.OverallScore >= 40 {
			relevant = append(relevant, r)
		}
	}

	if len(relevant) == 0 {
		return recommendations
	}

	// Mapear adoção de pacotes e padrões em repositórios similares
	packageAdoption := make(map[string][]string)
	for _, rel := range relevant {
		for _, repo := range repos {
			if repo.Ref.Key() == rel.CandidateRef.Key() && repo.Manifest != nil {
				for _, dep := range repo.Manifest.Dependencies {
					packageAdoption[dep.Name] = append(packageAdoption[dep.Name], repo.Ref.Key())
				}
			}
		}
	}

	totalRel := float64(len(relevant))

	// Recomendações por Categoria
	for pkg, adopters := range packageAdoption {
		rate := int(math.Round((float64(len(adopters)) / totalRel) * 100))
		if rate >= 50 {
			category := "General Library"
			title := fmt.Sprintf("Adotar '%s' para fortalecer a infraestrutura", pkg)
			pkgLower := strings.ToLower(pkg)

			if strings.Contains(pkgLower, "test") || strings.Contains(pkgLower, "check") {
				category = "Testing"
				title = fmt.Sprintf("Integrar '%s' para automação de testes", pkg)
			} else if strings.Contains(pkgLower, "zap") || strings.Contains(pkgLower, "log") || strings.Contains(pkgLower, "otel") {
				category = "Observability"
				title = fmt.Sprintf("Implementar observabilidade e logging estruturado com '%s'", pkg)
			} else if strings.Contains(pkgLower, "sql") || strings.Contains(pkgLower, "pgx") || strings.Contains(pkgLower, "orm") {
				category = "Database"
				title = fmt.Sprintf("Padronizar camada de persistência com '%s'", pkg)
			}

			recommendations = append(recommendations, domain.ImprovementRecommendation{
				Category:     category,
				Title:        title,
				AdoptionRate: rate,
				SuggestedBy:  adopters,
				Reasoning:    fmt.Sprintf("%d%% dos repositórios ontologicamente parecidos utilizam '%s'.", rate, pkg),
			})
		}
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].AdoptionRate > recommendations[j].AdoptionRate
	})

	return recommendations
}

func jaccardSimilarity(s1, s2 []string) (int, []string) {
	if len(s1) == 0 || len(s2) == 0 {
		return 0, nil
	}

	set1 := make(map[string]bool)
	for _, item := range s1 {
		set1[strings.ToLower(item)] = true
	}

	var intersection []string
	unionMap := make(map[string]bool)
	for item := range set1 {
		unionMap[item] = true
	}

	for _, item := range s2 {
		lower := strings.ToLower(item)
		unionMap[lower] = true
		if set1[lower] {
			intersection = append(intersection, lower)
		}
	}

	score := float64(len(intersection)) / float64(len(unionMap))
	return int(math.Round(score * 100)), intersection
}

func deriveProblemDomain(keywords []string) string {
	if len(keywords) == 0 {
		return "General Software"
	}
	kwStr := strings.Join(keywords, " ")
	if strings.Contains(kwStr, "search") || strings.Contains(kwStr, "funnel") || strings.Contains(kwStr, "code") {
		return "Code Intelligence & Search"
	}
	if strings.Contains(kwStr, "web") || strings.Contains(kwStr, "http") || strings.Contains(kwStr, "api") {
		return "Web Service & API"
	}
	return strings.Title(keywords[0]) + " System"
}
