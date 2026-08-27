# Passo 0006 — Endurecer e calibrar evidências

**Estado**: `IMPLEMENTADO`
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)
**Depende de**: [`Passo 0005 — Triar candidatos e comprovar capacidades`](0005-triar-candidatos-e-comprovar-capacidades.md)
**Onde roda**: Bateia
**Camadas previstas**: L0, L1, L2, L3 e L4
**Natureza**: endurecimento da coleta, redução de falsos positivos e calibração observada

## Intenção

Tornar `CandidateAssessment` confiável o suficiente para alimentar uma decisão de
reuso. O Passo 0005 criou os contratos e o primeiro pipeline; este passo fecha as
lacunas encontradas na auditoria antes de construir o motor de decisão.

O resultado esperado não é uma nova heurística de recomendação. É uma base de
evidências na qual:

- coleta incompleta nunca se apresenta como completa;
- `ABSENT` significa ausência dentro de um escopo realmente executado;
- símbolos e componentes apontam para implementação concreta;
- identidade lógica do repositório e revisão imutável não se confundem;
- métricas são calculadas sobre candidatos reais e julgamentos versionados.

## Restrição central: continuar sem LLM

Permanecem proibidos LLM, embeddings, rerankers neurais, classificação por prompt
e fallback para serviços de IA. Casos difíceis permanecem `INCONCLUSIVE` até que
haja evidência determinística suficiente.

## Incertezas reduzidas

- se todos os alvos necessários foram realmente tentados;
- se uma falha individual foi escondida por um arquivo coletado com sucesso;
- se a árvore remota cobre código abaixo da raiz;
- se um símbolo relacionado foi implementado ou apenas mencionado em comentário;
- se um componente emitido é reutilizável ou apenas possui nome interessante;
- qual é a precisão observada de `PROVEN` e `ABSENT` em candidatos reais;
- se os vetores de operação, licença e manutenção sustentam o próximo passo.

## Entrega

### 1. Identidade estável e revisão imutável

Separar explicitamente:

```go
type RepositoryIdentity struct {
    Owner string `json:"owner"`
    Repo  string `json:"repo"`
}

type ResolvedRepository struct {
    Identity RepositoryIdentity `json:"identity"`
    Commit   string             `json:"commit"`
}
```

Regras:

- joins com descoberta e peneira usam `owner/repo`;
- fontes e bundles usam `owner/repo@commit`;
- `HEAD` pode aparecer somente no plano como pendência de resolução;
- coleta resolve o SHA antes de ler árvore ou conteúdo;
- assessment preserva a `ScreeningDecision` original depois da resolução.

### 2. Resultado individual de coleta

Cada alvo deve registrar seu término:

```text
PLANNED → COLLECTED | NOT_FOUND | FAILED | SKIPPED_BUDGET | UNSUPPORTED
```

O bundle deve declarar:

- alvos planejados, tentados e coletados;
- bytes planejados e lidos;
- profundidade alcançada;
- falhas por path;
- orçamento esgotado;
- cobertura por capacidade.

`Complete` deixa de ser um booleano otimista e passa a derivar desses resultados.
Um arquivo bem-sucedido não pode esconder falhas nos demais.

### 3. Regra segura para `ABSENT`

`ABSENT` exige simultaneamente:

```text
regra conhecida
+ alvos relevantes planejados
+ árvore concluída no escopo
+ todos os alvos relevantes terminados
+ zero falha ou parser ausente
+ zero sinal suficiente
```

Qualquer violação produz `INCONCLUSIVE`. Testes de regressão devem cobrir falha de
árvore, arquivo individual ausente, timeout, limite de bytes e parser incompatível.

### 4. Árvore remota recursiva e limitada

O adaptador GitHub deve obter árvore no SHA resolvido e respeitar:

- profundidade máxima;
- quantidade máxima de nós;
- bytes e arquivos do plano;
- exclusão explícita de `vendor`, fixtures, traduções e dependências baixadas;
- falha observável, sem fabricar uma árvore substituta.

O adaptador local deve produzir o mesmo contrato para permitir replay comparável.

### 5. Planejamento orientado pela `IdeaSpec`

Alterar a boca de inspeção para:

```bash
bateia idea evidence plan \
  --idea ideia.json \
  --screening screening.json \
  --rules capability-rules.json \
  --output extraction-plan.json
```

Cada `ExtractionTarget` deve declarar quais capacidades cobre e qual regra o
escolheu. O plano deve priorizar paths e símbolos relacionados, mantendo uma
amostra geral pequena para evitar viés de confirmação.

### 6. Evidência estrutural conservadora

Substituir a aproximação “arquivo contém termo e alguma declaração” por extração
determinística de declarações suportadas.

Uma fonte de símbolo deve preencher:

```text
path · symbol · line · snippet · content_hash · rule_id · repo@commit
```

O termo precisa estar ligado ao nome, assinatura, corpo ou teste do símbolo, não a
um comentário distante. Linguagem sem extrator adequado gera `INCONCLUSIVE` para
prova estrutural, preservando evidências documentais auxiliares.

Arquivos de teste não serão descritos como testes executados. Sem executar código
remoto, sua classe será “implementação de teste encontrada”.

### 7. Componentes reutilizáveis concretos

Emitir `ReusableComponent` somente quando houver:

- arquivo de implementação;
- símbolo ou módulo identificável;
- capacidade comprovada;
- fonte forte escolhida por força, não por ordem alfabética;
- caveats de stack, operação e manutenção.

README, licença e manifesto isolados nunca viram componente reutilizável.

### 8. Vetores necessários para a decisão

Completar, sem colapsar num score único:

- **operação**: local comprovado, serviço obrigatório, upload de código, conta,
  token, banco, infraestrutura e runtimes;
- **licença**: detecção exata por texto/SPDX, ambiguidade e política configurável;
- **manutenção**: arquivamento, último commit, releases/tags, CI/testes declarados e
  aviso explícito de encerramento;
- **stack**: linguagens, manifestos e custo técnico observável de integração;
- **restrições e non-goals**: conflitos eliminatórios ligados à `IdeaSpec`.

Estrelas não são evidência suficiente de manutenção ou maturidade.

### 9. Coleta resiliente por candidato

Falha ao resolver ou coletar um candidato deve gerar bundle inconclusivo para ele,
sem abortar os demais. O relatório final preserva todas as falhas e permite replay
dos candidatos bem-sucedidos.

### 10. Corpus real e métricas

Congelar bundles e julgamentos para, no mínimo:

- Enola;
- Supermodel CLI;
- mllint;
- um componente parcial útil;
- um candidato encontrado somente por ruído vendorizado;
- cinco controles `ABSENT`;
- todos os casos com falha de coleta.

Calcular:

```text
evidence_precision
false_absence_rate
inconclusive_rate
source_trace_rate
budget_completion
component_precision
```

Regras são congeladas antes da revisão. Ajustes posteriores criam nova versão e
reexecutam todo o corpus, sem sobrescrever o resultado anterior.

## Critérios de aceitação

- `evidence_precision` ≥ 90%, com ao menos dez capacidades `PROVEN`;
- `false_absence_rate` = 0 na amostra revisada;
- `source_trace_rate` = 100%;
- `component_precision` = 100% na primeira amostra;
- toda coleta parcial produz `INCONCLUSIVE`, nunca `ABSENT`;
- toda fonte estrutural possui path, símbolo, linha, hash e commit;
- todos os candidatos geram resultado mesmo quando um deles falha;
- Enola, Supermodel, mllint e ruído vendorizado mordem controles distintos;
- métricas e julgamentos ficam versionados;
- `go test -race ./...`, `go vet ./...` e `crystalline-lint .` passam.

## Testes obrigatórios

1. falha ao listar árvore não produz árvore fictícia nem `ABSENT`;
2. um arquivo coletado não esconde falha em outro alvo;
3. árvore remota e local respeitam profundidade e orçamento equivalentes;
4. resolução de SHA preserva a decisão da peneira;
5. plano liga targets às capacidades da ideia;
6. comentário distante de declaração não produz evidência de símbolo;
7. símbolo real preserva nome e linha;
8. arquivo de teste encontrado não é descrito como teste executado;
9. fonte auxiliar alfabeticamente anterior não vira componente;
10. licença usa correspondência exata, sem substring acidental;
11. operação local pode resultar em `PROVEN` com evidência adequada;
12. falha de um candidato não aborta a coleta dos demais;
13. métricas são reproduzíveis a partir de bundles congelados;
14. avaliação offline não acessa rede nem lê variáveis de ambiente.

## Fora de escopo

- decisão `FORK`, `EXTEND`, `COMPOSE`, `INTEGRATE`, `REFERENCE` ou `BUILD`;
- executar ou instalar código não confiável;
- aconselhamento jurídico automático;
- cobertura estrutural de todas as linguagens;
- interface gráfica;
- LLM, embeddings ou inferência generativa.

## Próximo passo condicionado

O **Passo 0007** diagnosticará por que positivos reais permanecem `PARTIAL` ou
`INCONCLUSIVE`, ampliará a cobertura de forma controlada e demonstrará o gate antes
que o motor de decisão seja iniciado.

O **Passo 0008** construirá então a decisão `FORK · EXTEND · COMPOSE ·
INTEGRATE · REFERENCE · BUILD`, consumindo `CandidateAssessment` calibrado,
nunca descrições ou scores de descoberta como substitutos de evidência.
