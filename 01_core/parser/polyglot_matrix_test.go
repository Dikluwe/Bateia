/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/polyglot_matrix_test.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestEvaluatePolyglotMatrix(t *testing.T) {
	repos := []domain.ExtractedRepository{
		{
			Candidate: domain.RepositoryCandidate{Owner: "go-service", Repo: "app"},
			Manifest: &domain.Manifest{
				Dependencies: []domain.Dependency{
					{Name: "google.golang.org/grpc", Version: "v1.50.0"},
				},
			},
		},
	}

	matrix := EvaluatePolyglotMatrix(repos)
	if len(matrix) == 0 {
		t.Fatalf("Esperava encontrar contrato gRPC na matriz polyglot")
	}
	if matrix[0].Protocol != "gRPC / Protobuf" {
		t.Errorf("Esperava protocolo gRPC / Protobuf, obteve %s", matrix[0].Protocol)
	}
}
