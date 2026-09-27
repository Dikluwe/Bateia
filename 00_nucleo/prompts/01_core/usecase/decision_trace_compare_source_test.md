# 0100 — fontes tipadas na comparação

Fonte com EvidenceSource completa deve exigir fingerprint consistente mesmo
quando graph_hash é recalculado. Depois de selar corretamente o novo payload,
comparar deve registrar um único delta completo e iniciar revisão sem alterar a
fonte original nem promover a decisão. Usar referências e textos sintéticos.
