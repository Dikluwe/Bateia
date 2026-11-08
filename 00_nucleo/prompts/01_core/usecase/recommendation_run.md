# Durable recommendation run service

Implemente em L1 a orquestração pura do Passo 0026 conforme o contrato selado.
Exponha portas abstratas para store, reserva de fila, geração de ID e executor de
recomendação observável. Valide toda entrada antes de portas, normalize conjuntos,
compute hashes canônicos, aplique `CreateOrResolve` atômico e mantenha a máquina
de estados por CAS e versão monotônica.

Adicione observação controlada ao mesmo `RecommendationService`: `Recommend` e
`RecommendObserved` compartilham uma única implementação de pipeline/ranking e
produzem a mesma projeção certificada. O observer emite somente os oito estágios e
contadores congelados, sem texto externo. Não crie segundo pipeline.

Implemente start/get/artifact/cancel/execute/recover, projeções de artefatos e
provenance sem snippets. Preserve Unknown, `SUCCEEDED+INCONCLUSIVE`, corrida
terminal única e cancelamento que impede nova chamada após observação. Não faça
filesystem, rede, goroutine, mutex, log, env ou serialização MCP em L1.

## Passo 0129 — ownership de capacidade
O serviço durável conserva store, fila, IDs, execução e lifecycle único. Containers internos separados possuem registro de descoberta, configuração de extração e admissão de síntese. A extração e ordenação dos mecanismos pertence à capacidade; o runtime publica suas saídas sem reinterpretá-las. APIs Configure, hashes, artefatos e transições permanecem compatíveis.
