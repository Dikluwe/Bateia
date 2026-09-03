# Passo 0027 — Generalizar descoberta multilíngue e multifuente

**Estado**: `PLANEJADO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Depende de**: [`Passo 0026 — Persistir execuções MCP e consultar por run_id`](0026-persistir-execucoes-mcp-e-consultar-por-run-id.md)
**Baselines obrigatórios**: certificados `MCP_RECOMMENDATION_GATEWAY_READY` e `MCP_DURABLE_RUNS_READY`
**Onde roda**: planejamento puro em L1, conectores públicos em L3 e execução durável pelo MCP
**Camadas previstas**: L0, L1, L2, L3 e L4
**Natureza**: descoberta de produtos por múltiplas fontes, idiomas e mercados, sem IA embutida

## Problema

Os Passos 0025 e 0026 tornaram o Bateia utilizável por uma IA externa e seguro
para investigações longas. Entretanto, isso ampliou a forma de usar o motor, não
o universo que ele consegue observar. A recomendação continua centrada em
repositórios e a descoberta regional existente usa um catálogo controlado para
um domínio estreito.

Essa limitação impede a hipótese principal do produto: uma aplicação conhecida
em outro idioma ou mercado pode conter uma solução parcial, um fluxo ou uma
decisão de produto útil, mesmo sem possuir repositório público.

Também existe uma ambiguidade perigosa. Encontrar uma página com consulta
japonesa no mercado dos Estados Unidos não prova que o produto seja japonês,
esteja disponível no Japão, tenha interface japonesa ou tenha sido criado por
uma organização japonesa. Idioma, mercado observado e origem são fatos distintos.

## Incerteza que este passo reduz

Não sabemos se uma camada aberta de conectores públicos, combinada com consultas
multilíngues rastreáveis, aumenta o recall de produtos funcionalmente relevantes
sem:

- inventar região ou idioma do produto;
- fundir produtos homônimos;
- transformar metadado de catálogo em prova de implementação;
- reintroduzir conteúdo não confiável no contexto da IA;
- duplicar a fila, persistência ou o pipeline já certificados;
- esconder falhas de cobertura como ausência de alternativas.

## Objetivo

Materializar um registro de conectores observáveis e uma execução durável de
descoberta de produtos. Uma IA externa pode listar as fontes disponíveis, propor
consultas em vários idiomas, selecionar mercados e iniciar a investigação. O
Bateia valida, executa, congela respostas, resolve identidades com cautela e
publica candidatos e lacunas com proveniência.

```text
usuário
  → IA externa formula IdeaSpec e propostas de consulta
  → bateia_list_connectors
  → bateia_start_market_discovery
  ← run_id

worker durável do Passo 0026
  → plano locale × mercado × conector
  → catálogos de repositório + catálogo público de aplicações/produtos
  → snapshots por tentativa
  → resolução conservadora de identidade
  → market_candidates + market_gaps

IA externa
  → usa candidatos como alternativas ou inspiração
  → encaminha vínculos repo verificados ao funil de evidência existente
```

O Bateia não chama um modelo de IA. Propostas da IA entram como dados não
confiáveis, nunca como fatos ou evidências.

## Recorte da primeira materialização

Este passo entrega:

1. registro versionado de conectores e suas capacidades;
2. preservação dos conectores repo-backed existentes como baseline;
3. exatamente um novo conector público cuja unidade seja aplicação ou produto,
   não ocorrência de código;
4. planejamento explícito por idioma de consulta e mercado;
5. identidade de produto que não depende de `owner/repo`;
6. deduplicação conservadora dentro e entre fontes;
7. execução como novo tipo de job no runtime durável do Passo 0026;
8. artefatos MCP seguros para plano, snapshot, candidatos e lacunas;
9. benchmark multilíngue com replay offline.

A fonte não repo-backed deve ser escolhida e congelada no contrato antes da
implementação, mediante capacidade documentada, acesso permitido, resposta
congelável e cobertura da matriz de teste. Não vale implementar vários conectores
rasos para declarar cobertura nominal.

## Quatro dimensões que não podem ser confundidas

### 1. Idioma da consulta

Tag BCP 47 usada no texto enviado ao conector, por exemplo `pt-BR`, `ja-JP` ou
`ar-SA`. É proveniência da busca, não atributo do produto.

### 2. Idioma observado da listagem

Idioma declarado pela fonte ou detectado por mecanismo determinístico e
versionado sobre o campo retornado. Detecção incerta permanece desconhecida.

### 3. Mercado consultado ou disponibilidade observada

Código ISO 3166-1 alpha-2 passado ao filtro oficial da fonte ou declarado na
resposta. Significa apenas que o produto foi observado naquela superfície e
naquele instante.

### 4. Origem geográfica do produto ou publicador

Só existe quando uma fonte explícita e rastreável declara país do publicador ou
da organização. Não pode ser inferida de idioma, alfabeto, TLD, moeda, nome,
posição no ranking ou mercado consultado.

## Registro de conectores

L1 define uma descrição pura equivalente a:

```go
type DiscoveryConnectorKind string

const (
    ConnectorRepositoryCatalog DiscoveryConnectorKind = "REPOSITORY_CATALOG"
    ConnectorApplicationCatalog DiscoveryConnectorKind = "APPLICATION_CATALOG"
    ConnectorProductCatalog     DiscoveryConnectorKind = "PRODUCT_CATALOG"
)

type DiscoveryConnectorCapabilities struct {
    ID                   string                   `json:"id"`
    Version              string                   `json:"version"`
    Kind                 DiscoveryConnectorKind   `json:"kind"`
    QueryLocales         []string                 `json:"query_locales"`
    Markets              []string                 `json:"markets"`
    SupportsMarketFilter bool                     `json:"supports_market_filter"`
    IdentityLevels       []string                 `json:"identity_levels"`
    EvidenceSurfaces     []string                 `json:"evidence_surfaces"`
    Replayable           bool                     `json:"replayable"`
    PublicOnly           bool                     `json:"public_only"`
}
```

Capacidades são declaradas antes da primeira chamada e ordenadas canonicamente.
O registro rejeita:

- ID ou versão vazios;
- IDs duplicados;
- kind desconhecido;
- locale ou mercado inválido/duplicado;
- conector que não declara unidade de identidade;
- `PublicOnly=false` nesta materialização;
- conector que promete mercado sem explicar se há filtro oficial;
- conector não replayable no braço avaliado.

O adaptador implementa uma porta equivalente a:

```go
type MarketDiscoveryConnector interface {
    Capabilities() DiscoveryConnectorCapabilities
    Discover(context.Context, MarketDiscoveryRequest) (ConnectorDiscovery, error)
}
```

L1 conhece somente a porta e contratos. L3 implementa rede e parsing. L4 monta
allowlist, orçamento e instâncias. L2 não importa conectores concretos.

## Ferramentas MCP

As sete ferramentas do Passo 0026 permanecem nas mesmas posições, com seus
observáveis preservados. Duas ferramentas são acrescentadas ao final:

8. `bateia_list_connectors`;
9. `bateia_start_market_discovery`.

Essa ampliação exige recibo explícito de refinamento dos certificados 0025 e
0026; alterar silenciosamente o número de tools invalida a cadeia.

### `bateia_list_connectors`

Não recebe argumentos. Retorna somente capacidades estáticas dos conectores
montados em L4:

```json
{
  "schema_version": "1",
  "connectors": [
    {
      "id": "public-application-catalog",
      "version": "v1",
      "kind": "APPLICATION_CATALOG",
      "query_locales": ["en-US", "ja-JP"],
      "markets": ["US", "JP"],
      "supports_market_filter": true,
      "identity_levels": ["provider_product_id", "canonical_url"],
      "evidence_surfaces": ["name", "summary", "publisher"],
      "replayable": true,
      "public_only": true
    }
  ]
}
```

Listar não testa rede, não consome quota e não revela token, endpoint privado,
path local ou configuração sensível. Ausência de credencial é projetada apenas
como indisponibilidade pública quando uma fonte permitida exigir autenticação.

### `bateia_start_market_discovery`

Entrada:

```json
{
  "idempotency_key": "market-task-2026-09-01-0001",
  "idea": {
    "problem": "reduzir desperdício alimentar doméstico",
    "users": ["famílias"],
    "capabilities": ["inventário", "previsão de consumo"]
  },
  "query_proposals": [
    {
      "locale": "ja-JP",
      "query": "食品在庫 消費予測",
      "source_text": "inventário de alimentos e previsão de consumo",
      "derivation": "EXTERNAL_AGENT_PROPOSAL",
      "rationale": "buscar formulação funcional usada no Japão"
    }
  ],
  "markets": ["BR", "JP"],
  "connector_ids": ["public-application-catalog"],
  "budget": {
    "max_queries": 20,
    "max_results_per_query": 10,
    "max_total_candidates": 100
  }
}
```

Regras:

- `idempotency_key` preserva integralmente o contrato do Passo 0026;
- `idea` usa a mesma normalização do `RecommendationService`;
- cada locale é BCP 47 suportado por ao menos um conector selecionado;
- mercados são ISO 3166-1 alpha-2 explícitos e canônicos;
- `connector_ids` aceita somente IDs registrados em L4;
- a query é texto semântico, não expressão crua da API do provedor;
- adaptadores escapam ou codificam sintaxe específica da fonte;
- cada proposta possui texto-fonte, derivação e rationale não vazios;
- query, source e rationale têm limites de tamanho e rejeitam controles;
- limites pertencem a faixas congeladas no contrato e são validados antes da
  fila ou rede;
- ao menos uma proposta válida ou uma tradução controlada deve existir;
- nenhuma proposta constitui prova de funcionalidade, tradução correta ou
  origem regional.

Resposta:

```json
{
  "schema_version": "1",
  "run_id": "run_opaque_identifier",
  "run_kind": "MARKET_DISCOVERY",
  "lifecycle": "QUEUED",
  "version": 1
}
```

O runtime de jobs, idempotência, store, cancelamento e reinício é o mesmo do
Passo 0026. Não pode existir uma segunda fila específica para mercados.

## Propostas de consulta e proveniência

```go
type MarketQueryDerivation string

const (
    QueryControlledCatalog MarketQueryDerivation = "CONTROLLED_CATALOG"
    QueryExternalAgent     MarketQueryDerivation = "EXTERNAL_AGENT_PROPOSAL"
)

type MarketQueryProposal struct {
    ID             string                `json:"id"`
    Locale         string                `json:"locale"`
    Query          string                `json:"query"`
    SourceText     string                `json:"source_text"`
    Derivation     MarketQueryDerivation `json:"derivation"`
    CatalogVersion string                `json:"catalog_version,omitempty"`
    RuleID         string                `json:"rule_id,omitempty"`
    Rationale      string                `json:"rationale"`
    InputHash      string                `json:"input_hash"`
}
```

Propostas controladas registram catálogo, regra e versão. Propostas externas
registram hash da entrada e a autoridade `EXTERNAL_AGENT`, mas não nome, prompt
privado ou credencial de modelo. A origem externa é sempre marcada como não
verificada.

O Bateia pode validar forma, cobertura de termos e segurança; não pode declarar
que uma tradução externa preserva semântica. Quando não existe tradução
controlada, o status da derivação permanece `UNVERIFIED_PROPOSAL` no artefato.

## Planejamento da matriz

O planejador puro produz combinações ordenadas de:

```text
query proposal × market × compatible connector
```

Para cada combinação registra:

- locale da query;
- mercado solicitado;
- conector e versão;
- se o conector aplica filtro oficial de mercado;
- query original e representação efetiva sanitizada;
- custo planejado;
- motivo de execução ou gap;
- hash causal do pedido.

Se o conector suporta o idioma, mas não o mercado, a combinação vira gap. Se não
possui filtro de mercado, pode executar somente como descoberta global e deve
registrar `MARKET_NOT_FILTERED`; o resultado não prova disponibilidade no mercado.

O orçamento é alocado de forma determinística e justa entre conectores antes da
execução. Ordem de registro ou retorno não pode dar toda a quota ao primeiro
conector.

## Identidade de produto

`RepositoryCandidate` não pode representar sozinho aplicações sem código
público. O domínio recebe identidade própria:

```go
type ProductIdentity struct {
    ConnectorID string `json:"connector_id"`
    ExternalID  string `json:"external_id"`
    CanonicalURL string `json:"canonical_url,omitempty"`
}

type ProductCandidate struct {
    Identity             ProductIdentity       `json:"identity"`
    Name                 string                `json:"name"`
    Summary              string                `json:"summary,omitempty"`
    ProductKind          string                `json:"product_kind"`
    Publisher            string                `json:"publisher,omitempty"`
    ListingLocales       []string              `json:"listing_locales"`
    ObservedMarkets      []string              `json:"observed_markets"`
    PublisherCountries   []string              `json:"publisher_countries"`
    RepositoryLinks      []ProductRepositoryLink `json:"repository_links"`
    Origins              []MarketCandidateOrigin `json:"origins"`
}
```

`ExternalID` é obrigatório e possui namespace do conector. URL isolada não é
identidade suficiente. Nome, resumo, ícone, publisher ou posição no ranking não
viram chave canônica.

Vínculo com repositório possui status `VERIFIED`, `UNVERIFIED` ou `CONFLICT` e
fonte. Somente `VERIFIED` pode ser encaminhado como candidato ao funil de
evidência existente. Produto sem repositório permanece no relatório; o Bateia
não inventa `owner/repo`, commit ou ausência de implementação.

## Resolução entre fontes

Dentro do mesmo conector, igualdade exige o mesmo `ExternalID` canônico. Entre
conectores, resolver como o mesmo produto exige ao menos uma testemunha forte:

- identificador oficial compartilhado;
- URL oficial recíproca;
- package/application ID declarado por ambas as fontes;
- repositório oficial verificado por ambas;
- vínculo publicado pelo próprio fornecedor e observado nas duas superfícies.

Nome semelhante, mesma descrição, mesmo ícone, mesmo publisher textual ou mesma
query não bastam. O resultado de resolução é `SAME`, `DISTINCT`, `POSSIBLE` ou
`UNKNOWN`, com testemunhas. `POSSIBLE` e `UNKNOWN` permanecem grupos separados.

Merge automático falso é falha crítica. Duplicata preservada é preferível a
fusão sem prova.

## Snapshot e artefatos duráveis

Cada tentativa registra conector, versão, proposta, locale, mercado, filtro
aplicado, status, contagens, classe pública de erro e hash da resposta congelada.
Timeout, rate-limit e resposta inválida são distintos de `EMPTY`.

`bateia_get_run` passa a incluir `run_kind`. Para `MARKET_DISCOVERY`, o campo
`market_discovery_status` aceita `COMPLETE`, `PARTIAL` ou `INCONCLUSIVE`; não é
projetado em `recommendation_status`.

`bateia_get_artifact` acrescenta à allowlist:

- `connector_capabilities`;
- `market_plan`;
- `market_snapshot`;
- `market_candidates`;
- `market_gaps`.

Hash, publicação atômica, paginação, cursor HMAC, quota, permissões e integridade
seguem o Passo 0026. O snapshot seguro não expõe corpo integral de páginas,
reviews, dados de usuário, HTML, token ou mensagens externas.

## Relação com recomendação e evidência

Este passo descobre produtos; não afirma que sejam alternativas comprovadas.

- metadado de catálogo pode sustentar `DISCOVERED` ou sinal parcial de produto;
- não sustenta `PROVEN` para capacidade implementada;
- popularidade, downloads, avaliações e ranking da loja não substituem fit;
- produto com repo `VERIFIED` pode alimentar posteriormente o funil existente;
- produto sem repo pode ser analisado pela IA como hipótese ou referência, sempre
  citando a superfície observada;
- nenhum candidato de mercado entra automaticamente em `top_choice`;
- a ausência de candidato repo-backed não apaga candidatos de produto.

A integração automática de candidatos de produto ao ranking de alternativas fica
fora deste passo porque exige política de evidência própria para software sem
código público.

## Corpus multilíngue congelado

Congelar antes da implementação seis ideias distribuídas por ao menos três
famílias de produto, sem nomes de concorrentes nas ideias ou queries. Usar splits
`CONTRACT`, `VALIDATION` e `TEST`, com duas ideias por split. O split `TEST`
permanece oculto do implementador até o selo.

A matriz primária contém:

| Query locale | Mercado |
|---|---|
| `pt-BR` | `BR` |
| `en-US` | `US` |
| `es-ES` | `ES` |
| `fr-FR` | `FR` |
| `de-DE` | `DE` |
| `ja-JP` | `JP` |
| `ko-KR` | `KR` |
| `ru-RU` | `RU` |
| `ar-SA` | `SA` |
| `hi-IN` | `IN` |

Adicionar ao menos quatro controles cruzados em que idioma e mercado não
coincidem, por exemplo uma consulta inglesa aplicada a mercado não anglófono.
Esses controles impedem implementação que derive mercado do sufixo do locale.

O registry de produtos esperados é fechado antes do patch, mas fica oculto do
implementador. Cada julgamento identifica fonte, ExternalID, mercado observado,
relevância funcional e confiança de correspondência. Casos sem identidade forte
permanecem opacos.

## Métricas e gates de produto

- 100% dos conectores declaram capacidades válidas antes da rede;
- 100% das tentativas e candidatos possuem proveniência completa;
- 100% dos gaps distinguem falta de capacidade, budget, erro e vazio observado;
- ao menos um conector configurado cobre cada par primário da matriz ou o passo
  falha antes do experimento;
- conclusão de queries públicas ≥ 80% por conector avaliado;
- precisão julgada no top 10 ≥ 70% por família de produto;
- recall de candidatos adequados ≥ 70% no conjunto `VALIDATION + TEST`;
- ganho marginal de recall do conector não repo-backed ≥ 20 pontos percentuais
  sobre o baseline repo-backed congelado;
- pelo menos quatro das seis ideias recuperam candidato relevante ausente no
  baseline repo-backed;
- precisão de merges automáticos = 100%; qualquer falso merge bloqueia promoção;
- taxa de inferência indevida de origem, idioma ou disponibilidade = 0%;
- produtos sem repo preservados = 100%;
- respostas congeladas reproduzem snapshot e consolidação byte a byte;
- popularidade nunca altera identidade, estado de evidência ou classe de gap;
- mutation score `1,0`.

Se exposição de rede ficar abaixo de 80%, recall e precisão não são julgados; o
veredito é inconclusivo por exposição, não falha automática da formulação.

## Política de `Unknown`

- locale sem query válida: `UNKNOWN_QUERY_DERIVATION`;
- locale ou mercado não suportado: `UNSUPPORTED_CAPABILITY`;
- conector sem filtro regional: `MARKET_NOT_FILTERED`;
- timeout, rate-limit ou falha do provedor: `UNKNOWN_PROVIDER`;
- consulta concluída sem resultados: `EMPTY_OBSERVED`, limitada à fonte e ao
  instante;
- idioma da listagem não declarado ou detecção ambígua: desconhecido;
- país do publicador ausente: desconhecido;
- produto visto num mercado: não prova origem nem disponibilidade global;
- tradução proposta pela IA: não verificada;
- produto sem repo: identidade de código desconhecida, não software ausente;
- repo sem commit imutável: não encaminhável à prova;
- alias entre fontes sem testemunha forte: `POSSIBLE` ou `UNKNOWN`, sem merge;
- budget esgotado: candidatos não observados permanecem desconhecidos;
- resposta congelada corrompida: erro de integridade, nunca `EMPTY`.

## Segurança, termos e conteúdo não confiável

- usar somente API pública, export permitido ou superfície cuja coleta tenha sido
  autorizada e congelada antes da implementação;
- não contornar autenticação, rate-limit, robots, paywall ou restrição regional;
- não coletar review individual, perfil, e-mail, identificador de usuário ou dado
  pessoal;
- não instalar, baixar nem executar aplicação ou binário encontrado;
- URLs retornadas são normalizadas e não são seguidas fora da allowlist do
  conector;
- nenhum campo do candidato controla ferramenta, prompt, política ou próxima
  chamada;
- descrições e nomes são dados não confiáveis, limitados e sanitizados;
- propostas externas não aceitam sintaxe de filtro, cabeçalho, path ou URL;
- erros públicos não ecoam resposta, token ou canário do adaptador;
- snapshots brutos permanecem fora da projeção MCP quando contêm conteúdo maior
  que o contrato seguro;
- credenciais, quando necessárias e já autorizadas, existem somente em L3/L4.

## Observáveis a congelar antes da implementação

1. as sete tools certificadas mantêm ordem, schema e comportamento;
2. `tools/list` acrescenta somente `bateia_list_connectors` e
   `bateia_start_market_discovery` ao final;
3. listar conectores não acessa rede e retorna capacidades canônicas;
4. conector incompatível falha antes da rede;
5. `start` responde `QUEUED` e reutiliza a mesma fila/store do Passo 0026;
6. idempotência inclui ideia, propostas, mercados, conectores e orçamento;
7. idioma, mercado, listing locale e publisher country nunca se copiam
   implicitamente entre campos;
8. controles cruzados não derivam mercado do locale;
9. proposta externa é executável como busca, mas permanece não verificada e não
   altera evidência;
10. adaptador recebe texto escapado, não expressão arbitrária do cliente;
11. budget é distribuído de forma independente da ordem dos conectores;
12. falha de um conector preserva resultados e gaps dos demais;
13. `EMPTY` só aparece após tentativa concluída com sucesso;
14. produto sem repo permanece no snapshot e no artefato;
15. repo não verificado nunca entra no funil de código;
16. homônimos não são fundidos; alias forte é fundido com testemunhas;
17. reordenação de conectores, propostas e respostas produz artefatos canônicos;
18. replay offline não chama rede e é byte a byte idêntico;
19. cancelamento impede novas chamadas em todos os conectores;
20. canários em descrição, erro, URL e proposta não atravessam artefatos ou logs;
21. o novo conector prova cobertura real da matriz e ganho marginal no corpus;
22. recommendation e durable runs preservam seus hashes ou recebem refinamento
    explícito e compatível.

## Mutações mínimas a rejeitar

- locale é usado como país do produto;
- mercado é derivado do sufixo BCP 47;
- listing locale vira idioma suportado pela aplicação sem fonte;
- TLD, moeda ou alfabeto vira publisher country;
- proposta da IA recebe status de tradução comprovada;
- texto de query é concatenado como filtro cru da API;
- conector não registrado é instanciado pela entrada do usuário;
- `bateia_list_connectors` faz health check de rede ou revela configuração;
- primeiro conector consome todo o orçamento por ordem incidental;
- erro/rate-limit vira `EMPTY`;
- resultados de um conector são descartados porque outro falhou;
- nome normalizado funde produtos de publishers diferentes;
- URL sozinha resolve identidade entre fontes;
- produto sem repositório desaparece;
- link repo não verificado é promovido a `owner/repo@commit`;
- metadado de loja produz capacidade `PROVEN`;
- estrelas, downloads ou avaliações promovem candidato;
- MARKET_DISCOVERY usa fila, store ou idempotência paralelos;
- snapshot de replay chama rede;
- artefato MCP contém review, HTML, corpo integral ou canário;
- tool nova é inserida antes das sete ferramentas congeladas;
- certificados 0025/0026 são tratados como preservados sem recibo de refinamento.

Mutation score exigido para selar o contrato: `1,0`.

## Gates de aceitação

- spec, baselines, matriz, conector escolhido, corpus, contrato, oráculos e
  mutações congelados antes do candidato;
- `crystalline-lint .` sem violações de gravidade;
- `go test -count=1 ./...` passa duas vezes;
- `go test -race` passa nos pacotes novos e nos durable runs afetados;
- `go build -o bateia ./04_wiring/main.go` passa;
- `go vet ./...` passa;
- testes de conectores usam fixtures congeladas; nenhuma suíte unitária exige
  rede real;
- execução pública autorizada respeita budget e gera receipt sem credenciais;
- replay offline repetido produz artefatos byte a byte idênticos;
- testes adversariais cobrem idioma/região, identidade, SSRF, injection, erro,
  cancelamento, budget e preservação das tools antigas;
- scan de snapshots e projeções não encontra segredo ou dado pessoal de usuário;
- métricas de produto satisfazem os gates no split oculto;
- mutation score `1,0` sem converter caso opaco em sucesso;
- receipts de refinamento demonstram compatibilidade observável com 0025 e 0026.

## Regime Tekt previsto

Aplicar protocolo completo. O usuário ocupa a autoria da intenção; a spec deste
passo e os certificados 0025/0026 formam a raiz causal. Autor do contrato recebe
spec, baselines e capacidades públicas da fonte selecionada, mas não o patch.
Autor dos oráculos recebe corpus sem implementação. Implementador não lê o split
`TEST`, registry esperado, mutações ou suites protegidas. Adversário pode permutar
idioma/mercado, adulterar identidades e respostas e simular falhas, mas não
corrigir o candidato. Verificador não escreve produção nem entradas seladas.

Qualquer troca de conector, matriz, corpus, política de merge ou fonte de
tradução invalida o selo a partir do primeiro artefato afetado.

Sem isolamento material forte, o laudo deve declarar
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`, ainda que todos os gates passem.

## Fora de escopo

- incorporar modelo de IA, embedding ou tradutor neural ao Bateia;
- declarar tradução proposta pelo agente como correta;
- integrar automaticamente produto sem repo ao ranking de alternativas;
- provar capacidade implementada somente com metadado de catálogo;
- extrair padrões arquiteturais ou gerar `ImprovementHypothesis`;
- adicionar mais de um novo conector não repo-backed nesta materialização;
- cobrir todas as lojas, países, idiomas ou tipos de software;
- coletar reviews, dados pessoais, telemetria ou conteúdo privado;
- instalar, executar, comprar ou modificar aplicações;
- inferir origem cultural, nacionalidade ou localização de pessoa;
- transporte MCP remoto ou autenticação multiusuário;
- resume de MARKET_DISCOVERY interrompido;
- retenção, deleção ou garbage collection de runs;
- alterar automaticamente política, query ou conector após observar o split
  oculto.

## Como pode falhar

- o conector escolhido não cobre a matriz ou não permite congelar respostas;
- resultados crescem, mas precisão e ganho marginal não passam;
- a mesma solução aparece várias vezes sem identidade resolvível;
- merges falsos elevam recall artificialmente;
- traduções externas produzem ruído e são tratadas como verificadas;
- mercado é confundido com origem do produto;
- produto sem repo é descartado para caber no modelo antigo;
- a nova execução duplica infraestrutura do Passo 0026;
- conteúdo de catálogo volta a ocupar o contexto integral da IA;
- o benchmark mede disponibilidade do provedor em vez de qualidade da descoberta;
- mudanças em tools e artefatos quebram clientes certificados sem recibo.

## Critério de conclusão

O passo termina quando uma IA externa consegue listar conectores, iniciar uma
descoberta durável com consultas em múltiplos idiomas e mercados, recuperar
produtos repo-backed e não repo-backed com identidade e proveniência, observar
gaps honestos e reproduzir o snapshot offline; e quando o conector novo demonstra
ganho marginal de recall no corpus oculto sem falso merge, inferência regional,
promoção indevida de evidência ou regressão silenciosa dos Passos 0025 e 0026.

O veredito positivo é `MULTISOURCE_MARKET_DISCOVERY_READY`. Ele vale somente para
o conector, matriz, corpus, versões e hashes certificados. Não implica cobertura
mundial, tradução correta para ideias arbitrárias nem recomendação comprovada de
produtos sem código público.
