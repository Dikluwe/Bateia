# Passo 0007 — Diagnosticar e ampliar cobertura de evidências

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0006 — Endurecer e calibrar evidências`](0006-endurecer-e-calibrar-evidencias.md)
**Onde roda**: Bateia
**Camadas previstas**: L0, L1, L2 e L3
**Natureza**: diagnóstico de cobertura e ampliação conservadora do corpus

## Intenção

Explicar por que capacidades que um humano reconhece em projetos reais ainda
terminam como `PARTIAL` ou `INCONCLUSIVE`, e corrigir a cobertura responsável por
essas perdas sem reduzir o rigor de `PROVEN`.

O passo não existe para fabricar dez resultados positivos. Ele deve distinguir
falta de evidência, falta de coleta, falta de suporte estrutural e erro de regra.
Somente correções gerais, reproduzíveis e protegidas por controles negativos podem
alterar a classificação.

## Situação inicial

A calibração real do Passo 0006 avaliou 22 candidatos em quatro ideias e emitiu
quatro capacidades `PROVEN`. Nas ideias em que a precisão foi mensurável, ela foi
de 100%, com falsa ausência de 0%, rastreabilidade de 100% e conclusão de orçamento
de 100%.

O gate exige pelo menos dez `PROVEN` revisados. Os casos `stackus/edat` e
`reaper47/recipya` demonstram a lacuna: são positivos intuitivos, mas o recorte
coletado trouxe somente evidência documental ou estrutural insuficiente. Promovê-los
por exceção destruiria o significado do estado `PROVEN`.

## Restrição central: sem LLM e sem afrouxar o gate

Permanecem proibidos no pipeline:

- LLM, embeddings, rerankers neurais ou classificação por prompt;
- regras especiais por nome de repositório;
- aumento de força porque um candidato “parece correto”;
- contagem de README isolado como implementação comprovada;
- exclusão de controles negativos que passaram a falhar;
- alteração retroativa de bundles, julgamentos ou versões publicadas.

Os limiares do Passo 0006 ficam congelados durante o diagnóstico. Se uma revisão
de limiar se mostrar necessária, ela exige hipótese documentada, nova versão de
regras e replay integral dos positivos e controles.

## Taxonomia obrigatória de falhas

Cada capacidade humana esperada como presente e não emitida como `PROVEN` recebe
uma causa primária:

```text
TARGET_NOT_PLANNED
TARGET_NOT_COLLECTED
BUDGET_EXHAUSTED
LANGUAGE_UNSUPPORTED
PATTERN_MISS
SYMBOL_LINKAGE_MISS
STRENGTH_UNDERESTIMATED
INSUFFICIENT_INDEPENDENCE
NEGATIVE_CONFLICT
GROUND_TRUTH_AMBIGUOUS
```

Definições:

- `TARGET_NOT_PLANNED`: o arquivo ou região que comprova a capacidade não entrou
  no plano;
- `TARGET_NOT_COLLECTED`: foi planejado, mas não chegou ao bundle;
- `BUDGET_EXHAUSTED`: limites impediram terminar alvos relevantes;
- `LANGUAGE_UNSUPPORTED`: o conteúdo existe, mas não há extrator conservador;
- `PATTERN_MISS`: regra conhecida não reconheceu vocabulário observável;
- `SYMBOL_LINKAGE_MISS`: o termo foi encontrado, mas não ligado ao símbolo que o
  implementa;
- `STRENGTH_UNDERESTIMATED`: a classe da fonte recebeu força incompatível com sua
  natureza verificável;
- `INSUFFICIENT_INDEPENDENCE`: existem sinais, mas não em fontes independentes
  suficientes;
- `NEGATIVE_CONFLICT`: negação contextual impede a promoção;
- `GROUND_TRUTH_AMBIGUOUS`: a revisão humana não sustenta um esperado inequívoco.

Ausência de diagnóstico produz `GROUND_TRUTH_AMBIGUOUS`; nunca autoriza correção
automática.

## Contratos

```go
type EvidenceMissCause string

type EvidenceMissDiagnosis struct {
    Repository     string            `json:"repository"`
    Capability     string            `json:"capability"`
    Expected       EvidenceStatus    `json:"expected"`
    Observed       EvidenceStatus    `json:"observed"`
    PrimaryCause   EvidenceMissCause `json:"primary_cause"`
    SupportingData []EvidenceSource  `json:"supporting_data,omitempty"`
    MissingTargets []string          `json:"missing_targets,omitempty"`
    RuleID         string            `json:"rule_id,omitempty"`
    Explanation    string            `json:"explanation"`
}

type EvidenceCoverageReport struct {
    SchemaVersion string                  `json:"schema_version"`
    RulesVersion  string                  `json:"rules_version"`
    RulesHash     string                  `json:"rules_hash"`
    Diagnoses     []EvidenceMissDiagnosis `json:"diagnoses"`
    CauseCounts   map[EvidenceMissCause]int `json:"cause_counts"`
}
```

O relatório é offline, determinístico e ordenado por ideia, repositório e
capacidade. Toda causa aponta para bundle, alvo, regra ou ausência verificável.

## Boca de inspeção

```bash
bateia idea evidence diagnose \
  --assessment assessment.json \
  --evidence bundles.json \
  --judgments judgments.json \
  --output coverage-report.json
```

O comando não acessa rede nem modifica o corpus.

## Execução em duas fases

### Fase A — Auditoria congelada

1. congelar regras, bundles, assessments e julgamentos atuais;
2. diagnosticar todos os esperados `PROVEN` observados como outro estado;
3. revisar manualmente as fontes concretas sem alterar regras;
4. publicar contagem por causa e exemplos reproduzíveis;
5. escolher correções pela frequência e impacto transversal, não pelo candidato.

### Fase B — Ampliação controlada

Para cada correção aceita:

1. escrever teste positivo com fixture mínima;
2. escrever ataque negativo que a generalização não pode aceitar;
3. implementar nova versão de regra, extrator ou planejamento;
4. reexecutar todo o corpus anterior;
5. rejeitar a correção se criar falsa ausência ou falso `PROVEN`;
6. coletar candidatos adicionais somente depois do replay passar;
7. manter lado a lado as métricas anterior e nova.

## Corpus mínimo adicional

Adicionar casos suficientes para cobrir, sem duplicar a mesma evidência:

- implementação comprovada por teste e símbolo;
- biblioteca cuja capacidade aparece em API pública;
- ferramenta cuja capacidade aparece em exemplo executável;
- projeto positivo numa linguagem ainda não coberta estruturalmente;
- documentação correta sem implementação, como controle `PARTIAL`;
- README que menciona ferramenta externa, como controle negativo;
- código vendorizado e fixture, como controles negativos;
- coleta parcial, timeout e estouro de budget;
- ao menos um componente reutilizável revisado e um falso componente potencial.

Os casos devem atravessar pelo menos quatro ideias e três ecossistemas de
linguagem. Repetir a mesma capacidade do mesmo repositório com nomes equivalentes
não aumenta artificialmente a amostra.

## Métricas

Publicar antes e depois:

```text
expected_proven_count
emitted_proven_count
proven_recall
evidence_precision
false_absence_rate
inconclusive_rate
source_trace_rate
budget_completion
component_precision
miss_count_by_cause
```

`proven_recall` torna visível a cobertura que as métricas do Passo 0006 ainda não
penalizam diretamente:

```text
PROVEN corretos emitidos / capacidades julgadas como PROVEN
```

Precisão e recall permanecem separados. Melhorar recall não justifica degradar
precisão.

## Critérios de aceitação

- todo esperado `PROVEN` não emitido possui causa primária revisável;
- ao menos dez capacidades independentes são emitidas como `PROVEN` e revisadas;
- `evidence_precision` ≥ 90%;
- `false_absence_rate` = 0;
- `source_trace_rate` = 100%;
- `component_precision` = 100% na primeira amostra revisada;
- nenhuma métrica usa aliases equivalentes para inflar a amostra;
- toda regra nova passa por positivo, ataque negativo e replay integral;
- nenhuma regra depende de nome de repositório;
- bundles, julgamentos, regras e relatórios ficam versionados por hash;
- `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.

## Testes obrigatórios

1. alvo relevante ausente do plano produz `TARGET_NOT_PLANNED`;
2. alvo planejado com falha produz `TARGET_NOT_COLLECTED`;
3. budget interrompido produz `BUDGET_EXHAUSTED`;
4. linguagem sem extrator produz `LANGUAGE_UNSUPPORTED`;
5. termo conhecido sem padrão produz `PATTERN_MISS`;
6. match distante de símbolo produz `SYMBOL_LINKAGE_MISS`;
7. sinais independentes insuficientes permanecem abaixo de `PROVEN`;
8. julgamento ambíguo não propõe mudança de regra;
9. ordem das entradas não altera diagnóstico ou métricas;
10. aliases de uma capacidade contam uma única unidade amostral;
11. regressão em controle negativo rejeita nova regra;
12. diagnóstico funciona sem rede e com ambiente vazio.

## Falhas esperadas

- ground truth insuficiente para decidir se o projeto realmente possui a
  capacidade;
- API remota não entrega arquivos necessários dentro do budget;
- linguagem ou construção não suportada pelo extrator;
- ganho de recall acompanhado de perda de precisão;
- dez positivos disponíveis, mas não independentes;
- nenhum componente concreto suficiente para medir `component_precision`.

Essas falhas mantêm o passo aberto; não são convertidas em sucesso por exceção.

## Fora de escopo

- decidir `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD`;
- executar código de terceiros;
- suportar todas as linguagens;
- estimar esforço em horas;
- interface gráfica;
- LLM, embeddings ou inferência generativa.

## Próximo passo condicionado

O [`Passo 0008`](0008-comparar-candidatos-e-decidir-reuso.md) construirá a matriz
comparativa e o motor determinístico de decisão somente depois que este passo
demonstrar, por replay versionado, que o gate de evidências foi atingido.
