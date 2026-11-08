# Prompt Spec L0: Roteamento de Subcomandos CLI

## Objetivo
Implementar a arquitetura de subcomandos na CLI do Bateia (`similarity`, `search`, `lab`, `mcp`, `tui`) mantendo retrocompatibilidade com as flags da raiz.

## Restrições
- Camada: L2 (Shell)
- Delega a lógica de negócio para os orquestradores em L1.
- Saídas locais usam nome neutro no diretório escolhido pelo usuário; a CLI não grava
  resultados de execução em `00_nucleo` nem cria um depósito estrutural `relatorios/`.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
