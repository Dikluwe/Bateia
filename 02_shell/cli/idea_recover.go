/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recover.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func RunIdeaRecoverCommand(args []string, stdout io.Writer, provider domain.ContentProvider) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: bateia idea recover <diagnose|plan|value|execute>")
	}
	switch args[0] {
	case "diagnose":
		return runRecoverDiagnose(args[1:], stdout)
	case "plan":
		return runRecoverPlan(args[1:], stdout)
	case "value":
		return runRecoverValue(args[1:], stdout)
	case "execute":
		return runRecoverExecute(args[1:], stdout, provider)
	default:
		return fmt.Errorf("subcomando recover desconhecido %q", args[0])
	}
}
func runRecoverDiagnose(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea recover diagnose", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	run := fs.String("run", "", "diretório do dossiê")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *run == "" {
		return fmt.Errorf("recover diagnose exige --run")
	}
	d, err := loadRecoveryDiagnosis(*run)
	if err != nil {
		return err
	}
	return writeIdeaJSON(d, *output, stdout)
}
func runRecoverPlan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("idea recover plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	diagnosisPath := fs.String("diagnosis", "", "diagnóstico JSON")
	run := fs.String("run", "", "diretório do dossiê original")
	policyPath := fs.String("policy", "", "política JSON")
	round := fs.Int("round", 1, "rodada")
	output := fs.String("output", "", "arquivo de saída; omita para stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *diagnosisPath == "" || *policyPath == "" || *run == "" {
		return fmt.Errorf("recover plan exige --run, --diagnosis e --policy")
	}
	var d domain.InconclusiveDiagnosis
	var p domain.RecoveryPolicy
	if err := readJSONFile(*diagnosisPath, &d); err != nil {
		return err
	}
	if err := readJSONFile(*policyPath, &p); err != nil {
		return err
	}
	history, err := recoveryOriginalTargets(*run)
	if err != nil {
		return err
	}
	plan, err := usecase.PlanEvidenceRecovery(d, p, *round, history)
	if err != nil {
		return err
	}
	return writeIdeaJSON(plan, *output, stdout)
}
func runRecoverExecute(args []string, stdout io.Writer, provider domain.ContentProvider) error {
	fs := flag.NewFlagSet("idea recover execute", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	run := fs.String("run", "", "diretório do dossiê")
	policyPath := fs.String("policy", "", "política de recuperação")
	resume := fs.Bool("resume", false, "continua índice existente")
	timeout := fs.Duration("timeout", 2*time.Minute, "limite da rodada")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *run == "" || *policyPath == "" {
		return fmt.Errorf("recover execute exige --run e --policy")
	}
	if provider == nil {
		return fmt.Errorf("recover execute exige provedor de conteúdo")
	}
	var policy domain.RecoveryPolicy
	if err := readJSONFile(*policyPath, &policy); err != nil {
		return err
	}
	indexPath := filepath.Join(*run, "recovery", "index.json")
	index := domain.RecoveryIndex{SchemaVersion: "1"}
	if h, err := fileHash(filepath.Join(*run, "run.json")); err == nil {
		index.OriginalRunHash = h
	} else {
		return err
	}
	if *resume {
		_ = readJSONFile(indexPath, &index)
	} else if _, err := os.Stat(indexPath); err == nil {
		return fmt.Errorf("recuperação já existe; use --resume")
	}
	round := len(index.Rounds) + 1
	if round > policy.MaxRounds {
		return fmt.Errorf("limite de %d rodadas atingido", policy.MaxRounds)
	}
	diagnosis, err := loadRecoveryDiagnosisFromLatest(*run, index)
	if err != nil {
		return err
	}
	previous, err := recoveryOriginalTargets(*run)
	if err != nil {
		return err
	}
	for i := 1; i < round; i++ {
		var old domain.EvidenceRecoveryPlan
		if readJSONFile(filepath.Join(*run, "recovery", fmt.Sprintf("round-%02d", i), "plan.json"), &old) == nil {
			previous = append(previous, old.Targets...)
		}
	}
	plan, err := usecase.PlanEvidenceRecovery(diagnosis, policy, round, previous)
	if err != nil {
		return err
	}
	if len(plan.Targets) == 0 {
		return fmt.Errorf("nenhum target novo recuperável")
	}
	idea, screening, bundles, before, decisionPolicy, err := loadRecoveryInputs(*run, index)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, err := usecase.ExecuteEvidenceRecovery(ctx, idea, screening, bundles, before, plan, provider, usecase.DefaultEvidenceRules(), decisionPolicy)
	if err != nil {
		return err
	}
	dir := filepath.Join(*run, "recovery", fmt.Sprintf("round-%02d", round))
	artifacts := []struct {
		name  string
		value any
	}{{"diagnosis.json", diagnosis}, {"plan.json", plan}, {"bundles.json", result.Bundles}, {"assessment-report.json", result.Assessment}, {"reuse-decision.json", result.Decision}, {"remaining-diagnosis.json", result.Diagnosis}, {"metrics.json", result.Metrics}}
	for _, a := range artifacts {
		if err := writeJSONAtomic(filepath.Join(dir, a.name), a.value); err != nil {
			return err
		}
	}
	entry := domain.RecoveryRoundIndex{Round: round, DiagnosisHash: hashJSON(diagnosis), PlanHash: hashJSON(plan), BundlesHash: hashJSON(result.Bundles), AssessmentHash: hashJSON(result.Assessment), DecisionHash: hashJSON(result.Decision), MetricsHash: hashJSON(result.Metrics)}
	index.Rounds = append(index.Rounds, entry)
	if err := writeJSONAtomic(indexPath, index); err != nil {
		return err
	}
	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "recuperação round-%02d: %s → %s, %d findings restantes\n", round, before.Decision, result.Decision.Decision, len(result.Diagnosis.Findings))
	}
	return nil
}

func recoveryOriginalTargets(run string) ([]domain.RecoveryTarget, error) {
	var bundles []domain.EvidenceBundle
	if err := readJSONFile(filepath.Join(run, "evidence", "bundles.json"), &bundles); err != nil {
		return nil, err
	}
	out := []domain.RecoveryTarget{}
	for _, b := range bundles {
		for _, target := range b.Plan.Targets {
			out = append(out, domain.RecoveryTarget{Candidate: b.Repository, Capability: firstString(target.Capabilities), Path: target.Path, Kind: target.Kind, RuleID: recoveryHistoryRule(target.Kind)})
		}
		for _, target := range b.Targets {
			out = append(out, domain.RecoveryTarget{Candidate: b.Repository, Capability: firstString(target.Target.Capabilities), Path: target.Target.Path, Kind: target.Target.Kind, RuleID: recoveryHistoryRule(target.Target.Kind)})
		}
	}
	return out, nil
}
func recoveryHistoryRule(kind domain.EvidenceKind) string {
	switch kind {
	case domain.EvidenceLicense:
		return "GATE-LICENSE-001"
	case domain.EvidenceManifest:
		return "GATE-OPERATION-001"
	}
	return "original"
}
func firstString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func loadRecoveryDiagnosis(run string) (domain.InconclusiveDiagnosis, error) {
	var manifest domain.IdeaRunManifest
	var decision domain.ReuseDecisionReport
	var assessment domain.AssessmentReport
	var bundles []domain.EvidenceBundle
	if err := readJSONFile(filepath.Join(run, "run.json"), &manifest); err != nil {
		return domain.InconclusiveDiagnosis{}, err
	}
	if err := readJSONFile(filepath.Join(run, "reuse-decision.json"), &decision); err != nil {
		return domain.InconclusiveDiagnosis{}, err
	}
	if err := readJSONFile(filepath.Join(run, "assessment-report.json"), &assessment); err != nil {
		return domain.InconclusiveDiagnosis{}, err
	}
	if err := readJSONFile(filepath.Join(run, "evidence", "bundles.json"), &bundles); err != nil {
		return domain.InconclusiveDiagnosis{}, err
	}
	hash, err := fileHash(filepath.Join(run, "reuse-decision.json"))
	if err != nil {
		return domain.InconclusiveDiagnosis{}, err
	}
	return usecase.DiagnoseInconclusive(manifest.RunID, hash, decision, assessment, bundles, usecase.DefaultEvidenceRules()), nil
}
func loadRecoveryDiagnosisFromLatest(run string, index domain.RecoveryIndex) (domain.InconclusiveDiagnosis, error) {
	if len(index.Rounds) == 0 {
		return loadRecoveryDiagnosis(run)
	}
	var d domain.InconclusiveDiagnosis
	path := filepath.Join(run, "recovery", fmt.Sprintf("round-%02d", len(index.Rounds)), "remaining-diagnosis.json")
	return d, readJSONFile(path, &d)
}
func loadRecoveryInputs(run string, index domain.RecoveryIndex) (domain.IdeaSpec, domain.ScreeningReport, []domain.EvidenceBundle, domain.ReuseDecisionReport, domain.ReuseDecisionPolicy, error) {
	var idea domain.IdeaSpec
	var screening domain.ScreeningReport
	var bundles []domain.EvidenceBundle
	var decision domain.ReuseDecisionReport
	var policy domain.ReuseDecisionPolicy
	if err := readJSONFile(filepath.Join(run, "idea.json"), &idea); err != nil {
		return idea, screening, nil, decision, policy, err
	}
	if err := readJSONFile(filepath.Join(run, "screening-report.json"), &screening); err != nil {
		return idea, screening, nil, decision, policy, err
	}
	base := run
	if len(index.Rounds) > 0 {
		base = filepath.Join(run, "recovery", fmt.Sprintf("round-%02d", len(index.Rounds)))
	}
	bundlePath := filepath.Join(base, "bundles.json")
	if len(index.Rounds) == 0 {
		bundlePath = filepath.Join(run, "evidence", "bundles.json")
	}
	if err := readJSONFile(bundlePath, &bundles); err != nil {
		return idea, screening, nil, decision, policy, err
	}
	if err := readJSONFile(filepath.Join(base, "reuse-decision.json"), &decision); err != nil {
		return idea, screening, nil, decision, policy, err
	}
	if err := readJSONFile(filepath.Join(run, "reuse-policy.json"), &policy); err != nil {
		return idea, screening, nil, decision, policy, err
	}
	return idea, screening, bundles, decision, policy, nil
}
