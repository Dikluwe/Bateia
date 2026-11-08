# Coleta e serialização de evidências v1

## Origem e autoridade

Contrato ex-ante do [passo0158](../passos/0158-validar-origem-e-estabilizar-bundles.md), sob [ADR0022](../ADR/ADR-0022-integridade-de-coleta-e-serializacao-estavel.md) e regras Tekt. O sufixo v1 identifica este contrato, sem introduzir outro schema público. L1 domínio possui serialização; L1 casos de uso possuem admissão e vínculo de origem; os efeitos concretos permanecem em L3.

## Admissão e pin

A coleta coordenada valida todos os planos antes de executar qualquer porta. Exige schema1, identidade de repositório válida e igual à do candidato, quotas não negativas e orçamento de bytes positivo. A API genérica compartilha essas regras e conserva somente a exceção de schema ausente para consumidores legados; versões desconhecidas são rejeitadas. Revisões abreviadas legadas da API genérica não são promovidas a SHA validado pela jornada de insights.

Plano com commit fixado governa o commit das leituras. Metadados de branch mais recente não substituem esse pin. HEAD ou vazio usam a resolução explícita existente. Owner/repo divergente nos metadados, resposta de extração nil ou Ref diferente da identidade e revisão efetivas retornam ErrIntegrity antes de admitir conteúdo. A coordenação propaga esse erro sem bundles/cartões; não o reclassifica como falha operacional de resolução. Falhas operacionais de resolução continuam explícitas, sem promover histórico ou revisão opaca a fonte lida. Cancelamento observado após uma porta prevalece antes da próxima.

## Serialização e histórico

Novas serializações de EvidenceBundle emitem requirement_coverage como array não nulo, inclusive quando nil ou vazio. O marshaler copia o estado necessário e não altera slices do chamador. Após roundtrips, bytes e identidade causal devem estabilizar; null legado continua legível e é projetado como array em nova serialização.

Registros e hashes históricos não são reescritos. Formas vazias antes alternantes convergem nas novas projeções, que podem diferir do ID antigo instável. Não substituir IDs já emitidos nem prometer igualdade com os dois resultados defeituosos. Cobertura populada, já estável, conserva bytes e IDs sob controles congelados antes do patch. Reprodução histórica continua no commit registrado pelo gate; hashes de arquivos preservados são conferidos separadamente da nova execução.

## Aceitação

| ID | Observável | Consumidor |
| --- | --- | --- |
| EVD-01 | Owner/repo/commit divergentes ou resposta nil não geram conteúdo; falha é de integridade | [Regressões de coleta](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-02 | Plano inválido em qualquer posição impede todas as portas; schema ausente é exclusivo da API legada | [Admissão do lote](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-03 | Pin fixado permanece; HEAD resolve explicitamente; falha operacional continua visível | [Metadados e compatibilidade](../../01_core/usecase/evidence_integrity_regression_test.go) |
| EVD-04 | Nil/vazio/null legado estabilizam como array; exemplo populado conserva hash literal | [Regressões de domínio](../../01_core/domain/evidence_integrity_regression_test.go) |
| EVD-05 | Formulação salva/reaberta mantém IDs; controle populado conserva ID literal anterior | [Regressões de formulação](../../01_core/usecase/evidence_integrity_regression_test.go) |

Esses controles verificam integridade declarada e estabilidade local, sem autenticar provedor remoto, comprovar comportamento, novidade ou efeito no alvo.

## Recorte congelado de recuperação — 0159

Na recuperação, os caminhos explícitos ou derivados dos padrões dos targets são congelados antes da coleta. O modo interno restrito reutiliza a admissão de origem deste contrato, mas não amplia caminhos pela árvore ou pelo buscador de implementação. Resposta com caminho extra/duplicado é ErrIntegrity e não gera resultado útil. Arquivos elegíveis conhecidos fora do recorte conservam incompletude; enumeração/metadados não significam leitura desses arquivos. O modo público de amostragem geral conserva seu comportamento. O [passo0159](../passos/0159-fechar-cancelamento-escopo-e-recibos.md) e as [regressões de fronteiras](../../01_core/usecase/investigation_boundaries_regression_test.go) delimitam essa correção sem modificar registros anteriores.

## Revisões, interrupção e orçamento — 0160

O [passo0160](../passos/0160-preservar-revisoes-e-orcamentos-de-coleta.md) exige identidade repo@commit completa no diagnóstico, preparação e união da recuperação. Bundles de revisões diferentes não são fundidos, mesmo com paths/hashes iguais. Leituras preparatórias usam o pin do target; resultados conservam a origem dos itens e não alteram as entradas.

A coleta geral observa contexto após árvore e antes/depois da busca de implementação. Cancelamento/deadline prevalecem antes de outra porta.

TotalBytes abrange os metadados: JSON que não couber é omitido integralmente com SKIPPED_BUDGET e Complete=false. Nenhum JSON truncado é admitido; BytesRead fica dentro do teto. AdditionalBytes é um saldo compartilhado entre as coletas da rodada, consumido em ordem canônica de repo@commit. Sem saldo, gerar omissão por orçamento, sem portas de coleta e sem contar esses targets como executados. Bytes de bundles originais não são cobrados novamente. Falta de orçamento conserva inconclusão.

Critérios EVD-06 (revisões e isolamento), EVD-07 (interrupção entre portas), EVD-08 (saldo por rodada) e EVD-09 (metadados dentro do teto) são consumidos pelas [regressões0160](../../01_core/usecase/evidence_budget_revision_regression_test.go). Schemas e registros históricos permanecem intactos; somente novas execuções seguem a correção.

## Recuperação ativa, histórico e custo — 0161

O [passo0161](../passos/0161-corrigir-transicoes-historico-e-revalidacao.md) distingue o dossiê preservado dos bundles ativos. Quando a própria tentativa de uma origem vazia/HEAD resolve explicitamente para SHA completo, retirar a origem sem pin da avaliação da nova rodada, sem migrar seus itens/falhas para o SHA. Falha, omissão e leitura de outra revisão fixa não autorizam retirada. O original continua no dossiê e nos objetos de entrada, intacto; revisões fixas distintas continuam separadas.

Histórico de execução inclui COLLECTED/NOT_FOUND/FAILED, exclui SKIPPED_BUDGET/UNSUPPORTED e não presume execução de plano sem resultado. Item legado com path correspondente é testemunho na ausência de resultado. A projeção ativa AttemptedOriginalRecoveryTargets mantém ordem e RuleIDs dos elegíveis, delegando a OriginalRecoveryTargets; esta conserva sua projeção legada para consumidores históricos e não governa a elegibilidade ativa. Rodadas anteriores usam a última ocorrência correspondente nos resultados para selecionar os targets do plano; a CLI lê os bundles e delega essa regra a L1. Não reescrever planos antigos para corrigir a projeção nova.

ProviderCalls contabiliza cada chamada real a ContentProvider, inclusive falhas e preparação de padrões; grupo omitido não chama porta nem aumenta contador. A métrica não mede requisições internas do adaptador ou custo financeiro.

Critérios EVD-10 (transição de origem), EVD-11 (histórico de tentativas) e EVD-12 (chamadas reais) são consumidos pelas [regressões de recuperação](../../01_core/usecase/evidence_recovery_followup_regression_test.go) e [retomada CLI](../../02_shell/cli/idea_recover_history_regression_test.go). A integridade de declarações na busca possui [contrato próprio](revalidacao-de-declaracoes-v1.md).
