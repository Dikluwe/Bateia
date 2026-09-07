/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/ast_test.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang

import "testing"

func TestEnumDefinition0036(t *testing.T) {
	if !definition(&astNode{Kind: "EnumDecl", Inner: []*astNode{{Kind: "EnumConstantDecl"}}}) {
		t.Fatal("enum body not a definition")
	}
	if definition(&astNode{Kind: "EnumDecl"}) {
		t.Fatal("forward enum became definition")
	}
	if definition(&astNode{Kind: "VarDecl", Inner: []*astNode{{Kind: "EnumConstantDecl"}}}) {
		t.Fatal("wrong kind became definition")
	}
}
