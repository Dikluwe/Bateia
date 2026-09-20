# Contrato observável — Passo 0026

**Revisão**: `R1-candidate`  
**Regime**: protocolo Tekt completo  
**Veredito positivo reservado**: `MCP_DURABLE_RUNS_READY`  
**Atestação máxima neste ambiente**: `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`

## 1. Autoridade, entradas congeladas e escopo

Este contrato foi escrito por `/root/contract_0026`, antes de qualquer leitura ou
escrita de implementação candidata, testes, ataques, recibos, selos ou laudo do
Passo 0026. A autoridade deste papel termina na obrigação observável e nos prompts
L0 de produção listados na seção 17.

Entradas permitidas e congeladas:

| Entrada | Identidade |
|---|---|
| Spec do Passo 0026 | `00_nucleo/passos/0026-persistir-execucoes-mcp-e-consultar-por-run-id.md`, SHA-256 `781fd3b7deae630ab6321e5dd010d83d85a439637d6a8e74a3017e2d11bb2656` |
| Baseline Git | commit `21c48e3571807777013710bcb6675c95c4130d9e` |
| Manifesto final 0025 | SHA-256 `9f4a06917139098b2a7a63f36cb438b936de9bdbb95d049dd9277b7f1296bc5f` |
| Certificado 0025 | SHA-256 `bb678dcd4e198347896e98683311148a35613eac99bedefee73fdff321ee5c4d`, veredito `MCP_RECOMMENDATION_GATEWAY_READY` |
| Recibo final 0025 | SHA-256 `a49e2733cbf47f03a7f537b8ec4efd548333302d36ad877a4422344f95f4bf4a` |
| ADR estrutural | `00_nucleo/ADR/ADR-0002-estrutura-de-pastas-e-camadas-tekt.md` |

O fragmento coberto é um servidor MCP local por `stdio`, um único processo escritor
por raiz persistente, fila local limitada, reinício honesto e leitura paginada das
seis projeções seguras definidas aqui. Não cobre resume, replay, remoção, expiração,
multiusuário, transporte remoto ou conteúdo bruto.

## 2. Compatibilidade obrigatória com o Passo 0025

`tools/list` passa de três para sete ferramentas, nesta ordem e sem outras:

1. `bateia_search`;
2. `bateia_benchmark`;
3. `bateia_recommend`;
4. `bateia_start_recommendation`;
5. `bateia_get_run`;
6. `bateia_get_artifact`;
7. `bateia_cancel_run`.

As três primeiras preservam byte por byte seus `inputSchema`, regras de validação,
envelope de sucesso e erros públicos certificados pelo Passo 0025. Em particular,
`bateia_recommend` continua síncrono, chama `RecommendationService.Recommend`, e
não passa a persistir, enfileirar ou devolver `run_id` implicitamente.

As mudanças inevitáveis nos hashes de implementação de `server.go`, `app.go`,
`recommendation.go` ou `main.go` exigem no certificado 0026 um recibo explícito de
**refinamento compatível** contra os hashes do certificado 0025. Nenhum prompt,
contrato, oráculo ou selo protegido de 0025 pode ser editado. A extensão observável
de progresso descrita abaixo deve manter idêntica a projeção de
`RecommendationService.Recommend` para a mesma entrada e os mesmos provedores.

## 3. Tipos de domínio congelados

Os nomes Go podem receber comentários, mas seus valores, campos JSON e semântica
não podem variar sem nova revisão do contrato.

```go
type RunID string // formato externo exato: ^run_[0-9a-f]{32}$

type RecommendationLifecycle string
const (
    RunQueued       RecommendationLifecycle = "QUEUED"
    RunRunning      RecommendationLifecycle = "RUNNING"
    RunSucceeded    RecommendationLifecycle = "SUCCEEDED"
    RunFailed       RecommendationLifecycle = "FAILED"
    RunCancelling   RecommendationLifecycle = "CANCELLING"
    RunCancelled    RecommendationLifecycle = "CANCELLED"
    RunInterrupted  RecommendationLifecycle = "INTERRUPTED"
)

type RecommendationRunStage string
const (
    RunStageQueued       RecommendationRunStage = "queued"
    RunStageExpand       RecommendationRunStage = "expand"
    RunStageDiscover     RecommendationRunStage = "discover"
    RunStageScreen       RecommendationRunStage = "screen"
    RunStageCollect      RecommendationRunStage = "collect"
    RunStageAssess       RecommendationRunStage = "assess"
    RunStageDecide       RecommendationRunStage = "decide"
    RunStageAlternatives RecommendationRunStage = "alternatives"
    RunStagePublish      RecommendationRunStage = "publish"
    RunStageTerminal     RecommendationRunStage = "terminal"
)

type RecommendationRunProgress struct {
    Planned   int `json:"planned"`
    Completed int `json:"completed"`
}

type RunArtifactName string
const (
    ArtifactRequest        RunArtifactName = "request"
    ArtifactRecommendation RunArtifactName = "recommendation"
    ArtifactDecision       RunArtifactName = "decision"
    ArtifactAlternatives   RunArtifactName = "alternatives"
    ArtifactFailures       RunArtifactName = "failures"
    ArtifactProvenance     RunArtifactName = "provenance"
)

type RunArtifactState string
const (
    ArtifactPending   RunArtifactState = "PENDING"
    ArtifactPublished RunArtifactState = "PUBLISHED"
)

type RunArtifactIndex struct {
    Name          RunArtifactName  `json:"name"`
    SchemaVersion string           `json:"schema_version"`
    Hash          string           `json:"artifact_hash"`
    SizeBytes     int64            `json:"size_bytes"`
    EntryCount    int              `json:"entry_count"`
    Version       uint64           `json:"version"`
    State         RunArtifactState `json:"state"`
}

type RecommendationRunFailure struct {
    Code    string                 `json:"code"`
    Stage   RecommendationRunStage `json:"stage"`
    Message string                 `json:"message"`
}

type RecommendationRun struct {
    SchemaVersion        string                    `json:"schema_version"`
    RunID                RunID                     `json:"run_id"`
    Request              RecommendationRequest     `json:"request"`
    RequestHash          string                    `json:"request_hash"`
    Lifecycle            RecommendationLifecycle   `json:"lifecycle"`
    Version              uint64                    `json:"version"`
    Stage                RecommendationRunStage    `json:"stage"`
    Progress             RecommendationRunProgress `json:"progress"`
    RecommendationStatus *IdeaRunStatus            `json:"recommendation_status"`
    Artifacts            []RunArtifactIndex        `json:"artifacts"`
    Failure              *RecommendationRunFailure `json:"failure"`
}

type IdempotencyRecord struct {
    KeyDigest   string `json:"key_digest"`
    RequestHash string `json:"request_hash"`
    RunID       RunID  `json:"run_id"`
}

type RecommendationStartReceipt struct {
    SchemaVersion string                  `json:"schema_version"`
    RunID         RunID                   `json:"run_id"`
    Lifecycle     RecommendationLifecycle `json:"lifecycle"`
    Version       uint64                  `json:"version"`
}

type RecommendationRunView struct {
    SchemaVersion        string                    `json:"schema_version"`
    RunID                RunID                     `json:"run_id"`
    Lifecycle            RecommendationLifecycle   `json:"lifecycle"`
    Version              uint64                    `json:"version"`
    Stage                RecommendationRunStage    `json:"stage"`
    Progress             RecommendationRunProgress `json:"progress"`
    RecommendationStatus *IdeaRunStatus            `json:"recommendation_status"`
    Artifacts            []RunArtifactIndex        `json:"artifacts"`
    Failure              *RecommendationRunFailure `json:"failure"`
}

type ArtifactReadRequest struct {
    RunID   RunID
    Artifact RunArtifactName
    Cursor  string
    Limit   int
}

type ArtifactPage struct {
    SchemaVersion string          `json:"schema_version"`
    RunID         RunID           `json:"run_id"`
    Artifact      RunArtifactName `json:"artifact"`
    ArtifactHash  string          `json:"artifact_hash"`
    Entries       [][]byte        `json:"-"`
    NextCursor    *string         `json:"next_cursor"`
}
```

`RecommendationRun` é o registro persistido, não a projeção MCP. A projeção
`RecommendationRunView` nunca inclui `Request`, `RequestHash`, digest de
idempotência nem artefatos `PENDING`. Arrays públicos são `[]`, nunca `null`.
Ponteiros JSON acima significam `null` explícito quando ausentes; os campos não
são omitidos.

`RunID` é gerado com pelo menos 128 bits de entropia por uma porta injetada, usa
32 dígitos hexadecimais minúsculos depois de `run_`, e nunca deriva de chave,
ideia, projeto, hash de entrada ou caminho. Entrada com outro formato é rejeitada
antes de qualquer acesso ao store.

## 4. Portas L1 congeladas

L1 pode coordenar portas, mas permanece sem filesystem, goroutine, mutex, timer,
fila concreta, variável de ambiente ou conhecimento de MCP.

```go
type RecommendationRunStore interface {
    ResolveIdempotency(context.Context, string) (IdempotencyRecord, bool, error)
    CreateOrResolve(context.Context, RecommendationRun, IdempotencyRecord) (RecommendationRun, bool, error)
    Get(context.Context, RunID) (RecommendationRun, error)
    Transition(context.Context, RunTransition) (RecommendationRun, error)
    PublishArtifact(context.Context, ArtifactPublication) (RecommendationRun, error)
    ReadArtifact(context.Context, ArtifactReadRequest) (ArtifactPage, error)
    InterruptNonTerminal(context.Context) ([]RecommendationRun, error)
}

type RecommendationQueueReservation interface {
    Commit(RunID) error // após Reserve bem-sucedido, Commit não pode falhar por capacidade
    Release()
}

type RecommendationJobQueue interface {
    Reserve(context.Context) (RecommendationQueueReservation, error)
    Cancel(context.Context, RunID) error
    Shutdown(context.Context) error
}

type RecommendationRunIDGenerator interface {
    NewRunID() (RunID, error)
}

type ObservableRecommendationExecutor interface {
    Recommend(context.Context, RecommendationRequest) (RecommendationResult, error)
    RecommendObserved(context.Context, RecommendationRequest, RecommendationProgressReporter) (RecommendationResult, error)
}

type RecommendationProgressReporter interface {
    ReportRecommendationProgress(RecommendationProgressEvent)
}
```

`RunTransition` contém `RunID`, `ExpectedVersion`, destino, estágio, progresso,
status de recomendação e falha pública. `ArtifactPublication` contém `RunID`,
`ExpectedVersion`, nome, schema `1` e o documento canônico. Todo método mutante é
CAS pela versão esperada. Conflito de CAS jamais é convertido em sucesso.

O caso de uso público expõe operações equivalentes a:

```go
Start(context.Context, string, RecommendationRequest) (RecommendationStartReceipt, error)
GetRun(context.Context, RunID) (RecommendationRunView, error)
GetArtifact(context.Context, ArtifactReadRequest) (ArtifactPage, error)
Cancel(context.Context, RunID) (RecommendationRunView, error)
Execute(context.Context, RunID) error
Recover(context.Context) error
```

L2 depende somente de uma interface com as quatro primeiras operações. O worker
concreto chama `Execute`; L4 chama `Recover` antes de aceitar novos jobs e antes
de iniciar workers.

## 5. Validação, normalização e hash idempotente

`idempotency_key` é string obrigatória de 1 a 128 bytes ASCII imprimíveis
`U+0020..U+007E`, exceto `/` e `\`, e não pode conter a substring `..`. Não é
trimada: espaços internos ou nas pontas participam da identidade. Chave inválida
não consulta store, não reserva fila e não gera ID.

A entrada de recomendação passa primeiro por
`ValidateRecommendationRequest`. O hash exclui `idempotency_key` e é:

```text
request_hash = "sha256:" + hex(SHA-256(json.Marshal(normalized_request)))
key_digest   = "sha256:" + hex(SHA-256(raw_idempotency_key_bytes))
```

Antes do marshal, os campos semanticamente conjuntos (`users`, `capabilities`,
`constraints`, `non_goals`, `preferred_stack`, `allowed_licenses`, `locales`) são
normalizados como no Passo 0025, deduplicados e ordenados por byte UTF-8. A ordem
de ranking nunca entra nesses campos. Objetos usam structs de campos fixos;
`json.Marshal` compacto de Go, sem indentação nem newline, é a codificação
canônica desta revisão. Repetição e reordenação equivalentes produzem o mesmo
hash; mudança semântica em qualquer campo produz outro hash.

Algoritmo obrigatório de `Start`, observável sob concorrência:

1. validar chave e pedido, normalizar e computar hashes sem portas;
2. resolver o digest existente; mesmo hash retorna o run existente, hash distinto
   retorna conflito;
3. se ausente, reservar capacidade da fila;
4. gerar ID e chamar `CreateOrResolve` atômico para digest + run;
5. se outro concorrente venceu, liberar a reserva e aplicar a regra do registro
   vencedor;
6. se criou, confirmar a reserva com o `run_id` e devolver o recibo original
   `QUEUED`, versão `1`, sem esperar o worker.

Fila cheia ocorre no passo 3 e não cria run nem associação. Validação falha antes
do passo 2. `CreateOrResolve` é a autoridade contra TOCTOU e cria exatamente um
run. Uma reserva garante espaço; implementação em que `Commit` volta a falhar por
capacidade viola o contrato.

Repetir a mesma chave e hash retorna o mesmo `run_id` e a lifecycle/version
corrente. Uma repetição após `INTERRUPTED` não retoma nem cria; nova execução
exige chave diferente.

## 6. Máquina de estados e versões

Transições permitidas, e somente estas:

| Origem | Destino | Autoridade |
|---|---|---|
| `QUEUED` | `RUNNING` | worker que obteve o job |
| `QUEUED` | `CANCELLED` | cancelamento que venceu a retirada da fila |
| `QUEUED` | `INTERRUPTED` | recuperação na abertura |
| `RUNNING` | `RUNNING` | progresso monotônico |
| `RUNNING` | `CANCELLING` | cancelamento |
| `RUNNING` | `SUCCEEDED` | worker após publicar todos os artefatos exigidos |
| `RUNNING` | `FAILED` | worker após falha pública persistível |
| `RUNNING` | `INTERRUPTED` | recuperação na abertura |
| `CANCELLING` | `CANCELLED` | worker após observar cancelamento e parar de iniciar chamadas |
| `CANCELLING` | `INTERRUPTED` | recuperação na abertura |

Estados terminais são `SUCCEEDED`, `FAILED`, `CANCELLED`, `INTERRUPTED`. Não há
self-transition terminal. Cancelar `CANCELLING` ou terminal devolve o registro
corrente sem alterar versão. Cada mutação pública bem-sucedida incrementa a versão
exatamente em um. Todo CAS perdedor relê e respeita o vencedor; não altera índice,
artefato ou terminal.

O estado inicial é exatamente: lifecycle `QUEUED`, versão `1`, stage `queued`,
progress `{planned:8,completed:0}`, status `null`, artifacts `[]`, failure `null`.
Os oito marcos controlados são `expand`, `discover`, `screen`, `collect`,
`assess`, `decide`, `alternatives`, `publish`. `planned` fica em `8` e
`completed` só cresce de `0` a `8`. Ao entrar em `RUNNING`, stage é `expand`.
Cada estágio concluído incrementa `completed` e move ao seguinte. Terminal usa
stage `terminal`; cancelamento/falha precoce preserva o último `completed`.
Eventos fora da allowlist, regressivos ou com texto livre são rejeitados.

`SUCCEEDED` exige um `RecommendationResult` projetado e válido. Seu
`recommendation_status` é exatamente `COMPLETE`, `PARTIAL` ou `INCONCLUSIVE` do
resultado, nunca inferido da lifecycle. `FAILED`, `CANCELLED` e `INTERRUPTED`
mantêm status `null`. Um resultado válido `INCONCLUSIVE` termina
`SUCCEEDED + INCONCLUSIVE`.

## 7. Cancelamento e corrida terminal

- `QUEUED`: a retirada da fila e o CAS para `CANCELLED` formam uma decisão única;
  se o worker venceu e publicou `RUNNING`, o fluxo segue como cancelamento em
  execução. Nenhum provedor é chamado se o cancelamento venceu.
- `RUNNING`: CAS para `CANCELLING`, depois `queue.Cancel` sinaliza o contexto. A
  resposta pode mostrar `CANCELLING`; somente o worker publica `CANCELLED` após
  observar o sinal.
- depois que o contexto cancelado é observado, nenhuma nova chamada a semantic,
  discovery ou content pode começar;
- corrida conclusão/cancelamento parte de `RUNNING`: somente um CAS para
  `SUCCEEDED` ou `CANCELLING` vence. Se `CANCELLING` vence, resultado final não é
  publicado; se `SUCCEEDED` vence, cancelamento devolve esse terminal imutável;
- repetir cancelamento não incrementa versão nem duplica sinais.

## 8. Progresso e uso do serviço certificado

A via assíncrona usa o mesmo objeto `RecommendationService` do Passo 0025. É
permitido refiná-lo com `RecommendObserved`; `Recommend` deve delegar à mesma
orquestração com observador nulo e preservar resposta/erro anteriores.
`RecommendObserved` não cria outro pipeline nem outro ranking. O observador recebe
somente estágio e contadores enumerados pelo Bateia, nunca mensagem de provider.

O worker é limitado e não cria goroutine por requisição. Há entre `1` e `4`
workers e capacidade entre `1` e `32`, ambas fixadas na construção em L4. Valores
fora desses intervalos impedem inicialização. O limite simultâneo observado de
chamadas a `RecommendationService` nunca excede workers.

## 9. Artefatos canônicos e publicação

Um documento armazenado tem bytes exatos:

```json
{"schema_version":"1","artifact":"alternatives","entries":[]}
```

Isto é, struct de campos nessa ordem, `json.Marshal` compacto, sem newline. O hash
é `sha256:` mais 64 hex minúsculos sobre esses bytes completos. `size_bytes` é o
número desses bytes. `entries` nunca é `null`.

Conteúdo por nome:

| Nome | Entradas, em ordem canônica |
|---|---|
| `request` | exatamente um `RecommendationRequest` normalizado, sem chave operacional |
| `recommendation` | exatamente um `ProjectRecommendationResult` |
| `decision` | exatamente o `ReuseDecisionReport` já projetado do resultado |
| `alternatives` | exatamente o `SoftwareAlternativeReport` já projetado do resultado |
| `failures` | zero ou mais `ExecutionFailure` sanitizadas, na ordem canônica do Passo 0025 |
| `provenance` | zero ou mais entradas achatadas `{locale,query,repository,kind,path,symbol,line,rule_id,content_hash}`, ordenadas por essa tupla |

`repository` em provenance é `owner/repo@commit` e só existe com commit imutável;
`kind`, caminho, símbolo, linha, rule id e hash vêm de `EvidenceSource` projetada.
Não há snippet, descrição bruta, README, código, teste, licença integral ou bundle.
Locale e query vêm do snapshot controlado; qualquer mensagem livre é excluída.

Publicação L3 segue: índice `PENDING` confirmado atomicamente; blob temporário
escrito, sincronizado, relido e validado por tamanho/hash; rename atômico;
diretório sincronizado; CAS final do índice para `PUBLISHED`. Leitura MCP enxerga
somente `PUBLISHED`. Temporário, órfão, `PENDING` e hash divergente nunca retornam
bytes.

Matriz mínima antes do terminal:

- `SUCCEEDED`: os seis artefatos estão `PUBLISHED`;
- `FAILED`: `request` e `failures` estão `PUBLISHED`;
- `CANCELLED`: `request` está `PUBLISHED`;
- `INTERRUPTED`: preserva todo publicado e não exige artefato além de `request`.

Ausência/corrupção de item exigido em terminal é `INTEGRITY_ERROR`. Item ainda
ausente e não exigido, inclusive após cancelamento/interrupção, é
`ARTIFACT_NOT_READY`; nunca vira lista vazia inventada.

## 10. Paginação e cursor

`limit` pertence a `[1,50]`. Cursor inicial é `""`. Cursor seguinte tem no máximo
512 bytes ASCII e formato externo `base64url(payload).base64url(mac)`, sem padding.
O payload canônico contém schema `1`, `run_id`, artifact, `artifact_version`, hash
e offset seguinte. `mac` é HMAC-SHA-256 com segredo local persistente de 256 bits,
criado com permissão privada e reutilizado após reopen.

Na leitura com cursor, L3 verifica MAC em tempo constante, formato e limites;
confere run, nome, versão e hash contra pedido e índice `PUBLISHED`; então fatia a
lista imutável pelo offset. Alteração de um byte, cursor de outro run/nome, offset
fora do intervalo, versão/hash divergente ou segredo diferente resulta
`CURSOR_INVALID`, sem entradas. O cursor não se vincula ao `limit`; o cliente pode
mudar o tamanho da próxima página sem mudar a identidade do snapshot.

`next_cursor` é `null` exatamente na última página. Leituras repetidas com mesmos
run, artifact, cursor e limit produzem os mesmos bytes JSON.

## 11. Store local, atomicidade, reopen e quotas

O adaptador L3 recebe a raiz apenas na construção por L4. Nenhuma ferramenta ou
tipo público aceita caminho. Identificadores externos nunca são concatenados em
caminho: nomes internos são digests hex SHA-256 após validação. Componentes da
raiz e arquivos existentes são verificados por `Lstat`; symlink, arquivo não
regular, hardlink inesperado (`nlink != 1`) e escape da raiz são recusados.
Diretórios são privados (`0700`) e arquivos (`0600`) por default.

Registros e índice de idempotência usam arquivo temporário exclusivo, fsync,
rename na mesma pasta e fsync do diretório. Falha antes da confirmação não é
promovida a sucesso. Registro ilegível ou hash divergente é integridade, não
ausência. O registro nunca contém token, cookie, header, variável de ambiente,
mensagem crua de adaptador ou chave de idempotência em claro.

`RunStoreLimits` é injetado por L4 e valida:

```text
MaxArtifactBytes > 0
MaxRunBytes >= MaxArtifactBytes
MaxTotalBytes >= MaxRunBytes
```

Quota mede bytes lógicos confirmados de registros, índices e blobs; temporários
incompletos não são publicados e são ignorados após validação/reopen. Nenhuma
operação pode ultrapassar o limite confirmado. Quota cheia rejeita novo start com
`QUOTA_EXCEEDED`, sem run/chave, e nunca remove item antigo. Artefato que excede
limite produz falha pública persistível e terminal `FAILED` somente se store e CAS
continuam comprovadamente graváveis; falha do próprio store deixa o último estado
confirmado não terminal.

Ao abrir, antes de workers e requests, `InterruptNonTerminal` faz CAS de todos os
`QUEUED|RUNNING|CANCELLING` para `INTERRUPTED`, stage `terminal`, status `null`,
preservando progresso e publicados. Terminais permanecem byte-estáveis. Não há
reexecução automática. O mesmo terminal é consultável após fechar e reabrir store
e servidor.

Encerramento para de aceitar reservas, cancela contextos ativos, aguarda workers
limitadamente e tenta confirmar `CANCELLED`; o que não for confirmado será
`INTERRUPTED` no reopen. Nunca presume terminal.

## 12. Schemas MCP exatos das quatro ferramentas

Todos usam `additionalProperties:false` em cada objeto. Validação runtime repete
as restrições mesmo se o cliente ignorar JSON Schema e ocorre antes de I/O.

### `bateia_start_recommendation`

```json
{
  "type":"object",
  "additionalProperties":false,
  "properties":{
    "idempotency_key":{"type":"string","minLength":1,"maxLength":128,"pattern":"^[\\u0020-\\u002E\\u0030-\\u005B\\u005D-\\u007E]+$"},
    "idea":{"type":"object","additionalProperties":false,"properties":{"problem":{"type":"string","minLength":1},"users":{"type":"array","items":{"type":"string"}},"capabilities":{"type":"array","minItems":1,"items":{"type":"string","minLength":1}},"constraints":{"type":"array","items":{"type":"string"}},"non_goals":{"type":"array","items":{"type":"string"}},"preferred_stack":{"type":"array","items":{"type":"string"}}},"required":["problem","capabilities"]},
    "policy":{"type":"object","additionalProperties":false,"properties":{"schema_version":{"type":"string","const":"1"},"allowed_licenses":{"type":"array","items":{"type":"string","minLength":1}},"require_local":{"type":"boolean"},"allow_archived":{"type":"boolean"},"max_components":{"type":"integer","minimum":2},"extend_min_coverage":{"type":"number","exclusiveMinimum":0,"maximum":1},"integration_min_value":{"type":"number","exclusiveMinimum":0,"maximum":1}},"required":["schema_version","allowed_licenses","require_local","allow_archived","max_components","extend_min_coverage","integration_min_value"]},
    "target":{"type":"object","additionalProperties":false,"properties":{"name":{"type":"string"}}},
    "search":{"type":"object","additionalProperties":false,"properties":{"strategy":{"type":"string","enum":["literal","lexical"]},"language":{"type":"string"},"locales":{"type":"array","uniqueItems":true,"items":{"type":"string","enum":["pt-BR","en-US","es-ES","fr-FR","de-DE","ja-JP","ko-KR","ru-RU","ar-SA","hi-IN"]}},"max_per_query":{"type":"integer","minimum":1,"maximum":25},"max_candidates":{"type":"integer","minimum":1,"maximum":20}},"required":["strategy","locales","max_per_query","max_candidates"]}
  },
  "required":["idempotency_key","idea","policy","target","search"]
}
```

Além do pattern, runtime rejeita `..`. Os quatro objetos de recomendação são
deep-equal aos de `bateia_recommend`; só a propriedade e requisito raiz
`idempotency_key` são acrescentados.

### `bateia_get_run`

```json
{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"}},"required":["run_id"]}
```

### `bateia_get_artifact`

```json
{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"},"artifact":{"type":"string","enum":["request","recommendation","decision","alternatives","failures","provenance"]},"cursor":{"anyOf":[{"const":""},{"type":"string","maxLength":512,"pattern":"^[A-Za-z0-9_-]+\\.[A-Za-z0-9_-]+$"}]},"limit":{"type":"integer","minimum":1,"maximum":50}},"required":["run_id","artifact","cursor","limit"]}
```

### `bateia_cancel_run`

```json
{"type":"object","additionalProperties":false,"properties":{"run_id":{"type":"string","pattern":"^run_[0-9a-f]{32}$"}},"required":["run_id"]}
```

## 13. Respostas MCP exatas

Cada sucesso usa o envelope MCP já certificado:

```json
{"content":[{"type":"text","text":"<JSON compacto>"}]}
```

O JSON textual de `start` contém exatamente, nessa ordem lógica:

```json
{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"QUEUED","version":1}
```

O JSON de `get_run` e `cancel_run` é a mesma projeção e sempre contém todos os
campos abaixo; `recommendation_status` e `failure` são `null` quando ausentes:

```json
{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","lifecycle":"RUNNING","version":4,"stage":"collect","progress":{"planned":8,"completed":3},"recommendation_status":null,"artifacts":[],"failure":null}
```

Índices em `artifacts` expõem somente `PUBLISHED` e contêm exatamente
`name,schema_version,artifact_hash,size_bytes,entry_count,version,state`, em ordem
da allowlist. O JSON de artefato contém exatamente:

```json
{"schema_version":"1","run_id":"run_0123456789abcdef0123456789abcdef","artifact":"alternatives","artifact_hash":"sha256:<64-hex>","entries":[],"next_cursor":null}
```

`entries` injeta valores JSON já validados, não strings JSON e não base64.

## 14. Erros públicos exatos

Erros são JSON-RPC, nunca resultado vazio. `RPCError` não ganha detalhes físicos.

| Código JSON-RPC | `message` exata | Situação |
|---:|---|---|
| `-32602` | `Invalid params` | JSON, campo extra/ausente, chave, ID, nome, enum ou limite inválido |
| `-32020` | `RUN_NOT_FOUND` | ID válido sem registro |
| `-32021` | `ARTIFACT_NOT_READY` | artifact permitido ainda não publicado/não exigido |
| `-32022` | `IDEMPOTENCY_CONFLICT` | chave existente com hash diferente |
| `-32023` | `QUEUE_FULL` | nenhuma reserva disponível |
| `-32024` | `QUOTA_EXCEEDED` | quota impede criar run |
| `-32025` | `INTEGRITY_ERROR` | registro/hash/artefato terminal exigido corrompido ou ilegível |
| `-32026` | `CURSOR_INVALID` | cursor sintaticamente válido mas adulterado, cruzado, expirado ou fora do intervalo |
| `-32027` | `RUN_STORE_UNAVAILABLE` | I/O não confirmável sem detalhe seguro |
| `-32603` | `Internal error` | falha interna não classificada |

Nome de artifact fora da allowlist é `Invalid params`, não `NOT_READY`. Erro
físico, path, segredo, chave, erro externo ou canário nunca aparece em `message`.
As ferramentas 0025 mantêm seus códigos e mensagens anteriores.

Falha terminal persistida usa somente as categorias/mensagens controladas:

| `code` | `message` |
|---|---|
| `RECOMMENDATION_FAILED` | `recommendation execution failed` |
| `RESULT_INVALID` | `recommendation result was invalid` |
| `ARTIFACT_LIMIT_EXCEEDED` | `artifact publication exceeded configured quota` |

O campo `stage` é um estágio canônico. Texto de adaptador nunca é copiado.

## 15. Segurança, privacidade e política de Unknown

- tool algum aceita path, diretório, store URL, token, cookie ou header;
- logs não incluem pedido completo, chave/digest, cursor, blob, erro externo ou
  conteúdo de terceiro; permitem apenas run id válido, lifecycle, stage e código
  público;
- a projeção segura do Passo 0025 é a única fonte de recommendation, decision,
  alternatives, failures e evidência;
- canários injetados em token, erro, README, código, teste, licença ou bundle não
  aparecem em respostas, arquivos públicos projetados, logs ou erro MCP;
- processo desaparecido não terminal vira `INTERRUPTED`;
- cancelamento não confirmado fica `CANCELLING` e vira `INTERRUPTED` no reopen;
- falha de persistência conserva o último estado confirmado, sem inventar
  `FAILED`;
- corrupção não vira `RUN_NOT_FOUND`, artifact pendente não vira `[]`, e
  `SUCCEEDED` não vira `COMPLETE` por default;
- licença desconhecida mantém exatamente o warning não bloqueante do Passo 0025.

Construções deliberadamente opacas de filesystem/runtime podem resultar
`Unknown` somente nos oráculos marcados como opacos. Em execução normal, toda
incerteza acima é materializada como estado ou erro explícito.

## 16. Critérios verificáveis e gate discriminatório

O autor dos oráculos deve cobrir ao menos:

1. lista exata de sete tools e deep-equality das três schemas legadas;
2. schemas/runtime estritos das quatro novas tools;
3. start novo devolve `QUEUED/v1` enquanto provider falso permanece bloqueado;
4. inválido, fila cheia e quota cheia não criam run/chave nem chamam porta;
5. mesma chave/hash sequencial e concorrente cria um run e uma execução;
6. mesma chave/hash diferente é conflito imutável;
7. hash igual sob reordenação equivalente e diferente sob mudança semântica;
8. todos os arcos permitidos, rejeição dos demais, CAS e versão monotônica;
9. `SUCCEEDED+INCONCLUSIVE` distinto de `COMPLETE`;
10. cancelamento queued sem provider, running sem nova chamada após observação,
    idempotência e corrida com conclusão;
11. workers/concurrency/capacidade nos limites, sem goroutine ilimitada;
12. publicação em duas fases, hash/tamanho, imutabilidade terminal e matriz de
    artefatos;
13. paginação repetível, mudança de limit, cursor cruzado/adulterado e hash/version
    binding;
14. traversal, symlink, hardlink, ID excessivo e nome inválido antes de disco;
15. quotas sem overshoot, remoção implícita ou reserva residual;
16. crash em create, transição, blob e commit; reopen preserva terminais e
    interrompe todos os órfãos honestamente;
17. scan de canários em MCP, store e logs;
18. e2e offline start→get→paginate→cancel→reopen com o mesmo
    `RecommendationService` injetado em L4;
19. `go test -race` nos pacotes novos;
20. refinamento compatível das três tools e resposta síncrona 0025.

As mutações mínimas da spec são obrigatórias. Acrescentar pelo menos mutantes de:
digest em vez de reserva atômica; commit de fila falível; progress regressivo;
status omitido em JSON; cursor sem MAC; cursor preso ao limit; artifact hash só
dos entries; `PENDING` legível; terminal antes dos artefatos; quota calculada após
write; reopen depois de iniciar worker; chave clara em disco/log; e array público
`null`. Toda mutação válida deve ser `Violated`, casos positivos `Preserved`, e
somente opacos designados `Unknown`. `mutation_score` requerido é `1.0`; repetição
e reordenação do gate produzem o mesmo veredito.

Gates finais: hashes selados intactos; `crystalline-lint .`; duas execuções
independentes de `go test -count=1 ./...`; `go test -race` nos pacotes novos;
`go vet ./...`; build de `04_wiring/main.go`; crash/reopen em diretório temporário
sem rede; e2e MCP offline; auditoria de imports e de I/O; scan de canários; recibo
de refinamento 0025; `git diff --check`.

## 17. Ownership por camada e prompts causais

| Camada | Autoridade | Prompt novo |
|---|---|---|
| L1 domain | enums, registros, transições, artefatos e erros puros | `00_nucleo/prompts/01_core/domain/recommendation_run.md` |
| L1 usecase | validação, hash, idempotência, state machine, projeção e execução; zero I/O concreto | `00_nucleo/prompts/01_core/usecase/recommendation_run.md` |
| L3 store | filesystem, atomicidade, integridade, cursor e quotas | `00_nucleo/prompts/03_infra/runstore/file_store.md` |
| L3 worker | fila, reservas, workers, contextos e shutdown | `00_nucleo/prompts/03_infra/runworker/local_queue.md` |
| L2 MCP | schemas, envelopes e mapeamento de erros | `00_nucleo/prompts/02_shell/mcp/durable_runs.md` |
| L2 CLI | somente transporte da dependência pronta até o servidor MCP | `00_nucleo/prompts/02_shell/cli/durable_mcp.md` |
| L4 wiring | raiz, limites, segredo, store, recovery, fila, serviço e shutdown | `00_nucleo/prompts/04_wiring/durable_runs.md` |

Os prompts certificados de 0025 permanecem byte-idênticos. Os novos prompts são
um refinamento aditivo: podem governar novos arquivos e as alterações mínimas de
linhagem/composição nos arquivos existentes. L4 não contém transição, hash,
validação ou política; L2 não importa L3; L3 não importa L2/L4; L1 não faz I/O.

Hashes dos prompts produzidos por esta autoridade:

| Prompt | SHA-256 |
|---|---|
| `01_core/domain/recommendation_run.md` | `4b38241173be2ad8fac884463d05b328bb8d1a0f327b64340cb1c41cfebe6c12` |
| `01_core/usecase/recommendation_run.md` | `c216d025c134d660c75de6af56faa9d62e30ae7e51709d7ec9a3f66e1fed722f` |
| `03_infra/runstore/file_store.md` | `20527537a9a4a58c4fda10abd539bd978a7358cf738682fa8ba42e0cd5ead2c3` |
| `03_infra/runworker/local_queue.md` | `c596346b247b26c4807d93a2ec0193c6913364ecbf3429fb556c98c6deef1d66` |
| `02_shell/mcp/durable_runs.md` | `a5f3980ca333fa4839f9938a9a88bd6c64876c6882e618a375a8e112aa675126` |
| `02_shell/cli/durable_mcp.md` | `661a282ab8dc8a93f74f812bc2be744d08290e6b1cc68601d4aaa367017f8251` |
| `04_wiring/durable_runs.md` | `b91c4419b3be0d448e3cda3e6dd1daf65ec7a503c7f4d609eaf239d238c182d5` |

Hashes de prompts 0025 deliberadamente preservados: MCP server
`82fbdf503a125a50da615fb85534440200593bebfca43b86ec5e7629be839f1e`,
domain recommendation
`36900286dddaddc682a7b3805d77915ab5617d5cf24249a0d817d4edbdd12eb6` e
usecase recommendation
`b5115cde4c609e5ce263e399cba154f2010387d28ff7e4cd9b5df7c48044e7b6`.

## 18. Ambiguidades resolvidas e remanescentes

Resoluções normativas desta revisão:

- retry idempotente devolve lifecycle/version corrente, não força `QUEUED/v1`;
- cursor não se vincula a `limit`, mas se vincula a run/nome/versão/hash;
- `recommendation`, `decision` e `alternatives` são documentos de uma entrada;
  paginação material ocorre sobretudo em failures/provenance;
- artifact não exigido em terminal cancelado/interrompido continua
  `ARTIFACT_NOT_READY`; artifact exigido ausente é integridade;
- a chave não é persistida em claro; persiste apenas digest SHA-256;
- progresso possui oito marcos exatos e não transporta mensagem livre;
- recuperação acontece antes de iniciar workers, portanto todo não terminal
  herdado é órfão e vira `INTERRUPTED`.

Ambiguidade ambiental real, não convertida em sucesso: atomicidade/durabilidade de
`rename` e `fsync` em filesystem que não ofereça as garantias usuais deve ser caso
opaco `Unknown`, registrado com tipo de filesystem e runtime. O certificado não
generaliza além do filesystem exercitado.

## 19. Critério de conclusão

O verificador pode emitir `MCP_DURABLE_RUNS_READY` somente se o selo com score
`1.0`, os gates finais e o recibo de refinamento 0025 passarem; o cliente então
inicia sem bloquear, recupera identidade idempotente, observa lifecycle e
progresso monotônicos, cancela sem novas chamadas, reabre terminais e pagina bytes
seguros verificados. O veredito cobre apenas as versões e ambiente registrados e,
sem isolamento material, usa obrigatoriamente
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`.
