# Avaliação B — passo0048

2026-09-08. Regime A/B executado sem atestação de isolamento físico.

Aprovado no recorte observável congelado. Testes novos escritos antes da candidata a partir do passo0048, AGENTS e tipos públicos anteriores; o autor B não leu ranking/implementação candidata ou testes/gates antigos para autoria de casos. Os testes anteriores e runners foram consultados somente após congelamento dos casos novos. Hashes do prompt, spec e testes permanecem intactos.

Os testes0048 passam na primeira candidata. Baseline falha por ausência da nova API: isto registra distinção estrutural, não força discriminatória contra mutações semânticas. Controles0045,0046 e 0047 passaram em baseline e candidata.

Corpus real0046 mantém recall rotulado 2/6@3 e 4/6@12; holdout0047 mantém 2/2 em ambos limites. Negativas com indícios permanecem 1 em cada corpus. Mesmo conjunto de snapshots, hashes, perguntas, rótulos e budgets. Sem ganho artificial exigido e sem ajuste dos esperados. candidate-comparison.json registra maintenance_passed=true.

Runners0047 foram copiados para 0048 e adaptados somente para seleção da raiz baseline e cache0048; congelamento local registra adaptação e referência ao congelamento anterior. Baseline isolada core coincide com todos os hashes core do contract-baseline0048. Artefatos0047 não foram alterados.

Limites: corpus conveniente e holdout pequeno; indícios lexicais não verificam capacidade nem negação; top-k representa caminhos distintos nos 24 hits, não execução de extração. Testes0048 cobrem interpretação, ranking público e projeção de requisitos inclusive relatório antigo. Renderer/CLI, revalidação no usecase, suite completa, lint e build são verificados pelo integrador, fora deste veredito B. Nenhuma atestação física ou equivalência funcional geral é alegada.

Revisão mecânica r1: lint V14 rejeitou import slices; substituído por loops locais Equal/Contains, sort.Strings e append para cópia. Consultas/asserts/esperados preservados; freeze original preservado e freeze-r1.json encadeia hashes anteriores e novos antes de reexecutar. Recorte TestQuery0048 passou novamente (execution-r1.json). Produção não foi alterada.

Extensão independente r2: contraexemplo público s/use gerou novos controles de identificador isolado/backticks e fonte exata. Escrita sem ler candidata; congelada após implementação inicial e antes da próxima execução. A correção já estava pronta na execução, portanto não houve execução falha registrada da extensão contra a candidata inicial. Todos os testes0048 e controles0045/46/47 passaram. Reexecução final dos mesmos corpora/rótulos/budgets manteve integralmente métricas: final-comparison.json maintenance_passed=true. Hashes anteriores preservados em cadeia até freeze-r2.json; nenhum esperado anterior ajustado.
