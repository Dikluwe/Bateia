# Testes contratuais: 01_core/usecase/idea_radar_plan

Papel autor independente. Derivar testes somente do passo, contrato, baseline e budget
0067 congelados, sem ler o candidato.

Cobrir C01–C10; precedência exclusiva de conflito/ambiguidade; combinação das lacunas de
kind/origem; seleção mínima de Unknown; origem insuficiente; binding único e conservador;
controle C09 sem promessa de convergência; propagação de `UNKNOWN`; limites; operações;
imutabilidade; arrays não nulos; canonicalização e 20 permutações/nil-vazio idênticas.

O harness externo aplica as intervenções simuladas e compara taxa de remoção do bloqueador,
convergência e esforço ao baseline ingênuo. Exigir 11 contra 15 probes elegíveis, 8/8
bloqueadores removidos, 7/8 convergentes, redução 4/15 e os vereditos `OBSERVED`,
`NOT_OBSERVED` e `INCONCLUSIVE`. Eliminar M01–M16 com mutation score 1.0.

O futuro arquivo é `01_core/usecase/idea_radar_plan_test.go`, preferencialmente em package
externo. Congelar seu hash antes da primeira execução contra o candidato.

