# Contratos funcionais consolidados

Esta pasta reúne regras compartilhadas entre prompts proprietários de uma capacidade. A autoridade estrutural permanece nas [regras canônicas](../../.agentrules), nos ADRs e no [mapa de responsabilidades](../arquitetura/README.md). Um contrato funcional não concede exceções às camadas Tekt.

| Capacidade | Contrato vigente | Aceitação | Origem da consolidação |
| --- | --- | --- | --- |
| Memória, aprofundamento e expansão de vizinhanças | [Exploração v1](exploracao-de-vizinhancas-v1.md) | [Matriz v1](exploracao-de-vizinhancas-aceitacao-v1.md) | [Passo0146](../passos/0146-consolidar-prompts-de-exploracao.md), após0142–0145; correção ex-ante no [passo0147](../passos/0147-corrigir-limites-e-verificar-contratos-de-exploracao.md) |

O [contrato de coleta e serialização de evidências v1](coleta-e-serializacao-de-evidencias-v1.md) reúne a admissão de planos, vínculo de origem e estabilidade de bundles do passo0158, com critérios EVD-01–05 e limites de compatibilidade histórica.

O sufixo v1 identifica o recorte dos formatos públicos existentes; não cria um novo schema. Consolidar texto não prova que todo o código nasceu desse texto ou que todos os critérios estão cobertos. Divergências conhecidas devem aparecer no contrato e na matriz. Novas mudanças de comportamento exigem intenção ex-ante e verificação próprias. Consulte o [guia de prompts](../documentacao/escrever-prompts.md).
