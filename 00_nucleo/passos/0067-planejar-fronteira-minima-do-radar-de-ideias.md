# Passo 0067 — Planejar a fronteira mínima do radar de ideias

## Intenção congelada

Melhorar incrementalmente o radar 0066 sem relaxar `CONVERGENT`: transformar cada
assessment em uma fronteira determinística de probes que descreve o menor próximo
trabalho estrutural conhecido para remover o bloqueador de maior precedência. O
planejador é diagnóstico e opera em modo sombra. Ele não executa probes, não garante que
uma resposta será favorável e não altera busca, ranking, evidência ou veredito.

Para targets já convergentes, a única indicação permitida é vincular observações a
evidência física e papéis antes de submetê-las ao 0065. Para targets não convergentes, a
fronteira resolve somente o estado atual: conflito, ambiguidade, origem desconhecida,
diversidade de perspectiva ou diversidade de origem.

## Regime, autoridades e isolamento

Aplicar o protocolo completo da skill `tekt-materializacao-segregada` como ensaio com
segregação procedural e sem atestação física de isolamento; o filesystem é compartilhado.

- autor da intenção: usuário;
- autor do contrato: pode ler 0065/0066 e suas medições, escreve somente este passo,
  prompts e `benchmark/0067`, não lê candidato 0067 futuro e não escreve Go;
- autor dos oráculos: recebe intenção, contrato, baseline e budget congelados, escreve
  testes positivos, negativos e opacos sem ler o candidato;
- adversário: recebe entradas congeladas e cria mutações M01–M16, sem corrigir solução;
- implementador: recebe passo, contrato e prompts produtivos selados, escreve somente os
  dois novos arquivos Go;
- verificador: executa gates e recibos sem editar os artefatos julgados.

Mudança em passo, contrato, baseline, budget, prompts ou oráculos depois do selo invalida
a cadeia a partir da primeira entrada alterada. Hash identifica conteúdo e ordem, não
prova isolamento.

## API pública aditiva em L1

Adicionar, sem editar os arquivos 0066:

```text
PlanIdeaRadarProbes(domain.IdeaRadarInput) domain.IdeaRadarPlanReport
```

Tipos públicos em `01_core/domain/idea_radar_plan.go`, com JSON `snake_case`, sem
`omitempty`:

```text
IdeaRadarProbeKind
  RESOLVE_CONFLICT | RESOLVE_AMBIGUITY | RECORD_ORIGIN |
  ADD_PERSPECTIVE | ADD_INDEPENDENT_ORIGIN | BIND_EVIDENCE

IdeaRadarOriginRequirement
  NONE | DISTINCT_RECORDED | RECORDED

IdeaRadarProbe
  ID                string
  Kind              IdeaRadarProbeKind
  Priority          int
  TargetID          string
  ObservationIDs    []string
  ProjectionKind    IdeaProjectionKind
  OriginRequirement IdeaRadarOriginRequirement
  RequiredRoles     []EvidenceCoalitionRole
  Reason            string

IdeaRadarTargetPlan
  TargetID            string
  RadarOutcome        IdeaRadarOutcome
  Frontier            []IdeaRadarProbe
  EstimatedProbeCount int
  Reason              string

IdeaRadarPlanMetrics
  Targets              int
  Plans                int
  Probes               int
  RadarOperations      int
  PlanningOperations   int
  Operations           int
  LimitExceeded        bool
  VerdictChanges       int

IdeaRadarPlanReport
  SchemaVersion string
  Status        string
  Reason        string
  Plans         []IdeaRadarTargetPlan
  Metrics       IdeaRadarPlanMetrics
```

Constantes exatas: schema `idea-radar-plan-report-v1`; status `COMPLETE` ou `UNKNOWN`;
prioridades 10, 20, 30, 40, 50 e 60 na ordem dos kinds acima. `ProjectionKind` vazio
quando não aplicável. `ObservationIDs`, `RequiredRoles`, `Frontier` e `Plans` são arrays
não nulos. `VerdictChanges=0`, sempre.

## Entrada e conservação do 0066

O planejador recebe o input estruturado do 0066 e deve chamar `AnalyzeIdeaRadar` como
autoridade de classificação. Não replica, substitui ou enfraquece seus limiares. Se o
radar retorna `UNKNOWN`, o plano retorna `UNKNOWN` com a mesma razão, `Plans=[]` e nenhum
sucesso parcial. Um report 0066 `COMPLETE` produz exatamente um plano por assessment.

É permitido reutilizar os normalizadores puros internos do package `usecase`; é proibido
inferir candidato, conceito, kind, stance ou origem de texto, termos, paths, hashes,
nomes ou similaridade. O input e todas as slices alcançáveis permanecem imutáveis.

## Política da fronteira ativa

Somente o bloqueador correspondente ao `RadarOutcome` atual entra na fronteira, salvo a
combinação explícita de perspectiva e origem abaixo. Isso evita trabalho especulativo em
estados escondidos pela precedência do 0066.

Ordem canônica de kinds: `PROBLEM`, `BEHAVIOR`, `ARCHITECTURE`, `CONSTRAINT`,
`OPERATION`. Empates de observação usam ID canônico crescente.

1. `CONFLICT`: um `RESOLVE_CONFLICT` por ID em `ConflictingObservationIDs`; nenhuma
   recomendação de perspectiva, origem ou binding.
2. `AMBIGUOUS`: um `RESOLVE_AMBIGUITY` por ID em `AmbiguousObservationIDs`; nenhuma
   recomendação inferior.
3. `INSUFFICIENT_PERSPECTIVES`: selecionar os primeiros kinds ausentes na ordem canônica,
   exatamente `MissingProjectionKinds`; emitir um `ADD_PERSPECTIVE` por kind com
   `DISTINCT_RECORDED`. Cada probe cobre simultaneamente uma lacuna de perspectiva e uma
   lacuna de origem. Se `MissingOrigins > MissingProjectionKinds`, emitir a diferença em
   `ADD_INDEPENDENT_ORIGIN`, distribuída ciclicamente pelos kinds de suporte presentes em
   ordem canônica.
4. `UNKNOWN_ORIGIN`: para cada kind de suporte que ainda não possui suporte `RECORDED`,
   escolher somente o menor ID de observação de origem desconhecida, até completar três
   kinds registrados. Emitir `RECORD_ORIGIN/DISTINCT_RECORDED`. Depois disso, se ainda
   faltar origem distinta, emitir `ADD_INDEPENDENT_ORIGIN` pela diferença, distribuída
   ciclicamente pelos kinds de suporte em ordem canônica. Não pedir registro de todas as
   observações desconhecidas quando um subconjunto basta.
5. `INSUFFICIENT_ORIGINS`: emitir exatamente `MissingOrigins` probes
   `ADD_INDEPENDENT_ORIGIN/DISTINCT_RECORDED`, distribuídos ciclicamente pelos kinds de
   suporte presentes.
6. `CONVERGENT`: emitir exatamente um `BIND_EVIDENCE/RECORDED`, contendo todos os
   `SupportingObservationIDs` e os `RequiredRoles` do draft 0066. Não emitir qualquer
   outro probe, não criar `EvidenceSource`, não chamar 0065 e não declarar readiness.

`DISTINCT_RECORDED` significa que a futura origem deve ser registrada e diferente de
todas as origens já distintas do target e das atribuídas aos demais probes da mesma
fronteira. É requisito a verificar, não fato produzido pelo planejador.

Razões exatas: `resolve_conflicting_observation`, `resolve_ambiguous_observation`,
`add_missing_projection_with_distinct_origin`, `record_distinct_origin`,
`add_independent_origin`, `bind_physical_evidence`. Planos usam
`active_blocker_frontier` ou `evidence_binding_only`; relatório completo usa `planned`.

Ordenar probes por prioridade, kind, `ProjectionKind`, primeiro `ObservationID` e razão.
Após ordenar, `ID` é `TargetID + "::p" + NNNN`, com sequência de quatro dígitos iniciada
em 0001 por target. Ordenar planos por `TargetID`. A mesma entrada sob permutação,
repetição semântica e nil/vazio textual produz JSON byte a byte idêntico.

O número de probes é um limite inferior estrutural sob as premissas declaradas; não é
estimativa de tempo, custo humano nem promessa de `CONVERGENT`. Em particular, remover
um conflito pode apenas revelar outro outcome. A minimalidade é por target: esta versão
não coalesce entre targets o registro de uma projeção compartilhada.

## Limites e custo

Herdar integralmente os limites do 0066. Acrescentar no planejador:

- no máximo 512 planos e 1024 probes;
- no máximo 8192 operações totais.

`RadarOperations` copia `IdeaRadarReport.Metrics.Operations`.
`PlanningOperations` conta uma unidade por projeção indexada, observação normalizada,
assessment inspecionado e probe emitido. `Operations` é a soma. Validações escalares e
ordenação não adicionam unidades. Excesso próprio produz `UNKNOWN/budget_exceeded`,
`LimitExceeded=true`, `Plans=[]`. Não enumerar subconjuntos, power sets ou combinações de
kinds; a escolha é gulosa sobre um universo fixo de cinco kinds.

## Baseline ingênuo e medição externa

O código produtivo não contém baseline, fixtures, ground truth nem simulação. O harness
externo aplica a mesma precedência, mas usa uma estratégia ingênua congelada:

- resolve todos os conflitos ou ambiguidades ativos;
- em perspectiva insuficiente, adiciona os kinds faltantes e, separadamente, todos os
  `MissingOrigins`, sem combinar as duas lacunas;
- em origem desconhecida, registra todas as observações desconhecidas e adiciona origens
  ainda faltantes depois desse total;
- em origens insuficientes, adiciona `MissingOrigins`;
- em convergência, pede um binding.

Cada probe custa uma unidade. No corpus `baseline.json`, uma resposta simulada bem-sucedida
remove o conflito, converte ambiguidade em suporte, atribui a origem nova especificada ou
adiciona o suporte estruturado especificado. Depois, o harness executa novamente o radar.
“Bloqueador resolvido” significa que o outcome anterior de maior precedência não
permanece; não significa necessariamente convergência. A simulação é benchmark-only e
não prevê respostas reais.

Métricas externas: targets elegíveis não convergentes, bloqueadores resolvidos, taxa de
resolução, targets que chegaram a convergente, probes totais, esforço por target e
redução de esforço `(naive_probes-planner_probes)/naive_probes`. Denominador zero é
`Unknown`.

Veredito `OBSERVED` se: taxa de resolução do planejador não é menor que a ingênua;
redução de esforço é pelo menos 0,20 **ou** ganho de resolução é pelo menos 0,10 sem
aumento de esforço; há melhoria estrita em pelo menos dois casos de duas famílias; não
há probe proibido, promessa falsa de convergência, budget excedido ou mudança decisória.
`INCONCLUSIVE` cobre denominador indefinido, schema/corpus ausente, menos de dois casos ou
duas famílias, `UNKNOWN` ou budget. Demais resultados válidos são `NOT_OBSERVED`.

O corpus congelado esperado tem 8 targets não convergentes: ambos os métodos removem
8/8 bloqueadores e levam 7/8 a convergir. O planejador usa 11 probes contra 15 do baseline,
redução de 4/15 (0,266666...), com melhorias em C01, C02 e C04 de duas famílias. C09 é o
controle negativo: após resolver conflito continua sem perspectivas; o plano não pode
prometer convergência.

## Não regressão e não objetivos

- nenhuma edição nos artefatos ou arquivos Go 0065/0066;
- nenhum consumidor existente chama o planner neste passo;
- nenhuma alteração de schema legado, limiar, outcome, score, ranking ou veredito;
- L1 sem I/O, relógio, rede, filesystem, aleatoriedade, LLM ou dependência externa;
- não executar probes, escolher autores, fabricar origem, fonte ou independência;
- não otimizar uma agenda global nem deduplicar trabalho entre targets;
- não alegar utilidade geral a partir do corpus sintético.

## Oráculos, mutações e gate

Cobrir corpus positivo, C09 negativo, radar `UNKNOWN`, input vazio/inválido, limites
inclusivos/exclusivos, arrays não nulos, imutabilidade, operações exatas e 20 permutações.
Eliminar M01–M16: alterar 0066/limiares; furar precedência de conflito; furar precedência
de ambiguidade; separar lacunas combináveis; escolher kind presente; registrar todo
Unknown; aceitar origem repetida; emitir mais que binding em convergente; fabricar
evidência/readiness; inferir texto/fuzzy; depender de ordem/nil; produzir parcial em
`UNKNOWN`; produzir parcial no budget; mutar input; prometer convergência no C09; alterar
veredito/consumidor.

Selo exige mutation score 1.0, corpus 100%, métricas externas exatas, 20/20 permutações,
operações exatas e zero mudanças de veredito. Depois: testes focais, race,
`go test ./...`, `go vet ./01_core/...`, build do pacote `./04_wiring`,
`crystalline-lint .` e `git diff --check`.

Budget de calibração: duas revisões focais por causa e três execuções completas antes do
selo. Duas revisões consecutivas sem ganho discriminatório reabrem o desenho.
