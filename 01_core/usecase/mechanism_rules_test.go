/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/mechanism_rules_test.md
 * @layer L1
 * @updated 2026-09-18
 */

package usecase_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"bateia/01_core/usecase"
)

const mechanismRuleCatalogSHA256 = "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"

type mechanismRuleCatalogEnvelope struct {
	SchemaVersion         string `json:"schema_version"`
	CatalogID             string `json:"catalog_id"`
	CatalogVersion        int    `json:"catalog_version"`
	ClosedWorld           bool   `json:"closed_world"`
	UnknownClassification string `json:"unknown_classification"`
	Rules                 []struct {
		ID                    string `json:"id"`
		AcceptedEvidenceRules []struct {
			RuleID string `json:"rule_id"`
		} `json:"accepted_evidence_rules"`
	} `json:"rules"`
}

func TestActiveMechanismRuleCatalogMatchesFrozenV3Catalog(t *testing.T) {
	catalog, claimedHash := usecase.ActiveMechanismRuleCatalog()
	if len(catalog) != 8587 {
		t.Fatalf("catalog length = %d, want 8587", len(catalog))
	}
	if claimedHash != mechanismRuleCatalogSHA256 {
		t.Fatalf("claimed hash = %q, want %q", claimedHash, mechanismRuleCatalogSHA256)
	}
	if len(claimedHash) != sha256.Size*2 {
		t.Fatalf("claimed hash length = %d, want %d", len(claimedHash), sha256.Size*2)
	}
	if decoded, err := hex.DecodeString(claimedHash); err != nil || len(decoded) != sha256.Size {
		t.Fatalf("claimed hash is not canonical lowercase hexadecimal: bytes=%d err=%v", len(decoded), err)
	}
	digest := sha256.Sum256(catalog)
	if computed := hex.EncodeToString(digest[:]); computed != mechanismRuleCatalogSHA256 {
		t.Fatalf("computed hash = %q, want %q", computed, mechanismRuleCatalogSHA256)
	}

	var envelope mechanismRuleCatalogEnvelope
	if err := json.Unmarshal(catalog, &envelope); err != nil {
		t.Fatalf("catalog is not valid JSON: %v", err)
	}
	if envelope.SchemaVersion != "tekt.mechanism-rules.v3" {
		t.Fatalf("schema version = %q, want tekt.mechanism-rules.v3", envelope.SchemaVersion)
	}
	if envelope.CatalogID != "bateia-step0028-initial" || envelope.CatalogVersion != 3 {
		t.Fatalf("catalog identity = %q/v%d, want bateia-step0028-initial/v3", envelope.CatalogID, envelope.CatalogVersion)
	}
	if !envelope.ClosedWorld || envelope.UnknownClassification != "MECHANISM_CLASS_UNKNOWN" {
		t.Fatalf("closed-world policy changed: closed=%v unknown=%q", envelope.ClosedWorld, envelope.UnknownClassification)
	}
	if len(envelope.Rules) != 4 {
		t.Fatalf("rule count = %d, want 4", len(envelope.Rules))
	}
	wantRuleIDs := []string{
		"MR-APPEND-ONLY-LOG-V1",
		"MR-BOUNDED-QUEUE-V1",
		"MR-CONTENT-ADDRESS-V1",
		"MR-ADAPTER-BOUNDARY-V1",
	}
	bindings := 0
	for index, rule := range envelope.Rules {
		if rule.ID != wantRuleIDs[index] {
			t.Fatalf("rule[%d] id = %q, want %q", index, rule.ID, wantRuleIDs[index])
		}
		bindings += len(rule.AcceptedEvidenceRules)
	}
	if bindings != 32 {
		t.Fatalf("accepted evidence binding count = %d, want 32", bindings)
	}
}

func TestActiveMechanismRuleCatalogReturnsFreshDefensiveClone(t *testing.T) {
	first, firstHash := usecase.ActiveMechanismRuleCatalog()
	second, secondHash := usecase.ActiveMechanismRuleCatalog()
	if firstHash != mechanismRuleCatalogSHA256 || secondHash != mechanismRuleCatalogSHA256 {
		t.Fatalf("repeated hashes = %q and %q, want frozen hash", firstHash, secondHash)
	}
	if len(first) == 0 || len(second) == 0 {
		t.Fatal("catalog snapshots must not be empty")
	}
	if &first[0] == &second[0] {
		t.Fatal("two calls returned the same backing array")
	}
	if !bytes.Equal(first, second) {
		t.Fatal("two fresh snapshots differ before mutation")
	}
	pristine := append([]byte(nil), second...)
	first[0] ^= 0xff
	first[len(first)-1] ^= 0xff
	if !bytes.Equal(second, pristine) {
		t.Fatal("mutating one snapshot changed an earlier independent snapshot")
	}
	third, thirdHash := usecase.ActiveMechanismRuleCatalog()
	if thirdHash != mechanismRuleCatalogSHA256 || !bytes.Equal(third, pristine) {
		t.Fatal("mutating a returned snapshot contaminated a later call")
	}
}

func TestActiveMechanismRuleCatalogRepeatedCallsRemainIndependent(t *testing.T) {
	canonical, canonicalHash := usecase.ActiveMechanismRuleCatalog()
	if canonicalHash != mechanismRuleCatalogSHA256 {
		t.Fatalf("canonical hash = %q, want frozen hash", canonicalHash)
	}
	for iteration := 0; iteration < 100; iteration++ {
		got, gotHash := usecase.ActiveMechanismRuleCatalog()
		if gotHash != canonicalHash || !bytes.Equal(got, canonical) {
			t.Fatalf("call %d differs from canonical snapshot", iteration)
		}
		if len(got) > 0 {
			got[iteration%len(got)] ^= byte(iteration + 1)
		}
	}
	after, afterHash := usecase.ActiveMechanismRuleCatalog()
	if afterHash != canonicalHash || !bytes.Equal(after, canonical) {
		t.Fatal("mutations across repeated calls contaminated catalog state")
	}
}

func TestActiveMechanismRuleCatalogConcurrentCallsAreIndependent(t *testing.T) {
	canonical, canonicalHash := usecase.ActiveMechanismRuleCatalog()
	if canonicalHash != mechanismRuleCatalogSHA256 || len(canonical) != 8587 {
		t.Fatal("canonical catalog precondition failed")
	}

	const workers = 64
	const callsPerWorker = 32
	start := make(chan struct{})
	errors := make(chan error, workers)
	var wait sync.WaitGroup
	wait.Add(workers)
	for worker := 0; worker < workers; worker++ {
		worker := worker
		go func() {
			defer wait.Done()
			<-start
			for call := 0; call < callsPerWorker; call++ {
				catalog, claimedHash := usecase.ActiveMechanismRuleCatalog()
				if claimedHash != canonicalHash || !bytes.Equal(catalog, canonical) {
					errors <- fmt.Errorf("worker %d call %d observed a divergent snapshot", worker, call)
					return
				}
				catalog[(worker+call)%len(catalog)] ^= 0xff
			}
		}()
	}
	close(start)
	wait.Wait()
	close(errors)
	for err := range errors {
		t.Error(err)
	}
	after, afterHash := usecase.ActiveMechanismRuleCatalog()
	if afterHash != canonicalHash || !bytes.Equal(after, canonical) {
		t.Fatal("concurrent mutations contaminated the embedded catalog")
	}
}
