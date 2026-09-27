# 0102 — contexto documental das tarefas de revisão

Entrada schema_version 1 com review_input (0101) e build_input original (0098).
Saída decision-trace-review-context-v1: graph_hash, plan_hash, report_fingerprint,
scope RECORDED_DEPENDENCIES_ONLY, context_scope TASK_AND_CHANGED_ANCESTORS,
ledger completo de 0101, task_contexts, entries de 0099 e limitations.

Cada task_context tem node_id da tarefa e context_node_ids ordenados/deduplicados
da própria tarefa mais seus changed_ancestor_ids registrados. Entries é união
normalizada desses contextos: cada declaração e seus localizadores aparecem
uma vez, mesmo quando são usados por muitas tarefas. Ponteiros relativos aos
documentos lógicos report/idea/rule_versions; não são arquivos a abrir.

Não afirmar caminho causal completo, verdade/autoria, execução ou aprovação.
Zero tarefas implica task_contexts e entries vazios, sem dispensar validação.
Status, contagens, hashes, recibos e decisão do ledger preservados sem promoção.
