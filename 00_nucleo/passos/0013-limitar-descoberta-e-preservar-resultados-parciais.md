# Passo 0013 — Limitar descoberta e preservar resultados parciais

**Estado**: `EM VALIDAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0012 — Executar autofagia competitiva e formar um Katamari mínimo`](0012-executar-autofagia-competitiva-e-formar-katamari-minimo.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2 e L3  
**Natureza**: controle operacional, cancelamento estruturado e observabilidade da descoberta sem LLM

## Evidência que motiva este passo

A execução E-0012 tentou descobrir concorrentes do Bateia com:

- 16 expansões;
- dois fallbacks por expansão sem filtro de linguagem;
- concorrência 2;
- timeout HTTP de 45 segundos por requisição.

No pior caso, essa configuração pode consumir aproximadamente:

```text
16 expansões × 2 tentativas × 45s ÷ concorrência 2 = 12 minutos
```

O comando `idea discover` usa `context.Background()` e não oferece deadline
total. A execução autorizada permaneceu cerca de cinco minutos sem produzir
progresso, snapshot parcial ou erro e foi interrompida externamente. Todos os
resultados eventualmente concluídos foram perdidos porque a CLI somente grava o
snapshot depois que `DiscoverExpansionPlan` retorna com sucesso.

O primeiro run dentro do sandbox teve comportamento diferente: as dezesseis
expansões falharam imediatamente por bloqueio de rede. Essa diferença confirma
que precisamos distinguir:

```text
falha de ambiente
    ≠
timeout de tentativa
    ≠
deadline total
    ≠
cancelamento humano
    ≠
ausência de candidatos
```

## Incerteza que este passo reduz

Não sabemos se a descoberta competitiva falhou por baixa qualidade das queries,
indisponibilidade do provedor ou apenas porque o comando não terminou dentro do
tempo esperado. Sem receipts parciais, essas hipóteses permanecem misturadas.

Este passo deve tornar toda execução finita e preservar evidência suficiente para
atribuir a perda antes de alterar vocabulário, provedor ou regras de evidência.

## Intenção

Transformar descoberta remota numa operação limitada e observável:

```text
plano congelado
    → deadline total
    → tentativas com timeout individual
    → progresso por expansão
    → cancelamento estruturado
    → snapshot COMPLETE | PARTIAL | FAILED
```

Uma execução parcial é um artefato válido quando declara precisamente o que foi e
o que não foi tentado. Ela não é sucesso de descoberta nem deve ocultar perdas.

## Restrição central

Este passo não pode melhorar artificialmente recall. É proibido:

- adicionar nomes de concorrentes às queries;
- alterar o vocabulário lexical;
- trocar Sourcegraph por outro provedor;
- reduzir evidência exigida pela peneira;
- repetir silenciosamente queries expiradas;
- transformar timeout ou cancelamento em lista vazia bem-sucedida;
- contar candidato parcial como correspondente sem commit e origem;
- descartar attempts concluídos porque outra expansão falhou;
- usar LLM, embeddings ou geração probabilística.

O resultado esperado é melhor diagnóstico e finitude, não mais candidatos.

## Contratos

### Política de execução

```go
type DiscoveryExecutionPolicy struct {
    SchemaVersion  string        `json:"schema_version"`
    TotalTimeout   time.Duration `json:"total_timeout"`
    AttemptTimeout time.Duration `json:"attempt_timeout"`
    Concurrency    int           `json:"concurrency"`
    MaxAttempts    int           `json:"max_attempts"`
}
```

Defaults do primeiro incremento:

```text
TotalTimeout   = 3 minutos
AttemptTimeout = 20 segundos
Concurrency    = 2
MaxAttempts    = fallbacks já declarados pelo provedor, no máximo 3
```

`AttemptTimeout` nunca pode exceder `TotalTimeout`. Durações zero, negativas ou
acima dos limites de segurança são inválidas.

### Estado da execução

```go
type DiscoveryRunStatus string

const (
    DiscoveryRunComplete DiscoveryRunStatus = "COMPLETE"
    DiscoveryRunPartial  DiscoveryRunStatus = "PARTIAL"
    DiscoveryRunFailed   DiscoveryRunStatus = "FAILED"
)

type DiscoveryCompletionCause string

const (
    CompletionFinished       DiscoveryCompletionCause = "FINISHED"
    CompletionTotalTimeout   DiscoveryCompletionCause = "TOTAL_TIMEOUT"
    CompletionCancelled      DiscoveryCompletionCause = "CANCELLED"
    CompletionProviderFailed DiscoveryCompletionCause = "PROVIDER_FAILED"
    CompletionInvalidInput   DiscoveryCompletionCause = "INVALID_INPUT"
)
```

### Receipt por tentativa

```go
type ProviderAttemptReceipt struct {
    ExpansionID   string        `json:"expansion_id"`
    AttemptIndex  int           `json:"attempt_index"`
    PlannedQuery  string        `json:"planned_query"`
    EffectiveQuery string       `json:"effective_query"`
    StartedAt     time.Time     `json:"started_at"`
    Duration      time.Duration `json:"duration"`
    Status        string        `json:"status"`
    ResultCount   int           `json:"result_count"`
    ErrorClass    string        `json:"error_class,omitempty"`
    SafeMessage   string        `json:"safe_message,omitempty"`
}
```

Estados mínimos de attempt:

```text
SUCCEEDED · EMPTY · TIMED_OUT · CANCELLED · FAILED · NOT_STARTED
```

`SafeMessage` deve continuar removendo tokens e nunca serializar headers ou URLs
assinadas.

### Snapshot finito

Estender o snapshot sem quebrar replays anteriores:

```go
type DiscoveryRunReceipt struct {
    SchemaVersion   string                   `json:"schema_version"`
    Status          DiscoveryRunStatus       `json:"status"`
    CompletionCause DiscoveryCompletionCause `json:"completion_cause"`
    Policy          DiscoveryExecutionPolicy `json:"policy"`
    Planned         int                      `json:"planned"`
    Started         int                      `json:"started"`
    Completed       int                      `json:"completed"`
    TimedOut        int                      `json:"timed_out"`
    Cancelled       int                      `json:"cancelled"`
    Failed          int                      `json:"failed"`
    Attempts        []ProviderAttemptReceipt `json:"attempts"`
}
```

O `DiscoverySnapshot` recebe um receipt opcional. Snapshots legados sem receipt
continuam legíveis, mas são `UNKNOWN` quando um benchmark exige prova operacional.

## Semântica de resultado

### `COMPLETE`

Todas as expansões planejadas terminaram em sucesso, vazio ou falha terminal
registrada antes do deadline. Um run pode ser `COMPLETE` e conter falhas; complete
significa que o plano acabou, não que encontrou candidatos.

### `PARTIAL`

Pelo menos uma expansão produziu attempt ou candidato preservado, mas o deadline
total ou cancelamento impediu concluir o plano.

### `FAILED`

Nenhuma expansão produziu receipt consumível ou todos os attempts falharam antes
de qualquer resultado. Mesmo assim, o receipt deve ser escrito.

Lista vazia só representa ausência observada quando todas as expansões foram
concluídas e seus attempts terminaram `EMPTY`. Timeout nunca equivale a vazio.

## Cancelamento e concorrência

Cada expansão recebe o contexto do run. Cada attempt recebe um contexto filho
com `AttemptTimeout`, limitado pelo deadline restante do pai.

Ao cancelar:

1. impedir início de novas expansões;
2. propagar cancelamento para requests ativos;
3. aguardar somente a finalização estruturada dos workers;
4. registrar expansões não iniciadas como `NOT_STARTED`;
5. ordenar receipts por expansão e attempt;
6. escrever snapshot parcial atomicamente;
7. retornar erro tipado depois da escrita, quando apropriado.

Nenhuma goroutine pode permanecer ativa após o retorno.

## Progresso

A porta L1 deve aceitar um observador opcional:

```go
type DiscoveryProgressEvent struct {
    ExpansionID  string `json:"expansion_id"`
    AttemptIndex int    `json:"attempt_index"`
    Status       string `json:"status"`
    Completed    int    `json:"completed"`
    Planned      int    `json:"planned"`
}

type DiscoveryProgressReporter interface {
    ReportDiscoveryProgress(DiscoveryProgressEvent)
}
```

A CLI escreve progresso conciso em `stderr`; JSON permanece em arquivo ou
`stdout`, sem mistura. O reporter não participa da decisão nem pode modificar o
run.

## CLI

```bash
bateia idea discover \
  --input ideia.json \
  --strategy lexical \
  --max-per-query 3 \
  --concurrency 2 \
  --timeout 3m \
  --attempt-timeout 20s \
  --output discovery.json
```

Regras:

- `--output` continua obrigatório antes de acessar a rede;
- snapshot é escrito para `COMPLETE`, `PARTIAL` e `FAILED`;
- `PARTIAL` pode retornar exit code não zero, mas nunca perder o arquivo;
- erro de validação anterior à rede não cria snapshot enganoso;
- Ctrl-C usa o mesmo caminho estruturado de cancelamento;
- política efetiva aparece no receipt.

## Preflight semântico

Antes da rede, registrar:

- quantidade de expansões literais;
- propostas lexicais;
- duplicatas removidas;
- expansões finais;
- capacidades sem expansão alternativa.

`ProposedCount=0` não impede a execução, mas gera aviso explícito
`NO_SEMANTIC_EXPANSION`. Esse sinal permite atribuir recall baixo ao vocabulário
sem misturá-lo com falha do provedor.

Este passo apenas torna a perda visível; não adiciona termos ao vocabulário.

## Métricas operacionais

```go
type DiscoveryOperationalMetrics struct {
    CompletionRate       float64  `json:"completion_rate"`
    AttemptTimeoutRate   float64  `json:"attempt_timeout_rate"`
    ProviderFailureRate  float64  `json:"provider_failure_rate"`
    EmptyAttemptRate     float64  `json:"empty_attempt_rate"`
    CandidateYield       *float64 `json:"candidate_yield"`
    PartialPreserveRate  *float64 `json:"partial_preserve_rate"`
}
```

Denominadores vazios permanecem `null`. Essas métricas não medem relevância.

## Experimento de correção

Reexecutar exatamente a Fase A do E-0012:

- mesma IdeaSpec;
- mesma estratégia lexical;
- mesmas 16 expansões;
- `max-per-query=3`;
- concorrência 2;
- Sourcegraph;
- sem abrir ou injetar o registry competitivo.

O sucesso deste passo não exige encontrar Replen ou qualquer concorrente. Exige:

1. terminar em até 3 minutos mais tolerância de finalização de 5 segundos;
2. escrever snapshot e receipt;
3. preservar toda tentativa concluída;
4. classificar corretamente timeout, falha, vazio e cancelamento;
5. registrar `NO_SEMANTIC_EXPANSION` se o plano continuar literal.

O recall competitivo só será reavaliado depois de corrigirmos a ontologia lexical
num passo separado.

## Gates

- `100%` dos runs terminam dentro do deadline mais 5 segundos;
- `100%` dos attempts ativos respeitam o menor deadline aplicável;
- `100%` dos resultados concluídos antes do cancelamento aparecem no snapshot;
- `0` goroutines ou requests observáveis após o retorno;
- `0` timeout contado como vazio;
- `0` segredo serializado;
- ordem de conclusão não altera JSON canônico;
- Ctrl-C e timeout total produzem receipts semanticamente equivalentes;
- falha de todos os providers ainda escreve receipt `FAILED`;
- replay do snapshot parcial não toca rede;
- linter V15 permanece sem violações e todo consumidor novo possui prompt próprio.

## Testes obrigatórios

1. run rápido completo produz `COMPLETE/FINISHED`;
2. attempt lento respeita `AttemptTimeout`;
3. deadline total cancela attempts ativos;
4. expansões ainda não iniciadas viram `NOT_STARTED`;
5. candidato concluído antes do timeout permanece no snapshot;
6. run com resultado e timeout produz `PARTIAL`;
7. todos os attempts falhos produzem `FAILED/PROVIDER_FAILED`;
8. todos os attempts vazios produzem `COMPLETE`, não `FAILED`;
9. timeout não incrementa contador de vazio;
10. cancelamento humano produz `CANCELLED`;
11. nenhum worker reporta progresso depois do retorno;
12. reporter nulo não causa panic;
13. reporter lento não bloqueia deadline;
14. receipts são ordenados deterministicamente;
15. mensagens removem token do provedor;
16. output parcial é escrito atomicamente antes do erro da CLI;
17. output obrigatório é validado antes da rede;
18. timeout inválido é rejeitado antes da rede;
19. attempt timeout maior que total é rejeitado;
20. snapshot legado continua desserializável;
21. snapshot legado é opaco para gate operacional estrito;
22. `ProposedCount=0` gera `NO_SEMANTIC_EXPANSION`;
23. repetição com as mesmas respostas produz o mesmo conteúdo canônico;
24. suíte completa e linter atualizado permanecem verdes.

## Critério de produto

Uma pessoa deve conseguir responder, usando apenas o snapshot:

- quanto do plano foi executado;
- quais queries chegaram ao provedor;
- quais terminaram vazias, falharam ou expiraram;
- que candidatos foram preservados;
- por que o run terminou;
- se o problema seguinte está no provedor, deadline, query ou vocabulário.

O passo não é aprovado apenas porque “agora termina”. Ele precisa preservar os
resultados úteis e impedir conclusões falsas sobre ausência.

## Critério de parada

Interromper se:

- for necessário alterar queries para fazer o teste passar;
- o snapshot parcial não puder ser reproduzido offline;
- algum worker sobreviver ao retorno;
- cancelamento perder candidato já concluído;
- erro de rede for reclassificado como vazio;
- o output for escrito antes de validar caminho e política;
- receipts incluírem token, header ou conteúdo sensível;
- os defaults exigirem mais que três minutos no E-0012.

## Fora de escopo

- ampliar o vocabulário lexical;
- criar regras de evidência para capacidades do Bateia;
- preencher oráculos humanos;
- adicionar outro provedor;
- executar Katamari competitivo;
- melhorar ranking ou screening;
- aumentar recall;
- retry com backoff além dos fallbacks declarados;
- persistência distribuída de jobs;
- interface gráfica.

## Próximo passo previsto

Com a descoberta finita e diagnosticável, o próximo passo deve congelar oráculos
das capacidades internas do Bateia e criar uma ontologia lexical/evidencial
pequena para seleção de software. Somente depois repetiremos o recall competitivo
e, havendo bundles correspondentes, o Katamari real.
