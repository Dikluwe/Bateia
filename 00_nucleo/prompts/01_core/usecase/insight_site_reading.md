# insight_site_reading

**Contrato:** [0139](../../../passos/0139-ler-readmes-de-gitee-e-atomgit.md).

Admissão de leitores por sítio, validação de referência resolvida, coleta reutilizada e contexto isolado; namespace de cartões e tarefas. GitHub legado inalterado. Orçamento compartilhado e erros seguros.

## Evolução0140 (precedência sobre lista de fontes0139)

Aplicar passo0140: RepositoryReader também Codeberg/GitLab. Preservar namespace GitLab incluindo subgrupos, usar seletor puro diretamente apenas para seus bundles validados. Reader legado só GitHub; fontes anteriores preservadas. Composição passa a quatro RepositoryReader próprios, sem novas credenciais.

Para subgrupos, o seletor legado usa immutableRef que recusa barras em Owner. Antes de selar o cartão, restaurar Source a partir de identidade/revisão validadas do wrapper e recalcular ID com Source/Kind/Reference, depois aplicar namespace por sítio. Nunca publicar Source vazio ou inventar owner.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.
