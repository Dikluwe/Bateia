/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/algorithm_input_bytes_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFormulationRejectsOriginalInvalidBytes(t *testing.T) {
	for name, mutate := range map[string]func(*domain.AlgorithmInsightRequest){
		"question":   func(r *domain.AlgorithmInsightRequest) { r.Question += "\xff" },
		"title":      func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Title += "\xff" },
		"principle":  func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Principle += "\xff" },
		"step":       func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Steps[0] += "\xff" },
		"condition":  func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].ApplicationConditions[0] += "\xff" },
		"failure":    func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].FailureConditions[0] += "\xff" },
		"decision":   func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].Decision.Current += "\xff" },
		"observable": func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].ValidationObservable += "\xff" },
		"reject":     func(r *domain.AlgorithmInsightRequest) { r.Ideas[0].RejectIf += "\xff" },
		"source repaired hash": func(r *domain.AlgorithmInsightRequest) {
			r.Sources[0].Items[0].Content += "\xff"
			h := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ReplaceAll(r.Sources[0].Items[0].Content, "\xff", "�"))))
			r.Sources[0].Items[0].ContentHash, r.Ideas[0].Evidence[0].ContentHash = h, h
		},
		"target repaired hash": func(r *domain.AlgorithmInsightRequest) {
			r.Target.Items[0].Content += "\xff"
			h := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ReplaceAll(r.Target.Items[0].Content, "\xff", "�"))))
			r.Target.Items[0].ContentHash, r.Ideas[0].TargetEvidence[0].ContentHash = h, h
		},
		"uncited source": func(r *domain.AlgorithmInsightRequest) {
			r.Sources[0].Items = append(r.Sources[0].Items, domain.CollectedEvidence{Kind: domain.EvidenceFile, Ref: r.Sources[0].Repository, Path: "uncited.go", Content: "\xff", ContentHash: fmt.Sprintf("%x", sha256.Sum256([]byte("�")))})
		},
	} {
		t.Run(name, func(t *testing.T) {
			r, before := algorithmContractFixture(), algorithmContractFixture()
			mutate(&r)
			mutate(&before)
			got, err := FormulateAlgorithmInsights(context.Background(), r)
			if err == nil || got.ID != "" {
				t.Errorf("invalid original bytes accepted: %v %s", err, got.ID)
			}
			if !reflect.DeepEqual(r, before) {
				t.Fatal("caller mutated")
			}
		})
	}
}

func TestFormulationValidBytesKeepIdentity(t *testing.T) {
	r := algorithmContractFixture()
	got, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "algorithms_16a7f584dcdd7468cb65d8502de60f2fad2b35e0a8b62d9ab286fa7aa7506e77" {
		t.Fatal("valid identity changed", got.ID)
	}
	// A replacement character supplied as valid UTF-8 is ordinary literal evidence.
	r.Sources[0].Items[0].Content = strings.ReplaceAll(r.Sources[0].Items[0].Content, "reference", "refer�ence")
	r.Sources[0].Items[0].ContentHash = fmt.Sprintf("%x", sha256.Sum256([]byte(r.Sources[0].Items[0].Content)))
	r.Ideas[0].Evidence[0] = algorithmContractReference(r.Sources[0], 1)
	r.Ideas[0].Principle += " �"
	valid, err := FormulateAlgorithmInsights(context.Background(), r)
	if err != nil || valid.ID == got.ID || valid.Ideas[0].Principle != r.Ideas[0].Principle {
		t.Fatal("literal replacement character rejected or altered", err)
	}
}
