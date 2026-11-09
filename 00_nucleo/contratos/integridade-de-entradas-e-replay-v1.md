# Integridade de entradas e replay v1

Intenção ex-ante do [passo0166](../passos/0166-validar-bytes-testemunhos-e-politicas.md). Responsáveis L1: formulação de algoritmos e decisão de reúso; as correções de recuperação/releitura técnica estão no [contrato de coleta](coleta-e-serializacao-de-evidencias-v1.md#testemunhos-pertinentes-e-codificação-técnica--0166). Preservar formatos públicos, interpretações declaradas e registros/hashes históricos.

## ALG-01 — Validar os bytes originais

FormulateAlgorithmInsights aplica suas validações aos textos e conteúdos originais, antes de qualquer cópia capaz de substituir UTF-8 inválido. Serialização limitada pode medir o pedido, mas não reparar os bytes que serão validados. Conteúdo inválido em qualquer item de fonte/alvo gera ErrIntegrity; campos textuais inválidos continuam rejeitados pelas guardas existentes. Hash de conteúdo substituído não autentica o original inválido. U+FFFD originalmente válido permanece permitido; não proibir esse caractere nem normalizar conteúdo de evidência.

Copiar os slices que serão normalizados/ordenados sem alterar strings, caller ou registros anteriores. Preservar quotas, contexto, papéis, referências, canonicalização e IDs válidos. A saída mantém DECLARED_NOT_REASSESSED/UNVERIFIED; validar vínculo não comprova interpretação. Consumidor: [regressões de bytes da formulação](../../01_core/usecase/algorithm_input_bytes_regression_test.go).

## DEC-01 — Política finita e hash serializável

DecideReuse rejeita NaN e infinitos em ambos os limiares; os valores finitos continuam em (0,1]. Tratar erro de serialização da política como erro, sem produzir relatório útil nem hash de bytes vazios. Uma política válida mantém o hash SHA256 de seu JSON e os resultados anteriores; não modificar avaliações nem relaxar regras de decisão. O controle mínimo sem assessments deve retornar relatório serializável em todo sucesso. O recorte não promete admitir valores não finitos em outros campos públicos.

Consumidor: [regressões de política finita](../../01_core/usecase/reuse_policy_finite_regression_test.go). A falha foi reproduzida na API Go; JSON normal não representa NaN e não foi demonstrado bypass da entrada CLI.

## CAL-01 — Uma identidade por grupo de calibração

O [passo0167](../passos/0167-isolar-calibracao-e-validar-fontes-de-sintese.md) exige ao menos dois grupos com IdeaID não vazio, sem espaços nas extremidades e único por igualdade exata. Rejeitar antes de treinar, com erro e relatório zero; não agregar, renomear ou completar identidades automaticamente. Um registro representa uma ideia inteira no leave-one-idea-out. A distinção de caixa permanece significativa; a identidade é declarada, sem autenticar diversidade semântica.

Conservar a calibração dos datasets válidos, ordem dos folds, tratamento de uncertain, priors aprendidos só no treino, métricas e inputs. Não reclassificar resultados históricos. Consumidor: [regressões de grupos](../../01_core/usecase/calibration_group_admission_regression_test.go).

## CMP-01 — Conexões repetidas sem escolha implícita

Para cada par ordenado From/To, admitir repetições somente se as declarações forem estruturalmente idênticas: status, fontes na ordem fornecida e custo incluindo ausência, unidade e intervalo. Rejeitar divergência com erro e relatório zero antes da enumeração; não selecionar o primeiro custo, o menor ou inventar uma agregação. A direção oposta é outro par ordenado e conserva a política existente de conectividade.

Duplicatas idênticas mantêm comportamento idempotente, sem dupla cobrança ou evidência extra. Custo ausente continua desconhecido; igualdade estrutural não promove evidência nem dispensa o ensaio do conjunto. Preservar entradas, schemas e hashes históricos. Consumidor: [regressões de conexões duplicadas](../../01_core/usecase/composition_duplicate_connections_regression_test.go).

## SYN-01 — Identidade exata da execução lida

Após Get(source.RunID), exigir run.RunID == source.RunID antes de ler índices/páginas ou iniciar idempotência/fila. Ausência ou divergência é ErrIntegrity e base/recibo zero. Manter erros das portas, tipo, terminalidade, hashes e degradação existentes. A guarda pertence a improvement_sources em L1 e vale para seus dois consumidores duráveis.

## SYN-02 — Cursor nunca revisitado

Por artefato, lembrar todos os cursores de continuação já vistos. Recusar cursor vazio ou repetido, adjacente ou não, antes de solicitar novamente a página. Falha retorna ErrIntegrity e base zero, sem liberar itens parciais nem iniciar idempotência/fila. Conservar páginas de50, teto de1000, conferência de RunID/Artifact/ArtifactHash, cópia dos bytes, erros e término normal. Não compartilhar o conjunto de cursores entre artefatos distintos.

SYN-01/02 são consumidos pelas [regressões de identidade e paginação](../../01_core/usecase/synthesis_source_identity_pagination_regression_test.go), incluindo a API pública StartImprovementSynthesis e a porta restrita do resolvedor. Nenhuma guarda autentica origem remota ou recalcula hashes de registros anteriores.

## ASS-01 — Snapshot íntegro antes da união cega

O [passo0168](../passos/0168-preservar-integridade-e-determinismo-do-oraculo.md) introduziu admissão com replay íntegro de cada snapshot antes de admitir conteúdo: ReplayWebSnapshot sem erro, Replayable e sem MissingFields. Selo ausente/divergente, resposta adulterada ou material inválido gera ErrIntegrity e pacote zero; não reparar selos na admissão. Lista de snapshots vazia mantém o pacote vazio anterior. Fixtures positivas ativas antes não seladas passam a ser seladas explicitamente na preparação, sem alterar assertivas ou acervo histórico.

## ASS-02 — Variantes cegas sem sobrescrita

Agrupar por CandidateID derivado da URL canônica, sem expor URL/braço/query/rank. Por candidato, desduplicar os conteúdos de cada campo após TrimSpace, ordenar os textos e emitir campos na ordem title, description. Primeiro texto de cada campo usa o FragmentID legado; demais usam campo + hífen + SHA256 hexadecimal integral do conteúdo. Assim todas as variantes permanecem citáveis, com ordem/hash estáveis. Ocorrência única ou idêntica conserva pacote/IDs anteriores. Reordenar snapshots/respostas/resultados não altera a união. Conteúdo hostil continua dado, sem instrução operacional.

## ASS-03 — Gate de evidência prevalece

Revisão humana inválida mantém Effective=UNKNOWN, mas não substitui FAIL_EVIDENCE_GROUNDING já identificado por INCONCLUSIVE_HUMAN_REVIEW. Evidência inválida/ausente prevalece sobre pendência de revisão, independentemente da ordem dos candidatos. Revisões válidas mantêm a política existente; AIHistory e HumanHistory permanecem preservados.

## ASS-04 — Cobertura mede vínculo válido

EvidenceGroundingRate conta candidatos cujo julgamento selecionado passa judgmentEvidenceValid contra o próprio pacote; denominador é o número de candidatos do relatório. Julgamento ausente ou referência inexistente/inválida não conta. Corpus vazio permanece zero, sem NaN. Métrica de vínculo não autentica origem nem comprova sustentação semântica. Consumidores ASS-01–04: [regressões do oráculo](../../01_core/usecase/assisted_oracle_integrity_regression_test.go); transporte também rejeita snapshot inválido sem publicação.

## MCDA-01 — Normalização finita em escalas extremas

Manter a fórmula existente para spans finitos. Quando a diferença de extremos finitos transborda, dividir extremos e valores por2 antes das subtrações, sem alterar a entrada ou o significado da escala. Maximização/minimização mantêm endpoints0/1, pontos internos e intervalos; medida ausente continua incerteza na escala inteira, nunca zero. Sucesso deve produzir utilidades/contribuições/scores finitos e serializáveis, preservando hashes dos controles ordinários. Consumidor: [regressões de escala](../../01_core/usecase/multicriteria_extreme_scale_regression_test.go).

## MKT-01 — Representante e procedência canônicos

Para a mesma identidade exata connector_id/external_id, escolher a ocorrência de menor JSON canônico após canonicalização; não usar ordem de chegada. Preservar seus metadados escalares/ponteiros como conjunto, sem fabricar combinação semântica. Agregar conjuntos declarados de locales, mercados, países, repo links, origens e witnesses de todas as ocorrências, com cópias, desduplicação e ordem total. Witnesses canonicalizam ConnectorIDs e ordenam pela declaração completa; origens/repo links mantêm a ordenação anterior com desempate pela declaração completa. Duplicata continua SAME; ocorrência única conserva sua resolução padrão. Não ampliar equivalência entre conectores nem promover metadado a evidência verificada.

Consumidor: [regressões de consolidação](../../01_core/usecase/market_consolidation_order_regression_test.go). Manter inputs e hashes dos controles canônicos, e preservar todos os registros/recibos históricos. Os seis critérios0168 regem novas execuções, sem reparar resultados antigos.

## DSC-01 — União direta antes da escolha

O [passo0169](../passos/0169-canonicalizar-descoberta-e-validar-manifestos.md) canonicaliza cópias dos candidatos/matches antes da união. Representante mínimo pela ordem anterior de metadados, com desempate pelo JSON completo; conservar o representante inteiro, sem fabricar combinação de campos. Agregar todos os matches do repositório como origens, com evidências ordenadas por declaração completa e origens desduplicadas/ordenadas pela regra regional existente. Não descartar matches por serem do mesmo repo/attempt/query. Sem match, preservar fallback; match único mantém a projeção canônica anterior. Slices da porta não são modificadas; nil e listas vazias conservam sua representação. Não modificar snapshots antigos nem políticas de elegibilidade.

## DOS-01 — Hash obrigatório antes de anunciar estágio

BuildRecommendationRunStages mantém a assinatura pública. Erro do leitor de hashes ou hash vazio/em branco torna a entrada/saída indisponível: estágio afetado fica FAILED, sem Outputs, com motivo determinístico. Nenhum estágio com entrada/saída indisponível pode ser COMPLETE. A fachada grava manifesto FAILED e falhas por estágio, retorna erro e não chama Announce nem devolve Published=true. Erro ao persistir manifesto também retorna erro. Caminho válido conserva ordem, campos, versões e hashes anteriores; status operacional não promove evidência.

## DOS-02 — Locales canônicos antes da identidade

DossierEvidenceInputHash e RunRecommendationDossier canonicalizam SearchLocales antes de pré-imagens/identidade/checkpoint. Permutações, repetições e all equivalentes geram os mesmos hashes e RunID; locale inválido falha antes de Prepare ou publicação. Não modificar a entrada do caller. Perfil vazio continua omitido, com hash histórico intacto. A CLI já normaliza seus locales; a regra também pertence à admissão pública em L1.

## REG-01 — Chave de catálogo consistente

Concept deve ser não vazio, único e idêntico ao próprio TrimSpace. Chave com espaços nas extremidades é catálogo inconsistente, rejeitado antes do lookup; não reparar nem converter a inconsistência em lacuna. Tradução controlada realmente ausente para conceito válido continua Unknown. Nenhuma consulta livre ou promoção de cobertura é adicionada.

Consumidores DSC-01/REG-01: [regressões da descoberta](../../01_core/usecase/discovery_canonical_response_regression_test.go). DOS-01/02: [regressões do dossiê](../../01_core/usecase/dossier_integrity_regression_test.go), com falhas em todos os hashes obrigatórios, equivalência/retomada pela fachada e controles congelados de snapshot, pré-imagem, RunID e estágios. Regras0169 valem para novas execuções; registros, selos e hashes históricos não são recalculados.


## ASS-05 — Compatibilidade entre produtor e consumidor

A regra vigente admite WebSearchSnapshot schemas1 e2 em PrepareBlindCandidateUnion, mantendo replay íntegro, sem MissingFields e pacote zero/ErrIntegrity em falha. Schema2 é a saída do experimento autenticado0020; a restrição a schema1 introduzida em0168 interrompia a jornada0020→0021. O pacote cego permanece schema1. Não converter versões nem reparar selos existentes. Testar a ligação com saída real de ExecuteWebExperiment, incluindo hits públicos, e preservar controles schema1.

## WEB-01 — Retomada vinculada ao plano e aos selos

Validar toda a entrada de retomada antes de SearchOne: schemas1/2, replay íntegro, campos materiais, RunID/Arm/Provider/GrammarHash/QuerySetHash esperados. Schema2 exige Authorization/Budget iguais aos atuais; schema1 permite ausência legada desses dois metadados, mas exige igualdade se presentes. Respostas têm QueryID único conhecido, request exatamente igual ao planejado, Provider/Arm/Position corretos e AttemptNumber1. Hash recalculado não demonstra vínculo à consulta do plano. Braços N/F1/F2 são únicos; query IDs são não vazios/únicos por braço, textos não vazios. A entrada é copiada antes de ordenar.

Falha retorna erro de integridade/reprodutibilidade e resultado zero, sem busca adicional nem publicação. Subconjuntos válidos retomam SUCCESS/EMPTY, preservando buscas de ausentes/Unknown, sequência, limites e contadores. Não promover Unknown nem reescrever snapshots. O contrato protege a ligação declarada, sem autenticar resposta externa.

## ASS-06 — Históricos completos sem precedência implícita

Preservar todas as declarações de IA e revisão do candidato em ordem determinística pela declaração completa. Repetições idênticas são idempotentes para efetivo/métricas, mas permanecem auditáveis no histórico. Declarações distintas de qualquer dos dois tipos tornam o candidato ambíguo: escalonar, Effective=UNKNOWN e INCONCLUSIVE_HUMAN_REVIEW, com causas AMBIGUOUS_AI_JUDGMENTS/AMBIGUOUS_HUMAN_REVIEWS no campo existente. Qualquer vínculo inválido conserva precedência FAIL_EVIDENCE_GROUNDING.

Não atribuir cronologia a arrays sem prioridade/timestamp contratados. Ambiguidade não incrementa Reviewed, DefinitiveReviewed, DefinitiveAgreements, DefinitiveOverrides ou ResolvedUncertain; Escalated conta o caso. Grounding exige todos os julgamentos do candidato válidos. Preservar rubrica, amostra residual, formato e controles únicos; não alterar históricos já publicados.

## CAL-02 — Verdadeiro positivo exige fonte válida

TrueProven conta somente Expected e Observed PROVEN com validImmutableCalibrationSources aprovado. Manter denominadores de estados esperados/observados, InvalidSources e regras de veredito/diagnóstico. Fonte ausente ou de outro commit reduz precisão/recall no split e agregado; não fabricar100% de acerto numa avaliação inválida. Não recalcular métricas de registros históricos.

## ALT-01 — Diagnóstico estável em empates de metadados

RankSoftwareAlternatives usa ordem total dos snapshots: precedência de candidateMetadataLess seguida pelo JSON completo, conforme directCandidateMetadataLess. Escolher representante inteiro sem combinar campos e manter duplicatas INCONCLUSIVE. Stale/Archived e demais campos antes omitidos do comparador não podem alterar motivos ou hash pela ordem de chegada. Preservar política, inputs e controles únicos anteriores.

Consumidores de ASS-05/WEB-01: [integração e retomada web](../../01_core/usecase/integrity_0170_web_test.go). ASS-06: [históricos e ambiguidades](../../01_core/usecase/integrity_0170_judgment_test.go). CAL-02/ALT-01: [métricas e ordem de metadados](../../01_core/usecase/integrity_0170_metrics_test.go). Esses critérios regem novas execuções; não alteram registros, hashes, schemas publicados ou recibos de gates anteriores.


## CMP-02 — Declarações de partes sem escolha por chegada

Por compositionMemberID, comparar a declaração completa da parte válida/admissível antes de projetar capacidades comprovadas. Repetição estruturalmente idêntica é idempotente; declaração distinta é erro antes da enumeração, sem relatório útil nas APIs de decisão/multicritério. Capabilities, Evidence, Caveats e demais campos participam da comparação, preservando ordem declarada dos vetores. Não adicionar campos à identidade, unir alegações ou escolher representante por menor JSON/chegada. Preservar guardas de partes/fontes, caso único, inputs, IDs, custo e hashes do controle válido. A regra rege novas entradas, sem reparar composições históricas.

## WEB-02 — Margem persistida antes de novas buscas

Após admitir integralmente os snapshots da retomada, qualquer resposta admitida com Remaining<=SafeRemaining e reset não demonstradamente expirado impede novas buscas. Reset ausente ou relógio RuntimeNow zero não demonstram expiração. Com relógio conhecido, Reset<=Now.Unix já expirou. Conferir todas as respostas admitidas, inclusive Unknown, sem depender de sua ordem; não inferir quota nem chamar rede/espera para restaurá-la.

Preservar SUCCESS/EMPTY reutilizadas e não as repetir. Com restrição vigente, as demais queries ficam SKIPPED_BUDGET, Started/Retried=0 e StopReason=SAFE_RATE_LIMIT_MARGIN, com contagens reais e snapshots replayáveis. Após reset comprovado, faltantes/Unknown seguem o contrato de busca anterior. WEB-01 precede esta classificação: snapshot inválido retorna erro/output zero, mesmo diante de quota esgotada. Nenhum selo ou registro anterior é modificado.

## WEB-03 — Admissão dos limites na API

ExecuteWebExperiment admite limites antes de Capabilities, PreviousSnapshot e SearchOne. Autenticado exige MaxPerQuery=10 e SafeRemaining>=1, além da autorização válida; não autenticado exige limite positivo e conserva clamp de positivos acima10 para10. Erro retorna output zero e nenhuma chamada a portas. Compartilhar a política sem exigir flags de CLI na API. Preservar defaults e formatos dos pedidos válidos; helper de transporte não substitui a admissão de L1.

## ASS-07 — Métricas exigem revisão válida

Uma regra pura compartilhada governa efetivação e elegibilidade das métricas: CONFIRMED é válida e mantém o julgamento da IA; OVERRIDDEN/RESOLVED exigem veredito aceito pela rubrica e motivo não vazio após TrimSpace. Ação desconhecida é inválida. Revisão inválida não incrementa Reviewed, DefinitiveReviewed, DefinitiveAgreements, DefinitiveOverrides ou ResolvedUncertain; UNKNOWN não é acordo/override/resolução. Grounding e Escalated permanecem independentes; ambiguidades ASS-06 continuam excluídas e repetições idênticas válidas são idempotentes. Preservar histórias, precedência de falhas, efetivo, controles/hash válidos e resultados históricos.

Consumidores de CMP-02: [regressões de partes e integração com multicritério](../../01_core/usecase/integrity_0171_composition_test.go). WEB-02/03: [quota persistida e admissão da API](../../01_core/usecase/integrity_0171_web_test.go). ASS-07: [validade e denominadores de revisão](../../01_core/usecase/integrity_0171_metrics_test.go). As fixtures positivas ativas seguem os limites vigentes; registros e hashes publicados antes desta revisão permanecem inalterados.
