# 0117 — Modelo executável de checagem seguida de escrita

Lab, modelo autoral, não importar/executar candidato. Model representa somente
sentinela initialized e valor configurado limit, ambos atômicos. Configure(n,
afterCheck) rejeita por panic quando initialized=true; caso contrário chama hook
opcional de bancada após Load e antes de Store de limit. Initialize faz Store
true. Snapshot retorna initialized/limit. Sem garantir snapshot atomicamente
conjunto em presença de concorrência; lê-lo somente após rendezvous/joins.

Abstractização: substituir canal/sync.Once/worker por marcador; atomicidade por
operação remove data races do modelo, mas não torna check+write uma transação.
Não tratar resultado como comportamento executado de conc, performance, defeito
ou precisão independente. Hook existe só para agendar intercalação explícita.
