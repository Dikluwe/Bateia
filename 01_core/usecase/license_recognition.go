/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/license_recognition.md
 * @layer L1
 * @updated 2026-09-26
 */
package usecase

import (
	"regexp"
	"strings"
)

var simpleMITCopyright = regexp.MustCompile(`^copyright (?:\(c\) |© )?[12][0-9]{3}(?:-[12][0-9]{3})? [a-z0-9][a-z0-9_-]{0,127}$`)

// Exact normalized terms, not a substring marker, establish this narrow
// recognition. Other licenses remain advisory until full-text verifiers exist.
const canonicalMITTerms = `Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`

func recognizedLicenseText(content string) string {
	// Titles and copyright holders vary. Only these bounded header lines may
	// precede the complete terms; extra clauses, negation and trailing text fail.
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	start := -1
	headerLines := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "permission is hereby granted,") {
			start = i
			break
		}
		headerLines++
		if headerLines > 10 || len(trimmed) > 512 {
			return ""
		}
		lower := strings.ToLower(trimmed)
		if lower == "mit license" || lower == "the mit license (mit)" || lower == "spdx-license-identifier: mit" {
			continue
		}
		// A free-form holder suffix can conceal a restrictive clause. Recognize
		// only one syntactic holder token; multiword/non-ASCII notices remain
		// advisory rather than attempting legal-language classification.
		if !simpleMITCopyright.MatchString(strings.Join(strings.Fields(lower), " ")) {
			return ""
		}
	}
	if start < 0 {
		return ""
	}
	normalize := func(text string) string { return strings.ToLower(strings.Join(strings.Fields(text), " ")) }
	if normalize(strings.Join(lines[start:], "\n")) != normalize(canonicalMITTerms) {
		return ""
	}
	return "MIT"
}

func detectedLicenseMarkers(content string) []string {
	normalized := strings.ToLower(content)
	var ids []string
	for _, candidate := range []struct {
		id      string
		markers []string
	}{
		{"MIT", []string{"mit license", "permission is hereby granted, free of charge"}},
		{"Apache-2.0", []string{"apache license", "version 2.0"}},
		{"MPL-2.0", []string{"mozilla public license version 2.0", "mpl-2.0"}},
		{"GPL-3.0", []string{"gnu general public license", "version 3"}},
	} {
		if allMarkers(normalized, candidate.markers) {
			ids = append(ids, candidate.id)
		}
	}
	return ids
}
