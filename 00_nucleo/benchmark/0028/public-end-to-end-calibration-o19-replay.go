package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	o "bateia/04_wiring/step0028_gate/o5blackbox"
)

type document struct {
	SchemaVersion string           `json:"schema_version"`
	Artifact      string           `json:"artifact"`
	Entries       []map[string]any `json:"entries"`
}

type source struct {
	SemanticState    string          `json:"semantic_state"`
	ContractInput    json.RawMessage `json:"contract_input"`
	ArtifactDocument document        `json:"artifact_document"`
}

func main() {
	var corpus struct {
		Exchanges []source `json:"exchanges"`
	}
	raw, err := os.ReadFile("00_nucleo/benchmark/0028/public-end-to-end-calibration-o19.json")
	if err != nil || json.Unmarshal(raw, &corpus) != nil {
		panic("invalid public calibration corpus")
	}
	counts := map[o.Classification]int{}
	for _, item := range corpus.Exchanges {
		artifactBytes, err := json.Marshal(item.ArtifactDocument)
		if err != nil {
			panic(err)
		}
		sum := sha256.Sum256(artifactBytes)
		returned, _ := json.Marshal(map[string]any{
			"artifact":       item.ArtifactDocument.Artifact,
			"sha256":         hex.EncodeToString(sum[:]),
			"size":           len(artifactBytes),
			"canonical_json": base64.StdEncoding.EncodeToString(artifactBytes),
			"next_cursor":    nil,
		})
		output, _ := json.Marshal(o.ContractOutput{
			ReturnValue: returned,
			EffectObservation: o.EffectObservation{
				StoreInstanceIDs: []string{}, QueueInstanceIDs: []string{}, RootInstanceIDs: []string{},
				CursorSecretInstanceIDs: []string{}, WorkerInstanceIDs: []string{},
			},
			PublicationObservation: o.PublicationObservation{
				ArtifactNames: []string{}, ArtifactHashes: map[string]string{}, IndexHashes: map[string]string{},
				LeakedFieldNames: []string{}, LayerFindings: []string{},
			},
		})
		finding := o.Evaluate(item.ContractInput, output)
		counts[finding.Classification]++
		fmt.Printf("semantic_state=%s classification=%s\n", item.SemanticState, finding.Classification)
	}
	fmt.Printf("totals Preserved=%d Unknown=%d Violated=%d\n", counts[o.Preserved], counts[o.Unknown], counts[o.Violated])
}
