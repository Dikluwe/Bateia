# Prompt L0 — CLI Idea Web Commands
Expor `discover-observable`, `replay` e `compare`. A descoberta exige um
`ObservableWebSearch` injetado por L4; L2 nunca instancia mock nem importa L3.
O subcomando `experiment` intercala N/F1/F2, exige `--no-token` e grava três
snapshots somente depois de preservar todas as respostas e falhas.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.
