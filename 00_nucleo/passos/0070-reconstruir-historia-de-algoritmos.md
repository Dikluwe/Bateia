# Passo 0070 — Reconstruir história de algoritmos

## Intenção congelada

Criar no Bateia um módulo científico, separado do radar 0066–0069, para analisar uma
reconstrução histórica declarada de algoritmos. O módulo recebe alegações tipadas e suas
bases rastreáveis, distingue evidência de inferência, torna controvérsia e incerteza
observáveis e nunca promove ausência de conhecimento a fato.

O passo é deliberadamente offline. Não busca fontes, não acessa rede, não consulta LLM,
não executa browser, não exporta Archify e não altera o radar. A pessoa ou processo que
monta o input continua responsável por localizar, citar e preservar as fontes.

Regime: protocolo completo Tekt, como ensaio procedural sem atestação física de
isolamento. Autor da intenção: usuário. Autor do contrato: executor `contract_0070`, com
leitura da intenção e regras Tekt e escrita somente nos artefatos L0 de 0070; sem acesso
ou autoria de futura implementação, testes, ataques, gates, selo ou veredito.

## Superfície futura

O MVP terá tipos em `01_core/domain/algorithm_history.go` e análise pura em
`01_core/usecase/algorithm_history.go`:

```text
usecase.AnalyzeAlgorithmHistory(domain.AlgorithmHistoryInput) domain.AlgorithmHistoryReport
```

A casca offline exporá:

```text
bateia algorithm history analyze --input FILE [--output FILE]
```

O dispatcher fará somente encaminhamento. Sem `--output`, stdout conterá exclusivamente
um relatório JSON. Com `--output`, a escrita será atômica e stdout ficará vazio.

## Schema público de entrada

`schema_version` é exatamente `algorithm-history-input-v1`. O JSON contém as coleções
`algorithms`, `provenance`, `claims` e `bases`:

```json
{
  "schema_version": "algorithm-history-input-v1",
  "algorithms": [{
    "id": "beta-sort",
    "name": "Beta Sort",
    "chronology": "APPROXIMATE",
    "year_start": 1970,
    "year_end": 1972,
    "chronology_reason": "bounded by the cited archive records",
    "provenance_ids": ["archive-beta"]
  }],
  "provenance": [{
    "id": "archive-beta",
    "kind": "PRIMARY_SOURCE",
    "origin_id": "archive:collection-beta",
    "citation": "Beta laboratory notebook, item 7",
    "locator": "archive://beta/item-7"
  }],
  "claims": [{
    "id": "beta-from-alpha",
    "subject_id": "beta-sort",
    "relation": "DERIVED_FROM",
    "object_id": "alpha-sort",
    "statement": "Beta Sort derived its partition schedule from Alpha Sort"
  }],
  "bases": [{
    "id": "notebook-support",
    "claim_id": "beta-from-alpha",
    "kind": "EVIDENCE",
    "stance": "SUPPORT",
    "provenance_id": "archive-beta",
    "summary": "The notebook explicitly attributes the schedule to Alpha Sort",
    "derived_from_ids": []
  }]
}
```

### Algoritmos e cronologia

Cada algoritmo tem `id`, `name`, `chronology`, `year_start`, `year_end`,
`chronology_reason` e `provenance_ids`.

- `chronology` é `KNOWN`, `APPROXIMATE` ou `UNKNOWN`.
- `KNOWN` e `APPROXIMATE` exigem anos inteiros no intervalo `[-9999,9999]`,
  `year_start <= year_end` e `chronology_reason` não vazio.
- `UNKNOWN` exige ambos os anos ausentes (`null`) e razão pública não vazia.
- `provenance_ids` é não vazio em todos os casos. Cronologia desconhecida continua
  rastreável à fonte que documenta a lacuna; não vira data zero.

### Proveniência

Cada item de proveniência tem `id`, `kind`, `origin_id`, `citation` e `locator` não
vazios. `kind` é `PRIMARY_SOURCE`, `SECONDARY_SOURCE` ou `EMPIRICAL_ARTIFACT`.
`origin_id` é identidade opaca exata da origem física ou editorial: duas citações do
mesmo livro, arquivo ou artefato compartilham a origem e não contam como independentes.
`locator` é apenas dado; o módulo não o abre nem valida remotamente.

### Alegações e relações

Cada alegação tem `id`, `subject_id`, `relation`, `object_id` e `statement`. Sujeito e
objeto referenciam algoritmos distintos. Relações admitidas:

- `DERIVED_FROM`: o sujeito foi derivado do objeto;
- `INFLUENCED_BY`: o objeto influenciou o desenho do sujeito sem derivação integral;
- `IMPROVES`: o sujeito alega melhorar uma propriedade declarada do objeto;
- `GENERALIZES`: o sujeito amplia o domínio ou abstração do objeto;
- `SPECIALIZES`: o sujeito restringe o objeto a um domínio ou hipótese;
- `COMBINES`: o sujeito incorpora material do objeto com ao menos outra contribuição;
- `REDISCOVERS`: o sujeito reencontra independentemente resultado associado ao objeto;
- `RESEMBLES`: existe semelhança observável, sem alegar causalidade histórica.

O analisador não infere arestas transitivas, reciprocidade, cronologia causal, autoria ou
independência. Ciclos são preservados como alegações explícitas, não usados como prova.

### Bases, evidência e inferência

Cada base tem `id`, `claim_id`, `kind`, `stance`, `provenance_id`, `summary` e
`derived_from_ids`. `kind` é `EVIDENCE` ou `INFERENCE`; `stance` é `SUPPORT`, `CONFLICT`
ou `AMBIGUOUS`.

- Uma base `EVIDENCE` representa conteúdo observável na fonte citada e exige
  `derived_from_ids` vazio.
- Uma base `INFERENCE` representa interpretação do autor do dataset. Ela exige um ou
  mais `derived_from_ids`, todos referindo bases canônicas `EVIDENCE` da mesma alegação.
- Toda base exige proveniência e resumo. Toda alegação exige ao menos uma base
  `EVIDENCE` canônica. Inferência nunca incrementa contagem de evidência, origem
  independente ou corroboração.

## Canonicalização e deduplicação

IDs de algoritmo, proveniência, alegação e base são aparados em whitespace ASCII nas
bordas, convertidos para ASCII minúsculo e aceitam somente `[a-z0-9._:/-]`, com 1 a 128
bytes. Referências usam essa forma canônica. `origin_id`, textos e locators são opacos,
válidos UTF-8, não vazios e sem whitespace ASCII nas bordas; não sofrem case folding.

Repetição byte-semântica do mesmo ID é idempotente. O mesmo ID canônico com payload
distinto é `ambiguous_input`. Além disso:

- proveniências com a mesma tupla `(kind, origin_id, citation, locator)` retêm o menor ID;
- alegações com a mesma tupla `(subject_id, relation, object_id, statement)` retêm o
  menor ID;
- bases com a mesma tupla `(claim_id, kind, stance, provenance_id, summary,
  derived_from_ids)` retêm o menor ID.

Referências são remapeadas para os representantes mínimos. `provenance_ids` e
`derived_from_ids` são conjuntos ordenados: nil, vazio e duplicatas produzem a mesma
forma canônica. Após o remapeamento, referências inexistentes, inferência derivada de
outra inferência ou de evidência de outra alegação falham fechadas como `invalid_input`.

## Outcomes públicos por alegação

Somente bases `EVIDENCE` determinam o outcome. Após deduplicação, conte evidências por
stance e origens exatas distintas:

1. `CONTESTED`: há ao menos uma evidência `SUPPORT` e uma `CONFLICT`;
2. `UNKNOWN`: não se aplica 1 e há evidência `AMBIGUOUS`;
3. `CONTRADICTED`: há `CONFLICT` e nenhuma `SUPPORT`;
4. `CORROBORATED`: há `SUPPORT` proveniente de ao menos duas `origin_id` distintas e
   nenhuma `CONFLICT` ou `AMBIGUOUS`;
5. `SUPPORTED`: há `SUPPORT` de exatamente uma origem e nenhuma `CONFLICT` ou
   `AMBIGUOUS`;
6. `UNKNOWN`: qualquer restante, sem promoção implícita a sucesso.

O relatório torna inferências visíveis em listas separadas por stance. Elas podem explicar
uma reconstrução, mas não alteram o outcome. `RESEMBLES` obedece aos mesmos critérios e
nunca é convertido automaticamente em `INFLUENCED_BY` ou `DERIVED_FROM`.

## Schema público de saída

`AlgorithmHistoryReport` tem:

- `schema_version: "algorithm-history-report-v1"`;
- `status`: `COMPLETE` ou `UNKNOWN`;
- `reason`: reason code público;
- `assessments`: uma por alegação canônica, ordenada por `claim_id`;
- `metrics`: contagens e operações.

Cada assessment contém `claim_id`, `subject_id`, `relation`, `object_id`, `statement`,
`outcome`, `supporting_evidence_ids`, `conflicting_evidence_ids`,
`ambiguous_evidence_ids`, `supporting_inference_ids`, `conflicting_inference_ids`,
`ambiguous_inference_ids`, `distinct_supporting_origins` e `reason`. Todas as listas são
ordenadas lexicograficamente e não nulas.

Reasons de outcome são, respectivamente, `evidence_contested`, `evidence_ambiguous`,
`evidence_contradicts_claim`, `independent_origins_corroborate`,
`single_origin_support` e `insufficient_knowledge`. Relatório `COMPLETE` usa reason
`analyzed`; relatório global `UNKNOWN` tem assessments vazios e reason pública de
validação ou budget.

Metrics contém exatamente `input_algorithms`, `input_provenance`, `input_claims`,
`input_bases`, `algorithms`, `provenance`, `claims`, `bases`, `duplicate_algorithms`,
`duplicate_provenance`, `duplicate_claims`, `duplicate_bases`, `evidence`, `inferences`,
os outcomes `corroborated`, `supported`, `contested`, `contradicted`, `unknown`, além de
`operations`, `limit_exceeded` e `verdict_changes`. Este último é sempre zero no passo
aditivo.

## Validação, budgets e imutabilidade

Ordem de falha observável: `invalid_schema`, `budget_exceeded`, `no_algorithms`,
`no_provenance`, `no_claims`, `no_bases`, `invalid_input`, `ambiguous_input`. A primeira
classe aplicável encerra a análise sem assessment parcial. `limit_exceeded` só é true
para `budget_exceeded`.

Limites: 256 algorithms, 1024 provenance, 512 claims, 2048 bases, 64 referências por
lista, 128 bytes por ID/origin, 4096 bytes por texto/locator, 1 MiB agregado de texto e
16384 operações. Verificar limites declarados antes de alocação proporcional. Excesso,
overflow de contagem ou estrutura opaca resulta em relatório `UNKNOWN` sem parcial.

O contador soma uma operação por item de entrada visitado; uma por referência estrangeira
visitada (`provenance_ids`, os dois endpoints de claim, `claim_id` e `provenance_id` de
base, e `derived_from_ids`); uma por base canônica distribuída em stance; e uma por claim
canônica classificada. O exemplar consome 53 operações.

O trabalho é `O(n log n)` pelas ordenações limitadas e `O(n)` fora delas. Não há busca de
caminhos, fecho transitivo, power set ou recursão dependente do grafo. O input e todas as
slices alcançáveis permanecem imutáveis. Permutações equivalentes, duplicatas e nil/vazio
em listas opcionais produzem relatório profundamente idêntico.

## CLI offline e falha fechada

`RunAlgorithmHistoryCommand(args, stdout)` aceita exatamente `history analyze --input
FILE [--output FILE]` (o dispatcher `algorithm` consome somente `algorithm`). Flags
ausentes, repetidas, desconhecidas ou argumentos posicionais extras são erro. O decoder
aceita um único JSON, rejeita campos desconhecidos e trailing JSON.

Status `UNKNOWN` é serializado normalmente como resultado científico quando a entrada foi
decodificada; erros de uso, leitura ou JSON não produzem stdout nem destino. Sem output,
um único JSON indentado em dois espaços e newline final vai para stdout. Com output,
gravação atômica somente após análise completa, stdout vazio, preservação do destino em
erro e rejeição de alias lexical, symlink ou hardlink do input. Nenhum efeito de rede,
relógio, aleatoriedade ou execução externa é permitido.

## Exemplar mensurável

O baseline sintético contém três algoritmos, quatro proveniências de quatro origens, três
alegações e sete bases (cinco evidências, duas inferências):

- `beta-from-alpha` / `DERIVED_FROM`: dois supports de origens distintas e uma inferência
  support -> `CORROBORATED`;
- `gamma-improves-beta` / `IMPROVES`: um support e um conflict de evidência ->
  `CONTESTED`;
- `gamma-resembles-alpha` / `RESEMBLES`: uma evidência ambiguous e uma inferência support
  -> `UNKNOWN`.

Esperado: status `COMPLETE`; 3 assessments; métricas `corroborated=1`, `contested=1`,
`unknown=1`, `supported=0`, `contradicted=0`, `evidence=5`, `inferences=2`. O fixture é
declaradamente sintético e não constitui afirmação sobre algoritmos históricos reais.

## Materialização futura e gates

Arquivos Go futuros: `01_core/domain/algorithm_history.go`,
`01_core/usecase/algorithm_history.go`, seu teste independente, e
`02_shell/cli/algorithm_history.go`, seu teste independente, além do encaminhamento
aditivo mínimo no dispatcher. Não editar 0066–0069.

Gates futuros: oráculos positivos/negativos/opacos, mutation score 1.0, testes focais,
20 permutações, imutabilidade profunda, race, `go test ./...`, `go vet ./...`, build,
`crystalline-lint .` e `git diff --check`. Nenhum candidato ou gate é executado nesta
autoria contratual.
