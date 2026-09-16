# Contrato canônico 0065 — Coalizões de evidência em modo sombra

Autoridade normativa: `../../passos/0065-analisar-coalizoes-de-evidencia-em-modo-sombra.md`.
Este recorte congela os observáveis usados pelo gate antes de existir candidato.

## Fragmento observável

Entrada `evidence-coalition-input-v1` contém catálogo limitado e claims com
`BaselineStatus` e testemunhas explicitamente anotadas com papel, polaridade e
proveniência. A função pública pura é `usecase.AnalyzeEvidenceCoalitions`. Saída
`evidence-coalition-report-v1` contém assessment por claim e métricas; ordenação e JSON
são canônicos.

Resultados possíveis por claim, em precedência: `INVALID_INPUT`, `CONFLICT`,
`MISSING_ROLES`, `UNKNOWN_INDEPENDENCE`, `INSUFFICIENT_INDEPENDENCE`, `SATISFIED`.
`SUPPORT` cobre somente o próprio papel. Origem desconhecida não satisfaz independência;
hash/caminho/commit/ID distintos não a provam. Contraevidência prevalece e é preservada.
Duplicata física não multiplica evidência e divergência semântica entre duplicatas é
inválida.

O analisador ecoa `BaselineStatus`, emite `DecisionImpact=NONE` e
`Metrics.VerdictChanges=0`. Nenhum consumidor existente é alterado. Limites: 128 regras,
128 claims, 64 testemunhas por claim, 4096 total, quatro papéis por regra e 8192
operações. Excesso produz relatório `UNKNOWN/budget_exceeded`, sem resultado parcial.
Uma operação é uma testemunha normalizada ou um par claim/papel requerido inspecionado.
Power set é proibido.

Cada papel exigido precisa de ao menos um suporte `RECORDED`; origens conhecidas apenas
em outros papéis não removem `UNKNOWN_INDEPENDENCE`. Razões de assessment são exatamente
`coalition_satisfied`, `conflicting_evidence`, `missing_required_roles`,
`origin_unrecorded`, `distinct_origins_below_minimum` e `invalid_claim`. Razões de
relatório são `analyzed`, `analyzed_with_invalid_claims`, `invalid_schema`, `no_claims`,
`ambiguous_input` e `budget_exceeded`.

## Corpus mínimo a materializar pelo autor dos oráculos

| Caso | Construção | Resultado obrigatório |
|---|---|---|
| P01 | implementação A + teste B, ambas RECORDED, limiar 2 | SATISFIED / DISTINCT_ORIGIN |
| P02 | implementação A, limiar 1 | SATISFIED / RECORDED |
| N01 | implementação A + teste A, limiar 2 | INSUFFICIENT_INDEPENDENCE / SAME_ORIGIN |
| O01 | implementação A + teste UNKNOWN, limiar 2 | UNKNOWN_INDEPENDENCE / UNKNOWN_ORIGIN |
| N02 | apenas teste, exige implementação + teste | MISSING_ROLES |
| N03 | documentação + teste, exige implementação + teste | MISSING_ROLES |
| N04 | P01 mais uma testemunha CONFLICT | CONFLICT, suporte e conflito preservados |
| P03 | mesma fonte repetida semanticamente igual | uma testemunha canônica |
| N05 | mesma fonte repetida com polaridade/origem/papel divergente | INVALID_INPUT |
| O02 | claims vazios ou schema desconhecido | relatório UNKNOWN, nenhum assessment satisfeito |
| O03 | qualquer limite excedido | UNKNOWN / budget_exceeded, sem parcial |
| P04 | 20 permutações e nil/vazio | JSON byte a byte idêntico |

Todos os cinco valores de `EvidenceStatus` devem aparecer no corpus como baseline e ser
ecoado sem alteração. O corpus registra expectativa pública, mas não algoritmo candidato.

## Gate e medição

O adversário deve materializar M01–M12 enumeradas no passo. Selo somente com 12/12
rejeitadas (`mutation_score=1.0`), 100% do corpus correto, 20/20 permutações idênticas,
zero mudanças de veredito e contagem de operações exata. Regressão selecionada compara
snapshots de avaliações existentes antes/depois; o novo analisador não é chamado por
produção legada.

O contrato cobre somente classificação explícita em memória. Não mede qualidade da
inferência de papéis/proveniência, independência causal real, utilidade geral na busca ou
desempenho sobre repositórios reais.
