# Passo 0071 — Expor cobertura do catálogo sem atribuir lacunas ao candidato

## Estado e autoridade

- Estado: contrato candidato pré-implementação.
- Regime: materialização segregada Tekt completa.
- Papel desta autoria: autor do contrato e dos oráculos observáveis; sem acesso a
  implementação candidata, que ainda não existe.
- Baseline: commit `fa86a30be391e5a101aa708201858a96ca72e377` mais os arquivos
  públicos enumerados em **Entradas congeladas**.
- Alterações autorizadas nesta fase: este passo, os prompts proprietários dos testes e
  os testes de contrato novos. Implementação L1/L2 existente é somente leitura.
- Sucessor causal: uma implementação candidata só pode começar depois de congelados os
  hashes dos artefatos desta fase.

## Problema congelado

O catálogo embutido cobre apenas um conjunto delimitado de capacidades. Hoje o fluxo
consegue projetar uma capacidade sem regra como se tivesse investigado o candidato e
produzir `ABSENT`, `NOT_VERIFIED` ou uma lacuna funcional atribuída ao repositório. Essa
projeção mistura duas perguntas diferentes:

1. o Bateia possui uma regra segura para avaliar a capacidade?;
2. o candidato possui ou não a capacidade?

O passo separa essas perguntas antes de qualquer coleta. Uma lacuna da primeira não é
evidência sobre a segunda.

## Política explícita de Unknown

`CATALOG_GAP` é a causa pública e estável de Unknown quando nenhuma regra declarada e
segura cobre o requisito. Sua projeção funcional é obrigatoriamente
`EvidenceInconclusive`; nunca `EvidenceAbsent`, `EvidencePartial`, `EvidenceProven` ou
`EvidenceConflict`.

- `CATALOG_GAP` não entra em `CandidateAssessment.Gaps`, pois esse campo descreve lacunas
  atribuíveis ao candidato.
- `CATALOG_GAP` não reduz cobertura do candidato e não autoriza `BUILD`, `DROP` nem outra
  decisão conclusiva.
- a investigação usa `Classification == "CATALOG_GAP"`, e não `NOT_VERIFIED`, para essa
  causa;
- silêncio do catálogo, termo opaco, correspondência ambígua e fragmento lexical
  inseguro permanecem Unknown;
- falha ou incompletude de coleta continua `INCONCLUSIVE` por sua causa própria e não é
  rebatizada como `CATALOG_GAP`;
- nenhum default converte Unknown em sucesso.

## Contrato público aditivo

### Domínio

Adicionar em `01_core/domain/evidence.go`, preservando compatibilidade de leitura:

```go
type EvidenceRuleBasis string

const (
    EvidenceRuleDeclared   EvidenceRuleBasis = "DECLARED"
    EvidenceRuleHypothesis EvidenceRuleBasis = "HYPOTHESIS"
)

type CatalogCoverageStatus string

const (
    CatalogCoverageCovered    CatalogCoverageStatus = "COVERED"
    CatalogCoverageGap        CatalogCoverageStatus = "CATALOG_GAP"
    CatalogCoverageHypothesis CatalogCoverageStatus = "HYPOTHESIS"
)

type CatalogCoverageEntry struct {
    RequirementKind string                `json:"requirement_kind"`
    Requirement     string                `json:"requirement"`
    Status          CatalogCoverageStatus `json:"status"`
    RuleID          string                `json:"rule_id,omitempty"`
    RuleBasis       EvidenceRuleBasis     `json:"rule_basis,omitempty"`
    ReasonCode      string                `json:"reason_code"`
}

type EvidenceCatalogCoverage struct {
    SchemaVersion string                 `json:"schema_version"`
    RulesVersion  string                 `json:"rules_version"`
    RulesHash     string                 `json:"rules_hash"`
    Entries       []CatalogCoverageEntry `json:"entries"`
}
```

`CapabilityRule` ganha `Basis EvidenceRuleBasis` com JSON `basis,omitempty`. Ausência de
`basis` em catálogo estático legado validado equivale a `DECLARED`; qualquer regra criada
dinamicamente deve escrever `HYPOTHESIS` explicitamente. Valores desconhecidos são
inválidos.

`CapabilityEvidence` ganha `ReasonCode string` com JSON `reason_code,omitempty`.
`AssessmentReport` e `IdeaRunResult` ganham
`CatalogCoverage EvidenceCatalogCoverage` com JSON `catalog_coverage`. Arrays públicos
novos serializam como `[]`, nunca `null`.

### Fronteira pura anterior à coleta

Expor:

```go
func AnalyzeEvidenceCatalogCoverage(
    idea domain.IdeaSpec,
    catalog domain.CapabilityRuleCatalog,
) (domain.EvidenceCatalogCoverage, error)
```

A função valida o catálogo, não realiza I/O, não usa relógio/aleatoriedade/estado global,
não modifica entradas e devolve entradas canônicas para `problem` e `capability`.

Uma regra declarada é segura somente quando um de seus `CapabilityTerms`, após a mesma
normalização pública usada pelo avaliador, é igual ao requisito normalizado completo.
Substring reversa, fragmento genérico e semelhança lexical não são regras seguras. Zero
ou mais de uma regra aplicável produz `CATALOG_GAP`; uma regra hipotética produz
`HYPOTHESIS`, não `COVERED`.

Ordem canônica: `requirement_kind`, requisito normalizado, requisito original, `rule_id`.
Reordenar capacidades, regras, aliases ou itens de bundle não muda JSON, hash ou
classificação.

`PlanEvidenceForIdea*` deve obter essa cobertura antes de produzir qualquer plano, e
`RunIdeaPipeline` deve anexá-la ao resultado antes de chamar `CollectEvidence`. O mesmo
objeto, versão e hash atravessa planejamento, assessment, replay e apresentação, sem ser
recalculado a partir do conteúdo do candidato.

### Avaliação e atribuição

- Requisito `COVERED`: mantém as regras atuais de autoridade; correspondência lexical
  continua limitada a `PARTIAL` ou `INCONCLUSIVE` e jamais promove `PROVEN`.
- Requisito `CATALOG_GAP`: gera `INCONCLUSIVE`, fontes vazias, confiança zero e
  `ReasonCode == "CATALOG_GAP"`; não examina conteúdo para fabricar regra e não gera
  ausência do candidato.
- Requisito `HYPOTHESIS`: pode preservar indícios localizados, mas nunca gera `PROVEN`,
  `ABSENT` ou componente reutilizável. O `ReasonCode` é `HYPOTHESIS_ONLY`. Sem indício,
  permanece `INCONCLUSIVE`.
- `CandidateAssessment.Gaps` contém somente capacidades cobertas por regra declarada para
  as quais a avaliação admissível atribuiu uma lacuna ao candidato. Itens
  `CATALOG_GAP`/`HYPOTHESIS` ficam fora desse campo.
- Sanitização, decisão, ranking e investigação preservam os códigos; nenhum consumidor
  os reduz a `NOT_VERIFIED` genérico nem os conta como falso negativo do candidato.

### Catálogo fornecido no fluxo integrado

`IdeaRunDependencies.EvidenceRules` continua a porta de domínio. Catálogo não vazio deve
ser validado e usado na cobertura, planejamento, assessment, decisão e replay; o catálogo
embutido só é escolhido quando nenhum catálogo foi fornecido.

O gateway `idea recommend` aceita `--rules <catalog.json>`. Ele lê o catálogo antes de
provedores, injeta exatamente esse valor tanto no `RunIdeaPipeline` quanto no
`RecommendationService`, registra versão/hash em `run.json`, inclui o hash causal em
`input_hashes`, invalida resume com catálogo divergente e grava
`evidence/catalog-coverage.json`. Catálogo inválido falha antes de I/O externo.

Adicionar uma construção pública do serviço integrado que receba o catálogo, sem expor
dependências internas:

```go
func NewRecommendationServiceWithEvidenceRules(
    semantic SemanticExpander,
    discovery domain.ObservableDiscoveryProvider,
    content domain.ContentProvider,
    catalog domain.CapabilityRuleCatalog,
) (*RecommendationService, error)
```

O construtor legado mantém o default embutido. Não sintetizar regras é uma implementação
válida. Se qualquer implementação optar por síntese, todas as regras sintetizadas devem
ser `HYPOTHESIS` e obedecer os limites acima.

## Escopo observável

Coberto por este contrato:

- valores, tags JSON, arrays não nulos e round-trip dos novos tipos;
- análise de cobertura antes de coleta, sem I/O e sem mutação;
- distinção observável entre cobertura, lacuna do catálogo e hipótese;
- não atribuição de `ABSENT`, `NOT_VERIFIED`, gap ou componente ao candidato quando a
  limitação pertence ao catálogo;
- uso de catálogo fornecido pelo pipeline, serviço integrado e CLI;
- versão/hash causal, artifact de cobertura e invalidação de resume;
- determinismo sob repetição e pelo menos vinte reordenações;
- preservação conservadora de construções opacas.

Fora de escopo:

- aumentar o catálogo embutido;
- escolher mecanismo de síntese, LLM, embedding ou serviço remoto;
- comprovar comportamento do candidato por correspondência lexical;
- mudar políticas de licença, manutenção, operação, screening ou descoberta;
- provar equivalência funcional geral além deste fragmento.

## Corpus contratual pré-candidato

### Positivos

1. alias exato de regra declarada customizada gera `COVERED` e conserva `rule_id`;
2. pipeline/replay e `idea recommend --rules` usam versão/hash customizados, sem voltar ao
   default;
3. artifact de cobertura e assessment referem o mesmo catálogo.

### Negativos

1. capacidade fora do catálogo em bundle completo continua `INCONCLUSIVE/CATALOG_GAP`,
   nunca `ABSENT`, não entra em `CandidateAssessment.Gaps` e não vira `NOT_VERIFIED`;
2. fragmento `graph` não casa por reversão com alias `dependency graph`;
3. catálogo inválido ou hash divergente falha antes de provedor/coleta;
4. regra hipotética com símbolo e teste lexicais nunca gera `PROVEN`, `ABSENT` ou parte
   reutilizável.

### Opacos

1. requisito vazio, apenas pontuação, identidade ambígua ou sem alias exato permanece
   `CATALOG_GAP`;
2. regra hipotética sem testemunho permanece Unknown;
3. bundle incompleto continua causa de coleta, sem promover nem mascarar como cobertura.

### Determinismo

Executar no mínimo vinte permutações determinísticas de capacidades, regras, aliases e
itens. O JSON de cobertura e assessment deve ser byte a byte idêntico; inputs e slices
aninhadas permanecem imutáveis.

## Mutantes que o contrato deve eliminar

- M01: catálogo default substitui silenciosamente catálogo fornecido;
- M02: ausência de regra vira `ABSENT`;
- M03: ausência de regra vira `NOT_VERIFIED` na investigação;
- M04: `CATALOG_GAP` entra em `CandidateAssessment.Gaps`;
- M05: substring reversa escolhe regra;
- M06: primeira regra em ordem vence ambiguidade;
- M07: hipótese é serializada como declarada;
- M08: hipótese lexical promove `PROVEN`;
- M09: hipótese sem fonte produz `ABSENT`;
- M10: artifact/hash usa default enquanto assessment usa custom;
- M11: resume ignora mudança do catálogo;
- M12: reordenação altera saída;
- M13: nil é serializado no lugar de array vazio;
- M14: cobertura é inferida depois da coleta a partir do conteúdo do candidato.

Mutation score exigido para selo: `14/14 = 1.0`. `Unknown` só satisfaz os casos opacos
explicitamente declarados.

## Gates reproduzíveis

Pré-candidato, os três gates focais devem falhar por símbolos/comportamentos ainda não
materializados; essa falha é esperada e não constitui selo:

```sh
go test ./01_core/domain -run '^TestStep0071' -count=1
go test ./01_core/usecase -run '^TestStep0071' -count=1
go test ./02_shell/cli -run '^TestStep0071' -count=1
```

Após implementação candidata, executar na ordem:

```sh
go test ./01_core/domain -run '^TestStep0071' -count=1
go test ./01_core/usecase -run '^TestStep0071' -count=1
go test ./02_shell/cli -run '^TestStep0071' -count=1
go test ./01_core/domain ./01_core/usecase ./02_shell/cli -count=1
go test ./... -count=1
go test -race ./01_core/domain ./01_core/usecase ./02_shell/cli
crystalline-lint .
go build -o /tmp/bateia-step0071 ./04_wiring/main.go
```

Repetir os gates focais e comparar hashes/JSON. Qualquer alteração deste passo, prompts,
testes, baseline ou catálogo congelado invalida o selo e reinicia a cadeia na primeira
fase afetada.

## Entradas congeladas

Entradas normativas principais e SHA-256:

| caminho | SHA-256 |
|---|---|
| `.agents/AGENTS.md` | `0c8154925d3587fae78ec4ce82753fda16018424cbe6276bcb760401774de528` |
| `00_nucleo/prompts/01_core/domain/evidence.md` | `bd784f9fde373992b44590201c51055010dd5ba0d243f2a09fe110f953670f88` |
| `00_nucleo/prompts/01_core/usecase/evidence.md` | `57554d6c0f9d6c17b136943dbb7c6922b3d85f6305a8f8e6256703e4a01c9320` |
| `00_nucleo/prompts/01_core/usecase/evidence_authority.md` | `11ab46be889a94ce591c60d79eb69de369383d92705f92775f6e76b2e2aa82b5` |
| `00_nucleo/prompts/01_core/usecase/idea_run.md` | `7fa87ac532a8138d142276f2c1182952494f66b3fb4a0a4d0b3cd07ba9bddfa9` |
| `00_nucleo/prompts/02_shell/cli/idea_recommend.md` | `560feb3a3a1ce5f25c5f699841aaa40bfc1b48d14a9f7ef3546ac16cdd785e66` |
| `01_core/domain/evidence.go` | `cf91e64a08a514475c4fe4ba483bf4ec9b16a52eea7a8f43b136e192262a4f0a` |
| `01_core/domain/idea_run.go` | `5f2d648bd2d6dd9d5fbcb0b8b734ac1233896fd10a5a73713c7741aa21e0ab8c` |
| `01_core/domain/investigation_plan.go` | `98a80f4c72289e56ea3af38b7a3302118fc3ccc7354586b458e636fd33ce8d10` |
| `01_core/usecase/evidence.go` | `d1d19d88d6b1262b17ee848f69930bda1aec51986b73491a3cfd3adbf2572b03` |
| `01_core/usecase/idea_run.go` | `4104fe6178eca98662ab05728b884f8aca643e8860235196f70834169dcd339c` |
| `01_core/usecase/investigation_plan.go` | `05a642cd705c7daaad9559e8d0d25b6c6c82bad6205d5133333d7020676a32fa` |
| `01_core/usecase/recommendation.go` | `ad6badcdbf33019ffb0c81274961e68524da933a832c7cbbc6415b5c35491e1b` |
| `02_shell/cli/idea_recommend.go` | `f138b57ff2e96e12e5352b892c01eb3aad3fe9c5f4af1db0f5d8af43397bb953` |

Também pertencem à allowlist congelada todos os ADRs e todos os prompts
`00_nucleo/prompts/01_core/{domain,usecase}/evidence*.md` lidos nesta fase; seus hashes
individuais constam do recibo da autoria. Nenhum arquivo candidato futuro pertence a essa
allowlist.

## Limites de atestação

Este artefato e seus testes são segregados por autoridade e foram produzidos antes da
implementação candidata. Ainda não há selo, mutation score executado, veredito de
implementação ou certificado final. O que pode ser alegado nesta fase é apenas
**contrato pré-candidato segregado; atestação funcional pendente**.
