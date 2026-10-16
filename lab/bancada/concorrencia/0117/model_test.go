/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0117/model_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package phaseexperiment

import (
	"encoding/json"
	"testing"
	"testing/synctest"
)

type witness struct {
	Scenario               string `json:"scenario"`
	Scope                  string `json:"scope"`
	Accepted               bool   `json:"accepted"`
	Initialized            bool   `json:"initialized"`
	Limit                  int64  `json:"limit"`
	ArtificialInterleaving bool   `json:"artificial_interleaving"`
}

func record(t *testing.T, m *Model, scenario string, accepted, interleaving bool) {
	t.Helper()
	initialized, limit := m.Snapshot()
	data, err := json.Marshal(witness{scenario, "AUTHOR_MODEL_ONLY", accepted, initialized, limit, interleaving})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("WITNESS " + string(data))
}
func rejected(fn func()) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			if r != "already initialized" {
				panic(r)
			}
			result = true
		}
	}()
	fn()
	return false
}
func TestConfigureBeforeInitialize(t *testing.T) {
	var m Model
	m.Configure(3, nil)
	m.Initialize()
	initialized, limit := m.Snapshot()
	if !initialized || limit != 3 {
		t.Fatal("sequential configuration failed")
	}
	record(t, &m, "CONFIGURE_THEN_INITIALIZE", true, false)
}
func TestConfigureAfterInitialize(t *testing.T) {
	var m Model
	m.Configure(3, nil)
	m.Initialize()
	if !rejected(func() { m.Configure(7, nil) }) {
		t.Fatal("guard did not reject")
	}
	initialized, limit := m.Snapshot()
	if !initialized || limit != 3 {
		t.Fatal("rejection changed limit")
	}
	record(t, &m, "INITIALIZE_THEN_CONFIGURE", false, false)
}
func TestInitializeBetweenCheckAndWrite(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var m Model
		checked := make(chan struct{})
		resume := make(chan struct{})
		done := make(chan bool, 1)
		go func() {
			wasRejected := rejected(func() { m.Configure(7, func() { close(checked); <-resume }) })
			done <- wasRejected
		}()
		synctest.Wait()
		select {
		case <-checked:
		default:
			t.Fatal("writer did not rendezvous after check")
		}
		initialized, limit := m.Snapshot()
		if initialized || limit != 0 {
			t.Fatal("write occurred before release")
		}
		m.Initialize()
		close(resume)
		if <-done {
			t.Fatal("model unexpectedly rechecked initialization")
		}
		initialized, limit = m.Snapshot()
		if !initialized || limit != 7 {
			t.Fatal("interleaving not witnessed")
		}
		record(t, &m, "CHECK_THEN_INITIALIZE_THEN_WRITE", true, true)
	})
}
func TestInitializeBeforeConcurrentCheck(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var m Model
		m.Initialize()
		done := make(chan bool, 1)
		called := false
		go func() { done <- rejected(func() { m.Configure(7, func() { called = true }) }) }()
		if !<-done || called {
			t.Fatal("prior initialization did not block before hook")
		}
		initialized, limit := m.Snapshot()
		if !initialized || limit != 0 {
			t.Fatal("rejected write changed model")
		}
		record(t, &m, "INITIALIZE_BEFORE_CONCURRENT_CHECK", false, false)
	})
}
