/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/jsonpath_extractor.md
 * @layer L1
 * @updated 2026-08-08
 */

package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// QueryJSONPath avalia uma busca de campo estilo JSONPath (ex: "$.dependencies.express" ou "dependencies.express") em L1.
func QueryJSONPath(jsonContent string, path string) (string, bool) {
	if jsonContent == "" || path == "" {
		return "", false
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonContent), &data); err != nil {
		return "", false
	}

	cleanPath := strings.TrimPrefix(path, "$.")
	parts := strings.Split(cleanPath, ".")

	var current interface{} = data
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return "", false
		}
		val, exists := m[part]
		if !exists {
			return "", false
		}
		current = val
	}

	return fmt.Sprintf("%v", current), true
}
