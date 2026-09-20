# Prompt espelhado — `step0029_experiment` / harness do Passo 0029

Implemente futuramente um harness experimental determinístico conforme o contrato
R1 em `00_nucleo/benchmark/0029/contract.md`. Este prompt e o contrato formam uma
única obrigação: divergência entre eles invalida a materialização. Não implemente
dispatcher, integração de host, chamada de modelo, chamada MCP, filesystem,
processo, rede, credencial ou cobrança.

## Linhagem e fronteira

- pacote público: `step0029_experiment`;
- natureza: consumidor externo/benchmark, não componente do produto;
- não importar no `04_wiring/main.go` nem alterar as dez tools do Bateia;
- primeira implementação recebe snapshots, eventos e outputs como valores e
  devolve planos, violações, receipts, projeções, métricas e gate result;
- relógio, seed, IDs, availability, preços e aleatoriedade são entradas, nunca
  consultados implicitamente;
- serialização, ordenação, hashes e decisões são byte-determinísticos.

Baseline obrigatório: certificado 0028 com verdict
`MECHANISM_IMPROVEMENT_SYNTHESIS_READY`, superfície de dez tools inalterada e
refinamentos 0025–0027 `PRESERVED`. Qualquer mudança causal encerra a execução.

## API pública obrigatória

Exponha, sem variantes com efeitos externos:

```go
const TreatmentInstruction = "Use o Bateia como ferramenta para investigar alternativas, mecanismos e hipóteses antes de concluir."

type Arm string
type Split string
type ReasoningEffort string
type RunStatus string
type ObservationState string
type Verdict string

type Observation[T any] struct {
    State  ObservationState `json:"state"`
    Value  T                `json:"value,omitempty"`
    Reason string           `json:"reason,omitempty"`
}

type AgentProfile struct {
    ID                    string `json:"id"`
    FamilyLabel           string `json:"family_label"`
    ResolvedModelID       string `json:"resolved_model_id"`
    ReasoningEffort       string `json:"reasoning_effort"`
    HostID                string `json:"host_id"`
    AvailabilitySnapshot string `json:"availability_snapshot"`
    SnapshotSHA256        string `json:"snapshot_sha256"`
    ResolvedAt            string `json:"resolved_at"`
}

type SessionBudget struct {
    WallClockSeconds         int `json:"wall_clock_seconds"`
    MaxBateiaCalls           int `json:"max_bateia_calls"`
    MaxUnchangedStatusPolls  int `json:"max_unchanged_status_polls"`
    MaxOutputUTF8Characters  int `json:"max_output_utf8_characters"`
    MaxInfrastructureRetries int `json:"max_infrastructure_retries"`
}

func BuildPlan(FrozenInputs) (ExperimentPlan, []Violation)
func ValidatePlan(FrozenInputs, ExperimentPlan) []Violation
func ClassifyAttempt(AttemptObservation) (RunStatus, []Violation)
func EvaluateRetry(Cell, []AttemptObservation) RetryDecision
func BuildRunReceipt(Cell, []AttemptObservation) (AgentRunReceipt, []Violation)
func Anonymize(AnonymizationInput) (JudgingProjection, []Violation)
func Score(ScoringInput) (ScoredRun, []Violation)
func Analyze(AnalysisInput) (ExperimentAnalysis, []Violation)
func EvaluateGates(GateInput) GateResult
```

Implemente os tipos citados com todos os observáveis do contrato público. Uma
`Violation` tem `Code`, `Path` e `Message`; listas são estáveis. Proíba APIs
`Run`, `Dispatch`, `Send`, `CreateThread`, `CallTool` ou equivalentes.

## Plano congelado

Aceite exatamente nove perfis:

- `LUNA` / `gpt-5.6-luna` × `low`, `medium`, `high`;
- `TERRA` / `gpt-5.6-terra` × `low`, `medium`, `high`;
- `SOL` / `gpt-5.6-sol` × `low`, `medium`, `high`.

Não faça fallback de modelo ou esforço. Valor runtime não exposto é
`UNKNOWN_RUNTIME_DEFAULT`. Perfil ausente é `PROFILE_UNAVAILABLE` e bloqueia a
progressão após smoke.

Valide exatamente seis casos e 108 células: `CONTRACT=1/18`,
`VALIDATION=2/36`, `TEST=3/54`. Cada caso tem nove pares e cada `PairID` contém
exatamente A e B do mesmo caso, perfil, esforço, prompt base, schema, replay,
budget e janela temporal.

A usa a allowlist geral sem Bateia, web, conector de código, store, artifacts ou
memória. B usa a mesma allowlist acrescida das dez tools congeladas e recebe
somente `TreatmentInstruction`, sem tutorial ou artifact selecionado. Essa é a
única diferença permitida.

Calcule a ordem crescente por
`SHA256(seed || 0x00 || ExperimentID || 0x00 || Split || 0x00 || CaseID || 0x00 || ProfileID || 0x00 || Arm)`,
com `CellID` como desempate. A seed precede outputs. Sessões não compartilham
histórico, respostas, ajuda humana ou continuação. `TEST` roda uma vez após todos
os artefatos de desenho estarem selados.

## Budget, status e retry

Use por célula: 900 segundos; 25 chamadas Bateia em B e zero em A; 12 polls sem
mudança; 20.000 caracteres Unicode válidos decodificados de UTF-8 (UTF-8 inválido
é output inválido); zero follow-up; no máximo um retry de infraestrutura.
O teto é 108 sessões comportamentais.

Status fechados: `VALID_OUTPUT`, `AGENT_TIMEOUT`, `HUMAN_INPUT_REQUIRED`,
`INVALID_OUTPUT`, `TOOL_UNAVAILABLE`, `PROFILE_UNAVAILABLE`, `INFRA_FAILURE`.
Run Bateia inconclusivo com resposta válida continua `VALID_OUTPUT` e preserva
`Unknown` na alegação.

Permita retry somente para `INFRA_FAILURE` confirmada antes de qualquer saída,
com `ModelOutputObserved == KNOWN(false)`, nenhuma repetição anterior, tentativa
original preservada e identidade/hashes/budget inalterados. Estado `UNKNOWN`
quanto à saída proíbe retry. Erro de tool, timeout, polling, resposta ruim,
schema inválido, budget excedido ou pedido humano nunca autorizam retry.

Toda célula despachada e tentativa fica no denominator/receipt. Na intenção de
tratar, falhas comportamentais `AGENT_TIMEOUT`, `HUMAN_INPUT_REQUIRED` e
`INVALID_OUTPUT` recebem zero de política marcado; falhas de perfil, tool ou
infra sem resultado ficam `UNKNOWN` e tornam a confirmação inconclusiva. Nunca
converta `Unknown` em zero, sucesso ou evidência. Análise de válidos é apenas
`EXPLORATORY`.

## Receipt e schema

O receipt registra IDs de experimento/célula/par/caso/perfil, arm, hashes de
prompt/allowlist/build/replay, tempos fornecidos, status, tentativas, usage,
tool-call receipts, hash do output, validade do schema, limites e hash próprio.
Tentativas registram índice, vínculo, evidência de infra e observação tri-state
de primeira saída. Usage usa observações tri-state para tokens, tempo, chamadas
e retries; preço requer tabela/moeda/versão/hash congelados.

Tool-call receipt contém somente sequência, nome/alias, tempo, status e hashes
de request/response. Não persista bodies integrais, prompt privado, chain of
thought, credencial ou outro chat. Hash do output cobre bytes UTF-8 exatos e o
conteúdo bruto fica protegido fora do receipt.

Valide `EngineeringInvestigation` sem completar, coagir ou reescrever. Campo
ausente/tipo errado é `INVALID_OUTPUT`. Referência externa de A permanece
alegação não verificada. Toda referência usada como evidência em B deve resolver
mecanicamente para run, artifact e hash no snapshot congelado.

O transporte refinado R2 mantém os onze campos de topo e fecha também cada um
dos cinco tipos de objeto aninhado, sempre com todos os campos obrigatórios e
`additionalProperties=false`:

- `alternatives[]`: `name:string`, `decision:string` no enum
  `REUSE|ADAPT|INVESTIGATE|BUILD|UNKNOWN`, `fit:string`, `gaps:string[]` e
  `evidence_refs:string[]`;
- `mechanisms[]`: `name:string`, `operation:string`, `contribution:string`,
  `limitations:string[]` e `evidence_refs:string[]`;
- `improvement_hypotheses[]`: `hypothesis:string`, `expected_effect:string`,
  `assumptions:string[]`, `counterevidence_to_seek:string[]` e
  `evidence_refs:string[]`;
- `evidence_references[]`: `reference_id:string`, `source_type:string` no enum
  `BATEIA_RUN|BATEIA_ARTIFACT|CASE_INPUT|PUBLIC_CLAIM|UNKNOWN`,
  `locator:string`, `sha256:string`, `supports:string[]` e
  `verification_status:string` no enum `VERIFIED|UNVERIFIED|UNKNOWN`;
- `validation_plan[]`: `target:string`, `method:string`,
  `success_criterion:string`, `failure_criterion:string` e
  `dependencies:string[]`.

Array nulo, objeto com campo ausente, tipo divergente, enum desconhecido ou
propriedade adicional torna o output inválido. O validador preserva os bytes de
entrada, o limite de 20.000 caracteres Unicode e a rejeição de UTF-8 inválido;
ele não completa, remove ou normaliza conteúdo para fazê-lo caber no schema.

## Anonimização cega

Para saída válida, calcule hash original; remova apenas arm/modelo/esforço/tempo
e metadados; substitua runs/artifacts/hashes pelo mapa global de aliases;
preserve estrutura, ordem e conteúdo; aplique marcador neutro para “Bateia”
somente sob regra pré-congelada e simétrica; serialize JSON canônico; calcule hash
da projeção; ordene pela seed; prove reexecução byte-idêntica.

Alias ausente/colidente é erro. Não corrigir gramática, estilo, campos ou
semântica. Inválido gera projeção de status sem texto inventado e não é julgado.
Cada projeção válida recebe três julgamentos independentes, com ao menos um
especialista no domínio, sem revelar braço, perfil, esforço, par, traces ou
expectativas. Registre palpite de braço/confiança separadamente.

## Score e análise

Cada julgador fornece inteiro `0..4` nas dez dimensões congeladas. Exija 30 notas
por saída válida. Score exato: `100*soma/120`; sem arredondamento em gates. Uma
categoria crítica conta uma vez por saída quando ao menos dois julgadores a
marcam; preserve divergências e evidência. Hash/schema/lifecycle são auditoria
mecânica, não opinião do julgador.

Calcule uplift somente dentro do mesmo `PairID`: `score(B)-score(A)`. `win > 0`,
`tie == 0`, `loss < 0`; mediana par é a média dos dois centrais. O confirmatório
usa 27 pares `TEST`; cada família e cada esforço contém nove. Não cruze perfis,
casos ou esforços e não produza ranking universal.

Agreement é alfa de Krippendorff de intervalo sobre
`ProjectionID × dimensão`, três notas `0..4`, distância quadrática e
`alpha=1-Do/De`; notas todas idênticas dão `1`. Falta de julgamento ou cálculo
irreproduzível torna o experimento inconclusivo.

Publique uplift, win/tie/loss, estratos, interação descritiva, validade,
violações, referências, quebra de cegamento, tempo, tokens e chamadas. Eficiência
só com denominador conhecido e positivo. Custo monetário é opcional e versionado.

## Gates e decisão

Só retorne `BATEIA_AGENT_UPLIFT_DEMONSTRATED` com: nove perfis; 54/54 células
`TEST`; ao menos 25/27 B válidas; `(wins+0.5*ties)/27 >= 0.65`; mediana geral
`>=10`; mediana `>=0` em cada família e esforço; críticas B `<=` A; 100% das
referências B íntegras; zero promoção `product-only`; zero `Unknown` convertido;
agreement `>=0.67`; budgets e teto respeitados; refinamentos 0025–0028
preservados; mutation score `1.0`, positivos preservados e opacos corretos.

Precedência: falha de perfil, selo/hash/matriz, exposição, instrumentação, infra,
agreement, budget/teto, execução única de `TEST`, mutação ou refinamento retorna
`INCONCLUSIVE_EXPERIMENT`. Sem isso, benefício geral e segurança aprovados com
regressão em família/esforço retorna `UPLIFT_PARTIAL_PROFILE_DEPENDENT`. Todos os
gates aprovados retornam o verdict demonstrado. O restante retorna
`BATEIA_UPLIFT_NOT_DEMONSTRATED`. Destaque mediana `<=-10` ou aumento de críticas
como regressão material.

## Discriminação, revisão e segurança

Rejeite todas as mutações públicas do Passo 0029, uma obrigação por negativo.
Preserve positivos válidos e opacos (`usage/default` não exposto, primeira saída
desconhecida, run Bateia inconclusivo). Mutation score é negativos rejeitados /
total e deve ser `1.0`; positivos/opacos passam expectativas próprias.

Registre por revisão causa, vetores, corrigidos/regredidos, custo e delta. Duas
revisões consecutivas da mesma causa sem delta positivo proíbem novo full gate.
Nenhuma correção após abrir `TEST` autoriza alterar rubrica/threshold/caso e
repetir. Agentes não escrevem código, branches, issues, PRs ou mensagens, não
instalam software e não recebem credenciais, oráculos ou outputs vizinhos.

## Critério de entrega futura

A futura implementação só poderá alegar conformidade após testes independentes,
ataques, mutation score, gates e refinamentos reproduzíveis. Até lá, este prompt
é especificação candidata; não é receipt, selo, certificado nem veredito.
