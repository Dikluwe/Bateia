/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0121/cmd/declarations/main.md
 * @layer lab
 * @updated 2026-10-07
 */
package main

import (
	declarations "bateia/lab/bancada/concorrencia/0121"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("declarations", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "frozen declaration request")
	output := fs.String("output", "", "new output file only")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *input == "" || fs.NArg() != 0 {
		return fmt.Errorf("--input JSON required; --output NEW_FILE optional")
	}
	var request declarations.Request
	if err := read(*input, &request); err != nil {
		return err
	}
	report, err := declarations.Inspect(request)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if *output != "" {
		return publish(*output, data)
	}
	n, err := out.Write(data)
	if err == nil && n != len(data) {
		return io.ErrShortWrite
	}
	return err
}
func read(path string, value any) error { return readLimit(path, value, 64<<20) }
func readLimit(path string, value any, limit int64) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("input must be a regular file")
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return err
	}
	if int64(len(data)) > limit {
		return fmt.Errorf("input exceeds byte limit")
	}
	check := json.NewDecoder(bytes.NewReader(data))
	if err := uniqueJSON(check, 0); err != nil {
		return err
	}
	if _, err := check.Token(); !errors.Is(err, io.EOF) {
		return fmt.Errorf("multiple JSON documents")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(value)
}

func uniqueJSON(d *json.Decoder, depth int) error {
	if depth > 64 {
		return fmt.Errorf("JSON nesting limit")
	}
	t, err := d.Token()
	if err != nil {
		return err
	}
	delimiter, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	keys := map[string]bool{}
	for d.More() {
		if delimiter == '{' {
			t, err := d.Token()
			if err != nil {
				return err
			}
			key, ok := t.(string)
			if !ok || keys[key] {
				return fmt.Errorf("duplicate or invalid JSON key")
			}
			keys[key] = true
		}
		if err := uniqueJSON(d, depth+1); err != nil {
			return err
		}
	}
	_, err = d.Token()
	return err
}

func publish(path string, data []byte) error {
	// Linking a completed temporary file publishes exclusively, without replacing any alias.
	f, err := os.CreateTemp(filepath.Dir(path), ".observer-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Link(f.Name(), path)
}
