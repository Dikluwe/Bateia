/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/lexical/evaluation_test.md
 * @layer L3
 * @updated 2026-08-24
 */

package lexical

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"bateia/01_core/domain"
	"bateia/01_core/usecase"
)

func TestDefaultVocabularyBeatsLiteralBaselineOnVersionedCorpus(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "02_shell", "cli", "testdata", "evaluation_corpus.json"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus domain.EvaluationCorpus
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	expander, err := NewDefault(2)
	if err != nil {
		t.Fatal(err)
	}
	literal, err := usecase.EvaluateExpansionStrategy(corpus, "literal", nil)
	if err != nil {
		t.Fatal(err)
	}
	lexical, err := usecase.EvaluateExpansionStrategy(corpus, "lexical", expander)
	if err != nil {
		t.Fatal(err)
	}
	repeated, err := usecase.EvaluateExpansionStrategy(corpus, "lexical", expander)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(lexical)
	repeatedJSON, _ := json.Marshal(repeated)
	if string(firstJSON) != string(repeatedJSON) {
		t.Fatal("avaliação lexical não é determinística")
	}
	if lexical.Metrics.Coverage <= literal.Metrics.Coverage {
		t.Fatalf("lexical não superou literal: literal=%+v lexical=%+v", literal.Metrics, lexical.Metrics)
	}
	if lexical.Metrics.Contamination != 0 {
		t.Fatalf("baseline contaminado: %+v", lexical.Metrics)
	}
	withNovelty := 0
	for _, result := range lexical.Cases {
		if result.NewExpansionCount > 0 {
			withNovelty++
		}
	}
	if withNovelty*2 < len(lexical.Cases) {
		t.Fatalf("menos da metade dos casos recebeu expansão: %d/%d", withNovelty, len(lexical.Cases))
	}
}
