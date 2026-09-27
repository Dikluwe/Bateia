# Contrato JSON de comparação de snapshots

Verificar roundtrip dos tipos do domínio comparison-v1, com decisão histórica
distinta da candidata e deltas de adição/remoção/alteração. Ponteiros before/after
omitidos devem continuar ausentes; DecisionImpact permanece NONE. Apenas
serialização tipada em memória, sem executar decisões ou consultar fontes.
