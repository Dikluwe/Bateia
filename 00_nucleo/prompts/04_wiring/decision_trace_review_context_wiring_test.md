# 0102 — despacho offline real da revisão contextual

Invocar main em subprocesso com idea trace review-context. Pedido une review_input
e build_input originais, com recibo IN_PROGRESS contendo marcador sintético que
também aparece em GITHUB_TOKEN no .env do diretório. Remover tokens do ambiente
herdado. O comando deve emitir contexto JSON contendo nota, hashes e ledger
consultivo, provando despacho sem carregar .env. Sem rede/candidato executado,
deadline de dez segundos.
