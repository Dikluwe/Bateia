/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/resume_diagnostic.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/usecase"
	"fmt"
	"io"
)

func writeResumeDiagnostic(stdout io.Writer, assessment *usecase.EvidenceCheckpointAssessment) {
	if assessment == nil {
		return
	}
	messages := map[usecase.EvidenceCheckpointReason]string{
		usecase.CheckpointReusable:            "checkpoint compatível e hashes conferidos",
		usecase.CheckpointManifestUnavailable: "não foi possível ler ou interpretar o manifesto anterior",
		usecase.CheckpointStatusIneligible:    "estado da execução anterior não permite retomada",
		usecase.CheckpointSchemaUnsupported:   "schema do manifesto anterior incompatível",
		usecase.CheckpointEvidenceChanged:     "entradas da evidência mudaram",
		usecase.CheckpointRulesChanged:        "catálogo de regras da evidência mudou",
		usecase.CheckpointModelChanged:        "presença ou hash do modelo de triagem mudou",
		usecase.CheckpointVersionChanged:      "versão de regra mudou",
		usecase.CheckpointReaderUnavailable:   "leitor dos artefatos indisponível",
		usecase.CheckpointArtifactUnreadable:  "não foi possível ler o artefato anterior",
		usecase.CheckpointArtifactMismatch:    "hash do artefato diverge do manifesto anterior",
		usecase.CheckpointArtifactUnrecorded:  "registro com hash do artefato obrigatório ausente",
	}
	code := string(assessment.Reason)
	message, known := messages[assessment.Reason]
	if !known {
		code, message = "UNKNOWN_CHECKPOINT_REASON", "motivo indisponível nesta versão da apresentação"
	}
	state := "não reutilizada"
	if assessment.Reusable {
		state = "reutilizada"
	}
	subject := ""
	if known {
		switch assessment.Subject {
		case "run.json", "status", "schema_version", "evidence", "evidence_rules", "screening_model",
			"expansion", "screening", "analyzer", "implementation_search", "discovery_query", "local_snapshot", "local_collection",
			"discovery-snapshot.json", "evidence/bundles.json":
			subject = "; " + assessment.Subject
		}
	}
	_, _ = fmt.Fprintf(stdout, "retomada: evidência anterior %s — %s (%s%s)\n", state, message, code, subject)
}
