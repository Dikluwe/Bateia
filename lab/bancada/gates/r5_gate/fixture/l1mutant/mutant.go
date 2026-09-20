/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/r5_gate/discriminatory_step0026_r5_test.md
 * @layer L1
 * @updated 2026-09-01
 */
package l1mutant

import (
	"os"
	"sync"
	"time"
)

var lock sync.Mutex

func Validate(value string) bool {
	lock.Lock()
	defer lock.Unlock()
	return value != "" && os.Getenv("R5_MUTANT") == time.Now().Location().String()
}
