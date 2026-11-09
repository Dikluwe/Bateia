# Regressões de admissão0167

L1. Origem ex-ante: [passo0167](../../../passos/0167-isolar-calibracao-e-validar-fontes-de-sintese.md). Contrato: [integridade de entradas e replay](../../../contratos/integridade-de-entradas-e-replay-v1.md).

CMP-01: permutar conexões do mesmo par ordenado com custos/unidades/ausência, fontes ou status divergentes; ambas as ordens devem retornar erro e relatório zero. Conexão única e duplicatas idênticas conservam relatório e hash anteriores; duplicata de custo nil continua desconhecida e não é zero. Não mutar input. Usar DecideReuse e testemunhos admissíveis existentes. Sem I/O concreto ou alterações em testes, registros e hashes históricos.
