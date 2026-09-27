# 0102 — integração CLI da revisão com contexto

Testar comando review-context stdout/disco equivalentes, ledger inalterado e
localizadores históricos, entrada preservada e temporários limpos. Recusar
documentos adulterados, hash ausente, schema, JSON desconhecido/repetido/casefold,
trailing e segredo sintético configurado antes de emitir stdout/substituir saída.
Proteger aliases/hardlinks/symlinks e FIFO com guardas trace existentes. Sem rede.
