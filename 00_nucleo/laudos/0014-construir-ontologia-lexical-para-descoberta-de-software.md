# Laudo 0014 — Ontologia lexical para descoberta de software

**Estado**: `IMPLEMENTADO`  
**Data**: 2026-08-26

## Materializado

- ontologia determinística com sete conceitos, versão e hash canônico;
- resolução de triggers com normalização de caixa e diacríticos;
- queries conceituais limitadas a duas por conceito, com origem rastreável;
- auditoria pré-rede de discriminator, tamanho, proveniência, genericidade e
  termos proibidos;
- corpus offline de doze IdeaSpecs, incluindo três controles negativos;
- comandos `idea ontology expand`, `evaluate` e `audit`;
- materialização separada de plano auditado e consumo por `idea discover
  --plan`;
- propagação dos discriminadores até o Sourcegraph e fallback que os preserva.

A ontologia congelada tem hash:

```text
2ce75ee5d09d117c383d7e022dd98c16556e95fe15021e97698838aa05211140
```

## Gates offline

| Métrica | Resultado |
|---|---:|
| concept precision | 100% |
| concept recall | 100% |
| query acceptance | 100% |
| discriminator retention | 100% |
| generic fallback | 0% |
| forbidden leak | 0% |
| control contamination | 0% |
| uncovered sources | 22,22% |

Os gates passaram sem rede. A IdeaSpec do Bateia ativou seis conceitos e
materializou doze queries, todas aceitas e abaixo do teto de quatorze.

## Ataque encontrado durante o live

O primeiro braço B expôs uma violação real: embora as queries originais fossem
auditadas, o fallback legado do adaptador descartava seus discriminadores. Por
exemplo, `internal code reuse search` virava `internal code` e `GitHub
repository recommendation` virava `GitHub repository`.

Esse run foi preservado como `discovery-0014-invalid-fallback.json`, não contado
como resultado. A correção adicionou discriminadores ao contrato da expansão e
da requisição; o fallback passou a usar a expressão discriminatória declarada.
Uma regressão prova essa propriedade antes da rede.

## Experimento público A/B

Política comum: `max-per-query=3`, concorrência 2, timeout total de três minutos
e attempt timeout de vinte segundos.

| Observável | A — lexical 0013 | B — ontologia 0014 |
|---|---:|---:|
| queries planejadas | 16 | 12 |
| queries concluídas | 10 | 12 |
| completion rate | 62,5% | 100% |
| attempts expirados | 18 | 13 |
| candidates únicos | 3 | 31 |
| sementes repo-backed | 0/4 | 0/4 |
| recall competitivo | 0% | 0% |
| identidade imutável | 0 | 0 |

O braço B melhorou execução e yield, mas não encontrou Replen, OSS Compass,
ARIS ou Starboard. Todos os 31 candidatos continuaram sem revisão imutável e os
resultados observados eram matches incidentais em código, não descoberta de
produto. Nenhum nome ou owner do registry estava nas queries.

## Veredito

A hipótese lexical passou offline, mas Sourcegraph falhou como catálogo de
produtos: recall abaixo de 40%, matches incidentais e 100% dos candidatos sem
identidade resolvida. Isso aciona explicitamente a decisão prevista de separar
catálogo e comprovação.

O Sourcegraph deve permanecer como provedor de prova em código depois que um
catálogo identificar repositórios. O próximo passo útil é uma porta de catálogo
que pesquise descrição, README, tópicos e metadados, mantendo a ontologia e a
auditoria construídas aqui como geradoras de intenção.

Artefatos reproduzíveis ficam em `00_nucleo/benchmark/0012/run/`: expansão,
plano, auditoria, avaliação offline, run inválido preservado e snapshot live
final.
