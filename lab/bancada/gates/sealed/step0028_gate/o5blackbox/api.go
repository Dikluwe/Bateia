/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/api.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import "encoding/json"

type Classification string

const (
	Preserved Classification = "Preserved"
	Violated  Classification = "Violated"
	Unknown   Classification = "Unknown"
)

type Port string

const (
	MechanismExtract Port = "MECHANISM_EXTRACT"
	SynthesisStart   Port = "SYNTHESIS_START"
	SynthesisExecute Port = "SYNTHESIS_EXECUTE"
	RunRead          Port = "RUN_READ"
	ArtifactRead     Port = "ARTIFACT_READ"
	RunCancel        Port = "RUN_CANCEL"
	ToolsList        Port = "TOOLS_LIST"
	Replay           Port = "REPLAY"
	RefinementAudit  Port = "REFINEMENT_AUDIT"
)

type ContractInput struct {
	Port    Port            `json:"port"`
	Request json.RawMessage `json:"request"`
}

type ControlledError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EffectObservation struct {
	NetworkCalls            int      `json:"network_calls"`
	FilesystemReads         int      `json:"filesystem_reads"`
	FilesystemWrites        int      `json:"filesystem_writes"`
	ProcessSpawns           int      `json:"process_spawns"`
	ExternalMutations       int      `json:"external_mutations"`
	ResolveCalls            int      `json:"resolve_calls"`
	ReserveCalls            int      `json:"reserve_calls"`
	IDCalls                 int      `json:"id_calls"`
	CreateCalls             int      `json:"create_calls"`
	StoreInstanceIDs        []string `json:"store_instance_ids"`
	QueueInstanceIDs        []string `json:"queue_instance_ids"`
	RootInstanceIDs         []string `json:"root_instance_ids"`
	CursorSecretInstanceIDs []string `json:"cursor_secret_instance_ids"`
	WorkerInstanceIDs       []string `json:"worker_instance_ids"`
}

type PublicationObservation struct {
	ArtifactNames    []string          `json:"artifact_names"`
	ArtifactHashes   map[string]string `json:"artifact_hashes"`
	IndexHashes      map[string]string `json:"index_hashes"`
	LeakedFieldNames []string          `json:"leaked_field_names"`
	LayerFindings    []string          `json:"layer_findings"`
}

type ContractOutput struct {
	ReturnValue            json.RawMessage        `json:"return_value"`
	Error                  *ControlledError       `json:"error"`
	EffectObservation      EffectObservation      `json:"effect_observation"`
	PublicationObservation PublicationObservation `json:"publication_observation"`
}

type Exchange struct {
	ContractInput  json.RawMessage `json:"contract_input"`
	ContractOutput json.RawMessage `json:"mutant_contract_output"`
}

type Finding struct {
	Classification Classification `json:"classification"`
	Port           Port           `json:"port"`
	Predicates     int            `json:"predicates_evaluated"`
	Failures       []string       `json:"failures"`
	Unknowns       []string       `json:"unknowns"`
}

func Evaluate(contractInput, contractOutput []byte) Finding {
	return evaluate(contractInput, contractOutput)
}

func EvaluateSequence(exchanges []Exchange) []Finding {
	return evaluateSequence(exchanges)
}
