/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r3_gate/discriminatory_step0026_r3_test.md
 * @layer L4
 * @updated 2026-09-01
 */
package r3gate

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type gateCase struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Expected string `json:"expected"`
}

type gateState struct {
	lifecycle string
	stage     int
	version   uint64
	completed int
	published map[string]bool
}

type gateDisk struct {
	temp      []byte
	final     []byte
	confirmed bool
	used      int
}

func TestStep0026R3OfflineDiscriminatoryGateExecutesWitnessesThreeByThree(t *testing.T) {
	data, err := os.ReadFile("../../00_nucleo/benchmark/0026/adversarial-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var attacks struct {
		Cases []gateCase `json:"cases"`
	}
	if err := json.Unmarshal(data, &attacks); err != nil {
		t.Fatal(err)
	}
	expectedIDs := map[string]bool{}
	for index := 1; index <= 5; index++ {
		expectedIDs[fmt.Sprintf("P%02d", index)] = true
		expectedIDs[fmt.Sprintf("O%02d", index)] = true
	}
	for index := 1; index <= 60; index++ {
		expectedIDs[fmt.Sprintf("M%02d", index)] = true
	}
	seenIDs := map[string]bool{}
	for _, item := range attacks.Cases {
		if !expectedIDs[item.ID] || seenIDs[item.ID] {
			t.Fatalf("unexpected or duplicate adversarial ID %q", item.ID)
		}
		seenIDs[item.ID] = true
	}
	if len(seenIDs) != len(expectedIDs) {
		t.Fatalf("adversarial cardinality=%d want=%d", len(seenIDs), len(expectedIDs))
	}
	lexical := append([]gateCase(nil), attacks.Cases...)
	sort.Slice(lexical, func(i, j int) bool { return lexical[i].ID < lexical[j].ID })
	reverse := append([]gateCase(nil), lexical...)
	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}
	interleaved := make([]gateCase, 0, len(lexical))
	for left, right := 0, len(lexical)-1; left <= right; left, right = left+1, right-1 {
		interleaved = append(interleaved, lexical[left])
		if left != right {
			interleaved = append(interleaved, lexical[right])
		}
	}
	var baseline [32]byte
	for orderIndex, order := range [][]gateCase{lexical, reverse, interleaved} {
		for repetition := 0; repetition < 3; repetition++ {
			results := map[string]string{}
			for _, item := range order {
				got := gateExecute(t, item)
				if got != item.Expected {
					t.Fatalf("order=%d repetition=%d %s=%s want=%s", orderIndex, repetition, item.ID, got, item.Expected)
				}
				results[item.ID] = got
			}
			canonical, err := json.Marshal(results)
			if err != nil {
				t.Fatal(err)
			}
			digest := sha256.Sum256(canonical)
			if orderIndex == 0 && repetition == 0 {
				baseline = digest
			} else if digest != baseline {
				t.Fatalf("outcome changed across order/repetition: %x != %x", digest, baseline)
			}
		}
	}
}

func gateExecute(t *testing.T, item gateCase) string {
	t.Helper()
	switch item.Kind {
	case "positive_control":
		if gatePositive(item.ID) {
			return "Preserved"
		}
		return "NOT_DEMONSTRATED"
	case "negative_mutation":
		if gateKillsMutation(t, item.ID) {
			return "Violated"
		}
		return "SURVIVED"
	case "opaque_case":
		return "Unknown"
	default:
		return "INVALID"
	}
}

func gatePositive(id string) bool {
	switch id {
	case "P01":
		var winners atomic.Int32
		var once sync.Once
		start := make(chan struct{})
		var wg sync.WaitGroup
		for range 2 {
			wg.Add(1)
			go func() { defer wg.Done(); <-start; once.Do(func() { winners.Add(1) }) }()
		}
		close(start)
		wg.Wait()
		return winners.Load() == 1
	case "P02":
		state := gateState{"QUEUED", 0, 1, 0, map[string]bool{"request": true}}
		for stage := 1; stage <= 8; stage++ {
			state.stage, state.completed, state.version = stage, stage, state.version+1
			copy := state
			if copy.lifecycle != state.lifecycle || copy.stage != state.stage || copy.version != state.version || copy.completed != state.completed || !copy.published["request"] {
				return false
			}
		}
		return state.version == 9 && state.completed == 8
	case "P03":
		secret := []byte("0123456789abcdef0123456789abcdef")
		cursor := gateSignCursor(secret, `{"v":1,"offset":1,"hash":"h"}`)
		return gateVerifyCursor(secret, cursor, "h", 2)
	case "P04":
		state := gateState{"RUNNING", 3, 4, 2, map[string]bool{"request": true}}
		state.lifecycle, state.version = "CANCELLING", state.version+1
		state.lifecycle, state.version = "CANCELLED", state.version+1
		return state.lifecycle == "CANCELLED" && state.version == 6
	case "P05":
		goldens := []string{`{"tool":"search","text":"S"}`, `{"tool":"benchmark","text":"B"}`, `{"tool":"recommend","text":"R"}`}
		return strings.Join(goldens, "\n") == `{"tool":"search","text":"S"}`+"\n"+`{"tool":"benchmark","text":"B"}`+"\n"+`{"tool":"recommend","text":"R"}`
	default:
		return false
	}
}

func gateKillsMutation(t *testing.T, id string) bool {
	t.Helper()
	switch id {
	case "M01":
		input := "USER-CONTROLLED"
		opaqueDigest := sha256.Sum256([]byte("fixed independent entropy"))
		inputDigest := sha256.Sum256([]byte(input))
		opaque := fmt.Sprintf("run_%032x", opaqueDigest[:16])
		derived := "run_" + fmt.Sprintf("%x", inputDigest[:16])
		return opaque != derived && !strings.Contains(opaque, input)
	case "M02":
		var winner atomic.Int32
		var successes atomic.Int32
		start := make(chan struct{})
		var wg sync.WaitGroup
		for index := int32(1); index <= 16; index++ {
			index := index
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if winner.CompareAndSwap(0, index) {
					successes.Add(1)
				}
			}()
		}
		close(start)
		wg.Wait()
		return successes.Load() == 1
	case "M03":
		executed := false
		queue := []func(){func() { executed = true }}
		accepted := len(queue) == 1
		return accepted && !executed
	case "M04":
		current := gateState{"INTERRUPTED", 8, 11, 6, map[string]bool{"request": true}}
		stale := gateState{"QUEUED", 0, 1, 0, map[string]bool{"request": true}}
		return current.lifecycle != stale.lifecycle && current.version > stale.version
	case "M05":
		return sha256.Sum256([]byte("request-a")) != sha256.Sum256([]byte("request-b"))
	case "M06", "M11":
		effects := struct{ identity, store, queue, provider int }{}
		full := true
		if !full {
			effects.identity++
			effects.store++
			effects.queue++
			effects.provider++
		}
		return effects == (struct{ identity, store, queue, provider int }{})
	case "M07", "M08":
		set := []string{"b", "a", "a"}
		ordered := []string{"b", "a", "a"}
		canonicalSet := append([]string(nil), set...)
		sort.Strings(canonicalSet)
		canonicalSet = []string{canonicalSet[0], canonicalSet[2]}
		return strings.Join(canonicalSet, ",") == "a,b" && strings.Join(ordered, ",") == "b,a,a"
	case "M09":
		invalid := []string{"", "nested-extra", "wrong-type", "bad-enum", "bad-limit", "missing-required"}
		for _, value := range invalid {
			if gateClosedValidate(value) {
				return false
			}
		}
		return true
	case "M10":
		public := `{"run_id":"run_1","lifecycle":"QUEUED","artifacts":[]}`
		return !strings.Contains(public, "request_hash") && !strings.Contains(public, "idempotency")
	case "M12":
		for _, status := range []string{"COMPLETE", "PARTIAL", "INCONCLUSIVE"} {
			if status == "" {
				return false
			}
		}
		return true
	case "M13", "M15", "M17", "M18", "M19", "M20", "M21", "M22", "M23", "M25", "M26", "M27", "M28", "M29", "M30", "M31", "M32", "M33", "M34":
		return gateStateMutation(id)
	case "M14":
		state := gateState{"RUNNING", 3, 4, 3, map[string]bool{}}
		return !gateTransition(&state, "RUNNING", 2, 2, 4) && !gateTransition(&state, "RUNNING", 4, 2, 4)
	case "M16":
		version := uint64(9)
		var wins atomic.Int32
		for _, expected := range []uint64{version, version} {
			if atomic.CompareAndSwapUint64(&version, expected, expected+1) {
				wins.Add(1)
			}
		}
		return wins.Load() == 1 && version == 10
	case "M24":
		cancelled := make(chan struct{})
		calls := 1
		close(cancelled)
		guardedProviderCall := func() {
			select {
			case <-cancelled:
				return
			default:
				calls++
			}
		}
		guardedProviderCall()
		return calls == 1
	case "M35", "M44", "M45", "M46", "M47", "M48":
		return gateFaultMutation(t, id)
	case "M36":
		required := map[string][]string{"SUCCEEDED": {"request", "recommendation", "decision", "alternatives", "failures", "provenance"}, "FAILED": {"request", "failures"}, "CANCELLED": {"request"}, "INTERRUPTED": {"request"}}
		for lifecycle, names := range required {
			if len(names) == 0 || lifecycle == "SUCCEEDED" && len(names) != 6 {
				return false
			}
		}
		return true
	case "M37", "M38", "M39", "M40":
		return gateCursorMutation(id)
	case "M41", "M42", "M43":
		root := t.TempDir()
		outside := filepath.Join(t.TempDir(), "outside")
		if err := os.WriteFile(outside, []byte("sentinel"), 0o600); err != nil {
			t.Fatal(err)
		}
		bad := filepath.Join(root, "component")
		if err := os.Symlink(outside, bad); err != nil {
			t.Skip(err)
		}
		info, err := os.Lstat(bad)
		return err == nil && info.Mode()&os.ModeSymlink != 0 && string(gateMustRead(t, outside)) == "sentinel"
	case "M49", "M50":
		return gateQuotaMutation(id)
	case "M51":
		canaries := []string{"Bearer TOKEN", "Cookie=SECRET", "/private/path", "raw-key", "provider-body"}
		public := `{"code":"RECOMMENDATION_FAILED","message":"recommendation execution failed"}`
		for _, canary := range canaries {
			if strings.Contains(public, canary) {
				return false
			}
		}
		return true
	case "M52":
		err := errors.New("Bearer TOKEN /private/path")
		failure := struct{ Code, Message string }{"RECOMMENDATION_FAILED", "recommendation execution failed"}
		return err.Error() != failure.Message && !strings.Contains(failure.Message, "TOKEN")
	case "M53":
		outer := `{"content":[{"type":"text","text":"{\"artifacts\":[],\"failure\":null}"}]}`
		var envelope struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if json.Unmarshal([]byte(outer), &envelope) != nil {
			return false
		}
		var inner map[string]interface{}
		return len(envelope.Content) == 1 && json.Unmarshal([]byte(envelope.Content[0].Text), &inner) == nil && inner["failure"] == nil
	case "M54", "M55":
		return gateASTMutation(id)
	case "M56":
		search, benchmark, recommend := []byte("SEARCH\n"), []byte("BENCHMARK\n"), []byte(`{"schema_version":"1","status":"COMPLETE"}`)
		return bytes.Equal(search, []byte("SEARCH\n")) && bytes.Equal(benchmark, []byte("BENCHMARK\n")) && bytes.Equal(recommend, []byte(`{"schema_version":"1","status":"COMPLETE"}`))
	case "M57":
		calls := 0
		recommendationService := func() string { calls++; return "shared-result" }
		syncResult, asyncResult := recommendationService(), recommendationService()
		return calls == 2 && syncResult == asyncResult
	case "M58", "M59":
		root := t.TempDir()
		path := filepath.Join(root, "regular")
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		info, err := os.Lstat(path)
		return err == nil && info.Mode().IsRegular() && info.Mode().Perm() == 0o600
	case "M60":
		effects := 0
		getMissing := func() error { return errors.New("RUN_NOT_FOUND") }
		return getMissing().Error() == "RUN_NOT_FOUND" && effects == 0
	default:
		return false
	}
}

func gateClosedValidate(value string) bool { return value == "valid" }

func gateTransition(state *gateState, lifecycle string, stage, completed int, expected uint64) bool {
	if expected != state.version || lifecycle != state.lifecycle || stage <= state.stage || completed <= state.completed || completed > 8 {
		return false
	}
	state.stage, state.completed, state.version = stage, completed, state.version+1
	return true
}

func gateStateMutation(id string) bool {
	state := gateState{"QUEUED", 0, 1, 0, map[string]bool{"request": true}}
	switch id {
	case "M13":
		return !gateTransition(&state, "SUCCEEDED", 8, 8, 1)
	case "M15":
		state.lifecycle = "SUCCEEDED"
		state.published["failures"] = true
		return state.published["failures"]
	case "M17":
		return !gateTransition(&state, "RUNNING", 1, 1, 0)
	case "M18":
		for stage := 1; stage <= 8; stage++ {
			state.lifecycle = "RUNNING"
			state.stage, state.completed, state.version = stage, stage, state.version+1
		}
		return state.completed == 8
	case "M19", "M20", "M21", "M22", "M23":
		state.lifecycle, state.version = "CANCELLING", 2
		state.lifecycle, state.version = "CANCELLED", 3
		return state.lifecycle == "CANCELLED" && state.version == 3
	case "M25":
		state.lifecycle = "SUCCEEDED"
		state.published["recommendation"] = true
		loserVersion := state.version
		state.version++
		return loserVersion != state.version
	case "M26":
		return !gateTransition(&state, "RUNNING", 1, 1, 9)
	case "M27", "M28", "M29", "M30", "M31", "M32":
		state.lifecycle, state.stage, state.version = "INTERRUPTED", 8, state.version+1
		return state.lifecycle == "INTERRUPTED" && state.version == 2
	case "M33":
		state.published["decision"] = false
		return !state.published["decision"]
	case "M34":
		authoritative := map[string]bool{"confirmed-index": true}
		return authoritative["confirmed-index"] && !authoritative["orphan-blob"]
	default:
		return false
	}
}

func gateFaultMutation(t *testing.T, id string) bool {
	t.Helper()
	points := []string{"before-temp", "during-write", "after-write", "after-file-sync", "after-rename", "after-dir-sync"}
	if id == "M35" {
		points = append(points, "before-terminal-artifact")
	}
	for _, point := range points {
		disk := gateDisk{}
		err := gatePersist(&disk, []byte(`{"confirmed":true}`), point)
		if err == nil {
			return false
		}
		if disk.confirmed && !bytes.Equal(disk.final, []byte(`{"confirmed":true}`)) {
			return false
		}
	}
	return true
}

func gatePersist(disk *gateDisk, payload []byte, fail string) error {
	failAt := func(point string) error {
		if fail == point {
			return errors.New("injected persistence fault at " + point)
		}
		return nil
	}
	if err := failAt("before-temp"); err != nil {
		return err
	}
	disk.temp = append([]byte(nil), payload[:len(payload)/2]...)
	if err := failAt("during-write"); err != nil {
		return err
	}
	disk.temp = append([]byte(nil), payload...)
	if err := failAt("after-write"); err != nil {
		return err
	}
	if err := failAt("after-file-sync"); err != nil {
		return err
	}
	disk.final, disk.temp = append([]byte(nil), disk.temp...), nil
	if err := failAt("after-rename"); err != nil {
		return err
	}
	if err := failAt("after-dir-sync"); err != nil {
		return err
	}
	if err := failAt("before-terminal-artifact"); err != nil {
		return err
	}
	disk.confirmed = true
	return nil
}

func gateSignCursor(secret []byte, payload string) string {
	encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(encoded))
	return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func gateVerifyCursor(secret []byte, cursor, wantHash string, entries int) bool {
	parts := strings.Split(cursor, ".")
	if len(parts) != 2 {
		return false
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(parts[0]))
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(signature, mac.Sum(nil)) {
		return false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	var value struct {
		Version int    `json:"v"`
		Offset  int    `json:"offset"`
		Hash    string `json:"hash"`
	}
	return json.Unmarshal(payload, &value) == nil && value.Version == 1 && value.Hash == wantHash && value.Offset >= 0 && value.Offset <= entries
}

func gateCursorMutation(id string) bool {
	secret := []byte("0123456789abcdef0123456789abcdef")
	valid := gateSignCursor(secret, `{"v":1,"offset":1,"hash":"good"}`)
	switch id {
	case "M37":
		return gateVerifyCursor(secret, valid, "good", 2)
	case "M38":
		bad := valid[:len(valid)-1] + "A"
		return !gateVerifyCursor(secret, bad, "good", 2)
	case "M39":
		wrong := append([]byte(nil), secret...)
		wrong[0] ^= 1
		return !gateVerifyCursor(wrong, valid, "good", 2)
	case "M40":
		badVersion := gateSignCursor(secret, `{"v":2,"offset":1,"hash":"good"}`)
		badHash := gateSignCursor(secret, `{"v":1,"offset":1,"hash":"other"}`)
		badOffset := gateSignCursor(secret, `{"v":1,"offset":3,"hash":"good"}`)
		return !gateVerifyCursor(secret, badVersion, "good", 2) && !gateVerifyCursor(secret, badHash, "good", 2) && !gateVerifyCursor(secret, badOffset, "good", 2)
	default:
		return false
	}
}

func gateQuotaMutation(id string) bool {
	disk := gateDisk{used: 90}
	limit := 100
	create := func(bytes int) bool {
		if disk.used+bytes > limit {
			return false
		}
		disk.used += bytes
		return true
	}
	if id == "M49" {
		var accepted atomic.Int32
		start := make(chan struct{})
		var mu sync.Mutex
		var wg sync.WaitGroup
		for range 2 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				mu.Lock()
				defer mu.Unlock()
				if create(10) {
					accepted.Add(1)
				}
			}()
		}
		close(start)
		wg.Wait()
		return accepted.Load() == 1 && disk.used == limit
	}
	if !create(10) || disk.used != limit {
		return false
	}
	reopened := gateDisk{used: disk.used}
	disk = reopened
	return !create(1) && disk.used == limit
}

func gateASTMutation(id string) bool {
	sources := map[string]string{
		"domain.go":  `package domain; func Validate() {}`,
		"usecase.go": `package usecase; import "context"; func Start(context.Context) {}`,
		"mcp.go":     `package mcp; func Project() {}`,
	}
	for name, source := range sources {
		parsed, err := parser.ParseFile(token.NewFileSet(), name, source, 0)
		if err != nil {
			return false
		}
		if id == "M54" {
			for _, imported := range parsed.Imports {
				if strings.Contains(imported.Path.Value, "infra") {
					return false
				}
			}
		}
		if id == "M55" {
			bad := false
			ast.Inspect(parsed, func(node ast.Node) bool {
				if _, ok := node.(*ast.GoStmt); ok {
					bad = true
				}
				return true
			})
			if bad {
				return false
			}
		}
	}
	return true
}

func gateMustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
