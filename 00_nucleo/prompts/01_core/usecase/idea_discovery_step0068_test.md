# Prompt proprietário — `01_core/usecase/idea_discovery_step0068_test.go`

## Obrigação do oráculo 0068

Verificar, a partir do contrato congelado, que a descoberta copia profundamente a
derivação da expansão para a requisição, não muta plano ou slices, executa todas as
expansões técnicas mesmo quando uma anterior retorna resultados e deduplica candidatos
por `owner/repo` conservando expansion IDs, queries e origens distintas em ordem
canônica. Cobrir repetição/reordenação e impedir que união apague provenance.

Este arquivo é origem L0 somente do teste nomeado. O contrato produtivo permanece em
`00_nucleo/prompts/01_core/usecase/idea_discovery.md`.
