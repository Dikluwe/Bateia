# Contrato de oráculo black-box independente — Passo 0028 O13

**Revisão:** O13 vinculada a R4, catálogo v3 e protocolo O8
**Regime:** protocolo Tekt completo
**Papel:** autoridade segregada de ORÁCULOS
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

O avaliador preserva as APIs públicas:

```go
func Evaluate(contractInput, contractOutput []byte) Finding
func EvaluateSequence(exchanges []Exchange) []Finding
```

Somente `contract_input` e `contract_output` participam da decisão. IDs de caso,
attack hashes, records/configs, expected, markers, candidato e artefatos de
veredito permanecem fora da API e do evaluator.

## Entradas públicas congeladas

| Entrada | SHA-256 |
|---|---|
| Contrato R4 | `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322` |
| Catálogo v3 | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| Protocolo O8 | `c86b04d028bd6c0f33ddcff78ca323f70ab2b72d96bc86a6046b2ec6598755bb` |
| Coverage gap público O12 | `cd0007f7d8156a3fa824b57227ae2cee3fee5b339a2df1bda0f1460957b00039` |
| Falha discriminatória pública O12 | `ef03801e9fab6b6138eba9040a8704b506caf18f3e76a98b9d57cc45c443d30b` |
| Manifest protegido predecessor O11 | `0c7766e895338fb1cdba077e04b29531dd920819247e9d6fbb1179e64be20a7a` |

## Refinamentos normativos O13

`R4.MECHANISM_EXTRACT.CANONICAL_CAUSAL_MECHANISM_ID` reconstrói exatamente o
preimage fixo da seção 6: struct com campos na ordem
`schema_version,source,rule_id,mechanism_class,purpose,strategy,witness_hashes`,
source com ordem `repository,commit`, JSON compacto UTF-8, HTML escaping
desabilitado e sem newline. Source, RuleID, classe, propósito e estratégia usam
normalização local NFC/whitespace. `witness_hashes` contém todos e somente os
hashes de evidence e counter_evidence, deduplicados e ordenados por bytes. O
RuleID documental efetivo continua obrigatório quando toda evidência positiva é
`DOCUMENTARY`; evidência técnica ou mista conserva o RuleID pai.

`R4.ARTIFACT_READ.REAL_OPACITY_CLASSIFICATION` não usa licença
`UNRECOGNIZED`, warning `LICENSE_NOT_RECOGNIZED` nem
`integration_unknowns` explícitos conformes como proxy suficiente de opacidade.
Eles preservam a hipótese/artifact e seus limites, mas não tornam o gate
`Unknown` isoladamente. Baseline alvo realmente ausente, classe/identidade
insuficiente, budget esgotado, source indisponível ou disponibilidade não
evidenciada continuam opacidades concretas.

`R4.MECHANISM_EXTRACT.REAL_OPACITY_CLASSIFICATION` classifica `Unknown` quando
um baseline convencional ausente ou inválido é o motivo observável de um
mecanismo `PROVEN` conservar `novelty=UNKNOWN`. Para mecanismo `PARTIAL`, a
novidade já é `UNKNOWN` independentemente do baseline e não cria opacidade extra
por si. Toda violação decidível permanece precedente a `Unknown`.

Todos os predicados anteriores de catálogo fechado, binding uno, provenance,
status, novelty, fixed bytes, hash/size/paginação, admission, lifecycle, replay,
runtime, refinamento e tools permanecem ativos. O13 não emite receipt, seal,
certificado nem o veredito reservado.
