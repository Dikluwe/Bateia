# Contrato do oráculo independente — Passo 0027

**Autoridade:** `/root/oracle_0027`  
**Regime:** protocolo Tekt completo  
**Estado:** congelado antes do candidato; sem veredito  
**Atestação máxima possível:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

Este artefato especifica somente a superfície mínima que permite compilar as
suítes independentes. Ele não autoriza o implementador a ler o split `TEST`, o
registry esperado, as fixtures protegidas ou os testes-oráculo.

## API pública mínima congelada

Além dos nomes, valores e estruturas literais da revisão ativa de
`00_nucleo/benchmark/0027/contract.md`, a materialização expõe:

- constantes Go `RunKindRecommendation`, `RunKindMarketDiscovery`,
  `MarketDiscoveryComplete`, `MarketDiscoveryPartial`,
  `MarketDiscoveryInconclusive`, `RunStageMarketPlan`,
  `RunStageMarketDiscover`, `RunStageMarketResolve`, `RunStageMarketPublish`,
  `ArtifactConnectorCapabilities`, `ArtifactMarketPlan`,
  `ArtifactMarketSnapshot`, `ArtifactMarketCandidates`, `ArtifactMarketGaps`,
  `AttemptPlanned`, `AttemptSucceeded`, `AttemptEmptyObserved`,
  `AttemptUnknownProvider`, `AttemptMarketNotFiltered`, `AttemptCancelled` e
  `AttemptIntegrityError`, com os valores literais do contrato;
- `domain.MarketDiscoveryStartRequest`, contendo `Idea`, `QueryProposals`,
  `Markets`, `ConnectorIDs` e `Budget`;
- `domain.MarketDiscoveryRequest`, a requisição de uma tentativa, contendo
  `QueryID`, `QueryLocale`, `Query`, `Market` e `Limit`;
- `domain.MarketDiscoveryBudget`, `MarketDiscoveryAttempt`,
  `ConnectorDiscovery`, `IdentityWitness` e `MarketDiscoveryStartReceipt`;
- `usecase.NewDiscoveryConnectorRegistry`, `Capabilities` e `Resolve`;
- `usecase.ValidateMarketDiscoveryStartRequest`;
- `usecase.CanonicalMarketDiscoveryStartRequestHash`;
- `usecase.PlanMarketDiscovery` e `AllocateMarketDiscoveryBudget`;
- `usecase.ResolveProductIdentity` e `ConsolidateProductCandidates`;
- `usecase.ProjectMarketDiscoveryArtifacts`;
- `usecase.NewRecommendationRunServiceWithMarketDiscovery`, que estende o
  serviço durável único e conserva `Start`, `StartMarketDiscovery`, `GetRun`,
  `GetArtifact`, `Cancel`, `Execute` e `Recover` sobre as mesmas portas;
- `mcp.MarketDiscoveryStarter` e
  `(*MCPServer).WithMarketDiscovery(connectors, starter)`;
- `itunes.NewMarketDiscoveryConnector(*http.Client)`;
- `marketreplay.NewConnector([]byte)`;
- `main.NewMarketDiscoveryMCPServer`, extensão aditiva da composição durável.

Nenhum nome privado, layout físico de store, algoritmo auxiliar ou formato de
arquivo interno integra o contrato.

As estruturas auxiliares têm campos exportados e tags snake_case equivalentes:

```go
type MarketDiscoveryBudget struct {
    MaxQueries int; MaxResultsPerQuery int; MaxTotalCandidates int
}
type MarketDiscoveryStartRequest struct {
    Idea IdeaSpec; QueryProposals []MarketQueryProposal; Markets []string
    ConnectorIDs []string; Budget MarketDiscoveryBudget
}
type MarketDiscoveryRequest struct {
    QueryID string; QueryLocale string; Query string; Market string; Limit int
}
type MarketDiscoveryAttempt struct {
    AttemptID string; QueryID string; QueryLocale string; Market string
    ConnectorID string; ConnectorVersion string; MarketFilterApplied bool
    QueryHash string; EffectiveQueryHash string; RequestHash string
    PlannedCost int; Status MarketAttemptStatus; Reason string
}
type ConnectorDiscovery struct {
    Status MarketAttemptStatus; Candidates []ProductCandidate
    ResponseHash string; ErrorClass string
}
type MarketDiscoveryStartReceipt struct {
    SchemaVersion string; RunID RunID; RunKind RunKind
    Lifecycle RecommendationLifecycle; Version uint64
}
```

`RunTransition` e `RecommendationRun/View` recebem `RunKind` e
`MarketDiscoveryStatus`; as interfaces duráveis 0026 são estendidas de forma
aditiva, nunca duplicadas.

## Controles positivos

1. registry canônico lista três conectores sem rede;
2. a matriz primária e os quatro controles cruzados preservam locale e mercado;
3. orçamento round-robin independe da ordem de conectores/propostas;
4. alias forte bilateral permite `SAME`; homônimo e URL isolada não permitem;
5. Apple aceita somente a requisição fixa allowlisted, projeta `trackId` e deixa
   listing locale e publisher country vazios;
6. replay repetido produz bytes idênticos sem transporte;
7. as tools 1–5 e 7 permanecem deep-equal; a tool 6 permanece deep-equal após
   remover os cinco sufixos aditivos contratados de `artifact.enum`; tools 8–9
   são exatas;
8. recommendation e market discovery compartilham store, queue, gerador de ID,
   lifecycle, cancelamento, recovery e idempotência.

## Casos negativos que as suítes discriminam

- registry inválido, duplicado ou incompatível antes de qualquer efeito;
- mercado derivado do locale ou dimensão copiada para produto;
- budget consumido integralmente pelo primeiro conector;
- merge por nome, publisher, URL isolada, query, mercado ou popularidade;
- promoção de repo sem commit verificado ou metadado de catálogo para `PROVEN`;
- host/path/parâmetro/redirect controlado, query não escapada, SSRF e callback;
- `429`, timeout, non-2xx, oversized e JSON inválido classificados como vazio;
- item Apple sem `trackId`, URL divergente, canário, HTML/review/body/erro externo;
- replay com versão/hash divergente ou fallback de rede;
- cancelamento que inicia chamada posterior;
- segunda fila/store/idempotência ou regressão das sete tools certificadas;
- campos desconhecidos, JSON múltiplo ou valores fora do schema chegando à porta;
- artefato com query, source text, rationale, corpo, review, PII ou configuração.

## Opacos preservados

Sem testemunha forte entre fontes, a resolução permanece `POSSIBLE` ou
`UNKNOWN`. Falha do provedor permanece `UNKNOWN_PROVIDER`; budget não executado
permanece `BUDGET_EXHAUSTED`; integridade divergente permanece
`INTEGRITY_ERROR`. Nenhum desses casos conta como sucesso, ausência ou mutação
rejeitada.

## Métricas e gate de referência

As frações são inteiras e sem arredondamento. Exposição abaixo de 80% torna
precisão, recall e ganho `INCONCLUSIVE_EXPOSURE`. Caso contrário: precisão top 10
por família >= 70%, recall em `VALIDATION+TEST` >= 70%, ganho Apple >= 20 pontos,
ao menos 4/6 ideias com candidato incremental, merge 100%, inferência dimensional
0%, produto sem repo 100% e replay byte-idêntico 100%.

O gate de referência compara saídas candidatas com um modelo independente de
ordenação, produto cartesiano, round-robin, identidade e cálculo de métricas. O
modelo não recebe ID de mutação nem seleciona defeitos por metadata.

## Segregação

- leitura: spec/contrato/source selection/manifest, ADR-0002, certificados e
  receipts 0025/0026 e superfícies públicas basais certificadas;
- escrita: este contrato, corpus/registry/fixtures, três suítes Go e três prompts;
- proibido: produção, ataques, fixes, receipts, selos, certificado, laudo ou
  veredito 0027;
- o coordenador/implementador não deve ler `oracle/protected/**` nem as suítes
  `*_step0027_oracle_test.go` antes de congelado o candidato.
