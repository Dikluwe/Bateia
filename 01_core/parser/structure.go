/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/structure.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"path/filepath"
	"strings"

	"bateia/01_core/domain"
)

// StructuralScore analisa a lista de caminhos de arquivos para classificar a qualidade da arquitetura.
type StructuralScore struct {
	HasCmd      bool     `json:"has_cmd"`
	HasInternal bool     `json:"has_internal"`
	HasPkg      bool     `json:"has_pkg"`
	FoundPaths  []string `json:"found_paths"`
	Approved    bool     `json:"approved"`
	Reason      string   `json:"reason"`
}

// EvaluateStructure analisa puramente a lista de nós de diretório sem I/O.
func EvaluateStructure(nodes []domain.DirectoryNode, mustHave []string) StructuralScore {
	score := StructuralScore{
		FoundPaths: make([]string, 0),
	}

	pathSet := make(map[string]bool)
	for _, node := range nodes {
		pathSet[node.Path] = true

		parts := strings.Split(node.Path, "/")
		if len(parts) > 0 {
			pathSet[parts[0]] = true
		}
	}

	score.HasCmd = pathSet["cmd"]
	score.HasInternal = pathSet["internal"]
	score.HasPkg = pathSet["pkg"]

	missing := make([]string, 0)
	for _, req := range mustHave {
		if !matchGlobOrExact(req, pathSet, nodes) {
			missing = append(missing, req)
		}
	}

	if len(missing) > 0 {
		score.Approved = false
		score.Reason = "Faltam diretórios requeridos: " + strings.Join(missing, ", ")
	} else {
		score.Approved = true
		score.Reason = "Estrutura validada com sucesso"
	}

	return score
}

// MatchGlobRecursive verifica se um caminho de arquivo satisfaz um padrão de globbing com suporte a '**'.
func MatchGlobRecursive(pattern, path string) bool {
	if pattern == "" || path == "" {
		return pattern == path
	}

	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	patternSegs := strings.Split(pattern, "/")
	pathSegs := strings.Split(path, "/")

	return matchSegments(patternSegs, pathSegs)
}

func matchSegments(pSegs, pathSegs []string) bool {
	if len(pSegs) == 0 && len(pathSegs) == 0 {
		return true
	}
	if len(pSegs) == 0 {
		return false
	}

	if pSegs[0] == "**" {
		if len(pSegs) == 1 {
			return true
		}
		for i := 0; i <= len(pathSegs); i++ {
			if matchSegments(pSegs[1:], pathSegs[i:]) {
				return true
			}
		}
		return false
	}

	if len(pathSegs) == 0 {
		return false
	}

	matched, err := filepath.Match(pSegs[0], pathSegs[0])
	if err != nil || !matched {
		return false
	}

	return matchSegments(pSegs[1:], pathSegs[1:])
}

func matchGlobOrExact(req string, pathSet map[string]bool, nodes []domain.DirectoryNode) bool {
	if pathSet[req] {
		return true
	}

	for _, node := range nodes {
		if MatchGlobRecursive(req, node.Path) {
			return true
		}
		matchedBase, _ := filepath.Match(req, filepath.Base(node.Path))
		if matchedBase {
			return true
		}
	}

	return false
}
