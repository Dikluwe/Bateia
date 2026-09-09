# Veredito B — 0052 r1

Testes independentes B1 passam: catálogo completo/capitalização, proveniência e IDs de família/tipo, símbolos explícitos, sete marcadores de restrição0048, fronteiras de token, não cascata, repetição, determinismo, JSON e limites Markdown. API funcional recupera Decode de descompressão sintética e Decode de XOR (controle de mecanismo errado), não Compare. Revalidação de requisitos reconstrói Original e derivados usando consulta original externa.

Dois erros no primeiro teste foram adjudicados pelo integrador e corrigidos antes do congelamento B1: exigir ordem visual da tabela embora contrato diga lista ordenada, e usar ValidateImplementationHits sem argumento de consulta externa para tentar recuperar Original adulterado. Versão inicial congelada e falhas preservadas em tests-freeze-initial.json e tests-initial.log; B1 congelado em tests-freeze.json antes execução, log tests-B1.log. Nenhuma implementação lida para corrigir os testes.

Limites: A/B lógico, filesystem compartilhado, sem atestação física e sem selo de refinamento. Busca ampla inicial por assinaturas também revelou acidentalmente nomes/assinaturas de testes anteriores, nunca corpos ou corpora. Fixtures em memória são novas. Busca positiva é somente evidência lexical, e o controle XOR demonstra falso positivo possível. Não é benchmark de precisão geral. Testes limitados de apresentação não atestam toda possibilidade de escape Markdown.

# Veredito B — adendo r2

Suíte r2 passou integralmente na primeira execução da candidata r2: 18 casos de ação/contexto/famílias/exceção de símbolo, quatro casos de separação entre documentação e código e multilinha/requisitos. Catálogo/proveniência/negativas/Markdown B1 preservados também passaram. Congelamento tests-freeze-r2.json precedeu execução; log tests-r2.log.

Supersessão semântica explícita, adjudicada pelo integrador: adendo r2 exige contexto solicitado para mechanism de compressão/descompressão, contrariando a admissão irrestrita de Decode sem contexto no teste funcional B1. B1 integral preservado em query_expansion_B1.go.txt, e teste TestExpansion0052FunctionalRetrievalAndWrongMechanism retirado da suíte ativa. Demais esperados B1 intactos. A revisão não foi chamada de mecânica nem induzida por falha de teste r2: responde à obrigação alterada após falso positivo exploratório público do integrador. Versão inicial anterior à correção B1 também preservada em query_expansion_initial.go.txt.

Nenhum corpo da implementação foi lido. Antes da execução r2, integrador comunicou resultados exploratórios públicos (UUID0, Snappy positivo); casos e congelamento já estavam escritos. Continua sem atestação física ou alegação de precisão geral. O filtro exige coocorrência lexical, e mesmo Decode+block não demonstra mecanismo correto. Uma implementação semanticamente errada com essas palavras pode passar. O teste de multilinha controla MatchedTerms e projeção de requisito, não cada componente numérico interno do score.
