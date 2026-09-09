# Controles independentes do passo 0060

Papel B, regime Tekt A/B sem atestação física. Raiz normativa: `00_nucleo/passos/0060-traducao-local-de-caminho.md`.

Criar testes externos em `01_core/parser/path_query_step0060_test.go`, package parser_test, usando somente APIs públicas identificadas nos tipos e testes legados. Não ler implementação candidata, corpus, benchmarks ou relatórios. Não executar os testes antes do congelamento e entrega ao integrador.

Verificar singular e plural naturais, capitalização inicial, fonte distinta, provenance exata, termos anteriores preservados, ordenação e deduplicação, estabilidade entre invocações e linhas. Excluir símbolo isolado, backticks, identificadores compostos, modo mixed, todos os marcadores de negação legados, palavras parecidas e termos derivados. Não traduzir parar para Remove nem criar mecanismos; preservar expansões legadas em consulta conjunta.

Fixture neutra: função com nome e arquivo neutros e variável local path, sem caminho nos índices de código; recuperar via consulta natural portuguesa e manter retrieval lexical. Controle com símbolo explícito português não deve adquirir tradução pelo glossário de indexação. Projeção de requisito pode localizar evidência, nunca verificar satisfação; inspeção sintática não verifica comportamento. Não fixar pesos ou posição de ranking desnecessariamente.

Congelar prompt e teste com SHA-256 em /tmp/bateia0060-b.sha256; nenhuma execução pelo papel B.
