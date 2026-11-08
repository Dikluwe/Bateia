/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/domain/command_files_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package domain

import "testing"

func TestPortableFileOptionsAreIndependent(t *testing.T) {
	var seen FileOpenFlags
	for _, f := range []FileOpenFlags{OpenWriteOnly, OpenReadWrite, OpenAppend, OpenCreate, OpenExclusive, OpenTruncate, OpenNonblock, OpenNoFollow} {
		if f == 0 || seen&f != 0 {
			t.Fatal("overlapping file options")
		}
		seen |= f
	}
	if OpenReadOnly != 0 {
		t.Fatal("read-only base changed")
	}
}
