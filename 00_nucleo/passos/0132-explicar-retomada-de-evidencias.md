# 0132 — Explicar a retomada de evidências

**Data:** 2026-10-28. **Prioridade do usuário:** confiabilidade e clareza dos resultados.

## Contrato anterior à materialização

Hoje `idea recommend --resume` pode recusar um checkpoint e executar a jornada novamente sem explicar a recusa. Acrescentar diagnóstico da elegibilidade, mantendo a política dos passos 0130/0131. L1 possui avaliação e motivo; L2 somente apresenta o resultado. L3 e L4 conservam suas responsabilidades.

`EvaluateEvidenceCheckpoint` retorna `EvidenceCheckpointAssessment` com `Reusable`, código tipado `Reason` e `Subject` opcional limitado a nomes lógicos constantes. `CanReuseEvidenceCheckpoint` delega à mesma avaliação. Retornar o primeiro impedimento na ordem atual: status, schema, hash da evidência, catálogo, presença/hash do modelo opt-in, versões de regras na ordem existente, leitor e artefatos na ordem registrada. Distinguir artefato ilegível, hash divergente e registro obrigatório ausente. Na ausência de ambos os registros, informar `discovery-snapshot.json` primeiro, sem leitura extra. Preservar short-circuit, primeira ocorrência com hash não vazio e tratamento das duplicatas. Não comparar policy/target para invalidar evidência.

A jornada oferece `Checkpoint` opcional no resultado em memória, somente quando tentou a retomada. Erro na leitura/interpretação do manifesto produz `MANIFEST_UNAVAILABLE`, sem copiar o erro ou supor que o arquivo esteja ausente. O diagnóstico fica disponível também se uma etapa posterior falhar, quando já tiver sido calculado. A decisão de reutilizar vem dessa mesma avaliação, sem leitura adicional. Sucesso operacional e elegibilidade do checkpoint não promovem a força probatória de nenhuma conclusão.

Ao anunciar um dossiê concluído com `--resume`, a CLI acrescenta uma linha em português indicando reutilização ou recusa, código e campo/artefato lógico pertinente. Explicar a recusa sem afirmar que uma nova coleta foi necessariamente executada: replay pode continuar usando sua fonte congelada. Sem `--resume`, manter a saída atual. Falhas continuam sem anúncio de conclusão. Preservar callback antes do cleanup, ordem dos efeitos, erros, schema dos artefatos, preimages, RunID, versões e política de evidência. Nenhum sidecar novo, timestamp ou input causal.

Não incluir no diagnóstico valores de entradas, versões recebidas, hashes, caminhos físicos ou texto de erros do leitor. O renderizador aceita somente os campos lógicos conhecidos ao mostrar `Subject`. Motivo desconhecido recebe mensagem genérica. Não acrescentar consultas remotas, filtros ou relaxamento de critérios.

## Aceite

Testes puros cobrem todos os motivos, precedência, ordem/short-circuit, duplicatas, determinismo, imutabilidade e equivalência do predicado. Testes da jornada cobrem ausência de tentativa, sucesso, recusa, manifesto indisponível e falha posterior sem promoção; testes CLI verificam mensagens, silêncio sem flag, proteção de valores e preservação de causalidade/artefatos. Executar `make verify`. Baseline e resultados ficam em `lab/bancada/confiabilidade/0132`; o laudo registra cobertura e limites.

## Registro da execução

Consultar o [laudo 0132](../laudos/0132-diagnostico-de-retomada.md) e a [bancada](../../lab/bancada/confiabilidade/0132/README.md). O diagnóstico é uma melhoria de clareza com preservação da política, sem promoção dos resultados históricos. O laudo distingue a melhoria entregue da oscilação legada de serialização observada no ensaio.
