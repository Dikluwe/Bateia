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
	Query      string   `json:"query"`
	Language   string   `json:"language"`
	MaxResults int      `json:"max_results"`
	MustHave   []string `json:"must_have"` // ex: ["cmd", "internal"]
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

// ExtractedFile representa o conteúdo lido do repositório final.
type ExtractedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

// ExtractedRepository é o resultado final da etapa 3 (Fundo do funil - Extração).
type ExtractedRepository struct {
	Ref             CommitRef           `json:"ref"`
	Candidate       RepositoryCandidate `json:"candidate"`
	Manifest        *Manifest           `json:"manifest,omitempty"`
	Files           []ExtractedFile     `json:"files"`
	ConsultadoEm    time.Time           `json:"consultado_em"`
	Stale           bool                `json:"stale"`
	RelevanceScore  int                 `json:"relevance_score"`
	RelevanceLabel  string              `json:"relevance_label"`
	RelevanceReason string              `json:"relevance_reason"`
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

// SearchReport é o relatório compilado final entregue ao usuário.
type SearchReport struct {
	Query          QueryFilter            `json:"query"`
	GeneratedAt    time.Time              `json:"generated_at"`
	TopCandidates  int                    `json:"top_candidates"`
	TriagedCount   int                    `json:"triaged_count"`
	ExtractedCount int                    `json:"extracted_count"`
	DiscardedCount int                    `json:"discarded_count"`
	DiscardReason  string                 `json:"discard_reason"`
	Results        []ExtractedRepository  `json:"results"`
	Insights       *ArchitecturalInsights `json:"insights,omitempty"`
}
