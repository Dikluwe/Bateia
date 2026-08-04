# Specification: Polyglot Communication Matrix (L1 Core)

## Purpose
Analyze cross-language project dependencies to detect shared communication protocols, schemas, and transport contracts (gRPC, Protobuf, GraphQL, Redis, Kafka, RabbitMQ).

## Rules
- L1 package (`01_core/parser`).
- ZERO I/O (pure evaluation in memory).
- Function: `EvaluatePolyglotMatrix(results []domain.ExtractedRepository) []domain.PolyglotContract`.
- Returns list of `domain.PolyglotContract`.
