# Contrato de oráculo black-box independente — Passo 0028 O16

**Revisão:** O16 vinculada a R4, catálogo v3 e protocolo O8  
**Regime:** protocolo Tekt completo  
**Papel:** autoridade segregada de ORÁCULOS  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

As APIs permanecem exatamente `Evaluate(contractInput, contractOutput []byte)`
e `EvaluateSequence(exchanges []Exchange)`. A decisão continua binding-free e
usa somente os dois lados públicos do contrato; IDs de caso, hashes de ataque,
records/configs, expected, markers, texto livre, candidato e artefatos reservados
não controlam predicados.

## Entradas públicas congeladas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap público O15 | `6346e52b7c46bfb452eac348d810b178bda51a5ca065117cb6328e35747771e5` |
| Falha discriminatória pública O15 | `3b8ca43ff64bcb85ca1382cc1ce86a7fcc4c3d3334d61b2d0a19189b85d2d9fc` |
| Manifest protegido predecessor O15 | `ba502012af763a47420ca1d7b80b493d13620aea5cace203eac6ec4602b8f25b` |

## Refinamento normativo O16

Somente quatro limites são alterados. Em `MECHANISM_EXTRACT`, witness técnico
com assessment `PROVEN` que resulta em agregado `PARTIAL` é uma conclusão
decidível e permanece `Preserved`; insuficiência semântica requer assessment
técnico `PARTIAL` ou gap controlado coerente.

Em `ARTIFACT_READ`, target `UNKNOWN` com plano não prova baseline ausente, e a
lista obrigatória `integration_unknowns` não prova por si composição opaca.
Esses estados conformes permanecem `Preserved`. Licença `UNRECOGNIZED` numa
hipótese segue as relações de transfer/status de R4; já o warning controlado
`LICENSE_NOT_RECOGNIZED` dentro do artifact tipado `hypothesis_gaps` materializa
o estado opaco congelado e resulta em `Unknown`.

Toda violação decidível precede `Unknown`. Os 160 predicados negativos e os
demais nove controles positivos e nove opacos preservam a semântica O15. O16
não emite receipt, seal, certificado nem veredito.
