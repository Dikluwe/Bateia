# Prompt espelhado — oráculos A/B L1 do Passo 0025

Derive exclusivamente do Passo 0025, do contrato congelado e dos ataques congelados uma
suíte independente para a porta pura de recomendação. Congele os tipos públicos
`RecommendationSearch`, `RecommendationRequest`, `RecommendationSummary` e
`RecommendationResult`, o construtor `NewRecommendationService`, sua operação `Recommend`
e as projeções puras `ProjectRecommendationResult` e `MarshalRecommendationResult`.

Exija validação e normalização antes de expansor, descoberta ou conteúdo; execução `live`
offline com fakes; limites exatos; locales solicitados na descoberta e na proveniência;
licença desconhecida como aviso não bloqueante; preservação de `INCONCLUSIVE`; ausência de
`top_choice` impróprio; arrays não nulos; ordenação/deduplicação canônica; remoção de snippet,
conteúdo bruto, URL insegura e segredo-canário. Texto de repositório semelhante a instrução
permanece dado e nunca altera política, limites, locale ou autoridade de evidência.

Na revisão R2, cubra separadamente corpos integrais de README, fonte, teste e licença;
compare blockers conhecidos/desconhecidos; injete canários de token, Authorization, header,
cookie, ambiente, configuração de provedor e metadado do adaptador; cancele após uma
barreira de provedor já iniciada; e permute conjuntos de domínio sem mudar bytes públicos.

Na revisão R3, faça a instrução hostil atravessar descoberta e coleta reais do
`RecommendationService` e demonstre que não promove o estado observado a `PROVEN` nem a
execução a `COMPLETE`. Execute também duas recomendações equivalentes permutando coleções
set-like da request, candidatos e matches do provedor, ordem dos assessments/capabilities
derivados e ordem das fontes coletadas; os bytes canônicos e locales observados são iguais.

Na revisão R5, a fixture nominal usa exclusivamente regras já versionadas do catálogo:
`organização de receitas` e `importação de receitas`, ambas sob `recipe-import`. Fonte e teste
contêm padrões fortes reconhecidos; o candidato tem commit imutável, timestamps fixos não
nulos e metadados de manutenção; a execução conhecida contém o texto MIT integral. A execução
de licença desconhecida difere somente no conteúdo de `LICENSE`, mantendo todos os demais
gates conclusivos e as asserções de locale, proveniência e segurança.

Não leia, adapte nem corrija implementação candidata. `Unknown` não satisfaz caso negativo.

## 0158 — Double de conteúdo com identidade

No recorte de integridade0158, os doubles que retornam extração bem-sucedida devem preencher Ref com owner/repo/commit do candidato efetivamente solicitado. Conservar conteúdo, falhas, quotas, testemunhos, mutantes e assertivas funcionais anteriores. Ref ausente ou divergente passa a ser controle negativo explícito nos testes próprios de integridade, sem dispensar a validação produtiva.
