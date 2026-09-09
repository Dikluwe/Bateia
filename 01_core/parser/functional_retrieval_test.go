/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/functional_retrieval_test.md
 * @layer L1
 * @updated 2026-09-07
 */
package parser_test

import (
	"bateia/01_core/domain"
	"bateia/01_core/parser"
	"encoding/json"
	"reflect"
	"testing"
)

type retrieval0045Case struct {
	ID       string                 `json:"id"`
	Query    string                 `json:"query"`
	Files    []domain.ExtractedFile `json:"files"`
	Expected []string               `json:"expected"`
}

const retrieval0045Corpus = `[
  {
    "id": "queue",
    "query": "enfileirar tarefas",
    "files": [
      {
        "Path": "queue/dispatch.go",
        "Content": "package queue\nfunc EnqueueTask(q chan string, task string) { q <- task }\n"
      },
      {
        "Path": "queue/doc0.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description0 struct {}\n"
      },
      {
        "Path": "queue/doc1.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description1 struct {}\n"
      },
      {
        "Path": "queue/doc2.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description2 struct {}\n"
      },
      {
        "Path": "queue/doc3.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description3 struct {}\n"
      },
      {
        "Path": "queue/literal.go",
        "Content": "package queue\nfunc Message() string { return \"enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \" }\n"
      },
      {
        "Path": "queue/irrelevant.go",
        "Content": "package queue\nfunc Add(a,b int) int { return a+b }\n"
      }
    ],
    "expected": [
      "queue/dispatch.go"
    ]
  },
  {
    "id": "storage",
    "query": "gravar arquivo",
    "files": [
      {
        "Path": "storage/persist.go",
        "Content": "package storage\nfunc WriteFile(dst interface{ Write([]byte)(int,error) }, data []byte) { dst.Write(data) }\n"
      },
      {
        "Path": "storage/doc0.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description0 struct {}\n"
      },
      {
        "Path": "storage/doc1.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description1 struct {}\n"
      },
      {
        "Path": "storage/doc2.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description2 struct {}\n"
      },
      {
        "Path": "storage/doc3.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description3 struct {}\n"
      },
      {
        "Path": "storage/literal.go",
        "Content": "package storage\nfunc Message() string { return \"gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \" }\n"
      },
      {
        "Path": "storage/irrelevant.go",
        "Content": "package storage\nfunc Add(a,b int) int { return a+b }\n"
      }
    ],
    "expected": [
      "storage/persist.go"
    ]
  },
  {
    "id": "auth",
    "query": "validar senha",
    "files": [
      {
        "Path": "auth/check.go",
        "Content": "package auth\nfunc ValidatePassword(password, expected string) bool { return password == expected }\n"
      },
      {
        "Path": "auth/doc0.go",
        "Content": "package auth\n// validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha \ntype Description0 struct {}\n"
      },
      {
        "Path": "auth/doc1.go",
        "Content": "package auth\n// validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha \ntype Description1 struct {}\n"
      },
      {
        "Path": "auth/doc2.go",
        "Content": "package auth\n// validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha \ntype Description2 struct {}\n"
      },
      {
        "Path": "auth/doc3.go",
        "Content": "package auth\n// validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha \ntype Description3 struct {}\n"
      },
      {
        "Path": "auth/literal.go",
        "Content": "package auth\nfunc Message() string { return \"validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha validar senha \" }\n"
      },
      {
        "Path": "auth/irrelevant.go",
        "Content": "package auth\nfunc Add(a,b int) int { return a+b }\n"
      }
    ],
    "expected": [
      "auth/check.go"
    ]
  },
  {
    "id": "exact",
    "query": "EnqueueTask",
    "files": [
      {
        "Path": "queue/dispatch.go",
        "Content": "package queue\nfunc EnqueueTask(q chan string, task string) { q <- task }\n"
      },
      {
        "Path": "queue/doc0.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description0 struct {}\n"
      },
      {
        "Path": "queue/doc1.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description1 struct {}\n"
      },
      {
        "Path": "queue/doc2.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description2 struct {}\n"
      },
      {
        "Path": "queue/doc3.go",
        "Content": "package queue\n// enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \ntype Description3 struct {}\n"
      },
      {
        "Path": "queue/literal.go",
        "Content": "package queue\nfunc Message() string { return \"enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas enfileirar tarefas \" }\n"
      },
      {
        "Path": "queue/irrelevant.go",
        "Content": "package queue\nfunc Add(a,b int) int { return a+b }\n"
      }
    ],
    "expected": [
      "queue/dispatch.go"
    ]
  },
  {
    "id": "absent",
    "query": "zygomorphic",
    "files": [
      {
        "Path": "storage/persist.go",
        "Content": "package storage\nfunc WriteFile(dst interface{ Write([]byte)(int,error) }, data []byte) { dst.Write(data) }\n"
      },
      {
        "Path": "storage/doc0.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description0 struct {}\n"
      },
      {
        "Path": "storage/doc1.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description1 struct {}\n"
      },
      {
        "Path": "storage/doc2.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description2 struct {}\n"
      },
      {
        "Path": "storage/doc3.go",
        "Content": "package storage\n// gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \ntype Description3 struct {}\n"
      },
      {
        "Path": "storage/literal.go",
        "Content": "package storage\nfunc Message() string { return \"gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo gravar arquivo \" }\n"
      },
      {
        "Path": "storage/irrelevant.go",
        "Content": "package storage\nfunc Add(a,b int) int { return a+b }\n"
      }
    ],
    "expected": []
  },
  {
    "id": "invalid",
    "query": "ValidatePassword",
    "files": [
      {
        "Path": "invalid.go",
        "Content": "package auth\nfunc ValidatePassword( invalid !!!"
      }
    ],
    "expected": []
  }
]
`

func TestFunctional0045Evaluation(t *testing.T) {
	var cases []retrieval0045Case
	if err := json.Unmarshal([]byte(retrieval0045Corpus), &cases); err != nil {
		t.Fatal(err)
	}
	totalTP, totalExpected, totalReturned := 0, 0, 0
	for _, c := range cases {
		report := parser.RankGoImplementations(c.Query, c.Files, 9)
		paths := []string{}
		seen := map[string]bool{}
		tp := 0
		for _, h := range report.Hits {
			if !seen[h.Path] && len(paths) < 3 {
				seen[h.Path] = true
				paths = append(paths, h.Path)
			}
		}
		for _, p := range c.Expected {
			if seen[p] {
				tp++
			}
		}
		totalTP += tp
		totalExpected += len(c.Expected)
		totalReturned += len(paths)
		out, _ := json.Marshal(map[string]any{"id": c.ID, "paths": paths, "tp": tp, "expected": len(c.Expected), "returned": len(paths)})
		t.Log("METRIC0045 " + string(out))
		if (c.ID == "absent" || c.ID == "invalid") && len(paths) != 0 {
			t.Errorf("negative %s returned %v", c.ID, paths)
		}
		if c.ID == "exact" && (len(report.Hits) == 0 || report.Hits[0].Symbol != "EnqueueTask") {
			t.Error("exact identifier lost priority")
		}
	}
	out, _ := json.Marshal(map[string]any{"id": "aggregate", "tp": totalTP, "expected": totalExpected, "returned": totalReturned})
	t.Log("METRIC0045 " + string(out))
}
func TestFunctional0045Determinism(t *testing.T) {
	var cases []retrieval0045Case
	json.Unmarshal([]byte(retrieval0045Corpus), &cases)
	for _, c := range cases {
		a := parser.RankGoImplementations(c.Query, c.Files, 9)
		reversed := append([]domain.ExtractedFile{}, c.Files...)
		for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
			reversed[i], reversed[j] = reversed[j], reversed[i]
		}
		reversed = append(reversed, c.Files...)
		b := parser.RankGoImplementations(c.Query, reversed, 9)
		if !reflect.DeepEqual(a.Hits, b.Hits) {
			t.Errorf("%s order/duplicate instability", c.ID)
		}
	}
}

func retrieval0045Files(sources map[string]string) []domain.ExtractedFile {
	files := []domain.ExtractedFile{}
	for path, content := range sources {
		files = append(files, domain.ExtractedFile{Path: path, Content: content})
	}
	return files
}
func TestFunctional0045OperationalPriority(t *testing.T) {
	var cases []retrieval0045Case
	json.Unmarshal([]byte(retrieval0045Corpus), &cases)
	for _, c := range cases[:3] {
		r := parser.RankGoImplementations(c.Query, c.Files, 9)
		if len(r.Hits) == 0 || r.Hits[0].Path != c.Expected[0] {
			t.Errorf("%s functional definition must outrank repeated documentation and strings", c.ID)
		}
	}
}
func TestFunctional0045DirectCalls(t *testing.T) {
	base := map[string]string{"q/entry.go": "package q\nfunc QueueTask() {\n alpha()\n}\n", "q/alpha.go": "package q\nfunc alpha() { beta() }\n", "q/beta.go": "package q\nfunc beta() {}\n"}
	r := parser.RankGoImplementations("QueueTask", retrieval0045Files(base), 9)
	encoded, _ := json.Marshal(r)
	var wire struct {
		Hits []struct {
			Symbol    string `json:"symbol"`
			Retrieval string `json:"retrieval"`
			Related   []struct {
				Path   string `json:"path"`
				Symbol string `json:"symbol"`
				Line   int    `json:"line"`
				SHA256 string `json:"sha256"`
			} `json:"related_to"`
		}
	}
	json.Unmarshal(encoded, &wire)
	found := false
	for _, h := range wire.Hits {
		if h.Symbol == "beta" {
			t.Error("second-hop helper leaked")
		}
		if h.Symbol == "alpha" {
			found = true
			if h.Retrieval != "direct_call" || len(h.Related) != 1 {
				t.Errorf("helper needs direct_call and one reference: %+v", h)
			} else if h.Related[0].Path != "q/entry.go" || h.Related[0].Symbol != "QueueTask" || h.Related[0].Line != 3 || len(h.Related[0].SHA256) != 64 {
				t.Errorf("invalid caller identity %+v", h.Related)
			}
		}
	}
	if !found {
		t.Error("unique one-hop helper without query terms not retrieved")
	}
	for _, limit := range []int{1, 2, 3, 6, 9} {
		r := parser.RankGoImplementations("QueueTask", retrieval0045Files(base), limit)
		if len(r.Hits) > limit {
			t.Errorf("hit limit %d exceeded", limit)
		}
		b, _ := json.Marshal(r)
		json.Unmarshal(b, &wire)
		n := 0
		for _, h := range wire.Hits {
			if h.Retrieval == "direct_call" {
				n++
			}
		}
		if n > limit/3 {
			t.Errorf("expansion consumed %d/%d", n, limit)
		}
	}
}
func TestFunctional0045UnsafeCallResolution(t *testing.T) {
	cases := map[string]map[string]string{
		"shadow-local":     {"q/entry.go": "package q\nfunc QueueTask() { alpha := func() {}; alpha() }", "q/helper.go": "package q\nfunc alpha() {}"},
		"shadow-parameter": {"q/entry.go": "package q\nfunc QueueTask(alpha func()) { alpha() }", "q/helper.go": "package q\nfunc alpha() {}"},
		"ambiguous":        {"q/entry.go": "package q\nfunc QueueTask() { alpha() }", "q/helper.go": "package q\nfunc alpha() {}", "q/other.go": "package q\nfunc alpha() {}"},
		"selector":         {"q/entry.go": "package q\nfunc QueueTask() { x.alpha() }", "q/helper.go": "package q\nfunc alpha() {}"},
		"function-value":   {"q/entry.go": "package q\nfunc QueueTask() { f := alpha; f() }", "q/helper.go": "package q\nfunc alpha() {}"},
		"other-dir":        {"q/entry.go": "package q\nfunc QueueTask() { alpha() }", "other/helper.go": "package q\nfunc alpha() {}"},
		"other-package":    {"q/entry.go": "package q\nfunc QueueTask() { alpha() }", "q/helper.go": "package other\nfunc alpha() {}"},
	}
	for name, sources := range cases {
		r := parser.RankGoImplementations("QueueTask", retrieval0045Files(sources), 9)
		for _, h := range r.Hits {
			if h.Symbol == "alpha" {
				t.Errorf("%s resolved unsupported helper", name)
			}
		}
	}
}
func TestFunctional0045MultiRequirement(t *testing.T) {
	files := retrieval0045Files(map[string]string{"q/a.go": "package q\nfunc QueueTaskOne() {}\nfunc QueueTaskTwo() {}\nfunc QueueTaskThree() {}", "q/b.go": "package q\nfunc ValidatePassword() {}"})
	r := parser.RankGoImplementations("QueueTask\nValidatePassword", files, 3)
	seen := map[string]bool{}
	for _, h := range r.Hits {
		seen[h.Path] = true
	}
	if !seen["q/a.go"] || !seen["q/b.go"] {
		t.Errorf("requirements not distributed: %+v", r.Hits)
	}
}
