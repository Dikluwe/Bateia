# Gate discriminatório pré-candidato 0066

Materializar um gate L4 autocontido e *binding-free* para o contrato congelado do radar
de convergência. O gate pode ler somente a fixture adversarial pública de 0066 e deve
avaliar snapshots declarativos; é proibido importar ou ler a futura implementação
`idea_radar.go`, escrever testes produtivos ou adaptar expectativas ao candidato.

O gate deve demonstrar controles positivos como `Preserved`, mutações M01–M14 como
`Violated` com testemunha, e controles deliberadamente opacos como `Unknown`. `Unknown`
em negativo é sobrevivente. Deve ainda discriminar externamente `OBSERVED`,
`NOT_OBSERVED` e `INCONCLUSIVE`, executar vinte ordens/repetições determinísticas com o
mesmo vetor e exigir `mutation_score=1.0`.

O arquivo espelhado é `04_wiring/step0066_gate/discriminatory_test.go`. Ele não integra
produção, não contém ground truth consumível pelo radar e não emite veredito sobre uma
implementação candidata inexistente. Autoridade: adversário/gate pré-candidato; entradas
congeladas são passo, contrato, baseline, budget e manifesto 0066.
