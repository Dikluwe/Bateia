# Specification: Architecture Decision Record Generator (L1 Core)

## Purpose
Synthesize domain SearchReport into a formal Architecture Decision Record (ADR) document (Context, Decision, Status, Consequences).

## Rules
- L1 package (`01_core/parser`).
- ZERO I/O (pure transformation in memory).
- Function: `GenerateADRDocument(report *domain.SearchReport) string`.
- Returns Markdown ADR compliant with Michael Nygard / C4 standards.
