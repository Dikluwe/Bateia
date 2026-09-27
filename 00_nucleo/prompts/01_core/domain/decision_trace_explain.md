# Índice consultivo de declarações — 0099

Definir contrato separado do grafo histórico decision-trace-v1, sem mudar seu
schema, seus bytes ou hash. Entrada schema1 liga grafo congelado e build_input
original; node_ids opcionais selecionam nós existentes, não fatos a verificar.

Saída decision-trace-explanation-v1 inclui graph_hash, report_fingerprint,
scope RECORDED_DEPENDENCIES_ONLY e entries com node_id/kind/label/fingerprint,
statement textual limitado, statement_truncated explícito e locations. Location
identifica documento report/idea/rule_versions e JSON Pointer RFC6901 relativo
ao documento; ponteiro vazio significa sua raiz. Nenhum path/URL a consultar.

Statements são projeções das declarações existentes, não novos juízos. Índice
exige documentos originais para leitura completa e não autentica fatos, autoria
ou validade das fontes. Limitations são obrigatórias na saída do núcleo.
