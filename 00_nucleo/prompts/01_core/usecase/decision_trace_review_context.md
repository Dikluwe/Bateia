# 0102 — join puro de progresso e declarações históricas

BuildDecisionTraceReviewContext valida schema e orçamento agregado antes de
serializar entrada completa. Reutiliza TrackDecisionTraceReview e
ExplainDecisionTrace, ligando build_input ao grafo congelado de review_input.
Recusar mismatch mesmo sem tarefas; não aceitar índices/planos externos fornecidos.

Selecionar a união ordenada das próprias tarefas e changed_ancestor_ids. Cada
contexto de tarefa contém IDs ordenados únicos. Verificar cobertura de todas as
declarações selecionadas; saída sem duplicação de entries/localizadores.
Para zero tarefas, a API explain valida o vínculo documental completo, mas a
projeção descarta todas as entradas: seletor vazio de explain significa todos,
não projetar todos acidentalmente. Limites conservadores de explain continuam.

Preservar ledger de 0101, inclusive plan_hash/task_hash e COMPLETED declarativo.
Entradas só históricas, não ler after nem interpretar fonte ausente como revogada.
Statements/truncamentos e fingerprints de 0099 permanecem inalterados. Não
modificar entrada, não fazer I/O, não executar decisões/fontes/testes/candidatos.

Limites agregados de 0098 (2 MiB texto/JSON, 100000 valores, profundidade 64)
para entrada/saída. Guardar saída antes da serialização integral. Expansão do
join também pode exceder quotas, mesmo quando ledger e índice cabem separados.
Erro => resultado zero; nunca truncar tarefas, ancestrais ou localizadores.

A admissão final boundedReviewContextProjection é uma função privada pura de
orçamento de saída, não validador de coerência semântica. Testá-la diretamente
com projeção tipada expandida: ledger/índice isolados cabem, junção excede; retorno
zero em erro. O facade público sempre passa seu resultado por esta admissão.
