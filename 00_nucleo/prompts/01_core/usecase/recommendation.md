# Recommendation service

Exponha a porta pura e síncrona `RecommendationService.Recommend` definida em
`00_nucleo/benchmark/0025/contract.md`. Valide e normalize toda a entrada antes
da primeira chamada a dependências. Execute somente `live`, componha
`RunIdeaPipeline` e `RankSoftwareAlternatives`, e produza uma única projeção
canônica compartilhável por CLI e MCP. Não duplique expansão, descoberta,
assessment, decisão, licença ou ranking. Não faça I/O, não aceite caminhos, não
chame IA e não execute conteúdo de candidatos.

Preserve falhas e Unknown, derive `COMPLETE|PARTIAL|INCONCLUSIVE` com a precedência
contratada, propague locales até a proveniência e trate licença desconhecida como
aviso não bloqueante para fit. Canonicalize conjuntos sem destruir ordem semântica
de ranking. Remova snippets e conteúdo bruto, publique somente mensagens seguras
e nunca serialize credenciais, headers, valores de ambiente ou erros crus de
adaptadores. Cancelamento impede novas coletas.
