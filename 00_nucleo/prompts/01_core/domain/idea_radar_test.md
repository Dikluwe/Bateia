# Testes contratuais: 01_core/domain/idea_radar

Papel autor independente dos oráculos. Usar somente intenção, contrato, baseline e tipos
públicos anteriores ao candidato; não ler implementação futura nem adaptar expectativas.

Verificar constantes exatas, tags `snake_case`, arrays não nulos, round-trip, cinco kinds,
três stances, seis outcomes, estados do relatório, draft 0065 conservador e
`VerdictChanges=0`. Verificar que os tipos reutilizam `EvidenceOriginStatus` e
`EvidenceCoalitionRole` sem alterar schemas existentes.

O futuro arquivo é `01_core/domain/idea_radar_test.go`, preferencialmente em package
externo. Congelar o hash antes de executar contra o candidato.
