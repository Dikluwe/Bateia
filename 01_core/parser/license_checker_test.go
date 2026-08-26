/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/license_checker_test.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import "testing"

func TestCategorizeLicense(t *testing.T) {
	mitText := "Permission is hereby granted, free of charge, to any person obtaining a copy of this software (MIT License)"
	if CategorizeLicense(mitText) != Permissive {
		t.Errorf("Esperava categoria PERMISSIVE para MIT, obteve %s", CategorizeLicense(mitText))
	}

	gplText := "GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007"
	if CategorizeLicense(gplText) != Copyleft {
		t.Errorf("Esperava categoria COPYLEFT para GPL, obteve %s", CategorizeLicense(gplText))
	}

	unknownText := "Custom Proprietary License"
	if CategorizeLicense(unknownText) != Unknown {
		t.Errorf("Esperava categoria UNKNOWN para licença desconhecida, obteve %s", CategorizeLicense(unknownText))
	}
}
