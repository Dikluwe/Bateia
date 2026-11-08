# Regressões de revisões e orçamento de coleta

L1, consumidor evidence_budget_revision_regression_test.go. Origem ex-ante: [passo0160](../../../passos/0160-preservar-revisoes-e-orcamentos-de-coleta.md); critérios EVD-06–09 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#revisões-interrupção-e-orçamento--0160).

Usar APIs públicas e portas contadas em memória. Planejar dois findings UNKNOWN_LICENSE para commits diferentes do mesmo repositório: ambos são lidos no próprio pin e permanecem em bundles separados, mesmo com conteúdo/hash iguais. Incluir bundles originais, verificar JSON de entrada intacto e diagnóstico sem transferência de falhas entre revisões. Árvore preparatória deve receber o pin do target.

Cancelar e expirar contexto durante árvore ou busca opcional: retornar erro de contexto, sem próxima porta. Controle sem interrupção executa busca e extração.

Com política pequena e dois/três repositórios, somar BytesRead novos, exigir teto por rodada e controle positivo de ambos lidos quando cabem. Quando o primeiro esgota o saldo, os demais não chamam portas, são SKIPPED_BUDGET e não aumentam TargetsExecuted. Rever ordem das entradas e exigir distribuição canônica. Não cobrar bytes do original nem modificar seus registros. Orçamento esgotado não produz ABSENT.

Metadados com limite1, tamanho-1, tamanho exato e tamanho+1: bytes admitidos nunca excedem limite; JSON integral e hash válido quando cabe; item ausente e SKIPPED_BUDGET quando não cabe. Preservar testes anteriores e hashes congelados.
