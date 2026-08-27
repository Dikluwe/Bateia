/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_competitive.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
	"strings"
)

func runIdeaAutophagy(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea autofagia", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ideaPath := fs.String("idea", "", "IdeaSpec JSON")
	bundlesPath := fs.String("bundles", "", "lista JSON de EvidenceBundle internos")
	rulesPath := fs.String("rules", "", "CapabilityRuleCatalog JSON")
	targetValue := fs.String("target", "", "alvo imutável owner/repo@commit")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea autofagia inválidos: %w", err)
	}
	if *ideaPath == "" || *bundlesPath == "" || *rulesPath == "" || *targetValue == "" {
		return fmt.Errorf("idea autofagia exige --idea, --bundles, --rules e --target")
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var bundles []domain.EvidenceBundle
	var rules domain.CapabilityRuleCatalog
	if err := readDecisionJSON(*bundlesPath, &bundles); err != nil {
		return err
	}
	if err := readDecisionJSON(*rulesPath, &rules); err != nil {
		return err
	}
	target, err := parseResolvedTarget(*targetValue)
	if err != nil {
		return err
	}
	report, err := usecase.AnalyzeAutophagy(idea, target, bundles, rules)
	if err != nil {
		return err
	}
	return writeIdeaJSON(report, *output, stdout)
}

func runIdeaKatamari(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea katamari", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	autophagyPath := fs.String("autofagia", "", "AutophagyReport JSON")
	assessmentsPath := fs.String("assessments", "", "AssessmentReport competitivo JSON")
	policyPath := fs.String("policy", "", "KatamariPolicy JSON")
	reusePolicyPath := fs.String("reuse-policy", "", "ReuseDecisionPolicy JSON")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea katamari inválidos: %w", err)
	}
	if *autophagyPath == "" || *assessmentsPath == "" || *policyPath == "" || *reusePolicyPath == "" {
		return fmt.Errorf("idea katamari exige --autofagia, --assessments, --policy e --reuse-policy")
	}
	var autophagy domain.AutophagyReport
	var assessments domain.AssessmentReport
	var policy domain.KatamariPolicy
	var reusePolicy domain.ReuseDecisionPolicy
	for _, input := range []struct {
		path   string
		target any
	}{{*autophagyPath, &autophagy}, {*assessmentsPath, &assessments}, {*policyPath, &policy}, {*reusePolicyPath, &reusePolicy}} {
		if err := readDecisionJSON(input.path, input.target); err != nil {
			return err
		}
	}
	proposal, err := usecase.ProposeKatamari(autophagy, assessments, policy, reusePolicy)
	if err != nil {
		return err
	}
	return writeIdeaJSON(proposal, *output, stdout)
}

func parseResolvedTarget(value string) (domain.ResolvedRepository, error) {
	parts := strings.Split(value, "@")
	if len(parts) != 2 || parts[1] == "" || strings.EqualFold(parts[1], "HEAD") {
		return domain.ResolvedRepository{}, fmt.Errorf("--target exige owner/repo@commit imutável")
	}
	repo := strings.Split(parts[0], "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		return domain.ResolvedRepository{}, fmt.Errorf("--target exige owner/repo@commit imutável")
	}
	return domain.ResolvedRepository{Identity: domain.RepositoryIdentity{Owner: repo[0], Repo: repo[1]}, Commit: parts[1]}, nil
}
