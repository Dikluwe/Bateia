# Contrato de oráculo black-box independente — Passo 0028 O20

**Revisão:** O20 sobre R4/catálogo v3/protocolo O8  
**Regime:** Tekt completo  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

`Evaluate` e `EvaluateSequence` permanecem binding-free. Identidade de caso,
hash externo, record/config, marker, expected, `semantic_state`, texto humano,
candidato e artefatos reservados não controlam a classificação.

Entradas públicas adicionais: failure O19
`26baf6c4a015afc5081bc712d6cadd789cb7990baf7d3da7e5f13dd4c8a7f0b3`,
gap O19 `1b508d8938f73d9e9df826da37480efdfd56f60abbf7bddfb6822ff06e4cee92`,
corpus e2e `75d701d8a9f916dd770cd1b851f7f6ecb22c04191f7065ef16731b77cc465870`,
replay `cb0201c1c053375951289dbb9156f65ba385380d6c2c13bffa9368b0083be50b`
e documentação `103b0107af7e67edf98ffc0d2d0f6a71d3e79a19747dc2cf8c6ac70a7bc7e9dc`.

Na fronteira end-to-end, uma hipótese `PROPOSED` sem validation plan preserva
opacidade quando contém uma das quatro relações públicas: composição múltipla
com integration unknowns; licença `UNRECOGNIZED` em `REFERENCE_ONLY`; busca
`NONE_FOUND_IN_SCOPE` com escopo não vazio e nenhum resultado; ou target
`UNKNOWN` com target gap refs vazio. Um plano não vazio torna o vizinho
decidível. O01–O03, os 160 negativos, os 12 positivos e a precedência
`Violated > Unknown > Preserved` permanecem ativos. Nenhum receipt, seal,
certificado ou veredito é emitido.
