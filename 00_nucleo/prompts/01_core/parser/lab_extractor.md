# Especificação L0 — Extrator do Lab de Projetos & Algoritmos (01_core/parser/lab_extractor.go)

## Responsabilidade
Implementar funções puras em Go em L1 para analisar o relatório de busca e os arquivos extraídos (Zero I/O), classificando funções de algoritmos (sorting, graph, tree, DP, crypto) e pareando com suas suítes de teste de borda e benchmarks.

## Restrições Tekt
- **Zero I/O**: Operações estritamente em memória a partir dos structs do domínio.
- **Camada**: L1 (Núcleo Puro).
- Sem dependências de L2, L3 ou L4.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.
