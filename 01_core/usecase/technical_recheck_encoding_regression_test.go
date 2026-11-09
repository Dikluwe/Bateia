/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/technical_recheck_encoding_regression_test.md
 * @layer L1
 * @updated 2026-11-08
 */
package usecase

import (
	"bateia/01_core/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestTechnicalRecheckEncodingAdmissionAndReplay(t *testing.T) {
	for _, description := range []string{"valid", "literal �", "invalid\xff"} {
		t.Run(fmt.Sprintf("%q", description), func(t *testing.T) {
			req, p := recheckFixture()
			content := `{"bomFormat":"CycloneDX","specVersion":"1.6","components":[{"name":"widget","version":"1.2.3","description":"` + description + `"}]}`
			digest := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
			f := &req.Batch.Findings[0]
			f.Technique, f.Predicate, f.Value, f.Claim = domain.TechniqueSBOM, domain.TechnicalSBOMComponent, "widget@1.2.3", "SBOM declares widget"
			f.Source.Ref, f.Source.Digest = "sbom.json", "sha256:"+digest
			p.tree = []domain.DirectoryNode{{Path: "sbom.json", Type: "file", Size: int64(len(content))}}
			p.extracted.Files = []domain.ExtractedFile{{Path: "sbom.json", Content: content, Size: int64(len(content))}}
			before, _ := json.Marshal(req)
			_, _, _, asOf := verificationFixture()
			got, err := RecheckTechnicalSources(context.Background(), req, p, asOf)
			if err != nil || len(got.Bundles) != 1 {
				t.Fatal(err, got)
			}
			valid := description != "invalid\xff"
			b, v := got.Bundles[0], got.Review.Verifications[0]
			if !valid {
				if b.Complete || len(b.Items) != 0 || len(b.Failures) == 0 || len(got.Failures) == 0 || v.IntegrityVerified || v.Status == domain.TechnicalVerified {
					t.Errorf("invalid bytes admitted: bundle=%+v verification=%+v", b, v)
				}
			} else if !b.Complete || len(b.Items) != 1 || b.Items[0].Content != content || b.Items[0].ContentHash != digest || !v.IntegrityVerified || v.Status != domain.TechnicalVerified {
				t.Fatal("valid control changed", got)
			}
			for round := 0; round < 3; round++ {
				raw, e := json.Marshal(got)
				if e != nil {
					t.Fatal(e)
				}
				var reopened domain.TechnicalSourceRecheckReport
				if e = json.Unmarshal(raw, &reopened); e != nil {
					t.Fatal(e)
				}
				checked, e := VerifyTechnicalEvidence(req.Batch, reopened.Bundles, asOf)
				if e != nil || checked.Verifications[0].IntegrityVerified != valid {
					t.Errorf("replay integrity differs: %v %+v", e, checked)
				}
				if valid && (reopened.Bundles[0].Items[0].Content != content || reopened.Bundles[0].Items[0].ContentHash != digest) {
					t.Fatal("replay changed bytes/hash")
				}
				got = reopened
			}
			after, _ := json.Marshal(req)
			if string(before) != string(after) || p.extracted.Files[0].Content != content || p.calls != 2 || !reflect.DeepEqual(p.paths, []string{"sbom.json"}) {
				t.Fatal("inputs or reads changed")
			}
		})
	}
}
