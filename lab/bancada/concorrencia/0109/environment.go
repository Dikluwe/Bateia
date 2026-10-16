/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/environment.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
)

type PackageBinding struct {
	ImportPath string `json:"import_path"`
	Directory  string `json:"directory"`
}
type ExportData struct {
	ImportPath string `json:"import_path"`
	Origin     string `json:"origin"`
	Data       []byte `json:"data"`
	SHA256     string `json:"sha256"`
}
type Environment struct {
	SchemaVersion string           `json:"schema_version"`
	Toolchain     string           `json:"toolchain"`
	GOOS          string           `json:"goos"`
	GOARCH        string           `json:"goarch"`
	Packages      []PackageBinding `json:"packages"`
	Exports       []ExportData     `json:"exports"`
}

func canonicalEnvironment(input Environment) (Environment, string, error) {
	if input.SchemaVersion != "1" || !safeText(input.Toolchain, 100) || input.GOOS != "linux" || input.GOARCH != "amd64" || len(input.Packages) > 32 || len(input.Exports) > 128 {
		return Environment{}, "", fmt.Errorf("invalid supplied environment")
	}
	input.Packages = append([]PackageBinding(nil), input.Packages...)
	input.Exports = append([]ExportData(nil), input.Exports...)
	seen, dirs := map[string]bool{}, map[string]bool{}
	for _, p := range input.Packages {
		if !validImportPath(p.ImportPath) || seen[p.ImportPath] || dirs[p.Directory] || !safeText(p.Directory, 300) || p.Directory == ".." || path.IsAbs(p.Directory) || path.Clean(p.Directory) != p.Directory || strings.HasPrefix(p.Directory, "../") || strings.Contains(p.Directory, "\\") {
			return Environment{}, "", fmt.Errorf("invalid package binding")
		}
		seen[p.ImportPath] = true
		dirs[p.Directory] = true
	}
	total := 0
	for i, p := range input.Exports {
		total += len(p.Data)
		if !validImportPath(p.ImportPath) || seen[p.ImportPath] || !safeText(p.Origin, 500) || len(p.Data) == 0 || total > 32<<20 || p.SHA256 != Digest(string(p.Data)) {
			return Environment{}, "", fmt.Errorf("invalid export data")
		}
		seen[p.ImportPath] = true
		input.Exports[i].Data = append([]byte(nil), p.Data...)
	}
	sort.Slice(input.Packages, func(i, j int) bool { return input.Packages[i].ImportPath < input.Packages[j].ImportPath })
	sort.Slice(input.Exports, func(i, j int) bool { return input.Exports[i].ImportPath < input.Exports[j].ImportPath })
	data, err := json.Marshal(input)
	if err != nil || len(data) > 64<<20 {
		return Environment{}, "", fmt.Errorf("environment byte limit")
	}
	return input, Digest(string(data)), nil
}
func validImportPath(s string) bool {
	return safeText(s, 300) && s != "." && s != ".." && s != "unsafe" && path.Clean(s) == s && !path.IsAbs(s) && !strings.HasPrefix(s, "../") && !strings.ContainsAny(s, "\\ :")
}
