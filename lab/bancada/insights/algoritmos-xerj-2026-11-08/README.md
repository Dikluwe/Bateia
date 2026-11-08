# 0155 — Ideias de algoritmos e decisões: investigação do XERJ

**O Bateia foi ajustado e usado para investigar o XERJ. O resultado utilizável é uma hipótese condicional de formato/leitura e uma confirmação que evita uma recomendação duplicada. Nenhuma melhoria de desempenho foi demonstrada.** Cada registro liga princípio, passos, evidência da fonte e do alvo, condições/falhas, decisão e experimento de rejeição.

Alvo: `xerj-org/xerj` no commit `1557de2f5be981d9db441969201c92a8af6d2de3`, verificado como main no início. O usuário forneceu o alvo; referências externas precisavam surgir nativamente no Bateia. A pergunta foi aberta, sem algoritmos, arquivos ou fontes externas esperados.

## O que mudou no produto

- `--algorithm-ideas` em `discover`, `explore` e `evidence collect --question` acrescenta tarefas de entradas/estado/passos/invariantes, condições/custos/contraexemplos e decisão no alvo. Ativa a revisão de alternativas e passagens adicionais.
- `idea insights formulate --input ... --output ...` registra raciocínio fornecido pelo agente contra bundles congelados de **fonte** e **alvo**, com papéis separados. Exige princípio, passos, evidências, condições/falhas, decisão atual/proposta, benefício/custo e critério de rejeição.
- O núcleo confere identidade, commit, caminho, hash, linhas e literal. Rejeita trechos inventados, papéis trocados e campos essenciais ausentes. Não usa catálogo fechado de algoritmos nem chama modelo.
- O relatório conserva interpretação/relação com alvo como `DECLARED_NOT_REASSESSED` e efeito `UNVERIFIED`. Vínculo válido não significa ideia verdadeira, nova ou benéfica. O agente interpreta; revisão e experimento continuam necessários.

Regras e tipos ficam em [L1](../../../../01_core/usecase/algorithm_insight.go); L2 apenas adapta flags/JSON e publica. Prompts e [passo 0155](../../../../00_nucleo/passos/0155-descobrir-principios-e-decisoes-no-xerj.md) precederam a implementação. Os modos anteriores são preservados: o [replay](legacy-replay/result.json) teve saída padrão idêntica byte a byte. O [gate completo](verification/verification.json) passou em lint, arquitetura, testes, vet, build e diff. Testes novos usam agendamento fictício, não algoritmos ou marcas esperados desta investigação. O build verificado coincide com o binário congelado do ensaio.

## Resultado: evidência → condições → decisão

| Registro | Evidência e princípio | Condições e limites | Decisão que muda |
|---|---|---|---|
| **I1 — leitura seletiva de postings** | Fonte mantém diretórios residentes e offsets/comprimentos de blocos, lendo payloads por faixa. No recorte XERJ, o leitor abre/descomprime o envelope `.post` do campo e depois usa slices. | Só investigar benefício se abertura fria/miss tiver custo material e houver acesso seletivo. XERJ já tem cache FTS e `open_stats_only`. Não há ganho por consulta demonstrado. Framing, formatos antigos, scores, posições, frequências, checksums e total exato precisam ser preservados. | **Perfilar primeiro** a abertura/misses reais; só depois decidir por protótipo de framing e leitor por faixa, comparando com o leitor/cache existentes. Não aprovar migração. |
| **I2 — poda de score por limite superior** | Fonte externa explicita limites e testes de paridade; XERJ já encaminha um shape elegível de consulta a `wand_should_bool`, com comparação de limite superior e contagem exata. | Shape bool-should com termos simples, cap finito, min_should_match efetivo 1 e sem explain. A variante block-max não é automaticamente equivalente; pesos/limites BM25 e total exato precisam de contrato próprio. | **Descartar “adicionar WAND” como novidade.** Preservar o caminho existente e não propor a variante por blocos sem investigação adicional. |

### I1: fonte, comportamento atual e custo

A referência [arclabs561/postings, raw-segments.md:3](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/docs/raw-segments.md#L3) descreve diretórios residentes e leitura por faixa dos postings de termos. [raw.rs:609](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/src/raw.rs#L609) materializa base/último doc, offset e comprimento. [raw.rs:3308](https://github.com/arclabs561/postings/blob/edb02683c6359efe2034a6cc71b64c9aaed90cb0/src/raw.rs#L3308) usa esses metadados na leitura/decodificação. Essa organização é uma estratégia concreta de acesso, não apenas indicação de biblioteca.

No alvo, [xerj-fts/index.rs:2700](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/index.rs#L2700) separa abertura/descompressão e uso posterior por slices. [index.rs:2597](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/index.rs#L2597) já oferece abertura apenas de estatísticas. [postings.rs:74](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/postings.rs#L74) delimita a tabela de saltos não persistida nesse formato. Persistir esse vetor sozinho não tornaria o envelope comprimido inteiro acessível aleatoriamente: a proposta envolve framing independente, escritor/leitor/merge e compatibilidade.

**Correção trazida pela revisão:** o pesquisador declarou caches/rota de produção desconhecidos. O revisor encontrou [fts_reader_for](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-engine/src/index.rs#L33420): hit reutiliza Arc; miss abre o leitor sob orçamento. O loop FTS o chama em [index.rs:27061](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-engine/src/index.rs#L27061). O agente principal conferiu esses trechos em Git e registrou a [anotação adicional](review-annotation.json), sem reescrever o resultado congelado do pesquisador. Essa descoberta é da revisão, não atribuída à descoberta nativa do Bateia.

Assim, o ganho hipotético é reduzir bytes lidos/descomprimidos e memória retida **onde o leitor precise ser aberto e só parte do payload seja pertinente**. Pode perder por diretórios adicionais, mais I/O, compressão pior, código de migração e leitores já aquecidos. Não foram medidos hit/miss, orçamento efetivo, pressão de memória, frequência de reabertura ou distribuição de consultas. Comentários históricos de desempenho no alvo não são medições deste ensaio.

Experimento proposto, **não executado**: corpus representativo, consultas seletivas e amplas, cache de leitor e cache de páginas separados, baseline atual versus protótipo; medir bytes efetivos, memória transitória/retida, p50/p95, ingestão e espaço. Exigir paridade de IDs/scores/posições/total e leitura de legado/corrupção. Rejeitar a migração se o perfil não mostrar custo relevante ou se regressões/custos superarem o benefício. Os percentuais sugeridos pelo pesquisador são limiares negociáveis de um ensaio futuro, não estimativas de ganho.

### I2: um resultado útil também pode ser uma confirmação

A cadeia [search → search_bounded](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/search.rs#L594) já seleciona o [fast path WAND](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/search.rs#L622). O código [compara o limite superior com o pior retido](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/search.rs#L788) e [conta/avança os documentos](https://github.com/xerj-org/xerj/blob/1557de2f5be981d9db441969201c92a8af6d2de3/engine/crates/xerj-fts/src/search.rs#L824). Foram lidos testes específicos de paridade, sem executá-los.

A descoberta da fonte não autoriza anunciar WAND como algoritmo ausente no XERJ. Também não autoriza pular documentos e conservar automaticamente o total exato. A variante por blocos permanece não estabelecida como transferência pronta. O insight útil foi eliminar uma recomendação genérica e explicitar o contrato que uma mudança futura precisaria preservar.

## Como foi investigado e revisto

Um pesquisador novo, sem histórico ou referências prontas, escreveu [baseline](research/baseline.json) antes da primeira ferramenta. O alvo fornecido pelo usuário veio de um clone limpo no pin verificado. Toda leitura foi feita pelo Bateia; fontes externas precisavam de retorno nativo anterior à aquisição. Foram nove chamadas de pesquisa e duas aquisições externas, Rust e Java; o registro final usa evidência externa Rust. Não se executou fonte, não se alterou XERJ, não houve publicação, timer, feedback ou email enviado.

Houve um checkpoint genérico do agente principal pedindo verificar chamadores e alcance do caminho atual, sem fornecer arquivos, algoritmos ou fontes. Isso foi incluído no pacote de revisão. O ensaio é uma investigação assistida, não uma comparação cega entre ferramentas.

O [registro produzido pelo comando](research/formulate-output.json) contém dois itens e foi [reproduzido byte a byte](formulate-replay.json). Os [15 vínculos de fonte/alvo](evidence-check.json) conferem com blobs Git fixados e snapshots dos retornos reais, inclusive hashes e linhas. Todos os audits preservam o hash do baseline. Isso verifica a cadeia de evidência; não mede transferência ou desempenho.

Um revisor independente usou a [rubrica anterior aos resultados](review-rubric.md), o pacote congelado e os clones autorizados para conferir condições/chamadores. O [parecer](review.md) classificou:

- I1: evidência **3/3**, transferência **2/3**, decisão **2/2** — utilizável como hipótese condicional, com cache existente explícito.
- I2: evidência **3/3**, transferência **3/3**, decisão **2/2** — confirmação ALREADY_PRESENT, não melhoria nova.

As notas não foram somadas. O revisor não executou fontes, não verificou toda rota HTTP nem estabeleceu prevalência do caso frio. Zero melhorias demonstradas. A inspeção de cache corrigiu uma lacuna de cobertura e estreitou a conclusão em vez de promovê-la.

## Dados que podem entrar no email

Formulação fiel ao ensaio:

> Estou fazendo experimentos enquanto desenvolvo o Bateia para investigar ideias de algoritmos e outros caminhos em fontes reais. Usei o XERJ como alvo e registrei cada ideia com evidência, condições e a decisão que ela mudaria. Um exemplo foi uma hipótese de leitura seletiva de postings por blocos, que precisa ser avaliada frente ao cache FTS já existente; ainda não há ganho medido. Também confirmei que poda WAND já está presente para um tipo de consulta e descartei apresentá-la como novidade. O objetivo é ajudar a formular e filtrar hipóteses verificáveis a partir das implementações.

Não afirmar que o Bateia comprovou acelerar o XERJ, encontrou um defeito, criou algoritmo novo ou superou outra ferramenta. O dado positivo é a passagem rastreável de fonte para hipótese/condição/decisão, incluindo a rejeição de uma recomendação redundante e a correção pela revisão.

## Artefatos e limites

[Relatório original](research/report.md), [insights](research/insights.json), [fontes](research/sources.json), [cobertura](research/coverage.json), [pedido de formulação](research/formulate-input.json), [resultado](research/formulate-output.json), [síntese estruturada](comparison.json). Planos, retornos e snapshots ficam em `research`; clones e binários permanecem temporários e não entram no repositório. Setup registra pin, origem user do alvo e SHA-256 do binário; fingerprints ligam o código do Bateia ao ensaio.

É um alvo e uma rodada. Princípios amplos já estavam parcialmente no baseline; o efeito da pesquisa foi concretizar, delimitar e descartar. Não há estimativa de criatividade, eficiência em tokens, novidade no mundo ou superioridade geral. As condições e decisões são interpretações verificáveis por revisão, não conclusões semânticas automáticas do núcleo. O manifesto SHA-256 conserva os artefatos finais.
