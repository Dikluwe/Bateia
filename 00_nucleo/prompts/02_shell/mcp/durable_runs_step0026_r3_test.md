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

## Sucessor histórico O3

Preserve o digest certificado de `server.go` como identidade do predecessor no
artefato de invalidação e exija simultaneamente o digest exato do sucessor, os
dois artefatos de linhagem e os hashes inalterados do core. Os byte-goldens
independentes de benchmark e recommendation continuam sendo a prova semântica;
uma mudança futura de hash sem novo sucessor explícito deve falhar fechada.

O Passo 0031 fornece novo sucessor explícito para `server.go` e
`recommendation.go`. Exija seu artefato de supersessão e preserve os byte-goldens
de resposta; somente decode/default/schema de entrada pode mudar.
