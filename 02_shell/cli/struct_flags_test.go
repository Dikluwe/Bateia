/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/struct_flags_test.md
 * @layer L2
 * @updated 2026-08-08
 */

package cli

import "testing"

type TestOptions struct {
	Query string `flag:"query" help:"Query de busca"`
	Max   int    `flag:"max" help:"Número máximo"`
	Lab   bool   `flag:"lab" help:"Modo lab"`
}

func TestParseStructFlags(t *testing.T) {
	opts := TestOptions{
		Query: "default",
		Max:   10,
		Lab:   false,
	}

	args := []string{"-query", "quicksort", "-max", "50", "-lab"}
	err := ParseStructFlags(&opts, args)
	if err != nil {
		t.Fatalf("Erro inesperado ao parsear struct flags: %v", err)
	}

	if opts.Query != "quicksort" {
		t.Errorf("Esperava Query 'quicksort', obteve %s", opts.Query)
	}

	if opts.Max != 50 {
		t.Errorf("Esperava Max 50, obteve %d", opts.Max)
	}

	if !opts.Lab {
		t.Errorf("Esperava Lab true, obteve false")
	}
}
