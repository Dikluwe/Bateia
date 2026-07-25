/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/structure.md
 * @layer L1
 * @updated 2026-07-25
 */

package parser

import (
	"strings"

	"bateia/01_core/domain"
)

// StructuralScore analisa a lista de caminhos de arquivos para classificar a qualidade da arquitetura Go.
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
		if !pathSet[req] {
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
