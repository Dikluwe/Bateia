# CLI transport for durable MCP dependencies

Permita que `CLIApp` receba a porta durable já composta e a repasse, junto com o
recommender síncrono, às duas formas existentes de iniciar MCP. Preserve todos os
comandos, flags e comportamentos anteriores. Não instancie store, fila, worker,
limite ou adapter L3; não contenha regra de run e não importe L3.

Esta é somente uma extensão de composição L2 conforme o contrato 0026.

