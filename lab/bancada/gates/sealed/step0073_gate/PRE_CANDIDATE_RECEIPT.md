# Recibo pré-candidato — gate 0073

- Data: 2026-09-16.
- Papel: adversário/oráculo `cpp_locator_adversary`.
- Regime: protocolo Tekt completo como ensaio, com segregação procedural no mesmo
  filesystem e sem atestação física de isolamento.
- Entradas: passo 0073, prompts/testes contratuais congelados e baseline público
  `03fd5f1` acrescido do passo 0072 materializado.
- Leitura proibida e não realizada: implementação candidata do passo 0073.
- Escrita: prompt proprietário e `04_wiring/step0073_gate/`.
- Budget: duas revisões focais e no máximo três execuções do corpus.

## Vetor pré-candidato esperado

O baseline pré-candidato não declara `EvidenceMatchClass`, seus quatro valores nem os
campos `MatchedPattern` e `MatchClass`; portanto o gate deve falhar na compilação antes
de observar o localizador. Esse vetor discrimina o baseline, mas não constitui selo,
mutation score ou veredito sobre candidato futuro.

## Execução observada

O comando foi disparado uma vez após a materialização do oráculo, mas o filesystem
compartilhado já continha uma edição candidata concorrente incompleta. A compilação
falhou por referências ainda indefinidas a `firstAdmissibleContentMatch`,
`sourceFromMatch`, `locateRuleContentMatches` e `opaqueRuleSignal` no caso de uso.

Essa execução confirma rejeição, mas **não** é apresentada como replay puro do baseline
nem como selo pré-candidato: a fronteira temporal já havia sido atravessada por outro
papel. O gate e suas expectativas não foram adaptados a essa saída.

## Revisão focal R1

O primeiro replay compilável mostrou duas expectativas mais fortes que o contrato:
exigia componente reutilizável mesmo sem o conjunto de provas que o cria e exigia
`INCONCLUSIVE` para hipótese com indício, embora o passo 0071 permita preservar o indício
sem chegar a `PROVEN` ou `ABSENT`. R1 passou a condicionar a precedência ao componente
quando existente e, para hipótese, proíbe somente os estados autoritativos e o componente.
O localizador, as linhas, as classes e todos os mutantes 0073 permaneceram inalterados.
