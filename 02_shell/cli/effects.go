/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/effects.md
 * @layer L2
 * @updated 2026-09-08
 */

package cli

import (
	"bateia/01_core/domain"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"
	"time"
)

const effectHelp = `Usage: bateia effects verify --source FILE --symbol NAME [--output FILE] [--timeout 5s] [--seed 1] [--expect-sha256 HEX]
Offline sort-preserve-v1 pilot: package candidate, exported func([]int64) []int64;
only functions and imports sort/slices without aliases. Fixed cases plus 32 seeded vectors.
Checks returned ordering, multiplicities and input values after return, only on exercised inputs.
No general termination, transient-write, memory-sharing or external-effect proof.
Separate process is NOT a security sandbox or rigid memory limit; do not use hostile code.
Execution timeout >0 and <=30s; build budget 60s; stdout/stderr each <=1MiB.
Durations are tool telemetry, not candidate benchmarks. Unknown never means success.
--expect-sha256 binds the complete source file to 64 ASCII hexadecimal digits (case-insensitive).
A mismatch returns unknown without execution. Omission does not verify prior source identity.
No automatic link to search evidence, repository/commit identity or external dependencies.
`

func RunEffectsCommand(args []string, out io.Writer, verifier domain.EffectVerifier) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		_, err := io.WriteString(out, effectHelp)
		return err
	}
	if len(args) == 0 || args[0] != "verify" {
		return fmt.Errorf("expected effects verify")
	}
	fs := flag.NewFlagSet("effects verify", flag.ContinueOnError)
	fs.SetOutput(out)
	fs.Usage = func() { fmt.Fprint(out, effectHelp) }
	source := fs.String("source", "", "source Go file")
	symbol := fs.String("symbol", "", "exported function")
	output := fs.String("output", "", "JSON output (default stdout)")
	timeout := fs.Duration("timeout", 5*time.Second, "execution timeout")
	expected := fs.String("expect-sha256", "", "expected SHA-256 of the complete source file")
	seed := fs.Uint64("seed", 1, "deterministic seed")
	if err := fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() != 0 || strings.TrimSpace(*source) == "" || !token.IsIdentifier(*symbol) || !ast.IsExported(*symbol) || *timeout <= 0 || *timeout > 30*time.Second {
		return fmt.Errorf("invalid source, exported symbol, timeout or trailing arguments")
	}
	expectedProvided := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "expect-sha256" {
			expectedProvided = true
		}
	})
	if expectedProvided {
		if _, err := hex.DecodeString(*expected); len(*expected) != 64 || err != nil {
			return fmt.Errorf("expect-sha256 must contain exactly 64 ASCII hexadecimal digits")
		}
		*expected = strings.ToLower(*expected)
	}
	if verifier == nil {
		return fmt.Errorf("effect verifier unavailable")
	}
	report, err := verifier.Verify(context.Background(), domain.EffectRequest{ExpectedSourceSHA256: *expected, SourcePath: *source, Symbol: *symbol, Seed: *seed, Timeout: *timeout})
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if *output != "" {
		return os.WriteFile(*output, data, 0600)
	}
	_, err = out.Write(data)
	return err
}
