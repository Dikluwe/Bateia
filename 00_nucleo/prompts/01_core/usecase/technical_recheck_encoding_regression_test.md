# Regressões de integridade0166

L1. Intenção ex-ante: [passo0166](../../../passos/0166-validar-bytes-testemunhos-e-politicas.md). Contratos: [integridade de entradas](../../../contratos/integridade-de-entradas-e-replay-v1.md) e [coleta](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#testemunhos-pertinentes-e-codificação-técnica--0166).

EVD-21: usar provedor puro contado e RecheckTechnicalSources com SBOM válido, U+FFFD literal e UTF-8 bruto inválido, tamanhos corretos. Inválido gera falha/incompletude, sem item/IntegrityVerified. Controle válido conserva bytes/hash/verificação após JSON e VerifyTechnicalEvidence. Não alterar request/provedor ou testes anteriores. Sem I/O concreto, dependência do produto no Lab ou reescrita de registros/hashes históricos.
