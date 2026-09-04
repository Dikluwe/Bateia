# Laudo do Passo 0029 — harness local do experimento de ganho

**Regime:** protocolo Tekt completo  
**Executor do veredito:** `/root/verificador_0029`  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Veredito local:** `HARNESS_LOCAL_READY_WITH_ARCHITECTURAL_WARNINGS_EXTERNAL_EXPERIMENT_PENDING`  
**Veredito de produto:** `NOT_EVALUATED`

## Conclusão

O harness determinístico local do Passo 0029 está executável e discriminatório
para o fragmento congelado: as duas execuções focais passaram, o race passou,
as 25 mutações negativas foram rejeitadas (`mutation_score = 25/25 = 1,0`), os
nove positivos foram preservados e os cinco opacos permaneceram `Unknown`.
`go vet`, `go build`, JSON estrito, hashes, cadeia R1–R5 e refinamentos 0025–0028
também passaram no escopo descrito nos recibos.

Isso não demonstra ganho do Bateia. Nenhuma sessão externa/paga foi autorizada:
o smoke `CONTRACT` permanece `0/18` e a execução comportamental `0/108`, ambos
`BLOCKED_BY_EXPLICIT_BUDGET_AUTHORIZATION`. Portanto não há outputs, julgamentos,
custos observados ou análise causal e este laudo não declara
`BATEIA_AGENT_UPLIFT_DEMONSTRATED`.

## Identidades congeladas

- manifest ativo R5: `639f7e5c7d1e4431bbeb9082edbd71caa85f8742bbfd0ec5dd943c77e8264456`;
- input seal R5: `bae46a813d04e52cc71fd08be0ab82005ed307dceff3e1b03c88b6d4353679e3`;
- selo local pós-gate: `b45259ad04e6d0deb3c3c01724c856931377b434434ab0c0e3d60131b6ec470d`;
- contrato: `48f0f3ff313caf2cd1c17f0dacb3a0143bf7908928e6922d3838d25e4f1118e7`;
- baseline 0028: `7b1098c57c6425a2adb2fb6d0d15fea9917502d201e3fb46ae0dfe5dc5ba6506`;
- candidato local: `7c10330c11f5c98dcf3f0ade780ae0f42cddb8e37879ec592ce4bf6c6f029dea`;
- suíte oráculo: `d377ddb3c8ad4aa008e05938f3893595588d274074f02e9ed086ceda8564fb7a`;
- suíte adversarial: `ee1ee904c71d7e9f7c9ea73a1b86792802fce42f976f56dccf1ebd9cf40bcf06`;
- plano `BuildPlan`: `5215fc6764efea3f5b1a079c854c5a0b0837546f44aed4b9a1093eb13b4f2bf3`.

O plano foi reproduzido pelo verificador com 108 células, 54 pares, splits
18/36/54, braços 54/54 e zero violações de `ValidatePlan`. O campo de replay
liga o receipt certificado do Passo 0028 (`983246…bdc1f`), não um novo capture
do 0029.

## Comandos executados pelo verificador

- `go test -count=1 -v ./04_wiring/step0029_experiment` — PASS, duas vezes;
- `go test -race -count=1 ./04_wiring/step0029_experiment` — PASS;
- `go vet ./...` — PASS;
- `go build ./...` — PASS;
- `go test -tags step0028_candidate -count=1 ./01_core/domain ./01_core/usecase ./02_shell/mcp ./04_wiring` — PASS;
- parse dos 28 JSONs 0029 — PASS;
- validação da cadeia de manifests/invalidações/selos R1–R5 — PASS, zero mismatch;
- reprodução temporária independente de `BuildPlan` — PASS; o arquivo temporário foi removido;
- `crystalline-lint .` — exit 0 com warnings;
- `git diff --check` e `git diff --cached --check` — PASS.

Uma tentativa global dentro do sandbox foi impedida por loopback de `httptest`.
Na reexecução autorizada, um gate legado 0026 R5 apresentou sobrevivência
transitória de M28; a repetição isolada imediata passou. O coordenador reportou
uma execução global completa posterior como PASS. O recibo distingue essa
evidência recebida das provas executadas diretamente por este verificador.

## Refinamentos

Os hashes certificados de 0025, 0026 e 0027 permanecem, respectivamente,
`f8817c…5562`, `ebe445…d6a` e `ba1bbb…f197`. O certificado 0028 preserva
`MECHANISM_IMPROVEMENT_SYNTHESIS_READY` sob `7b1098…6506`. O 0029 é aditivo,
não acrescenta tool MCP e não altera runtime, store, queue, connectors ou
produção; as suítes candidatas 0028 passaram. Resultado: 0025–0028
`PRESERVED` no fragmento observado.

## Limitações e aviso arquitetural

O workspace e o contexto são compartilhados; hashes provam identidade e ordem,
não isolamento forte. O linteiro retornou sucesso, porém reportou warnings V12
para lógica no diretório L4, incluindo o harness 0029. Assim, a prontidão local
é certificada com avisos, sem alegar conformidade arquitetural limpa nem
equivalência funcional geral.

Próximo gate, condicionado a autorização e orçamento explícitos: executar as
18 células `CONTRACT` em sessões limpas. Apenas depois de smoke, validação, TEST,
julgamento cego e análise sob os mesmos hashes poderá existir um veredito de
produto.
