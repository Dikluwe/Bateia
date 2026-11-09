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

## Quotas de caminhos e completude técnica — 0162

O [passo0162](../passos/0162-aplicar-quotas-e-coerencia-de-fontes.md) aplica quantidade/profundidade aos targets explícitos, derivados congelados e amostra, antes de solicitar arquivos. Documents compartilha documentação/licença; Tests compartilha testes/exemplos; Manifests é próprio e demais arquivos usam CodeFiles. Contar paths únicos, mesmo com conteúdo igual. Zero desabilita a categoria; profundidade raiz é1. Omitidos recebem SKIPPED_BUDGET e Complete=false, sem readmissão por árvore/reserva/extras. Extras do modo geral respeitam as quotas; extras fora da seleção congelada continuam ErrIntegrity. Lista vazia não inicia extração. Métricas de execução usam os resultados de tentativa, excluindo omissões.

Na releitura técnica, tamanhos positivos conhecidos da árvore/extração devem coincidir com len(Content). Zero preserva compatibilidade desconhecida. Tamanho negativo ou declarações positivas conflitantes tornam o bundle incompleto; não admitir o arquivo divergente nem promover licença/SBOM a VERIFIED. Não exigir tamanho conhecido onde não foi informado. Preservar refs, hashes históricos e os demais limites.

EVD-13 (quotas de paths/recuperação) e EVD-14 (tamanhos técnicos) são consumidos pelas [regressões de quotas e tamanhos](../../01_core/usecase/evidence_limits_profile_regression_test.go). O perfil textual das fontes usa o [contrato de exploração](exploracao-de-vizinhancas-v1.md#perfil-textual-compartilhado--0162).

## Completude geral, saldo e unicidade de lotes — 0163

O [passo0163](../passos/0163-fechar-completude-bytes-e-admissao-de-lotes.md) estende a coerência de tamanhos ao coletor geral. Conferir bytes brutos contra tamanho positivo da extração/árvore antes do corte por orçamento. Zero é desconhecido; negativo ou conflito positivo invalida o arquivo. Conservar conflitos da árvore antes de canonicalizá-la. Omitir arquivo divergente com FAILED/falha e Complete=false, sem licença PROVEN ou cobertura examinada. O predicado de comparação é compartilhado com a releitura técnica.

Arquivo devolvido mas limitado por saldo gera SKIPPED_BUDGET, nunca NOT_FOUND: Bytes0 na omissão integral e Bytes admitidos no truncamento. Fragmento pode permanecer como dado incompleto com seu hash, sem promoção de força. Ambos mantêm cobertura incompleta e causa de orçamento; o histórico ativo admite nova tentativa. TargetsExecuted inclui fragmento admitido com Bytes>0 e exclui omissão integral. NOT_FOUND é reservado ao arquivo que realmente faltou na resposta.

A coleta coordenada recusa origens planejadas repetidas antes de qualquer porta, com ErrInvalidRunInput e relatório vazio. Chave usa owner/repo sem distinção de caixa e commit; vazio/HEAD são a mesma origem não resolvida. Não unir implicitamente planos ou quotas. Pins distintos permanecem válidos; colisão descoberta só ao resolver HEAD pode ser rejeitada depois dos efeitos necessários para conhecê-la.

EVD-15 (tamanho na coleta geral), EVD-16 (saldo/diagnóstico/histórico) e EVD-17 (unicidade prévia) são consumidos pelas [regressões de completude e lotes](../../01_core/usecase/evidence_completeness_batch_regression_test.go). Registros e hashes já emitidos permanecem intactos.

## Prefixos UTF-8 e replay — 0164

O [passo0164](../passos/0164-preservar-utf8-no-truncamento-e-replay.md) exige UTF-8 válido no conteúdo textual bruto admitido. Arquivo recebido inválido gera FAILED/falha de codificação, sem substituir bytes ou criar item. A comparação de tamanhos brutos permanece anterior ao corte.

Truncamento conserva o maior prefixo terminado em fronteira de caractere UTF-8 dentro do saldo. BytesRead, Target.Bytes e ContentHash descrevem os bytes efetivamente preservados. Sem espaço para o primeiro caractere, omitir com SKIPPED_BUDGET/Bytes0, sem item vazio e sem cobrar bytes descartados; o saldo restante pode atender outro arquivo. Preservar incompletude, cobertura, métricas e retomada0163. Não alterar ASCII ou arquivos completos, normalizar texto nem emitir U+FFFD para reparar bytes.

Após serializar/reabrir um bundle novo, conteúdo e hashes devem continuar coerentes e aceitos pela reavaliação, inclusive fragmentos. Registros e hashes antigos permanecem intactos. EVD-18 (prefixo íntegro/replay) é consumido pelas [regressões UTF-8](../../01_core/usecase/evidence_utf8_roundtrip_regression_test.go). A fronteira de caractere não comprova completude semântica do trecho.

## Execução agregada de targets — 0165

O [passo0165](../passos/0165-contar-execucao-agregada-de-targets.md) explicita a agregação de TargetsExecuted: cada target do plano conta no máximo uma vez quando qualquer resultado correspondente da coleta da rodada é COLLECTED/NOT_FOUND/FAILED ou SKIPPED_BUDGET com Bytes>0. Um target derivado pode produzir vários paths; a omissão de um deles não apaga a execução observada de outro. Somente omissões integrais/UNSUPPORTED não contam. Preservar o testemunho legado na ausência de resultado correspondente e não contar conteúdo original como execução da rodada.

A métrica agregada não altera a projeção de histórico 0161: a última ocorrência correspondente continua governando a retomada, inclusive quando outro path foi lido. Preservar schemas, identidade do target, bytes, chamadas reais e limites de força de evidência. Registros, hashes e testes anteriores permanecem intactos.

EVD-19 (execução agregada sem bloquear retomada) é consumido pelas [regressões de métricas derivadas](../../01_core/usecase/evidence_recovery_metrics_regression_test.go), com controles de saldo adjacentes, omissão integral e leitura completa.

## Testemunhos pertinentes e codificação técnica — 0166

O [passo0166](../passos/0166-validar-bytes-testemunhos-e-politicas.md) limita o fallback sem resultado da recuperação: @repository-metadata só testemunha target que pede metadados ou é de espécie repository; não substitui caminho pertinente de código/teste derivado. Sem resultado nem arquivo correspondente, não registrar execução nem bloquear retomada. Preservar itens legados pertinentes, resultados explícitos NOT_FOUND/FAILED, última ocorrência e agregação0165; OriginalRecoveryTargets mantém a projeção legada. EVD-20 é consumido pelas [regressões de testemunhos](../../01_core/usecase/evidence_recovery_witness_regression_test.go).

RecheckTechnicalSources exige UTF-8 válido nos bytes brutos antes de criar item/hash. Arquivo inválido gera falha e Complete=false, sem item ou IntegrityVerified/VERIFIED; não substituir bytes nem recalcular hash de conteúdo reparado. Fontes válidas, inclusive U+FFFD literal, conservam bytes/hashes e resultado de integridade após JSON. Manter tamanho bruto, limites e cancelamento. EVD-21 é consumido pelas [regressões de codificação técnica](../../01_core/usecase/technical_recheck_encoding_regression_test.go). O novo recorte não altera bundles históricos.
