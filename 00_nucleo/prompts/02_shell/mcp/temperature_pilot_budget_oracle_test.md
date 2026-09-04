# Prompt espelhado — orçamento MCP do piloto de temperatura

## Linhagem

- teste: `02_shell/mcp/temperature_pilot_budget_oracle_test.go`;
- camada: L2;
- papel: `CONTRACT_ORACLE_AND_ADVERSARY`;
- contrato: `00_nucleo/benchmark/0029/temperature-pilot-contract-r1.json`;
- candidato lido ou executado: não;
- alegação: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

## Superfície exigida

O pacote público `bateia/02_shell/mcp` expõe `JSONRPCError`,
`ToolCallBudget`, `NewToolCallBudget(*int)` e
`(*ToolCallBudget).Dispatch(method, decodeAndForward)`. Estado permanece privado;
o callback representa conjuntamente decoding, entrada no handler e efeitos
posteriores.

## Controles positivos

- `nil` preserva o modo legado ilimitado;
- limite `N=3` encaminha exatamente as três primeiras `tools/call`;
- `initialize`, `notifications/initialized` e `tools/list` são encaminhados sem
  consumir nem reiniciar o limite, inclusive quando intercalados após uma
  `tools/call` admitida;
- duas instâncias possuem contadores independentes.

## Ataques

- off-by-one: a quarta chamada deve receber exatamente código `-32010` e mensagem
  `Bateia tool call budget exhausted`;
- bloqueio tardio: o callback da quarta chamada não pode ser invocado;
- contador por tool/método: a sequência focal `recommend → tools/list → search →
  notifications/initialized → initialize → benchmark` admite somente as duas
  primeiras calls; labels downstream diferentes compartilham o mesmo limite de
  `tools/call`, e métodos não-call intercalados não o reiniciam;
- contador global: consumo em uma instância não afeta outra;
- construção com zero ou negativo falha fechada.

Os testes não abrem stdio, rede, store, modelo ou sessão MCP real. A verificação
de integração do servidor e o mutation score executado pertencem ao verificador
independente após implementação e selo.
