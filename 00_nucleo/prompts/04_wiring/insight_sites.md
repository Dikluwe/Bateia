# Composição da descoberta multifonte

**Camada:** L4. **Contrato:** passo 0138.

Factory lazy usa conectores federados existentes; substituir capability-only GitHub/Sourcegraph pelas bridges L1 sobre clientes reais. Só binding GitHub recebe reader. Preservar Bitbucket condicional configurado, tokens em seus clientes e ausência de credenciais em metadata. Sem regras de seleção/resultados na composição.

## Evolução 0139 (precedência sobre limites de fontes do recorte0138)

Aplicar o passo0139: RepositoryReader opcional exclusivo Gitee/AtomGit, resolução imutável e contextos de leitura com identidade por sítio. Reader legado permanece só GitHub. Preservar resultados antigos quando o novo leitor não está vinculado. Erros do leitor classificados sem mensagens remotas; budgets compartilhados; tarefas agregam IDs dos cartões com origem.

## Evolução0140 (precedência sobre lista de fontes0139)

Aplicar passo0140: RepositoryReader também Codeberg/GitLab. Preservar namespace GitLab incluindo subgrupos, usar seletor puro diretamente apenas para seus bundles validados. Reader legado só GitHub; fontes anteriores preservadas. Composição passa a quatro RepositoryReader próprios, sem novas credenciais.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.
