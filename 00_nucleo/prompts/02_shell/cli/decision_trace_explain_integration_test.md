# Índice vinculado aos documentos originais — 0099

Fixture com duas capacidades e payloads de candidato idênticos deve preservar
IDs/ponteiros distintos; fingerprint sozinho não identifica a alegação. Fonte
compartilhada mantém todas as ocorrências e path '~/' como valor literal, não
segmento de ponteiro. Resolver cada locator no JSON original em RawMessage e
compactar sem reordenar chaves: SHA deve igualar fingerprint do nó. Verificar
CLI explain offline, cobertura de todos nós, ausência de mutação de inputs e
compatibilidade do grafo histórico com Build/Impact. Teste não autentica fatos.
