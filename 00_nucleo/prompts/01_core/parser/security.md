# Specification: Security Advisory Scanner (L1 Core)

## Purpose
Analyze project manifest dependencies against known security advisories and vulnerable package versions.

## Rules
- Pure L1 package (`01_core/parser`).
- ZERO I/O (no file reads, no network requests).
- Evaluates `*domain.Manifest`.
- Returns `[]domain.SecurityAdvisory`.
