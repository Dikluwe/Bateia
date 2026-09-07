/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/evidence.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"fmt"
	"path"
	"sort"
	"strings"
	"unicode"
)

const EvidenceVersion = "evidence-v5-flow"

func QueryTerms(query string) []string {
	seen := map[string]bool{}
	var terms []string
	for _, term := range strings.FieldsFunc(strings.ToLower(NormalizeQuery(query)), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		if len(term) >= 3 && !seen[term] {
			terms = append(terms, term)
			seen[term] = true
		}
	}
	return terms
}
func MatchingQueryTerms(query, content string) []string {
	words := map[string]bool{}
	for _, word := range strings.FieldsFunc(strings.ToLower(content), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		words[word] = true
	}
	var matched []string
	for _, term := range QueryTerms(query) {
		if words[term] {
			matched = append(matched, term)
		}
	}
	return matched
}
func sourcePath(p string) bool {
	switch strings.ToLower(path.Ext(p)) {
	case ".go", ".rs", ".ts", ".tsx", ".js", ".py", ".java", ".c", ".cc", ".cpp", ".cxx", ".h", ".hh", ".hpp", ".hxx", ".glsl", ".vert", ".frag":
		return true
	}
	return false
}

// SelectEvidenceFiles is a deterministic path ranking, not a semantic implementation locator.
func SelectEvidenceFiles(query string, tree []domain.DirectoryNode, limit int) []string {
	type ranked struct {
		p     string
		score int
	}
	var code []ranked
	var selected []string
	seen := map[string]bool{}
	for _, n := range tree {
		if n.Type != "file" && n.Type != "blob" {
			continue
		}
		p := strings.ReplaceAll(n.Path, "\\", "/")
		if seen[p] {
			continue
		}
		seen[p] = true
		if !strings.Contains(p, "/") && (DetectManifestType(p) != "" || strings.EqualFold(p, "README.md")) {
			selected = append(selected, p)
			continue
		}
		if !sourcePath(p) {
			continue
		}
		score := 0
		lower := strings.ToLower(p)
		for _, term := range QueryTerms(query) {
			if strings.Contains(lower, term) {
				score += 20
			}
		}
		for _, dir := range []string{"tmp", "temp", "vendor", "third_party", "external", "generated", "build", "dist", "node_modules"} {
			if strings.Contains("/"+lower, "/"+dir+"/") && !strings.Contains(strings.ToLower(query), dir) {
				score -= 200
			}
		}
		if strings.Contains(lower, "test") {
			score -= 10
		}
		code = append(code, ranked{p, score})
	}
	sort.Strings(selected)
	sort.Slice(code, func(i, j int) bool {
		if code[i].score != code[j].score {
			return code[i].score > code[j].score
		}
		return code[i].p < code[j].p
	})
	for i, n := range code {
		if i >= limit {
			break
		}
		if len(QueryTerms(query)) > 0 && n.score <= 0 {
			continue
		}
		selected = append(selected, n.p)
	}
	return selected
}
func AnnotateEvidence(repo *domain.ExtractedRepository, query string, requested []string) {
	repo.Coverage = domain.AnalysisCoverage{Status: "insufficient", Reason: "Amostra limitada por ranking de caminhos; relação semântica da implementação com a consulta não verificada.", RequestedFiles: requested, SupportedLanguages: []string{"go (AST)", "rust (símbolos heurísticos)"}, ChecksNotRun: []string{"execução de testes e benchmarks", "validação semântica da solução", "AST de linguagens sem provedor; Clang reportado separadamente", "isolamento arquitetural por fluxo de chamadas"}}
	present := map[string]bool{}
	var hashes []string
	for i := range repo.Files {
		f := &repo.Files[i]
		present[f.Path] = true
		ev := domain.FileEvidence{SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(f.Content))), State: "content_inspected", Method: "busca lexical; não comprova implementação", MatchedTerms: MatchingQueryTerms(query, f.Content)}
		if f.Content != "" {
			ev.FirstLine = 1
			ev.LastLine = strings.Count(strings.TrimSuffix(f.Content, "\n"), "\n") + 1
		}
		for line, text := range strings.Split(f.Content, "\n") {
			if len(MatchingQueryTerms(query, text)) > 0 {
				ev.MatchLines = append(ev.MatchLines, line+1)
			}
		}
		if path.Ext(f.Path) == ".go" || path.Ext(f.Path) == ".rs" {
			ev.Method += "; tentativa de extração de símbolos"
		}
		f.Evidence = ev
		hashes = append(hashes, f.Path+"\x00"+ev.SHA256)
	}
	for _, p := range requested {
		if !present[p] {
			repo.Coverage.MissingFiles = append(repo.Coverage.MissingFiles, p)
		}
	}
	sort.Strings(hashes)
	repo.SnapshotHash = fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(hashes, "\n"))))
	repo.Coverage.ChecksRun = []string{"busca lexical no conteúdo extraído"}
	if repo.Manifest != nil {
		repo.Coverage.ChecksRun = append(repo.Coverage.ChecksRun, "heurísticas sobre manifesto e árvore parcial")
	}
	// A penalty tally cannot establish architectural health, even with no findings.
	repo.HealthScore = -1
}

// DependencyIdentityMatch rejects substrings such as ent in opentelemetry/client.
func DependencyIdentityMatch(name, pattern string) bool {
	name, pattern = strings.ToLower(name), strings.ToLower(pattern)
	if name == pattern || strings.HasPrefix(name, pattern+"/") {
		return true
	}
	for _, segment := range strings.Split(name, "/") {
		if segment == pattern {
			return true
		}
	}
	return false
}
