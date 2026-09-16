# Prompt proprietário — `01_core/domain/idea_discovery_step0068_test.go`

## Obrigação do oráculo 0068

Verificar exclusivamente a extensão aditiva de `DiscoveryRequest` congelada no passo
0068: o campo opcional `Derivation *ExpansionDerivation` deve possuir tag JSON
`derivation,omitempty`, serializar provenance quando presente e permanecer ausente em
requests legados. O teste não autoriza alterar outros campos, schemas persistidos ou
converter derivação ausente em sucesso técnico.

Este arquivo é origem L0 somente do teste nomeado. O contrato produtivo permanece em
`00_nucleo/prompts/01_core/domain/idea_discovery.md`.
