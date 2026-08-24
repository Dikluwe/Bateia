# Passo 0004 — Construir peneira probabilística explicável

**Estado**: `IMPLEMENTADO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0003 — Executar expansões e medir precisão da descoberta`](0003-executar-expansoes-e-medir-precisao-da-descoberta.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: peneira explicável, calibração empírica e auditoria de descartes

## Intenção

Explicar por que cada projeto apareceu na descoberta e estimar, com incerteza
declarada, se vale submetê-lo à triagem profunda.

O passo deve responder duas perguntas diferentes:

1. **Por que este candidato entrou?**
2. **Com a evidência disponível, qual a chance de ele ser relevante para a ideia?**

Não queremos substituir os julgamentos do Passo 0003 por uma fórmula opaca. Queremos
usar esses julgamentos como dados de calibração e produzir uma peneira cuja decisão
possa ser reconstruída sinal por sinal.

## Correção conceitual

Os 49 projetos classificados como irrelevantes no benchmark anterior não haviam
passado pela triagem do `SearchFunnel`. Eles passaram apenas pela correspondência
textual do Sourcegraph.

A trajetória observada foi:

```text
IdeaSpec
    ↓
expansão literal ou lexical
    ↓
query completa
    ↓ quando vazia
fallback reduzido com linguagem
    ↓ quando vazio
fallback reduzido sem linguagem
    ↓
match textual
    ↓
candidato ainda não triado
```

Exemplos do benchmark:

| Query efetiva | Relevantes julgados |
| :------------ | -------------------: |
| `mostrar impacto` | 0/10 |
| `code impact` | 1/10 |
| `change blast radius` | 1/1 |
| `analisar código` | 0/10 |
| `static code analysis` | 2/10 |
| `source code analyzer` | 1/3 |
| `organizar receitas` | 0/2 |
| `recipe manager` | 1/6 |
| `meal planning` | 0/3 |

O observável sugere que a perda de especificidade nos fallbacks é uma causa central
do ruído. O passo deve testar essa hipótese, não assumi-la como verdade definitiva.

## Incertezas reduzidas

- qual expansão e qual tentativa efetiva trouxeram cada candidato;
- quanto cada fallback se afastou da intenção original;
- quais sinais preliminares distinguem relevante, irrelevante e incerto;
- quanto uma estimativa numérica está calibrada ou apenas ordena candidatos;
- quantos descartes automáticos seriam falsos negativos;
- qual limiar reduz custo sem eliminar projetos úteis.

## Princípio estatístico

Um score entre 0 e 1 não é automaticamente uma probabilidade.

O Bateia só pode chamar uma saída de `probability` quando houver calibração
empírica verificável. Antes disso, deve usar `ranking_score` e declarar
`calibration_status: "uncalibrated"`.

Quando calibrado, o contrato deve preservar:

```text
P(relevante | sinais observados, população e protocolo de julgamento)
```

A interpretação vale somente para a população representada pelos dados. Os 84
resultados observados no Passo 0003 são um início pequeno, não autorização para
generalização universal.

## Entrega

### 1. Traço causal completo em L1

Estender o snapshot para ligar cada candidato à tentativa exata que o encontrou.
Uma forma de referência:

```go
type MatchEvidence struct {
    MatchType string `json:"match_type"`          // repository, file, symbol
    Path      string `json:"path,omitempty"`
    Snippet   string `json:"snippet,omitempty"`
}

type CandidateOrigin struct {
    ExpansionID  string          `json:"expansion_id"`
    PlannedQuery string          `json:"planned_query"`
    AttemptIndex int             `json:"attempt_index"`
    EffectiveQuery string        `json:"effective_query"`
    Evidence     []MatchEvidence `json:"evidence,omitempty"`
}
```

Cada `DiscoveredCandidate` passa a possuir `origins`. O traço deve permitir
responder sem inferência posterior:

- qual campo da `IdeaSpec` originou a expansão;
- qual regra lexical foi aplicada;
- qual query foi planejada;
- qual fallback encontrou o projeto;
- onde ocorreu a correspondência;
- quais outras expansões convergiram no mesmo candidato.

Snapshots antigos de schema `1` continuam legíveis. O novo formato usa
`schema_version: "2"` e uma migração pura e testada. Dados inexistentes no schema
antigo permanecem desconhecidos, não são inventados.

### 2. Evidência de match no Sourcegraph em L3

O adaptador atualmente reduz `Repository` e `FileMatch` ao mesmo
`RepositoryCandidate`. Preservar metadados suficientes para a peneira:

- tipo do resultado;
- caminho para `FileMatch`;
- trecho curto quando a API o fornecer;
- query efetiva responsável pelo conjunto;
- índice da tentativa;
- ausência explícita quando o Sourcegraph não fornecer o dado.

Restrições:

- no máximo 240 caracteres de snippet normalizado;
- nunca serializar token, cabeçalho ou resposta GraphQL integral;
- não contar o mesmo arquivo repetidamente como evidência independente;
- não buscar README ou código adicional neste estágio;
- manter o `SearchProvider` existente compatível.

### 3. Vetor de sinais explicáveis em L1

Extrair sinais puros, separados e serializáveis:

```go
type ScreeningSignals struct {
    ExpansionPrior       float64 `json:"expansion_prior"`
    QueryRetention       float64 `json:"query_retention"`
    FallbackDepth        int     `json:"fallback_depth"`
    MatchSpecificity     float64 `json:"match_specificity"`
    DescriptionAlignment float64 `json:"description_alignment"`
    OriginConvergence    int     `json:"origin_convergence"`
    CapabilityCoverage   float64 `json:"capability_coverage"`
    NonGoalConflict      float64 `json:"non_goal_conflict"`
    EvidenceCompleteness float64 `json:"evidence_completeness"`
}
```

Definições mínimas:

- `ExpansionPrior`: produtividade julgada da expansão ou família de regra;
- `QueryRetention`: proporção ponderada dos termos informativos preservados entre
  query planejada e efetiva;
- `FallbackDepth`: zero para query completa, um para reduzida com linguagem, dois
  para reduzida sem linguagem;
- `MatchSpecificity`: peso explícito para match de repositório, arquivo ou símbolo;
- `DescriptionAlignment`: alinhamento lexical entre ideia e descrição;
- `OriginConvergence`: número de expansões independentes que encontraram o projeto;
- `CapabilityCoverage`: fração de capacidades com algum sinal preliminar;
- `NonGoalConflict`: presença de conceitos negativos;
- `EvidenceCompleteness`: quanto dos campos necessários está disponível.

Nenhum sinal pode realizar I/O. Toda normalização, tokenização, stopword e peso
devem estar especificados em L0 e testados. Campos ausentes não recebem o mesmo
valor de evidência negativa.

### 4. Baseline interpretável antes do probabilístico

Implementar primeiro um modelo linear explicável:

```text
ranking_score = clamp(intercept + Σ pesoᵢ × sinalᵢ, 0, 1)
```

Os pesos não podem nascer de intuição escondida no código. Há duas formas aceitas:

1. pesos declarados num arquivo versionado e rotulados como `uncalibrated`; ou
2. pesos ajustados sobre dados julgados, com protocolo, seed e artefato de treino
   registrados.

Cada resultado deve trazer contribuições:

```go
type SignalContribution struct {
    Signal string  `json:"signal"`
    Value  float64 `json:"value"`
    Weight float64 `json:"weight"`
    Effect float64 `json:"effect"`
    Reason string  `json:"reason"`
}
```

O usuário deve conseguir somar as contribuições e reproduzir o score.

Não usar modelo de linguagem, embedding ou classificador de caixa-preta neste
passo. O objetivo é criar o padrão de evidência contra o qual modelos futuros serão
comparados.

### 5. Calibração com validação fora da amostra

Preparar um dataset a partir dos snapshots e julgamentos, preservando grupo por
ideia. `uncertain` não vira negativo automaticamente.

Como existem apenas quatro ideias, usar validação `leave-one-idea-out`:

```text
treina em 3 ideias → avalia na ideia restante
repete para as 4 ideias
```

Isso impede que candidatos da mesma ideia apareçam simultaneamente em treino e
teste e reduz a ilusão causada por memorizar uma regra lexical.

Comparar pelo menos:

- baseline constante pela prevalência do treino;
- regra simples por profundidade de fallback;
- modelo linear com todos os sinais disponíveis.

Métricas:

```text
Brier score       = média((p - y)²)
log loss          = perda logarítmica com probabilidades limitadas
precision@k       = relevantes nos k primeiros / julgados nos k primeiros
coverage@threshold = candidatos encaminhados / candidatos totais
false-negative rate auditada
```

Também produzir tabela de calibração por faixas quando houver pelo menos cinco
observações julgadas na faixa. Faixas menores devem aparecer como insuficientes.

Se o modelo não superar o baseline de prevalência em Brier score fora da amostra,
o valor continua sendo `ranking_score`; `probability` permanece ausente.

### 6. Decisão de peneira com três saídas

Evitar decisão binária prematura:

```text
ADVANCE  → enviar para triagem profunda
AUDIT    → evidência insuficiente ou amostra de controle
DROP     → baixo sinal, sujeito a auditoria
```

O resultado deve conter:

```go
type ScreeningDecision struct {
    Repository        string               `json:"repository"`
    RankingScore      float64              `json:"ranking_score"`
    Probability       *float64             `json:"probability"`
    CalibrationStatus string               `json:"calibration_status"`
    Decision          string               `json:"decision"`
    Contributions     []SignalContribution `json:"contributions"`
    UncertaintyReasons []string             `json:"uncertainty_reasons"`
    Origins           []CandidateOrigin    `json:"origins"`
}
```

Política inicial:

- `NonGoalConflict` forte nunca é ocultado pelo total;
- ausência de evidência aumenta `AUDIT`, não `DROP` automático;
- todo relevante conhecido deve ser `ADVANCE` durante replay de diagnóstico, mas
  não pode ser usado para decidir o próprio score fora do treino;
- uma amostra determinística de `DROP` sempre vira `AUDIT` para medir falsos
  negativos;
- limiares e taxa de auditoria ficam em configuração versionada.

### 7. Boca de inspeção em L2

Adicionar comandos offline:

```bash
bateia idea screen \
  --snapshot snapshot-v2.json \
  --model screening-model.json \
  --output screening.json

bateia idea screening calibrate \
  --snapshots benchmark/ \
  --judgments benchmark/ \
  --output calibration.json

bateia idea screening audit \
  --screening screening.json \
  --judgments judgments.json
```

Requisitos:

- calibração e replay completamente offline;
- seed explícita quando houver amostragem;
- ordenação estável por score e identidade;
- `stdout` somente JSON quando não houver `--output`;
- erro acionável para schema, modelo, limiar ou julgamento inválido;
- nenhuma inicialização de Sourcegraph, GitHub, cache ou `.env`.

## Auditoria de falsos negativos

O benchmark anterior julgou por descrição. Para medir a peneira, revisar:

1. todos os oito relevantes conhecidos;
2. todos os cinco incertos;
3. pelo menos dez irrelevantes escolhidos deterministicamente entre os menores
   scores;
4. pelo menos dez irrelevantes próximos ao limiar de avanço.

Para essa auditoria, ler pelo menos README e descrição, registrando evidência e
referência. Se um `irrelevant` revelar capacidade central, alterar o julgamento com
histórico explícito e contar o caso como falso negativo da classificação anterior.

Não apagar julgamentos antigos. Criar revisão com `previous_label`, `new_label`,
`evidence` e `reason`.

## Prova de mordida

Criar quatro candidatos controlados:

```text
A: query completa + match específico + duas origens + descrição alinhada
B: fallback sem idioma + match genérico + descrição divergente
C: pouco conteúdo; nenhum conflito; evidência incompleta
D: score lexical alto, mas conflito explícito com NonGoal
```

O teste deve demonstrar:

- A fica acima de B por contribuições reproduzíveis;
- B explica que entrou por fallback genérico;
- C vira `AUDIT`, não `DROP` por ausência de dados;
- D expõe o conflito e não avança silenciosamente pelo score total;
- remover uma origem de A reduz somente os sinais dependentes de convergência;
- migrar snapshot v1 não fabrica path, snippet ou match type;
- duas execuções geram JSON idêntico;
- alterar um peso muda somente as contribuições dependentes dele.

Se o score não puder ser reconstruído pelas contribuições, o passo não mordeu.

## Testes obrigatórios

### Unidade — L1

1. migração v1 → v2 preserva desconhecidos como ausentes;
2. retenção de query diminui com perda de termos informativos;
3. profundidade de fallback é derivada da tentativa real;
4. matches duplicados não inflam evidência;
5. convergência conta expansões independentes;
6. ausência difere de conflito negativo;
7. contribuições somam o score antes do `clamp`;
8. ranking possui desempate determinístico;
9. precisão não calibrada não é serializada como probabilidade;
10. folds de validação não misturam a mesma ideia;
11. `uncertain` não entra silenciosamente como classe negativa;
12. Brier, log loss e precision@k tratam conjuntos vazios;
13. amostra de auditoria é reproduzível pela seed;
14. entrada, snapshot e modelo não são modificados.

### Integração — L2/L3/L4

1. Sourcegraph preserva tipo, path e tentativa responsável;
2. snippets são normalizados e truncados;
3. token e resposta bruta não aparecem no snapshot;
4. captura nova gera schema v2;
5. snapshots v1 continuam avaliáveis após migração;
6. calibração leave-one-idea-out gera artefato reproduzível;
7. `idea screen`, `calibrate` e `audit` funcionam com `env -i`;
8. comandos offline anteriores continuam sem infraestrutura remota;
9. golden literal do Passo 0001 permanece intacto;
10. auditoria registra revisões sem sobrescrever julgamento anterior.

## Critérios de verificação

- [x] Cada candidato possui trajetória causal até a tentativa efetiva.
- [x] Match type, path e snippet são preservados quando disponíveis.
- [x] Sinais são separados, documentados e puros.
- [x] Score pode ser reconstruído integralmente pelas contribuições.
- [x] Valor não calibrado não é chamado de probabilidade.
- [x] Validação fora da amostra agrupa por ideia.
- [x] Modelo é comparado ao baseline de prevalência.
- [x] `ADVANCE`, `AUDIT` e `DROP` tratam ausência e conflito distintamente.
- [x] Amostra de descartes é auditada para falsos negativos.
- [x] Resultados anteriores podem ser reproduzidos por snapshots.
- [x] Prova de mordida explica A, B, C e D.
- [x] Novos componentes possuem prompts L0 e linhagem válida.
- [x] `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.
- [x] Um laudo registra calibração, falsos negativos, limites e decisão.

## Como este passo pode falhar

- pesos intuitivos são apresentados como probabilidade;
- o modelo memoriza regras ou ideias presentes no treino;
- fallback efetivo continua desconectado do candidato;
- ausência de descrição é interpretada como irrelevância;
- correspondências repetidas inflam artificialmente a confiança;
- score total esconde conflito com `NonGoals`;
- auditoria revisa apenas casos fáceis ou favoráveis;
- limiar é escolhido depois de observar o teste;
- precisão melhora apenas porque quase tudo vira `AUDIT`;
- snapshots v1 recebem evidência fabricada na migração.

## Fora de escopo

- afirmar capacidade comprovada somente por descrição ou match textual;
- usar LLM, embedding ou reranker opaco;
- baixar código de todos os candidatos;
- calcular recall no universo aberto;
- automatizar julgamento humano;
- substituir a triagem do `SearchFunnel`;
- emitir decisão final de reuso;
- recomendar arquitetura ou integração entre projetos;
- calibrar para domínios não representados sem novos dados.

## Próximo passo condicionado

Se o modelo superar o baseline fora da amostra e a auditoria mostrar taxa aceitável
de falsos negativos, executar triagem profunda em:

1. todos os `ADVANCE`;
2. todos os `AUDIT` por incerteza;
3. amostra determinística dos `DROP`.

Se não houver calibração suficiente, usar apenas o ranking explicável para ordenar
o trabalho humano e ampliar os julgamentos. Se os fallbacks dominarem o ruído,
corrigir ou penalizar a política de fallback antes de sofisticar o modelo.
