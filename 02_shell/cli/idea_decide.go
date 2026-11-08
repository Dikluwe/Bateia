/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_decide.md
 * @layer L2
 * @updated 2026-10-28
 */

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func runIdeaDecide(args []string, stdout io.Writer) error {
	if len(args) > 0 && args[0] == "frontier" {
		return runIdeaDecideFrontier(args[1:], stdout)
	}
	if len(args) > 0 && args[0] == "evaluate" {
		return runIdeaDecideEvaluate(args[1:], stdout)
	}
	fs := flag.NewFlagSet("idea decide", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ideaPath := fs.String("idea", "", "IdeaSpec JSON")
	assessmentPath := fs.String("assessment", "", "AssessmentReport JSON")
	policyPath := fs.String("policy", "", "ReuseDecisionPolicy JSON")
	compositionPath := fs.String("composition", "", "CompositionInput JSON opcional")
	technicalEvidencePath := fs.String("technical-evidence", "", "Achados técnicos especializados JSON opcionais; apenas consultivos")
	technicalBundlesPath := fs.String("technical-bundles", "", "Bundles coletados para verificar achados técnicos estruturados")
	evidenceAsOf := fs.String("evidence-as-of", "", "Data RFC3339 para avaliar validade dos achados; omita para agora")
	outputPath := fs.String("output", "", "arquivo de saída; omita para stdout")
	tracePath := fs.String("trace-output", "", "grafo consultivo de dependências; não altera o veredito")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea decide inválidos: %w", err)
	}
	if *ideaPath == "" || *assessmentPath == "" || *policyPath == "" {
		return fmt.Errorf("idea decide exige --idea, --assessment e --policy")
	}
	traceInputs := []string{*ideaPath, *assessmentPath, *policyPath, *compositionPath, *technicalEvidencePath, *technicalBundlesPath, *outputPath}
	if *tracePath != "" {
		if err := validateDecisionTraceOutput(*tracePath, traceInputs...); err != nil {
			return err
		}
		if *outputPath != "" {
			if err := validateDecisionTraceOutput(*outputPath, traceInputs[:len(traceInputs)-1]...); err != nil {
				return err
			}
		}
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var assessment domain.AssessmentReport
	if err := readDecisionJSON(*assessmentPath, &assessment); err != nil {
		return err
	}
	var policy domain.ReuseDecisionPolicy
	if err := readDecisionJSON(*policyPath, &policy); err != nil {
		return err
	}
	var composition *domain.CompositionInput
	if *compositionPath != "" {
		composition = &domain.CompositionInput{}
		if err := readDecisionJSON(*compositionPath, composition); err != nil {
			return err
		}
	}
	var technicalReview *domain.TechnicalEvidenceReview
	if *technicalEvidencePath != "" {
		var batch domain.TechnicalEvidenceBatch
		if err := readDecisionJSON(*technicalEvidencePath, &batch); err != nil {
			return err
		}
		asOf, err := parseTechnicalEvidenceAsOf(*evidenceAsOf)
		if err != nil {
			return err
		}
		var review domain.TechnicalEvidenceReview
		if *technicalBundlesPath == "" {
			review, err = usecase.ReviewTechnicalEvidence(batch, asOf)
		} else {
			bundles, readErr := readTechnicalBundles(*technicalBundlesPath)
			if readErr != nil {
				return readErr
			}
			review, err = usecase.VerifyTechnicalEvidence(batch, bundles, asOf)
			if err == nil {
				var applied bool
				assessment.Assessments, applied = usecase.ApplyVerifiedTechnicalLicenses(assessment.Assessments, review)
				if applied {
					review.DecisionImpact = "VERIFIED_LICENSE_GATE"
				}
			}
		}
		if err != nil {
			return fmt.Errorf("achados técnicos inválidos: %w", err)
		}
		technicalReview = &review
	} else if *evidenceAsOf != "" || *technicalBundlesPath != "" {
		return fmt.Errorf("--evidence-as-of e --technical-bundles exigem --technical-evidence")
	}
	report, err := usecase.DecideReuse(domain.ReuseDecisionRequest{Idea: idea, Assessments: assessment.Assessments, Policy: policy, Composition: composition})
	if err != nil {
		return fmt.Errorf("não foi possível decidir reúso: %w", err)
	}
	report.TechnicalEvidence = technicalReview
	if *tracePath != "" {
		graph, err := usecase.BuildDecisionTrace(domain.DecisionTraceBuildInput{SchemaVersion: "1", Idea: &idea, Report: report})
		if err != nil {
			return fmt.Errorf("construir rastreabilidade: %w", err)
		}
		if err := writeDecisionTraceJSON(graph, *tracePath, nil, traceInputs...); err != nil {
			return err
		}
	}
	return writeIdeaJSON(report, *outputPath, stdout)
}

func parseTechnicalEvidenceAsOf(value string) (time.Time, error) {
	if value == "" {
		return commandHost().Now().UTC(), nil
	}
	asOf, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("--evidence-as-of exige RFC3339: %w", err)
	}
	return asOf, nil
}

func readDecisionJSON(path string, target any) error {
	data, err := readIdeaInput(path)
	if err != nil {
		return fmt.Errorf("não foi possível ler %s: %w", path, err)
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	return nil
}
