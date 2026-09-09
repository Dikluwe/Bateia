# Testes independentes — contrato 0054

Autor B, sem leitura de produção, patch candidato ou benchmark 0053. Materializar testes HTTP externos pelas portas SearchCandidatesWithTrace e Discover, verificando consulta integral, termos tardios, negação e símbolos, conversão de linguagem, ausência de lang vazio, escopo e discriminador genérico. Uma chamada em sucesso, vazio válido e falha; nenhuma em entrada vazia ou contexto cancelado. Verificar rastros e OriginalQuery, contagens e erros sem token. Exercitar HTTP 500, erro GraphQL, JSON inválido, deadline e cancelamento durante requisição. Não inferir relevância de candidatos nem recall.

Supersessão explícita: testes legados que esperavam três relaxamentos, discriminador isolado, duas tentativas em timeout ou sucesso após resposta vazia passam a exigir uma única tentativa. Preservar as demais verificações. Baseline SHA256 client_test.go: eb52ebfdfd7e6bb5613791fdb4d79e4fa13320122558d22bf44630edac33d03d; observable_discovery_test.go: 744cd9fce987bf498f7e69ba417974ccf270207194318fc1f36763d363898cd0.

Congelar hashes antes de executar. Servidor local prova preservação e tratamento de respostas, sem provar qualidade da busca remota, invalidação de cache ou isolamento físico.
