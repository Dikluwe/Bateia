# ADR-0007: Container do runtime durável pertence a L3

**Status**: ACEITO
**Data**: 2026-09-19
**Relacionado**: ADR-0002, ADR-0004, ADR-0005, passos 0026–0028, 0077–0079

## Contexto

`04_wiring/durable_runs.go` ainda constrói `FileStore`, `LocalQueue`, o gerador de
IDs e `RecommendationRunService`, além de fechar diretamente store e fila. Essa
ligação forma um agregado de infraestrutura process-local e deixa L4 responsável
por detalhes concretos que não são decisões de composição de produto.

Testes históricos também verificavam por texto a localização concreta em L4. Essa
localização não é uma obrigação pública e conflita com a fronteira atual.

## Decisão

1. Criar `03_infra/runcontainer` como proprietário único do agregado concreto
   store, fila, gerador e serviço de runs.
2. Expor `New`, `Service`, `Recover`, `Start` e `Shutdown`; campos e componentes
   concretos permanecem privados.
3. Manter em L4 as quatro fachadas existentes, suas assinaturas, configuração de
   features, criação do servidor/CLI e a ordem explícita `Recover` antes de `Start`.
4. Preservar as ordens e precedências de cleanup.
5. Arquivar os testes históricos incompatíveis no Lab e cobrir o novo container com
   testes funcionais ativos.

## Consequências

L4 escolhe e injeta configuração, conectores, catálogo e transportes, mas deixa de
possuir store/fila. L3 passa a concentrar o lifecycle do agregado sem importar L2
ou L4. Tipos de limites, dispatch, defaults, ferramentas MCP, mensagens e formatos
ficam fora deste corte. Consumidores externos e variantes não executadas continuam
`Unknown`.
