# Contrato de oráculo black-box independente — Passo 0028 O18

**Revisão:** O18 sobre R4/catálogo v3/protocolo O8  
**Regime:** Tekt completo  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

`Evaluate` e `EvaluateSequence` permanecem binding-free. O avaliador não usa
identidades de caso, hashes externos, records/configs, markers, expected, texto
humano, candidato ou artefatos reservados.

Entradas adicionais: failure O17
`3af90e4e7c65cca9b4e797b9cbe9f34315215375de055cfcb211e23e1da21379`,
coverage gap O17
`33af523f2c017e1d3d0f85b196afbd906aaf0f22325e5fbd50b78ff582a1c6c8`
e manifest O17
`736019da8a8b087bedfbae7d92b8c13fd83ac2aa69a6e50489abdcdca4686090`.

O18 implementa literalmente as sete linhas públicas O04–O10. Relações de
hipótese exigem todas as conjunções publicadas, inclusive status
`INCONCLUSIVE` e plano não vazio. Relações de gap exigem artifact tipado e as
conjunções exatas de tokens controlados. O01–O03, P01–P12, as 160 violações e a
precedência `Violated > Unknown > Preserved` são preservados. Nenhum receipt,
seal, certificado ou veredito é emitido.
