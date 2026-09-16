# Prompt: 01_core/usecase/idea_radar

Implementar a função pura
`AnalyzeIdeaRadar(domain.IdeaRadarInput) domain.IdeaRadarReport` exatamente segundo o
contrato 0066 congelado. Canonicalizar somente IDs ASCII, validar projeções e origens,
deduplicar observações pela tupla contratual, agrupar por alvo e classificar pela
precedência pública. Exigir simultaneamente três kinds e três OriginIDs registrados.

Produzir draft explícito, sempre não pronto, para vinculação futura ao analisador 0065;
não fabricar `EvidenceSource` nem chamar coalizões com observações conceituais. Preservar
conflito, ambiguidade e Unknown. Ordenar toda saída, não mutar input e respeitar 64
projeções, 1024 observações, 512 alvos, 4096 operações e limites textuais. Usar índices e
passagens lineares; não enumerar combinações nem comparar observações aos pares.

Sem I/O, relógio, aleatoriedade, LLM, embeddings, dependências externas, ground truth ou
integração com busca/ranking/decisão. O futuro arquivo é
`01_core/usecase/idea_radar.go`, com cabeçalho de linhagem deste prompt e camada L1.
Autoridade do implementador após selo; não editar entradas protegidas nem testes.
