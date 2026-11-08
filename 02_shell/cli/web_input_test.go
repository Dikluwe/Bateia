/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/web_input_test.md
 * @layer L2
 * @updated 2026-10-29
 */
package cli

import "testing"

func TestWebTokenRequestedPreservesHistoricalEarlyProbe(t *testing.T) {
	for _, test := range []struct {
		args []string
		want bool
	}{
		{nil, true},
		{[]string{"web", "experiment"}, true},
		{[]string{"web", "experiment", "--no-token"}, false},
		{[]string{"--no-token", "web", "--no-token"}, false},
		{[]string{"--no-token=true"}, true},
		{[]string{"--no-token=false"}, true},
		{[]string{"--NO-TOKEN"}, true},
	} {
		if got := WebTokenRequested(test.args); got != test.want {
			t.Errorf("probe %q: got %v want %v", test.args, got, test.want)
		}
	}
}
