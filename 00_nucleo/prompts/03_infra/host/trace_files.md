# Proteção física dos artefatos trace

**Camada:** L3
**Proprietário:** `03_infra/host/trace_files.go`
**Contrato:** 0129. Extrair validação de saída/aliases, resolução de caminhos protegidos e publicação por raiz retida. Preservar SameFile, recusa de symlinks e folhas ausentes ambíguas, nonce aleatório, temporário exclusivo, Sync e rename pela raiz. Receber bytes já serializados; JSON/stdout/segredos permanecem na adaptação L2. Não reabrir caminho absoluto para publicar após validar uma raiz.
