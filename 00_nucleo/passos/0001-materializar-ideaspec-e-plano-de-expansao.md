# Passo 0001 — Materializar `IdeaSpec → ExpansionPlan`

**Estado**: `IMPLEMENTADO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1 e L2  
**Natureza**: primeiro corte vertical; contrato e observável, ainda sem descoberta

## Intenção

Dar ao Bateia uma representação explícita da ideia e provar que ela pode ser
convertida em **múltiplas expansões pequenas, independentes e rastreáveis**.

Ao final deste passo, ainda não queremos encontrar repositórios. Queremos conseguir
inspecionar e testar a pergunta que será feita ao mundo antes de atribuir ruído ao
Sourcegraph, ao GitHub ou ao ranking.

## Por que este é o primeiro passo

O fluxo atual mistura três coisas num único texto:

1. a ideia que define o alvo;
2. a query usada para descobrir candidatos;
3. os termos usados posteriormente para pontuar similaridade.

Na execução de 2026-08-23, `-target-idea` deixou de orientar a comparação quando
`-query` também foi informada. O resultado foi mensurável: 26 candidatos, score
máximo de 18%, stack zerada e muitos falsos positivos com relevância alta.

Antes de melhorar busca ou score, precisamos tornar a separação observável:

```text
IdeaSpec (o que queremos)
    ↓
ExpansionPlan (quais perguntas tentaremos)
    ↓ passo futuro
Discovery (o que o mundo devolveu)
```

## Entrega

### 1. Vocabulário de domínio em L1

Criar tipos puros equivalentes a:

```go
type IdeaSpec struct {
    Problem        string   `json:"problem"`
    Users          []string `json:"users,omitempty"`
    Capabilities   []string `json:"capabilities"`
    Constraints    []string `json:"constraints,omitempty"`
    NonGoals       []string `json:"non_goals,omitempty"`
    PreferredStack []string `json:"preferred_stack,omitempty"`
}

type ExpansionFacet string

const (
    FacetProblem      ExpansionFacet = "problem"
    FacetCapability   ExpansionFacet = "capability"
    FacetArchitecture ExpansionFacet = "architecture"
    FacetConstraint   ExpansionFacet = "constraint"
)

type SearchExpansion struct {
    ID        string         `json:"id"`
    Facet     ExpansionFacet `json:"facet"`
    Query     string         `json:"query"`
    Rationale string         `json:"rationale"`
    Covers    []string       `json:"covers"`
}

type ExpansionPlan struct {
    Idea       IdeaSpec          `json:"idea"`
    Expansions []SearchExpansion `json:"expansions"`
}
```

Os nomes podem mudar se a implementação revelar um vocabulário mais preciso, mas
as distinções não podem desaparecer.

### 2. Expansor puro e determinístico em L1

Criar um caso de uso puro que recebe `IdeaSpec` já estruturada e produz
`ExpansionPlan`.

Regras mínimas deste passo:

- normalizar espaços e remover entradas vazias;
- rejeitar ideia sem problema e sem capacidades;
- produzir ao menos uma expansão para o problema;
- produzir uma expansão independente por capacidade;
- produzir expansões de arquitetura ou restrição somente quando houver dados;
- cada expansão deve declarar `Facet`, `Rationale` e `Covers`;
- queries idênticas após normalização devem ser deduplicadas;
- IDs devem ser determinísticos para a mesma entrada;
- a ordem da saída deve ser estável;
- `NonGoals` não viram queries positivas;
- `PreferredStack` não define a identidade da ideia e não pode dominar o plano.

Não criar sinônimos “inteligentes” neste passo. Se a entrada disser apenas
`impact analysis`, não inventar `blast radius`. O objetivo inicial é provar o
contrato e a rastreabilidade; expansão semântica entra depois sobre uma base
mensurável.

### 3. Boca de inspeção em L2

Adicionar um comando que somente materializa o plano, sem rede:

```bash
bateia idea plan --input ideia.json
bateia idea plan --input ideia.json --output plano.json
```

Requisitos:

- entrada JSON contendo `IdeaSpec`;
- saída JSON contendo `ExpansionPlan`;
- `stdout` quando `--output` não for informado;
- erro claro para JSON inválido ou ideia vazia;
- nenhuma inicialização de Sourcegraph, GitHub, Bleve ou cache;
- execução determinística e adequada para golden tests.

Se a arquitetura atual obrigar o binário a inicializar infraestrutura antes de
reconhecer o comando, corrigir a composição para que `idea plan` permaneça
realmente offline. Não mascarar o acoplamento com mocks.

### 4. Prompt e linhagem

Criar as especificações L0 correspondentes aos novos componentes antes de fechar o
passo. Cada arquivo de produção novo deve apontar para seu prompt real por
`Crystalline Lineage`, preservando a trava arquitetural existente.

## Exemplo de referência

Entrada:

```json
{
  "problem": "mostrar a consequência estrutural de alterar código",
  "users": ["agentes de programação", "mantenedores"],
  "capabilities": [
    "encontrar dependentes transitivos de um símbolo",
    "visualizar dependências e ciclos entre módulos",
    "analisar somente o impacto de um diff"
  ],
  "constraints": ["execução local", "evidência rastreável"],
  "non_goals": ["monitoramento de comportamento em runtime"],
  "preferred_stack": ["Rust"]
}
```

Forma esperada da saída — os textos exatos podem ser refinados pelos testes, mas a
separação precisa sobreviver:

```json
{
  "expansions": [
    {
      "id": "exp-problem-...",
      "facet": "problem",
      "query": "mostrar a consequência estrutural de alterar código",
      "rationale": "descobrir soluções para o problema central",
      "covers": ["problem"]
    },
    {
      "id": "exp-capability-...",
      "facet": "capability",
      "query": "encontrar dependentes transitivos de um símbolo",
      "rationale": "buscar implementação da capacidade isoladamente",
      "covers": ["capabilities[0]"]
    }
  ]
}
```

Não exigir que o plano já contenha `blast radius`, `repository memory` ou outras
formulações ausentes da entrada. Esses termos serão a matéria do próximo
experimento de expansão semântica.

## Testes obrigatórios

### Unidade — L1

1. ideia completa produz problema + uma expansão por capacidade;
2. mesma entrada produz JSON idêntico em execuções repetidas;
3. capacidades duplicadas por caixa/espaço produzem uma expansão;
4. ideia sem problema e sem capacidades é rejeitada;
5. `NonGoals` nunca aparecem como expansão positiva;
6. stack preferida ausente não impede o plano;
7. somente stack, sem problema/capacidade, é inválido;
8. ordem das expansões independe da ordem acidental de mapas;
9. campos da entrada não são modificados pelo expansor;
10. cancelamento, rede, relógio e filesystem não fazem parte do caso de uso.

### Integração — L2

1. `idea plan --input fixture.json` produz JSON parseável;
2. sem `--output`, escreve apenas o plano em `stdout`;
3. com `--output`, grava o arquivo e não inicializa infraestrutura remota;
4. JSON inválido retorna código de erro e mensagem acionável;
5. golden test fixa a forma pública do plano;
6. o comando funciona sem `.env`, token, rede e diretório de cache gravável.

## Prova de mordida

Criar uma fixture de controle que reproduza a confusão atual:

```text
target idea: ferramenta de impacto estrutural para agentes
query: code intelligence AST dependency graph
```

O teste deve demonstrar que o novo modelo preserva os dois dados em lugares
distintos: a `IdeaSpec` continua sendo o alvo; as expansões são derivadas dela e
não são substituídas silenciosamente por uma query externa.

Se alvo e query ainda puderem ocupar o mesmo campo sem que um teste falhe, o passo
não mordeu e não está concluído.

## Critérios de verificação

- [x] `IdeaSpec`, `SearchExpansion` e `ExpansionPlan` existem em L1 sem I/O.
- [x] Expansor é determinístico, deduplica e preserva rastreabilidade por faceta.
- [x] O plano distingue problema, capacidades, restrições e stack.
- [x] `NonGoals` não são promovidos a busca positiva.
- [x] `bateia idea plan` funciona completamente offline.
- [x] A prova de mordida impede a regressão alvo/query.
- [x] Golden test torna a saída pública deliberada.
- [x] Novos componentes possuem prompt e linhagem válidos.
- [x] `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.
- [x] Um laudo registra o que foi implementado, medido e deixado aberto.

## Fora de escopo

- interpretar texto livre e preencher `IdeaSpec` automaticamente;
- LLM, embeddings ou provedor de expansão semântica;
- dicionário de sinônimos;
- executar as expansões no Sourcegraph;
- agregar ou deduplicar candidatos de várias buscas;
- modificar o score de similaridade;
- extrair evidências dos candidatos;
- emitir decisões `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou
  `BUILD`;
- interface TUI ou MCP para este fluxo.

## Próximo passo condicionado

Somente após o plano ser observável e estável, comparar estratégias de expansão
semântica sobre um corpus pequeno de ideias:

1. regras/dicionário local;
2. embeddings ou recuperação lexical;
3. modelo de linguagem local ou remoto;
4. combinação híbrida.

A estratégia seguinte deve ser escolhida por precisão e diversidade das expansões,
não por sofisticação presumida.
