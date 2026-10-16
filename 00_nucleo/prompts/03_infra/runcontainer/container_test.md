# Prompt — oráculo do container durável L3

A partir apenas do contrato 0079, teste API, validação/cleanup observável, mesma
instância de serviço, recovery antes de start, execução real pela fila, shutdown e
reabertura. Inclua repetição e race. Não leia o candidato antes de congelar o teste;
declare como `Unknown` erros concretos de Close que não possam ser produzidos.

## Adendo 0106 — Espera funcional sob concorrência da suíte

SUCCEEDED continua obrigatório antes de reabrir o store. O prazo de espera é
um watchdog de diagnóstico, não uma promessa de latência de dois segundos.
Use até 30 segundos, respeitando o deadline do teste, reduza polling para 10 ms
e falhe imediatamente em estado terminal diferente do esperado. Reporte a
última versão/etapa ao expirar; não aceite RUNNING como sucesso.
