# Contrato observável — Passo 0027

**Revisão:** 1  
**Estado:** candidato a selo, anterior à implementação  
**Regime:** protocolo Tekt completo; `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Veredito positivo reservado:** `MULTISOURCE_MARKET_DISCOVERY_READY`

`DEVE`, `NÃO DEVE`, `EXATAMENTE` e `SOMENTE` são normativos. Este contrato cobre
somente o fragmento observável aqui descrito. Não autoriza emitir o veredito reservado.

## 1. Raiz causal e escopo

| Entrada | SHA-256 |
|---|---|
| `00_nucleo/passos/0027-generalizar-descoberta-multilingue-e-multifuente.md` | `83f80a9603bc566690ddfb98a94ca18f4bd57b528a177c56d4ba5873b74f4c15` |
| `00_nucleo/benchmark/0027/manifest.json` | `7661f40f1b1b01cf0c99196bec9f81f7c462c9a9c529888c54dea8485844c2a3` |
| `00_nucleo/benchmark/0027/source-selection.json` | `cd833b03cd84a933b587da9a034fdb6d10c023fee423049a27e41de9d0539460` |
| `00_nucleo/ADR/ADR-0002-estrutura-de-pastas-e-camadas-tekt.md` | `b457fadd55d6328c2afce969b92a527b461edab9094b1bef66fcfebe71b79a68` |
| certificado 0025 | `bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d` |
| receipt de verificação 0025 | `a49e2733cbf47f03a7f537b8ec4efd548333302d36ad877a4422344f95f4bf4a` |
| certificado 0026 | `a04d6b07a939c1ce5fb7fb184e5cb56cc3cc34ed0b2b54c8087415a18a4a99ed` |
| receipt de verificação 0026 | `76c0c20aa66d00252a96d1cd78a78d0faa519b0113c63ff187d94dad3cff38d2` |

O Passo 0027 adiciona descoberta durável de produtos e EXATAMENTE um conector
não repo-backed: `apple-itunes-search-software`. Não adiciona IA, tradutor,
ranking automático de produtos, segunda fila, segundo store nem segundo domínio
de idempotência.

## 2. Refinamento obrigatório de 0025 e 0026

Os primeiros sete elementos de `tools/list`, seus nomes, descrições e schemas
devem ser deep-equal à superfície certificada pelos hashes abaixo:

| Superfície basal | SHA-256 |
|---|---|
| `02_shell/mcp/server.go` (tools 1–3) | `cbd0188fadeb982cea910fce2392b4ce4b384dc7e4a00fa53a3a3434de8fa616` |
| `02_shell/mcp/durable_runs.go` (tools 4–7) | `a4f461bcb7a37fe4f35527e0e90b3f69b1bfc57bfe1816918e3285c958808e83` |
| `01_core/domain/recommendation_run.go` | `b17fbc504485b4620f884f4d84c01fac1120fc835664ff0a9821d3ee515a639d` |
| `01_core/usecase/recommendation_run.go` | `fcbf5145ced2b39e3ff884e43b2fd1bf85bdd13d4a3d62aed6df3aa9c604f57f` |

Mudança necessária nos tipos públicos de run, na allowlist de artefatos ou na
fiação DEVE gerar receipts de refinamento separados para `MCP_RECOMMENDATION_GATEWAY_READY`
e `MCP_DURABLE_RUNS_READY`. Preservação não pode ser alegada apenas porque testes
novos passam. `bateia_recommend` e `bateia_start_recommendation` continuam usando
o mesmo `RecommendationService`; todos os observáveis de lifecycle, CAS,
cancelamento, cursor HMAC, quota, recovery, publicação atômica e erros públicos
de 0026 permanecem.

## 3. Ordem e schemas MCP

`tools/list` retorna EXATAMENTE nove tools nesta ordem:

1. `bateia_search`
2. `bateia_benchmark`
3. `bateia_recommend`
4. `bateia_start_recommendation`
5. `bateia_get_run`
6. `bateia_get_artifact`
7. `bateia_cancel_run`
8. `bateia_list_connectors`
9. `bateia_start_market_discovery`

As tools 1–7 são deep-equal ao baseline certificado. Tool 8 usa exatamente:

```json
{"type":"object","additionalProperties":false,"properties":{}}
```

Tool 9 usa exatamente o schema seguinte. Validação runtime rejeita unknown fields,
JSON múltiplo, controles, duplicatas depois de normalização e valores fora dos
limites antes de store, fila ou rede.

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "idempotency_key": {
      "type": "string",
      "minLength": 1,
      "maxLength": 128,
      "pattern": "^[\\u0020-\\u002E\\u0030-\\u005B\\u005D-\\u007E]+$"
    },
    "idea": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "problem": {"type":"string","minLength":1,"maxLength":1000},
        "users": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200}},
        "capabilities": {"type":"array","minItems":1,"maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200}},
        "constraints": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200}},
        "non_goals": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200}},
        "preferred_stack": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200}}
      },
      "required": ["problem","capabilities"]
    },
    "query_proposals": {
      "type": "array",
      "minItems": 1,
      "maxItems": 20,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "locale": {"type":"string","enum":["ar-SA","de-DE","en-US","es-ES","fr-FR","hi-IN","ja-JP","ko-KR","pt-BR","ru-RU"]},
          "query": {"type":"string","minLength":1,"maxLength":256,"pattern":"^[^\\u0000-\\u001F\\u007F]+$"},
          "source_text": {"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F]+$"},
          "derivation": {"type":"string","const":"EXTERNAL_AGENT_PROPOSAL"},
          "rationale": {"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F]+$"}
        },
        "required": ["locale","query","source_text","derivation","rationale"]
      }
    },
    "markets": {
      "type": "array",
      "minItems": 1,
      "maxItems": 10,
      "uniqueItems": true,
      "items": {"type":"string","enum":["BR","DE","ES","FR","IN","JP","KR","RU","SA","US"]}
    },
    "connector_ids": {
      "type": "array",
      "minItems": 1,
      "maxItems": 3,
      "uniqueItems": true,
      "items": {"type":"string","enum":["apple-itunes-search-software","github-repository-catalog","sourcegraph-repository-catalog"]}
    },
    "budget": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "max_queries": {"type":"integer","minimum":1,"maximum":20},
        "max_results_per_query": {"type":"integer","minimum":1,"maximum":10},
        "max_total_candidates": {"type":"integer","minimum":1,"maximum":100}
      },
      "required": ["max_queries","max_results_per_query","max_total_candidates"]
    }
  },
  "required": ["idempotency_key","idea","query_proposals","markets","connector_ids","budget"]
}
```

Tool 8 não acessa rede e responde num único item textual MCP contendo JSON compacto:

```json
{"schema_version":"1","connectors":[{"id":"apple-itunes-search-software","version":"itunes-search-doc-1.0.3-response-v2","kind":"APPLICATION_CATALOG","query_locales":["ar-SA","de-DE","en-US","es-ES","fr-FR","hi-IN","ja-JP","ko-KR","pt-BR","ru-RU"],"markets":["BR","DE","ES","FR","IN","JP","KR","RU","SA","US"],"supports_market_filter":true,"identity_levels":["provider_product_id","canonical_url"],"evidence_surfaces":["name","summary","publisher","provider_declared_language_codes","storefront_observation"],"replayable":true,"public_only":true},{"id":"github-repository-catalog","version":"baseline-21c48e357180","kind":"REPOSITORY_CATALOG","query_locales":["ar-SA","de-DE","en-US","es-ES","fr-FR","hi-IN","ja-JP","ko-KR","pt-BR","ru-RU"],"markets":[],"supports_market_filter":false,"identity_levels":["owner_repo","commit"],"evidence_surfaces":["repository_identity","immutable_commit","code_occurrence"],"replayable":true,"public_only":true},{"id":"sourcegraph-repository-catalog","version":"baseline-21c48e357180","kind":"REPOSITORY_CATALOG","query_locales":["ar-SA","de-DE","en-US","es-ES","fr-FR","hi-IN","ja-JP","ko-KR","pt-BR","ru-RU"],"markets":[],"supports_market_filter":false,"identity_levels":["owner_repo","commit"],"evidence_surfaces":["repository_identity","immutable_commit","code_occurrence"],"replayable":true,"public_only":true}]}
```

O array é ordenado por `id`. O suporte textual dos conectores repo-backed aos dez
locales não constitui filtro regional.

Tool 9 responde:

```json
{"schema_version":"1","run_id":"run_00000000000000000000000000000000","run_kind":"MARKET_DISCOVERY","lifecycle":"QUEUED","version":1}
```

`run_id` segue `^run_[0-9a-f]{32}$`; o valor do exemplo não é literal.

`bateia_get_artifact` preserva seu schema exceto pela enum aditiva, que passa a
ser exatamente `request`, `recommendation`, `decision`, `alternatives`,
`failures`, `provenance`, `connector_capabilities`, `market_plan`,
`market_snapshot`, `market_candidates`, `market_gaps`, nesta ordem.

## 4. API pública L1

L1 DEVE expor, sem I/O, os tipos e valores seguintes (nomes e strings exatos):

- `RunKind`: `RECOMMENDATION`, `MARKET_DISCOVERY`;
- `MarketDiscoveryStatus`: `COMPLETE`, `PARTIAL`, `INCONCLUSIVE`;
- estágios aditivos: `market_plan`, `market_discover`, `market_resolve`,
  `market_publish`; estados 0026 são preservados;
- `DiscoveryConnectorKind`: `REPOSITORY_CATALOG`, `APPLICATION_CATALOG`,
  `PRODUCT_CATALOG`;
- `MarketQueryDerivation`: `CONTROLLED_CATALOG`, `EXTERNAL_AGENT_PROPOSAL`;
- `QueryDerivationStatus`: `CONTROLLED`, `UNVERIFIED_PROPOSAL`;
- `RepositoryLinkStatus`: `VERIFIED`, `UNVERIFIED`, `CONFLICT`;
- `IdentityResolution`: `SAME`, `DISTINCT`, `POSSIBLE`, `UNKNOWN`;
- `MarketAttemptStatus`: `PLANNED`, `SUCCEEDED`, `EMPTY_OBSERVED`,
  `UNKNOWN_QUERY_DERIVATION`, `UNSUPPORTED_CAPABILITY`, `MARKET_NOT_FILTERED`,
  `UNKNOWN_PROVIDER`, `BUDGET_EXHAUSTED`, `CANCELLED`, `INTEGRITY_ERROR`.

`RecommendationRun` e sua projeção pública recebem `run_kind`. Para
`RECOMMENDATION`, `recommendation_status` mantém a semântica anterior e
`market_discovery_status` é `null`; para `MARKET_DISCOVERY`, ocorre o inverso.
Lifecycle, versão CAS e formato de `RunID` são compartilhados.

As estruturas públicas mínimas, com tags JSON homônimas, são:

```go
type DiscoveryConnectorCapabilities struct {
    ID string; Version string; Kind DiscoveryConnectorKind
    QueryLocales []string; Markets []string; SupportsMarketFilter bool
    IdentityLevels []string; EvidenceSurfaces []string
    Replayable bool; PublicOnly bool
}
type ProductIdentity struct { ConnectorID string; ExternalID string; CanonicalURL string }
type ProductRepositoryLink struct { OwnerRepo string; Commit string; Status RepositoryLinkStatus; Source string }
type MarketCandidateOrigin struct { ConnectorID string; ConnectorVersion string; QueryID string; QueryLocale string; ObservedMarket string; ResponseHash string }
type ProductCandidate struct {
    Identity ProductIdentity; Name string; Summary string; ProductKind string; Publisher string
    ListingLocales []string; ObservedMarkets []string; PublisherCountries []string
    RepositoryLinks []ProductRepositoryLink; Origins []MarketCandidateOrigin
}
```

Tags JSON são snake_case; slices públicas nunca são `null`. `ExternalID` é
obrigatório e namespaced por conector. L1 expõe `MarketDiscoveryConnector` com
`Capabilities()` sem rede e `Discover(context.Context, MarketDiscoveryRequest)`;
um registry puro valida e resolve somente IDs montados. L1 expõe operações puras
de validação, hash canônico, planejamento, distribuição justa de budget,
consolidação, resolução de identidade e projeção dos cinco artefatos.

O plano é o produto canônico de
`query_proposal × market × compatible_connector`, ordenado por
`query_id, market, connector_id`. IDs de proposta e tentativa são SHA-256
causais, não contadores de ordem. Reordenar entradas ou respostas não altera
hashes, plano ou artefatos. O budget é alocado round-robin sobre a ordem canônica;
nenhum conector recebe segunda unidade antes de todos os elegíveis receberem a
primeira.

## 5. Runtime durável único

`MARKET_DISCOVERY` usa EXATAMENTE as mesmas instâncias de
`RecommendationRunStore`, `RecommendationJobQueue`, gerador de run ID, raiz,
cursor secret, quota, workers, recovery e shutdown de 0026. É proibido criar
segunda fila, store, raiz, tabela/índice de idempotência ou chave de cursor.

O hash de idempotência inclui ideia normalizada, propostas normalizadas,
mercados e conectores ordenados, budget e versões de capacidades. Mesma chave e
mesmo hash reutilizam `run_id`; mesma chave com hash diferente produz
`IDEMPOTENCY_CONFLICT`. O dispatcher único lê `run_kind` do registro persistido e
encaminha ao executor correspondente. Recovery converte qualquer run não terminal,
inclusive de mercado, em `INTERRUPTED`; resume permanece fora de escopo.

`COMPLETE` exige todas as tentativas selecionadas concluídas e artefatos íntegros.
`PARTIAL` exige ao menos uma tentativa concluída e ao menos um gap/erro/budget.
`INCONCLUSIVE` cobre exposição inferior a 80%, nenhuma tentativa concluída,
integridade desconhecida ou apenas outcomes Unknown. Nenhum desses status altera
lifecycle: um run pode ser `SUCCEEDED + INCONCLUSIVE`.

Os erros MCP públicos de 0026 permanecem exatos (`-32020` a `-32027`). Entrada
inválida, conector inexistente, incompatibilidade integral ou orçamento incoerente
falha como `-32602 Invalid params` antes de I/O. Falhas por tentativa são gaps
controlados, não mensagens arbitrárias do provedor.

## 6. Contrato L3 Apple e replay

O adaptador Apple obedece integralmente `source-selection.json`: somente
`GET https://itunes.apple.com/search`, sem redirect, credencial, host/path/URL do
usuário ou DNS override; parâmetros permitidos são somente `term`, `country`,
`media=software`, `entity=software`, `limit` 1–10, `lang` em `en_us|ja_jp`,
`version=2`, `explicit=No`; `callback` e qualquer outro parâmetro são proibidos.

Limites: uma chamada concorrente, no máximo 15 chamadas/minuto, resposta até
2 MiB, strict JSON UTF-8. `429`, timeout, non-2xx, oversized ou JSON inválido são
`UNKNOWN_PROVIDER`, nunca vazio. `EMPTY_OBSERVED` exige 2xx válido com
`resultCount=0` e `results=[]`. Cancelamento observado impede nova chamada ou
retry. O texto de `term` é percent-encoded como um único valor.

`trackId` é a identidade Apple; item sem ele é rejeitado. `trackViewUrl` só é
aceito como testemunha se for HTTPS `apps.apple.com`, contiver `/id<trackId>` e
tiver query/fragment removidos; nunca é seguido. `country` registra apenas o
storefront observado. Para Apple, `listing_locales=[]` e
`publisher_countries=[]`; `lang` e `languageCodesISO2A` não preenchem esses campos.

O capture retém somente a projeção allowlisted descrita em
`source-selection.json`, metadados da tentativa e hash do body; artefatos
promocionais e body integral são descartados. O adapter de replay recebe essa
projeção congelada e NUNCA possui transporte de rede. Hash ou versão divergente é
`INTEGRITY_ERROR`. Duas execuções de replay produzem bytes idênticos para snapshot,
candidatos e gaps.

## 7. Identidade e merge

Dentro de um conector, somente `(connector_id, external_id)` igual permite merge.
Entre conectores, `SAME` exige ao menos uma testemunha forte bilateral:
identificador oficial compartilhado, URL oficial recíproca, package/application
ID declarado em ambas, repo oficial `owner/repo@commit` verificado em ambas, ou
vínculo do fornecedor observado nas duas superfícies. URL isolada, nome, resumo,
ícone, publisher textual, idioma, mercado, ranking ou query nunca bastam.

`POSSIBLE` e `UNKNOWN` permanecem registros separados. Merge automático falso
bloqueia promoção. Vínculo repo só entra no funil de evidência quando
`VERIFIED` e contém commit imutável; `UNVERIFIED` e `CONFLICT` ficam no produto.
Produto sem repo é preservado. Metadado de catálogo tem teto `DISCOVERED` e nunca
gera `PROVEN`, `top_choice` ou mudança de fit.

## 8. Artefatos e safe content

Todos usam `schema_version:"1"`, publicação atômica, hash do documento canônico,
paginação/cursor/quota de 0026 e arrays vazios. Ordenação:

| Artefato | Unidade/ordem | Campos obrigatórios por entrada |
|---|---|---|
| `connector_capabilities` | conector por `id` | todos os campos de capabilities |
| `market_plan` | tentativa por `query_id,market,connector_id` | `attempt_id`, `query_id`, `query_locale`, `market`, `connector_id`, `connector_version`, `market_filter_applied`, `query_hash`, `effective_query_hash`, `planned_cost`, `status`, `reason`, `request_hash` |
| `market_snapshot` | tentativa por `attempt_id` | campos do plano mais `result_count`, `accepted_count`, `error_class`, `response_hash`, `observed_at` |
| `market_candidates` | produto por `connector_id,external_id` | `identity`, `name`, `product_kind`, `publisher`, quatro arrays dimensionais, `repository_links`, `origins`, `resolution` e `witnesses` |
| `market_gaps` | gap por `attempt_id,status` | `attempt_id`, `query_locale`, `market`, `connector_id`, `status`, `reason_code`, `response_hash` |

`market_plan` não expõe query, source text ou rationale: somente hashes.
`market_snapshot` não expõe body, HTML, review, perfil, e-mail, screenshot,
artwork, preço, rating, device list, token, erro externo ou configuração. Para o
conector Apple, `summary` da projeção MCP fica vazio; a descrição participa apenas
da projeção congelada protegida e de seu hash. Nome/publisher são NFC, whitespace
colapsado, sem controles, limitados a 200 caracteres; URLs a 2048; erro público é
enum controlada. Query, descrição, erro ou URL não controlam prompt, ferramenta,
política, path, header, host ou chamada subsequente.

## 9. Quatro dimensões independentes e Unknown

- `query_locale`: BCP 47 da proposta; não é atributo do produto;
- `listing_locale`: somente declaração explícita ou detector determinístico
  versionado; caso Apple, vazio;
- `observed_market`: country explícito de tentativa bem-sucedida; não é origem nem
  disponibilidade global;
- `publisher_country`: somente declaração explícita da fonte; caso Apple, vazio.

É proibido copiar ou derivar qualquer dimensão das outras, inclusive por sufixo
BCP 47, alfabeto, TLD, moeda, nome ou ranking. Policies exatas:

| Situação | Outcome |
|---|---|
| locale sem proposta válida | `UNKNOWN_QUERY_DERIVATION` |
| locale/mercado/conector incompatível | `UNSUPPORTED_CAPABILITY` |
| conector global sem filtro | `MARKET_NOT_FILTERED` |
| timeout/rate-limit/provider/parser | `UNKNOWN_PROVIDER` |
| 2xx válido e zero resultados | `EMPTY_OBSERVED` |
| budget sem tentativa | `BUDGET_EXHAUSTED` |
| fixture/hash inválido | `INTEGRITY_ERROR` |
| alias sem testemunha forte | `POSSIBLE` ou `UNKNOWN`, sem merge |

`Unknown` nunca conta como sucesso, ausência, candidato negativo ou mutação
rejeitada, exceto em caso explicitamente opaco.

## 10. Corpus e métricas

O corpus tem EXATAMENTE seis ideias, três famílias e dois itens por split:

| ID | Split | Família | Conteúdo visível ao implementador |
|---|---|---|---|
| `C01` | `CONTRACT` | gestão de recursos domésticos | problema: reduzir desperdício alimentar doméstico; capacidades: inventário, validade, previsão de consumo |
| `C02` | `CONTRACT` | conhecimento e coordenação de equipes | problema: preservar decisões e handoffs de equipe; capacidades: registro de decisão, responsabilidade, acompanhamento |
| `V01` | `VALIDATION` | manutenção de ativos | problema: coordenar manutenção de equipamentos compartilhados; capacidades: cadastro, agenda recorrente, histórico |
| `V02` | `VALIDATION` | gestão de recursos domésticos | problema: acompanhar consumo doméstico de energia; capacidades: leituras, orçamento, alerta de anomalia |
| `T01` | `TEST` | conhecimento e coordenação de equipes | oculto; hash selado antes do candidato |
| `T02` | `TEST` | manutenção de ativos | oculto; hash selado antes do candidato |

Ideias e queries não contêm nomes de concorrentes. O autor de oráculos congela o
conteúdo integral, queries, expected registry e hashes de `T01/T02` antes do
candidato; o implementador recebe somente a tabela acima e nunca o split TEST ou
registry esperado.

A matriz primária é exatamente `pt-BR/BR`, `en-US/US`, `es-ES/ES`, `fr-FR/FR`,
`de-DE/DE`, `ja-JP/JP`, `ko-KR/KR`, `ru-RU/RU`, `ar-SA/SA`, `hi-IN/IN`.
Controles cruzados mínimos e obrigatórios: `en-US/JP`, `ja-JP/US`, `pt-BR/DE`,
`de-DE/BR`. Locale e mercado são entradas independentes.

Métricas usam frações inteiras, sem arredondamento:

- capabilities válidas antes da rede, proveniência e gaps classificados: 100%;
- exposição por conector = respostas 2xx strict JSON / tentativas públicas
  selecionadas após capability e budget; mínimo 80%;
- se exposição <80%, precisão/recall/gain são `INCONCLUSIVE_EXPOSURE`, não fail;
- precisão top 10 por família ≥70%;
- recall de candidatos adequados em `VALIDATION+TEST` ≥70%;
- ganho marginal de recall Apple sobre baseline repo-backed ≥20 pontos percentuais;
- pelo menos 4/6 ideias recuperam candidato relevante ausente do baseline;
- precisão de merge automático 100%; inferência dimensional indevida 0%;
- produtos sem repo preservados 100%; replay byte a byte 100%;
- mutation score = mutações válidas rejeitadas / mutações válidas = `1.0`;
- popularidade, ratings e downloads alteram identidade/evidência/gap em 0 casos.

## 11. Camadas e composição

- L1: tipos, portas, validação, registry, plano, budget, identidade, merge,
  idempotência, status e projeções; zero rede/filesystem/env/goroutine/mutex.
- L2: schemas/decoding/envelopes MCP e apresentação; depende somente de portas L1,
  nunca de L3/L4.
- L3: Apple HTTP, parsing/capture, replay, store e worker; importa somente L1,
  nunca L2/L4.
- L4: monta os três conectores, dispatcher, o mesmo store/queue/idempotência e
  injeta portas; zero regra de negócio.

## 12. Artefatos esperados e lineage única

| Fase | Artefato | Autoridade |
|---|---|---|
| raiz | passo, manifest, source selection, este contrato e prompts L0 | intenção/contrato |
| oráculo | corpus/registry/fixtures congelados, oracle contract e testes independentes | autor de oráculos |
| ataque | casos adversariais | adversário |
| discriminação | receipt 1.0 e selo | verificador |
| candidato | produção e receipt do candidato | implementador/coordenador |
| refinamento | receipts explícitos 0025 e 0026 | verificador |
| final | verification receipt, certificado e laudo | verificador |

Cada arquivo de produção possui exatamente um prompt espelhado conforme ADR-0002.
Prompts não são testes nem alternativa ao contrato. Alterar passo, source selection,
corpus, matriz, merge, tradução, endpoint, contrato ou prompt depois do selo
invalida a cadeia a partir da primeira entrada afetada.

## 13. Critérios executáveis

Antes de selo: JSON snippets parseiam estritamente; links locais existem; hashes
da seção 1 conferem; casos positivos são `Preserved`, mutações válidas são
`Violated`, opacos são `Unknown`; mutation score `1.0`.

Depois do candidato, todos devem passar:

```text
go test -count=1 ./...                         # duas execuções separadas
go test -race -count=1 ./01_core/... ./02_shell/mcp ./03_infra/... ./04_wiring
go vet ./...
go build -o /tmp/bateia-step0027 ./04_wiring
crystalline-lint .
git diff --check
```

Além disso: `tools/list` deep-equal para posições 1–7 e exatas posições 8–9;
nenhuma rede em listagem ou testes unitários; incompatibilidade antes da rede;
mesmo store/queue/root/secret por identidade de instância; quatro controles
cruzados; cancelamento sem novas chamadas; replay duas vezes comparado por bytes;
scan de artefatos/logs sem canários, segredo, body, review ou dado pessoal;
métricas da seção 10; hashes protegidos intactos; receipts de refinamento 0025/0026.

Somente então o verificador pode emitir `MULTISOURCE_MARKET_DISCOVERY_READY`,
limitado aos hashes, conector, matriz, corpus e versões deste contrato.
