/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/query_admission_test.md
 * @layer L1
 * @updated 2026-09-08
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"bateia/01_core/usecase"
	"strings"
	"testing"
)

func runAdmission52(query, code string) domain.ImplementationSearchReport {
	return parser.RankGoImplementations(query, []domain.ExtractedFile{{Path: "unit.go", Content: "package neutral\n" + code + "\n"}}, 12)
}
func codeHits52(r domain.ImplementationSearchReport) []domain.ImplementationHit {
	var out []domain.ImplementationHit
	for _, h := range r.Hits {
		if h.DocumentationLine == 0 {
			out = append(out, h)
		}
	}
	return out
}
func TestAdmission0052ActionAndContext(t *testing.T) {
	cases := []struct {
		name, query, code string
		want              bool
	}{
		{"secondary only", "descomprimir dados bloco", "func Data(block []byte) []byte { return block }", false},
		{"ambiguous no context", "descomprimir conteúdo", "func Decode(src []byte) []byte { return src }", false},
		{"context unrequested", "descomprimir conteúdo", "func Decode(block []byte) []byte { return block }", false},
		{"requested block", "descomprimir bloco", "func Decode(block []byte) []byte { return block }", true},
		{"requested corrupt", "descomprimir corrompido", "func Decode(corrupt []byte) []byte { return corrupt }", true},
		{"wrong context", "descomprimir corrompido", "func Decode(block []byte) []byte { return block }", false},
		{"literal translation", "descomprimir conteúdo", "func Decompress(src []byte) []byte { return src }", true},
		{"compression ambiguous", "comprimir dados", "func Encode(data []byte) []byte { return data }", false},
		{"compression contextual", "comprimir bloco", "func Encode(block []byte) []byte { return block }", true},
		{"compression literal", "comprimir conteúdo", "func Compress(src []byte) []byte { return src }", true},
		{"serialization alternative", "serializar conteúdo", "func Marshal(src []byte) []byte { return src }", true},
		{"deserialization alternative", "desserializar conteúdo", "func Unmarshal(src []byte) []byte { return src }", true},
		{"sorting literal", "ordenar conteúdo", "func Sort(src []byte) []byte { return src }", true},
		{"comparison literal", "comparar conteúdo", "func Compare(src []byte) []byte { return src }", true},
		{"all families required", "ordenar serializar dados", "func Sort(data []byte) []byte { return data }", false},
		{"all families satisfied", "ordenar serializar dados", "func Sort(data []byte) []byte { marshal := data; return marshal }", true},
		{"no action expansion", "dados conteúdo", "func Data(src []byte) []byte { return src }", true},
		{"exact symbol exception", "`Decode` descomprimir conteúdo", "func Decode(src []byte) []byte { return src }", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := runAdmission52(tc.query, tc.code)
			hits := codeHits52(r)
			if (len(hits) > 0) != tc.want {
				t.Fatalf("want code hit %v, got %#v", tc.want, r.Hits)
			}
			if tc.want && tc.name != "no action expansion" && tc.name != "exact symbol exception" {
				reason := strings.ToLower(r.Reason)
				for _, h := range hits {
					reason += " " + strings.ToLower(h.Reason)
				}
				if !strings.Contains(reason, "filtro") && !strings.Contains(reason, "filter") && !strings.Contains(reason, "admiss") {
					t.Errorf("filter not explicit in reasons: %s", reason)
				}
			}
		})
	}
}
func TestAdmission0052DocumentAndCodeSeparate(t *testing.T) {
	cases := []struct {
		name, query, code string
		codeWant, docWant bool
	}{
		{"context only doc", "descomprimir bloco", "// Decode processes a block.\nfunc Decode(src []byte) []byte { return src }", false, true},
		{"context only body", "descomprimir bloco", "// Decode processes input.\nfunc Decode(block []byte) []byte { return block }", true, false},
		{"translation only doc", "descomprimir conteúdo", "// Decompress input.\nfunc Work(src []byte) []byte { return src }", false, true},
		{"secondary only doc", "descomprimir dados bloco", "// Data block.\nfunc Work(src []byte) []byte { return src }", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := runAdmission52(tc.query, tc.code)
			codeHit, docHit := false, false
			for _, h := range r.Hits {
				if h.DocumentationLine > 0 {
					docHit = true
				} else {
					codeHit = true
				}
			}
			if codeHit != tc.codeWant || docHit != tc.docWant {
				t.Fatalf("want code/doc %v/%v got %v/%v hits=%#v", tc.codeWant, tc.docWant, codeHit, docHit, r.Hits)
			}
		})
	}
}
func TestAdmission0052MultilineAndRequirementScope(t *testing.T) {
	query := "descomprimir dados\nordenar conteúdo"
	files := []domain.ExtractedFile{{Path: "unit.go", Content: "package neutral\nfunc Sort(data []byte) []byte { return data }\n"}}
	r := parser.RankGoImplementations(query, files, 12)
	hits := codeHits52(r)
	if len(hits) == 0 {
		t.Fatal("rejected line blocked admissible sort line")
	}
	for _, h := range hits {
		if has52(h.MatchedTerms, "data") {
			t.Fatalf("rejected line leaked data into matched terms: %#v", h)
		}
	}
	r = usecase.ValidateImplementationRequirements(query, r, files)
	seen := false
	for _, req := range r.Requirements {
		if req.Requirement == "descomprimir dados" {
			seen = true
			if len(req.Sources) > 0 || len(req.DocumentationSources) > 0 || req.Status == "evidence_located" {
				t.Fatalf("rejected action became requirement evidence: %#v", req)
			}
		}
	}
	if !seen {
		t.Fatal("missing rejected requirement projection")
	}
}
