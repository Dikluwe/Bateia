# Laudo 0034 — Separação entre produto e bancada

**Resultado**: aprovado

Experimentos, benchmarks e verificações históricas foram retirados das camadas do
produto e organizados em `lab/arena` e `lab/bancada`. L0 ficou declarativo, L4 deixou de
conter subdiretórios de bancada e nenhuma camada produtiva passou a importar `lab/`.

A reorganização preservou o comportamento produtivo. As exceções antigas do linteiro
foram reduzidas e a arquitetura passou a expressar melhor a diferença entre produto e
material de investigação.
