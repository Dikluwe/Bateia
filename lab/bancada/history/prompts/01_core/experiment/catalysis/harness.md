# Prompt — Núcleo puro do experimento de catálise 0030

Materialize em L1 um contrato puro para o Passo 0030.

O componente deve:

- decodificar JSONL do Codex preservando chaves duplicadas e a ordem lexical;
- reconhecer `SINGLE_EXEC_ID`, `DUAL_ID_CODEX_0_149` e `OPAQUE`;
- parear `item.started` e `item.completed` pelo vetor integral de IDs;
- preservar `Unknown` para formato ou elo causal opaco;
- classificar saída separando validade estrutural, meta editorial de 18.000
  caracteres e teto rígido de 65.536 bytes;
- produzir a fila selada de seis células somente após preflight `PASS` e nunca
  cancelar células posteriores por falha comportamental anterior;
- validar controle com web e zero MCP, tratamento com uma
  `bateia_recommend` concluída antes de uma a três chamadas web;
- manter `INCONCLUSIVE` neutro;
- atribuir catálise somente quando direção, web resolvível e uso final formarem
  caminho observável.

L1 não lê arquivos, relógio, rede, ambiente ou processos. Entradas e saídas são
valores ou bytes fornecidos por portas externas. Violações usam reason codes
estáveis do contrato selado em
`00_nucleo/benchmark/0030/contract-seal-r2.json`.
