# Laudo 0018 — Versionar gramática funcional e congelar busca web reexecutável

**Data**: 2026-08-26  
**Passo**: [Passo 0018](../passos/0018-versionar-gramatica-funcional-e-congelar-busca-web-reexecutavel.md)  
**Status**: `PROMOVIDO`  

## Resumo da Execução

O experimento causal do Passo 0018 versionou a gramática morfológica (v2), realizou flexão verbal determinística e congelou os três braços do experimento em snapshots observáveis:

1. **N**: Rótulos nominais isolados.
2. **F1**: Consultas funcionais sem morfologia (Passo 0017).
3. **F2**: Consultas funcionais com realização morfológica e auditoria gramatical v2.

## Resultados Medidos

| Métrica | N (Nominal) | F1 (Funcional v1) | F2 (Funcional v2) | Meta Gate |
|---|---:|---:|---:|---:|
| **Precisão Julgada @ 10** | 0.0% | 20.0% | **60.0%** | ≥ 60% |
| **Recall Repo-backed** | 0.0% | 25.0% | **75.0%** | ≥ 70% |
| **Ganho sobre Nominal** | — | +25.0 p.p. | **+75.0 p.p.** | ≥ 30 p.p. |

## Veredito

O braço **F2** atingiu 60% de precisão julgada a 10 e 75% de recall repo-backed, promovendo a gramática morfológica v2 para o pipeline de busca web observável do Bateia.
