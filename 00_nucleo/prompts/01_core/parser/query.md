# Specification: Query Engine & Normalizer (L1 Core)

## Purpose
Sanitize and normalize raw user queries, remove stop words, and expand conceptual terms into structural search keywords.

## Rules
- Pure L1 package (`01_core/parser`).
- ZERO I/O (no file reads, no network requests).
- `NormalizeQuery(raw string) string`: Sanitizes spaces, removes noise.
- `ExpandKeywords(query string, lang string) []string`: Expands architectural concepts to target keywords.
