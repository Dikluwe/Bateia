# Passo 0008 — Comparar candidatos e decidir reuso

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0007 — Diagnosticar e ampliar cobertura de evidências`](0007-diagnosticar-e-ampliar-cobertura-de-evidencias.md)
**Onde roda**: Bateia
**Camadas previstas**: L0, L1 e L2
**Natureza**: matriz comparativa e decisão determinística de ponto de partida

## Condição de início

Este passo só pode ser implementado depois que o Passo 0007 demonstrar que o
Passo 0006 atingiu:

- `evidence_precision` ≥ 90% com ao menos dez `PROVEN` revisados;
- `false_absence_rate` = 0;
- `source_trace_rate` = 100%;
- `component_precision` = 100% na primeira amostra.

Enquanto o gate não for satisfeito, este documento é especificação, não autorização
para transformar assessments ainda não calibrados em recomendação de produto.

## Intenção

Transformar uma `IdeaSpec` e seus `CandidateAssessment` calibrados numa decisão de
partida explicável:

```text
FORK · EXTEND · COMPOSE · INTEGRATE · REFERENCE · BUILD
```

O resultado precisa dizer o que reutilizar, o que falta, quais conflitos impedem
uma opção e por que começar do zero é ou não justificável.

## Restrição central: decisão sem LLM

O motor usa somente regras versionadas, estados de evidência, cobertura, gates e
pesos transparentes. LLM, embeddings, prompts e geração probabilística permanecem
fora do pipeline. Falta de evidência produz `INCONCLUSIVE`, não uma recomendação
inventada.

## Entrada

```go
type ReuseDecisionRequest struct {
    Idea        IdeaSpec              `json:"idea"`
    Assessments []CandidateAssessment `json:"assessments"`
    Policy      ReuseDecisionPolicy   `json:"policy"`
}
```

No primeiro corte, todas as capacidades declaradas em `IdeaSpec.Capabilities` são
obrigatórias. Uma evolução posterior pode adicionar prioridade explícita sem
inferi-la pela ordem ou linguagem natural.

## Matriz comparativa

Produzir linhas por capacidade e colunas por candidato:

```text
capacidade × PROVEN | PARTIAL | ABSENT | INCONCLUSIVE | CONFLICT
```

Manter fora do score funcional:

- licença;
- manutenção;
- operação e privacidade;
- stack e runtimes;
- custo observável de integração;
- non-goals e conflitos eliminatórios.

Um total alto nunca pode esconder um conflito obrigatório.

## Contrato de decisão

```go
type ReuseDecision string

const (
    DecisionFork      ReuseDecision = "FORK"
    DecisionExtend    ReuseDecision = "EXTEND"
    DecisionCompose   ReuseDecision = "COMPOSE"
    DecisionIntegrate ReuseDecision = "INTEGRATE"
    DecisionReference ReuseDecision = "REFERENCE"
    DecisionBuild     ReuseDecision = "BUILD"
    DecisionUnknown   ReuseDecision = "INCONCLUSIVE"
)

type DecisionReason struct {
    RuleID     string           `json:"rule_id"`
    Effect     string           `json:"effect"`
    Candidate string           `json:"candidate,omitempty"`
    Evidence  []EvidenceSource `json:"evidence,omitempty"`
    Message   string           `json:"message"`
}

type ReuseDecisionReport struct {
    SchemaVersion string                `json:"schema_version"`
    Decision      ReuseDecision         `json:"decision"`
    Primary       *ResolvedRepository   `json:"primary,omitempty"`
    Components    []ReusableComponent   `json:"components,omitempty"`
    Matrix        []CapabilityComparison `json:"matrix"`
    Gaps          []string              `json:"gaps"`
    Blockers      []DecisionReason      `json:"blockers,omitempty"`
    Reasons       []DecisionReason      `json:"reasons"`
    Alternatives  []DecisionAlternative `json:"alternatives"`
}
```

## Regras iniciais

As regras são avaliadas na ordem abaixo e registram seus IDs.

### `FORK`

- um candidato cobre todas as capacidades obrigatórias;
- não possui conflito eliminatório;
- entrega projeto executável, não somente biblioteca;
- licença e modo de operação passam pela política;
- lacunas restantes são somente adaptação de produto.

### `EXTEND`

- um candidato cobre a capacidade central e a maioria das obrigatórias;
- possui núcleo ou arquitetura concretamente reutilizável;
- lacunas podem ser adicionadas sem substituir o núcleo;
- não há conflito eliminatório.

### `COMPOSE`

- nenhum candidato isolado cobre a ideia;
- dois ou mais componentes concretos cobrem capacidades complementares;
- interfaces, licenças e operação são compatíveis;
- a composição cobre mais que qualquer candidato isolado.

### `INTEGRATE`

- uma biblioteca, módulo ou binário resolve capacidade delimitada importante;
- o produto principal ainda precisa existir;
- componente possui símbolo, path, licença e caveats verificáveis.

### `REFERENCE`

- há solução conceitualmente útil, mas licença, manutenção, operação, stack ou
  ausência de componente impedem reutilização direta;
- evidências ainda justificam estudar arquitetura, algoritmo ou padrão.

### `BUILD`

- coleta e avaliação são conclusivas;
- nenhum candidato ou composição cobre a capacidade central de forma reutilizável;
- não existe componente compatível suficiente;
- o relatório registra quais buscas e evidências sustentam começar do zero.

### `INCONCLUSIVE`

- capacidade obrigatória possui coleta incompleta;
- licença ou conflito eliminatório permanece ambíguo;
- o corpus/política não sustenta uma das decisões anteriores.

`BUILD` nunca é fallback para falha operacional; falha produz `INCONCLUSIVE`.

## Política versionada

```go
type ReuseDecisionPolicy struct {
    SchemaVersion       string   `json:"schema_version"`
    AllowedLicenses     []string `json:"allowed_licenses"`
    RequireLocal        bool     `json:"require_local"`
    AllowArchived       bool     `json:"allow_archived"`
    MaxComponents       int      `json:"max_components"`
    ExtendMinCoverage   float64  `json:"extend_min_coverage"`
    IntegrationMinValue float64  `json:"integration_min_value"`
}
```

A política e seu hash aparecem no relatório. Alterar política não modifica os
assessments e permite replay offline.

## Custo observável

Não estimar horas fictícias. Classificar esforço em `LOW`, `MEDIUM`, `HIGH` ou
`INCONCLUSIVE` a partir de fatores visíveis:

- mesma linguagem/runtime;
- quantidade de componentes;
- serviço ou banco obrigatório;
- API pública versus símbolos internos;
- projeto arquivado;
- licença e distribuição;
- lacunas funcionais restantes.

Cada fator aparece separadamente na explicação.

## Boca de inspeção

```bash
bateia idea decide \
  --idea ideia.json \
  --assessment assessment.json \
  --policy reuse-policy.json \
  --output decision.json
```

O comando é completamente offline e determinístico.

## Testes obrigatórios

1. conflito eliminatório impede `FORK` apesar de cobertura total;
2. falha ou `INCONCLUSIVE` nunca vira `BUILD`;
3. candidato completo e compatível produz `FORK`;
4. núcleo reutilizável com lacunas adicionáveis produz `EXTEND`;
5. componentes complementares e compatíveis produzem `COMPOSE`;
6. biblioteca delimitada produz `INTEGRATE`;
7. licença incompatível preserva valor como `REFERENCE`;
8. ausência conclusiva de base reutilizável produz `BUILD`;
9. ordem dos candidatos não altera decisão;
10. score da descoberta não altera decisão;
11. toda razão aponta para regra e evidência;
12. política alterada muda somente o resultado esperado;
13. entrada não é modificada;
14. comando funciona com ambiente vazio e sem rede.

## Critério de produto

Executar o fluxo completo sobre pelo menos cinco ideias distintas. Em revisão
manual, a decisão deve ser considerada útil e sustentada pelas evidências em pelo
menos 80% dos casos. Divergências geram revisão de regra versionada, nunca ajuste
silencioso por candidato.

## Fora de escopo

- clonar, forkar ou integrar automaticamente um projeto;
- aconselhamento jurídico;
- estimativa de prazo em horas;
- executar código de terceiros;
- UI/TUI;
- LLM ou embeddings.

## Evidência de implementação

Implementado em 2026-08-26 com contratos L1, motor determinístico, boca de
inspeção L2 e prompts L0 próprios. O replay versionado cobre cinco ideias
distintas e obteve `5/5` decisões alinhadas aos julgamentos congelados. Os testes
adversariais exercitam as 14 obrigações deste passo; o laudo reproduz comandos,
resultados e limitações.

Ver: [`Laudo 0008`](../laudos/0008-comparar-candidatos-e-decidir-reuso.md).
