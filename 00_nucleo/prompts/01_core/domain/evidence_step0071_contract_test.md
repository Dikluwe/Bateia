# Contrato independente — domínio de cobertura do catálogo, passo 0071

Derivar testes apenas do passo
`00_nucleo/passos/0071-expor-cobertura-do-catalogo-sem-atribuir-lacunas-ao-candidato.md`
e do baseline público congelado. Não ler implementação candidata.

Verificar valores exatos e tags JSON de `EvidenceRuleBasis`, `CatalogCoverageStatus`,
`CatalogCoverageEntry` e `EvidenceCatalogCoverage`; arrays vazios não podem serializar
como `null`. Verificar round-trip de cobertura, `CapabilityRule.Basis`,
`CapabilityEvidence.ReasonCode` e os campos de cobertura em assessment/run.

O teste deve falhar pré-candidato e permanecer independente da solução escolhida. Não
aceitar aliases de constantes, omissão do reason code ou projeção de hipótese como regra
declarada.
