# Testes contratuais: 01_core/domain/evidence_coalition

Papel autor independente dos oráculos. Usar somente o contrato 0065 e os tipos públicos
baseline; não ler implementação candidata, testes do implementador ou recibos.

Verificar valores serializados exatos de schemas, papéis, polaridades, proveniência,
outcomes, impacto e status de relatório. Verificar JSON `snake_case`, arrays não nulos e
round-trip dos tipos sem alterar `EvidenceStatus`. Os testes não podem criar convenções
fora do contrato nem aceitar `omitempty` nas slices observáveis.

O futuro arquivo espelhado é `01_core/domain/evidence_coalition_test.go`, preferivelmente
em package externo. Congelar seu hash antes de executar contra o candidato.

