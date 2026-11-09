## Passo 0094 — releitura limitada

Na releitura de licença, documento raiz adicional não regular também torna a
coleta incompleta; não omitir symlink/submodule para fabricar unicidade. Releitura
de bytes íntegros registra PROVIDER_OBSERVED, nunca assinatura/autenticação de
exportação. As guardas não introduzem execução de terceiros.

# Reler fontes pelo provedor

Não aceitar bundles fornecidos. Validar candidatos GitHub com SHA completo,
caminhos e associação ao achado. Reler somente suporte atual de predicados
implementados; excluir conflitos estruturados mesmo com redações diferentes.
Conferir árvore e outras licenças raiz, revisão retornada, completude e
orçamentos. Recalcular hashes dos bytes observados e reutilizar o verificador
determinístico. Falhas de leitura não podem promover evidência. Separar
observação do provedor de autenticação de exportações. Preservar busca simples.

## 0159 — Cancelamento, targets congelados e recibos

Aplicar o [passo0159](../../../passos/0159-fechar-cancelamento-escopo-e-recibos.md). O responsável recheck verifica contexto antes/depois das portas e da conclusão; cancelamento não publica VERIFIED. Recuperação usa coleta interna restrita a targets congelados, sem ampliar caminhos ou aceitar extras da resposta; busca geral conserva seu modo. Fontes de vizinhança exigem correspondência um-para-um dos recibos COLLECTED primários por connector/candidate, rejeitando órfãos, duplicatas, nil e estados incoerentes antes dos efeitos. Preservar schemas, hashes de registros existentes e limites de força de evidência.

## 0162 — Tamanhos coerentes

Aplicar EVD-14 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#quotas-de-caminhos-e-completude-técnica--0162). Tamanhos positivos conhecidos da árvore e da extração devem coincidir com bytes entregues; zero é desconhecido legado. Negativos, conflitos positivos e divergências impedem completude/promoção, sem admitir arquivo divergente. Conferir paths solicitados e licenças adicionais. Preservar budgets, origem observada, schemas e hashes anteriores. Origem: [passo0162](../../../passos/0162-aplicar-quotas-e-coerencia-de-fontes.md).

## 0163 — Predicado compartilhado de tamanhos

Reutilizar a mesma comparação pura de tamanho dos bytes brutos aplicada no coletor geral, preservando EVD-14 e seus testes de zero/exato/negativos/conflitos. A nova obrigação EVD-15 pertence à coleta geral; não afrouxar a releitura técnica. Origem ex-ante: [passo0163](../../../passos/0163-fechar-completude-bytes-e-admissao-de-lotes.md).
