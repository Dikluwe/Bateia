# Ordenação de chaves de dependências de manifestos

**Camada:** L1
**Proprietário:** `01_core/parser/manifest_keys.go`
**Contrato:** passo 0128 e ADR-0011.

Fornecer sortedManifestKeys genérico, retornando uma nova lista de chaves em ordem lexicográfica. Não mutar o mapa ou consultar efeitos externos. Consumidores JSON usam essa lista dentro de cada grupo existente, preservando precedência direct/dev/peer e regras próprias de nomes/versões. Não alterar a ordem de listas originalmente ordenadas pelo conteúdo.
