# Prompt Spec L0: Peneira Probabilística Explicável

## Objetivo

Extrair sinais causais de snapshots, calcular ranking reconstruível, calibrar fora
da amostra e decidir `ADVANCE`, `AUDIT` ou `DROP` sem confundir ausência com
evidência negativa.

## Restrições

- L1 puro, determinístico e sem I/O.
- Todo efeito do score é serializado como contribuição.
- Probabilidade somente com calibrador que supere prevalência fora da amostra.
- Validação agrupa por ideia; `uncertain` não vira negativo.
- Conflito com `NonGoals` permanece separado do total.
- Auditoria de descartes é determinística por seed.

## Verificação

- Soma das contribuições reproduz score pré-clamp.
- Snapshot v1 migra sem fabricar evidência.
- Ausência leva a auditoria; conflito pode bloquear avanço.
- Folds nunca misturam a mesma ideia.
