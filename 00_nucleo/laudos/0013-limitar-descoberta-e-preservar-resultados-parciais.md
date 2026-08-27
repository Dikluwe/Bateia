# Laudo 0013 — Descoberta limitada e resultados parciais

**Estado**: `EM VALIDAÇÃO`  
**Data**: 2026-08-26

## Materializado

- `DiscoveryExecutionPolicy` com deadline total, attempt timeout, concorrência e
  máximo de tentativas;
- estados `COMPLETE`, `PARTIAL` e `FAILED` e causas de término;
- receipt agregado e receipts de attempts;
- distinção entre `SUCCEEDED`, `EMPTY`, `TIMED_OUT`, `CANCELLED`, `FAILED` e
  `NOT_STARTED`;
- propagação de timeout individual ao adaptador Sourcegraph;
- cancelamento estruturado por deadline, SIGINT e SIGTERM;
- preservação de candidatos e attempts concluídos antes do cancelamento;
- progresso conciso em stderr;
- escrita atômica do snapshot antes de retornar erro parcial ou terminal;
- warning `NO_SEMANTIC_EXPANSION` quando o expansor não propõe queries;
- compatibilidade de leitura: receipt é opcional em snapshots legados.

## Ataques executados

- candidato concluído antes do deadline não é perdido;
- expansão não iniciada recebe `NOT_STARTED`;
- timeout recebe classe `TIMEOUT`, nunca `EMPTY`;
- attempt timeout maior que o total falha antes de tocar o provedor;
- output parcial existe antes do erro retornado pela CLI;
- falha total continua sendo erro e agora preserva receipt `FAILED`;
- token continua removido das mensagens do Sourcegraph;
- suíte completa passou, inclusive `httptest` loopback fora do sandbox;
- linter V15 e `git diff --check` passaram.

## Reexecução E-0012

A repetição foi preparada com os parâmetros congelados:

```text
strategy=lexical
max-per-query=3
concurrency=2
timeout=3m
attempt-timeout=20s
output=benchmark/0012/run/discovery-0013.json
```

Uma primeira solicitação externa foi rejeitada até haver aprovação explícita do
payload e destino. Depois da autorização do usuário, a execução foi enviada ao
Sourcegraph público sem workaround, outro provedor ou inserção manual.

O run terminou no deadline com `PARTIAL/TOTAL_TIMEOUT`: 10 de 16 expansões foram
iniciadas e concluídas, 6 ficaram `NOT_STARTED`, 18 attempts expiraram e um
fallback terminou com 3 candidatos preservados. O snapshot foi escrito antes do
erro e contém `NO_SEMANTIC_EXPANSION`.

Assim, o gate operacional real passou. Os três candidatos eram ruído produzido
por `decidir entre`, estavam em `HEAD` e nenhum correspondia às quatro sementes
públicas verificáveis; o recall competitivo foi `0/4`. Isso mantém o E-0012 em
`FAIL_DISCOVERY`, agora atribuído ao vocabulário/query e ao baixo yield do
provedor dentro do budget, não à falta de observabilidade.

## Veredito

O fragmento local e determinístico está implementado e atacado, e a Fase A
demonstrou em ambiente real término dentro de três minutos com snapshot, receipt,
progresso e resultado parcial preservados. O passo permanece `EM VALIDAÇÃO`
somente pelos gates ainda não exercitados em ambiente real, como Ctrl-C e
repetição canônica; a correção que motivou o passo foi comprovada.

Este laudo não afirma melhora de recall. Vocabulário e regras de autofagia
permaneceram inalterados conforme o escopo.
