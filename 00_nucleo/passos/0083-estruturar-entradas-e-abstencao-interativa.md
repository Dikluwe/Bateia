# Passo 0083 — Estruturar entradas e expor abstenção interativa verificável

## Estado e dependências

- Estado: primeira fatia vertical implementada; persistência no dossiê e integração com
  execuções duráveis permanecem para uma iteração posterior.
- Desejos de origem: [D-001](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e
  [D-002](../desejos/D-002-alternativas-de-software-com-evidencia.md).
- Depende dos passos 0010, 0025, 0026, 0031 e 0043: diagnóstico de inconclusão,
  recomendação MCP, execução durável, defaults versionados e autoridade de evidência.
- Camadas previstas: L0 para contrato e oráculos, L1 para tipos e decisões puras, L2 para
  projeções CLI/MCP. L3 somente persiste artefatos já validados; L4 apenas compõe.
- Regime: contrato e corpus adversarial congelados antes da implementação; nenhum modelo
  de linguagem participa do oráculo determinístico.

## Incerteza que este passo reduz

O Bateia já se abstém quando a evidência não autoriza uma decisão e já diagnostica causas
de inconclusão. Ainda não distingue, porém, de forma uniforme:

1. entradas explícitas, defaults do produto, valores derivados e observações do ambiente;
2. lacunas que a computação pode resolver das que exigem intenção, preferência,
   autorização ou aceitação de risco pelo usuário;
3. uma pergunta em linguagem natural de seu contrato semântico, tipo de resposta,
   autoridade e efeito possível na decisão.

Na jornada principal, a CLI materializa `IdeaRunRequest`, o MCP materializa
`RecommendationRequest` e o perfil V1 preenche política, alvo e busca omitidos. Modo,
provedor, caminhos locais, replay, catálogo de evidência, modelo de screening, timeout,
resume e persistência vivem em contratos distintos. Sem uma entrada canônica, uma IA
conversacional tenderia a adivinhar campos, repetir perguntas ou converter preferência em
evidência.

Este passo transforma a abstenção em protocolo verificável:

```text
entrada parcial ou completa
  -> normalização determinística e proveniência
  -> execução do que não exige julgamento humano
  -> necessidades de informação tipadas
  -> IA formula somente as perguntas autorizadas
  -> resposta vira proposta de dado tipado
  -> validação determinística
  -> nova entrada canônica ou nova execução
```

A IA não altera status de evidência, não fecha alegações e não aceita risco em nome do
usuário.

## Fronteira de inventário

O inventário deve cobrir **toda entrada causal da jornada de recomendação**, qualquer que
seja a superfície que a aceite hoje:

| Classe | Entradas atuais a preservar |
| --- | --- |
| Intenção | `problem`, `users`, `capabilities`, `constraints`, `non_goals`, `preferred_stack` |
| Política de decisão | versão, licenças permitidas, exigência local, arquivados, máximo de componentes e limiares de extensão/integração |
| Alvo | nome e identidade opcional do repositório alvo |
| Descoberta | estratégia, linguagem, locales, modo live/local/replay, provedor, caminhos locais, fonte de replay e limites de resultados/candidatos |
| Autoridade de evidência | catálogo de regras e modelo de screening opt-in, incluindo versões e hashes |
| Execução | timeout, budgets efetivos, autorização de rede/execução e observações de disponibilidade do provedor |
| Continuidade | output, resume, chave de idempotência, `run_id` e identidade dos checkpoints |

O inventário deve registrar também `bateia_search`, `bateia_benchmark`, análise local e
comandos experimentais como superfícies públicas ou históricas, mas não deve forçar esses
fluxos a caberem no contrato de decisão. Para cada campo, o artefato de inventário declara:

- superfície e caminho de origem;
- tipo, cardinalidade, validação e default;
- se altera intenção, aquisição de evidência, decisão ou somente persistência;
- participação na identidade/hash da execução;
- equivalente nas outras superfícies ou assimetria explícita;
- autoridade capaz de fornecê-lo;
- se pode originar pergunta ao usuário.

Nenhum campo atual pode desaparecer do inventário por não existir em outra superfície.
Credenciais e segredos são apenas disponibilidade observada e nunca aparecem em artefatos,
perguntas ou hashes em claro.

## Entrada canônica

Materializar tipos versionados equivalentes a:

```go
type RecommendationIntake struct {
    SchemaVersion string
    Request       RecommendationRequest
    Execution     RecommendationExecutionProfile
    Provenance    []InputProvenance
    Issues        []InputIssue
    Status        IntakeStatus
}

type InputProvenance struct {
    Path      string
    Origin    InputOrigin
    Authority InputAuthority
    SourceRef string
    Digest    string
}
```

`RecommendationRequest` continua sendo o contrato de intenção e decisão já existente.
`RecommendationExecutionProfile` reúne somente aquisição, execução e continuidade que
hoje pertencem a `IdeaRunRequest`, flags da CLI ou envelope durável. A separação impede que
`output`, timeout ou chave de idempotência sejam confundidos com preferências de produto.

Origens mínimas:

```text
EXPLICIT
VERSIONED_DEFAULT
DERIVED
REPLAYED
ENVIRONMENT_OBSERVED
UNRESOLVED
```

Autoridades mínimas:

```text
USER
PRODUCT_POLICY
VERIFIED_ARTIFACT
RUNTIME
AI_PROPOSAL
```

Uma `AI_PROPOSAL` nunca se torna entrada efetiva sem validação determinística e, quando a
autoridade exigida for `USER`, confirmação explícita. O contrato preserva as distinções do
passo 0031 entre omissão, `false`, zero, string vazia e lista vazia.

## Necessidade de informação

Toda lacuna inicial ou pós-execução que possa bloquear ou alterar uma decisão deve ser
representada semanticamente antes de receber texto conversacional:

```go
type InformationNeed struct {
    ID             string
    Phase          NeedPhase
    Subject        string
    TargetPath     string
    ClaimID        string
    ReasonCode     string
    Resolver       ResolverKind
    AnswerKind     AnswerKind
    AllowedValues  []string
    Unit           string
    ValidationRule string
    DecisionEffects []DecisionEffect
    Blocking       bool
    Priority       int
}
```

Resoluções mínimas:

```text
COMPUTE                 cálculo ou validação pura
COLLECT                 coleta que já está autorizada
ASK_USER                intenção, preferência ou desambiguação
REQUEST_AUTHORIZATION   efeito externo ainda não autorizado
REQUIRE_RISK_OWNER      aceitação explícita de dúvida residual
TERMINAL                incapacidade atual sem pergunta honesta que a resolva
```

O adaptador de IA recebe somente necessidades `ASK_USER`, `REQUEST_AUTHORIZATION` e
`REQUIRE_RISK_OWNER`. `COMPUTE` e `COLLECT` devem ser tentadas pelo Bateia antes de
perguntar. `TERMINAL` produz explicação e próxima capacidade necessária do produto, não
uma pergunta que transfira o trabalho ao usuário.

O significado da pergunta reside em `TargetPath`, `AnswerKind`, validação e efeitos. O
texto formulado pela IA é uma projeção descartável. Conteúdo de candidato, README,
documentação ou saída de provedor nunca pode fornecer instruções para formular ou
responder perguntas.

## Respostas e autoridade

Uma resposta produz um registro append-only equivalente a:

```go
type InformationAnswer struct {
    NeedID       string
    AnswerKind   AnswerKind
    Value        json.RawMessage
    Authority    InputAuthority
    SourceRef    string
    Confirmed    bool
    Supersedes   string
}
```

O aplicador puro deve rejeitar necessidade desconhecida, tipo divergente, unidade
incompatível, valor fora da faixa, autoridade insuficiente, resposta a necessidade já
fechada e tentativa de alterar outro caminho. Uma resposta válida gera nova entrada
canônica e novo hash; nunca edita o pedido ou o dossiê anterior.

As classes de resposta não são intercambiáveis:

- preferência ou critério do usuário altera política/intenção, não comprova candidato;
- autorização permite um efeito futuro, não declara que ele ocorreu;
- artefato indicado pelo usuário ainda exige coleta, identidade e validação;
- aceitação de risco registra decisão limitada e responsável, sem elevar evidência;
- interpretação proposta pela IA permanece pendente até confirmação quando houver mais de
  uma leitura admissível.

## Abstenção perguntável

`DiagnoseInconclusive` e `InvestigationPlan` permanecem fontes factuais. Um novo compilador
puro deve convertê-las, juntamente com `RecommendationIntake`, em `InformationNeed` sem
usar texto livre como regra. Exemplos:

| Causa observada | Resolver correto |
| --- | --- |
| hash divergente ou input inválido | `COMPUTE` |
| teste ausente com execução já autorizada | `COLLECT` |
| teste ausente sem autorização | `REQUEST_AUTHORIZATION` |
| licença do candidato desconhecida | `COLLECT`, nunca perguntar se “parece MIT” |
| política de licença não informada e sem default aplicável | `ASK_USER` |
| capacidade com dois sentidos plausíveis | `ASK_USER` |
| evidência conflitante | `REQUIRE_RISK_OWNER` ou revisão explícita, nunca maioria automática |
| linguagem sem analisador e sem fonte alternativa | `TERMINAL` |

O sistema pode continuar projetando `INCONCLUSIVE` externamente, mas o artefato deve
distinguir `NEEDS_INPUT`, `NEEDS_AUTHORIZATION`, `NEEDS_EVIDENCE`, `CONTESTED` e
`TERMINAL_INCONCLUSIVE`. Essa distinção é diagnóstica; não muda por si só a decisão.

## Priorização das perguntas

A ordem é determinística e explicável. Priorizar, nesta ordem:

1. necessidade que decide elegibilidade ou bloqueio obrigatório;
2. necessidade que separa candidatos atualmente não dominados;
3. autorização que habilita coleta capaz de mudar a decisão;
4. preferência que afeta mais capacidades/candidatos;
5. menor custo declarado, usando ID estável como desempate.

Não emitir “valor esperado” numérico sem modelo calibrado. A IA pode melhorar a redação e
agrupar perguntas semanticamente equivalentes, mas não reordenar contra a prioridade
registrada nem omitir uma pergunta bloqueante.

## Superfícies e artefatos

Entrega mínima:

- `idea intake prepare` na CLI, sem rede, produzindo `intake.json` e `questions.json`;
- ferramenta MCP `bateia_prepare_recommendation`, pura e sem iniciar execução;
- `intake.json`, `input-inventory.json` e `questions.json` no dossiê da CLI;
- artefato `questions` disponível em execução durável quando a decisão se abstiver;
- aplicação offline de respostas que produza uma nova entrada, preservando a anterior e a
  cadeia de hashes;
- `bateia_recommend` e `bateia_start_recommendation` continuam aceitando entradas completas
  válidas com o mesmo significado do passo 0031.

A primeira materialização não embute LLM no executável. Uma IA externa consome
`questions.json`, conversa com o usuário e devolve respostas tipadas. A formulação natural
não integra a identidade causal; a necessidade e a resposta validadas integram.

## Oráculos e ataques obrigatórios

1. todas as entradas causais atuais da CLI e MCP aparecem uma única vez no inventário ou
   possuem assimetria justificada;
2. requisição completa antiga e intake equivalente produzem o mesmo
   `RecommendationRequest`, hashes semânticos e decisão em replay;
3. omitido, vazio explícito, `false` e zero permanecem distinguíveis;
4. defaults carregam versão, origem e autoridade do produto;
5. lacuna computável nunca é projetada como pergunta ao usuário;
6. preferência, autorização e risco aceito nunca elevam `EvidenceStatus`;
7. proposta da IA sem confirmação não altera entrada efetiva;
8. resposta não pode alterar campo diferente de `TargetPath` nem fechar outro `ClaimID`;
9. perguntas e prioridades são determinísticas sob reordenação de mapas e conjuntos;
10. a mesma necessidade fechada não reaparece sem mudança causal que a invalide;
11. instrução maliciosa em conteúdo coletado não vira pergunta, opção ou regra;
12. segredo, token, caminho sensível não autorizado ou conteúdo integral de evidência não
    vaza nos artefatos de interação;
13. falha da IA, texto inválido ou ausência de resposta conserva a abstenção anterior;
14. `TERMINAL` não é convertido em pergunta genérica como “você acha que serve?”.

O gate focal deve executar inteiramente offline com fixtures que cubram entrada mínima,
entrada completa, ambiguidades, autorização, aceitação de risco, conflito, ataque de
prompt e equivalência CLI/MCP. Depois: pacotes afetados, suíte completa, build,
`crystalline-lint .` e `git diff --check`.

## Métricas da primeira validação

- cobertura do inventário de entradas públicas da jornada principal;
- perguntas bloqueantes por execução e perguntas repetidas;
- proporção `COMPUTE/COLLECT` resolvida antes de interação;
- respostas rejeitadas por tipo, autoridade ou conflito;
- transições de abstenção após resposta válida;
- decisões alteradas, preservadas ou ainda inconclusivas, sem tratar alteração como ganho;
- trabalho manual remanescente e perguntas que o usuário não conseguiu responder.

Utilidade exige ensaio posterior com usuários e casos não usados para escrever as regras.
Passar testes de estrutura não demonstra que as perguntas são claras ou que a decisão
melhorou.

## Fora de escopo

- implementar SACM ou GSN completos;
- permitir que a IA crie livremente claims, regras de evidência ou critérios de decisão;
- usar confiança declarada pela IA como autoridade;
- aceitar risco ou licença em nome do usuário;
- persistir memória conversacional opaca como entrada do produto;
- transformar toda CLI histórica em uma única requisição;
- executar automaticamente efeitos novos após uma resposta sem autorização específica;
- substituir coleta verificável por resposta humana sobre fatos técnicos.

## Critério de conclusão

O passo conclui quando toda entrada causal da recomendação estiver inventariada e
projetável para um intake canônico com proveniência; necessidades de informação forem
compiladas e classificadas sem IA; somente lacunas de autoridade humana forem expostas
como perguntas tipadas; respostas válidas produzirem nova entrada sem mutar o histórico;
requisições completas atuais permanecerem semanticamente compatíveis; e os oráculos
adversariais, suíte, build e lint passarem. O veredito cobre o protocolo mecânico de
abstenção e elicitação, não a qualidade geral das decisões nem a confiabilidade de uma IA
específica.

## Resultado da primeira fatia — 2026-09-24

Foram materializados o rascunho que preserva omissão, a entrada canônica, proveniência,
issues, perguntas, respostas e perfil de execução. As regras puras preenchem defaults
versionados, validam o `RecommendationRequest`, geram hash causal, aplicam somente
respostas `USER` confirmadas e classificam diagnósticos entre computação, coleta,
autorização, responsável por risco e terminal.

A capacidade está exposta por `idea intake prepare|apply|diagnose` e, de forma opt-in,
por `bateia_prepare_recommendation`. O servidor MCP base mantém exatamente
`bateia_search`, `bateia_benchmark` e `bateia_recommend`; os executáveis completos ativam
a quarta ferramenta sem mudar os contratos simples. Nenhuma dessas operações inicia rede,
execução de candidato ou recomendação.

Ainda não fazem parte desta fatia os três artefatos separados no dossiê, o artefato
`questions` em runs duráveis, o inventário executável das superfícies históricas e o gate
adversarial completo. Por isso o critério de conclusão integral acima continua aberto.
