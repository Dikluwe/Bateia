# Contrato de falhas do estágio de rastreabilidade

Dar prompt proprietário exclusivo aos testes de falha existentes, sem mudar
suas expectativas. Falha de construção, explicação, emissão ou hashing deve
marcar manifest FAILED e estágio trace FAILED sem anunciar outputs incompletos.
Também verificar que erros antes da execução não fabricam estágio válido e que
quotas/proteções impedem emissão parcial anunciada. Sem executar candidatos ou
rede. Preservar os contratos de pertencimento aos artefatos de 0099.
