# Prompt espelhado — tipos públicos do harness do Passo 0029

Materialize os tipos de dados públicos e determinísticos usados pelo pacote
`step0029_experiment`, conforme o contrato R1 em
`00_nucleo/benchmark/0029/contract.md` e a obrigação operacional de
`00_nucleo/prompts/04_wiring/step0029_experiment/harness.md`.

Este arquivo é a semente proprietária de `types.go`. Ele separa a declaração da
superfície de valores da implementação das operações do harness, evitando que
dois consumers Go reivindiquem o mesmo prompt cristalino.

## Fronteira

- camada `L4`, composição externa de benchmark;
- nenhuma função de I/O, dispatcher, rede, filesystem, processo ou chamada de
  agente/tool;
- declarar constantes, enums, observações tri-state, perfis, casos, células,
  tentativas, usage, receipts, anonimização, julgamentos, scores, análise e
  entradas/saídas de gates;
- manter aliases JSON públicos determinísticos sem completar dados ausentes nem
  converter `Unknown` em sucesso;
- manter conteúdo privado fora de receipts e projeções;
- relógio, hashes, seed, IDs, preços e disponibilidade são sempre valores de
  entrada.

## Compatibilidade observável

Os tipos devem sustentar exatamente as nove funções públicas declaradas pelo
contrato: `BuildPlan`, `ValidatePlan`, `ClassifyAttempt`, `EvaluateRetry`,
`BuildRunReceipt`, `Anonymize`, `Score`, `Analyze` e `EvaluateGates`.

Decoders de compatibilidade podem aceitar nomes públicos `snake_case` e nomes
exportados Go, desde que campos canônicos tenham precedência, colisões sejam
observáveis e a representação interna preserve as políticas congeladas. A
serialização pública pode projetar números conhecidos de score e mediana junto
da observação tri-state correspondente.

## Segurança causal

Não incorporar casos protegidos, expectativas, mutantes, oráculos, resultados,
receipts de execução, selos, certificados ou vereditos do Passo 0029. Este
prompt define somente o vocabulário público consumido pela implementação e por
autoridades independentes.
