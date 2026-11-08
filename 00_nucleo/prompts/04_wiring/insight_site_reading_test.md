# insight_site_reading_test

**Contrato:** [0139](../../passos/0139-ler-readmes-de-gitee-e-atomgit.md).

Testar factory real: GitHub Reader legado, novos RepositoryReader só Gitee/AtomGit; plano admitido sem HTTP.

## Evolução0140 (precedência sobre lista de fontes0139)

Aplicar passo0140: RepositoryReader também Codeberg/GitLab. Preservar namespace GitLab incluindo subgrupos, usar seletor puro diretamente apenas para seus bundles validados. Reader legado só GitHub; fontes anteriores preservadas. Composição passa a quatro RepositoryReader próprios, sem novas credenciais.
