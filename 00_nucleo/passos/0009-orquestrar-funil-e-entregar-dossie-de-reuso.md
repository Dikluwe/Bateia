# Passo 0009 — Orquestrar o funil e entregar um dossiê de reúso

**Estado**: `EM VALIDAÇÃO`  
**Desejo de origem**: [`D-001 — Da ideia à decisão de reuso`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md)  
**Depende de**: [`Passo 0008 — Comparar candidatos e decidir reuso`](0008-comparar-candidatos-e-decidir-reuso.md)  
**Onde roda**: Bateia  
**Camadas previstas**: L0, L1, L2, L3 e L4  
**Natureza**: orquestração ponta a ponta, checkpoints e entrega de produto

## Incerteza que este passo reduz

O Bateia já materializa separadamente:

```text
IdeaSpec
  → ExpansionPlan
  → DiscoverySnapshot
  → ScreeningReport
  → EvidenceExtractionPlan / EvidenceBundle
  → AssessmentReport
  → ReuseDecisionReport
```

Ainda não sabemos se uma pessoa consegue atravessar esse funil como produto sem
conhecer os comandos internos, transportar caminhos manualmente ou descobrir qual
artefato alimentar no estágio seguinte.

Este passo deve provar que uma única execução coordenada consegue partir de uma
`IdeaSpec` e entregar um dossiê de decisão reproduzível, preservando todos os
artefatos intermediários para auditoria.

## Intenção

Criar um caso de uso de orquestração que invoque os contratos existentes sem
reimplementar suas regras. A execução deve:

1. validar a ideia, a política e a configuração antes de consultar provedores;
2. expandir a ideia;
3. descobrir e deduplicar candidatos;
4. aplicar a peneira calibrada;
5. planejar e coletar evidências somente dos candidatos aprovados ou auditados;
6. produzir assessments;
7. decidir o ponto de partida;
8. gravar um dossiê navegável e um manifesto de execução.

O produto deste passo não é “mais um score”. É a redução do funil inteiro a uma
boca de inspeção segura, retomável e explicável.

## Restrição central: continuar sem LLM

Nenhum estágio pode chamar LLM, embeddings ou geração remota. Expansão lexical,
descoberta por provedores observáveis, regras de screening, extração de conteúdo,
assessment e decisão continuam usando os contratos versionados já implementados.

O orquestrador não pode inventar dados ausentes para fazer a execução avançar.
Falha parcial permanece registrada; ausência de condição para decidir termina em
`INCONCLUSIVE`.

## Modos de execução

### `live`

Executa descoberta e coleta contra um provedor explicitamente configurado. Rede,
credenciais e limites pertencem aos adaptadores L3. O núcleo recebe interfaces,
não variáveis de ambiente ou clientes HTTP.

### `replay`

Reexecuta screening, assessment e decisão a partir de snapshots e bundles
congelados. Não acessa rede e deve produzir a mesma decisão para os mesmos bytes
de entrada, versões de regra e política.

### `local`

Analisa uma ou mais worktrees informadas explicitamente, inclusive repositórios
sem remoto e alterações ainda não publicadas. O caminho local nunca vira
identidade pública fictícia; o manifesto registra identidade local e hash do
conteúdo observado.

## Contrato da execução

```go
type IdeaRunMode string

const (
    IdeaRunLive   IdeaRunMode = "live"
    IdeaRunReplay IdeaRunMode = "replay"
    IdeaRunLocal  IdeaRunMode = "local"
)

type IdeaRunRequest struct {
    Idea            IdeaSpec            `json:"idea"`
    Policy          ReuseDecisionPolicy `json:"policy"`
    Mode            IdeaRunMode         `json:"mode"`
    Strategy        string              `json:"strategy"`
    Provider        string              `json:"provider,omitempty"`
    Language        string              `json:"language,omitempty"`
    LocalPaths      []string            `json:"local_paths,omitempty"`
    ReplayDirectory string              `json:"replay_directory,omitempty"`
    OutputDirectory string              `json:"output_directory"`
    Resume          bool                `json:"resume"`
}
```

O contrato não transporta tokens. Credenciais são resolvidas somente no shell ou
adaptador correspondente e nunca são serializadas.

## Manifesto e checkpoints

Cada execução produz `run.json` com:

```go
type IdeaRunManifest struct {
    SchemaVersion string             `json:"schema_version"`
    RunID         string             `json:"run_id"`
    Mode          IdeaRunMode        `json:"mode"`
    Status        IdeaRunStatus      `json:"status"`
    InputHashes   map[string]string  `json:"input_hashes"`
    RuleVersions  map[string]string  `json:"rule_versions"`
    Stages        []IdeaRunStage     `json:"stages"`
    DecisionPath  string             `json:"decision_path,omitempty"`
    Failures      []ExecutionFailure `json:"failures,omitempty"`
}
```

Estados do run:

```text
RUNNING · COMPLETE · PARTIAL · INCONCLUSIVE · FAILED
```

Estados de estágio:

```text
PENDING · RUNNING · COMPLETE · SKIPPED · FAILED
```

Cada estágio registra:

- nome e versão do contrato;
- hash das entradas efetivamente consumidas;
- caminhos e hashes das saídas;
- início e fim;
- contagens observáveis;
- falhas e motivo de `SKIPPED`.

Um estágio só pode ser reutilizado por `--resume` quando contrato, entradas,
regras, política e hashes das saídas ainda coincidirem. Caso contrário, ele e
todos os dependentes são recomputados. Nunca reutilizar um arquivo apenas porque
o nome existe.

## Diretório do dossiê

```text
run.json
idea.json
expansion-plan.json
discovery-snapshot.json
screening-report.json
evidence/
  plans/
  bundles/
assessment-report.json
reuse-policy.json
reuse-decision.json
decision.md
```

`decision.md` é uma projeção humana do JSON, não uma fonte de verdade paralela.
Deve apontar para candidato, componentes, gaps, bloqueios, regras e evidências
por caminhos relativos dentro do dossiê.

## Boca de inspeção

Execução por provedor:

```bash
bateia idea recommend \
  --idea ideia.json \
  --policy reuse-policy.json \
  --mode live \
  --provider github \
  --strategy lexical \
  --output .bateia/runs/minha-ideia
```

Repositório local:

```bash
bateia idea recommend \
  --idea ideia.json \
  --policy reuse-policy.json \
  --mode local \
  --local /caminho/projeto-a \
  --local /caminho/projeto-b \
  --output .bateia/runs/minha-ideia
```

Replay offline:

```bash
bateia idea recommend \
  --idea ideia.json \
  --policy reuse-policy.json \
  --mode replay \
  --replay .bateia/runs/minha-ideia \
  --output .bateia/replays/minha-ideia
```

Retomada:

```bash
bateia idea recommend ... --resume
```

## Regras de segurança e honestidade

- validar `--output` antes da primeira consulta externa;
- não sobrescrever dossiê incompatível sem erro explícito;
- gravar artefatos de forma atômica;
- nunca registrar tokens, cabeçalhos de autenticação ou conteúdo secreto de
  variáveis de ambiente;
- respeitar budget, rate limit e cancelamento em todos os estágios;
- não executar código encontrado nos candidatos;
- não publicar, clonar, forkar ou modificar repositórios candidatos;
- preservar falhas parciais por candidato;
- não promover run `PARTIAL` a `COMPLETE` apenas porque existe uma decisão;
- não emitir `BUILD` quando a cadeia necessária estiver incompleta.

## Testes obrigatórios

1. uma execução completa produz todos os artefatos e termina `COMPLETE`;
2. decisão orquestrada é idêntica à execução manual dos mesmos estágios;
3. replay não chama rede nem provedor;
4. duas execuções replay sobre os mesmos bytes produzem JSON semanticamente
   idêntico, exceto campos temporais explicitamente excluídos;
5. `--resume` reutiliza somente estágios cujos hashes permanecem válidos;
6. alterar a ideia invalida todos os estágios dependentes;
7. alterar somente a política reutiliza evidências e refaz a decisão;
8. alterar regra de assessment invalida assessment e decisão, não descoberta;
9. bundle ausente ou corrompido impede decisão conclusiva;
10. falha parcial aparece no manifesto, no dossiê e na decisão quando relevante;
11. saída inválida falha antes de chamar o provedor;
12. token e cabeçalhos secretos não aparecem em nenhum artefato;
13. modo local funciona sem remoto Git e inclui alterações não publicadas;
14. caminhos locais diferentes com conteúdo igual não recebem identidade pública
    inventada;
15. interrupção durante escrita não deixa JSON final truncado;
16. ordem de resultados do provedor não altera o dossiê semântico;
17. orçamento ou cancelamento interrompe novos trabalhos e preserva checkpoints;
18. `decision.md` referencia evidências existentes no dossiê;
19. o comando não executa código de terceiros;
20. o linter V15 permanece sem violações e cada novo consumidor possui prompt
    próprio.

## Critério de produto

Executar o comando ponta a ponta sobre pelo menos:

- duas ideias com descoberta `live`;
- duas ideias em modo `local`, incluindo um repositório sem remoto;
- as cinco ideias congeladas no replay do Passo 0008.

Para aprovação:

- `100%` dos runs devem possuir cadeia de hashes verificável;
- `100%` das decisões devem apontar para o dossiê que as sustenta;
- `0` segredo serializado;
- `0` divergência entre replay e execução manual equivalente;
- ao menos `80%` dos dossiês devem permitir que um revisor escolha o ponto de
  partida sem abrir previamente todos os repositórios candidatos;
- tempo, bytes, candidatos e chamadas externas devem ser reportados por estágio,
  sem impor neste passo um limite de desempenho ainda não medido.

## Critério de parada

Interromper a promoção a produto e voltar ao estágio responsável se:

- o orquestrador precisar duplicar regra de domínio existente;
- a retomada aceitar artefato com hash ou versão divergente;
- replay tocar a rede;
- uma falha desaparecer entre estágio e manifesto;
- a renderização humana acrescentar afirmação ausente no JSON;
- um caso local depender de GitHub ou de remoto configurado.

## Fora de escopo

- interface gráfica ou TUI;
- execução automática da recomendação (`fork`, clone, patch ou integração);
- escolha interativa entre alternativas durante o run;
- aconselhamento jurídico sobre licenças;
- execução de testes ou binários de terceiros;
- LLM, embeddings ou geração probabilística;
- daemon, fila distribuída ou agendamento recorrente;
- otimização de performance antes da primeira medição ponta a ponta.

## Próximo passo previsto

Depois que a orquestração provar uma entrega utilizável, o passo seguinte poderá
avaliar a qualidade do dossiê com usuários e transformar divergências observadas
em calibração versionada do funil completo. Ele não deve começar antes de termos
runs reais, custos medidos e pontos de abandono registrados.

## Evidência de execução

Implementação materializada em 2026-08-26 com contrato L1, orquestrador, CLI,
fiação live/local, escrita atômica, manifesto com cadeia de hashes, budgets,
timeout, retomada e projeção Markdown. Foram executados dois runs live, dois
locais — um deles sem remoto — e cinco replays.

Todos os nove dossiês terminaram honestamente em `INCONCLUSIVE`. Isso comprova a
execução e a preservação de falhas, mas não satisfaz o gate de utilidade de 80%.
Por isso o passo permanece `EM VALIDAÇÃO`: a implementação está presente e os
gates técnicos passaram, porém a promoção a produto aguarda um corpus de
evidências ponta a ponta capaz de sustentar decisões conclusivas úteis.

Ver: [`Laudo 0009`](../laudos/0009-orquestrar-funil-e-entregar-dossie-de-reuso.md).
