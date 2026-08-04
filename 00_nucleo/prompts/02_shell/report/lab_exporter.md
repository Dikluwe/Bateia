# Especificação L0 — Exportador do Lab de Projetos (02_shell/report/lab_exporter.go)

## Responsabilidade
Implementar formatadores de saída em L2 para transformar o dataset do Lab de Projetos em arquivos `JSON` estruturados e cadernos de exercícios/benchmarks em `Markdown`.

## Restrições Tekt
- **Camada**: L2 (Casca).
- Conhece apenas L0 e L1 (`domain`, `parser`).
- Sem dependências diretas de adaptadores concretos de L3.
