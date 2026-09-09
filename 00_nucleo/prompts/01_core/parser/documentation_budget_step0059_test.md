# Testes independentes do orçamento documental — passo 0059

Papel B, regime Tekt A/B executado sem atestação de isolamento. Entradas permitidas: contrato `00_nucleo/passos/0059-preservar-canal-documental.md` incluindo esclarecimento anterior ao candidato, diretrizes `.agents/AGENTS.md`, tipos de domínio e testes legados. Não consultar produção candidata, corpus, benchmarks ou relatórios; não executar antes do congelamento entregue ao integrador.

Gerar `01_core/parser/documentation_budget_step0059_test.go`, pacote externo `parser_test`, com fixtures Go sintéticas em memória e APIs públicas observadas nos testes legados. Não alterar testes existentes.

Oráculos: saturar seleção com declarações lexicais e documentação exclusiva, conferir limites 1/2/8 e fração de um quarto limitada por candidatos únicos e exatos. Identidades lexicais do conjunto completo não podem reaparecer como documentos depois do corte. Sem documentos exclusivos não reduzir código. Sem código preencher todas as vagas documentais. Preservar scores, ordem, campos e hashes comparando hits com a seleção de orçamento amplo sobre os mesmos bytes; permutações e duplicações de entrada não podem mudar hits ou contabilidade. Exigir partial e Limited quando houver omissão, e razão de reserva documental sem alegação de prova de comportamento.

Relações devem partir somente de sementes de código preservadas: comparar seleção mista com seleção de código no orçamento reduzido, verificar todas as referências RelatedTo e ausência de expansão da documentação. Validar documentos retornados com releitura pública, comportamento not_evaluated sem fatos, projeção documentation_located e restrições not_verified. Fixtures não usam nomes do corpus.

Congelar SHA-256 do contrato, deste prompt e do teste formatado antes de executar. Hash identifica artefatos e ordem, não comprova isolamento físico. Integrador executa testes e gates; divergências do contrato não autorizam ajustar oráculos à implementação.

Revisão r1: corrigir somente a fixture baseline do teste de relações. O integrador informou que r0 comparava quatro declarações com seis, alterando o denominador lexical (scores 1232/1306); isso não viola o contrato. Manter todas as declarações, caminhos e linhas nos dois cenários, neutralizando somente o comentário documental em z.go no baseline. Preservar as expectativas de orçamento, scores, referências e ausência de expansão documental, sem mudanças nos outros testes. Não ler produção candidata nem executar antes do novo congelamento.
