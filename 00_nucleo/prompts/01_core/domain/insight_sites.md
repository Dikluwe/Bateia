# Contratos de descoberta multifonte

**Camada:** L1. **Contrato:** passo 0138.

Declarar bindings de MarketDiscoveryConnector/leitor opcional, factory, pedido multifonte, plano fonte×query/budgets, tentativas, candidatos ProductCandidate com origem/seleção, recibos e contexto separado entre catálogo e leitura. Preservar identidade por conector e formato legado. Sem lógica/I/O.

## Evolução 0139 (precedência sobre limites de fontes do recorte0138)

Aplicar o passo0139: RepositoryReader opcional exclusivo Gitee/AtomGit, resolução imutável e contextos de leitura com identidade por sítio. Reader legado permanece só GitHub. Preservar resultados antigos quando o novo leitor não está vinculado. Erros do leitor classificados sem mensagens remotas; budgets compartilhados; tarefas agregam IDs dos cartões com origem.

## Evolução0141

Passo0141 tem precedência na agenda: alternar perspectivas e fontes com passo coprimo, tentativas em ordem canônica com execution_order, declarar scheduling_policy no plano. Extensões opcionais de cobertura e reference_leads; tarefas separam reference_lead_refs de catálogo/cartões. Trilhas derivadas de bytes reais, destinos REFERENCE_NOT_READ, zero HTTP adicional e zero promoção de novidade. Preservar modo GitHub legado e artefatos anteriores; atualizar testes apenas nas expectativas de agenda modificadas.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.


## Evolução0143

Extensão opcional neighborhood do passo0143, sem alterar JSON anterior quando ausente.
