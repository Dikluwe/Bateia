# Prompt espelhado — RecommendationService R3 do Passo 0026

Teste a instância pública real de `RecommendationService`, não apenas identidade de um
fake em L4. Injete as mesmas portas offline certificadas do Passo 0025. Execute
`Recommend` e `RecommendObserved` na mesma instância, exija resultados deep-equal,
contagens exatamente duplicadas nas mesmas portas e progresso controlado não vazio.

Para cancelamento entre stages reais, use uma porta discovery que cancela o contexto
dentro da primeira chamada e retorna `context.Canceled`. Exija exatamente uma chamada
discovery e zero chamadas posteriores de content. Não use o ramo impossível que testa
`ctx.Err()==nil` depois de já aguardar `ctx.Done()`.
