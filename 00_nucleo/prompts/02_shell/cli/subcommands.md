# Prompt Spec L0: Roteamento de Subcomandos CLI

## Objetivo
Implementar a arquitetura de subcomandos na CLI do Bateia (`similarity`, `search`, `lab`, `mcp`, `tui`) mantendo retrocompatibilidade com as flags da raiz.

## Restrições
- Camada: L2 (Shell)
- Delega a lógica de negócio para os orquestradores em L1.
- Saídas locais usam nome neutro no diretório escolhido pelo usuário; a CLI não grava
  resultados de execução em `00_nucleo` nem cria um depósito estrutural `relatorios/`.
