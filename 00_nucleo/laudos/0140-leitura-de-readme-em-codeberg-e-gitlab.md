# Laudo 0140 — Leitura de README em Codeberg e GitLab

**Data:** 2026-10-28. **Resultado:** materialização e verificação aprovadas para o [contrato0140](../passos/0140-ler-readmes-de-codeberg-e-gitlab.md), especializando o ADR-0020.

## Entrega

A composição multifonte recebe leitores anônimos Codeberg/GitLab. Ambos confirmam metadata pública e identidade, resolvem branch padrão para commit e requisitam árvore/README nesse SHA. Conteúdo passa por limites, base64, tamanho e hash de blob; GitLab confere file_path, commit_id e content_sha256 quando presente. Sem novos tokens/flags, redirects, download_url ou execução de candidatos.

L1 conserva admissão de fontes, quotas e contexto; L3 implementa protocolos em endpoints fixos; L4 injeta leitores. A seleção de passagem para subgrupos GitLab usa o seletor puro existente com bundle validado, Source restaurado a partir da identidade/revisão e ID recalculado antes de aplicar namespace do sítio. O normalizador público legado e domain.CommitRef não foram alterados. Namespace completo permanece em Owner e Source.Repository; wrapper site_contexts conserva origem. As tarefas da IA recebem os cartões das novas fontes.

Os limites existentes continuam: README regular raiz, quatro GETs por referência, máximo32KiB por documento/128KiB agregados, deadlines e cartões globais. Árvores truncadas/paginadas ou100 entradas são recusadas conservadoramente. Outras fontes permanecem catálogo.

## Evidência

make verify passou com lint e11 controles negativos, auditoria de379 arquivos produtivos, testes do módulo raiz com loopback autorizado, vet, build e diff-check. As fixtures usam os conectores/leitores reais em transporte controlado. Cada fonte executa cinco buscas e quatro GETs de leitura; Codeberg/GitLab geram contextos/cartões distintos mesmo para owner/repo/commit/texto idênticos. Subgrupos GitLab, encoding do ID do projeto, SHA fixado, limites exatos32KiB/SHA256 e referências das tarefas passaram. Casos negativos cobrem identidade/publicidade/branch/SHA, paginação/symlink, corpos, base64/tamanho/blob, commit/path/hash, auth/ausência/quota, redirect e cancelamento. Falha GitLab conserva contexto Codeberg.

Os seis replays históricos0106,0134–0138 passaram com o binário final. Os224 arquivos de evidência anteriores permanecem idênticos. A única mudança em teste anterior é a expectativa de composição de2 para4 leitores, verificada por reconstrução/hash; demais testes anteriores são intactos. Foram acrescentados um arquivo produtivo e dois de testes; três arquivos produtivos existentes e um teste de composição foram alterados. Prompts e hashes atuais estão na linhagem nova.

## Limites e referências

Não houve ensaio online: fixtures não certificam disponibilidade remota ou utilidade/novidade. Codeberg usa o protocolo Forgejo compatível documentado, mas seu OpenAPI não pôde ser acessado pelo navegador de pesquisa. O GitLab restringe-se à instância pública gitlab.com. Identidades ASCII, README raiz e integridade de bytes continuam limites explícitos; documentos convertidos pela API com blob divergente são recusados. Efeito permanece UNVERIFIED e novidade NOT_ASSESSED.

Uso e protocolos estão na [bancada0140](../../lab/bancada/insights/0140/README.md); [recibo](../../lab/bancada/insights/0140/verification.json) e [linhagem](../../lab/bancada/insights/0140/lineage.json) registram o recorte. Nenhum software candidato foi executado e nenhum artefato histórico foi reescrito.
