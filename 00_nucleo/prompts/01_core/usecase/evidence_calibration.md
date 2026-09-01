# Evidence calibration use case

Calcule métricas determinísticas por split e agregadas. `PROVEN` só é acerto
quando coincide com o oráculo e cada fonte existe no bundle do mesmo commit.
Promoção de `INCONCLUSIVE`, fonte inválida, perda sem diagnóstico, precisão menor
que 90% ou recall menor que 85% impedem promoção. Falha específica no split
`TEST` produz `FAIL_GENERALIZATION`. Reordenação não altera o relatório.
