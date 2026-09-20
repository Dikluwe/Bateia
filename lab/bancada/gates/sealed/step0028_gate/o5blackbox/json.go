/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/04_wiring/step0028_gate/o5blackbox/json.md
 * @layer L4
 * @updated 2026-09-02
 */

package o5blackbox

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func decodeStrict(document []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return fmt.Errorf("trailing JSON value")
	}
	return nil
}

func exactTopLevelKeys(document []byte, required ...string) bool {
	var object map[string]json.RawMessage
	if err := decodeStrict(document, &object); err != nil || len(object) != len(required) {
		return false
	}
	for _, key := range required {
		if _, exists := object[key]; !exists {
			return false
		}
	}
	return true
}

func rawIsNull(raw json.RawMessage) bool {
	return len(bytes.TrimSpace(raw)) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

func canonicalJSON(raw []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("trailing JSON value")
	}
	return marshalCanonical(value)
}

func marshalCanonical(value any) ([]byte, error) {
	switch typed := value.(type) {
	case nil, bool, string, json.Number:
		return json.Marshal(typed)
	case []any:
		var buffer bytes.Buffer
		buffer.WriteByte('[')
		for index, item := range typed {
			if index > 0 {
				buffer.WriteByte(',')
			}
			encoded, err := marshalCanonical(item)
			if err != nil {
				return nil, err
			}
			buffer.Write(encoded)
		}
		buffer.WriteByte(']')
		return buffer.Bytes(), nil
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var buffer bytes.Buffer
		buffer.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				buffer.WriteByte(',')
			}
			encodedKey, _ := json.Marshal(key)
			buffer.Write(encodedKey)
			buffer.WriteByte(':')
			encoded, err := marshalCanonical(typed[key])
			if err != nil {
				return nil, err
			}
			buffer.Write(encoded)
		}
		buffer.WriteByte('}')
		return buffer.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported JSON value %T", value)
	}
}

func digestBytes(contents []byte) string {
	digest := sha256.Sum256(contents)
	return hex.EncodeToString(digest[:])
}

func isLowerHex(value string, size int) bool {
	if len(value) != size {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == string(bytes.ToLower([]byte(value)))
}

func isSHA256(value string) bool {
	return len(value) == 71 && value[:7] == "sha256:" && isLowerHex(value[7:], 64)
}
