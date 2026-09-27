# Specification: Mermaid Diagram Generator (L2 Shell)

## Purpose
Render domain SearchReport as a Mermaid.js (`graph TD`) dependency graph for visualization in Markdown and HTML presentations.

## Rules
- L2 package (`02_shell/report`).
- Pure formatting of `domain.SearchReport`.
- Function: `GenerateMermaidGraph(report *domain.SearchReport) string`.
- Returns valid Mermaid `graph TD` syntax.
- Escape label delimiters, entity markers, markup and Markdown fences with Mermaid numeric entities, not HTML escaping. Normalize controls/newlines to spaces so data cannot introduce directives or edges.
- Node IDs are generated from repository/dependency positions, never input labels. Ordinary text remains readable.
- HTML separately escapes the resulting Mermaid source. It does not execute a Mermaid runtime; external consumers must use strict Mermaid security configuration.
