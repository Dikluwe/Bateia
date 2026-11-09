# Regressões de integridade0168

L1. Intenção ex-ante: [passo0168](../../../passos/0168-preservar-integridade-e-determinismo-do-oraculo.md). Contrato: [integridade e replay](../../../contratos/integridade-de-entradas-e-replay-v1.md#ass-01--snapshot-íntegro-antes-da-união-cega).

MCDA-01: testar AnalyzeMulticriteriaDecision com escala[-MaxFloat64,MaxFloat64] e[-1e308,1e308], ambos sentidos, endpoints/interior/intervalos e ausência. Esperados independentes0/0.25/0.5/0.75/1; conferir utilidade, contribuições, scores e JSON finitos. Controle ordinário conserva hash anterior e caller intacto. Sem I/O concreto ou atualização de registros/hashes históricos.
