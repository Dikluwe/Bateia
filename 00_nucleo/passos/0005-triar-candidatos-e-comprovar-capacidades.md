# Passo 0005 — Triar candidatos e comprovar capacidades

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0004 — Construir peneira probabilística explicável`](0004-construir-peneira-probabilistica-explicavel.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: triagem profunda determinística e matriz de evidências por capacidade

## Intenção

Transformar candidatos priorizados em avaliações verificáveis de reutilização.

O passo deve responder, para cada projeto:

1. qual parte do problema ele realmente resolve;
2. quais capacidades da `IdeaSpec` possuem evidência;
3. quais capacidades estão ausentes ou inconclusivas;
4. quais restrições são compatíveis ou conflitantes;
5. quais componentes concretos parecem reutilizáveis.

A saída deixa de ser “este repositório parece relevante” e passa a ser:

```text
esta capacidade foi encontrada
→ neste arquivo
→ neste símbolo, trecho, teste ou manifesto
→ neste commit
→ por esta regra determinística
```

## Restrição central: pipeline sem LLM

Este passo e o pipeline de produto do Bateia permanecem **sem LLM**.

São proibidos durante descoberta, triagem, extração e classificação:

- modelos de linguagem locais ou remotos;
- embeddings ou APIs de embeddings;
- rerankers neurais;
- geração automática de explicações por modelo;
- classificação por prompts;
- chamadas indiretas a serviços de IA;
- fallback silencioso para qualquer mecanismo probabilístico não versionado.

São permitidos:

- parsing determinístico;
- tokenização e normalização explícitas;
- regras lexicais versionadas;
- expressões regulares limitadas e testadas;
- análise de manifestos;
- árvores de arquivos;
- tabelas de símbolos e ASTs já suportidas pelo Bateia;
- testes, exemplos e documentação do candidato;
- métricas e modelos estatísticos transparentes, reproduzíveis e calibrados;
- heurísticas com pesos e proveniência declarados.

O pipeline pode produzir `INCONCLUSIVE`. Nunca deve inventar uma capacidade para
preencher uma lacuna que somente interpretação humana ou um modelo poderia inferir.

## Por que este passo vem agora

O Passo 0004 mostrou que a peneira consegue explicar o ruído e ordenar candidatos,
mas seus sinais ainda são preliminares. Eles derivam principalmente de queries,
descrições e matches textuais.

Isso é insuficiente para decidir reutilização:

- `supermodeltools/cli` declara blast radius, mas envia o repositório para uma API;
- `mllint` realiza análise estática, mas declara manutenção encerrada;
- `reviewdog/errorformat` não é um analisador completo, mas pode ser um componente
  reutilizável para integrar diagnósticos;
- `enola` declara grafo, impacto estrutural, execução local e MCP, mas essas
  capacidades precisam ser ligadas a arquivos e símbolos concretos.

O próximo observável é uma matriz de capacidades sustentada por evidências em
referências imutáveis.

## Incertezas reduzidas

- se a descrição do repositório corresponde à implementação;
- se uma capacidade está implementada, apenas planejada ou somente mencionada;
- qual parte do candidato pode ser reutilizada isoladamente;
- se licença, manutenção, stack ou operação impedem reuso;
- se o candidato resolve o problema ou apenas compartilha vocabulário;
- quanto da avaliação permanece inconclusivo após extração determinística.

## Entrega

### 1. Contrato de evidência em L1

Criar tipos puros equivalentes a:

```go
type EvidenceStatus string

const (
    EvidenceProven       EvidenceStatus = "PROVEN"
    EvidencePartial      EvidenceStatus = "PARTIAL"
    EvidenceAbsent       EvidenceStatus = "ABSENT"
    EvidenceInconclusive EvidenceStatus = "INCONCLUSIVE"
    EvidenceConflict     EvidenceStatus = "CONFLICT"
)

type EvidenceKind string

const (
    EvidenceDocumentation EvidenceKind = "documentation"
    EvidenceManifest      EvidenceKind = "manifest"
    EvidenceFile          EvidenceKind = "file"
    EvidenceSymbol        EvidenceKind = "symbol"
    EvidenceTest          EvidenceKind = "test"
    EvidenceExample       EvidenceKind = "example"
    EvidenceLicense       EvidenceKind = "license"
    EvidenceRepository    EvidenceKind = "repository_metadata"
)

type EvidenceSource struct {
    Kind       EvidenceKind `json:"kind"`
    Ref        CommitRef    `json:"ref"`
    Path       string       `json:"path,omitempty"`
    Symbol     string       `json:"symbol,omitempty"`
    Line       int          `json:"line,omitempty"`
    Snippet    string       `json:"snippet,omitempty"`
    RuleID     string       `json:"rule_id"`
    ContentHash string      `json:"content_hash,omitempty"`
}

type CapabilityEvidence struct {
    Capability  string           `json:"capability"`
    Status      EvidenceStatus   `json:"status"`
    Sources     []EvidenceSource `json:"sources"`
    Confidence  float64          `json:"confidence"`
    Explanation string           `json:"explanation"`
}
```

Regras:

- `PROVEN` exige ao menos uma evidência forte;
- `PARTIAL` exige evidência concreta de somente parte da capacidade;
- `ABSENT` exige busca executada e ausência segundo escopo declarado;
- ausência de conteúdo, falha de fetch ou parser incompatível gera
  `INCONCLUSIVE`, nunca `ABSENT`;
- `CONFLICT` registra incompatibilidade explícita com uma restrição ou `NonGoal`;
- toda fonte aponta para `repo@commit` real;
- snippet é curto, normalizado e acompanhado de hash do conteúdo quando aplicável;
- explicações são montadas por templates determinísticos, não por geração livre.

### 2. Avaliação completa do candidato

Criar uma saída equivalente a:

```go
type CandidateAssessment struct {
    Repository        CommitRef             `json:"repository"`
    ScreeningDecision ScreeningDecision     `json:"screening_decision"`
    ProblemAlignment  CapabilityEvidence    `json:"problem_alignment"`
    Capabilities      []CapabilityEvidence  `json:"capabilities"`
    Constraints       []ConstraintEvidence  `json:"constraints"`
    License           LicenseAssessment     `json:"license"`
    Maintenance       MaintenanceAssessment `json:"maintenance"`
    ReusableParts     []ReusableComponent   `json:"reusable_parts"`
    Gaps              []string              `json:"gaps"`
    Failures          []ExecutionFailure    `json:"failures,omitempty"`
}
```

O contrato deve manter vetores separados. Uma licença incompatível não pode ser
escondida por alta cobertura funcional. Um projeto sem manutenção não perde sua
capacidade técnica, mas recebe risco operacional separado.

### 3. Plano determinístico de extração em L1

Gerar um plano por candidato antes de ler conteúdo:

```go
type EvidenceExtractionPlan struct {
    Repository CommitRef          `json:"repository"`
    Targets    []ExtractionTarget `json:"targets"`
    Budget     ExtractionBudget   `json:"budget"`
}
```

Ordem mínima de busca:

1. README principal e documentação indexada;
2. licença;
3. manifestos reconhecidos;
4. árvore até profundidade limitada;
5. arquivos e símbolos ligados às capacidades;
6. testes e exemplos relacionados;
7. metadados de manutenção.

O plano deve declarar por que cada alvo foi escolhido e qual capacidade cobre.
Não permitir crawling irrestrito.

Orçamento inicial por candidato:

```text
documentos: até 8
manifestos: até 4
arquivos de código: até 20
testes/exemplos: até 12
bytes totais: até 2 MiB
profundidade de árvore: até 4
```

Os limites são configuráveis dentro de faixas seguras e aparecem no relatório.

### 4. Regras versionadas de capacidades

Criar um catálogo declarativo, auditável e pequeno:

```json
{
  "id": "dependency-graph",
  "capability_terms": ["grafo de dependências", "dependency graph"],
  "strong_evidence": [
    {"kind":"symbol", "patterns":["DependencyGraph", "BuildGraph"]},
    {"kind":"test", "patterns":["dependency graph", "transitive dependency"]}
  ],
  "supporting_evidence": [
    {"kind":"documentation", "patterns":["dependency graph"]},
    {"kind":"path", "patterns":["graph", "dependency"]}
  ],
  "negative_evidence": [
    {"patterns":["planned", "roadmap", "not implemented"]}
  ]
}
```

Requisitos:

- IDs estáveis;
- correspondência normalizada e limitada por palavra/frase;
- regras diferentes para evidência forte, auxiliar e negativa;
- nenhuma regex fornecida remotamente;
- regex local validada por tamanho e complexidade;
- regra registra versão e hash no relatório;
- capacidade desconhecida permanece `INCONCLUSIVE` até possuir regra adequada;
- regras não são ajustadas silenciosamente depois de observar o candidato.

O catálogo inicial cobre somente as capacidades do benchmark:

- grafo de dependências;
- dependentes transitivos;
- impacto ou blast radius;
- comparação estrutural de mudanças;
- análise estática;
- execução local/offline;
- integração com agentes ou MCP;
- processamento assíncrono de eventos;
- organização e importação de receitas.

### 5. Hierarquia de força da evidência

Adotar uma tabela explícita:

| Evidência | Força inicial | Pode provar sozinha? |
| :-------- | ------------: | :------------------- |
| teste executável relacionado | 1,00 | sim |
| símbolo + implementação | 0,90 | sim |
| exemplo funcional + símbolo | 0,85 | sim |
| manifesto/dependência específica | 0,70 | somente certas restrições |
| documentação + arquivo relacionado | 0,65 | parcial; pode provar com convergência |
| documentação isolada | 0,40 | não |
| path ou nome de arquivo | 0,20 | não |
| descrição do repositório | 0,10 | não |
| match em `vendor`, fixture ou tradução | 0,05 | não |

Os valores são baseline versionado, não probabilidades universais. A combinação
de evidências deve evitar dupla contagem do mesmo arquivo ou origem.

Critérios iniciais:

```text
PROVEN       força independente ≥ 0,85
             ou duas fontes independentes com força combinada ≥ 0,90
PARTIAL      força combinada ≥ 0,50, sem prova completa
ABSENT       busca completa + zero sinal + zero falha relevante
INCONCLUSIVE orçamento incompleto, parser ausente ou falha de acesso
CONFLICT     evidência negativa explícita ligada à restrição
```

### 6. Licença, manutenção e modo de operação

Avaliar separadamente:

#### Licença

- arquivo de licença e identificador reconhecido;
- compatibilidade declarada por política configurável;
- licença ausente ou ambígua gera `INCONCLUSIVE`;
- não oferecer aconselhamento jurídico automático.

#### Manutenção

- commit mais recente quando disponível;
- estado arquivado;
- releases ou tags;
- aviso explícito de manutenção encerrada;
- testes/CI declarados;
- número de estrelas nunca serve sozinho como maturidade.

#### Operação

- execução local versus serviço obrigatório;
- envio de código para API externa;
- necessidade de token ou conta;
- dependência de banco ou infraestrutura;
- linguagens e runtimes exigidos.

Exemplo esperado:

```text
supermodeltools/cli
capacidade de grafo: PROVEN
blast radius: PROVEN
execução estritamente local: CONFLICT
evidência: README declara upload do repositório para Supermodel API
```

### 7. Componentes reutilizáveis

Um `ReusableComponent` deve apontar para algo concreto:

```go
type ReusableComponent struct {
    Name         string           `json:"name"`
    Kind         string           `json:"kind"` // package, module, binary, algorithm, pattern
    Path         string           `json:"path,omitempty"`
    Symbols      []string         `json:"symbols,omitempty"`
    Capabilities []string         `json:"capabilities"`
    Evidence     []EvidenceSource `json:"evidence"`
    Caveats      []string         `json:"caveats,omitempty"`
}
```

Não emitir componente apenas porque uma pasta possui nome interessante. Exigir
símbolo, arquivo, manifesto, teste ou documentação convergente.

### 8. Seleção de candidatos

Entrada da triagem:

- todos os `ADVANCE`;
- `AUDIT` ordenados por probabilidade/score até o orçamento;
- amostra determinística de `DROP` para medir falsos negativos.

Configuração inicial:

```text
ADVANCE: todos
AUDIT: até 10 por ideia
DROP auditado: 2 por ideia
máximo total: 16 por ideia
```

Julgamentos conhecidos não podem orientar as regras durante a execução normal.
Eles entram somente em replay de avaliação.

### 9. Boca de inspeção em L2

Adicionar:

```bash
bateia idea evidence plan \
  --screening screening.json \
  --output extraction-plan.json

bateia idea evidence collect \
  --plan extraction-plan.json \
  --output evidence-bundle.json

bateia idea evidence assess \
  --idea ideia.json \
  --screening screening.json \
  --evidence evidence-bundle.json \
  --rules capability-rules.json \
  --output assessment.json
```

Separar planejamento, coleta e avaliação permite:

- inspecionar o orçamento antes da rede;
- congelar evidências para replay offline;
- alterar regras sem baixar novamente os projetos;
- testar avaliação sem GitHub ou Sourcegraph;
- atribuir falhas à coleta ou à classificação corretamente.

`plan` e `assess` são completamente offline. Somente `collect` inicializa o
provedor de conteúdo.

## Prova de mordida

Usar pelo menos estes controles:

### Enola

Deve encontrar evidência em `repo@commit` para:

- grafo de dependências;
- comparação de mudança estrutural;
- execução local;
- integração MCP;
- símbolos ou testes relacionados.

### Supermodel CLI

Deve preservar capacidade de grafo/impacto e, simultaneamente, registrar conflito
com a restrição de execução estritamente local por upload para API.

### mllint

Deve reconhecer análise estática, mas registrar aviso explícito de manutenção
encerrada sem transformar isso em ausência funcional.

### Ruído vendorizado

Um candidato encontrado apenas por comentário em `vendor/k8s.io/...` deve ficar
`ABSENT` ou `INCONCLUSIVE`, nunca `PROVEN`.

O teste deve falhar se uma descrição isolada provar capacidade, se `vendor` receber
força alta ou se uma falha de fetch virar ausência.

## Avaliação

Aplicar o pipeline aos candidatos de impacto estrutural e análise estática do
benchmark anterior.

Revisar manualmente:

- todos os `PROVEN`;
- todos os `CONFLICT`;
- amostra de cinco `ABSENT`;
- todos os `INCONCLUSIVE` causados por falha;
- componentes reutilizáveis emitidos.

Métricas:

```text
evidence_precision = capacidades confirmadas na revisão / capacidades marcadas PROVEN
false_absence_rate = ABSENT incorretos / ABSENT revisados
inconclusive_rate  = avaliações INCONCLUSIVE / avaliações totais
source_trace_rate  = fontes válidas em repo@commit / fontes emitidas
budget_completion  = alvos coletados / alvos planejados
```

Critérios iniciais:

- `evidence_precision` ≥ 90%;
- `false_absence_rate` = 0 na amostra;
- `source_trace_rate` = 100%;
- nenhuma capacidade provada somente por descrição;
- nenhuma falha convertida em `ABSENT`;
- todos os componentes reutilizáveis possuem fonte concreta.

Se menos de dez capacidades forem marcadas `PROVEN`, a precisão deve ser mostrada,
mas considerada amostra insuficiente.

## Testes obrigatórios

### Unidade — L1

1. documentação isolada não produz `PROVEN`;
2. símbolo implementado ou teste relacionado pode produzir `PROVEN`;
3. duas evidências auxiliares independentes podem convergir;
4. evidências do mesmo arquivo não contam como independentes;
5. match em `vendor`, fixture ou tradução recebe força mínima;
6. evidência negativa explícita gera `CONFLICT` quando ligada à restrição;
7. falha, orçamento esgotado e parser ausente geram `INCONCLUSIVE`;
8. `ABSENT` exige busca concluída;
9. ordem de arquivos não altera resultado;
10. snippets e hashes são determinísticos;
11. explicações vêm de templates conhecidos;
12. regra desconhecida não inventa capacidade;
13. entrada, plano e bundle não são modificados;
14. seleção de `DROP` é reproduzível pela seed.

### Integração — L2/L3/L4

1. candidato é resolvido para commit imutável antes da coleta;
2. README, licença, manifesto e árvore preservam path e commit;
3. limites de arquivo, bytes e profundidade são respeitados;
4. conteúdo duplicado é deduplicado por hash;
5. falha parcial aparece no bundle e no assessment;
6. `evidence plan` e `evidence assess` funcionam com `env -i`;
7. replay do bundle não acessa rede;
8. somente `evidence collect` constrói o provedor de conteúdo;
9. tokens e headers nunca aparecem nos artefatos;
10. schemas e regras inválidos produzem erros acionáveis;
11. golden literal do Passo 0001 permanece intacto;
12. `crystalline-lint` preserva as camadas.

## Critérios de verificação

- [ ] Pipeline de evidência não usa LLM, embeddings ou reranker neural.
- [ ] Toda evidência aponta para `repo@commit` real.
- [ ] Estados `PROVEN`, `PARTIAL`, `ABSENT`, `INCONCLUSIVE` e `CONFLICT` são distintos.
- [ ] Planejamento, coleta e avaliação permanecem etapas separadas.
- [ ] Coleta respeita orçamento explícito e produz replay offline.
- [ ] Regras de capacidade são versionadas e auditáveis.
- [ ] Descrição isolada nunca prova capacidade.
- [ ] Falha ou parser ausente nunca vira `ABSENT`.
- [ ] Licença, manutenção e operação permanecem vetores separados.
- [ ] Componentes reutilizáveis possuem evidência concreta.
- [ ] Enola, Supermodel, mllint e ruído vendorizado mordem casos diferentes.
- [ ] Métricas de evidência são calculadas e revisadas.
- [ ] Novos componentes possuem prompts L0 e linhagem válida.
- [ ] `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.
- [ ] Um laudo registra precisão, ausências falsas, falhas e limites.

## Como este passo pode falhar

- README continua sendo tratado como implementação comprovada;
- termos em dependências vendorizadas provam capacidade do projeto;
- arquivos duplicados inflam força da evidência;
- busca incompleta é interpretada como ausência;
- regra nasce ajustada especificamente para um único candidato;
- commit não é resolvido e a evidência muda depois;
- licença ou requisito de API é escondido pelo score funcional;
- componente reutilizável não aponta para módulo ou símbolo concreto;
- orçamento irrestrito torna a coleta lenta e cara;
- uma biblioteca útil é descartada por não resolver o produto inteiro;
- explicação textual diverge das evidências por usar template incorreto;
- um serviço de IA entra como dependência indireta do pipeline.

## Fora de escopo

- LLM, embeddings ou qualquer inferência generativa no pipeline;
- decisão final `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD`;
- composição automática entre múltiplos candidatos;
- cálculo jurídico de compatibilidade de licença;
- executar código não confiável dos candidatos;
- instalar dependências dos repositórios analisados;
- crawling integral de repositórios;
- suportar toda linguagem ou manifesto neste primeiro corte;
- modificar o score probabilístico durante a avaliação;
- interface gráfica ou TUI.

## Próximo passo condicionado

O [Passo 0006](0006-endurecer-e-calibrar-evidencias.md) corrigirá as lacunas de
completude, árvore remota, identidade, símbolos, componentes e vetores operacionais
encontradas na auditoria. Depois executará o corpus real e medirá precisão,
ausências falsas, rastreabilidade e componentes.

Somente após essa calibração, o **Passo 0007** construirá a matriz comparativa e o
motor determinístico `FORK · EXTEND · COMPOSE · INTEGRATE · REFERENCE · BUILD`.
Casos difíceis continuarão `INCONCLUSIVE`; a resposta não será adicionar LLM.
