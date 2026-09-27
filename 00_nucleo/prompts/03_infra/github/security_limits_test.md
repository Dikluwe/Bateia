# Prompt: limites por endpoint GitHub

Camada L3. Gerar `03_infra/github/security_limits_test.go`.

Verificar que Content-Length acima do teto bloqueia árvore, metadata, manifest root, raw, catálogo e superfícies sem decodificação nem fallback. Cobrir bytes chunked, classificação explícita de symlink/submodule, caminhos de árvore inválidos, raw acima de 2 MiB terminal, e HEAD preservado como mutável no fluxo simples (sem inventar SHA fixado).

Usar exclusivamente transport HTTP sintético, sem chamadas externas. Limites não garantem ausência de todas falhas de rede, nem assinatura criptográfica do relatório exportado.
