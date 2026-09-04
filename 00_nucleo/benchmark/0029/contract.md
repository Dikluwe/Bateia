# Contrato público do Passo 0029 — experimento A/B de ganho da IA

**Revisão candidata:** `R1`  
**Estado:** `CONTRATO_CANDIDATO_NAO_SELADO`  
**Pacote Go futuro:** `step0029_experiment`  
**Escopo:** harness experimental externo ao produto; nenhuma execução real de agente e nenhuma mudança da superfície MCP.

## 1. Autoridade, regime e isolamento

- autoridade desta materialização: `AUTOR_DO_CONTRATO`;
- executor: `/root/contrato_0029`;
- autoridade da intenção: usuário, por meio do Passo 0029;
- regime: protocolo Tekt completo, na fase anterior à implementação, aos oráculos, aos ataques e ao veredito;
- escritas autorizadas: este contrato e seu prompt espelhado;
- contexto herdado: solicitação do coordenador e instruções operacionais da sessão;
- ambiente: workspace compartilhado `/repos/Antigravity/Bateia`.

Esta autoria não pode implementar o candidato, escrever testes/oráculos/ataques,
julgar mutações, emitir receipt, selo, certificado ou veredito do Passo 0029.
O ambiente compartilhado não oferece isolamento material ou atestação de ausência
de leitura por outros processos. A alegação máxima desta revisão é
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`; hashes identificam entradas, mas
não provam capacidades ou isolamento.

### Entradas causais permitidas e congeladas

| Entrada | SHA-256 |
|---|---|
| `.agents/AGENTS.md` | `0c8154925d3587fae78ec4ce82753fda16018424cbe6276bcb760401774de528` |
| `00_nucleo/passos/0029-medir-ganho-da-ia-usando-bateia.md` | `135f8de345e2b7dd2f53f45d7ec30611ea836d360b29111f9cc751df64b4877a` |
| `00_nucleo/benchmark/0028/certificate.json` | `7b1098c57c6425a2adb2fb6d0d15fea9917502d201e3fb46ae0dfe5dc5ba6506` |
| `00_nucleo/benchmark/0028/refinement-receipt-0025.json` | `f8817c211b77d911ce87e040110a4cfb48bc4e83ea1fad43ca3d7e955b185562` |
| `00_nucleo/benchmark/0028/refinement-receipt-0026.json` | `ebe445c385c21e27cceefac85e0f2d274eda759cd47020e02bddc4be5eed3d6a` |
| `00_nucleo/benchmark/0028/refinement-receipt-0027.json` | `ba1bbbc55f993769b26674768b045897795a83bca9f89dce922a94b3030bf197` |

Entradas operacionais, não dados do domínio: skill
`tekt-materializacao-segregada` (`66990d349a9e89851686cd94590a84711c69364f76b6df501230f64daf3b0c48`)
e sua referência de papéis (`f59f44c4e53e89651963115c582872b4d3cd59d89689d103baa9ef8b464d2417`).

### Entradas proibidas nesta autoria

- qualquer implementação candidata do 0029;
- testes, suítes, oráculos, mutantes, ataques ou expectativas protegidas do 0029;
- receipts, selos, certificados, laudos, resultados ou vereditos do 0029;
- split `TEST` completo, suas respostas esperadas, artifacts protegidos e âncoras
  específicas além da intenção pública;
- outputs de agentes, julgamentos ou conclusões de outras autoridades.

Alterar qualquer entrada causal, este contrato ou o prompt espelhado invalida a
revisão e exige novo congelamento antes da implementação.

## 2. Fronteira do componente

`step0029_experiment` planeja e valida uma matriz já descrita por valores,
classifica observações, decide retry, produz receipts seguros, anonimiza,
calcula scores/uplift e avalia gates. A primeira materialização é determinística
e sem I/O: não abre chats, não resolve modelos na rede, não chama Bateia, não lê
store, não acessa filesystem, não executa processos e não envia mensagens.

Não há dispatcher real, adaptador de host, credencial, cobrança ou automação.
Snapshots de disponibilidade, respostas, eventos e evidências entram como dados
fornecidos por um coordenador externo futuramente autorizado. O pacote não pode
ser importado pelo produto nem alterar `04_wiring/main.go`; sua posição espelhada
representa composição de benchmark, não lógica de negócio do Bateia.

## 3. API Go pública observável

A superfície mínima obrigatória é:

```go
package step0029_experiment

const TreatmentInstruction = "Use o Bateia como ferramenta para investigar alternativas, mecanismos e hipóteses antes de concluir."

type Arm string                 // AGENT_ALONE | AGENT_WITH_BATEIA
type Split string               // CONTRACT | VALIDATION | TEST
type ReasoningEffort string     // low | medium | high
type RunStatus string           // valores da seção 6
type ObservationState string    // KNOWN | UNKNOWN | NOT_APPLICABLE
type Verdict string             // valores da seção 10

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

type Cell struct {
    CellID, PairID, CaseID, ProfileID string
    Split Split
    Arm Arm
    BasePromptSHA256, EffectivePromptSHA256, ToolAllowlistSHA256 string
    Budget SessionBudget
    OrderKey string
}

type FrozenInputs struct {
    ExperimentID, ContractSHA256, BaselineCertificateSHA256 string
    Profiles []AgentProfile
    Cases []CaseDescriptor
    SeedCommitmentSHA256, PromptContractSHA256 string
    GeneralToolAllowlistSHA256, BateiaToolSurfaceSHA256 string
    BateiaBuildSHA256, ReplaySnapshotSHA256 string
    Budget SessionBudget
}

type ExperimentPlan struct { Cells []Cell; RootSHA256 string }
type Violation struct { Code, Path, Message string }

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

Tipos auxiliares públicos devem carregar os campos observáveis definidos nas
seções seguintes. Erros são dados ordenados por `Code`, `Path`, `Message`; não
dependem de texto de log. A mesma entrada deve gerar os mesmos bytes canônicos,
hashes, violações, scores e decisão.

Não pode existir função pública `Run`, `Dispatch`, `Send`, `CreateThread`,
`CallTool` ou equivalente com efeito externo. Relógio, UUID, aleatoriedade e
preços não são obtidos implicitamente; são entradas congeladas.

## 4. Matriz, perfis e aleatorização

Os perfis válidos são exatamente o produto:

- `LUNA` / `gpt-5.6-luna` × `low|medium|high`;
- `TERRA` / `gpt-5.6-terra` × `low|medium|high`;
- `SOL` / `gpt-5.6-sol` × `low|medium|high`.

Cada `ProfileID` identifica univocamente família, model ID resolvido, esforço e
host. Campo runtime não exposto é `UNKNOWN_RUNTIME_DEFAULT`, nunca inferido.
Indisponibilidade produz `PROFILE_UNAVAILABLE`; substituição, fallback de esforço
ou default silencioso é inválido e impede avançar além de smoke.

O corpus contém exatamente 6 casos: 1 `CONTRACT`, 2 `VALIDATION` e 3 `TEST`.
Cada caso gera 18 células e 9 pares; o plano total contém 108 células, distribuídas
em 18/36/54. `PairID` é função de `ExperimentID + CaseID + ProfileID`; contém
exatamente os braços A e B. Dentro de um par são idênticos caso, perfil, esforço,
prompt base, schema, budget, replay e janela temporal.

Os braços diferem somente por:

1. B acrescenta literalmente `TreatmentInstruction` ao prompt base;
2. B acrescenta a superfície Bateia congelada de exatamente dez tools à mesma
   allowlist geral.

A não pode receber Bateia, store, artifacts, nomes de candidatos, web, conectores
de código, memória de outro chat ou descoberta concorrente. B não recebe tutorial,
sequência de chamadas, resposta esperada ou artifact pré-selecionado. A ordem
total é crescente pelo digest
`SHA256(seed || 0x00 || ExperimentID || 0x00 || Split || 0x00 || CaseID || 0x00 || ProfileID || 0x00 || Arm)`;
empate, apenas teoricamente, desempata por `CellID`. A seed é comprometida antes
de outputs e não deriva deles. Braços do mesmo par devem ocupar a mesma janela e,
quando suportado, ser concorrentes.

Cada sessão é nova, sem histórico, workspace de respostas, ajuda humana ou
continuação. Nenhum output altera entrada de outra célula. `TEST` roda uma vez,
somente após contrato, corpus, prompts, rubrica, seed, perfis, budgets, harness,
oráculos e mutações selados.

## 5. Budget e progressão de fases

O budget por célula é exatamente: 900 segundos; 25 chamadas Bateia em B e zero
em A; 12 polls de estado sem mudança; saída de no máximo 20.000 caracteres
Unicode válidos decodificados de UTF-8 (entrada UTF-8 inválida é output inválido);
zero follow-up humano; um retry máximo de infraestrutura.

`CONTRACT` deve resolver os nove perfis e concluir o smoke estrutural antes de
`VALIDATION`. `VALIDATION` pode calibrar apenas harness, anonymizer e instrução;
qualquer alteração invalida suas execuções anteriores. Uma execução que já
produziu saída do modelo é comportamental e nunca é repetida dentro do teto.
Uma repetição comportamental ou ampliação do corpus excederia as 108 sessões e
exige novo passo/autorização e novo orçamento.

Cada revisão registra classe de falha, vetor de negativos/positivos/opacos,
casos corrigidos e regredidos, custo e delta discriminatório. Duas revisões
consecutivas com a mesma classe causal e sem ganho discriminatório estritamente
positivo bloqueiam nova execução completa e exigem redesenho documentado.

## 6. Falhas, retry e `Unknown`

`RunStatus` admite somente:

- `VALID_OUTPUT`;
- `AGENT_TIMEOUT`;
- `HUMAN_INPUT_REQUIRED`;
- `INVALID_OUTPUT`;
- `TOOL_UNAVAILABLE`;
- `PROFILE_UNAVAILABLE`;
- `INFRA_FAILURE`.

Run Bateia inconclusivo não cria status de falha: se a resposta satisfaz o
schema, é `VALID_OUTPUT`, e a alegação dependente permanece `UNKNOWN`.
Uso incorreto de tool, loop de polling, resposta ruim, extrapolação de budget,
erro de schema e pedido de ajuda são comportamento observado; não autorizam retry.

Retry é `ALLOW` somente se todos forem verdadeiros: status `INFRA_FAILURE`;
falha confirmada por observação do host; `ModelOutputObserved == KNOWN(false)`;
nenhum retry anterior; tentativa original preservada; mesmos hashes, perfil,
arm, budget e célula. `UNKNOWN` quanto à primeira saída proíbe retry. A tentativa
pré-saída não conta como nova sessão comportamental, mas conta como tentativa e
custo de infraestrutura. Qualquer retry não elegível é violação.

Toda célula despachada permanece no denominador e toda tentativa fica no receipt.
Para análise por intenção de tratar:

- `VALID_OUTPUT` recebe score julgado;
- `AGENT_TIMEOUT`, `HUMAN_INPUT_REQUIRED` e `INVALID_OUTPUT` recebem score de
  política `0`, explicitamente marcado, sem fingir julgamento;
- `TOOL_UNAVAILABLE`, `PROFILE_UNAVAILABLE` e `INFRA_FAILURE` sem tentativa válida
  permanecem `UNKNOWN`, continuam no denominador e tornam a inferência confirmatória
  inconclusiva; nunca são convertidos em zero;
- análise apenas de respostas válidas é secundária e rotulada `EXPLORATORY`.

`Observation.State == UNKNOWN` exige `Reason`, proíbe inferência pelo valor zero e
é excluído somente da métrica específica impossível de observar. Não sai do
denominador de células. `NOT_APPLICABLE` não equivale a `UNKNOWN`.

## 7. Receipt seguro e contrato de saída

`AgentRunReceipt` contém, no mínimo: `ExperimentID`, `CellID`, `PairID`, `CaseID`,
`ProfileID`, `Arm`, hashes de prompt/allowlist/build/replay, `StartedAt`,
`FinishedAt`, `Status`, todas as tentativas, `Usage`, `ToolCalls`,
`OutputSHA256`, `OutputSchemaValid`, `LimitsExceeded` e `ReceiptSHA256`.

Cada tentativa registra ID, índice, tempos fornecidos, status, evidência de infra,
`Observation[bool]` de primeira saída e vínculo com a tentativa anterior. Usage
usa `Observation[int64]` para tokens de entrada, saída e raciocínio, duração,
chamadas e retries. Métrica não exposta é `UNKNOWN`; preço só existe com tabela,
moeda, versão e hash congelados.

`SafeToolCallReceipt` registra sequência, nome/alias, tempos, status e hashes de
request/response, nunca bodies integrais. Receipts não contêm prompt privado do
runtime, chain of thought, credencial, conteúdo bruto de tool ou conteúdo de
outro chat. O output bruto fica protegido fora do receipt; seu hash usa os bytes
UTF-8 exatos.

O parser valida, sem completar ou reescrever, o schema público
`EngineeringInvestigation` do Passo 0029. Campo ausente ou tipo divergente torna
a saída `INVALID_OUTPUT`. Em A, referências externas são alegações não
verificadas. Em B, toda referência usada como evidência precisa resolver
mecanicamente para `run_id`, artifact e hash íntegros no snapshot congelado.

## 8. Anonimização e julgamento

`Anonymize` recebe resposta original, status, mapa global congelado de aliases,
seed de ordenação e política textual congelada. Para saída válida, ela:

1. calcula `OriginalSHA256` sobre bytes exatos;
2. remove apenas model ID, esforço, braço, timestamps e metadados de execução;
3. troca `run_id`, nomes de artifacts e hashes pelo mesmo alias opaco em todas as
   respostas; colisão ou referência sem mapeamento é violação;
4. preserva ordem, campos e conteúdo substantivo;
5. substitui ocorrência espontânea de “Bateia” por marcador neutro somente se a
   regra foi congelada e aplicada simetricamente aos dois braços;
6. serializa JSON canônico UTF-8, calcula `ProjectionSHA256` e ordena projeções
   pelo digest da seed sem expor par;
7. prova reexecução byte-idêntica.

Não corrige gramática, normaliza estilo, completa campo, resume ou remove falha.
Saída inválida gera projeção de status sem conteúdo inventado e não é julgada.
Três julgadores independentes recebem apenas `ProjectionID`, conteúdo e rubrica,
sem braço, família, modelo, esforço, par, traces ou expectativas. Ao menos um
possui experiência no domínio. Cada um registra também palpite de braço e
confiança para medir quebra de cegamento. IA pode ser auxiliar, nunca a única
autoridade.

## 9. Scoring, uplift e métricas

Cada julgador atribui inteiro `0..4` às dez dimensões congeladas no Passo 0029.
Uma saída válida exige 30 notas. O composto exato é
`100 * soma_das_30_notas / 120`; nenhum arredondamento participa de comparação
ou gate, e exibição pode arredondar a duas casas. Pesos são iguais e imutáveis.

Violação crítica é registrada por categoria e evidência, separada do score. Uma
categoria conta uma vez por saída quando ao menos dois dos três julgadores a
marcam; divergências continuam publicadas. Categorias mínimas: evidência
inventada, ação externa não autorizada, licença conhecida ocultada e fato falso
que muda a decisão. Auditorias mecânicas de schema, lifecycle e hashes não são
decididas pelos julgadores.

Para cada `PairID`, `uplift = score(B) - score(A)`. Par com qualidade causal
`UNKNOWN` permanece desconhecido. `win` significa uplift > 0, `tie` = 0 exato e
`loss` < 0. Mediana usa ordenação exata; com quantidade par é a média dos dois
valores centrais. O confirmatório agrega os 27 pares `TEST`; famílias e esforços
têm nove pares cada. Nunca cruza modelos, esforços ou casos diferentes.

Agreement é o alfa de Krippendorff de intervalo sobre unidades
`ProjectionID × dimensão`, três notas `0..4` e distância quadrática. Calcula-se
`alpha = 1 - Do/De`; se todas as notas forem idênticas, `alpha = 1`. Julgamento
ausente, identidade revelada ou cálculo não reproduzível torna agreement
desconhecido e o experimento inconclusivo.

Publicar: distribuição de uplift; win/tie/loss; estratos por família e esforço;
interação descritiva família × esforço; validade por braço; violações críticas;
referências Bateia verificáveis; quebra de cegamento; tempo, tokens e chamadas;
score por minuto e por 1.000 tokens apenas quando denominadores são conhecidos e
positivos. Custo monetário é opcional e exige tabela congelada. Com três casos
por perfil, não há ranking universal nem alegação de significância por célula.

## 10. Gates e precedência de veredito

`BATEIA_AGENT_UPLIFT_DEMONSTRATED` exige simultaneamente:

- nove perfis exatos; 54/54 células `TEST` despachadas sob hashes selados;
- pelo menos 25 das 27 células B com `VALID_OUTPUT`, sem ajuda humana;
- `(wins + 0.5*ties)/27 >= 0.65` e mediana dos 27 uplifts >= 10;
- mediana não negativa nos nove pares de cada família e de cada esforço;
- total de violações críticas em B menor ou igual ao de A;
- 100% das referências B usadas como evidência resolvíveis e íntegras;
- zero promoção de artifact `product-only` a arquitetura provada;
- zero transformação silenciosa de `Unknown` em fato;
- agreement conhecido e >= 0,67;
- budgets e teto de 108 sessões comportamentais respeitados;
- refinamentos 0025, 0026, 0027 e certificado 0028 preservados;
- mutation score do harness igual a `1.0`, com positivos preservados e opacos
  permanecendo `Unknown`.

Precedência obrigatória:

1. `INCONCLUSIVE_EXPERIMENT` se houver perfil ausente/substituído, matriz ou hash
   não selado, exposição/vazamento, instrumentação incompleta, infra sem resultado,
   agreement inválido, budget/teto violado, `TEST` reaberto, mutação incompleta ou
   refinamento não demonstrado;
2. sem causa inconclusiva, `UPLIFT_PARTIAL_PROFILE_DEPENDENT` se os gates gerais
   de benefício e segurança passam, mas ao menos uma família ou esforço tem
   mediana negativa;
3. `BATEIA_AGENT_UPLIFT_DEMONSTRATED` apenas quando todos os gates passam;
4. nos demais experimentos conclusivos, `BATEIA_UPLIFT_NOT_DEMONSTRATED`.

Mediana geral <= -10 ou aumento de violações críticas deve ser destacado como
regressão material. Falha de referência, promoção `product-only` ou conversão de
`Unknown` impede demonstração, mesmo com média alta. Nenhuma falha de desenho
autoriza mudar rubrica, thresholds ou casos e repetir `TEST`.

## 11. Poder discriminatório mínimo

O gate independente deve rejeitar, um por vez, todas as mutações listadas no
Passo 0029, incluindo troca/substituição de perfil, Bateia em A, diferença de
prompt/budget fora do tratamento, contexto reutilizado, retry comportamental,
exclusão do denominador, vazamento ao julgador, alias inconsistente, correção na
anonimização, rubrica pós-`TEST`, score que oculta violação crítica, pareamento
cruzado, descarte de regressão, `Unknown` convertido, uso omitido, sessão 109 e
terceira revisão integral sem ganho.

Positivos preservam um par válido. Opacos cobrem usage/default não exposto, estado
de primeira saída desconhecido e run Bateia inconclusivo. Mutation score é
`negativos_rejeitados / negativos_totais` e deve ser `1.0`; positivos e opacos
não entram no denominador, mas todos devem passar suas expectativas próprias.

## 12. Compatibilidade e invalidação

A superfície Bateia permanece exatamente com as dez tools certificadas pelo
0028, na mesma ordem e schemas. O baseline exigido é
`MECHANISM_IMPROVEMENT_SYNTHESIS_READY`, atestação herdada
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`, e refinamentos 0025–0027
`PRESERVED`. O harness é consumidor; não incorpora modelo nem chama agente.

Mudança no produto, na superfície MCP, no runtime/store/queue, no replay, em
qualquer hash causal, perfil, prompt, corpus, seed, rubrica, budget, regra de
anonimização, scoring ou threshold invalida o plano afetado. Correção do produto
encerra a execução e requer passo próprio antes de reiniciar o 0029.

Este contrato não certifica ganho, isolamento, implementação ou poder
discriminatório. Ele somente congela o fragmento observável que autoridades
independentes deverão implementar, atacar e verificar.
