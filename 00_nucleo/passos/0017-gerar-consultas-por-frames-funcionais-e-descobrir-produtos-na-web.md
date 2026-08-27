# Passo 0017 — Gerar consultas por frames funcionais e descobrir produtos na web

**Estado**: `PRONTO PARA IMPLEMENTAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0016 — Comparar superfícies públicas de catálogo e recuperar recall`](0016-comparar-superficies-publicas-de-catalogo-e-recuperar-recall.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: decomposição funcional determinística, busca web de produtos e avaliação sem LLM no pipeline

## Evidência que motiva este passo

Os Passos 0014 a 0016 melhoraram a arquitetura, a identidade dos candidatos e a
segurança da peneira, mas não recuperaram os concorrentes conhecidos:

| Braço | Superfície | Recall competitivo |
|---|---|---:|
| A | Sourcegraph + expansão lexical | 0/4 |
| B | Sourcegraph + ontologia | 0/4 |
| C | GitHub nome/descrição/README | 0/4 |
| D | README documental + tópicos | 0/4 |

O braço D produziu 35 candidatos e a união chegou a 69, mas nenhum apareceu em
duas superfícies independentes. O problema não foi falta de resultados: foi
falta de correspondência funcional.

A ontologia produz rótulos corretos, porém abstratos:

```text
software reuse discovery
repository recommendation
component selection
adoption decision
software composition
internal reuse
```

Uma reconstrução controlada da busca originalmente usada para encontrar os
exemplos mostrou outro padrão. Consultas que descrevem o trabalho do produto
recuperaram imediatamente projetos relevantes:

```text
tool recommends GitHub repositories from natural language project idea
find existing open source projects before building software
AI tool discover similar GitHub repositories based on requirements
open source software adoption decision tool
```

Essas consultas recuperaram Replen, ARIS, Starboard, GitMVP e OSS Compass, além
de novos análogos como Reposeek, RepoWise, RepoInspo e Fixseek.

Os artefatos antigos não preservaram as consultas que originaram o registry.
Portanto essa reconstrução é uma **hipótese reproduzida**, não provenance
retroativa comprovada. O novo passo deve corrigir justamente essa lacuna:
consultas, resultados, seleção e instante precisam ser registrados desde a
origem.

## Incerteza que este passo reduz

Não sabemos se consultas determinísticas compostas como frases funcionais
recuperam produtos relevantes melhor que rótulos ontológicos.

Também não sabemos quanto do ganho observado vem:

- da formulação comportamental;
- da cobertura ampla de um índice web;
- do ranking semântico e opaco do mecanismo externo;
- da atualidade das páginas indexadas.

O passo precisa separar esses efeitos e nunca atribuir ao Bateia uma capacidade
que pertence ao ranking externo.

## Hipótese

Uma consulta construída a partir de quatro elementos observáveis recupera
produtos que executam o trabalho desejado:

```text
entrada + ação + objeto + saída
```

Exemplo:

```text
entrada: natural language project idea
ação: find and rank
objeto: existing GitHub repositories
saída: shortlist before building from scratch
```

Materialização:

```text
find and rank existing GitHub repositories from a natural language project idea
return a shortlist before building from scratch
```

O pipeline previsto é:

```text
IdeaSpec
   ↓
ontologia lexical
   ↓
frames funcionais determinísticos
   ↓
consultas comportamentais auditadas
   ↓
índice web / catálogo de páginas de produto
   ↓
resolução para produto e owner/repo
   ↓
peneira multiview
   ↓
Sourcegraph com repository_scope
```

## Restrição central: descrever trabalho, não respostas

É proibido:

- inserir nome, owner, domínio ou frase exclusiva de uma semente;
- derivar templates a partir dos textos das sementes;
- abrir o registry antes de congelar frames, templates e queries;
- ajustar uma query porque ela recuperou ou perdeu um concorrente específico;
- usar LLM, embeddings, tradução automática ou reranker neural dentro do
  pipeline do Bateia;
- fingir que o ranking web externo é determinístico ou explicável;
- converter página sem identidade em repositório;
- tratar resultado patrocinado, agregador ou lista como produto sem evidência;
- usar popularidade como substituto de correspondência funcional;
- promover `Unknown` para aumentar recall;
- remover restrições documentais do Sourcegraph ou da prova.

## Frame funcional

```go
type FunctionalFrameID string

type FunctionalFrame struct {
    ID           FunctionalFrameID `json:"id"`
    ConceptID    DiscoveryConceptID `json:"concept_id"`
    SourceIDs    []string           `json:"source_ids"`
    Actor        []string           `json:"actor,omitempty"`
    Inputs       []string           `json:"inputs"`
    Actions      []string           `json:"actions"`
    Objects      []string           `json:"objects"`
    Outputs      []string           `json:"outputs"`
    Timing       []string           `json:"timing,omitempty"`
    Constraints  []string           `json:"constraints,omitempty"`
    NegativeJobs []string           `json:"negative_jobs,omitempty"`
    RuleID       string             `json:"rule_id"`
}
```

Cada elemento vem de regra declarada e versionada. O frame nunca é completado
por inferência probabilística.

### Exemplos iniciais

#### Descoberta antes da construção

```text
concept: software-reuse-discovery
input: software idea | product requirement
action: find | search | check
object: existing open-source projects | existing software
output: reusable starting point | existing alternative
timing: before building from scratch
```

#### Recomendação de repositórios

```text
concept: repository-recommendation
input: natural-language project idea | software requirements
action: recommend | find and rank
object: GitHub repositories | open-source projects
output: ranked matches | repository shortlist
```

#### Decisão de adoção

```text
concept: adoption-decision
input: GitHub repository | open-source dependency
action: evaluate | compare | assess
object: software tool | dependency
output: adoption decision | adopt trial hold avoid verdict
```

#### Seleção de componentes

```text
concept: component-selection
input: software requirement | required capability
action: compare | select
object: open-source tools | software components
output: best-fit component | evidence-backed comparison
```

#### Composição

```text
concept: software-composition
input: capabilities not covered by one project
action: combine | compose
object: multiple open-source projects | reusable components
output: minimal capability coverage | integration plan
```

#### Reuso interno

```text
concept: internal-reuse
input: implementation need | code capability
action: find | locate
object: existing implementation across repositories
output: reusable internal component | previous solution
```

## Ontologia de verbos e papéis

Congelar uma tabela pequena e explícita:

```go
type FunctionalLexicon struct {
    SchemaVersion string                 `json:"schema_version"`
    Version       string                 `json:"version"`
    Actions       map[string][]string    `json:"actions"`
    Objects       map[string][]string    `json:"objects"`
    Outputs       map[string][]string    `json:"outputs"`
    Inputs        map[string][]string    `json:"inputs"`
    Timing        map[string][]string    `json:"timing"`
    Forbidden     []string               `json:"forbidden"`
    Hash          string                 `json:"hash"`
}
```

O léxico não contém produtos conhecidos. Alteração muda versão e hash.

## Templates determinísticos

```go
type FunctionalQueryTemplate struct {
    ID             string   `json:"id"`
    Pattern        string   `json:"pattern"`
    RequiredSlots  []string `json:"required_slots"`
    MaxVariants    int      `json:"max_variants"`
    Surface        string   `json:"surface"`
}
```

Templates mínimos:

```text
T1: {action} {object} from {input}
T2: tool that {action} {object} and returns {output}
T3: {action} {object} {timing}
T4: software that takes {input} and produces {output}
T5: open-source tool for {action} {object}
```

Regras:

1. no máximo três queries por frame;
2. no máximo duas variantes por slot;
3. cada query contém ação, objeto e entrada ou saída;
4. timing é opcional, mas preservado quando distingue a intenção;
5. termos genéricos isolados são rejeitados;
6. a mesma combinação canônica é emitida uma vez;
7. queries têm entre cinco e dezesseis tokens relevantes;
8. toda palavra vem do frame, template ou léxico versionado;
9. nenhuma query copia sentença longa da IdeaSpec;
10. ordenação é estável por conceito, frame e template.

## Contrato da consulta funcional

```go
type FunctionalQuery struct {
    ID             string            `json:"id"`
    FrameID        FunctionalFrameID `json:"frame_id"`
    ConceptID      DiscoveryConceptID `json:"concept_id"`
    TemplateID     string            `json:"template_id"`
    Text           string            `json:"text"`
    Inputs         []string          `json:"inputs"`
    Actions        []string          `json:"actions"`
    Objects        []string          `json:"objects"`
    Outputs        []string          `json:"outputs"`
    Timing         []string          `json:"timing,omitempty"`
    SourceIDs      []string          `json:"source_ids"`
    LexiconHash    string            `json:"lexicon_hash"`
}
```

Antes da rede, uma pessoa deve conseguir responder:

- qual trabalho a query descreve;
- quais slots foram usados;
- de qual campo da IdeaSpec e conceito ela veio;
- qual template gerou a frase;
- quais variantes foram descartadas;
- se algum termo proibido apareceu.

## Auditoria de qualidade

```go
type FunctionalQueryFinding struct {
    QueryID       string   `json:"query_id"`
    Verdict       string   `json:"verdict"`
    MissingSlots  []string `json:"missing_slots,omitempty"`
    GenericTerms  []string `json:"generic_terms,omitempty"`
    ForbiddenHits []string `json:"forbidden_hits,omitempty"`
    Reasons       []string `json:"reasons"`
}
```

Uma query é `ACCEPTED` somente quando:

- contém ação e objeto;
- contém ao menos uma entrada ou saída;
- preserva o conceito e origens;
- não contém nome esperado;
- não é categoria nominal isolada;
- não excede os budgets;
- possui léxico e template congelados;
- descreve comportamento verificável.

`REJECTED` nunca acessa rede. Construção insuficiente fica `UNKNOWN`, não é
completada por palpite.

## Provedor de índice web

```go
type WebProductSearch interface {
    Capabilities() WebSearchCapabilities
    SearchProducts(context.Context, WebProductRequest) (WebProductDiscovery, error)
}

type WebProductRequest struct {
    QueryID    string `json:"query_id"`
    Query      string `json:"query"`
    MaxResults int    `json:"max_results"`
    Locale     string `json:"locale,omitempty"`
}
```

O provedor pesquisa páginas de produto e documentação pública. Seu ranking é
externo e potencialmente semântico; o Bateia não afirma reproduzir o ranking.
Para compensar:

- congelar query exata e posição;
- registrar URL, título, snippet e domínio;
- registrar instante e provedor;
- preservar resultados brutos permitidos pelo contrato;
- resolver redirects antes da identidade;
- classificar anúncios separadamente;
- registrar timeout, bloqueio e rate-limit;
- permitir replay offline do snapshot;
- nunca expor token em receipt;
- não fazer scraping quando termos ou robots impedirem.

## Resultado web

```go
type WebProductHit struct {
    Rank        int      `json:"rank"`
    URL         string   `json:"url"`
    CanonicalURL string  `json:"canonical_url,omitempty"`
    Domain      string   `json:"domain"`
    Title       string   `json:"title"`
    Snippet     string   `json:"snippet"`
    QueryID     string   `json:"query_id"`
    MatchedJobs []string `json:"matched_jobs"`
    Sponsored   bool     `json:"sponsored"`
}
```

Um hit web é somente uma pista de produto. Ele não é automaticamente candidato
de código.

## Resolução produto → repositório

```go
type ProductIdentity struct {
    ProductURL     string              `json:"product_url"`
    ProductDomain  string              `json:"product_domain"`
    Repository     *RepositoryCandidate `json:"repository,omitempty"`
    ResolutionPath []string            `json:"resolution_path"`
    Confidence     string              `json:"confidence"`
    Status         string              `json:"status"`
    Reason         string              `json:"reason"`
}
```

Resolução permitida:

1. link explícito da página para GitHub/GitLab;
2. metadado `sameAs`, `repository` ou `source`;
3. pacote associado com source canônico;
4. owner/repo já presente na URL do hit;
5. busca secundária exata pelo nome e domínio, somente para confirmação.

Sem vínculo observável, o produto permanece `PRODUCT_ONLY`. Produtos sem
repositório podem contar como concorrentes de produto, mas ficam fora do
denominador repo-backed e não seguem ao Sourcegraph.

## Peneira de produto

### `ADVANCE_REPOSITORY`

- comportamento funcional aparece no título ou snippet;
- página descreve produto ou projeto reutilizável;
- identidade resolve para owner/repo@commit;
- não é anúncio, lista, artigo acadêmico ou agregador;
- ao menos duas partes do frame são observáveis na página.

### `ADVANCE_PRODUCT_ONLY`

- correspondência funcional forte;
- produto existe e possui URL canônica;
- repositório não é público ou não pôde ser resolvido;
- status fica fora da prova em código.

### `REVIEW`

- correspondência funcional parcial;
- página opaca ou snippet insuficiente;
- identidade ambígua;
- ranking retornou resultado sem conteúdo verificável.

### `REJECT`

- termo incidental;
- conteúdo educacional, artigo ou lista sem produto;
- controle negativo;
- anúncio sem correspondência orgânica;
- página copia palavras da query sem descrever o trabalho;
- domínio ou redirect suspeito.

## Corpus offline

Congelar no mínimo 24 IdeaSpecs:

- seis descoberta/recomendação de projetos;
- quatro decisão de adoção;
- quatro seleção de componentes;
- três composição de capacidades;
- três reuso interno;
- quatro controles negativos.

Cada caso declara:

- conceitos esperados;
- frames esperados;
- slots obrigatórios;
- trabalhos proibidos;
- famílias de query aceitáveis;
- queries nominais que devem ser rejeitadas;
- grupo causal.

O caso Bateia permanece num grupo separado. Textos de Replen, ARIS, Starboard,
GitMVP, OSS Compass e novos análogos não entram no corpus de construção.

## Experimento causal das queries

Antes do experimento competitivo, comparar dois braços no mesmo provedor web:

```text
N — rótulos nominais da ontologia
F — frames funcionais
```

Mesmos limites:

- no máximo doze queries por braço;
- máximo cinco resultados por query;
- mesma localidade;
- mesma janela de execução;
- timeout total de três minutos;
- sem registry durante a busca;
- deduplicação por URL canônica;
- snapshots congelados antes do julgamento.

Isso separa o efeito da formulação do efeito do índice web.

## Experimento competitivo

Depois de congelar N e F:

1. abrir o registry;
2. resolver produtos e repositórios;
3. calcular recall produto e repo-backed separadamente;
4. julgar os dez primeiros hits únicos de cada braço;
5. executar Sourcegraph apenas nos `ADVANCE_REPOSITORY`;
6. preservar produtos sem repositório como `PRODUCT_ONLY`;
7. comparar ganho de F sobre N.

Métricas:

```text
functional_query_acceptance
frame_slot_coverage
nominal_query_rejection
web_completion_rate
unique_product_yield
product_identity_resolution_rate
repo_identity_resolution_rate
product_competitive_recall
repo_backed_competitive_recall
judged_precision_at_10
functional_over_nominal_recall_gain
control_contamination
sponsored_result_rate
unknown_rate
proof_conversion_rate
```

Denominadores vazios permanecem `null`.

## CLI

```bash
bateia idea functional frame \
  --input ideia.json \
  --ontology discovery-ontology.json \
  --lexicon functional-lexicon.json \
  --output functional-frames.json

bateia idea functional expand \
  --frames functional-frames.json \
  --templates functional-query-templates.json \
  --output functional-queries.json

bateia idea functional audit \
  --queries functional-queries.json \
  --lexicon functional-lexicon.json \
  --output functional-query-quality.json

bateia idea web discover \
  --queries functional-queries.json \
  --provider web \
  --max-per-query 5 \
  --timeout 3m \
  --output web-product-snapshot.json

bateia idea web resolve \
  --snapshot web-product-snapshot.json \
  --output product-identities.json

bateia idea web evaluate \
  --snapshot web-product-snapshot.json \
  --judgments web-product-judgments.json \
  --output web-product-evaluation.json
```

`frame`, `expand`, `audit`, `resolve` e `evaluate` são reexecutáveis offline.
Somente `web discover` acessa o índice externo.

## Gates de aprovação

### Frames e queries offline

- concept-to-frame recall ≥ 90%;
- slot coverage = 100% nos slots obrigatórios;
- ação e objeto presentes em 100% das queries aceitas;
- entrada ou saída presente em 100%;
- query acceptance ≥ 90%;
- rótulos nominais isolados rejeitados = 100%;
- forbidden leak = 0%;
- controles negativos não geram trabalho de software;
- ordem da IdeaSpec, léxico e templates não altera saída;
- hash muda quando regra, slot ou template muda;
- mutation score = 1,0.

### Web

- ao menos metade das queries conclui no budget;
- todo hit preserva query, rank, URL, título e snippet ou fica `Unknown`;
- redirects são resolvidos e registrados;
- anúncios nunca contam como hit orgânico;
- nenhum token aparece em snapshots;
- replay offline reproduz deduplicação e métricas;
- produtos sem repositório não são descartados nem enviados à prova.

### Produto

- braço F encontra ao menos 3/4 sementes repo-backed;
- recall repo-backed de F ≥ 70%;
- recall de produto inclui corretamente sementes sem repositório verificável;
- `judged_precision_at_10 ≥ 60%`;
- ganho de recall de F sobre N ≥ 30 pontos percentuais;
- nenhum controle negativo recebe `ADVANCE`;
- nenhuma query contém nome, owner ou domínio do registry;
- Sourcegraph recebe somente `owner/repo@commit` resolvido e scoped.

### Decisão

- se F superar N e passar gates, frames funcionais tornam-se a estratégia
  principal de descoberta de produtos;
- se F superar N mas ficar entre 40% e 70%, ampliar corpus e frames sem alterar
  expected outcomes;
- se N e F obtiverem resultados equivalentes, o ganho veio do índice web e não
  da decomposição; manter o provedor como experimento, não como solução causal;
- se F ficar abaixo de 40%, revisar a decomposição da IdeaSpec em jobs e
  restrições, não afrouxar a peneira;
- se o provedor web for opaco ou não congelável, resultado é `INCONCLUSIVE`, não
  autorização para scraping.

## Testes obrigatórios

1. `software-reuse-discovery` produz frame de busca antes da construção;
2. `repository-recommendation` inclui ideia/requisito como entrada;
3. decisão de adoção inclui veredito como saída;
4. seleção inclui comparação e componente;
5. composição inclui múltiplos projetos e cobertura;
6. reuso interno inclui implementação existente;
7. frame preserva conceito, regra e source IDs;
8. frame insuficiente permanece `Unknown`;
9. controle Katamari não gera composição de software;
10. autofagia biológica não gera reuso interno;
11. peneiramento físico não gera seleção de componentes;
12. template exige ação;
13. template exige objeto;
14. template exige entrada ou saída;
15. no máximo três queries por frame;
16. variantes de slot respeitam budget;
17. query nominal isolada é rejeitada;
18. query comportamental válida é aceita;
19. frase longa da IdeaSpec não é copiada;
20. nomes e owners proibidos são rejeitados;
21. reordenação não altera queries;
22. mudança no léxico altera hash;
23. mudança no template altera ID;
24. provider ausente falha antes da rede;
25. output ausente bloqueia a rede;
26. timeout produz snapshot parcial;
27. rate-limit permanece `Unknown`;
28. anúncio é marcado separadamente;
29. URL canônica deduplica resultados;
30. domínios diferentes não convergem por título;
31. GitHub URL resolve owner/repo;
32. página com link explícito resolve repositório;
33. produto sem repo permanece `PRODUCT_ONLY`;
34. identidade ambígua permanece `REVIEW`;
35. artigo e lista não avançam como produto;
36. duas execuções offline reproduzem métricas;
37. braço nominal e funcional usam mesmo provedor e budgets;
38. registry não é legível antes dos snapshots;
39. nenhum expected outcome aparece nas queries;
40. Sourcegraph recebe somente repositórios resolvidos e scoped;
41. cada novo consumidor possui prompt próprio;
42. suíte completa, linter V15 e `git diff --check` permanecem verdes.

## Fora de escopo

- usar LLM ou embeddings no pipeline;
- reproduzir internamente o ranking semântico do mecanismo web;
- depender de um único provedor sem porta substituível;
- adaptar templates às sementes depois do experimento;
- usar conteúdo de concorrentes para construir o léxico;
- promover produtos sem identidade para prova em código;
- substituir evidência funcional por SEO;
- aceitar anúncio como relevância;
- raspar páginas contra termos ou robots;
- alterar Autofagia, Katamari ou políticas de decisão;
- afirmar provenance histórica que não foi registrada.

## Critério de produto

O Bateia deve passar de:

```text
procurar a categoria “software reuse discovery”
```

para:

```text
procurar uma ferramenta que recebe uma ideia de software,
encontra projetos existentes, compara sua adequação
e devolve um ponto de partida antes de construir do zero
```

Cada palavra dessa formulação precisa ter origem declarada. O índice web pode
ajudar a encontrar produtos, mas a explicação do trabalho procurado continua
pertencendo ao Bateia.
