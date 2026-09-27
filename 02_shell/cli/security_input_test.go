/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_input_test.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const securityMITLicense = `MIT License

Copyright (c) 2026 Audit

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`

func TestIdeaInputRejectsOversizedAndNonregularFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxIdeaInputBytes + 1); err != nil {
		t.Fatal(err)
	}
	f.Close()
	if _, err := readIdeaInput(path); err == nil {
		t.Fatal("oversized input accepted")
	}
	if _, err := readIdeaInput(dir); err == nil {
		t.Fatal("directory input accepted")
	}
	if err := os.WriteFile(path, []byte(`{"schema_version":"1"}`), 0600); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := readDecisionJSON(path, &out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Repeat("[", 65)+"0"+strings.Repeat("]", 65)), 0600); err != nil {
		t.Fatal(err)
	}
	if err := readDecisionJSON(path, &out); err == nil {
		t.Fatal("deep JSON accepted")
	}
	if err := validateIdeaJSONStructure([]byte(`{"text":"[[[[{}"}`)); err != nil {
		t.Fatal("quoted delimiters miscounted")
	}
}
