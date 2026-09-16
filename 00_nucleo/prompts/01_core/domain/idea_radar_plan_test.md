# Testes contratuais: 01_core/domain/idea_radar_plan

Papel autor independente dos oráculos. Usar somente passo, contrato, baseline e tipos
públicos anteriores ao candidato; não ler implementação futura nem adaptar expectativas.

Verificar schema/status, seis kinds de probe, três requisitos de origem, prioridades,
campos e tags `snake_case`, arrays não nulos, round-trip e reutilização dos tipos 0065/0066.
Exigir `VerdictChanges=0` e ausência de campos que afirmem readiness, execução ou mudança
de decisão.

O futuro arquivo é `01_core/domain/idea_radar_plan_test.go`, preferencialmente em package
externo. Congelar o hash antes da primeira execução contra o candidato.

