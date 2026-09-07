# Specification: Architectural Violation Linter (L1 Core)

## Purpose
Analyze project structure and dependencies to detect architectural anti-patterns and layer leakage.

## Rules
- Pure L1 package (`01_core/parser`).
- ZERO I/O (no file reads, no network requests).
- Evaluates `*domain.Manifest` and `[]domain.DirectoryNode`.
- Returns `[]domain.ArchViolation`.

## Revisão 2026-09-07 — Evidência proporcional

Consulta original é a autoridade de relevância; identidade candidato/módulo não prova pertinência. README exige origem explícita e termos literais. Extração depende da consulta, suporta caminhos C/C++/shaders, penaliza temporários e declara cobertura insuficiente. Registrar tentativas e relaxamentos, hashes dos bytes e intervalos examinados, parâmetros, versão e commit-base local. Pontuação heurística não demonstra solução nem saúde arquitetural. Separar fato observado, inferência, hipótese e não verificado; telemetria é do Bateia. Não inventar padrões com zero candidatos nem convenções Go para outras linguagens. Classificar dependências por identidade. Testar OIT/USalign/BinDiff, Coin, cobertura ausente e mudanças rastreadas/não rastreadas do worktree.
