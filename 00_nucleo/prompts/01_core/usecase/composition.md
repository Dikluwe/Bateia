# Resolução limitada de composições

Enumerar conjuntos de componentes válidos com limites explícitos, ordem canônica e sem I/O. Gates individuais de licença, manutenção, operação e evidência precedem a composição. Capacidades de partes só contam quando a capacidade está comprovada no assessment do mesmo repositório. Aplicar dependências e exclusões, verificar conectividade de interfaces e ensaio da composição exata. Sem testemunho admissível, devolver hipótese PROPOSED e abster-se de COMPOSE. Registrar truncamento do espaço de busca e não converter resultado parcial em conclusão negativa. Somar intervalos de custo apenas quando todos os vínculos necessários possuem custo na mesma unidade. Não inferir integração de cobertura isolada.

O estado VERIFIED significa atestado externo admissível de execução e aprovação, não que o Bateia executou o ensaio. Testemunhos de interface devem ser de teste de implementação. Custo desconhecido não se torna zero.

Uma parte só herda capacidade PROVEN do assessment quando pelo menos uma fonte da parte coincide com o testemunho daquela capacidade; prova em outro módulo não basta.

## 0167 — Conexões repetidas

Aplicar CMP-01 do [contrato](../../../contratos/integridade-de-entradas-e-replay-v1.md#cmp-01--conexões-repetidas-sem-escolha-implícita), conforme [passo0167](../../../passos/0167-isolar-calibracao-e-validar-fontes-de-sintese.md). Por par ordenado From/To, permitir somente duplicatas estruturalmente idênticas (status, fontes ordenadas e custo com ausência/unidade/intervalo). Rejeitar divergentes antes de enumerar, sem escolher custo por ordem ou mínimo. Preservar controle único, idempotência de duplicatas idênticas, custo desconhecido, direção oposta e inputs/registros anteriores.

## Admissão de partes repetidas — 0171

Antes de projetar as capacidades comprovadas de partes válidas/admissíveis, conferir a declaração completa por compositionMemberID. Duplicata estruturalmente idêntica é idempotente; declaração distinta do mesmo ID é erro, independentemente da ordem. Não unir capacidades/evidências, acrescentá-las à identidade nem escolher representante por chegada ou menor JSON. A comparação inclui Capabilities, Evidence, Caveats e demais campos, na ordem declarada de seus vetores. Preservar identidade existente, guardas de fontes/partes, caller e hashes de controles únicos/idênticos. DecideReuse e AnalyzeMulticriteriaDecision não devolvem relatório útil para ambiguidade admitida.
