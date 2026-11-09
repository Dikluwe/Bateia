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

## 0161 — Transição, tentativas e chamadas

Aplicar EVD-10–12 do [contrato ativo](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#recuperação-ativa-histórico-e-custo--0161), conforme [passo0161](../../../passos/0161-corrigir-transicoes-historico-e-revalidacao.md). Retirar da avaliação ativa o registro sem pin somente após a própria tentativa resolver para SHA completo, sem migrar conteúdo/falhas históricos ou fundir revisões fixas. AttemptedOriginalRecoveryTargets filtra os não tentados, mantendo ordem/RuleIDs dos demais por delegação a OriginalRecoveryTargets; a projeção legada e seus consumidores permanecem intactos; resultados SKIPPED_BUDGET/UNSUPPORTED não bloqueiam retomada. Projetar histórico de rodadas a partir do plano e da última ocorrência nos bundles, conservando identidade/padrões/RuleID do target. Plano sem testemunho não fabrica execução; item legado pode testemunhar na ausência de resultado. Contador de ContentProvider mede cada chamada real, inclusive preparação/falhas, sem I/O concreto. Preservar registros/hashes e restantes limites0160.

## 0162 — Omissões por quota

Aplicar EVD-13 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#quotas-de-caminhos-e-completude-técnica--0162): paths derivados congelados continuam sujeitos às quotas do coletor. Usar resultados observados para TargetsExecuted, excluindo targets omitidos por quantidade/profundidade e preservando identidade do plano. Manter saldo de bytes, métricas de portas e histórico0160/0161. Origem: [passo0162](../../../passos/0162-aplicar-quotas-e-coerencia-de-fontes.md).

## 0163 — Fragmento e retomada

Aplicar EVD-16 do [contrato](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#completude-geral-saldo-e-unicidade-de-lotes--0163). Histórico ativo não bloqueia target SKIPPED_BUDGET, inclusive fragmento. TargetsExecuted registra tentativa com Bytes>0 mesmo quando o fragmento não conclui a admissão; omissão integral não conta. Preservar contagem de portas e consumidores anteriores. Origem ex-ante: [passo0163](../../../passos/0163-fechar-completude-bytes-e-admissao-de-lotes.md).
