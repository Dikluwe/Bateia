/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/envfile/load.md
 * @layer L3
 * @updated 2026-09-02
 */

package envfile

import (
	"os"
	"strings"
)

func Load(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if os.Getenv(key) == "" && value != "" {
			_ = os.Setenv(key, value)
		}
	}
}
