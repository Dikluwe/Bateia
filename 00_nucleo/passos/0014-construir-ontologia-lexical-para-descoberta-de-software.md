# Passo 0014 — Construir ontologia lexical para descoberta de software

**Estado**: `IMPLEMENTADO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0013 — Limitar descoberta e preservar resultados parciais`](0013-limitar-descoberta-e-preservar-resultados-parciais.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2 e L3  
**Natureza**: ontologia lexical determinística, qualidade de queries e avaliação de recall sem LLM

## Evidência que motiva este passo

A repetição pública E-0012 após o Passo 0013 tornou a falha mensurável:

| Observável | Resultado |
|---|---:|
| expansões planejadas | 16 |
| expansões lexicais novas | 0 |
| attempts expirados | 18 |
| attempts bem-sucedidos | 1 |
| candidatos preservados | 3 |
| concorrentes conhecidos encontrados | 0/4 |
| recall competitivo | 0% |

O único fallback bem-sucedido reduziu:

```text
decidir entre fork extend compose integrate reference build e inconclusive
```

a:

```text
decidir entre
```

Essa query encontrou um addon de jogo, material educacional sobre REST e um guia
de tradução. Os três resultados eram funcionalmente irrelevantes e estavam em
`HEAD`.

O comportamento prova duas lacunas distintas:

1. o vocabulário atual não conhece a categoria de descoberta e seleção de
   software;
2. o fallback escolhe os primeiros termos não triviais, não os conceitos mais
   discriminatórios.

O Bateia já limita, registra e preserva a descoberta. Agora precisamos melhorar
o que ele pergunta sem inserir as respostas conhecidas nas queries.

## Incerteza que este passo reduz

Não sabemos se um vocabulário pequeno, explícito e determinístico é suficiente
para encontrar projetos funcionalmente próximos a partir de uma IdeaSpec de
produto, sem LLM ou embeddings.

Também não sabemos se Sourcegraph consegue sustentar a primeira fase de
descoberta por produto ou se deve ficar restrito à comprovação de capacidades
depois que outro catálogo identifica os repositórios.

Este passo deve responder essas perguntas com um corpus e gates congelados.

## Hipótese

Uma ontologia lexical versionada consegue transformar intenções como:

```text
não começar do zero
descobrir projetos por intenção
decidir um ponto de partida
compor capacidades existentes
avaliar adoção com evidências
```

em queries curtas e discriminatórias:

```text
software reuse discovery
open source project discovery
GitHub repository recommendation
software component selection
dependency adoption decision
build versus reuse
evidence backed tool evaluation
software composition planning
```

sem mencionar Replen, OSS Compass, ARIS, Starboard ou qualquer outro expected
outcome.

## Restrição central: conceitos, não respostas

É proibido:

- usar nomes, owners, URLs ou trechos exclusivos dos concorrentes esperados;
- gerar queries a partir dos repositórios do registry;
- observar o resultado público antes de congelar vocabulário e corpus;
- adicionar sinônimo apenas porque ele recuperou uma semente conhecida;
- usar LLM, embeddings, reranker neural ou tradução automática;
- substituir avaliação funcional por estrelas, tópicos ou popularidade;
- remover controles negativos depois de observar contaminação;
- alterar screening ou regras de evidência neste passo;
- tratar mais candidatos como melhor descoberta;
- abrir o Katamari antes de comprovar correspondência.

As sementes existem somente no avaliador de recall, depois do snapshot.

## Unidade ontológica

```go
type DiscoveryConceptID string

type DiscoveryConcept struct {
    ID           DiscoveryConceptID   `json:"id"`
    Triggers     []string             `json:"triggers"`
    Queries      []ConceptQuery       `json:"queries"`
    Related      []DiscoveryConceptID `json:"related,omitempty"`
    Forbidden    []string             `json:"forbidden,omitempty"`
    Rationale    string               `json:"rationale"`
}

type ConceptQuery struct {
    Text          string   `json:"text"`
    Discriminators []string `json:"discriminators"`
    NegativeTerms []string `json:"negative_terms,omitempty"`
    Language      string   `json:"language,omitempty"`
}

type DiscoveryOntology struct {
    SchemaVersion string             `json:"schema_version"`
    Version       string             `json:"version"`
    Concepts      []DiscoveryConcept `json:"concepts"`
    Hash          string             `json:"hash"`
}
```

`Triggers` reconhecem formulações da IdeaSpec. `Queries` são expressões públicas
do conceito, não concatenação da entrada. `Discriminators` identificam os termos
que não podem desaparecer no fallback.

## Conceitos mínimos da primeira versão

### `software-reuse-discovery`

Intenção: encontrar software existente antes de construir.

```text
software reuse discovery
existing software search
build versus reuse
```

### `repository-recommendation`

Intenção: transformar necessidade funcional em projetos candidatos.

```text
GitHub repository recommendation
natural language repository search
open source project recommendation
```

### `component-selection`

Intenção: comparar componentes para adoção ou integração.

```text
software component selection
open source component evaluation
dependency selection tool
```

### `adoption-decision`

Intenção: decidir se e como adotar uma ferramenta.

```text
dependency adoption decision
software adoption assessment
evidence backed tool evaluation
```

### `software-composition`

Intenção: cobrir uma necessidade com partes de vários projetos.

```text
software component composition
open source integration planning
capability coverage components
```

### `internal-reuse`

Intenção: descobrir capacidades já presentes no projeto ou portfólio.

```text
internal code reuse search
cross repository code discovery
developer knowledge reuse
```

### `technology-scouting`

Intenção: acompanhar alternativas e capacidades externas relevantes.

```text
open source technology scouting
developer tool discovery
software landscape analysis
```

As queries acima são candidatas pré-registradas. Sua permanência depende dos
testes offline de cobertura e contaminação, não do resultado live individual.

## Resolução de conceitos

```go
type ConceptMatch struct {
    ConceptID DiscoveryConceptID `json:"concept_id"`
    SourceID  string             `json:"source_id"`
    Trigger   string             `json:"trigger"`
    Strength  float64            `json:"strength"`
}

type OntologyExpansionResult struct {
    SchemaVersion string            `json:"schema_version"`
    OntologyVersion string          `json:"ontology_version"`
    OntologyHash  string            `json:"ontology_hash"`
    Matches       []ConceptMatch    `json:"matches"`
    Expansions    []SearchExpansion `json:"expansions"`
    Uncovered     []string          `json:"uncovered"`
}
```

Regras iniciais:

1. normalizar caixa, espaços e diacríticos;
2. reconhecer frase completa e variantes declaradas;
3. não usar stemmer probabilístico;
4. limitar a duas queries por conceito e origem;
5. deduplicar queries canonicamente;
6. preservar origem, trigger, conceito e hash da ontologia;
7. ordenar por `ConceptID`, query e `SourceID`;
8. registrar capacidades sem conceito em `Uncovered`.

Relações `Related` não expandem recursivamente na primeira versão. Elas servem
somente para diagnóstico e evolução futura.

## Qualidade de query antes da rede

```go
type QueryQualityVerdict string

const (
    QueryAccepted QueryQualityVerdict = "ACCEPTED"
    QueryWeak     QueryQualityVerdict = "WEAK"
    QueryRejected QueryQualityVerdict = "REJECTED"
)

type QueryQualityFinding struct {
    ExpansionID   string              `json:"expansion_id"`
    Query         string              `json:"query"`
    Verdict       QueryQualityVerdict `json:"verdict"`
    Discriminators []string           `json:"discriminators"`
    Missing       []string            `json:"missing,omitempty"`
    Reasons       []string            `json:"reasons"`
}
```

Uma query conceitual é `ACCEPTED` quando:

- possui entre 2 e 6 tokens relevantes;
- contém pelo menos um discriminator declarado;
- não contém nome ou owner proibido;
- não é apenas verbo genérico mais preposição;
- não repete literalmente uma sentença longa da IdeaSpec;
- possui rastreabilidade ontológica completa.

Queries `REJECTED` não acessam a rede. Queries `WEAK` podem entrar apenas no
subcorpus de controle e nunca como fallback automático.

## Fallback discriminatório

Substituir o algoritmo “primeiros dois termos” por:

```go
type QueryFallback struct {
    Original       string   `json:"original"`
    Effective      string   `json:"effective"`
    Preserved      []string `json:"preserved"`
    Removed        []string `json:"removed"`
    Discriminators []string `json:"discriminators"`
    Valid          bool     `json:"valid"`
}
```

Regras:

1. queries da ontologia já devem ser curtas e preferencialmente não precisam de
   fallback;
2. se necessário, preservar primeiro discriminators e substantivos declarados;
3. remover conectivos e qualificadores somente depois;
4. nunca produzir query sem discriminator;
5. nunca usar somente `decidir entre`, `encontrar projeto`, `comparar software`
   ou equivalentes genéricos;
6. fallback inválido é registrado e não enviado ao provedor.

O provedor continua responsável por adaptar sintaxe, como `lang:Go`, mas não por
decidir quais conceitos sobrevivem.

## Corpus offline obrigatório

Congelar ao menos doze IdeaSpecs em quatro grupos causais:

### Descoberta e reúso

- ideia do Bateia;
- busca de biblioteca de processamento de planilhas;
- escolha de fila de eventos local;

### Composição

- autenticação + autorização + auditoria;
- coleta + armazenamento + visualização;
- parser + grafo + análise de impacto;

### Reúso interno

- localizar implementação anterior de refresh token;
- encontrar padrão de retry usado no portfólio;
- consolidar dois componentes duplicados;

### Controles negativos

- jogo Katamari;
- autofagia biológica;
- peneiramento/mineração física.

Para cada case congelar:

- conceitos esperados;
- conceitos proibidos;
- discriminators obrigatórios;
- queries aceitáveis equivalentes;
- queries genéricas proibidas;
- grupo causal.

Casos derivados do Bateia e seus concorrentes permanecem no mesmo grupo para não
vazar entre treino e validação.

## Métricas offline

```go
type OntologyMetrics struct {
    ConceptPrecision      *float64 `json:"concept_precision"`
    ConceptRecall         *float64 `json:"concept_recall"`
    QueryAcceptanceRate   float64  `json:"query_acceptance_rate"`
    DiscriminatorRetention float64 `json:"discriminator_retention"`
    GenericFallbackRate   float64  `json:"generic_fallback_rate"`
    ForbiddenLeakRate     float64  `json:"forbidden_leak_rate"`
    ControlContamination  float64  `json:"control_contamination"`
    UncoveredRate         float64  `json:"uncovered_rate"`
}
```

Denominadores vazios permanecem `null`.

## CLI

```bash
bateia idea ontology expand \
  --input ideia.json \
  --ontology discovery-ontology.json \
  --output ontology-expansion.json

bateia idea ontology evaluate \
  --corpus ontology-corpus.json \
  --ontology discovery-ontology.json \
  --output ontology-evaluation.json

bateia idea ontology audit \
  --plan expansion-plan.json \
  --ontology discovery-ontology.json \
  --output query-quality.json
```

Os três comandos são offline. `idea discover` poderá consumir um plano já
auditado, mas continuará proibido de ler o registry de expected outcomes.

## Experimento live

Depois dos gates offline, repetir a Fase A do E-0012 com:

- mesma `bateia-idea.json`;
- ontologia e hash congelados;
- no máximo 14 queries ontológicas;
- `max-per-query=3`;
- concorrência 2;
- timeout total 3 minutos;
- attempt timeout 20 segundos;
- Sourcegraph público;
- registry fechado até o snapshot.

Executar dois braços:

```text
A — literal atual limitado pelo Passo 0013
B — ontologia lexical auditada
```

Não alterar screening. Comparar:

- completion rate;
- candidate yield;
- recall competitivo;
- contaminação nos controles;
- queries expiradas;
- candidatos em revisão imutável;
- origens por concorrente encontrado.

## Gate de decisão sobre o provedor

Após o braço B:

### Manter Sourcegraph na descoberta

Se:

- ao menos 3/4 sementes públicas forem encontradas;
- recall competitivo ≥ 70%;
- pelo menos metade das queries aceitas concluir no budget;
- controles negativos não avançarem;
- candidatos possuírem identidade resolvível.

### Separar catálogo e comprovação

Se os gates offline passarem, mas o Sourcegraph obtiver:

- recall competitivo < 40%; ou
- attempt timeout rate > 50%; ou
- predominância de matches incidentais em arquivos;
- maioria dos candidatos em `HEAD` sem identidade resolvida.

Nesse caso, o próximo passo deverá introduzir uma porta de catálogo baseada em
descrição, README, tópicos e metadados para descobrir projetos. Sourcegraph
continuará como provedor de comprovação em código, não será removido.

Resultados entre 40% e 70% permanecem `INCONCLUSIVE` e exigem ampliar corpus,
não ajustar thresholds retrospectivamente.

## Gates de aprovação

### Offline

- `concept_precision ≥ 90%`;
- `concept_recall ≥ 85%`;
- `discriminator_retention = 100%`;
- `generic_fallback_rate = 0%`;
- `forbidden_leak_rate = 0%`;
- `control_contamination = 0%`;
- `query_acceptance_rate ≥ 90%`;
- ordem da IdeaSpec e da ontologia não altera saída canônica;
- toda expansão aponta conceito, trigger, origem e hash.

### Live

- braço B supera o braço A em recall ou produz evidência suficiente para separar
  catálogo e comprovação;
- nenhum nome do registry aparece nas queries;
- nenhum fallback genérico é enviado;
- resultados parciais e timeouts continuam preservados pelo Passo 0013;
- nenhum candidato `HEAD` entra no benchmark de correspondência como sucesso.

O passo pode ser aprovado com decisão “Sourcegraph não serve para descoberta de
produto” se os gates offline passarem e o experimento live sustentar a separação.
Descobrir uma limitação do provedor é um resultado útil, não falha do passo.

## Testes obrigatórios

1. intenção de não começar do zero ativa `software-reuse-discovery`;
2. busca por projeto ativa `repository-recommendation`;
3. decisão de dependência ativa `adoption-decision`;
4. composição ativa `software-composition`;
5. busca no portfólio ativa `internal-reuse`;
6. controle Katamari não ativa composição de software;
7. autofagia biológica não ativa reúso interno;
8. trigger parcial dentro de outra palavra não casa;
9. normalização de diacríticos preserva resultado;
10. duas origens iguais deduplicam query e preservam ambas;
11. limite de duas queries por conceito é respeitado;
12. relação `Related` não expande recursivamente;
13. query aceita contém discriminator;
14. query longa literal é rejeitada;
15. query apenas verbal é rejeitada;
16. nome de concorrente em query é rejeitado;
17. fallback preserva `repository recommendation`;
18. fallback de `decidir entre ...` é rejeitado;
19. fallback nunca produz string vazia enviada à rede;
20. capacidades sem conceito aparecem em `Uncovered`;
21. hash muda quando trigger, query ou discriminator muda;
22. reordenação não altera hash nem expansões;
23. corpus mantém grupos causais juntos;
24. expected outcomes não são legíveis pelo expansor;
25. métricas vazias permanecem `null`;
26. braço live usa mesma política operacional do Passo 0013;
27. snapshot parcial continua reexecutável offline;
28. Sourcegraph não recebe query rejeitada;
29. nenhum token ou registry aparece em receipts;
30. suíte completa e linter V15 permanecem verdes.

## Critério de produto

Antes de acessar a rede, uma pessoa deve conseguir verificar:

- que conceitos da ideia foram reconhecidos;
- quais capacidades continuam sem linguagem de busca;
- por que cada query foi gerada;
- quais termos discriminatórios ela preserva;
- quais queries foram rejeitadas e por quê;
- que nenhum concorrente esperado foi usado para fabricar a busca.

Depois do run, devemos conseguir decidir com evidência se:

1. a linguagem de busca era inadequada;
2. o provedor não indexa bem descoberta por produto;
3. o budget operacional foi insuficiente;
4. a categoria não possui candidatos correspondentes;
5. precisamos de um catálogo diferente antes da comprovação em código.

## Critério de parada

Interromper se:

- uma query contiver nome ou owner do registry;
- vocabulário for alterado depois do snapshot;
- expected outcome orientar trigger ou discriminator;
- controle negativo ativar conceito de software sem justificativa;
- fallback perder todos os discriminators;
- recall live for calculado sobre candidato `HEAD` como correspondente;
- for necessário afrouxar screening para elevar recall;
- braço A e B usarem budgets ou políticas diferentes;
- o provider for trocado durante a comparação;
- `UNKNOWN` for contado como encontrado.

## Fora de escopo

- criar regras de evidência para capacidades internas do Bateia;
- preencher os oráculos de autofagia;
- executar Katamari;
- adicionar GitHub Search ou outro catálogo nesta implementação;
- embeddings ou busca vetorial;
- tradução automática;
- aprendizagem por feedback;
- ranking por popularidade;
- alterar timeout, screening ou decisões de reúso;
- executar código dos candidatos.

## Próximo passo previsto

Se a ontologia passar offline, mas Sourcegraph continuar com baixo recall e alto
timeout, o próximo passo separará formalmente:

```text
catálogo de projetos
    → descoberta por descrição, README, tópicos e metadados

Sourcegraph/conteúdo
    → comprovação de capacidades no repo@commit
```

Se Sourcegraph atingir o gate, o próximo passo congelará oráculos das capacidades
internas do Bateia e criará regras evidenciais correspondentes antes de retomar
autofagia e Katamari.
