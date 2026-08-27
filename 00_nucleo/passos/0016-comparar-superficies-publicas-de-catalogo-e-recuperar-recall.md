# Passo 0016 — Comparar superfícies públicas de catálogo e recuperar recall

**Estado**: `PRONTO PARA IMPLEMENTAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0015 — Separar catálogo de projetos e comprovação em código`](0015-separar-catalogo-de-projetos-e-comprovacao-em-codigo.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: torneio determinístico de catálogos, recuperação multiview e avaliação pública sem LLM

## Evidência que motiva este passo

O Passo 0015 corrigiu a arquitetura, mas não recuperou os produtos esperados:

| Observável | Braço C — GitHub repository search |
|---|---:|
| queries ontológicas | 12 |
| queries concluídas | 12/12 |
| candidatos únicos | 34 |
| identidades imutáveis | 34/34 |
| sementes repo-backed | 0/4 |
| recall competitivo | 0% |
| `ADVANCE` | 0 |
| `REVIEW` | 34 |

O primeiro run também mostrou por que não podemos apenas ampliar resultados: um
matcher permissivo por token promoveu 24 candidatos irrelevantes. Quando o
contrato passou a exigir o discriminator completo, nenhum candidato avançou.

Isso separa dois problemas:

1. identidade e operação do catálogo GitHub funcionam;
2. nome e descrição do repositório não expressam necessariamente a proposta do
   produto com o vocabulário usado na ideia.

Um projeto pode explicar sua intenção somente no README, documentação, tópico,
manifesto de pacote ou página de produto. Pesquisar apenas a visão “repositório”
é insuficiente, mas voltar a procurar ocorrências arbitrárias em código
reintroduziria o erro corrigido no Passo 0015.

## Incerteza que este passo reduz

Não sabemos qual superfície pública contém linguagem suficiente para descobrir
projetos por intenção funcional sem LLM:

- nome e descrição de repositório;
- README e documentação de apresentação;
- tópicos;
- registros de pacotes;
- catálogos alternativos de projetos.

Também não sabemos se a união controlada dessas superfícies aumenta recall sem
reintroduzir a contaminação que apareceu no Sourcegraph e no matcher por token.

## Hipótese

Uma recuperação multiview, em que cada superfície mantém identidade e peso
próprios, encontra projetos que a busca de repositórios não recupera:

```text
query ontológica auditada
        │
        ├── visão R: nome + descrição do repositório
        ├── visão D: README + documentação de apresentação
        ├── visão T: tópicos declarados
        ├── visão P: registros de pacotes
        └── visão X: outro catálogo público de projetos
                ↓
       consolidação por owner/repo
                ↓
       peneira multiview explicável
                ↓
       Sourcegraph restrito à prova em código
```

O ganho deve vir da concordância entre superfícies ou de um match forte em uma
superfície de produto, nunca de popularidade ou token genérico.

## Restrição central: superfície de produto, não ocorrência arbitrária

Uma superfície pode participar da descoberta somente se descreve o projeto
como unidade reutilizável.

São superfícies de catálogo válidas:

- nome e descrição do repositório;
- tópicos mantidos pelo projeto;
- README raiz;
- documentação inicial vinculada pelo repositório;
- nome, resumo e metadados de pacote;
- registro público cuja unidade seja projeto, produto ou pacote.

Não são superfícies de catálogo:

- qualquer arquivo de código que contém a palavra pesquisada;
- issue, comentário ou commit isolado;
- código vendored ou gerado;
- forks duplicados sem proposta própria;
- ranking baseado apenas em estrelas ou atividade;
- resultado web sem identidade resolvível.

Sourcegraph permanece fora da descoberta global. Ele só recebe candidatos
`ADVANCE` com `repository_scope`.

## Contrato de superfícies

```go
type CatalogSurfaceID string

const (
    SurfaceRepository CatalogSurfaceID = "REPOSITORY"
    SurfaceReadme     CatalogSurfaceID = "README"
    SurfaceTopics     CatalogSurfaceID = "TOPICS"
    SurfacePackage    CatalogSurfaceID = "PACKAGE_REGISTRY"
    SurfaceExternal   CatalogSurfaceID = "EXTERNAL_PROJECT_CATALOG"
)

type CatalogSurfaceCapabilities struct {
    ID                   CatalogSurfaceID `json:"id"`
    Provider             string           `json:"provider"`
    Unit                 string           `json:"unit"`
    SupportsExactPhrase  bool             `json:"supports_exact_phrase"`
    SupportsNegativeTerm bool             `json:"supports_negative_term"`
    ResolvesRepository   bool             `json:"resolves_repository"`
    ResolvesRevision     bool             `json:"resolves_revision"`
}

type SurfaceSearchRequest struct {
    ExpansionID    string           `json:"expansion_id"`
    Surface        CatalogSurfaceID `json:"surface"`
    Query          string           `json:"query"`
    Discriminators []string         `json:"discriminators"`
    NegativeTerms  []string         `json:"negative_terms,omitempty"`
    MaxResults     int              `json:"max_results"`
}
```

Cada adaptador declara sua superfície, unidade e capacidade de resolver
identidade antes da rede. Incompatibilidade falha antes da chamada externa.

## Resultado por superfície

```go
type SurfaceMatch struct {
    Surface       CatalogSurfaceID `json:"surface"`
    Provider      string           `json:"provider"`
    Repository    string           `json:"repository,omitempty"`
    Package       string           `json:"package,omitempty"`
    Field         string           `json:"field"`
    Snippet       string           `json:"snippet,omitempty"`
    Discriminators []string        `json:"discriminators"`
    Query         string           `json:"query"`
    ExpansionID   string           `json:"expansion_id"`
}
```

Regras:

1. o snippet é curto, normalizado e rastreável;
2. o match registra discriminator completo encontrado;
3. ausência de snippet ou campo observável produz `Unknown`, não match forte;
4. pacote sem repositório permanece separado até resolver seu `source`;
5. aliases convergem somente com evidência canônica;
6. toda identidade final recebe commit imutável;
7. cada superfície mantém receipt e rate-limit independentes.

## Superfícies iniciais

### R — GitHub repository

Preservar o braço C do Passo 0015 sem alteração como baseline. Não repetir ou
reinterpretar seus 34 candidatos.

### D — GitHub README/documentação

Pesquisar exclusivamente arquivos de apresentação:

```text
README
README.md
docs/index
docs/README
```

O adaptador pode usar busca de código do GitHub somente com restrição explícita
de caminho e tipo documental. Um match deve ser promovido para a unidade
`owner/repo`; o trecho do README explica a descoberta, mas não comprova
implementação.

Proibido remover a restrição documental como fallback.

### T — tópicos

Consultar tópicos como campos estruturados. Quando o provedor não suporta frase
exata em tópicos, decompor somente por termos declarados na ontologia e exigir
conjunção completa na resposta.

### P — registros de pacotes

Adicionar inicialmente os ecossistemas presentes na `preferred_stack` ou nas
restrições da IdeaSpec. Para uma ideia sem stack preferida, usar somente um
conjunto congelado e pequeno:

```text
npm
PyPI
crates.io
Go modules
```

O pacote só vira candidato de repositório depois que sua URL de source é
resolvida e normalizada. Downloads não contam como relevância.

### X — catálogo externo

Selecionar no máximo um catálogo público adicional antes do experimento. A
seleção deve ser baseada em capacidades documentadas:

- pesquisa projetos ou pacotes, não ocorrências arbitrárias;
- acesso público ou credencial já disponível;
- identidade e termos de uso observáveis;
- resposta congelável;
- rate-limit conhecido;
- vínculo resolvível com repositório.

A escolha e sua justificativa devem ser congeladas antes de abrir o registry.

## Planejador de consultas por superfície

A mesma query nem sempre possui sintaxe válida em todas as superfícies. O
planejador pode adaptar sintaxe, mas não semântica:

```go
type SurfaceQueryPlan struct {
    ExpansionID           string           `json:"expansion_id"`
    Surface               CatalogSurfaceID `json:"surface"`
    Original              string           `json:"original"`
    Effective             string           `json:"effective"`
    PreservedDiscriminators []string       `json:"preserved_discriminators"`
    Removed               []string         `json:"removed,omitempty"`
    Valid                 bool             `json:"valid"`
    Reason                string           `json:"reason"`
}
```

O plano é válido somente quando preserva ao menos um discriminator completo.
Fallback inválido não acessa rede.

## Consolidação multiview

```go
type MultiViewCandidate struct {
    Repository RepositoryCandidate `json:"repository"`
    Matches    []SurfaceMatch       `json:"matches"`
    Surfaces   []CatalogSurfaceID   `json:"surfaces"`
    Concepts   []DiscoveryConceptID `json:"concepts"`
}
```

Regras de consolidação:

1. chave primária `owner/repo` normalizada;
2. pacote converge somente por URL canônica de source;
3. fork não converge automaticamente com upstream;
4. queries repetidas na mesma superfície contam uma vez;
5. superfícies diferentes permanecem evidências independentes;
6. conceito é derivado da ontologia, não inferido do candidato;
7. ordem de respostas não altera o resultado canônico.

## Peneira multiview

### `ADVANCE`

Um candidato avança quando possui revisão imutável e satisfaz ao menos uma das
condições:

- discriminator completo em duas superfícies independentes;
- dois conceitos independentes em uma superfície de produto forte;
- match exato em tópico estruturado e confirmação no README;
- pacote com resumo discriminatório e README do repositório correspondente.

### `REVIEW`

- uma única superfície de produto com match completo;
- identidade resolvida, mas snippet indisponível;
- pacote relevante cuja origem não pôde ser resolvida;
- budget ou rate-limit impediu confirmação cruzada.

### `REJECT`

- somente token parcial ou genérico;
- apenas popularidade;
- controle negativo explícito;
- material educacional, lista ou dataset sem componente reutilizável;
- match apenas fora de README/documentação inicial;
- identidade divergente entre superfícies.

Nenhum `REVIEW` é enviado ao Sourcegraph automaticamente.

## Corpus offline

Congelar ao menos vinte casos de superfície:

- cinco READMEs de produtos com intenção explícita;
- três READMEs com termos incidentais;
- três tópicos estruturados positivos;
- três tópicos com apenas token parcial;
- três pacotes com source resolvível;
- três pacotes sem source ou com source divergente;
- controles Katamari, autofagia e peneiramento físico distribuídos entre as
  superfícies.

O corpus precisa conter casos em português e inglês, hífen, caixa e diacríticos.
Casos derivados das sementes competitivas permanecem no split de validação e
não podem alimentar regras.

## CLI

```bash
bateia idea catalog plan-surfaces \
  --plan ontology-plan.json \
  --surfaces repository,readme,topics,package,external \
  --output surface-query-plan.json

bateia idea catalog discover-surfaces \
  --plan surface-query-plan.json \
  --max-per-query 3 \
  --concurrency 2 \
  --timeout 3m \
  --output multiview-snapshot.json

bateia idea catalog screen-multiview \
  --snapshot multiview-snapshot.json \
  --output multiview-screening.json \
  --proof-plan code-proof-plan.json

bateia idea catalog evaluate-multiview \
  --snapshot multiview-snapshot.json \
  --judgments multiview-judgments.json \
  --output multiview-evaluation.json
```

Planejamento, screening e avaliação são offline. Somente
`discover-surfaces` acessa rede.

## Experimento público

Manter imutáveis:

- IdeaSpec E-0012;
- ontologia e hash do Passo 0014;
- registry competitivo;
- máximo de três resultados por query e superfície;
- concorrência 2;
- timeout total de três minutos por braço;
- attempt timeout de vinte segundos;
- braços A, B e C já capturados.

Executar incrementalmente:

```text
D1 — README/documentação
D2 — tópicos estruturados
D3 — registros de pacote aplicáveis
D4 — catálogo externo escolhido
U  — união canônica de R + D1 + D2 + D3 + D4
```

Parar uma superfície quando rate-limit ou acesso impedir observação confiável;
registrar `Unknown`, não zero resultados.

Abrir o registry somente depois de congelar todos os snapshots. Medir recall
individual e ganho marginal de cada superfície:

```text
marginal_recall(surface) = recall(union_with_surface) - recall(union_without_surface)
```

Somente candidatos `ADVANCE` da união seguem para prova restrita no
Sourcegraph.

## Métricas

- recall competitivo por superfície;
- recall da união;
- precisão julgada por superfície;
- ganho marginal de recall;
- overlap entre superfícies;
- taxa de resolução para `owner/repo`;
- taxa de revisão imutável;
- discriminator retention;
- taxa de matches por token parcial;
- contaminação dos controles;
- rate-limit e timeout por provedor;
- `ADVANCE`, `REVIEW` e `REJECT` por superfície;
- custo de consultas por candidato que avança;
- conversão catálogo → prova;
- escape de repositório na prova, obrigatoriamente zero.

Denominadores vazios permanecem `null`.

## Gates de aprovação

### Offline

- discriminator retention = 100%;
- token parcial nunca vira match forte;
- controles negativos nunca recebem `ADVANCE`;
- pacote sem source nunca recebe identidade fabricada;
- ordem das superfícies e respostas não altera a união;
- `Unknown` não é convertido em ausência ou sucesso;
- toda query efetiva aponta query original, conceito e superfície;
- mutation score do contrato = 1,0.

### Live

- união encontra ao menos 3/4 sementes repo-backed;
- recall competitivo da união ≥ 70%;
- precisão julgada dos dez primeiros `ADVANCE` ≥ 60%;
- pelo menos uma superfície acrescenta recall marginal positivo;
- ≥ 80% dos candidatos `ADVANCE` possuem commit imutável;
- nenhum controle negativo recebe `ADVANCE`;
- nenhum nome ou owner do registry aparece nas queries;
- pelo menos metade das consultas planejadas conclui no budget;
- Sourcegraph recebe somente repositórios `ADVANCE` e sempre com scope.

### Decisão

- se README recuperar recall com precisão aceitável, incorporá-lo como segunda
  visão permanente do GitHub;
- se um registro de pacotes contribuir apenas para bibliotecas, habilitá-lo por
  stack, não globalmente;
- se o catálogo externo dominar o ganho, mantê-lo atrás da porta
  `ProjectCatalog` sem contaminar L1;
- se a união permanecer abaixo de 40% de recall, concluir que a ontologia
  lexical isolada não expressa intenção suficiente para descoberta pública e
  escrever um passo de reformulação de queries por decomposição funcional,
  ainda sem LLM;
- resultados entre 40% e 70% permanecem `INCONCLUSIVE` e exigem ampliar corpus,
  não baixar gates.

## Testes obrigatórios

1. cada superfície declara unidade e capacidades;
2. superfície incompatível falha antes da rede;
3. README search mantém restrição de caminho em todo fallback;
4. ocorrência fora de README é rejeitada como catálogo;
5. tópico exige discriminator completo;
6. token parcial em tópico não conta;
7. package source resolve owner/repo canônico;
8. pacote sem source permanece `Unknown`;
9. source divergente permanece `REVIEW`;
10. catálogo externo sem identidade não avança;
11. query efetiva preserva discriminator;
12. fallback inválido não toca provedor;
13. duas superfícies consolidam um candidato;
14. duas queries na mesma superfície não simulam duas superfícies;
15. aliases convergem somente com URL canônica;
16. fork não converge silenciosamente com upstream;
17. reordenação produz snapshot idêntico;
18. match em duas superfícies pode avançar;
19. uma superfície isolada permanece `REVIEW`;
20. popularidade isolada não avança;
21. material educacional é rejeitado;
22. listas `awesome-*` não avançam sem componente reutilizável;
23. controles Katamari, autofagia e peneiramento não avançam;
24. commit ausente impede sucesso definitivo;
25. rate-limit produz `Unknown` e receipt parcial;
26. timeout preserva resultados concluídos;
27. registry não é legível por planejadores ou adaptadores;
28. queries não contêm nomes ou owners esperados;
29. braço C permanece byte a byte inalterado;
30. Sourcegraph recebe somente `ADVANCE` com `repository_scope`;
31. resultado fora do scope é descartado;
32. replay offline reproduz métricas;
33. cada novo consumidor possui prompt próprio;
34. suíte completa, linter V15 e `git diff --check` permanecem verdes.

## Fora de escopo

- usar LLM, embeddings ou reranking neural;
- adaptar queries depois de abrir o registry;
- aceitar ocorrência arbitrária em código como catálogo;
- promover `REVIEW` para aumentar recall reportado;
- usar estrelas, downloads ou atividade como substitutos de relevância;
- consultar todas as linguagens e registros sem relação com a IdeaSpec;
- executar Sourcegraph globalmente;
- alterar decisões de reuso, Autofagia ou Katamari;
- esconder rate-limit como resultado vazio.

## Critério de produto

Ao final, o Bateia deve explicar não apenas **qual projeto encontrou**, mas
**onde a proposta desse projeto foi encontrada**:

```text
ideia
→ conceito ontológico
→ superfície de produto
→ discriminator observado
→ owner/repo@commit
→ concordância entre superfícies
→ decisão de avançar ou revisar
→ prova restrita em código
```

O passo será útil mesmo com resultado negativo se identificar, com ganho
marginal mensurável, quais superfícies não ajudam. O que não é aceitável é
voltar a confundir volume de ocorrências com descoberta de produtos.
