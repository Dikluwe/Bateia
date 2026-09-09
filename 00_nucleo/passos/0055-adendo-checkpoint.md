# Adendo 0055 — Identidade de retomada

Inspeção de integração encontrou que EvidenceVersion invalida o cache do funil, mas a CLI idea recommend --resume compara uma lista independente de RuleVersions. Para modo live, incluir `discovery_query: 0055-v1` nessa identidade e comparar discovery_query na elegibilidade do checkpoint. Checkpoint live antigo sem essa versão não deve ser reutilizado. Modos local/replay continuam com campo ausente. Não reinterpretar replays históricos. Este adendo é anterior à implementação desse trecho e não altera catálogo, perguntas ou expectativas de recall.
