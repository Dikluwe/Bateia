# 0101 — Registrar o progresso consultivo das tarefas de revisão

## Resultado

idea trace review projeta recibos declarados num plano recomposto, sem promover
a decisão histórica. Não modifica recommend, a pesquisa simples, os grafos v1
ou os contratos de comparação/impacto anteriores. É offline: não abre referências,
executa candidatos, carrega .env ou consulta provedores.

Uso: `bateia idea trace review --input pedido.json --output revisao.json`.
Sem output, emitir JSON em stdout. O diretório de saída deve existir.
Entrada schema_version "1", impact_input completo no formato de 0098
(schema_version "1", graph congelado e changes), records opcional.

Para um plano de compare 0100, usar before como graph e seus
impact.changed_node_ids como changes com reason INPUT_CHANGED. A comparação
usa exatamente esse motivo consultivo. Não inferir REVOKED ou EXPIRED de ausência
de um nó. Para eventos reais explicitamente declarados, manter os motivos e
new_fingerprint do pedido original de impact.

## Fluxo

Primeira execução sem records: todas as tarefas registradas são PENDING.
A saída decision-trace-review-v1 apresenta plan_hash, graph_hash, tasks, counts,
status, scope, decision_before/after, decision_impact e limitations. Cada item
tasks tem task (a tarefa de impacto original), task_hash, status, note e evidence.

Para registrar andamento, adicionar records ao pedido original. Cada record
contém node_id (task.node_id), plan_hash, task_hash, status, note e evidence.
Usar os hashes da primeira saída, não preenchê-los por inferência.

Estados: PENDING, IN_PROGRESS, BLOCKED e COMPLETED. BLOCKED exige nota não vazia.
COMPLETED exige nota e ao menos uma referência evidence com ref não vazia e
content_hash SHA256 hexadecimal minúsculo, com 64 caracteres e sem prefixo.
Referências são identificadores lógicos declarados, não instruções de acesso.
Seu hash não é conferido contra arquivos, URLs ou testes executados.

Status geral: NO_REVIEW_TASKS se vazio; REVIEW_BLOCKED se há bloqueio;
RECORDED_REVIEW_COMPLETE se todas as tarefas têm conclusão declarada;
REVIEW_IN_PROGRESS se há alguma tarefa em andamento/concluída e outras abertas;
REVIEW_PENDING caso contrário. Zero tarefas não fabrica revisão concluída.
Contagens não medem risco, cobertura factual ou aprovação.

Decisão permanece histórica: decision_before igual a decision_after e
decision_impact NONE, inclusive em RECORDED_REVIEW_COMPLETE. É necessário produzir
e avaliar novas evidências/decisão pelo fluxo separado de Bateia. Conclusão
declarada não resolve automaticamente INCONCLUSIVE.

## Vinculação e limites

O núcleo recompõe o plano com PlanDecisionTraceImpact, em vez de aceitar tarefas
arbitrárias. PlanHash é SHA256 do JSON tipado de objeto interno {GraphHash,Changes,
Plan}: Changes copiado e ordenado por NodeID, Plan recomposto canonicamente.
TaskHash é SHA256 de {PlanHash,Task}. Esses objetos internos usam nomes de campos
Go sem tags, e incluem o pedido de alteração completo, não apenas o plano.
Assim, mudar new_fingerprint invalida recibos mesmo se o conjunto de tarefas
resultante continuar igual. Hashes vinculam representações, não autenticação.

Registros desconhecidos, duplicados, obsoletos ou com estado inválido falham
sem resultado parcial. Até 1000 registros e 20 referências por registro; limites
agregados de 0098 (2 MiB texto/JSON, 100000 valores, profundidade 64) antes da
serialização integral, tanto para entrada quanto para saída. Não truncar recibos
nem tarefas para caber. Notas e referências herdam limite de 8192 bytes UTF8.

Reutilizar leitura JSON estrita, proteção de entradas/aliases/hardlinks e escrita
atômica individual da CLI trace. Preservar saída existente em erro de validação.
Sem flag/novo comando, as saídas do produto continuam iguais.

Não é um log append-only: trata-se de snapshot de progresso, sem controle de
transições, identidade autenticada, relógio ou concorrência. Um novo pedido sem
recibos volta a PENDING; conservar os recibos fora da projeção é responsabilidade
do chamador. Não alegar auditoria temporal ou execução verificada.

## Verificação

Testes cobrem estados/contagens, hashes obsoletos, canonicidade, não mutação,
1000/1001 registros, 20/21 referências, entrada/saída expandida e falhas sem
resultado parcial. Integração cobre JSON estrito, aliases/FIFO, disco/stdout,
entrada intacta e despacho real offline ignorando .env.

Validação final: crystalline-lint --fail-on warning . sem violações,
go test -p 1 ./..., go vet ./..., git diff --check e testes de corrida
direcionados de domínio/núcleo/CLI/wiring passaram. Compilação de 04_wiring
concluída em diretório temporário. Gates, prazos e goldens legados preservados.
A checagem estrutural encontrou três pendências anteriores, corrigidas
separadamente no commit f9b7925: teste espelhado do domínio de comparação
e prompts exclusivos de dois testes, sem alterar seus oráculos.
