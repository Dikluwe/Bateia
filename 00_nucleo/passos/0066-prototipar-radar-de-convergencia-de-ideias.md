# Passo 0066 — Prototipar radar de convergência de ideias em modo sombra

## Intenção congelada

Materializar em L1 um protótipo puro e determinístico que recebe projeções já
estruturadas de uma ideia e observações conceituais já declaradas, identifica quando
perspectivas diferentes convergem sobre o mesmo par candidato/conceito e mede, fora do
código produtivo, se essa convergência acrescenta utilidade sobre a melhor projeção
isolada.

“Radar” e “triangulação” são metáforas. O protótipo não localiza objetos, não infere
semântica e não alega comprovar uma implementação. Uma convergência é somente uma
hipótese conceitual em modo sombra; não altera busca, ranking, score, evidência,
recomendação ou decisão de reúso.

## Regime e autoridades

Aplicar o protocolo completo da skill `tekt-materializacao-segregada` como ensaio sem
atestação física de isolamento. O filesystem é compartilhado.

- autor da intenção: usuário;
- autor do contrato: escreve apenas este passo, prompts e `benchmark/0066`, antes do
  candidato, sem escrever Go nem ler implementação futura;
- autor dos oráculos: recebe intenção, contrato, baseline e corpus congelados, escreve
  testes positivos, negativos e opacos sem ler o candidato;
- adversário: materializa M01–M14 sem corrigir a solução;
- implementador: recebe passo, contrato e prompts produtivos selados e escreve somente
  os dois novos arquivos Go;
- verificador: executa os gates sobre artefatos congelados sem editá-los.

Mudança neste passo, contrato, baseline, budget, corpus/ground truth ou oráculos depois
do selo invalida a cadeia a partir da primeira entrada alterada.

## API pública mínima em L1

Criar tipos em `01_core/domain/idea_radar.go` e uma função pura em
`01_core/usecase/idea_radar.go`:

```text
AnalyzeIdeaRadar(domain.IdeaRadarInput) domain.IdeaRadarReport
```

Tipos públicos, com JSON `snake_case` e sem `omitempty`:

```text
IdeaRadarInput
  SchemaVersion string
  Projections   []IdeaRadarProjection
  Observations  []IdeaRadarObservation

IdeaRadarProjection
  ID           string
  Kind         IdeaProjectionKind
  OriginStatus EvidenceOriginStatus
  OriginID     string
  Terms        []string
  Actions      []string
  Objects      []string
  Required     []string

IdeaRadarObservation
  ID           string
  ProjectionID string
  CandidateID  string
  ConceptID    string
  Stance       IdeaRadarStance

IdeaRadarAssessment
  TargetID                    string
  CandidateID                 string
  ConceptID                   string
  Outcome                     IdeaRadarOutcome
  SupportingObservationIDs    []string
  ConflictingObservationIDs   []string
  AmbiguousObservationIDs     []string
  UnknownOriginObservationIDs []string
  ProjectionKinds             []IdeaProjectionKind
  DistinctOrigins             []string
  MissingProjectionKinds      int
  MissingOrigins              int
  CoalitionDraft              IdeaRadarCoalitionDraft
  Reason                      string

IdeaRadarCoalitionDraft
  ClaimID          string
  ObservationIDs  []string
  RequiredRoles   []EvidenceCoalitionRole
  ReadyFor0065    bool
  Reason          string

IdeaRadarMetrics
  Projections                 int
  Observations                int
  CanonicalObservations       int
  DuplicateObservations       int
  Targets                     int
  Convergent                  int
  Conflicts                   int
  Ambiguous                   int
  InsufficientPerspectives    int
  UnknownOrigins              int
  InsufficientOrigins         int
  Operations                  int
  LimitExceeded               bool
  VerdictChanges              int

IdeaRadarReport
  SchemaVersion string
  Status        string
  Reason        string
  Assessments   []IdeaRadarAssessment
  Metrics       IdeaRadarMetrics
```

Constantes serializadas exatas:

- schemas `idea-radar-input-v1` e `idea-radar-report-v1`;
- `IdeaProjectionKind`: `PROBLEM`, `BEHAVIOR`, `ARCHITECTURE`, `CONSTRAINT`,
  `OPERATION`;
- stance: `SUPPORT`, `CONFLICT`, `AMBIGUOUS`;
- outcomes: `CONVERGENT`, `CONFLICT`, `AMBIGUOUS`,
  `INSUFFICIENT_PERSPECTIVES`, `UNKNOWN_ORIGIN`, `INSUFFICIENT_ORIGINS`;
- relatório: `COMPLETE` ou `UNKNOWN`;
- `ReadyFor0065=false`, sempre; `VerdictChanges=0`, sempre.

Aliases Go idiomáticos são permitidos, mas campos, valores serializados e assinatura são
contrato. Todas as slices do relatório são arrays não nulos.

## Entrada estruturada e identidade

As cinco projeções representam perguntas diferentes: problema, comportamento,
arquitetura, restrição e operação. `Terms`, `Actions`, `Objects` e `Required` preservam
a estrutura recebida para auditoria, mas o analisador não extrai conceitos desses textos.
Observações são declarações explícitas ligando uma projeção a um `CandidateID` e a um
`ConceptID` canônico.

IDs canônicos de projeção, observação, candidato e conceito são normalizados somente por
remoção de whitespace ASCII nas pontas e conversão ASCII `A-Z` para `a-z`. Depois disso
devem ser não vazios, ter no máximo 128 bytes e conter apenas `a-z`, `0-9`, `.`, `_`,
`:`, `/` ou `-`. Não há Unicode folding, stemming, tokenização, sinônimos, fuzzy match ou
equivalência por palavras. `OriginID` é opaco, ASCII não vazio, sem whitespace nas pontas,
com até 128 bytes, e não é convertido para minúsculas.

Cada projeção tem ID canônico único e kind válido. Na entrada, `OriginStatus` aceita
somente `RECORDED` com `OriginID` válido ou `UNKNOWN_ORIGIN` com `OriginID` vazio;
`SAME_ORIGIN` e `DISTINCT_ORIGIN` são estados derivados e inválidos na entrada. Duas
projeções com kinds diferentes podem ter o mesmo `OriginID`; isso prova diversidade de
perspectiva, mas não diversidade de origem.

Cada observação referencia exatamente uma projeção. IDs de observação são únicos depois
da canonicalização. Reusar o mesmo ID para conteúdo diferente torna o relatório
`UNKNOWN/ambiguous_input`. Observações semanticamente idênticas pela tupla canônica
`(projection_id,candidate_id,concept_id,stance)` são deduplicadas, conservando o menor ID;
IDs, termos, hashes, caminhos e nomes diferentes não criam origem independente.

`TargetID` e `ClaimID` são materializados exatamente como
`candidate_id + "#" + concept_id`; os caracteres permitidos tornam a separação não
ambígua porque `#` não é aceito nos IDs de entrada.

## Semântica de convergência

Agrupar observações pelo alvo canônico `(CandidateID, ConceptID)`. Somente `SUPPORT`
contribui para os limiares. Triangulação v1 exige, de forma fixa:

- pelo menos três `ProjectionKind` distintos entre suportes; e
- pelo menos três `OriginID` distintos entre suportes de origem `RECORDED`.

Diversidade de kind e de origem são contadas separadamente. Três termos, arquivos,
hashes, queries, observações ou IDs nunca substituem três kinds e três origens explícitas.
Não enumerar trincas: basta acumular os conjuntos distintos por alvo.

Precedência observável por alvo válido:

1. qualquer stance `CONFLICT` produz `CONFLICT` e preserva todos os sinais;
2. sem conflito, qualquer stance `AMBIGUOUS` produz `AMBIGUOUS`;
3. menos de três kinds entre todos os suportes produz `INSUFFICIENT_PERSPECTIVES`;
4. três kinds existem, mas suportes `RECORDED` cobrem menos de três kinds produz
   `UNKNOWN_ORIGIN`;
5. ao menos três kinds têm origem registrada, mas há menos de três `OriginID` distintos,
   produz `INSUFFICIENT_ORIGINS`;
6. os dois limiares demonstrados produzem `CONVERGENT`.

Suportes de origem desconhecida permanecem em `UnknownOriginObservationIDs`. Eles não
invalidam uma convergência já demonstrada por outros suportes registrados. Conflito não
desaparece por maioria. `MissingProjectionKinds=max(0,3-kinds_distintos_de_suporte)` e
`MissingOrigins=max(0,3-origens_distintas_registradas)`.

Razões exatas de assessment: `idea_convergent`, `conflicting_projection`,
`ambiguous_projection`, `projection_diversity_below_minimum`, `origin_unrecorded` e
`distinct_origins_below_minimum`. Relatório `COMPLETE` usa `analyzed`; relatório
`UNKNOWN` usa `invalid_schema`, `no_projections`, `no_observations`, `ambiguous_input`,
`invalid_input` ou `budget_exceeded`. Entrada inválida ou associação ambígua não produz
sucesso parcial.

Ordenar projeções, observações canônicas, assessments, IDs, kinds e origens pelos valores
canônicos. A função não muta o input. Permutações equivalentes e `nil` versus vazio em
slices textuais produzem JSON byte a byte idêntico após `encoding/json.Marshal`.

## Ponte conservadora para o passo 0065

O radar não fabrica `EvidenceSource`, papéis de evidência ou independência causal.
Portanto, ele não chama `AnalyzeEvidenceCoalitions` nesta versão. Cada assessment inclui
um `IdeaRadarCoalitionDraft` explícito:

- `ClaimID` é o TargetID;
- `ObservationIDs` contém os suportes canônicos;
- `RequiredRoles` é sempre `[IMPLEMENTATION, INDEPENDENT_TEST]`, em ordem canônica;
- `ReadyFor0065` é sempre `false`;
- `Reason` é sempre `requires_evidence_binding`.

Um passo futuro poderá resolver cada observação para fonte física, polaridade, papel e
proveniência e então criar `EvidenceCoalitionInput`. A convergência conceitual nunca é
convertida diretamente em testemunha 0065.

## Limites e custo

Limites inclusivos da versão 1:

- 64 projeções;
- 1024 observações;
- 512 alvos canônicos;
- 128 bytes por ID e por OriginID;
- 64 itens e 4096 bytes somados em cada dimensão textual de uma projeção;
- 4096 operações lógicas.

`Operations` conta uma unidade por projeção validada, uma por observação normalizada e
uma por associação de observação canônica ao alvo. Canonicalização escalar, ordenação e
emissão não adicionam unidades. O algoritmo deve usar índices e uma passagem limitada;
power set, combinação de perspectivas e comparação quadrática entre observações são
proibidos. Qualquer limite excedido produz relatório `UNKNOWN/budget_exceeded`,
`LimitExceeded=true`, assessments vazios e nenhuma utilidade presumida.

## Medição externa de utilidade

O código produtivo não recebe ground truth, não calcula precisão/recall e não escolhe o
baseline. O corpus congelado em `benchmark/0066/baseline.json`, consumido somente por
oráculos/harness de teste, contém famílias, alvos esperados, alvos proibidos e predições
de cada projeção isolada. A melhor projeção isolada é escolhida globalmente por maior F1,
desempate por precisão, recall e nome canônico do kind.

Predição radar é somente assessment `CONVERGENT`. Por correspondência exata de TargetID,
calcular externamente micro e por caso: precisão, recall, F1, falsa convergência
`FP/(TP+FP)` e contaminação `preditos_proibidos/preditos`; denominador zero torna a
métrica `Unknown`, nunca sucesso. Ganho incremental é `radar - melhor_isolada` para
precisão, recall e F1.

Veredito de utilidade:

- `OBSERVED` somente se precisão micro do radar não regredir, recall ou F1 micro melhorar
  em pelo menos `0.10`, e houver melhoria de F1 por caso em pelo menos dois casos de pelo
  menos duas famílias diferentes;
- `INCONCLUSIVE` se há denominador indefinido, ground truth/resultado ausente, menos de
  dois casos ou duas famílias elegíveis, schema divergente ou budget excedido;
- `NOT_OBSERVED` em qualquer outro resultado válido.

Falsa convergência e contaminação são sempre reportadas, mas não substituem o gate de
precisão. O corpus e thresholds não podem ser ajustados após observar o candidato.
`OBSERVED` atesta utilidade apenas nesse corpus, não utilidade geral.

## Invariantes e não objetivos

- nenhuma função existente chama o radar;
- nenhum arquivo ou schema existente é editado;
- `VerdictChanges=0`; nenhuma saída possui autoridade decisória;
- L1 permanece sem I/O, rede, filesystem, relógio, aleatoriedade, LLM, embeddings ou
  dependência externa;
- não derivar projeções de `IdeaSpec`, executar busca, inferir conceitos, resolver
  sinônimos, provar causalidade, alterar ranking ou promover `PROVEN`;
- não alegar que repetição lexical é triangulação;
- não usar o benchmark como tabela interna de produção.

## Corpus, ataques e gates

O corpus congelado tem seis casos e três famílias. Deve cobrir: convergência verdadeira,
ganho que a melhor projeção isolada perde, mesma origem em três kinds, três origens em
apenas dois kinds, origem desconhecida, conflito, ambiguidade, duplicata sem ganho,
canonicalização ASCII, IDs inválidos, vazio, schema desconhecido, limites e pelo menos 20
permutações.

O gate adversarial deve eliminar:

- M01: dois kinds bastam;
- M02: três kinds da mesma origem convergem;
- M03: três OriginIDs no mesmo kind convergem;
- M04: origem desconhecida conta como distinta;
- M05: termos/hash/caminho/ID diferentes provam independência;
- M06: conflito é ignorado por maioria;
- M07: ambiguidade é promovida;
- M08: duplicata multiplica kinds/origens;
- M09: fuzzy match ou lowercasing Unicode funde conceitos;
- M10: ordem ou nil/vazio altera JSON;
- M11: input vazio/inválido/schema desconhecido vira `COMPLETE`;
- M12: budget excedido produz parcial ou convergência;
- M13: draft 0065 fica ready ou ganha testemunha fabricada;
- M14: radar altera veredito/ranking ou `VerdictChanges` deixa zero.

Selo exige `mutation_score=1.0`, corpus contratual 100%, 20/20 permutações idênticas,
contagem de operações exata e zero mudanças de veredito. Depois: testes focais, race,
`go test ./...`, `go vet ./01_core/...`, build do pacote `./04_wiring`,
`crystalline-lint .` e `git diff --check`.

Budget de calibração: duas revisões focais por causa e no máximo três execuções completas
antes do selo. Duas revisões consecutivas sem ganho discriminatório reabrem o desenho.

## Escopo futuro de arquivos

Implementador: somente `01_core/domain/idea_radar.go` e
`01_core/usecase/idea_radar.go`, depois dos prompts produtivos selados. Autor de testes:
somente os dois `_test.go` espelhados, preferencialmente em packages externos. Nenhum
arquivo Go existente, inclusive IdeaSpec e coalizões 0065, pode ser editado.
