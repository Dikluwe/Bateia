/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ast_analyzer.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ASTSymbol representa um símbolo extraído de um arquivo fonte Go via AST nativa.
type ASTSymbol struct {
	Name string `json:"name"`
	Kind string `json:"kind"` // "struct", "interface", "function", "method"
	Line int    `json:"line"`
	Doc  string `json:"doc,omitempty"`
}

// ExtractGoSymbols analisa o conteúdo de código fonte Go em memória (Zero I/O) e retorna símbolos declarados.
func ExtractGoSymbols(filename string, code string) ([]ASTSymbol, error) {
	if strings.TrimSpace(code) == "" {
		return nil, nil
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, code, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("falha no parse da AST Go para %s: %w", filename, err)
	}

	var symbols []ASTSymbol

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					var kind string
					switch typeSpec.Type.(type) {
					case *ast.StructType:
						kind = "struct"
					case *ast.InterfaceType:
						kind = "interface"
					default:
						kind = "type"
					}

					pos := fset.Position(typeSpec.Pos())
					doc := ""
					if d.Doc != nil {
						doc = strings.TrimSpace(d.Doc.Text())
					}

					symbols = append(symbols, ASTSymbol{
						Name: typeSpec.Name.Name,
						Kind: kind,
						Line: pos.Line,
						Doc:  doc,
					})
				}
			}
		case *ast.FuncDecl:
			kind := "function"
			if d.Recv != nil && len(d.Recv.List) > 0 {
				kind = "method"
			}

			pos := fset.Position(d.Pos())
			doc := ""
			if d.Doc != nil {
				doc = strings.TrimSpace(d.Doc.Text())
			}

			symbols = append(symbols, ASTSymbol{
				Name: d.Name.Name,
				Kind: kind,
				Line: pos.Line,
				Doc:  doc,
			})
		}
	}

	return symbols, nil
}
