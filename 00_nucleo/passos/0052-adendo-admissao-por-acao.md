# Adendo 0052 r2 — Admissão lexical por ação

A candidata inicial recuperou Decode no Snappy, mas trouxe 12 hits no UUID para a mesma pergunta de descompressão. Resultado exploratório público em benchmark/0052/cli-comparison.json (preservar como r1). Não promover a expansão irrestrita. Este adendo complementa, não apaga, o contrato inicial.

## Obrigação adicional

Para cada linha que recebeu expansões, reconhecer como famílias de ação: compression, decompression, serialization, deserialization, sorting, comparison (pelos IDs versionados). Se houver famílias de ação, um candidato a essa linha só pontua quando satisfaz TODAS as famílias ativas no campo examinado. A presença de data/dst/block sozinha não satisfaz uma ação.

Cada família pode ser satisfeita por um termo translation correspondente presente, ou por um termo mechanism correspondente presente. Exceção obrigatória: mechanism das famílias compression/decompression exige também ao menos um contexto da consulta presente no MESMO campo examinado: block, blocks, corrupt ou corrupted. Os termos de contexto devem estar nos Terms da linha (não podem ser inventados pelo projeto). Sem contexto, essas alternativas técnicas ambíguas não são admitidas. Isso é filtro lexical conservador, não equivalência funcional.

Campo examinado para hit de código: união dos termos já existentes de nome/corpo. Campo para hit documental: palavras do grupo Doc. Não misturar documentação com corpo para satisfazer o filtro de código. Linhas sem expansão de ação preservam comportamento anterior. Identificador exato preserva prioridade/recuperação. Em consultas multilinha, considerar termos e pontuação somente das linhas admissíveis para aquele candidato; falha de uma linha não impede outra linha admissível. A pontuação por requisito segue a mesma admissão. Hits e Reason devem explicitar a aplicação do filtro quando relevante. Não alterar IDF nem transformar negação em operação avaliada.

B acrescenta testes congelados para: termos secundários sem ação rejeitados; decode sem contexto rejeitado em pergunta de descompressão; decode+block ou decode+corrupt aceito quando contexto solicitado; contexto existente só no projeto não basta; literal translation satisfaz ação; documentação/code não se misturam; multilinha não reaproveita termos da linha rejeitada. Conservar testes iniciais, registrando revisões mecânicas de observáveis separadamente.

Integrador repete mesmos casos sem trocar perguntas: Snappy deve recuperar Decode; UUID na mesma pergunta não deve produzir hit de compressão/descompressão. Se falhar ou se uma segunda revisão não melhorar o mesmo gargalo, interromper promoção e revisar obrigação, sem maquiar resultados.
