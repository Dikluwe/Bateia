# Especificação L0 — Interface TUI Interativa (02_shell/cli/tui.md)

## Responsabilidade
Implementar uma interface TUI (Terminal User Interface) interativa em L2 utilizando `charmbracelet/bubbletea` e `charmbracelet/lipgloss` para busca interativa, navegação pelos estágios do funil e exibição formatada dos relatórios do Bateia.

## Restrições Tekt
- **Camada**: L2 (Casca CLI/TUI).
- Conhece apenas L0 e L1 (`domain`, `usecase`, `parser`).
- Não conhece adaptadores concretos de L3.
- Comunica-se exclusivamente via `SearchFunnel` recebido por injeção.
