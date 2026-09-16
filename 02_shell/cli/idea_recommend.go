/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_recommend.md
 * @layer L2
 * @updated 2026-09-16
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"bateia/02_shell/report"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type repeatedPaths []string

func (r *repeatedPaths) String() string     { return strings.Join(*r, ",") }
func (r *repeatedPaths) Set(v string) error { *r = append(*r, v); return nil }

func RunIdeaRecommendCommand(args []string, stdout io.Writer, semantic usecase.SemanticExpander, discovery domain.ObservableDiscoveryProvider, content domain.ContentProvider, recommendationServices ...*usecase.RecommendationService) error {
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
	maxPerQuery := fs.Int("max-per-query", 5, "máximo de resultados por expansão")
	maxCandidates := fs.Int("max-candidates", 5, "máximo de candidatos para coleta")
	targetName := fs.String("target-name", "", "nome do software alvo; usa idea.problem quando omitido")
	targetRepo := fs.String("target-repo", "", "repositório alvo owner/repo ou URL GitHub")
	timeout := fs.Duration("timeout", 2*time.Minute, "limite total da execução")
	var locals repeatedPaths
	var searchLocales repeatedPaths
	fs.Var(&locals, "local", "worktree local; repetível")
	fs.Var(&searchLocales, "search-locale", "idioma natural/região BCP 47 para descoberta; repetível ou all")
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
	catalogCoverage, err := usecase.AnalyzeEvidenceCatalogCoverage(idea, evidenceRules)
	if err != nil {
		return err
	}
	if err := prepareDossier(*output, *resume); err != nil {
		return err
	}
	req := domain.IdeaRunRequest{Idea: idea, Policy: policy, Mode: runMode, Strategy: strings.ToLower(*strategy), Provider: *provider, Language: *language, SearchLocales: normalizedSearchLocales, LocalPaths: append([]string(nil), locals...), ReplayDirectory: *replay, OutputDirectory: *output, Resume: *resume, MaxPerQuery: *maxPerQuery, MaxCandidates: *maxCandidates}
	deps := usecase.IdeaRunDependencies{Semantic: semantic, Discovery: discovery, Content: content, EvidenceRules: evidenceRules}
	if runMode == domain.IdeaRunLocal {
		deps.LocalCandidates, err = localCandidatesExcluding(idea, locals, *output)
		if err != nil {
			return err
		}
	}
	evidenceHash, err := evidenceInputHash(req, deps.LocalCandidates)
	if err != nil {
		return err
	}
	inputHashes := map[string]string{"idea": hashJSON(idea), "policy": hashJSON(policy), "mode": hashBytes([]byte(runMode)), "evidence": evidenceHash, "evidence_rules": catalogCoverage.RulesHash, "target": hashJSON(target)}
	if len(normalizedSearchLocales) > 0 {
		inputHashes["search_locales"] = hashJSON(normalizedSearchLocales)
	}
	ruleVersions := map[string]string{"expansion": "1", "screening": "1", "evidence": catalogCoverage.RulesVersion, "decision": "0043", "alternatives": "0043", "analyzer": "0043"}
	if runMode == domain.IdeaRunLocal || runMode == domain.IdeaRunLive {
		ruleVersions["implementation_search"] = "0061"
	}
	if runMode == domain.IdeaRunLive {
		ruleVersions["discovery_query"] = "live0056-v1"
	}
	var previous domain.IdeaRunManifest
	canReuseEvidence := false
	if *resume {
		if e := readJSONFile(filepath.Join(*output, "run.json"), &previous); e == nil {
			canReuseEvidence = reusableEvidenceCheckpoint(*output, previous, evidenceHash, catalogCoverage.RulesHash, ruleVersions)
		}
	}
	runIdentity := struct {
		EvidenceHash string            `json:"evidence_hash"`
		PolicyHash   string            `json:"policy_hash"`
		TargetHash   string            `json:"target_hash"`
		RuleVersions map[string]string `json:"rule_versions"`
	}{EvidenceHash: evidenceHash, PolicyHash: inputHashes["policy"], TargetHash: inputHashes["target"], RuleVersions: ruleVersions}
	manifest := domain.IdeaRunManifest{SchemaVersion: "1", RunID: usecase.StableRunID(hashJSON(runIdentity)), Mode: runMode, Status: domain.IdeaRunRunning, InputHashes: inputHashes, RuleVersions: ruleVersions, Stages: []domain.IdeaRunStage{}}
	if err = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest); err != nil {
		return err
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var result domain.IdeaRunResult
	var alternatives domain.SoftwareAlternativeReport
	var runErr error
	usesRecommendationService := runMode == domain.IdeaRunLive && targetRepository == nil && !canReuseEvidence
	if usesRecommendationService {
		recommendationService, constructorErr := usecase.NewRecommendationServiceWithEvidenceRules(semantic, discovery, content, evidenceRules)
		if constructorErr != nil {
			return constructorErr
		}
		if strings.TrimSpace(*rulesPath) == "" && len(recommendationServices) > 0 && recommendationServices[0] != nil {
			recommendationService = recommendationServices[0]
		}
		execution, recommendErr := recommendationService.RecommendWithRun(ctx, domain.RecommendationRequest{
			Idea: idea, Policy: policy, Target: target,
			Search: domain.RecommendationSearch{Strategy: strings.ToLower(*strategy), Language: *language, Locales: normalizedSearchLocales, MaxPerQuery: *maxPerQuery, MaxCandidates: *maxCandidates},
		})
		result, alternatives, runErr = execution.Run, execution.Recommendation.Alternatives, recommendErr
	} else {
		result, runErr = usecase.RunIdeaPipeline(ctx, req, deps)
	}
	if runErr != nil {
		manifest.Status = domain.IdeaRunFailed
		manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "pipeline", Message: runErr.Error()}}
		_ = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest)
		return runErr
	}
	if !usesRecommendationService {
		alternativesStarted := time.Now()
		alternativeDiscovery, alternativeAssessments := alternativeInputsForTarget(runMode, result.Discovery, result.Assessment.Assessments, targetRepository)
		alternatives, err = usecase.RankSoftwareAlternatives(target, idea, alternativeDiscovery, alternativeAssessments, policy)
		result.StageDurations["alternatives"] = time.Since(alternativesStarted).Milliseconds()
		if err != nil {
			manifest.Status = domain.IdeaRunFailed
			manifest.Failures = []domain.ExecutionFailure{{Scope: "run", Stage: "alternatives", Message: err.Error()}}
			_ = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest)
			return err
		}
	}
	artifacts := []struct {
		path  string
		value any
	}{{"idea.json", idea}, {"reuse-policy.json", policy}, {"expansion-plan.json", result.Expansion}, {"discovery-snapshot.json", result.Discovery}, {"screening-report.json", result.Screening}, {filepath.Join("evidence", "catalog-coverage.json"), result.CatalogCoverage}, {filepath.Join("evidence", "plans.json"), result.Plans}, {filepath.Join("evidence", "bundles.json"), result.Bundles}, {"assessment-report.json", result.Assessment}, {"reuse-decision.json", result.Decision}, {"alternatives.json", alternatives}}
	for _, a := range artifacts {
		if err = writeJSONAtomic(filepath.Join(*output, a.path), a.value); err != nil {
			return err
		}
	}
	if err = writeDecisionMarkdown(filepath.Join(*output, "decision.md"), result.Decision, result.Bundles); err != nil {
		return err
	}
	if err = writeAlternativesMarkdown(filepath.Join(*output, "alternatives.md"), alternatives); err != nil {
		return err
	}
	manifest.Status = usecase.IdeaRunStatusFor(result)
	manifest.DecisionPath = "reuse-decision.json"
	manifest.Failures = append(manifest.Failures, result.Failures...)
	manifest.Stages = buildRunStages(*output, result, alternatives)
	if err = writeJSONAtomic(filepath.Join(*output, "run.json"), manifest); err != nil {
		return err
	}
	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "dossiê %s: %s (%s)\n", manifest.RunID, result.Decision.Decision, *output)
	}
	return nil
}

func parseAlternativeTargetRepository(value string) (*domain.RepositoryIdentity, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	repository := value
	lower := strings.ToLower(value)
	switch {
	case strings.HasPrefix(lower, "git@github.com:"):
		repository = value[len("git@github.com:"):]
	case strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "ssh://"):
		parsed, err := url.Parse(value)
		if err != nil || !strings.EqualFold(parsed.Hostname(), "github.com") || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("target-repo inválido %q; use owner/repo ou URL GitHub", value)
		}
		repository = parsed.Path
	}
	repository = strings.Trim(strings.TrimSpace(repository), "/")
	if strings.HasSuffix(strings.ToLower(repository), ".git") {
		repository = repository[:len(repository)-4]
	}
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || invalidTargetRepositoryPart(parts[0]) || invalidTargetRepositoryPart(parts[1]) {
		return nil, fmt.Errorf("target-repo inválido %q; use owner/repo ou URL GitHub", value)
	}
	return &domain.RepositoryIdentity{Owner: strings.ToLower(parts[0]), Repo: strings.ToLower(parts[1])}, nil
}

type localAlternativeBinding struct {
	repository *domain.RepositoryIdentity
}

func alternativeInputsForTarget(mode domain.IdeaRunMode, discovery domain.DiscoverySnapshot, assessments []domain.CandidateAssessment, target *domain.RepositoryIdentity) (domain.DiscoverySnapshot, []domain.CandidateAssessment) {
	discovery.Candidates = append([]domain.DiscoveredCandidate(nil), discovery.Candidates...)
	assessments = append([]domain.CandidateAssessment(nil), assessments...)
	if mode != domain.IdeaRunLocal || target == nil {
		return discovery, assessments
	}
	bindings := map[string]localAlternativeBinding{}
	for i := range discovery.Candidates {
		candidate := discovery.Candidates[i].Repository
		key := alternativeIdentityKey(candidate.Owner, candidate.Repo)
		identities := localGitHubRepositoryIdentities(candidate.URL)
		if len(identities) == 0 {
			discovery.Candidates[i].Repository.Stale = true
			continue
		}
		for _, identity := range identities {
			if alternativeIdentityKey(identity.Owner, identity.Repo) != alternativeIdentityKey(target.Owner, target.Repo) {
				continue
			}
			bound := domain.RepositoryIdentity{Owner: strings.ToLower(strings.TrimSpace(target.Owner)), Repo: strings.ToLower(strings.TrimSpace(target.Repo))}
			discovery.Candidates[i].Repository.Owner = bound.Owner
			discovery.Candidates[i].Repository.Repo = bound.Repo
			bindings[key] = localAlternativeBinding{repository: &bound}
			break
		}
	}
	for i := range assessments {
		key := alternativeIdentityKey(assessments[i].Repository.Owner, assessments[i].Repository.Repo)
		if binding := bindings[key]; binding.repository != nil {
			assessments[i] = rebindAlternativeAssessment(assessments[i], *binding.repository)
		}
	}
	return discovery, assessments
}

func alternativeIdentityKey(owner, repo string) string {
	return strings.ToLower(strings.TrimSpace(owner)) + "/" + strings.ToLower(strings.TrimSpace(repo))
}

func localGitHubRepositoryIdentities(candidateURL string) []domain.RepositoryIdentity {
	parsed, err := url.Parse(candidateURL)
	if err != nil || !strings.EqualFold(parsed.Scheme, "file") || parsed.Path == "" {
		return nil
	}
	configPath := localGitConfigPath(filepath.FromSlash(parsed.Path))
	if configPath == "" {
		return nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}
	seen := map[string]domain.RepositoryIdentity{}
	remoteSection := false
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			remoteSection = strings.HasPrefix(strings.ToLower(line), "[remote ")
			continue
		}
		if !remoteSection {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || !strings.EqualFold(strings.TrimSpace(key), "url") {
			continue
		}
		identity, parseErr := parseAlternativeTargetRepository(strings.TrimSpace(value))
		if parseErr != nil || identity == nil {
			continue
		}
		seen[alternativeIdentityKey(identity.Owner, identity.Repo)] = *identity
	}
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]domain.RepositoryIdentity, 0, len(keys))
	for _, key := range keys {
		out = append(out, seen[key])
	}
	return out
}

func localGitConfigPath(root string) string {
	dotGit := filepath.Join(root, ".git")
	info, err := os.Stat(dotGit)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		return filepath.Join(dotGit, "config")
	}
	data, err := os.ReadFile(dotGit)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(strings.ToLower(line), "gitdir:") {
		return ""
	}
	gitDir := strings.TrimSpace(line[len("gitdir:"):])
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	configRoot := filepath.Clean(gitDir)
	if common, readErr := os.ReadFile(filepath.Join(configRoot, "commondir")); readErr == nil {
		commonDir := strings.TrimSpace(string(common))
		if !filepath.IsAbs(commonDir) {
			commonDir = filepath.Join(configRoot, commonDir)
		}
		configRoot = filepath.Clean(commonDir)
	}
	return filepath.Join(configRoot, "config")
}

func rebindAlternativeAssessment(assessment domain.CandidateAssessment, repository domain.RepositoryIdentity) domain.CandidateAssessment {
	assessment.Repository.Owner = repository.Owner
	assessment.Repository.Repo = repository.Repo
	assessment.ScreeningDecision.Repository = assessment.Repository.Key()
	assessment.ProblemAlignment.Sources = rebindAlternativeSources(assessment.ProblemAlignment.Sources, repository)
	assessment.Capabilities = append([]domain.CapabilityEvidence(nil), assessment.Capabilities...)
	for i := range assessment.Capabilities {
		assessment.Capabilities[i].Sources = rebindAlternativeSources(assessment.Capabilities[i].Sources, repository)
	}
	assessment.License.Sources = rebindAlternativeSources(assessment.License.Sources, repository)
	assessment.Maintenance.Sources = rebindAlternativeSources(assessment.Maintenance.Sources, repository)
	assessment.Operation.Sources = rebindAlternativeSources(assessment.Operation.Sources, repository)
	return assessment
}

func rebindAlternativeSources(sources []domain.EvidenceSource, repository domain.RepositoryIdentity) []domain.EvidenceSource {
	out := append([]domain.EvidenceSource(nil), sources...)
	for i := range out {
		out[i].Ref.Owner = repository.Owner
		out[i].Ref.Repo = repository.Repo
	}
	return out
}

func invalidTargetRepositoryPart(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || strings.ContainsAny(value, " \t\r\n@:#?")
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

type evidenceInputs struct {
	Idea            domain.IdeaSpec      `json:"idea"`
	Mode            domain.IdeaRunMode   `json:"mode"`
	Strategy        string               `json:"strategy"`
	Provider        string               `json:"provider,omitempty"`
	Language        string               `json:"language,omitempty"`
	SearchLocales   []string             `json:"search_locales,omitempty"`
	MaxPerQuery     int                  `json:"max_per_query"`
	MaxCandidates   int                  `json:"max_candidates"`
	LocalCandidates []localEvidenceInput `json:"local_candidates,omitempty"`
	Replay          *replayEvidenceInput `json:"replay,omitempty"`
}

type localEvidenceInput struct {
	URL      string `json:"url"`
	CommitID string `json:"commit_id"`
}

type replayEvidenceInput struct {
	Directory    string `json:"directory"`
	SnapshotHash string `json:"snapshot_hash"`
	BundlesHash  string `json:"bundles_hash"`
}

func evidenceInputHash(req domain.IdeaRunRequest, local []domain.RepositoryCandidate) (string, error) {
	inputs := evidenceInputs{Idea: req.Idea, Mode: req.Mode, Strategy: req.Strategy, Provider: req.Provider, Language: req.Language, SearchLocales: append([]string(nil), req.SearchLocales...), MaxPerQuery: req.MaxPerQuery, MaxCandidates: req.MaxCandidates}
	for _, candidate := range local {
		inputs.LocalCandidates = append(inputs.LocalCandidates, localEvidenceInput{URL: candidate.URL, CommitID: candidate.CommitID})
	}
	if req.Mode == domain.IdeaRunReplay {
		directory, err := filepath.Abs(req.ReplayDirectory)
		if err != nil {
			return "", fmt.Errorf("diretório de replay inválido: %w", err)
		}
		snapshotHash, err := fileHash(filepath.Join(directory, "discovery-snapshot.json"))
		if err != nil {
			return "", fmt.Errorf("snapshot de replay inválido: %w", err)
		}
		bundlesHash, err := fileHash(filepath.Join(directory, "evidence", "bundles.json"))
		if err != nil {
			return "", fmt.Errorf("bundles de replay inválidos: %w", err)
		}
		inputs.Replay = &replayEvidenceInput{Directory: directory, SnapshotHash: snapshotHash, BundlesHash: bundlesHash}
	}
	return hashJSON(inputs), nil
}

func reusableEvidenceCheckpoint(root string, previous domain.IdeaRunManifest, evidenceHash, evidenceRulesHash string, ruleVersions map[string]string) bool {
	switch previous.Status {
	case domain.IdeaRunComplete, domain.IdeaRunPartial, domain.IdeaRunInconclusive:
	default:
		return false
	}
	if previous.SchemaVersion != "1" || previous.InputHashes["evidence"] != evidenceHash || previous.InputHashes["evidence_rules"] != evidenceRulesHash {
		return false
	}
	for _, key := range []string{"expansion", "screening", "evidence", "analyzer", "implementation_search", "discovery_query"} {
		if previous.RuleVersions[key] != ruleVersions[key] {
			return false
		}
	}
	required := map[string]bool{"discovery-snapshot.json": false, filepath.Join("evidence", "bundles.json"): false}
	for _, stage := range previous.Stages {
		for _, artifact := range stage.Outputs {
			expected, ok := required[artifact.Path]
			if !ok || expected || artifact.Hash == "" {
				continue
			}
			actual, err := fileHash(filepath.Join(root, artifact.Path))
			if err != nil || actual != artifact.Hash {
				return false
			}
			required[artifact.Path] = true
		}
	}
	for _, valid := range required {
		if !valid {
			return false
		}
	}
	return true
}

func localCandidates(idea domain.IdeaSpec, paths []string) ([]domain.RepositoryCandidate, error) {
	return localCandidatesExcluding(idea, paths, "")
}

func localCandidatesExcluding(idea domain.IdeaSpec, paths []string, excludedPath string) ([]domain.RepositoryCandidate, error) {
	var out []domain.RepositoryCandidate
	for _, path := range paths {
		abs, e := filepath.Abs(path)
		if e != nil {
			return nil, e
		}
		h, e := treeHashExcluding(abs, excludedPath)
		if e != nil {
			return nil, e
		}
		base := strings.ToLower(filepath.Base(abs))
		pathHash := hashBytes([]byte(abs))
		repo := base + "-" + pathHash[:8]
		out = append(out, domain.RepositoryCandidate{Owner: "local", Repo: repo, URL: "file://" + abs, CommitID: h})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out, nil
}
func treeHash(root string) (string, error) {
	return treeHashExcluding(root, "")
}

func treeHashExcluding(root, excludedPath string) (string, error) {
	h := sha256.New()
	excluded := ""
	if strings.TrimSpace(excludedPath) != "" {
		var err error
		excluded, err = filepath.Abs(excludedPath)
		if err != nil {
			return "", fmt.Errorf("caminho de saída inválido: %w", err)
		}
	}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "target" || d.Name() == "node_modules") {
			if path != root {
				return filepath.SkipDir
			}
		}
		if excluded != "" && path == excluded {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
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

func buildRunStages(root string, r domain.IdeaRunResult, alternatives ...domain.SoftwareAlternativeReport) []domain.IdeaRunStage {
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
	if len(alternatives) > 0 {
		jsonHash, _ := fileHash(filepath.Join(root, "alternatives.json"))
		markdownHash, _ := fileHash(filepath.Join(root, "alternatives.md"))
		assessmentHash, _ := fileHash(filepath.Join(root, "assessment-report.json"))
		discoveryHash, _ := fileHash(filepath.Join(root, "discovery-snapshot.json"))
		policyHash, _ := fileHash(filepath.Join(root, "reuse-policy.json"))
		out = append(out, domain.IdeaRunStage{
			Name:            "alternatives",
			ContractVersion: "1",
			Status:          domain.IdeaStageComplete,
			InputHashes: map[string]string{
				"assessment-report.json":  assessmentHash,
				"discovery-snapshot.json": discoveryHash,
				"reuse-policy.json":       policyHash,
				"target":                  hashJSON(alternatives[0].Target),
			},
			Outputs:        []domain.IdeaRunArtifact{{Path: "alternatives.json", Hash: jsonHash}, {Path: "alternatives.md", Hash: markdownHash}},
			Counts:         map[string]int64{"items": int64(len(alternatives[0].Alternatives) + len(alternatives[0].Excluded))},
			DurationMillis: r.StageDurations["alternatives"],
		})
	}
	return out
}
