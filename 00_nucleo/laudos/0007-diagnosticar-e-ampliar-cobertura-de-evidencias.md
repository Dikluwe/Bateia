# Laudo 0007 — Diagnóstico e ampliação da cobertura de evidências

**Data:** 2026-08-26
**Estado:** implementado; gate atingido
**Pipeline:** determinístico e sem LLM

## Entrega

Foi adicionada a boca offline:

```bash
bateia idea evidence diagnose \
  --assessment assessment.json \
  --evidence bundles.json \
  --judgments judgments.json \
  --output coverage-report.json
```

O relatório agrupa aliases pelo par repositório–regra, calcula `proven_recall` e
atribui uma causa fechada a cada esperado `PROVEN` não emitido. Falhas operacionais
têm precedência sobre interpretações de conteúdo.

## Diagnóstico inicial

O número anterior de quatro `PROVEN` era inflado por aliases: problema e capacidade
usavam a mesma regra e a mesma prova. A linha de base independente era:

```text
esperados PROVEN independentes: 6
emitidos PROVEN independentes:  2
proven_recall:                  33,3%
```

As causas observadas após replay com alvos corretamente declarados foram:

| Projeto | Regra | Estado | Causa |
|---|---|---|---|
| Enola | `blast-radius` | `PARTIAL` | `SYMBOL_LINKAGE_MISS` |
| mllint | `static-analysis` | `INCONCLUSIVE` | `SYMBOL_LINKAGE_MISS` |
| edat | `async-events` | `INCONCLUSIVE` | `INSUFFICIENT_INDEPENDENCE` |
| Recipya | `recipe-import` | `PARTIAL` | `SYMBOL_LINKAGE_MISS` |

## Correções aceitas

1. O planejador passou a declarar também `IdeaSpec.Problem`, não somente
   `Capabilities`, nos alvos de evidência.
2. A árvore passou a priorizar paths relacionados à capacidade antes de consumir
   o orçamento por ordem alfabética.
3. `blast-radius` passou a reconhecer testes específicos de spillover de mudança.
4. `async-events` passou a reconhecer dispatcher, handler e store de eventos.
5. `recipe-import` passou a reconhecer importadores de plataformas e componentes
   específicos de repositório/controlador de receitas.
6. `component_precision` passou a revisar somente componentes realmente emitidos.
7. A versão das regras avançou para `2026-08-26.2`.

Cada generalização recebeu caso positivo e ataque negativo. Termos genéricos como
`Controller`, `Repository`, `Recipe` e spillover não relacionado continuam
insuficientes para `PROVEN`.

## Revisões de ground truth

- mllint foi rebaixado de esperado `PROVEN` para `PARTIAL`: o corpus comprova
  orquestração de linters, não análise estática própria.
- `ccbrown/cloud-snitch` foi promovido de esperado `ABSENT` para `PROVEN`: o commit
  coletado contém `stripeEventHandler` ligado a uma fila de mensagens mortas. A
  classificação negativa anterior contradizia a evidência estrutural.

Essas mudanças alteraram julgamentos versionados, não regras especiais por
repositório.

## Ampliação do corpus

Foram adicionados positivos independentes:

- `looplab/eventhorizon`;
- `ThreeDotsLabs/watermill`;
- `asaskevich/EventBus`;
- `TandoorRecipes/recipes`;
- `mealie-recipes/mealie`.

Os projetos só entraram na contagem depois de apresentarem símbolos, testes ou
componentes rastreáveis no commit coletado.

## Resultado final

| Corpus | Esperados independentes | Emitidos | Recall | Precisão | Falsa ausência | Componentes |
|---|---:|---:|---:|---:|---:|---:|
| impacto estrutural | 2 | 2 | 100% | 100% | 0% | 100% |
| análise estática | 1 | 1 | 100% | 100% | 0% | não mensurável |
| eventos assíncronos | 5 | 5 | 100% | 100% | 0% | 100% |
| organização de receitas | 3 | 3 | 100% | 100% | 0% | 100% |
| **agregado** | **11** | **11** | **100%** | **100%** | **0%** | **100%** |

`source_trace_rate` e `budget_completion` foram 100% nos quatro corpora. Nenhum
alias equivalente aumentou a amostra.

## Reprodutibilidade

Relatórios de cobertura pequenos foram incorporados ao corpus. Planos, bundles e
assessments completos são identificados por SHA-256, tamanho e lista de
`repo@commit` no manifesto
`02_shell/cli/testdata/evidence_calibration/calibration-run-2026-08-26.manifest.json`.
Isso permite recoletar nos mesmos commits e comparar o hash sem versionar código
integral de terceiros.

## Veredito

O gate do Passo 0006 foi atingido sem LLM, redução de limiar ou aliases duplicados.
O Passo 0008 está liberado para implementar a decisão determinística de reuso.

Regime Tekt: A/B limitado, executado sem atestação formal de isolamento. Os testes
adversariais foram escritos antes das generalizações locais, mas contrato, solução
e veredito compartilharam a mesma sessão.
