# Prompt R2 — oráculos ampliados da fila local

## Testemunhos

- cancelamento que vence ainda em fila impede entrada no handler;
- cancelamento concorrente de execução ativa é entregue exatamente uma vez ao handler; a fronteira de chamadas posteriores é exercitada no serviço real pela suíte R3 L1;
- 32 jobs bloqueados com limites 1 e 4 não excedem workers nem crescem uma goroutine por start;
- shutdown com worker que não reconhece cancelamento respeita deadline e não presume confirmação.

Todos os eventos têm canais/barreiras explícitos; a suite é offline e destinada também a `go test -race`.
