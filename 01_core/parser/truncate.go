/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/truncate.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

const (
	DefaultMaxChars = 1500
	DefaultMaxLines = 40
)

// ShouldTruncate verifica se o conteúdo excede os limites de caracteres (runes) ou linhas (OR lógico).
func ShouldTruncate(content string, maxChars, maxLines int) bool {
	if maxChars > 0 && utf8.RuneCountInString(content) > maxChars {
		return true
	}
	if maxLines > 0 {
		lineCount := strings.Count(content, "\n") + 1
		if lineCount > maxLines {
			return true
		}
	}
	return false
}

// TruncateSmart corta o conteúdo ao final de uma linha completa respeitando runes UTF-8 e limites de linhas/caracteres.
func TruncateSmart(content string, maxChars, maxLines int) (string, bool) {
	if !ShouldTruncate(content, maxChars, maxLines) {
		return content, false
	}

	var (
		runes     []rune
		totalR    int
		lineCount int
	)

	for _, r := range content {
		if maxChars > 0 && totalR >= maxChars {
			break
		}
		if maxLines > 0 && lineCount >= maxLines {
			break
		}

		runes = append(runes, r)
		totalR++

		if r == '\n' {
			lineCount++
		}
	}

	result := string(runes)

	// Se não terminou exatamente em \n e há quebras de linha no resultado, recuar até o último \n
	if !strings.HasSuffix(result, "\n") {
		lastNewline := strings.LastIndex(result, "\n")
		if lastNewline != -1 {
			result = result[:lastNewline+1]
		}
	}

	return result, true
}

// SanitizeCommitRef garante que a referência a commit seja válida para URLs no GitHub.
func SanitizeCommitRef(commit string) string {
	c := strings.TrimSpace(commit)
	if c == "" || c == "HEAD" {
		return "main"
	}
	return c
}

// GenerateGitHubFileURL produz uma URL válida e escapada para o arquivo no repositório GitHub.
func GenerateGitHubFileURL(owner, repo, commit, filePath string) string {
	cleanPath := strings.TrimPrefix(filePath, "/")
	pathParts := strings.Split(cleanPath, "/")
	escapedParts := make([]string, len(pathParts))
	for i, p := range pathParts {
		escapedParts[i] = url.PathEscape(p)
	}
	escapedPath := strings.Join(escapedParts, "/")

	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s",
		owner, repo, SanitizeCommitRef(commit), escapedPath)
}
