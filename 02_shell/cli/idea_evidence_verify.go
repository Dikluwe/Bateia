/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_evidence_verify.md
 * @layer L2
 * @updated 2026-09-24
 */

package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaEvidenceVerify(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea evidence verify", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "TechnicalEvidenceBatch JSON")
	bundlesPath := fs.String("bundles", "", "EvidenceBundle JSON ou array")
	asOfText := fs.String("as-of", "", "Data RFC3339; omita para agora")
	outputPath := fs.String("output", "", "Arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de evidence verify inválidos: %w", err)
	}
	if *inputPath == "" || *bundlesPath == "" {
		return fmt.Errorf("evidence verify exige --input e --bundles")
	}
	var batch domain.TechnicalEvidenceBatch
	if err := readDecisionJSON(*inputPath, &batch); err != nil {
		return err
	}
	bundles, err := readTechnicalBundles(*bundlesPath)
	if err != nil {
		return err
	}
	asOf, err := parseTechnicalEvidenceAsOf(*asOfText)
	if err != nil {
		return err
	}
	review, err := usecase.VerifyTechnicalEvidence(batch, bundles, asOf)
	if err != nil {
		return fmt.Errorf("achados técnicos inválidos: %w", err)
	}
	return writeIdeaJSON(review, *outputPath, stdout)
}

func readTechnicalBundles(path string) ([]domain.EvidenceBundle, error) {
	data, err := readIdeaInput(path)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, fmt.Errorf("bundle JSON vazio em %s", path)
	}
	if data[0] == '[' {
		var bundles []domain.EvidenceBundle
		if err := json.Unmarshal(data, &bundles); err != nil {
			return nil, fmt.Errorf("bundles JSON inválidos em %s: %w", path, err)
		}
		return bundles, nil
	}
	var bundle domain.EvidenceBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("bundle JSON inválido em %s: %w", path, err)
	}
	return []domain.EvidenceBundle{bundle}, nil
}
