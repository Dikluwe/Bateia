# Busca pura de implementações Go — passo 0044

Origem exclusiva de `01_core/parser/implementation_search.go`. Contrato congelado em
`00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md`.

Implementar RankGoImplementations(query, files, limit) sem I/O: parse AST íntegro,
funções, métodos, tipos e valores declarados; linha, hash SHA256 dos bytes e motivo.
Declaração exata vence qualquer coincidência lexical. Separar identificadores CamelCase,
acrônimos e snake_case. Remover stopwords PT/EN; glossário explícito genérico PT/EN
quando necessário é expansão lexical, nunca inferência semântica. Ordenação é heurística.
Canonicalizar caminhos, entradas e empates; duplicatas não alteram resultados. Reportar
parse inválido e limites; ausência é somente no escopo indexado. Excluir metadados,
oráculos, gates, fixtures, diretórios gerados e testes salvo busca explícita por testes.
Expor predicado puro de escopo para adaptador local usar antes de abrir arquivos.
Usar status complete, partial e not_evaluated; partial cobre erros e truncamento.

Consultas multilinha distribuem resultados: declarações exatas globais primeiro,
melhor declaração por requisito em ordem canônica, depois ranking global. Calcular
a cobertura lexical por linha no mesmo AST, sem nova análise sintática ou I/O.

## Passo 0045 — recuperação funcional
Contrato adicional congelado em `00_nucleo/passos/0045-recuperar-perguntas-funcionais.md`.
Método go_ast_functional_v2. Identificadores AST do corpo, sem comentários ou literais,
fornecem pertinência operacional lexical. Funções com corpo e termos pertinentes têm
preferência sobre tipos/valores; identificadores exatos mantêm prioridade absoluta.
Glossário PT/EN finito e genérico; conjuntos evitam inflação por repetição. Perguntas
multilinha continuam distribuídas em ordem canônica. Retornar Retrieval lexical.
Após sementes lexicais, auxiliar functional_retrieval expande chamadas diretas seguras
por somente um salto, com no máximo floor(MaxResults/3) hits exclusivos da expansão.

Identificador exato: linha inteira com um identificador, ou identificador explícito
com fronteira CamelCase/snake_case/acrônimo em pergunta. Palavra minúscula isolada
no meio de frase não recebe prioridade exata sobre relevância funcional.

0046: diversificar caminhos dos hits em rodadas, após distribuição de requisitos e antes expansão/truncamento. Identificadores exatos mantêm prioridade. Método go_ast_functional_v3. Correção do consumo de orçamento por símbolos duplicando caminhos observada no corpus real.

0046 — Veredito experimental: a tentativa v3/diversificação acima foi descartada após comparação congelada sem ganho de recall@12 e com mais retornos não rotulados. Contrato ativo de ranking permanece v2/0045; histórico experimental preservado, não implementação ativa.

0047: método go_ast_corpus_v1 substitui apenas pontuação v2. Coletar documentos
nomeados durante o único parse, incluindo não matches, e após coleta calcular
estatística/scoring com corpus_scoring. Preservar exatos, distribuição multilinha,
chamadas, conflitos, hashes e orçamento. Hit.TermWeights contém contribuições
numéricas usadas; Reason explica método e estatística limitada ao corpus lido.

0048: método go_ast_query_v1. InterpretImplementationQuery fornece Terms/Symbols
canônicos para pontuação global e distribuição de requisitos. Terms reúne expansão
lexical interpretada; Symbols controla bônus exato preservado. Expor a interpretação
no relatório, sem reescrever query original nem alterar corpus_scoring,
functional_retrieval, tokenização/glossário AST ou frequência documental.

0051: durante o parse coletar canal separado de grupos Doc anexados a FuncDecl, somente leitura integral (Size não maior que bytes disponíveis). Termos interpretados de consultas natural/mixed; símbolo puro excluído. Usar implementationTerms, sem alterar documentos/IDF/código ou expansão de chamadas. Após ranking e expansão existentes, acrescentar documentação em slots restantes sem duplicar path/line/symbol; score número de termos distintos até 99, empate path/line/symbol. Hit definition/documentation, intervalo declaração e Doc, SHA256 e razão documental sem prova. Candidatos omitidos pelo limite marcam partial/Limited.

0052 r2: aplicar admissão por ação do adendo congelado separadamente a cada
linha/candidato. Campo de código une termos existentes de nome/corpo; Doc usa
somente palavras do grupo anexado. Unir apenas Terms de linhas admissíveis para
pontuação global/documental e zerar contribuição por requisito rejeitado. Preservar
estatísticas IDF, prioridade exata e termos originais; não misturar canais. Report
e hits explicam filtro lexical por ação/contexto quando houver famílias de ação.

## Passo 0059 — preservar canal documental
Aplicar 00_nucleo/passos/0059-preservar-canal-documental.md: preparar candidatos
documentais únicos contra todos os lexicais admitidos e calcular reserva antes de
relations.expand. Expandir e cortar com orçamento de código limit-reserva, mantendo
exatos e impedindo relações órfãs de sementes cortadas. Depois preencher documentação
com o procedimento existente até limite total original, sem alterar ranking/pesos,
consulta, AST ou comprovação. Omissões continuam Limited/partial. Reserva positiva
acrescenta Reason explícita de reserva documental limitada sem comprovar comportamento.
