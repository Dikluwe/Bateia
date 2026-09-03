# Contrato de oráculo black-box independente — Passo 0028 O19

**Revisão:** O19 sobre R4/catálogo v3/protocolo O8  
**Regime:** Tekt completo  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

`Evaluate` e `EvaluateSequence` permanecem binding-free. Identidade de caso,
hash externo, record/config, marker, expected, texto humano, candidato e
artefatos reservados não controlam a classificação.

Entradas públicas adicionais: failure O18
`2f2c1b9efa519198e83400c21a71c6bc44dae440887ae24f47dfb08949f0f1ed`,
gap O18 `06b1e154d0f2b13229539e2f159006698c0f89dc31a007527bf4e312b64c56b8`
e projeção sanitizada
`7b6a7f10487b7ff599c73f48f2c83b4c5ee6aec5a29e59b117e7de22ae114da8`.

O04–O10 são reconhecidos somente após validar o documento canônico schema `1`,
artifact `hypothesis_gaps` e entry com `hypothesis_id`, `code` e `source_ref`
strings não vazias. `code` deve ser igualdade exata a um dos sete enums
normativos da projeção pública. A ausência de qualquer campo obrigatório não
ativa esse predicado. O01–O03, os 160 negativos, os 12 positivos e a precedência
`Violated > Unknown > Preserved` permanecem ativos. Nenhum receipt, seal,
certificado ou veredito é emitido.
