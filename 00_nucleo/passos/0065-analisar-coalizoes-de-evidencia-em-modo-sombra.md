# Passo 0065 — Analisar coalizões de evidência em modo sombra

## Intenção congelada

Materializar em L1 um experimento puro, conservador e mensurável que descreve quando
fontes de evidência formam uma coalizão de suporte, quando há contraevidência, quais
papéis faltam e se a proveniência permite afirmar origens distintas. O resultado é um
artefato derivado em modo sombra: ele não substitui, promove, rebaixa nem recalcula
`EvidenceStatus`, scores, confiança, `CandidateAssessment` ou qualquer veredito atual.

O nome normativo é **analisador de coalizões de evidência**. A implementação interna não
precisa ser, nem alegar ser, um complexo simplicial. O modelo inicial é uma relação
limitada entre alegação, regra e testemunhas; nenhuma biblioteca de grafos, hipergrafos,
probabilidade ou topologia é necessária.

## Regime e autoridades

Aplicar o protocolo completo da skill `tekt-materializacao-segregada` como ensaio sem
atestação física de isolamento. O filesystem é compartilhado.

- autor da intenção: usuário; congela a obrigação acima;
- autor do contrato: escreve este passo, prompts e `benchmark/0065`; não lê uma
  implementação candidata futura e não escreve Go;
- autor dos oráculos: recebe somente intenção, contrato e baseline selados; escreve
  casos positivos, negativos e opacos sem ler o candidato;
- adversário: recebe intenção, contrato e baseline; materializa as mutações exigidas,
  sem corrigir a solução;
- implementador: recebe este passo, contrato e prompts selados; não edita contrato,
  baseline, oráculos ou mutantes;
- verificador: executa gates sobre artefatos congelados e não corrige o que julga.

Qualquer mudança na intenção, neste passo, em `benchmark/0065/contract.md`, no baseline
ou nos oráculos após o selo invalida a cadeia a partir da primeira entrada alterada.

## API pública L1

Criar tipos em `01_core/domain/evidence_coalition.go` e a função pura em
`01_core/usecase/evidence_coalition.go`:

```text
AnalyzeEvidenceCoalitions(domain.EvidenceCoalitionInput) domain.EvidenceCoalitionReport
```

Tipos públicos, com JSON `snake_case` e sem `omitempty`:

```text
EvidenceCoalitionInput
  SchemaVersion string
  Rules         []EvidenceCoalitionRule
  Claims        []EvidenceCoalitionClaim

EvidenceCoalitionRule
  ID                     string
  RequiredRoles          []EvidenceCoalitionRole
  MinimumDistinctOrigins int

EvidenceCoalitionClaim
  ID             string
  RuleID         string
  BaselineStatus EvidenceStatus
  Witnesses      []EvidenceCoalitionWitness

EvidenceCoalitionWitness
  ID           string
  Role         EvidenceCoalitionRole
  Polarity     EvidenceCoalitionPolarity
  OriginStatus EvidenceOriginStatus
  OriginID     string
  Source       EvidenceSource

EvidenceCoalitionAssessment
  ClaimID                string
  RuleID                 string
  BaselineStatus         EvidenceStatus
  DecisionImpact         string
  Outcome                EvidenceCoalitionOutcome
  Independence           EvidenceOriginStatus
  SupportingWitnessIDs   []string
  ConflictingWitnessIDs  []string
  MissingRoles           []EvidenceCoalitionRole
  DistinctOrigins        []string
  Reason                 string

EvidenceCoalitionMetrics
  Claims                   int
  Witnesses                int
  Satisfied                int
  Conflicts                int
  MissingRoles             int
  UnknownIndependence      int
  InsufficientIndependence int
  Invalid                  int
  Operations               int
  LimitExceeded            bool
  VerdictChanges           int

EvidenceCoalitionReport
  SchemaVersion string
  Status        string
  Reason        string
  Assessments   []EvidenceCoalitionAssessment
  Metrics       EvidenceCoalitionMetrics
```

Constantes públicas exatas:

- input: `evidence-coalition-input-v1`; output: `evidence-coalition-report-v1`;
- papéis: `IMPLEMENTATION`, `INDEPENDENT_TEST`, `EXECUTION`, `DOCUMENTATION`;
- polaridades: `SUPPORT`, `CONFLICT`;
- proveniência: `RECORDED`, `SAME_ORIGIN`, `DISTINCT_ORIGIN`, `UNKNOWN_ORIGIN`;
- outcomes: `SATISFIED`, `CONFLICT`, `MISSING_ROLES`,
  `UNKNOWN_INDEPENDENCE`, `INSUFFICIENT_INDEPENDENCE`, `INVALID_INPUT`;
- impacto: sempre `NONE` nesta versão;
- status de relatório: `COMPLETE` ou `UNKNOWN`.

Nomes podem receber aliases Go idiomáticos, mas valores serializados, campos e função
acima são contrato. Slices públicas no relatório são sempre não nulas; ausência é `[]`.

## Semântica observável

Cada claim referencia exatamente uma regra existente. IDs de claim, regra e testemunha
são strings ASCII não vazias, sem espaços nas pontas, únicas no respectivo escopo.
`BaselineStatus` aceita somente os cinco valores atuais de `EvidenceStatus`. Uma regra
contém de um a quatro papéis, sem repetição, e
`1 <= MinimumDistinctOrigins <= len(RequiredRoles)`.

Uma testemunha `RECORDED` exige `OriginID` não vazio; uma testemunha de origem não
registrada exige `OriginStatus=UNKNOWN_ORIGIN` e `OriginID` vazio. Estados derivados
`SAME_ORIGIN` e `DISTINCT_ORIGIN` não são entradas válidas. O analisador não valida a
verdade externa de `OriginID`: ele apenas conserva a proveniência declarada. Hash,
caminho, símbolo, commit ou IDs diferentes nunca provam por si sós origem distinta.

Testemunhas são deduplicadas pela identidade física canônica de `EvidenceSource`
(`kind`, owner, repo, commit, path, symbol, line, rule_id e content_hash). Repetir a
mesma fonte com outro ID não cria papel ou origem adicional. Se as duplicatas divergem
em papel, polaridade ou origem, o claim é `INVALID_INPUT`; duplicatas semanticamente
iguais conservam somente o menor ID lexicográfico.

Somente `SUPPORT` cobre papéis exigidos. `DOCUMENTATION` não cobre `IMPLEMENTATION`, e
nenhum papel cobre outro. Toda testemunha `CONFLICT` válida é preservada como
contraevidência, mesmo se a coalizão de suporte estiver completa.

Para cada claim válido, aplicar esta precedência:

1. qualquer contraevidência: `CONFLICT`;
2. papel exigido sem testemunha de suporte: `MISSING_ROLES`;
3. todos os papéis existem, mas ao menos um papel não possui qualquer suporte com origem
   `RECORDED`: `UNKNOWN_INDEPENDENCE`, ainda que outros papéis já exibam o número do
   limiar; toda função exigida precisa participar com proveniência registrada;
4. todos os papéis possuem suporte com origem registrada, mas a união de `OriginID`
   entre suportes dos papéis exigidos é menor que `MinimumDistinctOrigins`:
   `INSUFFICIENT_INDEPENDENCE` e independência `SAME_ORIGIN`;
5. papéis completos e limiar demonstrado: `SATISFIED` e independência
   `DISTINCT_ORIGIN`, exceto limiar 1, cujo valor derivado é `RECORDED`.

Em `MISSING_ROLES`, independência é `UNKNOWN_ORIGIN`. Em `CONFLICT`, a independência
continua sendo calculada sobre os suportes e será `DISTINCT_ORIGIN`, `RECORDED`,
`SAME_ORIGIN` ou `UNKNOWN_ORIGIN`; conflito jamais desaparece por maioria. `MissingRoles`
segue a ordem canônica dos valores de papel, não a ordem de entrada.

`SupportingWitnessIDs`, `ConflictingWitnessIDs`, `DistinctOrigins`, regras, claims e
assessments são ordenados lexicograficamente pelos seus IDs/valores canônicos. A função
não muta input nem slices/objetos alcançáveis. Duas permutações semanticamente iguais,
incluindo `nil` versus slice vazia, produzem relatórios JSON byte a byte iguais após
`encoding/json.Marshal`.

Entrada estruturalmente inválida gera assessment `INVALID_INPUT` quando o claim pode
ser identificado. Schema desconhecido, lista vazia de claims, duplicação que torna a
associação ambígua ou limite global excedido gera relatório `UNKNOWN`, razão pública
estável e nenhum sucesso parcial. `Unknown` nunca é convertido em `SATISFIED`.

Razões públicas exatas: assessments usam `coalition_satisfied`,
`conflicting_evidence`, `missing_required_roles`, `origin_unrecorded`,
`distinct_origins_below_minimum` ou `invalid_claim`. Relatório `COMPLETE` usa
`analyzed` ou `analyzed_with_invalid_claims`; relatório `UNKNOWN` usa
`invalid_schema`, `no_claims`, `ambiguous_input` ou `budget_exceeded`. Não incluir
conteúdo de mapa, ponteiro ou ordem incidental nas razões.

## Limites e custo

Limites fixos da versão 1:

- até 128 regras e 128 claims;
- até 64 testemunhas por claim e 4096 no total;
- até quatro papéis por regra;
- até 8192 operações lógicas observadas.

`Metrics.Operations` conta uma unidade por testemunha normalizada e uma por par
claim/papel exigido inspecionado. Ordenação e validações escalares não adicionam
unidades. A implementação faz uma passagem indexada por testemunhas e papéis; é
proibido enumerar subconjuntos ou o power set. Se qualquer limite for excedido, o
relatório é `UNKNOWN`, `reason=budget_exceeded`, `limit_exceeded=true`, sem assessments
parciais. Esse estado é resultado válido, não erro nem sucesso.

## Invariantes de não regressão

- `BaselineStatus` é ecoado sem mudança em cada assessment;
- `DecisionImpact` é sempre `NONE` e `Metrics.VerdictChanges` é sempre zero;
- nenhuma função existente chama o analisador nesta primeira materialização;
- nenhum schema legado recebe campo novo;
- nenhuma fórmula de confiança, agregação, ordenação de candidatos ou decisão muda;
- L1 permanece sem I/O, relógio, rede, filesystem, aleatoriedade ou estado global;
- nenhuma alegação de independência causal nasce de hashes ou caminhos distintos;
- conflito e origem desconhecida permanecem explícitos;
- output não vazio não equivale a veredito de produto.

## Não objetivos

Não integrar o relatório em CLI, persistência, busca, recomendação ou relatório final;
não inferir papéis a partir de texto; não descobrir proveniência; não executar código;
não calcular probabilidades ou novo score; não promover `PROVEN`; não construir TDA,
complexo simplicial ou hipergrafo geral; não provar independência de autores/processos;
não enumerar todas as coalizões mínimas; não alterar o extrator de mecanismos existente.

## Medição, oráculos e mutações obrigatórias

O corpus contratual deve conter, no mínimo: suporte completo com duas origens; limiar 1;
mesma origem insuficiente; origem desconhecida; papel ausente; documentação incapaz de
substituir implementação; conflito junto de suporte completo; duplicata idêntica;
duplicata contraditória; input vazio; budget excedido; e permutações/repetições.

O gate discriminatório exige positivos válidos classificados como esperado, opacos como
`UNKNOWN_INDEPENDENCE` ou relatório `UNKNOWN`, e rejeição de todas as mutações:

- M01: coalizão satisfeita promove/rebaixa `BaselineStatus` ou impacto deixa de `NONE`;
- M02: `UNKNOWN_ORIGIN` conta como origem distinta;
- M03: hash, caminho, commit ou ID distinto prova independência;
- M04: conflito é ignorado quando há suportes suficientes;
- M05: documentação preenche papel de implementação/teste/execução;
- M06: ordem de entrada altera relatório serializado;
- M07: duplicata física cria testemunha, papel ou origem adicional;
- M08: duplicata física contraditória é aceita silenciosamente;
- M09: input vazio ou schema desconhecido vira `COMPLETE`/sucesso;
- M10: limite excedido produz resultado parcial ou `SATISFIED`;
- M11: `nil` e vazio produzem JSON semanticamente diferente;
- M12: implementação enumera subconjuntos ou excede a contagem contratual de operações.

`mutation_score = mutações válidas rejeitadas / mutações válidas` deve ser 1.0 antes do
selo. Medir ainda: acurácia 100% no corpus congelado, `VerdictChanges=0`, igualdade JSON
em pelo menos 20 permutações determinísticas, operações dentro da fórmula e ausência de
delta nos snapshots/vereditos legados selecionados. Tempos, se coletados pelo runner de
teste, são telemetria e não gate de correção.

## Escopo futuro e validação

Implementador: somente os dois novos `.go` espelhados pelos prompts deste passo.
Autor de testes: somente os dois novos `_test.go` espelhados, externos aos packages
quando possível. Nenhum arquivo Go existente deve ser editado.

Após congelar contrato, baseline, oráculos e mutantes: executar testes focais, mutantes,
20 permutações, `go test ./...`, `go build ./04_wiring`, `crystalline-lint .` e
`git diff --check`. Budget de calibração: duas revisões focais por mesma causa e no
máximo três execuções do corpus completo. Duas revisões consecutivas sem ganho no vetor
de classificação suspendem o gate e reabrem o desenho; nunca afrouxar `Unknown`.
