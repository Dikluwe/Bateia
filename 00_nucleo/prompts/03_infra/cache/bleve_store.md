# Especificação L0 — Store de Cache e Busca Local Bleve (03_infra/cache/bleve_store.go)

## Responsabilidade
Implementar a interface `domain.CacheStore` utilizando Bleve v2 (`github.com/blevesearch/bleve/v2`) para persistência embarcada, indexação full-text e consulta rápida de relatórios, manifestos e descobertas em disco.

## Restrições Tekt
- **Camada**: L3 (Infraestrutura).
- Conhece apenas L0 e L1 (`domain`).
- Proibido importar L2 ou L4.
