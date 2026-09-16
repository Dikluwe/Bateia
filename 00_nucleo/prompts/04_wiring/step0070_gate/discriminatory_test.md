# Prompt L0 — `04_wiring/step0070_gate/discriminatory_test.go`

Implemente, antes de existir candidato, um gate discriminatório binding-free para o passo
0070. Leia `00_nucleo/benchmark/0070/adversarial-fixtures.json`, confira os hashes das
cinco autoridades congeladas e classifique somente observáveis públicos como `Preserved`,
`Violated` ou `Unknown`. Não importe, leia nem execute implementação ou oráculos 0070.

Os controles devem preservar separação entre evidência e inferência, precedência de
controvérsia e `RESEMBLES` epistemicamente desconhecido. Opacidade global válida deve
permanecer `Unknown`. Rejeite promoção de inferência, dupla contagem da mesma origem,
transitividade inventada, promoção causal de `RESEMBLES`, controvérsia apagada,
`UNKNOWN` promovido, referência órfã, dependência de ordem, parcial em falha, mutação do
input e I/O. Exija score de mutações `1.0`, duas repetições e vinte ordens determinísticas.
Registre custo no fixture; o gate prova apenas poder discriminatório do fragmento público,
não verdade histórica nem conformidade de uma candidata futura.
