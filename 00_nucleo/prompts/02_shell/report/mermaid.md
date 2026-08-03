# Specification: Mermaid Diagram Generator (L2 Shell)

## Purpose
Render domain SearchReport as a Mermaid.js (`graph TD`) dependency graph for visualization in Markdown and HTML presentations.

## Rules
- L2 package (`02_shell/report`).
- Pure formatting of `domain.SearchReport`.
- Function: `GenerateMermaidGraph(report *domain.SearchReport) string`.
- Returns valid Mermaid `graph TD` syntax.
