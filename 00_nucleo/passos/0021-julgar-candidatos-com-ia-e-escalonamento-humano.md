# Passo 0021 — Julgar candidatos com IA e escalonamento humano

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0020 — Executar busca autenticada com orçamento observável`](0020-executar-busca-autenticada-com-orcamento-observavel.md)
**Onde roda**: avaliação externa assistida por IA + revisão humana; métricas no Bateia em replay offline
**Camadas previstas**: L0, L1 e L2
**Natureza**: oráculo assistido, cego e versionado; IA fora do pipeline determinístico

## Incerteza que este passo reduz

Não sabemos se os resultados encontrados correspondem de fato à ideia nem se uma
IA de desenvolvimento consegue triar sua qualidade com concordância suficiente
para reduzir o trabalho humano sem fabricar certeza.

## Princípio de autoridade

A IA emite o primeiro veredito. O usuário confirma dúvidas e casos de borda. O
julgamento humano prevalece e toda substituição permanece auditável.

```text
snapshot selado
  → candidatos anonimizados
  → julgamento inicial da IA
  → escalonamento determinístico
  → confirmação ou substituição humana
  → oráculo final selado
  → métricas offline
```

A IA não recebe braço, query, rank original, métrica parcial, registry protegido
ou expectativa de promoção. Ela não executa nova busca nem altera a rubrica.

## Rubrica congelada

Cada candidato recebe:

```text
RELEVANT · PARTIAL · IRRELEVANT · UNCERTAIN · CONFLICT
```

O registro inclui: correspondência ao problema, capacidades sustentadas,
restrições conflitantes, fontes citadas, justificativa curta, confiança de 0 a 1
e versão da rubrica. Afirmação sem fonte não pode elevar `PARTIAL` para
`RELEVANT`.

## Escalonamento humano obrigatório

Escalonar quando ocorrer qualquer condição:

- veredito `UNCERTAIN` ou `CONFLICT`;
- confiança menor que `0,80`;
- evidências citadas não sustentam a justificativa;
- candidato está no limite do top 10 ou altera recall esperado;
- o julgamento mudaria o braço vencedor ou a promoção da gramática;
- identidade, URL ou natureza repo-backed é ambígua;
- amostra determinística de 10% dos casos restantes para controle de excesso de confiança.

O usuário registra `CONFIRMED`, `OVERRIDDEN` ou `RESOLVED`, com motivo obrigatório
em substituições. Caso escalonado sem resposta permanece `UNKNOWN`, nunca herda o
veredito da IA silenciosamente.

## Contratos públicos

Materializar tipos versionados equivalentes a:

```go
type AssistedJudgment struct {
    CandidateID    string
    Verdict        string
    Confidence     float64
    Reasons        []string
    Evidence       []EvidencePointer
    RubricVersion  string
}

type HumanReview struct {
    CandidateID string
    Action      string
    Verdict     string
    Reason      string
}
```

O oráculo final deve preservar julgamento original, revisão e veredito efetivo;
não sobrescrever histórico.

## Regime Tekt

- autor da rubrica não lê o braço de origem;
- executor da IA recebe somente candidatos cegos e rubrica selada;
- revisor humano não vê métricas agregadas antes de resolver escalonamentos;
- verificador calcula concordância e métricas a partir dos artefatos selados;
- adversário testa evidência inexistente, confiança extrema, identidade ambígua,
  instruções maliciosas vindas de README e tentativa de revelar o braço.

Conteúdo de repositório é dado não confiável. Instruções nele contidas nunca
alteram a tarefa, a rubrica ou as ferramentas disponíveis ao julgador.

## Métricas e gates

- 100% dos julgamentos citam evidência presente no snapshot;
- 100% dos casos escalonados possuem revisão ou ficam `UNKNOWN`;
- concordância bruta IA/humano de pelo menos 80% na amostra revisada;
- taxa de override menor ou igual a 20%;
- zero vazamento de braço/query/rank no pacote cego;
- zero instrução de candidato obedecida como comando;
- precisão julgada em 10 e recall repo-backed calculáveis;
- mutation score `1,0` para violações da rubrica e do cegamento;
- repetição offline produz o mesmo oráculo efetivo e as mesmas métricas.

Concordância insuficiente não autoriza reduzir a amostra: produz
`FAIL_AI_RELIABILITY` e exige ampliar revisão humana.

## Veredito

- gates passam: `ORACLE_ACCEPTED`;
- muita divergência: `FAIL_AI_RELIABILITY`;
- revisão humana incompleta: `INCONCLUSIVE_HUMAN_REVIEW`;
- cegamento violado: `FAIL_BLINDING`;
- fontes insuficientes: `FAIL_EVIDENCE_GROUNDING`.

## Fora de escopo

- colocar LLM no executável Bateia;
- permitir que a IA altere queries, regras ou expected outcomes;
- julgamento jurídico de licença;
- substituir revisão humana por confiança declarada pela própria IA;
- medir capacidades em código; este passo julga qualidade do retorno da descoberta.

## Critério de conclusão

O passo termina com um oráculo final selado, histórico IA/humano preservado e
métricas de qualidade da descoberta calculáveis sem rede.
