# Regressões de completude e admissão de lotes

L1. Origem ex-ante [passo0163](../../../passos/0163-fechar-completude-bytes-e-admissao-de-lotes.md); EVD-15–17 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#completude-geral-saldo-e-unicidade-de-lotes--0163).

Portas em memória contadas. Template MIT devolvido com tamanho desconhecido/exato mantém o controle; tamanho negativo/maior/menor, árvore divergente e conflito positivo não geram item de licença nem PROVEN. Comparar bytes brutos antes de corte legítimo; cobertura dirigida não vira examinada em arquivo recusado.

Teto calculado para metadados e primeiro arquivo: segundo existe na resposta mas é SKIPPED_BUDGET/Bytes0, sem mensagem de ausência, com BUDGET_EXHAUSTED e LICENSE retomável. Arquivo parcial também marca orçamento, mantém hash/bytes do fragmento e não bloqueia histórico; ausência real é NOT_FOUND. Métrica de recuperação mantém tentativa de fragmento e não conta omissão integral. Preservar testes anteriores.

Lote com duplicatas conhecidas em qualquer posição/caixa ou vazio/HEAD falha sem chamadas e sem relatório útil; um plano ou dois pins distintos funcionam. Verificar que entradas não foram alteradas. Nenhum I/O concreto ou atualização de registros/hashes históricos.
