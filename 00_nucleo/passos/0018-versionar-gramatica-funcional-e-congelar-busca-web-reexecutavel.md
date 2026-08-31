# Passo 0018 — Versionar gramática funcional e congelar busca web reexecutável

**Estado**: `EM VALIDAÇÃO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0017 — Gerar consultas por frames funcionais e descobrir produtos na web`](0017-gerar-consultas-por-frames-funcionais-e-descobrir-produtos-na-web.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: gramática determinística, porta de busca observável e experimento causal sem LLM no pipeline

## Evidência que motiva este passo

O Passo 0017 confirmou que descrever o trabalho do produto é melhor que emitir
somente categorias nominais, mas não passou o gate de produto:

| Métrica | Nominal | Funcional | Gate |
|---|---:|---:|---:|
| consultas concluídas | 12/12 | 12/12 | ≥ 50% |
| recall competitivo de produto | 0/4 | 1/4 | ganho ≥ 30 p.p. |
| recall competitivo repo-backed | 0/4 | 0/4 | ≥ 70% |
| precisão julgada em 10 | `null` | `null` | ≥ 60% |

O braço funcional encontrou GitMVP e pistas novas, como IdeaProof e Reposeek,
porém o ganho observado ficou em 25 pontos percentuais. Duas causas ficaram
visíveis:

1. o template concatena slots sem flexão, produzindo frases como
   `tool that evaluate`;
2. o provedor usado no experimento agregou resultados por lote e não preservou
   rank, snippet e resposta bruta por consulta.

O primeiro problema reduz a qualidade da consulta. O segundo impede replay,
deduplicação comprovável e cálculo honesto de precisão.

## Incerteza que este passo reduz

Não sabemos se o baixo recall veio principalmente da formulação gramatical ou
da baixa observabilidade do índice externo.

O passo deve separar essas causas:

- comparar a gramática congelada do Passo 0017 com uma gramática nova;
- manter frames, orçamento e provedor constantes;
- preservar a resposta individual de cada consulta antes de qualquer juízo;
- repetir deduplicação e métricas inteiramente offline.

## Hipótese

Uma gramática pequena, flexionada por regras explícitas e especializada pelo
tipo de trabalho produzirá consultas mais naturais sem inferência probabilística.
Se o provedor preservar cada resultado bruto, será possível atribuir o ganho à
formulação, e não a uma inspeção manual posterior.

```text
frame congelado
   ↓
seleção determinística da família de intenção
   ↓
realização morfológica por tabela versionada
   ↓
auditoria lexical e estrutural
   ↓
uma requisição observável por query
   ↓
snapshot bruto imutável
   ↓
normalização e julgamento offline
```

## Proibição de aprendizado pelas respostas esperadas

Replen, ARIS, Starboard, GitMVP, OSS Compass e outros resultados conhecidos
podem existir somente no registry aberto **depois** dos snapshots.

É proibido:

- extrair palavras, slogans ou descrições desses produtos para o léxico;
- criar uma regra porque ela recupera uma semente específica;
- escolher tempos verbais depois de observar o ranking;
- adicionar domínio, nome, owner ou trecho exclusivo de concorrente;
- remover uma consulta por ter produzido resultado ruim;
- usar LLM, embeddings, tradução neural ou reranker neural no Bateia;
- atribuir ao Bateia a semântica interna do índice externo.

As novas regras devem nascer exclusivamente das classes de erro do Passo 0017:
concordância verbal, intenção, ausência de entrada, saída pouco natural,
ambiguidade nominal e perda de timing.

## Gramática funcional versionada

```go
type FunctionalGrammar struct {
    SchemaVersion string                   `json:"schema_version"`
    Version       string                   `json:"version"`
    Rules         []FunctionalGrammarRule  `json:"rules"`
    Morphology    []MorphologyRule         `json:"morphology"`
    Hash          string                   `json:"hash"`
}

type FunctionalGrammarRule struct {
    ID            string   `json:"id"`
    Intent        string   `json:"intent"`
    Pattern       string   `json:"pattern"`
    RequiredSlots []string `json:"required_slots"`
    OptionalSlots []string `json:"optional_slots,omitempty"`
    Priority      int      `json:"priority"`
    MaxVariants   int      `json:"max_variants"`
}

type MorphologyRule struct {
    ID          string `json:"id"`
    SubjectKind string `json:"subject_kind"`
    Source      string `json:"source"`
    Realized    string `json:"realized"`
}
```

Não haverá conjugador geral. A morfologia será uma tabela finita, auditável e
com falha fechada: verbo sem realização conhecida gera `Unknown`.

## Famílias de intenção

Cada frame escolhe uma família por `concept_id`, sem classificador aprendido:

| Família | Conceitos iniciais | Forma funcional |
|---|---|---|
| descoberta | software-reuse-discovery, internal-reuse | `find {object} for {input} {timing}` |
| recomendação | repository-recommendation | `recommend {object} for {input} and return {output}` |
| avaliação | adoption-decision | `evaluate {object} from {input} and produce {output}` |
| seleção | component-selection | `compare {object} for {input} and select {output}` |
| composição | software-composition | `combine {object} to produce {output} for {input}` |

Uma família pode ter no máximo três padrões congelados. Sinônimos continuam
vindos do léxico versionado; a gramática apenas realiza relações entre slots.

## Auditoria pré-rede

```go
type GrammarFinding struct {
    QueryID          string   `json:"query_id"`
    GrammarVersion   string   `json:"grammar_version"`
    RuleID           string   `json:"rule_id"`
    Verdict          string   `json:"verdict"`
    UnresolvedVerbs  []string `json:"unresolved_verbs,omitempty"`
    MissingSlots     []string `json:"missing_slots,omitempty"`
    ForbiddenHits    []string `json:"forbidden_hits,omitempty"`
    Reasons          []string `json:"reasons"`
}
```

Uma consulta somente acessa a rede quando:

- todos os verbos foram realizados por regra identificável;
- ação e objeto permanecem observáveis;
- entrada, saída ou timing preserva a finalidade;
- possui entre cinco e dezoito tokens relevantes;
- não contém esperado, domínio ou owner proibido;
- não repete um n-grama longo da IdeaSpec;
- regra, léxico, frame e hashes estão registrados.

## Porta de busca observável

```go
type ObservableWebSearch interface {
    Capabilities(context.Context) WebSearchCapabilities
    SearchOne(context.Context, ObservableWebRequest) (RawWebResponse, error)
}

type ObservableWebRequest struct {
    RequestID  string `json:"request_id"`
    QueryID    string `json:"query_id"`
    Query      string `json:"query"`
    MaxResults int    `json:"max_results"`
    Locale     string `json:"locale"`
}

type RawWebResponse struct {
    RequestID    string          `json:"request_id"`
    QueryID      string          `json:"query_id"`
    Provider     string          `json:"provider"`
    RetrievedAt  string          `json:"retrieved_at"`
    Status       string          `json:"status"`
    HTTPStatus   int             `json:"http_status,omitempty"`
    Results      []RawWebResult  `json:"results"`
    Receipt      WebReceipt      `json:"receipt"`
}
```

Cada chamada corresponde a exatamente uma query. O adaptador deve preservar:

- rank original;
- URL retornada e URL canônica quando disponível;
- título e snippet sem reescrita;
- indicador de anúncio, quando oferecido;
- instante, status, timeout e rate-limit;
- hash da requisição e da resposta persistida;
- identificação do provedor e versão do adaptador;
- receipt sanitizado, nunca token ou cabeçalho secreto.

Se o contrato do provedor proibir armazenamento de conteúdo, o snapshot registra
somente campos permitidos e declara `NON_REPLAYABLE`. Esse resultado não pode
aprovar o gate causal.

## Snapshot e replay offline

```go
type WebSearchSnapshot struct {
    SchemaVersion string            `json:"schema_version"`
    RunID         string            `json:"run_id"`
    Arm           string            `json:"arm"`
    GrammarHash   string            `json:"grammar_hash"`
    QuerySetHash  string            `json:"query_set_hash"`
    Provider      string            `json:"provider"`
    Responses     []RawWebResponse  `json:"responses"`
    SnapshotHash  string            `json:"snapshot_hash"`
}
```

O replay recebe snapshots, nunca rede. Ele reproduz:

1. canonicalização de URL;
2. deduplicação dentro e entre queries;
3. preservação da primeira posição observada;
4. separação de anúncio;
5. união dos trabalhos correspondentes;
6. métricas com os mesmos denominadores;
7. lista de campos ausentes e `Unknown`.

O hash do snapshot exclui somente o próprio campo `snapshot_hash`. Alterar
query, rank, URL, título, snippet, status ou receipt muda o hash.

## Experimento causal

Congelar três braços antes da rede:

```text
N  — doze rótulos nominais do Passo 0017
F1 — doze queries funcionais exatamente como executadas no Passo 0017
F2 — doze queries da gramática versionada deste passo
```

Controles:

- mesmos seis frames para F1 e F2;
- doze queries por braço;
- cinco resultados por query;
- mesma localidade e provedor;
- chamadas intercaladas N/F1/F2 para reduzir efeito temporal;
- timeout por chamada e total idênticos;
- registry inacessível até os três snapshots serem selados;
- julgamento cego por URL canônica, sem indicação do braço;
- produto e repo-backed medidos separadamente.

F1 deve ser regenerado a partir de seus artefatos congelados e comparado ao hash
do Passo 0017. Divergência encerra o experimento antes da rede.

## Classificação dos erros

Cada resultado rejeitado recebe uma causa primária:

```text
GRAMMAR_AGREEMENT
MISSING_JOB_CONTEXT
ACADEMIC_HOMONYM
NON_SOFTWARE_HOMONYM
CATEGORY_ONLY
AGGREGATOR
PRODUCT_WITHOUT_FUNCTIONAL_MATCH
AMBIGUOUS_IDENTITY
PROVIDER_TRUNCATION
MISSING_OBSERVABILITY
```

Essa taxonomia pode orientar a próxima revisão, mas não autoriza alteração do
léxico durante a execução.

## CLI

```bash
bateia idea functional realize \
  --expansion functional-expansion.json \
  --grammar functional-grammar-v2.json \
  --output functional-queries-v2.json

bateia idea functional grammar-audit \
  --queries functional-queries-v2.json \
  --grammar functional-grammar-v2.json \
  --output grammar-findings.json

bateia idea web discover-observable \
  --queries functional-queries-v2.json \
  --provider web \
  --max-per-query 5 \
  --output web-search-snapshot.json

bateia idea web replay \
  --snapshot web-search-snapshot.json \
  --output web-search-replay.json

bateia idea web compare \
  --nominal nominal-snapshot.json \
  --functional-v1 functional-v1-snapshot.json \
  --functional-v2 functional-v2-snapshot.json \
  --judgments blinded-judgments.json \
  --output grammar-experiment.json
```

`realize`, `grammar-audit`, `replay` e `compare` são puros em L1. Acesso à rede
existe apenas no adaptador L3 chamado pela composição em L4.

## Gates de aprovação

### Gramática

- 100% dos verbos das queries aceitas possuem regra morfológica identificada;
- zero ocorrência de `tool that evaluate`, `tool that find` ou equivalentes
  inválidos definidos no corpus;
- cobertura dos slots obrigatórios = 100%;
- forbidden leak = 0%;
- queries aceitas ≥ 90%;
- reordenação de frames, regras e léxico não muda a saída;
- mudança em padrão ou morfologia muda hashes e IDs;
- controles negativos não geram consultas de software;
- mutation score = 1,0.

### Observabilidade e replay

- 100% das respostas preservam `request_id`, `query_id`, status e instante;
- ≥ 95% dos hits orgânicos preservam rank, URL, título e snippet;
- resposta de uma query nunca é atribuída a outra;
- snapshots não contêm tokens ou cabeçalhos secretos;
- duas repetições do replay geram bytes idênticos;
- adulterar qualquer campo material invalida o hash;
- timeout e rate-limit permanecem `Unknown`, nunca conjunto vazio;
- `NON_REPLAYABLE` não pode aprovar o experimento causal.

### Produto

- `judged_precision_at_10` de F2 ≥ 60%;
- F2 encontra ao menos 3/4 sementes repo-backed;
- recall repo-backed de F2 ≥ 70%;
- ganho de recall de F2 sobre N ≥ 30 pontos percentuais;
- F2 não piora precisão em mais de 5 pontos contra F1;
- nenhum controle negativo recebe `ADVANCE`;
- produto sem repositório continua `PRODUCT_ONLY`;
- Sourcegraph recebe apenas identidade `owner/repo@commit` resolvida.

### Decisão

- se F2 superar N e F1, promover a gramática v2;
- se F2 superar F1 mas não passar recall, manter `EM VALIDAÇÃO` e ampliar o
  corpus de jobs, não o registry;
- se F1 e F2 forem equivalentes, morfologia não explica a perda;
- se todos os braços falharem igualmente, investigar índice e identidade;
- se o snapshot não for reexecutável, declarar `INCONCLUSIVE`;
- se F2 piorar, preservar F1 e registrar a mutação regressiva.

## Testes obrigatórios

1. `evaluate` vira `evaluates` somente quando o sujeito exige terceira pessoa;
2. forma imperativa preserva `evaluate`;
3. ação composta mantém coordenação verbal válida;
4. verbo desconhecido produz `Unknown`;
5. regra morfológica aplicada fica rastreável;
6. descoberta preserva entrada e timing;
7. recomendação preserva saída;
8. avaliação preserva veredito;
9. seleção preserva comparação e escolha;
10. composição preserva múltiplos projetos e cobertura;
11. reuso interno não vira descoberta pública;
12. nomes e owners proibidos bloqueiam a query;
13. registry não pode ser lido antes dos snapshots;
14. F1 reproduz o hash congelado do Passo 0017;
15. cada request contém uma única query;
16. respostas concorrentes não trocam `query_id`;
17. rank original é preservado;
18. URL canônica deduplica sem fundir domínios diferentes;
19. anúncio permanece separado;
20. snippet ausente gera `Unknown` observável;
21. timeout preserva snapshot parcial;
22. rate-limit não vira ausência de resultados;
23. receipt remove credenciais;
24. replay não acessa rede;
25. replay repetido é byte a byte idêntico;
26. alteração de resultado invalida o hash;
27. julgamento cego não recebe nome do braço;
28. métricas separam produto de repo-backed;
29. denominador vazio permanece `null`;
30. erro acadêmico recebe causa própria;
31. homônimo não software recebe causa própria;
32. produto sem identidade permanece `PRODUCT_ONLY`;
33. Sourcegraph recusa produto sem `owner/repo@commit`;
34. cada novo consumidor Go possui prompt próprio;
35. suíte completa, build, linter V15 e `git diff --check` permanecem verdes.

## Artefatos esperados

```text
03_infra/lexical/functional-grammar-v2.json
00_nucleo/benchmark/0018/grammar-corpus.json
00_nucleo/benchmark/0018/queries-n.json
00_nucleo/benchmark/0018/queries-f1.json
00_nucleo/benchmark/0018/queries-f2.json
00_nucleo/benchmark/0018/snapshot-n.json
00_nucleo/benchmark/0018/snapshot-f1.json
00_nucleo/benchmark/0018/snapshot-f2.json
00_nucleo/benchmark/0018/blinded-judgments.json
00_nucleo/benchmark/0018/grammar-experiment.json
00_nucleo/laudos/0018-versionar-gramatica-funcional-e-congelar-busca-web-reexecutavel.md
```

## Fora de escopo

- usar LLM, embeddings ou reranking neural no pipeline;
- implementar conjugação linguística geral;
- adaptar regras às páginas dos concorrentes;
- substituir o índice externo por crawler indiscriminado;
- armazenar conteúdo proibido pelo contrato do provedor;
- resolver identidade por semelhança de nome apenas;
- promover produto sem repositório para prova em código;
- alterar os snapshots do Passo 0017;
- afirmar determinismo do ranking externo;
- medir qualidade sem julgamento congelado.

## Critério de produto

O Bateia terá avançado quando conseguir transformar o mesmo frame em consultas
gramaticalmente válidas, demonstrar que elas aumentam recall com julgamento
cego e reproduzir toda normalização e métrica a partir de snapshots, sem rede e
sem LLM.
