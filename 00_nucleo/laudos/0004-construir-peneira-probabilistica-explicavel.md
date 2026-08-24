# Laudo 0004 — Peneira probabilística explicável

**Data**: 2026-08-24  
**Passo verificado**: [0004 — Construir peneira probabilística explicável](../passos/0004-construir-peneira-probabilistica-explicavel.md)  
**Resultado**: `APROVADO_COM_RESTRIÇÕES`

## O que foi implementado

- snapshot schema v2 com origem por candidato, tentativa, query efetiva e evidência;
- migração pura v1 → v2 sem fabricar tipo, path ou snippet;
- captura Sourcegraph de `Repository` e `FileMatch`, com path e snippet normalizado
  limitado a 240 caracteres;
- nove sinais separados: prior da expansão, retenção da query, fallback,
  especificidade, alinhamento da descrição, convergência, capacidades, `NonGoals`
  e completude;
- modelo linear cujas contribuições reproduzem integralmente o ranking;
- decisões `ADVANCE`, `AUDIT` e `DROP`, com ausência diferente de conflito;
- amostragem determinística de descartes;
- calibrador Platt com validação `leave-one-idea-out`;
- comparação fora da amostra com prevalência, Brier, log loss, precision@k e bins;
- replay diagnóstico que encaminha relevantes conhecidos sem alterar seu score;
- comandos offline `idea screen`, `screening calibrate` e `screening audit`;
- seleção e revisão de 33 candidatos por README.

## O que a trajetória causal revelou

Uma captura v2 de impacto estrutural mostrou três mecanismos concretos de ruído:

1. **tradução/localização** — `nirholas/crypto-vision` apareceu por textos em
   `agents/locales/...` que continham “mostrar impacto” fora do domínio desejado;
2. **coincidência incidental** — `aventuras-en-megadrive` apareceu por
   `activa_impacto // para mostrar impactos` num jogo;
3. **código vendorizado** — projetos Kubernetes apareceram pelo mesmo comentário
   em `vendor/k8s.io/apiserver/...`: “changes to this code impact API...”.

Em contraste, `dejo1307/enola` apareceu na primeira tentativa de
`change blast radius`, no arquivo `pkg/explain/explain.go`, com um símbolo que
descreve raio de impacto por acoplamento.

O adaptador agora preserva essa diferença. Paths em `vendor`, `node_modules` e
localizações recebem especificidade reduzida, em vez de o mero `FileMatch` contar
como evidência forte.

## Calibração

Foram usadas 57 observações conclusivas; os cinco `uncertain` foram excluídos da
classe negativa. Cada fold reteve uma ideia inteira.

| Métrica fora da amostra | Resultado |
| :---------------------- | --------: |
| Brier do modelo | 0,11274 |
| Brier da prevalência | 0,12311 |
| Log loss | 0,35526 |
| Precision@10 | 10% |

O modelo superou a prevalência em Brier agregado por margem de 0,01037 e, conforme
o critério declarado, recebeu `calibration_status: calibrated`. Isso autoriza
`probability` somente para a população estreita deste benchmark.

A generalização ainda é fraca:

- o fold de receitas foi pior que a prevalência: 0,10418 contra 0,09239;
- precision@10 fora da amostra foi apenas 10%;
- 35 observações ficaram na faixa 0–20%, com probabilidade média de 4,00% e taxa
  observada de 11,43%, indicando subestimação;
- 22 ficaram na faixa 20–40%, com média de 22,64% e taxa observada de 18,18%;
- nenhuma faixa acima de 40% teve cinco observações.

Portanto, a probabilidade é uma estimativa experimental local, não garantia nem
modelo pronto para outros domínios.

## Prova de mordida

Os testes controlados demonstram:

- A, com query completa, match específico e duas origens, supera B;
- B registra fallback genérico;
- C, sem conteúdo suficiente, vira `AUDIT`;
- D, com conflito de `NonGoals`, vira `DROP` apesar do sinal positivo;
- remover uma origem altera convergência sem contaminar os demais sinais;
- um modelo não calibrado serializa `probability: null`;
- score pré-clamp é reconstruído pela soma das contribuições.

Na captura v2 real, `dejo1307/enola` recebeu:

```text
decision: ADVANCE
probability: 0,64338
```

Os seis ruídos da mesma captura receberam `DROP`, entre 2,79% e 16,14%.

## Auditoria dos descartes

A seleção determinística incluiu:

- oito relevantes conhecidos;
- cinco incertos;
- dez irrelevantes entre os menores scores;
- dez irrelevantes mais próximos do avanço.

Os 33 READMEs foram consultados e suas referências preservadas. Resultado:

```text
revisados: 33
falsos negativos entre os controles descartados: 0
incertos resolvidos como irrelevantes: 5
```

Achados laterais importantes:

- `mllint` é relevante para análise estática, mas declara manutenção encerrada;
- `supermodeltools/cli` é relevante para grafos e blast radius, mas envia o
  repositório para uma API;
- `reviewdog/errorformat` é relevante como componente de integração de
  diagnósticos, não como analisador completo.

Esses fatores pertencem à compatibilidade e ao custo de reuso, não devem apagar a
relevância funcional.

## Artefatos reproduzíveis

Em `02_shell/cli/testdata/discovery_benchmark/` foram registrados:

- `screening.calibration.json`;
- quatro relatórios `*.screening.json`;
- `structural-impact.v2.snapshot.json` e sua peneira;
- `screening.audit-selection.json`;
- `screening.audit-revisions.json`;
- `screening.audit-report.json`.

Calibração, classificação e auditoria foram executadas com `env -i`, sem `.env`,
cache ou provedor remoto. Apenas a captura v2 e a leitura explícita de READMEs
usaram rede.

## Verificações

```text
go test -race ./...    PASS
go vet ./...           PASS
go build ...           PASS
crystalline-lint .     PASS — No violations found
git diff --check       PASS
```

## Próxima decisão

A peneira já é útil para explicar e ordenar, mas o baixo `precision@10` impede
automação agressiva. O próximo passo deve extrair README, manifesto, árvore e
referência imutável de `ADVANCE`, `AUDIT` e uma amostra de `DROP`, transformando
sinais textuais em evidência por capacidade.

O modelo deve ser recalibrado somente depois de novos julgamentos por ideia. Não
adicionar LLM ou embeddings antes de testar se evidência de conteúdo melhora os
folds e corrige a subestimação observada.
