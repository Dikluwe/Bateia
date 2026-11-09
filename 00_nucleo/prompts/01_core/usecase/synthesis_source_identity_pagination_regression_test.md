# Regressões de admissão0167

L1. Origem ex-ante: [passo0167](../../../passos/0167-isolar-calibracao-e-validar-fontes-de-sintese.md). Contrato: [integridade de entradas e replay](../../../contratos/integridade-de-entradas-e-replay-v1.md).

SYN-01/02: testar porta restrita e StartImprovementSynthesis com Get devolvendo RunID divergente/vazio (zero páginas/idempotência/fila), repetição A→B→A e ciclos maiores (recusa antes da releitura), controle normal, cursores iguais em artefatos distintos, hash fallback, erros e teto1000. Falhas retornam ErrIntegrity/base ou recibo zero, sem efeitos posteriores; controles válidos chegam à próxima etapa por sentinela contada, sem criar execução real. Sem I/O concreto ou alterações em testes, registros e hashes históricos.
