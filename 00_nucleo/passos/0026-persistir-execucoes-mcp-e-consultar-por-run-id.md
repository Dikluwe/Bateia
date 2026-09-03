# Passo 0026 — Persistir execuções MCP e consultar por `run_id`

**Estado**: `PLANEJADO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Depende de**: [`Passo 0025 — Expor recomendação completa como ferramenta MCP`](0025-expor-recomendacao-completa-como-ferramenta-mcp.md)
**Baseline obrigatório**: certificado `MCP_RECOMMENDATION_GATEWAY_READY` do Passo 0025
**Onde roda**: servidor MCP local por `stdio`, com fila e armazenamento injetados em L4
**Camadas previstas**: L0, L1, L2, L3 e L4
**Natureza**: execução assíncrona local, persistência crash-safe e leitura seletiva

## Problema

O Passo 0025 permite que uma IA externa execute o funil completo por
`bateia_recommend`, mas a chamada permanece síncrona. Enquanto descoberta e
coleta estão em andamento, o agente não recebe uma identidade persistente, não
consegue distinguir fila de execução, não pode consultar progresso e precisa
receber todo o resultado final numa única resposta.

Se o processo MCP terminar, o cliente também não possui um contrato para saber se
a execução concluiu, foi interrompida ou nunca começou. Repetir uma chamada após
timeout pode iniciar trabalho duplicado e consumir novamente o orçamento dos
provedores.

## Incerteza que este passo reduz

Não sabemos se o Bateia pode sustentar investigações longas para agentes externos
com identidade estável, estados monotônicos, idempotência e artefatos consultáveis,
sem mover I/O para L1, sem usar caminhos como API e sem confundir término
operacional com conclusão da recomendação.

## Objetivo

Adicionar uma via assíncrona e persistente ao MCP. O agente inicia uma
recomendação, recebe imediatamente um `run_id`, acompanha o ciclo de vida, lê
somente projeções seguras do dossiê e pode solicitar cancelamento.

O `bateia_recommend` síncrono do Passo 0025 permanece disponível e
retrocompatível. A via nova chama o mesmo `RecommendationService`; não cria um
segundo pipeline.

```text
IA externa
  → bateia_start_recommendation
  ← run_id + QUEUED

worker limitado
  → RecommendationService
  → artefatos canônicos em RunStore

IA externa
  → bateia_get_run
  → bateia_get_artifact
  → bateia_cancel_run
```

## Distinção obrigatória de estados

O ciclo operacional do job não reutiliza `IdeaRunStatus`:

```text
QUEUED → RUNNING → SUCCEEDED
   │        ├────→ FAILED
   │        ├────→ CANCELLING → CANCELLED
   │        └────→ INTERRUPTED
   └─────────────→ CANCELLED
```

Estados terminais: `SUCCEEDED`, `FAILED`, `CANCELLED` e `INTERRUPTED`.

`SUCCEEDED` significa apenas que o serviço produziu um
`RecommendationResult` válido. O campo separado `recommendation_status` preserva
`COMPLETE`, `PARTIAL` ou `INCONCLUSIVE`. Portanto:

- `SUCCEEDED + COMPLETE`: execução íntegra e decisão conclusiva;
- `SUCCEEDED + PARTIAL`: resultado válido com falhas preservadas;
- `SUCCEEDED + INCONCLUSIVE`: resultado válido sem evidência para concluir;
- `FAILED`: não foi possível produzir um contrato final válido;
- `CANCELLED`: cancelamento foi reconhecido e nenhuma nova chamada pode começar;
- `INTERRUPTED`: o processo terminou ou o estado ficou opaco antes de um terminal.

É proibido projetar `SUCCEEDED` como `COMPLETE` por default. Estados terminais e
seus artefatos publicados são imutáveis.

## Ferramentas MCP

`tools/list` passa a anunciar, em ordem, as três ferramentas certificadas no
Passo 0025 e quatro ferramentas novas:

1. `bateia_search`;
2. `bateia_benchmark`;
3. `bateia_recommend`;
4. `bateia_start_recommendation`;
5. `bateia_get_run`;
6. `bateia_get_artifact`;
7. `bateia_cancel_run`.

### `bateia_start_recommendation`

Recebe exatamente o contrato aceito por `bateia_recommend`, acrescido de:

```json
{
  "idempotency_key": "agent-task-2026-09-01-0001",
  "idea": {},
  "policy": {},
  "target": {},
  "search": {}
}
```

`idempotency_key` é obrigatório, possui de 1 a 128 caracteres ASCII imprimíveis
e não pode conter barra, contrabarra ou sequência `..`. Não é `run_id` nem caminho.

Resposta imediata:

```json
{
  "schema_version": "1",
  "run_id": "run_opaque_identifier",
  "lifecycle": "QUEUED",
  "version": 1
}
```

A entrada inteira é validada antes de criar o run, reservar fila ou chamar uma
porta externa. A resposta deve ser produzida sem esperar descoberta ou coleta.

### `bateia_get_run`

Entrada:

```json
{
  "run_id": "run_opaque_identifier"
}
```

Saída mínima:

```json
{
  "schema_version": "1",
  "run_id": "run_opaque_identifier",
  "lifecycle": "RUNNING",
  "version": 4,
  "stage": "collect",
  "progress": {
    "planned": 5,
    "completed": 2
  },
  "recommendation_status": null,
  "artifacts": [],
  "failure": null
}
```

`stage` aceita somente nomes canônicos versionados: `queued`, `expand`,
`discover`, `screen`, `collect`, `assess`, `decide`, `alternatives`, `publish` e
`terminal`. Contadores são monotônicos e nunca contêm mensagens livres vindas de
provedores.

Em estado terminal, a resposta inclui `recommendation_status`, índice de
artefatos publicados e somente falha pública sanitizada quando aplicável. Um run
inexistente produz `RUN_NOT_FOUND`; não é materializado como execução vazia.

### `bateia_get_artifact`

Entrada:

```json
{
  "run_id": "run_opaque_identifier",
  "artifact": "alternatives",
  "cursor": "",
  "limit": 20
}
```

`artifact` pertence à allowlist:

- `request`: entrada canônica sem campos operacionais internos;
- `recommendation`: projeção segura do `RecommendationResult`;
- `decision`: `ReuseDecisionReport` projetado;
- `alternatives`: `SoftwareAlternativeReport` projetado;
- `failures`: falhas públicas sanitizadas;
- `provenance`: locales, consultas, identidades `repo@commit` e referências de
  evidência, sem conteúdo bruto.

`limit` pertence a `[1,50]`. O cursor é opaco, vinculado a `run_id`, nome,
versão e hash do artefato. Cursor alterado, pertencente a outro run ou invalidado
por corrupção é rejeitado. A resposta contém:

```json
{
  "schema_version": "1",
  "run_id": "run_opaque_identifier",
  "artifact": "alternatives",
  "artifact_hash": "sha256:...",
  "entries": [],
  "next_cursor": null
}
```

Artefato ainda não publicado produz `ARTIFACT_NOT_READY`, não lista vazia.
Artefato desconhecido produz erro de validação. Conteúdo integral de README,
código, teste, licença ou bundle não pertence a esta superfície.

### `bateia_cancel_run`

Entrada:

```json
{
  "run_id": "run_opaque_identifier"
}
```

Regras:

- `QUEUED → CANCELLED` sem chamar provedores;
- `RUNNING → CANCELLING`, cancela o contexto e só então publica `CANCELLED`;
- após a observação do cancelamento, nenhuma nova chamada externa começa;
- cancelar `CANCELLING` é idempotente;
- cancelar estado terminal não o reescreve e retorna o terminal existente;
- se cancelamento e conclusão competirem, uma única transição terminal vence por
  compare-and-swap; o perdedor não altera artefatos.

## Idempotência

O servidor mantém associação durável entre `idempotency_key`, hash da entrada
canônica e `run_id`:

- mesma chave e mesmo hash retornam o mesmo `run_id`, independentemente do estado;
- mesma chave e hash diferente produzem `IDEMPOTENCY_CONFLICT`;
- duas chamadas concorrentes com a mesma chave criam exatamente um run;
- erro de validação ou fila cheia não reserva a chave;
- `run_id` é opaco, gerado pelo servidor e não derivado diretamente de texto do
  usuário, nome de projeto ou caminho;
- repetir `get` ou `cancel` não cria nova execução.

## Contratos de domínio e portas

L1 define dados puros para identidade, ciclo de vida, progresso, índice de
artefatos e transições. I/O permanece atrás de portas equivalentes a:

```go
type RecommendationRunStore interface {
    Create(context.Context, RecommendationRun) error
    ResolveIdempotency(context.Context, string) (IdempotencyRecord, error)
    Get(context.Context, RunID) (RecommendationRun, error)
    Transition(context.Context, RunTransition) (RecommendationRun, error)
    PublishArtifact(context.Context, RunID, RunArtifact) error
    ReadArtifact(context.Context, ArtifactReadRequest) (ArtifactPage, error)
}

type RecommendationJobQueue interface {
    Enqueue(context.Context, RunID) error
    Cancel(context.Context, RunID) error
}
```

Os nomes internos podem variar; as autoridades e observáveis não. O serviço de
aplicação não conhece diretório, goroutine, mutex, formato físico, token ou MCP.

Eventos de progresso são dados controlados pelo Bateia. A persistência não aceita
mensagem arbitrária de adaptador. L4 injeta store, fila, limites e
`RecommendationService` já certificado.

## Adaptador persistente L3

Materializar primeiro um store local em disco, com raiz configurada em L4:

- cada `run_id` é mapeado internamente sem interpolar entrada do usuário;
- criação e transições usam compare-and-swap por `version` monotônica;
- registros e índices são escritos atomicamente;
- artefatos são publicados somente depois de hash e tamanho confirmados;
- artefatos terminais nunca são sobrescritos;
- symlink, hardlink inesperado e traversal são rejeitados;
- permissões não tornam o store público por default;
- um registro nunca contém token, cookie, cabeçalho ou valor de ambiente;
- corrupção de hash é erro de integridade, não ausência nem sucesso;
- quota por run, quota total e tamanho máximo de artefato são configurados em L4.

Este passo suporta um único processo MCP escritor por raiz. Lock distribuído,
múltiplas instâncias e armazenamento remoto ficam fora do escopo.

## Fila e reinício

- número de workers pertence a `[1,4]` e é fixado em L4;
- capacidade da fila é finita e configurada, com máximo inicial de 32;
- fila cheia rejeita `start` com erro público e não cria run;
- workers não são criados por requisição sem limite;
- encerramento cancela contextos ativos e tenta publicar estado antes de sair;
- ao iniciar, o servidor preserva todos os terminais já gravados;
- qualquer `QUEUED`, `RUNNING` ou `CANCELLING` herdado sem worker comprovadamente
  ativo transita para `INTERRUPTED`;
- `INTERRUPTED` preserva progresso e artefatos parciais indexados, mas não inventa
  `recommendation_status` nem `top_choice`;
- retomada automática não pertence a este passo.

Essa política privilegia honestidade depois de crash. Um run interrompido continua
consultável; nova execução exige nova `idempotency_key` até existir contrato de
resume explícito.

## Publicação de artefatos

Cada artefato possui nome lógico, schema, SHA-256, tamanho e estado
`PENDING|PUBLISHED`. O índice é gravado antes como pendente; o blob validado é
publicado atomicamente; só depois o índice muda para publicado.

O MCP lê exclusivamente `PUBLISHED`. Bytes órfãos, temporários ou com hash
divergente não são retornados. A projeção segura do Passo 0025 continua sendo a
autoridade para decisão, alternativas, URLs, falhas e evidências. Persistir não
reintroduz snippets ou mensagens externas removidas por aquela projeção.

## Política de `Unknown`

- processo desaparecido durante estado não terminal produz `INTERRUPTED`;
- falta de artefato em run ativo produz `ARTIFACT_NOT_READY`;
- artefato esperado ausente em terminal produz erro de integridade;
- hash divergente ou registro ilegível não é convertido em `RUN_NOT_FOUND`;
- falha de store antes da criação não produz `run_id` fictício;
- falha depois da criação preserva o último estado confirmado e não promove
  `FAILED` quando a persistência do terminal não pôde ser comprovada;
- cancelamento sem confirmação permanece `CANCELLING` ou vira `INTERRUPTED` após
  reinício, nunca `CANCELLED` presumido;
- resultado `INCONCLUSIVE` continua sendo um resultado válido e pode coexistir
  com lifecycle `SUCCEEDED`;
- licença desconhecida mantém o aviso não bloqueante certificado no Passo 0025.

## Segurança e privacidade

- nenhum tool aceita caminho, diretório, URL de store, token ou cabeçalho;
- `run_id`, chave de idempotência e cursor passam por validação de tamanho e
  alfabeto antes de qualquer acesso ao disco;
- conteúdo de repositório continua sendo dado não confiável;
- erros físicos e mensagens de adaptadores são projetados para categorias
  públicas estáveis;
- artefatos MCP não incluem conteúdo bruto de terceiros;
- logs não incluem entrada completa, chave de idempotência ou conteúdo de
  artefatos;
- quotas impedem crescimento ilimitado e fila limita uso concorrente;
- este store é local e de usuário único; não promete isolamento entre usuários;
- não existe ferramenta destrutiva ou política automática de remoção neste passo.

Quando a quota total for atingida, novas execuções são rejeitadas. Expiração,
deleção e garbage collection exigirão contrato posterior e não podem apagar runs
silenciosamente.

## Observáveis a congelar antes da implementação

1. `tools/list` anuncia exatamente sete ferramentas na ordem especificada;
2. as três ferramentas do Passo 0025 preservam schema, resultado e erros públicos;
3. `start` válido responde `QUEUED` antes que um provedor falso bloqueado retorne;
4. entrada inválida e fila cheia não criam run, não reservam idempotência e não
   chamam portas externas;
5. mesma chave/mesmo hash, inclusive sob concorrência, produz exatamente um run;
6. mesma chave/hash diferente produz conflito e não altera o run original;
7. versões e estados seguem somente transições permitidas e nunca retrocedem;
8. `SUCCEEDED + INCONCLUSIVE` permanece distinto de `COMPLETE`;
9. `get_run` lê o mesmo terminal depois de reiniciar store e servidor;
10. não terminais órfãos viram `INTERRUPTED`, preservando progresso sem resultado
    inventado;
11. cancelamento em fila não chama provedor; cancelamento em execução não inicia
    nova chamada depois de observado;
12. corrida cancelamento/conclusão publica um único terminal imutável;
13. artefatos publicados possuem hash válido, paginação canônica e cursor preso ao
    run e à versão;
14. artefato pendente, temporário, órfão ou corrompido nunca é devolvido;
15. traversal, symlink, identificador excessivo e cursor adulterado são rejeitados;
16. canários em token, erro externo, README e código não aparecem em status,
    artefatos seguros, logs ou erros MCP;
17. fila e workers nunca excedem os limites configurados;
18. reordenar entradas equivalentes preserva hash canônico e idempotência.

## Mutações mínimas a rejeitar

- `run_id` usa diretamente `idempotency_key` ou texto do usuário;
- mesma chave cria dois runs concorrentes;
- chave igual com pedido diferente reaproveita o run;
- lifecycle `SUCCEEDED` força `recommendation_status=COMPLETE`;
- estado terminal volta para `RUNNING` ou artefato terminal é sobrescrito;
- resposta `QUEUED` espera o provedor concluir;
- cada `start` cria goroutine sem fila ou limite;
- cancelamento publica `CANCELLED` antes de o worker parar de iniciar chamadas;
- restart trata `RUNNING` como `FAILED`, `SUCCEEDED` ou `CANCELLED` presumido;
- `ARTIFACT_NOT_READY` vira lista vazia;
- cursor de outro run consegue ler artefato;
- hash divergente é ignorado;
- store concatena `run_id` ou artifact name num caminho sem validação;
- arquivo por symlink escapa da raiz;
- resposta reintroduz conteúdo bruto ou canário externo;
- L1 importa filesystem, banco, fila concreta, shell ou infraestrutura;
- L2 instancia store ou worker concreto;
- ferramenta antiga muda para satisfazer a via assíncrona;
- quota cheia remove automaticamente o run mais antigo;
- run interrompido inventa resultado ou `top_choice`.

Mutation score exigido para selar o contrato: `1,0`.

## Gates de aceitação

- manifesto, contrato, oráculos positivos/opacos e mutações congelados antes do
  patch candidato;
- `crystalline-lint .` sem violações de camada;
- `go test -count=1 ./...` passa duas vezes;
- `go build -o bateia ./04_wiring/main.go` passa;
- `go vet ./...` passa;
- teste de crash/reabertura usa diretório temporário e nenhuma rede real;
- teste de concorrência passa também com `go test -race` nos pacotes novos;
- teste MCP ponta a ponta inicia, consulta, pagina, cancela e reabre runs com
  provedores falsos;
- mutation score `1,0` para transições, idempotência, integridade, cancelamento,
  quotas e fronteiras de segurança;
- auditoria confirma L1 sem I/O, L2 sem L3 e L4 sem regra de negócio;
- scan de canários não encontra credenciais, erro externo ou conteúdo bruto;
- os hashes e observáveis certificados do Passo 0025 permanecem preservados ou
  recebem recibo explícito de refinamento compatível.

## Regime Tekt previsto

Aplicar protocolo completo. A spec deste passo, o certificado do Passo 0025 e a
política de `Unknown` formam as entradas congeladas. Autor do contrato não lê o
patch candidato; autor de oráculos e adversário não adaptam casos ao store
escolhido; implementador não altera contrato, mutações ou testes protegidos;
verificador não corrige produção nem artefatos julgados.

O adversário recebe capacidade para interromper escrita, reabrir store, competir
transições, adulterar cursor, trocar hashes, criar symlinks e injetar canários,
mas não para corrigir o candidato. O verificador registra versões do runtime,
filesystem e comandos de repetição.

Sem isolamento material, o resultado deve ser declarado
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`, mesmo que todos os gates passem.

## Fora de escopo

- retomar automaticamente runs interrompidos;
- modos MCP `replay`, `resume` ou análise de worktree local;
- leitura de bundles, código, README, testes ou licenças em conteúdo bruto;
- inspeção aprofundada de um candidato por nova coleta;
- notificações push, long polling ou recursos MCP além de tools;
- transporte MCP por rede, autenticação multiusuário ou autorização por run;
- múltiplos processos escrevendo na mesma raiz;
- store remoto, banco distribuído ou object storage;
- deleção, expiração e garbage collection;
- novos conectores de app stores, SaaS, registries ou mercados regionais;
- tradução automática ou incorporação de IA ao Bateia.

## Como pode falhar

- persistência vira um espelho do diretório privado da CLI;
- `run_id` ou cursor se tornam caminhos disfarçados;
- estado em memória diverge do estado confirmado no store;
- fila responde rápido, mas inicia trabalho sem limite;
- progresso livre ecoa mensagens externas ou segredos;
- cancelamento responde antes de impedir novas chamadas;
- crash produz terminal positivo sem testemunha persistida;
- paginação muda entre leituras do mesmo artefato;
- quotas são aplicadas depois de gravar parcialmente;
- abstrações de store e fila transportam detalhes de L3 para L1;
- a via assíncrona altera o resultado semântico certificado da via síncrona.

## Critério de conclusão

O passo termina quando um cliente MCP consegue iniciar uma recomendação sem
bloquear, recuperar o mesmo `run_id` por idempotência, observar transições
monotônicas, cancelar com segurança, consultar um terminal depois de reiniciar o
servidor e paginar artefatos seguros cujo hash foi verificado; tudo isso usando o
mesmo `RecommendationService` do Passo 0025 e sem expor caminhos, conteúdo bruto,
credenciais ou I/O em L1.

O veredito positivo é `MCP_DURABLE_RUNS_READY`. Ele atesta somente execução local
assíncrona em um processo escritor por store e não implica retomada automática,
multiusuário, transporte remoto, cobertura de novos mercados ou inspeção bruta de
repositórios.
