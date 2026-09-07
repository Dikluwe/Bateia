/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/idea_step0040_test.md
 * @layer L2
 * @updated 2026-09-07
 */
package cli

import (
	"bateia/01_core/domain"
	"reflect"
	"testing"
)

func TestStep0040LocalMetadataDoesNotDependOnIdea(t *testing.T) {
	paths := []string{t.TempDir(), t.TempDir()}
	a, err := localCandidates(domain.IdeaSpec{Problem: "asynchronous event processing", Capabilities: []string{"event queue"}}, paths)
	if err != nil {
		t.Fatal(err)
	}
	b, err := localCandidates(domain.IdeaSpec{Problem: "export notes to Markdown", Capabilities: []string{"markdown"}}, paths)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("query contaminated candidate metadata")
	}
	for _, c := range a {
		if c.Description != "" {
			t.Fatalf("unobserved description: %q", c.Description)
		}
	}
}
