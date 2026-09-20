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

## Passo 0040

Queries são rastreabilidade da busca, não conteúdo do candidato. Cobertura preliminar e não objetivos usam descrição e snippets retornados; nunca os termos consultados. A cobertura permanece sinal lexical preliminar, não comprovação.

## Autorização para descarte

`DROP` exige `CalibrationStatus == "calibrated"`, parâmetros de calibração presentes e
ambos finitos. O estado explícito `uncalibrated + nil` continua usando o ranking para
ordenar e pode emitir `ADVANCE`, mas converte qualquer `DROP` proposto em `AUDIT`, com a
razão pública `modelo sem calibração válida não autoriza DROP`.

Qualquer combinação opaca ou contraditória de status e parâmetros produz `AUDIT`,
`Probability == nil` e a razão pública
`estado de calibração opaco ou malformado exige AUDIT`. Não transformar opacidade em erro
ou sucesso. Com calibração válida, preservar integralmente os caminhos atuais de decisão,
probabilidade e auditoria. Não alterar sinais, pesos, priors, thresholds, seed/taxa,
calibrador, schema/API, descoberta, evidência ou decisão final de reúso.
