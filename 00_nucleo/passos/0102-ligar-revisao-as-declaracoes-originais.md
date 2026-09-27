# 0102 — Ligar tarefas de revisão às declarações originais

## Resultado

idea trace review-context combina o progresso consultivo de 0101 com os
localizadores históricos de 0099. Mostra onde a declaração da tarefa e seus
ancestrais alterados estão nos documentos originais. Não executa investigações,
autentica evidências ou aprova a decisão. Pesquisa simples e saídas anteriores
permanecem iguais; o novo comando é opt-in e offline.

Uso: `bateia idea trace review-context --input pedido.json --output contexto.json`.
Sem output, JSON em stdout. Diretório de saída já existente.
Entrada schema_version "1", review_input (pedido completo de 0101, incluindo
seu próprio schema_version "1") e build_input original (pedido de construção
0098, com schema_version "1", report e idea/rule_versions quando usados).

Não fornecer um índice, ledger ou plano pronto para ser confiado. O núcleo
recompõe o plano/progresso e o índice. build_input deve corresponder exatamente
ao grafo congelado de review_input.impact_input, inclusive política, ideia e
versões presentes na construção original. Incompatibilidade falha sem saída.

## Projeção normalizada

Saída decision-trace-review-context-v1 tem graph_hash, plan_hash,
report_fingerprint, scope RECORDED_DEPENDENCIES_ONLY,
context_scope TASK_AND_CHANGED_ANCESTORS, ledger, task_contexts, entries e limitations.
ledger é a saída inalterada de 0101, com tarefas, recibos, contagens e decisão
histórica. TaskHash/PlanHash permanecem iguais: acrescentar contexto documental
não conclui tarefas nem exige novos recibos.

Cada task_context liga node_id da tarefa a context_node_ids: ela própria e seus
changed_ancestor_ids registrados, ordenados e sem repetição. Não representar
isso como caminho causal completo ou como todas as dependências da tarefa.
Intermediários que não são a tarefa ou raízes alteradas não são inventados.

entries é a união desses IDs, cada um uma única vez. O consumidor resolve IDs
de cada tarefa nessa tabela. Entradas e localizadores de fontes compartilhadas
não são multiplicados por quantidade de tarefas. Cada entrada preserva o
fingerprint, statement, statement_truncated e locations de 0099.

Locations usam nomes lógicos report/idea/rule_versions e JSON Pointer RFC6901.
Ponteiro vazio é raiz do documento. Não são caminhos de arquivo ou instruções
para rede. Statements truncados continuam ligados ao registro integral, sem
inferir conclusão a partir do texto resumido. ReportFingerprint é da
representação tipada, não dos bytes arbitrários do arquivo (ver 0099).

Zero tarefas resulta em task_contexts e entries vazios, e não em índice completo
ou revisão concluída. O vínculo com os documentos ainda é validado. Como a API
explain usa seletor vazio para todos, a validação completa ocorre antes de
descartar suas entradas nesse caso; seus limites conservadores continuam valendo.

Os trechos são do snapshot histórico, não do candidato after. Remoção de fonte
não vira revogação factual. COMPLETED continua declaração de progresso e a
decisão histórica não é recalculada, substituída ou promovida pelo contexto.

## Segurança e limites

Entrada e saída agregadas seguem 0098: 2 MiB texto/JSON, 100000 valores,
profundidade 64 e demais limites de grafos/registros/localizadores. Uma requisição
com review_input e build_input pode exceder limites mesmo se cada parte couber
separadamente. A projeção de saída também passa por admissão agregada final;
não truncar tarefas, raízes, recibos ou localizadores para caber.

JSON estrito e writer trace compartilhados: hashes/documentos/recibos inválidos
não emitem stdout nem substituem saída existente. Recusar saídas que identificam
entradas por path, hardlink ou symlink; FIFO não é entrada válida. Segredos
configurados continuam bloqueados na emissão. Despacho não carrega .env/provedores.
Nenhuma leitura adicional de fontes, execução de candidatos ou I/O no núcleo.

## Verificação

Testes cobrem fontes compartilhadas, tarefas/raízes deduplicadas, ponteiros,
truncamento, ledger/decisão preservados, zero tarefas, canonicidade, não mutação,
documentos/recibos adulterados e orçamento agregado. A admissão final de saída
é testada diretamente com projeção tipada sintética expandida e vizinho positivo;
esse teste mede limites, não coerência documental (validada nos testes públicos).
CLI cobre JSON ambíguo, aliases/FIFO, segredo sintético configurado, disco/stdout
e despacho real sem .env.

Validação final: crystalline-lint --fail-on warning . sem violações,
go test -p 1 ./..., go vet ./..., git diff --check e testes de corrida
direcionados de domínio/núcleo/CLI/wiring passaram. 04_wiring compilou em
diretório temporário. Gates R5/R6, prazos, oráculos e goldens preservados.
