/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/clang/client.md
 * @layer L3
 * @updated 2026-10-09
 */
package clang

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

const maxUnits = 12
const maxOutput = 32 * 1024 * 1024
const timeout = 20 * time.Second
const maxDatabaseBytes = 16 * 1024 * 1024

type Client struct {
	binary   string
	flowTool string
}

func NewClient(binary string) *Client {
	if binary == "" {
		binary = "clang"
	}
	return &Client{binary: binary}
}

type compilation struct {
	Directory string   `json:"directory"`
	File      string   `json:"file"`
	Arguments []string `json:"arguments"`
	Command   string   `json:"command"`
}

func (c *Client) Analyze(ctx context.Context, req domain.SemanticRequest) (domain.SemanticReport, error) {
	report := domain.SemanticReport{Status: "not_evaluated", Reason: "Configuração de compilação indisponível; solução e comportamento não verificados.", Analyzer: c.binary, MaxUnits: maxUnits, TimeoutSeconds: int(timeout.Seconds()), MaxOutputBytes: maxOutput}
	if ctx.Err() != nil {
		return report, ctx.Err()
	}
	root, err := filepath.Abs(req.Root)
	if err != nil || req.Root == "" {
		return report, fmt.Errorf("raiz local inválida")
	}
	var database []byte
	var databasePath string
	for _, p := range []string{"compile_commands.json", "build/compile_commands.json", "../build/compile_commands.json"} {
		file := filepath.Join(root, p)
		data, e := readLimited(file, maxDatabaseBytes)
		if e == nil {
			database = data
			databasePath = file
			break
		}
		if !os.IsNotExist(e) {
			report.Reason = "Configuração de compilação não lida: " + file + ": " + e.Error()
			return report, nil
		}
	}
	if database == nil {
		return report, nil
	}
	report.ConfigurationHash = digest(database)
	var entries []compilation
	if err = json.Unmarshal(database, &entries); err != nil {
		report.Reason = "compile_commands.json inválido: " + err.Error()
		return report, nil
	}
	binary, err := exec.LookPath(c.binary)
	if err != nil {
		report.Reason = "Clang indisponível: " + err.Error()
		return report, nil
	}
	report.Analyzer = binary
	version, _, err := run(ctx, binary, root, []string{"--version"}, 64*1024)
	if err != nil {
		report.Reason = "Versão Clang indisponível: " + err.Error()
		return report, nil
	}
	report.Version = strings.TrimSpace(string(version))
	paths := map[string]bool{}
	for _, p := range req.Paths {
		if !filepath.IsAbs(p) {
			p = filepath.Join(root, p)
		}
		paths[filepath.Clean(p)] = true
	}
	var selected []compilation
	for _, entry := range entries {
		if !filepath.IsAbs(entry.Directory) {
			entry.Directory = filepath.Join(filepath.Dir(databasePath), entry.Directory)
		}
		if !filepath.IsAbs(entry.File) {
			entry.File = filepath.Join(entry.Directory, entry.File)
		}
		entry.File = filepath.Clean(entry.File)
		if !inside(root, entry.File) || !isSource(entry.File) {
			continue
		}
		report.TotalUnits++
		if len(paths) > 0 && !paths[entry.File] {
			continue
		}
		selected = append(selected, entry)
	}

	// Search source text within a fixed budget before choosing compilation units.
	// This is discovery evidence only; the compiler still validates the selected units.
	scores := map[string]int{}
	scannedBytes := 0
	for _, entry := range selected {
		score := 10 * pathScore(req.Query, entry.File)
		if scannedBytes < 64*1024*1024 {
			if content, e := readLimited(entry.File, 2*1024*1024); e == nil {
				scannedBytes += len(content)
				score += 2 * len(parser.MatchingQueryTerms(req.Query, string(content)))
			}
		}
		scores[entry.File] = score
	}
	sort.SliceStable(selected, func(i, j int) bool {
		a, b := scores[selected[i].File], scores[selected[j].File]
		if a != b {
			return a > b
		}
		return selected[i].File < selected[j].File
	})

	if len(selected) == 0 {
		report.Reason = "Nenhuma unidade C/C++ selecionada possui comando de compilação."
		return report, nil
	}
	report.Status = "complete"
	report.Reason = "Resolução estática apenas das unidades selecionadas; adequação à consulta, execução e saúde não verificadas. Chamadas virtuais/indiretas não comprovam destino em execução."
	if len(selected) > maxUnits {
		selected = selected[:maxUnits]
		report.Status = "partial"
		report.Reason += " Orçamento de unidades excedido."
	}
	report.SelectedUnits = len(selected)
	completed := 0
	for _, entry := range selected {
		if ctx.Err() != nil {
			report.Status = "partial"
			report.Reason += " Cancelado."
			return report, ctx.Err()
		}
		unit := c.analyzeUnit(ctx, binary, root, entry, databasePath, database)
		if unit.Status != "complete" {
			report.Status = "partial"
		} else {
			completed++
		}
		report.Units = append(report.Units, unit)
	}
	if completed == 0 {
		report.Status = "not_evaluated"
		report.Reason = "Nenhuma unidade concluiu a análise estática; configuração e inputs capturados não comprovam resolução. Consulte os diagnósticos das unidades."
	}
	return report, nil
}
func (c *Client) analyzeUnit(ctx context.Context, binary, root string, entry compilation, dbPath string, db []byte) domain.SemanticUnit {
	unit := domain.SemanticUnit{Path: entry.File, Status: "not_evaluated", Inputs: []domain.SemanticInput{{Path: dbPath, SHA256: digest(db)}}}
	flags, err := safeFlags(entry)
	if err != nil {
		unit.Diagnostics = err.Error()
		return unit
	}
	// Preflight the real compiler dependency closure, including macro-only and system headers.
	depsArgs := append(append([]string{}, flags...), "-M", "-MT", "bateia", entry.File)
	deps, diagnostics, err := run(ctx, binary, entry.Directory, depsArgs, 4*1024*1024)
	unit.Arguments = append(append([]string{}, flags...), "-fsyntax-only", "-Xclang", "-ast-dump=json", entry.File)
	if err != nil {
		unit.Diagnostics = string(diagnostics) + "\n" + err.Error()
		return unit
	}
	files, err := dependencyPaths(string(deps), entry.Directory)
	if err != nil {
		unit.Diagnostics = err.Error()
		return unit
	}
	files = append(files, entry.File)
	data := map[string][]byte{}
	totalInputBytes := 0
	for _, p := range files {
		if _, ok := data[p]; ok {
			continue
		}
		content, e := readLimited(p, 16*1024*1024)
		if e != nil {
			unit.Diagnostics = "Entrada não capturada: " + p + ": " + e.Error()
			return unit
		}
		totalInputBytes += len(content)
		if totalInputBytes > 128*1024*1024 {
			unit.Diagnostics = "Limite total de entradas excedido"
			return unit
		}
		data[p] = content
		unit.Inputs = append(unit.Inputs, domain.SemanticInput{Path: p, SHA256: digest(content)})
	}
	ast, diagnostics, err := run(ctx, binary, entry.Directory, unit.Arguments, maxOutput)
	unit.Diagnostics = string(diagnostics)
	if err != nil {
		unit.Diagnostics += "\n" + err.Error()
		return unit
	}
	unit.Flow = c.analyzeFlow(ctx, entry, flags)
	for p, before := range data {
		after, e := readLimited(p, 16*1024*1024)
		if e != nil || digest(before) != digest(after) {
			unit.Diagnostics += "\nEntrada mudou durante análise: " + p
			if unit.Flow != nil {
				unit.Flow.Status = "not_evaluated"
				unit.Flow.Reason = "Entradas mudaram durante análise"
				unit.Flow.Functions = nil
				unit.Flow.Skipped = nil
			}
			return unit
		}
	}
	after, e := readLimited(dbPath, maxDatabaseBytes)
	if e != nil || digest(after) != digest(db) {
		unit.Diagnostics += "\nConfiguração mudou durante análise"
		if unit.Flow != nil {
			unit.Flow.Status = "not_evaluated"
			unit.Flow.Reason = "Configuração mudou durante análise"
			unit.Flow.Functions = nil
			unit.Flow.Skipped = nil
		}
		return unit
	}
	symbols, calls, err := decodeAST(ast, root, entry.File, data)
	if err != nil {
		unit.Diagnostics += "\nAST inválida: " + err.Error()
		return unit
	}
	unit.Symbols, unit.Calls = symbols, calls
	sort.Slice(unit.Inputs, func(i, j int) bool { return unit.Inputs[i].Path < unit.Inputs[j].Path })
	unit.Status = "complete"
	return unit
}
func pathScore(query, path string) int {
	score := 0
	for _, term := range parser.QueryTerms(query) {
		if strings.Contains(strings.ToLower(path), term) {
			score++
		}
	}
	return score
}
func isSource(p string) bool {
	switch strings.ToLower(filepath.Ext(p)) {
	case ".c", ".cc", ".cpp", ".cxx":
		return true
	}
	return false
}
func inside(root, p string) bool {
	rel, err := filepath.Rel(root, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
func digest(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }
func readLimited(p string, limit int64) ([]byte, error) {
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if len(b) > int(limit) {
		return nil, fmt.Errorf("limite de bytes excedido")
	}
	return b, e
}

type bounded struct {
	bytes.Buffer
	limit    int
	exceeded bool
}

// Hide bytes.Buffer.ReadFrom: os/exec's io.Copy would otherwise bypass Write's limit.
func (b *bounded) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(struct{ io.Writer }{b}, r)
}

func (b *bounded) Write(p []byte) (int, error) {
	if b.Len()+len(p) > b.limit {
		b.exceeded = true
		return 0, fmt.Errorf("limite de saída excedido")
	}
	return b.Buffer.Write(p)
}
func run(parent context.Context, binary, dir string, args []string, limit int) ([]byte, []byte, error) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = dir
	cmd.WaitDelay = time.Second
	// Avoid ambient compiler include injection and arbitrary compiler/tool wrappers.
	cmd.Env = []string{"PATH=/usr/bin:/bin", "LANG=C", "LC_ALL=C"}
	out := &bounded{limit: limit}
	stderr := &bounded{limit: 128 * 1024}
	cmd.Stdout = out
	cmd.Stderr = stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if out.exceeded || stderr.exceeded {
		err = errors.New("limite de saída excedido")
	}
	return out.Bytes(), stderr.Bytes(), err
}

// safeFlags accepts literal compilation context, never repository compiler execution or action flags.
func safeFlags(entry compilation) ([]string, error) {
	if len(entry.Arguments) == 0 {
		var err error
		entry.Arguments, err = splitCompilationCommand(entry.Command)
		if err != nil {
			return nil, err
		}
	}
	args := entry.Arguments[1:]
	var out []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.ContainsRune(a, 0) || strings.HasPrefix(a, "@") {
			return nil, fmt.Errorf("opção não permitida: %q", a)
		}
		if a == "-c" {
			continue
		}
		if a == "-o" || a == "-MF" || a == "-MT" || a == "-MQ" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("argumento ausente: %s", a)
			}
			i++
			continue
		}
		if a == "-MD" || a == "-MMD" || a == "-MP" {
			continue
		}
		if a == "-I" || a == "-isystem" || a == "-iquote" || a == "-include" || a == "-D" || a == "-U" || a == "-x" || a == "-isysroot" {
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "@") || strings.HasPrefix(args[i+1], "-") {
				return nil, fmt.Errorf("argumento inválido: %s", a)
			}
			out = append(out, a, args[i+1])
			i++
			continue
		}
		if strings.HasPrefix(a, "-I") || strings.HasPrefix(a, "-D") || strings.HasPrefix(a, "-U") || strings.HasPrefix(a, "-std=") || strings.HasPrefix(a, "--sysroot=") {
			out = append(out, a)
			continue
		}
		switch a {
		case "-mtls-dialect=gnu2", "-flifetime-dse=1":
			continue // GCC code generation only; no compiler action or AST context.
		case "-pthread", "-fPIC", "-fpic", "-fPIE", "-fpie", "-fno-exceptions", "-fexceptions", "-fno-rtti", "-frtti", "-m32", "-m64", "-ansi", "-fvisibility=hidden", "-fdiagnostics-color=always", "-fno-math-errno", "-fno-trapping-math", "-fno-common", "-ffunction-sections", "-fdata-sections", "-msse4.1", "-fno-lto":
			out = append(out, a)
			continue
		}
		// Ordinary diagnostics/optimization/debug flags do not select executable compiler extensions.
		if a == "-g" || a == "-g0" || a == "-g1" || a == "-g2" || a == "-g3" || a == "-O0" || a == "-O1" || a == "-O2" || a == "-O3" || a == "-Os" || a == "-DNDEBUG" || strings.HasPrefix(a, "-W") && !strings.HasPrefix(a, "-Wl,") && !strings.HasPrefix(a, "-Wp,") {
			out = append(out, a)
			continue
		}
		p := a
		if !filepath.IsAbs(p) {
			p = filepath.Join(entry.Directory, p)
		}
		if filepath.Clean(p) == filepath.Clean(entry.File) {
			continue
		}
		return nil, fmt.Errorf("opção de compilação não suportada: %q", a)
	}
	return out, nil
}

// splitCompilationCommand decodes a literal argv, never a shell program.
func splitCompilationCommand(command string) ([]string, error) {
	var args []string
	var token strings.Builder
	var quote rune
	started, escaped := false, false
	for _, r := range command {
		if r == 0 || r == '\n' || r == '\r' {
			return nil, errors.New("command contém controle não permitido")
		}
		if escaped {
			// Within double quotes POSIX only escapes these characters.
			if quote == '"' && !strings.ContainsRune("$`\"\\", r) {
				token.WriteRune('\\')
			}
			token.WriteRune(r)
			escaped = false
			continue
		}
		if quote == '\'' {
			if r == quote {
				quote = 0
			} else {
				token.WriteRune(r)
			}
			continue
		}
		if r == '\\' {
			escaped = true
			started = true
			continue
		}
		if r == '$' || r == '`' || quote == 0 && strings.ContainsRune(";|&<>()*?~#", r) {
			return nil, errors.New("command contém sintaxe de shell não permitida")
		}
		if quote == '"' {
			if r == quote {
				quote = 0
			} else {
				token.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			started = true
			continue
		}
		if r == ' ' || r == '\t' {
			if started {
				args = append(args, token.String())
				token.Reset()
				started = false
			}
			continue
		}
		started = true
		token.WriteRune(r)
	}
	if quote != 0 || escaped {
		return nil, errors.New("command com aspas ou escape incompleto")
	}
	if started {
		args = append(args, token.String())
	}
	if len(args) == 0 || args[0] == "" || strings.Contains(args[0], "=") {
		return nil, errors.New("arguments ausente e command literal inválido")
	}
	return args, nil
}
func dependencyPaths(text, dir string) ([]string, error) {
	colon := strings.Index(text, ":")
	if colon < 0 {
		return nil, errors.New("dependências Clang sem alvo")
	}
	text = strings.ReplaceAll(text[colon+1:], "\\\n", "")
	var paths []string
	var token strings.Builder
	escaped := false
	flush := func() {
		if token.Len() > 0 {
			p := token.String()
			if !filepath.IsAbs(p) {
				p = filepath.Join(dir, p)
			}
			paths = append(paths, filepath.Clean(p))
			token.Reset()
		}
	}
	for _, r := range text {
		if escaped {
			token.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			flush()
		} else {
			token.WriteRune(r)
		}
	}
	flush()
	if len(paths) == 0 {
		return nil, errors.New("dependências vazias")
	}
	return paths, nil
}
