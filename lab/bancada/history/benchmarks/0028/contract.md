# Contrato observável — Passo 0028

**Revisão:** R3 (substitui R2 antes de oracle e selo)  
**Estado:** candidato a selo, anterior à implementação  
**Regime:** protocolo Tekt completo; alvo de atestação
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Papel:** autoridade segregada de CONTRATO  
**Implementação candidata lida:** não  
**Veredito positivo reservado ao verificador:**
`MECHANISM_IMPROVEMENT_SYNTHESIS_READY`

`DEVE`, `NÃO DEVE`, `EXATAMENTE` e `SOMENTE` são normativos. Este contrato
cobre apenas o fragmento observável e os hashes registrados. Ele não emite o
veredito reservado, não autoriza decisão de adoção e não prova equivalência
funcional geral.

## 1. Raiz causal, capacidades e invalidação

| Entrada congelada | SHA-256 |
|---|---|
| `00_nucleo/passos/0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md` | `9105b27e0e8e9165b0b133a647669e858b977ea2ab8c7d951ed0f217a26ac1a5` |
| `00_nucleo/benchmark/0028/manifest.json` | `0b8fbcebf9f3af491ae1851f59715c248fa63181617cc336dc39b6969fb555ad` |
| `00_nucleo/benchmark/0028/mechanism-rules.json` | `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5` |
| `00_nucleo/ADR/ADR-0002-estrutura-de-pastas-e-camadas-tekt.md` | `b457fadd55d6328c2afce969b92a527b461edab9094b1bef66fcfebe71b79a68` |
| certificado 0025 | `bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d` |
| certificado 0026 | `a04d6b07a939c1ce5fb7fb184e5cb56cc3cc34ed0b2b54c8087415a18a4a99ed` |
| certificado 0027 | `acf05df728b05ccdc5bc15f78a18e49f19e743b2016ab1141e553ca2417c7e43` |

O autor deste contrato recebeu spec, manifesto, ADR, baseline público e
certificados predecessores. Pôde escrever somente este contrato, o catálogo de
regras e prompts L0 proprietários. Não leu nem pode escrever implementação
candidata, corpus/suites protegidos, ataques, receipts, seals, certificado,
laudo ou veredito do Passo 0028. Trabalho em workspace compartilhado impede
alegação de isolamento forte mesmo que todos os hashes e gates passem.

Qualquer alteração na spec, manifesto, catálogo, contrato, prompts, taxonomia,
política de licença, baseline convencional ou corpus invalida a cadeia desde a
primeira entrada afetada. Hash identifica entrada; não substitui registro das
capacidades nem prova isolamento.

## 2. Escopo e separações epistemológicas

O Passo 0028 adiciona:

1. extração pura, determinística e versionada de `mechanism_catalog` e
   `mechanism_gaps` durante novos runs `RECOMMENDATION`, usando somente o
   `EvidenceBundle` congelado desse run;
2. run durável `IMPROVEMENT_SYNTHESIS`, no runtime único já certificado;
3. tool MCP 10 `bateia_start_improvement_synthesis`;
4. cinco artefatos seguros de síntese e dois artefatos seguros de mecanismo.

Não adiciona modelo de IA embutido, embedding online, crawler, clone, download,
execução de candidato, segundo store, segunda fila, novo cursor secret, escrita no
alvo, implementação automática nem decisão de adoção.

As categorias são disjuntas:

| Categoria | Autoridade factual |
|---|---|
| `Observation` | fato limitado a uma superfície e hash congelados |
| `ArchitecturalMechanism` | estratégia implementada, qualificada pelas testemunhas da mesma source imutável |
| `ImprovementHypothesis` | proposição falsificável, ainda não é adoção mesmo em `SUPPORTED` |
| `AdoptionDecision` | decisão posterior, fora do schema e deste passo |

Produto não repo-backed permanece observação. Proposta externa permanece
proposta. Nenhuma conversão entre essas categorias ocorre por nome, descrição,
popularidade, novidade ou autoridade do emissor.

## 3. Refinamento obrigatório dos Passos 0025–0027

`tools/list` retorna EXATAMENTE dez tools, nesta ordem:

1. `bateia_search`
2. `bateia_benchmark`
3. `bateia_recommend`
4. `bateia_start_recommendation`
5. `bateia_get_run`
6. `bateia_get_artifact`
7. `bateia_cancel_run`
8. `bateia_list_connectors`
9. `bateia_start_market_discovery`
10. `bateia_start_improvement_synthesis`

Nome, descrição, schema e comportamento das posições 1–9 são deep-equal aos
certificados anteriores, salvo estas extensões aditivas e somente estas:

- `bateia_get_run` preserva seu schema de entrada e acrescenta a variante
  `IMPROVEMENT_SYNTHESIS` na projeção persistida de `run_kind`, mais o campo
  sempre presente `improvement_status` descrito na seção 9;
- `bateia_get_artifact` preserva shape, cursor e limites e amplia somente sua
  enum na ordem exata da seção 10;
- novos runs `RECOMMENDATION` publicados sob 0028 exigem também os dois
  artefatos de mecanismo; runs históricos permanecem byte-estáveis e não são
  reinterpretados;
- o dispatcher único aceita `IMPROVEMENT_SYNTHESIS` sem alterar store, fila,
  idempotência, CAS, quota, cursor, recovery ou cancelamento.

São obrigatórios três receipts de refinamento separados:

- `MCP_RECOMMENDATION_GATEWAY_READY`: tools 1–3 e recomendação síncrona
  preservadas; mecanismo não altera ranking, decisão, `AdoptionRate` ou resposta
  síncrona certificada;
- `MCP_DURABLE_RUNS_READY`: tools 4–7, lifecycle, CAS, idempotência, quota,
  cancelamento, cursor HMAC, publicação atômica e recovery preservados;
- `MULTISOURCE_MARKET_DISCOVERY_READY`: tools 8–9, conectores, quatro dimensões,
  identidade, budget, replay e artefatos de mercado preservados.

Teste novo passando não substitui receipts de refinamento. A tool 10 não pode ser
inserida antes das nove existentes.

## 4. API pública L1 de mecanismo

L1 define nomes, strings e tags JSON exatos:

```go
type MechanismEvidenceStatus string
const (
    MechanismProven       MechanismEvidenceStatus = "PROVEN"
    MechanismPartial      MechanismEvidenceStatus = "PARTIAL"
    MechanismInconclusive MechanismEvidenceStatus = "INCONCLUSIVE"
)

type MechanismNovelty string
const (
    NoveltyConventional MechanismNovelty = "CONVENTIONAL"
    NoveltyDivergent    MechanismNovelty = "DIVERGENT"
    NoveltyUnknown      MechanismNovelty = "UNKNOWN"
)

type MechanismWitnessKind string
const (
    WitnessImplementation MechanismWitnessKind = "IMPLEMENTATION"
    WitnessTest           MechanismWitnessKind = "TEST"
    WitnessExample        MechanismWitnessKind = "EXAMPLE"
    WitnessManifest       MechanismWitnessKind = "MANIFEST"
    WitnessSchema         MechanismWitnessKind = "SCHEMA"
    WitnessConfiguration  MechanismWitnessKind = "CONFIGURATION"
    WitnessDocumentation  MechanismWitnessKind = "DOCUMENTATION"
    WitnessCounter        MechanismWitnessKind = "COUNTEREVIDENCE"
)

type ImmutableCandidateRef struct {
    Repository string `json:"repository"`
    Commit     string `json:"commit"`
}

type MechanismWitness struct {
    Repository  string               `json:"repository"`
    Commit      string               `json:"commit"`
    Path        string               `json:"path"`
    Symbol      string               `json:"symbol"`
    Line        int                  `json:"line"`
    SourceKind  EvidenceKind         `json:"source_kind"`
    Kind        MechanismWitnessKind `json:"kind"`
    AssessmentStatus EvidenceStatus  `json:"assessment_status"`
    RuleID      string               `json:"rule_id"`
    ContentHash string               `json:"content_hash"`
}

type MechanismComponentRef struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Kind string `json:"kind"`
}

type ArchitecturalMechanism struct {
    ID                string                  `json:"id"`
    MechanismClass    string                  `json:"mechanism_class"`
    Name              string                  `json:"name"`
    Purpose           string                  `json:"purpose"`
    Trigger           string                  `json:"trigger"`
    Strategy          string                  `json:"strategy"`
    Components        []MechanismComponentRef `json:"components"`
    DataFlow          []string                `json:"data_flow"`
    StateModel        []string                `json:"state_model"`
    FailureModel      []string                `json:"failure_model"`
    Boundaries        []string                `json:"boundaries"`
    Preconditions     []string                `json:"preconditions"`
    Tradeoffs         []string                `json:"tradeoffs"`
    Source            ImmutableCandidateRef   `json:"source"`
    Evidence          []MechanismWitness      `json:"evidence"`
    CounterEvidence   []MechanismWitness      `json:"counter_evidence"`
    Status            MechanismEvidenceStatus `json:"status"`
    Novelty           MechanismNovelty        `json:"novelty"`
    ExtractionRuleIDs []string                `json:"extraction_rule_ids"`
}

type ConventionalBaseline struct {
    BaselineID       string   `json:"baseline_id"`
    IdeaHash         string   `json:"idea_hash"`
    MechanismClasses []string `json:"mechanism_classes"`
}
```

`Repository` é `owner/repo`, sem URL, e `Commit` é exatamente 40 ou 64 hex
minúsculos. `ContentHash` é `sha256:` + 64 hex minúsculos. `Path` é relativo,
usa `/`, não começa por `/`, não contém segmento vazio, `.` ou `..`; não é lido
durante extração. `Symbol` pode ser vazio; quando uma alegação depende de símbolo,
ele é obrigatório. `Line >= 0`. Todo texto público é UTF-8 NFC, sem C0/C1/DEL,
com whitespace Unicode colapsado e limites da seção 12.

Slices públicas são `[]`, nunca `null`. A projeção nunca inclui snippet, bytes de
arquivo, README, licença integral, teste integral, body, checkout path ou URL.

## 5. Catálogo fechado e algoritmo de extração

O único catálogo desta revisão é `tekt.mechanism-rules.v2`, versão `2`, hash da
seção 1. Ele contém EXATAMENTE quatro classes:

1. `APPEND_ONLY_LOG_RECOVERY`;
2. `BOUNDED_QUEUE_BACKPRESSURE`;
3. `CONTENT_ADDRESSED_STORAGE`;
4. `PORT_ADAPTER_ISOLATION`.

Não existe fallback lexical, fuzzy match, LLM, regra embutida extra, alias ou
"classe mais próxima". Uma classe/regra não catalogada produz gap
`RULE_NOT_IN_MECHANISM_CATALOG`, nunca mecanismo inventado. O arquivo JSON é
validado estritamente: campos desconhecidos, regra duplicada, versão divergente,
`rule_id` duplicado entre papéis incompatíveis ou hash divergente impedem
inicialização.

O extrator R3 recebe conjuntamente o `EvidenceBundle` congelado e os
`CandidateAssessment` já avaliados do MESMO novo run `RECOMMENDATION`, antes de
ambos saírem de escopo. `EvidenceBundle.Items` são `CollectedEvidence`; eles não
possuem RuleID, Symbol, Line, Strength ou papel de mecanismo. R3 não atribui esses
campos ao bundle.

Fontes estruturadas são SOMENTE `EvidenceSource` aninhadas em
`CandidateAssessment.ProblemAlignment` ou `CandidateAssessment.Capabilities` sob
seu `CapabilityEvidence.Status`. O extrator:

1. valida que `CandidateAssessment.Repository`, `EvidenceSource.Ref` e o
   `CollectedEvidence.Ref` são a mesma identidade imutável `owner/repo@commit`;
2. resolve cada source para EXATAMENTE um item do mesmo bundle pela igualdade de
   `(kind,ref,path,content_hash)`; zero matches gera
   `UNBOUND_EVALUATED_SOURCE`, mais de um gera `AMBIGUOUS_BUNDLE_BINDING`;
3. recomputa `sha256:` + SHA-256 de `CollectedEvidence.Content` apenas para
   conferir `CollectedEvidence.ContentHash` e `EvidenceSource.ContentHash`;
   divergência gera `EVIDENCE_CONTENT_HASH_MISMATCH` e a source não participa;
4. ignora integralmente `EvidenceSource.Snippet` e nunca procura em `Content`
   termos, patterns, symbol, line, regra, mecanismo, instrução ou classificação;
5. usa `EvidenceSource.RuleID` como a única identidade de regra e exige entrada
   literal correspondente em `accepted_evidence_rules` do catálogo;
6. obtém `role` e `witness_kind` SOMENTE dessa entrada catalogada; preserva
   `EvidenceSource.Kind`, Symbol, Line e o `CapabilityEvidence.Status` pai;
7. valida compatibilidade estrutural: `IMPLEMENTATION` aceita source kind
   `file|symbol`; `TEST`, `EXAMPLE`, `MANIFEST` e `DOCUMENTATION` exigem o kind
   homônimo; `SCHEMA|CONFIGURATION` aceitam `file|manifest`;
   `COUNTEREVIDENCE` preserva o source kind avaliado e nunca vira positivo;
8. agrupa por `(repository,commit,mechanism-rule.id)`; nunca cruza repositório ou
   commit; deduplica pela tupla
   `(repository,commit,path,symbol,line,source_kind,kind,assessment_status,rule_id,content_hash)`;
9. preenche somente `fillable_fields` com templates do catálogo e dados
   estruturais seguros; texto bruto não vira alegação;
10. calcula status, novidade e IDs; ordena mecanismos por `id` e gaps por
    `(code,source_repository,source_commit,rule_id,witness_hash)`.

Uma `EvidenceSource` pode conter campos adicionais na evolução já certificada;
eles não participam salvo se este contrato os nomeia. Em particular, Strength da
source, quando presente, não eleva o teto imposto pelo `EvidenceStatus` pai.

Status possui esta precedência exata:

- source catalogada com role `CONFLICT` e status pai `CONFLICT`, ou identidade
  inconsistente, produz `INCONCLUSIVE`; conflito com outro status produz gap de
  avaliação inconsistente;
- source positiva sob status pai `ABSENT|INCONCLUSIVE|CONFLICT` não conta e gera
  gap; source sob `PARTIAL` impõe teto `PARTIAL` ao mecanismo;
- `PROVEN` exige ao menos uma role `IMPLEMENTATION` e uma `INDEPENDENT`, ambas sob
  status pai `PROVEN`, dois `content_hash` distintos, mesma source imutável, todos
  os campos centrais `purpose,strategy,boundaries` preenchidos e nenhuma
  contraevidência central;
- `PARTIAL` exige ao menos uma source positiva ligada ao bundle, status pai
  `PROVEN|PARTIAL` e campos centrais preenchidos, mas não satisfaz toda a prova;
  documentação isolada fica no máximo aqui;
- qualquer outro caso não produz mecanismo público e produz gap Unknown.

`INDEPENDENT` pode ser teste, exemplo executável, manifesto, schema ou configuração
causalmente relacionada. Documentação nunca conta como implementação nem
independência para `PROVEN`. Duas ocorrências com mesmo content hash contam uma.
Nome de pacote/produto, README promocional, tag, estrela, download, ranking ou
`AdoptionRate` não são regras avaliadas e jamais são inferidos de `Content`.

Novidade é calculada depois do status, contra um `ConventionalBaseline` congelado:

- classe presente no baseline: `CONVENTIONAL`;
- classe ausente e mecanismo `PROVEN`: `DIVERGENT`;
- baseline ausente/inválido ou mecanismo não provado: `UNKNOWN`.

`DIVERGENT` não muda status, score, ordenação, licença, fit ou hipótese.

## 6. Canonicalização e identidades causais

Todas as funções abaixo são puras e usam SHA-256, hex minúsculo. Objetos são
structs de campos fixos na ordem mostrada; a codificação é JSON compacto UTF-8,
sem indentação, HTML escaping desabilitado e sem newline. Mapas são proibidos nos
preimages. Conjuntos são NFC-normalizados, deduplicados por bytes exatos e
ordenados por bytes UTF-8 antes do marshal.

`mechanism_id` é `mech_` + 64 hex sobre:

```json
{"schema_version":"1","source":{"repository":"owner/repo","commit":"<immutable>"},"rule_id":"MR-...-V1","mechanism_class":"CLASS","purpose":"normalized","strategy":"normalized","witness_hashes":["sha256:..."]}
```

`witness_hashes` contém todos e somente os content hashes de evidência e
contraevidência, deduplicados e ordenados. Horário, ordem de descoberta, checkout,
PID, locale incidental e posição em array não participam.

`hypothesis_id` é `hyp_` + 64 hex sobre:

```json
{"schema_version":"1","title":"normalized","target_property":"normalized","target_gap_refs":[],"expected_effect":"normalized","mechanism_refs":[],"product_observation_refs":[],"derivation":"RULE_DERIVED|EXTERNAL_AGENT_PROPOSAL","source_hashes":[]}
```

Todos os arrays desse preimage são conjuntos canônicos. Duas propostas que
normalizam para o mesmo preimage produzem uma hipótese, preservando as duas
derivações somente se forem iguais; derivações diferentes permanecem hipóteses
distintas. Colisão de ID com bytes canônicos diferentes é `INTEGRITY_ERROR`.

## 7. API pública L1 de síntese

```go
type ImprovementHypothesisStatus string
const (
    HypothesisProposed     ImprovementHypothesisStatus = "PROPOSED"
    HypothesisSupported    ImprovementHypothesisStatus = "SUPPORTED"
    HypothesisConflicting  ImprovementHypothesisStatus = "CONFLICTING"
    HypothesisInconclusive ImprovementHypothesisStatus = "INCONCLUSIVE"
)

type ImprovementSynthesisStatus string
const (
    ImprovementComplete     ImprovementSynthesisStatus = "COMPLETE"
    ImprovementPartial      ImprovementSynthesisStatus = "PARTIAL"
    ImprovementInconclusive ImprovementSynthesisStatus = "INCONCLUSIVE"
)

type TransferMode string
const (
    TransferComponentReuse TransferMode = "COMPONENT_REUSE"
    TransferPatternAdapt   TransferMode = "PATTERN_ADAPTATION"
    TransferReferenceOnly  TransferMode = "REFERENCE_ONLY"
    TransferUnknown        TransferMode = "UNKNOWN"
)

type HypothesisDerivation string
const (
    DerivationRuleDerived HypothesisDerivation = "RULE_DERIVED"
    DerivationExternal    HypothesisDerivation = "EXTERNAL_AGENT_PROPOSAL"
)

type TargetEvidenceState string
const (
    TargetPresent TargetEvidenceState = "PRESENT"
    TargetAbsent  TargetEvidenceState = "ABSENT"
    TargetPartial TargetEvidenceState = "PARTIAL"
    TargetUnknown TargetEvidenceState = "UNKNOWN"
)

type ImprovementLicenseAssessment string
const (
    ImprovementLicenseCompatible   ImprovementLicenseAssessment = "COMPATIBLE"
    ImprovementLicenseIncompatible ImprovementLicenseAssessment = "INCOMPATIBLE"
    ImprovementLicenseUnrecognized ImprovementLicenseAssessment = "UNRECOGNIZED"
)

type CounterEvidenceSearchStatus string
const (
    CounterEvidenceFound           CounterEvidenceSearchStatus = "FOUND"
    CounterEvidenceNoneInScope     CounterEvidenceSearchStatus = "NONE_FOUND_IN_SCOPE"
    CounterEvidenceBudgetExhausted CounterEvidenceSearchStatus = "BUDGET_EXHAUSTED"
)

type SourceArtifactRef struct {
    RunID    RunID          `json:"run_id"`
    RunKind  RunKind        `json:"run_kind"`
    Artifact RunArtifactName `json:"artifact"`
    SHA256   string         `json:"sha256"`
}

type ImprovementProposal struct {
    ProposalID            string               `json:"proposal_id"`
    Title                 string               `json:"title"`
    TargetProperty        string               `json:"target_property"`
    ExpectedEffect        string               `json:"expected_effect"`
    MechanismRefs         []string             `json:"mechanism_refs"`
    ProductObservationRefs []string            `json:"product_observation_refs"`
    Assumptions           []string             `json:"assumptions"`
    Derivation            HypothesisDerivation `json:"derivation"`
    Rationale             string               `json:"rationale"`
}

type ValidationStep struct {
    Question          string `json:"question"`
    Observable        string `json:"observable"`
    SafeProcedure     string `json:"safe_procedure"`
    PassCriterion     string `json:"pass_criterion"`
    RejectCriterion   string `json:"reject_criterion"`
    RelativeCost      string `json:"relative_cost"`
    MissingEvidence   string `json:"missing_evidence"`
}

type CounterEvidenceSearch struct {
    Status     CounterEvidenceSearchStatus `json:"status"`
    ScopeRefs  []string                    `json:"scope_refs"`
    ResultRefs []string                    `json:"result_refs"`
}

type ImprovementHypothesis struct {
    ID                     string                      `json:"id"`
    Title                  string                      `json:"title"`
    TargetProperty         string                      `json:"target_property"`
    TargetState            TargetEvidenceState         `json:"target_state"`
    TargetGapRefs          []string                    `json:"target_gap_refs"`
    ExpectedEffect         string                      `json:"expected_effect"`
    MechanismRefs          []string                    `json:"mechanism_refs"`
    ProductObservationRefs []string                    `json:"product_observation_refs"`
    TransferMode           TransferMode                `json:"transfer_mode"`
    License                ImprovementLicenseAssessment `json:"license"`
    Assumptions            []string                    `json:"assumptions"`
    Preconditions          []string                    `json:"preconditions"`
    IntegrationUnknowns    []string                    `json:"integration_unknowns"`
    Risks                  []string                    `json:"risks"`
    Tradeoffs              []string                    `json:"tradeoffs"`
    EvidenceRefs           []string                    `json:"evidence_refs"`
    CounterEvidenceRefs    []string                    `json:"counter_evidence_refs"`
    CounterEvidenceSearch  CounterEvidenceSearch       `json:"counter_evidence_search"`
    ValidationPlan         []ValidationStep             `json:"validation_plan"`
    Derivation             HypothesisDerivation         `json:"derivation"`
    Status                 ImprovementHypothesisStatus  `json:"status"`
    AdoptionDisposition    string                       `json:"adoption_disposition"`
}

type ImprovementSynthesisPolicy struct {
    MaxHypotheses                 int  `json:"max_hypotheses"`
    MaxMechanismsPerHypothesis    int  `json:"max_mechanisms_per_hypothesis"`
    RequireCounterEvidenceSearch  bool `json:"require_counterevidence_search"`
}

type ImprovementSynthesisRequest struct {
    Idea            IdeaSpec                    `json:"idea"`
    SourceArtifacts []SourceArtifactRef         `json:"source_artifacts"`
    Proposals       []ImprovementProposal       `json:"proposals"`
    Policy          ImprovementSynthesisPolicy  `json:"policy"`
}
```

`AdoptionDisposition` é sempre `NOT_AN_ADOPTION_DECISION`. `RelativeCost` aceita
somente `LOW`, `MEDIUM`, `HIGH`, `UNKNOWN`. O plano é dado, não executor: não cria
branch, arquivo, código, experimento, PR, instalação ou chamada externa.

Política inicial não permite elevar mecanismo `PARTIAL` a hipótese `SUPPORTED`.
Mudar esse teto exige nova revisão do contrato e novo selo.

### 7.1 Portas e operações públicas exatas

```go
type MechanismGap struct {
    Code          string                `json:"code"`
    Source        ImmutableCandidateRef `json:"source"`
    RuleID        string                `json:"rule_id"`
    WitnessHash   string                `json:"witness_hash"`
    StatusCeiling MechanismEvidenceStatus `json:"status_ceiling"`
}

type MechanismExtractionInput struct {
    Bundle               EvidenceBundle
    Assessments          []CandidateAssessment
    ConventionalBaseline *ConventionalBaseline
    RuleCatalogJSON      []byte
    RuleCatalogSHA256    string
}

type MechanismExtractionResult struct {
    Mechanisms []ArchitecturalMechanism `json:"mechanisms"`
    Gaps       []MechanismGap           `json:"gaps"`
}

type MechanismExtractor interface {
    ExtractMechanisms(context.Context, MechanismExtractionInput) (MechanismExtractionResult, error)
}

type ImprovementSynthesisStartReceipt struct {
    SchemaVersion string                  `json:"schema_version"`
    RunID         RunID                   `json:"run_id"`
    RunKind       RunKind                 `json:"run_kind"`
    Lifecycle     RecommendationLifecycle `json:"lifecycle"`
    Version       uint64                  `json:"version"`
}

type ImprovementSynthesisStarter interface {
    StartImprovementSynthesis(context.Context, string, ImprovementSynthesisRequest) (ImprovementSynthesisStartReceipt, error)
}

type ImprovementSynthesisExecutor interface {
    ExecuteImprovementSynthesis(context.Context, RunID) error
}
```

L2 depende somente de `ImprovementSynthesisStarter`; o dispatcher/worker depende
somente de `ImprovementSynthesisExecutor`. Get, artifact e cancel continuam nas
portas já certificadas. `RecommendationRunStore.Get` e `ReadArtifact` são as únicas
operações necessárias para resolver sources; não há nova porta ou implementação
L3. `Bundle` e `Assessments` pertencem ao mesmo run e são consumidos em memória;
assessment ausente, repository divergente ou source não ligada produz gap e nunca
fallback em conteúdo bruto. `RuleCatalogJSON` é tratado como input imutável,
validado pelo hash exato e nunca publicado cru. O nome de structs concretos, seus
campos privados e construtores não são observáveis; as assinaturas das interfaces
acima são.

## 8. Source refs, validação pré-fila e hash do pedido

`SourceArtifactRef.SHA256` é exatamente 64 hex minúsculos sem prefixo. Ele é
comparado aos 64 hex após `sha256:` do índice publicado. Pares permitidos:

| `run_kind` | `artifact` |
|---|---|
| `RECOMMENDATION` | `mechanism_catalog`, `mechanism_gaps` |
| `MARKET_DISCOVERY` | `market_candidates`, `market_snapshot` |

Ao menos uma source e no máximo 12 são obrigatórias. Duplicata da tupla
`(run_id,run_kind,artifact)` é inválida, mesmo com hash diferente. Síntese não
aceita source `IMPROVEMENT_SYNTHESIS`, path, URL, texto substituto nem referência
a outro store. A autoridade local é a presença no MESMO store injetado para o
servidor corrente; não existe lookup global, fallback remoto ou credencial trazida
pelo pedido.

Antes de `ResolveIdempotency`, reserva de fila, geração de ID ou criação de run,
L1 executa nesta ordem:

1. decode estrito, tipos, limites, campos desconhecidos, controles e enums;
2. normalização de ideia, propostas, sets e política;
3. validação dos pares da allowlist e ordenação de sources por
   `(run_id,run_kind,artifact,sha256)`;
4. para cada source, no MESMO `RecommendationRunStore` injetado: `Get`, conferir
   `run_kind`, conferir lifecycle terminal, localizar índice `PUBLISHED`, comparar
   hash informado, ler todas as páginas pelo snapshot hash/version e validar o
   documento canônico;
5. construir índices locais de mechanism IDs, gap IDs e product observation IDs;
   toda referência de proposta deve existir em uma source declarada;
6. validar que conteúdo produto-only aparece apenas em
   `product_observation_refs` e mecanismo apenas em `mechanism_refs`;
7. computar o hash causal e só então aplicar o algoritmo certificado de
   idempotência/reserva/criação do Passo 0026.

Terminais são `SUCCEEDED|FAILED|CANCELLED|INTERRUPTED`. Source terminal diferente
de `SUCCEEDED` pode ser usada somente se o artifact está `PUBLISHED` e íntegro;
ela produz warning `SOURCE_RUN_NOT_SUCCEEDED` e limita toda hipótese dependente a
`INCONCLUSIVE`. Run não terminal ou artifact histórico ausente/não publicado
produz `SOURCE_BASIS_UNAVAILABLE` antes da fila. Em particular, run
`RECOMMENDATION` histórico sem `mechanism_catalog` nunca relê checkout, HEAD,
provider ou bundle.

Run inexistente produz `RUN_NOT_FOUND`; kind declarado diferente do persistido ou
par não allowlisted é `Invalid params`; hash cliente divergente, blob corrompido,
schema inválido ou referência interna forjada é `INTEGRITY_ERROR`. Nenhum desses
casos cria chave, run ou reserva.

Normalização de `IdeaSpec` é a já certificada. Todo outro texto usa NFC,
whitespace colapsado, rejeição de C0/C1/DEL e limites do schema. Sets ordenam por
bytes; propostas ordenam por `proposal_id`; IDs duplicados após normalização são
inválidos. Ordem original não é semântica.

O hash exclui `idempotency_key` e inclui o catálogo:

```text
request_hash = "sha256:" + hex(SHA-256(json.Marshal(canonical_request_envelope)))
canonical_request_envelope = {
  schema_version: "1",
  idea: normalized_idea,
  source_artifacts: sorted_sources,
  proposals: sorted_normalized_proposals,
  policy: normalized_policy,
  mechanism_rules_sha256: "a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5"
}
```

A mesma chave e o mesmo hash devolvem o mesmo `run_id` e estado corrente. Mesma
chave com source, proposta, ideia, policy ou catálogo diferente produz
`IDEMPOTENCY_CONFLICT`. A chave conserva exatamente as regras de bytes/digest do
Passo 0026.

## 9. Tool MCP 10 e schema exato

A entrada anunciada em `tools/list` é EXATAMENTE:

```json
{
  "name": "bateia_start_improvement_synthesis",
  "description": "Inicia síntese durável de hipóteses de melhoria a partir de mecanismos e observações imutáveis",
  "inputSchema": {
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
          "problem": {"type":"string","minLength":1,"maxLength":1000,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},
          "users": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},
          "capabilities": {"type":"array","minItems":1,"maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},
          "constraints": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},
          "non_goals": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},
          "preferred_stack": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":200,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}}
        },
        "required": ["problem","capabilities"]
      },
      "source_artifacts": {
        "type": "array",
        "minItems": 1,
        "maxItems": 12,
        "items": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "run_id": {"type":"string","pattern":"^run_[0-9a-f]{32}$"},
            "run_kind": {"type":"string","enum":["RECOMMENDATION","MARKET_DISCOVERY"]},
            "artifact": {"type":"string","enum":["mechanism_catalog","mechanism_gaps","market_candidates","market_snapshot"]},
            "sha256": {"type":"string","pattern":"^[0-9a-f]{64}$"}
          },
          "required": ["run_id","run_kind","artifact","sha256"]
        }
      },
      "proposals": {
        "type": "array",
        "maxItems": 20,
        "items": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "proposal_id": {"type":"string","minLength":1,"maxLength":128,"pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$"},
            "title": {"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},
            "target_property": {"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},
            "expected_effect": {"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"},
            "mechanism_refs": {"type":"array","maxItems":4,"items":{"type":"string","pattern":"^mech_[0-9a-f]{64}$"}},
            "product_observation_refs": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":256,"pattern":"^[A-Za-z0-9][A-Za-z0-9._:-]{0,255}$"}},
            "assumptions": {"type":"array","maxItems":20,"items":{"type":"string","minLength":1,"maxLength":512,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}},
            "derivation": {"type":"string","const":"EXTERNAL_AGENT_PROPOSAL"},
            "rationale": {"type":"string","minLength":1,"maxLength":1000,"pattern":"^[^\\u0000-\\u001F\\u007F-\\u009F]+$"}
          },
          "required": ["proposal_id","title","target_property","expected_effect","mechanism_refs","product_observation_refs","assumptions","derivation","rationale"]
        }
      },
      "policy": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "max_hypotheses": {"type":"integer","minimum":1,"maximum":12},
          "max_mechanisms_per_hypothesis": {"type":"integer","minimum":1,"maximum":4},
          "require_counterevidence_search": {"type":"boolean","const":true}
        },
        "required": ["max_hypotheses","max_mechanisms_per_hypothesis","require_counterevidence_search"]
      }
    },
    "required": ["idempotency_key","idea","source_artifacts","proposals","policy"]
  }
}
```

Além do schema, runtime rejeita `..` na chave, duplicatas normalizadas, pares de
source incompatíveis, refs não resolvidas e qualquer URL, path, prompt, expressão
de ferramenta, peso de popularidade, comando ou campo adicional. `proposals:[]`
é válido e habilita somente derivação pura `RULE_DERIVED`.

Sucesso usa um item textual MCP com JSON compacto exato:

```json
{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","run_kind":"IMPROVEMENT_SYNTHESIS","lifecycle":"QUEUED","version":1}
```

O envelope é `{"content":[{"type":"text","text":"<JSON compacto>"}]}`. A tool
não espera worker nem retorna hipótese inline.

Erros 0026 permanecem. Adição exata:

| Código | `message` exata | Situação |
|---:|---|---|
| `-32028` | `SOURCE_BASIS_UNAVAILABLE` | run não terminal ou artifact source permitido ausente/não publicado |

`RUN_NOT_FOUND`, `IDEMPOTENCY_CONFLICT`, `QUEUE_FULL`, `QUOTA_EXCEEDED`,
`INTEGRITY_ERROR`, `RUN_STORE_UNAVAILABLE` e `Invalid params` conservam códigos e
mensagens anteriores. Erro não ecoa proposal, content, hash cliente, path,
canário nem detalhe físico.

## 10. Runtime durável, lifecycle e artefatos

`RunKind` acrescenta `IMPROVEMENT_SYNTHESIS`. `RecommendationRunView` contém
sempre os três campos nullable:

```go
RecommendationStatus *IdeaRunStatus              `json:"recommendation_status"`
MarketDiscoveryStatus *MarketDiscoveryStatus     `json:"market_discovery_status"`
ImprovementStatus     *ImprovementSynthesisStatus `json:"improvement_status"`
```

Exatamente um corresponde ao `run_kind`; os outros são `null`. Antes do terminal,
todos são `null`. `SUCCEEDED` significa apenas execução concluída e publicação
íntegra; não implica `SUPPORTED` nem `COMPLETE`.

Síntese usa cinco marcos, nesta ordem:

1. `synthesis_compare`;
2. `synthesis_derive`;
3. `synthesis_counterevidence`;
4. `synthesis_validate`;
5. `synthesis_publish`.

Estado inicial conserva `QUEUED`, versão 1 e stage `queued`, mas progress é
`{planned:5,completed:0}`. Ao rodar, progress cresce monotonicamente até 5; stage
terminal é `terminal`. A máquina de lifecycle, CAS, cancelamento, corrida terminal,
recovery e retry é exatamente a do Passo 0026. Cancelamento observado impede novo
passo e nenhuma análise chama rede.

Runs `RECOMMENDATION` mantêm os oito estágios certificados. Extração ocorre
dentro de `publish`, sem alterar `planned`, ranking ou status. Para run novo sob
0028, `SUCCEEDED` requer também `mechanism_catalog` e `mechanism_gaps` publicados;
um catálogo vazio com gaps explícitos é válido. Run histórico não muda.

Publicação de recommendation preserva os seis artifacts basais e acrescenta, no
fim, `mechanism_catalog` e depois `mechanism_gaps`. Publicação de synthesis segue
EXATAMENTE `architecture_comparison`, `improvement_hypotheses`, `hypothesis_gaps`,
`validation_plans`, `synthesis_provenance`. Índices MCP são sempre projetados na
ordem global da allowlist abaixo, independentemente da ordem física de escrita.

A enum de `bateia_get_artifact` fica EXATAMENTE, nesta ordem:

1. `request`
2. `recommendation`
3. `decision`
4. `alternatives`
5. `failures`
6. `provenance`
7. `connector_capabilities`
8. `market_plan`
9. `market_snapshot`
10. `market_candidates`
11. `market_gaps`
12. `mechanism_catalog`
13. `mechanism_gaps`
14. `architecture_comparison`
15. `improvement_hypotheses`
16. `hypothesis_gaps`
17. `validation_plans`
18. `synthesis_provenance`

Documentos conservam bytes canônicos:

```json
{"schema_version":"1","artifact":"improvement_hypotheses","entries":[]}
```

Hash, size, índice, paginação, cursor HMAC, quota e publicação atômica são os de
0026. Cada artifact só é válido para seu `run_kind`; nome globalmente allowlisted
mas incompatível com o run retorna `ARTIFACT_NOT_READY`, nunca artifact vazio
inventado. Ordenação e conteúdo:

| Artifact | Entradas/ordem | Conteúdo seguro obrigatório |
|---|---|---|
| `mechanism_catalog` | mecanismo por `id` | `ArchitecturalMechanism` integral seguro |
| `mechanism_gaps` | `(code,source_repository,source_commit,rule_id,witness_hash)` | IDs/hashes, reason code e teto; sem texto bruto |
| `architecture_comparison` | `(target_property,mechanism_id,source_run_id)` | target state/ref, mecanismo/source preservados, novelty e unknowns; sem score de popularidade |
| `improvement_hypotheses` | hipótese por `id` | estrutura da seção 7 e disclaimer machine-readable |
| `hypothesis_gaps` | `(hypothesis_id,code,source_ref)` | falta de evidência, alvo, licença, compatibilidade, budget ou proposta inválida |
| `validation_plans` | `(hypothesis_id,step_index)` | passos na ordem causal de cada hipótese |
| `synthesis_provenance` | `(source_run_id,artifact,sha256)` e depois `(hypothesis_id,derivation)` | refs resolvidas, catálogo hash, rule IDs e escopos de busca |

`architecture_comparison` nunca achata mecanismos de sources distintas.
`synthesis_provenance` não contém idempotency key, proposta/rationale integral,
query, body, snippet, prompt ou configuração privada.

Antes de `SUCCEEDED`, os cinco artifacts de síntese estão `PUBLISHED`.
`ImprovementSynthesisStatus` é:

- `COMPLETE`: todas as sources/propostas admissíveis foram avaliadas, todas as
  buscas exigidas concluíram dentro do escopo, nenhum gap de budget/integridade e
  todos os artifacts íntegros; pode haver zero hipótese `SUPPORTED`;
- `PARTIAL`: ao menos uma hipótese válida foi materializada e há gap não fatal de
  source terminal degradada, proposta, alvo, licença, compatibilidade ou budget;
- `INCONCLUSIVE`: nenhuma hipótese avaliável, integridade epistemológica
  insuficiente, somente outcomes Unknown ou counterevidence não procurada por
  budget.

## 11. Regras exatas de hipótese

Propostas externas são conteúdo não confiável. O serviço resolve refs e recalcula
status; nunca aceita status, evidência, licença, compatibilidade ou adoção trazidos
pelo proponente. `RULE_DERIVED` usa somente correspondência literal entre
propriedade/lacuna taxonômica congelada e `purpose`/classe do catálogo; não usa
frequência, ranking, embedding ou similaridade opaca.

Precedência do status de hipótese:

1. contraevidência que derruba lacuna, causalidade ou compatibilidade central,
   ou licença `INCOMPATIBLE` com `COMPONENT_REUSE`, produz `CONFLICTING`;
2. source terminal não sucedida, mecanismo `INCONCLUSIVE`, busca de
   contraevidência `BUDGET_EXHAUSTED`, ref epistemologicamente opaca ou baseline
   de propriedade necessário e ausente produz `INCONCLUSIVE`;
3. `SUPPORTED` exige simultaneamente: todo mecanismo referido `PROVEN`; target
   `ABSENT` comprovado ou `PARTIAL/PRESENT` com medição basal própria; relação
   causal explícita; ao menos uma evidence ref do alvo e uma de cada mecanismo;
   assumptions, preconditions, risks, tradeoffs e busca de contraevidência
   registrados; nenhuma contraevidência central; todas as refs/hash resolvidas;
   plano de validação não vazio e falsificável; licença compatível quando modo é
   `COMPONENT_REUSE`;
4. caso admissível restante é `PROPOSED`.

Target `UNKNOWN` nunca vira `ABSENT` e limita a `PROPOSED`, salvo opacidade/budget
que exige `INCONCLUSIVE`. Sem repositório/relatório alvo, a hipótese de produto é
no máximo `PROPOSED` e o primeiro ValidationStep mede o estado atual.

Composição de dois ou mais mecanismos sempre preserva source de cada um e inclui,
até testemunha específica, pelo menos `interface compatibility`, `state model
compatibility`, `failure model compatibility`, `performance compatibility` e
`license compatibility` em `integration_unknowns`. Partes individualmente
`PROVEN` não provam composição.

Licença:

- `UNRECOGNIZED` conserva mecanismo/hipótese, emite
  `LICENSE_NOT_RECOGNIZED`, proíbe promoção automática a `COMPONENT_REUSE` e
  permite no máximo `PATTERN_ADAPTATION|REFERENCE_ONLY|UNKNOWN`;
- `INCOMPATIBLE` fica visível, bloqueia reúso direto e conflita hipótese que o
  exija;
- `COMPATIBLE` apenas permite avaliar `COMPONENT_REUSE`; não aprova adoção;
- abstração de padrão não autoriza copiar código, texto, asset ou interface.

`NONE_FOUND_IN_SCOPE` significa busca limitada sem resultado, nunca prova de
segurança ou ausência. `ScopeRefs` é não vazio e lista superfícies realmente
consultadas. Cada hipótese preservada tem assumptions, risks e validation plan
não vazios; preconditions, tradeoffs, unknowns e counterevidence refs são arrays
explícitos, ainda que vazios.

## 12. Limites, segurança e conteúdo não confiável

- strings de nome/título/propriedade/efeito: 512 bytes UTF-8; propósito,
  estratégia, rationale e critérios: 1000; itens de sets: 512; arrays: no máximo
  20, salvo limites menores expressos;
- mecanismo tem no máximo 64 witnesses e 32 itens por campo estrutural; hipótese
  tem no máximo 4 mecanismos e 20 validation steps;
- texto de repo, produto ou proposta nunca controla tool, próxima chamada,
  prompt, política, host, header, path, filesystem ou classificação;
- URL não é seguida durante síntese; path de witness é referência opaca e não é
  aberto;
- nenhum clone, fetch, instalação, compra, build ou execução de terceiro;
- nenhuma escrita no software alvo e nenhum efeito externo;
- nenhum artifact/log/erro contém segredo, token, cookie, header, dado pessoal,
  review, body, screenshot, artwork, prompt privado, snippet, arquivo integral ou
  canário;
- erros externos viram enum/mensagem pública controlada; conteúdo contraditório
  não é removido para melhorar status;
- `MineOntologicalImprovements`, `ImprovementRecommendation`, `AdoptionRate`,
  estrelas, downloads, ratings e ranking legado não entram em evidência, hash,
  novelty, status ou ordenação;
- L1 não faz rede/filesystem/env/relógio/goroutine/mutex; L2 não importa L3/L4;
  L3 existente não importa L2/L4; L4 somente compõe.

## 13. Política total de Unknown

| Situação | Resultado obrigatório |
|---|---|
| run histórico sem catálogo | `SOURCE_BASIS_UNAVAILABLE`; sem releitura |
| source não terminal | `SOURCE_BASIS_UNAVAILABLE`; sem fila |
| hash divergente/corrupção | `INTEGRITY_ERROR`; sem fila |
| identidade sem commit imutável | gap; nunca mecanismo |
| classe fora do catálogo | `MECHANISM_CLASS_UNKNOWN` gap |
| testemunha mínima incompleta | `PARTIAL` ou gap; nunca `PROVEN` |
| documentação sem implementação | no máximo `PARTIAL` |
| produto sem repo | observação; nunca mecanismo |
| target `UNKNOWN` | nunca `ABSENT`; hipótese no máximo `PROPOSED` |
| medição basal ausente | efeito ainda não mensurável e plano obrigatório |
| compatibilidade de composição desconhecida | `integration_unknowns`; nunca compatível |
| licença não reconhecida | warning e bloqueio de reúso direto; sem exclusão |
| counterevidence não procurada por budget | `INCONCLUSIVE` |
| proposta com ref inválida | `Invalid params`; não alucinar ref |
| zero hipótese suportada | resultado válido, não sucesso fabricado |
| lifecycle `SUCCEEDED` | não implica `SUPPORTED` ou `COMPLETE` |

`Unknown` nunca é convertido por default em ausência, prova, compatibilidade,
segurança, permissão, empty result, mecanismo negativo ou mutação rejeitada.
Somente casos opacos marcados antes do candidato podem terminar `Unknown` no gate,
e são excluídos do numerador e denominador do mutation score.

## 14. Replay determinístico

Replay recebe somente captures/artifacts congelados e hashes registrados, injeta
o catálogo exato e não possui transporte de rede. Duas execuções com as mesmas
entradas produzem bytes idênticos para todos os sete artifacts novos, seus
índices, IDs, hashes, status, gaps e provenance. Reordenação de bundles, witnesses,
sources, proposals ou respostas semanticamente conjuntistas não altera bytes.

Relógio, `observed_at`, path local, map iteration, PID, worker scheduling, ordem de
discovery e erro livre não entram. Se capture, artifact, catálogo ou versão
divergir, replay produz `INTEGRITY_ERROR`, nunca tenta reparar, reler rede ou usar
HEAD atual.

## 15. Observáveis e gate discriminatório

Oráculos independentes devem observar no mínimo:

1. dez tools exatas e posições 1–9 preservadas;
2. schema e runtime estritos da tool 10;
3. mecanismo `PROVEN` somente com implementação + witness independente na mesma
   source imutável;
4. documentação isolada no máximo `PARTIAL`;
5. produto-only nunca mecanismo e proposta externa nunca evidência;
6. cross-repo/cross-commit nunca combinado;
7. regra, source e witness hashes em todo mecanismo;
8. toda source avaliada liga exatamente um item/hash do bundle; RuleID/Symbol/Line
   nunca são inferidos de Content;
9. ordem de bundle/assessments/sources não altera ID/catálogo;
10. histórico sem catálogo não relê fonte;
11. source kind/lifecycle/allowlist/hash verificados antes da fila;
12. mesma chave/pedido reutiliza run e mudança causal conflita;
13. target Unknown não vira gap e ausência de target não vira `SUPPORTED`;
14. composição não vira compatível sem testemunha;
15. licença desconhecida preserva e bloqueia somente reúso direto;
16. licença incompatível e contraevidência permanecem visíveis;
17. popularidade, `AdoptionRate` e novelty não alteram status/ordem;
18. ausência de contraevidência é busca limitada;
19. lifecycle/status/hypothesis status não se confundem;
20. cancelamento impede novo passo;
21. artifacts, logs e erros não vazam conteúdo/canários;
22. nenhuma saída executa ou modifica alvo;
23. replay offline é byte-idêntico;
24. mesmo store/queue/root/secret/worker/idempotência por identidade de instância;
25. três refinamentos predecessores têm receipts próprios.

Casos positivos válidos devem resultar `Preserved`; todas as mutações negativas
válidas devem resultar `Violated` com testemunha; somente opacos previamente
designados podem resultar `Unknown`. Exigir:

```text
mutation_score = mutações negativas válidas rejeitadas / mutações negativas válidas = 1.0
```

Qualquer sobrevivente impede selo. Repetir e reordenar o gate deve produzir o
mesmo resultado. `Unknown` nunca satisfaz mutação negativa.

Mutações mínimas incluem todas as listadas na spec, em especial: README sozinho
`PROVEN`; produto vira arquitetura; proposta cria evidence; cross-source; commit
mutável; Unknown vira absent; frequência/novelty/AdoptionRate elevam status;
ausência de contraevidência vira segurança; composição vira compatível; licença
desconhecida exclui ou permite reuse; licença incompatível some; source hash não
é verificado; histórico relê HEAD; `SUCCEEDED` vira `SUPPORTED`; raw/canary vaza;
URL/path provoca I/O; runtime paralelo; cancelamento trabalha; replay acessa rede;
tool 10 inserida antes; refinamento alegado sem receipt; RuleID inventado a partir
de `CollectedEvidence.Content`; source avaliada sem binding/mesmo hash aceita;
source `PARTIAL` promove mecanismo a `PROVEN`.

## 16. Gates finais reservados ao verificador

Depois do selo e do candidato, e não por esta autoridade, devem passar:

```text
go test -count=1 ./...                         # duas execuções separadas
go test -race -count=1 ./01_core/... ./02_shell/mcp ./03_infra/... ./04_wiring
go vet ./...
go build -o /tmp/bateia-step0028 ./04_wiring
crystalline-lint .
git diff --check
```

O build é deliberadamente do pacote multifile `./04_wiring`; build por
`./04_wiring/main.go` é gate incorreto para 0028.

Também: validação strict JSON/hashes; suites protegidas; métricas da spec;
precisão `PROVEN` 100%; recall esperado e divergente >=80% onde aplicável; zero
produto-only, Unknown-as-absence ou compatibilidade inventada; 100% provenance,
planos e buscas; replay duas vezes por comparação de bytes; testes unitários sem
rede/modelo; scan de canários/raw content; auditoria de imports e I/O; identidade
única de runtime; hashes selados intactos; receipts de refinamento 0025, 0026 e
0027.

Somente o verificador pode emitir `MECHANISM_IMPROVEMENT_SYNTHESIS_READY`, limitado
aos hashes, versões, corpus, rules e ambiente registrados. Sem isolamento forte,
a atestação final continua obrigatoriamente
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.

## 17. Ownership e prompts causais

Cada arquivo de produção novo possui um prompt proprietário único:

| Camada/arquivo futuro | Prompt L0 | SHA-256 |
|---|---|---|
| `01_core/domain/architectural_mechanism.go` | `00_nucleo/prompts/01_core/domain/architectural_mechanism.md` | `cc9d9b7f7af45410a6de03ac90becdad51549e6b78cbd6f188882ec41bad7540` |
| `01_core/domain/improvement_synthesis.go` | `00_nucleo/prompts/01_core/domain/improvement_synthesis.md` | `eaf502e438f1d79cd95306986a1e1af5e74b943e73d7aa77f7f52a8c6736c9d0` |
| `01_core/usecase/mechanism_extraction.go` | `00_nucleo/prompts/01_core/usecase/mechanism_extraction.md` | `d3f3d58d9dbb5609b3d946442326ccb5474c1dda3910f98d84de90e40b6d62e2` |
| `01_core/usecase/improvement_synthesis.go` | `00_nucleo/prompts/01_core/usecase/improvement_synthesis.md` | `ffd9d4dca7d5d8a63b25872a7605fa719a920405ce10f10541c27fab8eee040a` |
| `02_shell/mcp/improvement_synthesis.go` | `00_nucleo/prompts/02_shell/mcp/improvement_synthesis.md` | `43c5be8d6aac9f251c88caa3cdd1152070343bbde9cb5d564210067901e9c263` |
| `04_wiring/improvement_synthesis.go` | `00_nucleo/prompts/04_wiring/improvement_synthesis.md` | `7f98d11406321a1a402d1b95e2b56ea6310a8cd0a63bac0b56ab73e0f41ae890` |

Arquivos existentes de run, recommendation publish, registry/getters MCP e
dispatcher conservam sua linhagem L0 anterior byte a byte; a implementação deve
ser mínima e aditiva, referenciar o contrato selado no receipt e obter receipts de
refinamento. Nenhum prompt ou artifact dos Passos 0025–0027 é alterado por esta
autoridade. Prompts não são testes nem alternativa ao contrato.

## 18. Ambiguidades resolvidas e remanescentes

Resolvido normativamente nesta revisão:

- R3 renomeia somente o enum escalar de licença da síntese para
  `ImprovementLicenseAssessment` e prefixa seus identificadores Go; as strings
  JSON e toda a semântica permanecem iguais, enquanto o struct basal
  `LicenseAssessment` continua intocado;
- R2 separa o bundle bruto da avaliação: `CollectedEvidence` fornece somente a
  âncora/hash; RuleID, Symbol, Line, source kind e status vêm de
  `CandidateAssessment`/`EvidenceSource` já avaliados e ligados univocamente;
- `CollectedEvidence.Content` serve somente para conferir SHA-256 e nunca para
  inferir signal, mecanismo, regra ou testemunha;
- source `sha256` usa 64 hex sem prefixo; artifact index continua
  `sha256:<64-hex>`;
- terminal degradado com artifact íntegro é aceito com teto `INCONCLUSIVE`;
  não terminal falha pré-fila;
- histórico sem catálogo retorna `SOURCE_BASIS_UNAVAILABLE` e nunca relê fonte;
- política inicial não eleva mecanismo `PARTIAL` a `SUPPORTED`;
- novidade `DIVERGENT` exige mecanismo `PROVEN` e nunca dá bônus;
- regra desconhecida produz gap, não classificação aproximada;
- extração em recommendation ocupa o estágio `publish`, preservando os oito
  marcos 0026;
- build correto é o pacote multifile `./04_wiring`.

Ambiguidades epistemológicas preservadas, nunca convertidas em sucesso:

- suficiência semântica de testemunhas que usam regras avaliadas catalogadas mas não
  expõem estrutura bastante fica `PARTIAL|INCONCLUSIVE`;
- compatibilidade real entre mecanismos, custo operacional, segurança e efeito no
  alvo permanecem para validação;
- interpretação jurídica de licença permanece fora do Bateia;
- guarantees de fsync/rename em filesystem não exercitado permanecem caso opaco
  do contrato 0026.

## 19. Suplemento normativo append-only R4

**Revisão ativa:** R4 (substitui R3 e o catálogo v2 antes de candidato e selo)

**Modo da revisão:** append-only; as seções R3 acima permanecem como registro
histórico e são normativas somente onde este suplemento não as substitui

**Motivo:** resolver a impossibilidade pública
`DOCUMENTATION_ONLY -> PARTIAL` registrada antes do candidato

**API pública alterada:** nenhuma

**Implementação candidata lida:** não

**Receipt, selo, certificado ou veredito emitido por esta autoridade:** não

### 19.1 Raiz causal R4 e precedência

Esta tabela substitui a tabela da seção 1 para qualquer execução R4:

| Entrada causal R4 | SHA-256 |
|---|---|
| `00_nucleo/passos/0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md` | `9105b27e0e8e9165b0b133a647669e858b977ea2ab8c7d951ed0f217a26ac1a5` |
| `00_nucleo/benchmark/0028/manifest.json` | `0b8fbcebf9f3af491ae1851f59715c248fa63181617cc336dc39b6969fb555ad` |
| contrato predecessor R3 | `f8d8649a9bb3dd1707ae01f3cedc93648b1bd3c501d62a506f745b66a03c928b` |
| catálogo predecessor `tekt.mechanism-rules.v2` | `a09c6d7cc086848914e9d9c8f2f83f83ad91c35db7746608a377f31e6ef88cb5` |
| `00_nucleo/benchmark/0028/contract-inconsistency-o7.json` | `04e8a93a43cd447d82f9a1761d986ca1c7d3217e6588e4204f1f5e2343bc7a79` |
| catálogo ativo `00_nucleo/benchmark/0028/mechanism-rules.json`, `tekt.mechanism-rules.v3` | `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142` |
| `00_nucleo/ADR/ADR-0002-estrutura-de-pastas-e-camadas-tekt.md` | `b457fadd55d6328c2afce969b92a527b461edab9094b1bef66fcfebe71b79a68` |
| certificado 0025 | `bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d` |
| certificado 0026 | `a04d6b07a939c1ce5fb7fb184e5cb56cc3cc34ed0b2b54c8087415a18a4a99ed` |
| certificado 0027 | `acf05df728b05ccdc5bc15f78a18e49f19e743b2016ab1141e553ca2417c7e43` |
| prompt proprietário `00_nucleo/prompts/01_core/usecase/mechanism_extraction.md` | `18423bbbf5b6240f821a252b514887373f6bfb82659f8c89d9bdb0fa7f3376bb` |

Em caso de conflito textual, esta seção 19 prevalece sobre as seções 1, 5, 8,
15, 17 e 18. Todo o restante de R3 continua normativo. O catálogo ativo é
SOMENTE v3 e seu hash acima. O par R3/v2 não pode ser selado nem usado como
entrada de candidato, protocolo, ataque executável, oracle, replay ou gate.

### 19.2 Catálogo fechado v3 e schema estrito

`tekt.mechanism-rules.v3`, versão `3`, conserva EXATAMENTE as quatro classes da
seção 5, seus templates, limites, conflitos e requisitos de prova. Não adiciona
classe, alias, fallback lexical, fuzzy match, LLM, inferência por conteúdo ou
regra embutida. Contém 32 `accepted_evidence_rules` com `rule_id` globalmente
único, distribuídas nos papéis:

- 8 `IMPLEMENTATION`;
- 12 `INDEPENDENT`;
- 4 `DOCUMENTARY`;
- 8 `CONFLICT`.

Os únicos RuleIDs `DOCUMENTARY` são:

| Classe | RuleID explícito |
|---|---|
| `APPEND_ONLY_LOG_RECOVERY` | `append-only-architecture-documentation` |
| `BOUNDED_QUEUE_BACKPRESSURE` | `bounded-queue-architecture-documentation` |
| `CONTENT_ADDRESSED_STORAGE` | `content-address-architecture-documentation` |
| `PORT_ADAPTER_ISOLATION` | `adapter-boundary-architecture-documentation` |

O parser do catálogo é closed-world e DEVE rejeitar campo desconhecido,
ausente, tipo divergente, enum desconhecida, versão diferente, classe duplicada,
RuleID duplicado globalmente ou hash divergente. Cada item de
`accepted_evidence_rules` satisfaz EXATAMENTE uma destas formas:

1. `IMPLEMENTATION|INDEPENDENT|CONFLICT`: campos exatos `rule_id`, `role` e
   `witness_kind`, sem `status_max`;
2. `DOCUMENTARY`: campos exatos `rule_id`, `role`, `witness_kind` e
   `status_max`, com `witness_kind` EXATAMENTE `DOCUMENTATION` e `status_max`
   EXATAMENTE `PARTIAL`.

Nenhum outro papel é válido. `DOCUMENTARY` é positivo apenas no sentido limitado
da seção seguinte; não é alias nem subconjunto de `IMPLEMENTATION` ou
`INDEPENDENT`.

### 19.3 Regra normativa de documentação e status

O passo 7 do algoritmo da seção 5 passa a aceitar também a forma exata:
`role=DOCUMENTARY`, `witness_kind=DOCUMENTATION`,
`EvidenceSource.Kind=documentation`, `status_max=PARTIAL`. Qualquer outra
combinação desses quatro valores é incompatível e produz gap; não sofre coerção.

Depois do binding exato ao bundle, verificação de hash e preservação do status
pai, uma source `DOCUMENTARY` catalogada:

- conta como source positiva para satisfazer o mínimo `partial: 1` da sua classe;
- pode preencher apenas os templates seguros em `fillable_fields` do catálogo;
- materializa `PARTIAL` quando o status pai é `PROVEN|PARTIAL`, os campos centrais
  `purpose,strategy,boundaries` estão preenchidos e não há conflito central;
- conserva source imutável, RuleID e witness hash na proveniência;
- nunca conta como testemunha `IMPLEMENTATION`;
- nunca conta como testemunha `INDEPENDENT`;
- nunca participa dos mínimos `proven_implementation`, `proven_independent` ou
  `distinct_content_hashes` de `PROVEN`;
- nunca concede status maior que `PARTIAL`, ainda que o status pai seja `PROVEN`
  ou existam várias fontes documentais distintas.

Assim, `PROVEN` continua exigindo ao menos uma entrada `IMPLEMENTATION` e uma
entrada `INDEPENDENT` catalogadas, ambas sob status pai `PROVEN`, com dois hashes
distintos e todas as demais condições R3. Documentação junto de provas técnicas
pode permanecer na evidência segura, mas é neutra para esses dois mínimos.
Documentação não catalogada, README promocional, nome, tag, descrição, estrela,
download ou frequência continuam sem autoridade positiva e produzem gap quando
avaliados como regra desconhecida. Nenhuma classificação é inferida de
`CollectedEvidence.Content`, `Snippet` ou texto livre.

Precedência R4 de status, substituindo somente a lista correspondente na seção 5:

1. conflito central catalogado ou identidade inconsistente produz
   `INCONCLUSIVE` conforme R3;
2. source positiva sob `ABSENT|INCONCLUSIVE|CONFLICT` não conta; source positiva
   sob `PARTIAL` impõe teto `PARTIAL`;
3. `PROVEN` usa somente roles `IMPLEMENTATION` e `INDEPENDENT` e satisfaz todos
   os mínimos R3;
4. `PARTIAL` exige ao menos uma source positiva catalogada com role
   `IMPLEMENTATION|INDEPENDENT|DOCUMENTARY`, status pai `PROVEN|PARTIAL`, campos
   centrais preenchidos e nenhuma contraevidência central; `DOCUMENTARY` traz
   adicionalmente seu teto literal `PARTIAL`;
5. qualquer outro caso produz gap Unknown e nenhum mecanismo público.

As quatro classes, política total de Unknown, novidade, licença, síntese e teto
que proíbe mecanismo `PARTIAL` de sustentar hipótese `SUPPORTED` permanecem
inalterados.

### 19.4 Canonicalização, hashes e compatibilidade

A canonicalização e todos os preimages das seções 6 e 8 permanecem byte a byte,
salvo o valor causal do catálogo. O envelope do pedido R4 contém:

```text
mechanism_rules_sha256: "52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142"
```

Logo, uma idempotency key anteriormente associada a v2 não pode ser reutilizada
como se a semântica fosse a mesma: o hash causal difere e deve produzir
`IDEMPOTENCY_CONFLICT`. `mechanism_id` continua incorporando o RuleID efetivamente
usado; um mecanismo documental usa um dos quatro RuleIDs R4 explícitos.

Nenhum nome, tipo, enum, campo, tag JSON, interface, assinatura, tool, schema MCP,
ordem de tool, artifact, lifecycle, fila, store ou porta pública muda em R4.
Compatibilidade observável dos Passos 0025–0027 continua sujeita aos três receipts
de refinamento reservados ao verificador. O único prompt L0 de produção cujo
conteúdo normativo muda é:

| Prompt | R3 SHA-256 | R4 SHA-256 |
|---|---|---|
| `00_nucleo/prompts/01_core/usecase/mechanism_extraction.md` | `d3f3d58d9dbb5609b3d946442326ccb5474c1dda3910f98d84de90e40b6d62e2` | `18423bbbf5b6240f821a252b514887373f6bfb82659f8c89d9bdb0fa7f3376bb` |

Os outros cinco prompts proprietários da seção 17 permanecem com os hashes R3.

### 19.5 Controle positivo P02 e escopo causal de invalidação

O controle público P02 é satisfatível em R4: uma source documentation-only,
ligada exatamente ao bundle e a um dos quatro RuleIDs `DOCUMENTARY`, com status
pai `PROVEN|PARTIAL` e campos centrais do template preenchidos, deve materializar
um mecanismo `PARTIAL`, nunca `PROVEN`, e ser classificada `Preserved` pelo gate.
O controle não autoriza documentação genérica, RuleID inventado, inferência
lexical ou substituição de witness técnico.

R4/v3 invalida causalmente, desde o primeiro input alterado:

1. qualquer protocolo executável que fixe o contrato R3, catálogo v2, hash v2,
   enum de papéis sem `DOCUMENTARY` ou shape sem `status_max` documental;
2. todos os registros, subjects, configs, manifests e hashes de ataques
   executáveis derivados desse protocolo, inclusive P02;
3. toda calibração, corpus derivado, fixture, evaluator, suite, prompt de teste,
   matriz de cobertura, autoria, manifest e oracle que tenha lido ou codificado
   R3/v2;
4. todo resultado discriminatório, replay, diagnóstico, receipt ou seal derivado
   dessas entradas;
5. qualquer candidato futuro que tenha sido iniciado antes de novo selo R4/v3.

Protocolo, adversário e oracle devem ser regenerados de suas autoridades
segregadas a partir do contrato R4 final e do hash v3. O adversário deve manter
seus ataques independentes e reconstruir os subjects afetados; o oracle deve
reavaliar P02 sem ler ataques ou candidato; o verificador deve executar novamente
positivos, negativos e opacos completos. Nenhum resultado anterior pode ser
transportado por equivalência alegada. Mutation score exigido continua `1.0`,
`Unknown` continua fora do numerador/denominador somente nos opacos previamente
designados e nenhum selo pode preceder o novo gate completo.

Esta revisão resolve somente a inconsistência estrutural pública registrada em
`contract-inconsistency-o7.json`. Não emite ataque, oracle, receipt, selo,
certificado, laudo ou o veredito reservado.
