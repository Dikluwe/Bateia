/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/security_input.md
 * @layer L2
 * @updated 2026-09-26
 */
package cli

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"bateia/01_core/parser"
)

const maxIdeaInputBytes = 16 << 20

// Explicit CLI input may be anywhere the user selected; this is a resource
// bound, not the repository confinement boundary. Nonregular files are refused.
func readIdeaInput(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readIdeaDescriptor(f)
}

func readIdeaDescriptor(f *os.File) ([]byte, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxIdeaInputBytes {
		return nil, fmt.Errorf("entrada deve ser arquivo regular de até %d bytes", maxIdeaInputBytes)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxIdeaInputBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxIdeaInputBytes {
		return nil, fmt.Errorf("entrada excedeu limite durante leitura")
	}
	return data, nil
}

// Bound structural work before json.Unmarshal. Quoted delimiters are data.
func validateIdeaJSONStructure(data []byte) error {
	return parser.ValidateBoundedJSON(data)
}
