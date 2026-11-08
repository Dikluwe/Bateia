# Especificação L0 — Interface TUI Interativa (02_shell/cli/tui.md)

## Responsabilidade
Implementar uma interface TUI (Terminal User Interface) interativa em L2 utilizando `charmbracelet/bubbletea` e `charmbracelet/lipgloss` para busca interativa, navegação pelos estágios do funil e exibição formatada dos relatórios do Bateia.

## Restrições Tekt
- **Camada**: L2 (Casca CLI/TUI).
- Conhece apenas L0 e L1 (`domain`, `usecase`, `parser`).
- Não conhece adaptadores concretos de L3.
- Comunica-se exclusivamente via `SearchFunnel` recebido por injeção.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
