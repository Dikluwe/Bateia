# Revisão sistemática viva do funil — passo 0084

## Resultado

O Bateia ganhou uma superfície avançada, offline e aditiva para tratar uma execução
como revisão auditável. Ela não substitui `search` nem o fluxo simples de pesquisa.

| Passo | Artefato implementado | Propriedade principal |
|---|---|---|
| 1 | `ReviewProtocol` | regras congeladas e hash verificável antes da seleção |
| 2 | `CandidateDispositionLedger` | todo registro e candidato conserva seu destino |
| 3 | `FunnelFlowReport` | contagens derivadas; identificação e terminais fecham |
| 4 | `CandidateDispositionReason` | vocabulário estável e obrigatório para o destino final |
| 5 | `ReviewCertaintyReport` | certeza por claim separada de `EvidenceStatus` e decisão |
| 6 | `LivingReviewIndex` | linhagens por protocolo e deltas append-only |
| 7 | `ActiveLearningPlan` | prioriza leitura, mas nunca exclui nem encerra a revisão |

## Superfície operacional

```text
bateia idea review protocol   --input protocol-draft.json
bateia idea review synthesize --input synthesis.json
bateia idea review update     --input living-update.json
bateia idea review prioritize --input active-learning-request.json
bateia idea review evaluate   --input corpus.json
```

Todos os comandos usam JSON estrito e rejeitam campos desconhecidos. O protocolo e o
pacote são recusados se seu conteúdo não corresponder ao hash declarado.

## Avaliação independente do runtime

O corpus versionado em `lab/bancada/corpora/step0084-active-learning.json` fica fora do
núcleo e contém dois casos, doze candidatos e seis relevantes. Os rótulos completos são
usados apenas pelo avaliador: a cada iteração, o priorizador recebe somente julgamentos
já revelados.

Resultado observado em 2026-09-24:

- recall em metade do orçamento: `0,6667` nos dois casos;
- itens lidos para atingir ao menos 95% de recall: `4/6` nos dois casos;
- trabalho poupado nesse ponto: `33,33%` nos dois casos;
- autoridade para stopping: `false`.

Esses números validam o mecanismo e o isolamento do gold; não demonstram validade
externa e não autorizam exclusão automática. Um corpus maior, rotulado por revisores
independentes e representativo dos domínios reais continua necessário antes de usar as
métricas como expectativa de produto.

## Limites deliberados

- o relatório de certeza usa vocabulário próprio; não alega conformidade GRADE;
- o fluxo é inspirado por práticas de revisões sistemáticas, sem alegar conformidade
  PRISMA;
- `provider_rejected_observations` é telemetria do último attempt de cada expansão e
  não participa da equação de conservação dos registros retornados;
- active learning é uma fila adaptativa de revisão humana, não uma decisão automática.
