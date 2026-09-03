# Contrato de oráculo black-box independente — Passo 0028 O23

**Revisão:** O23 sobre R4/catálogo v3/protocolo O22  
**Regime:** Tekt completo  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

`Evaluate` e `EvaluateSequence` permanecem binding-free. A classificação usa
somente `ContractInput`, `ContractOutput`, shapes tipados e relações
contratuais. Identidade de caso ou hipótese, valores de hash, marker,
`expected`, `semantic_state`, tamanho-fingerprint, texto humano, candidato e
artefatos reservados não controlam decisões.

Entradas novas congeladas: protocolo O22
`7436cb8007ce0b9ef96a31ea22321ba3ba260ab06656ca414798f2715f215daf`,
invalidação O8→O22
`52c8807d996665e1a303c6a060c8435762e8fab2fee5a2eb12cbe95388bb64b0`
e inconsistência pública O20
`ebd2fa45bde1c8e672b54c22fc737b2a20b9664637ec0164ab6e194c8e6b2eec`.

Para `ARTIFACT_READ` de `hypothesis_gaps`, `Unknown` requer simultaneamente:
retorno sem erro; resposta/documento do mesmo artifact; documento schema `1`
compacto, com ordem externa `schema_version,artifact,entries`; exatamente uma
entry fechada na ordem `hypothesis_id,code,source_ref`; identidade apenas
validada pelo formato público; um dos sete pares `code/source_ref` normativos
O22; bytes, size e SHA coerentes; e observação de publicação que exponha o
artifact com o mesmo SHA. Um único requisito diferente não ativa opacidade.

Todos os predicados negativos, P01–P12 e O01–O03 predecessores permanecem.
Precedência total: `Violated > Unknown > Preserved`. O23 não emite receipt,
seal, certificado ou veredito.
