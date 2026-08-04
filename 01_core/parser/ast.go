/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/ast.md
 * @layer L1
 * @updated 2026-08-03
 */

package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"bateia/01_core/domain"
)

// ExtractASTSymbols extrai os símbolos exported/estruturais de um conteúdo de arquivo de código em memória.
func ExtractASTSymbols(path string, content []byte) []domain.ASTSymbol {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return extractGoSymbols(path, content)
	case ".rs":
		return ExtractRustSymbols(path, content)
	default:
		return nil
	}
}

func extractGoSymbols(path string, content []byte) []domain.ASTSymbol {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil || node == nil {
		return nil
	}

	symbols := make([]domain.ASTSymbol, 0)

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						position := fset.Position(typeSpec.Pos())
						kind := "type"
						switch typeSpec.Type.(type) {
						case *ast.StructType:
							kind = "struct"
						case *ast.InterfaceType:
							kind = "interface"
						}
						symbols = append(symbols, domain.ASTSymbol{
							Name: typeSpec.Name.Name,
							Kind: kind,
							File: path,
							Line: position.Line,
						})
					}
				}
			}
		case *ast.FuncDecl:
			position := fset.Position(d.Pos())
			kind := "function"
			name := d.Name.Name
			if d.Recv != nil && len(d.Recv.List) > 0 {
				kind = "method"
			}
			symbols = append(symbols, domain.ASTSymbol{
				Name: name,
				Kind: kind,
				File: path,
				Line: position.Line,
			})
		}
	}

	return symbols
}
