# Passo 0076 — Separar produto e bancada

## Objetivo

Deixar as camadas numeradas somente com código e documentação do produto. Experimentos,
benchmarks e verificações históricas pertencem a `lab/` e nunca podem ser dependência
do produto.

## Mudança

- mover experimentos para `lab/arena`;
- mover benchmarks e verificações para `lab/bancada`;
- retirar código executável de `00_nucleo`;
- retirar bancadas de `04_wiring`;
- ensinar o linteiro a reconhecer `lab/` como camada isolada.

## Aceite

- nenhuma camada produtiva importa `lab/`;
- L0 permanece declarativo;
- L4 contém somente composição e ciclo de vida;
- lint, testes, análise estática e build passam sem mudar o comportamento do produto.
