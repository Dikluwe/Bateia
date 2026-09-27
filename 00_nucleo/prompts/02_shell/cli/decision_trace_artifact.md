# Artefato opcional de rastreabilidade — 0098

Montar estágio trace apenas quando explicitamente solicitado, depois dos artefatos
causais. Hash de saída deve ser dos bytes do arquivo; graph_hash é selo canônico
interno distinto. Entradas: idea.json, reuse-decision.json e versões declaradas.
Não adicionar versão trace a RuleVersions/RunID ou critérios de retomada de
evidência. Contar nós/arestas e declarar limites, sem validar fatos ou veredito.
Propagar erros de leitura em vez de registrar hash vazio.
