# Passo 0074 — Alinhar regra de precondição ao vocabulário do GSL

## Regime

Calibração Tekt A/B limitada ao catálogo experimental
`relatorios/coin5-pathcopy-rules.json`. O motor, os contratos 0072/0073 e seus gates
permanecem inalterados. Não se alega equivalência geral nem isolamento físico.

## Baseline observável

No dossiê `coin5-pathcopy-dossie-v3`, a capacidade
`reject an out-of-range start position unconditionally, not only in debug builds`
recebe somente a fonte documental `README.md:46`, padrão `precondition`, confiança 0.40.
O mesmo bundle contém `include/gsl/span:609` e `include/gsl/span:747`, mas a regra
`unconditional-start-validation` não inclui o vocabulário `Expects` do GSL.

## Mudança autorizada

Adicionar à evidência forte de símbolo somente os padrões compostos:

- `Expects((size() >= Offset)`;
- `Expects(size() >= offset)`.

Não adicionar o padrão genérico `Expects(`, pois ele casaria precondições não relacionadas.
Atualizar a versão do catálogo para registrar a calibração.

## Oráculos A/B

- positivo: a avaliação do mesmo bundle cita `include/gsl/span:609` e/ou linha 747 como
  `STRONG_CONTENT`, com o padrão composto correspondente;
- negativo: ocorrências como `Expects(count <= size())` e `Expects(idx < size())` não são
  atribuídas a essa capacidade por estes novos padrões;
- preservação: a fonte documental pode permanecer, mas não deve ser a única fonte;
- limite: a correspondência lexical não comprova execução nem transforma o resultado em
  `PROVEN`.

