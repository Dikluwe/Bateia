# Resultado 0075 — Expressões probatórias compostas

A capacidade de guard de autoatribuição passou de `INCONCLUSIVE` para `PARTIAL` nos dois
candidatos do estudo SbHash. A causa observada não era perda da busca dirigida: os
arquivos já constavam como candidatos examinados. O localizador intrarquivo é que
rejeitava padrões de `symbol` formados por uma expressão com operadores.

## Correção

`supportedSymbolOccurrence` agora admite uma expressão composta quando:

- o padrão declarado contém ao menos dois identificadores e operador de código;
- o casamento literal já ocorreu numa única linha limpa;
- comentários e literais já foram removidos pelo lexer conservador.

Expressões multiline, texto natural sem operador, comentários, strings e substrings de
identificador continuam sem autoridade. Os padrões fortes não foram adicionados às
consultas dirigidas: `EvidenceRequirementQuery.Terms` continua contendo somente
`hashtable.h` e `robin_hood.h`, derivados de `TargetHints`.

## Replay observado

O replay reutilizou os bundles imutáveis de
`lab/bancada/history/reports/0075-coin-sbhash/dossier` e o catálogo
`lab/bancada/history/reports/0075-coin-sbhash/rules.json`:

| Candidato | Antes | Depois | Fonte exata |
|---|---|---|---|
| libstdc++ 13 | INCONCLUSIVE | PARTIAL (0.60) | `bits/hashtable.h:1257`, `if (&__ht == this)` |
| robin-hood-hashing | INCONCLUSIVE | PARTIAL (0.60) | `src/include/robin_hood.h:1614`, `if (&o == this) {` |

O assessment reproduzido teve SHA-256
`f93c01660c1bc25974265a30ffddcd81120c4a452ec6fbee37bf7d9df8097663`. O estado geral
do dossiê permaneceu `INCONCLUSIVE`, pois outras capacidades continuam não verificadas.

## Verificação e limites

- gate focal passou duas vezes;
- `go test ./01_core/usecase` passou;
- `go test -race ./01_core/usecase` passou;
- `go test ./...` passou;
- build de `./04_wiring`, `crystalline-lint .` e `git diff --check` passaram.

O teste pré-candidato testemunhou a falha dos quatro positivos antes da alteração, mas
contrato, testes, implementação e veredito foram integrados pela mesma autoridade no
mesmo filesystem. Portanto, o regime Tekt foi executado sem atestação de isolamento e
sem mutation score independente. O resultado atesta apenas localização lexical; não
prova que o guard seja a primeira instrução efetiva nem que nenhuma mutação destrutiva o
preceda após expansão de macros ou análise semântica.
