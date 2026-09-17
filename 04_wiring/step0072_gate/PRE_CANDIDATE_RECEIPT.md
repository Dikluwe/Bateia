# Recibo pré-candidato — gate 0072

- Data: 2026-09-16.
- Papel: adversário/oráculo `cpp_collection_adversary`.
- Regime: protocolo completo Tekt como ensaio; segregação procedural, filesystem
  compartilhado e sem atestação física de isolamento.
- Entradas: passo 0072, prompts/testes contratuais congelados e baseline Git `03fd5f1`.
- Leitura proibida e não realizada: implementação candidata do passo 0072.
- Escrita: prompt proprietário e `04_wiring/step0072_gate/`.
- Budget de calibração: duas revisões focais e no máximo três execuções completas; foram
  feitas duas execuções completas e uma revisão focal, sem adaptar expectativas à solução.

## Comando

```sh
GOCACHE=/tmp/bateia-step0072-gate-cache go test ./04_wiring/step0072_gate -run '^TestStep0072Gate' -count=1
```

## Vetor observado

O gate compilou e rejeitou o estado pré-candidato. Falharam:

- reserva dirigida: somente `algorithm/a_distractor_000.hpp` foi coletado;
- diagnósticos dirigidos: coberturas de caminho omitido, compartilhado, erro de árvore e
  busca sem candidato ainda estavam vazias;
- enumeração local: `entries=20000 final=false`, testemunhando o corte silencioso.

Os demais controles que chegaram a uma observação compatível não constituem aprovação
parcial nem selo. Estado: **gate discriminatório materializado; selo pendente do candidato
e de mutation score 14/14**.

## Revisão focal R1

O primeiro oráculo exigia `SearchComplete=false` quando um caminho já enumerado era
omitido por bytes. Isso confundia completude da **busca de caminhos** com completude da
**leitura**. O passo permite enumeração completa com `skipped_paths` não vazio e exige,
nesse caso, `Status=COLLECTION_INCOMPLETE`. R1 removeu apenas a expectativa indevida sobre
`SearchComplete`; preservou candidatos, examinados, pulados, status, reason code,
confiança e gaps. No replay, permaneceram duas violações materiais: M01 (corte em 20.000)
e M07 (path match promovido a `PARTIAL`).
