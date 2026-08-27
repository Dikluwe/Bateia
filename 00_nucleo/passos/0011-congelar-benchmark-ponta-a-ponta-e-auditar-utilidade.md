# Passo 0011 — Congelar benchmark ponta a ponta e auditar utilidade

**Estado**: `EM VALIDAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0010 — Diagnosticar inconclusões e coletar adaptativamente`](0010-diagnosticar-inconclusoes-e-coletar-adaptativamente.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1 e L2  
**Natureza**: benchmark causal, oráculos versionados e auditoria de produto sem LLM  
**Regime previsto**: protocolo completo de materialização segregada Tekt

## Evidência que motiva este passo

Os Passos 0009 e 0010 produziram uma descoberta importante:

- o funil e a recuperação preservaram `INCONCLUSIVE` corretamente;
- nos dois runs live, a recuperação resolveu `8/16` findings;
- no agregado, resolveu somente `8/28 = 28,6%`;
- nenhuma decisão final mudou;
- os cinco replays usaram bundles que não correspondiam às respectivas ideias.

Os replays não correspondentes são controles úteis para provar que o Bateia não
inventa certeza. Porém, não podem sustentar uma alegação de baixa ou alta
utilidade: quando ideia, candidatos e bundles não pertencem à mesma cadeia
causal, `INCONCLUSIVE` é o único resultado honesto.

Antes de alterar regras de operação, targets ou budgets, precisamos de um corpus
que consiga distinguir:

```text
motor conservador correto
    ≠
motor incapaz de concluir quando existe evidência suficiente
    ≠
corpus contaminado ou não correspondente
```

## Incerteza que este passo reduz

Não sabemos ainda se o Bateia consegue produzir decisões úteis quando recebe:

1. uma ideia bem formada;
2. descoberta correspondente à ideia;
3. candidatos realmente relacionados;
4. bundles completos no commit declarado;
5. evidência positiva, negativa e opaca suficiente;
6. política de reúso conhecida;
7. julgamento humano congelado antes da execução candidata.

Este passo deve construir essa base e impedir que mudanças futuras sejam
calibradas sobre exemplos trocados, conhecidos pela implementação ou escolhidos
depois de observar o resultado.

## Intenção

Materializar um benchmark ponta a ponta, completamente reexecutável offline, que
avalie separadamente:

- correspondência entre ideia, descoberta, candidato e bundle;
- precisão dos estados de evidência;
- precisão da decisão de reúso;
- cobertura conclusiva;
- preservação correta de `INCONCLUSIVE`;
- utilidade do dossiê para uma pessoa escolher um ponto de partida;
- ganho e custo da recuperação adaptativa;
- resistência a mutações semanticamente negativas.

O benchmark não é um conjunto de arquivos soltos. É uma cadeia causal congelada
com hashes, proveniência, oráculos e divisão de avaliação.

## Restrição central: benchmark antes de heurística

Durante este passo, é proibido alterar para elevar métricas:

- regras de `PROVEN`, `PARTIAL`, `ABSENT`, `CONFLICT` ou `INCONCLUSIVE`;
- regras de `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD`;
- thresholds de screening;
- targets de recuperação;
- budgets padrão;
- expected outcomes depois de executar a implementação candidata.

Se o benchmark revelar uma falha de regra, ela vira finding versionado para um
passo posterior. O benchmark não corrige o motor que está julgando.

## Unidade do benchmark

```go
type BenchmarkCaseID string

type EndToEndBenchmarkCase struct {
    SchemaVersion      string                    `json:"schema_version"`
    ID                 BenchmarkCaseID           `json:"id"`
    Idea               IdeaSpec                  `json:"idea"`
    Discovery          DiscoverySnapshot         `json:"discovery"`
    ScreeningModel     ScreeningModel            `json:"screening_model"`
    Bundles            []EvidenceBundle          `json:"bundles"`
    EvidenceRules      CapabilityRuleCatalog     `json:"evidence_rules"`
    ReusePolicy        ReuseDecisionPolicy       `json:"reuse_policy"`
    RecoveryPolicy     RecoveryPolicy            `json:"recovery_policy"`
    ExpectedEvidence   []EvidenceJudgment        `json:"expected_evidence"`
    ExpectedDecision   ReuseDecisionJudgment     `json:"expected_decision"`
    ExpectedRecovery   RecoveryJudgment          `json:"expected_recovery"`
    UtilityJudgment    DossierUtilityJudgment    `json:"utility_judgment"`
    Provenance         BenchmarkProvenance       `json:"provenance"`
    Split              BenchmarkSplit            `json:"split"`
}
```

Os objetos embutidos podem ser armazenados por referência relativa e hash no
manifesto físico, desde que a leitura resolva um caso imutável equivalente.

## Julgamento de decisão

```go
type ReuseDecisionJudgment struct {
    Expected       ReuseDecision   `json:"expected"`
    Acceptable     []ReuseDecision `json:"acceptable,omitempty"`
    Primary        string          `json:"primary,omitempty"`
    Components     []string        `json:"components,omitempty"`
    RequiredRules  []string        `json:"required_rules"`
    RequiredSources []string       `json:"required_sources"`
    Forbidden      []ReuseDecision `json:"forbidden"`
    Rationale      string          `json:"rationale"`
}
```

`Acceptable` existe somente para equivalências justificadas antes do replay. Não
pode ser preenchido depois para absorver uma divergência.

Exemplos:

- `EXTEND` e `INTEGRATE` só podem ser ambos aceitáveis quando a fronteira entre
  núcleo e componente estiver explicitamente ambígua no caso;
- `BUILD` nunca é aceitável quando há falha, budget incompleto ou revisão opaca;
- `INCONCLUSIVE` é obrigatório nos casos construídos para opacidade.

## Julgamento de recuperação

```go
type RecoveryJudgment struct {
    RecoverableCauses   []InconclusiveCause `json:"recoverable_causes"`
    TerminalCauses      []InconclusiveCause `json:"terminal_causes"`
    ExpectedNewKinds    []EvidenceKind      `json:"expected_new_kinds"`
    ForbiddenTargets    []string            `json:"forbidden_targets"`
    MaxRounds           int                 `json:"max_rounds"`
    MaxAdditionalBytes  int64               `json:"max_additional_bytes"`
    ExpectedFinal       ReuseDecision       `json:"expected_final"`
}
```

O julgamento não prescreve paths internos quando isso revelaria a solução. Ele
prescreve causas, tipos de evidência, limites e alvos proibidos suficientes para
rejeitar coleta genérica ou repetida.

## Julgamento de utilidade

```go
type DossierUtilityLabel string

const (
    UtilityUseful      DossierUtilityLabel = "USEFUL"
    UtilityNotUseful   DossierUtilityLabel = "NOT_USEFUL"
    UtilityUnjudgeable DossierUtilityLabel = "UNJUDGEABLE"
)

type DossierUtilityJudgment struct {
    Label              DossierUtilityLabel `json:"label"`
    DecisionActionable bool                `json:"decision_actionable"`
    EvidenceSufficient bool                `json:"evidence_sufficient"`
    BlockersVisible    bool                `json:"blockers_visible"`
    RepositoriesOpened int                 `json:"repositories_opened"`
    Rationale          string              `json:"rationale"`
}
```

Um dossiê é útil quando permite aceitar, rejeitar ou adiar conscientemente um
ponto de partida. Portanto, um `INCONCLUSIVE` pode ser útil quando identifica
precisamente o bloqueio e evita investigação desnecessária. Ele não é útil quando
apenas repete “faltam evidências” sem indicar causa ou próxima ação.

## Proveniência

```go
type BenchmarkProvenance struct {
    CapturedAt       string            `json:"captured_at"`
    SourceMode       IdeaRunMode       `json:"source_mode"`
    Provider         string            `json:"provider,omitempty"`
    OriginalRunID    string            `json:"original_run_id"`
    OriginalRunHash  string            `json:"original_run_hash"`
    ArtifactHashes   map[string]string `json:"artifact_hashes"`
    CandidateCommits map[string]string `json:"candidate_commits"`
    CollectorVersion string            `json:"collector_version"`
    Notes            []string          `json:"notes,omitempty"`
}
```

Não registrar tokens, cabeçalhos ou URLs assinadas. Cases remotos precisam usar
commits imutáveis. Cases locais registram identidade `local/...`, hash da
worktree e conteúdo necessário ao replay; não inventam remoto público.

## Gate de correspondência

Antes de avaliar o motor, validar o corpus:

```go
type CorrespondenceVerdict string

const (
    CorrespondencePreserved CorrespondenceVerdict = "PRESERVED"
    CorrespondenceViolated  CorrespondenceVerdict = "VIOLATED"
    CorrespondenceUnknown   CorrespondenceVerdict = "UNKNOWN"
)
```

Um caso é `PRESERVED` somente quando:

- `DiscoverySnapshot.Idea` corresponde byte a byte à ideia do caso após
  canonicalização definida;
- todo candidato julgado aparece na descoberta;
- todo bundle pertence a candidato descoberto;
- `repo@commit` do bundle coincide com a proveniência congelada;
- todo item pertence ao commit do bundle e possui hash válido;
- screening, assessment e decisão podem ser recomputados sem artefato externo;
- expected evidence referencia capacidades da ideia ou gates declarados;
- política, regras e modelo possuem versão e hash;
- não existe dependência de rede para replay.

É `VIOLATED` quando há troca comprovada de ideia, candidato, commit, bundle,
política ou regra. É `UNKNOWN` quando a identidade não pode ser provada — por
exemplo, bundle em `HEAD`, hash ausente ou snapshot legado sem origem suficiente.

`UNKNOWN` nunca entra nas métricas de precisão como sucesso. Ele aparece numa
fila explícita de recaptura.

## Mutação negativa obrigatória

O validador precisa rejeitar, no mínimo:

1. bundle de outro case com vocabulário parecido;
2. bundle correto em commit diferente;
3. candidato presente no bundle, ausente na descoberta;
4. `DiscoverySnapshot.Idea` trocada depois da coleta;
5. conteúdo alterado sem atualizar hash;
6. regra alterada sem atualizar manifesto;
7. política alterada mantendo expected decision antigo incompatível;
8. remoção do teste que sustentava `PROVEN`;
9. falha operacional convertida em bundle completo;
10. revisão `HEAD` apresentada como imutável;
11. caso opaco rotulado como decisão conclusiva;
12. expected outcome modificado depois do recibo de execução.

O gate discriminatório exige:

```text
mutation_score = mutações negativas rejeitadas / mutações válidas = 1,0
```

Qualquer mutação sobrevivente impede selar o benchmark.

## Política de `UNKNOWN`

`UNKNOWN` é resultado obrigatório para:

- identidade de candidato ambígua;
- commit não resolvido;
- parser ou linguagem sem suporte suficiente;
- artifact hash ausente ou inválido sem cópia recuperável;
- cadeia causal incompleta;
- julgamento humano sem justificativa ou proveniência;
- caso que dependa de rede para ser compreendido.

Casos `UNKNOWN` não são descartados silenciosamente. Devem aparecer em
`correspondence-unknown.json` com razão e ação de recaptura.

## Composição inicial do corpus

Congelar no mínimo 15 cases correspondentes:

| Classe | Quantidade mínima |
|---|---:|
| `FORK` | 2 |
| `EXTEND` | 2 |
| `COMPOSE` | 2 |
| `INTEGRATE` | 2 |
| `REFERENCE` | 2 |
| `BUILD` conclusivo | 1 |
| `INCONCLUSIVE` recuperável | 2 |
| `INCONCLUSIVE` terminal | 2 |

Distribuição adicional:

- pelo menos cinco ideias funcionalmente distintas;
- pelo menos quatro linguagens;
- pelo menos três cases locais, um sem Git/remoto;
- pelo menos três cases com política de operação local;
- pelo menos dois conflitos de licença;
- pelo menos dois projetos arquivados ou com manutenção problemática;
- pelo menos quatro cases com componentes reutilizáveis verificáveis;
- pelo menos três controles lexicalmente parecidos e funcionalmente irrelevantes;
- pelo menos dois cases que exijam recuperação e dois em que recuperar seja
  proibido.

Um único case pode satisfazer mais de uma dimensão, mas não duas classes de
decisão incompatíveis.

## Evidência de operação: subcorpus obrigatório

Como `UNKNOWN_OPERATION` permaneceu aberto nos runs do Passo 0010, o benchmark
deve conter exemplos positivos, negativos e opacos:

### Positivo local

- comando documentado ou entrypoint verificável;
- manifesto/runtime suficiente;
- ausência de envio obrigatório de código comprovada somente dentro do escopo
  declarado;
- teste ou exemplo local quando disponível.

### Negativo

- upload obrigatório de repositório ou código;
- dependência exclusiva de serviço remoto;
- autenticação remota indispensável;
- fonte exata que comprova o conflito.

### Opaco

- manifesto presente sem afirmação de localidade;
- README sem instrução operacional;
- chamadas externas ambíguas;
- coleta incompleta ou budget esgotado.

Manifesto isolado nunca é oráculo de operação local. O subcorpus deve rejeitar
explicitamente essa mutação.

## Divisão causal do corpus

```go
type BenchmarkSplit string

const (
    SplitContract  BenchmarkSplit = "CONTRACT"
    SplitTrain     BenchmarkSplit = "TRAIN"
    SplitValidation BenchmarkSplit = "VALIDATION"
    SplitTest      BenchmarkSplit = "TEST"
)
```

- `CONTRACT`: casos mínimos positivos, negativos e opacos que definem semântica;
- `TRAIN`: pode orientar políticas futuras, nunca expected outcomes;
- `VALIDATION`: escolhe entre políticas previamente declaradas;
- `TEST`: permanece fechado para implementação e calibração até o veredito.

Casos derivados da mesma ideia, repositório ou fixture devem permanecer no mesmo
grupo para evitar vazamento. O split é feito por grupo causal, não por arquivo.

## Manifesto Tekt do benchmark

Antes de implementar o runner, materializar:

```go
type BenchmarkManifest struct {
    ProtocolVersion   string                    `json:"protocol_version"`
    PromptPath        string                    `json:"prompt_path"`
    PromptHash        string                    `json:"prompt_hash"`
    BaselineCommit    string                    `json:"baseline_commit"`
    ContractPath      string                    `json:"contract_path"`
    ContractHash      string                    `json:"contract_hash"`
    UnknownPolicyHash string                    `json:"unknown_policy_hash"`
    CaseHashes        map[BenchmarkCaseID]string `json:"case_hashes"`
    SplitGroups       map[BenchmarkCaseID]string `json:"split_groups"`
    RoleCapabilities  []RoleCapabilityReceipt   `json:"role_capabilities"`
}
```

O manifesto congela prompt, baseline, contrato, política de `UNKNOWN`, cases e
splits antes de existir resultado candidato.

## Papéis previstos

### Autor da obrigação

- lê desejo, passos e laudos anteriores;
- escreve este passo e o prompt L0;
- não aprova a implementação.

### Autor do contrato

- lê prompt, baseline e tipos públicos necessários;
- escreve contratos e validador de correspondência;
- não lê patch do runner candidato.

### Autor dos oráculos

- recebe prompt e contrato congelados;
- escreve expected evidence, decisão e utilidade;
- não recebe saídas do motor para ajustar julgamentos.

### Adversário

- recebe manifesto, contrato e baseline;
- produz mutações negativas e opacas;
- não corrige a implementação.

### Implementador

- recebe prompt, contrato e oráculos selados;
- escreve runner e relatórios;
- não altera cases, splits ou expected outcomes.

### Verificador

- executa correspondence gate, mutation gate e benchmark;
- emite recibos e veredito;
- não edita o material avaliado.

Se o ambiente não conseguir isolar essas autoridades, a execução deve declarar
“executado sem atestação de isolamento” e não usar o selo como prova formal.

## Runner offline

```bash
bateia benchmark validate \
  --manifest benchmark/manifest.json \
  --output benchmark/validation.json

bateia benchmark mutate \
  --manifest benchmark/manifest.json \
  --mutations benchmark/mutations \
  --output benchmark/mutation-report.json

bateia benchmark run \
  --manifest benchmark/manifest.json \
  --split test \
  --output benchmark/results/test
```

Os três comandos são offline. `run` reexecuta screening, assessment, decisão e,
quando julgado recuperável, recuperação sobre os bytes congelados. Não consulta
provedor nem usa cache externo.

## Resultado por case

```go
type BenchmarkCaseResult struct {
    CaseID                  BenchmarkCaseID      `json:"case_id"`
    Correspondence         CorrespondenceVerdict `json:"correspondence"`
    EvidenceMatches        int                  `json:"evidence_matches"`
    EvidenceDivergences    int                  `json:"evidence_divergences"`
    DecisionObserved       ReuseDecision        `json:"decision_observed"`
    DecisionAccepted       bool                 `json:"decision_accepted"`
    ForbiddenDecision      bool                 `json:"forbidden_decision"`
    RecoveryExecuted       bool                 `json:"recovery_executed"`
    RecoveryFinal          ReuseDecision        `json:"recovery_final,omitempty"`
    Utility                DossierUtilityLabel  `json:"utility"`
    Witnesses              []string             `json:"witnesses"`
}
```

Toda divergência possui testemunha: regra, capability, source, hash ou gate que
permite reproduzi-la.

## Métricas agregadas

```go
type EndToEndBenchmarkMetrics struct {
    CorrespondencePreservedRate float64  `json:"correspondence_preserved_rate"`
    CorrespondenceUnknownRate   float64  `json:"correspondence_unknown_rate"`
    EvidencePrecision           *float64 `json:"evidence_precision"`
    EvidenceRecall              *float64 `json:"evidence_recall"`
    DecisionAccuracy            *float64 `json:"decision_accuracy"`
    ForbiddenDecisionRate       float64  `json:"forbidden_decision_rate"`
    ConclusiveCoverage          float64  `json:"conclusive_coverage"`
    CorrectInconclusiveRate     *float64 `json:"correct_inconclusive_rate"`
    RecoveryResolutionRate      *float64 `json:"recovery_resolution_rate"`
    RecoveryUsefulEvidenceRate  *float64 `json:"recovery_useful_evidence_rate"`
    UtilityRate                 *float64 `json:"utility_rate"`
}
```

Não calcular precisão sobre cases `VIOLATED` ou `UNKNOWN`. O relatório mostra
quantos casos foram excluídos e por quê; nunca transforma denominador vazio em
100%.

## Gates de aprovação

### Correspondência

- `correspondence_preserved_rate = 100%` nos cases selados;
- `correspondence_unknown_rate = 0%` no split `TEST`;
- `100%` das mutações de troca de bundle, commit, ideia e política rejeitadas.

### Poder discriminatório

- `mutation_score = 1,0`;
- positivos válidos → `PRESERVED`;
- mutações negativas → `VIOLATED` com testemunha;
- construções deliberadamente opacas → `UNKNOWN`;
- reordenação e repetição → mesmo veredito.

### Evidência e decisão

- `evidence_precision ≥ 90%`;
- `evidence_recall ≥ 85%`;
- `decision_accuracy ≥ 80%`;
- `forbidden_decision_rate = 0%`;
- `correct_inconclusive_rate = 100%`;
- nenhum falso `BUILD` por falha, budget ou identidade opaca.

### Recuperação

- `recovery_resolution_rate ≥ 40%` nos findings julgados recuperáveis;
- `recovery_useful_evidence_rate ≥ 30%`;
- `0` target repetido;
- `0` terceira rodada;
- `100%` das transições sustentadas por evidência nova.

### Produto

- `utility_rate ≥ 80%`;
- `100%` dos blockers eliminatórios visíveis;
- mediana de repositórios abertos pelo revisor ≤ 1;
- `UNJUDGEABLE` não conta como sucesso e exige motivo.

## Testes obrigatórios

1. case correspondente válido produz `PRESERVED`;
2. ideia trocada produz `VIOLATED`;
3. bundle de outro case produz `VIOLATED`;
4. commit trocado produz `VIOLATED`;
5. conteúdo com hash inválido produz `VIOLATED`;
6. `HEAD` produz `UNKNOWN`, nunca `PRESERVED`;
7. candidato não descoberto produz `VIOLATED`;
8. regra sem versão/hash produz `UNKNOWN`;
9. expected evidence fora da ideia/gates produz `VIOLATED`;
10. caso opaco exige `INCONCLUSIVE`;
11. falha operacional nunca aceita `BUILD`;
12. bundle completo sem implementação não prova capacidade;
13. implementação sem teste respeita independência exigida;
14. manifesto isolado não prova operação local;
15. upload obrigatório produz conflito de operação;
16. evidência local positiva exige fonte rastreável;
17. licença incompatível não desaparece no score funcional;
18. alteração de política invalida expected incompatível;
19. remoção de evidência forte derruba o resultado esperado;
20. ordem de candidatos, bundles e evidências não altera resultado;
21. duplicação de conteúdo não aumenta cobertura;
22. split mantém grupos causais juntos;
23. runner não lê casos `TEST` durante calibração;
24. replay não toca rede, cache externo ou relógio não controlado;
25. denominador vazio serializa métrica `null`;
26. toda divergência possui testemunha reproduzível;
27. utilidade não pode ser inferida da decisão automaticamente;
28. julgamento alterado após recibo invalida o selo;
29. mutation score abaixo de 1 impede certificação;
30. linter V15 permanece sem violações e cada consumidor possui prompt próprio.

## Critério de produto

O benchmark só promove os Passos 0009 e 0010 de `EM VALIDAÇÃO` quando:

1. o manifesto e o contrato estiverem selados antes do split `TEST`;
2. todos os cases `TEST` forem correspondentes;
3. os gates de evidência, decisão, recuperação e utilidade forem atingidos;
4. revisores não precisarem conhecer a implementação para justificar o veredito;
5. resultados puderem ser reproduzidos offline a partir do repositório;
6. nenhuma divergência tiver sido absorvida por edição retrospectiva do oráculo.

Se os gates falharem, 0009 e 0010 permanecem `EM VALIDAÇÃO`. O relatório deve
atribuir cada perda a corpus, regra, coleta, decisão, recuperação ou utilidade.

## Critério de parada

Interromper o passo se:

- não for possível provar correspondência entre ideia e bundles;
- o split `TEST` vazar para implementação ou calibração;
- expected outcome for alterado depois de observar o motor;
- alguma mutação negativa sobreviver;
- `UNKNOWN` for contado como sucesso;
- utilidade for inferida automaticamente sem julgamento;
- cases remotos dependerem da rede para replay;
- o benchmark exigir afrouxar estados ou decisões existentes.

## Fora de escopo

- alterar regras de evidência ou decisão;
- ajustar targets ou budgets da recuperação;
- LLM, embeddings ou geração probabilística;
- executar código ou testes dos candidatos;
- aconselhamento jurídico;
- interface gráfica;
- benchmark público hospedado ou serviço de leaderboard;
- inferir julgamento humano a partir do output do Bateia;
- promover 0009/0010 apenas porque os gates técnicos passaram.

## Próximo passo previsto

Depois do benchmark selado, as divergências do split `VALIDATION` poderão originar
passos pequenos e específicos — por exemplo, evidência de operação local ou
targeting de testes — enquanto o split `TEST` permanece fechado. Somente depois
dessas correções o benchmark será reexecutado para decidir a promoção de 0009 e
0010.
