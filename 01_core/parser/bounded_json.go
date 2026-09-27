/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/bounded_json.md
 * @layer L1
 * @updated 2026-09-26
 */
package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxStructuredJSONBytes  = 16 << 20
	maxStructuredJSONDepth  = 64
	maxStructuredJSONTokens = 100000
)

// ValidateBoundedJSON validates exactly one JSON value before typed decoding.
// Object members may not repeat, including Unicode case-fold aliases: Go's
// struct decoder accepts case-insensitive fields whereas other consumers may
// interpret them differently. This intentionally rejects some valid JSON.
func ValidateBoundedJSON(data []byte) error {
	if len(data) > MaxStructuredJSONBytes || !utf8.Valid(data) {
		return fmt.Errorf("JSON excede orçamento de bytes ou contém UTF-8 inválido")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	tokens := 0
	next := func() (json.Token, error) {
		if tokens >= maxStructuredJSONTokens {
			return nil, fmt.Errorf("JSON excede orçamento de tokens")
		}
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("JSON inválido ou incompleto")
		}
		tokens++
		return token, nil
	}
	var value func(int) error
	value = func(depth int) error {
		token, err := next()
		if err != nil {
			return err
		}
		delimiter, composite := token.(json.Delim)
		if !composite {
			return nil
		}
		if depth >= maxStructuredJSONDepth {
			return fmt.Errorf("JSON excede profundidade máxima")
		}
		switch delimiter {
		case '{':
			keys := map[string]bool{}
			for decoder.More() {
				keyToken, err := next()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return fmt.Errorf("JSON contém membro inválido")
				}
				folded := foldJSONKey(key)
				if keys[folded] {
					return fmt.Errorf("JSON contém chaves duplicadas ou aliases de capitalização")
				}
				keys[folded] = true
				if err := value(depth + 1); err != nil {
					return err
				}
			}
		case '[':
			for decoder.More() {
				if err := value(depth + 1); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("JSON contém delimitador inesperado")
		}
		end, err := next()
		if err != nil {
			return err
		}
		if delimiter == '{' && end != json.Delim('}') || delimiter == '[' && end != json.Delim(']') {
			return fmt.Errorf("JSON contém fechamento inválido")
		}
		return nil
	}
	if err := value(0); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("JSON deve conter exatamente um valor")
	}
	return nil
}

// Choose the same canonical rune for an entire unicode.SimpleFold cycle,
// including sigma/final-sigma and Kelvin K. strings.ToLower alone misses some
// aliases recognized by encoding/json's field matching.
func foldJSONKey(key string) string {
	var folded strings.Builder
	for _, r := range key {
		minimum := r
		for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
			if next < minimum {
				minimum = next
			}
		}
		folded.WriteRune(minimum)
	}
	return folded.String()
}
