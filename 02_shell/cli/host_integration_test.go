/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/host_integration_test.md
 * @layer L2
 * @updated 2026-10-28
 */
package cli_test

import (
	"bateia/02_shell/cli"
	"bateia/03_infra/host"
	"os"
	"testing"
)

func TestMain(m *testing.M) { cli.ConfigureCommandHost(host.New()); os.Exit(m.Run()) }
