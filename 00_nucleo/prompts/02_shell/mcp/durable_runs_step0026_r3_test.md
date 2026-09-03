# Prompt espelhado — compatibilidade legada MCP R3 do Passo 0026

Congele a compatibilidade 0025 por testemunho independente, sem chamar o serializer
candidato para construir o esperado. Verifique os hashes certificados dos arquivos
imutáveis do predecessor e invoque de verdade `bateia_benchmark` e `bateia_recommend`
com fakes offline. Não invente sucesso ausente: a implementação 0025 certificada retorna
`Invalid params: depth deve estar entre 1 e 5` para o input de benchmark declarado pelo
schema, portanto congele os bytes realmente observados desse erro legado e registre a
incompatibilidade causal. Compare o texto da recomendação síncrona contra bytes JSON
literais produzidos antes do candidato 0026.

Decodifique sempre o item `content[0].text` antes de verificar o JSON interno. O teste
não pode derivar o golden de `MarshalRecommendationResult`, nem aceitar campos durable
nas respostas legadas.
