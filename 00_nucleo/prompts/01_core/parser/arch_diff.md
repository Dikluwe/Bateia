# Specification: Architectural Diff Evaluator (L1 Core)

## Purpose
Compute semantic architectural diff between two repository snapshots (added/removed dependencies, version upgrades, AST symbol delta).

## Rules
- L1 package (`01_core/parser`).
- ZERO I/O (pure evaluation in memory).
- Function: `ComputeArchDiff(oldRepo, newRepo *domain.ExtractedRepository) domain.ArchDiff`.
- Returns `domain.ArchDiff`.
