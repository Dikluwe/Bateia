# Testes independentes CLI — passo 0063

Autor B, contexto novo, sem atestação física; fontes autorizadas: contrato congelado 0063, AGENTS e assinaturas públicas fornecidas pelo integrador. Não ler produção candidata ou testes existentes.

Usar adaptador fake da interface domain.EffectVerifier. Verificar que `effects verify --expect-sha256` transmite hash integral normalizado para minúsculas. Argumentos vazios explícitos, comprimentos errados, caracteres não hexadecimais e Unicode devem ser rejeitados antes de chamar adaptador. Pedidos sem flag preservam hash vazio. Conferir ajuda sobre arquivo completo e mismatch unknown sem execução, sem importar infraestrutura em L2.

Escrever este prompt antes do código. Não executar testes durante autoria; gofmt permitido; congelar arquivos para o integrador.
