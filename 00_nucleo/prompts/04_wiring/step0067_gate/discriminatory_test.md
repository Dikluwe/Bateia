# Gate discriminatório pré-candidato 0067

Materializar um gate L4 autocontido e *binding-free* para o contrato congelado do
planejador de fronteira mínima. O gate lê somente a fixture adversarial pública de 0067
e avalia snapshots declarativos; é proibido importar ou ler a futura implementação
`idea_radar_plan.go`, escrever oráculos produtivos ou adaptar expectativas ao candidato.

O gate deve preservar controles P, violar M01–M16 com testemunha e manter controles
deliberadamente opacos como `Unknown`. `Unknown` em caso negativo é mutante sobrevivente.
Deve discriminar bloqueador ativo, minimalidade, combinação de lacunas, escolha canônica,
binding-only, atomicidade de `UNKNOWN`/budget, imutabilidade, ausência de promessa e zero
impacto decisório. Deve ainda verificar as métricas externas exatas (8/8, 7/8, 11 contra
15, redução 4/15), os vereditos `OBSERVED`, `NOT_OBSERVED` e `INCONCLUSIVE`, e vinte
ordens com vetor idêntico.

O arquivo espelhado é `04_wiring/step0067_gate/discriminatory_test.go`. Ele não integra
produção, não contém lógica consumível pelo planejador e não julga candidato inexistente.
Autoridade: adversário/gate pré-candidato; entradas congeladas são passo, contrato,
baseline, budget, manifesto e APIs públicas 0066.
