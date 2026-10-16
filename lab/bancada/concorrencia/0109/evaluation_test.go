/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/evaluation_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"encoding/json"
	"os"
	"testing"
)

func loadCorpus(t *testing.T) Corpus {
	t.Helper()
	data, err := os.ReadFile("../0110/testdata/baseline-v2.json")
	if err != nil {
		t.Fatal(err)
	}
	var c Corpus
	if err = json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestExploratoryCorpusAndDiscriminatingNegativeControl(t *testing.T) {
	c := loadCorpus(t)
	r, err := Evaluate(c)
	if err != nil {
		t.Fatal(err)
	}
	if r.Cases != 15 || r.Passed != r.Cases || r.Families >= r.Cases || r.IndependentSemanticAccuracyMeasured {
		t.Fatalf("evaluation: %+v", r)
	}
	c.Cases[2].Expected[0].Value = "CAPACITY_LITERAL:99"
	r, err = Evaluate(c)
	if err != nil || r.Passed != r.Cases-1 {
		t.Fatalf("negative control survived: %+v %v", r, err)
	}
	c = loadCorpus(t)
	c.Cases[2].Expected = nil
	r, err = Evaluate(c)
	if err != nil || r.Passed != r.Cases-1 {
		t.Fatal("unexpected observation was ignored")
	}
	c = loadCorpus(t)
	c.Cases[0].RequiredGaps = []string{"SOME_UNKNOWN_LABEL"}
	if _, err = Evaluate(c); err == nil {
		t.Fatal("unknown label accepted")
	}
}

func TestTypeContrastCorpus(t *testing.T) {
	data, err := os.ReadFile("../0110/testdata/types-corpus.json")
	if err != nil {
		t.Fatal(err)
	}
	var c Corpus
	if err = json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	r, err := Evaluate(c)
	if err != nil || r.Passed != r.Cases || r.IndependentSemanticAccuracyMeasured {
		t.Fatalf("types corpus: %+v %v", r, err)
	}
}
