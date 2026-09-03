# Contrato de oráculo black-box independente — Passo 0028 O15

**Revisão:** O15 vinculada a R4, catálogo v3 e protocolo O8  
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
| Coverage gap público O14 | `5ea3436e717f536746d251af7f0921d3705b4b2fdc3d9a30c05d55f53fe4be9a` |
| Falha discriminatória pública O14 | `fad0a998b7fd007cffe450e63a927ada51dd056a208caac95adb1f0db7f66e85` |
| Manifest protegido predecessor O14 | `ac8d10da663d301f4f10beb3c35113db7b27f61a50b2e037086642b7c63b1343` |

## Refinamento normativo O15

`R4.MECHANISM_EXTRACT.REAL_OPACITY_CLASSIFICATION` distingue o `PARTIAL`
documental aceito de um grupo técnico que não alcança a relação mínima de
witnesses congelada no catálogo. O segundo é um limite epistemológico real;
nenhum gap inventado nem identificador de caso é necessário para reconhecê-lo.

`R4.ARTIFACT_READ.REAL_OPACITY_CLASSIFICATION` passa a considerar o tipo
declarado do artefato. Hipóteses completas são decididas pelas relações
`target_state`/evidência ausente, composição, license/transfer/status e busca de
counterevidence. Em `hypothesis_gaps` e `mechanism_gaps`, somente o campo
controlado `code` pode materializar ausência ou ambiguidade de baseline,
compatibilidade, witness, identidade e provider, ou budget esgotado. Títulos,
descrições, expected effects e demais texto livre nunca controlam a classe.

`UNRECOGNIZED` isolado, `PATTERN_ADAPTATION` ou `REFERENCE_ONLY` com status
`PROPOSED`, `integration_unknowns` isolado e `PARTIAL` documentary-only
permanecem `Preserved`. Uma falha decidível de integridade, provenance, schema,
ordem, safety, status ou qualquer predicado anterior sempre tem precedência e
resulta em `Violated`.

Todos os demais predicados O14 permanecem ativos. O15 não emite receipt, seal,
certificado nem veredito reservado.
