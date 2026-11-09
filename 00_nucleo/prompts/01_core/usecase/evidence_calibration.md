# Evidence calibration use case

Calcule métricas determinísticas por split e agregadas. `PROVEN` só é acerto
quando coincide com o oráculo e cada fonte existe no bundle do mesmo commit.
Promoção de `INCONCLUSIVE`, fonte inválida, perda sem diagnóstico, precisão menor
que 90% ou recall menor que 85% impedem promoção. Falha específica no split
`TEST` produz `FAIL_GENERALIZATION`. Reordenação não altera o relatório.

## Acerto condicionado à fonte — 0170

TrueProven exige Expected=PROVEN, Observed=PROVEN e aprovação de validImmutableCalibrationSources. Fonte ausente ou commit divergente não é acerto; preservar denominadores ExpectedProven/ObservedProven e InvalidSources. Calcular taxas coerentes no split e agregado sem alterar precedência de veredictos, promoção opaca ou diagnósticos. Controles válidos conservam relatórios e hashes; métricas corrigidas valem para novas avaliações, sem recalcular resultados anteriores.
