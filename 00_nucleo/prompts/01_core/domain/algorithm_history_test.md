# Prompt proprietário — `01_core/domain/algorithm_history_test.go`

Teste somente o contrato público de tipos de `01_core/domain/algorithm_history.go` do
passo 0070, sem ler nem usar detalhes da implementação candidata como oráculo.

Materialize testes table-driven que comprovem a existência e o valor exato dos tokens
públicos:

- cronologias `KNOWN`, `APPROXIMATE` e `UNKNOWN`;
- kinds de proveniência `PRIMARY_SOURCE`, `SECONDARY_SOURCE` e
  `EMPIRICAL_ARTIFACT`;
- relações `DERIVED_FROM`, `INFLUENCED_BY`, `IMPROVES`, `GENERALIZES`, `SPECIALIZES`,
  `COMBINES`, `REDISCOVERS` e `RESEMBLES`;
- kinds de base `EVIDENCE` e `INFERENCE`;
- stances `SUPPORT`, `CONFLICT` e `AMBIGUOUS`;
- outcomes `CONTESTED`, `UNKNOWN`, `CONTRADICTED`, `CORROBORATED` e `SUPPORTED`;
- schemas `algorithm-history-input-v1` e `algorithm-history-report-v1`, e statuses
  `COMPLETE` e `UNKNOWN`.

Valide por `encoding/json`, usando apenas tipos e constantes exportados, o shape público
snake_case de input, algorithm, provenance, claim, basis, report, assessment e metrics.
Exija que `year_start` e `year_end` representem número ou `null`, que enums sejam strings
e que nenhum campo público contratual desapareça ou seja renomeado silenciosamente.

Construa um `AlgorithmHistoryReport` público mínimo e comprove que `assessments` e todas
as listas do assessment — evidências e inferências separadas por stance, além de origens
de suporte — serializam como `[]`, nunca `null`. Faça o mesmo para listas públicas de
entrada quando inicializadas vazias. Não aceite inferência misturada nas listas de
evidência nem campos JSON privados/não contratuais.

Os testes devem ser puros, determinísticos e independentes: zero filesystem, rede, LLM,
relógio, aleatoriedade, reflexão sobre fonte, `unsafe` ou import de L2/L3. Não teste
algoritmos de análise neste arquivo; essa responsabilidade pertence ao oráculo segregado
de `01_core/usecase/algorithm_history_test.go`.
