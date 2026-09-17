# Passo 0072 — gate adversarial segregado

Autoridade: adversário/oráculo, posterior ao contrato e anterior à implementação
candidata. Derive este gate exclusivamente do passo 0072, de seus testes contratuais e
do baseline público congelado no commit `03fd5f1`. É proibido ler, inferir ou adaptar
expectativas à implementação candidata.

Materialize um gate de caixa-preta sobre as APIs públicas que elimine M01–M14. Cubra
casos positivos, negativos e opacos; execute vinte permutações de árvores, consultas,
termos, catálogo e requisitos; compare observações canônicas byte a byte; e prove que
entradas permanecem imutáveis.

Observáveis obrigatórios:

- consultas existem somente para regras declaradas e cobertas;
- `span` e `bits/basic_string.h` recebem reserva justa antes de distratores;
- `spanner`, extensões documentais e arquivos extensionless irrelevantes não casam;
- nome de caminho nunca constitui prova de capacidade;
- quota ou bytes insuficientes preservam `COLLECTION_INCOMPLETE`, confiança zero e
  nenhum gap atribuído ao candidato;
- caminhos candidatos, examinados e pulados permanecem explícitos e ordenados;
- um arquivo compartilhado é lido/cobrado uma vez, mas cobre todos os requisitos;
- falha de árvore preserva `TreeComplete=false` e Unknown;
- `FetchTree` local não devolve sucesso truncado em 20.000 entradas;
- árvore Go pequena conserva a coleta legada.

`Unknown` somente satisfaz os casos opacos declarados. Um mutante válido sobrevivente
impede o selo. O gate não corrige a solução e não escreve fora de
`04_wiring/step0072_gate/` e deste prompt proprietário.
