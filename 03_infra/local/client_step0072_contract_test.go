/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/local/client_step0072_contract_test.md
 * @layer L3
 * @updated 2026-09-16
 */

package local_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
	localinfra "bateia/03_infra/local"
)

func TestStep0072LocalTreeDoesNotSilentlyStopAtTwentyThousandEntries(t *testing.T) {
	root := t.TempDir()
	const distractors = 20_005
	for i := 0; i < distractors; i++ {
		name := filepath.Join(root, fmt.Sprintf("distractor_%05d.hpp", i))
		file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "span"), []byte("template<class T> class span {};"), 0o600); err != nil {
		t.Fatal(err)
	}

	tree, err := localinfra.NewClient().FetchTree(context.Background(), domain.RepositoryCandidate{
		Owner: "local", Repo: "libstdcxx", CommitID: "immutable-13", URL: "file://" + root,
	})
	if err == nil {
		t.Fatalf("árvore acima do teto devolvida como sucesso: entries=%d", len(tree))
	}
	if tree != nil {
		t.Fatalf("árvore incompleta não deve parecer resultado íntegro: entries=%d err=%v", len(tree), err)
	}
}
