# Verificação independente da configuração durável local

Derive os testes somente do contrato do passo 0077. Cubra raiz default e explícita,
criação e persistência da chave, modos, limites, symlinks, hardlinks, arquivos
especiais, tamanhos inválidos, normalização de modo e chamadas concorrentes.

Casos indisponíveis no filesystem devem ser `Skip` explícito e tornam a cobertura do
observável inconclusiva; não adapte expectativas à implementação candidata.
