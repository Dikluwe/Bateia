# Laudo do Passo 0030 — correção do teste de catálise da busca

**Regime:** protocolo Tekt completo  
**Executor do veredito:** `/root/verificador_final_0030`  
**Atestação máxima:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Veredito do harness:** `EXECUTADO_COM_EVIDENCIA_REPRODUZIVEL`  
**Validade do estimando:** `NAO_IDENTIFICADO_POR_FALHA_DE_ADERENCIA_DO_TRATAMENTO`  
**Veredito de produto:** `BATEIA_CATALYSIS_NOT_ESTIMABLE`

## Conclusão

O harness corrigido completou as seis células comprometidas, uma vez cada, na
ordem selada, sem retry ou substituição, com uma sessão global em voo e
intervalos de início de `125, 151, 121, 129 e 131` segundos. Os seis outputs
são objetos JSON válidos no schema compacto, ficam abaixo dos limites editorial
e rígido e são byte-idênticos à última mensagem registrada. O decoder
duplicate-aware preservou as duas ocorrências de `item.id` e formou 14 máquinas
`DUAL_ID_CODEX_0_149`, sete por braço.

O estimando pretendido, contudo, não foi observado. As três chamadas
`bateia_recommend` terminaram `failed` com `MCP -32602 Invalid params`; nenhuma
produziu resultado, recomendação ou direção material. Isso viola a obrigação
`TREATMENT_BATEIA_NOT_COMPLETED_BEFORE_WEB`. Uma chamada terminal falha não pode
ser reclassificada como direção confirmada, rejeitada, estreitada, ignorada ou
usada sem suporte. Assim, há zero caminhos observáveis
direção Bateia → web posterior → resultado e não se pode alegar catálise.

Essa é uma falha metodológica e de aderência do tratamento, não vitória
automática do controle. O experimento prova que a correção de despacho por
intenção de tratar preserva ambos os braços e torna a falha auditável; não prova
o efeito de disponibilizar um Bateia funcional.

## Julgamento cego após o reveal

O commitment foi criado antes do julgamento, o receipt cego foi congelado antes
do reveal e o SHA-256 comprometido (`9186bd…6f53f`) coincide exatamente com os
bytes do arquivo revelado. A ordem observável por birth time foi commitment em
`21:46:37Z`, receipt em `21:50:44Z` e reveal em `21:51:11Z`. Os seis hashes das
cópias opacas coincidem com os outputs-fonte. Essa cadeia sustenta a ordem
registrada, mas não equivale a timestamp externo nem a isolamento forte.

Após revelar o mapping:

| Caso | Controle | Tratamento | Delta T−C | Veredito cego revelado |
|---|---:|---:|---:|---|
| `EXP29-C01` | 24 | 18 | -6 | output rotulado controle melhor |
| `EXP29-V01` | 17 | 23 | +6 | output rotulado tratamento melhor |
| `EXP29-V02` | 23 | 22 | -1 | output rotulado controle melhor |

As somas são apenas uma descrição das sete notas `0..4`: controle `64`,
tratamento `63`, delta total `-1` e média pareada `-0,3333`. Houve dois
vereditos a favor do output rotulado controle e um a favor do rotulado
tratamento. Com três pares exploratórios e zero recomendações Bateia, esses
números não são estimativa confirmatória, não demonstram superioridade do
controle e não podem ser atribuídos a catálise.

## Catálise e custo

Cada tratamento executou uma tentativa Bateia antes da web, mas as três
tentativas falharam e emitiram zero direções. Por isso, a unidade congelada na
rubrica de trace — cada direção material — tem cardinalidade zero. O gate de
atribuição falha por ausência do elo observável, e não por score ruim ou por
`INCONCLUSIVE`.

Os controles consumiram `216.432` tokens de entrada, `10.537` tokens de saída e
`5.363` tokens de raciocínio; os tratamentos, `284.906`, `12.709` e `6.423`,
respectivamente. As somas de duração superior observável foram `234 s` e
`291 s`. O custo monetário não foi registrado, logo ganho por unidade de custo
permanece `Unknown`. Volume de tokens, texto ou chamadas não conta como ganho.

## Integridade e gates

Foram recomputados sem mismatch os hashes do selo R2, contratos e rubricas
preservados, oráculos, fixtures R3, implementação, prompts, invocações e todos os
artefatos das seis sessões. A derivação do mapping confere para os três casos.
O budget terminou em `22/108`, a partir de `16/108`, com seis sessões e zero
retry.

Reexecuções locais:

- `go test` dos pacotes focais L1–L4 do Passo 0030 — `PASS`;
- `go test ./04_wiring/step0029_experiment` — `PASS`;
- `go build ./04_wiring/...` — `PASS`;
- `git diff --check` — `PASS`;
- `crystalline-lint .` — exit `0`, somente warnings preexistentes e nenhum do
  Passo 0030;
- `go test ./...` — reproduziu a falha semântica global congelada do gate R4
  0026 (`retrocompatibility got=SURVIVED want=Violated`). Três pacotes de
  `httptest` também foram impedidos pelo sandbox de abrir socket loopback; o
  preflight registra esses pacotes passando fora do sandbox. Nenhuma falha nova
  pertence ao candidato 0030.

## O que os três pares permitem concluir

- o harness corrige o cancelamento do segundo braço e preserva evidência mesmo
  após falha comportamental;
- ordem, intervalos, zero retry, budgets, schema, telemetria web e blinding são
  reproduzíveis no fragmento observado;
- a integração Bateia exata do piloto falhou em `3/3` tratamentos com parâmetros
  inválidos;
- as seis respostas continuam julgáveis e exibem qualidade cega mista.

Eles não permitem concluir uplift do Bateia, catálise útil, superioridade do
controle, equivalência funcional geral nem razão de ganho por custo. Mesmo com
aderência perfeita, três pares de calibração não sustentariam generalização de
produto.

## Limites de atestação

O workspace, o contexto e o filesystem foram compartilhados. A segregação de
papéis, os hashes e a ordem causal registrada são evidência de disciplina e
identidade, não prova forte de isolamento de capacidades. Conforme a skill de
materialização segregada, `Unknown` foi preservado e o certificado distingue a
execução bem-sucedida do harness da validade falha do estimando.

Artefatos canônicos: verificação pós-reveal `9e0e4f…72ee0d`; resultados
`259e3b…f82ed8`. Qualquer alteração nas entradas seladas, sessões, cópias cegas,
receipt ou reveal invalida este laudo.
