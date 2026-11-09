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
