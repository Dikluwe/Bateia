/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/truncate.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncateSmart_ShortContent(t *testing.T) {
	content := "module github.com/test/app\n\ngo 1.22\n"
	result, truncated := TruncateSmart(content, 1500, 40)

	if truncated {
		t.Error("esperava truncated == false para conteúdo curto")
	}
	if result != content {
		t.Errorf("esperava %q, obteve %q", content, result)
	}
}

func TestTruncateSmart_BreakAtNewline(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 50; i++ {
		sb.WriteString(fmt.Sprintf("linha %d: require github.com/pkg/dep%d v1.0.0\n", i, i))
	}
	content := sb.String()

	result, truncated := TruncateSmart(content, 500, 20)

	if !truncated {
		t.Error("esperava truncated == true")
	}
	if !strings.HasSuffix(result, "\n") {
		t.Error("esperava que o resultado terminasse com quebra de linha \\n")
	}
	if strings.HasSuffix(result, "require github.com/pkg/de") {
		t.Error("não deve truncar no meio de palavras")
	}
}

func TestTruncateSmart_PreservesUTF8(t *testing.T) {
	content := "require (\n    github.com/user/pkg v1.0.0 // Teste de comentário: 🚀 日本語\n    github.com/another/pkg v2.0.0\n)\n"

	result, _ := TruncateSmart(content, 70, 10)

	if !utf8.ValidString(result) {
		t.Error("o resultado do truncamento deve ser uma string UTF-8 válida")
	}
}

func TestGenerateGitHubFileURL(t *testing.T) {
	url := GenerateGitHubFileURL("werbot", "werbot", "HEAD", "go.mod")
	expected := "https://github.com/werbot/werbot/blob/main/go.mod"
	if url != expected {
		t.Errorf("esperava %q, obteve %q", expected, url)
	}

	urlSub := GenerateGitHubFileURL("owner", "repo", "v1.0.0", "internal/api (v1)/app.go")
	expectedSub := "https://github.com/owner/repo/blob/v1.0.0/internal/api%20%28v1%29/app.go"
	if urlSub != expectedSub {
		t.Errorf("esperava %q, obteve %q", expectedSub, urlSub)
	}
}
