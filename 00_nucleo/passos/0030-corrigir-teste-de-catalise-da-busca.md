# Passo 0030 — Corrigir o teste de catálise da busca

**Estado**: `PLANEJADO`
**Desejos de origem**: [`D-001`](../desejos/D-001-da-ideia-a-decisao-de-reuso.md) e [`D-002`](../desejos/D-002-alternativas-de-software-com-evidencia.md)
**Corrige**: a trilha ecológica do [`Passo 0029 — Medir o ganho da IA usando o Bateia`](0029-medir-ganho-da-ia-usando-bateia.md)
**Evidência causal**: tentativa R3 `PAIR_INVALID_NO_COMPARISON` e veredito independente `APROVADO_FAIL_CLOSED`
**Baseline obrigatório**: artefatos R3 imutáveis, Bateia MCP certificado e contagem global `16/108`
**Onde roda**: harness experimental segregado, sessões externas limpas, pesquisa nativa e Bateia MCP
**Camadas previstas**: L0, L1, L2, L3 e L4; L4 permanece somente como composition root e nenhuma ferramenta MCP de produção é alterada
**Natureza**: correção metodológica e de observabilidade antes de novo teste comportamental

## Problema

O teste ecológico que deveria comparar:

```text
A — IA + pesquisa nativa na internet
B — IA + Bateia + a mesma pesquisa nativa na internet
```

não chegou a produzir um par comparável. O tratamento executou uma chamada
`bateia_recommend` e depois três pares físicos de eventos `web_search`, mas o
controle não foi despachado porque duas obrigações operacionais falharam:

1. a saída tinha `22.578` caracteres Unicode, `4.578` acima do teto editorial
   de `18.000`, embora fosse JSON válido, aderente ao schema e idêntica à última
   mensagem do agente;
2. cada item web continha duas ocorrências da chave `item.id`. Um parser comum
   descartava silenciosamente a primeira e fazia parecer que havia três IDs
   normalizados inequívocos, apesar de nenhuma política para chaves duplicadas
   ter sido congelada.

O contrato fail-closed encerrou corretamente aquela tentativa. Porém, ao fazer
o despacho do segundo braço depender da validade comportamental do primeiro,
o desenho impediu justamente a comparação que pretendia medir.

`INCONCLUSIVE` não causou a invalidação. A recomendação do Bateia foi
inconclusiva, a IA rejeitou candidatos desalinhados e usou a web para procurar
alternativas e reduzir lacunas. Uma abstenção calibrada pode orientar busca e
deve permanecer resultado neutro, nunca derrota automática do Bateia.

## Evidência congelada da falha

Esta correção parte dos bytes já auditados e não os reinterpreta nem substitui:

| Artefato | SHA-256 | Observação protegida |
|---|---|---|
| `00_nucleo/benchmark/0029/catalysis-pilot-r3-results.json` | `3239ac23443e451c63027c93867aacfa55ec07e22c87709af97c8666ac2f8524` | par inválido, sem comparação |
| `00_nucleo/benchmark/0029/catalysis-pilot-r3-final-verification.json` | `7fd942124b0f3045e19b3c4e377b7e7c87371fd283aad4c90df06130888d82e2` | `APROVADO_FAIL_CLOSED` |
| `00_nucleo/benchmark/0029/sessions-catalysis-r3/treatment/events.jsonl` | `f1e61bf6d04de5c13426f8b09652d471f35f281d9b9c34e3dacd8e3383410a73` | uma recomendação seguida de três pares web físicos e IDs duplicados |
| `00_nucleo/benchmark/0029/sessions-catalysis-r3/treatment/output.json` | `42c36b20e94915efca624febd1eb58a03253f48d3954bf1c39b757def096b9a8` | JSON/schema válidos, `22.578` caracteres |

Qualquer diferença nesses hashes invalida o uso da tentativa R3 como fixture de
regressão. O novo passo cria sucessores versionados; não edita contrato, prompt,
recibo, output ou veredito R3.

## Incerteza que este passo reduz

Antes de medir se o Bateia melhora a pesquisa, precisamos saber se o harness:

- completa os dois braços mesmo quando o primeiro produz resposta ruim;
- separa elegibilidade técnica de qualidade, concisão e custo;
- lê eventos web sem perder chaves JSON duplicadas;
- distingue chamada física observada de identidade normalizada desconhecida;
- preserva `INCONCLUSIVE` como abstenção válida;
- mede se o Bateia muda, estreita ou corrige a pesquisa posterior;
- mantém iguais modelo, esforço, caso, pesquisa nativa e orçamento entre braços;
- continua dentro do teto global já autorizado pelo Passo 0029.

## Objetivo

Materializar e validar um harness corrigido e, somente após novo selo e nova
autorização de execução, realizar três comparações pareadas exploratórias:

```text
mesmo caso + mesmo modelo + mesmo esforço + mesmo budget web

A — IA + 1..3 pesquisas nativas
B — IA + exatamente 1 bateia_recommend
         + 1..3 pesquisas nativas para confirmar, rejeitar ou estreitar direções
```

O estimando primário é o ganho incremental de disponibilizar o Bateia quando a
IA já pode pesquisar a internet. O teste não compara Bateia contra internet e
não proíbe que o próprio Bateia use seus conectores externos certificados.

Três pares são calibração exploratória do método, não evidência confirmatória
geral do produto.

## Igualdade e diferença entre braços

| Dimensão | Controle `IA_INTERNET` | Tratamento `IA_INTERNET_BATEIA` |
|---|---|---|
| Caso, modelo e esforço | idênticos | idênticos |
| Pesquisa nativa | `1..3` chamadas | `1..3` chamadas |
| Limites de tempo e saída | idênticos | idênticos |
| Outras ferramentas | nenhuma | nenhuma |
| Bateia MCP | ausente | somente `bateia_recommend` |
| Ordem obrigatória | pesquisa e conclusão | Bateia concluído antes da primeira pesquisa |
| Contrato de resposta | idêntico | idêntico |

O prompt base é byte-idêntico. A única diferença textual permitida é a
instrução de tratamento que exige a chamada ao Bateia e o uso crítico de suas
direções. O agente pode confirmar, rejeitar ou não usar uma direção; não é
obrigado a concordar com a recomendação.

## Correção 1 — completar o par por intenção de tratar

Depois que o preflight global passa, a validade comportamental do primeiro
braço não controla o despacho do segundo:

- ambos os braços são despachados exatamente uma vez conforme a ordem selada;
- timeout, resposta inválida, excesso editorial, uso incorreto da ferramenta ou
  `INCONCLUSIVE` permanecem no denominador;
- falha do primeiro braço não cancela nem substitui o segundo;
- não há retry comportamental;
- execução continua com uma sessão global em voo e intervalo mínimo de `120`
  segundos entre inícios;
- somente falha global anterior à primeira chamada de modelo pode impedir os
  dois despachos.

Uma falha de infraestrutura que torne o segundo despacho impossível produz
`PAIR_PARTIAL_INFRA`, preserva tudo que já ocorreu e não autoriza comparação de
eficácia.

## Correção 2 — separar validade, concisão e segurança

O contrato define três superfícies diferentes:

### Validade estrutural

A saída é estruturalmente válida quando:

- contém um único objeto JSON;
- passa o schema congelado;
- é byte-idêntica à última mensagem registrada;
- respeita o limite rígido demonstrável do schema e do runtime.

### Meta editorial

`18.000` caracteres Unicode é meta de concisão e custo, não condição para
apagar o par. Excedê-la:

- registra `EDITORIAL_TARGET_EXCEEDED`;
- reduz a dimensão de utilidade/concisão congelada na rubrica;
- permanece observável em caracteres, bytes e tokens;
- não impede despacho, anonimização ou julgamento do outro braço.

### Limite rígido

O schema sucessor deve possuir `maxLength`, `maxItems` e
`additionalProperties: false` suficientes para um analisador determinístico
provar um teto serializado máximo antes das sessões. O manifesto congela esse
teto, no máximo `65.536` bytes, e a mesma prova vale para ambos os braços.

Se o runtime não conseguir impor ou demonstrar o teto, o preflight falha com
zero sessão. Não truncar, resumir ou reparar output depois da execução.

O tamanho integral da resposta do Bateia é comportamento e custo do tratamento:
deve ser medido, não ocultado por uma projeção criada apenas para favorecer o
teste. Uma futura mudança de produto para compactar essa resposta exige outro
baseline e refinamento próprio.

## Correção 3 — decodificação web que preserva duplicatas

O harness não pode desserializar a linha bruta diretamente em `map[string]any`.
Ele usa um scanner JSON tokenizado que preserva ordem, quantidade e valor de
todas as ocorrências de cada chave.

Para cada item `web_search`, o decoder classifica:

```go
type WebEventEncoding string

const (
    WebEventSingleExecID WebEventEncoding = "SINGLE_EXEC_ID"
    WebEventDualIDV0149  WebEventEncoding = "DUAL_ID_CODEX_0_149"
    WebEventOpaque       WebEventEncoding = "OPAQUE"
)
```

`DUAL_ID_CODEX_0_149` exige exatamente duas chaves `item.id`, na ordem:

1. identidade do item no formato `item_<n>`;
2. identidade física da chamada no formato `exec-<uuid>`.

A segunda ocorrência é `normalized_web_call_id` somente porque esta regra foi
congelada antes do novo candidato. Os eventos `started` e `completed` precisam
ter o mesmo par de IDs e formar uma máquina de estados válida. Uma ocorrência,
ordem, prefixo ou pareamento diferente resulta em `OPAQUE`, preserva os bytes e
produz `Unknown`; nunca aplica política implícita de primeiro ou último valor.

Opacidade de telemetria invalida a métrica de aderência de tools, mas não apaga
a resposta nem impede a conclusão do par por intenção de tratar.

## Correção 4 — `INCONCLUSIVE` neutro e útil

Existem duas abstenções diferentes:

- `bateia_recommend.status == INCONCLUSIVE`: resultado válido da ferramenta,
  que pode fornecer lacunas, candidatos rejeitados ou direções de verificação;
- decisão final da IA inconclusiva: resposta válida quando explica a evidência
  insuficiente e propõe validação falsificável.

Nenhuma recebe derrota, vitória ou score zero por default. Julgadores aplicam
as mesmas âncoras a uma decisão conclusiva ou inconclusiva. Abstenção calibrada
pode superar uma conclusão inventada; abstenção vazia pode perder por baixa
utilidade.

## Observável de catálise

Qualidade final e mecanismo de catálise são avaliados separadamente.

O julgamento cego recebe somente as duas respostas anonimizadas. Depois de os
scores e o receipt cego serem congelados, um auditor de trace pode classificar
cada direção material do Bateia:

- `CONFIRMED_BY_WEB`: a pesquisa encontrou suporte verificável;
- `REJECTED_BY_WEB`: a pesquisa encontrou contraevidência ou desalinhamento;
- `NARROWED_BY_WEB`: a direção ampla virou hipótese ou consulta mais precisa;
- `USED_WITHOUT_WEB_SUPPORT`: entrou na resposta sem confirmação suficiente;
- `NOT_USED`: não aparece na pesquisa ou na resposta;
- `UNKNOWN`: o trace não permite ligar direção, consulta e afirmação.

Rejeitar uma direção ruim com evidência é catálise potencialmente útil. Contar
somente candidatos aceitos premiaria concordância e jogaria `INCONCLUSIVE`
contra o tratamento.

O laudo publica, por par:

- delta de qualidade cega;
- cobertura das perguntas de evidência;
- validade e força das fontes;
- direções confirmadas, rejeitadas, estreitadas e ignoradas;
- alegações materiais sem suporte;
- tempo, chamadas, tamanho dos payloads e tokens expostos;
- ganho por unidade de custo quando o denominador for conhecido;
- `Unknown` explícito quando ligação causal ou uso não forem observáveis.

Mais consultas, candidatos ou texto não contam como ganho por si sós.

## Preflight sem sessão comportamental

Antes do selo e novamente antes do primeiro despacho, o preflight deve:

1. verificar hashes do caso, prompts, schema, Bateia, CLI e perfil de execução;
2. provar `CODEX_HOME` contendo somente autenticação e diretório de trabalho
   vazio, sem registrar bytes de credencial;
3. provar a mesma pesquisa nativa e os mesmos budgets nos dois braços;
4. provar zero configuração MCP no controle e exatamente um servidor/uma tool
   no tratamento;
5. executar somente `initialize` e `tools/list` no Bateia;
6. executar o decoder sobre o `events.jsonl` R3 congelado e recuperar três
   máquinas start/completion `DUAL_ID_CODEX_0_149` sem parser last-wins;
7. classificar o output R3 de `22.578` caracteres como
   `VALID_SCHEMA_EDITORIAL_TARGET_EXCEEDED`, não como output apagável;
8. simular primeiro braço inválido e provar que o segundo continua agendado;
9. provar o teto rígido máximo do novo schema;
10. produzir recibo `PASS` antes de qualquer modelo, web ou `tools/call`.

Qualquer falha produz `FAIL_CLOSED_ZERO_BEHAVIORAL_SESSIONS`.

## Piloto de três pares

Após contrato, parser, oráculos, mutantes e preflight selados, reservar três
casos de `VALIDATION`, de famílias diferentes, para um único perfil congelado.

```text
3 casos × 2 braços = 6 sessões comportamentais máximas
```

Contabilidade obrigatória:

| Estado | Sessões globais |
|---|---:|
| Antes deste passo | 16 |
| Reserva máxima do piloto | 6 |
| Máximo após o piloto | 22 |
| Teto herdado do Passo 0029 | 108 |

Escrever ou aprovar este passo não autoriza os seis despachos. A execução exige
autorização explícita posterior e não aumenta o teto de `108`.

## Poder discriminatório

Antes do selo, positivos, negativos e opacos devem provar que o contrato:

- preserva um par completo com controle web e tratamento Bateia+web;
- rejeita controle com qualquer MCP ou artifact Bateia;
- rejeita tratamento sem uma recomendação concluída antes da primeira web;
- rejeita `0` ou mais de `3` chamadas web em qualquer braço;
- rejeita diferença de modelo, esforço, caso, schema ou budget dentro do par;
- rejeita parser JSON convencional que colapse chaves `item.id` duplicadas;
- recupera corretamente o formato `DUAL_ID_CODEX_0_149` da fixture R3;
- classifica formato de ID não congelado como `Unknown`, não como sucesso;
- rejeita pareamento start/completion com IDs ou ordem divergentes;
- rejeita cancelar o segundo braço por falha comportamental do primeiro;
- rejeita retry ou substituição de resposta ruim;
- preserva output válido acima de `18.000` como excesso editorial;
- rejeita output acima do teto rígido demonstrado pelo schema;
- rejeita converter `INCONCLUSIVE` em derrota ou vitória automática;
- rejeita dar ao julgador cego tool trace, braço, ordem ou identidade;
- rejeita atribuir catálise sem elo observável direção→web→resultado;
- rejeita promover quantidade de busca, texto ou candidatos a ganho;
- rejeita executar sessão `23` sem novo estado de contagem compatível;
- preserva `Unknown` em telemetria, custo ou ligação causal opaca.

Mutation score obrigatório: `1.0`, com repetição em ordem declarada e reversa.
Duas revisões consecutivas sem novo mutante rejeitado ou sem mudança da causa
dominante interrompem a calibração e reabrem o desenho.

## Artefatos previstos

Materializar sob `00_nucleo/benchmark/0030/`:

- `intent.json` e hashes das evidências R3;
- `experiment-manifest.json`;
- `compact-output-schema.json` e prova do teto rígido;
- `web-event-decoder-contract.json`;
- fixtures positivas, negativas e opacas;
- `pair-dispatch-contract.json`;
- `prompt-control.txt` e `prompt-treatment.txt`;
- `execution-profile.json` e `randomization-plan.json`;
- `rubric.json` e `catalysis-trace-rubric.json`;
- contrato, verificação discriminatória e selo;
- preflight e verificação independente;
- receipts, eventos e outputs por sessão;
- projeções cegas e julgamentos;
- auditoria pós-reveal de catálise;
- resultado, laudo e certificado.

O código sucessor respeita o reticulado:

- tipos, scanner JSON duplicate-aware, máquina de estados, políticas de
  pareamento e análise pura ficam em `01_core/experiment/catalysis/`;
- CLI, dispatcher abstrato e apresentação de receipts ficam em
  `02_shell/experiment/catalysis/`;
- adaptadores concretos de processo, relógio, filesystem e captura JSONL ficam
  em `03_infra/experiment/catalysis/` e implementam portas de L1;
- `04_wiring/step0030_catalysis/` instancia adaptadores e injeta dependências,
  sem parser, política experimental ou lógica de negócio.

Cada arquivo recebe linhagem para o prompt espelhado correspondente sob
`00_nucleo/prompts/`.

## Gates de aceitação

- os quatro hashes R3 protegidos permanecem iguais;
- decoder preserva chaves duplicadas e passa as fixtures reais e mutantes;
- schema possui prova determinística de teto rígido compartilhado;
- excesso editorial não invalida nem suprime o par;
- falha comportamental do primeiro braço não cancela o segundo;
- ambos os braços recebem pesquisa nativa `1..3` sob o mesmo budget;
- tratamento recebe somente `bateia_recommend` e a conclui antes da web;
- `INCONCLUSIVE` permanece neutro em contrato, rubrica e análise;
- anonymizer é determinístico e não corrige conteúdo;
- julgamento de qualidade congela antes do acesso ao trace;
- auditoria de catálise exige ligação observável e aceita rejeição útil;
- mutation score `1.0` permanece estável sob repetição e ordem reversa;
- contagem máxima termina em `22/108`, sem retry ou replacement;
- `crystalline-lint .` não introduz nova violação, `go test ./...` não introduz
  nova falha e o pacote `./04_wiring` compila;
- verificador independente emite certificado sem editar entradas julgadas.

## Regime Tekt previsto

Aplicar protocolo completo.

- o usuário é autor da intenção de medir catálise sobre pesquisa web;
- autor do contrato lê este passo e a fixture R3, não o candidato;
- autor dos oráculos congela casos positivos, negativos e opacos antes do código;
- implementador recebe somente prompt, contrato e selo;
- dispatcher recebe perfil e fila selados, não rubrica protegida;
- julgador cego não recebe mapping, traces ou identidade dos braços;
- auditor de trace atua somente após o receipt cego;
- adversário ataca parser, pareamento, budgets, neutralidade e análise;
- verificador não corrige os artefatos que julga.

Sem isolamento material dessas entradas e capacidades, o laudo usa
`EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO` mesmo que todos os gates passem.

## Fora de escopo

- executar as seis sessões apenas porque este passo foi escrito;
- reabrir, corrigir ou renomear a tentativa R3 como comparação válida;
- provar ganho geral do Bateia com somente três pares;
- comparar Bateia sem internet contra IA com internet;
- premiar automaticamente decisão conclusiva ou penalizar `INCONCLUSIVE`;
- compactar silenciosamente o payload atual do Bateia no harness;
- alterar ferramentas MCP ou código de produção do Bateia;
- usar mais de três pesquisas nativas por braço;
- adicionar modelos ou esforços ao piloto focal;
- ampliar o teto global de `108`;
- publicar credenciais, chain of thought ou conteúdo protegido de tools.
