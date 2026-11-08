# CLI idea recover

Exponha diagnose e plan offline, e execute somente com provedor injetado. Leia o
dossiê sem modificar artefatos originais; grave cada rodada atomicamente sob
recovery/round-NN, mantenha índice com hashes, respeite max_rounds, histórico de
targets e --resume, e rejeite execução sem target novo.

## Fronteiras explícitas — 0129

Delegar efeitos concretos às portas de runtime/arquivos do ADR-0012. Preservar formatos, erro, validações, limites e ordens de chamadas/cleanup. Configuração real em L4; nenhum fallback de I/O concreto em L1/L2. Testes de integração fornecem host explicitamente.

0131: a projeção e os RuleIDs do histórico original pertencem a OriginalRecoveryTargets em L1. A casca somente lê os bundles e delega.

## 0161 — Histórico observado

Aplicar [EVD-11](../../../contratos/coleta-e-serializacao-de-evidencias-v1.md#recuperação-ativa-histórico-e-custo--0161): carregar plano e bundles de cada rodada anterior, delegando a L1 a seleção dos targets tentados. Não anexar todos os targets planejados nem decidir status na casca. Falha ao ler resultado exigido é erro antes de novas portas/publicação. Preservar artefatos e hashes antigos; origem ex-ante no [passo0161](../../../passos/0161-corrigir-transicoes-historico-e-revalidacao.md).
