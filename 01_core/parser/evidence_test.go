/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/evidence_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser

import (
	"bateia/01_core/domain"
	"strings"
	"testing"
)

func TestOITRejectsCircularModuleAndPrefix(t *testing.T) {
	for _, name := range []string{"USalign", "bindiff"} {
		m := &domain.Manifest{ModuleName: name, Language: "cpp", Content: []byte("independently align protein structures"), ContentPath: "README.md"}
		r := EvaluateRelevance("order independent transparency opengl", m, domain.RepositoryCandidate{Repo: name}, nil)
		if r.Score >= ScoreReadmeKeyword {
			t.Fatalf("unrelated %s scored %+v", name, r)
		}
	}
	m := &domain.Manifest{Content: []byte("OpenGL transparency"), ContentPath: "README.md"}
	r := EvaluateRelevance("order independent transparency opengl", m, domain.RepositoryCandidate{}, nil)
	if r.Score != 32 || !strings.Contains(r.Reason, "transparency") || !strings.Contains(r.Reason, "opengl") {
		t.Fatalf("missing auditable partial match: %+v", r)
	}
	m.ContentPath = "CMakeLists.txt"
	if EvaluateRelevance("opengl", m, domain.RepositoryCandidate{}, nil).Label == LabelReadmeKeyword {
		t.Fatal("manifest mislabeled README")
	}
}
func TestCoinSelectionAndCoverage(t *testing.T) {
	tree := []domain.DirectoryNode{{Path: "README.md", Type: "file"}, {Path: "tmp/a.py", Type: "file"}, {Path: "tmp/b.py", Type: "file"}, {Path: "vendor/Transparency.cpp", Type: "file"}, {Path: "src/rendering/SoTransparency.cpp", Type: "file"}, {Path: "src/rendering/SoGLRenderAction.cpp", Type: "file"}}
	paths := SelectEvidenceFiles("Coin rendering transparency opengl", tree, 2)
	if strings.Join(paths, ",") != "README.md,src/rendering/SoTransparency.cpp,src/rendering/SoGLRenderAction.cpp" {
		t.Fatalf("selection: %v", paths)
	}
	content := strings.Repeat("// filler\n", 60) + "void render() { /* transparency */ }\n"
	repo := domain.ExtractedRepository{Files: []domain.ExtractedFile{{Path: paths[1], Content: content}}}
	AnnotateEvidence(&repo, "transparency", paths)
	ev := repo.Files[0].Evidence
	if ev.LastLine != 61 || len(ev.MatchLines) != 1 || ev.MatchLines[0] != 61 || ev.SHA256 == "" || repo.Coverage.Status != "insufficient" {
		t.Fatalf("coverage %+v", repo)
	}
	if len(repo.Coverage.MissingFiles) != 2 || repo.HealthScore != -1 {
		t.Fatal("missing coverage hidden")
	}
	old := repo.SnapshotHash
	repo.Files[0].Content += "changed"
	AnnotateEvidence(&repo, "transparency", paths)
	if old == repo.SnapshotHash {
		t.Fatal("snapshot ignores edits")
	}
}
func TestHealthAndLanguageEvidence(t *testing.T) {
	if CalculateHealthScore(nil) != -1 {
		t.Fatal("empty findings establish health")
	}
	for _, lang := range []string{"c", "cpp"} {
		vs := EvaluateArchViolations(&domain.Manifest{Language: lang}, []domain.DirectoryNode{{Path: "src/main.c", Type: "file"}})
		for _, v := range vs {
			if strings.Contains(v.Suggestion, "_test.go") || strings.Contains(v.CodeSnippet, "_test.go") {
				t.Fatal("Go convention leaked")
			}
		}
	}
	m := &domain.Manifest{Language: "go", Dependencies: []domain.Dependency{{Name: "go.opentelemetry.io/otel"}, {Name: "github.com/prometheus/client_golang"}}}
	for _, v := range EvaluateArchViolations(m, []domain.DirectoryNode{{Path: "main.go", Type: "file"}}) {
		if v.Rule == "MISSING_INFRASTRUCTURE_ISOLATION" {
			t.Fatalf("observability is not ORM: %+v", v)
		}
	}
	if !DependencyIdentityMatch("entgo.io/ent", "entgo.io/ent") || DependencyIdentityMatch("github.com/prometheus/client_golang", "ent") {
		t.Fatal("dependency identities")
	}
}
