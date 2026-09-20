# Contrato observável 0025 — recomendação completa por MCP

**Regime**: protocolo completo Tekt  
**Papel deste artefato**: Autor do Contrato  
**Spec congelada**: `00_nucleo/passos/0025-expor-recomendacao-completa-como-ferramenta-mcp.md`  
**SHA-256 da spec**: `8ddd8785f874ed36df41fbcedfa3fc580c72740dee1adeb1449ba7037f232a6f`  
**Baseline congelado**: `21c48e3571807777013710bcb6675c95c4130d9e`  
**Implementação candidata lida pelo autor**: não  

## 1. Autoridade e escopo

Este contrato fixa somente a superfície observável do corte síncrono `live`.
Ele não autoriza persistência, `run_id`, replay, resume, execução local, diretórios
temporários, chamadas a modelos de IA nem execução de conteúdo de candidatos.

O Autor do Contrato pôde ler a spec, o baseline e os contratos existentes de
domínio, use case, CLI, MCP e wiring. Pôde escrever apenas este contrato e prompts
L0 dos componentes de produção necessários. Não pôde ler patch candidato nem
escrever Go, testes, mutações, manifesto, selo, recibo ou laudo.

As obrigações abaixo são observáveis. Nomes, campos e funções explicitamente
declarados como públicos são estáveis; composição interna, helpers, arquivos
privados e estratégia de alocação não são congelados.

## 2. Porta pública mínima de L1

L1 deve oferecer os seguintes contratos de dados, com tags JSON indicadas:

```go
type RecommendationSearch struct {
    Strategy      string   `json:"strategy"`
    Language      string   `json:"language,omitempty"`
    Locales       []string `json:"locales"`
    MaxPerQuery   int      `json:"max_per_query"`
    MaxCandidates int      `json:"max_candidates"`
}

type RecommendationRequest struct {
    Idea   IdeaSpec                  `json:"idea"`
    Policy ReuseDecisionPolicy       `json:"policy"`
    Target SoftwareAlternativeTarget `json:"target"`
    Search RecommendationSearch      `json:"search"`
}

type RecommendationSummary struct {
    ExpansionCount int `json:"expansion_count"`
    CandidateCount int `json:"candidate_count"`
    AssessedCount  int `json:"assessed_count"`
}

type RecommendationResult struct {
    SchemaVersion string                    `json:"schema_version"`
    Status        IdeaRunStatus             `json:"status"`
    Summary       RecommendationSummary     `json:"summary"`
    Decision      ReuseDecisionReport       `json:"decision"`
    Alternatives  SoftwareAlternativeReport `json:"alternatives"`
    Failures      []ExecutionFailure        `json:"failures"`
    Warnings      []string                  `json:"warnings"`
}
```

O caso de uso oferece uma operação equivalente a:

```go
type RecommendationService struct { /* somente portas permitidas em L1 */ }

func NewRecommendationService(
    semantic SemanticExpander,
    discovery domain.ObservableDiscoveryProvider,
    content domain.ContentProvider,
) *RecommendationService

func (s *RecommendationService) Recommend(
    context.Context,
    domain.RecommendationRequest,
) (domain.RecommendationResult, error)
```

O adaptador MCP depende somente da forma estrutural mínima:

```go
type RecommendationRecommender interface {
    Recommend(context.Context, domain.RecommendationRequest) (domain.RecommendationResult, error)
}
```

O nome exato do campo privado que guarda essa porta não é contrato. Para preservar
compatibilidade de fonte dos chamadores atuais, `NewMCPServer(funnel, reader,
writer)` continua válido. A injeção opcional do recommender pode ser acrescentada
ao construtor de modo compatível (por exemplo, argumento variádico) ou por um
construtor adicional; a escolha privada não é congelada. O servidor montado por
L4, porém, **deve** receber uma instância funcional do serviço.

`RecommendationService.Recommend` valida e normaliza antes da primeira chamada a
`SemanticExpander`, `ObservableDiscoveryProvider` ou `ContentProvider`; executa o
modo `live`; compõe `RunIdeaPipeline` e `RankSoftwareAlternatives`; e não reimplementa
expansão, descoberta, assessment, decisão, licença ou ranking. A exigência de
`output_directory` é responsabilidade da CLI e não pode impedir a porta pura de
L1. L1 não cria, lê nem escreve arquivos.

## 3. Lista MCP e schema exato

`tools/list` anuncia, nesta ordem, `bateia_search`, `bateia_benchmark` e
`bateia_recommend`. Nome, descrição, schema, resultados e erros públicos das duas
ferramentas antigas permanecem iguais aos do baseline. A nova entrada anunciada é
exatamente:

```json
{
  "name": "bateia_recommend",
  "description": "Transforma uma ideia em decisão de reúso e alternativas de software com evidências rastreáveis",
  "inputSchema": {
    "type": "object",
    "additionalProperties": false,
    "properties": {
      "idea": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "problem": { "type": "string", "minLength": 1 },
          "users": { "type": "array", "items": { "type": "string" } },
          "capabilities": {
            "type": "array",
            "minItems": 1,
            "items": { "type": "string", "minLength": 1 }
          },
          "constraints": { "type": "array", "items": { "type": "string" } },
          "non_goals": { "type": "array", "items": { "type": "string" } },
          "preferred_stack": { "type": "array", "items": { "type": "string" } }
        },
        "required": ["problem", "capabilities"]
      },
      "policy": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "schema_version": { "type": "string", "const": "1" },
          "allowed_licenses": { "type": "array", "items": { "type": "string", "minLength": 1 } },
          "require_local": { "type": "boolean" },
          "allow_archived": { "type": "boolean" },
          "max_components": { "type": "integer", "minimum": 2 },
          "extend_min_coverage": {
            "type": "number",
            "exclusiveMinimum": 0,
            "maximum": 1
          },
          "integration_min_value": {
            "type": "number",
            "exclusiveMinimum": 0,
            "maximum": 1
          }
        },
        "required": [
          "schema_version",
          "allowed_licenses",
          "require_local",
          "allow_archived",
          "max_components",
          "extend_min_coverage",
          "integration_min_value"
        ]
      },
      "target": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "name": { "type": "string" }
        }
      },
      "search": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "strategy": { "type": "string", "enum": ["literal", "lexical"] },
          "language": { "type": "string" },
          "locales": {
            "type": "array",
            "uniqueItems": true,
            "items": {
              "type": "string",
              "enum": [
                "pt-BR", "en-US", "es-ES", "fr-FR", "de-DE",
                "ja-JP", "ko-KR", "ru-RU", "ar-SA", "hi-IN"
              ]
            }
          },
          "max_per_query": { "type": "integer", "minimum": 1, "maximum": 25 },
          "max_candidates": { "type": "integer", "minimum": 1, "maximum": 20 }
        },
        "required": ["strategy", "locales", "max_per_query", "max_candidates"]
      }
    },
    "required": ["idea", "policy", "target", "search"]
  }
}
```

O schema deliberadamente não expõe `provider`: L4 seleciona e injeta o provedor
permitido. Também não existem campos para token, cabeçalho, cookie, comando, URL,
`mode`, caminho local, `output_directory`, `replay_directory`, `resume` ou `run_id`.
`additionalProperties: false` vale em todos os objetos da entrada; o decoder deve
aplicar a mesma rejeição em runtime, pois o schema anunciado não substitui a
validação do servidor.

## 4. Validação e normalização

Antes de qualquer dependência externa:

1. `idea.problem`, após colapsar whitespace Unicode nas bordas e no interior, não
   pode ser vazio.
2. `idea.capabilities` deve conter ao menos um item não vazio após a mesma
   normalização. Listas de `IdeaSpec` usam a normalização já definida por L1:
   removem vazios, deduplicam sem diferenciar caixa e preservam a primeira grafia.
3. `policy.schema_version` deve ser exatamente `"1"`;
   `max_components >= 2`; os dois limiares pertencem a `(0,1]`; e
   `integration_min_value <= extend_min_coverage`.
4. `allowed_licenses` remove vazios, colapsa whitespace, deduplica sem diferenciar
   caixa e assume a ordenação canônica usada por `RankSoftwareAlternatives`.
5. `target.name` vazio após normalização recebe o `idea.problem` já normalizado.
6. `search.strategy` aceita no baseline somente `literal` ou `lexical`, em caixa
   exata conforme o schema. `lexical` exige expansor injetado. Estratégia ausente,
   desconhecida ou não configurada é inválida.
7. `search.locales` aceita somente as dez tags anunciadas. Duplicata, tag em outra
   caixa, `all` ou locale desconhecido é inválido; não há fallback para inglês.
   A ordem canônica é `pt-BR`, `en-US`, `es-ES`, `fr-FR`, `de-DE`, `ja-JP`,
   `ko-KR`, `ru-RU`, `ar-SA`, `hi-IN`.
8. `max_per_query` deve estar em `[1,25]` e `max_candidates` em `[1,20]`; ausência,
   zero e número fracionário são inválidos. Nenhum default implícito atravessa a
   fronteira MCP.
9. Campos desconhecidos, tipos incompatíveis e mais de um valor JSON no texto de
   `arguments` são inválidos.

Uma falha em qualquer item retorna JSON-RPC `-32602` e nenhum provedor, expansor ou
coletor é chamado. Cancelamento do contexto impede o início de novas chamadas.

## 5. Projeção canônica de sucesso

Uma chamada bem-sucedida a `bateia_recommend` usa a resposta MCP:

```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"schema_version\":\"1\",\"status\":\"COMPLETE\",\"summary\":{\"expansion_count\":0,\"candidate_count\":0,\"assessed_count\":0},\"decision\":{},\"alternatives\":{},\"failures\":[],\"warnings\":[]}"
    }
  ]
}
```

O valor de `text` é um único objeto JSON UTF-8 compacto, sem prefixo Markdown,
code fence, comentário ou newline final. Seus membros de primeiro nível aparecem
exatamente na ordem mostrada e têm estes significados:

- `schema_version`: sempre `"1"`;
- `status`: somente `COMPLETE`, `PARTIAL` ou `INCONCLUSIVE`;
- `summary.expansion_count`: `len(IdeaRunResult.Expansion.Expansions)`;
- `summary.candidate_count`: `len(IdeaRunResult.Discovery.Candidates)` após
  deduplicação/canonicalização do pipeline;
- `summary.assessed_count`: `len(IdeaRunResult.Assessment.Assessments)`;
- `decision`: a projeção pública completa do `ReuseDecisionReport` produzido por
  `DecideReuse`;
- `alternatives`: a projeção pública completa do `SoftwareAlternativeReport`
  produzido por `RankSoftwareAlternatives`, incluindo `top_choice`, elegíveis,
  excluídos, fontes e `discovery_locales` quando existirem;
- `failures`: união deduplicada das falhas de descoberta, coleta e assessment que
  pertencem à execução;
- `warnings`: avisos públicos deduplicados, incluindo licença não reconhecida e
  lacunas regionais controladas.

`failures` é sempre array, inclusive vazio, e cada item usa somente `scope`,
`stage` e `message`. `warnings`, `alternatives.alternatives` e
`alternatives.excluded` também são arrays, nunca `null`.

### 5.1 Status

A precedência é:

1. `INCONCLUSIVE` quando `decision.decision == "INCONCLUSIVE"`;
2. `PARTIAL` quando a decisão não é inconclusiva, mas houve falha de provedor,
   descoberta, coleta ou assessment, budget esgotado, cancelamento parcial ou
   outra degradação registrada;
3. `COMPLETE` somente quando a execução terminou sem essas degradações e a decisão
   não é inconclusiva.

Aviso isolado de licença desconhecida não força `PARTIAL`, não altera a classe de
fit e não exclui candidato. Licença comprovadamente incompatível continua sendo
bloqueio. Identidade mutável, commit divergente, coleta incompleta ou budget
esgotado não podem produzir `COMPLETE`. `top_choice` só pode apontar para uma
alternativa `RECOMMENDED` ou `PARTIAL` com identidade imutável e consistente e
sem opacidade ou bloqueio referente a ela; caso contrário deve estar ausente.

### 5.2 Determinismo

As ordenações de domínio seladas por `DecideReuse` e
`RankSoftwareAlternatives` são preservadas. Nos conjuntos sem ordem semântica:

- falhas ordenam por `stage`, `scope`, `message` e são deduplicadas pela tupla;
- warnings normalizam whitespace, ordenam lexicalmente e são deduplicados;
- locales usam a ordem canônica da seção 4;
- fontes de evidência são deduplicadas e ordenadas pela representação JSON de
  seus campos públicos depois da sanitização;
- coleções de strings que representam conjuntos são normalizadas, deduplicadas e
  ordenadas; rankings e sequências com precedência de domínio não são
  reordenados genericamente.

Objetos semanticamente equivalentes e fixtures equivalentes, ainda que suas
coleções de conjunto cheguem reordenadas, geram bytes idênticos em `text`.
Tempos de execução, relógio, mapas não ordenados, caminhos e identificadores de
processo não participam da projeção.

## 6. Sanitização e conteúdo não confiável

A projeção ocorre em L1 e é consumida igualmente por CLI e MCP. Antes da
serialização pública:

- `EvidenceSource.snippet` é removido de `decision` e `alternatives`; permanecem
  referência `repo@commit`, tipo, caminho, símbolo, linha, regra, hash e força;
- nenhum `EvidenceBundle`, conteúdo de README, licença, código, teste, manifesto
  ou arquivo coletado é incluído;
- URLs de candidato, quando preservadas pelo relatório de alternativas, não podem
  conter userinfo, query, fragmento, credencial ou esquema diferente de `https`;
  um valor inseguro vira string vazia, nunca é ecoado parcialmente;
- `message`, `warnings`, razões e caveats provenientes de fronteiras externas são
  mensagens públicas controladas ou sanitizadas; `error.Error()` de adaptador não
  é serializado diretamente;
- nomes de variáveis de ambiente, tokens, valores de `Authorization`, cookies,
  headers e segredos-canário nunca aparecem em sucesso nem em erro;
- texto de repositório que se pareça com instrução continua sendo dado. Ele não
  altera schema, estratégia, locales, política, limites, ferramenta escolhida ou
  autoridade de evidência.

Sanitização não inventa evidência nem transforma falha em ausência. Se não for
possível publicar uma mensagem com segurança, usa-se uma mensagem genérica
estável que preserve apenas estágio e classe da falha.

## 7. Erros MCP

- JSON-RPC malformado: `-32700`, mensagem pública `Parse error`;
- envelope `tools/call`, schema, validação ou campo desconhecido inválido:
  `-32602`, mensagem iniciada por `Invalid params` e sem eco de conteúdo sensível;
- ferramenta desconhecida: comportamento `-32601` preservado do baseline;
- recommender não injetado ou falha total sem resultado parcial seguro: `-32603`,
  mensagem pública genérica e sanitizada.

Falha parcial representável não vira JSON-RPC error: retorna
`RecommendationResult` com `PARTIAL` ou `INCONCLUSIVE`, falhas preservadas e
mensagens seguras. Erro de provedor nunca é convertido em “nenhum candidato”.

## 8. Igualdade entre CLI e MCP

Para uma mesma `RecommendationRequest` `live`, as mesmas dependências falsas e o
mesmo contexto não cancelado, CLI e MCP devem invocar a mesma instância lógica de
`RecommendationService` e obter o mesmo `RecommendationResult` canônico. A
igualdade exigida é entre:

1. o objeto retornado pelo serviço à CLI antes de qualquer escrita de dossiê; e
2. o objeto decodificado de `content[0].text` no MCP.

Depois da mesma sanitização/canonicalização, os dois objetos e seus bytes JSON
são iguais. A CLI pode continuar gravando artefatos adicionais e Markdown; esses
efeitos de casca não pertencem à comparação. O MCP não chama
`RunIdeaRecommendCommand`, não interpreta flags, não lê seus arquivos e não cria
diretórios. A implementação de ranking existe somente em L1.

CLI `local`, `replay` e `resume` permanecem compatíveis, mas não ampliam o modo
MCP deste passo. A extração da porta compartilhada não autoriza mudança semântica
nas ferramentas MCP antigas.

## 9. Obrigações discriminatórias para os autores seguintes

Oráculos devem observar ao menos: três ferramentas e o schema exato; sucesso
offline decodificável; igualdade CLI/MCP; propagação e proveniência de locales;
licença desconhecida não bloqueante; validação antes das portas; retrocompatibilidade
das duas ferramentas antigas; invariância a reordenação; neutralidade de conteúdo
hostil; ausência de canários; e wiring L4 funcional.

O contrato só pode ser aceito se rejeitar todas as mutações listadas na spec
congelada, com mutation score `1,0`. Casos genuinamente opacos preservam
`INCONCLUSIVE`/`Unknown`; `Unknown` nunca conta como sucesso. Este documento não é
selo, recibo nem veredito de implementação.
