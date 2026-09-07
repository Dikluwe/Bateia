# Qualidade da análise após o passo 0035

Avaliação de 2026-09-07. Escopo: saída integrada Coin/OIT salva no passo 0035, JSON integral, Markdown, implementação do Bateia e intervalos do Coin correspondentes. Não é uma nova busca pública, benchmark de modelos ou auditoria independente. Nenhum código foi corrigido nesta avaliação.

## Veredito

A coleta e a rastreabilidade melhoraram. A resposta analítica entregue ao consumidor ainda é fraca: o Bateia obtém evidência pertinente, mas perde parte decisiva na seleção/apresentação e não a incorpora na classificação de pertinência. Não há base para declarar qualidade geral de análise validada.

Os testes anteriores demonstram propriedades de infraestrutura (AST, hashes, limites, classificação de candidato), não que o relatório responda adequadamente à consulta do usuário.

## Achados prioritários

### P1 — Scores continuam desconectados das implementações

O relatório atribui Coin 16/100 por um termo do README, 8/100 de similaridade e `Divergent`, com intenção/stack/capacidades em zero. Entretanto, seu próprio JSON tem as definições e chamadas de implementação de transparência por camadas em `SoGLRenderAction.cpp`.

`01_core/usecase/ontology_matcher.go:28` usa descrição/nome para intenção; a função não consome `repo.Semantic`. `funnel.go` calcula relevância antes do Clang e não a reavalia com as evidências semânticas. A nota geral 8 resulta do arredondamento de 50 de topologia × 15%, com demais vetores em zero; não é uma medição de incompatibilidade com OIT. Ausência de sinal no modelo atual acaba apresentada como divergência do projeto.

Além disso, `deriveProblemDomain` retorna primeira palavra + `System` quando as heurísticas conhecidas não se aplicam. Assim, a consulta de gráficos é apresentada como `Order System`, sem evidência de domínio. Fonte: `ontology_matcher.go:305`.

### P1 — A IA que lê o Markdown perde as funções centrais

Amostra de seis definições confirmadas no JSON e nos intervalos do fonte:

| Função | Linhas em SoGLRenderAction.cpp | JSON | Markdown gerado |
| --- | --- | --- | --- |
| handleTransparency | 1205–1311 | Sim | Sim |
| doSortedLayersBlendRendering | 2042–2075 | Sim | Não |
| renderOneBlendLayer | 2117–2200 | Sim | Não |
| initSortedLayersBlendRendering | 2203–2254 | Sim | Não |
| renderSortedLayersFP | 2525–2580 | Sim | Não |
| renderSortedLayersNV | 2583–2678 | Sim | Não |

`doSortedLayersBlendRendering` chama inicialização na linha 2046, configuração de texturas na 2047, renderização da camada na 2063 e composição FP/NV nas linhas 2071/2073. O hash do arquivo presente hoje foi conferido contra o hash registrado no JSON: corresponde.

`02_shell/report/generator.go` ordena por substring de nomes/caminhos e exibe só os primeiros 12 símbolos e chamadas por unidade. Não segue a cadeia de chamadas. Operadores de acesso e declarações de headers ocupam espaço que poderia explicar o mecanismo.

As seis funções são uma amostra dirigida deste caso, não um denominador universal de recall. A descrição delas no README de validação foi uma análise adicional feita pelo agente; não equivale à síntese automática do relatório do Bateia.

### P1 — Extração textual e análise Clang seguem caminhos separados

`Files` contém 14 arquivos: dois documentos/manifestos e 12 arquivos de código. Oito desses 12 são headers XML. O corpo de `SoGLRenderAction.cpp` não está em `Files`, embora sua AST esteja em `Semantic.Units`.

`01_core/parser/evidence.go:78` usa termos no caminho e desempate alfabético, preenchendo o restante da cota com caminhos sem correspondência. `funnel.go` executa a extração antes do Clang; os arquivos/intervalos descobertos por ele não realimentam a extração textual. Portanto ainda há trechos iniciais de licença/XML, sem o corpo relevante do renderizador no conteúdo fornecido.

### P2 — Relatório pesado e repetitivo

Medições do Markdown e JSON do passo 0035:

- 476.063 bytes e 4.139 linhas no Markdown.
- 3.073 linhas contendo hashes, somando 394.697 bytes: aproximadamente 82,9% do documento. Isso mede volume, não tokens nem relevância de cada linha.
- 30.495 registros de símbolos; 10.145 combinações únicas de arquivo/linha/nome/tipo/estado de definição.
- 8.157 registros de chamadas; 4.958 combinações únicas de arquivo/linha/identidades de origem/destino.
- 3.059 registros de inputs por unidade; 473 caminhos únicos.

Headers repetidos entre unidades não são necessariamente um erro do Clang; tratá-los como informação nova no relatório aumenta o custo de leitura. `<details>` colapsa a interface visual, mas mantém o texto enviado ao consumidor de Markdown. Hashes devem permanecer disponíveis em manifesto deduplicado/anexo, com referência no resumo.

### P2 — Metadado de definição não é correto para todos os tipos

A enumeração `TransparencyType` em `include/Inventor/actions/SoGLRenderAction.h:61–69`, com corpo, aparece com `definition=false`. O adaptador só reconhece `completeDefinition`, `CompoundStmt` e `CXXTryStmt`; a classificação precisa ser específica para funções, enums, variáveis etc. Fonte: `03_infra/clang/ast.go:82`. Isso não invalida as seis definições de funções verificadas, mas impede confiar indiscriminadamente no campo para todos os símbolos.

## Melhorias efetivamente demonstradas

- O pipeline C/C++ agora encontra unidades pertinentes e resolve tipos/chamadas com a configuração de compilação registrada.
- Os seis mecanismos acima existem no JSON; antes o fluxo de extração não contemplava C/C++.
- Cobertura parcial 12/898, hashes, parâmetros e diagnóstico estão explícitos.
- Não há promoção automática a solução exata ou saúde 100; a comparação circular foi removida.

Nada disso demonstra qualidade geral em outros repositórios ou tarefas, nem execução gráfica correta do Coin. Os 195 avisos eliminados são um resultado de correção do lint, não uma métrica de qualidade da análise do Bateia.

## Próxima correção que mais agrega valor

1. Produzir um conjunto pequeno de evidências por pergunta: implementação principal, encadeamento de chamadas e trechos exatos. Usar esse conjunto para realimentar a extração textual e a apresentação.
2. Separar score de descoberta e avaliação de pertinência da implementação. Quando a evidência semântica ainda não foi avaliada, reportar pertinência indeterminada; não inferir divergência de vetores vazios. Retirar domínio fabricado pela primeira palavra.
3. Deduplicar símbolos/inputs, reduzir operadores e boilerplate na resposta principal e mover recibos completos para artefatos vinculados.
4. Validar a projeção destinada à IA: no caso Coin, o Markdown deve explicar a cadeia por camadas, citar os intervalos do renderer e preservar limites de execução. Acrescentar outros casos e negativos para evitar uma solução específica para o Coin.

Resposta mínima que o conteúdo coletado já permite, mas que o relatório ainda não sintetiza automaticamente: “Há uma implementação pertinente de transparência por camadas em SoGLRenderAction.cpp: doSortedLayersBlendRendering inicializa estado/texturas, itera renderOneBlendLayer e compõe via renderSortedLayersFP/NV. Esses caminhos foram analisados estaticamente na configuração registrada. Correção visual, GPU e adequação completa à consulta permanecem não verificadas.”
