# Contrato canônico 0066 — Radar de convergência de ideias

Autoridade normativa:
`../../passos/0066-prototipar-radar-de-convergencia-de-ideias.md`. Este recorte congela
os observáveis e a medição antes da implementação candidata.

## Fragmento observável produtivo

`usecase.AnalyzeIdeaRadar` recebe `idea-radar-input-v1` e devolve
`idea-radar-report-v1`. Projeções e observações já vêm estruturadas; a função não deriva
semântica de texto. O alvo é a chave canônica exata `candidate_id#concept_id`.

Convergência exige ao mesmo tempo três `IdeaProjectionKind` distintos e três `OriginID`
distintos, registrados. Diversidade de termos, observações, IDs, hashes, caminhos ou
queries não substitui nenhum limiar. A classificação segue: conflito, ambiguidade,
perspectivas insuficientes, origem desconhecida, origens insuficientes e convergência.

Canonicalização de IDs limita-se a trim ASCII e lowercase ASCII; caracteres aceitos são
`[a-z0-9._:/-]`, tamanho 1..128. `OriginID` é opaco, case-sensitive, trimado apenas para
validação e nunca reescrito. Observação duplicada pela tupla semântica não multiplica
contagens. Outputs são ordenados, arrays não nulos e independentes da ordem do input.

O draft 0065 registra a hipótese, os suportes e os papéis
`IMPLEMENTATION,INDEPENDENT_TEST`, mas permanece sempre
`ReadyFor0065=false/requires_evidence_binding`. `VerdictChanges=0`; nenhum consumidor
legado é chamado ou alterado.

Limites: 64 projeções, 1024 observações, 512 alvos, 128 bytes por ID/origem, 64 itens e
4096 bytes por dimensão textual de projeção e 4096 operações. Operações são projeções
validadas + observações normalizadas + associações canônicas. Excesso é relatório
`UNKNOWN/budget_exceeded` sem parcial. O algoritmo não enumera trincas ou power set.

## Corpus mínimo congelado

O autor do oráculo materializa observações a partir desta tabela. `P/B/A/C/O` abreviam
`PROBLEM/BEHAVIOR/ARCHITECTURE/CONSTRAINT/OPERATION`. Cada letra listada representa um
suporte; sufixos são OriginIDs. Alvos auxiliares são proibidos no ground truth.

| Caso | Família | Alvo principal e construção radar | Auxiliar obrigatório | Outcome principal |
|---|---|---|---|---|
| C01 | event-recovery | `cand-log#event-log-recovery`: P/or-a, B/or-b, A/or-c | duplicata semântica do suporte B | CONVERGENT |
| C02 | event-recovery | `cand-journal#write-ahead-log`: B/or-d, C/or-e, O/or-f | `noise-cache#cache`: P/or-x, B/or-y | CONVERGENT; auxiliar INSUFFICIENT_PERSPECTIVES |
| C03 | authorization | input com caixa/whitespace que canoniza para `cand-policy#policy-enforcement`: P/or-g, A/or-h, C/or-i | `noise-auth#acl`: P/or-j, B/UNKNOWN, A/or-k | CONVERGENT; auxiliar UNKNOWN_ORIGIN |
| C04 | authorization | `cand-rbac#role-based-access`: P/shared, B/shared, A/shared | `noise-role#role-map`: três suportes e um CONFLICT | INSUFFICIENT_ORIGINS; auxiliar CONFLICT |
| C05 | offline-sync | `cand-sync#operation-log-sync`: B/or-l, C/or-m, O/or-n | repetição de termo e IDs físicos diferentes não cria origem | CONVERGENT |
| C06 | offline-sync | `cand-crdt#conflict-free-merge`: P/or-p, B/or-q, A/or-r e uma observação AMBIGUOUS | `noise-merge#last-write-wins`: três suportes e um AMBIGUOUS | AMBIGUOUS nos dois alvos |

Acrescentar controles separados para ID Unicode/inválido, ID duplicado contraditório,
schema desconhecido, projeções/observações vazias, cada limite inclusivo/exclusivo,
imutabilidade e 20 permutações com nil/vazio textual.

## Medição externa congelada

O ground truth e as predições isoladas vivem em `baseline.json`, nunca em produção.
Usar correspondência exata de TargetID. A predição radar contém somente assessments
`CONVERGENT`; conflito, ambiguidade e Unknown não são positivos.

Para cada conjunto, `TP=preditos∩esperados`, `FP=preditos-esperados`,
`FN=esperados-preditos`; precisão `TP/(TP+FP)`, recall `TP/(TP+FN)`, F1
`2TP/(2TP+FP+FN)`, falsa convergência `FP/(TP+FP)` e contaminação
`preditos∩proibidos / preditos`. Divisão por zero produz `Unknown`. F1 por caso usa a
forma direta e é definida porque todo caso congelado possui ao menos um esperado.

Escolher a melhor projeção isolada globalmente por F1 micro, depois precisão, recall e
nome canônico do kind. O baseline congelado resulta em `ARCHITECTURE`: TP=3, FP=1,
FN=3, precisão=0.75, recall=0.50, F1=0.60, falsa convergência=0.25 e
contaminação=0.25.

O resultado radar esperado pelo contrato tem TP=4, FP=0, FN=2, precisão=1.00,
recall=2/3, F1=0.80, falsa convergência=0 e contaminação=0. Ganhos: precisão +0.25,
recall +1/6 e F1 +0.20. C02 e C05 melhoram F1 em duas famílias diferentes.

O veredito externo é `OBSERVED` somente com precisão sem regressão, ganho >=0.10 em
recall ou F1 micro e melhoria de F1 em >=2 casos de >=2 famílias. É `INCONCLUSIVE` se
qualquer métrica necessária ao gate é indefinida, entradas/ground truth estão ausentes,
há menos de dois casos ou famílias, schema diverge ou o radar excede budget. Demais
resultados válidos são `NOT_OBSERVED`. O oráculo deve testar também um vetor
`NOT_OBSERVED` e um `INCONCLUSIVE`; não basta aprovar o vetor esperado.

## Gate discriminatório

Eliminar M01–M14 do passo, com testemunha pública, `mutation_score=1.0`. Exigir corpus
produtivo 100%, métricas externas exatas, seleção de baseline exata, 20/20 permutações,
operações exatas, input imutável, draft sempre não pronto e zero mudanças em vereditos.

O contrato cobre apenas associação conceitual explicitamente anotada e utilidade nesse
corpus. Não mede descoberta real, qualidade de geração de projeções, independência causal
real ou melhoria geral de decisões.
