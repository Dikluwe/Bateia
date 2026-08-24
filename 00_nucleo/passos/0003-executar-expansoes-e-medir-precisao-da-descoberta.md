# Passo 0003 — Executar expansões e medir precisão da descoberta

**Estado**: `IMPLEMENTADO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0002 — Expandir vocabulário com baseline mensurável`](0002-expandir-vocabulario-com-baseline-mensuravel.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: segundo corte vertical; descoberta rastreável e benchmark com julgamento

## Intenção

Executar cada expansão de uma ideia como uma busca independente, reunir e
deduplicar os repositórios encontrados e medir se o vocabulário lexical melhora a
**qualidade dos candidatos**, não apenas a aparência das queries.

O passo deve responder:

> As expansões lexicais encontram mais projetos relevantes ou apenas produzem mais
> resultados e mais chamadas ao provedor?

Ainda não queremos extrair arquivos, calcular similaridade final ou recomendar
reuso. Este passo termina na descoberta julgada.

## Por que este passo vem agora

O Passo 0002 mediu 100% de cobertura conceitual no corpus controlado e nenhuma
contaminação. Esse resultado prova o contrato, mas pode superestimar o valor real:
o vocabulário foi escrito junto do corpus.

A próxima prova precisa incluir o mundo externo:

```text
IdeaSpec
    ↓
ExpansionPlan literal ou lexical
    ↓
uma chamada rastreável por expansão
    ↓
candidatos deduplicados, ainda sem triagem profunda
    ↓
julgamento versionado
    ↓
comparação literal × lexical
```

Sem essa separação, um resultado ruim pode ser atribuído incorretamente à ideia,
à expansão, ao fallback do Sourcegraph ou ao ranking.

## Incertezas reduzidas

- quais expansões realmente produzem candidatos;
- quais expansões produzem somente ruído;
- quanto as expansões se sobrepõem nos mesmos repositórios;
- se o baseline lexical melhora precisão ou apenas volume;
- quais consultas efetivas o provedor executou após seus fallbacks;
- quanto da avaliação permanece inconclusivo por falta de julgamento.

## Entrega

### 1. Contrato de descoberta rastreável em L1

Criar tipos puros equivalentes a:

```go
type DiscoveryRequest struct {
    ExpansionID string `json:"expansion_id"`
    Query       string `json:"query"`
    Language    string `json:"language,omitempty"`
    MaxResults int    `json:"max_results"`
}

type ProviderAttempt struct {
    Query      string `json:"query"`
    ResultCount int   `json:"result_count"`
    Error      string `json:"error,omitempty"`
}

type ProviderDiscovery struct {
    Provider   string                `json:"provider"`
    Attempts   []ProviderAttempt     `json:"attempts"`
    Candidates []RepositoryCandidate `json:"candidates"`
}

type DiscoveredCandidate struct {
    Repository   RepositoryCandidate `json:"repository"`
    ExpansionIDs []string            `json:"expansion_ids"`
    Queries      []string            `json:"queries"`
}

type ExpansionDiscovery struct {
    Expansion SearchExpansion   `json:"expansion"`
    Provider  ProviderDiscovery `json:"provider"`
}

type DiscoverySnapshot struct {
    SchemaVersion string               `json:"schema_version"`
    Idea          IdeaSpec             `json:"idea"`
    Strategy      string               `json:"strategy"`
    Language      string               `json:"language,omitempty"`
    PerExpansion  []ExpansionDiscovery `json:"per_expansion"`
    Candidates    []DiscoveredCandidate `json:"candidates"`
    Failures      []ExecutionFailure   `json:"failures,omitempty"`
}
```

Os nomes podem mudar, mas o snapshot precisa preservar:

- a ideia-alvo completa;
- a estratégia que produziu o plano;
- cada expansão e sua proveniência semântica;
- cada consulta efetivamente tentada pelo provedor;
- resultados por expansão antes da deduplicação;
- todas as expansões que encontraram cada repositório deduplicado;
- falhas parciais sem apagar resultados válidos.

Datas operacionais podem existir nos candidatos, mas não devem ser usadas como
identidade nem impedir uma forma canônica para golden tests.

### 2. Porta de descoberta observável em L1

O `SearchProvider` atual retorna somente candidatos e oculta os fallbacks internos.
Não quebrar esse contrato apenas para o experimento. Criar uma porta complementar:

```go
type ObservableDiscoveryProvider interface {
    Discover(ctx context.Context, request DiscoveryRequest) (ProviderDiscovery, error)
}
```

O caso de uso de descoberta do plano deve depender dessa porta e:

- validar plano, linguagem, limite e IDs;
- executar exatamente uma unidade lógica por expansão;
- limitar concorrência, com valor padrão pequeno;
- respeitar cancelamento de contexto;
- continuar após falhas parciais e registrá-las;
- retornar erro se todas as expansões falharem;
- deduplicar por identidade canônica `owner/repo`, sem diferenciar caixa;
- mesclar IDs e queries de origem na ordem do plano;
- produzir ordem determinística independentemente da conclusão concorrente;
- não executar triagem, manifesto, árvore, conteúdo ou score.

O relógio, a rede, o provedor concreto e o filesystem permanecem fora de L1.

### 3. Observabilidade real do Sourcegraph em L3

Adaptar o cliente Sourcegraph à nova porta preservando o método existente. A nova
operação deve expor cada tentativa já realizada pelo adaptador, incluindo:

```text
lang:<linguagem> <query completa>
lang:<linguagem> <termos curtos>
<termos curtos sem linguagem>
```

Requisitos:

- registrar query efetiva, quantidade retornada e erro resumido por tentativa;
- parar no primeiro resultado não vazio, mantendo as tentativas anteriores;
- não incluir token, cabeçalhos ou corpo GraphQL sensível no snapshot;
- distinguir resposta vazia de falha de transporte/API;
- preservar limite total de candidatos por expansão;
- usar o contexto recebido e os timeouts já existentes;
- manter compatibilidade com `SearchProvider` e seus testes.

Essa mudança é parte da prova: sem as tentativas efetivas não é possível saber se
o resultado veio da expansão planejada ou de um fallback reduzido e genérico.

### 4. Comando explícito de descoberta em L2/L4

Adicionar:

```bash
bateia idea discover \
  --input ideia.json \
  --strategy lexical \
  --provider sourcegraph \
  --language go \
  --max-per-query 10 \
  --concurrency 2 \
  --output snapshot.json
```

Requisitos:

- `--strategy literal|lexical`;
- `--provider sourcegraph` explícito neste passo;
- `--max-per-query` entre 1 e 25;
- `--concurrency` entre 1 e 4;
- `--output` obrigatório para evitar perder uma captura custosa;
- erro claro para estratégia, provedor ou limites inválidos;
- confirmação no `stderr` apenas depois de o arquivo ser gravado;
- JSON do snapshot nunca misturado com logs;
- `idea plan`, `idea expand` e `idea evaluate` continuam offline e não constroem o
  provedor remoto;
- somente `idea discover` inicializa Sourcegraph.

O comando deve funcionar sem token quando o endpoint público permitir, mas não
deve prometer que rede ou cota estarão disponíveis.

### 5. Captura reproduzível e modo replay

Resultados remotos mudam. A avaliação não pode consultar a rede toda vez que os
testes rodam. O snapshot gravado pelo comando é a fronteira reproduzível.

Adicionar inspeção offline:

```bash
bateia idea discovery evaluate \
  --snapshot snapshot.json \
  --judgments judgments.json
```

Esse comando não inicializa rede. Testes e métricas usam fixtures congeladas,
enquanto uma captura real e datada fica registrada no laudo.

O snapshot deve possuir `schema_version`. Mudança incompatível exige nova versão,
não interpretação silenciosa.

### 6. Julgamentos explícitos

Criar um formato versionado que não confunda ausência de julgamento com
irrelevância:

```go
type RelevanceLabel string

const (
    Relevant   RelevanceLabel = "relevant"
    Irrelevant RelevanceLabel = "irrelevant"
    Uncertain  RelevanceLabel = "uncertain"
)

type CandidateJudgment struct {
    Repository string         `json:"repository"`
    Label      RelevanceLabel `json:"label"`
    Reason     string         `json:"reason"`
}
```

Cada julgamento exige motivo curto e deve ser independente do score atual do
Bateia. Repositórios sem entrada permanecem `unjudged`; não viram `irrelevant`.
Duplicatas e julgamentos conflitantes são erros.

## Benchmark mínimo

Executar `literal` e `lexical` sob o mesmo provedor e orçamento para pelo menos
quatro ideias:

1. impacto estrutural em código;
2. análise estática;
3. processamento de eventos;
4. organização de receitas, como controle fora do domínio principal.

Para cada estratégia e ideia:

- mesmo idioma/filtro quando aplicável;
- mesmo `max-per-query`;
- mesma janela de captura sempre que possível;
- pelo menos os dez primeiros candidatos únicos julgados, ou todos quando houver
  menos de dez;
- motivo registrado para cada julgamento.

Se a rede pública impedir a captura completa, o passo não deve fabricar dados.
Registrar a limitação e manter fixtures sintéticas somente para testar o contrato;
o estado permanece `BLOQUEADO_POR_MEDIÇÃO`, não `IMPLEMENTADO`.

## Métricas

Calcular por ideia, estratégia, expansão e agregado:

```text
judged_precision = relevantes / (relevantes + irrelevantes)
uncertain_rate    = incertos / candidatos únicos
unjudged_rate     = não julgados / candidatos únicos
expansion_hit_rate = expansões com ao menos um relevante / expansões executadas
unique_yield      = candidatos únicos / expansões executadas
overlap_rate      = candidatos encontrados por mais de uma expansão / candidatos únicos
failure_rate      = expansões com falha / expansões planejadas
```

Sem julgamentos relevantes ou irrelevantes, `judged_precision` deve ser `null`, não
zero. Não calcular recall: a população total de projetos relevantes é desconhecida.

O relatório também deve mostrar volume absoluto. Uma precisão maior sobre apenas
um candidato não pode parecer automaticamente superior a dez resultados úteis.

## Critério experimental

O lexical é considerado promissor para descoberta se, no benchmark mínimo:

- não reduzir a precisão julgada agregada em mais de 10 pontos percentuais;
- aumentar o número absoluto de candidatos relevantes em pelo menos duas das
  quatro ideias;
- possuir `expansion_hit_rate` maior que zero em pelo menos três ideias;
- não exceder 25% de falhas de expansão;
- manter 100% dos candidatos rastreáveis a uma ou mais expansões;
- nenhum resultado depender de uma query efetiva ausente do snapshot.

Esse critério não declara o lexical definitivo. Ele decide somente se vale levá-lo
à triagem e extração de evidências.

## Prova de mordida

Construir uma fixture em que:

```text
exp-problem   → repo A, repo B
exp-blast     → repo B, repo C
exp-runtime   → falha do provedor
```

O teste deve provar que:

- a saída única é `A, B, C` em ordem determinística;
- B mantém as duas expansões e queries de origem;
- a falha de `exp-runtime` permanece visível;
- A e C não herdam origens que não os encontraram;
- o caso de uso retorna resultado parcial, não erro total;
- se todas as expansões falharem, retorna erro;
- conclusões concorrentes em ordem diferente geram JSON canônico idêntico.

Se a deduplicação apagar a origem de B, o passo não mordeu.

## Testes obrigatórios

### Unidade — L1

1. uma chamada lógica é feita por expansão;
2. concorrência nunca supera o limite;
3. cancelamento interrompe novas chamadas;
4. deduplicação ignora caixa em `owner/repo`;
5. todas as origens são mescladas sem duplicata;
6. ordem não depende da conclusão das goroutines;
7. falha parcial é registrada e preserva resultados;
8. falha total retorna erro acionável;
9. plano de entrada não é modificado;
10. métricas tratam denominadores vazios e precisão desconhecida;
11. julgamento ausente difere de `irrelevant`;
12. julgamento duplicado ou conflitante é rejeitado.

### Integração — L2/L3/L4

1. Sourcegraph registra query completa e fallbacks na ordem real;
2. erro de API e resposta vazia permanecem distinguíveis;
3. token nunca aparece na serialização ou mensagem de erro;
4. `idea discover` grava snapshot válido e versionado;
5. `--output` ausente impede a chamada remota;
6. `idea discovery evaluate` funciona com ambiente vazio;
7. fixtures congeladas não acessam rede;
8. comandos offline anteriores continuam sem construir Sourcegraph;
9. golden literal do Passo 0001 continua intacto;
10. snapshot e julgamentos inválidos produzem erros acionáveis.

## Critérios de verificação

- [x] Contratos de descoberta e julgamento existem em L1 sem I/O.
- [x] Cada expansão gera descoberta independente e rastreável.
- [x] Deduplicação preserva todas as origens do candidato.
- [x] Tentativas e fallbacks reais do provedor ficam observáveis.
- [x] Falhas parciais não apagam resultados válidos.
- [x] Snapshot versionado permite avaliação sem rede.
- [x] Ausência de julgamento não é convertida em irrelevância.
- [x] Métricas evitam recall e falsa precisão.
- [x] Benchmark literal × lexical contém quatro ideias e julgamentos justificados.
- [x] Critério experimental é calculado, não decidido por impressão.
- [x] Prova de mordida cobre deduplicação, concorrência e falha parcial.
- [x] Novos componentes possuem prompts L0 e linhagem válida.
- [x] `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.
- [x] Um laudo registra captura, métricas, falhas e decisão experimental.

## Como este passo pode falhar

- o fallback curto encontra resultados, mas eles são atribuídos à query original;
- mais queries aumentam volume e custo sem aumentar candidatos relevantes;
- deduplicação apaga quais capacidades encontraram um projeto;
- resultados ao vivo tornam os testes instáveis;
- não julgados são tratados como negativos e distorcem precisão;
- julgamentos são guiados pelo score existente em vez da ideia;
- diferenças de orçamento tornam literal e lexical incomparáveis;
- concorrência dispara cota ou torna o snapshot não determinístico.

## Fora de escopo

- baixar README, manifesto, árvore ou código dos candidatos;
- usar o funil completo `SearchFunnel` para cada expansão;
- pontuar capacidades como comprovadas;
- extrair evidência ou verificar licença;
- modificar o score ontológico existente;
- embeddings, LLM ou reranking semântico;
- combinar Sourcegraph e GitHub no mesmo benchmark;
- calcular recall sobre o universo aberto;
- recomendar `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD`;
- automatizar julgamentos humanos.

## Próximo passo condicionado

Se o lexical for promissor, executar triagem e extração somente nos candidatos
relevantes ou incertos, produzindo evidência por capacidade da `IdeaSpec`.

Se aumentar volume sem relevância, revisar ou podar regras usando produtividade
por expansão. Se o problema estiver nos fallbacks do provedor, corrigir a tradução
de queries antes de sofisticar a expansão. Se a medição ficar inconclusiva por
falta de julgamentos, ampliar o benchmark sem alterar o algoritmo.
