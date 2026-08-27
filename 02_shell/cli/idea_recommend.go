/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend.md
 * @layer L2
 * @updated 2026-08-26
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type repeatedPaths []string

func (r *repeatedPaths) String() string     { return strings.Join(*r, ",") }
func (r *repeatedPaths) Set(v string) error { *r = append(*r, v); return nil }

func RunIdeaRecommendCommand(args []string, stdout io.Writer, semantic usecase.SemanticExpander, discovery domain.ObservableDiscoveryProvider, content domain.ContentProvider) error {
	fs := flag.NewFlagSet("idea recommend", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	ideaPath := fs.String("idea", "", "IdeaSpec JSON")
	policyPath := fs.String("policy", "", "ReuseDecisionPolicy JSON")
	mode := fs.String("mode", "replay", "live, local ou replay")
	provider := fs.String("provider", "sourcegraph", "provedor live")
	strategy := fs.String("strategy", "lexical", "literal ou lexical")
	language := fs.String("language", "", "linguagem opcional")
	replay := fs.String("replay", "", "dossiê fonte para replay")
	output := fs.String("output", "", "diretório do dossiê")
	resume := fs.Bool("resume", false, "retoma checkpoints válidos")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo de resultados por expansão")
	maxCandidates := fs.Int("max-candidates", 5, "máximo de candidatos para coleta")
	timeout := fs.Duration("timeout", 2*time.Minute, "limite total da execução")
	var locals repeatedPaths
	fs.Var(&locals, "local", "worktree local; repetível")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea recommend inválidos: %w", err)
	}
	if *ideaPath == "" || *policyPath == "" || *output == "" {
		return fmt.Errorf("idea recommend exige --idea, --policy e --output")
	}
	runMode := domain.IdeaRunMode(strings.ToLower(*mode))
	if runMode == domain.IdeaRunLive && strings.ToLower(*provider) != "sourcegraph" {
		return fmt.Errorf("provedor %q desconhecido", *provider)
	}
	if runMode == domain.IdeaRunReplay && *replay == "" {
		return fmt.Errorf("modo replay exige --replay")
	}
	if runMode == domain.IdeaRunLocal && len(locals) == 0 {
		return fmt.Errorf("modo local exige ao menos um --local")
	}
	if *maxPerQuery < 1 || *maxPerQuery > 25 || *maxCandidates < 1 || *maxCandidates > 20 {
		return fmt.Errorf("limites de candidatos inválidos")
	}
	if err := prepareDossier(*output, *resume); err != nil {
		return err
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	var policy domain.ReuseDecisionPolicy
	if err = readJSONFile(*policyPath, &policy); err != nil {
		return err
	}
	req := domain.IdeaRunRequest{Idea: idea, Policy: policy, Mode: runMode, Strategy: strings.ToLower(*strategy), Provider: *provider, Language: *language, LocalPaths: append([]string(nil), locals...), ReplayDirectory: *replay, OutputDirectory: *output, Resume: *resume, MaxPerQuery: *maxPerQuery, MaxCandidates: *maxCandidates}
	inputHashes := map[string]string{"idea": hashJSON(idea), "policy": hashJSON(policy), "mode": hashBytes([]byte(runMode))}
	var previous domain.IdeaRunManifest
	canReuseEvidence := false
	if *resume {
		if e := readJSONFile(filepath.Join(*output, "run.json"), &previous); e == nil {
			_, snapshotErr := os.Stat(filepath.Join(*output, "discovery-snapshot.json"))
			_, bundlesErr := os.Stat(filepath.Join(*output, "evidence", "bundles.json"))
			canReuseEvidence = previous.InputHashes["idea"] == inputHashes["idea"] && snapshotErr == nil && bundlesErr == nil
		}
	}
	manifest := domain.IdeaRunManifest{SchemaVersion: "1", RunID: usecase.StableRunID(hashJSON(req)), Mode: runMode, Status: domain.IdeaRunRunning, InputHashes: inputHashes, RuleVersions: map[string]string{"expansion": "1", "screening": "1", "evidence": usecase.DefaultEvidenceRules().Version, "decision": "1"}, Stages: []domain.IdeaRunStage{}}
	if err = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest); err != nil {
		return err
	}
	deps := usecase.IdeaRunDependencies{Semantic: semantic, Discovery: discovery, Content: content}
	source := *replay
	if canReuseEvidence {
		source = *output
	}
	if runMode == domain.IdeaRunReplay || canReuseEvidence {
		var snap domain.DiscoverySnapshot
		var bundles []domain.EvidenceBundle
		if e := readJSONFile(filepath.Join(source, "discovery-snapshot.json"), &snap); e != nil {
			manifest.Status = domain.IdeaRunFailed
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "replay", Message: e.Error()}}
			_ = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest)
			return e
		}
		if e := readJSONFile(filepath.Join(source, "evidence", "bundles.json"), &bundles); e != nil {
			manifest.Status = domain.IdeaRunInconclusive
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "replay", Message: e.Error()}}
			_ = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest)
			return e
		}
		deps.ReplaySnapshot = &snap
		deps.ReplayBundles = bundles
	}
	if runMode == domain.IdeaRunLocal {
		deps.LocalCandidates, err = localCandidates(idea, locals)
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	result, runErr := usecase.RunIdeaPipeline(ctx, req, deps)
	if runErr != nil {
		manifest.Status = domain.IdeaRunFailed
		manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "pipeline", Message: runErr.Error()}}
		_ = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest)
		return runErr
	}
	artifacts := []struct {
		path  string
		value any
	}{{"idea.json", idea}, {"reuse-policy.json", policy}, {"expansion-plan.json", result.Expansion}, {"discovery-snapshot.json", result.Discovery}, {"screening-report.json", result.Screening}, {filepath.Join("evidence", "plans.json"), result.Plans}, {filepath.Join("evidence", "bundles.json"), result.Bundles}, {"assessment-report.json", result.Assessment}, {"reuse-decision.json", result.Decision}}
	for _, a := range artifacts {
		if err = writeJSONAtomic(filepath.Join(*output, a.path), a.value); err != nil {
			return err
		}
	}
	if err = writeDecisionMarkdown(filepath.Join(*output, "decision.md"), result.Decision); err != nil {
		return err
	}
	manifest.Status = usecase.IdeaRunStatusFor(result)
	manifest.DecisionPath = "reuse-decision.json"
	manifest.Failures = append(manifest.Failures, result.Failures...)
	manifest.Stages = buildRunStages(*output, result)
	if err = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest); err != nil {
		return err
	}
	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "dossiê %s: %s (%s)\n", manifest.RunID, result.Decision.Decision, *output)
	}
	return nil
}

func prepareDossier(path string, resume bool) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("diretório de saída vazio")
	}
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("diretório de saída não pode ser symlink")
		}
		entries, e := os.ReadDir(path)
		if e != nil {
			return e
		}
		if len(entries) > 0 && !resume {
			return fmt.Errorf("diretório de saída já contém artefatos; use --resume")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, 0755)
}
func writeJSONAtomic(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeAtomic(path, data, 0644)
}
func writeAtomic(path string, data []byte, mode os.FileMode) error {
	if err := rejectSerializedSecrets(data); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".bateia-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err = tmp.Chmod(mode); err == nil {
		_, err = tmp.Write(data)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(name, path)
}
func rejectSerializedSecrets(data []byte) error {
	for _, key := range []string{"GITHUB_TOKEN", "SOURCEGRAPH_TOKEN"} {
		secret := os.Getenv(key)
		if len(secret) >= 8 && strings.Contains(string(data), secret) {
			return fmt.Errorf("artefato contém segredo de %s", key)
		}
	}
	return nil
}
func hashJSON(v any) string     { b, _ := json.Marshal(v); return hashBytes(b) }
func hashBytes(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func fileHash(path string) (string, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return "", e
	}
	return hashBytes(b), nil
}
func localCandidates(idea domain.IdeaSpec, paths []string) ([]domain.RepositoryCandidate, error) {
	var out []domain.RepositoryCandidate
	for _, path := range paths {
		abs, e := filepath.Abs(path)
		if e != nil {
			return nil, e
		}
		h, e := treeHash(abs)
		if e != nil {
			return nil, e
		}
		base := strings.ToLower(filepath.Base(abs))
		pathHash := hashBytes([]byte(abs))
		repo := base + "-" + pathHash[:8]
		out = append(out, domain.RepositoryCandidate{Owner: "local", Repo: repo, URL: "file://" + abs, Description: idea.Problem + " " + strings.Join(idea.Capabilities, " "), CommitID: h})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out, nil
}
func treeHash(root string) (string, error) {
	h := sha256.New()
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "target" || d.Name() == "node_modules") {
			if path != root {
				return filepath.SkipDir
			}
		}
		if d.IsDir() {
			return nil
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		_, _ = h.Write([]byte(filepath.ToSlash(rel)))
		_, _ = h.Write([]byte{0})
		_, _ = h.Write(b)
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func writeDecisionMarkdown(path string, r domain.ReuseDecisionReport) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Decisão de reúso: %s\n\n", r.Decision)
	fmt.Fprintf(&b, "Fonte canônica: [reuse-decision.json](reuse-decision.json)\n\n")
	if r.Primary != nil {
		fmt.Fprintf(&b, "## Ponto de partida\n\n%s/%s@%s\n\n", r.Primary.Identity.Owner, r.Primary.Identity.Repo, r.Primary.Commit)
	}
	fmt.Fprintln(&b, "## Razões")
	fmt.Fprintln(&b)
	for _, x := range r.Reasons {
		fmt.Fprintf(&b, "- `%s`: %s\n", x.RuleID, x.Message)
	}
	fmt.Fprintln(&b, "\n## Lacunas")
	fmt.Fprintln(&b)
	for _, x := range r.Gaps {
		fmt.Fprintf(&b, "- %s\n", x)
	}
	fmt.Fprintln(&b, "\n## Evidências\n\nConsulte [assessment-report.json](assessment-report.json) e [evidence/bundles.json](evidence/bundles.json).")
	return writeAtomic(path, []byte(b.String()), 0644)
}
func buildRunStages(root string, r domain.IdeaRunResult) []domain.IdeaRunStage {
	defs := []struct {
		name, path string
		count      int64
	}{{"expand", "expansion-plan.json", int64(len(r.Expansion.Expansions))}, {"discover", "discovery-snapshot.json", int64(len(r.Discovery.Candidates))}, {"screen", "screening-report.json", int64(len(r.Screening.Decisions))}, {"collect", filepath.Join("evidence", "bundles.json"), int64(len(r.Bundles))}, {"assess", "assessment-report.json", int64(len(r.Assessment.Assessments))}, {"decide", "reuse-decision.json", 1}}
	out := make([]domain.IdeaRunStage, 0, len(defs))
	ideaHash, _ := fileHash(filepath.Join(root, "idea.json"))
	prev := map[string]string{"idea.json": ideaHash}
	for _, d := range defs {
		hash, _ := fileHash(filepath.Join(root, d.path))
		if d.name == "decide" {
			policyHash, _ := fileHash(filepath.Join(root, "reuse-policy.json"))
			prev["reuse-policy.json"] = policyHash
		}
		out = append(out, domain.IdeaRunStage{Name: d.name, ContractVersion: "1", Status: domain.IdeaStageComplete, InputHashes: prev, Outputs: []domain.IdeaRunArtifact{{Path: d.path, Hash: hash}}, Counts: map[string]int64{"items": d.count}})
		out[len(out)-1].DurationMillis = r.StageDurations[d.name]
		prev = map[string]string{d.path: hash}
	}
	return out
}
