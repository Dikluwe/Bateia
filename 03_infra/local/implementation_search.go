/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/implementation_search.md
 * @layer L3
 * @updated 2026-09-08
 */
package local

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	"syscall"

	"bateia/01_core/domain"
	"bateia/01_core/parser"
)

func (c *Client) SearchImplementations(ctx context.Context, candidate domain.RepositoryCandidate, query string, budget domain.ImplementationSearchBudget) (domain.ImplementationSearchReport, error) {
	if budget.MaxFiles <= 0 {
		budget.MaxFiles = 2048
	}
	if budget.MaxBytes <= 0 {
		budget.MaxBytes = 32 * 1024 * 1024
	}
	if budget.MaxFileBytes <= 0 {
		budget.MaxFileBytes = 512 * 1024
	}
	if budget.MaxResults <= 0 {
		budget.MaxResults = 24
	}
	report := domain.ImplementationSearchReport{Status: "partial", Method: "go_ast_query_v1", Budget: budget, Hits: []domain.ImplementationHit{}, Examined: []domain.ImplementationInput{}}
	root, err := candidateRoot(candidate)
	if err != nil {
		report.Issues = []string{err.Error()}
		return report, err
	}
	scoped, err := os.OpenRoot(root)
	if err != nil {
		report.Issues = []string{err.Error()}
		return report, err
	}
	defer scoped.Close()
	var files []domain.ExtractedFile
	var examined []domain.ImplementationInput
	var issues []string
	limited := false
	var bytesRead int64
	attempted := 0
	err = fs.WalkDir(scoped.FS(), ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if walkErr != nil {
			issues = append(issues, fmt.Sprintf("varredura %s: %v", name, walkErr))
			return nil
		}
		if name == "." {
			return nil
		}
		if !parser.ImplementationSearchPathAllowed(name, query) {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		// Never follow aliases: confinement alone would permit aliases into protected corpus.
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			return nil
		}
		if attempted >= budget.MaxFiles || bytesRead >= budget.MaxBytes {
			limited = true
			issues = append(issues, "limite de varredura atingido")
			return fs.SkipAll
		}
		attempted++
		f, openErr := scoped.OpenFile(name, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
		if openErr != nil {
			issues = append(issues, fmt.Sprintf("leitura %s: %v", name, openErr))
			return nil
		}
		defer f.Close()
		info, statErr := f.Stat()
		if statErr != nil {
			issues = append(issues, fmt.Sprintf("stat %s: %v", name, statErr))
			return nil
		}
		if !info.Mode().IsRegular() {
			issues = append(issues, "arquivo não regular: "+name)
			return nil
		}
		cap := budget.MaxFileBytes
		if remaining := budget.MaxBytes - bytesRead; remaining < cap {
			cap = remaining
		}
		if info.Size() > cap {
			limited = true
			issues = append(issues, "limite de bytes impede leitura: "+name)
			return nil
		}
		data, readErr := io.ReadAll(io.LimitReader(f, cap))
		bytesRead += int64(len(data))
		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		examined = append(examined, domain.ImplementationInput{Path: name, SHA256: hash, Bytes: int64(len(data))})
		after, afterErr := f.Stat()
		if readErr != nil || afterErr != nil {
			issues = append(issues, fmt.Sprintf("leitura incompleta %s: read=%v stat=%v", name, readErr, afterErr))
			return nil
		}
		if after.Size() != int64(len(data)) || !after.ModTime().Equal(info.ModTime()) {
			limited = true
			issues = append(issues, "arquivo alterado ou truncado durante leitura: "+name)
			return nil
		}
		files = append(files, domain.ExtractedFile{Path: name, Content: string(data), Size: int64(len(data))})
		return nil
	})
	report = parser.RankGoImplementations(query, files, budget.MaxResults)
	report.Budget = budget
	report.Examined = examined
	report.FilesRead = len(examined)
	report.BytesRead = bytesRead
	report.Limited = report.Limited || limited
	report.Issues = append(report.Issues, issues...)
	if err != nil {
		report.Issues = append(report.Issues, err.Error())
	}
	sort.Strings(report.Issues)
	if len(report.Issues) > 0 || report.Limited {
		report.Status = "partial"
	}
	return report, err
}
