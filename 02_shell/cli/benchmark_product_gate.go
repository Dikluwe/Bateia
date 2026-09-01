/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/benchmark_product_gate.md
 * @layer L2
 * @updated 2026-08-31
 */
package cli

import (
	"bateia/01_core/domain"
	"bateia/01_core/usecase"
	"flag"
	"fmt"
	"io"
)

func runBenchmarkProductGate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("benchmark product-gate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	inputPath := fs.String("input", "", "entrada congelada do benchmark de produto")
	outputPath := fs.String("output", "", "relatório de produto")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parâmetros de product-gate inválidos: %w", err)
	}
	if *inputPath == "" {
		return fmt.Errorf("product-gate exige --input")
	}
	var input domain.ProductValidationInput
	if err := readDecisionJSON(*inputPath, &input); err != nil {
		return err
	}
	return writeIdeaJSON(usecase.ValidateSoftwareSearchProduct(input), *outputPath, stdout)
}
