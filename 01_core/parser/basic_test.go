/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/basic_test.md
 * @layer L1
 * @updated 2026-07-31
 */

package parser

import (
	"testing"

	"bateia/01_core/domain"
)

func TestParseBasicVbp_Basic(t *testing.T) {
	content := `Type=Exe
Reference=*\G{00020430-0000-0000-C000-000000000466}#2.0#0#C:\Windows\SysWOW64\stdole2.tlb#OLE Automation
Reference=*\G{00000205-0000-0010-8000-00AA006D2EA4}#2.5#0#C:\Program Files\Common Files\system\ado\msado25.tlb#Microsoft ActiveX Data Objects 2.5 Library
Object={831F9261-6C0C-11D2-A9FC-0000F8754DA1}#2.0#0; MSCOMCTL.OCX`

	ref := domain.CommitRef{Owner: "legacy-vb", Repo: "VB6App"}
	manifest, err := ParseBasicVbp(ref, content)
	if err != nil {
		t.Fatalf("ParseBasicVbp falhou: %v", err)
	}

	if manifest.Language != "basic" {
		t.Errorf("Esperado Language 'basic', obtido '%s'", manifest.Language)
	}

	if len(manifest.Dependencies) != 2 {
		t.Fatalf("Esperado 2 referências OLE/ADO, obtido %d", len(manifest.Dependencies))
	}
}
