/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/discovery_query_checkpoint_test.md
 * @layer L2
 * @updated 2026-09-08
 */
package cli

import (
	"bateia/01_core/domain"
	"path/filepath"
	"testing"
)

func TestDiscovery0056CheckpointVersionControlsResume(t *testing.T) {
	for _, version := range []string{"", "0054", "0055-v1", "0055-v2", "live0056-v1"} {
		t.Run("version="+version, func(t *testing.T) {
			f := newCheckpointFixture(t)
			path := filepath.Join(f.outputPath, "run.json")
			var manifest domain.IdeaRunManifest
			if err := readJSONFile(path, &manifest); err != nil {
				t.Fatal(err)
			}
			if manifest.RuleVersions["discovery_query"] != "live0056-v1" {
				t.Fatalf("live identity=%+v", manifest.RuleVersions)
			}
			if version == "" {
				delete(manifest.RuleVersions, "discovery_query")
			} else {
				manifest.RuleVersions["discovery_query"] = version
			}
			if err := writeJSONAtomic(path, manifest); err != nil {
				t.Fatal(err)
			}
			before := f.provider.calls
			f.run(t, "go", true)
			if version == "live0056-v1" {
				if f.provider.calls != before {
					t.Fatal("unchanged version did not reuse checkpoint")
				}
			} else if f.provider.calls <= before {
				t.Fatal("old discovery strategy reused checkpoint")
			}
		})
	}
}
