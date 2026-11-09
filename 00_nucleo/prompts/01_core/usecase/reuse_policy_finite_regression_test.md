# Regressões de integridade0166

L1. Intenção ex-ante: [passo0166](../../../passos/0166-validar-bytes-testemunhos-e-politicas.md). Contratos: [integridade de entradas](../../../contratos/integridade-de-entradas-e-replay-v1.md) e [coleta](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#testemunhos-pertinentes-e-codificação-técnica--0166).

DEC-01: testar DecideReuse por API com NaN/infinitos em cada limiar, limites inválidos e finitos válidos. Rejeitar sem relatório útil, preservar input; todo sucesso do controle sem assessments deve ser serializável e ter SHA256 do JSON efetivo da política. Hash válido anterior literal permanece. Sem I/O concreto, dependência do produto no Lab ou reescrita de registros/hashes históricos.
