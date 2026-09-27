# 0101 — registro consultivo de progresso de revisão

Entrada schema_version 1 com impact_input (DecisionTraceImpactInput) e records
opcional. Cada registro tem node_id, plan_hash, task_hash, status, note e evidence.
Evidence contém referências lógicas ref e content_hash SHA256 hexadecimal minúsculo
de 64 caracteres, declarados pelo autor, sem acesso à fonte. Os estados permitidos
são PENDING, IN_PROGRESS, BLOCKED e COMPLETED; ausência significa PENDING.

Saída decision-trace-review-v1 contém plan_hash, graph_hash, scope, decision_before,
decision_after, decision_impact NONE, status, counts, tasks e limitations. Cada task
contém task (tarefa de impacto), task_hash, status, note e evidence. Counts distingue
pending, in_progress, blocked e completed. Status geral NO_REVIEW_TASKS,
REVIEW_PENDING, REVIEW_IN_PROGRESS, REVIEW_BLOCKED ou RECORDED_REVIEW_COMPLETE.
Conclusão é declaração, nunca aprovação, validade factual ou execução comprovada.
Sem I/O, transições automáticas, identidade autenticada, relógio ou histórico de eventos.
