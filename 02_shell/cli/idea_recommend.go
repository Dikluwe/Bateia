/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/report"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type repeatedPaths []string

func (r *repeatedPaths) String() string     { return strings.Join(*r, ",") }
func (r *repeatedPaths) Set(v string) error { *r = append(*r, v); return nil }

type optionalPath struct {
	value string
	set   bool
}

func (p *optionalPath) String() string { return p.value }
func (p *optionalPath) Set(v string) error {
	p.value = v
	p.set = true
	return nil
}

func RunIdeaRecommendCommand(args []string, stdout io.Writer, semantic usecase.SemanticExpander, discovery domain.ObservableDiscoveryProvider, content domain.ContentProvider, recommendationServices ...*usecase.RecommendationService) (returnErr error) {
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
	rulesPath := fs.String("rules", "", "catálogo JSON de regras de evidência")
	resume := fs.Bool("resume", false, "retoma checkpoints válidos")
	withTrace := fs.Bool("trace", false, "registra grafo consultivo de dependências da decisão")
	maxPerQuery := fs.Int("max-per-query", 5, "máximo de resultados por expansão")
	maxCandidates := fs.Int("max-candidates", 5, "máximo de candidatos para coleta")
	targetName := fs.String("target-name", "", "nome do software alvo; usa idea.problem quando omitido")
	targetRepo := fs.String("target-repo", "", "repositório alvo owner/repo ou URL GitHub")
	timeout := fs.Duration("timeout", 2*time.Minute, "limite total da execução")
	var locals repeatedPaths
	var searchLocales repeatedPaths
	var screeningModelPath optionalPath
	fs.Var(&locals, "local", "worktree local; repetível")
	fs.Var(&searchLocales, "search-locale", "idioma natural/região BCP 47 para descoberta; repetível ou all")
	fs.Var(&screeningModelPath, "screening-model", "modelo JSON de screening opt-in")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de idea recommend inválidos: %w", err)
	}
	if *ideaPath == "" || *policyPath == "" || *output == "" {
		return fmt.Errorf("idea recommend exige --idea, --policy e --output")
	}
	traceInputs := []string{*ideaPath, *policyPath, *rulesPath, screeningModelPath.value}
	if *withTrace {
		if strings.EqualFold(*mode, string(domain.IdeaRunReplay)) && *replay != "" {
			traceInputs = append(traceInputs, filepath.Join(*replay, "discovery-snapshot.json"), filepath.Join(*replay, "evidence", "bundles.json"))
		}
		if *resume {
			traceInputs = append(traceInputs, filepath.Join(*output, "discovery-snapshot.json"), filepath.Join(*output, "evidence", "bundles.json"))
		}
		graphPath, indexPath := filepath.Join(*output, "trace-graph.json"), filepath.Join(*output, "trace-index.json")
		for _, paths := range [][2]string{{graphPath, indexPath}, {indexPath, graphPath}} {
			protected := append(append([]string(nil), traceInputs...), paths[1])
			if err := validateDecisionTraceOutput(paths[0], protected...); err != nil {
				return err
			}
		}
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
	if *timeout <= 0 {
		return fmt.Errorf("timeout deve ser positivo")
	}
	ctx, cancel := commandHost().WithTimeout(commandContext(context.Background()), *timeout)
	defer cancel()
	var screeningModel *domain.ScreeningModel
	screeningModelHash := ""
	if screeningModelPath.set {
		model, hash, modelErr := readOptInScreeningModel(screeningModelPath.value)
		if modelErr != nil {
			return modelErr
		}
		screeningModel = &model
		screeningModelHash = hash
	}
	regionalInput, err := usecase.PlanRegionalDiscovery(usecase.DefaultRegionalQueryCatalog(), "recipe-manager", searchLocales)
	if err != nil {
		return err
	}
	normalizedSearchLocales := append([]string(nil), regionalInput.Requested...)
	targetRepository, err := parseAlternativeTargetRepository(*targetRepo)
	if err != nil {
		return err
	}
	idea, err := readIdeaSpec(*ideaPath)
	if err != nil {
		return err
	}
	resolvedTargetName := strings.TrimSpace(*targetName)
	if resolvedTargetName == "" {
		resolvedTargetName = idea.Problem
	}
	resolvedTargetName = strings.Join(strings.Fields(resolvedTargetName), " ")
	target := domain.SoftwareAlternativeTarget{Name: resolvedTargetName, Repository: targetRepository}
	var policy domain.ReuseDecisionPolicy
	if err = readJSONFile(*policyPath, &policy); err != nil {
		return err
	}
	var evidenceRules domain.CapabilityRuleCatalog
	if strings.TrimSpace(*rulesPath) != "" {
		if err = readJSONFile(*rulesPath, &evidenceRules); err != nil {
			return err
		}
	} else {
		evidenceRules = usecase.DefaultEvidenceRules()
	}
	req := domain.IdeaRunRequest{Idea: idea, Policy: policy, Mode: runMode, Strategy: strings.ToLower(*strategy), Provider: *provider, Language: *language, SearchLocales: normalizedSearchLocales, LocalPaths: append([]string(nil), locals...), ReplayDirectory: *replay, OutputDirectory: *output, Resume: *resume, MaxPerQuery: *maxPerQuery, MaxCandidates: *maxCandidates}
	input := usecase.RecommendationDossierInput{Request: req, Target: target, Rules: evidenceRules, CustomRules: strings.TrimSpace(*rulesPath) != "", ScreeningModel: screeningModel, ScreeningModelHash: screeningModelHash, Trace: *withTrace, Dependencies: usecase.IdeaRunDependencies{Semantic: semantic, Discovery: discovery, Content: content, EvidenceRules: evidenceRules}}
	if len(recommendationServices) > 0 {
		input.Recommendation = recommendationServices[0]
	}
	_, runErr := usecase.RunRecommendationDossier(ctx, input, cliDossierFiles{output: *output, replay: *replay, traceInputs: traceInputs, stdout: stdout})
	return runErr
}

func parseAlternativeTargetRepository(value string) (*domain.RepositoryIdentity, error) {
	return usecase.ParseRepositoryIdentity(value)
}

type alternativeIdentityLookup func(string) []domain.RepositoryIdentity

func (lookup alternativeIdentityLookup) RepositoryIdentities(url string) []domain.RepositoryIdentity {
	return lookup(url)
}
func alternativeInputsForTarget(mode domain.IdeaRunMode, discovery domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, target *domain.RepositoryIdentity) (domain.DiscoverySnapshot, []domain.CandidateAssessment) {
	return alternativeInputsForTargetUsing(mode, discovery, assessments, target, localGitHubRepositoryIdentities)
}
func alternativeInputsForTargetUsing(mode domain.IdeaRunMode, discovery domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, target *domain.RepositoryIdentity, identitiesFor func(string) []domain.RepositoryIdentity) (domain.DiscoverySnapshot, []domain.CandidateAssessment) {
	return usecase.AlternativeInputsForTarget(mode, discovery, assessments, target, alternativeIdentityLookup(identitiesFor))
}

func localGitHubRepositoryIdentities(candidateURL string) []domain.RepositoryIdentity {
	return localGitHubRepositoryIdentitiesLimit(candidateURL, maxIdeaInputBytes)
}

func localGitHubRepositoryIdentitiesLimit(candidateURL string, limit int64) []domain.RepositoryIdentity {
	return commandHost().LocalGitRepositoryIdentities(candidateURL, limit)
}

func prepareDossier(path string, resume bool) error {
	return commandHost().PrepareDossier(path, resume)
}
func writeJSONAtomic(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeAtomic(path, data, 0644)
}
func writeAtomic(path string, data []byte, mode domain.FileMode) error {
	if err := rejectSerializedSecrets(data); err != nil {
		return err
	}
	return commandHost().WriteAtomic(path, data, mode)
}
func rejectSerializedSecrets(data []byte) error {
	for _, key := range []string{"GITHUB_TOKEN", "SOURCEGRAPH_TOKEN", "BATEIA_GITEE_TOKEN", "BATEIA_ATOMGIT_TOKEN"} {
		secret := commandHost().Getenv(key)
		if len(secret) >= 8 {
			encoded, _ := json.Marshal(secret)
			if strings.Contains(string(data), secret) || strings.Contains(string(data), string(encoded[1:len(encoded)-1])) {
				return fmt.Errorf("artefato contém segredo de %s", key)
			}
		}
	}
	return nil
}
func hashJSON(v any) string     { b, _ := json.Marshal(v); return hashBytes(b) }
func hashBytes(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func fileHash(path string) (string, error) {
	b, e := readIdeaInput(path)
	if e != nil {
		return "", e
	}
	return hashBytes(b), nil
}

func readOptInScreeningModel(path string) (domain.ScreeningModel, string, error) {
	data, err := readIdeaInput(path)
	if err != nil {
		return domain.ScreeningModel{}, "", fmt.Errorf("não foi possível ler o modelo de screening")
	}
	if err := validateIdeaJSONStructure(data); err != nil {
		return domain.ScreeningModel{}, "", fmt.Errorf("modelo de screening inválido")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var model domain.ScreeningModel
	if err := decoder.Decode(&model); err != nil {
		return domain.ScreeningModel{}, "", fmt.Errorf("modelo de screening inválido")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return domain.ScreeningModel{}, "", fmt.Errorf("modelo de screening inválido")
	}
	if !usecase.IsValidOptInScreeningModel(model) {
		return domain.ScreeningModel{}, "", fmt.Errorf("modelo de screening inválido")
	}
	return model, hashBytes(data), nil
}

func evidenceInputHash(req domain.IdeaRunRequest, local []domain.RepositoryCandidate) (string, error) {
	return usecase.DossierEvidenceInputHash(req, local, cliDossierFiles{})
}

type checkpointArtifactHashes struct{ root string }

func (files checkpointArtifactHashes) ArtifactSHA256(path string) (string, error) {
	return fileHash(filepath.Join(files.root, path))
}
func reusableEvidenceCheckpoint(root string, previous domain.IdeaRunManifest, evidenceHash, evidenceRulesHash string, ruleVersions map[string]string, screeningModelHashes ...string) bool {
	modelHash := ""
	if len(screeningModelHashes) > 0 {
		modelHash = screeningModelHashes[0]
	}
	return usecase.CanReuseEvidenceCheckpoint(previous, usecase.EvidenceCheckpointCriteria{EvidenceHash: evidenceHash, EvidenceRulesHash: evidenceRulesHash, RuleVersions: ruleVersions, ScreeningModelHash: modelHash}, checkpointArtifactHashes{root: root})
}

func localCandidates(idea domain.IdeaSpec, paths []string) ([]domain.RepositoryCandidate, error) {
	return localCandidatesExcluding(idea, paths, "")
}

func localCandidatesExcluding(idea domain.IdeaSpec, paths []string, excludedPath string) ([]domain.RepositoryCandidate, error) {
	return localCandidatesExcludingContext(commandContext(context.Background()), idea, paths, excludedPath)
}

func localCandidatesExcludingContext(ctx context.Context, idea domain.IdeaSpec, paths []string, excludedPath string) ([]domain.RepositoryCandidate, error) {
	roots, err := normalizedLocalHashRoots(ctx, paths, excludedPath)
	if err != nil {
		return nil, err
	}
	budget := newTreeHashBudget()
	var out []domain.RepositoryCandidate
	for _, abs := range roots {
		h, e := treeHashExcludingContext(ctx, abs, excludedPath, budget)
		if e != nil {
			return nil, e
		}
		out = append(out, localHashCandidate(abs, h))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out, nil
}

func localHashCandidate(abs, digest string) domain.RepositoryCandidate {
	return usecase.LocalRepositoryCandidate(abs, filepath.Base(abs), digest)
}

func treeHash(root string) (string, error) {
	return treeHashExcluding(root, "")
}

func treeHashExcluding(root, excludedPath string) (string, error) {
	return treeHashExcludingContext(commandContext(context.Background()), root, excludedPath, newTreeHashBudget())
}
func writeDecisionMarkdown(path string, r domain.ReuseDecisionReport, evidence ...[]domain.EvidenceBundle) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Decisão de reúso: %s\n\n", r.Decision)
	fmt.Fprintf(&b, "Fonte canônica: [reuse-decision.json](reuse-decision.json)\n\n")
	if r.Primary != nil {
		fmt.Fprintf(&b, "## Ponto de partida\n\n%s/%s@%s\n\n", r.Primary.Identity.Owner, r.Primary.Identity.Repo, r.Primary.Commit)
	}
	fmt.Fprintln(&b, "## Razões")
	fmt.Fprintln(&b)
	for _, x := range r.Reasons {
		fmt.Fprintf(&b, "- `%s` %s: %s\n", x.RuleID, x.Candidate, x.Message)
	}
	fmt.Fprintln(&b, "\n## Capacidades ainda não comprovadas")
	fmt.Fprintln(&b)
	for _, x := range r.Gaps {
		fmt.Fprintf(&b, "- %s\n", x)
	}
	if r.Investigation != nil {
		fmt.Fprintf(&b, "\n## Plano de investigação\n\n%s\n\nEncerramento: `%s`. %s\n\n", markdownCell(r.Investigation.Scope), r.Investigation.StopReason, markdownCell(r.Investigation.TelemetryLabel))
		for _, action := range r.Investigation.Actions {
			fmt.Fprintf(&b, "- **%s** (`%s`) — %s/%s. Lacuna: %s Próxima ação: %s Aceitação: %s\n", markdownCell(action.Question), action.RequirementID, action.Classification, action.Status, markdownCell(action.Gap), markdownCell(action.NextAction), markdownCell(action.Acceptance))
			for _, evidence := range action.Evidence {
				fmt.Fprintf(&b, "  - Fonte: `%s` em `%s:%d` (%s).\n", evidence.Ref.Key(), markdownCell(evidence.Path), evidence.Line, evidence.Kind)
			}
		}
	}
	fmt.Fprintf(&b, "\nEsforço: %s. Método: %s. %s\n", r.Effort.Level, markdownCell(r.Effort.Method), markdownCell(r.Effort.Limitations))
	if len(evidence) > 0 {
		for i, bundle := range evidence[0] {
			if i >= 5 {
				fmt.Fprintf(&b, "\nMais %d candidatos em evidence/bundles.json.\n", len(evidence[0])-i)
				break
			}
			if bundle.ImplementationSearch == nil || len(bundle.ImplementationSearch.Requirements) == 0 {
				continue
			}
			fmt.Fprintf(&b, "\n## Recuperação em %s\n\n", markdownCell(bundle.Repository.Key()))
			report.WriteQueryInterpretation(&b, bundle.ImplementationSearch.QueryInterpretation)
			report.WriteRequirementRetrieval(&b, bundle.ImplementationSearch.Requirements)
			report.WriteActionAssessments(&b, bundle.ImplementationSearch.ActionAssessments)
		}
	}
	fmt.Fprintln(&b, "\n## Evidências\n\nConsulte [assessment-report.json](assessment-report.json) e [evidence/bundles.json](evidence/bundles.json).")
	return writeAtomic(path, []byte(b.String()), 0644)
}

func writeAlternativesMarkdown(path string, report domain.SoftwareAlternativeReport) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Alternativas para %s\n\n", report.Target.Name)
	fmt.Fprintln(&b, "Fonte canônica: [alternatives.json](alternatives.json)")
	if report.TopChoice != nil {
		fmt.Fprintf(&b, "\n## Primeira escolha\n\n`%s`\n", report.TopChoice.Key())
	}
	fmt.Fprintln(&b, "\n## Alternativas")
	fmt.Fprintln(&b)
	regional := alternativesHaveDiscoveryLocales(report)
	if regional {
		fmt.Fprintln(&b, "| Repositório | Fit | Cobertura | Estrelas | Licença | Idiomas/regiões de descoberta |")
		fmt.Fprintln(&b, "|---|---:|---:|---:|---|---|")
	} else {
		fmt.Fprintln(&b, "| Repositório | Fit | Cobertura | Estrelas | Licença |")
		fmt.Fprintln(&b, "|---|---:|---:|---:|---|")
	}
	for _, alternative := range report.Alternatives {
		if regional {
			fmt.Fprintf(&b, "| `%s` | %s | %.0f%% | %d | %s | %s |\n", alternative.Repository.Key(), alternative.Fit, alternative.Coverage*100, alternative.Stars, alternativeLicenseMarkdown(alternative), markdownCell(strings.Join(alternative.DiscoveryLocales, ", ")))
		} else {
			fmt.Fprintf(&b, "| `%s` | %s | %.0f%% | %d | %s |\n", alternative.Repository.Key(), alternative.Fit, alternative.Coverage*100, alternative.Stars, alternativeLicenseMarkdown(alternative))
		}
	}
	fmt.Fprintln(&b, "\n## Excluídos")
	fmt.Fprintln(&b)
	if regional {
		fmt.Fprintln(&b, "| Repositório | Fit | Razões e bloqueios | Idiomas/regiões de descoberta |")
		fmt.Fprintln(&b, "|---|---:|---|---|")
	} else {
		fmt.Fprintln(&b, "| Repositório | Fit | Razões e bloqueios |")
		fmt.Fprintln(&b, "|---|---:|---|")
	}
	for _, alternative := range report.Excluded {
		details := append(append([]string(nil), alternative.Blockers...), alternative.Reasons...)
		if regional {
			fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n", alternative.Repository.Key(), alternative.Fit, markdownCell(strings.Join(details, "; ")), markdownCell(strings.Join(alternative.DiscoveryLocales, ", ")))
		} else {
			fmt.Fprintf(&b, "| `%s` | %s | %s |\n", alternative.Repository.Key(), alternative.Fit, markdownCell(strings.Join(details, "; ")))
		}
	}
	return writeAtomic(path, []byte(b.String()), 0644)
}

func alternativesHaveDiscoveryLocales(report domain.SoftwareAlternativeReport) bool {
	for _, alternative := range report.Alternatives {
		if len(alternative.DiscoveryLocales) > 0 {
			return true
		}
	}
	for _, alternative := range report.Excluded {
		if len(alternative.DiscoveryLocales) > 0 {
			return true
		}
	}
	return false
}

func alternativeLicenseMarkdown(alternative domain.SoftwareAlternative) string {
	value := alternative.License.Identifier
	for _, reason := range alternative.Reasons {
		if reason != "licença não reconhecida; valide os termos antes da adoção" {
			continue
		}
		if strings.TrimSpace(value) != "" {
			value += " — "
		}
		value += "⚠ " + reason
		break
	}
	return markdownCell(value)
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	return strings.Join(strings.Fields(value), " ")
}

func buildRunStages(root string, result domain.IdeaRunResult, alternatives ...domain.SoftwareAlternativeReport) []domain.IdeaRunStage {
	return usecase.BuildRecommendationRunStages(checkpointArtifactHashes{root: root}, result, alternatives...)
}

type cliDossierFiles struct {
	output, replay string
	traceInputs    []string
	stdout         io.Writer
}

func (f cliDossierFiles) Prepare(resume bool) error { return prepareDossier(f.output, resume) }
func (f cliDossierFiles) ArtifactSHA256(path string) (string, error) {
	return fileHash(filepath.Join(f.output, path))
}
func (f cliDossierFiles) RepositoryIdentities(url string) []domain.RepositoryIdentity {
	return localGitHubRepositoryIdentities(url)
}
func (f cliDossierFiles) ReplayEvidenceIdentity(root string) (usecase.DossierReplayEvidenceIdentity, error) {
	directory, err := commandHost().Abs(root)
	if err != nil {
		return usecase.DossierReplayEvidenceIdentity{}, fmt.Errorf("diretório de replay inválido: %w", err)
	}
	snapshotHash, err := fileHash(filepath.Join(directory, "discovery-snapshot.json"))
	if err != nil {
		return usecase.DossierReplayEvidenceIdentity{}, fmt.Errorf("snapshot de replay inválido: %w", err)
	}
	bundlesHash, err := fileHash(filepath.Join(directory, "evidence", "bundles.json"))
	if err != nil {
		return usecase.DossierReplayEvidenceIdentity{}, fmt.Errorf("bundles de replay inválidos: %w", err)
	}
	return usecase.DossierReplayEvidenceIdentity{Directory: directory, SnapshotHash: snapshotHash, BundlesHash: bundlesHash}, nil
}
func (f cliDossierFiles) CaptureLocal(ctx context.Context, paths []string, metadata bool) (domain.CapturedWorkspace, error) {
	return commandHost().CaptureLocal(ctx, paths, f.output, metadata)
}
func (f cliDossierFiles) SnapshotProvider(fallback domain.ContentProvider, bindings []domain.SnapshotBinding) (domain.ContentProvider, error) {
	return buildSnapshotProvider(fallback, bindings)
}
func (f cliDossierFiles) HasArtifact(path string) bool {
	_, err := commandHost().Lstat(filepath.Join(f.output, path))
	return err == nil
}
func (f cliDossierFiles) ReadJSON(source usecase.DossierSource, path string, value any) error {
	root := f.replay
	if source == usecase.DossierOutputSource {
		root = f.output
	}
	return readJSONFile(filepath.Join(root, path), value)
}
func (f cliDossierFiles) WriteJSON(path string, value any) error {
	return writeJSONAtomic(filepath.Join(f.output, path), value)
}
func (f cliDossierFiles) WriteDecision(report domain.ReuseDecisionReport, bundles []domain.EvidenceBundle) error {
	return writeDecisionMarkdown(filepath.Join(f.output, "decision.md"), report, bundles)
}
func (f cliDossierFiles) WriteAlternatives(report domain.SoftwareAlternativeReport) error {
	return writeAlternativesMarkdown(filepath.Join(f.output, "alternatives.md"), report)
}
func (f cliDossierFiles) WriteTrace(value any, path, other string) error {
	return writeDecisionTraceJSON(value, filepath.Join(f.output, path), nil, append(append([]string(nil), f.traceInputs...), filepath.Join(f.output, other))...)
}

func (f cliDossierFiles) Announce(execution usecase.RecommendationDossierOutput) {
	if f.stdout == nil {
		return
	}
	_, _ = fmt.Fprintf(f.stdout, "dossiê %s: %s (%s)\n", execution.Manifest.RunID, execution.Result.Decision.Decision, f.output)
	writeResumeDiagnostic(f.stdout, execution.Checkpoint)
	if execution.StaleTrace {
		_, _ = fmt.Fprintln(f.stdout, "aviso: projeção anterior (trace-graph.json/trace-index.json) permanece no disco, mas não pertence à execução atual; consulte os estágios e hashes em run.json")
	}
}
