/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/ast.md
 * @layer L3
 * @updated 2026-09-07
 */
package clang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"bateia/01_core/domain"
)

type location struct {
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Offset    *int      `json:"offset"`
	Spelling  *location `json:"spellingLoc"`
	Expansion *location `json:"expansionLoc"`
}
type astNode struct {
	ID                 string   `json:"id"`
	Kind               string   `json:"kind"`
	Name               string   `json:"name"`
	Mangled            string   `json:"mangledName"`
	Implicit           bool     `json:"isImplicit"`
	CompleteDefinition bool     `json:"completeDefinition"`
	Loc                location `json:"loc"`
	Range              struct {
		Begin location `json:"begin"`
		End   location `json:"end"`
	} `json:"range"`
	Type struct {
		Qual string `json:"qualType"`
	} `json:"type"`
	Ref    *astNode   `json:"referencedDecl"`
	Member string     `json:"referencedMemberDecl"`
	Inner  []*astNode `json:"inner"`
}

func flat(l location) location {
	if l.Expansion != nil {
		return flat(*l.Expansion)
	}
	if l.Spelling != nil {
		return flat(*l.Spelling)
	}
	return l
}
func lineAt(l location, file string, data map[string][]byte) int {
	l = flat(l)
	if l.Line > 0 {
		return l.Line
	}
	if l.Offset != nil {
		b, ok := data[file]
		if ok && *l.Offset >= 0 && *l.Offset <= len(b) {
			return bytes.Count(b[:*l.Offset], []byte("\n")) + 1
		}
	}
	return 0
}
func declaration(kind string) bool {
	switch kind {
	case "FunctionDecl", "CXXMethodDecl", "CXXConstructorDecl", "CXXDestructorDecl", "CXXConversionDecl", "RecordDecl", "CXXRecordDecl", "EnumDecl", "TypedefDecl", "TypeAliasDecl", "VarDecl", "FieldDecl":
		return true
	}
	return false
}
func function(kind string) bool {
	switch kind {
	case "FunctionDecl", "CXXMethodDecl", "CXXConstructorDecl", "CXXDestructorDecl", "CXXConversionDecl":
		return true
	}
	return false
}
func definition(n *astNode) bool {
	if n.CompleteDefinition {
		return true
	}
	for _, child := range n.Inner {
		if (function(n.Kind) && (child.Kind == "CompoundStmt" || child.Kind == "CXXTryStmt")) || (n.Kind == "EnumDecl" && child.Kind == "EnumConstantDecl") {
			return true
		}
	}
	return false
}

func decodeAST(raw []byte, root, main string, data map[string][]byte) ([]domain.SemanticSymbol, []domain.SemanticCall, error) {
	var tree astNode
	if err := json.Unmarshal(raw, &tree); err != nil {
		return nil, nil, err
	}
	if tree.Kind != "TranslationUnitDecl" {
		return nil, nil, fmt.Errorf("raiz AST não é unidade de tradução")
	}
	nodes := map[string]*astNode{}
	files := map[*astNode]string{}
	// Clang elides repeated file names across adjacent declarations in its JSON printer.
	last := main
	var index func(*astNode, string)
	index = func(n *astNode, inherited string) {
		loc := flat(n.Loc)
		begin := flat(n.Range.Begin)
		explicit := loc.File
		if explicit == "" {
			explicit = begin.File
		}
		file := inherited
		if explicit != "" {
			file = explicit
			last = explicit
		} else if file == "" {
			file = last
		}
		if !filepath.IsAbs(file) {
			file = filepath.Join(root, file)
		}
		files[n] = file
		nodes[n.ID] = n
		for _, child := range n.Inner {
			childFile := file
			if n.Kind == "TranslationUnitDecl" {
				childFile = ""
			}
			index(child, childFile)
		}
	}
	index(&tree, main)
	var symbols []domain.SemanticSymbol
	var walk func(*astNode)
	walk = func(n *astNode) {
		file := files[n]
		line := lineAt(n.Loc, file, data)
		if line == 0 {
			line = lineAt(n.Range.Begin, file, data)
		}
		if declaration(n.Kind) && !n.Implicit && n.Name != "" && inside(root, file) && line > 0 {
			end := lineAt(n.Range.End, file, data)
			if end < line {
				end = line
			}
			id := n.Mangled
			if id == "" {
				id = fmt.Sprintf("%s:%d:%s:%s", file, line, n.Kind, n.Name)
			}
			symbols = append(symbols, domain.SemanticSymbol{ID: id, Name: n.Name, Kind: n.Kind, File: file, Type: n.Type.Qual, Line: line, EndLine: end, Definition: definition(n)})
		}
		for _, child := range n.Inner {
			walk(child)
		}
	}
	walk(&tree)
	var calls []domain.SemanticCall
	stableID := func(n *astNode) string {
		if n == nil {
			return ""
		}
		if n.Mangled != "" {
			return n.Mangled
		}
		return fmt.Sprintf("%s:%d:%s:%s", files[n], lineAt(n.Loc, files[n], data), n.Kind, n.Name)
	}
	var target func(*astNode) *astNode
	target = func(n *astNode) *astNode {
		if n.Ref != nil && function(n.Ref.Kind) {
			if full := nodes[n.Ref.ID]; full != nil {
				return full
			}
			return n.Ref
		}
		if n.Member != "" {
			if decl := nodes[n.Member]; decl != nil && function(decl.Kind) {
				return decl
			}
		}
		for _, child := range n.Inner {
			if name := target(child); name != nil {
				return name
			}
		}
		return nil
	}
	var callWalk func(*astNode, *astNode)
	callWalk = func(n *astNode, caller *astNode) {
		if function(n.Kind) {
			caller = n
		}
		if (n.Kind == "CallExpr" || n.Kind == "CXXMemberCallExpr" || n.Kind == "CXXOperatorCallExpr") && inside(root, files[n]) && len(n.Inner) > 0 {
			callee := target(n.Inner[0])
			line := lineAt(n.Range.Begin, files[n], data)
			if callee != nil && caller != nil && line > 0 {
				calls = append(calls, domain.SemanticCall{Caller: caller.Name, Callee: callee.Name, CallerID: stableID(caller), CalleeID: stableID(callee), File: files[n], Line: line})
			}
		}
		for _, child := range n.Inner {
			callWalk(child, caller)
		}
	}
	callWalk(&tree, nil)
	sort.SliceStable(symbols, func(i, j int) bool {
		if symbols[i].File != symbols[j].File {
			return symbols[i].File < symbols[j].File
		}
		return symbols[i].Line < symbols[j].Line
	})
	return symbols, calls, nil
}
