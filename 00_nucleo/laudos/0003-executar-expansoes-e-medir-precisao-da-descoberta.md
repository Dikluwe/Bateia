# Laudo 0003 — Precisão real da descoberta por expansões

**Data**: 2026-08-24  
**Passo verificado**: [0003 — Executar expansões e medir precisão da descoberta](../passos/0003-executar-expansoes-e-medir-precisao-da-descoberta.md)  
**Resultado**: `APROVADO_COM_RUÍDO`

## O que foi implementado

- contratos L1 para requisição, tentativas do provedor, snapshot, candidato
  deduplicado, julgamentos e métricas;
- porta `ObservableDiscoveryProvider` paralela ao `SearchProvider` existente;
- executor puro com concorrência limitada, cancelamento, falhas parciais e ordem
  canônica pelo plano;
- deduplicação sem diferença de caixa em `owner/repo`, preservando todas as queries
  e expansões de origem;
- avaliador offline com precisão anulável, incerteza, não julgados, hit rate,
  produtividade, overlap e falhas;
- adaptador Sourcegraph que registra query completa, fallbacks, contagem e erros;
- comando remoto `idea discover`, com snapshot obrigatório;
- comando offline `idea discovery evaluate`;
- oito snapshots reais: quatro ideias por duas estratégias;
- quatro conjuntos de julgamentos justificados e oito avaliações congeladas.

## Captura

As capturas consultaram o Sourcegraph público com:

```text
max-per-query = 10
concurrency = 2
language = go
provider = sourcegraph
```

Ideias avaliadas:

1. impacto estrutural;
2. análise estática;
3. processamento de eventos;
4. organização de receitas.

As capturas e julgamentos estão versionados em
`02_shell/cli/testdata/discovery_benchmark/`. Testes e replay não repetem chamadas
remotas.

## Resultado agregado

| Métrica | Literal | Lexical |
| :------ | ------: | ------: |
| Candidatos únicos | 22 | 62 |
| Relevantes | 0 | 8 |
| Irrelevantes | 19 | 49 |
| Incertos | 3 | 5 |
| Não julgados | 0 | 0 |
| Precisão julgada | 0% | 14,04% |
| Expansões com relevante | 0/4 | 8/12 |
| Expansion hit rate | 0% | 66,67% |
| Unique yield | 5,50 | 5,17 |
| Overlap observado | 0% | 0% |
| Falhas de expansão | 0% | 0% |

O lexical encontrou candidatos relevantes nas quatro ideias. O literal não
encontrou nenhum candidato julgado relevante.

## Resultado por ideia

| Ideia | Estratégia | Precisão julgada | Relevantes | Hit rate |
| :---- | :--------- | ----------------: | ---------: | -------: |
| Impacto estrutural | literal | 0% | 0 | 0% |
| Impacto estrutural | lexical | 9,52% | 2 | 66,67% |
| Análise estática | literal | 0% | 0 | 0% |
| Análise estática | lexical | 14,29% | 3 | 66,67% |
| Eventos | literal | `null` | 0 | 0% |
| Eventos | lexical | 40% | 2 | 66,67% |
| Receitas | literal | 0% | 0 | 0% |
| Receitas | lexical | 10% | 1 | 33,33% |

A precisão literal de eventos é `null`: seu único candidato foi julgado incerto.
O avaliador preservou corretamente a diferença entre “sem evidência suficiente” e
“irrelevante”.

## Achados concretos

Para impacto estrutural, o lexical encontrou:

- `supermodeltools/cli`, cuja descrição declara grafos de chamadas, dependências,
  detecção de código morto e blast radius;
- `dejo1307/enola`, que declara snapshot arquitetural e grafo de conhecimento.

Para análise estática, apareceram `bvobart/mllint`, `Halleck45/ast-metrics` e
`reviewdog/errorformat`. Para eventos, `stackus/edat` e `etkecc/postmoogle`. No
controle externo, `reaper47/recipya` corresponde diretamente a gerenciamento de
receitas para famílias.

Esses julgamentos usam metadados de descoberta. Eles demonstram sinal suficiente
para triagem, não capacidades comprovadas por código.

## Fallbacks observados

A consulta literal em português para impacto estrutural produziu zero resultados
com query completa e com idioma. Somente o fallback genérico `mostrar impacto`,
sem idioma, retornou nove candidatos — todos julgados irrelevantes.

Isso confirma que atribuir apenas a query planejada aos resultados seria enganoso.
O snapshot agora registra a cadeia efetiva de tentativas.

## Decisão experimental

O baseline lexical foi considerado **promissor** porque:

- aumentou, em vez de reduzir, a precisão julgada agregada;
- aumentou relevantes em 4/4 ideias, acima do mínimo de 2/4;
- teve hit rate positivo em 4/4 ideias;
- teve zero falhas de expansão;
- todos os candidatos mantiveram uma ou mais expansões de origem;
- todas as tentativas efetivas ficaram presentes no snapshot.

O resultado não autoriza executar extração em todos os 62 candidatos. O ruído de
49 irrelevantes mostra que a próxima fase deve limitar triagem aos relevantes e
incertos e medir evidência por capacidade.

## Verificações

```text
go test -race ./...    PASS
go vet ./...           PASS
go build ...           PASS
crystalline-lint .     PASS — No violations found
git diff --check       PASS
```

A prova sintética também verificou A/B/C, múltiplas origens em B, falha parcial,
falha total, limite de concorrência e ordem independente das goroutines.

## Limites

- julgamentos foram feitos pelas descrições retornadas, sem leitura de README ou
  código;
- resultados do Sourcegraph são uma fotografia de 2026-08-24;
- não houve sobreposição entre expansões nesta captura, portanto a mesclagem foi
  provada por teste sintético, não pelo benchmark;
- o filtro `language=go` pode excluir soluções relevantes em outras stacks;
- precisão permanece baixa em termos absolutos, apesar do ganho contra o literal;
- não há recall mensurável num universo aberto.

## Próximo passo

Triar os 13 candidatos relevantes ou incertos, buscar README, manifesto, árvore e
referência imutável e produzir evidência vinculada às capacidades da `IdeaSpec`.
Somente depois dessa prova será possível distinguir projeto reutilizável de mera
descrição coincidente.
