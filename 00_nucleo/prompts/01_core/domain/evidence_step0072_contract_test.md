# Contrato independente — diagnóstico de recuperação, passo 0072

Derivar testes somente do passo 0072 e do baseline público congelado. Não ler nem adaptar
expectativas à implementação candidata.

Verificar valores exatos, tags JSON e round-trip de `RequirementCollectionStatus`,
`EvidenceRequirementQuery` e `RequirementCollectionCoverage`, além dos novos campos em
budget, plano e bundle. Arrays públicos vazios devem serializar como `[]`, nunca `null`.

O teste deve falhar pré-candidato. Unknown por coleta incompleta não pode compartilhar um
valor ambíguo com ausência nem ser omitido por `omitempty`.
