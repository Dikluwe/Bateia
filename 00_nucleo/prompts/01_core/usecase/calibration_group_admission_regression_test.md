# Regressões de admissão0167

L1. Origem ex-ante: [passo0167](../../../passos/0167-isolar-calibracao-e-validar-fontes-de-sintese.md). Contrato: [integridade de entradas e replay](../../../contratos/integridade-de-entradas-e-replay-v1.md).

CAL-01: rejeitar IdeaID duplicado adjacente ou separado por outro grupo, vazio e com espaços nas extremidades, com relatório zero e input intacto. Controle de IDs distintos mantém hash de relatório obtido antes do patch e uncertain fora dos rótulos conclusivos; folds e treino preservados. Sem I/O concreto ou alterações em testes, registros e hashes históricos.
