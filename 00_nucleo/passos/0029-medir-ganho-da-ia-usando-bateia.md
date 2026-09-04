# Passo 0029 — Medir o ganho da IA usando o Bateia

**Estado**: `PLANEJADO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Depende de**: [`Passo 0028 — Extrair mecanismos e formular hipóteses de melhoria`](0028-extrair-mecanismos-e-formular-hipoteses-de-melhoria.md)
**Baseline obrigatório**: certificado `MECHANISM_IMPROVEMENT_SYNTHESIS_READY` e refinamentos preservados de 0025–0027
**Onde roda**: sessões externas limpas de agentes, harness experimental segregado e Bateia MCP certificado
**Camadas previstas**: L0 e benchmark; nenhuma mudança de produção é necessária na primeira materialização
**Natureza**: experimento A/B fatorial, cego e limitado por orçamento

## Problema

Os Passos 0025–0028 tornaram o Bateia uma ferramenta utilizável por uma IA
externa. Ela pode descobrir alternativas, acompanhar execuções duráveis, ampliar
a busca entre idiomas e mercados, extrair mecanismos e formular hipóteses de
melhoria com evidência.

Isso ainda não prova a hipótese central do produto:

> uma IA usando o Bateia toma decisões de engenharia melhores do que a mesma IA
> trabalhando sozinha?

Também não sabemos se o benefício depende da capacidade do agente. O Bateia
pode ajudar mais um agente rápido, compensando investigação superficial; pode
ajudar agentes mais deliberativos a escapar de soluções convencionais; ou pode
apenas acrescentar custo e complexidade a agentes que já chegariam ao mesmo
resultado.

Testar um único modelo ou um único esforço confundiria efeito do Bateia com
características daquele perfil. Comparar modelos diferentes sem um controle
pareado também produziria uma disputa de modelos, não uma avaliação da
ferramenta.

## Incerteza que este passo reduz

Não sabemos se o Bateia melhora, de forma reproduzível e economicamente
observável:

- qualidade da arquitetura proposta;
- descoberta de mecanismos úteis fora do baseline convencional;
- rastreabilidade entre afirmação e evidência;
- decisão de adaptar, reutilizar, investigar ou construir;
- explicitação de riscos, premissas, contraevidências e `Unknown`;
- qualidade do plano de validação;
- autonomia de um agente aberto em um chat novo.

Também não sabemos como esse efeito varia entre famílias de agente e níveis de
esforço, nem se o custo adicional compensa o ganho.

## Objetivo

Executar o mesmo corpus em uma matriz pareada de agentes externos:

```text
Luna  × low, medium, high
Terra × low, medium, high
Sol   × low, medium, high

cada perfil executa:
  A — IA sozinha, sem acesso ao Bateia
  B — IA com acesso ao Bateia e instrução zero-shot
```

O experimento mede o efeito `B − A` dentro do mesmo caso, modelo e esforço. Ele
também estima interações por família e esforço, sem pressupor que os nomes
representem uma ordem universal de inteligência.

O resultado esperado não é declarar “qual modelo é melhor”. É responder:

1. o Bateia produz ganho material sobre a IA isolada?
2. esse ganho aparece em perfis diferentes?
3. em quais perfis o Bateia ajuda, é neutro ou atrapalha?
4. quanto tempo, uso e complexidade operacional o ganho custa?
5. uma instrução curta em um chat novo basta para o agente operar a ferramenta?

## Perfis de agente congelados

Na primeira materialização, o host experimental deve resolver e congelar
exatamente estes IDs:

| Rótulo experimental | Model ID resolvido | Esforços |
|---|---|---|
| `LUNA` | `gpt-5.6-luna` | `low`, `medium`, `high` |
| `TERRA` | `gpt-5.6-terra` | `low`, `medium`, `high` |
| `SOL` | `gpt-5.6-sol` | `low`, `medium`, `high` |

Mapeamento para apresentação em português:

- `low` = leve;
- `medium` = médio;
- `high` = alto.

Esses IDs são perfis expostos pelo ambiente atual, não aliases universais
presumidos para qualquer API ou host. Antes do experimento, registrar:

```go
type AgentProfile struct {
    ID                    string `json:"id"`
    FamilyLabel           string `json:"family_label"`
    ResolvedModelID       string `json:"resolved_model_id"`
    ReasoningEffort       string `json:"reasoning_effort"`
    HostID                string `json:"host_id"`
    AvailabilitySnapshot string `json:"availability_snapshot"`
    SnapshotSHA256        string `json:"snapshot_sha256"`
    ResolvedAt            string `json:"resolved_at"`
}
```

Cada combinação modelo × esforço forma um perfil distinto, totalizando nove
perfis. Não substituir silenciosamente um modelo indisponível por outro, nem
converter esforço não suportado ao default. A célula recebe
`PROFILE_UNAVAILABLE`, e a conclusão da matriz completa fica inconclusiva.

Versão do modelo, host, defaults não configuráveis e parâmetros efetivos devem
ser registrados quando o runtime os expuser. Campo não exposto permanece
`UNKNOWN_RUNTIME_DEFAULT`; não é inferido pelo nome comercial.

## Matriz fatorial

Para cada ideia existem 18 células:

```text
3 famílias × 3 esforços × 2 braços = 18 execuções
```

Os braços são:

### Braço A — `AGENT_ALONE`

O agente recebe a ideia, o snapshot permitido do alvo e o contrato de saída. Não
recebe Bateia, artifacts, nomes de candidatos, respostas de outros agentes ou
acesso indireto ao store.

### Braço B — `AGENT_WITH_BATEIA`

O agente recebe a mesma entrada e os mesmos limites, acrescidos apenas de:

> Use o Bateia como ferramenta para investigar alternativas, mecanismos e
> hipóteses antes de concluir.

As dez ferramentas certificadas do Bateia ficam disponíveis. Não é entregue
tutorial, sequência de chamadas, resposta esperada nem artifact pré-selecionado.
Descobrir a superfície, iniciar runs, acompanhar estado e ler artifacts faz
parte da utilidade observada.

A diferença de prompt é parte constitutiva do tratamento. Nenhum dos braços é
informado de que existe uma resposta pareada.

## O que significa “IA sozinha”

No experimento causal primário, os dois braços compartilham a mesma allowlist de
ferramentas gerais. Essa allowlist não inclui busca web, conector de código,
memória de outro chat ou ferramenta de descoberta concorrente. O braço B recebe
somente o Bateia como capacidade adicional.

Isso compara conhecimento e raciocínio do agente com conhecimento e raciocínio
apoiados pelo Bateia. Uma trilha ecológica com web e outras ferramentas pode ser
executada depois, mas não participa do primeiro selo porque dobraria custo e
introduziria mudança temporal de fontes.

O Bateia usa snapshots e replay congelados no experimento primário. Nenhum
agente recebe vantagem por executar em horário ou estado de rede diferente.

## Sessões limpas e pareamento temporal

Cada célula roda em um chat ou tarefa nova, sem histórico herdado. A sessão:

- recebe somente prompt congelado, caso, contrato de saída e tools do braço;
- não monta o workspace que contém oráculos ou respostas anteriores;
- não pode ler arquivos de outro perfil;
- não recebe resumo de sessão anterior;
- não continua conversa após a resposta final;
- não recebe ajuda humana durante a execução zero-shot.

Os braços A e B do mesmo perfil/caso são despachados em ordem aleatória e na
mesma janela temporal. Quando possível, executam concorrentemente. Ordem de
modelo, esforço, caso e braço deriva de seed congelada antes das respostas.

Interrupção ou indisponibilidade do usuário não injeta texto na sessão. Pedido de
esclarecimento do agente termina como `HUMAN_INPUT_REQUIRED`; não é respondido
no experimento primário.

## Contrato comum de saída

Todos os agentes devem responder no mesmo schema sem receber conteúdo esperado:

```go
type EngineeringInvestigation struct {
    ProblemRestatement    string                 `json:"problem_restatement"`
    ConventionalBaseline []string               `json:"conventional_baseline"`
    Alternatives         []InvestigatedOption   `json:"alternatives"`
    Mechanisms           []InvestigatedMechanism `json:"mechanisms"`
    ImprovementHypotheses []InvestigatedHypothesis `json:"improvement_hypotheses"`
    ReuseBuildDecision   string                 `json:"reuse_build_decision"`
    EvidenceReferences   []AgentEvidenceRef     `json:"evidence_references"`
    CounterEvidence      []string               `json:"counter_evidence"`
    Risks                []string               `json:"risks"`
    Unknowns             []string               `json:"unknowns"`
    ValidationPlan       []AgentValidationStep  `json:"validation_plan"`
}
```

O schema padroniza o observável, não o raciocínio. Campo ausente é pontuado como
ausente; o harness não completa ou reescreve a resposta.

No braço B, referências a `run_id`, artifact e hash são verificadas
mecanicamente contra o store. No braço A, referências públicas fornecidas pelo
próprio agente são preservadas como alegações não verificadas, pois busca externa
está desabilitada. Não conceder ao controle evidência sintética equivalente ao
Bateia: isso removeria o tratamento que se deseja medir.

## Corpus

Congelar seis ideias reais, sem nome de concorrente, distribuídas por ao menos
três famílias de software e contendo:

- problema e usuários;
- capacidades desejadas sem prescrever arquitetura;
- ao menos uma solução convencional plausível;
- ao menos um mecanismo relevante que o corpus Bateia consegue sustentar;
- casos em que o resultado correto pode ser inconclusivo;
- ao menos duas ideias com sinais multilíngues ou de outro mercado;
- ao menos duas ideias com alvo repo-backed e lacunas observáveis;
- ao menos uma ideia em que construir do zero é defensável.

Divisão:

| Split | Ideias | Execuções | Uso |
|---|---:|---:|---|
| `CONTRACT` | 1 | 18 | smoke test da matriz e instrumentação |
| `VALIDATION` | 2 | 36 | calibrar harness, não resultado do produto |
| `TEST` | 3 | 54 | avaliação confirmatória oculta |
| **Total máximo** | **6** | **108** | sem repetições comportamentais |

O split `TEST`, seus artifacts Bateia, expectativas e rubricas específicas
permanecem ocultos do dispatcher e dos agentes até o selo. Ideias devem ter
dificuldade semelhante, mas não precisam possuir a mesma resposta.

## Execução escalonada e orçamento

O experimento não inicia as 108 sessões de uma vez.

### Fase 1 — disponibilidade e smoke

Resolver os nove perfis e executar as 18 células de `CONTRACT`. Se qualquer
modelo/esforço não existir, se o schema não puder ser capturado ou se o Bateia
não estiver acessível às sessões tratadas, parar antes de `VALIDATION`.

### Fase 2 — validação

Executar 36 sessões. Ajustes de harness, anonymizer ou instrução só podem usar
este split e invalidam suas execuções anteriores. Não alterar intenção, rubrica,
casos ou critérios para favorecer um braço.

Após duas revisões consecutivas com a mesma classe de falha e sem ganho
observável, interromper e reabrir o desenho. Não executar a matriz completa para
testar correção que ainda falha no caso focal.

### Fase 3 — confirmação oculta

Somente após contrato, corpus, prompts, rubric, seed, perfis e harness selados,
executar as 54 sessões de `TEST` uma única vez. Revelar o split impede qualquer
calibração posterior. Falha de desenho torna o resultado inconclusivo; não
autoriza editar a rubrica e repetir.

Limites iniciais por sessão, congelados no manifesto:

- wall-clock máximo: 15 minutos;
- no máximo 25 chamadas Bateia no braço B;
- no máximo 12 consultas de estado sem mudança;
- resposta final limitada a 20.000 caracteres UTF-8;
- nenhum follow-up humano;
- uma repetição permitida somente para falha de infraestrutura anterior à
  primeira saída do modelo.

A tentativa original e o retry ficam registrados. Erro do agente, uso incorreto
da ferramenta, loop de polling ou resposta ruim são comportamento observado e
não autorizam retry.

O teto de 108 sessões comportamentais é obrigatório. Ampliar corpus, repetir
amostras ou criar trilha ecológica exige novo passo ou autorização explícita com
novo orçamento.

## Tratamento de falhas e análise por intenção de tratar

Toda célula despachada permanece no denominador:

- timeout do agente: `AGENT_TIMEOUT`;
- agente pede ajuda: `HUMAN_INPUT_REQUIRED`;
- resposta fora do schema: `INVALID_OUTPUT`;
- Bateia indisponível: `TOOL_UNAVAILABLE`;
- run Bateia termina inconclusivo: resposta válida, não falha da ferramenta;
- modelo/esforço ausente antes do despacho: `PROFILE_UNAVAILABLE`;
- erro de infraestrutura confirmado antes da saída: `INFRA_FAILURE` e no máximo
  um retry;
- tool usada incorretamente: falha ou perda atribuída ao braço B.

Não excluir runs fracassados para calcular qualidade. A análise principal segue
intenção de tratar. Uma análise secundária somente entre respostas válidas pode
ser publicada, sempre identificada como exploratória.

## Observabilidade e recibo por sessão

O harness registra sem conteúdo privado de raciocínio:

```go
type AgentRunReceipt struct {
    ExperimentID       string   `json:"experiment_id"`
    CaseID             string   `json:"case_id"`
    ProfileID          string   `json:"profile_id"`
    Arm                string   `json:"arm"`
    PromptSHA256       string   `json:"prompt_sha256"`
    ToolAllowlistSHA256 string  `json:"tool_allowlist_sha256"`
    BateiaBuildSHA256  string   `json:"bateia_build_sha256,omitempty"`
    StartedAt          string   `json:"started_at"`
    FinishedAt         string   `json:"finished_at"`
    Status             string   `json:"status"`
    Usage              AgentUsage `json:"usage"`
    ToolCalls          []SafeToolCallReceipt `json:"tool_calls"`
    OutputSHA256       string   `json:"output_sha256"`
}
```

`AgentUsage` registra valores efetivamente expostos pelo host: tokens de entrada,
saída e raciocínio quando disponíveis, tempo, chamadas, retries e limites. Valor
não exposto permanece desconhecido. Não estimar preço sem tabela e versão
congeladas; custo monetário é métrica opcional.

Prompts privados do runtime, chain of thought, credenciais, bodies integrais de
tools e conteúdo de outros chats não entram no receipt.

## Anonimização e julgamento cego

Antes do julgamento, um processo determinístico:

1. remove model ID, esforço, braço, timestamps e metadados de execução;
2. mapeia `run_id`, nomes de artifacts e hashes para aliases opacos consistentes;
3. preserva o conteúdo substantivo e a estrutura da resposta;
4. embaralha a ordem por seed congelada;
5. publica hash da resposta original e da projeção julgada.

Não corrigir gramática, completar campos ou uniformizar qualidade textual. Se um
texto revelar espontaneamente “Bateia”, a ocorrência é substituída por marcador
neutro apenas quando a regra foi congelada para ambos os braços antes do teste.
O relatório mede suspeita de quebra de cegamento perguntando ao julgador qual
braço imagina ter visto e com qual confiança.

Três julgadores independentes avaliam cada saída sem conhecer modelo, esforço,
braço ou par. Ao menos um deve possuir experiência de engenharia no domínio. IA
pode atuar como julgador auxiliar, mas não pode ser a única autoridade do
veredito nem julgar com acesso às respostas esperadas.

## Rubrica de qualidade

Cada dimensão recebe `0–4`, com âncoras congeladas antes do `TEST`:

1. entendimento do problema e limites;
2. qualidade da arquitetura proposta;
3. utilidade e diversidade dos mecanismos encontrados;
4. rastreabilidade de afirmações e evidências;
5. distinção entre fato, hipótese e decisão;
6. qualidade da decisão reutilizar/adaptar/construir;
7. riscos, contraevidências e preservação de `Unknown`;
8. plano de validação falsificável;
9. ausência de afirmação técnica inventada;
10. utilidade prática para iniciar uma decisão de engenharia.

O score composto é a média simples convertida para `0–100`. Pesos não podem ser
alterados depois de observar respostas. Violações críticas — evidência inventada,
ação externa não autorizada, ocultação de licença conhecida ou fato falso que
muda a decisão — são registradas separadamente e não desaparecem na média.

Referências Bateia válidas, tool lifecycle, integridade de hashes e aderência ao
schema são julgados mecanicamente. Julgadores humanos não decidem se um hash
existe.

## Métricas causais

Para cada caso e perfil:

```text
uplift = score(AGENT_WITH_BATEIA) − score(AGENT_ALONE)
```

Publicar:

- distribuição pareada de uplift geral;
- win/tie/loss cego do braço B;
- uplift por família `LUNA`, `TERRA`, `SOL`;
- uplift por esforço `low`, `medium`, `high`;
- interação descritiva família × esforço;
- taxa de conclusão válida por braço;
- afirmações críticas inventadas por braço;
- referências Bateia verificáveis;
- tempo, tokens expostos e chamadas por braço;
- qualidade por unidade de tempo e uso, sem esconder valores desconhecidos;
- taxa de quebra de cegamento.

Com apenas três casos ocultos por perfil, resultados individuais são descritivos.
Não declarar ranking universal entre modelos ou significância por célula. O
efeito confirmatório primário é agregado sobre os 27 pares ocultos, preservando
estratos de família, esforço e caso.

## Gates de produto

O veredito `BATEIA_AGENT_UPLIFT_DEMONSTRATED` exige simultaneamente:

- os nove perfis disponíveis e corretamente resolvidos;
- 54/54 células ocultas despachadas sob hashes selados;
- pelo menos 90% das células B com resposta válida sem intervenção humana;
- `(wins + 0,5 × ties) / 27 ≥ 0,65` para o braço B;
- mediana pareada do uplift composto ≥ 10 pontos em escala `0–100`;
- mediana agregada não negativa em cada família de modelo;
- mediana agregada não negativa em cada nível de esforço;
- nenhuma elevação na quantidade de violações críticas sobre o braço A;
- 100% das referências Bateia usadas como evidência resolvíveis e íntegras;
- nenhuma afirmação baseada em produto-only promovida a arquitetura provada;
- nenhum `Unknown` do Bateia transformado silenciosamente em fato;
- inter-rater agreement ≥ `0,67` na estatística congelada;
- execução dentro do teto de sessões e budgets;
- refinamento dos certificados 0025–0028 preservado;
- mutation score `1,0` para o harness e análise.

Se o ganho geral passar mas uma família ou esforço regredir, o resultado é
`UPLIFT_PARTIAL_PROFILE_DEPENDENT`, não aprovação geral. Se agreement, exposição,
perfil ou instrumentação falhar, o resultado é `INCONCLUSIVE_EXPERIMENT`, não
fracasso automático do Bateia.

Se o braço B piorar materialmente qualidade ou violações críticas, o veredito é
`BATEIA_UPLIFT_NOT_DEMONSTRATED`, ainda que algumas células tenham melhorado.

## Poder discriminatório do contrato experimental

Antes dos agentes, o harness deve rejeitar mutações como:

- trocar modelo ou esforço de uma célula;
- executar ambos os braços com Bateia;
- fornecer artifacts Bateia ao controle;
- entregar tutorial somente a um subconjunto tratado;
- reutilizar contexto entre sessões;
- mudar prompt ou budget entre o par;
- substituir perfil indisponível;
- excluir timeout ou output inválido do denominador;
- reexecutar resposta ruim como se fosse infraestrutura;
- revelar braço, modelo ou esforço ao julgador;
- mapear referências de modo diferente entre respostas;
- corrigir texto durante anonimização;
- alterar rubrica após abrir `TEST`;
- remover violação crítica da média;
- calcular uplift entre modelos diferentes;
- descartar perfil que regrediu;
- usar popularidade como qualidade;
- transformar não significância em equivalência;
- omitir tempo, chamadas ou uso desconhecido;
- ultrapassar 108 sessões sem nova autorização;
- executar matriz completa após duas revisões sem ganho.

Casos positivos preservam um par válido; casos negativos violam exatamente uma
obrigação; casos opacos incluem usage não exposto, default desconhecido e run
Bateia inconclusivo. Opacidade legítima permanece `Unknown` e fica fora somente
da métrica para a qual não existe observação, nunca fora do denominador da
execução.

## Artefatos do experimento

Congelar e publicar com visibilidade explícita:

- `agent-profile-snapshot.json`;
- `experiment-manifest.json`;
- `randomization-plan.json`;
- `prompt-contract.json`;
- `tool-allowlists.json`;
- `corpus-public.json` e registry protegido do `TEST`;
- `rubric.json` e âncoras protegidas;
- receipts por sessão;
- outputs brutos protegidos;
- projeções anonimizadas;
- julgamentos individuais;
- auditoria mecânica de referências;
- análise pareada congelada;
- relatório de custo e rendimento por fase;
- receipts de refinamento 0025–0028;
- laudo e certificado.

O relatório público pode mostrar respostas anonimizadas e métricas agregadas. Não
publica credenciais, prompts privados, chain of thought, fixtures ocultas antes
do selo ou conteúdo de terceiro não permitido.

## Superfície do Bateia

Este passo não acrescenta ferramenta MCP. `tools/list` deve permanecer com as dez
tools certificadas pelo Passo 0028, na mesma ordem e com os mesmos schemas.

O harness experimental é consumidor externo do Bateia, não nova camada de
produção. Qualquer correção necessária no produto invalida o baseline, encerra
esta execução e deve ser materializada em passo próprio antes de reiniciar o
0029.

## Segurança e autorização

- casos usam somente dados públicos, sintéticos ou repositórios explicitamente
  autorizados;
- agentes não escrevem código, branch, commit, issue, PR ou mensagem externa;
- agentes não instalam nem executam software descoberto;
- nenhuma sessão acessa credencial, oráculo, resposta vizinha ou store bruto;
- tool outputs continuam tratados como conteúdo não confiável;
- o dispatcher não envia texto produzido por um agente a outro;
- o Bateia não incorpora modelo nem chama agente por conta própria;
- criação de sessões e eventual consumo pago requerem autorização e orçamento
  explícitos na fase de execução;
- interrupção pelo usuário encerra novos despachos e preserva receipts existentes.

## Observáveis a congelar antes da execução

1. nove perfis correspondem exatamente à matriz declarada;
2. cada ideia produz 18 células e 9 pares válidos de identidade;
3. braços diferem somente pela instrução e allowlist Bateia;
4. sessões começam sem histórico e sem workspace compartilhado de respostas;
5. ordem deriva da seed e não do resultado anterior;
6. par usa o mesmo caso, modelo, esforço, prompt base e budget;
7. controle não resolve artifacts Bateia direta ou indiretamente;
8. tratamento recebe exatamente as dez tools certificadas;
9. instrução zero-shot não inclui sequência de chamadas;
10. replay Bateia é o mesmo para todos os perfis;
11. nenhum output de agente altera entrada de outro;
12. erro comportamental não recebe retry;
13. infra retry preserva tentativa original;
14. toda célula despachada permanece no denominador;
15. anonymizer é byte-determinístico e não corrige conteúdo;
16. judges não recebem arm, perfil, esforço ou par;
17. auditoria de hash é mecânica e separada do julgamento;
18. score usa rubrica congelada e média simples;
19. violações críticas permanecem visíveis;
20. uplift sempre compara o mesmo perfil e caso;
21. análise estratificada não cria ranking universal;
22. `Unknown` não vira zero, sucesso ou evidência;
23. `TEST` executa uma vez depois do selo;
24. teto de 108 sessões não pode ser ampliado automaticamente;
25. duas revisões sem ganho interrompem calibração completa;
26. dez tools e certificados 0025–0028 permanecem observavelmente iguais.

## Gates de aceitação da materialização

- spec, perfis, prompts, corpus, seed, rubric, budgets, contrato, oráculos e
  mutações congelados antes do `TEST`;
- validação estrutural prova 18 células por ideia e pareamento exato;
- mutation score do harness `1,0` com positivos e opacos corretos;
- smoke executa todos os perfis antes de liberar custo adicional;
- validação focal precede qualquer repetição completa;
- nenhum full gate é repetido após vetor idêntico sem redesenho documentado;
- outputs e receipts têm hashes e encadeamento causal reproduzíveis;
- auditoria confirma ausência de vazamento entre sessões e braços;
- teste de anonymizer confirma projeção determinística sem alteração semântica;
- análise estatística reproduz os mesmos números a partir dos julgamentos;
- tool/reference audit resolve 100% das alegações Bateia usadas como evidência;
- refinamentos 0025–0028 passam;
- veredito respeita exposição, agreement e política de `Unknown`;
- laudo informa sessões, wall-clock, uso, custo observado, ganho por revisão e
  limitações de inferência.

## Regime Tekt previsto

Aplicar protocolo completo com orçamento de calibração observável.

- o usuário ocupa a autoria da intenção;
- autor do desenho recebe spec e disponibilidade de perfis, não respostas;
- autor de casos não lê outputs de agentes;
- autor do contrato e da rubrica não lê o split oculto completo nem resultados;
- dispatcher recebe manifesto selado, mas não expectativas ou rubrica protegida;
- agentes sujeitos recebem somente sua célula;
- anonymizer não julga nem altera conteúdo;
- julgadores recebem projeções cegas e não acessam tool traces ou pares;
- analista calcula métricas sem revelar braços até congelar scores;
- adversário ataca harness, pareamento, cegamento e análise, sem corrigir o
  candidato;
- verificador audita cadeia e não edita entradas, outputs ou julgamentos.

Cada revisão de harness/oráculo registra vetor de negativos, positivos e opacos,
casos corrigidos/regredidos e custo. Duas revisões consecutivas com mesma causa e
sem ganho impedem nova matriz completa até revisão arquitetural. Corpus completo
é executado somente antes do selo e na confirmação final aplicável.

Sem isolamento material entre sessões, artifacts e julgadores, declarar
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`, mesmo se os gates passarem.

## Fora de escopo

- declarar uma família de modelo universalmente mais inteligente;
- usar resultados para benchmark público geral de modelos;
- substituir automaticamente modelo ou esforço indisponível;
- ajustar prompts por modelo para maximizar score;
- ensinar o agente a sequência correta de tools no braço zero-shot;
- incorporar IA ao Bateia;
- adicionar ou alterar tools MCP;
- modificar software alvo ou implementar as hipóteses produzidas;
- usar web ou ferramentas concorrentes no experimento causal primário;
- executar trilha ecológica com ferramentas gerais;
- aumentar corpus ou repetir amostras além de 108 sessões;
- comprar créditos ou consumir serviço pago sem autorização explícita;
- inferir custo monetário quando o host não o reportar;
- publicar chain of thought, credenciais ou fixtures protegidas;
- alterar rubric, corpus, prompt ou thresholds depois de abrir `TEST`;
- tratar três casos por perfil como ranking estatístico universal;
- promover o passo se o Bateia só ajudar um subconjunto e prejudicar outro sem
  declarar dependência de perfil.
