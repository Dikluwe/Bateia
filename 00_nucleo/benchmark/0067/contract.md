# Contrato canônico 0067 — Planejador de fronteira mínima

Autoridade normativa:
`../../passos/0067-planejar-fronteira-minima-do-radar-de-ideias.md`. Este recorte congela
observáveis e medição antes do candidato.

## Fragmento produtivo

`PlanIdeaRadarProbes(IdeaRadarInput) IdeaRadarPlanReport` chama o radar 0066 e produz um
plano por assessment. Não muda input, limiar, classificação, draft, consumidor ou
veredito 0066. Se o radar é `UNKNOWN`, propaga status/razão sem plano parcial.

Probes: `RESOLVE_CONFLICT`/10, `RESOLVE_AMBIGUITY`/20, `RECORD_ORIGIN`/30,
`ADD_PERSPECTIVE`/40, `ADD_INDEPENDENT_ORIGIN`/50 e `BIND_EVIDENCE`/60. Requisitos de
origem: `NONE`, `DISTINCT_RECORDED`, `RECORDED`. JSON é `snake_case`, sem `omitempty`;
arrays de output nunca são nulos. Report schema `idea-radar-plan-report-v1`, status
`COMPLETE|UNKNOWN`, razão completa `planned`, `VerdictChanges=0`.

Somente o bloqueador ativo gera probes. Conflito e ambiguidade emitem um probe por sinal
canônico e ocultam tarefas inferiores. Perspectivas insuficientes escolhem os primeiros
kinds ausentes na ordem P/B/A/C/O e combinam cada adição com nova origem distinta;
somente o déficit residual gera probes de origem. Origem desconhecida escolhe o menor ID
por kind ainda não registrado até três kinds, depois cobre déficit residual. Origem
insuficiente gera exatamente `MissingOrigins`. Kinds para origem adicional são escolhidos
ciclicamente entre kinds de suporte na ordem P/B/A/C/O.

Convergência gera somente um `BIND_EVIDENCE`, com suportes e papéis do draft 0066;
continua sem `EvidenceSource`, chamada ao 0065 ou readiness. O planner não promete
resultado do probe. Razões, ordenação e IDs `target::pNNNN` são os definidos no passo.
A minimalidade é local a cada target; coalescer um probe compartilhado entre targets não
faz parte deste fragmento.

## Corpus congelado

O autor dos oráculos materializa os inputs a partir de `baseline.json`. Todo novo suporte
simulado usa ID e OriginID únicos. Respostas bem-sucedidas removem conflito, convertem
ambiguidade em suporte ou adicionam/registram o suporte indicado, e então o radar é
reexecutado.

| Caso | Família | Estado | Planner | Ingênuo | Depois |
|---|---|---|---:|---:|---|
| C01 | coverage-gap | P/or-a | 2 add-perspective distintas | 4 | convergente |
| C02 | coverage-gap | P/or-a, B/or-b | 1 add-perspective distinta | 2 | convergente |
| C03 | provenance-gap | P/a, B/b, A/Unknown | 1 record-origin | 1 | convergente |
| C04 | provenance-gap | P/a, B/b, duas A/Unknown | 1 record-origin | 2 | convergente |
| C05 | provenance-gap | P/B/A em `shared` | 2 add-origin | 2 | convergente |
| C06 | signal-resolution | P/B/A distintos + 2 conflicts | 2 resolve-conflict | 2 | convergente |
| C07 | signal-resolution | P/B/A distintos + ambiguous | 1 resolve-ambiguity | 1 | convergente |
| C08 | evidence-binding | convergente | 1 bind-only | 1 | não aplicável |
| C09 | safety-control | P/a + conflict | 1 resolve-conflict | 1 | perspectivas insuficientes |
| C10 | opaque-control | schema inválido | 0/UNKNOWN | 0 | UNKNOWN |

Nos oito não convergentes elegíveis, planner e baseline removem 8/8 bloqueadores e 7/8
chegam a convergente. Esforço: 11 contra 15; redução 4/15. C01, C02 e C04 melhoram em
duas famílias. Incluindo o binding C08 são 12 contra 16, valor não usado no gate.

## Baseline e veredito externo

O baseline resolve sinais ativos, mas soma separadamente `MissingProjectionKinds` e
`MissingOrigins`, e registra toda observação Unknown. Cada probe vale uma unidade.
Medição e simulação ficam apenas no harness.

`OBSERVED`: resolução do planner >= baseline e (redução de esforço >=0,20 ou ganho de
resolução >=0,10 sem mais esforço), com melhoria estrita em >=2 casos de >=2 famílias,
sem ação proibida, falsa promessa, budget ou impacto decisório. Métrica indefinida,
entrada ausente, menos de dois casos/famílias ou `UNKNOWN` dá `INCONCLUSIVE`; demais
resultados válidos dão `NOT_OBSERVED`.

## Limites e gate

Herdar 64 projeções, 1024 observações, 512 targets e demais limites 0066; adicionar 512
planos, 1024 probes e 8192 operações totais. Operações = radar + uma por projeção
indexada, observação normalizada, assessment inspecionado e probe emitido. Excesso:
`UNKNOWN/budget_exceeded`, sem parcial.

Oráculos cobrem C01–C10, inclusivo/exclusivo, Unknown, imutabilidade, JSON não nulo,
operações e 20 permutações. O gate elimina M01–M16 do passo com mutation score 1.0,
confirma métricas exatas e zero impacto. A conclusão cobre somente planejamento
estrutural no corpus sintético, não resposta real, custo humano ou utilidade geral.
