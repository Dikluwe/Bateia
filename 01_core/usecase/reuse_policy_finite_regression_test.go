/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/reuse_policy_finite_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func TestReusePolicyRejectsNonFiniteAndOutOfRangeThresholds(t *testing.T) {
	for _, field := range []string{"extend", "integration"} {
		for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 0, -.1, 1.001, .0001, 1} {
			t.Run(fmt.Sprintf("%s/%v", field, value), func(t *testing.T) {
				p := decisionPolicy()
				if field == "extend" {
					p.ExtendMinCoverage = value
				} else {
					p.IntegrationMinValue = value
				}
				req := domain.ReuseDecisionRequest{Idea: domain.IdeaSpec{Capabilities: []string{"queue"}}, Policy: p}
				got, err := DecideReuse(req)
				valid := !math.IsNaN(value) && !math.IsInf(value, 0) && value > 0 && value <= 1
				if !valid {
					if err == nil || got.PolicyHash != "" || got.Decision != "" {
						t.Errorf("invalid policy admitted: %v %+v", err, got)
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}
					raw, e := json.Marshal(p)
					if e != nil || got.PolicyHash != fmt.Sprintf("%x", sha256.Sum256(raw)) {
						t.Fatal("policy hash changed", e)
					}
					if _, e = json.Marshal(got); e != nil {
						t.Fatal("successful result is not serializable", e)
					}
				}
				if math.Float64bits(req.Policy.ExtendMinCoverage) != math.Float64bits(p.ExtendMinCoverage) || math.Float64bits(req.Policy.IntegrationMinValue) != math.Float64bits(p.IntegrationMinValue) {
					t.Fatal("policy mutated")
				}
			})
		}
	}
	// Literal hash captured from the valid public-API control before this patch.
	p := domain.ReuseDecisionPolicy{SchemaVersion: "1", AllowedLicenses: []string{"MIT"}, MaxComponents: 2, ExtendMinCoverage: .6, IntegrationMinValue: .3}
	got, err := DecideReuse(domain.ReuseDecisionRequest{Idea: domain.IdeaSpec{Capabilities: []string{"queue"}}, Policy: p})
	if err != nil || got.PolicyHash != "ea94e37e39ef1a77b95f16e68b4ca2e455aac56b39f9af5a8bfbff26f5e61014" {
		t.Fatal("valid hash changed", err, got.PolicyHash)
	}
}
