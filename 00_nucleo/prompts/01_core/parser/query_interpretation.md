# Interpretação pura da consulta — passo 0048

Origem: `00_nucleo/passos/0048-distinguir-consulta-e-simbolos.md`.
Implementar InterpretImplementationQuery sem I/O. Canonicalizar linhas por trim,
deduplicação e ordenação, preservando Original. Usar token.IsIdentifier para
símbolos Go: linha inteira identificadora, conteúdo identificador entre um par
de backticks simples ou critério CamelCase/snake_case0047 no texto. Grafia original
em Symbols; ordenar e deduplicar. Modo symbol se só símbolos, mixed com texto,
natural no restante. Delimitadores malformados/conteúdo inválido são texto natural.
Em palavras naturais, remover sem/não/nao/nunca/not/without/never sem distinguir
maiúsculas antes da expansão lexical existente. Registrar Constraints por marcador
normalizado, Text linha integral, Status not_verified; nunca inferir escopo ou
cumprimento. Símbolos preservam seus tokens, inclusive operadores. Terms ordenados
e deduplicados usam implementationTerms existente, sem mudar vocabulário do corpus.

Símbolos explícitos conservam todos os segmentos normalizados de
implementationWords, inclusive os presentes nas stopwords (por exemplo s/use),
além da expansão lexical existente. Essa preservação vale somente para a consulta;
não modificar a tokenização ou os termos das declarações do corpus.

0052: cumprir catálogo finito de
00_nucleo/passos/0052-expandir-vocabulario-com-proveniencia.md. Coletar somente
tokens naturais originais normalizados; proteger identificadores explícitos e
conteúdo entre backticks. Inicial maiúscula natural permite expansão. Após
interpretar a linha, se não houver Constraints, aplicar regras exatas aos gatilhos
originais deduplicados e ordenados; nunca aplicar cascata ou substring. Acrescentar
derivações a Terms e registrar Expansions por source ordenado, translation antes
de mechanism, Terms de cada recibo ordenados e únicos, Status hypothesis e
RuleID 0052-v1:<família>:<tipo>. Preservar termos e símbolos anteriores e Original.
Não alterar glossário do corpus, ranking ou classificação de requisitos.

0052 r2: cumprir 00_nucleo/passos/0052-adendo-admissao-por-acao.md.
Helper puro de admissão por linha e campo: todas as famílias de ação presentes
nos RuleID 0052-v1 precisam de translation ou mechanism coincidente. Para
compression/decompression, mecanismo exige também block/blocks/corrupt/corrupted
presentes simultaneamente nos Terms da consulta e no campo examinado. Agrupar
recibos por família; linhas sem ação são admissíveis como antes. Não misturar
campos/canais nem inferir contexto ausente da consulta.

## Passo 0058 — Plano lexical local
Após determinar Mode, somente linhas natural consultam PlanDiscoveryQuery com a
linha original completa. Apenas planos hypothesis acrescentam ActionTerms e
ObjectTerms a Terms, deduplicados e ordenados com os termos existentes. Registrar
dois recibos adicionais, ação antes de objeto, Source igual à linha original,
Kind translation, Status hypothesis e RuleID 0058-v1:<plan.Version>:<RuleID do
grupo>. Cada recibo preserva exatamente os termos e a ordem do respectivo grupo
do plano. Não expandir recursivamente. Preservar Original, Symbols, Constraints,
expansões 0052, consulta do provedor, admissão, pesos, limites e seleção documental.
Planos literal_only e modos mixed/symbol não recebem a expansão 0058. Hipóteses
lexicais nunca comprovam satisfação de requisitos.

## Passo 0060 — Tradução local de caminho
Somente em Mode natural e sem marcadores de negação, traduzir as palavras naturais
originais elegíveis caminho e caminhos para path. Usar os sources normalizados,
deduplicados e ordenados já coletados; nunca expandir símbolos, conteúdo entre
backticks ou termos derivados. Para cada source distinto, acrescentar recibo com
Source minúsculo, Terms [path], RuleID 0060-v1:path:translation, Kind translation e
Status hypothesis. Preservar Original, Symbols, Constraints, expansões anteriores
e termos; deduplicar e ordenar Terms globais. Não alterar o glossário de indexação,
o planejamento 0055/0058, busca pública, ranking, orçamento ou seleção de arquivos.
Não inferir ações, funções, mecanismos ou satisfação comportamental da tradução.
