# Specification: Architectural Violation Linter (L1 Core)

## Purpose
Analyze project structure and dependencies to detect architectural anti-patterns and layer leakage.

## Rules
- Pure L1 package (`01_core/parser`).
- ZERO I/O (no file reads, no network requests).
- Evaluates `*domain.Manifest` and `[]domain.DirectoryNode`.
- Returns `[]domain.ArchViolation`.
