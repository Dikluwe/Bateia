# ADR-0016 — Baseline de conteúdo na comparação de insights

**Status:** ACEITO PARA O RECORTE 0135. **Data:** 2026-10-28. **Relacionados:** ADR-0008, ADR-0010, ADR-0015; passos 0133–0135.

## Contexto e decisão

Detalhes da fonte sem observação do alvo favorecem sugestões genéricas de capacidades que o alvo já pode possuir. O alvo local também pode ter alterações não commitadas. Usar HEAD como identidade desses bytes atribuiria uma revisão falsa ao conteúdo.

Adicionar `idea insights compare`, com contrato próprio, delegando a investigação existente de fontes v2. Receber baseline opcional de arquivos e observações do alvo como snapshot de conteúdo: digest do manifesto de paths/hashes, associação ao repositório declarada e escopo limitado aos arquivos fornecidos. Não apresentar esse digest como commit, classificação do alvo ou prova de completude.

Agrupar observações por fonte/classe/tipo e informar somente disponibilidade de declarações, com referências e próximos passos para conferir diferenças antes de alterar o alvo. Não inferir equivalência/diferença semântica por igualdade de strings, classe ou palavras-chave. Não projetar hipóteses genéricas de adição de capacidade. Um par de declarações permite revisão comparativa, sem demonstrar efeito ou aprovar adoção.

## Responsabilidades e consequências

L1 possui normalização, integridade, projeção e política consultiva; L2 possui transporte e apresentação. O normalizador de observações do ADR-0015 é reutilizado; o digest da baseline atua somente como chave interna, em namespace distinto das observações de fontes. Não há novo runtime ou efeitos. Analyze v1/v2 e a síntese durável mantêm seus contratos.

A comparação admite worktrees e snapshots sem exigir commit e preserva desconhecidos quando falta baseline ou observação. As declarações ainda exigem curadoria; sua associação ao alvo e sua interpretação não são autenticadas/reavaliadas. A necessidade de coleta/análise automática e de validar efeito no alvo permanece trabalho posterior.
