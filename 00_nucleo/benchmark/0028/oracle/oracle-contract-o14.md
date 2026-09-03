# Contrato de oráculo black-box independente — Passo 0028 O14

**Revisão:** O14 vinculada a R4, catálogo v3 e protocolo O8
**Regime:** protocolo Tekt completo
**Papel:** autoridade segregada de ORÁCULOS
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

As APIs permanecem exatamente:

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

O avaliador decide somente por `contract_input` e `contract_output`. IDs de
caso, hashes de ataque/binding, records/configs, expected, markers, candidato e
artefatos de veredito não entram na API nem nos predicados.

## Entradas públicas congeladas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap público O13 | `5cd383b5e1ab05d2c27859348f1be7414b0188bfda195adcb615d8f9b395b369` |
| Falha discriminatória pública O13 | `895f3a547df72fd5f04695a7e1011092f4940e82345572ea37704003305b570b` |
| Manifest protegido predecessor O13 | `06d8c7a88910b4ef2e99e9db83b5a63384546d7fbdf6d83b5c21506bfc76afb1` |

## Refinamentos normativos O14

No predicado `R4.MECHANISM_EXTRACT.CANONICAL_CAUSAL_MECHANISM_ID`, o campo
`rule_id` do preimage é sempre o `id` do registro `MechanismRule` resolvido pela
classe e source no catálogo v3. `ExtractionRuleIDs` e `MechanismWitness.rule_id`
continuam preservando os RuleIDs das EvidenceSources; eles não substituem a
identidade da regra de mecanismo. Ordem fixa, NFC local, HTML escaping
desabilitado e o conjunto completo de witness hashes de O13 permanecem ativos.

`R4.MECHANISM_EXTRACT.REAL_OPACITY_CLASSIFICATION` reconhece witness mínimo
semanticamente insuficiente quando um gap request-derived coerente preserva o
teto `PARTIAL|INCONCLUSIVE`. O gap não pode mascarar falha de provenance,
integridade, status, ID ou safety.

`R4.ARTIFACT_READ.REAL_OPACITY_CLASSIFICATION` reconhece relacionalmente os
estados literais do protocolo O8:

- target `UNKNOWN` junto de baseline ausente;
- composição de mecanismos não medida junto de status/transfer incapaz e plano
  ou contexto epistemológico;
- licença `UNRECOGNIZED` junto de status `INCONCLUSIVE` ou transfer
  `UNKNOWN`, mais plano/contexto;
- busca de counterevidence `BUDGET_EXHAUSTED` junto de status `INCONCLUSIVE` e
  plano/contexto;
- identidade imutável ausente/ambígua junto de teto/status incapaz;
- disponibilidade de provider não evidenciada e ausência de capture congelado,
  com contexto da source;
- efeito de utility cost ou security sem baseline, com plano/contexto.

Licença `UNRECOGNIZED` com `PATTERN_ADAPTATION` ou `REFERENCE_ONLY` explícito e
status `PROPOSED` permanece `Preserved`. `integration_unknowns` isolado ou
composição preservada como proposta não é proxy suficiente. Toda violação
decidível tem precedência sobre `Unknown`.

Todos os predicados anteriores de catálogo, binding uno, provenance, status,
novelty, bytes fixos, hash/size/paginação, admission, lifecycle, replay,
runtime, refinamento e tools permanecem ativos. O14 não emite receipt, seal,
certificado nem veredito reservado.
