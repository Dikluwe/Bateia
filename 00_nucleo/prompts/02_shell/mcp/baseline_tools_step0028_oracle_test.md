# Prompt espelhado — captura basal MCP do Passo 0028

Construa o servidor MCP basal diretamente pelo wiring público anterior ao Passo 0028, execute `tools/list` em processo e serialize a captura de modo canônico. Exija exatamente as nove tools congeladas, na ordem pública histórica, e o SHA-256 `ee630cff5eada0dee4afe82ad9513dea7b760f0beb7905e915d2714a9114a596`.

O teste deve falhar se aceitar captura fornecida externamente sem vínculo por hash. Não importe nem leia implementação candidata da décima tool.

No sucessor explícito do Passo 0031, leia a captura versionada em
`00_nucleo/benchmark/0031/basal-tools-successor.json`. Preserve os nove nomes e
a ordem; a única diferença autorizada no digest é o contrato de defaults das
tools de recomendação síncrona e durável.
