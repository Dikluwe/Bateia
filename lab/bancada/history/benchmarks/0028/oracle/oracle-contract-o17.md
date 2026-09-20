# Contrato de oráculo black-box independente — Passo 0028 O17

**Revisão:** O17, R4/catálogo v3/protocolo O8  
**Regime:** Tekt completo  
**Autoridade:** ORÁCULOS segregada  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

`Evaluate(contractInput, contractOutput []byte)` e `EvaluateSequence` permanecem
binding-free. Identidades de caso, hashes externos, records/configs, expected,
markers, texto humano, candidato e artefatos reservados não participam da
decisão.

## Entradas congeladas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap O16 | `93f18dfda98a96337394861a95625b6c5d59c2404c2da9b8fabc4cfdbfcf463b` |
| Falha discriminatória O16 | `fe5460b69349fca097aa9cf5b0705efe891dde7a214b8403fd1968e1438d9b50` |
| Manifest predecessor O16 | `f038cea1a4be217a0452c85bcc0abe808982388a4ef3e2995263c419d5772da7` |

## Refinamento O17

`R4.ARTIFACT_READ.REAL_OPACITY_CLASSIFICATION` é uma tabela estrutural total.
Em hipótese completa, target `UNKNOWN`, composição, licença e counterevidence
só são opacos quando os enums e relações de incapacidade de R4 estão presentes.
Em `mechanism_gaps|hypothesis_gaps`, somente `code` controlado reconhece
baseline, compatibility, witness, identidade/binding imutável,
provider/capture, utility/security, budget e licença não resolvidos. As formas
controladas missing/absent/unknown/unavailable/unmeasured/insufficient,
ambiguous/unbound/unresolved, without e `NOT_*` equivalentes são aceitas.

O01/O02 permanecem `Unknown`; os oito limites públicos O03–O10 passam a
`Unknown`; P01–P12 permanecem `Preserved`; as 160 violações e sua precedência
são preservadas. O17 não emite receipt, seal, certificado nem veredito.
