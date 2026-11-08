# Execução de descoberta entre sítios

**Camada:** L1. **Contrato:** passo 0138.

Implementar DiscoverInsightSites com plano e bindings, execução fonte×query sob budgets/deadlines e registros de operações omitidas. AUTH_REQUIRED bloqueia só a fonte; validar/sanitizar respostas, preservar parcial/cancelamento e identidade por sítio. Coletar somente GitHub pelo owner 0137 e contexto ExploreInsights; demais fontes explícitas CATALOG_ONLY. Tarefas/instruções novas referenciam catálogo e leitura separadamente, sem inventar evidência/novidade. Sem I/O concreto ou estado de credenciais.

## Evolução 0139 (precedência sobre limites de fontes do recorte0138)

Aplicar o passo0139: RepositoryReader opcional exclusivo Gitee/AtomGit, resolução imutável e contextos de leitura com identidade por sítio. Reader legado permanece só GitHub. Preservar resultados antigos quando o novo leitor não está vinculado. Erros do leitor classificados sem mensagens remotas; budgets compartilhados; tarefas agregam IDs dos cartões com origem.

## Evolução0141

Passo0141 tem precedência na agenda: alternar perspectivas e fontes com passo coprimo, tentativas em ordem canônica com execution_order, declarar scheduling_policy no plano. Extensões opcionais de cobertura e reference_leads; tarefas separam reference_lead_refs de catálogo/cartões. Trilhas derivadas de bytes reais, destinos REFERENCE_NOT_READ, zero HTTP adicional e zero promoção de novidade. Preservar modo GitHub legado e artefatos anteriores; atualizar testes apenas nas expectativas de agenda modificadas.

## Evolução0142

Aplicar passo0142: memória validada antes de efeitos, followups de um salto com limites e origem própria, --memory/--follow-links na casca e feedback offline. FollowReader apenas GitHub, demais forges reaproveitam RepositoryReader. Memória em snapshot único de saída, sem releitura automática de examinados/descartados, sem promover histórico a evidência atual. Modo sem --sources preservado.

Se catálogo não completar nenhuma consulta mas o backlog da memória produzir leitura válida, conservar PARTIAL com conteúdo atual; FAILED/ErrInsightDiscoveryUnavailable só quando também não houver followup coletado. Falha de busca não apaga progresso da memória.

## 0154 — Descoberta contextual
Propagar política do plano ao pedido do conector e pergunta à seleção apenas em context-v2, preservando schedule/budgets/memória.

## 0155 — Ideias de algoritmos
Encaminhar AlgorithmIdeas à exploração inicial multifonte; contextos com suas próprias identidades continuam com quotas vigentes. Após expand, o leitor pode preparar novo explore no modo desejado.
