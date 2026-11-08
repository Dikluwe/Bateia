# Idea run domain

Defina contratos versionados para uma execução ponta a ponta nos modos live,
local e replay. O manifesto deve manter status do run, status por estágio, hashes
de entrada e saída, versões de regras, contagens, falhas e caminho da decisão.
Credenciais, tokens e cabeçalhos não pertencem a nenhum contrato.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
