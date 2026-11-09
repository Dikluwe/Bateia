# Integridade de entradas e replay v1

Intenção ex-ante do [passo0166](../passos/0166-validar-bytes-testemunhos-e-politicas.md). Responsáveis L1: formulação de algoritmos e decisão de reúso; as correções de recuperação/releitura técnica estão no [contrato de coleta](coleta-e-serializacao-de-evidencias-v1.md#testemunhos-pertinentes-e-codificação-técnica--0166). Preservar formatos públicos, interpretações declaradas e registros/hashes históricos.

## ALG-01 — Validar os bytes originais

FormulateAlgorithmInsights aplica suas validações aos textos e conteúdos originais, antes de qualquer cópia capaz de substituir UTF-8 inválido. Serialização limitada pode medir o pedido, mas não reparar os bytes que serão validados. Conteúdo inválido em qualquer item de fonte/alvo gera ErrIntegrity; campos textuais inválidos continuam rejeitados pelas guardas existentes. Hash de conteúdo substituído não autentica o original inválido. U+FFFD originalmente válido permanece permitido; não proibir esse caractere nem normalizar conteúdo de evidência.

Copiar os slices que serão normalizados/ordenados sem alterar strings, caller ou registros anteriores. Preservar quotas, contexto, papéis, referências, canonicalização e IDs válidos. A saída mantém DECLARED_NOT_REASSESSED/UNVERIFIED; validar vínculo não comprova interpretação. Consumidor: [regressões de bytes da formulação](../../01_core/usecase/algorithm_input_bytes_regression_test.go).

## DEC-01 — Política finita e hash serializável

DecideReuse rejeita NaN e infinitos em ambos os limiares; os valores finitos continuam em (0,1]. Tratar erro de serialização da política como erro, sem produzir relatório útil nem hash de bytes vazios. Uma política válida mantém o hash SHA256 de seu JSON e os resultados anteriores; não modificar avaliações nem relaxar regras de decisão. O controle mínimo sem assessments deve retornar relatório serializável em todo sucesso. O recorte não promete admitir valores não finitos em outros campos públicos.

Consumidor: [regressões de política finita](../../01_core/usecase/reuse_policy_finite_regression_test.go). A falha foi reproduzida na API Go; JSON normal não representa NaN e não foi demonstrado bypass da entrada CLI.
