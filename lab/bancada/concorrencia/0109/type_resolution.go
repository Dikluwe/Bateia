/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/type_resolution.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"path"
	"sort"
)

type parsedSource struct {
	source SourceFile
	fset   *token.FileSet
	tree   *ast.File
}
type typeResolution struct {
	info        *types.Info
	complete    bool
	diagnostics []types.Error
}

type suppliedImporter struct {
	load func(string) (*types.Package, error)
}

func (s suppliedImporter) Import(p string) (*types.Package, error) { return s.load(p) }

func resolveTypes(files []parsedSource, fs *token.FileSet, broken map[string]bool, env *Environment) map[string]typeResolution {
	groups := map[string][]*ast.File{}
	for _, p := range files {
		key := path.Dir(p.source.Path) + "/" + p.tree.Name.Name
		groups[key] = append(groups[key], p.tree)
	}
	keys := []string{}
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := map[string]typeResolution{}
	bindings := map[string]string{}
	names := map[string]string{}
	var gc types.Importer
	if env != nil {
		for _, binding := range env.Packages {
			match := ""
			for _, key := range keys {
				if path.Dir(key) == binding.Directory {
					if match != "" {
						match = "ambiguous"
						break
					}
					match = key
				}
			}
			if match != "" && match != "ambiguous" {
				bindings[binding.ImportPath] = match
				names[match] = binding.ImportPath
			}
		}
		exports := map[string][]byte{}
		for _, e := range env.Exports {
			exports[e.ImportPath] = e.Data
		}
		gc = importer.ForCompiler(fs, "gc", func(p string) (io.ReadCloser, error) {
			if data, ok := exports[p]; ok {
				return io.NopCloser(bytes.NewReader(data)), nil
			}
			return nil, fmt.Errorf("export not supplied: %s", p)
		})
	}
	active := map[string]bool{}
	packages := map[string]*types.Package{}
	var check func(string) (*types.Package, error)
	var imports suppliedImporter
	imports.load = func(p string) (*types.Package, error) {
		if key, ok := bindings[p]; ok {
			return check(key)
		}
		if gc != nil {
			return gc.Import(p)
		}
		return nil, fmt.Errorf("dependency not supplied: %s", p)
	}
	check = func(key string) (*types.Package, error) {
		if r, ok := result[key]; ok {
			if r.complete {
				return packages[key], nil
			}
			return nil, fmt.Errorf("supplied package failed: %s", key)
		}
		if active[key] {
			return nil, fmt.Errorf("supplied import cycle: %s", key)
		}
		active[key] = true
		defer delete(active, key)
		info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Uses: map[*ast.Ident]types.Object{}, Defs: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}}
		failed := false
		diagnostics := []types.Error{}
		cfg := types.Config{Sizes: types.SizesFor("gc", "amd64"), Error: func(err error) {
			failed = true
			if e, ok := err.(types.Error); ok && len(diagnostics) < 100 {
				diagnostics = append(diagnostics, e)
			}
		}}
		if env != nil {
			cfg.Importer = imports
		}
		name := names[key]
		if name == "" {
			name = "supplied/" + key
		}
		pkg, err := cfg.Check(name, fs, groups[key], info)
		complete := err == nil && !failed && !broken[key] && !broken[path.Dir(key)+"/*"]
		result[key] = typeResolution{info: info, complete: complete, diagnostics: diagnostics}
		packages[key] = pkg
		if !complete {
			return nil, fmt.Errorf("supplied package failed: %s", key)
		}
		return pkg, nil
	}
	for _, key := range keys {
		_, _ = check(key)
	}
	return result
}

// knownOther suppresses a channel-specific gap for a successfully typed non-channel make.
func (r typeResolution) channelMake(call *ast.CallExpr, id *ast.Ident) (channel, knownOther bool) {
	if !r.complete || r.info == nil {
		return false, false
	}
	builtin, ok := r.info.Uses[id].(*types.Builtin)
	if !ok || builtin.Name() != "make" {
		return false, false
	}
	tv := r.info.Types[call.Args[0]]
	if !tv.IsType() || tv.Type == nil {
		return false, false
	}
	_, isChannel := tv.Type.Underlying().(*types.Chan)
	if !isChannel {
		return false, true
	}
	result := r.info.Types[call].Type
	if result == nil {
		return false, false
	}
	_, resultChannel := result.Underlying().(*types.Chan)
	return resultChannel, false
}
