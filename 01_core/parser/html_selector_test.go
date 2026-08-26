/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/html_selector_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestExtractHTMLTags(t *testing.T) {
	html := `
<html>
	<head><title>Bateia Documentation</title></head>
	<body>
		<h1>Title 1</h1>
		<p>First paragraph</p>
		<h1>Title 2</h1>
	</body>
</html>
`

	titles := ExtractHTMLTags(html, "h1")
	if len(titles) != 2 {
		t.Fatalf("Esperava 2 tags <h1>, obteve %d", len(titles))
	}

	if titles[0] != "Title 1" || titles[1] != "Title 2" {
		t.Errorf("Extração de h1 incorreta: %v", titles)
	}

	paras := ExtractHTMLTags(html, "p")
	if len(paras) != 1 || paras[0] != "First paragraph" {
		t.Errorf("Extração de p incorreta: %v", paras)
	}
}
