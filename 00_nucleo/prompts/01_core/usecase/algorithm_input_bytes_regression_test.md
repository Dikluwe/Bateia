# Regressões de integridade0166

L1. Intenção ex-ante: [passo0166](../../../passos/0166-validar-bytes-testemunhos-e-politicas.md). Contratos: [integridade de entradas](../../../contratos/integridade-de-entradas-e-replay-v1.md) e [coleta](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#testemunhos-pertinentes-e-codificação-técnica--0166).

ALG-01: testar API FormulateAlgorithmInsights com fixtures independentes, UTF-8 inválido em pergunta/raciocínio e bytes de fonte/alvo, inclusive hash dos bytes substituídos e item não citado. Rejeitar sem relatório útil e sem mutar caller. Aceitar U+FFFD literal; congelar ID válido anterior e preservar canonicalização/papéis. Sem I/O concreto, dependência do produto no Lab ou reescrita de registros/hashes históricos.
