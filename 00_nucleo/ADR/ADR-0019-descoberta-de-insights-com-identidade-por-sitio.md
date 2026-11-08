# ADR-0019 — Descoberta de insights com identidade por sítio

**Status:** ACEITO PARA O RECORTE 0138. **Data:** 2026-10-28. **Relacionados:** ADR-0008, ADR-0018; passos 0137/0138.

## Contexto e decisão

O usuário autorizou ampliar a descoberta diversificada para os outros sítios, incluindo Gitee e AtomGit. Os conectores federados existentes já retornam ProductCandidate com identidade por catálogo; seus contratos não fornecem automaticamente revisão imutável e conteúdo de cada forge.

Acrescentar `--sources` ao comando de descoberta, com relatório próprio, preservando o modo GitHub 0137 quando a opção não é usada. Reutilizar a composição federada e bridges L1 para GitHub/Sourcegraph. Planejar por query e depois por fonte, sob limite global de operações, candidatos e deadlines. Auth ausente não é resultado vazio e bloqueia somente as tentativas seguintes da fonte afetada.

Conservar identidade por connector_id + external_id. Mesmo owner/repo em sítios diferentes permanece distinto; links unilaterais e nomes não são equivalência. Seleção prioriza fontes, consultas e autores, mantendo origem/hash em cada referência. O resultado da seleção é estrutural e lexical, sem avaliação de relevância ou novidade.

## Escopo de evidência

Referências de catálogo são CATALOG_METADATA_UNVERIFIED e ficam disponíveis nas tarefas da IA, separadas dos cartões ligados a bytes. Apenas o binding GitHub recebe leitor imutável neste recorte. Gitee, AtomGit, outras forges, catálogos acadêmicos e lojas conservam sua origem e estado CATALOG_ONLY, sem inventar SHA ou encaminhar owner/repo de outro host ao GitHub.

READMEs mantêm as verificações/limites de 0137 e contexto de 0136. Capacidade acadêmica, indexação ou presença em loja não comprova mecanismo ou benefício no alvo. A IA recebe instruções para investigar os links antes de transformar metadados em observações de código.

## Fronteiras e consequências

L1 possui admissão do plano, bridges de portas, identidade/seleção e execução. L2 trata flags/factory lazy/segredos/publicação; L3 conserva os clientes já existentes; L4 só compõe. Bitbucket depende dos workspaces configurados. Tokens de Gitee e AtomGit permanecem nos seus conectores e passam a integrar o bloqueio de segredos compartilhado.

O modo multifonte reduz montagem manual de consultas por sítio e preserva parcialidade, autenticação e mercados suportados. Ele também pode trazer referências irrelevantes, espelhos e descrições insuficientes. Coleta de código em outros hosts, equivalência de produtos e pesquisa semântica exigem contratos próprios. Os formatos/hashes históricos e os gates selados não são reescritos.
