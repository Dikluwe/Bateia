/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/cmd/observer/main.md
 * @layer lab
 * @updated 2026-10-06
 */
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	concurrency "bateia/lab/bancada/concorrencia/0109"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 || args[0] != "observe" && args[0] != "evaluate" {
		return fmt.Errorf("usage: observer observe --input JSON | evaluate --corpus JSON [--output NEW_FILE]")
	}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "frozen observation request")
	corpus := fs.String("corpus", "", "exploratory corpus")
	environment := fs.String("environment", "", "supplied frozen type environment")
	output := fs.String("output", "", "new output file only")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected arguments")
	}
	var value any
	failed := false
	switch args[0] {
	case "observe":
		if *input == "" || *corpus != "" {
			return fmt.Errorf("observe requires --input only")
		}
		var request concurrency.Request
		if err := read(*input, &request); err != nil {
			return err
		}
		var r concurrency.Report
		var err error
		if *environment != "" {
			var env concurrency.Environment
			if err = readLimit(*environment, &env, 64<<20); err != nil {
				return err
			}
			r, err = concurrency.ObserveWithEnvironment(request, env)
		} else {
			r, err = concurrency.Observe(request)
		}
		if err != nil {
			return err
		}
		value = r
	case "evaluate":
		if *corpus == "" || *input != "" || *environment != "" {
			return fmt.Errorf("evaluate requires --corpus only")
		}
		var request concurrency.Corpus
		if err := read(*corpus, &request); err != nil {
			return err
		}
		r, err := concurrency.Evaluate(request)
		if err != nil {
			return err
		}
		value = r
		failed = r.Passed != r.Cases
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if *output != "" {
		err = publish(*output, data)
	} else {
		var n int
		n, err = out.Write(data)
		if err == nil && n != len(data) {
			err = io.ErrShortWrite
		}
	}
	if err != nil {
		return err
	}
	if failed {
		return fmt.Errorf("exploratory expectations diverged; inspect evaluation report")
	}
	return nil
}

func read(path string, value any) error { return readLimit(path, value, 2<<20) }
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
