# Passo 0068 — Compilar uma capacidade em consultas técnicas discriminatórias

## Intenção e causa raiz congeladas

Corrigir incrementalmente a perda observada no dossiê `coin-tessellation`: uma capacidade
natural longa em português foi enviada como consulta lexical e a redução posterior
`predicados geometricos` admitiu `Dante-Leoncini/Whisk3D-Editor`, cujo trecho apenas usa a
expressão genérica em comentário. O dossiê v2, escrito manualmente com vocabulário técnico
em inglês, encontrou candidatos mais plausíveis para outras capacidades, mas a consulta
longa de predicados robustos também terminou vazia. Já o controle dirigido `Shewchuk`
recuperou `VTREEM/Carve` e observou a opção `CARVE_USE_EXACT_PREDICATES`.

O recorte 0068 materializa um compilador **finito**, local e auditável somente para a
capacidade:

```text
predicados geometricos de orientacao e in-circle robustos a casos quase degenerados
```

A tradução inglesa canônica é exatamente:

```text
robust geometric orientation and in-circle predicates for nearly degenerate cases
```

Isso não é tradução geral. Entradas fora do catálogo continuam literais e desconhecidas.

## Regime, autoridades e isolamento

Aplicar o protocolo completo da skill `tekt-materializacao-segregada` como ensaio com
segregação procedural, filesystem compartilhado e sem atestação física de isolamento.

- autor da intenção: usuário;
- autor do contrato: lê diagnóstico, dossiês v1/v2, controle isolado e produção anterior;
  escreve somente este passo, prompts e `benchmark/0068`, sem ler candidato 0068 futuro;
- autor dos oráculos: recebe intenção, contrato, baseline e budget congelados; escreve
  testes sem ler o candidato;
- adversário: constrói mutações e controles opacos, sem corrigir a solução;
- implementador: recebe somente passo, contrato, prompts produtivos e selo; não edita
  baseline, contrato, oráculos ou dossiês;
- verificador: executa gates e produz recibos sem editar material julgado.

Mudança posterior em passo, contrato, baseline, budget, prompts ou oráculos invalida a
cadeia a partir da primeira entrada alterada. Hash registra conteúdo e ordem, não prova
isolamento.

## Compilação finita e preservação da intenção

Estender o vocabulário embutido do expansor lexical, preservando compatibilidade das
regras existentes. Uma regra pode declarar opcionalmente `strategy` e
`discriminators`. Ausência de `strategy` mantém `lexical`; discriminadores ausentes
mantêm o comportamento legado. Para regras técnicas, `strategy` deve ser
`technical-lexicon-v1` e exatamente dois discriminadores canônicos, não vazios e
distintos, devem acompanhar cada proposta.

O catálogo contém três regras, nesta ordem, todas acionadas somente quando a frase
portuguesa integral normalizada está contida em uma expansão de origem:

| RuleID | Query produzida | Discriminadores obrigatórios |
|---|---|---|
| `robust-predicates.translation.en` | `robust geometric orientation and in-circle predicates for nearly degenerate cases` | `orientation`, `in-circle` |
| `robust-predicates.symbol-pair` | `orient2d incircle` | `orient2d`, `incircle` |
| `robust-predicates.reference-pair` | `Shewchuk exact predicates` | `Shewchuk`, `predicates` |

Cada proposta conserva `Facet`, `Covers` e o ID da expansão literal em
`Derivation.SourceIDs`; usa seu RuleID em `Derivation.RuleIDs`; recebe IDs estáveis pelo
contrato atual. A expansão literal e `ExpansionPlan.Idea` permanecem intactas. O catálogo
nunca produz `predicados geometricos`, `geometric predicates` isolado, nomes de
repositórios, afirmações de capacidade ou consultas derivadas dos resultados.

O limite injetado do expansor passa de 2 para 3 nos três caminhos produtivos de
`04_wiring/main.go`, permitindo as três propostas deste catálogo. As regras legadas têm
no máximo duas saídas e, portanto, preservam seus resultados. Instanciações explícitas
com outro limite continuam respeitadas.

## Transporte de proveniência e admissão técnica

Adicionar a `DiscoveryRequest` o campo opcional
`Derivation *ExpansionDerivation json:"derivation,omitempty"`. O orquestrador L1 copia
profundamente a derivação da expansão para a requisição, sem mutar plano ou slices.

O adaptador Sourcegraph reconhece o caminho técnico somente quando:

1. `Derivation.Strategy == "technical-lexicon-v1"`;
2. existe exatamente um RuleID e um SourceID não vazios;
3. há exatamente dois discriminadores não vazios, distintos após normalização;
4. a consulta é não vazia.

Entrada técnica inválida falha fechada antes de I/O. Para entrada válida, executar uma
única consulta técnica, preservando `repo:`/`lang:` e acrescentando o mesmo `count:N`
limitado do 0056. Não chamar `PlanDiscoveryQuery`, não reduzir palavras, não remover
linguagem e não fazer fallback genérico.

Reutilizar a coleta limitada e a admissão por arquivo do 0056: o primeiro discriminador
forma o grupo A e o segundo o grupo B de um `DiscoveryQueryPlan` de compatibilidade,
versão `0068-technical-v1`, razão `technical_discriminators`, com a consulta integral em
`Original` e `Unverified`. Um candidato técnico é selecionado somente quando snippets
retidos do **mesmo arquivo** contêm ambos os discriminadores, case-insensitive. Nome,
descrição, path sem snippet, resultado `Repository`, dois arquivos diferentes e trecho
omitido não contam. Os campos legados `ActionFound` e `ObjectFound` representam, neste
caminho compatível, grupo A e grupo B; não alegam ação ou objeto semântico.

`QueryPlan`, `Coverage`, `CandidateDiagnostics`, `RawCandidates`,
`RejectedCandidates`, `Matches` e `AdmissionEvidence` permanecem observáveis como no
0056. `Completeness` continua `unknown`. Erro, timeout ou cancelamento não vira vazio nem
aciona outra tentativa. Uma resposta válida vazia é `EMPTY`. Resultado bruto não vazio,
mas sem cobertura dos dois discriminadores no mesmo arquivo, também é `EMPTY` com
diagnósticos — nunca sucesso útil.

Consultas literais e hipóteses 0055 não técnicas mantêm o comportamento atual. Cada
proposta técnica é uma expansão independente; `DiscoverExpansionPlanWithPolicy` executa
todas dentro dos budgets e une candidatos por `owner/repo`, preservando todas as origens,
queries e expansion IDs. Assim, um resultado em uma query não cancela as demais queries
técnicas. Não há coalescência entre requests nem paralelismo novo.

## Métricas e replay sem rede

Nenhum ground truth entra em código produtivo. O harness externo usa:

1. replay real congelado dos dossiês v1/v2 e do controle `Shewchuk`, somente como
   diagnóstico histórico;
2. corpus contrafactual congelado, com respostas de provedor determinísticas, para testar
   compilação, admissão e união sem rede.

Métricas por capacidade:

- `recall_at_5`: alvos relevantes recuperados nos cinco primeiros candidatos únicos /
  alvos relevantes congelados;
- `contamination_at_5`: candidatos julgados irrelevantes nos cinco primeiros / candidatos
  julgados nos cinco primeiros; sem candidato julgado produz `Unknown`, não zero;
- `evidence_addressable_capabilities`: capacidades com ao menos um candidato admitido por
  ambos discriminadores no mesmo arquivo. Isso significa somente que há evidência lexical
  endereçável para coleta; não significa capacidade comprovada;
- `provider_attempts`, `received_items`, `examined_repositories` e candidatos únicos como
  custo observável. Tempo histórico não é comparável ao replay e não entra no gate.

No contrafactual principal, o baseline literal recupera apenas `noise/Whisk` pelo termo
genérico; o catálogo técnico recupera `geometry/RobustKernel` em duas queries e deduplica
suas origens. O esperado é recall@5 `0 -> 1`, contaminação@5 `1 -> 0.5`, capacidades com
evidência endereçável `0 -> 1` e tentativas `1 -> 4` (literal mais três propostas).
Há controles fora de geometria, frase parcial, caixa/espaçamento, split entre arquivos,
metadata-only, um discriminador só, duplicação e ordem.

Veredito externo `OBSERVED` exige simultaneamente: recall@5 maior; contaminação@5 não
aumentada; ao menos uma capacidade adicional com evidência endereçável; resultado alvo
com duas origens técnicas distintas; nenhuma consulta genérica proibida; custo máximo de
quatro tentativas por origem; nenhum `Unknown`, erro, budget ou regressão. Métrica
indefinida, replay ausente, menos de duas famílias, opacidade ou budget produz
`INCONCLUSIVE`. Demais execuções válidas produzem `NOT_OBSERVED`.

O replay real não pode produzir `OBSERVED` sem nova execução pública: ele congela que v1
e v2 tiveram recall zero para esta capacidade e que o controle dirigido localizou Carve.
Resultados novos no Sourcegraph são temporalmente instáveis e ficam fora deste passo.

## Limites, Unknown e não regressão

- no máximo 3 propostas por expansão de origem e 64 expansões totais no plano enriquecido;
- no máximo 2 discriminadores técnicos, 64 bytes UTF-8 por discriminador e 256 bytes por
  query técnica;
- no máximo 4 tentativas derivadas da origem literal no experimento;
- budgets 0056 permanecem: `count` 20..200, até 32 snippets por repositório e 240 runes por
  snippet; `MaxResults` e `MaxAttempts` não são ampliados dentro de um request;
- limite excedido, derivação ausente/opaca ou métrica obrigatória indefinida é `Unknown`
  no harness e nunca conta como melhoria;
- entradas sem correspondência exata no catálogo, traduções gerais, sinônimos inventados
  e interpretação semântica permanecem desconhecidos.

Não alterar dossiês, schemas de saída existentes além do campo opcional de request,
ranking, screening, avaliação de capacidade, thresholds, radar 0066/0067, evidência,
veredito ou cache. Não adicionar LLM, embeddings, rede, filesystem, relógio ou I/O a L1.
Não afirmar que presença de `orient2d`, `incircle` ou `Shewchuk` comprova robustez,
correção numérica, desempenho, licença ou adequação ao Coin.

## Oráculos, mutações e gates

Oráculos independentes cobrem a tradução exata, as duas consultas curtas, preservação da
frase e origem, estratégia/discriminadores, limites, não mutação, caixa/espaçamento e
compatibilidade das regras antigas. Integração controlada cobre admissão same-file,
split-file, metadata-only, resultado vazio/erro, união/deduplicação com duas origens,
execução das três propostas apesar de resultado anterior e custo exato.

Eliminar ao menos estas mutações: M01 traduzir com texto diferente; M02 remover literal;
M03 emitir fallback genérico; M04 omitir uma consulta curta; M05 usar estratégia lexical
legada; M06 perder provenance; M07 aceitar um discriminador; M08 combinar arquivos; M09
usar nome/descrição/path; M10 parar após primeira proposta não vazia; M11 perder origem na
deduplicação; M12 exceder três propostas; M13 ativar em frase parcial; M14 alterar regra
legada; M15 converter erro/Unknown em vazio; M16 alegar capacidade comprovada.

Selo exige mutation score 1.0, corpus contrafactual 100%, métricas esperadas, 20/20
permutações, nenhum arquivo 0065–0067 alterado e hashes congelados intactos. Depois do
selo: testes focais, `-race`, suíte completa, `go vet` dos pacotes tocados, build de
`./04_wiring`, `crystalline-lint .` e `git diff --check`.

Budget de calibração: duas revisões focais por causa e três execuções completas antes do
selo. Duas revisões consecutivas sem ganho discriminatório reabrem o desenho.
