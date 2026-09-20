# Execução E-0012 — Autofagia competitiva e Katamari real

**Data**: 2026-08-26  
**Veredito**: `FAIL_DISCOVERY`  
**Autofagia local**: `INCONCLUSIVE`  
**Katamari real**: `NÃO EXECUTADO`  
**Atestação de cegamento**: indisponível; a mesma sessão conhecia o registry

## Entradas

- IdeaSpec: `../experimento/bateia-idea.json`;
- estratégia: `lexical`;
- expansões: 16;
- limite solicitado: 3 candidatos por expansão;
- concorrência: 2;
- provedor: Sourcegraph público;
- alvo local: `/repos/Antigravity/Bateia`;
- política: `../experimento/reuse-policy.json`.

O plano literal e o enriquecido foram idênticos:

```text
plan.json      2c85951689ca5757761c3f044c7f6a01335365df628c34ab3690e40b034ca945
expansion.json 2c85951689ca5757761c3f044c7f6a01335365df628c34ab3690e40b034ca945
```

Isso indica que o vocabulário lexical atual não acrescentou termos à IdeaSpec do
experimento.

## Fase A — descoberta pública

### Tentativa no sandbox

A execução falhou nas dezesseis expansões por bloqueio de DNS/rede:

```text
Post https://sourcegraph.com/.api/graphql:
dial udp 127.0.0.53:53: socket: operation not permitted
```

### Tentativa com rede autorizada

A mesma execução foi repetida fora do sandbox. O provedor não produziu snapshot,
resultado parcial ou erro durante aproximadamente cinco minutos. A execução foi
interrompida e terminou com `signal: interrupt`.

Nenhuma semente do registry foi inserida manualmente e nenhum candidato externo
foi coletado. Pelas regras pré-registradas, isso é `FAIL_DISCOVERY`; não é
evidência de que os concorrentes estejam ausentes.

## Fase C — execução local do Bateia sobre o Bateia

O modo local concluiu e gerou o dossiê
`local-dossier/` com identidade:

```text
local/bateia-a061b61e@287320b0297d664fa441b70b7b09f8ecf251896c79eba01f55b55bc9222f4f1b
```

Observáveis da coleta:

| Observável | Resultado |
|---|---:|
| bundles | 1 |
| coleta completa | sim |
| árvore completa | sim |
| itens | 40 |
| bytes | 95.991 |
| falhas | 0 |

Estados observados para as dez capacidades:

| Estado | Quantidade |
|---|---:|
| `PROVEN` | 0 |
| `PARTIAL` | 0 |
| `ABSENT` | 0 |
| `INCONCLUSIVE` | 10 |
| `CONFLICT` | 0 |

A decisão final foi `INCONCLUSIVE`. Como a coleta foi completa e não teve falhas,
a causa dominante observada é cobertura insuficiente do catálogo de regras para
as formulações de produto usadas na IdeaSpec. Por exemplo, existe regra para
“offline”, mas não há regras que reconheçam diretamente “composição mínima de
projetos”, “decisão de partida” ou “autofagia competitiva”. Isto é um finding de
regra/oráculo, não autorização para afrouxar `PROVEN`.

Hashes principais:

```text
assessment-report.json  bd66fac05682a611ffa442aff84aa18708d8233fb21a0891d25f72edf3aab1c6
discovery-snapshot.json ad15599b42a1152e8716a3b8371648102c218ff06f7a56e7c5ed86787a5f16ef
bundles.json            24f9df2b1dce2fcbc45aa79774960d7f5b2c90c36f7c674d9ea8e93739e756c1
reuse-decision.json     b4684a82af7401f7e0ed00a13a820df62af32edabfcf4de40ef8d0e67bb73d9f
```

## Fases não executadas

- controles negativos: dependiam de descoberta pública capturável;
- bundles dos concorrentes: nenhum candidato externo foi descoberto;
- oráculos: o formulário permanece `UNFILLED`;
- `idea autofagia`: sem oráculos/regras correspondentes, apenas repetiria as dez
  inconclusões já materializadas pelo assessment;
- `idea katamari`: não existe matriz externa correspondente para compor;
- métricas de recall, rejeição, cobertura adicional e utilidade: `null`.

Executar Katamari com a lista conhecida de concorrentes violaria o teste, pois os
inseriria depois da falha de descoberta.

## Findings

### F-E0012-01 — Provedor sem limite observável

`idea discover` não encerrou nem escreveu progresso parcial quando o provedor
permaneceu sem resposta. A fronteira precisa de timeout total e receipt das
tentativas já concluídas.

### F-E0012-02 — Expansão lexical nula

A estratégia lexical produziu o mesmo plano literal para esta IdeaSpec. O teste
não ganhou vocabulário alternativo para encontrar a categoria competitiva.

### F-E0012-03 — Autofagia não reconhece o produto

Mesmo com bundle local completo, nenhuma das dez capacidades foi comprovada. O
catálogo atual é adequado aos casos anteriores de arquitetura e operação, mas
não cobre ainda as capacidades do próprio Bateia em nível de produto.

## Próxima execução válida

Antes de repetir E-0012:

1. materializar timeout e receipts parciais em `idea discover`;
2. criar oráculos humanos para capacidades internas usando os bytes já coletados;
3. propor regras específicas em passo separado, sem alterar este resultado;
4. repetir descoberta a partir da mesma IdeaSpec e dos mesmos limites;
5. somente depois coletar concorrentes e executar Katamari.

O registry e os expected outcomes deste teste não devem ser ampliados para
absorver candidatos encontrados numa repetição.

## Reexecução após o Passo 0013

**Autorização externa**: concedida explicitamente pelo usuário  
**Snapshot**: `discovery-0013.json`  
**SHA-256**: `c55b9c0e8d2f038f8c28ce4c3f371884972e515107fa8238b5de85fcefd01671`

A Fase A foi repetida sem alterar IdeaSpec, queries, estratégia, limites ou
registry:

```text
strategy=lexical
max-per-query=3
concurrency=2
timeout=3m
attempt-timeout=20s
```

O comando terminou por `TOTAL_TIMEOUT`, gravou o snapshot antes do erro e
produziu:

| Observável | Resultado |
|---|---:|
| expansões planejadas | 16 |
| iniciadas | 10 |
| concluídas | 10 |
| não iniciadas | 6 |
| attempts expirados | 18 |
| attempts bem-sucedidos | 1 |
| candidatos preservados | 3 |
| falhas por expansão | 9 |
| warning | `NO_SEMANTIC_EXPANSION` |

O receipt ficou `PARTIAL/TOTAL_TIMEOUT`. A correção operacional funcionou: o run
terminou no deadline, mostrou progresso e não perdeu o único resultado concluído.

### Candidatos observados

Todos vieram do fallback genérico `decidir entre`:

| Repositório | Evidência recuperada | Correspondência |
|---|---|---|
| `VgerMods/Pawn@HEAD` | tradução sobre decidir entre objetos num jogo | irrelevante e opaco |
| `beduExpert/Java-Standar-Edition-2-2025@HEAD` | material sobre escolher REST ou mensageria | adjacente pontual, não concorrente, opaco |
| `mozilla-l10n/styleguides@HEAD` | guia de tradução sobre decidir entre abreviações | irrelevante e opaco |

Nenhum aparece no registry positivo. Entre as quatro sementes com repositório
público verificável, nenhuma foi encontrada:

```text
competitor_discovery_recall = 0 / 4 = 0%
```

Os candidatos não seguem para coleta porque `HEAD` viola a correspondência
imutável e nenhum é concorrente funcional. Não foram substituídos manualmente.

### Veredito atualizado

- controle operacional do Passo 0013: `PASS`;
- descoberta competitiva E-0012: `FAIL_DISCOVERY`;
- principal perda agora atribuível: expansão lexical nula e fallback curto
  semanticamente genérico;
- Katamari: continua `NÃO EXECUTADO` por ausência de candidato correspondente.

O novo receipt remove a ambiguidade da primeira execução: o problema não era
somente duração indefinida. Dentro do budget, 18 consultas expiraram e a única
que respondeu produziu ruído previsível a partir de “decidir entre”.
