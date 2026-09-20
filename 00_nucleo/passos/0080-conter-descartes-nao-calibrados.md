# Passo 0080 — Conter descartes não calibrados

## Problema

O ranking ajuda a ordenar candidatos, mas um score sem calibração válida não representa
probabilidade e não deve autorizar descarte.

## Regra

- somente `calibrated`, com parâmetros presentes e finitos, pode produzir probabilidade
  e `DROP`;
- `uncalibrated` sem parâmetros pode produzir `ADVANCE`, mas converte qualquer descarte
  em `AUDIT`;
- combinações contraditórias, desconhecidas ou não finitas sempre resultam em `AUDIT` e
  probabilidade ausente;
- ranking, sinais e ordenação não mudam por causa dessa contenção.

## Aceite

Testes normais de `probabilistic_screening` cobrem os três estados, o conflito com
`NonGoals` e a ausência de autoridade de parâmetros não finitos.
