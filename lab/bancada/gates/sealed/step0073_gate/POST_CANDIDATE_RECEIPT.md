# Recibo pós-candidato — gate 0073

- Data: 2026-09-16.
- Papel: adversário/oráculo `cpp_locator_adversary`.
- Regime: segregação procedural no mesmo filesystem, sem alegação de isolamento físico.
- Runner selado: `04_wiring/step0073_gate/run.sh`.
- Entradas congeladas: todas verificadas por SHA-256.
- Artefatos do oráculo: todos verificados por SHA-256.
- Baseline histórico: blobs públicos de domínio e caso de uso no commit `03fd5f1`
  verificados pelo runner.

## Resultado

O runner passou uma vez e o gate foi repetido com `-count=2`, também com sucesso.
`git diff --check` no escopo do oráculo passou. As testemunhas da matriz eliminam
M01–M14; mutation score observado: **14/14 = 1.0**.

R1 está registrada no recibo pré-candidato e removeu somente expectativas que excediam
os contratos 0071/0073. Não houve leitura nem ajuste da implementação candidata por este
papel.

