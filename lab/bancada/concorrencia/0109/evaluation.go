/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/evaluation.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"fmt"
	"sort"
)

type Expectation struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Count    int    `json:"count"`
}
type Case struct {
	ID           string        `json:"id"`
	Family       string        `json:"family"`
	Request      Request       `json:"request"`
	Expected     []Expectation `json:"expected"`
	RequiredGaps []string      `json:"required_gaps"`
}
type Corpus struct {
	SchemaVersion   string `json:"schema_version"`
	ReferenceOrigin string `json:"reference_origin"`
	Cases           []Case `json:"cases"`
}
type CaseResult struct {
	ID         string   `json:"id"`
	Family     string   `json:"family"`
	Passed     bool     `json:"passed"`
	Mismatches []string `json:"mismatches"`
}
type Evaluation struct {
	SchemaVersion                       string       `json:"schema_version"`
	Analyzer                            string       `json:"analyzer"`
	ReferenceOrigin                     string       `json:"reference_origin"`
	Cases                               int          `json:"cases"`
	Families                            int          `json:"families"`
	Passed                              int          `json:"passed"`
	SyntaxProfileAgreement              float64      `json:"syntax_profile_agreement"`
	IndependentSemanticAccuracyMeasured bool         `json:"independent_semantic_accuracy_measured"`
	Results                             []CaseResult `json:"results"`
}

func Evaluate(corpus Corpus) (Evaluation, error) {
	if corpus.SchemaVersion != "1" || corpus.ReferenceOrigin != "AUTHOR_REVIEWED_EXPLORATORY_NOT_INDEPENDENT" || len(corpus.Cases) == 0 || len(corpus.Cases) > 100 {
		return Evaluation{}, fmt.Errorf("invalid exploratory corpus")
	}
	r := Evaluation{SchemaVersion: "concurrency-evaluation-v1", Analyzer: Version, ReferenceOrigin: corpus.ReferenceOrigin, Results: []CaseResult{}}
	ids, families := map[string]bool{}, map[string]bool{}
	for _, c := range corpus.Cases {
		if !safeText(c.ID, 100) || !safeText(c.Family, 100) || ids[c.ID] {
			return Evaluation{}, fmt.Errorf("invalid corpus identity")
		}
		ids[c.ID] = true
		families[c.Family] = true
		expected := map[string]int{}
		for _, e := range c.Expected {
			if !knownProperty(e.Property) || !safeText(e.Value, 100) || e.Count < 1 || expected[e.Property+"="+e.Value] != 0 {
				return Evaluation{}, fmt.Errorf("invalid expectation in %s", c.ID)
			}
			expected[e.Property+"="+e.Value] = e.Count
		}
		for _, gap := range c.RequiredGaps {
			if !knownGap(gap) {
				return Evaluation{}, fmt.Errorf("unknown expected gap")
			}
		}
		observed, err := Observe(c.Request)
		if err != nil {
			return Evaluation{}, fmt.Errorf("case %s: %w", c.ID, err)
		}
		actual := map[string]int{}
		for _, o := range observed.Observations {
			actual[o.Property+"="+o.Value]++
		}
		result := CaseResult{ID: c.ID, Family: c.Family, Passed: true, Mismatches: []string{}}
		for key, want := range expected {
			if actual[key] != want {
				result.Mismatches = append(result.Mismatches, fmt.Sprintf("%s: expected %d, observed %d", key, want, actual[key]))
			}
		}
		for key, count := range actual {
			if expected[key] == 0 {
				result.Mismatches = append(result.Mismatches, fmt.Sprintf("unexpected %s: %d", key, count))
			}
		}
		gaps := map[string]bool{}
		for _, gap := range observed.Gaps {
			gaps[gap.Code] = true
		}
		for _, gap := range c.RequiredGaps {
			if !gaps[gap] {
				result.Mismatches = append(result.Mismatches, "missing gap: "+gap)
			}
		}
		sort.Strings(result.Mismatches)
		result.Passed = len(result.Mismatches) == 0
		if result.Passed {
			r.Passed++
		}
		r.Results = append(r.Results, result)
	}
	r.Cases = len(corpus.Cases)
	r.Families = len(families)
	r.SyntaxProfileAgreement = float64(r.Passed) / float64(r.Cases)
	return r, nil
}

func knownProperty(p string) bool {
	switch p {
	case "channel_type_identity", "channel_buffer_syntax", "channel_send", "channel_receive", "select_shape", "goroutine_launch", "cancellation_candidate":
		return true
	}
	return false
}
func knownGap(p string) bool {
	switch p {
	case "CALLBACK_VALUE_FLOW_UNRESOLVED", "FLOW_PATH_FEASIBILITY_UNRESOLVED", "CALLBACK_EXECUTION_UNRESOLVED", "CHANNEL_INSTANCE_IDENTITY_UNRESOLVED", "SYNTAX_PROPERTIES_NOT_TYPECHECKED", "SUPPLIED_PACKAGE_TYPECHECK_FAILED", "NAMED_CHANNEL_TYPE_UNRESOLVED", "BUILD_SELECTION_NOT_VERIFIED", "TYPECHECK_NOT_PERFORMED", "GO_PARSE_FAILED", "WORKER_BOUND_UNRESOLVED", "ADMISSION_BOUND_UNRESOLVED", "CANCELLATION_COMPLETION_UNRESOLVED", "DONE_RECEIVER_TYPE_UNRESOLVED", "MAKE_NAME_SHADOWED", "INVALID_CHANNEL_MAKE_ARITY", "CAPACITY_VALIDITY_NOT_TYPECHECKED":
		return true
	}
	return false
}
