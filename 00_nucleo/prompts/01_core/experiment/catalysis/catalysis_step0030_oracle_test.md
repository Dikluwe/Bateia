# Oráculo independente — Passo 0030 / catálise da busca

## Papel e regime

Este artefato foi escrito pelo testador A/B `/root/testador_0030`, sem leitura
do candidato em `01_core/experiment/catalysis/` e sem leitura da implementação
do Passo 0029. O regime aplicado é A/B: a independência material protegida é
entre o código candidato e este oráculo black-box. A tarefa restringiu a
leitura aos contratos enumerados e às fixtures R3 referenciadas; por isso não
foi procurada uma ADR adicional e não se alega atestação do protocolo Tekt
completo.

Nenhuma sessão comportamental, modelo, web, MCP ou conector foi executado.

## Entradas congeladas e hashes observados

| Entrada | SHA-256 |
|---|---|
| `00_nucleo/passos/0030-corrigir-teste-de-catalise-da-busca.md` | `2a7a5bd76847c6dcf8d40f05c7acbc6d502fb3c3e96e38be331d04ba3c60f03f` |
| `00_nucleo/benchmark/0030/contract-seal-r2.json` | `01147b3cf441d90c895f98eacfae95e2feb560d1541a0b256a05afd08e27fcc4` |
| `00_nucleo/benchmark/0030/web-event-decoder-contract.json` | `8a11f29cc3443009e06ce1b359c97ad9ccb257c3ac3f224c03169689ec1c9206` |
| `00_nucleo/benchmark/0030/pair-dispatch-contract.json` | `b4c742d32ac0219013be0c48cef7f4f4adee6a13ddacdcde4b74283c916787b1` |
| `00_nucleo/benchmark/0030/positive-opaque-fixtures.json` | `0cd9d62182928650615ff073c1e1fdbc7bd6854e502958fb35743d9d985ffa41` |
| `00_nucleo/benchmark/0030/adversarial-mutations.json` | `095e21fda17bf368459965708b93287f8c386195934b2218aea387f6ac40ef70` |
| `00_nucleo/benchmark/0030/output-bound-proof.json` | `dbc6a5766ece9ce767af9fcd60794cf5df7c4e2bbd881e533539356252ee069b` |
| `00_nucleo/benchmark/0030/randomization-plan-r2.json` | `89c1de4c946962abc662c3e3ca80b41150677e9c3ddbbe48dd662bb9ea6cfb9c` |
| `00_nucleo/benchmark/0029/catalysis-pilot-r3-results.json` | `3239ac23443e451c63027c93867aacfa55ec07e22c87709af97c8666ac2f8524` |
| `00_nucleo/benchmark/0029/catalysis-pilot-r3-final-verification.json` | `7fd942124b0f3045e19b3c4e377b7e7c87371fd283aad4c90df06130888d82e2` |
| `00_nucleo/benchmark/0029/sessions-catalysis-r3/treatment/events.jsonl` | `f1e61bf6d04de5c13426f8b09652d471f35f281d9b9c34e3dacd8e3383410a73` |
| `00_nucleo/benchmark/0029/sessions-catalysis-r3/treatment/output.json` | `42c36b20e94915efca624febd1eb58a03253f48d3954bf1c39b757def096b9a8` |

O hash observado do próprio `contract-seal-r2.json` é deliberadamente
registrado como observação de entrada; ele não aparece como hash autorreferente
no selo. Os demais hashes protegidos conferem com os valores declarados pelo
selo R2.

## API pública assumida

Import path: `bateia/01_core/experiment/catalysis`, pacote externo de teste
`catalysis_test`.

Funções:

```go
DecodeWebTrace(raw []byte) WebTrace
AssessOutput(raw, lastMessage []byte, schemaValid bool) OutputAssessment
BuildCommittedQueue(preflightPass bool) []DispatchCell
ValidateArmTrace(arm Arm, events []ToolEvent) ArmTraceAssessment
PreserveInconclusive() AnswerPolicy
ClassifyCatalysisPath(path CatalysisPath) CatalysisClass
```

Shapes públicos congelados pelo integrador:

```text
WebTrace{Records, ThreadID, Machines, PhysicalCallCount,
         NormalizedCallCount, OpaqueEvents, Violations}
WebMachine{StartSequence, FinishSequence, IDs, Encoding,
           NormalizedWebCallID, NormalizedEventID, Query,
           RawStartSHA256, RawCompletionSHA256}
OutputAssessment{StructuralValid, Editorial, UnicodeCharacters,
                 RawBytes, Violations}
DispatchCell{Sequence, CaseID, Arm, Required}
ToolEvent{Tool, Status, Ordinal}
ArmTraceAssessment{State, WebCalls, RecommendCalls, Violations}
AnswerPolicy{State, AutomaticWinner, AutomaticScore}
CatalysisPath{Direction, LaterWebQuery, WebLocator, OutputClaim,
              Used, Confirmed, Rejected, Narrowed}
```

Constantes: `ArmInternet`, `ArmInternetBateia`,
`CatalysisConfirmed`, `CatalysisRejected`, `CatalysisNarrowed`,
`CatalysisUsedWithoutSupport`, `CatalysisNotUsed` e `CatalysisUnknown`.
As observações genéricas devem expor campos exportados `State` e `Value`; o
oráculo exige `Known` ou `Unknown` pela representação textual de `State`.
`Violation` deve expor o código normativo no campo exportado `Code`.

## Obrigações discriminadas

- A fixture real R3 conserva 12 registros, três máquinas físicas e os pares
  duais de IDs na ordem congelada, com a segunda ocorrência `exec-*`
  normalizada.
- Um encoding `call-*` permanece `OPAQUE`, e tanto identidade quanto contagem
  normalizada ficam `Unknown`.
- `started` e `completed` com vetores divergentes não formam máquina e emitem
  `WEB_LIFECYCLE_ID_VECTOR_MISMATCH`.
- O output R3 permanece estruturalmente válido apesar de seus 22.578
  caracteres e recebe apenas `EDITORIAL_TARGET_EXCEEDED`; um JSON original de
  exatamente 65.537 bytes recebe `HARD_OUTPUT_CAP_EXCEEDED`.
- Preflight falho não compromete célula alguma. Preflight aprovado compromete
  exatamente as seis células na ordem selada; invalidade da primeira resposta
  não remove nem reordena o segundo braço de seu par.
- Controle aceita somente web; tratamento exige exatamente uma recomendação
  concluída antes de 1..3 chamadas web. Vazamento no controle, recomendação
  tardia e quatro chamadas web são violações.
- `INCONCLUSIVE` tem estado preservado e nenhum vencedor ou score automático.
- As seis classes de catálise ficam disponíveis, mas um caminho com consulta
  posterior sem locator/evidence causal permanece `CatalysisUnknown`.

## Limite da alegação

Estes testes discriminam apenas a superfície pública acima e o corpus
congelado do Passo 0030. Eles não executam nem certificam runtime, preflight
global, sessões comportamentais, julgamento cego, custo ou ganho geral do
Bateia.
