# 0101 — projeção pura do andamento

TrackDecisionTraceReview aceita entrada tipada, aplica limites agregados antes da
serialização e recompõe PlanDecisionTraceImpact em vez de confiar num plano fornecido.
Máximo 1000 registros, 20 referências por registro, 2 MiB texto/JSON e demais
guardas 0098. Hash do plano: SHA256 de JSON tipado de {GraphHash,Changes,Plan},
com Changes copiado/ordenado por NodeID e Plan canônico recomposto. Inclui motivos
e new_fingerprint, inclusive eventos sem impacto, para invalidar recibos obsoletos.
TaskHash: SHA256 do JSON tipado de {PlanHash,Task}. Nomes dos campos são os campos
Go sem tags desses dois objetos internos, fazendo parte do algoritmo v1.

Recusar registros desconhecidos, repetidos, hashes não correspondentes e estado
inválido. BLOCKED requer nota não vazia; COMPLETED requer nota e ao menos uma
referência. Referências fornecidas em qualquer estado requerem ref não vazia e
content_hash hexadecimal minúsculo SHA256. Não buscar, abrir ou verificar bytes.

Gerar tarefas na ordem canônica do plano. Recibos ausentes => PENDING. Status geral
prioriza BLOCKED, depois all COMPLETED, depois algum IN_PROGRESS/COMPLETED, depois
PENDING; zero tarefas => NO_REVIEW_TASKS, sem alegar revisão concluída. Preservar
decisão histórica sempre. Copiar fatias de referências; não alterar entrada.
Validar a saída agregada antes de retornar; erro => resultado zero, sem truncamento.
Reordenar grafos/eventos/registros não deve alterar hash ou resultado.
