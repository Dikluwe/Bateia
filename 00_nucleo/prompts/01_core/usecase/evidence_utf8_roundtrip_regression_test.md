# Regressões de prefixo UTF-8 e replay

L1. Origem ex-ante [passo0164](../../../passos/0164-preservar-utf8-no-truncamento-e-replay.md); EVD-18 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#prefixos-utf-8-e-replay--0164).

Usar portas contadas existentes e APIs públicas. Para texto com ASCII, café, símbolo de3 bytes, emoji de4 bytes e marca combinante, percorrer todos os limites de corte. Construir esperado por runes, sem chamar helper de produção. Conferir prefixo literal máximo, UTF-8 válido, hash SHA256 independente, orçamento/Target.Bytes e SKIPPED_BUDGET/incompletude; completo conserva controle. Reabrir JSON várias vezes e reavaliar por AssessEvidence, exigindo conteúdo/hash/bytes idênticos.

Primeiro caractere maior que o saldo não cria item vazio; arquivo seguinte pode consumir o saldo ASCII. Fonte bruta inválida é FAILED, não sanitizada. Coleta coordenada e recuperação conservam bundles serializáveis e histórico de retomada; inputs e testes anteriores ficam intactos. Sem I/O concreto ou alteração de recibos históricos.
