# Canal documental de declarações — 0051

Origem: contrato congelado 00_nucleo/passos/0051-recuperar-documentacao-de-declaracoes.md. Auxiliar puro recebe FuncDecl íntegra e Doc AST, termos interpretados e identidade do arquivo. Tokeniza Doc.Text com implementationTerms; interseção distinta ordenada e score até 99, sem pesos IDF. Retorna hit documental somente com correspondência e grupo anexado. Anexação ao relatório ocorre após código/chamadas, ordenada score decrescente e path/line/symbol crescente, deduplicada pela identidade da declaração; omissão marca Limited e issue. Nunca apresenta documentação como comportamento comprovado.

## Passo 0059 — orçamento documental
Antes da expansão de chamadas, ordenar candidatos por score desc/path/line/symbol,
deduplicar por path/line/symbol e excluir identidades de todos os hits lexicais
admitidos, inclusive os que serão cortados. Usar esse conjunto também no preenchimento
final, impedindo reintrodução documental de código omitido. Reserva nominal
max(1, limit/4), limitada aos candidatos únicos novos e às vagas após declarações
exatas (score>=100000); se houver código, deixar ao menos uma vaga para ele.
Reserva zero sem novos candidatos. Sem código, preenchimento documental pode usar
todo limite. Preservar campos e provenance, registro de omissões e ordem documental.
