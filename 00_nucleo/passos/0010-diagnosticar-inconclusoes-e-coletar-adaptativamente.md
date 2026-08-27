# Passo 0010 — Diagnosticar inconclusões e coletar adaptativamente

**Estado**: `EM VALIDAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0009 — Orquestrar o funil e entregar um dossiê de reúso`](0009-orquestrar-funil-e-entregar-dossie-de-reuso.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2 e L4  
**Natureza**: diagnóstico causal, recuperação limitada e reavaliação sem LLM

## Evidência que motiva este passo

O Passo 0009 materializou o funil ponta a ponta e executou:

- dois runs `live`;
- dois runs `local`;
- cinco runs `replay`.

Todos os nove dossiês terminaram `INCONCLUSIVE`. O comportamento é seguro —
nenhuma falha virou `BUILD` —, mas ainda não entrega uma decisão de partida útil.

A próxima incerteza não é se precisamos de outro score. Precisamos descobrir se
uma inconclusão contém informação suficiente para produzir uma coleta adicional
pequena, causal e capaz de resolvê-la.

## Intenção

Transformar uma decisão `INCONCLUSIVE` num diagnóstico acionável:

```text
dossiê inconclusivo
    ↓
causas separadas por candidato, capacidade e gate
    ↓
alvos de recuperação derivados das lacunas
    ↓
coleta adicional limitada
    ↓
novo assessment
    ↓
nova decisão ou INCONCLUSIVE final fundamentado
```

O motor não deve buscar “mais coisas” genericamente. Cada novo alvo precisa
declarar qual inconclusão tenta resolver, qual evidência seria suficiente e qual
regra consumirá o resultado.

## Restrição central: não afrouxar evidência

Este passo não pode:

- reduzir o limiar de `PROVEN`;
- promover documentação isolada a implementação;
- transformar falha de rede em `ABSENT` ou `BUILD`;
- alterar uma regra para beneficiar um candidato conhecido;
- gerar queries ou interpretações com LLM;
- repetir silenciosamente a mesma coleta;
- esconder custo adicional no budget original.

O sucesso ocorre quando adquirimos evidência nova e suficiente, não quando
mudamos o significado dos estados existentes.

## Taxonomia inicial de inconclusões

```go
type InconclusiveCause string

const (
    CauseUnresolvedRevision   InconclusiveCause = "UNRESOLVED_REVISION"
    CauseIncompleteCollection InconclusiveCause = "INCOMPLETE_COLLECTION"
    CauseMissingImplementation InconclusiveCause = "MISSING_IMPLEMENTATION_EVIDENCE"
    CauseMissingTest          InconclusiveCause = "MISSING_TEST_EVIDENCE"
    CauseUnknownLicense       InconclusiveCause = "UNKNOWN_LICENSE"
    CauseUnknownOperation     InconclusiveCause = "UNKNOWN_OPERATION"
    CauseConflictingEvidence  InconclusiveCause = "CONFLICTING_EVIDENCE"
    CauseBudgetExhausted      InconclusiveCause = "BUDGET_EXHAUSTED"
    CauseProviderFailure      InconclusiveCause = "PROVIDER_FAILURE"
    CauseUnsupportedArtifact  InconclusiveCause = "UNSUPPORTED_ARTIFACT"
    CauseRuleCoverageGap      InconclusiveCause = "RULE_COVERAGE_GAP"
    CauseInsufficientCorpus   InconclusiveCause = "INSUFFICIENT_CORPUS"
)
```

As causas não são mutuamente exclusivas. Um candidato pode ter licença
desconhecida, budget esgotado e ausência de teste ao mesmo tempo. Nenhuma soma ou
prioridade pode ocultar um gate eliminatório.

## Contrato do diagnóstico

```go
type InconclusiveFinding struct {
    ID          string             `json:"id"`
    Cause       InconclusiveCause  `json:"cause"`
    Candidate   string             `json:"candidate,omitempty"`
    Capability  string             `json:"capability,omitempty"`
    RuleID      string             `json:"rule_id"`
    Evidence    []EvidenceSource   `json:"evidence,omitempty"`
    Failures    []ExecutionFailure `json:"failures,omitempty"`
    Missing     []string           `json:"missing"`
    Recoverable bool               `json:"recoverable"`
    Message     string             `json:"message"`
}

type InconclusiveDiagnosis struct {
    SchemaVersion string                `json:"schema_version"`
    RunID         string                `json:"run_id"`
    DecisionHash  string                `json:"decision_hash"`
    Findings      []InconclusiveFinding `json:"findings"`
    Recoverable   int                   `json:"recoverable"`
    Terminal      int                   `json:"terminal"`
}
```

O diagnóstico é puro e offline. Ele consome `run.json`, `reuse-decision.json`,
`assessment-report.json`, planos e bundles do dossiê. Não chama provedor.

## Política de recuperação

```go
type RecoveryPolicy struct {
    SchemaVersion        string `json:"schema_version"`
    MaxRounds            int    `json:"max_rounds"`
    MaxCandidates        int    `json:"max_candidates"`
    AdditionalDocuments int    `json:"additional_documents"`
    AdditionalCodeFiles int    `json:"additional_code_files"`
    AdditionalTests     int    `json:"additional_tests"`
    AdditionalBytes     int64  `json:"additional_bytes"`
    AdditionalTreeDepth int    `json:"additional_tree_depth"`
    RetryProviderFailure bool  `json:"retry_provider_failure"`
}
```

Defaults do primeiro experimento:

```text
max_rounds             = 2
max_candidates         = 3
additional_documents   = 4
additional_code_files  = 8
additional_tests       = 6
additional_bytes       = 2 MiB
additional_tree_depth  = 2
retry_provider_failure = true
```

Alterar a política não modifica o dossiê original. Política e hash aparecem em
cada rodada de recuperação.

## Plano de recuperação

```go
type RecoveryTarget struct {
    FindingID    string           `json:"finding_id"`
    Candidate    CommitRef        `json:"candidate"`
    Capability   string           `json:"capability,omitempty"`
    Path         string           `json:"path,omitempty"`
    Patterns     []string         `json:"patterns,omitempty"`
    Kind         EvidenceKind     `json:"kind"`
    RuleID       string           `json:"rule_id"`
    Expected     []EvidenceStatus `json:"expected"`
    Reason       string           `json:"reason"`
}

type EvidenceRecoveryPlan struct {
    SchemaVersion string                `json:"schema_version"`
    Round         int                   `json:"round"`
    DiagnosisHash string                `json:"diagnosis_hash"`
    Policy        RecoveryPolicy        `json:"policy"`
    PolicyHash    string                `json:"policy_hash"`
    Targets       []RecoveryTarget      `json:"targets"`
    Skipped       []InconclusiveFinding `json:"skipped"`
}
```

Todo target aponta para um `FindingID`. Um finding não recuperável aparece em
`skipped` com razão; nunca desaparece do relatório.

## Regras iniciais de recuperação

### Revisão não resolvida

- repetir somente a resolução de metadata;
- exigir SHA imutável antes de qualquer coleta de código;
- registrar redirects, rate limit e falha do provedor;
- não usar `HEAD` como evidência.

### Implementação sem teste

- partir dos símbolos e paths já encontrados;
- procurar testes por nome de símbolo, pacote e convenções da linguagem;
- limitar busca aos paths novos;
- não reler testes cujo hash já esteja no bundle.

### Documentação sem implementação

- extrair termos concretos da regra versionada e da capacidade, não texto gerado;
- priorizar paths de código relacionados aos termos e à árvore já observada;
- exigir símbolo, implementação ou comportamento verificável para `PROVEN`.

### Licença desconhecida

- procurar nomes reconhecidos como `LICENSE`, `COPYING`, `NOTICE` e metadata do
  provedor;
- preservar ambiguidade entre licença do repositório e licença de componente;
- não emitir aconselhamento jurídico.

### Operação local desconhecida

- procurar manifesto, comandos documentados, serviços obrigatórios, bancos,
  chamadas externas e configuração de runtime;
- comprovar localidade somente com evidência positiva;
- ausência de documentação não comprova operação offline.

### Budget esgotado

- usar somente o budget adicional da política;
- priorizar targets ligados a gates e capacidades obrigatórias;
- registrar quais targets perderam a disputa pelo budget;
- budget novamente esgotado encerra a causa como inconclusiva na rodada.

### Falha de provedor

- repetir no máximo uma vez por rodada quando permitido;
- respeitar `Retry-After`, timeout e cancelamento;
- não trocar silenciosamente de provedor;
- a falha repetida permanece `PROVIDER_FAILURE`.

### Lacuna de cobertura de regra

- produzir diagnóstico e amostra para revisão humana;
- não criar ou ampliar regra automaticamente;
- não executar coleta sem padrão versionado capaz de consumir a evidência.

## Deduplicação e novidade

Antes de executar um target, calcular sua identidade:

```text
repo@commit + kind + path/pattern + capability + rule_id
```

Um target é novo somente se:

- não foi executado no run original;
- não apareceu em rodada anterior;
- não é semanticamente equivalente após normalização determinística;
- pode produzir um tipo de evidência ainda ausente.

Conteúdo retornado com hash já conhecido é registrado como duplicado e não conta
como evidência nova nem como consumo útil do budget.

## Rodadas

Cada rodada produz:

```text
recovery/round-01/diagnosis.json
recovery/round-01/plan.json
recovery/round-01/bundles.json
recovery/round-01/assessment-report.json
recovery/round-01/reuse-decision.json
recovery/round-01/metrics.json
```

Ao final da rodada:

1. unir evidências antigas e novas por identidade e hash;
2. reexecutar o mesmo assessment versionado;
3. reexecutar a mesma política de decisão;
4. comparar estados anteriores e posteriores;
5. encerrar quando houver decisão conclusiva, não houver target novo ou o limite
   de rodadas for atingido.

Uma decisão pode continuar `INCONCLUSIVE`. Isso é resultado válido quando todas
as tentativas e limitações permanecem observáveis.

## Transições permitidas

```text
INCONCLUSIVE → PROVEN
INCONCLUSIVE → PARTIAL
INCONCLUSIVE → ABSENT
INCONCLUSIVE → CONFLICT
INCONCLUSIVE → INCONCLUSIVE
```

Uma recuperação nunca altera diretamente o estado. Ela acrescenta evidências e o
assessment existente calcula a transição.

`ABSENT` só pode surgir quando a busca planejada foi concluída dentro do escopo
declarado. Falha, target pulado ou budget esgotado preserva `INCONCLUSIVE`.

## Métricas

```go
type RecoveryMetrics struct {
    Round                  int      `json:"round"`
    FindingsBefore         int      `json:"findings_before"`
    TargetsPlanned         int      `json:"targets_planned"`
    TargetsExecuted        int      `json:"targets_executed"`
    TargetsNew             int      `json:"targets_new"`
    DuplicateContent       int      `json:"duplicate_content"`
    AdditionalBytes        int64    `json:"additional_bytes"`
    ProviderCalls          int      `json:"provider_calls"`
    StatusTransitions      int      `json:"status_transitions"`
    ResolvedFindings       int      `json:"resolved_findings"`
    RemainingFindings      int      `json:"remaining_findings"`
    DecisionBefore         string   `json:"decision_before"`
    DecisionAfter          string   `json:"decision_after"`
    UsefulEvidenceRate     *float64 `json:"useful_evidence_rate"`
    InconclusiveResolution *float64 `json:"inconclusive_resolution_rate"`
}
```

`useful_evidence_rate` mede targets novos que alteraram ao menos um estado de
evidência. Quantidade de bytes ou arquivos, isoladamente, não representa sucesso.

## Boca de inspeção

Diagnóstico offline:

```bash
bateia idea recover diagnose \
  --run .bateia/runs/minha-ideia \
  --output diagnosis.json
```

Plano offline:

```bash
bateia idea recover plan \
  --run .bateia/runs/minha-ideia \
  --diagnosis diagnosis.json \
  --policy recovery-policy.json \
  --output recovery-plan.json
```

Execução:

```bash
bateia idea recover execute \
  --run .bateia/runs/minha-ideia \
  --policy recovery-policy.json \
  --resume
```

`diagnose` e `plan` são completamente offline. Somente `execute` pode tocar
provedores, e apenas para targets presentes no plano congelado.

## Integração com o dossiê

O `run.json` original não é reescrito. Acrescentar um índice separado:

```text
recovery/index.json
```

O índice registra hashes do run original, diagnósticos, planos, bundles,
assessments e decisões por rodada. A decisão mais recente é uma derivação; a
decisão original permanece disponível.

`decision.md` pode ganhar um link para a recuperação, mas não substituir seu
conteúdo histórico silenciosamente.

## Testes obrigatórios

1. decisão diferente de `INCONCLUSIVE` não gera plano de recuperação;
2. diagnóstico separa capacidade, licença, operação, conflito e falha;
3. cada finding possui regra, causa, mensagem e origem observável;
4. falha de resolução produz somente target de metadata;
5. documentação sem implementação procura código, não mais documentação igual;
6. implementação sem teste procura testes relacionados;
7. licença desconhecida procura apenas artefatos de licença e metadata;
8. lacuna de regra não cria regra nem target automaticamente;
9. targets já executados não são repetidos;
10. conteúdo com hash conhecido não é incorporado duas vezes;
11. mudança de ordem não altera diagnóstico ou plano;
12. alterar política muda somente budget e targets afetados;
13. policy hash aparece em todas as rodadas;
14. `max_rounds` impede terceira rodada;
15. cancelamento interrompe novas coletas e preserva a rodada parcial;
16. budget esgotado nunca produz `ABSENT`;
17. falha repetida nunca produz `BUILD`;
18. evidência nova é reavaliada pelas regras existentes;
19. decisão conclusiva encerra novas rodadas;
20. ausência de targets novos encerra com motivo explícito;
21. replay da recuperação não chama rede;
22. mesma entrada produz diagnóstico e plano idênticos;
23. run e artifacts originais não são modificados;
24. paths de evidência não escapam do diretório do dossiê;
25. tokens e headers não aparecem nos artefatos de recuperação;
26. linter V15 permanece sem violações e cada consumidor possui prompt próprio.

## Experimento de produto

Aplicar a recuperação aos nove dossiês produzidos no Passo 0009, sem substituir
seus artefatos originais.

Separar os resultados em:

```text
resolvido corretamente
permaneceu inconclusivo corretamente
decisão conclusiva incorreta
custo excessivo sem evidência útil
```

Cada decisão alterada precisa de revisão manual cega às regras internas, usando
somente a ideia, o dossiê e as evidências citadas.

## Critérios de aprovação

- `100%` dos findings possuem causa e regra rastreáveis;
- `100%` dos targets apontam para um finding;
- `0` repetição de target já concluído;
- `0` segredo serializado;
- `0` falso `BUILD` causado por falha ou budget;
- `100%` das transições apontam para evidência nova;
- resolver corretamente ao menos `40%` das inconclusões recuperáveis na primeira
  amostra;
- decisões alteradas mantêm precisão manual ≥ `90%`;
- `useful_evidence_rate` ≥ `30%` no agregado;
- nenhum run ultrapassa a política de duas rodadas e 2 MiB adicionais por rodada.

O limiar inicial de resolução é deliberadamente modesto. O objetivo é demonstrar
coleta causal útil antes de tentar maximizar cobertura.

## Critério de parada

Interromper a implementação ou rebaixar o experimento se:

- for necessário alterar `PROVEN`, `ABSENT` ou `BUILD` para elevar a métrica;
- a maioria dos targets repetir conteúdo conhecido;
- uma regra não conseguir explicar por que pediu um arquivo;
- decisões mudarem sem evidência nova;
- o custo adicional não puder ser atribuído a findings;
- recuperação local depender de GitHub ou remoto;
- replay tocar rede;
- a segunda rodada não acrescentar estados nem evidências úteis.

## Fora de escopo

- LLM, embeddings ou geração probabilística;
- criação automática de regras de evidência;
- aconselhamento jurídico;
- execução de código, testes ou binários dos candidatos;
- mudança automática de provedor;
- terceira rodada de recuperação;
- UI/TUI;
- alteração automática da `IdeaSpec`;
- clonar, forkar, integrar ou modificar projetos candidatos;
- transformar `REFERENCE` ou `BUILD` em inconclusão para forçar recuperação.

## Próximo passo previsto

Se a coleta adaptativa resolver inconclusões com precisão e custo aceitáveis, o
passo seguinte poderá congelar um benchmark ponta a ponta e calibrar políticas de
budget por classe de ideia. Se não resolver, o próximo trabalho deverá atacar as
lacunas específicas de regras ou provedores reveladas pelas métricas, sem criar
mais uma camada genérica de ranking.

## Evidência de execução

Implementação materializada em 2026-08-26 com diagnóstico offline, plano
canonicalizado, histórico contra o run original, duas rodadas máximas, merge por
hash, reavaliação pelos contratos existentes e CLI
`recover diagnose|plan|execute`.

O diagnóstico encontrou 28 findings nos nove dossiês do Passo 0009. A recuperação
foi executada sobre todos eles. Nos dois runs live, resolveu 4 de 8 findings em
cada (`50%`) sem retirar a decisão segura `INCONCLUSIVE`. Nos demais runs não
houve resolução útil. O agregado ficou em 8/28 (`28,6%`), abaixo dos gates de 40%
de resolução e 30% de evidência útil.

Por isso o passo permanece `EM VALIDAÇÃO`: contratos, execução e gates técnicos
estão materializados, mas o experimento de produto ainda não atingiu os limiares
agregados. Ver [`Laudo 0010`](../laudos/0010-diagnosticar-inconclusoes-e-coletar-adaptativamente.md).
