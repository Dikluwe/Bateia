# Bateia: ampliar a investigação por linguagem

**Foram pesquisadas e lidas fontes de sete linguagens:** Rust, Go, Java, JavaScript, TypeScript, Haskell e C++. A leitura trouxe mecanismos para tornar as propostas ao XERJ mais concretas, especialmente redução com precondições, persistência de falhas, replay e comparação entre implementações independentes.

Esta é uma etapa adicional, solicitada depois dos dois braços da [comparação completa](../README.md). As referências e resultados abaixo não foram fornecidos àqueles agentes nem contam a favor do braço Bateia original.

## Fontes efetivamente lidas

| Linguagem usada na consulta | Fonte encontrada pelo catálogo | Material lido | Mecanismo observado |
| --- | --- | --- | --- |
| Rust | proptest-rs/proptest | `proptest/README.md` | Estratégias de geração/redução e persistência de falhas em `proptest-regressions` |
| Go | flyingmutant/rapid | `README.md` | Minimização automática, persistência/reexecução e testes de máquina de estados |
| Java | jqwik-team/jqwik | README e dois arquivos de testes/exemplos | Redução de sequências com atenção às precondições das ações; asserção de caso mínimo |
| JavaScript | dubzzz/fuzz-rest-api | `README.md` | Geração de pedidos REST e reprodução de falhas HTTP com seed/path |
| TypeScript | dubzzz/fast-check | README do pacote, configuração, runner de modelo e teste de replay | Seed + caminho do contraexemplo; comandos condicionados ao modelo e replay verificável |
| Haskell | hedgehogqa/haskell-hedgehog | `README.md` | Shrinking integrado que declara preservar invariantes, com máquina de estados abstrata |
| C++ | div0rce/quant-systems-lab | `README.md` | Replay determinístico e comparação com modelo OCaml independente, declarados na documentação |

As linguagens são filtros do catálogo GitHub e contexto dos arquivos, não uma medição uniforme da capacidade semântica dos parsers do Bateia. O OCaml é uma segunda implementação descrita na fonte C++; não lemos seu código nem a executamos. Arquivos de testes mostram asserções, não resultados de uma execução nesta sessão.

## Ideias que ficaram mais específicas para o XERJ

1. **Reduzir uma falha preservando suas condições.** Gerar sequências de operações de evidência e diminuir a sequência sem remover a precondição que tornava o caso significativo. O teste de jqwik e as declarações de Hedgehog/rapid ajudam a formular o critério de rejeição: um caso menor inválido não é uma reprodução bem-sucedida.
2. **Conservar regressões mínimas com escopo de versão.** Proptest e rapid documentam persistência e reexecução de casos reduzidos. Aplicar esse padrão às propriedades específicas de integridade/calibração/revisão do XERJ, mantendo separados caso histórico, versão executada e resultado atual.
3. **Reproduzir com parâmetros suficientes.** Fast-check explicita seed e caminho de redução; seu teste verifica reconstrução do mesmo contraexemplo. Um recibo do XERJ precisaria também de versão do gerador, dados/pin, política e ambiente pertinentes. Seed sozinho não garante equivalência.
4. **Conferir a implementação com um modelo independente.** A referência C++ descreve um motor OCaml que reconstrói o mesmo estado. Para o XERJ, começar por um contrato pequeno de contagem, revisão ou admissibilidade e comparar decisões/estados. Acordo entre implementações não comprova por si só que a especificação está correta.
5. **Separar robustez HTTP de verdade da resposta.** A prova de conceito JavaScript usa ausência de HTTP500 como propriedade. Esse tipo de teste pode conferir contratos de erro do XERJ, mas não deve promover uma resposta semanticamente inadequada a resposta correta.

Cada proposta tem observação, hipótese, condições, experimento de rejeição e trechos com URL/commit/hash/linhas em [ideas.json](ideas.json). Nenhuma foi implementada ou teve benefício demonstrado. Antes de materializar, é necessário conferir a cobertura já existente no XERJ; a pesquisa não demonstra ausência desses mecanismos em todo o produto.

## Como a busca foi feita

O operador escolheu sete linguagens e a consulta `property testing language:<linguagem>`; não forneceu nomes de projetos. O Bateia executou `idea catalog discover` sobre o [plano declarado](catalog-plan.json), com até três candidatos por consulta. As sete buscas responderam. O filtro de linguagem foi preservado pelo catálogo, evitando tratar a palavra Rust/Java como simples tema lexical.

As primeiras respostas JavaScript/C++ incluíram bibliotecas genéricas, propriedades Qt e projetos gráficos pouco pertinentes. Duas [consultas refinadas](refined-plan.json) por tópico `property-based-testing` trouxeram candidatos mais específicos. Essas reformulações e a escolha de uma fonte por linguagem são trabalho do operador, não expansão semântica automática do produto. Os nove resultados de consulta e candidatos ficam em [catalog.json](catalog.json) e [refined-catalog.json](refined-catalog.json).

A coleta remota não resolveu os sete candidatos para commits imutáveis; os bundles de falha permanecem em [bundles.json](bundles.json). As fontes foram então adquiridas por URL exata a partir do catálogo, e o Bateia fez coleta local em revisões fixadas. Os comandos, tempos e pins ficam em [acquisitions.json](acquisitions.json). Essa aquisição manual não é apresentada como coleta remota bem-sucedida.

O README raiz do fast-check não pôde ser extraído; o inventário do coletor apontou o README real em `packages/fast-check/README.md`. A primeira leitura do jqwik continha pouca informação de mecanismo; arquivos de testes/exemplos foram selecionados a partir do mesmo inventário. O README raiz do proptest é um symlink; uma coleta adicional do caminho canônico `proptest/README.md` tornou o localizador compatível com os bytes do objeto Git. Os bundles originais e a primeira exploração foram preservados.

## Contexto, interpretação e verificação

O Bateia preparou 12 cartões de 12 arquivos, cobrindo sete fontes, com `CONTEXT_PREPARED`, novidade `NOT_ASSESSED` e efeito `UNVERIFIED`. Todos os bundles continuam incompletos: ler essa seleção não equivale a analisar integralmente os projetos.

Os arquivos de leitura foram unidos numa projeção por repositório/commit, removendo metadata e diagnósticos volumosos de árvore, preservando falhas e marcando incompletude. A [derivação](projection.json) registra hashes dos originais; não houve reescrita dos bundles brutos. O contexto está em [exploration.json](exploration.json). A IA interpretou os cartões e trechos adicionais dos mesmos arquivos coletados; o executável não gerou sozinho as cinco propostas.

A conferência em [verification.json](verification.json) validou os hashes, os 12 localizadores e a igualdade dos conteúdos selecionados com os objetos Git nos commits registrados. Isso vincula a observação aos bytes, sem autenticar as alegações documentais ou demonstrar execução. Não foram instaladas dependências nem executados os projetos ou seus testes.

As fontes repetem alguns mecanismos já vistos com Hypothesis; sete linguagens não significam sete conceitos novos. O acréscimo mais distinto é a analogia de replay diferencial C++/OCaml. Os demais resultados tornam explícitas condições que uma implementação futura deve preservar ao gerar, reduzir e reproduzir falhas.
