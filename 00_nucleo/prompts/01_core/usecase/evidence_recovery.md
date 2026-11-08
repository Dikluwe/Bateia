# Evidence recovery engine

Diagnostique INCONCLUSIVE offline, mantendo causas e gates separados. Planeje
targets determinísticos ligados a findings, aplique política e histórico para
evitar repetição, colete apenas targets congelados, una conteúdo por hash e
reexecute os mesmos assessment e decisor. Falha, budget e revisão opaca continuam
INCONCLUSIVE; nenhuma regra de evidência é criada ou afrouxada.

0131: OriginalRecoveryTargets projeta targets de planos e resultados na ordem recebida, mantendo primeira capability e RuleID original/GATE-LICENSE-001/GATE-OPERATION-001 conforme a regra antes na CLI. Sem efeitos ou política duplicada em L2.

## 0159 — Cancelamento, targets congelados e recibos

Aplicar o [passo0159](../../../passos/0159-fechar-cancelamento-escopo-e-recibos.md). O responsável recheck verifica contexto antes/depois das portas e da conclusão; cancelamento não publica VERIFIED. Recuperação usa coleta interna restrita a targets congelados, sem ampliar caminhos ou aceitar extras da resposta; busca geral conserva seu modo. Fontes de vizinhança exigem correspondência um-para-um dos recibos COLLECTED primários por connector/candidate, rejeitando órfãos, duplicatas, nil e estados incoerentes antes dos efeitos. Preservar schemas, hashes de registros existentes e limites de força de evidência.

## 0160 — Revisões distintas e saldo por rodada

Aplicar EVD-06/08 do [contrato de coleta](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#revisões-interrupção-e-orçamento--0160). Diagnóstico, planos e união usam repo@commit; paths/hashes iguais em commits distintos não colapsam. Fixar o candidato da árvore preparatória ao target; união não altera entradas nem mistura revisões. AdditionalBytes é saldo global da rodada, incluindo metadados, distribuído por chave canônica. Saldo zero gera targets SKIPPED_BUDGET sem portas de coleta, sem contá-los como executados. Metadados que não cabem seguem a omissão integral do coletor. Não cobrar histórico, afrouxar evidência ou recalcular hashes antigos. Origem ex-ante: [passo0160](../../../passos/0160-preservar-revisoes-e-orcamentos-de-coleta.md).
