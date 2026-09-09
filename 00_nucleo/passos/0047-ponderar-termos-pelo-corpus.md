# Passo0047 — Ranking por poder discriminatório dos termos

2026-09-07 (local). Contrato congelado antes da implementação. Baseline: worktree0046, incluindo rankingv2 herdado0045. Tekt A/B sem atestação física.

## Ideia investigada com o Bateia

Executar o Bateia para buscar referências de ranking. A consulta local `BM25 idf scoring terms query` no Blevev2.6.0 recuperou search/searcher/search_term.go:newTermSearcherFromReader e search/scorer/scorer_term.go:docScore. Inspeção dessas fontes mostrou cálculo de frequência documental e normalização de comprimento. Isso inspira um ensaio de pontuação; não demonstra adequação automática ao código ou autorização para chamar a nova heurística de BM25 completo. Guardar consultas, resultados e hashes. Busca pública tem resultado/erro registrado separadamente.

## Mudança delimitada

O ranking atual soma pesos quase constantes, favorecendo funções grandes e termos comuns. Implementar ponderação por frequência de termos no corpus elegível de declarações AST, com normalização de comprimento, preservando orçamento, campos, vocabulário e identidade. Não ampliar glossário com perguntas do benchmark; não adicionar LLM, embeddings, serviços ou dependências. Não pontuar comentários/literais isolados. Unidade documental é declaração AST; declarar denominador/escopo. Campo de nome deve ter peso próprio, corpo usa identificadores AST únicos, comprimento registra conjunto de termos. Funções pertinentes com corpo precedem tipos não exatos; identificador exato preserva prioridade absoluta. Duplicação de entradas não muda estatísticas. Orçamentos incompletos permanecem parciais; frequência é do recorte lido, não do ecossistema.

Adicionar domain.ImplementationTermWeight {Term, Field string; DocumentFrequency, CorpusDocuments int; Weight float64}, e ImplementationHit.TermWeights []ImplementationTermWeight com JSON snake_case/omitempty; Weight é contribuição numérica finita efetivamente usada, não probabilidade. Reason informa método/função do peso e limites. Score inteiro permanece heurístico; um denominador estável e desempate determinístico devem preservar reprodutibilidade. Expansão por chamada mantém contrato0045, sem inventar peso de pertinência para helper sem termos.

Não mudar casamento lexical, fontes de requisitos ou detectores para converter notas em comprovação. Markdown mostra até cinco contribuições por hit, restantes no JSON. Versionar método/analisador/checkpoint; revalidação conserva pesos apenas nas fontes de bytes confirmados. Usar função pura de cálculo em parser/corpus_scoring.go com prompt prévio; testes separados.

## Avaliação antes de promover

Autor B reusa corpus real0046 congelado, preserva perguntas/esperados/runner original, e executa controle0045. Adiciona pelo menos um projeto real não usado no0046 e duas perguntas funcionais mais uma negativa como holdout, com commit/origem/hashes, congelado antes da candidata e sem ler produção. Nenhum código externo executado. Medir baseline primeiro; autora implementação não recebe holdout/expected.

Aceitar mudança de ranking somente com ganho de recall@3 ou@12 no corpus real0046, sem queda de recall no outro corte, e sem queda de recall do holdout ou aumento de negativas com indícios nos dois corpora; preservar controles sintéticos/exatos/negativos0045. Precisão é rotulada e limites precisam ser explicados. Máximo duas revisões sem ganho na mesma causa; se não atingir critérios, descartar ranking experimental e entregar resultados sem alegar melhoria. Não ajustar casos para passar.

Testes independentes verificam termo raro versus comum, normalização de função muito longa, pesos/denominadores finitos e reprodutíveis, exato prioritário, escopo só AST, dados reordenados/duplicados, chamadas e hashes preservados. Sem exigir valores mágicos de pesos. Suite completa, lint, build ./04_wiring, CLI local/replay, corpus real e holdout. Não promover resultados conhecidos como evidência independente de generalização.

## Fora do escopo

Tradução semântica universal, equivalência de programas, avaliação de linguagens além Go, prova de capacidades, correções Bleve/Sourcegraph/gRPC, alterações de gates/certificados e commits.
