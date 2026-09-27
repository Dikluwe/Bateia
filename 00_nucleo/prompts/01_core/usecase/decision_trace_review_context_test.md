# 0102 — testes puros do contexto de revisão

Conferir joins próprios+ancestrais, shared source com várias ocorrências,
deduplicação/canonicidade e não mutação. Ledger/plan_hash/task_hash idênticos ao
0101; entries iguais à seleção de 0099, inclusive truncamento e fingerprints.
COMPLETED não promove decisão; zero tarefas não projeta índice inteiro.

Recusar mismatch de relatório/ideia/regras, schema, hash adulterado, recibos
obsoletos e UTF8 inválido, inclusive sem tarefas. Quotas agregadas antes de
serialização e expansão de join em output devem falhar com resultado zero;
vizinho positivo continua utilizável. Sem I/O/rede/relógio.
