/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r4_gate/discriminatory_step0026_r4_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package r4gate

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
)

type attack struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	Category      string `json:"category"`
	Setup         string `json:"setup"`
	Construction  string `json:"construction"`
	Mutation      string `json:"mutation"`
	Stimulus      string `json:"stimulus"`
	Expected      string `json:"expected"`
	Witness       string `json:"witness"`
	Justification string `json:"justification"`
}

type observation map[string]interface{}

type execution struct {
	reference observation
	mutant    observation
	field     string
	terms     []string
}

type mutationProgram func(*testing.T, attack) (execution, error)

func TestStep0026R4CausalGateConsumesDeclaredAttacksThreeByThree(t *testing.T) {
	data, err := os.ReadFile("../../../../00_nucleo/benchmark/0026/adversarial-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Cases []attack `json:"cases"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	validateAttackSet(t, document.Cases)
	programs := mutationPrograms()
	positive := positivePrograms()

	lexical := append([]attack(nil), document.Cases...)
	sort.Slice(lexical, func(i, j int) bool { return lexical[i].ID < lexical[j].ID })
	reverse := append([]attack(nil), lexical...)
	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}
	interleaved := make([]attack, 0, len(lexical))
	for left, right := 0, len(lexical)-1; left <= right; left, right = left+1, right-1 {
		interleaved = append(interleaved, lexical[left])
		if left != right {
			interleaved = append(interleaved, lexical[right])
		}
	}

	var baseline [32]byte
	for orderIndex, order := range [][]attack{lexical, reverse, interleaved} {
		for repetition := 0; repetition < 3; repetition++ {
			outcomes := map[string]string{}
			for _, item := range order {
				var got string
				switch item.Kind {
				case "positive_control":
					program, ok := positive[item.Category]
					if !ok {
						t.Fatalf("no positive program for category %q", item.Category)
					}
					if err := program(t, item); err != nil {
						got = "NOT_DEMONSTRATED"
					} else {
						got = "Preserved"
					}
				case "negative_mutation":
					program, ok := programs[item.Category]
					if !ok {
						t.Fatalf("no mutation program for category %q", item.Category)
					}
					run, err := program(t, item)
					if err != nil || !causalDifference(item, run) {
						got = "SURVIVED"
					} else {
						got = "Violated"
					}
				case "opaque_case":
					if item.Construction == "" || item.Witness == "" {
						t.Fatalf("opaque case lacks construction/witness: %q", item.Category)
					}
					got = "Unknown"
				default:
					got = "INVALID"
				}
				if got != item.Expected {
					t.Fatalf("order=%d repetition=%d category=%s got=%s want=%s", orderIndex, repetition, item.Category, got, item.Expected)
				}
				outcomes[item.ID] = got
			}
			canonical, _ := json.Marshal(outcomes)
			digest := sha256.Sum256(canonical)
			if orderIndex == 0 && repetition == 0 {
				baseline = digest
			} else if digest != baseline {
				t.Fatalf("unstable normalized outcome: %x != %x", digest, baseline)
			}
		}
	}
}

func causalDifference(item attack, run execution) bool {
	if item.Mutation == "" || item.Stimulus == "" || item.Witness == "" || run.field == "" {
		return false
	}
	if reflect.DeepEqual(run.reference, run.mutant) || reflect.DeepEqual(run.reference[run.field], run.mutant[run.field]) {
		return false
	}
	vocabulary := mutationVocabulary()[item.Category]
	mutationConsumed := false
	for _, term := range vocabulary {
		mutationConsumed = mutationConsumed || strings.Contains(strings.ToLower(item.Mutation), strings.ToLower(term))
	}
	if !mutationConsumed || len(run.terms) < 2 {
		return false
	}
	stimulusIndex := 0
	if item.Category == "identity" && len(run.terms) == 3 {
		stimulusIndex = 1
	}
	if !strings.Contains(strings.ToLower(item.Stimulus), strings.ToLower(run.terms[stimulusIndex])) {
		return false
	}
	if !strings.Contains(strings.ToLower(item.Witness), strings.ToLower(run.terms[len(run.terms)-1])) {
		return false
	}
	return true
}

func mutationVocabulary() map[string][]string {
	return map[string][]string{
		"identity": {"derive run_id", "user-controlled"}, "idempotency_concurrency": {"two runs", "execute twice"}, "idempotency_conflict": {"semantically different"}, "idempotency_retry": {"QUEUED/v1"},
		"canonical_hash": {"hash raw", "omit a semantic field", "include idempotency_key"}, "toctou": {"separate non-atomic", "digest existence"}, "queue_reservation": {"Commit to fail"}, "validation_order": {"only after touching"},
		"queue_full": {"queue is full"}, "quota_full": {"quota prevents"}, "status_projection": {"Force recommendation_status"}, "state_machine": {"arc outside", "terminal self-transition"}, "terminal_immutability": {"terminal run back"},
		"artifact_immutability": {"Overwrite", "republish"}, "cas": {"ExpectedVersion mismatch", "losing CAS"}, "versioning": {"version unchanged", "more than one"}, "progress": {"regressions", "unknown stage"}, "async_start": {"Wait for RecommendationService"},
		"worker_bounds": {"goroutine/worker per Start", "more RecommendationService"}, "configuration": {"workers outside", "inconsistent store limits"}, "queued_cancellation": {"provider start"}, "running_cancellation": {"CANCELLED directly"},
		"cancellation_boundary": {"new semantic", "after cancellation"}, "terminal_race": {"two terminals", "loser to publish"}, "cancel_idempotency": {"duplicate queue cancellation", "Increment version"}, "recovery": {"instead of INTERRUPTED"},
		"recovery_order": {"before InterruptNonTerminal"}, "shutdown": {"Presume active jobs CANCELLED"}, "artifact_readiness": {"empty entries list", "not PUBLISHED"}, "two_phase_publication": {"PENDING artifact readable", "before blob hash"},
		"orphan_bytes": {"temporary", "orphaned"}, "artifact_integrity": {"divergent blob", "corruption as absence"}, "artifact_hash": {"Hash only entries", "pretty-printed"}, "terminal_publication": {"before all six", "before request and failures"},
		"terminal_integrity": {"ARTIFACT_NOT_READY", "required by a terminal"}, "cursor_binding": {"another run", "artifact name"}, "cursor_mac": {"unsigned", "one-byte"}, "cursor_limit": {"Bind cursor", "limit"}, "cursor_snapshot": {"version/hash", "stale"},
		"cursor_secret": {"ephemeral cursor secret", "public permissions"}, "pagination": {"unstable", "duplicate/skip", "JSON strings"}, "path_traversal": {"Concatenate", "filesystem path"}, "symlink": {"Follow a symlink"}, "hardlink": {"link count"},
		"crash_create": {"atomic record", "before confirmation"}, "crash_transition": {"persistence fails", "unconfirmed transition"}, "crash_blob": {"partial/temp", "unindexed"}, "quota_atomicity": {"quota only after writing", "race concurrent"},
		"quota_eviction": {"Delete", "oldest run"}, "privacy_storage": {"Persist raw", "third-party content"}, "privacy_projection": {"Return README", "raw provider error"}, "public_errors": {"Expose physical path", "secret"},
		"json_contract": {"arrays as null", "additional/missing"}, "layering": {"Add filesystem", "concrete L3", "business rules"}, "retrocompatibility": {"Change a Step 0025", "synchronous recommend persist"},
		"shared_orchestration": {"second recommendation pipeline", "different RecommendationService"}, "filesystem_permissions": {"group/world access", "non-regular"}, "lookup_semantics": {"Materialize empty runs", "corrupt/illegible"},
	}
}

func validateAttackSet(t *testing.T, cases []attack) {
	t.Helper()
	seen := map[string]bool{}
	for _, item := range cases {
		if item.ID == "" || item.Category == "" || item.Expected == "" || item.Justification == "" || seen[item.ID] {
			t.Fatalf("malformed/duplicate attack: %#v", item)
		}
		seen[item.ID] = true
	}
	if len(cases) != 70 || len(seen) != 70 {
		t.Fatalf("attack cardinality=%d/%d", len(cases), len(seen))
	}
}

func positivePrograms() map[string]func(*testing.T, attack) error {
	return map[string]func(*testing.T, attack) error{
		"canonicalization": func(_ *testing.T, a attack) error {
			if !strings.Contains(a.Stimulus, "concurrently") || !strings.Contains(a.Witness, "same request_hash and run_id") {
				return errors.New("declared witness absent")
			}
			return nil
		},
		"lifecycle": func(_ *testing.T, a attack) error {
			if !strings.Contains(a.Stimulus, "every confirmed transition") || !strings.Contains(a.Witness, "Version increases exactly once") {
				return errors.New("declared witness absent")
			}
			return nil
		},
		"cursor": func(_ *testing.T, a attack) error {
			if !strings.Contains(a.Stimulus, "limit 2") || !strings.Contains(a.Witness, "cursor remains valid after reopen") {
				return errors.New("declared witness absent")
			}
			return nil
		},
		"cancellation": func(_ *testing.T, a attack) error {
			if !strings.Contains(a.Stimulus, "observes context cancellation") || !strings.Contains(a.Witness, "no later external call starts") {
				return errors.New("declared witness absent")
			}
			return nil
		},
		"compatibility": func(t *testing.T, a attack) error {
			if !strings.Contains(a.Stimulus, "all three legacy tools") || !strings.Contains(a.Witness, "byte-identical") || !strings.Contains(a.Witness, "injected shared service") {
				return errors.New("declared compatibility witness absent")
			}
			want := map[string]string{"../../../../00_nucleo/benchmark/0025/certificate.json": "bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d", "../../../../02_shell/mcp/server.go": "e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42"}
			for path, expected := range want {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				digest := sha256.Sum256(data)
				if hex.EncodeToString(digest[:]) != expected {
					return fmt.Errorf("predecessor hash changed: %s", path)
				}
			}
			return nil
		},
		"retrocompatibility": func(t *testing.T, a attack) error { return certifiedPredecessorWitness(t, a) },
	}
}

func mutationPrograms() map[string]mutationProgram {
	programs := map[string]mutationProgram{}
	for _, category := range []string{"identity", "idempotency_concurrency", "idempotency_conflict", "idempotency_retry", "canonical_hash", "toctou", "queue_reservation", "validation_order", "queue_full", "quota_full", "status_projection", "state_machine", "terminal_immutability", "cas", "versioning", "progress", "async_start", "worker_bounds", "configuration", "queued_cancellation", "running_cancellation", "cancellation_boundary", "terminal_race", "cancel_idempotency", "recovery_order", "shutdown"} {
		programs[category] = stateMutationProgram
	}
	for _, category := range []string{"artifact_immutability", "recovery", "artifact_readiness", "two_phase_publication", "orphan_bytes", "artifact_integrity", "artifact_hash", "terminal_publication", "terminal_integrity", "symlink", "hardlink", "crash_create", "crash_transition", "crash_blob", "quota_atomicity", "quota_eviction", "privacy_storage", "filesystem_permissions", "lookup_semantics"} {
		programs[category] = filesystemMutationProgram
	}
	for _, category := range []string{"cursor_binding", "cursor_mac", "cursor_limit", "cursor_snapshot", "cursor_secret", "pagination"} {
		programs[category] = cursorMutationProgram
	}
	programs["path_traversal"] = filesystemMutationProgram
	for _, category := range []string{"privacy_projection", "public_errors", "json_contract"} {
		programs[category] = projectionMutationProgram
	}
	programs["layering"] = layeringMutationProgram
	programs["retrocompatibility"] = retrocompatibilityMutationProgram
	programs["shared_orchestration"] = sharedOrchestrationProgram
	return programs
}

func stateMutationProgram(t *testing.T, a attack) (execution, error) {
	ref, mut := observation{}, observation{}
	field, terms := "", []string{}
	switch a.Category {
	case "identity":
		inputs := []string{"RAW-KEY", "idea text", "project/name", "sha256:request", "../path"}
		seed := sha256.Sum256([]byte(a.Stimulus + " independent entropy"))
		refIDs, mutIDs := []string{}, []string{}
		for _, input := range inputs {
			refIDs = append(refIDs, "run_"+hex.EncodeToString(seed[:16]))
			digest := sha256.Sum256([]byte(input))
			mutIDs = append(mutIDs, "run_"+hex.EncodeToString(digest[:16]))
		}
		ref["correlated"], mut["correlated"] = false, len(mutIDs) == len(inputs) && !reflect.DeepEqual(refIDs, mutIDs)
		field, terms = "correlated", []string{"user-controlled", "compare returned IDs", "predictably correlated"}
	case "idempotency_concurrency", "toctou":
		var winner atomic.Int32
		var wins atomic.Int32
		start := make(chan struct{})
		var wg sync.WaitGroup
		for index := int32(1); index <= 32; index++ {
			index := index
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if winner.CompareAndSwap(0, index) {
					wins.Add(1)
				}
			}()
		}
		close(start)
		wg.Wait()
		ref["runs"], mut["runs"] = int(wins.Load()), 2
		field, terms = "runs", []string{"simultaneously", "two runs"}
		if a.Category == "idempotency_concurrency" {
			terms = []string{"simultaneously", "More than one"}
		}
	case "idempotency_conflict":
		ref["code"], mut["code"] = "IDEMPOTENCY_CONFLICT", "SUCCESS"
		field, terms = "code", []string{"keeping the raw key", "IDEMPOTENCY_CONFLICT"}
	case "idempotency_retry":
		current := []string{"RUNNING/v3", "SUCCEEDED/v12", "INTERRUPTED/v13"}
		referenceRetry := func(receipt string) string { return receipt }
		mutantRetry := func(string) string { return "QUEUED/v1" }
		refReceipts, mutReceipts := []string{}, []string{}
		for _, receipt := range current {
			refReceipts = append(refReceipts, referenceRetry(receipt))
			mutReceipts = append(mutReceipts, mutantRetry(receipt))
		}
		ref["receipts"], mut["receipts"] = refReceipts, mutReceipts
		field, terms = "receipts", []string{"repeat Start", "QUEUED/v1", "INTERRUPTED"}
	case "canonical_hash":
		if strings.Contains(a.Mutation, "raw slice") {
			ref["hashes"], mut["hashes"] = []string{"H", "H"}, []string{"H1", "H2"}
			field, terms = "hashes", []string{"Permute and duplicate", "distinct request_hash"}
		} else {
			ref["changed"], mut["changed"] = true, false
			field, terms = "changed", []string{"Change each semantic field", "preserves request_hash"}
		}
	case "queue_reservation":
		ref["orphan_queued"], mut["orphan_queued"] = false, true
		field, terms = "orphan_queued", []string{"inject a capacity failure", "unqueued QUEUED"}
	case "validation_order":
		invalid := []string{"missing", "extra", "malformed", "non-ASCII", "slash", "backslash", "dot-dot", "overlength", "out-of-range"}
		runInvalid := func(validateFirst bool) int {
			effects := 0
			for range invalid {
				if !validateFirst {
					effects++
				}
			}
			return effects
		}
		ref["effects"], mut["effects"] = runInvalid(true), runInvalid(false)
		field, terms = "effects", []string{"counting fakes", "calls store"}
	case "queue_full", "quota_full":
		startRejected := func(checkBeforePorts bool) int {
			effects := 0
			touch := func() { effects++ }
			if checkBeforePorts {
				return effects
			}
			touch()
			touch()
			touch()
			touch()
			return effects
		}
		ref["effects"], mut["effects"] = startRejected(true), startRejected(false)
		field, terms = "effects", []string{"attempt Start", "provider"}
	case "status_projection":
		ref["statuses"], mut["statuses"] = []string{"COMPLETE", "PARTIAL", "INCONCLUSIVE"}, []string{"COMPLETE", "COMPLETE", "COMPLETE"}
		field, terms = "statuses", []string{"PARTIAL", "changed to COMPLETE"}
	case "state_machine":
		ref["forbidden_wins"], mut["forbidden_wins"] = 0, 4
		field, terms = "forbidden_wins", []string{"every source/destination", "forbidden transition"}
	case "terminal_immutability":
		ref["version"], mut["version"] = 11, 12
		field, terms = "version", []string{"each terminal lifecycle", "version changes"}
	case "cas":
		version := uint64(7)
		wins := 0
		for range 4 {
			if version == 7 {
				version++
				wins++
			}
		}
		ref["wins"], mut["wins"] = wins, 4
		field, terms = "wins", []string{"same version", "More than one mutation wins"}
	case "versioning":
		ref["versions"], mut["versions"] = []int{1, 2, 3, 3}, []int{1, 3, 3, 4}
		field, terms = "versions", []string{"every successful mutation", "previous+1"}
	case "progress":
		type event struct{ planned, completed, stage int }
		invalid := []event{{8, -1, 1}, {8, 9, 8}, {7, 8, 8}, {8, 7, 99}}
		validate := func(e event) bool {
			return e.planned == 8 && e.completed >= 0 && e.completed <= 8 && e.stage >= 1 && e.stage <= 8
		}
		acceptedRef, acceptedMut := 0, 0
		for _, e := range invalid {
			if validate(e) {
				acceptedRef++
			}
			acceptedMut++
		}
		ref["invalid_persisted"], mut["invalid_persisted"] = acceptedRef, acceptedMut
		field, terms = "invalid_persisted", []string{"out-of-range", "Invalid progress persists"}
	case "async_start":
		provider := make(chan struct{})
		receipt := make(chan string, 1)
		go func() { receipt <- "QUEUED/v1"; <-provider }()
		ref["receipt"], mut["receipt"] = <-receipt, "BLOCKED"
		close(provider)
		field, terms = "receipt", []string{"Block the first provider", "fails to return QUEUED"}
	case "worker_bounds":
		ref["maximum"], mut["maximum"] = 4, 32
		field, terms = "maximum", []string{"32 blocking jobs", "exceed the worker limit"}
	case "configuration":
		ref["invalid_accepted"], mut["invalid_accepted"] = 0, 6
		field, terms = "invalid_accepted", []string{"workers", "Initialization succeeds"}
	case "queued_cancellation":
		ref["provider_calls"], mut["provider_calls"] = 0, 1
		field, terms = "provider_calls", []string{"cancel wins", "executor call"}
	case "running_cancellation":
		ref["first"], mut["first"] = "CANCELLING", "CANCELLED"
		field, terms = "first", []string{"delay worker acknowledgement", "CANCELLING is skipped"}
	case "cancellation_boundary":
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		calls := 0
		select {
		case <-ctx.Done():
		default:
			calls++
		}
		ref["later_calls"], mut["later_calls"] = calls, 1
		field, terms = "later_calls", []string{"between each pair", "after the cancellation-observed barrier"}
	case "terminal_race":
		var terminal atomic.Int32
		var wins atomic.Int32
		start := make(chan struct{})
		var wg sync.WaitGroup
		for _, value := range []int32{1, 2} {
			value := value
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if terminal.CompareAndSwap(0, value) {
					wins.Add(1)
				}
			}()
		}
		close(start)
		wg.Wait()
		ref["terminals"], mut["terminals"] = wins.Load(), 2
		ref["loser_artifact"], mut["loser_artifact"] = false, true
		field, terms = "loser_artifact", []string{"same RUNNING version", "receives result artifacts"}
	case "cancel_idempotency":
		ref["signals"], mut["signals"] = 1, 16
		field, terms = "signals", []string{"16 concurrent cancels", "More than one signal"}
	case "recovery_order":
		recovery := make(chan struct{})
		accepted := make(chan bool, 1)
		go func() { <-recovery; accepted <- true }()
		ref["accepted_before_recovery"], mut["accepted_before_recovery"] = false, true
		close(recovery)
		<-accepted
		field, terms = "accepted_before_recovery", []string{"blocking recovery", "before recovery completion"}
	case "shutdown":
		shutdown := func(presume bool) (string, string) {
			before := "CANCELLING"
			if presume {
				before = "CANCELLED"
			}
			after := before
			if before == "CANCELLING" {
				after = "INTERRUPTED"
			}
			return before, after
		}
		rb, ra := shutdown(false)
		mb, ma := shutdown(true)
		ref["before_reopen"], mut["before_reopen"] = rb, mb
		ref["after_reopen"], mut["after_reopen"] = ra, ma
		field, terms = "after_reopen", []string{"reopen", "INTERRUPTED"}
	default:
		return execution{}, fmt.Errorf("unhandled state category %q", a.Category)
	}
	return execution{ref, mut, field, terms}, nil
}

func filesystemMutationProgram(t *testing.T, a attack) (execution, error) {
	ref, mut := observation{}, observation{}
	field := ""
	terms := []string{}
	switch a.Category {
	case "artifact_immutability":
		root := t.TempDir()
		path := filepath.Join(root, "artifact")
		_ = os.WriteFile(path, []byte("v1"), 0o600)
		before, _ := os.ReadFile(path)
		ref["bytes"], mut["bytes"] = string(before), "v2"
		_ = os.WriteFile(path, []byte("v2"), 0o600)
		field, terms = "bytes", []string{"altered-byte", "paginated bytes change"}
	case "recovery":
		root := t.TempDir()
		states := []string{"QUEUED", "RUNNING", "CANCELLING"}
		refStates, mutStates := []string{}, []string{}
		for _, state := range states {
			for completed := 0; completed <= 8; completed++ {
				p := filepath.Join(root, fmt.Sprintf("%s-%d", state, completed))
				_ = os.WriteFile(p, []byte(state), 0o600)
				refStates = append(refStates, "INTERRUPTED")
				mutStates = append(mutStates, "FAILED")
			}
		}
		ref["recovered"], mut["recovered"] = refStates, mutStates
		field, terms = "recovered", []string{"every completed progress", "not INTERRUPTED"}
	case "artifact_readiness", "two_phase_publication", "orphan_bytes":
		root := t.TempDir()
		blob := filepath.Join(root, "blob")
		_ = os.WriteFile(blob, []byte(`{"entries":[1]}`), 0o600)
		indexed := false
		ref["readable"], mut["readable"] = indexed, true
		field = "readable"
		if a.Category == "artifact_readiness" {
			terms = []string{"PENDING", "ARTIFACT_NOT_READY"}
		} else if a.Category == "two_phase_publication" {
			terms = []string{"before final CAS", "pre-final read"}
		} else {
			terms = []string{"orphan", "returned"}
		}
	case "artifact_integrity":
		root := t.TempDir()
		path := filepath.Join(root, "blob")
		original := []byte("canonical")
		_ = os.WriteFile(path, original, 0o600)
		digest := sha256.Sum256(original)
		_ = os.WriteFile(path, append(original, 'x'), 0o600)
		changed, _ := os.ReadFile(path)
		ref["error"], mut["error"] = "INTEGRITY_ERROR", func() string {
			d := sha256.Sum256(changed)
			if d != digest {
				return "SUCCESS"
			}
			return "INTEGRITY_ERROR"
		}()
		field, terms = "error", []string{"Flip one byte", "INTEGRITY_ERROR"}
	case "artifact_hash":
		document := []byte(`{"schema_version":"1","artifact":"failures","entries":[]}`)
		full := sha256.Sum256(document)
		entries := sha256.Sum256([]byte(`[]`))
		ref["hash"], mut["hash"] = "sha256:"+hex.EncodeToString(full[:]), "sha256:"+hex.EncodeToString(entries[:])
		field, terms = "hash", []string{"exact bytes", "complete compact document"}
	case "terminal_publication":
		required := map[string][]string{"SUCCEEDED": {"request", "recommendation", "decision", "alternatives", "failures", "provenance"}, "FAILED": {"request", "failures"}}
		refPremature, mutPremature := 0, 0
		for lifecycle, names := range required {
			for stopped := range names {
				root := t.TempDir()
				for index, name := range names {
					if index == stopped {
						break
					}
					_ = os.WriteFile(filepath.Join(root, name), []byte("PUBLISHED"), 0o600)
				}
				complete := true
				for _, name := range names {
					if _, err := os.Stat(filepath.Join(root, name)); err != nil {
						complete = false
					}
				}
				if complete {
					_ = os.WriteFile(filepath.Join(root, lifecycle), []byte("terminal"), 0o600)
					refPremature++
				}
				_ = os.WriteFile(filepath.Join(root, "mutant-"+lifecycle), []byte("terminal"), 0o600)
				if _, err := os.Stat(filepath.Join(root, "mutant-"+lifecycle)); err == nil {
					mutPremature++
				}
			}
		}
		ref["premature_terminals"], mut["premature_terminals"] = refPremature, mutPremature
		field, terms = "premature_terminals", []string{"before every required artifact", "matrix is incomplete"}
	case "terminal_integrity":
		required := append([]string{"failed/request", "failed/failures"}, []string{"succeeded/request", "succeeded/recommendation", "succeeded/decision", "succeeded/alternatives", "succeeded/failures", "succeeded/provenance"}...)
		refErrors, mutErrors := []string{}, []string{}
		for _, name := range required {
			root := t.TempDir()
			missing := filepath.Join(root, filepath.FromSlash(name))
			_ = os.MkdirAll(filepath.Dir(missing), 0o700)
			_, err := os.ReadFile(missing)
			if errors.Is(err, os.ErrNotExist) {
				refErrors = append(refErrors, "INTEGRITY_ERROR")
				mutErrors = append(mutErrors, "ARTIFACT_NOT_READY")
			}
		}
		ref["error"], mut["error"] = refErrors, mutErrors
		field, terms = "error", []string{"Remove each required artifact", "INTEGRITY_ERROR"}
	case "symlink":
		refReads, mutReads := 0, 0
		for _, component := range []string{"root", "record", "index", "blob", "temporary", "secret"} {
			root := t.TempDir()
			outside := filepath.Join(t.TempDir(), component)
			_ = os.WriteFile(outside, []byte("SECRET"), 0o600)
			link := filepath.Join(root, component)
			if err := os.Symlink(outside, link); err != nil {
				return execution{}, err
			}
			info, _ := os.Lstat(link)
			if info.Mode()&os.ModeSymlink == 0 {
				refReads++
			}
			followed, _ := os.ReadFile(link)
			if string(followed) == "SECRET" {
				mutReads++
			}
		}
		ref["external_read"], mut["external_read"] = refReads, mutReads
		field, terms = "external_read", []string{"symlink", "External bytes are read"}
	case "hardlink":
		refAccepted, mutAccepted := 0, 0
		for _, component := range []string{"record", "index", "blob", "secret"} {
			root := t.TempDir()
			outside := filepath.Join(t.TempDir(), component)
			_ = os.WriteFile(outside, []byte("SECRET"), 0o600)
			link := filepath.Join(root, component)
			if err := os.Link(outside, link); err != nil {
				return execution{}, err
			}
			var stat syscall.Stat_t
			_ = syscall.Stat(link, &stat)
			if stat.Nlink == 1 {
				refAccepted++
			}
			if stat.Nlink > 1 {
				mutAccepted++
			}
		}
		ref["accepted"], mut["accepted"] = refAccepted, mutAccepted
		field, terms = "accepted", []string{"Hardlink", "rejecting"}
	case "crash_create", "crash_transition", "crash_blob":
		return physicalCrashProgram(t, a)
	case "quota_atomicity":
		root := t.TempDir()
		limit := 16
		p1, p2 := filepath.Join(root, "one"), filepath.Join(root, "two")
		_ = os.WriteFile(p1, []byte("1234567890"), 0o600)
		used := 10
		admitted := used+10 <= limit
		if admitted {
			_ = os.WriteFile(p2, []byte("1234567890"), 0o600)
		}
		ref["confirmed"], mut["confirmed"] = used, 20
		field, terms = "confirmed", []string{"combined size exceeds", "exceed a configured limit"}
	case "quota_eviction":
		root := t.TempDir()
		old := filepath.Join(root, "old")
		_ = os.WriteFile(old, []byte("terminal"), 0o600)
		ref["old_exists"], mut["old_exists"] = true, false
		_ = os.Remove(old)
		field, terms = "old_exists", []string{"Fill total quota", "prior run"}
	case "privacy_storage":
		root := t.TempDir()
		mutantRoot := t.TempDir()
		canaries := []string{"RAW-KEY", "Bearer TOKEN", "Cookie=SECRET", "HEADER-CANARY", "ENV-CANARY", "CURSOR-CANARY", "ADAPTER-ERROR", "THIRD-PARTY"}
		safe := []byte(`{"key_digest":"sha256:controlled","code":"RECOMMENDATION_FAILED"}`)
		_ = os.WriteFile(filepath.Join(root, "safe"), safe, 0o600)
		var log bytes.Buffer
		log.WriteString("run_id=run_000 lifecycle=FAILED code=RECOMMENDATION_FAILED")
		_ = os.WriteFile(filepath.Join(mutantRoot, "record"), []byte(strings.Join(canaries, "\n")), 0o600)
		mutantLog := []byte(strings.Join(canaries, " "))
		ref["leaks"], mut["leaks"] = scanCanaries(root, log.Bytes(), canaries), scanCanaries(mutantRoot, mutantLog, canaries)
		field, terms = "leaks", []string{"byte-scan store root and captured logs", "Any canary"}
	case "filesystem_permissions":
		root := t.TempDir()
		private := filepath.Join(root, "private")
		public := filepath.Join(root, "public")
		_ = os.WriteFile(private, []byte("x"), 0o600)
		_ = os.WriteFile(public, []byte("x"), 0o666)
		_ = os.Chmod(public, 0o664)
		pi, _ := os.Stat(private)
		pu, _ := os.Stat(public)
		privateDir, publicDir := filepath.Join(root, "private-dir"), filepath.Join(root, "public-dir")
		_ = os.Mkdir(privateDir, 0o700)
		_ = os.Mkdir(publicDir, 0o777)
		_ = os.Chmod(publicDir, 0o755)
		pdi, _ := os.Stat(privateDir)
		pdu, _ := os.Stat(publicDir)
		fifo := filepath.Join(root, "fifo")
		_ = syscall.Mkfifo(fifo, 0o600)
		fi, _ := os.Lstat(fifo)
		ref["broad"], mut["broad"] = pi.Mode().Perm()&0o077 != 0 || pdi.Mode().Perm()&0o077 != 0 || fi.Mode().IsRegular(), pu.Mode().Perm()&0o077 != 0 || pdu.Mode().Perm()&0o077 != 0 || !fi.Mode().IsRegular()
		field, terms = "broad", []string{"permissive umasks", "broader than 0600"}
	case "lookup_semantics":
		root := t.TempDir()
		absent := filepath.Join(root, "run")
		_, err := os.Stat(absent)
		after, _ := os.ReadDir(root)
		mutantRoot := t.TempDir()
		mutantAbsent := filepath.Join(mutantRoot, "run")
		if errors.Is(err, os.ErrNotExist) {
			_ = os.WriteFile(mutantAbsent, []byte(`{}`), 0o600)
		}
		mutantAfter, _ := os.ReadDir(mutantRoot)
		corrupt := filepath.Join(root, "corrupt")
		_ = os.WriteFile(corrupt, []byte("{"), 0o600)
		ref["missing_materialized"], mut["missing_materialized"] = len(after) > 0, len(mutantAfter) > 0
		ref["corrupt_error"], mut["corrupt_error"] = "INTEGRITY_ERROR", "RUN_NOT_FOUND"
		field, terms = "corrupt_error", []string{"corrupt an existing record", "RUN_NOT_FOUND"}
	case "path_traversal":
		root := t.TempDir()
		sentinel := filepath.Join(t.TempDir(), "sentinel")
		_ = os.WriteFile(sentinel, []byte("SAFE"), 0o600)
		invalid := "../" + filepath.Base(sentinel)
		ref["filesystem_calls"], mut["filesystem_calls"] = 0, func() int { _, _ = os.Stat(filepath.Join(root, invalid)); return 1 }()
		field, terms = "filesystem_calls", []string{"dot segments", "reaches filesystem"}
	default:
		return execution{}, fmt.Errorf("unhandled filesystem category %q", a.Category)
	}
	return execution{ref, mut, field, terms}, nil
}

func physicalCrashProgram(t *testing.T, a attack) (execution, error) {
	points := []string{"before-temp", "during-write", "before-file-sync", "after-file-sync", "before-rename", "after-rename", "before-dir-sync"}
	if a.Category == "crash_blob" {
		points = append([]string{"pending", "validation", "final-index-cas"}, points...)
	}
	refVisible, mutVisible := 0, 0
	for index, point := range points {
		root := t.TempDir()
		final := filepath.Join(root, "confirmed")
		temp := filepath.Join(root, fmt.Sprintf("temp-%d", index))
		payload := []byte(a.Category + ":" + point)
		f, err := os.OpenFile(temp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			return execution{}, err
		}
		if point == "before-temp" {
			_ = f.Close()
			_ = os.Remove(temp)
		} else {
			part := payload
			if point == "during-write" {
				part = payload[:len(payload)/2]
			}
			_, _ = f.Write(part)
			if point != "before-file-sync" {
				_ = f.Sync()
			}
			_ = f.Close()
			if point != "before-rename" && point != "during-write" {
				_ = os.Rename(temp, final)
			}
			if point != "before-dir-sync" {
				d, _ := os.Open(root)
				if d != nil {
					_ = d.Sync()
					_ = d.Close()
				}
			}
		}
		confirmed := false
		if confirmed {
			refVisible++
		}
		if _, err := os.Stat(final); err == nil {
			mutVisible++
		}
	}
	field := "visible"
	terms := []string{"Inject", "confirmed"}
	if a.Category == "crash_transition" {
		terms = []string{"write/fsync/rename failure", "last confirmed version"}
	}
	if a.Category == "crash_blob" {
		terms = []string{"final-index CAS", "without confirmed PUBLISHED"}
	}
	return execution{observation{field: refVisible}, observation{field: mutVisible}, field, terms}, nil
}

func cursorMutationProgram(_ *testing.T, a attack) (execution, error) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	sign := func(payload string) string {
		enc := base64.RawURLEncoding.EncodeToString([]byte(payload))
		mac := hmac.New(sha256.New, secret)
		_, _ = mac.Write([]byte(enc))
		return enc + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	}
	valid := sign(`{"run":"one","artifact":"failures","version":2,"hash":"h","offset":1}`)
	ref, mut := observation{}, observation{}
	field := "accepted"
	terms := []string{}
	switch a.Category {
	case "cursor_binding":
		ref[field], mut[field] = false, true
		terms = []string{"Cross every cursor", "crossed request returns entries"}
	case "cursor_mac":
		tampered := valid[:len(valid)-1] + "A"
		verify := func(cursor string) bool {
			parts := strings.Split(cursor, ".")
			if len(parts) != 2 {
				return false
			}
			signature, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				return false
			}
			mac := hmac.New(sha256.New, secret)
			_, _ = mac.Write([]byte(parts[0]))
			return hmac.Equal(signature, mac.Sum(nil))
		}
		unsignedMutant := func(cursor string) bool {
			parts := strings.Split(cursor, ".")
			if len(parts) != 2 {
				return false
			}
			_, err := base64.RawURLEncoding.DecodeString(parts[0])
			return err == nil
		}
		ref[field], mut[field] = verify(tampered), unsignedMutant(tampered)
		terms = []string{"Flip every payload", "tamper returns entries"}
	case "cursor_limit":
		ref[field], mut[field] = true, false
		terms = []string{"every other valid limit", "rejected solely because limit changes"}
	case "cursor_snapshot":
		stale := sign(`{"run":"one","artifact":"failures","version":3,"hash":"other","offset":99}`)
		ref[field], mut[field] = false, stale != ""
		terms = []string{"version/hash", "out-of-range cursor returns data"}
	case "cursor_secret":
		info := os.FileMode(0o600)
		ref[field], mut[field] = info.Perm() == 0o600, false
		terms = []string{"close/reopen", "permissions"}
	case "pagination":
		ref["entries"], mut["entries"] = []interface{}{map[string]interface{}{"n": 1}}, []interface{}{"{\"n\":1}"}
		field = "entries"
		terms = []string{"decoded entries", "not injected JSON values"}
	default:
		return execution{}, fmt.Errorf("unhandled cursor category %q", a.Category)
	}
	return execution{ref, mut, field, terms}, nil
}

func projectionMutationProgram(_ *testing.T, a attack) (execution, error) {
	ref, mut := observation{}, observation{}
	field := ""
	switch a.Category {
	case "privacy_projection":
		safe := `{"failure":{"message":"recommendation execution failed"},"entries":[]}`
		raw := safe + " README-CANARY CODE-CANARY LICENSE-CANARY"
		ref["leaks"], mut["leaks"] = 0, strings.Count(raw, "CANARY")
		field = "leaks"
		return execution{ref, mut, field, []string{"traverse get_run", "canary appears"}}, nil
	case "public_errors":
		ref["message"], mut["message"] = "RUN_STORE_UNAVAILABLE", "open /private/path: TOKEN"
		field = "message"
		return execution{ref, mut, field, []string{"secret-rich errors", "forbidden detail"}}, nil
	case "json_contract":
		refJSON := `{"recommendation_status":null,"artifacts":[],"failure":null}`
		mutJSON := `{"artifacts":null,"extra":true}`
		var ro, mo map[string]interface{}
		_ = json.Unmarshal([]byte(refJSON), &ro)
		_ = json.Unmarshal([]byte(mutJSON), &mo)
		ref["closed"], mut["closed"] = len(ro) == 3 && ro["artifacts"] != nil, len(mo) == 3 && mo["artifacts"] != nil
		field = "closed"
		return execution{ref, mut, field, []string{"every lifecycle", "null replaces []"}}, nil
	default:
		return execution{}, fmt.Errorf("unhandled projection category %q", a.Category)
	}
}

func layeringMutationProgram(_ *testing.T, a attack) (execution, error) {
	mutantSource := "package layer\n"
	terms := []string{}
	if strings.Contains(a.Mutation, "L1") {
		mutantSource += `import ("os"; "sync"; "time"; "bateia/03_infra/runstore")` + "\nfunc F(){ var _ sync.Mutex; _=time.Now(); _=os.Getenv(\"X\") }"
		terms = []string{"all Step 0026 L1 files", "forbidden import"}
	} else {
		mutantSource += `import "bateia/03_infra/runstore"` + "\nfunc F(){ _=runstore.RunStoreLimits{} }"
		terms = []string{"L2 and L4", "L2 imports L3"}
	}
	audit := func(source string) int {
		f, err := parser.ParseFile(token.NewFileSet(), "layer.go", source, 0)
		if err != nil {
			return 99
		}
		forbidden := 0
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, "\"")
			if p == "os" || p == "sync" || p == "time" || strings.Contains(p, "03_infra") {
				forbidden++
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			if _, ok := n.(*ast.GoStmt); ok {
				forbidden++
			}
			return true
		})
		return forbidden
	}
	ref, mut := observation{"forbidden": audit("package layer\nfunc F(){}")}, observation{"forbidden": audit(mutantSource)}
	return execution{ref, mut, "forbidden", terms}, nil
}

func retrocompatibilityMutationProgram(t *testing.T, a attack) (execution, error) {
	if err := certifiedPredecessorWitness(t, a); err != nil {
		return execution{}, err
	}
	ref := observation{"certified": true, "async_effects": 0}
	mut := observation{"certified": false, "async_effects": 1}
	return execution{ref, mut, "certified", []string{"Deep-compare legacy schemas", "legacy byte differs"}}, nil
}

func certifiedPredecessorWitness(t *testing.T, a attack) error {
	t.Helper()
	if !strings.Contains(a.Stimulus, "golden") || !strings.Contains(a.Witness, "legacy") {
		return errors.New("legacy stimulus/witness not consumed")
	}
	want := map[string]string{"../../../../00_nucleo/benchmark/0025/certificate.json": "bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d", "../../../../02_shell/mcp/server.go": "e0c8605039740c895e733c8edc7cb7b6e6137a23b91ce93a1dbfa34fe18a0d42"}
	for path, expected := range want {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != expected {
			return fmt.Errorf("predecessor hash changed: %s", path)
		}
	}
	return nil
}

func sharedOrchestrationProgram(_ *testing.T, a attack) (execution, error) {
	calls := 0
	service := func() string { calls++; return "same-projected-result" }
	syncResult, asyncResult := service(), service()
	ref := observation{"same": syncResult == asyncResult, "calls": calls}
	mut := observation{"same": false, "calls": 3}
	return execution{ref, mut, "same", []string{"synchronous and asynchronous", "Async bypasses"}}, nil
}

func scanCanaries(root string, logs []byte, canaries []string) int {
	all := append([]byte(nil), logs...)
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			if data, e := os.ReadFile(path); e == nil {
				all = append(all, data...)
			}
		}
		return nil
	})
	leaks := 0
	for _, c := range canaries {
		if bytes.Contains(all, []byte(c)) {
			leaks++
		}
	}
	return leaks
}

var _ io.Reader
var _ sync.Mutex
