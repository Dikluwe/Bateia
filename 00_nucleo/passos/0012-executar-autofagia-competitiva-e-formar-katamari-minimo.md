# Passo 0012 — Executar autofagia competitiva e formar um Katamari mínimo

**Estado**: `EM VALIDAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0011 — Congelar benchmark ponta a ponta e auditar utilidade`](0011-congelar-benchmark-ponta-a-ponta-e-auditar-utilidade.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2 e L4  
**Natureza**: autofagia, análise competitiva, cobertura de capacidades e composição mínima sem LLM

## Evidência que motiva este passo

O Passo 0011 materializou o mecanismo offline do benchmark, mas interrompeu a
promoção de produto porque ainda não existem quinze casos correspondentes nem
oráculos humanos isolados. Ao mesmo tempo, a análise competitiva identificou uma
família real de produtos que cobre partes da intenção do Bateia:

- Replen: capacidades de projetos, monitoramento de OSS e vereditos de adoção;
- OSS Compass: recomendação funcional e avaliação de saúde, licença e segurança;
- GitMVP: ideia em linguagem natural para shortlist de repositórios;
- ARIS: decisão de adoção apoiada por dimensões e confiança;
- Starboard: projetos semelhantes com evidências visíveis;
- DRY.codes: descoberta de código e padrões já existentes;
- GitSurfer: composição de projetos e infraestrutura a partir de uma ideia.

Esses nomes são sementes de descoberta, não oráculos de qualidade. Precisamos
executar o próprio método do Bateia sobre essa categoria e observar duas coisas:

1. se ele reconhece corretamente o que já possui — **autofagia**;
2. se consegue propor a menor composição externa que cubra lacunas reais —
   **Katamari**.

## Hipótese

O Bateia será mais útil se responder:

> Qual é a menor composição de capacidades que já existem, dentro e fora do meu
> projeto, necessária para materializar esta ideia sem começar do zero?

Uma recomendação só é válida quando distingue capacidade interna, capacidade
externa comprovada, lacuna real e custo de absorção. Sem essa distinção, Katamari
é apenas acumulação e autofagia é apenas busca textual no próprio repositório.

## Definições operacionais

### Autofagia

Autofagia é aplicar a `IdeaSpec` do alvo ao próprio projeto ou portfólio local e
produzir um inventário comprovado:

```text
ideia/capacidade desejada
    → evidências internas
    → capacidade existente | parcial | ausente | opaca
    → redundância, consolidação ou lacuna
```

Ela não significa apagar ou reescrever código automaticamente. Neste passo,
autofagia é somente diagnóstico e classificação.

### Katamari

Katamari é selecionar uma composição mínima de capacidades externas verificadas
para cobrir apenas as lacunas deixadas pela autofagia:

```text
lacunas internas
    → capacidades externas comprovadas
    → compatibilidade e gates eliminatórios
    → menor conjunto suficiente de fontes
```

“Mínimo” considera número de projetos, capacidades duplicadas, conflitos de
licença, dependências, esforço de integração e risco operacional. Um projeto não
entra na composição apenas porque é popular ou lexicalmente semelhante.

## Restrições centrais

O pipeline principal permanece sem LLM. Este passo não pode:

- usar LLM, embeddings ou resumos gerados como evidência ou decisão;
- alterar regras dos Passos 0005–0011 para favorecer concorrentes conhecidos;
- copiar código de concorrente;
- inferir compatibilidade jurídica apenas do nome da licença;
- absorver capacidade já comprovada internamente;
- esconder `INCONCLUSIVE` para completar artificialmente a composição;
- usar estrelas, popularidade ou atividade como prova funcional;
- misturar documentação de um projeto com código de outro;
- promover uma sugestão sem `repo@commit`, paths e hashes reproduzíveis;
- executar código dos candidatos.

## Entrada congelada

```go
type CompetitiveAutophagyRequest struct {
    SchemaVersion string              `json:"schema_version"`
    Idea          IdeaSpec            `json:"idea"`
    Target        RepositoryIdentity  `json:"target"`
    Internal      []EvidenceBundle    `json:"internal"`
    Competitors   []EvidenceBundle    `json:"competitors"`
    Rules         CapabilityRuleCatalog `json:"rules"`
    ReusePolicy   ReuseDecisionPolicy `json:"reuse_policy"`
    KatamariPolicy KatamariPolicy     `json:"katamari_policy"`
}
```

Os concorrentes descobertos ao vivo precisam ser capturados como bundles
imutáveis antes da análise. O replay e a decisão são offline.

## Contrato de autofagia

```go
type InternalCapabilityDisposition string

const (
    InternalKeep        InternalCapabilityDisposition = "KEEP"
    InternalConsolidate InternalCapabilityDisposition = "CONSOLIDATE"
    InternalGap         InternalCapabilityDisposition = "GAP"
    InternalUnknown     InternalCapabilityDisposition = "UNKNOWN"
)

type InternalCapabilityFinding struct {
    Capability  string                        `json:"capability"`
    Status      EvidenceStatus                `json:"status"`
    Disposition InternalCapabilityDisposition `json:"disposition"`
    Sources     []EvidenceSource              `json:"sources"`
    RuleIDs     []string                      `json:"rule_ids"`
    Rationale   string                        `json:"rationale"`
}

type AutophagyReport struct {
    SchemaVersion string                      `json:"schema_version"`
    Target        RepositoryIdentity          `json:"target"`
    Findings      []InternalCapabilityFinding `json:"findings"`
    Existing      []string                    `json:"existing"`
    Partial       []string                    `json:"partial"`
    Gaps          []string                    `json:"gaps"`
    Unknown       []string                    `json:"unknown"`
    InputHash     string                      `json:"input_hash"`
}
```

Regras mínimas:

- `KEEP`: capacidade comprovada e não duplicada internamente;
- `CONSOLIDATE`: duas ou mais implementações internas comprovadas e sobrepostas;
- `GAP`: ausência sustentada por coleta completa e regra aplicável;
- `UNKNOWN`: coleta incompleta, regra insuficiente ou identidade opaca.

`UNKNOWN` nunca vira lacuna automaticamente, pois isso faria o Katamari absorver
dependências apenas por falta de informação interna.

## Contrato do Katamari

```go
type KatamariPolicy struct {
    SchemaVersion       string   `json:"schema_version"`
    MaxProjects         int      `json:"max_projects"`
    MaxComponents       int      `json:"max_components"`
    AllowedLicenses     []string `json:"allowed_licenses"`
    RequiredOperation   string   `json:"required_operation,omitempty"`
    RequireImmutableRef bool     `json:"require_immutable_ref"`
    AllowUnknown        bool     `json:"allow_unknown"`
}

type KatamariComponent struct {
    Repository   RepositoryIdentity `json:"repository"`
    Capabilities []string           `json:"capabilities"`
    Sources      []EvidenceSource   `json:"sources"`
    Decision     ReuseDecision      `json:"decision"`
    Cost         int                `json:"cost"`
    Risks        []string           `json:"risks"`
}

type KatamariProposal struct {
    SchemaVersion  string              `json:"schema_version"`
    Components     []KatamariComponent `json:"components"`
    CoveredGaps    []string            `json:"covered_gaps"`
    UncoveredGaps  []string            `json:"uncovered_gaps"`
    Rejected       []DecisionAlternative `json:"rejected"`
    Minimal        bool                `json:"minimal"`
    Feasible       bool                `json:"feasible"`
    InputHash      string              `json:"input_hash"`
}
```

O custo é uma ordenação determinística, não uma estimativa monetária. A primeira
versão soma penalidades versionadas por novo projeto, componente adicional,
adaptação de stack, operação incompatível, manutenção problemática e capacidade
duplicada. Gates de licença, identidade e evidência são eliminatórios e não podem
ser compensados por custo baixo.

## Problema de cobertura

O seletor recebe somente lacunas `GAP` da autofagia e uma matriz de capacidades
externas `PROVEN`. Ele deve encontrar, por enumeração determinística limitada, o
conjunto factível com a seguinte ordem:

1. maior número de lacunas cobertas;
2. menor número de projetos;
3. menor custo total;
4. menor número de capacidades redundantes;
5. identidade lexical de `repo@commit` como desempate estável.

Não havendo cobertura completa, o resultado pode ser parcial e deve listar as
lacunas descobertas. Não havendo composição segura, `Feasible=false` é um
resultado correto.

## Experimento competitivo

O protocolo executável e os inputs pré-registrados estão em
[`benchmark/0012/experimento`](../benchmark/0012/experimento/README.md). Eles não
foram executados nem constituem resultado de produto enquanto os bundles e
oráculos independentes permanecerem ausentes.

### Ideia alvo

Congelar uma `IdeaSpec` do Bateia que preserve, no mínimo:

- receber uma ideia antes de iniciar um projeto;
- descobrir projetos públicos, privados ou locais;
- decompor capacidades e restrições;
- eliminar semelhantes irrelevantes;
- comprovar capacidades em artefatos rastreáveis;
- decidir entre `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE`, `BUILD`
  e `INCONCLUSIVE`;
- operar sem LLM no pipeline principal;
- reproduzir análise offline.

### Casos obrigatórios

1. **Bateia → Bateia**: inventariar capacidades internas e lacunas.
2. **Bateia → concorrentes conhecidos**: avaliar as sementes competitivas sem
   assumir previamente que são relevantes.
3. **Bateia → controles negativos**: incluir ao menos três projetos com termos
   parecidos, mas função diferente.
4. **Katamari completo**: fixture em que dois ou mais projetos cobrem lacunas
   complementares.
5. **Katamari redundante**: projeto externo oferece somente o que já existe
   internamente e deve ser rejeitado.
6. **Katamari proibido**: boa cobertura funcional com gate eliminatório.
7. **Katamari opaco**: evidência insuficiente deve permanecer inconclusiva.

Os nomes conhecidos ajudam a medir recall, mas não podem ser inseridos à força no
resultado. Se o Bateia não os descobrir, isso é um finding de descoberta.

## Saídas e CLI

```bash
bateia idea autofagia \
  --idea benchmark/bateia-idea.json \
  --bundles benchmark/bateia-internal.json \
  --rules evidence-rules.json \
  --output autofagia.json

bateia idea katamari \
  --autofagia autofagia.json \
  --assessments competitors-assessment.json \
  --policy katamari-policy.json \
  --output katamari.json
```

Ambos os comandos são offline. Descoberta e coleta continuam em comandos
anteriores, preservando a fronteira entre aquisição e julgamento.

## Observáveis

O relatório competitivo deve produzir uma matriz verificável:

| Capacidade | Bateia | Concorrentes comprovados | Lacuna | Componente escolhido |
|---|---|---|---|---|

Além dela, deve registrar:

- cobertura interna antes da composição;
- cobertura adicional do Katamari;
- projetos considerados, rejeitados e razão;
- duplicações evitadas;
- gaps ainda abertos;
- quantidade de componentes e custo determinístico;
- hashes de todos os inputs e políticas.

## Métricas

```go
type CompetitiveAutophagyMetrics struct {
    InternalCapabilityPrecision *float64 `json:"internal_capability_precision"`
    InternalCapabilityRecall    *float64 `json:"internal_capability_recall"`
    CompetitorDiscoveryRecall   *float64 `json:"competitor_discovery_recall"`
    IrrelevantRejectionRate     *float64 `json:"irrelevant_rejection_rate"`
    GapCoverageBefore           *float64 `json:"gap_coverage_before"`
    GapCoverageAfter            *float64 `json:"gap_coverage_after"`
    RedundantAbsorptionRate     *float64 `json:"redundant_absorption_rate"`
    ForbiddenAbsorptions        int      `json:"forbidden_absorptions"`
    SelectedProjects            int      `json:"selected_projects"`
    MinimalityRegret            *float64 `json:"minimality_regret"`
}
```

Denominadores vazios permanecem `null`, nunca `100%`.

## Gates de aprovação

- `internal_capability_precision ≥ 90%`;
- `internal_capability_recall ≥ 85%`;
- `competitor_discovery_recall ≥ 70%` no conjunto conhecido congelado;
- `irrelevant_rejection_rate = 100%` nos controles negativos;
- cobertura após Katamari maior que antes em ao menos dois casos correspondentes;
- `redundant_absorption_rate = 0%`;
- `forbidden_absorptions = 0`;
- seletor encontra a solução mínima em `100%` das fixtures enumeráveis;
- reordenação de inputs não altera proposta nem hashes canônicos;
- nenhum comando offline acessa rede, cache externo ou credenciais;
- nenhuma regra existente é afrouxada para atingir o gate.

Os gates são hipóteses iniciais. Podem ser revisados somente antes de abrir o
split de avaliação, nunca depois de observar um resultado candidato.

## Testes obrigatórios

1. capacidade interna comprovada produz `KEEP`;
2. duplicação interna comprovada produz `CONSOLIDATE`;
3. ausência com coleta completa produz `GAP`;
4. coleta incompleta produz `UNKNOWN`, não `GAP`;
5. capacidade `KEEP` não entra no problema de cobertura;
6. capacidade externa apenas documentada não entra no Katamari;
7. capacidade externa `PROVEN` cobre a lacuna correspondente;
8. duas fontes complementares podem produzir `COMPOSE`;
9. fonte redundante é rejeitada;
10. conflito de licença elimina a fonte;
11. revisão `HEAD` elimina ou torna a fonte inconclusiva conforme política;
12. operação remota obrigatória viola política local;
13. projeto arquivado recebe custo, sem inventar ausência funcional;
14. solução com menos projetos vence quando a cobertura é igual;
15. menor custo vence quando número de projetos e cobertura são iguais;
16. desempate por identidade é determinístico;
17. limite de projetos pode produzir composição parcial;
18. composição inviável retorna `Feasible=false`;
19. lacuna não coberta permanece visível;
20. ordem e duplicação de candidatos não alteram o resultado;
21. hash de input muda quando política, bundle ou autofagia muda;
22. controle lexical irrelevante nunca é absorvido;
23. ausência de concorrente conhecido reduz recall, sem inserção artificial;
24. utilidade não é inferida automaticamente da cobertura;
25. execução offline não lê token, rede ou relógio;
26. cada consumidor Go possui prompt próprio e o linter V15 permanece verde.

## Critério de produto

O passo demonstra valor quando o relatório permite responder, com evidência:

1. quais capacidades o Bateia já possui;
2. quais são lacunas reais e quais são apenas opacidade de coleta;
3. que concorrente comprova cada capacidade ausente;
4. qual é a menor composição segura para aumentar cobertura;
5. por que alternativas populares, redundantes ou incompatíveis foram rejeitadas;
6. qual experimento isolado deve testar cada absorção antes de implementá-la.

O Katamari é um plano de composição, não autorização para incorporar código ou
dependência. Cada capacidade aprovada origina um passo pequeno próprio.

## Critério de parada

Interromper e manter `EM VALIDAÇÃO` se:

- os bundles do Bateia ou concorrentes não forem correspondentes à ideia;
- as capacidades internas forem rotuladas pela saída que deveriam avaliar;
- a descoberta depender de nomes inseridos depois do resultado;
- algum gate eliminatório for convertido em penalidade compensável;
- a solução “mínima” não puder ser reproduzida por enumeração;
- um caso `UNKNOWN` for tratado como lacuna ou cobertura;
- for necessário copiar código ou usar LLM para atingir os gates;
- não houver revisão humana das capacidades competitivas.

## Fora de escopo

- implementar capacidades copiadas dos concorrentes;
- instalar ou executar concorrentes;
- alterar heurísticas de evidência, decisão ou recuperação;
- monitoramento contínuo da internet;
- análise jurídica de compatibilidade;
- otimização para corpora grandes antes de provar a enumeração mínima;
- interface gráfica;
- ranking comercial ou leaderboard público;
- aprendizagem automática por feedback de adoção.

## Próximo passo previsto

Se a autofagia e o Katamari passarem pelos gates, escolher exatamente uma lacuna
de alto valor e baixo custo, escrever um passo independente e testá-la contra o
benchmark congelado. As primeiras candidatas prováveis são feedback explícito de
adoção, sinais objetivos de saúde ou comparação de componentes, mas o resultado
do experimento — não esta lista — decide qual será implementada.
