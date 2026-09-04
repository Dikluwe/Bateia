# Prompt espelhado — configuração do orçamento MCP no wiring

## Linhagem

- teste: `04_wiring/temperature_pilot_config_oracle_test.go`;
- camada: L4;
- papel: `CONTRACT_ORACLE_AND_ADVERSARY`;
- contrato: `00_nucleo/benchmark/0029/temperature-pilot-contract-r1.json`;
- candidato lido ou executado: não;
- alegação: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

## Superfície exigida

`parseMCPMaxToolCalls(raw string, present bool) (*int, error)` traduz somente a
configuração recebida por L4. `present=false` seleciona modo legado ilimitado;
`present=true` aceita exclusivamente inteiro decimal ASCII canônico positivo.
O chamador deve tratar erro antes de construir ou iniciar `Serve`.

## Controles e ataques

- ausência retorna `(nil, nil)`;
- `"3"` retorna ponteiro para `3`;
- vazio, zero, negativo, sinal explícito, whitespace, fração, texto e overflow
  retornam `(nil, erro)`.

O teste não inicia `main`, servidor, processo, MCP, modelo ou ação externa. Ele
congela somente a tradução fail-closed anterior a `Serve`.
