# Laudo 0039 — Qualidade da triagem

**Amostra**: 24 casos de uma captura OpenCascade

**Resultado**: insuficiente para automação

## Resultado observado

| Decisão | relevante | irrelevante | incerto |
|---|---:|---:|---:|
| `ADVANCE` | 0 | 0 | 0 |
| `AUDIT` | 13 | 8 | 3 |
| `DROP` | 0 | 0 | 0 |

- cobertura de automação: `0`;
- taxa de auditoria: `1`;
- recall retido: `1`;
- average precision do ranking: `0.8166871814645094`;
- precisão nos 10 primeiros: `0.8`;
- lift nos 10 primeiros: `1.2923076923076924`;
- Brier skill contra prevalência: `-0.2557295551001093`;
- ECE em cinco faixas: `0.2612041677780249`.

## Conclusão

O modelo preservou candidatos relevantes porque não automatizou nenhuma decisão. Isso é
contenção, não triagem útil. O ranking contém sinal, mas o ganho e a calibração observados
não justificam promover o modelo, reduzir auditoria ou alegar generalização. A amostra é
pequena, de um único domínio e serve apenas para orientar a próxima melhoria.
