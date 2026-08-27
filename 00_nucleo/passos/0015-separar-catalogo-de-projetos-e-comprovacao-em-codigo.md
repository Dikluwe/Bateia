# Passo 0015 — Separar catálogo de projetos e comprovação em código

**Estado**: `PRONTO PARA IMPLEMENTAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0014 — Construir ontologia lexical para descoberta de software`](0014-construir-ontologia-lexical-para-descoberta-de-software.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: arquitetura de provedores, descoberta por catálogo e experimento público sem LLM

## Evidência que motiva este passo

O Bateia inicialmente descobria repositórios com o GitHub CLI. Depois encontrou
no próprio processo a ideia de usar Sourcegraph e essa capacidade foi
implementada. A integração era válida, mas o Sourcegraph acabou ocupando o papel
de provedor principal de descoberta, embora sua especialidade seja localizar
ocorrências e estruturas dentro de código.

O Passo 0014 isolou o efeito dessa troca:

| Observável | GitHub/expectativa de catálogo | Sourcegraph live |
|---|---:|---:|
| unidade procurada | projeto/repositório | ocorrência em código |
| queries ontológicas | 12 | 12 |
| queries concluídas | — | 12/12 |
| candidatos únicos | — | 31 |
| sementes encontradas | — | 0/4 |
| identidade imutável | esperada | 0/31 |
| recall competitivo | — | 0% |

As queries deixaram de ser o principal limitador: passaram em todos os gates
offline e o braço público completou. Ainda assim, os resultados eram arquivos
que continham termos incidentais como `reuse`, `component` e `recommendation`,
não produtos funcionalmente relacionados.

O problema agora é arquitetural: usamos um provedor de comprovação como se
fosse um catálogo.

## Incerteza que este passo reduz

Não sabemos se restaurar GitHub como catálogo, alimentado pela ontologia do
Passo 0014, recupera projetos funcionalmente próximos com precisão suficiente
para tornar o restante do funil útil.

Também não sabemos se separar explicitamente os papéis evita que candidatos
incidentais em código avancem para screening e evidência.

## Hipótese

Uma descoberta em duas fases melhora recall e rastreabilidade:

```text
IdeaSpec
   ↓
ontologia lexical auditada
   ↓
catálogo GitHub: descrição + README + tópicos + metadados
   ↓
peneira de candidatos como produtos
   ↓
Sourcegraph: comprovação de capacidades no código
   ↓
evidência + decisão de reúso + Autofagia/Katamari
```

GitHub deve responder **quais projetos podem atender à ideia**. Sourcegraph deve
responder **onde e como uma capacidade alegada aparece no código**.

## Restrição central: papel antes da tecnologia

Todo provedor deve declarar seu papel semântico antes de ser conectado ao
pipeline. Um provedor não pode substituir outro apenas porque ambos oferecem
uma operação chamada “search”.

```go
type DiscoveryProviderRole string

const (
    RoleProjectCatalog DiscoveryProviderRole = "PROJECT_CATALOG"
    RoleCodeProof      DiscoveryProviderRole = "CODE_PROOF"
)

type ProviderCapabilities struct {
    Provider       string                  `json:"provider"`
    Roles          []DiscoveryProviderRole `json:"roles"`
    SearchUnits    []string                `json:"search_units"`
    IdentityLevels []string                `json:"identity_levels"`
}
```

Regras:

1. `PROJECT_CATALOG` retorna repositórios identificáveis e metadados de produto;
2. `CODE_PROOF` retorna ocorrências vinculadas a um repositório já candidato;
3. descoberta inicial aceita somente provedores com `PROJECT_CATALOG`;
4. comprovação aceita somente provedores com `CODE_PROOF`;
5. um provedor pode declarar ambos os papéis, mas deve ser avaliado separadamente
   em cada um;
6. incompatibilidade de papel falha antes de acessar a rede;
7. o receipt registra papel, unidade pesquisada e provedor efetivo.

## Porta de catálogo

Introduzir uma porta independente da implementação:

```go
type CatalogSearchRequest struct {
    ExpansionID   string   `json:"expansion_id"`
    Query         string   `json:"query"`
    Discriminators []string `json:"discriminators"`
    MaxResults    int      `json:"max_results"`
}

type CatalogCandidate struct {
    Repository     RepositoryCandidate `json:"repository"`
    Description    string              `json:"description,omitempty"`
    Topics         []string            `json:"topics,omitempty"`
    ReadmeSnippet  string              `json:"readme_snippet,omitempty"`
    License        string              `json:"license,omitempty"`
    Archived       bool                `json:"archived"`
    DefaultBranch  string              `json:"default_branch,omitempty"`
    UpdatedAt      string              `json:"updated_at,omitempty"`
    Stars          int                 `json:"stars,omitempty"`
    MatchedFields  []string            `json:"matched_fields"`
}

type ProjectCatalog interface {
    SearchCatalog(context.Context, CatalogSearchRequest) (CatalogDiscovery, error)
}
```

Estrelas e atualização são contexto, nunca prova de correspondência funcional.
O candidato precisa indicar em quais campos ocorreu o match.

## Adaptador GitHub

Restaurar descoberta pública por GitHub usando o mecanismo disponível e mais
reproduzível no ambiente:

1. preferir a API de busca via cliente já existente quando ela expuser todos os
   campos necessários;
2. admitir `gh search repos` como adaptador substituível;
3. não executar shell dentro de L1;
4. limitar paginação e concorrência pela mesma política operacional do Passo
   0013;
5. preservar resposta parcial e rate-limit em receipt;
6. normalizar `owner/repo`, URL e branch padrão;
7. resolver a revisão imutável antes de o candidato contar como recuperado;
8. registrar quando README não pôde ser consultado;
9. nunca usar o registry competitivo para gerar a busca.

### Estratégia de busca

Cada query ontológica aceita produz no máximo duas consultas de catálogo:

```text
A — nome, descrição e tópicos
B — README, somente se A não atingir o budget de candidatos distintos
```

O fallback deve preservar os discriminadores do Passo 0014. Consultas como
`decidir entre`, `internal code` ou `GitHub repository` permanecem proibidas.

## Peneira específica de catálogo

Antes do Sourcegraph, classificar cada candidato como:

```go
type CatalogVerdict string

const (
    CatalogAdvance CatalogVerdict = "ADVANCE"
    CatalogReject  CatalogVerdict = "REJECT"
    CatalogReview  CatalogVerdict = "REVIEW"
)
```

Sinais positivos explicáveis:

- discriminator em nome, descrição ou tópico;
- duas ou mais queries conceituais independentes apontam o mesmo repositório;
- README declara capacidade compatível com a IdeaSpec;
- repositório possui identidade e revisão resolvidas;
- licença e estado de arquivamento são observáveis.

Sinais negativos:

- match existe somente em texto incidental do README;
- descrição pertence a jogo, biologia, mineração ou outro controle;
- repositório é material educacional sem produto reutilizável;
- identidade não pode ser resolvida;
- candidato foi recuperado apenas por termo genérico;
- repositório está arquivado e não há razão explícita para referência.

O score não deve usar LLM e precisa expor contribuição de cada sinal.

## Sourcegraph reposicionado

Sourcegraph deixa de receber a IdeaSpec aberta como busca global. Para cada
candidato `ADVANCE`, recebe:

- `repo:^owner/repo$` obrigatório;
- revisão imutável quando suportada;
- uma capacidade ou discriminator por consulta;
- budget de ocorrências por capacidade;
- caminhos preferenciais derivados do plano de evidência;
- exclusões de vendored, generated e fixtures quando aplicável.

Nenhum repositório novo descoberto incidentalmente pelo Sourcegraph entra no
funil. Uma ocorrência fora do candidato é registrada como observação, não como
novo candidato.

## Contrato de proveniência

Cada candidato deve carregar a sequência completa:

```text
IdeaSpec field
→ ontology concept
→ catalog query
→ matched catalog field
→ owner/repo
→ immutable revision
→ catalog verdict
→ code-proof query
→ evidence item
```

Perder qualquer elo impede decisão afirmativa de `FORK`, `EXTEND`, `COMPOSE` ou
`INTEGRATE`.

## CLI

```bash
bateia idea catalog discover \
  --plan ontology-plan.json \
  --provider github \
  --max-per-query 3 \
  --concurrency 2 \
  --timeout 3m \
  --output catalog-snapshot.json

bateia idea catalog screen \
  --snapshot catalog-snapshot.json \
  --model catalog-screening.json \
  --output catalog-screening-report.json

bateia idea prove \
  --catalog catalog-screening-report.json \
  --provider sourcegraph \
  --output code-proof-snapshot.json
```

Os comandos de avaliação continuam offline:

```bash
bateia idea catalog evaluate \
  --snapshot catalog-snapshot.json \
  --judgments catalog-judgments.json \
  --output catalog-evaluation.json
```

## Experimento público

Reutilizar a mesma IdeaSpec, ontologia, hash, registry e budgets do Passo 0014.
Executar três braços sem alterar queries após observar resultados:

```text
A — baseline lexical/Sourcegraph preservado do Passo 0013
B — ontologia/Sourcegraph preservado do Passo 0014
C — ontologia/GitHub catálogo deste passo
```

Somente depois do snapshot C abrir o registry e calcular recall.

Para candidatos recuperados no braço C:

1. congelar revisão;
2. executar a peneira de catálogo;
3. enviar somente `ADVANCE` ao Sourcegraph com filtro de repositório;
4. medir quantas capacidades alegadas recebem prova em código;
5. manter os controles negativos no mesmo corpus.

## Métricas

### Catálogo

- `catalog_completion_rate`;
- `catalog_candidate_yield`;
- `catalog_unique_repository_rate`;
- `catalog_identity_resolution_rate`;
- `catalog_competitive_recall`;
- `catalog_judged_precision`;
- `catalog_control_contamination`;
- `catalog_multi_query_support_rate`;
- `catalog_rate_limit_rate`.

### Comprovação

- `proof_attempt_rate`: candidatos `ADVANCE` efetivamente consultados;
- `proof_capability_hit_rate`;
- `proof_incidental_match_rate`;
- `proof_repository_escape_rate`, obrigatoriamente zero;
- `proof_immutable_revision_rate`;
- `catalog_to_proof_conversion_rate`.

Denominadores vazios permanecem `null`.

## Gates de aprovação

### Arquitetura

- descoberta inicial rejeita provedor sem `PROJECT_CATALOG`;
- comprovação rejeita provedor sem `CODE_PROOF`;
- Sourcegraph recebe sempre filtro explícito de repositório na comprovação;
- achados fora do repositório não viram candidatos;
- GitHub e Sourcegraph implementam portas distintas;
- receipts registram papel e unidade pesquisada;
- L1 não conhece `gh`, HTTP, GraphQL ou detalhes de Sourcegraph.

### Catálogo live

- pelo menos 3/4 sementes repo-backed encontradas;
- recall competitivo ≥ 70%;
- ao menos 80% dos candidatos possuem identidade resolvida;
- ao menos metade das queries conclui dentro do budget;
- nenhum controle negativo recebe `ADVANCE`;
- precisão julgada dos dez primeiros candidatos ≥ 60%;
- nenhum nome do registry aparece nas queries.

### Comprovação live

- 100% das consultas possuem filtro `repo:`;
- `proof_repository_escape_rate = 0`;
- pelo menos uma capacidade é comprovada para cada semente recuperada e
  funcionalmente correspondente;
- nenhum candidato sem revisão imutável conta como sucesso;
- resultado parcial permanece reexecutável offline.

Se o GitHub atingir recall mas baixa precisão, o próximo passo deve calibrar a
peneira de catálogo. Se também falhar em recall abaixo de 40%, o próximo passo
deve comparar outro catálogo público, não rebaixar retrospectivamente os gates.

## Testes obrigatórios

1. GitHub declara `PROJECT_CATALOG`;
2. Sourcegraph declara `CODE_PROOF`;
3. Sourcegraph é recusado como catálogo;
4. GitHub é recusado como prova se não declarar esse papel;
5. incompatibilidade falha antes da primeira chamada externa;
6. query de catálogo preserva discriminator;
7. fallback genérico não acessa rede;
8. owner/repo é normalizado canonicamente;
9. duplicatas entre queries são consolidadas;
10. origens de todas as queries são preservadas;
11. nome, descrição, tópico e README são distinguíveis nos matches;
12. estrelas não bastam para `ADVANCE`;
13. match exclusivo em README incidental não basta para `ADVANCE`;
14. dois conceitos independentes aumentam suporte sem duplicar candidato;
15. controles Katamari, autofagia e peneiramento não avançam;
16. repositório arquivado é sinalizado;
17. revisão ausente impede sucesso definitivo;
18. rate-limit produz receipt parcial, não resultado vazio;
19. cancelamento preserva candidatos concluídos;
20. Sourcegraph recebe `repo:^owner/repo$`;
21. Sourcegraph recebe a revisão congelada quando suportada;
22. match externo não cria novo candidato;
23. cada prova referencia candidato, capacidade, query e origem de catálogo;
24. registry não é dependência do expansor, catálogo ou peneira;
25. nenhum token aparece em snapshot, erro ou receipt;
26. replay offline reproduz métricas;
27. braços A e B permanecem imutáveis;
28. braço C usa a mesma ontologia e os mesmos budgets;
29. cada consumidor novo possui prompt próprio;
30. suíte completa, linter V15 e `git diff --check` permanecem verdes.

## Fora de escopo

- usar LLM, embeddings ou reranking neural;
- adaptar ontologia depois de observar o registry;
- substituir GitHub por scraping HTML não versionado;
- alterar regras de decisão de reúso;
- formar Katamari antes da comprovação funcional;
- tratar popularidade como relevância;
- enviar todos os resultados do catálogo ao Sourcegraph sem peneira;
- remover Sourcegraph do Bateia.

## Critério de produto

Ao final, uma pessoa deve conseguir fornecer uma ideia e receber:

1. projetos candidatos descobertos como produtos, com explicação do match;
2. rejeições de catálogo justificadas antes de gastar prova em código;
3. capacidades comprovadas dentro dos repositórios que avançaram;
4. uma trilha reproduzível que separa claramente “parece atender” de “há
   evidência de que atende”.

O produto começa a ser útil quando encontra um ponto de partida real sem
confundir ocorrência lexical em código com projeto reutilizável.
