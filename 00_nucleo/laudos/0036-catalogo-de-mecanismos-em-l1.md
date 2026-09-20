# Laudo 0036 — Catálogo de mecanismos em L1

**Resultado**: aprovado

O catálogo ativo de mecanismos e sua API de leitura passaram para
`01_core/usecase`. O arquivo é incorporado em compilação, cada chamada retorna uma cópia
independente e L4 deixou de possuir os bytes ou interpretá-los.

O conteúdo e os consumidores semânticos foram preservados. A mudança reduz autoridade
duplicada sem introduzir I/O em L1.
