# Artefato opcional de rastreabilidade — 0098

Montar estágio trace apenas quando explicitamente solicitado, depois dos artefatos
causais. Hash de saída deve ser dos bytes do arquivo; graph_hash é selo canônico
interno distinto. Entradas: idea.json, reuse-decision.json e versões declaradas.
Não adicionar versão trace a RuleVersions/RunID ou critérios de retomada de
evidência. Contar nós/arestas e declarar limites, sem validar fatos ou veredito.
Propagar erros de leitura em vez de registrar hash vazio.
## Projeção 0099

Estágio trace usa contract_version decision-trace-projection-v2. Preservar schema
do grafo e todas as identidades causais. Outputs ordenados: trace-graph.json,
trace-index.json, cada qual hash dos bytes. Counts inclui nodes/edges/entries/
locations. O índice aponta a declarações do build input exato; não comprova fatos.
Falha de build/explain/escrita/hash registra estágio FAILED com outputs vazios,
sem anunciar projeção parcial como concluída. Escritas individuais são atômicas;
o par grafo/índice não é uma transação, e manifesto governa seu consumo.

0131: estágios e falhas de trace pertencem a recommendation_dossier_artifacts.go em L1. Este adaptador fornece hashes lógicos e mantém wrappers compatíveis para consumidores.
