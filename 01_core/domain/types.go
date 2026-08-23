/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/types.md
 * @layer L1
 * @updated 2026-07-25
 */

package domain

import (
	"fmt"
	"time"
)

// CommitRef representa a referência imutável a um repositório num dado commit.
type CommitRef struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Commit string `json:"commit"`
}

func (c CommitRef) Key() string {
	if c.Commit == "" || c.Commit == "HEAD" {
		return fmt.Sprintf("%s/%s", c.Owner, c.Repo)
	}
	return fmt.Sprintf("%s/%s@%s", c.Owner, c.Repo, c.Commit)
}

// RepositoryCandidate é um candidato retornado no topo do funil (Descoberta).
type RepositoryCandidate struct {
	Owner        string    `json:"owner"`
	Repo         string    `json:"repo"`
	URL          string    `json:"url"`
	Description  string    `json:"description"`
	Stars        int       `json:"stars"`
	CommitID     string    `json:"commit_id"`
	ConsultadoEm time.Time `json:"consultado_em"`
	Stale        bool      `json:"stale"`
}

// Dependency representa uma dependência declarada em manifesto.
type Dependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Indirect bool   `json:"indirect"`
	Type     string `json:"type,omitempty"` // "direct", "dev", "peer"
}

// Manifest representa o manifesto de dependências e módulo triado (ex: go.mod, package.json).
type Manifest struct {
	Ref          CommitRef    `json:"ref"`
	ModuleName   string       `json:"module_name"`
	Language     string       `json:"language,omitempty"` // "go", "typescript", etc.
	GoVersion    string       `json:"go_version,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
	ConsultadoEm time.Time    `json:"consultado_em"`
	Content      []byte       `json:"content,omitempty"`
}

// DirectoryNode representa um nó na árvore de arquivos do repositório.
type DirectoryNode struct {
	Path string `json:"path"`
	Type string `json:"type"` // "file" ou "dir"
	Size int64  `json:"size,omitempty"`
}

// QueryFilter define os filtros da busca executada pelo usuário.
type QueryFilter struct {
	Query       string   `json:"query"`
	Language    string   `json:"language"`
	MaxResults  int      `json:"max_results"`
	MustHave    []string `json:"must_have"`            // ex: ["cmd", "internal"]
	Depth       int      `json:"depth"`                // Profundidade da busca recursiva em grafo
	DirectRepos []string `json:"direct_repos"`         // Repositórios injetados diretamente no modo Benchmark
	LocalPath   string   `json:"local_path,omitempty"` // Worktree local analisada sem depender de provedor remoto
}

// TriagedRepository é o resultado da etapa 2 (Meio do funil - Triagem).
type TriagedRepository struct {
	Candidate       RepositoryCandidate `json:"candidate"`
	Manifest        *Manifest           `json:"manifest,omitempty"`
	Tree            []DirectoryNode     `json:"tree,omitempty"`
	Approved        bool                `json:"approved"`
	ScoreReason     string              `json:"score_reason"`
	RelevanceScore  int                 `json:"relevance_score"`
	RelevanceLabel  string              `json:"relevance_label"`
	RelevanceReason string              `json:"relevance_reason"`
}

// ASTSymbol representa um símbolo extraído de um arquivo de código (função, estrutura, interface, etc.).
type ASTSymbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"` // "function", "struct", "interface", "type"
	Signature string `json:"signature,omitempty"`
	File      string `json:"file,omitempty"`
	Line      int    `json:"line,omitempty"`
}

// ExtractedFile representa o conteúdo lido do repositório final.
type ExtractedFile struct {
	Path    string      `json:"path"`
	Content string      `json:"content"`
	Size    int64       `json:"size"`
	Symbols []ASTSymbol `json:"symbols,omitempty"`
}

// ArchViolation representa uma violação de arquitetura ou acoplamento.
type ArchViolation struct {
	Rule        string `json:"rule"`
	Severity    string `json:"severity"` // "HIGH", "MEDIUM", "LOW"
	Message     string `json:"message"`
	Target      string `json:"target,omitempty"`
	Line        int    `json:"line,omitempty"`
	CodeSnippet string `json:"code_snippet,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// SecurityAdvisory representa uma correspondência heurística com uma entrada da base local.
type SecurityAdvisory struct {
	Package           string `json:"package"`
	VulnerableVersion string `json:"vulnerable_version"`
	Severity          string `json:"severity"` // "CRITICAL", "HIGH", "MEDIUM"
	AdvisoryID        string `json:"advisory_id"`
	Message           string `json:"message"`
	Confidence        string `json:"confidence"`
	MatchType         string `json:"match_type"`
	Source            string `json:"source"`
	DatabaseUpdatedAt string `json:"database_updated_at"`
}

// ExtractedRepository é o resultado final da etapa 3 (Fundo do funil - Extração).
type ExtractedRepository struct {
	Ref                CommitRef           `json:"ref"`
	Candidate          RepositoryCandidate `json:"candidate"`
	Manifest           *Manifest           `json:"manifest,omitempty"`
	Files              []ExtractedFile     `json:"files"`
	Symbols            []ASTSymbol         `json:"symbols,omitempty"`
	ArchViolations     []ArchViolation     `json:"arch_violations,omitempty"`
	SecurityAdvisories []SecurityAdvisory  `json:"security_advisories,omitempty"`
	ConsultadoEm       time.Time           `json:"consultado_em"`
	Stale              bool                `json:"stale"`
	RelevanceScore     int                 `json:"relevance_score"`
	RelevanceLabel     string              `json:"relevance_label"`
	RelevanceReason    string              `json:"relevance_reason"`
	HealthScore        int                 `json:"health_score"`
}

// PatternMatch representa um padrão de pacote/ferramenta detectado nos repositórios.
type PatternMatch struct {
	Name       string   `json:"name"`
	Count      int      `json:"count"`
	Percentage float64  `json:"percentage"`
	FoundIn    []string `json:"found_in"`
	Versions   []string `json:"versions,omitempty"`
}

// FilePatternMatch representa um arquivo/diretório de devops detectado.
type FilePatternMatch struct {
	Name    string   `json:"name"`
	Count   int      `json:"count"`
	FoundIn []string `json:"found_in"`
}

// ArchitecturalInsights agrupa os achados das 5 categorias de necessidades ocultas.
type ArchitecturalInsights struct {
	Testing       []PatternMatch     `json:"testing"`
	Configuration []PatternMatch     `json:"configuration"`
	Database      []PatternMatch     `json:"database"`
	Observability []PatternMatch     `json:"observability"`
	DevOps        []FilePatternMatch `json:"devops"`
}

// ArchDiff representa o delta estrutural entre duas versões de um repositório.
type ArchDiff struct {
	OldRef         CommitRef    `json:"old_ref"`
	NewRef         CommitRef    `json:"new_ref"`
	AddedDeps      []Dependency `json:"added_deps,omitempty"`
	RemovedDeps    []Dependency `json:"removed_deps,omitempty"`
	AddedSymbols   []ASTSymbol  `json:"added_symbols,omitempty"`
	RemovedSymbols []ASTSymbol  `json:"removed_symbols,omitempty"`
}

// PolyglotContract representa um protocolo ou interface de comunicação compartilhada entre tecnologias.
type PolyglotContract struct {
	Protocol     string   `json:"protocol"` // "gRPC", "GraphQL", "MessageBroker", "Caching", "REST"
	Technology   string   `json:"technology"`
	Repositories []string `json:"repositories"`
}

// FunnelTelemetry estatísticas de execução e telemetria de desempenho dos estágios do funil.
type FunnelTelemetry struct {
	DiscoveryDurationMs  int64 `json:"discovery_duration_ms"`
	TriageDurationMs     int64 `json:"triage_duration_ms"`
	ExtractionDurationMs int64 `json:"extraction_duration_ms"`
	TotalDurationMs      int64 `json:"total_duration_ms"`
	CacheHits            int   `json:"cache_hits"`
	CacheMisses          int   `json:"cache_misses"`
}

// ProgressEvent descreve progresso operacional emitido pelo caso de uso.
type ProgressEvent struct {
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

// ExecutionFailure registra uma parte da execução que falhou sem invalidar resultados parciais.
type ExecutionFailure struct {
	Scope   string `json:"scope"`
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

// OperationalWarning registra uma degradação recuperável, como uma falha de cache.
type OperationalWarning struct {
	Stage     string `json:"stage"`
	Operation string `json:"operation"`
	Target    string `json:"target,omitempty"`
	Message   string `json:"message"`
}

// SearchReport é o relatório compilado final entregue ao usuário.
type SearchReport struct {
	Query                QueryFilter                 `json:"query"`
	GeneratedAt          time.Time                   `json:"generated_at"`
	TopCandidates        int                         `json:"top_candidates"`
	TriagedCount         int                         `json:"triaged_count"`
	ExtractedCount       int                         `json:"extracted_count"`
	DiscardedCount       int                         `json:"discarded_count"`
	DiscardReason        string                      `json:"discard_reason"`
	Results              []ExtractedRepository       `json:"results"`
	Insights             *ArchitecturalInsights      `json:"insights,omitempty"`
	PolyglotMatrix       []PolyglotContract          `json:"polyglot_matrix,omitempty"`
	Telemetry            *FunnelTelemetry            `json:"telemetry,omitempty"`
	Failures             []ExecutionFailure          `json:"failures,omitempty"`
	Warnings             []OperationalWarning        `json:"warnings,omitempty"`
	TargetOntology       *ProjectOntology            `json:"target_ontology,omitempty"`
	SimilarityResults    []IdeaSimilarityResult      `json:"similarity_results,omitempty"`
	Recommendations      []ImprovementRecommendation `json:"recommendations,omitempty"`
	LanguageDistribution []LanguageDistribution      `json:"language_distribution,omitempty"`
}

// LanguageDistribution estatísticas compiladas para uma linguagem na varredura poliglota.
type LanguageDistribution struct {
	Language      string   `json:"language"`
	Count         int      `json:"count"`
	Percentage    float64  `json:"percentage"`
	TopSimilarity int      `json:"top_similarity"`
	AvgStars      int      `json:"avg_stars"`
	DominantTech  []string `json:"dominant_tech,omitempty"`
}

// AlgorithmChallenge representa um algoritmo extraído com seus testes e benchmarks associados.
type AlgorithmChallenge struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"` // "Sorting", "Graph", "Tree", "DynamicProgramming", "Crypto", "Math", "General"
	Language    string   `json:"language"`
	Repository  string   `json:"repository"`
	FilePath    string   `json:"file_path"`
	CodeSnippet string   `json:"code_snippet"`
	TestCases   []string `json:"test_cases,omitempty"`
	Benchmarks  []string `json:"benchmarks,omitempty"`
	Complexity  string   `json:"complexity,omitempty"`
}

// LabDataset representa o conjunto de desafios e benchmarks extraídos para o Lab de Projetos.
type LabDataset struct {
	Query           string               `json:"query"`
	GeneratedAt     time.Time            `json:"generated_at"`
	TotalChallenges int                  `json:"total_challenges"`
	Challenges      []AlgorithmChallenge `json:"challenges"`
}

// ProjectOntology representa o perfil ontológico extraído de um repositório ou descrição de ideia.
type ProjectOntology struct {
	Ref                CommitRef `json:"ref"`
	ProblemDomain      string    `json:"problem_domain"`      // Categoria do problema (ex: "Code Intelligence", "Vector Search")
	IntentKeywords     []string  `json:"intent_keywords"`     // Palavras-chave de intenção/problema
	Capabilities       []string  `json:"capabilities"`        // Capacidades/Features detectadas
	EcosystemPackages  []string  `json:"ecosystem_packages"`  // Pacotes/dependências chave do manifesto
	ArchitecturalStyle string    `json:"architectural_style"` // Estilo arquitetural (ex: "Layered", "Clean", "Monolith")
	FileStructure      []string  `json:"file_structure"`      // Diretórios chave (ex: "cmd", "internal", "pkg")
}

// SimilarityVector detalha a pontuação decomposta nos 4 vetores ontológicos.
type SimilarityVector struct {
	DomainIntentScore int `json:"domain_intent_score"` // Vet1: 0-100 (Intenção & Semântica)
	StackMatchScore   int `json:"stack_match_score"`   // Vet2: 0-100 (Stack & Ecossistema)
	TopologyScore     int `json:"topology_score"`      // Vet3: 0-100 (Topologia Estrutural)
	CapabilityScore   int `json:"capability_score"`    // Vet4: 0-100 (Cobertura de Capacidades)
}

// IdeaSimilarityResult representa o resultado da comparação entre a ontologia alvo e um repositório candidato.
type IdeaSimilarityResult struct {
	TargetRef           CommitRef        `json:"target_ref"`
	CandidateRef        CommitRef        `json:"candidate_ref"`
	CandidateRepo       string           `json:"candidate_repo"`
	OverallScore        int              `json:"overall_score"`    // Score total ponderado (0-100%)
	SimilarityLabel     string           `json:"similarity_label"` // "Identical Idea", "High Similarity", "Moderate", "Divergent"
	Vectors             SimilarityVector `json:"vectors"`
	CommonKeywords      []string         `json:"common_keywords,omitempty"`
	CommonStack         []string         `json:"common_stack,omitempty"`
	MissingCapabilities []string         `json:"missing_capabilities,omitempty"`
}

// ImprovementRecommendation sugestão de melhoria para o projeto alvo derivada de repositórios parecidos.
type ImprovementRecommendation struct {
	Category     string   `json:"category"`      // "Testing", "Observability", "Database", "Security", "Architecture"
	Title        string   `json:"title"`         // Ex: "Adicionar Testcontainers para Testes de Integração"
	AdoptionRate int      `json:"adoption_rate"` // % de adoção entre projetos parecidos (> 60% similaridade)
	SuggestedBy  []string `json:"suggested_by"`  // Repositórios de referência que já adotam
	Reasoning    string   `json:"reasoning"`     // Motivação baseada na ontologia
}
